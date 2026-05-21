package conformance

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/Raj-glitch-max/atlas/internal/verify"
	jose "github.com/go-jose/go-jose/v3"
)

// Negative (adversarial) conformance vectors: malformed and attacking records
// that a conformant verifier MUST reject. A positive/verdict-space corpus is
// half a conformance suite — Wycheproof and Frankencerts are dominated by
// negative cases, because that is where independent implementations diverge
// (the lenient verifier accepts what the strict one rejects). These serialize
// the JWS attack families the Go integrity tests already cover, so a foreign
// verifier's integrity check is tested against identical attacks.
//
// Every negative vector holds the SAME valid trust key, a fresh not-revoked
// answer, and an in-window clock — so the record's malformation is the sole
// reason for rejection. EmitNegativeVectors runs the reference verifier to
// record the authoritative verdict and FAILS generation if any adversarial
// record is accepted (a silent-acceptance bug would surface here).

func b64url(s string) string { return base64.RawURLEncoding.EncodeToString([]byte(s)) }

func splitCompact(b []byte) ([3]string, bool) {
	parts := strings.Split(string(b), ".")
	if len(parts) != 3 {
		return [3]string{}, false
	}
	return [3]string{parts[0], parts[1], parts[2]}, true
}

func joinCompact(h, p, s string) []byte { return []byte(h + "." + p + "." + s) }

// signES256 signs a raw payload with the kit's real key (authentic signature,
// used to build authentic-but-malformed-payload attacks).
func (k *Kit) signES256(payload []byte) ([]byte, error) {
	signer, err := jose.NewSigner(
		jose.SigningKey{Algorithm: jose.ES256, Key: jose.JSONWebKey{Key: k.signer.Key, KeyID: k.signer.KeyID}},
		(&jose.SignerOptions{}).WithType("atlas-record+jws"))
	if err != nil {
		return nil, err
	}
	jws, err := signer.Sign(payload)
	if err != nil {
		return nil, err
	}
	c, err := jws.CompactSerialize()
	if err != nil {
		return nil, err
	}
	return []byte(c), nil
}

// hs256Confusion signs the payload with HS256 using the public key's raw
// coordinates as the HMAC secret — the classic algorithm-confusion attack.
func (k *Kit) hs256Confusion(payload []byte) ([]byte, error) {
	pub := &k.signer.Key.PublicKey
	secret := append(pub.X.Bytes(), pub.Y.Bytes()...)
	signer, err := jose.NewSigner(
		jose.SigningKey{Algorithm: jose.HS256, Key: jose.JSONWebKey{Key: secret, KeyID: k.signer.KeyID}},
		(&jose.SignerOptions{}).WithType("atlas-record+jws"))
	if err != nil {
		return nil, err
	}
	jws, err := signer.Sign(payload)
	if err != nil {
		return nil, err
	}
	c, err := jws.CompactSerialize()
	if err != nil {
		return nil, err
	}
	return []byte(c), nil
}

