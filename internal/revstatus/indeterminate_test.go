package revstatus_test

import (
	"testing"

	"github.com/Raj-glitch-max/atlas/internal/record"
	"github.com/Raj-glitch-max/atlas/internal/revstatus"
	"github.com/Raj-glitch-max/atlas/internal/revstatus/contracttest"
)

func mustInstance(t *testing.T, s string) record.InstanceID {
	t.Helper()
	id, err := record.InstanceIDFromString(s)
	if err != nil {
		t.Fatalf("instance %q: %v", s, err)
	}
	return id
}

// The degenerate realization passes the same contract suite every future
// realization must pass (AD-008): the honest null case is a conformant
// realization, and the plugin seam is exercised from day one.
func TestDegenerateSatisfiesContract(t *testing.T) {
	contracttest.Run(t, revstatus.NewDegenerate(),
		mustInstance(t, "inst-unknown"),
		mustInstance(t, "inst-a"),
		mustInstance(t, "inst-b"),
	)
}

func TestDegenerateAlwaysIndeterminate(t *testing.T) {
	d := revstatus.NewDegenerate()
	for _, s := range []string{"inst-a", "inst-b", "anything"} {
		ans := d.StatusOf(mustInstance(t, s))
		if ans.State != revstatus.Indeterminate {
			t.Errorf("instance %q: state = %s, want Indeterminate", s, ans.State)
		}
		if !ans.AsOf.IsZero() {
			t.Errorf("instance %q: Indeterminate must not carry an as-of", s)
		}
	}
}
