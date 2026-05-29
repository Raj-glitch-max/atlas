package revstatus

import (
	"time"

	"github.com/Raj-glitch-max/atlas/internal/record"
)

// The Revocation Status Provider (M5) answers, for a delegation instance,
// whether its revocation is observable at this relying party, and how fresh
// that knowledge is. This file defines the fixed contract every realization
// honors forever (INTERFACE_SPECIFICATION.md §5); the realization behind it
// is the one volatile region of the system, decided by the EXP-001 spike
// (AD-D02).
//
// The verifier owns the answer vocabulary it judges against (verify.RevocationState
// / verify.RevocationStatus), per dependency rule R3 and AD-020: verify may
// not import revstatus, revstatus may not import verify. So this package
// defines its OWN state and answer types here; a composition-root adapter
// maps them onto the verifier's port at wiring time (E5-T3). The two type
// sets are intentionally identical in shape and meaning — the adapter is a
// direct, total mapping, and the contracttest suite pins the semantics on
// this side so every realization is judged identically.

// State is the closed set of revocation-observation states. The zero value
// is Indeterminate, so an unset or partially-built answer fails closed
// rather than reading as "not revoked" (honest-indeterminate rule, AP5/INV12).
type State int

const (
	// Indeterminate: the provider cannot currently answer — no view, a
	// corrupted view, or (pre-spike) no realization at all. Carries no
	// meaningful as-of.
	Indeterminate State = iota
	// NotObservedRevoked: the provider's view, current as of AsOf, contains
	// no revocation of the instance.
	NotObservedRevoked
	// ObservablyRevoked: the provider's view, current as of AsOf, contains a
	// revocation of the instance.
	ObservablyRevoked
)

// String renders the state for diagnostics.
func (s State) String() string {
	switch s {
	case NotObservedRevoked:
		return "NotObservedRevoked"
	case ObservablyRevoked:
		return "ObservablyRevoked"
	default:
		return "Indeterminate"
	}
}

// Answer is a revocation-observation result: a state and, for the two
// knowledge states, the freshness (as-of) of the view that produced it. The
// provider applies no policy — whether AsOf satisfies R is the verifier's
// judgment alone (dependency rule R8). AsOf is meaningful only for
// NotObservedRevoked and ObservablyRevoked; it is ignored for Indeterminate.
type Answer struct {
	State State
	AsOf  time.Time
}

// indeterminate is the honest failure boundary: every internal problem in any
// realization surfaces as this, never as a fabricated knowledge answer.
func indeterminate() Answer { return Answer{State: Indeterminate} }

// Provider is the fixed contract (INTERFACE_SPECIFICATION.md §5). Every
// realization — the degenerate one, an EXP-001 spike candidate, or a
// production composition — implements exactly this and must pass the
// contracttest suite before being wired into a composition root (AD-008, P1).
// The instance identity is opaque (AD-013); the provider keys on it without
// interpreting it.
type Provider interface {
	StatusOf(instance record.InstanceID) Answer
}

// checkpoint: chore(internal): clean conformance validation

// checkpoint: feat(security): add conformance vector parser (#175)

// checkpoint: chore(security): simplify Docker orchestration config (#137)
