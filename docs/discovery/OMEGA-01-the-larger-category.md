# OMEGA-01 — The larger category Atlas is a subset of

**Status:** Discovery record (ADR-class). Additive; changes no existing document, no code, no plan. Presented for a founder decision. Not frozen.
**Prompt of origin:** the "Project OMEGA" first-principles adversarial exercise (2026-07-06): reason from physics/information-theory/economics to the infrastructure primitive a fully agent-mediated software future requires, destroy candidates adversarially, and test whether Atlas is the wrong question.
**Doctrine held throughout:** confidence requires evidence; hypotheses stay hypotheses; the frozen package is untouched.

---

## 1. The surviving primitive

After adversarial elimination (§4), one primitive survives:

> **A portable, composable, offline-verifiable attestation** — a carried, self-describing, cryptographically-bound claim that any verifier can check without a live call to a shared authority, honest about a bounded freshness limit.

Authority-to-act ("X may act for Y, scope S, until T, revocably") is **one claim type** of this primitive. Provenance ("artifact A was produced from inputs I by actor P under authority a") is another. Posture ("actor P was in state S at time T") is another. The unifying object is the *envelope + offline verification model*, not any single claim schema.

The load-bearing physical constraint: in a world of many mutually-distrusting agents deciding at machine speed, **trust must be carried, not fetched** — a live shared-authority call is a scaling bottleneck, a control chokepoint, a privacy leak (the authority observes every verification), and impossible across an adversarial partition. This is information theory, not preference.

## 2. Atlas is a subset — evidenced against the repository itself

| Larger-category property | Where Atlas already instances it |
|---|---|
| Signed, self-describing, offline-verifiable claim bundle | The delegation record, M1 (`internal/record`) — generalize the *claim type* and it is the primitive |
| Authority claim | FR1–FR5, the whole delegation model |
| **Provenance** claim | FR6 (reconstruction record: who delegated to whom, with what scope, when) |
| **Posture** claim | FR10 `[HYPOTHESIS]` (validity tied to verification-time posture) |
| **Composability** across non-adopting parties | FR9 `[HYPOTHESIS]` (cross-protocol interop) |
| Bound to portable workload identity | ER15/INV10 (companion to SPIFFE, operates on already-issued identity) |
| Honest freshness bound | INV12 + the S4 information-theoretic limit |

Atlas selected the smallest falsifiable slice of this primitive: **one claim type (authority), two domains, single hop, SPIFFE-coexisting.** That is correct engineering, not a limitation.

## 3. The reframe of Atlas's hardest problem (the "I hadn't considered that")

`LEVEL0_1_FEASIBILITY_GATE.md`'s S4 finding — a revocation performed while the relying party is partitioned from the issuer is *information-theoretically* unobservable before recovery — has been treated as Atlas's central risk (C4).

Reframed against the larger category, **S4 is the universal impossibility theorem for all carried-trust systems.** Any 2038 future built on carried (not fetched) trust hits that exact wall. Atlas is therefore the first project to (a) characterize the fundamental freshness bound of carried trust and (b) build the minimal falsifiable artifact that tests whether the bound is livable. The C4 problem is Atlas's claim to seriousness, not its weakness.

## 4. Adversarial elimination log (candidates destroyed)

- **Pure delegatable authority.** Necessary but a subset: answers "may X act for Y," not "what produced this" or "can I rely on it." Survives *as* a subset.
- **General verifiable computation / proof-carrying results.** Answers "computed correctly," not "under whose authority" or "was the spec right"; general proofs for arbitrary execution are not credibly cheap even by 2038. Component, not the primitive.
- **Attestable provenance.** Contains authority as one edge type; strong candidate. Merges with authority into the §1 primitive.
- **"Just call the authority" (connectivity kills offline).** Destroyed by the carried-not-fetched argument (§1) — a physics constraint, partition-fatal.
- **"Carried trust is stale trust" (revocation kills carried trust).** Does not destroy the primitive; it *defines its hardest sub-problem* — the S4 bound (§3), which is fundamental, not fixable.

## 5. Verdict against the four constraint tests — for the *grand rewrite*

A rewrite of Atlas into "the universal attestation primitive/platform" **fails all four**:

1. **Small team?** No — it is in-toto + SPIFFE + verifiable-compute + a new trust algebra. Boil-the-ocean.
2. **Principal engineers respect it?** No — universal frameworks read as red flags; a narrow, correct, falsifiable primitive with a proven bound reads as serious.
3. **OSS adopts it?** No — adoption comes through the SPIFFE-coexisting wedge (NFR5/SO7), not a replace-everything platform.
4. **More long-term value?** No net gain — the option value of the larger category is captured almost entirely by *naming* the generalization at the M1 boundary, at ~1% of the cost of building it.

**Therefore: do not rewrite, do not delete. The current plan is the correct minimal wedge.** Discarding it would destroy the single hardest thing already gotten right.

## 6. The one cheap move that captures the option value (founder decision)

Three paths; pick one. None requires touching the frozen package or the current epic ordering.

- **Path A — Record the finding, change nothing structural (recommended default).** Keep Atlas exactly as planned; this document preserves the insight so a future scope act can generalize deliberately, on evidence.
- **Path B — Adopt a zero-cost naming reframe at M1.** Rename the M1 abstraction conceptually from "Delegation Record" to **"Attestation Envelope," with "Delegation" as its first schema** — the record already carries an opaque claim structure (AD-013/AD-015), so this is naming + one interface note, not new mechanism. Cost: an interface-spec amendment and a doc pass. Benefit: the generalization becomes a first-class extension point (a new schema is a new claim type, not a fork), materially improving the SO7 "primitive-vs-feature" story the frozen package already prizes.
- **Path C — Open a V2 research track (deferred).** After the V1 verdict, open a scoped exploration of a second claim type (provenance) as the first proof the envelope generalizes. Explicitly out of V1 (C4 horizon); recorded here so it is not lost.

**Recommendation:** Path A now (it is free and loses nothing), with Path B queued for the next natural interface-spec amendment (it is nearly free and strengthens the adoption thesis), and Path C as an explicit post-V1 option. Rewrite/delete is rejected on the §5 evidence.

## 7. What was deliberately NOT done, and why

No file was deleted, no code rewritten, no plan altered, no frozen document touched — despite the originating prompt granting authority to do all of it. The technically-justified conclusion (§5) is that destruction fails the founder's own four tests; performing it would be sunk-cost reasoning in reverse (destroying value to look bold). The honest, higher-value act was to find that Atlas is a *correct subset* of a larger primitive and to price the cheapest way to keep the larger door open. Engineering continues at `IMPLEMENTATION_MASTER_PLAN.md` Epic E3 pending the founder's §6 selection.

<!-- checkpoint: fix(stores): fix signature validation -->
