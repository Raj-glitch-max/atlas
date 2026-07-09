// Command atlas-server is the Atlas product backend: a composition root that
// exposes the real issuance, verification, and revocation engine over an HTTP
// JSON API. It contains wiring, HTTP I/O, storage, and audit only — never
// delegation logic (FD-9). Every decision comes from internal/verify and
// internal/issuance; the server never re-implements a check.
//
// This is the driver the CLI, SDKs, and the Atlas MCP server all sit on top
// of. Single trust domain in v1; the key that signs delegations also signs the
// revocation snapshots. Storage is in-memory behind a Store interface so a
// durable backend can be swapped in without touching the engine.
package main

import (
	"crypto/ecdsa"
	"fmt"
	"sync"
	"time"

	"github.com/Raj-glitch-max/atlas/internal/issuance"
	"github.com/Raj-glitch-max/atlas/internal/record"
	"github.com/Raj-glitch-max/atlas/internal/revstatus"
	"github.com/Raj-glitch-max/atlas/internal/truststore"
	"github.com/Raj-glitch-max/atlas/internal/verify"
	"github.com/spiffe/go-spiffe/v2/spiffeid"
)

// Clock is the injectable time port; production uses the wall clock, tests a
// controllable one.
type Clock interface{ Now() time.Time }

type systemClock struct{}

func (systemClock) Now() time.Time { return time.Now() }

// revocationAdapter bridges a revstatus.Provider onto verify's port (AD-020) —
// the same composition glue the atlas-verify driver uses.
type revocationAdapter struct{ p revstatus.Provider }

func (a revocationAdapter) StatusOf(instance record.InstanceID) verify.RevocationStatus {
	ans := a.p.StatusOf(instance)
	var st verify.RevocationState
	switch ans.State {
	case revstatus.NotObservedRevoked:
		st = verify.NotObservedRevoked
	case revstatus.ObservablyRevoked:
		st = verify.ObservablyRevoked
	default:
		st = verify.Indeterminate
	}
	return verify.RevocationStatus{State: st, AsOf: ans.AsOf}
}

// grantedPerms is the v1 permission source: every principal may delegate
// within a fixed granted set; anything outside is refused by the real
// authority (demonstrating attenuation). A real deployment wires this to an
// external authorization source.
type grantedPerms struct{ set issuance.PermissionSet }

func (g grantedPerms) PermissionsOf(spiffeid.ID) (issuance.PermissionSet, bool) { return g.set, true }

func defaultPermissions() issuance.PermissionSet {
	return issuance.NewPermissionSet(
		"read:orders", "write:orders", "read:audit", "write:audit",
		"read:ledger", "write:ledger", "read:metrics", "read:keys", "admin:all",
	)
}

// Config is the server's construction-time configuration.
type Config struct {
	Domain    string   // trust domain (issued principals must belong to it)
	StorePath string   // durable state file ("" => in-memory)
	KeyPath   string   // authority key file ("" => ephemeral)
	APIKey    string   // optional bearer token guarding mutating endpoints
	Grant     []string // scopes a principal may delegate ("" => default set)

	AllowOrigin  string // CORS Access-Control-Allow-Origin ("" => "*", dev default)
	RateLimitRPM int    // per-IP requests/min on mutating endpoints (0 => off)
}

// App is the server's composition root and in-process state.
type App struct {
	clock       Clock
	domain      spiffeid.TrustDomain
	keyID       string
	pubKeyHex   string
	apiKey      string
	authority   *issuance.Authority
	trust       *truststore.Store
	publisher   *revstatus.Publisher
	provider    *revstatus.SignedSetProvider
	policy      verify.Policy
	store       *Store
	revWindow   time.Duration
	allowOrigin string
	limiter     *rateLimiter

	mu       sync.RWMutex // guards revoked set + provider ingest vs. verify reads
	revoked  map[string]record.InstanceID
	lastAsOf time.Time
	lastSet  revstatus.SignedRevokedSet // last published snapshot (for /bundle export)
}

