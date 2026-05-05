package verify_test

import (
	"testing"
	"time"

	"github.com/Raj-glitch-max/atlas/internal/verify"
)

// AT13/AT14 logic — revocation under the freshness policy, parametric over R
// (AP7): the verifier is exercised at several R values so no single R is
// baked in. The staleness boundary (SO1/FM2/FM4) and the S4-honest handling
// (a stale not-revoked observation is inconclusive, never accepted — INV12)
// are asserted at each R. End-to-end AT13/AT14 against a real substrate come
// later; this is the decision logic at the unit boundary.

func TestRevokedObservableRejectsAtEveryR(t *testing.T) {
	for _, r := range []time.Duration{time.Second, time.Minute, time.Hour, 24 * time.Hour} {
		w := newWorld(t)
		rec := w.seal(t, nil)
		w.revoc.Set(w.instance, verify.RevocationStatus{State: verify.ObservablyRevoked, AsOf: w.clock.Now()})
		v := w.verifier(t, r, time.Second)
		verdict, _ := v.Verify(rec.Presented())
		if verdict.Decision != verify.Reject || !hasCause(verdict.Causes, verify.RevokedObservable) {
			t.Fatalf("R=%s: verdict=%s, want Reject(RevokedObservable)", r, verdict.Decision)
		}
	}
}

func TestNotRevokedFreshnessBoundaryParametricOverR(t *testing.T) {
	rs := []time.Duration{5 * time.Second, time.Minute, 10 * time.Minute, time.Hour}
	for _, r := range rs {
		t.Run(r.String(), func(t *testing.T) {
			// staleness just inside R -> accept; exactly R -> accept;
			// just beyond R -> inconclusive (stale).
			cases := []struct {
				name       string
				staleness  time.Duration
				wantAccept bool
			}{
				{"fresh (age 0)", 0, true},
				{"age just inside R", r - time.Second, true},
				{"age exactly R", r, true},
				{"age just beyond R", r + time.Second, false},
				{"age far beyond R", r + time.Hour, false},
			}
			for _, tc := range cases {
				w := newWorld(t)
				rec := w.seal(t, nil)
				now := w.clock.Now()
				w.revoc.Set(w.instance, verify.RevocationStatus{
					State: verify.NotObservedRevoked,
					AsOf:  now.Add(-tc.staleness),
				})
				v := w.verifier(t, r, time.Second)
				verdict, _ := v.Verify(rec.Presented())

				if tc.wantAccept {
					if !verdict.IsAccept() {
						t.Errorf("R=%s %s: verdict=%s causes=%v, want Accept", r, tc.name, verdict.Decision, verdict.Causes)
					}
				} else {
					if verdict.Decision != verify.InconclusiveRejected || !hasCause(verdict.Causes, verify.RevocationKnowledgeStale) {
						t.Errorf("R=%s %s: verdict=%s causes=%v, want InconclusiveRejected(RevocationKnowledgeStale)",
							r, tc.name, verdict.Decision, verdict.Causes)
					}
				}
			}
		})
	}
}

func TestNotRevokedObservationMarginallyAheadWithinSkewIsFresh(t *testing.T) {
	// Regression (found via the live-clock driver smoke run): the revocation
	// as-of is read a hair after the verifier captures now, so it can be
	// marginally ahead. Within the skew tolerance that is fresh, not stale —
	// otherwise every real-clock verification would fail closed.
	w := newWorld(t)
	rec := w.seal(t, nil)
	w.revoc.Set(w.instance, verify.RevocationStatus{
		State: verify.NotObservedRevoked,
		AsOf:  w.clock.Now().Add(5 * time.Second), // ahead, but within skew below
	})
	v := w.verifier(t, time.Minute, 30*time.Second) // skew 30s > 5s
	verdict, _ := v.Verify(rec.Presented())
	if !verdict.IsAccept() {
		t.Fatalf("marginally-ahead as-of within skew must be fresh: verdict=%s causes=%v", verdict.Decision, verdict.Causes)
	}
}

func TestNotRevokedObservationDatedAheadIsInconclusive(t *testing.T) {
	// An as-of ahead of the verifier's clock claims currency the clock
	// cannot corroborate; the verifier must not accept on it (fail closed).
	w := newWorld(t)
	rec := w.seal(t, nil)
	w.revoc.Set(w.instance, verify.RevocationStatus{
		State: verify.NotObservedRevoked,
		AsOf:  w.clock.Now().Add(time.Hour),
	})
	v := w.verifier(t, time.Minute, time.Second)
	verdict, _ := v.Verify(rec.Presented())
	if verdict.Decision != verify.InconclusiveRejected || !hasCause(verdict.Causes, verify.RevocationKnowledgeStale) {
		t.Fatalf("verdict=%s causes=%v, want InconclusiveRejected(RevocationKnowledgeStale)", verdict.Decision, verdict.Causes)
	}
}

// TestDegenerateProviderAlwaysInconclusive mirrors the pre-spike default (M5
// degenerate realization): with no configured knowledge the provider answers
// Indeterminate, and the system fails closed rather than accepting (AP5).
func TestDegenerateProviderFailsClosed(t *testing.T) {
	w := newWorld(t)
	rec := w.seal(t, nil) // revoc left answering Indeterminate for all
	v := w.verifier(t, time.Minute, time.Second)
	verdict, _ := v.Verify(rec.Presented())
	if verdict.Decision != verify.InconclusiveRejected || !hasCause(verdict.Causes, verify.RevocationStatusIndeterminate) {
		t.Fatalf("verdict=%s causes=%v, want InconclusiveRejected(RevocationStatusIndeterminate)", verdict.Decision, verdict.Causes)
	}
}
