package verify_test

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"testing"
	"time"

	"github.com/Raj-glitch-max/atlas/internal/record"
	"github.com/Raj-glitch-max/atlas/internal/verify"
	"github.com/Raj-glitch-max/atlas/tests/harness"
	"github.com/spiffe/go-spiffe/v2/spiffeid"
)

// FuzzVerify is the coverage-guided fuzz target for the Verification Core —
// the OSS-Fuzz on-ramp that complements the static property tests with
// libFuzzer-style mutation. It feeds arbitrary bytes as a presented record and
// asserts the two invariants that must hold for ANY input:
//
//  1. Verify never panics and always emits a five-stage trace (structural).
//  2. Accept is sound: the core never accepts bytes that are not a genuinely
//     authentic, decodable record — the anti-silent-acceptance guarantee.
//
// Run continuously with:  go test ./internal/verify -run x -fuzz FuzzVerify
// In normal CI it executes the seed corpus as ordinary cases.
func FuzzVerify(f *testing.F) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		f.Fatal(err)
	}
	domain := spiffeid.RequireTrustDomainFromString("domain-a.test")
	tm, err := record.NewTrustMaterial(domain, map[string]*ecdsa.PublicKey{"k1": &key.PublicKey})
	if err != nil {
		f.Fatal(err)
	}
	now := time.Unix(1_800_000_300, 0).UTC()

	// Seed with a genuinely valid record and a spread of malformed inputs.
	inst, _ := record.InstanceIDFromString("inst-fuzz")
	valid, err := record.Seal(record.Assertions{
		Principal:  spiffeid.RequireFromString("spiffe://domain-a.test/principal"),
		Delegate:   spiffeid.RequireFromString("spiffe://domain-a.test/delegate"),
		Scope:      []string{"read:orders", "write:audit"},
		Expiration: now.Add(time.Hour),
		IssuedAt:   now.Add(-time.Minute),
		Instance:   inst,
	}, record.Signer{Key: key, KeyID: "k1"})
	if err != nil {
		f.Fatal(err)
	}
	f.Add(valid.Presented())
	f.Add([]byte(""))
	f.Add([]byte("not-a-jws"))
	f.Add([]byte("a.b.c"))
	f.Add([]byte("eyJhbGciOiJub25lIn0..")) // alg=none shape

	policy, err := verify.NewPolicy(time.Minute, 30*time.Second)
	if err != nil {
		f.Fatal(err)
	}
	v, err := verify.NewVerifier(
		policy,
		harness.NewTrustStore().Put(tm),
		harness.NewUniformRevocation(harness.NewClock(now), verify.NotObservedRevoked),
		harness.NewClock(now),
	)
	if err != nil {
		f.Fatal(err)
	}

	f.Fuzz(func(t *testing.T, data []byte) {
		verdict, trace := v.Verify(data) // must not panic on any input

		if len(trace.Entries) != 5 {
			t.Fatalf("trace has %d stages, want 5 (unconditional) for input %q", len(trace.Entries), data)
		}
		if verdict.IsAccept() {
			// Soundness: an Accept must correspond to a genuinely authentic,
			// decodable record. Anything else is a silent-acceptance bug.
			rec, outcome := record.ValidateIntegrity(data, tm)
			if outcome != record.Intact || rec == nil {
				t.Fatalf("ACCEPTED bytes that are not an authentic record: %q", data)
			}
			_ = rec.Read() // must not panic on an accepted record
		}
	})
}
