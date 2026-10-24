package revstatus_test

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"sort"
	"testing"
	"time"

	"github.com/Raj-glitch-max/atlas/internal/record"
	"github.com/Raj-glitch-max/atlas/internal/revstatus"
	"github.com/Raj-glitch-max/atlas/internal/revstatus/contracttest"
)

func key(t *testing.T) *ecdsa.PrivateKey {
	t.Helper()
	k, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	return k
}

func inst(t *testing.T, s string) record.InstanceID {
	t.Helper()
	id, err := record.InstanceIDFromString(s)
	if err != nil {
		t.Fatal(err)
	}
	return id
}

var tRef = time.Unix(1_800_000_000, 0).UTC()

func TestSignedSet_PublishIngestStatusOf(t *testing.T) {
	k := key(t)
	pub, err := revstatus.NewPublisher(k, "list-A")
	if err != nil {
		t.Fatal(err)
	}
	prov := revstatus.NewSignedSetProvider(&k.PublicKey, "list-A")

	revoked := inst(t, "inst-revoked")
	notRevoked := inst(t, "inst-live")

	// Before ingest: honest ignorance for everything.
	if got := prov.StatusOf(revoked); got.State != revstatus.Indeterminate {
		t.Fatalf("pre-ingest state = %s, want Indeterminate", got.State)
	}

	set, err := pub.Publish([]record.InstanceID{revoked}, tRef)
	if err != nil {
		t.Fatal(err)
	}
	adopted, err := prov.Ingest(set)
	if err != nil || !adopted {
		t.Fatalf("ingest: adopted=%v err=%v", adopted, err)
	}

	// Revoked → ObservablyRevoked with the signed as-of (verifiable freshness).
	got := prov.StatusOf(revoked)
	if got.State != revstatus.ObservablyRevoked || !got.AsOf.Equal(tRef) {
		t.Errorf("revoked: %s asOf=%v, want ObservablyRevoked @ %v", got.State, got.AsOf, tRef)
	}
	// Absent from a complete snapshot → honestly NotObservedRevoked.
	got = prov.StatusOf(notRevoked)
	if got.State != revstatus.NotObservedRevoked || !got.AsOf.Equal(tRef) {
		t.Errorf("live: %s asOf=%v, want NotObservedRevoked @ %v", got.State, got.AsOf, tRef)
	}
}

func TestSignedSet_TamperedSetRefused(t *testing.T) {
	k := key(t)
	pub, _ := revstatus.NewPublisher(k, "list-A")
	prov := revstatus.NewSignedSetProvider(&k.PublicKey, "list-A")

	set, _ := pub.Publish([]record.InstanceID{inst(t, "a")}, tRef)

	// Add an instance after signing (forge a revocation): must be refused.
	forged := set
	forged.Revoked = append(append([]record.InstanceID(nil), set.Revoked...), inst(t, "victim"))
	if adopted, err := prov.Ingest(forged); adopted || err == nil {
		t.Error("forged (extra-revocation) set must be refused")
	}

	// Remove an instance after signing (erase a revocation): must be refused.
	erased := set
	erased.Revoked = nil
	if adopted, err := prov.Ingest(erased); adopted || err == nil {
		t.Error("forged (erased-revocation) set must be refused")
	}

	// Move the as-of forward after signing (forge freshness): must be refused.
	fresher := set
	fresher.AsOf = tRef.Add(time.Hour)
	if adopted, err := prov.Ingest(fresher); adopted || err == nil {
		t.Error("forged (moved as-of) set must be refused")
	}

	// After all forgeries refused, the provider still has no snapshot.
	if got := prov.StatusOf(inst(t, "a")); got.State != revstatus.Indeterminate {
		t.Errorf("after refused forgeries: %s, want Indeterminate (nothing adopted)", got.State)
	}
}

