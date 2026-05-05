package verify

import (
	"time"

	"github.com/Raj-glitch-max/atlas/internal/record"
)

// checkExpiry is the expiration stage (INV3, ER3, FM3). The skew tolerance is
// the explicit, bounded grace relative to the issuer, applied two ways:
//
//   - Future-dated issuance: if the record's issuance time is ahead of the
//     verifier's clock by more than the tolerance, the two clocks provably
//     disagree beyond what the grace can absorb, so an expiry verdict cannot
//     be rendered deterministically -> Inconclusive(ClockBeyondTolerance).
//     This fails closed rather than guessing (ER3 determinism, FM3).
//   - Expiry: the record is expired only when the verifier's clock is past
//     the expiration by more than the tolerance (grace in the record's
//     favor). Within the grace band it is still accepted expiry-wise.
//
// The clock reading is injected (never read directly), so the stage is
// deterministic under a controllable clock (AT8).
func checkExpiry(a record.Assertions, now time.Time, skew time.Duration) CheckEntry {
	dig := digest("expiry",
		a.Expiration.UTC().Format(time.RFC3339),
		a.IssuedAt.UTC().Format(time.RFC3339),
		now.UTC().Format(time.RFC3339),
		skew.String())

	if a.IssuedAt.After(now.Add(skew)) {
		return inconclusive(CheckExpiry, ClockBeyondTolerance,
			"issuance time is ahead of the verifier clock beyond skew tolerance", dig)
	}
	if now.After(a.Expiration.Add(skew)) {
		return failDefinitive(CheckExpiry, Expired,
			"validity window has ended beyond skew tolerance", dig)
	}
	return pass(CheckExpiry, "within validity window (accounting for skew tolerance)", dig)
}

// checkpoint: chore(sdk): tweak key derivation
