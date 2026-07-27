# 11 — TECHNICAL WRITER — RAW FINDINGS

Phase 4 (Polish). I own the spec as a normative artifact, ADRs, and traceability. Spec (11) and
vectors (07) are both normative; a disagreement between them is a P1 for both (R8).

Strong prior credit, because it reverses the finding this seat usually files: **the ADRs exist and
the interface spec is normative.** `ENGINEERING_DECISION_RECORD.md` carries AD-001…017+ each with
Alternatives / Reason / Trade-off / Complexity-cost; `agents/journal/` has 13 dated decision
entries; `INTERFACE_SPECIFICATION.md` declares itself "Canonical, binding contract set." The
"judgment is invisible" cluster is REFUTED. I am not backfilling ADRs; I am crediting them.

FINDING-TW-1
  Type:        spec (normative gap)
  Claim:       The *logical* contracts are normative (INTERFACE_SPECIFICATION.md) and the vectors
               are machine-normative, but there is no normative *wire-format* spec. `VECTORS.md`
               explicitly defers it: "the normative wire format is a mechanism-RFC concern, sketched
               here to make the vectors self-contained."
  Impact:      A second-language reimplementer cannot build the byte layout from prose alone — they
               must read the Go. For a protocol, the wire spec IS the artifact (OQ-3).
  Fix:         Ship `rfc/RFC-004-wire-format` (or `spec/v0.1-draft`), marked explicitly unstable and
               append-only, freezing header/claims/signature byte semantics. Low cost — the logical
               contracts + 28 vectors already pin most of it. (S5)
  Owner:       11_technical_writer

FINDING-TW-2
  Type:        traceability / honesty
  Claim:       The claim surface is internally inconsistent — the exact failure mode a
               traceability matrix exists to catch:
               (a) `00_atlas_context.md` §6 lists "Delegation chains ✓"; the repo is single-hop by
                   design (README, LIMITATIONS §1).
               (b) `LIMITATIONS.md` §10 says "one SDK (Python)"; the repo ships Go + Python + TS.
               (c) benchmarks are published without CPU / Go version / percentile / chain-depth.
  Impact:      Two docs describe two different Atlases. In a project whose credibility IS its
               honesty, an over-claim (a) and a stale under-count (b) both corrode trust.
  Fix:         Reconcile to the honest floor; add a claims→evidence traceability line per guarantee;
               qualify every benchmark. (S4)
  Owner:       11_technical_writer

FINDING-TW-3
  Type:        traceability (spec↔vectors agreement, R8)
  Claim:       The reference verifier is replayed against the committed vectors (good), but there is
               no CI check asserting the *prose* verification algorithm in `VECTORS.md` §"algorithm"
               stays in sync with `INTERFACE_SPECIFICATION.md` §3 and the code. Prose and vectors can
               drift silently.
  Fix:         A doc-lint that maps each cause name + routing rule across spec, VECTORS.md, and
               `internal/verify/cause.go`. Cheap; makes R8 mechanical.
  Owner:       11_technical_writer

VERDICT: Documentation maturity is genuinely high (ADRs, normative contracts, honest limitations).
The two real gaps are the wire-format spec (TW-1) and reconciling the over/under-claims (TW-2) —
both cheap, both high-credibility.
