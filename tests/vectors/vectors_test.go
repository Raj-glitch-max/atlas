package vectors_test

import (
	"crypto/ecdsa"
	"encoding/json"
	"flag"
	"os"
	"sort"
	"testing"
	"time"

	"github.com/Raj-glitch-max/atlas/internal/record"
	"github.com/Raj-glitch-max/atlas/internal/verify"
	"github.com/Raj-glitch-max/atlas/tests/conformance"
	jose "github.com/go-jose/go-jose/v3"
	"github.com/spiffe/go-spiffe/v2/spiffeid"
)

// This test replays the committed language-neutral vectors against the
// reference verifier, reconstructing every input FROM THE JSON using only
// public APIs — exactly as an independent implementation in any language
// would. It is both the drift guard (the reference must match the published
// vectors) and the executable template a foreign implementer mirrors.

const (
	vectorPath   = "verdict-vectors.json"
	negativePath = "negative-vectors.json"
)

var update = flag.Bool("update", false, "regenerate the committed vector files from the conformance corpus")

// TestVectorsRegenerate rewrites the committed vector files from the corpus.
// Run deliberately (go test ./tests/vectors -run TestVectorsRegenerate
// -update) when the corpus changes; commit the result. Records carry fresh
// signatures each regeneration (ECDSA signing is randomized); verification is
// deterministic, so replay is stable against whatever is committed.
func TestVectorsRegenerate(t *testing.T) {
	if !*update {
		t.Skip("run with -update to regenerate committed vectors")
	}
	write := func(path string, vf conformance.VectorFile) {
		data, err := json.MarshalIndent(vf, "", "  ")
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, append(data, '\n'), 0o644); err != nil {
			t.Fatal(err)
		}
		t.Logf("wrote %d vectors to %s", len(vf.Vectors), path)
	}
	pos, err := conformance.EmitVectors()
	if err != nil {
		t.Fatalf("emit positive: %v", err)
	}
	write(vectorPath, pos)
	neg, err := conformance.EmitNegativeVectors()
	if err != nil {
		t.Fatalf("emit negative: %v", err)
	}
	write(negativePath, neg)
}

func TestVectorsReplayAgainstReferenceVerifier(t *testing.T) {
	data, err := os.ReadFile(vectorPath)
	if err != nil {
		t.Fatalf("read vectors (regenerate with -update): %v", err)
	}
	var vf conformance.VectorFile
	if err := json.Unmarshal(data, &vf); err != nil {
		t.Fatalf("parse vectors: %v", err)
	}
	if vf.Schema != conformance.SchemaVersion {
		t.Fatalf("vector schema %d, code expects %d", vf.Schema, conformance.SchemaVersion)
	}
	if len(vf.Vectors) == 0 {
		t.Fatal("no vectors")
	}

	for _, vec := range vf.Vectors {
		vec := vec
		t.Run(vec.Name, func(t *testing.T) {
			// Reconstruct inputs from JSON with public APIs only.
			domain := spiffeid.RequireTrustDomainFromString(vec.Trust.Domain)
			trust := buildTrust(t, domain, vec.Trust.Keys)
			rev := buildRevocation(t, vec.Revocation)
			now := parseTime(t, vec.Now)
			policy, err := verify.NewPolicy(seconds(vec.Policy.RSeconds), seconds(vec.Policy.SkewSeconds))
			if err != nil {
				t.Fatalf("policy: %v", err)
			}
			v, err := verify.NewVerifier(policy, trust, rev, clockPort{now})
			if err != nil {
				t.Fatalf("verifier: %v", err)
			}

			verdict, _ := v.Verify([]byte(vec.Record))

			if got := decisionToken(verdict.Decision); got != vec.Expect.Decision {
				t.Fatalf("decision = %s, vector expects %s", got, vec.Expect.Decision)
			}
			if !sameSet(causeNames(verdict.Causes), vec.Expect.Causes) {
				t.Fatalf("causes = %v, vector expects %v", causeNames(verdict.Causes), vec.Expect.Causes)
			}
		})
	}
}

