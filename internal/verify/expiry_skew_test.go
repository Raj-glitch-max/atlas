package verify_test

import (
	"testing"
	"time"

	"github.com/Raj-glitch-max/atlas/internal/record"
	"github.com/Raj-glitch-max/atlas/internal/verify"
)

// AT8 logic — expiry within an explicit, bounded clock-skew tolerance (ER3,
// FM3). The clock is injected, so these are deterministic unit assertions,
// not wall-time-dependent. Two boundaries are exercised: the expiry grace
// band (a record is expired only past expiration + skew) and the
// future-issuance tolerance (issuance ahead of now beyond skew is
// inconclusive, not a guess).

func TestExpiryGraceBandBoundaries(t *testing.T) {
	const skew = 30 * time.Second
	exp := baseExpiry

	cases := []struct {
		name       string
		now        time.Time
		wantAccept bool
	}{
		{"well within window", exp.Add(-time.Hour), true},
		{"just before expiry", exp.Add(-time.Second), true},
		{"exactly at expiry", exp, true},
		{"within grace band", exp.Add(skew - time.Second), true},
		{"exactly at grace edge", exp.Add(skew), true},
		{"just past grace edge", exp.Add(skew + time.Second), false},
		{"long past expiry", exp.Add(time.Hour), false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			w := newWorld(t)
			w.clock.Set(tc.now)
			w.revoc.Set(w.instance, verify.RevocationStatus{State: verify.NotObservedRevoked, AsOf: tc.now})
			rec := w.seal(t, nil)
			v := w.verifier(t, time.Minute, skew)

			verdict, trace := v.Verify(rec.Presented())
			entry, _ := trace.Find(verify.CheckExpiry)
			if tc.wantAccept {
				if !verdict.IsAccept() {
					t.Fatalf("now=%v: verdict=%s causes=%v, want Accept", tc.now, verdict.Decision, verdict.Causes)
				}
				if entry.Outcome != verify.OutcomePass {
					t.Errorf("expiry stage = %s, want Pass", entry.Outcome)
				}
			} else {
				if verdict.Decision != verify.Reject || !hasCause(verdict.Causes, verify.Expired) {
					t.Fatalf("now=%v: verdict=%s causes=%v, want Reject(Expired)", tc.now, verdict.Decision, verdict.Causes)
				}
			}
		})
	}
}

func TestFutureIssuanceToleranceBoundaries(t *testing.T) {
	const skew = 30 * time.Second
	now := baseNow

	cases := []struct {
		name             string
		issuedAt         time.Time
		wantInconclusive bool
	}{
		{"issued in the past", now.Add(-time.Hour), false},
		{"issued now", now, false},
		{"issued slightly ahead within skew", now.Add(skew - time.Second), false},
		{"issued exactly at skew edge", now.Add(skew), false},
		{"issued just beyond skew", now.Add(skew + time.Second), true},
		{"issued far in the future", now.Add(time.Hour), true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			w := newWorld(t)
			w.clock.Set(now)
			w.revoc.Set(w.instance, verify.RevocationStatus{State: verify.NotObservedRevoked, AsOf: now})
			rec := w.seal(t, func(a *record.Assertions) {
				a.IssuedAt = tc.issuedAt
				a.Expiration = tc.issuedAt.Add(time.Hour)
			})
			v := w.verifier(t, time.Minute, skew)

			verdict, _ := v.Verify(rec.Presented())
			if tc.wantInconclusive {
				if verdict.Decision != verify.InconclusiveRejected || !hasCause(verdict.Causes, verify.ClockBeyondTolerance) {
					t.Fatalf("iat=%v: verdict=%s causes=%v, want InconclusiveRejected(ClockBeyondTolerance)",
						tc.issuedAt, verdict.Decision, verdict.Causes)
				}
			} else if !verdict.IsAccept() {
				t.Fatalf("iat=%v: verdict=%s causes=%v, want Accept", tc.issuedAt, verdict.Decision, verdict.Causes)
			}
		})
	}
}
