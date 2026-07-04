# System Invariants

**Phase:** 8 — Engineering Requirements.
**Source authority:** the frozen Phase 7 Product Definition package. Each invariant traces to a Product Definition item.
**What this document is not:** no architecture, technology, repository, API, or protocol. Invariants state properties that **any** compliant implementation must hold true **at all times, in every state**, regardless of how it is built — and regardless of which outcome the C4 feasibility spike produces.

## Discipline

1. **Always-true, by construction.** An invariant here is asserted to hold for every compliant implementation, in every state, at all times. It is a constraint the system **cannot violate**, not a behavior it performs.
2. **Trace, don't invent.** Each invariant carries a `Traces to` line to a Product Definition item. Invariants are restatements of structural properties the Product Definition already commits to; they do not strengthen it.
3. **Spike-outcome-independent.** Because the C4 spike outcome (composition works / technology gap / logical impossibility / unresolvable) is not yet known, no invariant may depend on a spike result. Invariants hold whichever way S1–S5 are resolved and whichever composition (if any) the spike selects.
4. **Hypothesis properties are NOT invariants.** A property sourced from a `[HYPOTHESIS]` Product Definition item is not guaranteed-always-true until V1 confirms it. Such properties are listed under **Candidate invariants (not established)** — they become invariants only if V1 validates them.
5. **The S4 impossibility is honored.** No invariant claims revocation observability across a partition that isolated the relying party from the issuer at the moment of revocation; that is an information-theoretic limit (gate S4, deducible from C6/FR5). See INV12.

## Established invariants

**INV1 — Identity-binding integrity.** A delegation that verifies successfully binds exactly one principal identity and exactly one delegate identity, both deterministically recoverable from the delegation's content together with trust material the relying party already holds. No successful verification shall attribute a delegation to a principal or delegate other than those established at issuance.
*Traces to:* FR1, FR6, NFR6. *Always-true because:* a delegation is a single presentable unit binding both identities (FR1); tamper-evidence (NFR6) prevents silent rebinding.

**INV2 — No scope escalation.** No delegation shall grant a scope that is not a strict subset of the principal's own permissions. The system shall never produce a delegation whose scope exceeds the principal's permission set.
*Traces to:* FR2. *Always-true because:* FR2 requires the delegation scope to be a strict subset of the principal's permissions and issuance refuses otherwise.

**INV3 — Time-bounded validity monotonicity.** No delegation shall be treated as valid by a conformant verifier after its expiration time. Once expired, a delegation shall never return to a valid state.
*Traces to:* FR3. *Always-true because:* expiration (FR3) is monotone and terminal w.r.t. validity; a conformant verifier rejects an expired delegation.

**INV4 — Revocation monotonicity and terminality.** Once a specific delegation is revoked, it shall never return to a valid state. Revocation of a delegation is one-way and terminal for that delegation.
*Traces to:* FR4 ("a specific delegation … is subsequently rejected"). *Always-true because:* revocation, once effected, is permanent for the targeted delegation.

**INV5 — Revocation independence from underlying identity.** Revoking a delegation shall not invalidate the underlying identity of the principal or of the delegate.
*Traces to:* FR4(b). *Always-true because:* FR4 explicitly requires the underlying identities to remain "valid and unaffected" by the revocation.

**INV6 — Revocation target singularity.** A revocation targets exactly one delegation. Revoking one delegation is never, by itself, a sufficient cause for a conformant verifier to reject a different, unrevoked delegation.
*Traces to:* FR4 ("a specific delegation") and `PRODUCT_DEFINITION.md` headline ("independent revocability"). *Always-true because:* "specific" + "independent" mean per-delegation granularity; no spillover to siblings. *(This is the direct reading of "specific"/"independent," not a new capability — see Discipline §2.)*

**INV7 — No live call on the core verification path.** The core verification path — determining that a presented delegation is valid, unexpired, unrevoked, and scope-correct — shall not require a network call to a shared authority at the moment of verification.
*Traces to:* FR5, NFR2, C6. *Always-true because:* C6 makes the no-live-call property non-negotiable within scope; a design requiring one fails the Product Definition by definition.

