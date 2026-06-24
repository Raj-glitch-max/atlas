# Failure Model

**Phase:** 8 — Engineering Requirements.
**Source authority:** the frozen Phase 7 Product Definition package (especially `ASSUMPTIONS_AND_RISKS.md` R1–R6; `FUNCTIONAL_REQUIREMENTS.md` FR4, FR5; `NON_FUNCTIONAL_REQUIREMENTS.md` NFR2, NFR3, NFR6; `CONSTRAINTS.md` C1, C6). `LEVEL0_1_FEASIBILITY_GATE.md` supplies the S1–S5 scope parameters and the S4 impossibility; the gate's candidate compositions are referenced only as the spike's space of options, never as choices.
**What this document is not:** no architecture, technology, repository, API, or protocol. Every failure mode assumes **distributed-systems** conditions (partition, replication, partial failure, skew) **and** **adversarial** conditions (forgery, over-issuance, tampering, rollback attempts). Failure responses are stated as requirements the architecture must satisfy, parameterized by S1–S5 where the founder scope act has not yet fixed a value.

## Discipline

1. **Distributed + adversarial by construction.** Each Failure Mode (FM) names both a distributed-system trigger (network/timing/replication) and the adversarial exploitation it enables or the silent path it opens.
2. **Trace, don't invent.** Each FM's required response traces to a Product Definition requirement, an engineering requirement (`01_ENGINEERING_REQUIREMENTS.md` ER), an invariant (`03_SYSTEM_INVARIANTS.md` INV), or a security objective (`02_SECURITY_OBJECTIVES.md` SO).
3. **Spike-outcome-independent.** The C4 spike outcome is unknown. FMs are written to hold whichever composition the spike selects (or none); where a failure mode is composition-dependent, the requirement is stated as a bound the chosen composition must satisfy, not as a chosen mechanism.
4. **Honest about limits.** Where the Product Definition warrants no response (e.g., issuer key compromise, within-window replay), the FM records the mode as **unmitigated within current scope** and the system is forbidden from claiming resistance it does not have (per project doctrine: confidence without evidence is forbidden; `02_SECURITY_OBJECTIVES.md` Non-objectives).

## Scope parameters (carried from `01_ENGINEERING_REQUIREMENTS.md`)

- **R (S1):** revocation-observability latency; a revoked delegation must be rejected within R of the revocation becoming observable, R < delegation TTL.
- **S2/S3:** cached-pull admissibility / broker definition — whether a periodic signed-artifact pull or a passive cache is a forbidden "live call to a shared authority." FMs touching a live call use the strict reading; relaxation narrows the offline FMs, never strengthens them.
- **S4:** partition reading. A revocation performed while the RP is partitioned from the issuer is **not** observable to that RP before partition recovery — an information-theoretic limit (no live call ⇒ no fresh information). FMs honor this; the system is forbidden claims to the contrary.
- **S5:** per-hop authorization — out of V1 (single-hop). Not parameterized here.

## Failure modes

**FM1 — Revocation performed during a partition isolating the relying party from the issuer.**
*Trigger (distributed):* network partition separates the RP from the issuer/revocation source precisely at the moment a delegation is revoked, and persists into the RP's subsequent verification.
*Adversarial angle:* an attacker who can induce or extend the partition can keep a revoked delegation accepted for the partition's duration.
*Required response:* per the founder's S4 resolution — either (a) bound revocation observability to "eventual, within P seconds of partition recovery," or (b) explicitly exclude observability of revocations performed during such a partition. In either resolution, the system shall not accept a revoked delegation beyond the S4-defined bound, and shall not claim in-partition observability it cannot have. Under `[HYPOTHESIS]` NFR3 (ER11/SO4), once the bound cannot be satisfied the system rejects rather than accepts.
*Trace:* C6, FR5, NFR3[HYP], INV12; gate S4. *Spike-dependence:* none — S4 is an information-theoretic limit any composition must respect.

