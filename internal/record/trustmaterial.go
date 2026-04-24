package record

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"errors"
	"fmt"

	"github.com/spiffe/go-spiffe/v2/spiffeid"
)

// TrustMaterial is the verification-key material a relying party holds
// locally for one trust domain: the vocabulary type consumed by
// ValidateIntegrity and held by the Trust Material Store (M4). It is
// provisioned out-of-band (gate C1 pattern); nothing in this package or its
// consumers fetches it (FM9, structural).
//
// Keys are identified by key ID and are P-256 ECDSA public keys, matching
// the pinned record signature algorithm (ES256, AD-012/SR-2). Enforcing the
// curve at construction keeps algorithm confusion unrepresentable in held
// material, not merely detected at verification.
type TrustMaterial struct {
	domain spiffeid.TrustDomain
	keys   map[string]*ecdsa.PublicKey
}

// Construction refusal causes (closed set; INTERFACE_SPECIFICATION.md
// universal rule 1 — malformed material is refused, never half-trusted).
var (
	ErrTrustMaterialNoDomain = errors.New("record: trust material requires a trust domain")
	ErrTrustMaterialNoKeys   = errors.New("record: trust material requires at least one key")
	ErrTrustMaterialBadKey   = errors.New("record: trust material key invalid")
)

// NewTrustMaterial builds trust material for one trust domain from a keyID →
// public-key map. The map is copied; the caller's map is not retained.
func NewTrustMaterial(domain spiffeid.TrustDomain, keys map[string]*ecdsa.PublicKey) (TrustMaterial, error) {
	if domain.IsZero() {
		return TrustMaterial{}, ErrTrustMaterialNoDomain
	}
	if len(keys) == 0 {
		return TrustMaterial{}, ErrTrustMaterialNoKeys
	}
	held := make(map[string]*ecdsa.PublicKey, len(keys))
	for kid, key := range keys {
		if kid == "" {
			return TrustMaterial{}, fmt.Errorf("%w: empty key ID", ErrTrustMaterialBadKey)
		}
		if key == nil {
			return TrustMaterial{}, fmt.Errorf("%w: nil key for %q", ErrTrustMaterialBadKey, kid)
		}
		if key.Curve != elliptic.P256() {
			return TrustMaterial{}, fmt.Errorf("%w: key %q is not P-256 (ES256 is pinned, AD-012)", ErrTrustMaterialBadKey, kid)
		}
		held[kid] = key
	}
	return TrustMaterial{domain: domain, keys: held}, nil
}

// Domain returns the trust domain this material verifies for.
func (tm TrustMaterial) Domain() spiffeid.TrustDomain { return tm.domain }

// KeyCount returns the number of verification keys held. Used by the trust
// store's provisioning record (M4 observable); never exposes the keys.
func (tm TrustMaterial) KeyCount() int { return len(tm.keys) }

// keyFor resolves a key ID against the held material.
func (tm TrustMaterial) keyFor(kid string) (*ecdsa.PublicKey, bool) {
	key, ok := tm.keys[kid]
	return key, ok
}

// checkpoint: fix(stores): fix test assertions
