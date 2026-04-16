package verify_test

import (
	"testing"
	"time"

	"github.com/Raj-glitch-max/atlas/internal/record"
	"github.com/Raj-glitch-max/atlas/internal/verify"
	"github.com/Raj-glitch-max/atlas/tests/harness"
)

// The decision trace is the reviewer's instrument (SO8) and the substrate of
// single-check rollback (SO5): it must be emitted unconditionally — Accepts
// included — carry all five stages in canonical order, record the policy in
// force, and expose per-stage outcomes without leaking material.

func TestTraceEmittedOnAccept(t *testing.T) {
	w := newWorld(t)
	w.freshNotRevoked()
	rec := w.seal(t, nil)
	v := w.verifier(t, time.Minute, 30*time.Second)

	verdict, trace := v.Verify(rec.Presented())
	if !verdict.IsAccept() {
		t.Fatalf("precondition: want Accept, got %s", verdict.Decision)
	}
	if len(trace.Entries) != 5 {
		t.Fatalf("Accept must still carry all five stages, got %d", len(trace.Entries))
	}
	// An Accept trace must be as inspectable as a reject trace: every stage
	// present and passing, verdict recorded.
	if trace.Verdict.Decision != verify.Accept {
		t.Errorf("trace verdict = %s, want Accept", trace.Verdict.Decision)
	}
}

func TestTraceCanonicalOrderAndCompleteness(t *testing.T) {
	w := newWorld(t)
	w.freshNotRevoked()
	rec := w.seal(t, nil)
	v := w.verifier(t, time.Minute, 30*time.Second)
	_, trace := v.Verify(rec.Presented())

	want := []verify.CheckName{
		verify.CheckIdentityBinding,
		verify.CheckIntegrity,
		verify.CheckExpiry,
		verify.CheckScope,
		verify.CheckRevocation,
	}
	if len(trace.Entries) != len(want) {
		t.Fatalf("trace has %d entries, want %d", len(trace.Entries), len(want))
	}
	for i, name := range want {
		if trace.Entries[i].Check != name {
			t.Errorf("entry %d = %s, want %s", i, trace.Entries[i].Check, name)
		}
		if trace.Entries[i].InputsDigest == "" {
			t.Errorf("entry %s has no inputs digest (reproducibility)", name)
		}
	}
}

func TestTraceRecordsPolicyAndTime(t *testing.T) {
	w := newWorld(t)
	w.freshNotRevoked()
	rec := w.seal(t, nil)
	const r, skew = 42 * time.Second, 7 * time.Second
	v := w.verifier(t, r, skew)

	_, trace := v.Verify(rec.Presented())
	if trace.Policy.R != r || trace.Policy.SkewTolerance != skew {
		t.Errorf("trace policy = {R:%s skew:%s}, want {R:%s skew:%s}",
			trace.Policy.R, trace.Policy.SkewTolerance, r, skew)
	}
	if !trace.TimeReading.Equal(w.clock.Now()) {
		t.Errorf("trace time = %v, want injected clock %v", trace.TimeReading, w.clock.Now())
	}
}

func TestTraceGatedStagesMarkedNotEvaluated(t *testing.T) {
	// When integrity does not pass, read-dependent stages are recorded
	// NotEvaluated — honest, not omitted; the trace still has all five.
	w := newWorld(t)
	w.freshNotRevoked()
	rec := w.seal(t, nil)
	w.trust = harness.NewTrustStore() // force TrustMaterialAbsent at integrity
	v := w.verifier(t, time.Minute, time.Second)

	verdict, trace := v.Verify(rec.Presented())
	if verdict.Decision != verify.InconclusiveRejected {
		t.Fatalf("want InconclusiveRejected, got %s", verdict.Decision)
	}
	if len(trace.Entries) != 5 {
		t.Fatalf("gated trace must still carry five stages, got %d", len(trace.Entries))
	}
	integrity, _ := trace.Find(verify.CheckIntegrity)
	if integrity.Outcome != verify.OutcomeInconclusive || integrity.Cause != verify.TrustMaterialAbsent {
		t.Errorf("integrity = %s/%s, want Inconclusive/TrustMaterialAbsent", integrity.Outcome, integrity.Cause)
	}
	for _, name := range []verify.CheckName{verify.CheckIdentityBinding, verify.CheckExpiry, verify.CheckScope, verify.CheckRevocation} {
		e, _ := trace.Find(name)
		if e.Outcome != verify.OutcomeNotEvaluated {
			t.Errorf("stage %s = %s, want NotEvaluated (gated by integrity)", name, e.Outcome)
		}
	}
}

func TestTraceDigestsAreStableAndInputSensitive(t *testing.T) {
	// Same inputs -> same digests (reproducibility); different record ->
	// different integrity digest (the digest actually covers inputs).
	w := newWorld(t)
	w.freshNotRevoked()
	rec := w.seal(t, nil)
	v := w.verifier(t, time.Minute, 30*time.Second)

	_, t1 := v.Verify(rec.Presented())
	_, t2 := v.Verify(rec.Presented())
	e1, _ := t1.Find(verify.CheckIntegrity)
	e2, _ := t2.Find(verify.CheckIntegrity)
	if e1.InputsDigest != e2.InputsDigest {
		t.Error("identical verifications produced different integrity digests")
	}

	other := w.seal(t, func(a *record.Assertions) {}) // same content, fresh signature
	_, t3 := v.Verify(other.Presented())
	e3, _ := t3.Find(verify.CheckIntegrity)
	// other is byte-identical content but a fresh signature (ECDSA is
	// randomized), so the presented bytes differ -> digest differs.
	if e3.InputsDigest == e1.InputsDigest {
		t.Error("distinct presented bytes must yield distinct integrity digests")
	}
}

// checkpoint: feat(sdk): add key derivation
