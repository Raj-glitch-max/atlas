# RFC-001 — System Context

## 1. Header

| | |
|---|---|
| RFC | 001 |
| Title | System Context |
| Status | Draft |
| Author | founder |
| Date | 2026-07-04 |
| Supersedes | none |
| Superseded by | none |

This is the architectural context for all later architecture work. It defines the system's boundary, the external actors and systems it relates to, the trust boundaries across which its guarantees hold, the responsibilities it does and does not accept, and the high-level interactions in which it participates. It does not define repository structure, technologies, APIs, protocols, storage, deployment, or implementation — those are deferred to later RFCs (§3).

## 2. Source authority

- Frozen Phase 7 `SYSTEM_CONTEXT.md` (Product Definition package, frozen 2026-07-04) — the canonical system-context statement this RFC is the architectural restatement of.
- Frozen Phase 7 `PRODUCT_DEFINITION.md`, `USER_MODEL.md`, `USE_CASE_CATALOG.md`, `CONSTRAINTS.md`, `V1_SCOPE.md`, `DEFERRED.md` — boundary, actor, scope, and deferral source.
- Frozen Phase 8 — `01_ENGINEERING_REQUIREMENTS.md` (ER1–ER17), `02_SECURITY_OBJECTIVES.md` (SO1–SO8), `03_SYSTEM_INVARIANTS.md` (INV1–INV12, C-INV1 candidate), `04_FAILURE_MODEL.md` (FM1–FM11), `05_ACCEPTANCE_TEST_PLAN.md` (AT1–AT30) — the boundary-fixing requirements, objectives, invariants, and failure modes.
- Founder instruction, 2026-07-04, authorizing RFC-001 as exactly one new RFC, defining only the six content items in §3 and no others.

## 3. Scope of this RFC

This RFC defines:
- the **system boundary** — what is inside the system and what is outside it;
- the **external actors** the system relates to at runtime, and the governance actors it relates to off-runtime;
- the **external systems** the system depends on or coexists with;
- the **trust boundaries** across which the system's guarantees are defined;
- the **responsibilities** the system accepts and the responsibilities it explicitly declines;
- the **high-level interactions** (boundary-level flows) in which the system participates.

This RFC deliberately does **not** define (each deferred to a later, separately authorized RFC):
- **Repository structure.** No module decomposition, package layout, or source-tree shape.
- **Technologies.** No language, runtime, library, crypto primitive, or framework.
- **APIs.** No interface signatures, request/response shapes, or call contracts.
- **Protocols.** No wire format, transport, message sequencing, or encoding.
- **Storage.** No persistence substrate, data layout, or durability mechanism.
- **Deployment.** No topology, runtime placement, or deployment arrangement.
- **Implementation.** No code, no internal data structures, no algorithm specifics.

## 4. Traces

| Architectural decision | Forced by |
|---|---|
| The system operates on already-issued workload identity, never issues it | ER15, INV10, C1, AP2 |
| The system does not modify the existing workload-identity standard | ER16, INV11, C2, AP2 |
| Guarantees are defined at the conformant-verifier boundary only | FM10, SO2, ER7 |
| No live call to a shared authority is required to verify | ER7, INV7, SO2, AP1 |
| The system is bounded to exactly two trust domains, coexisting with the existing substrate | ER8, ER17, C3, AP6 |
| Specific revocation takes effect without workload restart | ER5, ER6, INV4, INV5, INV6, SO1 |
| Issuance scope is a subset of the principal's permissions, enforced at issuance | ER1, ER2, INV1, INV2, SO6, AP10 |
| Over-scoped issuance is refused at the boundary | FM6, AP10 |
| The recorded package is tamper-evident and sufficient for reconstruction | ER4, INV8, INV9, SO3, AP3 |
| Revocation observability is bounded by R and by the S4 partition limit | FM1, FM2, FM4, INV12, AP7 |
| Unmitigated-within-current-scope failures are carried honestly, not solved | FM5, FM8, INV12, AP5, AP8 |
| The system declines control of non-conformant / malicious relying parties | FM10 |
| The system declines to decide whether the relying party should grant | `SYSTEM_CONTEXT.md`, FM10 |

## 5. RFC-000 principle compliance

