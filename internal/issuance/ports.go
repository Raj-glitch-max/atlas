package issuance

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/Raj-glitch-max/atlas/internal/record"
	"github.com/spiffe/go-spiffe/v2/spiffeid"
)

// The Issuance Authority consumes four ports. It defines the contracts;
// realizations (a real permission source, the spike-selected revocation
// binding source) and the harness satisfy them. Time and instance-minting are
// injected so issuance is deterministic and testable.

// PermissionSource answers a principal's own permission set (E4/E5 realization
// deferred, P3). Absence is honest: an unavailable source causes issuance to
// refuse (PermissionsUnavailable), never to guess.
type PermissionSource interface {
	PermissionsOf(principal spiffeid.ID) (PermissionSet, bool)
}

// Clock supplies the issuance-time reading (AD-014): the same single-method
// shape the verifier's time port uses, so one clock serves both. Injected,
// never a direct wall-clock read.
type Clock interface {
	Now() time.Time
}

// RevBindingSource supplies the opaque, mechanism-specific revocation binding
// carried by a record (AD-015). Pre-spike it answers absent (NoBinding); the
// real source is decided by the EXP-001 outcome (E7). The binding is opaque
// to issuance — it is minted by the mechanism side and interpreted only by
// revocation-status realizations.
type RevBindingSource interface {
	RevocationBindingFor(instance record.InstanceID) record.RevBinding
}

// NoBinding is the pre-spike RevBindingSource: it carries no binding. Records
// issued with it have an absent revocation binding, which the degenerate
// revocation provider (and the verifier) handle by failing closed.
type NoBinding struct{}

// RevocationBindingFor always returns absent.
func (NoBinding) RevocationBindingFor(record.InstanceID) record.RevBinding { return nil }

// Minter mints a fresh, unique instance identity per issuance (AD-013 seam).
// Injectable so tests are deterministic; the future FM5-amendment semantics
// land here without changing any interface.
type Minter interface {
	Mint() (record.InstanceID, error)
}

// RandomMinter mints a 128-bit random instance identity — unique per issuance
// with overwhelming probability. This is the V1 default (AD-013): opaque and
// unique, with no interpretable structure.
type RandomMinter struct{}

// Mint returns a fresh random instance identity.
func (RandomMinter) Mint() (record.InstanceID, error) {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		return record.InstanceID{}, err
	}
	return record.InstanceIDFromString(hex.EncodeToString(b[:]))
}
