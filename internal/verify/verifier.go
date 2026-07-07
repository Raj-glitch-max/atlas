package verify

import "errors"

// Verifier is the conformant Verification Core (M3): the locus at which the
// system's guarantees are defined (FM10). Conformance IS this type's
// behavior — the five checks, the closed cause and verdict sets, the routing,
// and the unconditional trace. A verifier holds only injected configuration
// (policy and ports); it keeps no per-verification state, so a re-presentation
// is a fresh verification (RFC-002 §9.2) and concurrent Verify calls are safe
// (the core is pure over its inputs; AD-017).

// Construction refusal causes (closed set). An unparameterized or
// unwireable verifier must not exist.
var (
	ErrNoPolicy            = errors.New("verify: policy is unset (construct with NewPolicy)")
	ErrNoTrustMaterialPort = errors.New("verify: trust-material port is nil")
	ErrNoRevocationPort    = errors.New("verify: revocation-status port is nil")
	ErrNoTimePort          = errors.New("verify: time port is nil")
)

// Verifier verifies presented delegation records against injected inputs.
type Verifier struct {
	policy Policy
	trust  TrustMaterialPort
	revoc  RevocationStatusPort
	clock  TimePort
}

// NewVerifier constructs a Verifier, refusing an unset policy (FM2/FM4 have
// no defaults) and any nil port (the core performs no I/O, so every input
// source must be supplied). The refusal is total: there is no way to obtain
// a Verifier that could run without R, a skew tolerance, and all three ports.
func NewVerifier(policy Policy, trust TrustMaterialPort, revoc RevocationStatusPort, clock TimePort) (*Verifier, error) {
	if !policy.set {
		return nil, ErrNoPolicy
	}
	if trust == nil {
		return nil, ErrNoTrustMaterialPort
	}
	if revoc == nil {
		return nil, ErrNoRevocationPort
	}
	if clock == nil {
		return nil, ErrNoTimePort
	}
	return &Verifier{policy: policy, trust: trust, revoc: revoc, clock: clock}, nil
}

// Verify determines the verdict for a presented record and emits a decision
// trace unconditionally (Accepts included).
//
// Execution order: integrity runs first as the gate, because the other four
// stages require an authenticated read (a record's assertions cannot be
// trusted before its integrity is established). If integrity does not pass,
// the read-dependent stages are recorded NotEvaluated (honest, not omitted).
// The trace presents the five stages in canonical label order regardless.
//
// Verdict routing is order-independent: all produced causes are collected and
// a single precedence applies — any definitive cause yields Reject; otherwise
// any inconclusive cause yields InconclusiveRejected [HYPOTHESIS]; otherwise
// Accept. Thus forcing any single check to fail while others pass flips the
// verdict away from Accept (SO5).
func (v *Verifier) Verify(presented []byte) (Verdict, DecisionTrace) {
	now := v.clock.Now()

	rec, integrity := checkIntegrity(presented, v.trust)

	var (
		binding, expiry, scope, revocation CheckEntry
		status                             RevocationStatus
	)
	if rec != nil {
		a := rec.Read()
		binding = checkBinding(a)
		expiry = checkExpiry(a, now, v.policy.skew)
		scope = checkScope(a)
		status = v.revoc.StatusOf(a.Instance)
		revocation = checkRevocation(a.Instance, status, now, v.policy.r, v.policy.skew)
	} else {
		const gated = "not evaluated: gated by integrity"
		binding = notEvaluated(CheckIdentityBinding, gated)
		expiry = notEvaluated(CheckExpiry, gated)
		scope = notEvaluated(CheckScope, gated)
		revocation = notEvaluated(CheckRevocation, gated)
	}

	// Canonical label order for the trace: binding(1), integrity(2),
	// expiry(3), scope(4), revocation(5).
	entries := []CheckEntry{binding, integrity, expiry, scope, revocation}
	verdict := route(entries)

	trace := DecisionTrace{
		Entries:                entries,
		TimeReading:            now,
		RevocationState:        status.State.String(),
		RevocationObservedAsOf: status.AsOf,
		Policy:                 v.policy.summary(),
		Verdict:                verdict,
	}
	return verdict, trace
}

// route computes the verdict from the collected stage outcomes, applying the
// single precedence rule (definitive dominates inconclusive dominates pass).
// It is order-independent: it inspects causes, not sequence.
func route(entries []CheckEntry) Verdict {
	var definitive, inconclusiveList []Cause
	for _, e := range entries {
		switch e.Outcome {
		case OutcomeFailDefinitive:
			definitive = append(definitive, e.Cause)
		case OutcomeInconclusive:
			inconclusiveList = append(inconclusiveList, e.Cause)
		}
	}
	switch {
	case len(definitive) > 0:
		return Verdict{Decision: Reject, Causes: definitive}
	case len(inconclusiveList) > 0:
		return Verdict{Decision: InconclusiveRejected, Causes: inconclusiveList}
	default:
		return Verdict{Decision: Accept}
	}
}
