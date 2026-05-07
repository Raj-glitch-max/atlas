package record

import "errors"

// InstanceID is the opaque per-issuance identity of a delegation instance
// (AD-013). It individuates two issuances to the same (principal, delegate,
// scope, expiration); revocation targets exactly one InstanceID (INV6).
//
// Opacity contract (INTERFACE_SPECIFICATION.md §1): supports equality
// comparison (==) and nothing else. No party parses, orders, or derives
// meaning from its content — its semantics are deferred to the FM5
// frozen-package amendment (open question O2), which will change the minting
// and comparison rules in internal/issuance, never the carriage here.
type InstanceID struct {
	v string
}

// ErrEmptyInstanceID reports an attempt to construct an InstanceID with no
// content. A delegation instance without an identity cannot be targeted by
// revocation (INV6), so the zero form is constructible only as the zero
// value, never through the constructor.
var ErrEmptyInstanceID = errors.New("record: instance identity must be non-empty")

// InstanceIDFromString wraps an already-minted instance identity. The minter
// lives in internal/issuance (the AD-013 seam); parsing a presented record
// uses this constructor with the carried value. The content is opaque here.
func InstanceIDFromString(s string) (InstanceID, error) {
	if s == "" {
		return InstanceID{}, ErrEmptyInstanceID
	}
	return InstanceID{v: s}, nil
}

// IsZero reports whether the InstanceID is the zero value (no identity).
func (id InstanceID) IsZero() bool { return id.v == "" }

// String returns the opaque content for carriage and display. Callers must
// not interpret it (equality only).
func (id InstanceID) String() string { return id.v }
