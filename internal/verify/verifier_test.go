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

// --- shared fixtures ----------------------------------------------------

const (
	principalID = "spiffe://domain-a.test/principal"
	delegateID  = "spiffe://domain-a.test/delegate"
	keyID       = "authority-key-1"
)

var (
	baseIssued = time.Unix(1_800_000_000, 0).UTC()
	baseExpiry = time.Unix(1_800_003_600, 0).UTC() // +1h
	baseNow    = time.Unix(1_800_000_300, 0).UTC() // +5m, well within window
)

type world struct {
	signer   record.Signer
	domain   spiffeid.TrustDomain
	trust    *harness.TrustStore
	revoc    *harness.Revocation
	clock    *harness.Clock
	instance record.InstanceID
}

func newWorld(t *testing.T) *world {
	t.Helper()
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("key: %v", err)
	}
	domain := spiffeid.RequireTrustDomainFromString("domain-a.test")
	tm, err := record.NewTrustMaterial(domain, map[string]*ecdsa.PublicKey{keyID: &key.PublicKey})
	if err != nil {
		t.Fatalf("trust material: %v", err)
	}
	inst, err := record.InstanceIDFromString("inst-0001")
	if err != nil {
		t.Fatal(err)
	}
	return &world{
		signer:   record.Signer{Key: key, KeyID: keyID},
		domain:   domain,
		trust:    harness.NewTrustStore().Put(tm),
		revoc:    harness.NewRevocation(),
		clock:    harness.NewClock(baseNow),
		instance: inst,
	}
}

func (w *world) seal(t *testing.T, mutate func(*record.Assertions)) *record.Record {
	t.Helper()
	a := record.Assertions{
		Principal:  spiffeid.RequireFromString(principalID),
		Delegate:   spiffeid.RequireFromString(delegateID),
		Scope:      []string{"read:orders", "write:audit"},
		Expiration: baseExpiry,
		IssuedAt:   baseIssued,
		Instance:   w.instance,
	}
	if mutate != nil {
		mutate(&a)
	}
	rec, err := record.Seal(a, w.signer)
	if err != nil {
		t.Fatalf("seal: %v", err)
	}
	return rec
}

func (w *world) verifier(t *testing.T, r, skew time.Duration) *verify.Verifier {
	t.Helper()
	policy, err := verify.NewPolicy(r, skew)
	if err != nil {
		t.Fatalf("policy: %v", err)
	}
	v, err := verify.NewVerifier(policy, w.trust, w.revoc, w.clock)
	if err != nil {
		t.Fatalf("verifier: %v", err)
	}
	return v
}

// freshNotRevoked configures the revocation fake to report a fresh
// not-revoked observation as of the clock's now.
func (w *world) freshNotRevoked() {
	w.revoc.Set(w.instance, verify.RevocationStatus{
		State: verify.NotObservedRevoked,
		AsOf:  w.clock.Now(),
	})
}

// --- baseline accept ----------------------------------------------------

func TestVerifyBaselineAccept(t *testing.T) {
	w := newWorld(t)
	w.freshNotRevoked()
	rec := w.seal(t, nil)
	v := w.verifier(t, time.Minute, 30*time.Second)

	verdict, trace := v.Verify(rec.Presented())
	if !verdict.IsAccept() {
		t.Fatalf("verdict = %s causes=%v, want Accept", verdict.Decision, verdict.Causes)
	}
	if len(trace.Entries) != 5 {
		t.Fatalf("trace has %d entries, want 5", len(trace.Entries))
	}
	for _, e := range trace.Entries {
		if e.Outcome != verify.OutcomePass {
			t.Errorf("check %s did not pass: %s (%s)", e.Check, e.Outcome, e.Cause)
		}
	}
}

// --- definitive rejects (one per definitive cause) ----------------------

