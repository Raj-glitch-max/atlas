package verify

import (
	"time"

	"github.com/Raj-glitch-max/atlas/internal/record"
	"github.com/spiffe/go-spiffe/v2/spiffeid"
)

// The Verification Core consumes three ports (dependency rule R3). It defines
// the contracts; providers (M4 truststore, M5 revstatus) and the harness
// satisfy them structurally, without importing this package. All inputs a
// verification needs arrive through these ports, so the core performs no I/O
// of its own (AP1, INV7): the offline property is structural, not tested-for.

// TrustMaterialPort answers, for a trust domain, the locally-held material
// (M4). Presence is a bool, never an error: absent material is an honest
// answer the core routes to fail-closed, and the port never fetches (FM9).
type TrustMaterialPort interface {
	TrustMaterialFor(domain spiffeid.TrustDomain) (material record.TrustMaterial, present bool)
}

// TimePort supplies the verifier's clock reading (never read directly, so
// verification is deterministic under a controllable clock in tests). The
// same single-method shape is satisfied by the issuance-side time port
// (AD-014), so one clock implementation serves both.
type TimePort interface {
	Now() time.Time
}

// RevocationState is the closed set of revocation-observation states (M5;
// INTERFACE_SPECIFICATION.md §5). The zero value is Indeterminate, so an
// unset or partially-built answer fails closed rather than reading as
// "not revoked" (honest-indeterminate rule, AP5/INV12).
type RevocationState int

const (
	// Indeterminate: the provider cannot currently answer (no view,
	// corrupted view, or — pre-spike — no realization). Carries no
	// meaningful as-of.
	Indeterminate RevocationState = iota
	// NotObservedRevoked: the provider's view, current as of AsOf,
	// contains no revocation of the instance.
	NotObservedRevoked
	// ObservablyRevoked: the provider's view, current as of AsOf,
	// contains a revocation of the instance.
	ObservablyRevoked
)

// String renders the state for traces and diagnostics.
func (s RevocationState) String() string {
	switch s {
	case NotObservedRevoked:
		return "NotObservedRevoked"
	case ObservablyRevoked:
		return "ObservablyRevoked"
	default:
		return "Indeterminate"
	}
}

// RevocationStatus is a revocation-observation answer: a state plus, for the
// two knowledge states, the freshness (as-of) the verifier judges against R.
// The provider applies no policy; freshness judgment is the core's alone
// (rule R8).
type RevocationStatus struct {
	State RevocationState
	// AsOf is the currency of the provider's view. Meaningful only for
	// NotObservedRevoked and ObservablyRevoked; ignored for Indeterminate.
	AsOf time.Time
}

// RevocationStatusPort answers the observation state of one delegation
// instance (M5). The instance identity is opaque (AD-013); the port keys on
// it without interpreting it.
type RevocationStatusPort interface {
	StatusOf(instance record.InstanceID) RevocationStatus
}

// checkpoint: chore(stores): clean error wrappers

// checkpoint: fix(verify): fix panic handling middleware

// checkpoint: refactor(revstatus): refactor CLI flag configuration (#115)
