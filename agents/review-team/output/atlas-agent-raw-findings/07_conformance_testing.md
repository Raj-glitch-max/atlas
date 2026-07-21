# 07 — CONFORMANCE & TESTING ENGINEER — RAW FINDINGS

Phase 1 (Foundation). I asked "what's the oracle?" of every test and mutated the verifier in my
head against each property. The rigor here is real, not decorative — I will say so plainly, and
then name the two gaps that matter.

What exists and bites:
- **Property tests** (`tests/conformance/properties_test.go`) P1–P6: tamper-never-accepts-different-
  content (5000 bit-flips, oracle = decoded-content equality), garbage-never-accepts (3000),
  accept-implies-all-pass, routing consistency (4000, tampered + randomized states),
  absent-material-never-accepts, determinism. These are genuine property tests with real
  iteration counts, not example tests in costume.
- **Adversarial vectors** `negative-vectors.json` (18) — generation *fails closed* if any is
  accepted. **This is the crown-jewel discipline** and it is present.
- **Language-neutral vectors** `verdict-vectors.json` + `negative-vectors.json`, schema-versioned,
  documented in `VECTORS.md`, replayed against the reference verifier in CI
  (`TestVectorsReplayAgainstReferenceVerifier`). Regeneration is append-only by convention.
- **Fuzzing** `FuzzVerify` (821k execs claimed) for fail-closed totality.

---

FINDING-TEST-1
  Severity:      MAJOR
  Category:      differential
  Gap:           "Differential testing" is claimed but no CI job runs the three shipped SDKs
                 (Go/Python/TS) against the shared vectors and asserts identical verdicts. The
                 vectors are replayed only against the Go reference. Cross-impl parity — the whole
                 reason the vectors are language-neutral — is asserted, not tested.
  Invariant:     ∀ vector v, ∀ impl i,j: verdict_i(v) == verdict_j(v).
  Oracle:        differential (the vectors are the spec; divergence is a bug in ≥1 impl).
  Required test: A CI matrix job: Python + TS + Go each consume verdict/negative vectors, compare
                 decision + cause-set. A single mismatch fails the build. (S9)
  Owner:         07_conformance_testing

FINDING-TEST-2
  Severity:      MAJOR
  Category:      conformance
  Gap:           Duplicate-JSON-key handling is unspecified. Go's decode is last-wins; a first-wins
                 reimplementation diverges (CRYPTO-9). No vector pins it.
  Invariant:     A payload with a duplicate object key resolves identically (reject) in all impls.
  Oracle:        exact (the vector states the expected verdict).
  Required test: Add `ATLAS-CONF-DUP-KEY`: authentic-sig payload with duplicate `scope` → reject.
                 Add the corresponding spec sentence (→ 03, 11). (S9)
  Owner:         07_conformance_testing

FINDING-TEST-3
  Severity:      MINOR
  Category:      conformance
  Gap:           The vector set is `schema: 1` and append-by-convention, but there is no *suite
                 versioning* (`conformance/v1/…`) an SDK can declare conformance to, and no mutation-
                 testing job proving the suite would catch a verifier regression. Coverage is
                 reported; mutation score is not.
  Invariant:     n/a (process).
  Oracle:        n/a.
  Required test: (a) freeze a named suite version SDKs can cite; (b) a nightly `gremlins`/mutation
                 run on `internal/verify` with a tracked floor. Nightly, not per-PR (06's budget).
  Owner:         07_conformance_testing

**REFUTED red flags (checked, and the suite already handles them):**
- Attenuation monotonicity is an example test → **N/A**: single-hop, no multi-hop operator to
  property-test. The single-hop strict-subset property is unit-tested at issuance. Correct scope.
- Property tests with too-few iterations / fixed seed → the counts are 2k–5k and the seed is
  logged; adequate. Approve.
- Adversarial vectors not frozen → they are frozen and CI-gated. Approve.
- Fail-closed totality unfuzzed → `FuzzVerify` covers it. Approve.

## MUTATION-SCORE TARGET & PROPERTY-TEST GAPS
- Mutation-score floor for `internal/verify`: set an initial floor of 80% on a nightly run, ratchet
  up. Currently unmeasured — that is the one number missing from an otherwise strong story.
- Invariants currently *example/unit*-tested that are fine to leave so (single-hop subset), and
  invariants that are genuinely property-tested and should stay (tamper, routing, fail-closed).
  Nothing here needs promotion; the gap is *differential-in-CI*, not property coverage.

VERDICT: **APPROVE-WITH-CHANGES.** The property + adversarial-vector discipline is above bar. Close
the differential-in-CI gap (TEST-1) and pin duplicate-key semantics (TEST-2) — those are the two
tests that could catch the bug that will actually ship in a second SDK.
