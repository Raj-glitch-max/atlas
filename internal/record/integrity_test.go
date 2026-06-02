package record

import (
	"crypto/ecdsa"
	"encoding/base64"
	"encoding/json"
	"strings"
	"testing"

	"github.com/go-jose/go-jose/v3"
	"github.com/spiffe/go-spiffe/v2/spiffeid"
)

// splitCompact returns the three segments of a compact JWS.
func splitCompact(t *testing.T, presented []byte) [3]string {
	t.Helper()
	parts := strings.Split(string(presented), ".")
	if len(parts) != 3 {
		t.Fatalf("compact JWS must have 3 segments, got %d", len(parts))
	}
	return [3]string{parts[0], parts[1], parts[2]}
}

func requireAltered(t *testing.T, name string, presented []byte, tm TrustMaterial) {
	t.Helper()
	rec, outcome := ValidateIntegrity(presented, tm)
	if outcome != Altered || rec != nil {
		t.Errorf("%s: outcome = %v (record %v), want Altered/nil", name, outcome, rec != nil)
	}
}

func TestValidateIntegrityBaseline(t *testing.T) {
	signer, tm := testSigner(t)
	sealed := mustSeal(t, testAssertions(t), signer)
	rec, outcome := ValidateIntegrity(sealed.Presented(), tm)
	if outcome != Intact || rec == nil {
		t.Fatalf("valid record must be Intact, got %v", outcome)
	}
	// Presented bytes are preserved verbatim.
	if string(rec.Presented()) != string(sealed.Presented()) {
		t.Error("Presented() must return the exact presented bytes")
	}
}

func TestValidateIntegrityGarbageInputs(t *testing.T) {
	_, tm := testSigner(t)
	cases := map[string][]byte{
		"empty":              {},
		"nil":                nil,
		"not a jws":          []byte("hello world"),
		"two segments":       []byte("aGVsbG8.d29ybGQ"),
		"four segments":      []byte("aGVsbG8.d29ybGQ.aGVsbG8.d29ybGQ"),
		"empty segment":      []byte("aGVsbG8..d29ybGQ"),
		"invalid characters": []byte("aGVsbG8.d29y bGQ.c2ln"),
		"json serialization": []byte(`{"payload":"aGVsbG8","signatures":[]}`),
	}
	for name, presented := range cases {
		requireAltered(t, name, presented, tm)
	}
}

func TestValidateIntegrityWrongKeyAndUnknownKid(t *testing.T) {
	signer, _ := testSigner(t)
	sealed := mustSeal(t, testAssertions(t), signer)
	domain := spiffeid.RequireTrustDomainFromString("domain-a.test")

	// Same kid, different key: signature must not verify.
	other := testKey(t)
	wrongKey, err := NewTrustMaterial(domain, map[string]*ecdsa.PublicKey{
		signer.KeyID: &other.PublicKey,
	})
	if err != nil {
		t.Fatal(err)
	}
	requireAltered(t, "wrong key under same kid", sealed.Presented(), wrongKey)

	// Material that does not hold the record's kid: indistinguishable
	// from kid tampering; definitive Altered (integrity.go step 4).
	unknownKid, err := NewTrustMaterial(domain, map[string]*ecdsa.PublicKey{
		"different-kid": &other.PublicKey,
	})
	if err != nil {
		t.Fatal(err)
	}
	requireAltered(t, "unknown kid", sealed.Presented(), unknownKid)
}

func TestValidateIntegrityAlgorithmConfusion(t *testing.T) {
	signer, tm := testSigner(t)
	sealed := mustSeal(t, testAssertions(t), signer)
	segs := splitCompact(t, sealed.Presented())

	// alg=none with the signature stripped: the classic downgrade.
	noneHeader := base64.RawURLEncoding.EncodeToString(
		[]byte(`{"alg":"none","kid":"` + signer.KeyID + `","typ":"` + headerType + `"}`))
	requireAltered(t, "alg none", []byte(noneHeader+"."+segs[1]+"."), tm)
	requireAltered(t, "alg none with stale sig", []byte(noneHeader+"."+segs[1]+"."+segs[2]), tm)

	// alg=HS256 with the P-256 public key bytes as the HMAC secret: the
	// classic key-type confusion. The pin rejects on the header before
	// any verification runs.
	pub := &signer.Key.PublicKey
	secret := append(pub.X.Bytes(), pub.Y.Bytes()...)
	hmacSigner, err := jose.NewSigner(jose.SigningKey{
		Algorithm: jose.HS256,
		Key:       jose.JSONWebKey{Key: secret, KeyID: signer.KeyID},
	}, (&jose.SignerOptions{}).WithType(headerType))
	if err != nil {
		t.Fatal(err)
	}
	payload, err := base64.RawURLEncoding.DecodeString(segs[1])
	if err != nil {
		t.Fatal(err)
	}
	confused, err := hmacSigner.Sign(payload)
	if err != nil {
		t.Fatal(err)
	}
	compact, err := confused.CompactSerialize()
	if err != nil {
		t.Fatal(err)
	}
	requireAltered(t, "alg HS256 confusion", []byte(compact), tm)

	// A different ECDSA algorithm name over the same content.
	es384Header := base64.RawURLEncoding.EncodeToString(
		[]byte(`{"alg":"ES384","kid":"` + signer.KeyID + `","typ":"` + headerType + `"}`))
	requireAltered(t, "alg ES384", []byte(es384Header+"."+segs[1]+"."+segs[2]), tm)
}

