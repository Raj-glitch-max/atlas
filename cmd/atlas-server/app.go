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
	LogRequests  bool   // emit a structured access log line per request
	LogVerbose   bool   // also log noisy probe/scrape endpoints
}

// App is the server's composition root. It holds only state that is SHARED
// across callers — the authority key, trust material, issuance policy, and
// transport config. All mutable per-caller state (delegation store, audit log,
// revoked set, signed snapshot) lives in a Session.
type App struct {
	clock       Clock
	domain      spiffeid.TrustDomain
	keyID       string
	pubKey      *ecdsa.PublicKey
	pubKeyHex   string
	listID      string
	apiKey      string
	authority   *issuance.Authority
	trust       *truststore.Store
	publisher   *revstatus.Publisher
	policy      verify.Policy
	revWindow   time.Duration
	allowOrigin string
	limiter     *rateLimiter
	logRequests bool
	logVerbose  bool

	// def is the default session: the durable, single-graph behaviour used by
	// every client that does not send a session id (CLI, SDKs, MCP, tests).
	def *Session

	sessionsMu sync.Mutex
	sessions   map[string]*Session

	// activity is the cross-session aggregate exposed at /activity.
	activity *activityLog
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
		pubKey: &key.PublicKey, pubKeyHex: publicKeyHex(&key.PublicKey), listID: listID,
		authority: authority, trust: trust, publisher: publisher,
		policy: policy, revWindow: 2 * time.Second,
		allowOrigin: allowOrigin, limiter: limiter,
		logRequests: cfg.LogRequests, logVerbose: cfg.LogVerbose,
		sessions:    map[string]*Session{},
		activity:    newActivityLog(),
	}
	// The default session carries the durable store — this is the pre-existing
	// single-graph behaviour, unchanged for every client that does not opt into
	// a session.
	app.def = &Session{
		ID: "default", created: clock.Now(), lastSeen: clock.Now(),
		store: store, provider: provider, revoked: map[string]record.InstanceID{},
	}
	// Rebuild the revoked set from any persisted revoked delegations.
	for _, instStr := range store.RevokedInstances() {
		if inst, err := record.InstanceIDFromString(instStr); err == nil {
			app.def.revoked[instStr] = inst
		}
	}
	// Seed a signed snapshot (empty, or the restored revoked set) so
	// verification answers immediately with the correct revocation state.
	if err := app.republish(app.def); err != nil {
		return nil, err
	}
	return app, nil
}

// Flush persists the default session's store if durability is enabled. Only the
// default session is durable; visitor sessions are intentionally ephemeral.
func (a *App) Flush() error { return a.def.store.Flush() }

// nextAsOf returns a strictly-increasing signed timestamp for a session (the
// realization only adopts snapshots with a newer as-of). Caller holds s.mu.
func (a *App) nextAsOf(s *Session) time.Time {
	now := a.clock.Now()
	if !now.After(s.lastAsOf) {
		now = s.lastAsOf.Add(time.Millisecond)
	}
	s.lastAsOf = now
	return now
}

// republishLocked re-signs a session's revoked set with a fresh as-of and
// ingests it into that session's provider. Caller must hold s.mu for write.
func (a *App) republishLocked(s *Session) error {
	set := make([]record.InstanceID, 0, len(s.revoked))
	for _, v := range s.revoked {
		set = append(set, v)
	}
	snap, err := a.publisher.Publish(set, a.nextAsOf(s))
	if err != nil {
		return err
	}
	if _, err := s.provider.Ingest(snap); err != nil {
		return err
	}
	s.lastSet = snap
	return nil
}

// republish is republishLocked with the session lock taken.
func (a *App) republish(s *Session) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return a.republishLocked(s)
}

// Bundle exports the relying-party trust bundle: the trust material (public
// key) plus the latest signed revocation snapshot. A holder of this bundle can
// verify delegations fully offline; the snapshot's signature makes the bundle
// tamper-evident (a doctored bundle is refused at import).
func (a *App) Bundle(s *Session) BundleDTO {
	s.mu.RLock()
	defer s.mu.RUnlock()
	revoked := make([]string, 0, len(s.lastSet.Revoked))
	for _, r := range s.lastSet.Revoked {
		revoked = append(revoked, r.String())
	}
	return BundleDTO{
		Version:     1,
		TrustDomain: a.domain.Name(),
		Keys:        map[string]string{a.keyID: a.pubKeyHex},
		Revocation: BundleRevocation{
			ListID:  s.lastSet.ListID,
			AsOf:    s.lastSet.AsOf,
			Revoked: revoked,
			Sig:     s.lastSet.Sig,
		},
		ExportedAt: a.clock.Now().UTC(),
	}
}

