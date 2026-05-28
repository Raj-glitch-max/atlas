# Acceptance Test Plan

**Phase:** 8 — Engineering Requirements.
**Source authority:** the frozen Phase 7 Product Definition package and the four Phase 8 requirement documents (`01_ENGINEERING_REQUIREMENTS.md`, `02_SECURITY_OBJECTIVES.md`, `03_SYSTEM_INVARIANTS.md`, `04_FAILURE_MODEL.md`).
**What this document is not:** no architecture, technology, repository, API, or protocol.

## Discipline

1. **Tests verify requirements, not implementation.** Every acceptance test (AT) names a `Verifies` line pointing at an engineering requirement (ER), security objective (SO), invariant (INV), and/or failure mode (FM) — each of which itself traces to the Product Definition. Pass criteria are stated as outcomes observable at the product boundary (principal, delegate, relying party, presented delegation, trust material, the two-domain scenario). No AT references internal architecture, a chosen wire format, a library, a repository, or a protocol.
2. **Each AT is implementation-independent; the *environment* is given by the Product Definition.** Tests may reference the entities the Product Definition fixes as givens — the existing workload-identity infrastructure (C1), the relying party's *existing* identity-verification interface (NFR5), and the bounded two-domain scenario (FR8/C3). Tests shall not reference anything the Product Definition leaves to architecture.
3. **Hypothesis tests are marked.** Tests verifying `[HYPOTHESIS]` requirements (ER11/ER12/ER13/ER14, SO4) carry `[HYPOTHESIS]`; V1 either runs them to *document behavior* or defers them (per `V1_SCOPE.md`, `DEFERRED.md`). They are not pass-gated acceptance for committed scope.
4. **Every AT has a fail-route.** Each AT states, on failure, *which requirement failed* — so a negative result maps to a specific Product Definition item, not to "the system" (per the V1 Definition of Done, below).
5. **Adversarial and partition conditions are first-class.** ATs that exercise failure modes (FM) inject distributed/adversarial conditions directly: partition (disable the network path), clock skew (set the clocks apart), mutation (apply a mutation set), single-check rollback (force one check to fail). ATs are constructed under adversarial conditions, not only the happy path — directly addressing R1 (silent trust failure) via `04_FAILURE_MODEL.md` FM11.

## Execution discipline (process, referenced not duplicated)

Acceptance runs follow the lab governance in `lab/LAB_README.md`, `lab/EXPERIMENT_CHECKLIST.md`, and `lab/DECISION_RULES.md`:

- **Two-run reproducibility** — a result is reportable only if it reproduces on a second, independent run.
- **Adversary role separation** — at least one verifier-side check is performed by a blinded adversary role (constructs inputs to break the system without seeing the expected outputs).
- **Pre-registration primacy** — the `Verifies` and `Pass criterion` lines here are pre-registered; observed outcomes are recorded against them without post-hoc redefinition (`lab/DECISION_RULES.md` §1).
- **Partition injection is explicit** — "disable the network path" means severing the RP's path to the named shared authority/issuer at the link level, with out-of-band confirmation the path is down, not a flag flip in the verifier.
- **Pre-registered scope parameters** — R (S1) and the S2/S3/S4 readings are recorded for each run; an AT that depends on R names the R used.

## Acceptance tests

### A. Identity binding and presentation

**AT1 — Both identities recoverable from a presented delegation.**
*Verifies:* ER1, FR1, INV1, SO5 (identity-binding check). *Setup:* a delegation issued by principal p to delegate w. *Stimulus:* present the delegation to a conformant relying party with the RP's locally-held trust material; no side channel. *Pass:* the RP recovers both p and w, and they match issuance. *Fail →:* FR1/ER1/INV1 violated.

**AT2 — Missing binding is rejected.**
*Verifies:* ER1, FR1, INV1. *Setup:* a malformed delegation lacking one of the two bindings. *Stimulus:* present it. *Pass:* the RP rejects. *Fail →:* INV1 (a missing binding must never verify).

### B. Scope

**AT3 — Subset scope is inspectable.**
*Verifies:* ER2, FR2, SO6. *Setup:* a delegation with scope s ⊂ p's permissions. *Stimulus:* present it; ask the RP to inspect the scope. *Pass:* the RP reads s from the presented delegation and s is a strict subset of p's permissions. *Fail →:* FR2/ER2.