func TestVerifyDefinitiveRejects(t *testing.T) {
	t.Run("expired", func(t *testing.T) {
		w := newWorld(t)
		w.freshNotRevoked()
		rec := w.seal(t, nil)
		w.clock.Set(baseExpiry.Add(time.Hour)) // well past expiry
		w.revoc.Set(w.instance, verify.RevocationStatus{State: verify.NotObservedRevoked, AsOf: w.clock.Now()})
		v := w.verifier(t, time.Minute, time.Second)
		assertReject(t, v, rec, verify.Expired)
	})

	t.Run("revoked observable", func(t *testing.T) {
		w := newWorld(t)
		rec := w.seal(t, nil)
		w.revoc.Set(w.instance, verify.RevocationStatus{State: verify.ObservablyRevoked, AsOf: w.clock.Now()})
		v := w.verifier(t, time.Minute, time.Second)
		assertReject(t, v, rec, verify.RevokedObservable)
	})

	t.Run("binding mismatch (self-delegation)", func(t *testing.T) {
		w := newWorld(t)
		w.freshNotRevoked()
		rec := w.seal(t, func(a *record.Assertions) {
			a.Delegate = a.Principal // principal == delegate
		})
		v := w.verifier(t, time.Minute, time.Second)
		assertReject(t, v, rec, verify.BindingMismatch)
	})

	t.Run("integrity failed (tampered)", func(t *testing.T) {
		w := newWorld(t)
		w.freshNotRevoked()
		rec := w.seal(t, nil)
		tampered := append([]byte(nil), rec.Presented()...)
		// Flip a byte in the payload segment (middle of the token).
		tampered[len(tampered)/2] ^= 0x01
		v := w.verifier(t, time.Minute, time.Second)
		verdict, _ := v.Verify(tampered)
		if verdict.Decision != verify.Reject {
			t.Fatalf("verdict = %s, want Reject (IntegrityFailed)", verdict.Decision)
		}
	})
}

// --- inconclusive (fail-closed) rejects, one per inconclusive cause -----

func TestVerifyInconclusiveRejects(t *testing.T) {
	t.Run("trust material absent", func(t *testing.T) {
		w := newWorld(t)
		w.freshNotRevoked()
		rec := w.seal(t, nil)
		w.trust = harness.NewTrustStore() // hold nothing
		v := w.verifier(t, time.Minute, time.Second)
		assertInconclusive(t, v, rec, verify.TrustMaterialAbsent)
	})

	t.Run("revocation indeterminate (degenerate provider)", func(t *testing.T) {
		w := newWorld(t)
		// revoc left as fresh NewRevocation() -> Indeterminate for all
		rec := w.seal(t, nil)
		v := w.verifier(t, time.Minute, time.Second)
		assertInconclusive(t, v, rec, verify.RevocationStatusIndeterminate)
	})

	t.Run("revocation knowledge stale", func(t *testing.T) {
		w := newWorld(t)
		rec := w.seal(t, nil)
		// as-of far older than R
		w.revoc.Set(w.instance, verify.RevocationStatus{
			State: verify.NotObservedRevoked,
			AsOf:  w.clock.Now().Add(-time.Hour),
		})
		v := w.verifier(t, time.Minute, time.Second)
		assertInconclusive(t, v, rec, verify.RevocationKnowledgeStale)
	})

	t.Run("clock beyond tolerance (future issuance)", func(t *testing.T) {
		w := newWorld(t)
		w.freshNotRevoked()
		rec := w.seal(t, func(a *record.Assertions) {
			a.IssuedAt = baseNow.Add(time.Hour) // issued an hour ahead of now
			a.Expiration = baseNow.Add(2 * time.Hour)
		})
		v := w.verifier(t, time.Minute, time.Second) // tolerance 1s << 1h
		assertInconclusive(t, v, rec, verify.ClockBeyondTolerance)
	})
}

// --- no verdict memory / determinism ------------------------------------