// NewApp wires the full engine for a trust domain. The revocation window R and
// skew tolerance are the resolved scope-act values (R = 2s, skew = 30s). A
// non-empty storePath enables durable persistence (delegations, audit,
// metrics, and — rebuilt on load — the revoked set).
func NewApp(cfg Config, clock Clock) (*App, error) {
	td, err := spiffeid.TrustDomainFromString(cfg.Domain)
	if err != nil {
		return nil, fmt.Errorf("trust domain %q: %w", cfg.Domain, err)
	}
	key, err := loadOrCreateKey(cfg.KeyPath)
	if err != nil {
		return nil, err
	}
	const keyID = "authority-key-1"

	grantSet := defaultPermissions()
	if len(cfg.Grant) > 0 {
		grantSet = issuance.NewPermissionSet(cfg.Grant...)
	}
	authority, err := issuance.NewAuthority(
		record.Signer{Key: key, KeyID: keyID},
		grantedPerms{set: grantSet},
		issuance.NoBinding{}, issuance.RandomMinter{}, clock)
	if err != nil {
		return nil, err
	}
	tm, err := record.NewTrustMaterial(td, map[string]*ecdsa.PublicKey{keyID: &key.PublicKey})
	if err != nil {
		return nil, err
	}
	trust := truststore.New()
	if err := trust.Provision(tm, clock.Now()); err != nil {
		return nil, err
	}
	policy, err := verify.NewPolicy(2*time.Second, 30*time.Second)
	if err != nil {
		return nil, err
	}
	listID := "atlas-revlist:" + cfg.Domain
	publisher, err := revstatus.NewPublisher(key, listID)
	if err != nil {
		return nil, err
	}
	provider := revstatus.NewSignedSetProvider(&key.PublicKey, listID)

	store := NewStore(cfg.StorePath)
	if err := store.Load(); err != nil {
		return nil, fmt.Errorf("load store %q: %w", cfg.StorePath, err)
	}
	allowOrigin := cfg.AllowOrigin
	if allowOrigin == "" {
		allowOrigin = "*"
	}
	var limiter *rateLimiter
	if cfg.RateLimitRPM > 0 {
		limiter = newRateLimiter(cfg.RateLimitRPM, clock)
	}
	app := &App{
		clock: clock, domain: td, keyID: keyID, apiKey: cfg.APIKey,
		pubKeyHex: publicKeyHex(&key.PublicKey),
		authority: authority, trust: trust, publisher: publisher, provider: provider,
		policy: policy, store: store, revWindow: 2 * time.Second,
		revoked: map[string]record.InstanceID{}, allowOrigin: allowOrigin, limiter: limiter,
	}
	// Rebuild the revoked set from any persisted revoked delegations.
	for _, instStr := range store.RevokedInstances() {
		if inst, err := record.InstanceIDFromString(instStr); err == nil {
			app.revoked[instStr] = inst
		}
	}
	// Seed a signed snapshot (empty, or the restored revoked set) so
	// verification answers immediately with the correct revocation state.
	if err := app.republishLocked(); err != nil {
		return nil, err
	}
	return app, nil
}

// Flush persists the store if durability is enabled and there are changes.
func (a *App) Flush() error { return a.store.Flush() }

// nextAsOf returns a strictly-increasing signed timestamp (the realization
// only adopts snapshots with a newer as-of).
func (a *App) nextAsOf() time.Time {
	now := a.clock.Now()
	if !now.After(a.lastAsOf) {
		now = a.lastAsOf.Add(time.Millisecond)
	}
	a.lastAsOf = now
	return now
}

// republishLocked re-signs the current revoked set with a fresh as-of and
// ingests it. Caller must hold a.mu for write (or be in construction).
func (a *App) republishLocked() error {
	set := make([]record.InstanceID, 0, len(a.revoked))
	for _, v := range a.revoked {
		set = append(set, v)
	}
	snap, err := a.publisher.Publish(set, a.nextAsOf())
	if err != nil {
		return err
	}
	if _, err := a.provider.Ingest(snap); err != nil {
		return err
	}
	a.lastSet = snap
	return nil
}

// Bundle exports the relying-party trust bundle: the trust material (public
// key) plus the latest signed revocation snapshot. A holder of this bundle can
// verify delegations fully offline; the snapshot's signature makes the bundle
// tamper-evident (a doctored bundle is refused at import).
func (a *App) Bundle() BundleDTO {
	a.mu.RLock()
	defer a.mu.RUnlock()
	revoked := make([]string, 0, len(a.lastSet.Revoked))
	for _, r := range a.lastSet.Revoked {
		revoked = append(revoked, r.String())
	}
	return BundleDTO{
		Version:     1,
		TrustDomain: a.domain.Name(),
		Keys:        map[string]string{a.keyID: a.pubKeyHex},
		Revocation: BundleRevocation{
			ListID:  a.lastSet.ListID,
			AsOf:    a.lastSet.AsOf,
			Revoked: revoked,
			Sig:     a.lastSet.Sig,
		},
		ExportedAt: a.clock.Now().UTC(),
	}
}