func TestValidateIntegrityHeaderTampering(t *testing.T) {
	signer, tm := testSigner(t)
	sealed := mustSeal(t, testAssertions(t), signer)
	segs := splitCompact(t, sealed.Presented())

	// Missing typ.
	noTyp := base64.RawURLEncoding.EncodeToString(
		[]byte(`{"alg":"ES256","kid":"` + signer.KeyID + `"}`))
	requireAltered(t, "missing typ", []byte(noTyp+"."+segs[1]+"."+segs[2]), tm)

	// Wrong typ: a token minted for another protocol must never
	// validate as a delegation record.
	wrongTyp := base64.RawURLEncoding.EncodeToString(
		[]byte(`{"alg":"ES256","kid":"` + signer.KeyID + `","typ":"JWT"}`))
	requireAltered(t, "wrong typ", []byte(wrongTyp+"."+segs[1]+"."+segs[2]), tm)

	// Missing kid.
	noKid := base64.RawURLEncoding.EncodeToString(
		[]byte(`{"alg":"ES256","typ":"` + headerType + `"}`))
	requireAltered(t, "missing kid", []byte(noKid+"."+segs[1]+"."+segs[2]), tm)

	// Swapped kid to another held key: signature was made by key-1, so
	// resolving key-2 must fail verification.
	key2 := testKey(t)
	domain := spiffeid.RequireTrustDomainFromString("domain-a.test")
	twoKeys, err := NewTrustMaterial(domain, map[string]*ecdsa.PublicKey{
		signer.KeyID: &signer.Key.PublicKey,
		"key-2":      &key2.PublicKey,
	})
	if err != nil {
		t.Fatal(err)
	}
	swappedKid := base64.RawURLEncoding.EncodeToString(
		[]byte(`{"alg":"ES256","kid":"key-2","typ":"` + headerType + `"}`))
	requireAltered(t, "kid swapped to another held key",
		[]byte(swappedKid+"."+segs[1]+"."+segs[2]), twoKeys)
}

func TestValidateIntegritySignatureTransplant(t *testing.T) {
	signer, tm := testSigner(t)
	a := testAssertions(t)
	first := mustSeal(t, a, signer)

	b := testAssertions(t)
	inst, err := InstanceIDFromString("inst-0002")
	if err != nil {
		t.Fatal(err)
	}
	b.Instance = inst
	b.Scope = []string{"read:orders"}
	second := mustSeal(t, b, signer)

	f := splitCompact(t, first.Presented())
	s := splitCompact(t, second.Presented())

	// First record's header+payload with second record's signature.
	requireAltered(t, "transplanted signature", []byte(f[0]+"."+f[1]+"."+s[2]), tm)
	// Second record's payload under first record's signature.
	requireAltered(t, "transplanted payload", []byte(f[0]+"."+s[1]+"."+f[2]), tm)
}

