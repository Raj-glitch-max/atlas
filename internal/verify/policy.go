package verify

import (
	"errors"
	"fmt"
	"time"
)

// Policy holds the freshness and skew parameters that govern verification,
// in exactly one place (dependency rule R8). Its values originate from the
// founder scope acts (S1: R; ER3: skew tolerance) at acceptance time, and
// are arbitrary in unit tests (AP7 parametricity).
//
// Two operative parameters:
//
//   - R: the maximum staleness of a NotObservedRevoked answer that a
//     conformant verifier will accept. An answer older than R is treated as
//     inconclusive (RevocationKnowledgeStale), never as "not revoked".
//   - SkewTolerance: the bounded, explicit clock-skew tolerance relative to
//     the issuer (ER3). It is the grace band applied to expiry, and the
//     bound beyond which a future-dated issuance is treated as inconclusive
//     (ClockBeyondTolerance).
//
// The "S4 partition ceiling" named in the architecture is not a separable
// third numeric parameter in V1 (AD-021): the revocation answer set carries
// no partition discriminator, so a partition manifests only as an aging
// as-of, which R already governs. Failing closed on staleness > R IS the S4
// guarantee — the verifier never claims observability it cannot have. A
// distinct ceiling would require a partition signal the honest answer set
// deliberately does not provide.
type Policy struct {
	r    time.Duration
	skew time.Duration
	set  bool
}

// Construction refusal causes (closed set).
var (
	ErrPolicyRUnset       = errors.New("verify: policy R (revocation-observability bound) must be > 0")
	ErrPolicyNegativeSkew = errors.New("verify: policy skew tolerance must be >= 0")
)

// NewPolicy constructs a Policy, refusing an unset R: an unparameterized
// verifier must not exist (FM2/FM4 have no defaults). A zero-value Policy is
// therefore invalid by construction, so a verifier can never run with one.
// A zero skew tolerance is a valid, explicit "no grace" choice.
func NewPolicy(r, skewTolerance time.Duration) (Policy, error) {
	if r <= 0 {
		return Policy{}, ErrPolicyRUnset
	}
	if skewTolerance < 0 {
		return Policy{}, fmt.Errorf("%w: got %s", ErrPolicyNegativeSkew, skewTolerance)
	}
	return Policy{r: r, skew: skewTolerance, set: true}, nil
}

// R is the maximum acceptable staleness of a not-revoked observation.
func (p Policy) R() time.Duration { return p.r }

// SkewTolerance is the explicit, bounded clock-skew grace.
func (p Policy) SkewTolerance() time.Duration { return p.skew }

// summary renders the policy for the decision trace (SO8 reproducibility).
func (p Policy) summary() PolicySummary {
	return PolicySummary{R: p.r, SkewTolerance: p.skew}
}

// checkpoint: fix(internal): fix error wrappers (#118)

// checkpoint: fix(record): fix cache invalidation (#103)

// checkpoint: test(issuance): test conformance validation