func TestVerifyHasNoMemoryAndIsDeterministic(t *testing.T) {
	w := newWorld(t)
	w.freshNotRevoked()
	rec := w.seal(t, nil)
	v := w.verifier(t, time.Minute, 30*time.Second)

	for i := 0; i < 25; i++ {
		verdict, _ := v.Verify(rec.Presented())
		if !verdict.IsAccept() {
			t.Fatalf("iteration %d: verdict changed to %s", i, verdict.Decision)
		}
	}
	// A subsequent revoked presentation is judged fresh, uninfluenced by
	// the prior accepts (no memory).
	w.revoc.Set(w.instance, verify.RevocationStatus{State: verify.ObservablyRevoked, AsOf: w.clock.Now()})
	verdict, _ := v.Verify(rec.Presented())
	if verdict.Decision != verify.Reject {
		t.Fatalf("verdict = %s, want Reject after revocation", verdict.Decision)
	}
}

// --- construction refusals ----------------------------------------------

func TestNewPolicyRefusals(t *testing.T) {
	if _, err := verify.NewPolicy(0, time.Second); err == nil {
		t.Error("R = 0 must refuse (unparameterized verifier forbidden)")
	}
	if _, err := verify.NewPolicy(-time.Second, time.Second); err == nil {
		t.Error("negative R must refuse")
	}
	if _, err := verify.NewPolicy(time.Minute, -time.Second); err == nil {
		t.Error("negative skew must refuse")
	}
	if _, err := verify.NewPolicy(time.Minute, 0); err != nil {
		t.Errorf("zero skew is a valid explicit choice: %v", err)
	}
}

func TestNewVerifierRefusals(t *testing.T) {
	w := newWorld(t)
	good, _ := verify.NewPolicy(time.Minute, time.Second)

	if _, err := verify.NewVerifier(verify.Policy{}, w.trust, w.revoc, w.clock); err == nil {
		t.Error("zero-value (unset) policy must refuse")
	}
	if _, err := verify.NewVerifier(good, nil, w.revoc, w.clock); err == nil {
		t.Error("nil trust port must refuse")
	}
	if _, err := verify.NewVerifier(good, w.trust, nil, w.clock); err == nil {
		t.Error("nil revocation port must refuse")
	}
	if _, err := verify.NewVerifier(good, w.trust, w.revoc, nil); err == nil {
		t.Error("nil time port must refuse")
	}
	if _, err := verify.NewVerifier(good, w.trust, w.revoc, w.clock); err != nil {
		t.Errorf("fully-wired verifier must construct: %v", err)
	}
}

// --- shared assertion helpers -------------------------------------------

func assertReject(t *testing.T, v *verify.Verifier, rec *record.Record, want verify.Cause) {
	t.Helper()
	verdict, trace := v.Verify(rec.Presented())
	if verdict.Decision != verify.Reject {
		t.Fatalf("verdict = %s causes=%v, want Reject(%s)", verdict.Decision, verdict.Causes, want)
	}
	if !hasCause(verdict.Causes, want) {
		t.Fatalf("causes = %v, want to contain %s", verdict.Causes, want)
	}
	if !want.IsDefinitive() {
		t.Fatalf("%s is not a definitive cause", want)
	}
	_ = trace
}

func assertInconclusive(t *testing.T, v *verify.Verifier, rec *record.Record, want verify.Cause) {
	t.Helper()
	verdict, _ := v.Verify(rec.Presented())
	if verdict.Decision != verify.InconclusiveRejected {
		t.Fatalf("verdict = %s causes=%v, want InconclusiveRejected(%s)", verdict.Decision, verdict.Causes, want)
	}
	if !hasCause(verdict.Causes, want) {
		t.Fatalf("causes = %v, want to contain %s", verdict.Causes, want)
	}
	if !want.IsInconclusive() {
		t.Fatalf("%s is not an inconclusive cause", want)
	}
	if verdict.IsAccept() {
		t.Fatal("InconclusiveRejected must not report IsAccept")
	}
}

func hasCause(causes []verify.Cause, want verify.Cause) bool {
	for _, c := range causes {
		if c == want {
			return true
		}
	}
	return false
}