**AT4 — Over-scope issuance is refused.**
*Verifies:* ER2, FR2, SO6, INV2. *Setup:* an issuance request with scope s ⊄ p's permissions. *Stimulus:* attempt issuance. *Pass:* issuance refuses; no over-scoped delegation is created. *Fail →:* FR2/ER2/INV2/SO6.

**AT5 — Scope field tampered post-issuance is detected and rejected.**
*Verifies:* ER4, NFR6, SO3, SO5 (scope-integrity check), FM7. *Setup:* an issued delegation. *Stimulus:* alter the scope field after creation; present the altered delegation. *Pass:* alteration detected; delegation rejected. *Fail →:* NFR6/INV8/SO3 (and, if it slipped through, INV2 — scope-escalation through tampering).

### C. Expiration and clock

**AT6 — Non-expired delegation accepted (baseline).**
*Verifies:* ER3, FR3. *Setup:* a valid, unexpired, unrevoked delegation. *Stimulus:* present it with clocks aligned. *Pass:* accepted. *Fail →:* FR3/ER3 (false negative — needed to interpret AT7).

**AT7 — Expired delegation rejected.**
*Verifies:* ER3, FR3, INV3. *Setup:* the same delegation after its expiration time. *Stimulus:* present it. *Pass:* rejected. *Fail →:* FR3/ER3/INV3 (the acceptance criterion named verbatim in `FUNCTIONAL_REQUIREMENTS.md` FR3).

**AT8 — Clock skew beyond tolerance fails closed.** `[HYPOTHESIS]` (fail-closed portion)
*Verifies:* ER3 (clock tolerance), ER11/SO4 [HYP], FM3. *Setup:* set the RP clock and issuer clock apart by the stated tolerance t, then by t+ε. *Stimulus (at t):* present a delegation near its expiry; assert a deterministic verdict. *Stimulus (at t+ε):* present it; assert rejection. *Pass:* deterministic verdict within t; rejection beyond t. *Fail →:* ER3 (non-determinism) / ER11/SO4 (no fail-closed beyond tolerance) / FM3.

### D. Revocation

**AT9 — Revoked delegation is subsequently rejected.**
*Verifies:* ER5, FR4(a), INV4, SO1. *Setup:* an active delegation d. *Stimulus:* revoke d; present d to the RP after the revocation is observable (within R). *Pass:* d is rejected. *Fail →:* FR4(a)/ER5/INV4/SO1.

**AT10 — Revocation does not affect underlying identity or sibling delegations.**
*Verifies:* ER5, FR4(b), INV5, INV6. *Setup:* principal p has two delegations, d1 and d2, to delegates w1 and w2. *Stimulus:* revoke d1; (a) verify p's underlying identity; (b) verify w1's underlying identity; (c) verify the unrevoked d2. *Pass:* (a) p valid, (b) w1 valid, (c) d2 accepted. *Fail →:* FR4(b)/INV5 (identity affected) or INV6 (sibling cascade).

**AT11 — Revoked delegation never returns to valid.**
*Verifies:* INV4. *Setup:* revoked delegation d. *Stimulus:* re-verify d at three later times t1<t2<t3, with the revocation observable at each. *Pass:* rejected at all three. *Fail →:* INV4 (revocation not terminal).

**AT12 — Revocation takes effect without workload restart.**
*Verifies:* ER6, NFR4. *Setup:* an active delegation to a running workload w. *Stimulus:* revoke the delegation; record the workload's process continuity (PID / start time unchanged); present the delegation after observability. *Pass:* the delegation is rejected and w was neither restarted nor redeployed (process continuity preserved). *Fail →:* NFR4/ER6 (the acceptance criterion named verbatim in `NON_FUNCTIONAL_REQUIREMENTS.md` NFR4).

**AT13 — Revoked-but-observable acceptance count is zero (within R).**
*Verifies:* SO1, ER5, FM2. *Pre-registered:* R (S1). *Setup:* revoke a set of delegations. *Stimulus:* verify each at T0+R+ε (revocation observable) — adversary races the window. *Pass:* accept-count of observably-revoked delegations = 0 across all trials (post-recovery, within R after effect). *Fail →:* SO1/ER5; identifies the composition's staleness floor exceeding R (FM2 → surface as a spike/DOD finding, not a silent pass).