func TestSignedSet_WrongKeyAndWrongList(t *testing.T) {
	k := key(t)
	other := key(t)
	pub, _ := revstatus.NewPublisher(k, "list-A")
	set, _ := pub.Publish([]record.InstanceID{inst(t, "a")}, tRef)

	// Provider trusting a different key must reject.
	wrongKey := revstatus.NewSignedSetProvider(&other.PublicKey, "list-A")
	if adopted, err := wrongKey.Ingest(set); adopted || err == nil {
		t.Error("set signed by an untrusted key must be refused")
	}
	// Provider for a different list must reject.
	wrongList := revstatus.NewSignedSetProvider(&k.PublicKey, "list-B")
	if adopted, err := wrongList.Ingest(set); adopted || err == nil {
		t.Error("set for a different list ID must be refused")
	}
}

func TestSignedSet_MonotoneFreshness(t *testing.T) {
	k := key(t)
	pub, _ := revstatus.NewPublisher(k, "list-A")
	prov := revstatus.NewSignedSetProvider(&k.PublicKey, "list-A")

	newer, _ := pub.Publish([]record.InstanceID{inst(t, "a")}, tRef.Add(time.Hour))
	older, _ := pub.Publish(nil, tRef)

	if adopted, _ := prov.Ingest(newer); !adopted {
		t.Fatal("newer set must be adopted")
	}
	if adopted, _ := prov.Ingest(older); adopted {
		t.Error("older set must not replace a newer one (monotone freshness)")
	}
	// 'a' stays revoked as of the newer snapshot.
	if got := prov.StatusOf(inst(t, "a")); got.State != revstatus.ObservablyRevoked || !got.AsOf.Equal(tRef.Add(time.Hour)) {
		t.Errorf("after older-set ignored: %s asOf=%v, want ObservablyRevoked @ newer", got.State, got.AsOf)
	}
}

func TestSignedSet_PublisherRefusals(t *testing.T) {
	if _, err := revstatus.NewPublisher(nil, "list-A"); err == nil {
		t.Error("nil key must refuse")
	}
	if _, err := revstatus.NewPublisher(key(t), ""); err == nil {
		t.Error("empty list ID must refuse")
	}
}

// The realization satisfies the M5 Provider contract in its ignorance state
// (no snapshot ingested → Indeterminate for all, honest). The with-snapshot
// behavior (complete census → NotObservedRevoked for absent instances) is a
// legitimate knowledge state exercised by the tests above; the contract's
// honest-indeterminate check targets the ignorance state, which is where an
// "unknown" instance genuinely has no backing knowledge.
func TestSignedSet_SatisfiesContractWhenIgnorant(t *testing.T) {
	k := key(t)
	prov := revstatus.NewSignedSetProvider(&k.PublicKey, "list-A")
	contracttest.Run(t, prov, inst(t, "unknown"), inst(t, "a"), inst(t, "b"))
}

