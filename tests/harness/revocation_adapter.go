package harness

import (
	"github.com/Raj-glitch-max/atlas/internal/record"
	"github.com/Raj-glitch-max/atlas/internal/revstatus"
	"github.com/Raj-glitch-max/atlas/internal/verify"
)

// AdaptRevocation bridges a revstatus.Provider (M5) onto the verifier's
// RevocationStatusPort (M3), the composition-root glue described in AD-020.
// verify may not import revstatus and revstatus may not import verify
// (dependency rule R3), so the two carry structurally-identical but distinct
// answer vocabularies; this adapter maps one to the other explicitly (not by
// relying on iota coincidence). A composition root (a cmd driver, or a test)
// wires a chosen revstatus realization to a verifier through this.
func AdaptRevocation(p revstatus.Provider) verify.RevocationStatusPort {
	return revocationAdapter{provider: p}
}

type revocationAdapter struct {
	provider revstatus.Provider
}

func (a revocationAdapter) StatusOf(instance record.InstanceID) verify.RevocationStatus {
	ans := a.provider.StatusOf(instance)
	return verify.RevocationStatus{State: mapState(ans.State), AsOf: ans.AsOf}
}

// mapState is a total, explicit mapping between the two answer-state sets.
// A new member on either side breaks compilation here rather than silently
// mis-mapping (the default panics on an unknown source state).
func mapState(s revstatus.State) verify.RevocationState {
	switch s {
	case revstatus.Indeterminate:
		return verify.Indeterminate
	case revstatus.NotObservedRevoked:
		return verify.NotObservedRevoked
	case revstatus.ObservablyRevoked:
		return verify.ObservablyRevoked
	default:
		panic("harness: unmapped revstatus.State — adapter must be total")
	}
}