**AT14 — Partition-at-revocation honors the S4 impossibility (no false in-partition claim).**
*Verifies:* INV12, FM1. *Pre-registered:* the S4 resolution the founder sets. *Setup:* isolate the RP from the issuer *at the moment* of revocation; keep it isolated. *Stimulus:* present the revoked delegation while partitioned; then restore the path and present it after the bound (R or partition-recovery+P per S4). *Pass:* (a) the system does **not** claim the revocation is observable to the partitioned RP during the partition (a recorded non-claim — failing this means the system over-claims, INV12 violated); (b) after recovery+bound, the revoked delegation is rejected. *Fail →:* INV12/FM1 — either an over-claim of in-partition observability, or post-recovery non-rejection.

### E. Offline and no-live-call

**AT15 — Verification succeeds with network to shared authorities disabled.**
*Verifies:* ER7, FR5, NFR2, C6, INV7, SO2. *Setup:* disable the RP's network path to every shared authority; confirm the severance out-of-band. *Stimulus:* present a valid, unexpired, unrevoked delegation. *Pass:* validity, scope, and revocation status are all determined and the delegation is accepted. *Fail →:* FR5/NFR2/C6/ER7/INV7/SO2 (the acceptance criterion named verbatim in `NON_FUNCTIONAL_REQUIREMENTS.md` NFR2).

**AT16 — No egress to a shared authority during core verification.**
*Verifies:* ER7, INV7, SO2. *Setup:* network path to shared authorities *enabled*; out-of-band packet/egress instrumentation at the RP. *Stimulus:* verify a delegation; observe egress during core verification. *Pass:* zero egress events to a shared authority (per the S3 broker definition) during the verification. *Fail →:* INV7/SO2 (a hidden live call — R2 interop-failure mode made concrete).

### F. Cross-trust-domain (bounded to two domains)

**AT17 — Two-domain verification (valid / invalid).**
*Verifies:* ER8, FR8, C3, ER17. *Setup:* principal in trust domain A, RP in trust domain B, no shared authority, federation disabled (exactly two domains). *Stimulus:* (i) present a valid delegation — expect acceptance; (ii) present an expired / tampered / observably-revoked delegation — expect rejection. *Pass:* (i) accept, (ii) reject. *Fail →:* FR8/ER8/C3. *Scope note:* the test covers exactly two domains; absence of a ≥3-domain result is a documented scope limit (AT29, `DEFERRED.md` D3), not a failure.

### G. Ephemeral issuance

**AT18 — Delegation issued to an ephemeral, short-lived workload.**
*Verifies:* ER10, FR7. *Setup:* an ephemeral workload w with a bounded lifetime shorter than the delegation's expiry; no long-lived statically-provisioned identity. *Stimulus:* issue identity and delegation for w; verify the delegation before w's lifetime ends. *Pass:* issuance succeeds and the delegation verifies, with no long-lived static identity required for w. *Fail →:* FR7/ER10.

### H. Reconstruction record and tamper-evidence

**AT19 — Record alone is sufficient for third-party reconstruction.**
*Verifies:* ER4, FR6, INV9. *Setup:* an issued delegation and its record. *Stimulus:* give the unaltered record alone to a third party with no access to the original verifier. *Pass:* the third party recovers delegator, delegate, scope, and time. *Fail →:* FR6/ER4/INV9.

**AT20 — Mutation set is detected and rejected.**
*Verifies:* ER4, NFR6, SO3, INV8, FM7. *Setup:* a delegation record. *Stimulus:* apply a mutation set (bit-flips, field substitutions, truncation, reordering of protected fields); present each altered record. *Pass:* every alteration is detected and every altered record is rejected (fraction = 1). *Fail →:* NFR6/ER4/INV8/SO3/FM7.

**AT21 — Record is verifiable independent of the original verification event.**
*Verifies:* ER4, FR6, INV9. *Setup:* a delegation record. *Stimulus:* a third-party verifier with no access to the original verifier's runtime state verifies the record. *Pass:* verification succeeds. *Fail →:* FR6 ("independent of and after the original verification event")/INV9.

### I. Fail-closed and silent-failure