// Refresh re-publishes every live session's set with a fresh as-of; the
// background loop calls it every < R so no held snapshot ages past the
// freshness bound even when no revocations occur. A session whose snapshot went
// stale would start answering Indeterminate and fail closed — correct, but not
// what a visitor idling on the page should see.
func (a *App) Refresh() error {
	if err := a.republish(a.def); err != nil {
		return err
	}
	a.sessionsMu.Lock()
	live := make([]*Session, 0, len(a.sessions))
	for _, s := range a.sessions {
		live = append(live, s)
	}
	a.sessionsMu.Unlock()
	for _, s := range live {
		if err := a.republish(s); err != nil {
			return err
		}
	}
	return nil
}

// Issue creates a delegation via the real authority. Over-scope requests are
// Refused by the engine (attenuation), surfaced as a 422.
func (a *App) Issue(s *Session, principal, delegate string, scope []string, ttl time.Duration) (*IssueResult, *apiError) {
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
		s.store.Audit(AuditEvent{Time: a.clock.Now(), Type: "issue.refused", Principal: principal, Delegate: delegate, Detail: res.Refusal.String()})
		return nil, &apiError{Status: 422, Message: "issuance refused: " + res.Refusal.String(), Refused: true}
	}
	asrt := res.Record.Read()
	inst := asrt.Instance.String()
	s.store.AddDelegation(&Delegation{
		Instance: inst, Principal: principal, Delegate: delegate, Scope: scope,
		IssuedAt: asrt.IssuedAt, ExpiresAt: asrt.Expiration,
	})
	s.store.Audit(AuditEvent{Time: a.clock.Now(), Type: "issue", Principal: principal, Delegate: delegate, Instance: inst})
	a.activity.record(actIssued, a.clock.Now())
	return &IssueResult{
		Record: string(res.Record.Presented()), Instance: inst,
		Principal: principal, Delegate: delegate, Scope: scope,
		IssuedAt: asrt.IssuedAt, ExpiresAt: asrt.Expiration,
	}, nil
}

// Verify runs the real Verification Core against a presented record. The
// latency is measured around Verify (AT26), never inside it.
func (a *App) Verify(s *Session, rec string) *VerifyResult {
	s.mu.RLock()
	v, err := verify.NewVerifier(a.policy, a.trust, revocationAdapter{p: s.provider}, a.clock)
	if err != nil {
		s.mu.RUnlock()
		return &VerifyResult{Decision: "error", Causes: []string{err.Error()}}
	}
	start := time.Now()
	verdict, trace := v.Verify([]byte(rec))
	elapsed := time.Since(start)
	s.mu.RUnlock()

	dec := decisionString(verdict.Decision)
	s.store.RecordVerdict(dec)
	s.store.ObserveLatency(elapsed.Seconds())
	s.store.Audit(AuditEvent{Time: a.clock.Now(), Type: "verify", Decision: dec, Detail: joinCauses(verdict.Causes)})
	a.activity.record(actVerified, a.clock.Now())
	return &VerifyResult{
		Decision:      dec,
		Accept:        verdict.IsAccept(),
		Causes:        causeStrings(verdict.Causes),
		Trace:         traceDTO(trace),
		LatencyMicros: elapsed.Microseconds(),
	}
}

// Revoke adds an instance to this session's signed revoked set and republishes
// it. A revocation is scoped to the session that issued it: one visitor cannot
// revoke another visitor's capability, because the instance is not in the other
// session's revoked set and never enters its snapshot.
func (a *App) Revoke(s *Session, instanceStr string) *apiError {
	inst, err := record.InstanceIDFromString(instanceStr)
	if err != nil {
		return badRequest("instance is not a valid instance id: " + err.Error())
	}
	s.mu.Lock()
	s.revoked[instanceStr] = inst
	err = a.republishLocked(s)
	s.mu.Unlock()
	if err != nil {
		return serverError("failed to publish revocation snapshot: " + err.Error())
	}
	s.store.MarkRevoked(instanceStr)
	s.store.Audit(AuditEvent{Time: a.clock.Now(), Type: "revoke", Instance: instanceStr})
	a.activity.record(actRevoked, a.clock.Now())
	return nil
}

// snapshotAge reports how long ago a session's held revocation snapshot was
// signed — the metric operators watch against R.
func (a *App) snapshotAge(s *Session) time.Duration {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return a.clock.Now().Sub(s.lastAsOf)
}
