---
date: 2026-07-06
slug: phase-a-primitive-discovery
artifact: founder charter — "Phase A: discover whether Atlas is an entirely new infrastructure primitive" + a large external research compendium
decision: Phase A discovery concludes Atlas is NOT a new primitive but a known design point (SPIFFE-native offline attenuable delegation ~ Biscuit/UCAN) with an honest freshness-limit treatment; the compendium describes a superset platform; adopting it is a frozen-Product-Definition amendment (founder's call), not a discovery this document asserts. External research tagged. No code, no frozen edit, no scope change applied.
agents_consulted: [empiricist, cartographer, red-team, economist, operator]
overrides: false
related_entries: [discriminating-observation-invariant, currency-of-trust-composition, freshness-composition-falsification, omega-impact-conformance]
---

# Context

The founder supplied a large external research compendium (trust/delegation/authorization + CS foundations) and a "Phase A — Primitive Discovery" charter reframing Atlas as a candidate new infrastructure primitive: forget implementation bias, produce a research design review (no code) answering nine questions about what Atlas fundamentally is — primitives, math, execution, invariants, architecture derivation, gap analysis, novelty (prove it, be brutally skeptical), research programs, implementation readiness. The charter's framing and the compendium describe a much larger system (trust query languages, VMs, graph engines, 1M-delegation scale) than the frozen Phase 7 Product Definition, which fixes Atlas narrowly (offline single-hop SPIFFE delegation, two domains, no live authority, no store).

# Decision

Recorded (asserts nothing into scope; applies no change):

- **Atlas is not a new primitive.** It is a specific, well-chosen point in the known design space of offline, public-key, attenuable capability tokens — closest to Biscuit and UCAN — distinguished by SPIFFE-nativeness and an honest treatment of the offline-revocation-freshness limit. (`PHASE-A-primitive-discovery.md` §0, §6, §7.)
- **The execution model is proof-verification of one self-contained record** (proof-carrying authorization lineage), NOT the compendium's graph traversal / query planning / trust VM — those presuppose a store Atlas does not have. (§3.)
- **Brutal-skepticism novelty verdict:** no proven publishable primitive. Atlas is an integration + a rigorous engineering realization of a known design point, plus one modest, workshop-adjacent systematization (the freshness-limit statement, from OMEGA-02/04). (§7.)
- **The compendium is a superset.** Most of it (trust DB/query/graph/VM/indexing/economics) is irrelevant to frozen Atlas and becomes relevant only under a scope amendment. Tagged in `EXTERNAL-RESEARCH-2026-07-06.md`; the on-point subset is offline authorization + revocation/compression math (feeds the E7 spike). (§8.)
- **The real fork is governance, and it is the founder's:** Path 1 (stay the narrow primitive → next action is resolving S1–S4 to unblock the built kernel) or Path 2 (become the platform → a deliberate frozen-Product-Definition amendment opening a multi-year program). This document refuses to write Path 2 as if decided (that would be scope-expansion-by-document). (§"The decision this actually surfaces".)
- **No code, no frozen edit, no new package** — per the charter's own rule that nothing exists until justified; Phase A justified nothing new. (§9.)

# Evidence cited

- Frozen Phase 7/8 (Atlas's narrow definition) vs the compendium's superset framing (§0 table).
- OMEGA-02/03/04 (the math and novelty verdicts) — related entries.
- Gap analysis vs Zanzibar/OpenFGA/OPA/Cedar/SPIFFE/OAuth-8693/Macaroons/Biscuit/UCAN/PCA/CHERI (§6) — recall-based, flagged for live verification.
- `tests/conformance/` (the executable down-payment on formal invariants).

# Council positions

## The Empiricist
The gap analysis is grounded and checkable, but the positioning ("~ Biscuit/UCAN for SPIFFE") rests on recall of those specs, not a live re-read — cap it and flag a live check before it is asserted externally. The novelty verdict (no proven primitive) is the honest one and follows from OMEGA-02/03/04's already-recorded classifications. Confidence Atlas is a known design point: High. Confidence in the exact Biscuit/UCAN delta: Medium, needs verification.

## The Red Team
The charter is itself the risk: "forget implementation bias / discover a new primitive" invites **scope-expansion-by-document** — writing a grand research architecture that silently redefines Atlas away from its frozen, mostly-built scope. The document correctly refuses that and routes Path 2 through the amendment process. Sharper danger recorded: the compendium assumes a graph / 1M-delegation / query-engine system Atlas explicitly is NOT (frozen single-hop, no store); building toward it would abandon the frozen scope and the four-test verdicts. Do not let a research charter override hash-pinned governance.

## The Operator
The useful deliverable for an adopter is exactly the honest positioning: "SPIFFE-native, public-key, offline, attenuable delegation token, honest about revocation freshness — think Biscuit/UCAN for workloads." That one sentence tells someone what Atlas is and is not, which the grand framing obscures. Approve the positioning; it is usable.

## The Economist
Blunt, and consistent with prior entries: Phase A produced zero V1 value. The project's bottleneck is the S1–S4 scope acts that unblock the built kernel's spike — not a fifth discovery pass. The compendium's research programs are multi-year and unfunded by the frozen scope. Recommendation on record: **stop discovering; decide.** Either amend scope deliberately (Path 2) or resolve S1–S4 (Path 1). Continued discovery at this point is negative ROI.

## The Cartographer
Restate: the charter asked to redefine Atlas against its own frozen definition; that is a scope-amendment decision, not a document's call. Frame surfaced: this is the **fifth** expansive excursion (OMEGA-01..04 + Phase A) and all five converge on the same answer — the narrow kernel is correct; the grand vision is deferred, known-adjacent, and fails the four tests. Convergence this consistent is a signal the question is answered. Label the compendium's programs DEFERRED; do not relabel Atlas.

# Domain anchors consulted

Not consulted — the question was epistemic/positioning, covered by the council. A distributed-systems anchor would own any future N-domain-composition program (deferred).

# Dissent preserved

- **Empiricist vs the flat "no publishable contribution":** "the clean statement of the offline-revocation-freshness impossibility plus the verify-or-exclude currency, specialized to SPIFFE delegation, is a *combination of known work* that may be a workshop-grade systematization — modest, but not nothing." Recorded as the single most defensible contribution; still classified below "new primitive."
- **Cartographer vs continued exploration:** "five convergent passes is enough; a sixth is re-litigation." Recorded; the founder may of course choose Path 2.

No seat argued Atlas *is* a new primitive.

# Founder override (if applicable)

None. This records a discovery and surfaces a decision; it commits no change and overrides nothing. Path 2, if chosen, is a separate founder amendment act.

# Open questions

- **The Path 1 / Path 2 fork** — the actual decision, the founder's.
- **Live verification** of the Biscuit/UCAN/Macaroons positioning (Empiricist).
- **S1–S4 scope acts** — still the gating item for the built kernel regardless of Path (unchanged since the Architecture Readiness Review).
- Whether the freshness-limit systematization is worth writing up externally (modest; only if Path 1 and there is appetite).

# Status
- decided: 2026-07-06
