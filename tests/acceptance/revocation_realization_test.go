package acceptance_test

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"testing"
	"time"

	"github.com/Raj-glitch-max/atlas/internal/record"
	"github.com/Raj-glitch-max/atlas/internal/revstatus"
	"github.com/Raj-glitch-max/atlas/internal/verify"
	"github.com/Raj-glitch-max/atlas/tests/harness"
)

// End-to-end of the revocation α-path (Epic E7 composition 1, built under the
// S1–S4 scope act of 2026-07-06): the signed-revoked-set realization wired
// through the whole system, replacing the degenerate provider. This is the
// in-process realization; the two-domain SPIRE substrate run (AT13/AT14 with
// out-of-band instrumentation and two-run reproducibility) remains Epic E6/E7
// and is not claimed here.
//
// Resolved scope parameters used here: R = 2s (S1), signed-artifact pulls
// admissible (S2), passive distributor is not a broker (S3).
func TestE2E_RevocationRealization_SignedRevokedSet(t *testing.T) {
	const R = 2 * time.Second

	s := newSystem(t)
	s.clock.Set(nowMid)
	rec := s.issue(t, "read:orders", "write:audit")
	instance := rec.Read().Instance

	// Domain A's revocation authority (independent of the record signing key;
	// the RP holds its public key out of band, like a trust bundle).
	raKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	pub, err := revstatus.NewPublisher(raKey, "domain-a-revlist")
	if err != nil {
		t.Fatal(err)
	}
	prov := revstatus.NewSignedSetProvider(&raKey.PublicKey, "domain-a-revlist")

	// The propagation channel (composition-root wiring): read the authoritative
	// register (revorigin), publish a signed snapshot, deliver it to the RP,
	// which verifies and ingests it. revstatus never imports revorigin (R5);
	// the register's View crosses here, at the root.
	propagate := func(asOf time.Time) {
		var revoked []record.InstanceID
		for _, e := range s.register.View() {
			revoked = append(revoked, e.Instance)
		}
		set, err := pub.Publish(revoked, asOf)
		if err != nil {
			t.Fatalf("publish: %v", err)
		}
		if _, err := prov.Ingest(set); err != nil {
			t.Fatalf("ingest: %v", err)
		}
	}

	policy, err := verify.NewPolicy(R, 30*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	verifier, err := verify.NewVerifier(policy, s.trust, harness.AdaptRevocation(prov), s.clock)
	if err != nil {
		t.Fatal(err)
	}

	// 1. Fresh signed snapshot (empty register) → the record verifies Accept
	//    against the REAL provider (not a fake): freshness is verifiable (the
	//    snapshot's signed as-of), not asserted.
	propagate(s.clock.Now())
	if verdict, _ := verifier.Verify(rec.Presented()); !verdict.IsAccept() {
		t.Fatalf("before revocation: verdict=%s, want Accept", verdict.Decision)
	}

	// 2. Time advances (a real snapshot carries a newer signed as-of; the
	//    provider's monotone-freshness guard correctly refuses an equal-or-
	//    older snapshot, so the new one must be strictly newer). Revoke in the
	//    authoritative register, propagate a fresh signed snapshot → the
	//    record is now Reject(RevokedObservable).
	s.clock.Set(nowMid.Add(time.Second))
	s.register.Revoke(instance, s.clock.Now())
	propagate(s.clock.Now())
	verdict, _ := verifier.Verify(rec.Presented())
	if verdict.Decision != verify.Reject || !hasCause(verdict.Causes, verify.RevokedObservable) {
		t.Fatalf("after revocation: verdict=%s causes=%v, want Reject(RevokedObservable)", verdict.Decision, verdict.Causes)
	}

	// 3. Revocation is terminal: it stays rejected on re-presentation.
	if v2, _ := verifier.Verify(rec.Presented()); v2.Decision != verify.Reject {
		t.Fatalf("revocation not terminal: verdict=%s", v2.Decision)
	}
}

// The resolved R actually governs: a not-revoked snapshot older than R makes
// verification fail closed (RevocationKnowledgeStale), demonstrating the
// S1 scope act's bound in force rather than a hard-coded value.
func TestE2E_RevocationFreshnessBound_RGoverns(t *testing.T) {
	const R = 2 * time.Second

	s := newSystem(t)
	s.clock.Set(nowMid)
	rec := s.issue(t, "read:orders")

	raKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pub, _ := revstatus.NewPublisher(raKey, "domain-a-revlist")
	prov := revstatus.NewSignedSetProvider(&raKey.PublicKey, "domain-a-revlist")

	// Publish a fresh not-revoked snapshot as of now; ingest it.
	set, _ := pub.Publish(nil, s.clock.Now())
	if _, err := prov.Ingest(set); err != nil {
		t.Fatal(err)
	}

	policy, _ := verify.NewPolicy(R, 30*time.Second)
	verifier, _ := verify.NewVerifier(policy, s.trust, harness.AdaptRevocation(prov), s.clock)

	// Fresh: Accept.
	if verdict, _ := verifier.Verify(rec.Presented()); !verdict.IsAccept() {
		t.Fatalf("fresh snapshot: verdict=%s, want Accept", verdict.Decision)
	}

	// Advance the verifier clock past R without re-publishing: the held
	// snapshot's as-of is now stale (age > R) → fail closed.
	s.clock.Set(nowMid.Add(R + time.Second))
	verdict, _ := verifier.Verify(rec.Presented())
	if verdict.Decision != verify.InconclusiveRejected || !hasCause(verdict.Causes, verify.RevocationKnowledgeStale) {
		t.Fatalf("stale snapshot: verdict=%s causes=%v, want InconclusiveRejected(RevocationKnowledgeStale)", verdict.Decision, verdict.Causes)
	}
}
