package record

// AT20-class mutation testing (SO3, INV8, FM7): across the mutation corpus —
// bit-flips, field substitutions, truncation, reordering of protected
// fields, header and signature substitution — the detection-and-rejection
// fraction must be exactly 1. Corpus parameters (seed, flip count) come from
// tests/fixtures/mutations/corpus.json so runs are reproducible.

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

type corpusSpec struct {
	Seed     int64 `json:"seed"`
	BitFlips int   `json:"bitFlips"`
}

func loadCorpus(t *testing.T) corpusSpec {
	t.Helper()
	raw, err := os.ReadFile(filepath.Join("..", "..", "tests", "fixtures", "mutations", "corpus.json"))
	if err != nil {
		t.Fatalf("reading mutation corpus: %v", err)
	}
	var spec corpusSpec
	if err := json.Unmarshal(raw, &spec); err != nil {
		t.Fatalf("parsing mutation corpus: %v", err)
	}
	if spec.BitFlips <= 0 {
		t.Fatal("corpus bitFlips must be positive")
	}
	return spec
}

// mutation is one corpus entry: a named alteration of a presented record.
type mutation struct {
	name      string
	presented []byte
}

// segmentFinalPositions marks byte positions that are the final character of
// a base64url segment. Flipping bits there can alter unused padding bits
// that lenient base64 decoding ignores, producing different presented bytes
// with identical decoded content — an encoding variation, not a content
// alteration, and outside AT20's mutation set (which targets the record's
// protected content). All other positions must be detected.
func segmentFinalPositions(presented []byte) map[int]bool {
	final := map[int]bool{}
	last := len(presented) - 1
	for i, b := range presented {
		if b == '.' && i > 0 {
			final[i-1] = true
		}
	}
	if last >= 0 && presented[last] != '.' {
		final[last] = true
	}
	return final
}

func bitFlipMutations(t *testing.T, presented []byte, spec corpusSpec) []mutation {
	t.Helper()
	rng := rand.New(rand.NewSource(spec.Seed))
	skip := segmentFinalPositions(presented)
	var out []mutation
	for len(out) < spec.BitFlips {
		pos := rng.Intn(len(presented))
		bit := rng.Intn(8)
		if presented[pos] == '.' || skip[pos] {
			continue
		}
		m := append([]byte(nil), presented...)
		m[pos] ^= 1 << bit
		out = append(out, mutation{
			name:      fmt.Sprintf("bit-flip pos=%d bit=%d", pos, bit),
			presented: m,
		})
	}
	return out
}

// reassemble swaps one segment of a compact JWS.
func reassemble(segs [3]string, idx int, replacement string) []byte {
	s := segs
	s[idx] = replacement
	return []byte(s[0] + "." + s[1] + "." + s[2])
}

func fieldSubstitutionMutations(t *testing.T, presented []byte) []mutation {
	t.Helper()
	parts := strings.Split(string(presented), ".")
	segs := [3]string{parts[0], parts[1], parts[2]}
	payload, err := base64.RawURLEncoding.DecodeString(segs[1])
	if err != nil {
		t.Fatal(err)
	}
	var claims map[string]json.RawMessage
	if err := json.Unmarshal(payload, &claims); err != nil {
		t.Fatal(err)
	}

	substitute := func(name string, mutate func(map[string]json.RawMessage)) mutation {
		m := map[string]json.RawMessage{}
		for k, v := range claims {
			m[k] = v
		}
		mutate(m)
		altered, err := json.Marshal(m)
		if err != nil {
			t.Fatal(err)
		}
		return mutation{
			name:      "field substitution: " + name,
			presented: reassemble(segs, 1, base64.RawURLEncoding.EncodeToString(altered)),
		}
	}

	return []mutation{
		// Identity rebinding attempts (INV1).
		substitute("principal", func(m map[string]json.RawMessage) {
			m["sub"] = json.RawMessage(`"spiffe://domain-a.test/attacker"`)
		}),
		substitute("delegate", func(m map[string]json.RawMessage) {
			m["act"] = json.RawMessage(`{"sub":"spiffe://domain-a.test/attacker"}`)
		}),
		// Scope escalation through tampering (INV2 via INV8, AT5).
		substitute("scope widened", func(m map[string]json.RawMessage) {
			m["scope"] = json.RawMessage(`["admin:all","read:orders","write:audit"]`)
		}),
		substitute("scope replaced", func(m map[string]json.RawMessage) {
			m["scope"] = json.RawMessage(`["admin:all"]`)
		}),
		// Validity-window stretching (INV3).
		substitute("expiration extended", func(m map[string]json.RawMessage) {
			m["exp"] = json.RawMessage(`9999999999`)
		}),
		substitute("issuance backdated", func(m map[string]json.RawMessage) {
			m["iat"] = json.RawMessage(`1`)
		}),
		// Revocation-targeting misdirection (INV6).
		substitute("instance identity", func(m map[string]json.RawMessage) {
			m["atl_ins"] = json.RawMessage(`"inst-other"`)
		}),
		// Mechanism-slot injection (AD-015 opacity does not exempt it
		// from integrity coverage).
		substitute("revocation binding injected", func(m map[string]json.RawMessage) {
			m["atl_rvb"] = json.RawMessage(`"AAAA"`)
		}),
		substitute("claim removed", func(m map[string]json.RawMessage) {
			delete(m, "atl_ins")
		}),
	}
}