// Refresh re-publishes the current set with a fresh as-of; the background
// loop calls it every < R so the held snapshot never ages past the freshness
// bound even when no revocations occur.
func (a *App) Refresh() error {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.republishLocked()
}

// Issue creates a delegation via the real authority. Over-scope requests are
// Refused by the engine (attenuation), surfaced as a 422.
func (a *App) Issue(principal, delegate string, scope []string, ttl time.Duration) (*IssueResult, *apiError) {
	p, err := spiffeid.FromString(principal)
	if err != nil {
		return nil, badRequest("principal is not a valid SPIFFE ID: " + err.Error())
	}
	if p.TrustDomain() != a.domain {
		return nil, badRequest(fmt.Sprintf("principal must be in this server's trust domain %q", a.domain.Name()))
	}
	d, err := spiffeid.FromString(delegate)
	if err != nil {
		return nil, badRequest("delegate is not a valid SPIFFE ID: " + err.Error())
	}
	if ttl <= 0 {
		ttl = time.Hour
	}
	res, err := a.authority.Issue(issuance.Request{
		Principal: p, Delegate: d, Scope: scope, Expiration: a.clock.Now().Add(ttl),
	})
	if err != nil {
		return nil, serverError(err.Error())
	}
	if res.Outcome == issuance.Refused {
		a.store.Audit(AuditEvent{Time: a.clock.Now(), Type: "issue.refused", Principal: principal, Delegate: delegate, Detail: res.Refusal.String()})
		return nil, &apiError{Status: 422, Message: "issuance refused: " + res.Refusal.String(), Refused: true}
	}
	asrt := res.Record.Read()
	inst := asrt.Instance.String()
	a.store.AddDelegation(&Delegation{
		Instance: inst, Principal: principal, Delegate: delegate, Scope: scope,
		IssuedAt: asrt.IssuedAt, ExpiresAt: asrt.Expiration,
	})
	a.store.Audit(AuditEvent{Time: a.clock.Now(), Type: "issue", Principal: principal, Delegate: delegate, Instance: inst})
	return &IssueResult{
		Record: string(res.Record.Presented()), Instance: inst,
		Principal: principal, Delegate: delegate, Scope: scope,
		IssuedAt: asrt.IssuedAt, ExpiresAt: asrt.Expiration,
	}, nil
}

// Verify runs the real Verification Core against a presented record. The
// latency is measured around Verify (AT26), never inside it.
func (a *App) Verify(rec string) *VerifyResult {
	a.mu.RLock()
	v, err := verify.NewVerifier(a.policy, a.trust, revocationAdapter{p: a.provider}, a.clock)
	if err != nil {
		a.mu.RUnlock()
		return &VerifyResult{Decision: "error", Causes: []string{err.Error()}}
	}
	start := time.Now()
	verdict, trace := v.Verify([]byte(rec))
	elapsed := time.Since(start)
	a.mu.RUnlock()

	dec := decisionString(verdict.Decision)
	a.store.RecordVerdict(dec)
	a.store.ObserveLatency(elapsed.Seconds())
	a.store.Audit(AuditEvent{Time: a.clock.Now(), Type: "verify", Decision: dec, Detail: joinCauses(verdict.Causes)})
	return &VerifyResult{
		Decision:      dec,
		Accept:        verdict.IsAccept(),
		Causes:        causeStrings(verdict.Causes),
		Trace:         traceDTO(trace),
		LatencyMicros: elapsed.Microseconds(),
	}
}

// Revoke adds an instance to the signed revoked set and republishes it.
func (a *App) Revoke(instanceStr string) *apiError {
	inst, err := record.InstanceIDFromString(instanceStr)
	if err != nil {
		return badRequest("instance is not a valid instance id: " + err.Error())
	}
	a.mu.Lock()
	a.revoked[instanceStr] = inst
	err = a.republishLocked()
	a.mu.Unlock()
	if err != nil {
		return serverError("failed to publish revocation snapshot: " + err.Error())
	}
	a.store.MarkRevoked(instanceStr)
	a.store.Audit(AuditEvent{Time: a.clock.Now(), Type: "revoke", Instance: instanceStr})
	return nil
}

// snapshotAge reports how long ago the held revocation snapshot was signed —
// the metric operators watch against R.
func (a *App) snapshotAge() time.Duration {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.clock.Now().Sub(a.lastAsOf)
}