**FM2 — Revocation observability latency exceeds R (S1).**
*Trigger (distributed):* revocation propagation to the RP takes longer than R, so a revoked delegation is accepted in the window (revocation-effect, revocation-effect + R).
*Adversarial angle:* an attacker holding the soon-to-be-revoked delegation races verification against the propagation window.
*Required response:* a conformant verifier shall reject any delegation whose revocation is observable-and-effective for longer than R. Whichever revocation composition the C4 spike selects, its staleness floor must be ≤ R; if no candidate composition can meet R, that is a spike outcome the failure model must surface (it corresponds to gate outcome β — technology gap — or δ).
*Trace:* FR4, ER5, SO1, INV12; R (S1). *Spike-dependence:* the *mechanism* is spike-resolved; the *bound* (≤ R) is a hard requirement here.

**FM3 — Clock skew between issuer and verifier.**
*Trigger (distributed):* the verifier's notion of "now" differs from the issuer's, causing premature acceptance of a near-expiry delegation or premature rejection of a not-yet-expired one.
*Adversarial angle:* an attacker manipulates timing to push a delegation's apparent expiry past its real expiry (or vice versa).
*Required response:* the verifier's clock-skew tolerance relative to the issuer shall be explicit and bounded (ER3); within the tolerance the expired/not-expired verdict is deterministic; beyond the tolerance the system rejects (fail-closed, ER11/SO4 `[HYP]`).
*Trace:* FR3, ER3, NFR3[HYP]; gate U4 (the latency threshold is a budget, not the falsifier — but clock skew directly affects the expiry falsifier and so is load-bearing here).
*Spike-dependence:* none.

