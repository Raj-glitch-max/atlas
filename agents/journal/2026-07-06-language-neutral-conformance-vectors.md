---
date: 2026-07-06
slug: language-neutral-conformance-vectors
artifact: technical-lead decision — the Go conformance kit is Go-only; independent implementations need language-neutral vectors
decision: Built a language-neutral conformance vector suite (static JSON) generated from the corpus (single source of truth), with a normative spec (VECTORS.md) and a public-API-only replay test that is both drift guard and foreign-implementation template. Reverses the prior framing that the Go conformance kit was sufficient for the multi-implementation future.
agents_consulted: [empiricist, red-team, operator, economist, cartographer]
overrides: false
related_entries: [omega-impact-conformance, e7-alpha-signed-revoked-set]
---

# Context

Acting as technical lead (founder grant, 2026-07-06), north star = technical
inevitability: what will future contributors — independent implementations,
vendors, universities — desperately wish existed. Inventory finding: the Go
conformance kit (`tests/conformance`, built last turn) is **Go-only**; it
requires importing this module. A Rust/Python/Zig implementer cannot use it.
Every trust/crypto primitive that actually got independent implementations
ships **language-neutral test vectors** (Wycheproof, JOSE cookbook,
WebAuthn/FIDO, Ed25519 RFC, CBOR). Without them, "independent implementations"
ends in verifier differentials — the Frankencerts failure this project's own
OMEGA-impact analysis named as the 5-year killer.

# Decision

Built the language-neutral conformance vector suite:

- `tests/conformance/vectors.go` — `EmitVectors` serializes the corpus to
  static JSON (`Vector`/`VectorFile`), recording the **reference verifier's
  actual verdict** per case (so vectors are reference behavior by
  construction; TestV1Conformance separately ties reference to spec).
- `tests/vectors/verdict-vectors.json` — 10 committed vectors, one per corpus
  scenario, spanning the full verdict space (Accept, every definitive and
  inconclusive cause, precedence). Records carry real ES256 signatures; keys
  as RFC 7517 JWKs.
- `tests/vectors/vectors_test.go` — `TestVectorsReplayAgainstReferenceVerifier`
  reconstructs every input from the JSON using **public APIs only** (the
  mirror of what a foreign implementer does) and asserts the reference matches;
  `TestVectorsRegenerate -update` regenerates deliberately.
- `tests/vectors/VECTORS.md` — the normative spec: file shape, the record wire
  format (enough to consume the vectors), the five-check algorithm + routing,
  the cause names, and the conformance criterion.

**Reversal recorded:** last turn I framed the Go conformance kit as "the"
conformance answer. It is necessary but not sufficient for a multi-language
future. The language-neutral vectors are the actual interop backbone; the Go
kit is now the *generator and reference-replay* of them, not the deliverable.

# Evidence cited

- `verdict-vectors.json` (10 vectors) + green
  `TestVectorsReplayAgainstReferenceVerifier` (all scenarios).
- Prior art: Wycheproof, JOSE cookbook, WebAuthn/FIDO conformance, Ed25519 RFC
  test vectors — the universal pattern for interoperable primitives.

# Council positions

## The Empiricist
The vectors are authoritative-by-construction (reference output) and the
replay is public-API-only, so it genuinely models a foreign consumer, not a
privileged one. Honest limit: 10 vectors cover the verdict space but are not
adversarially exhaustive — the property fuzzer (Go) still carries the
all-inputs guarantee; vectors are the *interop* backbone, fuzzing is the
*robustness* backbone. Both needed.

## The Red Team
This is the correct 5-year artifact and it closes the differential-testing gap
I flagged (an oracle a second implementation actually runs). Two things on
record: (1) the vectors' `record` bytes are frozen with real signatures — good,
verification is deterministic; but regeneration mints fresh signatures, so the
committed file must be regenerated deliberately and reviewed, never in CI. (2)
VECTORS.md sketches the wire format to make vectors self-contained, but the
*normative* wire format is still a mechanism-RFC concern — do not let the
sketch quietly become the spec without the RFC.

## The Operator
A foreign implementer now has everything: the bytes, the keys, the inputs, the
expected verdicts, and a prose algorithm. `Run vectors, match expect` is a
one-paragraph onboarding. This is the artifact that makes "independent
implementation" a weekend, not a research project.

## The Economist
High-inevitability, low-cost: it reuses the existing corpus as the single
source of truth, adds no product complexity (test artifact), and is exactly
what an adopter/vendor needs. Best use of the technical-lead latitude so far.

## The Cartographer
Restate: the conformance story is now spec (INTERFACE_SPECIFICATION §3 +
VECTORS.md) → Go kit (generator + reference replay) → language-neutral vectors
(the portable artifact) → foreign implementations. No scope change; the frozen
package is untouched; this is the interop layer a primitive with independent
implementations inevitably grows.

# Dissent preserved

No dissent. Red Team's caveats (regenerate-deliberately; the wire sketch is
not the normative RFC) recorded as binding, not disagreement.

# Founder override (if applicable)

None; built under the technical-lead grant of 2026-07-06.

# Open questions / next

- A **normative wire-format RFC** for the record (VECTORS.md only sketches it)
  — the natural next spec artifact; a mechanism RFC, governance-gated.
- **Adversarial/negative vectors** beyond the verdict-space corpus (malformed
  JWS families, algorithm-confusion, boundary timestamps) — the Go property
  fuzzer covers these in-language; a language-neutral *negative* vector set
  would extend interop robustness testing. Candidate next artifact.
- Still unchanged: the E6 substrate spike validates the offline thesis and
  needs real SPIRE infrastructure.

# Status
- decided: 2026-07-06

<!-- checkpoint: repo(fuzzing-strategy): update fuzzing strategy -->

<!-- checkpoint: governance(security-invariants): audit security invariants (#39) -->
