package record

// The integrity envelope: a JWS-signed token carrying the delegation claims
// (AD-012; gate C2/C3/C7 — token format, scope/expiry claims, and JWS
// tamper-rejection are solved components). This file is the private codec
// between Assertions and the signed payload; nothing here is exported.
//
// Claim layout follows the solved compositions: RFC 8693's act (actor)
// claim expresses "delegate acting on behalf of principal" (sub = principal,
// act.sub = delegate), exp/iat are RFC 7519 NumericDates, scope is a JSON
// array (canonical: sorted, unique, non-empty entries), and the two
// atlas-specific elements ride as atl_ins (instance identity, AD-013) and
// atl_rvb (revocation binding, base64url, AD-015, omitted when absent).

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"github.com/go-jose/go-jose/v3"
	"github.com/spiffe/go-spiffe/v2/spiffeid"
)

const (
	// signatureAlgorithm is the pinned record signature algorithm
	// (AD-012, SR-2). Exactly one algorithm is accepted; a presented
	// record naming any other — including "none" — is Altered. Changing
	// it is an interface-spec amendment.
	signatureAlgorithm = jose.ES256

	// headerType pins the JWS typ header, so a token minted for another
	// protocol never validates as a delegation record (SR-2-adjacent
	// cross-protocol confusion guard).
	headerType = "atlas-record+jws"
)

// payloadClaims is the wire shape of the signed payload.
type payloadClaims struct {
	Subject    string     `json:"sub"`
	Actor      actorClaim `json:"act"`
	Scope      []string   `json:"scope"`
	Expiration int64      `json:"exp"`
	IssuedAt   int64      `json:"iat"`
	Instance   string     `json:"atl_ins"`
	RevBinding string     `json:"atl_rvb,omitempty"`
}

type actorClaim struct {
	Subject string `json:"sub"`
}

// canonicalScope sorts and deduplicates a scope, refusing empty entries.
// Canonical form makes encoding deterministic and inspection unambiguous.
func canonicalScope(scope []string) ([]string, error) {
	if len(scope) == 0 {
		return nil, fmt.Errorf("%w: scope must contain at least one permission", ErrIncompleteAssertions)
	}
	c := append([]string(nil), scope...)
	sort.Strings(c)
	out := c[:0]
	prev := ""
	for i, p := range c {
		if p == "" {
			return nil, fmt.Errorf("%w: scope entries must be non-empty", ErrIncompleteAssertions)
		}
		if i > 0 && p == prev {
			continue
		}
		out = append(out, p)
		prev = p
	}
	return out, nil
}

// isCanonicalScope reports whether a decoded scope is already in canonical
// form. Only Seal produces records, and Seal always writes canonical form,
// so a non-canonical scope in an authentic payload is malformed issuance
// output and the record is Altered (decode refuses).
func isCanonicalScope(scope []string) bool {
	if len(scope) == 0 {
		return false
	}
	for i, p := range scope {
		if p == "" {
			return false
		}
		if i > 0 && scope[i-1] >= p {
			return false
		}
	}
	return true
}

// encodeClaims maps canonicalized assertions to payload bytes. Struct-field
// order makes the JSON deterministic for identical input.
func encodeClaims(a Assertions) ([]byte, error) {
	claims := payloadClaims{
		Subject:    a.Principal.String(),
		Actor:      actorClaim{Subject: a.Delegate.String()},
		Scope:      a.Scope,
		Expiration: a.Expiration.Unix(),
		IssuedAt:   a.IssuedAt.Unix(),
		Instance:   a.Instance.String(),
	}
	if !a.RevocationBinding.IsAbsent() {
		claims.RevBinding = base64.RawURLEncoding.EncodeToString(a.RevocationBinding)
	}
	return json.Marshal(claims)
}

// decodeClaims maps an authenticated payload back to assertions. It runs
// only on signature-verified bytes; any required element that is missing or
// malformed means the issuer emitted something this model cannot vouch
// assertions for, and the caller reports Altered (the contract has no
// failure distinct from it). Unknown fields are tolerated: fields are
// append-only across versions (interface-spec universal rule 5), so an
// older reader must not reject a newer issuer's authentic record.
func decodeClaims(payload []byte) (Assertions, error) {
	var claims payloadClaims
	if err := json.Unmarshal(payload, &claims); err != nil {
		return Assertions{}, fmt.Errorf("record: payload not valid JSON: %w", err)
	}
	principal, err := spiffeid.FromString(claims.Subject)
	if err != nil {
		return Assertions{}, fmt.Errorf("record: sub is not a SPIFFE ID: %w", err)
	}
	delegate, err := spiffeid.FromString(claims.Actor.Subject)
	if err != nil {
		return Assertions{}, fmt.Errorf("record: act.sub is not a SPIFFE ID: %w", err)
	}
	if !isCanonicalScope(claims.Scope) {
		return Assertions{}, fmt.Errorf("record: scope is not in canonical form")
	}
	if claims.Expiration <= 0 || claims.IssuedAt <= 0 {
		return Assertions{}, fmt.Errorf("record: exp and iat must be positive NumericDates")
	}
	instance, err := InstanceIDFromString(claims.Instance)
	if err != nil {
		return Assertions{}, fmt.Errorf("record: atl_ins missing: %w", err)
	}
	var binding RevBinding
	if claims.RevBinding != "" {
		raw, err := base64.RawURLEncoding.DecodeString(claims.RevBinding)
		if err != nil {
			return Assertions{}, fmt.Errorf("record: atl_rvb is not base64url: %w", err)
		}
		binding = RevBinding(raw)
	}
	return Assertions{
		Principal:         principal,
		Delegate:          delegate,
		Scope:             claims.Scope,
		Expiration:        time.Unix(claims.Expiration, 0).UTC(),
		IssuedAt:          time.Unix(claims.IssuedAt, 0).UTC(),
		Instance:          instance,
		RevocationBinding: binding,
	}, nil
}

// isCompactJWS enforces the compact serialization shape — exactly three
// base64url segments separated by dots — before any parser runs. The JSON
// serializations jose would also accept are rejected up front: the
// presentable unit is one fixed shape, which shrinks the parse surface and
// keeps "the bytes presented" and "the bytes signed over" in one-to-one
// correspondence.
func isCompactJWS(presented []byte) bool {
	dots := 0
	segLen := 0
	for _, b := range presented {
		if b == '.' {
			if segLen == 0 {
				return false
			}
			dots++
			segLen = 0
			continue
		}
		if !isBase64URLByte(b) {
			return false
		}
		segLen++
	}
	return dots == 2 && segLen > 0
}

func isBase64URLByte(b byte) bool {
	switch {
	case b >= 'A' && b <= 'Z', b >= 'a' && b <= 'z', b >= '0' && b <= '9', b == '-', b == '_':
		return true
	}
	return false
}