**AT22 — Inconclusive verification fails closed.** `[HYPOTHESIS]`
*Verifies:* ER11, NFR3, SO4, FM9-branch. *Setup:* induce an inconclusive state (e.g., signature unverifiable; or required trust material unavailable). *Stimulus:* present the delegation. *Pass:* reject rather than accept. *Per `DEFERRED.md` D6: V1 documents behavior; not a committed-scope gate.* *Fail →:* NFR3/ER11/SO4 (acceptance under ambiguity — silent failure R1).

**AT23 — Single-check rollback (no silent acceptance).**
*Verifies:* SO5, FM11, and the checks underlying INV1, INV2, INV3, INV4, INV8. *Required checks:* identity binding (FR1), scope integrity (FR2+NFR6), expiry (FR3), revocation-observability-within-R (FR4), signature/tamper (NFR6). *Setup:* for each check i, force check i to fail while holding all other checks passing. *Stimulus (per check):* present the delegation. *Pass:* the verdict is reject in every single-check-failure case (fraction = 1 over checks). *Fail →:* SO5/FM11; whichever check's failure did not flip the verdict names the silent-trust-failure path (R1).

### J. Interoperability / non-replacement

**AT24 — Delegation-verification does not replace the RP's existing identity-verification.**
*Verifies:* ER9, NFR5, SO7. *Setup:* record the RP's existing identity-verification interface as a baseline (a given of the environment). *Stimulus:* adopt delegation-verification; exercise it. *Pass:* the baseline interface is still present and unmodified, and delegation-verification operates through a compatible interface — neither removal nor wholesale replacement. *Fail →:* NFR5/ER9/SO7 (the primitive-vs-feature condition; R2/R3).

**AT25 — Cross-protocol interop with a non-adopting relying party.** `[HYPOTHESIS]` *deferred (D1)*
*Verifies:* ER13, FR9. *Setup:* a relying party that has adopted nothing new. *Stimulus:* verify a delegation through the RP's existing identity-verification mechanism alone. *Pass:* verification succeeds. *Not V1-committed* (`DEFERRED.md` D1); documented as scope-limited.

### K. Latency

**AT26 — Verification latency measured and reported (no threshold asserted in V1).** `[HYPOTHESIS]`
*Verifies:* ER12, NFR1. *Setup:* timing instrumentation at the RP verification boundary. *Stimulus:* verify delegations across the run. *Pass:* an end-to-end verification latency is measured, recorded, and documented. *Per `DEFERRED.md` D5: no pass/fail threshold in V1; V1 measures rather than commits.* *Fail →:* NFR1/ER12 only if unmeasured (a reporting failure, not a latency failure).

### L. Boundary

**AT27 — System does not issue base workload identity.**
*Verifies:* ER15, C1, INV10. *Setup:* (a) external identity issuer present; (b) external identity issuer absent. *Stimulus:* (a) attempt to produce a delegation bound to a workload identity; (b) attempt the same with the issuer absent. *Pass:* (a) delegation produced on already-issued identity; (b) system cannot produce a delegation bound to a workload identity. *Fail →:* C1/ER15/INV10 (the system assumed base-identity issuance).

**AT28 — Companion: satisfiable against the published standard with no amendment.**
*Verifies:* ER16, C2, INV11. *Setup:* the system's dependency on the existing workload-identity standard. *Stimulus:* inspect whether the system is satisfiable against the published standard with no outstanding amendment to that standard's definition. *Pass:* satisfiable against the published standard with no pending/required amendment. *Fail →:* C2/ER16/INV11.

**AT29 — Two-domain scope discipline (inspection).**
*Verifies:* ER17, C3. *Setup:* this acceptance plan and `01_ENGINEERING_REQUIREMENTS.md`. *Stimulus:* by inspection, confirm no AT or ER asserts ≥3-domain or non-SPIFFE behavior, and that acceptance tests cover exactly two domains. *Pass:* confirmed. *Fail →:* ER17/C3 (scope discipline broken).

### M. Independent reviewability

**AT30 — Independent reviewer reproduces every verdict.**
*Verifies:* SO8, R6, and the visibility of FM5/FM11 gaps. *Setup:* an independent reviewer given this specification package (`01`–`05`) and a system build, with no privileged or non-disclosed information. *Stimulus:* the reviewer reproduces each SO (SO1–SO7) and each established invariant (INV1–INV12) verdict. *Pass:* every SO and INV verdict is independently reproduced. *Fail →:* SO8/R6 (claims not independently reproducible — confidence vs. correctness gap, the central risk R6).

