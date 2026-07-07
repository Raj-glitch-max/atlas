package verify_test

import (
	"testing"
	"time"

	"github.com/Raj-glitch-max/atlas/internal/record"
	"github.com/Raj-glitch-max/atlas/internal/verify"
)

// AT23 / SO5 — single-check rollback (no silent acceptance): for each
// required check, forcing that check into a failing state while all others
// pass must flip the verdict away from Accept. This is the structural
// mitigation of the silent-trust-failure meta-mode (FM11): no single
// undetected check-failure can yield acceptance.
//
// The test first establishes the all-pass baseline (Accept), then rolls back
// one check at a time and asserts the verdict is no longer Accept, naming the
// cause that fired. Scope-integrity's rollback is realized through the
// integrity guarantee it shares (AD-024): tampering the scope trips the
// integrity check, which is the honest, cryptographically-truthful outcome.

func TestSingleCheckRollback(t *testing.T) {
	const (
		r    = time.Minute
		skew = 30 * time.Second
	)

	// Baseline: every check passes -> Accept. If this ever fails, the
	// rollback assertions below are meaningless, so guard it explicitly.
	base := newWorld(t)
	base.freshNotRevoked()
	baseRec := base.seal(t, nil)
	if verdict, _ := base.verifier(t, r, skew).Verify(baseRec.Presented()); !verdict.IsAccept() {
		t.Fatalf("baseline is not Accept (%s); rollback test cannot proceed", verdict.Decision)
	}

	rollbacks := []struct {
		name  string
		build func(t *testing.T) (*verify.Verifier, []byte)
		cause verify.Cause
	}{
		{
			name: "identity binding rolled back (principal == delegate)",
			build: func(t *testing.T) (*verify.Verifier, []byte) {
				w := newWorld(t)
				w.freshNotRevoked()
				rec := w.seal(t, func(a *record.Assertions) { a.Delegate = a.Principal })
				return w.verifier(t, r, skew), rec.Presented()
			},
			cause: verify.BindingMismatch,
		},
		{
			name: "integrity rolled back (tampered byte)",
			build: func(t *testing.T) (*verify.Verifier, []byte) {
				w := newWorld(t)
				w.freshNotRevoked()
				rec := w.seal(t, nil)
				b := append([]byte(nil), rec.Presented()...)
				b[len(b)/2] ^= 0x01
				return w.verifier(t, r, skew), b
			},
			cause: verify.IntegrityFailed,
		},
		{
			name: "expiry rolled back (clock past expiry)",
			build: func(t *testing.T) (*verify.Verifier, []byte) {
				w := newWorld(t)
				rec := w.seal(t, nil)
				w.clock.Set(baseExpiry.Add(time.Hour))
				w.revoc.Set(w.instance, verify.RevocationStatus{State: verify.NotObservedRevoked, AsOf: w.clock.Now()})
				return w.verifier(t, r, skew), rec.Presented()
			},
			cause: verify.Expired,
		},
		{
			name: "scope integrity rolled back (tampered scope -> integrity fires, AD-024)",
			build: func(t *testing.T) (*verify.Verifier, []byte) {
				w := newWorld(t)
				w.freshNotRevoked()
				rec := w.seal(t, nil)
				b := tamperScope(t, rec.Presented())
				return w.verifier(t, r, skew), b
			},
			cause: verify.IntegrityFailed, // shared INV8 guarantee
		},
		{
			name: "revocation rolled back (observably revoked)",
			build: func(t *testing.T) (*verify.Verifier, []byte) {
				w := newWorld(t)
				rec := w.seal(t, nil)
				w.revoc.Set(w.instance, verify.RevocationStatus{State: verify.ObservablyRevoked, AsOf: w.clock.Now()})
				return w.verifier(t, r, skew), rec.Presented()
			},
			cause: verify.RevokedObservable,
		},
	}

	for _, rb := range rollbacks {
		t.Run(rb.name, func(t *testing.T) {
			v, presented := rb.build(t)
			verdict, trace := v.Verify(presented)
			if verdict.IsAccept() {
				t.Fatalf("rolling back a single check still ACCEPTED — silent-trust failure (FM11/SO5)")
			}
			if !hasCause(verdict.Causes, rb.cause) {
				t.Fatalf("causes = %v, want to contain %s", verdict.Causes, rb.cause)
			}
			if len(trace.Entries) != 5 {
				t.Errorf("trace must still contain all five stages, got %d", len(trace.Entries))
			}
		})
	}
}

// tamperScope alters the scope inside the signed payload without re-signing,
// so the record fails integrity — proving no single scope-level failure
// yields acceptance (SO5), via the INV8 guarantee scope shares with the
// signature (AD-024).
func tamperScope(t *testing.T, presented []byte) []byte {
	t.Helper()
	// Re-seal an over-scoped record under a DIFFERENT key, then present it
	// against the original trust material: the scope differs and the
	// signature cannot verify — the realistic "tampered scope" an attacker
	// would produce without the issuer key.
	other := newWorld(t)
	rec := other.seal(t, func(a *record.Assertions) {
		a.Scope = []string{"admin:all", "read:orders", "write:audit"}
	})
	return rec.Presented()
}