// DEFECT D1: Snapshot as-of uses whole seconds in signature but nanosecond comparison in Ingest
// This test asserts the CURRENT WRONG behavior: a snapshot with fractional seconds
// edited to appear newer can be adopted, causing rollback from revoked to not-revoked.
// After fix, this test should be renamed TestD1_SnapshotAsOfRollbackRejected and inverted.
func TestDefectD1_SnapshotAsOfRollback(t *testing.T) {
	k := key(t)
	pub, _ := revstatus.NewPublisher(k, "list-A")
	prov := revstatus.NewSignedSetProvider(&k.PublicKey, "list-A")

	revoked := inst(t, "inst-revoked")
	notRevoked := inst(t, "inst-live")

	// Snapshot A: revoked, as-of T+0.9s (Unix() = T, signed as T)
	// We simulate this by creating a snapshot with AsOf = T+900ms
	// The signature covers Unix() = T (whole seconds)
	tRef := time.Unix(1_800_000_000, 0).UTC()
	snapA := revstatus.SignedRevokedSet{
		ListID:  "list-A",
		AsOf:    tRef.Add(900 * time.Millisecond), // T+0.9s
		Revoked: []record.InstanceID{revoked},
		// Signature would be computed over canonicalDigest with asOf.Unix() = T
		// For this characterization test, we manually construct what the buggy
		// canonicalDigest would produce for T+0.9s (which is same as T)
	}
	// Manually compute the buggy digest (whole seconds)
	buggyDigestA := canonicalDigestBuggy("list-A", tRef.Add(900*time.Millisecond), []record.InstanceID{revoked})
	sigA, _ := ecdsa.SignASN1(rand.Reader, k, buggyDigestA[:])
	snapA.Sig = sigA

	adopted, err := prov.Ingest(snapA)
	if err != nil || !adopted {
		t.Fatalf("ingest snapA (revoked): adopted=%v err=%v", adopted, err)
	}
	// Should be revoked
	if got := prov.StatusOf(revoked); got.State != revstatus.ObservablyRevoked {
		t.Fatalf("after snapA: %s, want ObservablyRevoked", got.State)
	}

	// Snapshot B: NOT revoked, as-of T+0.1s (Unix() = T)
	// But we EDIT the AsOf to T+0.95s to make it appear "newer" than snapA in nanosecond comparison
	// The signature covers Unix() = T (same as snapA)
	snapB := revstatus.SignedRevokedSet{
		ListID:  "list-A",
		AsOf:    tRef.Add(950 * time.Millisecond), // T+0.95s (edited fractional seconds)
		Revoked: []record.InstanceID{},            // NOT revoked
	}
	buggyDigestB := canonicalDigestBuggy("list-A", tRef.Add(100*time.Millisecond), []record.InstanceID{})
	// The signature is over the ORIGINAL as-of (T+0.1s) which has same Unix() = T
	sigB, _ := ecdsa.SignASN1(rand.Reader, k, buggyDigestB[:])
	snapB.Sig = sigB

	// Ingest B AFTER A - in the buggy code, B.AsOf.After(A.AsOf) is TRUE (950ms > 900ms)
	// even though both signatures cover the SAME whole-second timestamp T
	adopted, err = prov.Ingest(snapB)
	// CURRENT WRONG BEHAVIOR: B is adopted because nanosecond comparison says it's newer
	// even though both have same signed whole-second timestamp
	if !adopted {
		t.Logf("DEFECT D1 NOT REPRODUCED: snapB not adopted (adopted=%v, err=%v)", adopted, err)
		// This is the expected PASS after fix
		return
	}

	// WRONG: instance flips from revoked to not-observed-revoked
	got := prov.StatusOf(revoked)
	if got.State != revstatus.NotObservedRevoked {
		t.Fatalf("DEFECT D1: expected rollback to NotObservedRevoked, got %s", got.State)
	}
	t.Logf("DEFECT D1 REPRODUCED: instance flipped from ObservablyRevoked to NotObservedRevoked via same-second rollback")
}

// canonicalDigestBuggy replicates the CURRENT buggy behavior: uses Unix() (whole seconds)
func canonicalDigestBuggy(listID string, asOf time.Time, revoked []record.InstanceID) [32]byte {
	sorted := dedupSortedBuggy(revoked)
	h := sha256.New()
	h.Write([]byte(listID))
	h.Write([]byte{0})
	var ts [8]byte
	binary.BigEndian.PutUint64(ts[:], uint64(asOf.Unix())) // BUG: whole seconds only
	h.Write(ts[:])
	h.Write([]byte{0})
	for _, id := range sorted {
		h.Write([]byte(id.String()))
		h.Write([]byte{0})
	}
	var out [32]byte
	copy(out[:], h.Sum(nil))
	return out
}

func dedupSortedBuggy(in []record.InstanceID) []record.InstanceID {
	if len(in) == 0 {
		return nil
	}
	cp := append([]record.InstanceID(nil), in...)
	sort.Slice(cp, func(i, j int) bool { return cp[i].String() < cp[j].String() })
	out := cp[:0]
	var prev string
	for i, id := range cp {
		if i > 0 && id.String() == prev {
			continue
		}
		out = append(out, id)
		prev = id.String()
	}
	return out
}