- **AP1 (offline):** the system's verification interaction (§12.3) has no live call to a shared authority.
- **AP2 (companion):** the system is downstream of, never a replacement for, base identity (§7, §9).
- **AP3 (tamper-evidence):** the system boundary includes a tamper-evident, reconstruction-sufficient record (§7, §12.1).
- **AP4 (fail-closed, `[HYPOTHESIS]`):** the system's behavior in the S4-inconclusive state is a hypothesis carried forward (§16), not a present commitment.
- **AP5 (honesty):** §14 carries the FM5/FM8/S4 limits into the system's claims.
- **AP6 (two-domain):** §10 trust boundaries; §15 scope statement.
- **AP7 (parametric):** revocation observability is stated for any parameter value of R and S4, with no fixed mechanism (§12.4).
- **AP8 (independent reviewability):** §17.
- **AP9 / AP13 (adversarial, observable):** §13 names the observable wrong-state consequence of each load-bearing boundary choice.
- **AP10 (no over-issuance):** §7 boundary; §12.1 issuance interaction.
- **AP11 (minimize accidental complexity):** this RFC defines boundary, not mechanism; no component is introduced beyond what a frozen requirement forces.
- **AP12 (stable module boundaries):** this RFC places the externally-exposed surface (the record-facing surface) as stable; mechanism volatility is isolated into later RFCs.
- **TP1–TP6:** no trade-off arises within a boundary-only RFC; the philosophy binds later mechanism RFCs.
- **DR1–DR10:** §4 traces (DR1); §3 deferral (DR9); §13 adversarial (DR4); §14 honest claims (DR5); §15 scope (DR6); §16 hypothesis (DR7); §17 reviewability (DR8); DR2/DR3 hold by the absence of any architecture-of-architecture or spike-coupling in this RFC; DR10 freeze-status in §19.

## 6. Spike-outcome analysis

This RFC is spike-outcome-independent by construction: it describes the *boundary* of the system, not the *mechanism* by which the boundary is realized. The scope parameters S1–S5 — R (revocation-observability latency), S2/S3 (cached-pull admissibility / broker definition), S4 (partition observability), S5 (per-hop authorization, out of V1) — concern *how* a behavior is realized inside the boundary, not *whether* the behavior is inside it.

Therefore this RFC supports all four spike outcomes without modification: α (composition-works — the boundary holds), β (technology-gap — the boundary holds; the gap is a mechanism problem for a later RFC), γ (logical-impossibility — already established at S4 — is carried as an honest limit, §14), and δ (unresolvable — the boundary holds; an unresolvable parameter is bounded, not solved, here). This RFC precludes nothing and defers every mechanism question.

## 7. System boundary

### Inside the system boundary
- **Delegation issuance** — binding an already-issued principal identity to a delegate identity, with scope ⊆ the principal's permissions and an expiration, refusing over-scope requests.
- **Delegation record production** — tamper-evident, reconstruction-sufficient records.
- **Verification logic at the conformant verifier** — the boundary at which the system's guarantees are defined.
- **Specific revocation** — revocation of an individual delegation, observable to the conformant relying party per the spike-selected composition, bounded by R and S4.
- **Honest degradation** — the system's behavior when the S4 partition makes a verdict inconclusive is part of the system (and is a hypothesis, §16).

### Outside the system boundary
- Base workload-identity issuance (C1, INV10, ER15) — the external identity infrastructure's responsibility.
- Modification of the existing workload-identity standard (C2, INV11, ER16).
- The relying party's grant decision (`SYSTEM_CONTEXT.md`, FM10).
- Control of non-conformant or malicious relying parties (FM10).
- Resistance to issuer key compromise (FM5) — out of warranted scope, carried honestly (§14).
- Resistance to within-window replay (FM8) — out of warranted scope, carried honestly (§14).
- Three-or-more-domain federation and non-coexisting substrates (ER17, AP6, DEFERRED D3–D4).

## 8. External actors

Runtime actors:
- **Principal** — the workload identity on whose behalf a delegate will act. Holds a permission set against which a delegation's scope is subset-checked. Originates delegation requests spanning a subset of its permissions.
- **Delegate** — the workload (or agent) that receives a delegation and presents it to a relying party. Exercises no authority beyond what the delegation grants.
- **Relying Party (RP)** — the verifier that examines a presented delegation and decides whether to grant. The system's guarantees are defined at the **conformant RP** boundary; a non-conformant RP is outside the trust boundary (§10.1, §13.3).

Off-runtime (governance) actors:
- **Founder / scope-setter** — resolves the scope parameters S1–S5; advances RFC states; authorizes freeze. Not a runtime participant.
- **Independent reviewer** — reproduces a verdict from the specification package (and, for mechanism RFCs, a build), per SO8 / AP8.

## 9. External systems

