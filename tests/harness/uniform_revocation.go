package harness

import (
	"time"

	"github.com/Raj-glitch-max/atlas/internal/record"
	"github.com/Raj-glitch-max/atlas/internal/verify"
)

// UniformRevocation answers the same revocation state for every instance,
// with a knowledge answer's as-of tracking a clock so it stays fresh as time
// advances. It models a deliberate scenario — e.g. "revocation is
// continuously observed as not-revoked" — for end-to-end tests where the
// instance identity is minted at issuance and not known in advance.
//
// This is a configured knowledge state, not laundered ignorance: it is used
// only where the test author asserts the RP's view holds that knowledge. For
// the ignorance case, use the degenerate provider (or NewRevocation with
// nothing configured), which answers Indeterminate (FD-4).
type UniformRevocation struct {
	clock interface{ Now() time.Time }
	state verify.RevocationState
}

// NewUniformRevocation returns a provider answering state for every instance.
// For a knowledge state the as-of is the clock's current reading at query
// time; Indeterminate carries no as-of.
func NewUniformRevocation(clock interface{ Now() time.Time }, state verify.RevocationState) UniformRevocation {
	return UniformRevocation{clock: clock, state: state}
}

// StatusOf implements verify.RevocationStatusPort.
func (u UniformRevocation) StatusOf(_ record.InstanceID) verify.RevocationStatus {
	if u.state == verify.Indeterminate {
		return verify.RevocationStatus{State: verify.Indeterminate}
	}
	return verify.RevocationStatus{State: u.state, AsOf: u.clock.Now()}
}

// checkpoint: chore(security): harden network partition test
