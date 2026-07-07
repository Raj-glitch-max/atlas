package contracttest

import (
	"testing"
	"time"

	"github.com/Raj-glitch-max/atlas/internal/record"
	"github.com/Raj-glitch-max/atlas/internal/revstatus"
)

// Run executes the realization-independent contract suite against a Provider
// (AD-008, plugin boundary P1; INTERFACE_SPECIFICATION.md §5). No realization
// — the degenerate one, an EXP-001 spike candidate, or a production
// composition — is wired into a composition root unless it passes this suite.
// It is the same seam for all of them, so every candidate the spike attempts
// is judged by identical, pre-registered criteria.
//
// It is a reusable helper (not a _test.go file) so each realization's test
// package can call it: contracttest.Run(t, myRealization, knownInstances...).
//
// The suite asserts the invariants every realization must honor:
//
//  1. Closed answer set — every answer's state is one of the three members.
//  2. Mandatory as-of on knowledge answers — NotObservedRevoked and
//     ObservablyRevoked carry a non-zero AsOf; Indeterminate does not claim one.
//  3. Honest-indeterminate — an unknown instance is never answered
//     NotObservedRevoked; ignorance is Indeterminate (AP5, FD-4).
//  4. Determinism — repeated queries against an unchanged view agree.
//
// knownInstances are instances the realization is expected to have knowledge
// about (may be empty, e.g. for the degenerate realization); unknownInstance
// is one the realization is NOT expected to know, used for the
// honest-indeterminate check.
func Run(t *testing.T, p revstatus.Provider, unknownInstance record.InstanceID, knownInstances ...record.InstanceID) {
	t.Helper()

	all := append([]record.InstanceID{unknownInstance}, knownInstances...)

	t.Run("closed answer set", func(t *testing.T) {
		for _, inst := range all {
			ans := p.StatusOf(inst)
			switch ans.State {
			case revstatus.Indeterminate, revstatus.NotObservedRevoked, revstatus.ObservablyRevoked:
			default:
				t.Fatalf("instance %q: state %d is outside the closed answer set", inst, ans.State)
			}
		}
	})

	t.Run("mandatory as-of on knowledge answers", func(t *testing.T) {
		for _, inst := range all {
			ans := p.StatusOf(inst)
			switch ans.State {
			case revstatus.NotObservedRevoked, revstatus.ObservablyRevoked:
				if ans.AsOf.IsZero() {
					t.Errorf("instance %q: %s answer must carry a non-zero as-of", inst, ans.State)
				}
			case revstatus.Indeterminate:
				if !ans.AsOf.IsZero() {
					t.Errorf("instance %q: Indeterminate must not claim an as-of (got %v)", inst, ans.AsOf)
				}
			}
		}
	})

	t.Run("honest-indeterminate for unknown instance", func(t *testing.T) {
		ans := p.StatusOf(unknownInstance)
		if ans.State == revstatus.NotObservedRevoked {
			t.Fatalf("unknown instance answered NotObservedRevoked — ignorance laundered into knowledge (AP5/FD-4)")
		}
	})

	t.Run("determinism under an unchanged view", func(t *testing.T) {
		for _, inst := range all {
			first := p.StatusOf(inst)
			for i := 0; i < 20; i++ {
				got := p.StatusOf(inst)
				if got.State != first.State || !got.AsOf.Equal(first.AsOf) {
					t.Fatalf("instance %q: non-deterministic answer (%v then %v)", inst, first, got)
				}
			}
		}
	})

	t.Run("as-of never dated in an implausible future relative to a knowledge answer", func(t *testing.T) {
		// A realization may legitimately report an as-of up to "now"; it must
		// not fabricate a future currency. This is a sanity floor on the
		// honesty of freshness, not a policy judgment (policy is the verifier's).
		ceiling := time.Now().Add(time.Minute)
		for _, inst := range all {
			ans := p.StatusOf(inst)
			if ans.State != revstatus.Indeterminate && ans.AsOf.After(ceiling) {
				t.Errorf("instance %q: as-of %v is implausibly in the future", inst, ans.AsOf)
			}
		}
	})
}