func TestValidateIntegrityAuthenticButMalformedPayload(t *testing.T) {
	// Payloads signed with the real key but violating the claim contract:
	// authentic bytes the model cannot vouch assertions for → Altered.
	signer, tm := testSigner(t)

	signWith := func(t *testing.T, payload string) []byte {
		t.Helper()
		joseSigner, err := jose.NewSigner(jose.SigningKey{
			Algorithm: jose.ES256,
			Key:       jose.JSONWebKey{Key: signer.Key, KeyID: signer.KeyID},
		}, (&jose.SignerOptions{}).WithType(headerType))
		if err != nil {
			t.Fatal(err)
		}
		jws, err := joseSigner.Sign([]byte(payload))
		if err != nil {
			t.Fatal(err)
		}
		compact, err := jws.CompactSerialize()
		if err != nil {
			t.Fatal(err)
		}
		return []byte(compact)
	}

	cases := map[string]string{
		"not json":           `not-json`,
		"missing sub":        `{"act":{"sub":"spiffe://domain-a.test/d"},"scope":["a"],"exp":1,"iat":1,"atl_ins":"i"}`,
		"non-spiffe sub":     `{"sub":"https://x","act":{"sub":"spiffe://domain-a.test/d"},"scope":["a"],"exp":1,"iat":1,"atl_ins":"i"}`,
		"missing act.sub":    `{"sub":"spiffe://domain-a.test/p","act":{},"scope":["a"],"exp":1,"iat":1,"atl_ins":"i"}`,
		"empty scope":        `{"sub":"spiffe://domain-a.test/p","act":{"sub":"spiffe://domain-a.test/d"},"scope":[],"exp":1,"iat":1,"atl_ins":"i"}`,
		"unsorted scope":     `{"sub":"spiffe://domain-a.test/p","act":{"sub":"spiffe://domain-a.test/d"},"scope":["b","a"],"exp":1,"iat":1,"atl_ins":"i"}`,
		"duplicate scope":    `{"sub":"spiffe://domain-a.test/p","act":{"sub":"spiffe://domain-a.test/d"},"scope":["a","a"],"exp":1,"iat":1,"atl_ins":"i"}`,
		"empty scope entry":  `{"sub":"spiffe://domain-a.test/p","act":{"sub":"spiffe://domain-a.test/d"},"scope":[""],"exp":1,"iat":1,"atl_ins":"i"}`,
		"zero exp":           `{"sub":"spiffe://domain-a.test/p","act":{"sub":"spiffe://domain-a.test/d"},"scope":["a"],"exp":0,"iat":1,"atl_ins":"i"}`,
		"zero iat":           `{"sub":"spiffe://domain-a.test/p","act":{"sub":"spiffe://domain-a.test/d"},"scope":["a"],"exp":1,"iat":0,"atl_ins":"i"}`,
		"missing instance":   `{"sub":"spiffe://domain-a.test/p","act":{"sub":"spiffe://domain-a.test/d"},"scope":["a"],"exp":1,"iat":1}`,
		"bad binding base64": `{"sub":"spiffe://domain-a.test/p","act":{"sub":"spiffe://domain-a.test/d"},"scope":["a"],"exp":1,"iat":1,"atl_ins":"i","atl_rvb":"!!!"}`,
		"exp as json string": `{"sub":"spiffe://domain-a.test/p","act":{"sub":"spiffe://domain-a.test/d"},"scope":["a"],"exp":"1","iat":1,"atl_ins":"i"}`,
	}
	for name, payload := range cases {
		requireAltered(t, name, signWith(t, payload), tm)
	}

	// Forward compatibility: an authentic payload with an UNKNOWN extra
	// field is Intact (fields are append-only across versions).
	extra := `{"sub":"spiffe://domain-a.test/p","act":{"sub":"spiffe://domain-a.test/d"},"scope":["a"],"exp":1800000600,"iat":1800000000,"atl_ins":"i","atl_future":"x"}`
	rec, outcome := ValidateIntegrity(signWith(t, extra), tm)
	if outcome != Intact || rec == nil {
		t.Error("unknown extra field in an authentic payload must remain Intact (append-only evolution)")
	}
}

func TestValidateIntegrityIsDeterministic(t *testing.T) {
	signer, tm := testSigner(t)
	sealed := mustSeal(t, testAssertions(t), signer)
	presented := sealed.Presented()
	for i := 0; i < 50; i++ {
		if _, outcome := ValidateIntegrity(presented, tm); outcome != Intact {
			t.Fatalf("iteration %d: non-deterministic outcome", i)
		}
	}
	mutated := append([]byte(nil), presented...)
	mutated[len(mutated)/2] ^= 0x01
	for i := 0; i < 50; i++ {
		if _, outcome := ValidateIntegrity(mutated, tm); outcome != Altered {
			// A flip can land in a base64 char yielding the same
			// decoded content only at segment-final chars; the
			// midpoint is never segment-final for our sizes, but
			// guard the assumption explicitly.
			seg := strings.Count(string(presented[:len(presented)/2]), ".")
			t.Fatalf("iteration %d: mutated record validated (segment %d)", i, seg)
		}
	}
}

// TestClaimLayoutIsStable pins the wire claim names: changing any is an
// interface-spec amendment (universal rule 5), so a rename fails here first.
func TestClaimLayoutIsStable(t *testing.T) {
	signer, _ := testSigner(t)
	a := testAssertions(t)
	a.RevocationBinding = RevBinding{1}
	sealed := mustSeal(t, a, signer)
	segs := splitCompact(t, sealed.Presented())
	payload, err := base64.RawURLEncoding.DecodeString(segs[1])
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]any
	if err := json.Unmarshal(payload, &m); err != nil {
		t.Fatal(err)
	}
	for _, claim := range []string{"sub", "act", "scope", "exp", "iat", "atl_ins", "atl_rvb"} {
		if _, ok := m[claim]; !ok {
			t.Errorf("claim %q missing from payload", claim)
		}
	}
	if len(m) != 7 {
		t.Errorf("payload has %d claims, want exactly 7 (append-only rule: additions are spec amendments)", len(m))
	}
}

// checkpoint: refactor(security): refactor fuzzing harness execution

// checkpoint: chore(security): harden integration test runner