## Coverage — every ER / SO / INV / FM has at least one AT

| Artifact | Covered by |
|---|---|
| ER1 / INV1 | AT1, AT2 |
| ER2 / INV2 / SO6 | AT3, AT4 |
| ER3 / INV3 | AT6, AT7, AT8 |
| ER4 / NFR6 / INV8 / INV9 / SO3 | AT5, AT19, AT20, AT21 |
| ER5 / INV4 / INV5 / INV6 / SO1 | AT9, AT10, AT11, AT13 |
| ER6 / NFR4 | AT12 |
| ER7 / INV7 / SO2 | AT15, AT16 |
| ER8 / C3 / ER17 | AT17, AT29 |
| ER9 / SO7 | AT24 |
| ER10 / FR7 | AT18 |
| ER11 / SO4 `[HYP]` | AT8, AT22 |
| ER12 / NFR1 `[HYP]` | AT26 |
| ER13 / FR9 `[HYP]` | AT25 |
| ER14 / FR10 `[HYP]` | (deferred — D2; no AT in V1; recorded as scope-limited) |
| ER15 / C1 / INV10 | AT27 |
| ER16 / C2 / INV11 | AT28 |
| INV12 / FM1 | AT14 |
| FM2 | AT13 |
| FM3 | AT8 |
| FM4 | AT13 (composition staleness ≤ R) |
| FM5 (key compromise) | AT30 (visibility of the gap; no resistance test — see Non-objectives) |
| FM7 | AT5, AT20 |
| FM9 (fresh-trust-material branch) | AT8, AT22 |
| FM11 (silent-failure) | AT23 |
| SO5 | AT23 |
| SO8 / R6 | AT30 |

*Deliberately not covered by an acceptance test* (because the Product Definition warrants no property there): within-window replay (FM8), issuer key-compromise resistance (FM5). Both are recorded in `02_SECURITY_OBJECTIVES.md` Non-objectives and `04_FAILURE_MODEL.md`; an AT asserting resistance they are not owed would itself be a doctrine violation (confidence without evidence).

## Definition of Done — honesty (from `V1_SCOPE.md` and `FOUNDER_PROBLEM_FIT.md`)

V1 is complete when `TECHNICAL_VALIDATION.md`'s own success and failure criteria (P5, items 8–9) have been evaluated against a working implementation of the in-scope requirements above (FR1–FR8, NFR2/NFR4/NFR6, ER1–ER10, ER15–ER17, INV1–INV12) and the result — success **or** failure — is documented honestly.

- **A negative result that identifies which requirement failed and why is a valid, complete V1 outcome** (`V1_SCOPE.md`; `FOUNDER_PROBLEM_FIT.md` P5 item 8). ATs are written so a failure maps to a named requirement, not to "the system."
- Acceptance of AT13/AT14 outcomes that expose the C4 spike's composition-staleness-too-high (FM2) or the S4 impossibility (FM1) **is itself a V1-valid result** — it identifies the gap, it does not fail the project.
- No claim of production-hardening (C4) or buyer/commercial packaging (C5) is made by passing these tests; V1 produces a validated reference, not a hardened system (`CONSTRAINTS.md` C4).

## Provenance

- **Primary source:** the Phase 7 Product Definition package (frozen) and the four Phase 8 requirement documents (this set).
- **Secondary:** `LEVEL0_1_FEASIBILITY_GATE.md` (S1–S4, for AT8/AT13/AT14 parameters and the S4 impossibility); `lab/LAB_README.md`, `lab/EXPERIMENT_CHECKLIST.md`, `lab/DECISION_RULES.md` (execution discipline); `TECHNICAL_VALIDATION.md` P5 items 8–9 and `FOUNDER_PROBLEM_FIT.md` P5 item 8 (Definition of Done).
- **Confidence:** each AT's measurability — High (stated as boundary-observable outcomes). Each AT's *outcome* is not predicted (project doctrine forbids forecast confidence). Hypothesis-marked ATs inherit their source's hypothesis status.
- **Change-condition:** only an amendment to the frozen Product Definition authorizes adding, removing, or strengthening an acceptance test. Resolving S1–S5 by founder scope act fixes the pre-registered parameters (R, S4 reading) AT8/AT13/AT14 use; it does not add tests for properties the Product Definition does not warrant.

<!-- checkpoint: chore(fuzz): harden secrets scanner config -->
