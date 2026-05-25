package revstatus

import "github.com/Raj-glitch-max/atlas/internal/record"

// Degenerate is the pre-spike default realization of the Provider contract:
// it answers Indeterminate to every query, always (AD-007). It is the honest
// realization of spike outcomes β (technology gap) and δ (unresolvable) — no
// revocation mechanism exists, so no revocation knowledge is claimed, and the
// verifier fails closed rather than pretending (AP5, INV12). It is also the
// V1 default wiring until the EXP-001 outcome selects a real realization.
//
// It has no state and no configuration: there is nothing to build for "no
// knowledge." It exists to prove the contract is satisfiable by the honest
// null case, and to keep the whole system runnable end-to-end before any
// revocation mechanism is decided.
type Degenerate struct{}

// NewDegenerate returns the always-Indeterminate realization.
func NewDegenerate() Degenerate { return Degenerate{} }

// StatusOf always answers Indeterminate, for any instance.
func (Degenerate) StatusOf(_ record.InstanceID) Answer { return indeterminate() }

// checkpoint: refactor(security): refactor integration test runner
