package verify_test

import (
	"testing"
	"time"

	"github.com/Raj-glitch-max/atlas/internal/record"
	"github.com/Raj-glitch-max/atlas/internal/verify"
	"github.com/Raj-glitch-max/atlas/tests/harness"
)

// AT22 [HYPOTHESIS] — fail-closed under inconclusive verification. This is
// the designed fail-closed posture (NFR3/ER11/SO4/C-INV1); V1 documents the
// behavior, it is NOT a warranted guarantee (DR7). These tests assert the
// observed routing, and their names carry the [HYPOTHESIS] marker so nothing
// here reads as a committed property.

func TestHypothesisFailClosed_InconclusiveNeverAccepts(t *testing.T) {
	// Every inconclusive condition, in isolation, must route away from
	// Accept to InconclusiveRejected (a rejection).
	conditions := []struct {
		name  string
		setup func(t *testing.T, w *world) *record.Record
		cause verify.Cause
	}{
		{
			name: "trust material unavailable",
			setup: func(t *testing.T, w *world) *record.Record {
				rec := w.seal(t, nil)
				w.trust = harness.NewTrustStore() // hold nothing
				return rec
			},
			cause: verify.TrustMaterialAbsent,
		},
		{
			name: "revocation indeterminate",
			setup: func(t *testing.T, w *world) *record.Record {
				return w.seal(t, nil) // provider answers Indeterminate
			},
			cause: verify.RevocationStatusIndeterminate,
		},
		{
			name: "clock beyond tolerance",
			setup: func(t *testing.T, w *world) *record.Record {
				w.freshNotRevoked()
				return w.seal(t, func(a *record.Assertions) {
					a.IssuedAt = baseNow.Add(time.Hour)
					a.Expiration = baseNow.Add(2 * time.Hour)
				})
			},
			cause: verify.ClockBeyondTolerance,
		},
	}
	for _, c := range conditions {
		t.Run(c.name, func(t *testing.T) {
			w := newWorld(t)
			rec := c.setup(t, w)
			v := w.verifier(t, time.Minute, time.Second)
			verdict, _ := v.Verify(rec.Presented())
			if verdict.IsAccept() {
				t.Fatalf("[HYPOTHESIS] fail-closed violated: inconclusive condition ACCEPTED")
			}
			if verdict.Decision != verify.InconclusiveRejected {
				t.Fatalf("verdict=%s, want InconclusiveRejected [HYPOTHESIS]", verdict.Decision)
			}
			if !hasCause(verdict.Causes, c.cause) {
				t.Fatalf("causes=%v, want %s", verdict.Causes, c.cause)
			}
		})
	}
}

func TestDefinitiveDominatesInconclusive(t *testing.T) {
	// When both a definitive failure and an inconclusive condition are
	// present, the verdict is Reject (definitive), not InconclusiveRejected:
	// routing is order-independent and definitive-dominant. Here: an expired
	// record (definitive) whose revocation is also indeterminate.
	w := newWorld(t)
	rec := w.seal(t, nil) // revocation left Indeterminate
	w.clock.Set(baseExpiry.Add(time.Hour))
	v := w.verifier(t, time.Minute, time.Second)

	verdict, _ := v.Verify(rec.Presented())
	if verdict.Decision != verify.Reject {
		t.Fatalf("verdict=%s, want Reject (definitive dominates inconclusive)", verdict.Decision)
	}
	if !hasCause(verdict.Causes, verify.Expired) {
		t.Fatalf("causes=%v, want to contain Expired", verdict.Causes)
	}
}

func TestCauseClassificationIsTotalAndDisjoint(t *testing.T) {
	// Every named cause (except CauseNone) is exactly one of definitive or
	// inconclusive — the closed set carries no unclassified or double-
	// classified member (FM11: no ambiguous outcome).
	all := []verify.Cause{
		verify.BindingMismatch, verify.IntegrityFailed, verify.Expired,
		verify.ScopeIntegrityFailed, verify.RevokedObservable,
		verify.TrustMaterialAbsent, verify.ClockBeyondTolerance,
		verify.RevocationStatusIndeterminate, verify.RevocationKnowledgeStale,
		verify.SignatureUnverifiable,
	}
	for _, c := range all {
		if c.IsDefinitive() == c.IsInconclusive() {
			t.Errorf("%s must be exactly one of definitive/inconclusive (def=%v inc=%v)",
				c, c.IsDefinitive(), c.IsInconclusive())
		}
	}
	if verify.CauseNone.IsDefinitive() || verify.CauseNone.IsInconclusive() {
		t.Error("CauseNone must be neither definitive nor inconclusive")
	}
}
