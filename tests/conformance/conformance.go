// Package conformance is the executable conformance kit for the Verification
// Core (M3). It turns the prose conformance definition ("a conformant
// verifier is exactly these five checks, these verdicts, this routing") into
// a runnable corpus that ANY implementation — the V1 verifier, a future
// competing implementation, or a third-party port — must pass.
//
// Why this exists (OMEGA-04 impact analysis, journal 2026-07-06): the
// verifier IS the discriminating channel. If two implementations disagree on
// an edge case, the observation is ill-defined and an adversary selects the
// verifier that accepts — the Frankencerts / verifier-differential failure
// that has burned every multi-implementation trust primitive (X.509, TLS,
// JWT). Prose conformance does not survive competing implementations; an
// executable oracle does. This kit is the oracle.
//
// It is test infrastructure (imports the harness), lives under tests/, and is
// imported by no product package. A future implementation runs
// conformance.Run(t, itsFactory) to prove itself conformant.
package conformance

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"testing"
	"time"

	"github.com/Raj-glitch-max/atlas/internal/record"
	"github.com/Raj-glitch-max/atlas/internal/verify"
	"github.com/spiffe/go-spiffe/v2/spiffeid"
)

// Verifier is the contract a conformant verification implementation exposes.
// verify.Verifier satisfies it; a competing implementation would too.
type Verifier interface {
	Verify(presented []byte) (verify.Verdict, verify.DecisionTrace)
}

// Factory builds a verifier from a policy and the three ports. Each
// implementation supplies its own; the V1 factory is verify.NewVerifier.
type Factory func(verify.Policy, verify.TrustMaterialPort, verify.RevocationStatusPort, verify.TimePort) (Verifier, error)

// V1Factory adapts verify.NewVerifier to the Factory contract.
func V1Factory(p verify.Policy, t verify.TrustMaterialPort, r verify.RevocationStatusPort, c verify.TimePort) (Verifier, error) {
	return verify.NewVerifier(p, t, r, c)
}

// Scenario is one fully-specified conformance case: inputs plus the verdict
// every conformant verifier must produce.
type Scenario struct {
	Name         string
	Presented    []byte
	Trust        verify.TrustMaterialPort
	Revocation   verify.RevocationStatusPort
	Clock        verify.TimePort
	Policy       verify.Policy
	WantDecision verify.Decision
	WantCauses   []verify.Cause // causes that MUST be present in the verdict
}

// --- fixed time base (deterministic; no wall clock) ---------------------

var (
	tIssued = time.Unix(1_800_000_000, 0).UTC()
	tExpiry = time.Unix(1_800_003_600, 0).UTC() // +1h
	tMid    = time.Unix(1_800_000_300, 0).UTC() // +5m, in-window
	tPast   = time.Unix(1_800_010_000, 0).UTC() // well past expiry
)

// --- port fakes (self-contained; the kit does not import the harness so it
// stays a minimal, dependency-light oracle a third party can vendor) -----

type trustStore struct {
	byDomain map[spiffeid.TrustDomain]record.TrustMaterial
}

func (s trustStore) TrustMaterialFor(d spiffeid.TrustDomain) (record.TrustMaterial, bool) {
	m, ok := s.byDomain[d]
	return m, ok
}

type fixedClock struct{ t time.Time }

func (c fixedClock) Now() time.Time { return c.t }

type uniformRevocation struct {
	state verify.RevocationState
	asOf  time.Time
}

func (r uniformRevocation) StatusOf(record.InstanceID) verify.RevocationStatus {
	if r.state == verify.Indeterminate {
		return verify.RevocationStatus{State: verify.Indeterminate}
	}
	return verify.RevocationStatus{State: r.state, AsOf: r.asOf}
}

// Kit mints records and trust material under one authority key, so scenarios
// are self-contained and reproducible.
type Kit struct {
	signer     record.Signer
	trust      trustStore
	emptyTrust trustStore
	domain     spiffeid.TrustDomain
	principal  spiffeid.ID
	delegate   spiffeid.ID
}

// NewKit builds a kit with a fresh P-256 authority key and matching trust
// material.
func NewKit() (*Kit, error) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	}
	const keyID = "conformance-key-1"
	domain := spiffeid.RequireTrustDomainFromString("domain-a.test")
	tm, err := record.NewTrustMaterial(domain, map[string]*ecdsa.PublicKey{keyID: &key.PublicKey})
	if err != nil {
		return nil, err
	}
	return &Kit{
		signer:     record.Signer{Key: key, KeyID: keyID},
		trust:      trustStore{byDomain: map[spiffeid.TrustDomain]record.TrustMaterial{domain: tm}},
		emptyTrust: trustStore{byDomain: map[spiffeid.TrustDomain]record.TrustMaterial{}},
		domain:     domain,
		principal:  spiffeid.RequireFromString("spiffe://domain-a.test/principal"),
		delegate:   spiffeid.RequireFromString("spiffe://domain-a.test/delegate"),
	}, nil
}

