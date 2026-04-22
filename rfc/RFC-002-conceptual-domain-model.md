# RFC-002 — Conceptual Domain Model

## 1. Header

| | |
|---|---|
| RFC | 002 |
| Title | Conceptual Domain Model |
| Status | Draft |
| Author | founder |
| Date | 2026-07-05 |
| Supersedes | none |
| Superseded by | none |

This RFC defines the conceptual domain model of the system: the concepts that exist in the problem domain, the relationships among them, the state transitions they undergo, the invariants that constrain them, and the concepts explicitly out of scope. It is implementation-independent: it names no classes, database schema, APIs, repository layout, technologies, storage, protocols, deployment, or implementation. Every concept traces to a frozen Product Definition or Engineering Requirements item. This is the conceptual substrate on which later, mechanism-bearing RFCs will be built; it adds no mechanism and asserts no runtime behavior.

## 2. Source authority

- Frozen Phase 7 Product Definition package: `FUNCTIONAL_REQUIREMENTS.md` (FR1–FR10), `NON_FUNCTIONAL_REQUIREMENTS.md` (NFR1–NFR6), `CONSTRAINTS.md` (C1–C6), `USE_CASE_CATALOG.md` (UC1–UC8), `USER_MODEL.md`, `SYSTEM_CONTEXT.md`, `V1_SCOPE.md`, `DEFERRED.md` (D1–D9), `PRODUCT_DEFINITION.md`.
- Frozen Phase 8 Engineering Requirements: `01_ENGINEERING_REQUIREMENTS.md` (ER1–ER17), `02_SECURITY_OBJECTIVES.md` (SO1–SO8), `03_SYSTEM_INVARIANTS.md` (INV1–INV12, C-INV1 candidate), `04_FAILURE_MODEL.md` (FM1–FM11), `05_ACCEPTANCE_TEST_PLAN.md` (AT1–AT30).
- Frozen `LEVEL0_1_FEASIBILITY_GATE.md` — for the scope parameters S1–S5 and the spike outcomes α (composition-works) / β (technology-gap) / γ (logical-impossibility at S4) / δ (unresolvable). The same frozen doc establishes the founder scope-act as the authority that resolves S1–S5 (an off-runtime governance role, not a problem-domain concept).
- RFC-000 (architectural principles + the Architecture Review Process envelope this RFC follows) and RFC-001 (System Context — the boundary this model is the inside of), both unfrozen drafts presented for founder review.
- Founder instruction, 2026-07-04: produce RFC-002 — Conceptual Domain Model — answering the five questions in §3, forbidding the categories in §3, tracing every concept, and stopping after.

## 3. Scope of this RFC

This RFC answers:

- **What concepts exist** (§7) — the named concepts of the problem domain, each with a forcing trace.
- **What relationships exist** (§8) — the named relations among concepts, with cardinality stated in words.
- **What state transitions exist** (§9) — the named states and labeled transitions of the lifecycle concepts.
- **What invariants constrain the model** (§10) — the conceptual constraints every compliant realization must respect.
- **What concepts are explicitly out of scope** (§11) — concepts considered and excluded, each with a trace.

This RFC deliberately does **not** define (each deferred to a later, separately authorized RFC, or out of scope per the frozen package):

- **Classes.** No class definitions, no field lists, no type hierarchies. Concepts are named and related in prose, not sketched as types.
- **Database schema.** No tables, columns, rows, keys, indices, or persistence layout. The model is silent on how concepts are represented at rest.
- **APIs.** No operation signatures, request/response shapes, or call contracts.
- **Repository layout.** No module or package decomposition.
- **Technologies.** No language, runtime, library, or crypto primitive.
- **Storage.** No persistence substrate or data-at-rest arrangement.
- **Protocols.** No wire format, transport, or message encoding. The revocation-propagation mechanism varieties (push, pull, cached-pull, brokered) are the frozen S2/S3 candidate categories, named only to declare their deferral, not adopted here.
- **Deployment.** No topology, runtime placement, or deployment arrangement.
- **Implementation.** No code, no internal data structures, no algorithm specifics.

Concepts here are described as domain vocabulary — what *exists in the problem* and what *relates to what* — not as constructs of any realization. Cardinality is stated in plain language ("exactly one", "zero or more"), never as schema notation.

## 4. Traces — each conceptual decision mapped to frozen forcing item