- **Existing workload-identity infrastructure** — issues the base workload identity the system consumes. The system is strictly downstream; it never replaces this infrastructure (C1, INV10, ER15) and coexists with it, extending it with delegation (C2, INV11, ER16). Resolves AP2.
- **The relying party's existing identity-verification baseline** — delegation-verification operates *alongside* this baseline, not as a replacement for it (NFR5, ER9, SO7, AP2). Whether the two are one combined mechanism or two coexisting ones is a mechanism question for a later RFC; here it is stated only as a coexistence requirement.
- **Revocation-information source** — wherever revocation decisions originate and are made observable to the RP. The relationship between the RP and this source — push, pull, cached-pull, or brokered (the S2/S3 range the spike examines) — is determined by the spike-selected composition (AP7, S2/S3, S4) and is deferred to a later RFC. This RFC names the source as an external-system concept, not its mechanism.

## 10. Trust boundaries

1. **The conformant-verifier boundary.** The line at which the system's guarantees are defined (FM10, SO2, ER7). Inside: conformant verifiers exercising the documented checks. Outside: non-conformant or malicious RPs; the system warrants nothing there.
2. **The two-domain boundary.** Domain A (containing the principal) and domain B (containing the relying party), no shared authority between them, federation disabled (C3, ER8, ER17, AP6). Delegation crosses this boundary.
3. **The partition boundary (S4).** Between the RP and the revocation-information source. A partition is an information-theoretic limit on observability: revocation performed during the partition is not observable to the RP before recovery (FM1, FM2, INV12, AP7). The system does not claim in-partition observability; this is an honest limit, not a gap.
4. **The issuance boundary.** Where scope-subset is enforced (SO6, ER1, ER2, AP10). An over-scoped request is refused here (FM6); issuance refusal is the contract (AP10).

## 11. Responsibilities

### The system accepts
- Issue delegations bound to already-issued identity, scope ⊆ the principal's permissions, with expiry (ER1, ER2, ER3, INV1, INV2, INV3, SO6, AP10).
- Produce tamper-evident, reconstruction-sufficient records (ER4, INV8, INV9, SO3, AP3).
- Enable offline verification by the conformant RP (ER7, INV7, SO2, AP1).
- Enable cross-domain (A→B) verification bounded to two domains (ER8, ER17, AP6).
- Allow specific revocation taking effect without workload restart (ER5, ER6, INV4, INV5, INV6, SO1).
- Allow issuance to ephemeral workloads (ER10).
- Degrade honestly under the S4 partition bound (FM1, INV12, AP5; the fail-closed posture itself is `[HYPOTHESIS]`, §16).
- Remain reconstructable by an independent reviewer (ER4, INV9, SO8, AP8).

### The system declines
- Issuing base workload identity (C1, ER15, INV10, AP2).
- Modifying the existing workload-identity standard (C2, ER16, INV11, AP2).
- Deciding whether the RP should grant the scope (`SYSTEM_CONTEXT.md`, FM10).
- Controlling non-conformant or malicious RPs (FM10).
- Resisting issuer key compromise (FM5) — unmitigated-within-current-scope.
- Resisting within-window replay (FM8) — unmitigated-within-current-scope.
- Three-or-more-domain or non-coexisting behavior (ER17, AP6, DEFERRED D3–D4).

### Actor responsibilities
- **Principal** — holds a permission set; originates delegation requests spanning a subset of its permissions; does not request scope it does not hold (else refused at the issuance boundary, §10.4).
- **Delegate** — receives and presents the delegation; exercises no authority beyond the delegation's scope and validity window.
- **Conformant RP** — verifies the delegation offline using locally-held trust material; makes the grant decision; degrades per the documented checks under the S4 partition.
- **Non-conformant RP** — outside the trust boundary; its behavior is its own responsibility (FM10).
- **Existing identity infrastructure** — supplies the base identity on which the system operates; the system never modifies or replaces it.
- **Independent reviewer** — reproduces verdicts from the package (and, for mechanism RFCs, a build) per SO8 / AP8.

## 12. High-level interactions

Boundary-level flows in which the system participates — no protocol, no API, no message format.

