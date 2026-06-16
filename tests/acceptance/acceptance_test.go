package acceptance_test

// In-process acceptance suite (IMPLEMENTATION_MASTER_PLAN.md E5-T4;
// docs/engineering/05_ACCEPTANCE_TEST_PLAN.md). It exercises the full system
// wired end to end — Issuance Authority (M2) -> Delegation Record (M1) ->
// Verification Core (M3) with the Trust Material Store (M4), a revocation
// provider (M5, via the AD-020 adapter), the clock, and the Revocation Origin
// register (M6) — at the product boundaries each AT names.
//
// This is the Sprint-1 exit milestone (M4): a complete
// issue -> verify -> tamper -> expire -> revoke -> reconstruct flow in one
// process. ATs that require the two-domain SPIRE substrate or the
// spike-selected revocation propagation channel skip with a named blocker
// (never a TODO); the map of blockers is in tests/README.md.

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"testing"
	"time"

	"github.com/Raj-glitch-max/atlas/internal/issuance"
	"github.com/Raj-glitch-max/atlas/internal/record"
	"github.com/Raj-glitch-max/atlas/internal/revorigin"
	"github.com/Raj-glitch-max/atlas/internal/revstatus"
	"github.com/Raj-glitch-max/atlas/internal/truststore"
	"github.com/Raj-glitch-max/atlas/internal/verify"
	"github.com/Raj-glitch-max/atlas/tests/harness"
	"github.com/spiffe/go-spiffe/v2/spiffeid"
)

const (
	principalID = "spiffe://domain-a.test/principal"
	delegateID  = "spiffe://domain-a.test/delegate"
	keyID       = "authority-key-1"
)

var (
	issuedAt = time.Unix(1_800_000_000, 0).UTC()
	expiry   = time.Unix(1_800_003_600, 0).UTC() // +1h
	nowMid   = time.Unix(1_800_000_300, 0).UTC() // +5m
)

// system is the whole wired product for one test.
type system struct {
	authority *issuance.Authority
	trust     *truststore.Store
	register  *revorigin.Register
	clock     *harness.Clock
	principal spiffeid.ID
	delegate  spiffeid.ID
}

func newSystem(t *testing.T) *system {
	t.Helper()
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	principal := spiffeid.RequireFromString(principalID)
	clock := harness.NewClock(issuedAt)

	perms := harness.NewPermissionSource().Grant(principal, "read:orders", "write:audit", "admin:all")
	authority, err := issuance.NewAuthority(
		record.Signer{Key: key, KeyID: keyID}, perms,
		issuance.NoBinding{}, issuance.RandomMinter{}, clock)
	if err != nil {
		t.Fatal(err)
	}

	// Relying party in domain B holds domain A's trust material, provisioned
	// out of band (the gate C1 manual-bundle-exchange pattern).
	tm, err := record.NewTrustMaterial(principal.TrustDomain(), map[string]*ecdsa.PublicKey{keyID: &key.PublicKey})
	if err != nil {
		t.Fatal(err)
	}
	trust := truststore.New()
	if err := trust.Provision(tm, issuedAt); err != nil {
		t.Fatal(err)
	}

	return &system{
		authority: authority,
		trust:     trust,
		register:  revorigin.New(),
		clock:     clock,
		principal: principal,
		delegate:  spiffeid.RequireFromString(delegateID),
	}
}

func (s *system) issue(t *testing.T, scope ...string) *record.Record {
	t.Helper()
	res, err := s.authority.Issue(issuance.Request{
		Principal: s.principal, Delegate: s.delegate, Scope: scope, Expiration: expiry,
	})
	if err != nil {
		t.Fatalf("issue: %v", err)
	}
	if res.Outcome != issuance.Issued {
		t.Fatalf("issue refused: %s", res.Refusal)
	}
	return res.Record
}