| Conceptual decision | Forced by |
|---|---|
| The model has Principal and Delegate as distinct identity concepts, both recoverable from a delegation | FR1, ER1, INV1 |
| The model has a Permission Set concept (the principal's own permissions) | FR2, ER2, INV2, SO6 |
| A Delegation's Scope is a strict subset of the principal's Permission Set | FR2, ER2, INV2, SO6 |
| A Delegation has an Expiration | FR3, ER3, INV3 |
| The model has a tamper-evident, self-sufficient Reconstruction Record | FR6, NFR6, ER4, INV8, INV9, SO3 |
| The model has Trust Domains, bounded to exactly two | C3, FR8, ER8, ER17 |
| The model has a Conformant Relying Party and a Non-conformant Relying Party (the latter outside the trust boundary) | ER7, FM10, SO2 |
| Verification uses locally-held Trust Material with no live call | FR5, NFR2, C6, ER7, INV7 |
| The model has a Revocation concept, one-way and targeting exactly one delegation | FR4, ER5, ER6, INV4, INV5, INV6, SO1, NFR4 |
| The model has a Revocation-Information Source and a Partition concept (S4) | FM1, FM2, FM4, INV12; S2, S3, S4 (gate) |
| The model has an Independent Reviewer / Third Party concept | FR6, ER4, INV9, SO8, R6 |
| The model has a Verification Verdict with Accept / Reject / Inconclusive[HYP] | SO5, FR3; NFR3, ER11, SO4, C-INV1 (Inconclusive[HYP]) |
| The model excludes base-identity issuance and standard modification | C1, C2, INV10, INV11, ER15, ER16 |
| The model excludes multi-hop, ≥3 domains, non-SPIFFE, replay-resistance, key-compromise-resistance, posture-attestation, cross-protocol interop, latency SLA, buyer, auditor, trust-bridging | DEFERRED D1–D9; SO Non-objectives; FM5, FM8; S5; ER17 |
| A concept is present only if a frozen item forces it | RFC-000 AP11, DR1, Discipline §"trace, don't invent" (03_SYSTEM_INVARIANTS.md, 04_FAILURE_MODEL.md) |

## 5. RFC-000 principle compliance

- **AP1 (offline):** the Verification concept's core path requires no live call to a shared authority (M-INV7; §9.2 Accept transition).
- **AP2 (companion):** Principal/Delegate identity are *consumed* by the model, never *issued* (M-INV10); the model requires no change to the existing standard (M-INV11).
- **AP3 (tamper-evidence):** the Reconstruction Record concept is tamper-evident and self-sufficient (M-INV8, M-INV9).
- **AP4 (fail-closed, `[HYPOTHESIS]`):** the Inconclusive state and the Inconclusive→Reject transition are carried as `[HYPOTHESIS]` in the model (C-M-INV1, §9.2), not established.
- **AP5 (honesty):** §13 carries FM5 / FM8 / S4 / C-INV1 limits into the model as named, honest labels — no concept claims resistance the package does not warrant.
- **AP6 (two-domain):** the Trust Domain concept is bounded to exactly two (§7, §14).
- **AP7 (parametric):** revocation-observability transitions (§9.3) are stated for any value of R and the S4 bound, with no fixed mechanism.
- **AP8 (independent reviewability):** the Independent Reviewer concept reconstructs from the Record alone (M-INV9).
- **AP9 / AP13 (adversarial, observable):** §12 names the observable wrong-state consequence of each load-bearing conceptual choice; each falsifies against a frozen item.
- **AP10 (no over-issuance):** the Issuance transition (§9.1) refuses over-scope (M-INV2, SO6, FM6); the verification-time scope check is integrity (INV8), not subset re-derivation.
- **AP11 (minimize accidental complexity):** no concept is introduced without a forcing trace (§4); the model is conceptual, adds no mechanism, no component, no layer. The founder scope-act role is referenced only via its frozen LEVEL0_1 basis, not modeled as a problem-domain concept.
- **AP12 (stable module boundaries):** the stable concepts (Delegation, Reconstruction Record, Verification) are isolated from the volatile mechanism concept (revocation-information mechanism — S2/S3) which the model names but does not fix.
- **TP1–TP6, DR1–DR10:** the trade-off philosophy and decision rules bind later mechanism RFCs; DR1 (traceability) is honored by §4; DR7 (hypothesis-preservation) by §15; DR10 (freeze-eligibility) is the posture in §18.

## 6. Spike-outcome analysis

The conceptual model is spike-outcome-independent by construction: its concepts, relationships, state transitions, and invariant constraints are stated independently of *how* they are realized. The model supports all four spike outcomes without modification:

- **α (composition-works):** the model's concepts and transitions hold; a mechanism RFC then realizes them.
- **β (technology-gap):** the model holds; the gap is a mechanism problem for a later RFC, not a conceptual one.
- **γ (logical-impossibility at S4):** already established and carried in the model as the honest Partition / Not-yet-Observable bound (§9.3, M-INV12).
- **δ (unresolvable):** the model holds; an unresolvable S-parameter is bounded (M-INV12), not solved, here.

The model precludes nothing conceptually and defers every mechanism question to later RFCs.

## 7. Concepts (what concepts exist)

Each concept names a thing in the problem domain. A concept is present only because a frozen item forces it (DR1, AP11). Concepts are grouped for readability; the grouping adds no semantics.

### A. Identity concepts

- **Principal Identity** — the identity (workload, agent, or organization) on whose behalf a delegation is made; sourced from existing infrastructure, never issued by the system. *Forces:* FR1, ER1, INV1, INV10, C1; `SYSTEM_CONTEXT.md` "External entities".
- **Delegate Identity** — the workload or agent that presents proof of delegation; may be long-running or ephemeral (the ephemeral case is a distinguishing requirement, not an edge case). *Forces:* FR1, ER1, FR7, ER10, UC6, INV1; `SYSTEM_CONTEXT.md`.
- **Permission Set** — the set of permissions held by a principal, against which a delegation's scope is subset-checked. *Forces:* FR2 ("the principal's own permissions"), ER2, INV2, SO6.

### B. Delegation concepts

- **Delegation** — the single presentable unit that binds a principal and a delegate, carries a scope and an expiration, and produces a reconstruction record. *Forces:* FR1–FR6, ER1–ER6, INV1–INV6, INV8, INV9 (INV7 forces the *Verification* concept, not the Delegation concept); `PRODUCT_DEFINITION.md` "what this product is". What individuates two issuances to the same (principal, delegate, scope, expiration) — i.e., what makes a revoked delegation a *distinct instance* from a freshly-issued one to the same logical grant — is not fixed by any frozen item (frozen FM5 explicitly calls this "a distinct question the Product Definition does not address"); the model therefore does not name a delegation-instance-identity concept, and M-INV4/M-INV6 (§10) apply to a delegation *instance* whose instance-identity is deferred to a later mechanism RFC (see §17).
- **Scope** — a permission set that is a **strict subset** of the bound principal's Permission Set; inspectable by the relying party. The subset relation is enforced at the **Issuance** transition; at **Verification** the scope is checked for integrity (tamper-evidence), not re-derived as a subset (the relying party need not hold the principal's full Permission Set). *Forces:* FR2, ER2, INV2, SO6 (issuance); INV8 (verification-time integrity).
- **Expiration** — a time after which a conformant verifier treats the delegation as invalid. *Forces:* FR3, ER3, INV3.
- **Reconstruction Record** (the record produced for each delegation) — a tamper-evident, self-sufficient unit from which a third party can determine which identity delegated to which, with what scope, and at what time, independent of and after the original verification event. *Forces:* FR6, NFR6, ER4, INV8, INV9, SO3.

