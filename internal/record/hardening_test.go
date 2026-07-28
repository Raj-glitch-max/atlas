package record

import (
	"encoding/base64"
	"encoding/json"
	"strings"
	"testing"
)

func TestKeyIDIsBoundedAndInert(t *testing.T) {
	for _, tc := range []struct {
		name string
		kid  string
		ok   bool
	}{
		{"normal", "authority-key-1", true},
		{"dotted", "key.v2", true},
		{"empty", "", false},
		{"traversal", "../../etc/passwd", false},
		{"absolute path", "/etc/passwd", false},
		{"windows path", `..\..\secret`, false},
		{"url", "https://evil.example/jwks.json", false},
		{"dot", ".", false},
		{"dotdot", "..", false},
		{"nul byte", "key\x00", false},
		{"newline", "key\nInjected: header", false},
		{"space", "key 1", false},
		{"unicode", "kéy", false},
		{"too long", strings.Repeat("k", maxKeyIDLen+1), false},
		{"at limit", strings.Repeat("k", maxKeyIDLen), true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := validKeyID(tc.kid); got != tc.ok {
				t.Fatalf("validKeyID(%q) = %v, want %v", tc.kid, got, tc.ok)
			}
		})
	}
}

func TestTimestampsMustBeJSONIntegers(t *testing.T) {
	for _, tc := range []struct {
		name    string
		exp     string
		wantErr bool
	}{
		{"integer", "1800000000", false},
		{"exponent", "1.8e9", true},
		{"float", "1800000000.0", true},
		{"tiny fraction", "1800000000.0000001", true},
		{"plus sign", "+1800000000", true},
		{"hex-ish", "0x6B49D200", true},
		{"string", `"1800000000"`, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			payload := []byte(`{"exp":` + tc.exp + `,"iat":1799999000}`)
			err := requireIntegerNumerics(payload, "exp", "iat")
			if tc.wantErr && err == nil {
				t.Fatalf("exp=%s must be rejected: a JavaScript verifier would accept it and could round, "+
					"so this is a cross-language expiry divergence over identical signed bytes", tc.exp)
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("exp=%s must be accepted, got %v", tc.exp, err)
			}
		})
	}
}

func TestNumericDateUpperBound(t *testing.T) {
	signer, _ := testSigner(t)
	a := testAssertions(t)
	sealed := mustSeal(t, a, signer)

	// Decode the authentic payload, push exp past the representable bound, and
	// re-sign. The signature is valid; only the parsing rule can catch this.
	parts := strings.Split(string(sealed.Presented()), ".")
	raw, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]json.RawMessage
	if err := json.Unmarshal(raw, &m); err != nil {
		t.Fatal(err)
	}
	m["exp"] = json.RawMessage("9223372036854775807") // MaxInt64
	mutated, err := json.Marshal(m)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := decodeClaims(mutated); err == nil {
		t.Fatal("exp = MaxInt64 must be refused: time.Unix overflows into a nonsense date, " +
			"which would make expiry — the one limit enforced with zero coordination — disableable by a large integer")
	}
}

// TestSignatureNonceIsNeverReused is the detector for the one catastrophic
// misuse mode of the pinned algorithm.
//
// ES256 is ECDSA, which needs a per-signature nonce k. If k EVER repeats across
// two different messages, the private key is recoverable in closed form from the
// two (r,s) pairs — not weakened, recovered. This has happened in production
// repeatedly (PS3, Android wallets) and is the strongest argument against ES256
// and for Ed25519, which derives k deterministically.
//
// The nonce is not observable, but its reuse IS: r = (k*G).x, so identical r
// across different messages means identical k. This test signs many distinct
// records and asserts every r is unique.
//
// It cannot prove an implementation is safe — a malicious signer can pass it and
// still leak keys — but it catches the accident, which is how this has always
// actually happened: a weak RNG or a naive library with a fixed or low-entropy k.
// Any reimplementation of Atlas should port this test.
func TestSignatureNonceIsNeverReused(t *testing.T) {
	signer, _ := testSigner(t)
	const n = 256

	seen := make(map[string]int, n)
	for i := 0; i < n; i++ {
		a := testAssertions(t)
		// Vary the payload so each signature covers different bytes. Reused k
		// over IDENTICAL messages is undetectable and harmless; reuse across
		// DIFFERENT messages is the key-recovery case.
		a.Scope = []string{"read:orders", "probe:" + strings.Repeat("x", i%17) + itoaTest(i)}
		canon, err := canonicalScope(a.Scope)
		if err != nil {
			t.Fatal(err)
		}
		a.Scope = canon

		sealed := mustSeal(t, a, signer)
		parts := strings.Split(string(sealed.Presented()), ".")
		sig, err := base64.RawURLEncoding.DecodeString(parts[2])
		if err != nil {
			t.Fatal(err)
		}
		if len(sig) != 64 {
			t.Fatalf("ES256 signature must be 64 bytes (r||s), got %d", len(sig))
		}
		r := string(sig[:32]) // r = (k*G).x — identical r ⇒ identical k

		if prev, dup := seen[r]; dup {
			t.Fatalf("ECDSA NONCE REUSE: signatures %d and %d share the same r, "+
				"so they share k, so the private key is recoverable from the two (r,s) pairs. "+
				"The signer must use RFC 6979 deterministic nonces or a CSPRNG.", prev, i)
		}
		seen[r] = i
	}
	if len(seen) != n {
		t.Fatalf("expected %d distinct nonces, saw %d", n, len(seen))
	}
}

func itoaTest(i int) string {
	if i == 0 {
		return "0"
	}
	var b [8]byte
	p := len(b)
	for i > 0 {
		p--
		b[p] = byte('0' + i%10)
		i /= 10
	}
	return string(b[p:])
}