// verifier builds a verifier wired with the given revocation port.
func (s *system) verifier(t *testing.T, rev verify.RevocationStatusPort) *verify.Verifier {
	t.Helper()
	policy, err := verify.NewPolicy(time.Minute, 30*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	v, err := verify.NewVerifier(policy, s.trust, rev, s.clock)
	if err != nil {
		t.Fatal(err)
	}
	return v
}

// freshNotRevoked is a revocation port reporting a fresh not-revoked
// observation for any instance, as of the clock's now (a deliberate
// "revocation known current" scenario, not laundered ignorance).
func (s *system) freshNotRevoked() verify.RevocationStatusPort {
	return harness.NewUniformRevocation(s.clock, verify.NotObservedRevoked)
}

// --- AT1 / AT3 : identity binding + scope inspectability, end to end -----

func TestAT1_AT3_IssueVerifyAcceptRecoverIdentitiesAndScope(t *testing.T) {
	s := newSystem(t)
	s.clock.Set(nowMid)
	rec := s.issue(t, "read:orders", "write:audit")

	verdict, trace := s.verifier(t, s.freshNotRevoked()).Verify(rec.Presented())
	if !verdict.IsAccept() {
		t.Fatalf("verdict = %s causes=%v, want Accept", verdict.Decision, verdict.Causes)
	}
	a := rec.Read()
	if a.Principal.String() != principalID || a.Delegate.String() != delegateID {
		t.Errorf("identities: %s / %s", a.Principal, a.Delegate)
	}
	if len(a.Scope) != 2 {
		t.Errorf("scope = %v", a.Scope)
	}
	if len(trace.Entries) != 5 {
		t.Errorf("trace must carry five stages, got %d", len(trace.Entries))
	}
}

// --- AT4 : over-scope issuance refused ----------------------------------

func TestAT4_OverScopeIssuanceRefused(t *testing.T) {
	s := newSystem(t)
	res, err := s.authority.Issue(issuance.Request{
		Principal: s.principal, Delegate: s.delegate,
		Scope: []string{"read:orders", "delete:everything"}, Expiration: expiry,
	})
	if err != nil {
		t.Fatal(err)
	}
	if res.Outcome != issuance.Refused || res.Refusal != issuance.OverScope {
		t.Fatalf("outcome=%s refusal=%s, want Refused/OverScope", res.Outcome, res.Refusal)
	}
	if res.Record != nil {
		t.Error("over-scope issuance must create nothing")
	}
}

// --- AT5 : scope tampered post-issuance is detected and rejected --------

func TestAT5_TamperedRecordRejected(t *testing.T) {
	s := newSystem(t)
	s.clock.Set(nowMid)
	rec := s.issue(t, "read:orders")
	tampered := append([]byte(nil), rec.Presented()...)
	tampered[len(tampered)/2] ^= 0x01

	verdict, _ := s.verifier(t, s.freshNotRevoked()).Verify(tampered)
	if verdict.Decision != verify.Reject || !hasCause(verdict.Causes, verify.IntegrityFailed) {
		t.Fatalf("verdict=%s causes=%v, want Reject(IntegrityFailed)", verdict.Decision, verdict.Causes)
	}
}

// --- AT6 / AT7 : non-expired accepted, expired rejected -----------------

func TestAT6_AT7_ExpiryBoundary(t *testing.T) {
	s := newSystem(t)
	rec := s.issue(t, "read:orders")

	s.clock.Set(nowMid) // within window
	if verdict, _ := s.verifier(t, s.freshNotRevoked()).Verify(rec.Presented()); !verdict.IsAccept() {
		t.Fatalf("within window: verdict=%s, want Accept", verdict.Decision)
	}

	s.clock.Set(expiry.Add(time.Hour)) // past window
	verdict, _ := s.verifier(t, s.freshNotRevoked()).Verify(rec.Presented())
	if verdict.Decision != verify.Reject || !hasCause(verdict.Causes, verify.Expired) {
		t.Fatalf("expired: verdict=%s causes=%v, want Reject(Expired)", verdict.Decision, verdict.Causes)
	}
}

// --- AT9 / AT11 : revoked observable rejected, and stays rejected -------

func TestAT9_AT11_RevokedObservableRejectedAndStaysRejected(t *testing.T) {
	s := newSystem(t)
	s.clock.Set(nowMid)
	rec := s.issue(t, "read:orders")
	instance := rec.Read().Instance

	// Author-side: record the revocation in the register (M6).
	s.register.Revoke(instance, s.clock.Now())
	if !s.register.IsRevoked(instance) {
		t.Fatal("register must hold the revocation")
	}

	// RP-side: a revocation provider that observes it (stands in for the
	// spike-selected propagation channel, which is deferred).
	observed := harness.NewRevocation().Set(instance, verify.RevocationStatus{
		State: verify.ObservablyRevoked, AsOf: s.clock.Now(),
	})
	v := s.verifier(t, observed)

	// AT11: rejected at three successive times.
	for i, dt := range []time.Duration{0, time.Minute, 10 * time.Minute} {
		s.clock.Set(nowMid.Add(dt))
		verdict, _ := v.Verify(rec.Presented())
		if verdict.Decision != verify.Reject || !hasCause(verdict.Causes, verify.RevokedObservable) {
			t.Fatalf("time %d: verdict=%s causes=%v, want Reject(RevokedObservable)", i, verdict.Decision, verdict.Causes)
		}
	}
}

// --- AT18 : issuance to an ephemeral delegate ---------------------------

func TestAT18_EphemeralDelegateIssuance(t *testing.T) {
	// The delegate is named, not contacted or provisioned: issuance requires
	// no long-lived static identity for it (ER10). Full ephemeral-lifetime
	// behavior against a running workload is a substrate concern (deferred).
	s := newSystem(t)
	s.clock.Set(nowMid)
	rec := s.issue(t, "read:orders")
	if verdict, _ := s.verifier(t, s.freshNotRevoked()).Verify(rec.Presented()); !verdict.IsAccept() {
		t.Fatalf("ephemeral issuance verify = %s, want Accept", verdict.Decision)
	}
}

// --- AT19 / AT21 : third-party reconstruction from the record alone -----

func TestAT19_AT21_ThirdPartyReconstruction(t *testing.T) {
	s := newSystem(t)
	s.clock.Set(nowMid)
	rec := s.issue(t, "read:orders", "write:audit")

	// A third party holds only the record bytes and the same public trust
	// material — no verifier state, no issuer state.
	tm, _ := s.trust.TrustMaterialFor(s.principal.TrustDomain())
	validated, outcome := record.ValidateIntegrity(rec.Presented(), tm)
	if outcome != record.Intact {
		t.Fatal("reconstruction failed on an authentic record")
	}
	a := validated.Read()
	if a.Principal.String() != principalID || a.Delegate.String() != delegateID || len(a.Scope) != 2 {
		t.Errorf("reconstructed assertions wrong: %+v", a)
	}
	if !a.IssuedAt.Equal(nowMid) {
		t.Errorf("reconstructed issuance time = %v, want %v", a.IssuedAt, nowMid)
	}
}

// --- AT22 : inconclusive verification fails closed [HYPOTHESIS] ----------

func TestAT22_FailClosed_DegenerateProvider(t *testing.T) {
	// The degenerate revocation realization (pre-spike default) answers
	// Indeterminate; wired through the real AD-020 adapter, the system fails
	// closed rather than accepting on absent revocation knowledge. This is
	// the designed fail-closed behavior; V1 documents it, does not warrant it.
	s := newSystem(t)
	s.clock.Set(nowMid)
	rec := s.issue(t, "read:orders")

	rev := harness.AdaptRevocation(revstatus.NewDegenerate())
	verdict, _ := s.verifier(t, rev).Verify(rec.Presented())
	if verdict.IsAccept() {
		t.Fatal("[HYPOTHESIS] fail-closed violated: accepted on indeterminate revocation")
	}
	if verdict.Decision != verify.InconclusiveRejected || !hasCause(verdict.Causes, verify.RevocationStatusIndeterminate) {
		t.Fatalf("verdict=%s causes=%v, want InconclusiveRejected(RevocationStatusIndeterminate)", verdict.Decision, verdict.Causes)
	}
}

// --- AT27 : the system does not issue base identity ---------------------

func TestAT27_SystemDoesNotIssueBaseIdentity(t *testing.T) {
	// Issuance consumes already-issued identities (the principal and delegate
	// SPIFFE IDs are request inputs) and binds a delegation to them; it never
	// mints a base workload identity. A request naming no principal is
	// malformed — the authority cannot conjure one.
	s := newSystem(t)
	res, err := s.authority.Issue(issuance.Request{
		Delegate: s.delegate, Scope: []string{"read:orders"}, Expiration: expiry,
	})
	if err != nil {
		t.Fatal(err)
	}
	if res.Outcome != issuance.Refused || res.Refusal != issuance.MalformedRequest {
		t.Fatalf("issuance without a principal must refuse (no base-identity issuance), got %s/%s", res.Outcome, res.Refusal)
	}
}

// --- substrate-blocked ATs: skip with a named blocker -------------------

func TestSubstrateBlockedATs(t *testing.T) {
	blocked := map[string]string{
		"AT10 revocation-does-not-affect-sibling (two delegations, cross-domain)": "substrate",
		"AT12 revocation without workload restart (running workload)":             "substrate",
		"AT15 verification with network to authorities disabled":                  "substrate",
		"AT16 no egress during verification (packet instrumentation)":             "substrate",
		"AT17 two-domain valid/invalid (two SPIRE domains)":                       "substrate",
		"AT13 revoked-but-observable acceptance-count zero within R (end-to-end)": "S1-scope-act + spike-outcome",
		"AT14 partition-at-revocation honors S4 (end-to-end)":                     "S1-scope-act + spike-outcome",
		"AT24 delegation-verification does not replace RP baseline":               "substrate (needs an existing RP identity-verification baseline)",
		"AT26 verification latency measured (driver measurement point)":           "substrate",
	}
	for name, blocker := range blocked {
		t.Run(name, func(t *testing.T) {
			t.Skipf("blocked on: %s (IMPLEMENTATION_MASTER_PLAN.md AT unblock map)", blocker)
		})
	}
}

// --- helpers ------------------------------------------------------------

func hasCause(causes []verify.Cause, want verify.Cause) bool {
	for _, c := range causes {
		if c == want {
			return true
		}
	}
	return false
}