### C. Boundary concepts

- **Trust Domain** — an independently-operated trust boundary. Exactly two are in scope: domain A (containing the principal) and domain B (containing the relying party), with no shared authority between them and federation disabled. *Forces:* C3, FR8, ER8, ER17.
- **Conformant Relying Party** — a verifier that exercises the documented checks; the locus at which the system's guarantees are defined. The conformant-verifier boundary — the surface across which a Conformant Relying Party verifies — is part of this concept: the system's guarantees are defined at this boundary (`04_FAILURE_MODEL.md` FM10). *Forces:* ER7, SO2, FM10 (boundary), RFC-001 §10.1.
- **Non-conformant Relying Party** — a verifier that does not exercise the documented checks, or is malicious; outside the system's trust boundary by definition. *Forces:* FM10.
- **Trust Material** — the material a relying party holds locally and uses to verify a delegation offline; never fetched live during core verification. *Forces:* FR5, ER1, ER7, INV7, NFR2, C6.

### D. Revocation concepts

- **Revocation** — the one-way, terminal invalidation of a *specific* delegation, targeting exactly one delegation and not affecting underlying identities. *Forces:* FR4, ER5, ER6, INV4, INV5, INV6, SO1, NFR4.
- **Revocation-Information Source** — the external locus from which a revocation becomes observable to the relying party. The model names the source as a concept; its mechanism (push, pull, cached-pull, or brokered — the frozen S2/S3 candidate categories, named here only to declare their deferral, not adopted) is the spike-selected composition (S2/S3) and is deferred to a later RFC. *Forces:* FM1, FM2, FM4, INV12; S2, S3 (gate).
- **Partition** — the condition in which the relying party is isolated from the revocation-information source at the moment of revocation; an information-theoretic bound on observability (no live call ⇒ no fresh information). *Forces:* FM1, INV12; S4 (gate).
- **Revocation-Observability State** — per revocation, the relying party's view of whether that revocation is observable: **Observable** or **Not-yet-Observable**, bounded by R in non-partitioned operation and by partition recovery when a partition is in effect. *Forces:* SO1, FM2, FM4, INV12; R (S1), S4 (gate).

### E. Verification concepts

- **Verification** — the relying party's determination of a presented delegation's validity, scope, expiry, and revocation status from the presented material together with locally-held trust material, with no live call to a shared authority. *Forces:* FR5, ER7, INV7, SO5.
- **Verification Verdict — Accept** — the outcome when every required check passes: identity binding (M-INV1), integrity of the presented Scope (INV8 — the tamper-evidence check on the scope; the scope-subset property itself is enforced at issuance, M-INV2/SO6, and is not a verification-time check), expiry (M-INV3), revocation-observability within R (FR4, SO1), signature/tamper (M-INV8). Offline, with no live call (M-INV7). *Forces:* SO5, FR3 baseline (AT6), INV1, INV3, INV8.
- **Verification Verdict — Reject** — the outcome when any single required check fails — single-check rollback (SO5, FM11). Expired (M-INV3), Revoked-and-Observable (M-INV4, SO1), tampered record or tampered Scope (M-INV8) each force Reject. Over-scope (M-INV2) is refused at the Issuance transition (§9.1), not re-derived at verification; the verification-time scope check is integrity (INV8), not subset re-derivation. *Forces:* SO5, FM11, INV3, INV4, INV8; R1.
- **Verification Verdict — Inconclusive** `[HYPOTHESIS]` — the outcome when a required check cannot be conclusively determined (e.g., signature unverifiable, required trust material unavailable, clock beyond stated tolerance). Under the fail-closed hypothesis this state transitions to Reject. *Forces:* NFR3, ER11, SO4, C-INV1, FM3, FM9.

