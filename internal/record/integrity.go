package record

import "github.com/go-jose/go-jose/v3"

// ValidateIntegrity determines whether presented bytes are exactly what the
// issuer created, verifiable with the supplied trust material
// (INTERFACE_SPECIFICATION.md §1; INV8). The answer set is closed: Intact or
// Altered — the operation is never inconclusive, because it holds everything
// it needs by construction; verification-key *absence* is the caller's prior
// condition (the verifier's check 2 handles absent material before this
// operation runs).
//
// On Intact the validated *Record is returned so Read is defined exactly on
// validated records; on Altered the record is nil.
//
// Any byte sequence is a legal input. The validation pipeline, every exit of
// which is Altered:
//
//  1. shape — compact JWS form only (three base64url segments);
//  2. parse — well-formed JWS with exactly one signature;
//  3. pin — header algorithm is exactly ES256 (AD-012/SR-2: an
//     attacker-selected algorithm, including "none", never reaches
//     verification), header typ is exactly the record type, kid present;
//  4. resolve — kid is held in the supplied trust material (an unknown kid
//     is indistinguishable from kid tampering, and the frozen package
//     warrants no key-rotation state that could make it legitimate — FM5
//     non-objective — so the definitive answer is the honest one);
//  5. verify — the ES256 signature over the payload with the resolved key;
//  6. decode — the authenticated payload yields complete, well-formed
//     assertions in canonical form.
//
// Deterministic: same presented bytes and material, same answer. Expired or
// revoked records still validate Intact — validity is the verifier's
// question, never this one.
func ValidateIntegrity(presented []byte, tm TrustMaterial) (*Record, Outcome) {
	if !isCompactJWS(presented) {
		return nil, Altered
	}
	jws, err := jose.ParseSigned(string(presented))
	if err != nil {
		return nil, Altered
	}
	if len(jws.Signatures) != 1 {
		return nil, Altered
	}
	header := jws.Signatures[0].Header
	if header.Algorithm != string(signatureAlgorithm) {
		return nil, Altered
	}
	if typ, _ := header.ExtraHeaders[jose.HeaderType].(string); typ != headerType {
		return nil, Altered
	}
	if header.KeyID == "" {
		return nil, Altered
	}
	key, held := tm.keyFor(header.KeyID)
	if !held {
		return nil, Altered
	}
	payload, err := jws.Verify(key)
	if err != nil {
		return nil, Altered
	}
	assertions, err := decodeClaims(payload)
	if err != nil {
		return nil, Altered
	}
	return &Record{compact: string(presented), assertions: assertions}, Intact
}