// Seal mints a record with the given shape. selfDelegate binds principal to
// itself (for the binding-mismatch case).
func (k *Kit) Seal(scope []string, iat, exp time.Time, instance string, selfDelegate bool) []byte {
	del := k.delegate
	if selfDelegate {
		del = k.principal
	}
	inst, err := record.InstanceIDFromString(instance)
	if err != nil {
		panic("conformance: bad instance id: " + err.Error())
	}
	rec, err := record.Seal(record.Assertions{
		Principal: k.principal, Delegate: del, Scope: scope,
		Expiration: exp, IssuedAt: iat, Instance: inst,
	}, k.signer)
	if err != nil {
		panic("conformance: seal failed: " + err.Error())
	}
	return rec.Presented()
}

// Valid mints the canonical valid record used as the mutation baseline.
func (k *Kit) Valid() []byte {
	return k.Seal([]string{"read:orders", "write:audit"}, tIssued, tExpiry, "inst-valid", false)
}

func mustPolicy() verify.Policy {
	p, err := verify.NewPolicy(time.Minute, 30*time.Second)
	if err != nil {
		panic("conformance: policy: " + err.Error())
	}
	return p
}

// Exported builders so property tests and future implementations' differential
// harnesses can construct the same ports the corpus uses.

// DefaultPolicy is the corpus policy (R=1m, skew=30s).
func DefaultPolicy() verify.Policy { return mustPolicy() }

// MidTime is the in-window reference instant.
func MidTime() time.Time { return tMid }

// Trust returns the kit's trust-material port (holds the authority key).
func (k *Kit) Trust() verify.TrustMaterialPort { return k.trust }

// EmptyTrust returns a trust-material port that holds nothing.
func (k *Kit) EmptyTrust() verify.TrustMaterialPort { return k.emptyTrust }

// Clock returns a fixed-time port.
func Clock(t time.Time) verify.TimePort { return fixedClock{t} }

// FreshRevocation returns a not-revoked port with the given as-of.
func FreshRevocation(asOf time.Time) verify.RevocationStatusPort {
	return uniformRevocation{state: verify.NotObservedRevoked, asOf: asOf}
}

// RevocationInState returns a uniform port in the given state/as-of.
func RevocationInState(s verify.RevocationState, asOf time.Time) verify.RevocationStatusPort {
	return uniformRevocation{state: s, asOf: asOf}
}

// Decode validates presented bytes against the kit's trust material and, if
// intact, returns the record's assertions. It is the content oracle for
// differential/property tests: two byte-different presentations that both
// decode to equal assertions are the same delegation (the record's transport
// encoding is non-canonical — base64 padding-bit malleable — so byte identity
// is weaker than content identity; content integrity is what INV8 guarantees).
func (k *Kit) Decode(presented []byte) (record.Assertions, bool) {
	tm, ok := k.trust.byDomain[k.domain]
	if !ok {
		return record.Assertions{}, false
	}
	rec, outcome := record.ValidateIntegrity(presented, tm)
	if outcome != record.Intact {
		return record.Assertions{}, false
	}
	return rec.Read(), true
}