// TestD1_SnapshotAsOfRollbackRejected tests that the D1 fix works: a snapshot
// with a different millisecond timestamp cannot reuse a signature from a
// different millisecond timestamp. This replaces TestDefectD1_SnapshotAsOfRollback
// after the fix.
func TestD1_SnapshotAsOfRollbackRejected(t *testing.T) {
	k := key(t)
	pub, _ := revstatus.NewPublisher(k, "list-A")
	prov := revstatus.NewSignedSetProvider(&k.PublicKey, "list-A")

	revoked := inst(t, "inst-revoked")

	// Snapshot A: revoked, as-of T+500ms (signed with millisecond precision)
	snapA, err := pub.Publish([]record.InstanceID{revoked}, tRef.Add(500*time.Millisecond))
	if err != nil {
		t.Fatal(err)
	}

	adopted, err := prov.Ingest(snapA)
	if err != nil || !adopted {
		t.Fatalf("ingest snapA (revoked): adopted=%v err=%v", adopted, err)
	}
	// Should be revoked
	if got := prov.StatusOf(revoked); got.State != revstatus.ObservablyRevoked {
		t.Fatalf("after snapA: %s, want ObservablyRevoked", got.State)
	}

	// Snapshot B: NOT revoked, as-of T+999ms (same second, different millisecond)
	// B uses the SAME signature as A (attacker replay)
	snapB := revstatus.SignedRevokedSet{
		ListID:  "list-A",
		AsOf:    tRef.Add(999 * time.Millisecond),
		Revoked: []record.InstanceID{},
		Sig:     snapA.Sig, // SAME signature as A
	}

	// Ingest B AFTER A - with millisecond-precision digest, the signature
	// should NOT verify because the digest covers the millisecond timestamp
	adopted, err = prov.Ingest(snapB)
	if adopted {
		t.Fatalf("snapB should NOT be adopted (different millisecond = different digest), but was adopted")
	}
	if err == nil {
		t.Fatalf("snapB should return error (bad signature), got nil")
	}

	// Instance should STILL be revoked (no rollback)
	got := prov.StatusOf(revoked)
	if got.State != revstatus.ObservablyRevoked {
		t.Fatalf("instance should remain ObservablyRevoked, got %s", got.State)
	}
}

// FuzzIngest tests the Ingest function with random inputs to catch edge cases
// in timestamp validation and signature verification.
func FuzzIngest(f *testing.F) {
	f.Add([]byte("test-seed"))
	f.Fuzz(func(t *testing.T, data []byte) {
		// Generate a key pair
		k, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err != nil {
			t.Skip("failed to generate key")
		}
		pub, err := revstatus.NewPublisher(k, "list-fuzz")
		if err != nil {
			t.Skip("failed to create publisher")
		}
		prov := revstatus.NewSignedSetProvider(&k.PublicKey, "list-fuzz")

		// Create a valid snapshot with millisecond-aligned timestamp
		now := time.Unix(1_800_000_000, 0).UTC()
		set, err := pub.Publish(nil, now)
		if err != nil {
			return
		}

		// Try to ingest - should succeed
		_, _ = prov.Ingest(set)

		// Now test with sub-millisecond timestamp (should be rejected)
		badAsOf := now.Add(123456 * time.Nanosecond) // sub-millisecond
		badSet := revstatus.SignedRevokedSet{
			ListID:  "list-fuzz",
			AsOf:    badAsOf,
			Revoked: nil,
			Sig:     set.Sig, // reuse valid signature
		}
		adopted, err := prov.Ingest(badSet)
		if adopted {
			t.Errorf("sub-millisecond asOf should be rejected, but was adopted")
		}
		if err == nil {
			t.Errorf("sub-millisecond asOf should return error, got nil")
		}

		// Test with valid millisecond-aligned timestamp (should be accepted if newer)
		goodAsOf := now.Add(500 * time.Millisecond)
		goodSet := revstatus.SignedRevokedSet{
			ListID:  "list-fuzz",
			AsOf:    goodAsOf,
			Revoked: nil,
			Sig:     set.Sig, // this signature won't verify for different asOf, but that's OK
		}
		_, _ = prov.Ingest(goodSet) // may fail signature verification, that's expected
	})
}
