package record

import (
	"encoding/json"

	"github.com/go-jose/go-jose/v3"
	"github.com/spiffe/go-spiffe/v2/spiffeid"
)

// PeekTrustDomainUnverified extracts the principal's trust domain from the
// presented bytes WITHOUT verifying the signature (AD-022). It exists for
// exactly one purpose: letting a relying party that may hold trust material
// for several domains select which material to attempt verification with.
//
// It is NOT authoritative and must never inform a trust decision.
// ValidateIntegrity is the sole authority on authenticity; a record that
// lies about its domain to select different material simply fails
// ValidateIntegrity (the wrong key yields Altered). The unverified peek can
// therefore only cause a verification to be attempted with the wrong
// material and fail — never to wrongly succeed.
//
// Returns false if the bytes are not a well-formed record or carry no
// parseable principal SPIFFE ID.
func PeekTrustDomainUnverified(presented []byte) (spiffeid.TrustDomain, bool) {
	if !isCompactJWS(presented) {
		return spiffeid.TrustDomain{}, false
	}
	jws, err := jose.ParseSigned(string(presented))
	if err != nil {
		return spiffeid.TrustDomain{}, false
	}
	payload := jws.UnsafePayloadWithoutVerification()
	var claims struct {
		Subject string `json:"sub"`
	}
	if err := json.Unmarshal(payload, &claims); err != nil {
		return spiffeid.TrustDomain{}, false
	}
	id, err := spiffeid.FromString(claims.Subject)
	if err != nil {
		return spiffeid.TrustDomain{}, false
	}
	return id.TrustDomain(), true
}