// BuildCorpus returns the full conformance corpus: exhaustive over the verdict
// space — Accept, every definitive cause, every inconclusive cause, and the
// definitive-dominates-inconclusive precedence. This IS the executable
// conformance definition; adding a case is an interface-spec amendment.
func BuildCorpus(k *Kit) []Scenario {
	p := mustPolicy()
	fresh := uniformRevocation{state: verify.NotObservedRevoked, asOf: tMid}
	indet := uniformRevocation{state: verify.Indeterminate}
	revoked := uniformRevocation{state: verify.ObservablyRevoked, asOf: tMid}
	stale := uniformRevocation{state: verify.NotObservedRevoked, asOf: tIssued.Add(-time.Hour)} // > R before tMid

	valid := k.Valid()
	tampered := append([]byte(nil), valid...)
	tampered[len(tampered)/2] ^= 0x01

	return []Scenario{
		{
			Name:      "accept: valid, fresh not-revoked, in-window",
			Presented: valid, Trust: k.trust, Revocation: fresh, Clock: fixedClock{tMid}, Policy: p,
			WantDecision: verify.Accept,
		},
		{
			Name:      "reject: expired",
			Presented: valid, Trust: k.trust, Revocation: fresh, Clock: fixedClock{tPast}, Policy: p,
			WantDecision: verify.Reject, WantCauses: []verify.Cause{verify.Expired},
		},
		{
			Name:      "reject: integrity (tampered byte)",
			Presented: tampered, Trust: k.trust, Revocation: fresh, Clock: fixedClock{tMid}, Policy: p,
			WantDecision: verify.Reject, WantCauses: []verify.Cause{verify.IntegrityFailed},
		},
		{
			Name:      "reject: binding mismatch (self-delegation)",
			Presented: k.Seal([]string{"read:orders"}, tIssued, tExpiry, "inst-self", true),
			Trust:     k.trust, Revocation: fresh, Clock: fixedClock{tMid}, Policy: p,
			WantDecision: verify.Reject, WantCauses: []verify.Cause{verify.BindingMismatch},
		},
		{
			Name:      "reject: observably revoked",
			Presented: valid, Trust: k.trust, Revocation: revoked, Clock: fixedClock{tMid}, Policy: p,
			WantDecision: verify.Reject, WantCauses: []verify.Cause{verify.RevokedObservable},
		},
		{
			Name:      "inconclusive: trust material absent",
			Presented: valid, Trust: k.emptyTrust, Revocation: fresh, Clock: fixedClock{tMid}, Policy: p,
			WantDecision: verify.InconclusiveRejected, WantCauses: []verify.Cause{verify.TrustMaterialAbsent},
		},
		{
			Name:      "inconclusive: revocation indeterminate (degenerate)",
			Presented: valid, Trust: k.trust, Revocation: indet, Clock: fixedClock{tMid}, Policy: p,
			WantDecision: verify.InconclusiveRejected, WantCauses: []verify.Cause{verify.RevocationStatusIndeterminate},
		},
		{
			Name:      "inconclusive: revocation knowledge stale (> R)",
			Presented: valid, Trust: k.trust, Revocation: stale, Clock: fixedClock{tMid}, Policy: p,
			WantDecision: verify.InconclusiveRejected, WantCauses: []verify.Cause{verify.RevocationKnowledgeStale},
		},
		{
			Name:      "inconclusive: clock beyond tolerance (future issuance)",
			Presented: k.Seal([]string{"read:orders"}, tMid.Add(time.Hour), tMid.Add(2*time.Hour), "inst-future", false),
			Trust:     k.trust, Revocation: fresh, Clock: fixedClock{tMid}, Policy: p,
			WantDecision: verify.InconclusiveRejected, WantCauses: []verify.Cause{verify.ClockBeyondTolerance},
		},
		{
			Name:      "precedence: definitive dominates inconclusive (expired + indeterminate)",
			Presented: valid, Trust: k.trust, Revocation: indet, Clock: fixedClock{tPast}, Policy: p,
			WantDecision: verify.Reject, WantCauses: []verify.Cause{verify.Expired},
		},
	}
}

// Run executes the full conformance corpus against a verifier implementation.
// It asserts, for every scenario: the exact decision, the presence of every
// required cause, a five-stage trace, and determinism (two verifications of
// the same input agree). A future implementation calls this to prove itself.
func Run(t *testing.T, factory Factory) {
	t.Helper()
	k, err := NewKit()
	if err != nil {
		t.Fatalf("conformance kit: %v", err)
	}
	for _, s := range BuildCorpus(k) {
		s := s
		t.Run(s.Name, func(t *testing.T) {
			v, err := factory(s.Policy, s.Trust, s.Revocation, s.Clock)
			if err != nil {
				t.Fatalf("factory: %v", err)
			}
			verdict, trace := v.Verify(s.Presented)

			if verdict.Decision != s.WantDecision {
				t.Fatalf("decision = %s (causes %v), want %s", verdict.Decision, verdict.Causes, s.WantDecision)
			}
			for _, want := range s.WantCauses {
				if !containsCause(verdict.Causes, want) {
					t.Fatalf("causes = %v, want to contain %s", verdict.Causes, want)
				}
			}
			if len(trace.Entries) != 5 {
				t.Errorf("trace = %d stages, want 5 (unconditional)", len(trace.Entries))
			}
			// Determinism: a re-presentation is a fresh verification with the
			// identical verdict (no memory) — required for the invariant to
			// hold across repeated observation.
			verdict2, _ := v.Verify(s.Presented)
			if verdict2.Decision != verdict.Decision || !sameCauses(verdict.Causes, verdict2.Causes) {
				t.Errorf("non-deterministic: %s%v then %s%v", verdict.Decision, verdict.Causes, verdict2.Decision, verdict2.Causes)
			}
		})
	}
}

func containsCause(cs []verify.Cause, want verify.Cause) bool {
	for _, c := range cs {
		if c == want {
			return true
		}
	}
	return false
}

func sameCauses(a, b []verify.Cause) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// PresentedEqual reports whether two presented byte slices are identical —
// a helper for differential harnesses comparing inputs across implementations.
func PresentedEqual(a, b []byte) bool { return bytes.Equal(a, b) }