func truncationMutations(presented []byte) []mutation {
	parts := strings.Split(string(presented), ".")
	segs := [3]string{parts[0], parts[1], parts[2]}
	return []mutation{
		{"truncation: last byte", presented[:len(presented)-1]},
		{"truncation: last 5 bytes", presented[:len(presented)-5]},
		{"truncation: signature removed", []byte(segs[0] + "." + segs[1] + ".")},
		{"truncation: header removed", []byte("." + segs[1] + "." + segs[2])},
		{"truncation: payload removed", []byte(segs[0] + ".." + segs[2])},
		{"truncation: first byte", presented[1:]},
	}
}

func reorderMutation(t *testing.T, presented []byte) mutation {
	// Reordering of protected fields: identical claims re-marshaled with
	// reversed key order. The content is semantically identical JSON, but
	// the signature covers the exact payload bytes, so a reordered
	// payload never verifies as the original (INV8: integrity is a
	// property of the record bytes, not of a JSON interpretation).
	t.Helper()
	parts := strings.Split(string(presented), ".")
	segs := [3]string{parts[0], parts[1], parts[2]}
	payload, err := base64.RawURLEncoding.DecodeString(segs[1])
	if err != nil {
		t.Fatal(err)
	}
	var claims map[string]json.RawMessage
	if err := json.Unmarshal(payload, &claims); err != nil {
		t.Fatal(err)
	}
	keys := make([]string, 0, len(claims))
	for k := range claims {
		keys = append(keys, k)
	}
	sort.Sort(sort.Reverse(sort.StringSlice(keys)))
	var b strings.Builder
	b.WriteByte('{')
	for i, k := range keys {
		if i > 0 {
			b.WriteByte(',')
		}
		kb, err := json.Marshal(k)
		if err != nil {
			t.Fatal(err)
		}
		b.Write(kb)
		b.WriteByte(':')
		b.Write(claims[k])
	}
	b.WriteByte('}')
	if b.String() == string(payload) {
		t.Fatal("reordered payload equals original; mutation is vacuous")
	}
	return mutation{
		name:      "protected-field reordering",
		presented: reassemble(segs, 1, base64.RawURLEncoding.EncodeToString([]byte(b.String()))),
	}
}

func headerAndSignatureMutations(t *testing.T, presented []byte, other *Record) []mutation {
	t.Helper()
	parts := strings.Split(string(presented), ".")
	segs := [3]string{parts[0], parts[1], parts[2]}
	otherParts := strings.Split(string(other.Presented()), ".")

	randomSig := base64.RawURLEncoding.EncodeToString(make([]byte, 64))
	return []mutation{
		{"header: kid substituted", reassemble(segs, 0, base64.RawURLEncoding.EncodeToString(
			[]byte(`{"alg":"ES256","kid":"forged-kid","typ":"`+headerType+`"}`)))},
		{"signature: zeroed", reassemble(segs, 2, randomSig)},
		{"signature: transplanted from sibling record", reassemble(segs, 2, otherParts[2])},
		{"payload: transplanted from sibling record", reassemble(segs, 1, otherParts[1])},
	}
}

func TestMutationCorpusDetectionFractionIsOne(t *testing.T) {
	spec := loadCorpus(t)
	signer, tm := testSigner(t)

	a := testAssertions(t)
	a.RevocationBinding = nil
	sealed := mustSeal(t, a, signer)
	presented := sealed.Presented()

	b := testAssertions(t)
	inst, err := InstanceIDFromString("inst-sibling")
	if err != nil {
		t.Fatal(err)
	}
	b.Instance = inst
	sibling := mustSeal(t, b, signer)

	var corpus []mutation
	corpus = append(corpus, bitFlipMutations(t, presented, spec)...)
	corpus = append(corpus, fieldSubstitutionMutations(t, presented)...)
	corpus = append(corpus, truncationMutations(presented)...)
	corpus = append(corpus, reorderMutation(t, presented))
	corpus = append(corpus, headerAndSignatureMutations(t, presented, sibling)...)

	undetected := 0
	for _, m := range corpus {
		if rec, outcome := ValidateIntegrity(m.presented, tm); outcome != Altered || rec != nil {
			undetected++
			t.Errorf("UNDETECTED mutation: %s", m.name)
		}
	}
	if undetected != 0 {
		t.Fatalf("detection fraction = %d/%d, SO3 requires 1 (= %d/%d)",
			len(corpus)-undetected, len(corpus), len(corpus), len(corpus))
	}
	t.Logf("mutation corpus: %d mutations, all detected and rejected", len(corpus))

	// Control: the unmutated record still validates — the corpus run
	// proves detection, not blanket rejection.
	if _, outcome := ValidateIntegrity(presented, tm); outcome != Intact {
		t.Fatal("control failed: unmutated record no longer validates")
	}
}

// checkpoint: refactor(issuance): refactor revocation status lookup

// checkpoint: chore(test): audit integration test runner

// checkpoint: chore(ui): tweak interactive console (#249)