1. **Issuance.** A principal (or a delegate on its behalf) requests a delegation. The system issues it bound to already-issued base identity, with scope ⊆ the principal's permissions and an expiration; an over-scoped request is refused at the issuance boundary (§10.4, AP10/SO6/FM6/ER1/ER2). A tamper-evident, reconstruction-sufficient record is produced (AP3/ER4/INV8/INV9). Ephemeral workloads may receive delegations (ER10).
2. **Presentation.** The delegate presents the delegation to the relying party in domain B, crossing the two-domain trust boundary (§10.2, AP6/ER8).
3. **Verification.** The RP — a conformant verifier — determines validity, scope, expiry, and revocation status from the presented material and locally-held trust material, with no live call to a shared authority (§10.1, AP1/ER7/INV7/SO2). Behavior in the inconclusive (partition) state is a `[HYPOTHESIS]` (§16), carried forward, not promoted.
4. **Revocation.** A specific delegation is revoked. The revocation becomes observable to the RP per the spike-selected composition, bounded by R and S4; thereafter the RP rejects the revoked delegation (AP7/ER5/SO1/INV4/INV5/INV6/FM1/FM2/FM4). The underlying principal and delegate identities remain valid (INV5). No workload restart is required for the new revocation state to take effect (ER6).
5. **Reconstruction.** A third party verifies the record independent of, and after, the original verification event, with no access to the original verifier's runtime state (AP3/ER4/INV9/SO3/SO8/AP8).
6. **Partition.** At any moment, the RP may be partitioned from the revocation-information source (§10.3). The system's guarantees degrade honestly within the S4 bound — the system rejects rather than over-claims; it does not assert revocation observability it cannot have (AP7/FM1/INV12/AP5). The fail-closed posture in this state is `[HYPOTHESIS]` (§16).

## 13. Adversarial review

Each load-bearing boundary choice, counter-held to its wrong state, and the observable consequence a reviewer would see:

1. **Boundary drawn inside base-identity issuance.** Wrong state: the system attempts to issue base identity. Observable consequence: the system violates C1/INV10/ER15 — a reviewer sees the system claiming an authority it does not have, and the companion-not-replacement invariant (INV10/INV11) fails its acceptance test. The boundary is therefore drawn outside issuance (§7).
2. **Boundary drawn to include the relying party's grant decision.** Wrong state: the system decides whether the RP grants. Observable consequence: the system warrants a grant it cannot enforce (FM10) — a reviewer sees a guarantee with no enforcement point, and the honest-claims check (DR5) fails. The boundary is drawn to exclude the grant decision; the system proves the delegation, the RP decides the grant (§7, §11).
3. **Boundary drawn to include non-conformant / malicious RPs.** Wrong state: the system warrants behavior at a non-conformant verifier. Observable consequence: the guarantee is unfalsifiable against a hostile RP (FM10) — a reviewer sees a claim that cannot be tested at the asserted boundary, and observability (AP13) fails. The boundary is drawn to exclude non-conformant RPs; the conformant-verifier boundary is where guarantees are defined (§10.1).
4. **Boundary drawn outside the S4 partition limit (claiming in-partition revocation observability).** Wrong state: the system claims to observe revocation performed during a partition. Observable consequence: the claim is information-theoretically false (FM1, INV12) — a reviewer sees a guarantee contradicted by the established S4 limit, and AP5/DR5 fail. The boundary is drawn to make the partition an explicit honest limit, not a claimed guarantee (§10.3, §14).
5. **Boundary omits the reconstruction record.** Wrong state: the system produces no tamper-evident reconstruction-sufficient record. Observable consequence: a third party cannot verify after the fact (ER4/INV9/SO8 fail) — a reviewer sees a verdict that cannot be reconstructed, and AP3/AP8/AP13 fail. The boundary is drawn to include the record (§7).

A choice whose wrongness produces no observable consequence fails DR4 and is rejected. All five choices above produce an observable, falsifiable consequence; none is dark.

## 14. Honest-claims statement

This RFC carries forward, unchanged in status, the following limits from the frozen package:
- **FM5 — issuer key compromise:** within current scope, the system does not warrant resistance to compromise of the issuing key. A later scope act may address it; until then, claims that assume resistance are out of scope.
- **FM8 — within-window replay:** within current scope, the system does not warrant against replay of a valid delegation within its validity window. A later scope act may address it; until then, claims that assume replay-resistance are out of scope.
- **S4 — partition observability limit:** the system does not warrant revocation observability for revocations performed during a partition of the RP from the revocation source. This is an information-theoretic limit (FM1), not a gap to be closed by a better mechanism.

The system's claims are bounded to conformant verifiers operating within these limits. No claim in this RFC exceeds them (AP5, DR5).

## 15. Scope statement

