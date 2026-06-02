package conformance_test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/Raj-glitch-max/atlas/internal/record"
	"github.com/Raj-glitch-max/atlas/internal/verify"
	"github.com/Raj-glitch-max/atlas/tests/conformance"
)

// assertionsEqual reports whether two decoded records assert identical
// content (the true unit of tamper-evidence — INV8 protects content, not the
// non-canonical transport encoding).
func assertionsEqual(a, b record.Assertions) bool {
	if a.Principal != b.Principal || a.Delegate != b.Delegate || a.Instance != b.Instance {
		return false
	}
	if !a.Expiration.Equal(b.Expiration) || !a.IssuedAt.Equal(b.IssuedAt) {
		return false
	}
	if len(a.Scope) != len(b.Scope) {
		return false
	}
	for i := range a.Scope {
		if a.Scope[i] != b.Scope[i] {
			return false
		}
	}
	return true
}

// Property / differential fuzzing of the verifier: the invariants that must
// hold for ALL inputs — the behavioral core of conformance a finite corpus
// cannot enumerate. Generators are Frankencerts-style (systematically
// malformed and adversarial presented records). A failure here is a genuine
// silent-acceptance bug; this is falsification, not demonstration.

func fuzzRNG() *rand.Rand { return rand.New(rand.NewSource(0xADACE)) }

func mustKit(t *testing.T) *conformance.Kit {
	t.Helper()
	k, err := conformance.NewKit()
	if err != nil {
		t.Fatalf("kit: %v", err)
	}
	return k
}

func freshRev() verify.RevocationStatusPort {
	return conformance.FreshRevocation(conformance.MidTime())
}

func randomRev(rng *rand.Rand) verify.RevocationStatusPort {
	switch rng.Intn(3) {
	case 0:
		return conformance.RevocationInState(verify.Indeterminate, time.Time{})
	case 1:
		return conformance.FreshRevocation(conformance.MidTime())
	default:
		return conformance.RevocationInState(verify.ObservablyRevoked, conformance.MidTime())
	}
}

func buildVerifier(t *testing.T, k *conformance.Kit, clock time.Time, rev verify.RevocationStatusPort) conformance.Verifier {
	t.Helper()
	v, err := conformance.V1Factory(conformance.DefaultPolicy(), k.Trust(), rev, conformance.Clock(clock))
	if err != nil {
		t.Fatalf("factory: %v", err)
	}
	return v
}

func buildVerifierEmptyTrust(t *testing.T, k *conformance.Kit) conformance.Verifier {
	t.Helper()
	v, err := conformance.V1Factory(conformance.DefaultPolicy(), k.EmptyTrust(), freshRev(), conformance.Clock(conformance.MidTime()))
	if err != nil {
		t.Fatalf("factory: %v", err)
	}
	return v
}

// P1 — Content tamper-evidence: flipping any single byte of a valid record
// never yields acceptance of DIFFERENT content. The signature channel must
// catch every content mutation. A flip may still Accept iff it is
// encoding-equivalent (base64 padding-bit malleability decoding to the same
// signature/content) — in which case the decoded assertions must be identical.
// Accepting differing content would be the real silent-acceptance bug.
func TestProperty_TamperNeverAcceptsDifferentContent(t *testing.T) {
	k := mustKit(t)
	valid := k.Valid()
	orig, ok := k.Decode(valid)
	if !ok {
		t.Fatal("baseline valid record did not decode")
	}
	v := buildVerifier(t, k, conformance.MidTime(), freshRev())
	if verdict, _ := v.Verify(valid); !verdict.IsAccept() {
		t.Fatalf("baseline valid record did not accept: %s", verdict.Decision)
	}

	rng := fuzzRNG()
	acceptedEquivalents := 0
	for i := 0; i < 5000; i++ {
		m := append([]byte(nil), valid...)
		pos := rng.Intn(len(m))
		m[pos] ^= 1 << uint(rng.Intn(8))
		verdict, _ := v.Verify(m)
		if !verdict.IsAccept() {
			continue
		}
		got, ok := k.Decode(m)
		if !ok || !assertionsEqual(got, orig) {
			t.Fatalf("ACCEPTED a record with different content (byte %d) — real silent-acceptance", pos)
		}
		acceptedEquivalents++ // encoding-equivalent reframe; benign (base64 malleability)
	}
	t.Logf("accepted %d encoding-equivalent reframings (base64 malleability, identical content)", acceptedEquivalents)
}

