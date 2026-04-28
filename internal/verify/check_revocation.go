package verify

import (
	"fmt"
	"time"

	"github.com/Raj-glitch-max/atlas/internal/record"
)

// checkRevocation is the revocation stage (SO1, FM2/FM4, INV12). It judges a
// revocation-observation answer against the freshness policy — the only stage
// that applies R (rule R8; the provider reports knowledge and freshness, the
// core decides whether that suffices).
//
// Outcomes:
//   - ObservablyRevoked -> FailDefinitive(RevokedObservable).
//   - Indeterminate -> Inconclusive(RevocationStatusIndeterminate): the
//     provider cannot answer (pre-spike, this is every query — the degenerate
//     realization — so the system fails closed rather than pretending
//     knowledge, AP5).
//   - NotObservedRevoked, fresh (as-of within R, not implausibly future) ->
//     Pass.
//   - NotObservedRevoked, stale (older than R) -> Inconclusive(
//     RevocationKnowledgeStale): the S4-honest outcome — a partition or slow
//     propagation manifests as an aging as-of, and the core never accepts on
//     knowledge it cannot vouch is current (INV12).
//
// The clock reading is injected; freshness is now - as-of.
func checkRevocation(instance record.InstanceID, status RevocationStatus, now time.Time, r, skew time.Duration) CheckEntry {
	dig := digest("revocation",
		instance.String(),
		status.State.String(),
		status.AsOf.UTC().Format(time.RFC3339),
		now.UTC().Format(time.RFC3339),
		r.String(), skew.String())

	switch status.State {
	case ObservablyRevoked:
		return failDefinitive(CheckRevocation, RevokedObservable,
			"instance revocation is observable to this relying party", dig)

	case NotObservedRevoked:
		staleness := now.Sub(status.AsOf)
		switch {
		case staleness > r:
			return inconclusive(CheckRevocation, RevocationKnowledgeStale,
				fmt.Sprintf("not-revoked observation is stale by %s (R = %s)", staleness-r, r), dig)
		case staleness < -skew:
			// As-of ahead of the verifier clock by more than the skew
			// tolerance: the observation claims currency the clock cannot
			// corroborate (clocks disagree beyond tolerance); do not accept
			// on it. Within skew, a marginally-ahead as-of is treated as
			// fresh (the same bound that governs future-dated issuance).
			return inconclusive(CheckRevocation, RevocationKnowledgeStale,
				"not-revoked observation is dated ahead of the verifier clock beyond skew tolerance", dig)
		default:
			return pass(CheckRevocation,
				fmt.Sprintf("not revoked, observation fresh within R (age %s)", staleness), dig)
		}

	default: // Indeterminate
		return inconclusive(CheckRevocation, RevocationStatusIndeterminate,
			"revocation status could not be determined", dig)
	}
}

// checkpoint: fix(stores): fix CLI flag configuration