This RFC asserts behavior only within the **two-domain, coexisting scenario**: domain A containing the principal, domain B containing the relying party, no shared authority between them, federation disabled, the system coexisting with (not replacing) the existing workload-identity substrate. Three-or-more-domain federation, non-coexisting identity substrates, and a shared-authority (single-domain) arrangement are out of V1 (ER17, AP6, DEFERRED D3–D4). Should a later RFC require behavior outside this scope, it must stop and declare itself out of V1, per DR6.

## 16. Hypothesis-preservation statement

This RFC introduces no new hypothesis commitments. It references one hypothesis property — fail-closed behavior under the S4 inconclusive state (AP4, ER11 `[HYPOTHESIS]`, C-INV1 candidate) — *as context*, not as an established behavior of this system. The fail-closed posture is a hypothesis inherited from the frozen package; it remains testable-but-unconfirmed until V1, and this RFC does not promote it (DR7). A later RFC that commits the system to a fail-closed mechanism shall carry the `[HYPOTHESIS]` label forward.

## 17. Reviewability statement

The boundary claims in this RFC are reviewable by direct reading against the frozen specification package (Phase 7 Product Definition + Phase 8 Engineering Requirements, §2): does each boundary decision trace to a frozen item (§4); does each wrong-state consequence falsify against a frozen invariant or failure mode (§13); do the deferrals in §3 cover every forbidden category. No build is required for this RFC, because it asserts no runtime mechanism. Later mechanism RFCs that assert runtime behavior will inherit the build-complemented form of SO8 (AP8, AP13). An independent reviewer needs the frozen package (and, for mechanism RFCs, a build) and nothing else to reproduce each claim's verdict.

## 18. Open questions

- **Revocation interaction mechanism.** Whether the revocation-observability interaction is push, pull, cached-pull, or brokered (S2/S3) is not decided here; it is the spike-selected composition, deferred to a later RFC. This RFC names the *interaction* (revocation becomes observable to the RP bounded by R and S4), not the mechanism.
- **Per-hop authorization (S5).** Out of V1; not modeled here; a later scope act or a V2 RFC may include it.
- **The system's internal decomposition into components.** Out of scope for this RFC; a later module-boundary RFC (under AP12) will place the stable record-facing surface and the volatile spike-facing surface.
- **Whether the founder scope act adds three-or-more-domain federation.** Out of V1; not modeled here.

## 19. Provenance

- **Primary source:** frozen Phase 7 `SYSTEM_CONTEXT.md` and the frozen Phase 8 Engineering Requirements package (26 documents, hash-pinned in `FROZEN.sha256`, frozen 2026-07-04 per `agents/agents/journal/2026-07-04-freeze-phase-7-phase-8.md`); RFC-000 (the architectural principles this RFC complies with in §5, including the Architecture Review Process structure this RFC follows).
- **Secondary (evidence, not principle source):** `LEVEL0_1_FEASIBILITY_GATE.md` for the scope parameters S1–S5 and the spike outcomes α/β/γ/δ; `PRODUCT_DEFINITION.md` line 3 for the package scope.
- **Confidence:** High that the boundary, actors, external systems, trust boundaries, responsibilities, and interactions stated here are traced restatements of frozen `SYSTEM_CONTEXT.md` and the frozen Phase 8 requirements, objectives, invariants, and failure modes. High that the exclusions (no repository structure, technologies, APIs, protocols, storage, deployment, or implementation) are honored — each forbidden category is named in §3 and absent from §7–§12. Medium that the revocation-interaction description (§12.4) is the right level of abstraction without choosing a mechanism — the description inherits S2/S3/S4 from the frozen gate and adds no commitment; what would change it is a founder scope act resolving S2/S3.
- **Change-condition:** only an amendment to frozen `SYSTEM_CONTEXT.md` or to the frozen Phase 8 package (per `CONTRIBUTING.md` §4 amend procedure: journal entry → dated change note → `make frozen-baseline` → commit `FROZEN.sha256`) authorizes changing the boundary, actors, trust boundaries, or responsibilities stated here. A founder scope act resolving S1–S5 adjusts the *mechanism* a later RFC may choose, not this RFC's boundary.
- **Freeze status of RFC-001:** not yet frozen. RFC-001 declares itself *eligible* for the frozen set (DR10) but is frozen only by a future, explicitly authorized act — not by this founding. Until then it is an uncommitted, unfrozen document on disk, presented for founder review.

<!-- checkpoint: rfc(problem-scoring-guidelines): audit problem scoring guidelines (#43) -->

<!-- checkpoint: feat(stores): add panic handling middleware (#107) -->