// TestNegativeVectorsReplay replays the adversarial vectors: every one must
// be rejected (never Accept), and the reference verifier must match the
// committed verdict exactly. This is the interop robustness backbone — a
// foreign verifier that accepts any of these has a differential.
func TestNegativeVectorsReplay(t *testing.T) {
	data, err := os.ReadFile(negativePath)
	if err != nil {
		t.Fatalf("read negative vectors (regenerate with -update): %v", err)
	}
	var vf conformance.VectorFile
	if err := json.Unmarshal(data, &vf); err != nil {
		t.Fatalf("parse negative vectors: %v", err)
	}
	if vf.Schema != conformance.SchemaVersion {
		t.Fatalf("negative vector schema %d, code expects %d", vf.Schema, conformance.SchemaVersion)
	}
	if len(vf.Vectors) == 0 {
		t.Fatal("no negative vectors")
	}
	for _, vec := range vf.Vectors {
		vec := vec
		t.Run(vec.Name, func(t *testing.T) {
			if vec.Expect.Decision == "Accept" {
				t.Fatalf("negative vector %q expects Accept — not adversarial", vec.Name)
			}
			domain := spiffeid.RequireTrustDomainFromString(vec.Trust.Domain)
			trust := buildTrust(t, domain, vec.Trust.Keys)
			rev := buildRevocation(t, vec.Revocation)
			now := parseTime(t, vec.Now)
			policy, err := verify.NewPolicy(seconds(vec.Policy.RSeconds), seconds(vec.Policy.SkewSeconds))
			if err != nil {
				t.Fatal(err)
			}
			v, err := verify.NewVerifier(policy, trust, rev, clockPort{now})
			if err != nil {
				t.Fatal(err)
			}
			verdict, _ := v.Verify([]byte(vec.Record))
			if verdict.IsAccept() {
				t.Fatalf("adversarial record ACCEPTED — silent-acceptance differential")
			}
			if got := decisionToken(verdict.Decision); got != vec.Expect.Decision {
				t.Fatalf("decision = %s, vector expects %s", got, vec.Expect.Decision)
			}
			if !sameSet(causeNames(verdict.Causes), vec.Expect.Causes) {
				t.Fatalf("causes = %v, vector expects %v", causeNames(verdict.Causes), vec.Expect.Causes)
			}
		})
	}
}

// --- reconstruction helpers (public APIs only; a foreign impl's mirror) ---

type trustPort struct {
	domain spiffeid.TrustDomain
	mat    record.TrustMaterial
	has    bool
}

func (p trustPort) TrustMaterialFor(d spiffeid.TrustDomain) (record.TrustMaterial, bool) {
	if p.has && d == p.domain {
		return p.mat, true
	}
	return record.TrustMaterial{}, false
}

func buildTrust(t *testing.T, domain spiffeid.TrustDomain, rawKeys []json.RawMessage) verify.TrustMaterialPort {
	t.Helper()
	if len(rawKeys) == 0 {
		return trustPort{domain: domain, has: false}
	}
	keys := make(map[string]*ecdsa.PublicKey, len(rawKeys))
	for _, raw := range rawKeys {
		var jwk jose.JSONWebKey
		if err := jwk.UnmarshalJSON(raw); err != nil {
			t.Fatalf("parse JWK: %v", err)
		}
		pub, ok := jwk.Key.(*ecdsa.PublicKey)
		if !ok {
			t.Fatalf("JWK %q is not an EC public key", jwk.KeyID)
		}
		keys[jwk.KeyID] = pub
	}
	mat, err := record.NewTrustMaterial(domain, keys)
	if err != nil {
		t.Fatalf("trust material: %v", err)
	}
	return trustPort{domain: domain, mat: mat, has: true}
}

type revPort struct {
	state verify.RevocationState
	asOf  time.Time
}

func (p revPort) StatusOf(record.InstanceID) verify.RevocationStatus {
	if p.state == verify.Indeterminate {
		return verify.RevocationStatus{State: verify.Indeterminate}
	}
	return verify.RevocationStatus{State: p.state, AsOf: p.asOf}
}

func buildRevocation(t *testing.T, vr conformance.VectorRevocation) verify.RevocationStatusPort {
	t.Helper()
	var st verify.RevocationState
	switch vr.State {
	case "NotObservedRevoked":
		st = verify.NotObservedRevoked
	case "ObservablyRevoked":
		st = verify.ObservablyRevoked
	case "Indeterminate":
		st = verify.Indeterminate
	default:
		t.Fatalf("unknown revocation state %q", vr.State)
	}
	var asOf time.Time
	if vr.AsOf != nil {
		asOf = parseTime(t, *vr.AsOf)
	}
	return revPort{state: st, asOf: asOf}
}

type clockPort struct{ t time.Time }

func (c clockPort) Now() time.Time { return c.t }

func parseTime(t *testing.T, s string) time.Time {
	t.Helper()
	ts, err := time.Parse(time.RFC3339, s)
	if err != nil {
		t.Fatalf("parse time %q: %v", s, err)
	}
	return ts
}

func seconds(f float64) time.Duration { return time.Duration(f * float64(time.Second)) }

func decisionToken(d verify.Decision) string {
	switch d {
	case verify.Accept:
		return "Accept"
	case verify.Reject:
		return "Reject"
	default:
		return "InconclusiveRejected"
	}
}

func causeNames(cs []verify.Cause) []string {
	out := make([]string, 0, len(cs))
	for _, c := range cs {
		out = append(out, c.String())
	}
	return out
}

func sameSet(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	x := append([]string(nil), a...)
	y := append([]string(nil), b...)
	sort.Strings(x)
	sort.Strings(y)
	for i := range x {
		if x[i] != y[i] {
			return false
		}
	}
	return true
}