**INV8 — Tamper-evidence of the delegation record.** A delegation record altered after creation shall not verify as the original record; tampering is always detectable by a conformant verifier.
*Traces to:* NFR6, FR6. *Always-true because:* NFR6 requires any alteration to be detectable; a presentation that passes verification therefore is, by construction, an unaltered record.

**INV9 — Reconstruction-record self-sufficiency.** A delegation's reconstruction record shall be verifiable by a third party without access to the original verification event or the original verifier's runtime state.
*Traces to:* FR6 ("independent of and after the original verification event"). *Always-true because:* FR6 requires the record alone to be sufficient to determine delegator, delegate, scope, and time.

**INV10 — Base-identity boundary.** The system shall not assume responsibility for issuing base workload identity; that responsibility remains with the external workload-identity infrastructure. The system operates on already-issued identity material.
*Traces to:* C1, `SYSTEM_CONTEXT.md` ("outside this boundary"). *Always-true because:* C1 fixes the boundary; the system is a companion, not a replacement.

**INV11 — Companion, not a spec modification.** The system's operation shall not require, and shall not depend on, any change to the definition of the existing workload-identity standard it extends.
*Traces to:* C2, NFR5. *Always-true because:* C2 forbids a change to the standard's core; the system must be satisfiable against the standard as published.

**INV12 — Observability-claim upper bound.** The system's revocation-observability claim shall not exceed what is achievable without a live call to a shared authority. In particular, no claim shall assert that a relying party observes a revocation performed while that relying party is partitioned from the issuer, earlier than partition recovery.
*Traces to:* C6, FR5; sharpened by `LEVEL0_1_FEASIBILITY_GATE.md` S4 (an information-theoretic limit deducible from "no live call ⇒ no fresh information"). *Always-true because:* the no-live-call constraint (C6) caps what the system can ever honestly claim; the S4 limit is a physical fact, not a design choice. *(Traced primarily to C6/FR5; S4 is the evidence that pins the bound.)*

## Candidate invariants (not established — pending V1 validation)

These are properties the Product Definition offers only as `[HYPOTHESIS]`. They are not asserted as invariants because they are not guaranteed-always-true until V1 confirms them. They are recorded so they are not silently lost and so `05_ACCEPTANCE_TEST_PLAN.md` can target them.

**C-INV1 — Fail-closed under inconclusive verification.** `[HYPOTHESIS]` When a conformant verifier cannot reach a conclusive validity determination, it shall reject rather than accept.
*Traces to:* NFR3, ER11, SO4. *Becomes an invariant only if* V1 (`DEFERRED.md` D6) confirms the behavior under adversarial/ambiguous verification conditions.

*(No other hypothesis property of the Product Definition is of invariant character: FR9 (cross-protocol interop) and FR10 (posture-tied validity) are capability/feature hypotheses, not always-true structural properties.)*

## Coverage

Every Product Definition requirement that has an always-true structural character is captured as an invariant:

| Source | Invariant |
|---|---|
| FR1 | INV1 |
| FR2 | INV2 |
| FR3 | INV3 |
| FR4(a) | INV4 |
| FR4(b) | INV5 |
| FR4 ("specific"/"independent") | INV6 |
| FR5, NFR2, C6 | INV7 |
| FR6, NFR6 | INV8, INV9 |
| C1 | INV10 |
| C2, NFR5 | INV11 |
| C6, FR5 (+ gate S4) | INV12 |
| NFR3 `[HYP]` | C-INV1 (candidate) |

## Provenance

- **Primary source:** the Phase 7 Product Definition package (frozen).
- **Secondary (evidence for INV12's bound):** `LEVEL0_1_FEASIBILITY_GATE.md` §"Strengthen" S4.
- **Confidence:** INV1–INV12 — High, as restatements of structural commitments in the (evidenced, V1-in-scope) Product Definition items FR1–FR8, NFR2/NFR4/NFR6, C1–C2/C6. C-INV1 — inherits NFR3's hypothesis status; not established.
- **Change-condition:** only an amendment to the frozen Product Definition authorizes adding, removing, or strengthening an invariant. Resolving S4 by scope act fixes the observability bound INV12 expresses; it does not remove INV12.