// EmitNegativeVectors builds the adversarial vector set.
func EmitNegativeVectors() (VectorFile, error) {
	k, err := NewKit()
	if err != nil {
		return VectorFile{}, err
	}
	jwk, err := k.publicJWK()
	if err != nil {
		return VectorFile{}, err
	}
	kid := k.signer.KeyID

	base := k.Valid()
	segs, ok := splitCompact(base)
	if !ok {
		return VectorFile{}, fmt.Errorf("base record is not compact JWS")
	}
	payloadBytes, err := base64.RawURLEncoding.DecodeString(segs[1])
	if err != nil {
		return VectorFile{}, err
	}
	sibling := k.Seal([]string{"read:orders"}, tIssued, tExpiry, "inst-sibling", false)
	ssegs, _ := splitCompact(sibling)

	type mutation struct {
		name string
		desc string
		rec  []byte
	}
	muts := []mutation{
		{"alg-none-stripped-sig", "header alg=none with the signature removed (downgrade)",
			joinCompact(b64url(`{"alg":"none","kid":"`+kid+`","typ":"atlas-record+jws"}`), segs[1], "")},
		{"alg-none-stale-sig", "header alg=none keeping a stale signature",
			joinCompact(b64url(`{"alg":"none","kid":"`+kid+`","typ":"atlas-record+jws"}`), segs[1], segs[2])},
		{"alg-es384", "header claims ES384 over an ES256 signature",
			joinCompact(b64url(`{"alg":"ES384","kid":"`+kid+`","typ":"atlas-record+jws"}`), segs[1], segs[2])},
		{"typ-missing", "header omits the pinned typ",
			joinCompact(b64url(`{"alg":"ES256","kid":"`+kid+`"}`), segs[1], segs[2])},
		{"typ-wrong", "header typ is a generic JWT, not atlas-record+jws",
			joinCompact(b64url(`{"alg":"ES256","kid":"`+kid+`","typ":"JWT"}`), segs[1], segs[2])},
		{"kid-missing", "header omits kid",
			joinCompact(b64url(`{"alg":"ES256","typ":"atlas-record+jws"}`), segs[1], segs[2])},
		{"kid-forged", "header names a key not held",
			joinCompact(b64url(`{"alg":"ES256","kid":"forged-kid","typ":"atlas-record+jws"}`), segs[1], segs[2])},
		{"truncated-signature-removed", "signature segment emptied",
			joinCompact(segs[0], segs[1], "")},
		{"truncated-last-byte", "last byte removed", base[:len(base)-1]},
		{"transplanted-signature", "sibling record's signature over this header+payload",
			joinCompact(segs[0], segs[1], ssegs[2])},
		{"transplanted-payload", "sibling record's payload under this signature",
			joinCompact(segs[0], ssegs[1], segs[2])},
		{"garbage-not-jws", "not a JWS at all", []byte("this is not a delegation record")},
		{"empty", "empty input", []byte("")},
	}

	// Authentic-but-malformed payloads: real ES256 signature, payload violates
	// the claim contract (integrity must still reject — signature authenticity
	// is necessary, not sufficient).
	authentic := []struct{ name, desc, payload string }{
		{"authentic-missing-sub", "authentic signature, payload omits the principal (sub)",
			`{"act":{"sub":"spiffe://domain-a.test/delegate"},"scope":["read:orders"],"exp":1800003600,"iat":1800000000,"atl_ins":"inst-x"}`},
		{"authentic-missing-instance", "authentic signature, payload omits the instance identity",
			`{"sub":"spiffe://domain-a.test/principal","act":{"sub":"spiffe://domain-a.test/delegate"},"scope":["read:orders"],"exp":1800003600,"iat":1800000000}`},
		{"authentic-unsorted-scope", "authentic signature, scope not in canonical (sorted) form",
			`{"sub":"spiffe://domain-a.test/principal","act":{"sub":"spiffe://domain-a.test/delegate"},"scope":["write:audit","read:orders"],"exp":1800003600,"iat":1800000000,"atl_ins":"inst-x"}`},
		{"authentic-nonspiffe-sub", "authentic signature, principal is not a SPIFFE ID",
			`{"sub":"https://example.test/p","act":{"sub":"spiffe://domain-a.test/delegate"},"scope":["read:orders"],"exp":1800003600,"iat":1800000000,"atl_ins":"inst-x"}`},
	}
	for _, a := range authentic {
		rec, err := k.signES256([]byte(a.payload))
		if err != nil {
			return VectorFile{}, err
		}
		muts = append(muts, mutation{a.name, a.desc, rec})
	}
	hs, err := k.hs256Confusion(payloadBytes)
	if err != nil {
		return VectorFile{}, err
	}
	muts = append(muts, mutation{"alg-hs256-confusion", "signed HS256 using the public key coordinates as the HMAC secret", hs})

	// Common inputs: valid trust key, fresh not-revoked, in-window clock, so
	// the malformation is the sole rejection cause.
	asOf := tMid.UTC().Format(time.RFC3339)
	trust := VectorTrust{Domain: k.domain.String(), Keys: []json.RawMessage{jwk}}
	rev := VectorRevocation{State: verify.NotObservedRevoked.String(), AsOf: &asOf}
	pol := VectorPolicy{RSeconds: mustPolicy().R().Seconds(), SkewSeconds: mustPolicy().SkewTolerance().Seconds()}

	freshRev := RevocationInState(verify.NotObservedRevoked, tMid)

	var vectors []Vector
	for _, m := range muts {
		v, err := verify.NewVerifier(mustPolicy(), k.trust, freshRev, fixedClock{tMid})
		if err != nil {
			return VectorFile{}, err
		}
		verdict, _ := v.Verify(m.rec)
		if verdict.IsAccept() {
			return VectorFile{}, fmt.Errorf("negative vector %q was ACCEPTED — silent-acceptance bug", m.name)
		}
		causes := make([]string, 0, len(verdict.Causes))
		for _, c := range verdict.Causes {
			causes = append(causes, c.String())
		}
		vectors = append(vectors, Vector{
			Name:        m.name,
			Description: m.desc,
			Record:      string(m.rec),
			Trust:       trust,
			Revocation:  rev,
			Now:         tMid.UTC().Format(time.RFC3339),
			Policy:      pol,
			Expect:      VectorExpect{Decision: decisionToken(verdict.Decision), Causes: causes},
		})
	}
	return VectorFile{
		Schema:      SchemaVersion,
		Description: "Atlas Verification Core (M3) NEGATIVE conformance vectors: adversarial/malformed records a conformant verifier must reject. Each holds a valid trust key, fresh not-revoked status, and an in-window clock, so the record's malformation is the sole rejection cause.",
		Vectors:     vectors,
	}, nil
}