### F. Reviewability concept

- **Independent Reviewer / Third Party** — an entity that reproduces a verdict from the Reconstruction Record alone, with no access to the original verifier's runtime state and no privileged information. *Forces:* FR6, ER4, INV9, SO8, R6.

(The founder scope-act — resolving S1–S5 — is an off-runtime governance role, not a problem-domain concept; its frozen basis is `LEVEL0_1_FEASIBILITY_GATE.md`. RFC-state advancement and freeze authorization are governed by RFC-000's Architecture Review Process and `agents/agents/GOVERNANCE.md`, both referenced in §17 as non-frozen envelopes, and are not modeled as domain concepts here.)

## 8. Relationships (what relationships exist)

Named relations among concepts, each with a forcing trace and cardinality stated in words.

- **binds** — a Delegation *binds* exactly one Principal Identity and exactly one Delegate Identity; both are deterministically recoverable from the Delegation together with the relying party's trust material. *(1 Delegation : 1 Principal : 1 Delegate.)* *Forces:* INV1, FR1, ER1.
- **scopes-to** — a Delegation *scopes-to* exactly one Scope, which is a strict subset of the bound Principal's Permission Set. *(1 Delegation : 1 Scope; Scope ⊂ Permission Set.)* *Forces:* INV2, FR2, ER2, SO6.
- **expires-at** — a Delegation *expires-at* exactly one Expiration time. *(1 Delegation : 1 Expiration.)* *Forces:* INV3, FR3, ER3.
- **produces** — a Delegation *produces* exactly one Reconstruction Record; a Reconstruction Record *reconstructs* exactly one Delegation. *(1 Delegation : 1 Record.)* *Forces:* FR6, ER4, INV8, INV9.
- **held-by** — a Principal *held-by* exactly one Permission Set. *(1 Principal : 1 Permission Set.)* *Forces:* FR2, ER2, INV2.
- **presents** — a Delegate *presents* a Delegation to a Relying Party, crossing the two-domain trust boundary. *(1 presentation : 1 Delegation : 1 RP.)* *Forces:* FR1, ER1, FR8, ER8.
- **contains** — Trust Domain A *contains* the Principal; Trust Domain B *contains* the relying party; exactly two domains, no shared authority. *(2 domains, fixed roles.)* *Forces:* C3, ER8, ER17.
- **holds-locally** — a Conformant Relying Party *holds-locally* Trust Material sufficient for offline verification. *(1 RP : its trust material.)* *Forces:* FR5, INV7, NFR2, C6.
- **verifies-with** — Verification *verifies-with* a presented Delegation together with locally-held Trust Material, with no live call to a shared authority. *(1 Verification : 1 Delegation, no live call.)* *Forces:* INV7, ER7, FR5.
- **targets** — a Revocation *targets* exactly one Delegation instance. *(1 Revocation : 1 Delegation instance.)* *Forces:* INV6, ER5, FR4. (What constitutes the instance-identity that a revocation keys to is deferred — see §7B note and §17.)
- **does-not-affect** — a Revocation *does-not-affect* the Principal's or the Delegate's underlying identity. *(1 Revocation : 0 affected underlying identities.)* *Forces:* INV5, ER5, FR4.
- **observes** — a Conformant Relying Party *observes* a Revocation's observability state, bounded by R in non-partitioned operation and by partition recovery when a partition is in effect; observability is never claimed for in-partition revocations. *Forces:* SO1, FM2, FM4, INV12.
- **isolates** — a Partition *isolates* the Conformant Relying Party from the Revocation-Information Source, capping observability at eventual-upon-recovery. *Forces:* FM1, INV12; S4.
- **reconstructs-from** — an Independent Reviewer *reconstructs-from* a Reconstruction Record alone, without the original verifier's runtime state and without privileged information. *Forces:* INV9, SO8, ER4, FR6.
- **issues-base (external)** — the existing workload-identity infrastructure *issues-base* Principal and Delegate identity; the system does **not** issue base identity. This is named as an *external* relationship to mark the boundary, not as a relationship the system performs. *Forces:* INV10, ER15, C1.
- **refuses-over-scope (guard)** — at the issuance boundary, an over-scoped request is *refused*; no Delegation is created. Named as the guard on the Issuance transition in §9.1. *Forces:* SO6, FM6, INV2, AP10.

## 9. State transitions (what state transitions exist)

Three labeled transition systems, expressed as named states and labeled transitions with guards. These are conceptual state machines about the domain, not runtime artifacts: no event queues, no timers, no concurrency primitives, no execution order is implied.

### 9.1 Delegation lifecycle

States: **NotIssued**, **Issued**, **Expired**, **Revoked**.

- **NotIssued → Issued** (*Issuance*): guard — the requested scope is a strict subset of the principal's Permission Set (M-INV2, SO6); an over-scoped request is Refused (FM6) and no Delegation is created. For an ephemeral delegate, this transition requires no long-lived statically-provisioned identity (ER10, FR7).
- **Issued → Expired** (*Expiry*): occurs at the delegation's Expiration time; monotone and terminal with respect to validity (M-INV3).
- **Issued → Revoked** (*Revocation*): occurs when a Revocation targeting this delegation instance takes effect; one-way and terminal for this instance (M-INV4, M-INV6).
- **Expired → Issued**: forbidden (M-INV3 terminality).
- **Revoked → Issued**: forbidden (M-INV4 terminality).

`Expired` and `Revoked` are terminal **for validity**. The Reconstruction Record persists in either state: reconstruction is independent of, and after, the original verification event, including after expiry or revocation (M-INV9, FR6). A freshly-issued delegation to the same logical grant as a revoked one is neither "the revoked instance resurrected" (forbidden by M-INV4) nor a model concept this RFC fixes; what individuates delegation instances is deferred (§7B, §17) — frozen FM5 explicitly calls this an open question.

### 9.2 Verification verdict (per presentation)

States: **NotPresented**, **Accept**, **Reject**, **Inconclusive** `[HYPOTHESIS]`.

- **NotPresented → Accept**: every required check passes — identity binding (M-INV1), integrity of the presented Scope (INV8 — the tamper-evidence check on the scope; the scope-subset property is enforced at issuance, M-INV2/SO6, and is not a verification-time check, the relying party need not hold the principal's full Permission Set), expiry (M-INV3), revocation-observability within R (FR4, SO1), signature/tamper (M-INV8). Offline, with no live call (M-INV7).
- **NotPresented → Reject**: any single required check fails — single-check rollback (SO5, FM11). Expired (M-INV3), Revoked-and-Observable (M-INV4, SO1), tampered record or tampered Scope (M-INV8) each force Reject. Over-scope (M-INV2) is refused at the Issuance transition (§9.1), not re-derived at verification; the verification-time scope check is integrity (INV8), not subset re-derivation, so an over-scope delegation that somehow reached presentation is caught by the integrity check, not by a subset re-derivation the RP may be unequipped to perform.
- **NotPresented → Inconclusive** `[HYPOTHESIS]`: a required check cannot be conclusively determined — signature unverifiable, required trust material unavailable, or clock beyond stated tolerance (FM3, FM9, ER3 clock tolerance).
- **Inconclusive → Reject** (*fail-closed*, `[HYPOTHESIS]`): under C-INV1, an inconclusive verification rejects rather than accepts. This transition is a hypothesis, not established (C-M-INV1, NFR3, ER11, SO4, `DEFERRED.md` D6).
- **Accept / Reject** are terminal for that presentation; a new presentation is a new Verification.

### 9.3 Revocation-observability (the relying party's view of a given Revocation)

States: **Not-yet-Observable**, **Observable**.

- **Not-yet-Observable → Observable (non-partitioned operation)**: the revocation propagates to the relying party within R of taking effect (SO1, FM2). The propagation *mechanism* is the spike-selected composition (S2/S3), deferred; the bound (≤ R in non-partitioned operation) is fixed here.
- **Partition ceiling (hard bound)**: when a Partition *isolates* the relying party from the Revocation-Information Source at the moment of the revocation, partition recovery — not R — is the hard ceiling. Observability is bounded by partition recovery (no earlier), per S4 / FM1 / M-INV12 — an information-theoretic limit. R is the target **within non-partitioned operation**, not a bound that overrides a standing partition. The model makes **no** claim of in-partition observability.
- **Observable** is terminal for that revocation (M-INV4 terminality of the underlying revocation).

## 10. Model invariants (what invariants constrain the model)

These are the conceptual constraints every compliant realization must respect. Each is a traced restatement of an established frozen invariant (`03_SYSTEM_INVARIANTS.md` INV1–INV12) or the candidate invariant C-INV1, restated at the conceptual level — **not new invariants** (Discipline §2 "trace, don't invent"). The parenthetical cites the frozen invariant each restates.

- **M-INV1 — Identity binding.** A Delegation that verifies successfully binds exactly one Principal and exactly one Delegate, both **deterministically** recoverable from the Delegation's content together with the relying party's trust material; no successful verification attributes a Delegation to identities other than those established at issuance. *(INV1.)*
- **M-INV2 — Scope subset.** A Delegation's Scope is a strict subset of the bound Principal's Permission Set, enforced at the Issuance transition; over-scope is refused. *(INV2.)*
- **M-INV3 — Expiry monotonicity.** Once a Delegation is Expired it never returns to Issued; a conformant verifier never treats an Expired delegation as valid. *(INV3.)*
- **M-INV4 — Revocation terminality.** Once a Delegation instance is Revoked it never returns to Issued. *(INV4.)*
- **M-INV5 — Revocation identity-independence.** A Revocation does not invalidate the Principal's or the Delegate's underlying identity. *(INV5.)*
- **M-INV6 — Revocation singularity.** A Revocation targets exactly one Delegation instance; revoking one Delegation is not, by itself, cause to reject a different, unrevoked Delegation. *(INV6.)*
- **M-INV7 — Offline verification.** Verification requires no network call to a shared authority at the moment of verification. *(INV7.)*
- **M-INV8 — Tamper-evidence.** A Reconstruction Record altered after creation never verifies as the original; tampering is always detectable by a Conformant Relying Party. *(INV8.)*
- **M-INV9 — Reconstruction self-sufficiency.** A Reconstruction Record is verifiable by a third party without access to the original verification event or the original verifier's runtime state. *(INV9.)*
- **M-INV10 — Base-identity boundary.** The model has no system-issued base-identity concept; the system operates on already-issued identity material. *(INV10.)*
- **M-INV11 — Companion, not a spec change.** The model requires no change to the existing workload-identity standard. *(INV11.)*
- **M-INV12 — Observability-claim bound.** The model's revocation-observability claims do not exceed what is achievable without a live call; in particular, no claim asserts in-partition observability. *(INV12.)*
- **C-M-INV1 — Fail-closed under inconclusive verification.** `[HYPOTHESIS]` An Inconclusive verification rejects rather than accepts. *(C-INV1, candidate — not established until V1, per `DEFERRED.md` D6.)*

Every frozen established invariant (INV1–INV12) and the candidate invariant (C-INV1) is restated. No model invariant strengthens a frozen invariant; each is the conceptual-level restatement of one. M-INV4 and M-INV6 apply to a delegation *instance*; what constitutes the instance-identity they keys to is not fixed by any frozen item and is deferred to a later mechanism RFC (§7B, §17).

## 11. Concepts explicitly out of scope

Each concept was considered and excluded, with a trace. The exclusions are listed so nothing is silently dropped.

- **Multi-hop Delegation / Delegation Chain (beyond one hop) and Per-hop Authorization.** Single-hop only in V1; the model covers exactly one delegator→delegate link. The word "chain" in `USE_CASE_CATALOG.md` UC5 and `FUNCTIONAL_REQUIREMENTS.md` FR6 is read as the single delegate link a reconstruction record carries, not a transitive chain. *Trace:* S5 (gate, out of V1); `01_ENGINEERING_REQUIREMENTS.md` "Requirements explicitly not added" (multi-hop); `DEFERRED.md`.
- **Three-or-more Trust Domains / Federation.** *Trace:* D3, ER17, C3.
- **Non-SPIFFE environments.** *Trace:* D4.
- **Within-Window Replay Resistance (as a model concept).** The model has no replay-resistance concept; a captured valid Delegation re-presented within its validity window is, in scope, a valid presentation (it would pass every check in §9.2 Accept). *Trace:* FM8; `02_SECURITY_OBJECTIVES.md` Non-objectives.
- **Issuer Key-Compromise Resistance.** The model has no key-rotation or compromise-recovery concept. *Trace:* FM5; Non-objectives.
- **Continuous Posture Re-attestation (validity tied to verification-time posture).** *Trace:* D2, FR10, ER14 (`[HYPOTHESIS]`).
- **Cross-Protocol Interoperability with a non-adopting Relying Party.** *Trace:* D1, FR9, ER13 (`[HYPOTHESIS]`).
- **Committed Latency SLA / Latency as a model concept.** The model has no latency concept. *Trace:* D5, NFR1, ER12 (`[HYPOTHESIS]`).
- **Fail-closed as an established/committed requirement.** The fail-closed posture is carried in the model as `[HYPOTHESIS]` **only** — the Inconclusive→Reject transition (§9.2) and C-M-INV1 (§10) — and as the honest label in §13; it is **not** established as a committed requirement. *Trace:* D6, NFR3, ER11, SO4, C-INV1 (`[HYPOTHESIS]`).
- **Delegation-instance identity.** What individuates two issuances to the same (principal, delegate, scope, expiration) is not fixed by any frozen item (FM5 explicitly calls it an open question). Carried as an open question (§17), not excluded outright and not modeled as a concept. *Trace:* FM5 (`04_FAILURE_MODEL.md`).
- **Base-Identity Issuance as a system concept.** Modeled as an external relationship only (`issues-base`, §8). *Trace:* INV10, ER15, C1.
- **Modification of the Existing Workload-Identity Standard.** *Trace:* INV11, ER16, C2.
- **Relying-Party Grant Decision.** The model proves the Delegation; it does not decide whether the scope should be granted. *Trace:* `SYSTEM_CONTEXT.md` "outside this boundary"; FM10.
- **Buyer / Commercial Packaging.** *Trace:* D7, C5.
- **Auditor Persona / Audit-Specific Tooling.** The Independent Reviewer concept covers record-reconstruction; the auditor-as-a-validated-user concept is a hypothesis. *Trace:* D8, `USER_MODEL.md` (auditor = hypothesis).
- **Trust Bridging Between Organizations With No Prior Federation.** *Trace:* D9, UC7.

## 12. Adversarial review

Each load-bearing conceptual choice, counter-held to its wrong state, and the observable consequence a reviewer would see. A choice whose wrongness produces no observable consequence fails DR4; all six below produce one.

1. **Scope modeled as a subset-relation concept enforced at issuance, plus an integrity check at verification (not a subset re-derivation).** Wrong state: scope is opaque/unchecked, or subset is re-derived at verification by an RP lacking the principal's Permission Set. Observable consequence: SO6 / AT4 — over-scope issuance would not be refused at the Issuance transition; or an RP would silently accept an over-scope delegation it cannot re-derive. M-INV2 binds the issuance guard; INV8 binds the verification-time integrity check; R1 silent privilege-escalation path closes.
2. **Revoked → Issued forbidden (revocation terminal).** Wrong state: the model allows Revoked → Issued. Observable consequence: INV4 / AT11 — a revoked delegation re-accepted; the monotonicity falsifies against INV4. The transition is forbidden (M-INV4).
3. **Inconclusive carried as a named verdict state.** Wrong state: only Accept/Reject exist. Observable consequence: C-INV1 / AT22 — an inconclusive condition forced to Accept is the canonical silent trust failure (FM11, R1); the model has no place to carry the fail-closed posture. Inconclusive `[HYPOTHESIS]` is kept as a named state.
4. **Partition modeled as an honest Not-yet-Observable bound (no in-partition observability claim).** Wrong state: the model claims a revocation is Observable to a partitioned RP. Observable consequence: INV12 / AT14 — an over-claim contradicted by the S4 information-theoretic limit. The model keeps the Partition as an honest bound, not a claimed guarantee.
5. **Base-identity issuance modeled as external-only.** Wrong state: the system issues base identity. Observable consequence: INV10 / AT27 — the system assumed a boundary it does not have; the companion-not-replacement invariant falsifies. Base-identity issuance is kept as an external relationship.
6. **Reconstruction Record modeled as a concept (tamper-evident, self-sufficient).** Wrong state: the model omits the record. Observable consequence: INV8 / INV9 / AT19–AT21 — no third-party reconstruction of delegator, delegate, scope, and time; AP8 / SO8 falsify. The record is kept.

## 13. Honest-claims statement

This model carries forward, unchanged in status, the following limits from the frozen package:

- **FM5 — issuer signing-key compromise:** the model has no concept warranting resistance; unmitigated within current scope. The "re-issue a revoked delegation" sub-case of FM5 is carried as an explicit open question (§17), not solved.
- **FM8 — within-window replay:** the model has no replay-resistance concept; unmitigated within current scope.
- **S4 — partition observability limit:** the model claims no in-partition revocation observability; an information-theoretic limit, not a design gap.
- **C-INV1 — fail-closed posture:** carried as `[HYPOTHESIS]` (the Inconclusive→Reject transition in §9.2, C-M-INV1 in §10); not established until V1.

No concept or invariant in this model claims resistance to FM5/FM8, observability beyond S4, or established fail-closed behavior (AP5, DR5).

## 14. Scope statement

This model asserts concepts and relationships only within the two-domain, SPIFFE-coexisting, single-hop scenario: domain A containing the principal, domain B containing the relying party, no shared authority, federation disabled, single delegator→delegate link. Three-or-more domains, non-SPIFFE environments, multi-hop delegation, and non-coexisting substrates are out of V1 (ER17, AP6, `DEFERRED.md` D3–D4, S5). A later RFC requiring behavior outside this scope must stop and declare itself out of V1 (DR6).

## 15. Hypothesis-preservation statement

This model references two hypothesis properties **as context, not as established behavior**: the fail-closed posture (C-INV1 / ER11 / SO4 — the Inconclusive→Reject transition in §9.2 and C-M-INV1 in §10) and the latency concept (NFR1 / ER12 `[HYPOTHESIS]` — absent from the model entirely). Both remain testable-but-unconfirmed until V1; neither is promoted to established (DR7, AP4). A later mechanism RFC that commits the system to a fail-closed behavior or to a latency target shall carry the `[HYPOTHESIS]` label forward.

## 16. Reviewability statement

The conceptual model is reviewable by direct reading against the frozen specification package (Phase 7 Product Definition + Phase 8 Engineering Requirements + the frozen `LEVEL0_1_FEASIBILITY_GATE.md`): does each concept trace to a forcing frozen item (§4, §7); does each model invariant restate an established frozen invariant without invention (§10); does each out-of-scope concept trace to a frozen deferral or non-objective (§11); does each wrong-state consequence in §12 falsify against a frozen invariant or failure mode. **No build is required** for this RFC, because it asserts no runtime mechanism. Later mechanism RFCs that assert runtime behavior will inherit the build-complemented form of SO8 (AP8, AP13). An independent reviewer needs the frozen package (and, for mechanism RFCs, a build) — and nothing else — to reproduce each conceptual claim's verdict.

## 17. Open questions

- **Revocation-Information Source mechanism.** Whether the propagation is push, pull, cached-pull, or brokered (S2/S3 — frozen candidate categories, named only to declare deferral) is deferred to a later RFC; here the source is named as a concept only.
- **Delegation-instance identity.** The model binds a Delegation by (Principal, Delegate, Scope, Expiration) but does not name what individuates two issuances to the same P+D+scope+expiration. Frozen FM5 explicitly flags the "re-issue a revoked delegation" sub-case as "a distinct question the Product Definition does not address" (`04_FAILURE_MODEL.md`). Consequently M-INV4 (revocation terminality) and M-INV6 (revocation targets exactly one) apply to a delegation *instance*; what constitutes the instance-identity that a revocation keys to is deferred to a later mechanism RFC. The model makes no claim here and carries the gap honestly (§13).
- **Internal decomposition.** The decomposition of the Delegation / Reconstruction Record / Verification concepts into modules is deferred to a later module-boundary RFC (AP12).
- **Multi-hop.** Whether a founder scope act returns multi-hop delegation to scope (S5) — out of V1; not modeled. If it returns, a "Delegation Chain" concept would enter the model.
- **Latency.** Whether latency becomes a model concept at V1 (D5, ER12 `[HYP]`) — not modeled here.
- **Governance envelope (non-frozen).** RFC-state advancement and freeze authorization are governed by RFC-000's Architecture Review Process (an unfrozen draft) and `agents/agents/GOVERNANCE.md` (an unfrozen framework); the founder's scope-act authority (resolving S1–S5) is the only such role supported by a frozen doc (`LEVEL0_1_FEASIBILITY_GATE.md`). Neither envelope is a problem-domain concept and neither is modeled in §7; they are referenced here for completeness.

## 18. Provenance

- **Primary source:** the frozen Phase 7 Product Definition package and the frozen Phase 8 Engineering Requirements set (26 documents, hash-pinned in `FROZEN.sha256`, frozen 2026-07-04 per `agents/agents/journal/2026-07-04-freeze-phase-7-phase-8.md`); frozen `LEVEL0_1_FEASIBILITY_GATE.md` for the scope parameters S1–S5, the spike outcomes α/β/γ/δ, and the founder scope-act role; RFC-000 (principles + Architecture Review Process envelope) and RFC-001 (System Context — the boundary this model is the inside of), both unfrozen drafts.
- **Secondary (evidence, not concept source):** `agents/agents/GOVERNANCE.md` for the honesty / confidence discipline (an unfrozen framework, referenced not as a concept source).
- **Confidence:** High that the concepts, relationships, state transitions, and model invariants are traced restatements of frozen FR/ER/INV/SO/FM/C items, with no invention (every concept in §7 carries a forcing trace in §4; every model invariant in §10 cites the frozen invariant it restates). High that the out-of-scope concept list in §11 covers every frozen deferral (D1–D9) and every Security-Non-objective. Medium that the single-hop reading of "delegation chain" (`USE_CASE_CATALOG.md` UC5; `FUNCTIONAL_REQUIREMENTS.md` FR6) is the right conceptualization — the model covers exactly one delegator→delegate link; what would change it is a founder scope act returning multi-hop delegation to scope (S5). Medium that delegating instance-identity to a later mechanism RFC (rather than modeling it here) is correct — frozen FM5 explicitly declines to address it; what would change it is a frozen-package amendment that fixes instance-identity.
- **Change-condition:** only an amendment to the frozen Product Definition or frozen Phase 8 package (per `CONTRIBUTING.md` §4 amend procedure: journal entry → dated change note → `make frozen-baseline` → commit `FROZEN.sha256`) authorizes adding, removing, or strengthening a concept or model invariant. A founder scope act resolving S1–S5 adjusts which concepts acquire a realized mechanism; it does not add concepts the frozen package does not force. Resolving the FM5 instance-identity open question requires a frozen-package amendment (it is a Product Definition gap, not a scope-parameter resolution).
- **Freeze status of RFC-002:** not yet frozen. RFC-002 declares itself *eligible* for the frozen set (DR10) but is frozen only by a future, explicitly authorized act — not by this founding. Until then it is an uncommitted, unfrozen document, presented for founder review.

<!-- checkpoint: governance(CI-testing-gates): improve CI testing gates -->

<!-- checkpoint: context(system-boundary-definition): update system boundary definition -->

<!-- checkpoint: docs(deployment-manual): finalize deployment manual (#42) -->

<!-- checkpoint: feat(sdk): implement attenuation rule engine -->

<!-- checkpoint: refactor(record): refactor key derivation -->