**FM4 — Stale revocation data (composition-dependent staleness floor).**
*Trigger (distributed):* the verifier holds stale revocation information and accepts a delegation whose revocation is already in effect at the (composition's) source but has not yet reached the verifier.
*Adversarial angle:* race the staleness window as in FM2; the distinction is that FM4 is the *composition's own* staleness floor (e.g., a periodic signed-artifact refresh interval), whereas FM2 is the latency ceiling R. The composition's staleness floor must not exceed R.
*Required response:* the chosen revocation composition's staleness is a hard bound, not a target: the verifier shall not accept a delegation whose revocation is observable-and-effective longer than R. If the candidate composition's design staleness exceeds R for the R the founder sets, the architecture is non-compliant and the spike must surface this (gate outcome β/δ).
*Trace:* FR4, ER5, SO1, R (S1). *Spike-dependence:* high — the failure exists *whichever* composition is chosen, but its exact floor is the spike's to characterize. The requirement binds the floor to R.

**FM5 — Issuer signing-key compromise (unmitigated within current scope).**
*Trigger (distributed/adversarial):* an attacker compromises the issuer's signing key. Two sub-cases:
  - *Forge new delegations:* the attacker mints delegations that pass signature verification by a conformant verifier.
  - *Re-issue a revoked delegation:* the attacker re-signs a delegation that was revoked, producing a "new" valid-looking unit (revocation is keyed to the delegation identity, not the signing event — per INV4 the original revoked delegation stays revoked, but a freshly-signed claim to the same logical grant is a distinct question the Product Definition does not address).
*Adversarial angle:* full issuer compromise is the strongest adversary in this model.
*Required response:* **none warranted by the Product Definition.** Tamper-evidence (NFR6/INV8/SO3) ensures *existing* records altered after creation are detected, but it does not protect against *newly forged* records signed with the compromised key. The system is therefore forbidden from claiming key-compromise resistance it does not have; `05_ACCEPTANCE_TEST_PLAN.md` shall not assert an acceptance test for a property the Product Definition does not warrant. This is recorded here (and in `02_SECURITY_OBJECTIVES.md` Non-objectives) so the gap is visible, not silently dropped.
*Trace:* R6 (no external validation / self-certification risk), NFR6, INV8. *Spike-dependence:* none; this is a product-scope gap, not a spike question.

**FM6 — Scope-escalation attempt at issuance.**
*Trigger (distributed):* a principal (or an attacker acting through one) requests a delegation whose scope is not a subset of the principal's own permissions.
*Adversarial angle:* privilege escalation through delegation.
*Required response:* issuance shall refuse; no over-scoped delegation shall be created (FR2/ER2/INV2/SO6, refusal fraction = 1).
*Trace:* FR2, ER2, INV2, SO6. *Spike-dependence:* none — this is an issuance-side property independent of the revocation composition.

**FM7 — Tampered delegation record (in-flight or at-rest).**
*Trigger (distributed):* an attacker alters a delegation record between creation and verification (in transit, in storage, or in the presented channel).
*Adversarial angle:* forge or downgrade a delegation by mutation.
*Required response:* the alteration is detected and the record is rejected; an altered record never verifies as the original (NFR6/ER4/INV8/SO3, detection-and-rejection fraction = 1 over the mutation set).
*Trace:* NFR6, FR6, ER4, INV8, SO3. *Spike-dependence:* none.

**FM8 — Replay of a captured delegation (within the validity window).**
*Trigger (distributed):* an attacker captures a valid delegation and re-presents it — to the same verifier or a different one — within its validity window.
*Adversarial angle:* reuse of legitimately-issued authority.
*Required response:*
  - *Post-expiry replay:* rejected by FR3/INV3 (the captured unit is expired) — fully mitigated by the Product Definition.
  - *Within-window replay:* **unmitigated within current scope.** The Product Definition specifies no replay-resistance within the validity window (no nonce, no audience binding, no single-use contract is traceable to FR1–FR10). The system is forbidden from claiming within-window replay resistance it does not have. Recording it here and in `02_SECURITY_OBJECTIVES.md` Non-objectives keeps the gap visible.
*Trace:* FR3 (post-expiry), INV3; non-objective (within-window). *Spike-dependence:* none; within-window replay is a product-scope gap. (The gate's H1 mentions "reject replay" — that is a different artifact; the Product Definition does not place it in scope, so it is not asserted here.)

**FM9 — Relying party lacks fresh trust material for signature verification.**
*Trigger (distributed):* the RP's locally-held trust material is stale or unavailable at verification time.
*Adversarial angle:* an attacker who can shake the RP's trust bundle could force either false rejects (denial) or — if the RP falls back insecurely — false accepts.
*Required response:* verification of an *already-issued, untampered, un-expired* token is offline (C1 — base identity from existing infrastructure; the gate confirms this is solved). For valid tokens the RP uses locally-held trust material and no fresh fetch is required (FR5/NFR2/INV7). For revocation status, the partition/staleness handling of FM1/FM2/FM4 applies. If signature verification cannot be concluded (trust material absent/corrupted), fail-closed (ER11/SO4 `[HYP]`) — reject rather than accept or fetch insecurely.
*Trace:* C1, FR5, NFR2, ER7, INV7, SO4[HYP]; gate C1/C5 (valid-token offline verification solved). *Spike-dependence:* none for valid-token verification; revocation-status branch inherits FM1/FM2/FM4.

**FM10 — Non-conformant (malicious) relying party.**
*Trigger (adversarial):* the relying party itself is malicious or non-conformant — it chooses to accept delegations its checks should reject, or to skip checks.
*Adversarial angle:* the relying party is the trust decision-maker; a malicious RP is not something the system can force to verify correctly.
*Required response:* the system's guarantees are defined at the **conformant-verifier** boundary. A non-conformant RP is outside the system's trust boundary by definition; the system cannot and does not warrant verification outcomes by a verifier that does not conform. This is a boundary statement (SYSTEM_CONTEXT), not an internal mitigation. The engineering consequence: the acceptance tests (`05_ACCEPTANCE_TEST_PLAN.md`) target conformant-verifier behavior only; no test asserts control over a malicious RP.
*Trace:* `SYSTEM_CONTEXT.md` (boundary — the product proves *that* a delegation exists and its stated scope, not whether the scope should be granted), R1. *Spike-dependence:* none.

**FM11 — Silent trust failure (the meta-failure mode).**
*Trigger (adversarial, by construction):* the central product risk (R1): a delegation mechanism that appears to work under tested conditions but has an undetected flaw exploitable under adversarial conditions. Unlike a conventional bug, it does not announce itself.
*Adversarial angle:* the adversary succeeds *because the test suite missed the conditions*, not because of a visible malfunction.
*Required response:* this is a **meta**-failure about the adequacy of verification, not a runtime mode. Its mitigation is structural across three layers, all required:
  - *Single-check rollback* (SO5): for every required check, forcing that check to fail while others pass must flip the verdict to reject — so no single undetected check-failure yields acceptance.
  - *Adversarial acceptance tests* (`05_ACCEPTANCE_TEST_PLAN.md`): every SO and INV has a test constructed under adversarial conditions (mutation, partition, skew), not just the happy path.
  - *Process discipline* (lab): two-run reproducibility, adversary-blinded role separation, pre-registration primacy, and the garden-of-forging-paths guard (U5) — referenced from the lab governance, not duplicated here.
*Trace:* R1, SO5, SO3, SO1; `lab/DECISION_RULES.md` (U5), `lab/EXPERIMENT_CHECKLIST.md`. *Spike-dependence:* none — this is invariant to the spike outcome; whichever composition is chosen, silent-failure detection must hold.

## Coverage — risks and requirements addressed

| Source | Failure modes | Response owners |
|---|---|---|
| R1 (silent trust failure) | FM7, FM6, FM11 | SO3, SO5, SO6 + lab discipline |
| R2 (interop failure) | (services the security objectives SO2/SO7; not a runtime FM) | SO2, SO7 |
| R3 (standards conflict) | (boundary; INV11/ER16) | INV11, ER16 |
| R4 (vendor fragmentation) | (ecosystem risk; not a system runtime FM) | — (out of model) |
| R5 (unresolved buyer) | (claim-discipline C5; not a runtime FM) | — (out of model) |
| R6 (no external validation) | FM5, via SO8 | SO8 |
| FR3 | FM3, FM8 (post-expiry) | ER3, INV3 |
| FR4 | FM1, FM2, FM4 | ER5, SO1, INV4, INV12 |
| FR5 / NFR2 / C6 | FM1, FM9 | ER7, INV7, INV12, SO2 |
| NFR6 | FM5 (existing-record), FM7 | SO3, INV8 |
| C1 | FM9 (valid-token offline) | INV10 |

Risks R2, R3, R4, R5 have no dedicated runtime failure-mode entry because they are adoption/ecosystem/claim-discipline risks, not system runtime behaviors; their disposition is in `02_SECURITY_OBJECTIVES.md` (SO7) and in the claim-discipline constraints (C4, C5, handled in `05_ACCEPTANCE_TEST_PLAN.md`).

## Provenance

- **Primary source:** the Phase 7 Product Definition package (frozen), especially `ASSUMPTIONS_AND_RISKS.md`, the security-relevant FR/NFR items, and `CONSTRAINTS.md`.
- **Secondary:** `LEVEL0_1_FEASIBILITY_GATE.md` (S1, S4, and the candidate-composition space for FM2/FM4); `lab/DECISION_RULES.md`, `lab/EXPERIMENT_CHECKLIST.md` (FM11 process discipline).
- **Confidence:** FMs traced to evidenced Product Definition items (FM3, FM6, FM7, FM9, FM10, FM11) — High. FM1/FM2/FM4 — High as statements of distributed-systems behavior; their *resolution* depends on the founder scope act (R, S4) and on the spike outcome (composition). FM5 (key-compromise) and FM8 within-window replay — High confidence that they are **unmitigated within current scope** (negative findings, per project doctrine: confidence without evidence is forbidden, and confidence in an absence is asserted only when deducible).
- **Change-condition:** only an amendment to the frozen Product Definition authorizes adding a required response to FM5 or FM8 (i.e., warranting resistance the product currently does not) or otherwise changing the FM set. Resolving S1–S5 changes parameters and the floor of FM2/FM4, not the existence of the modes.

<!-- checkpoint: feat(examples): implement verification results state -->
