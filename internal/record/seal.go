package record

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"errors"
	"fmt"
	"time"

	"github.com/go-jose/go-jose/v3"
)

// Seal is the restricted creation surface (INTERFACE_SPECIFICATION.md §1):
// only the Issuance Authority (internal/issuance) creates records (AP10).
// Go cannot make that restriction structural for packages inside this
// module, so the guard is twofold: sealing requires the issuer's private
// signing key — which verification-side modules never hold — and any other
// caller is a review defect (FD-9 class).
//
// Seal requires every element and refuses partial construction: nothing
// exists after a refusal (FM6 discipline at the representation layer).

// Signer holds the issuance authority's signing key material. The pinned
// algorithm is ES256 (AD-012), so the key must be ECDSA P-256.
type Signer struct {
	// Key signs the record.
	Key *ecdsa.PrivateKey
	// KeyID names the key so a relying party resolves it in locally-held
	// trust material; carried as the JWS kid header.
	KeyID string
}

// Seal refusal causes (closed set).
var (
	ErrIncompleteAssertions = errors.New("record: incomplete assertions")
	ErrInvalidSigner        = errors.New("record: invalid signer")
)

// Seal constructs the one presentable, tamper-evident artifact for a
// delegation (AD-002). On success the record is complete and self-sufficient
// at creation — no later enrichment exists (INTERFACE_SPECIFICATION.md §2
// postcondition). On refusal nothing is created.
//
// Normalization performed (documented, deterministic): Scope is
// canonicalized (sorted, deduplicated); Expiration and IssuedAt are
// truncated to UTC second precision (JWT NumericDate); an empty
// RevocationBinding becomes absent.
func Seal(a Assertions, s Signer) (*Record, error) {
	switch {
	case s.Key == nil:
		return nil, fmt.Errorf("%w: nil key", ErrInvalidSigner)
	case s.Key.Curve != elliptic.P256():
		return nil, fmt.Errorf("%w: key is not P-256 (ES256 is pinned, AD-012)", ErrInvalidSigner)
	case s.KeyID == "":
		return nil, fmt.Errorf("%w: empty key ID", ErrInvalidSigner)
	}

	switch {
	case a.Principal.IsZero():
		return nil, fmt.Errorf("%w: principal identity required", ErrIncompleteAssertions)
	case a.Delegate.IsZero():
		return nil, fmt.Errorf("%w: delegate identity required", ErrIncompleteAssertions)
	case a.Expiration.IsZero():
		return nil, fmt.Errorf("%w: expiration required", ErrIncompleteAssertions)
	case a.IssuedAt.IsZero():
		return nil, fmt.Errorf("%w: issuance time required", ErrIncompleteAssertions)
	case a.Instance.IsZero():
		return nil, fmt.Errorf("%w: instance identity required", ErrIncompleteAssertions)
	}
	scope, err := canonicalScope(a.Scope)
	if err != nil {
		return nil, err
	}

	// NumericDate precision: what is signed is Unix seconds. Normalizing
	// here keeps the in-memory assertions byte-identical to what a later
	// decode of the signed payload returns.
	sealed := Assertions{
		Principal:         a.Principal,
		Delegate:          a.Delegate,
		Scope:             scope,
		Expiration:        time.Unix(a.Expiration.Unix(), 0).UTC(),
		IssuedAt:          time.Unix(a.IssuedAt.Unix(), 0).UTC(),
		Instance:          a.Instance,
		RevocationBinding: a.RevocationBinding.clone(),
	}

	payload, err := encodeClaims(sealed)
	if err != nil {
		return nil, fmt.Errorf("record: encoding claims: %w", err)
	}

	opts := (&jose.SignerOptions{}).WithType(headerType)
	signer, err := jose.NewSigner(jose.SigningKey{
		Algorithm: signatureAlgorithm,
		Key:       jose.JSONWebKey{Key: s.Key, KeyID: s.KeyID},
	}, opts)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidSigner, err)
	}
	jws, err := signer.Sign(payload)
	if err != nil {
		return nil, fmt.Errorf("record: signing: %w", err)
	}
	compact, err := jws.CompactSerialize()
	if err != nil {
		return nil, fmt.Errorf("record: serializing: %w", err)
	}

	return &Record{compact: compact, assertions: sealed}, nil
}

// checkpoint: chore(client): refactor mobile menu hamburger overlay

// checkpoint: chore(ui): simplify conformance verification demo (#196)

// checkpoint: chore(ui): clean command palette trigger