// P2 — Random garbage never accepts, and the trace is always five stages.
func TestProperty_GarbageNeverAccepts(t *testing.T) {
	k := mustKit(t)
	v := buildVerifier(t, k, conformance.MidTime(), freshRev())
	rng := fuzzRNG()
	for i := 0; i < 3000; i++ {
		b := make([]byte, rng.Intn(600))
		rng.Read(b)
		verdict, trace := v.Verify(b)
		if verdict.IsAccept() {
			t.Fatalf("random bytes ACCEPTED (len %d)", len(b))
		}
		if len(trace.Entries) != 5 {
			t.Fatalf("trace must have 5 stages even on garbage, got %d", len(trace.Entries))
		}
	}
}

// P3 — No Accept without every stage passing (structural SO5).
func TestProperty_AcceptImpliesAllPass(t *testing.T) {
	k := mustKit(t)
	rng := fuzzRNG()
	for i := 0; i < 2000; i++ {
		clock := conformance.MidTime().Add(time.Duration(rng.Intn(20000)-10000) * time.Second)
		v := buildVerifier(t, k, clock, randomRev(rng))
		verdict, trace := v.Verify(k.Valid())
		if verdict.Decision == verify.Accept {
			for _, e := range trace.Entries {
				if e.Outcome != verify.OutcomePass {
					t.Fatalf("Accept but stage %s outcome=%s", e.Check, e.Outcome)
				}
			}
		}
	}
}

// P4 — Routing consistency: a definitive cause always yields Reject; no
// definitive but an inconclusive yields InconclusiveRejected; Accept carries
// no causes. Holds for all randomized states, tampered or not.
func TestProperty_RoutingConsistency(t *testing.T) {
	k := mustKit(t)
	rng := fuzzRNG()
	for i := 0; i < 4000; i++ {
		clock := conformance.MidTime().Add(time.Duration(rng.Intn(20000)-10000) * time.Second)
		v := buildVerifier(t, k, clock, randomRev(rng))
		rec := k.Valid()
		if rng.Intn(3) == 0 {
			rec[rng.Intn(len(rec))] ^= 1 << uint(rng.Intn(8))
		}
		verdict, _ := v.Verify(rec)
		anyDef, anyInc := false, false
		for _, c := range verdict.Causes {
			if c.IsDefinitive() {
				anyDef = true
			}
			if c.IsInconclusive() {
				anyInc = true
			}
		}
		switch verdict.Decision {
		case verify.Accept:
			if len(verdict.Causes) != 0 {
				t.Fatalf("Accept carried causes %v", verdict.Causes)
			}
		case verify.Reject:
			if !anyDef {
				t.Fatalf("Reject without a definitive cause: %v", verdict.Causes)
			}
		case verify.InconclusiveRejected:
			if anyDef {
				t.Fatalf("InconclusiveRejected with a definitive cause: %v", verdict.Causes)
			}
			if !anyInc {
				t.Fatalf("InconclusiveRejected without an inconclusive cause: %v", verdict.Causes)
			}
		}
	}
}

// P5 — Absent trust material never accepts (never fetches; fails closed).
func TestProperty_AbsentMaterialNeverAccepts(t *testing.T) {
	k := mustKit(t)
	v := buildVerifierEmptyTrust(t, k)
	for i := 0; i < 500; i++ {
		if verdict, _ := v.Verify(k.Valid()); verdict.IsAccept() {
			t.Fatal("accepted with absent trust material")
		}
	}
}

// P6 — Determinism under repetition (no verdict memory).
func TestProperty_Deterministic(t *testing.T) {
	k := mustKit(t)
	v := buildVerifier(t, k, conformance.MidTime(), freshRev())
	valid := k.Valid()
	first, _ := v.Verify(valid)
	for i := 0; i < 500; i++ {
		if got, _ := v.Verify(valid); got.Decision != first.Decision {
			t.Fatalf("iteration %d: verdict changed %s -> %s", i, first.Decision, got.Decision)
		}
	}
}
