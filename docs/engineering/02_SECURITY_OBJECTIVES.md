# Security Objectives

**Phase:** 8 — Engineering Requirements.
**Source authority:** the frozen Phase 7 Product Definition package (security-relevant items: `FUNCTIONAL_REQUIREMENTS.md` FR2, FR4, FR5, FR6; `NON_FUNCTIONAL_REQUIREMENTS.md` NFR2, NFR3, NFR5, NFR6; `CONSTRAINTS.md` C1, C2, C6; `ASSUMPTIONS_AND_RISKS.md` risks R1, R2, R6). `01_ENGINEERING_REQUIREMENTS.md` ERs are cited where a security objective operationalizes an ER.
**What this document is not:** no architecture, technology, repository, API, or protocol. Each objective is **measurable** (stated metric, threshold, and test locus) and traces to a Product Definition source.

## Discipline

1. **Measurable by construction.** Each Security Objective (SO) states an *Objective*, a *Metric* (the quantity measured), a *Threshold* (the pass value), and a *Test locus* (where in the product boundary the measurement is taken). `05_ACCEPTANCE_TEST_PLAN.md` operationalizes each.
2. **Trace, don't invent.** Each SO carries a `Traces to` line to a Product Definition requirement, constraint, or risk. Security properties not traceable to the Product Definition are recorded explicitly under **Non-objectives** rather than smuggled in.
3. **Hypothesis status inherited.** SOs derived from `[HYPOTHESIS]` Product Definition items are marked `[HYPOTHESIS]` and are not committed (per `V1_SCOPE.md`, `DEFERRED.md`).
4. **Observability scope (gate S1/S4).** Revocation objectives are scoped to *observably-revoked* delegations — those whose revocation has become observable to the verifier per the system's propagation model (parameter R from `01_ENGINEERING_REQUIREMENTS.md` §Scope parameters, bounded by partition recovery per S4). Objectives never claim in-partition observability the Product Definition does not warrant.

## Security objectives

**SO1 — Revocation completeness (post-observability).** After a delegation's revocation has become observable to a relying party per the system's propagation model (within R of the revocation taking effect, and no later than partition recovery per S4), a conformant verifier shall reject that delegation; a revoked-and-observable delegation shall not be accepted.
*Traces to:* FR4, NFR4, ER5, ER6; R from gate S1, S4. *Risks addressed:* R1 (silent failure), in part.
*Metric:* across N revival attempts of delegations whose revocation is observable to the verifier, count the number accepted by a conformant verifier.
*Threshold:* accept-count = 0 for all observably-revoked delegations, across all trials and both post-recovery and within-R-after-effect regimes.
*Test locus:* relying-party verification boundary.

**SO2 — Offline verification independence.** Verification of a valid, unexpired, unrevoked delegation shall succeed with the relying party's network path to every shared authority disabled, and shall require no egress to a shared authority during core verification even when the path is enabled.
*Traces to:* FR5, NFR2, C6, ER7. *Risks addressed:* R2 (interoperability failure — a verification that secretly phones home is the interop failure mode).
*Metric:* (a) binary — does a valid delegation verify with the network path to all shared authorities disabled; (b) count of egress events to a shared authority during core verification with the path enabled.
*Threshold:* (a) succeeds; (b) egress count = 0.
*Test locus:* relying-party verification boundary, network path (disabled / instrumented).

**SO3 — Tamper-evidence.** Any alteration of a delegation record after creation shall be detectable by, and rejected by, a conformant verifier.
*Traces to:* NFR6, FR6, ER4. *Risks addressed:* R1 (an undetected alteration is the canonical silent-trust failure).
*Metric:* across a mutation set M (bit-flips, field substitutions, truncation, reordering of protected fields), the fraction of altered records that a conformant verifier both detects-as-altered and rejects.
*Threshold:* detection-and-rejection fraction = 1 over M.
*Test locus:* relying-party verification boundary, presented record.

**SO4 — Fail-closed on inconclusive verification.** `[HYPOTHESIS]` When a conformant verifier cannot reach a conclusive validity determination, it shall reject the delegation rather than accept it.
*Traces to:* NFR3, ER11. *Risks addressed:* R1 (acceptance under ambiguity is silent failure). *Not committed scope — `DEFERRED.md` D6.*
*Metric:* across the enumerated set C of inconclusive conditions (e.g., signature unverifiable, required trust material unavailable, clock beyond stated tolerance), the fraction of conditions under which the verifier rejects.
*Threshold:* rejection fraction = 1 over C.
*Test locus:* relying-party verification boundary, induced-inconclusive-state mechanics.

**SO5 — No silent acceptance (single-check rollback).** A conformant verifier shall not accept a delegation that fails any single required check while the others pass. Rolling back any one required check to a failing state shall flip the verdict to reject.
*Traces to:* FR1, FR2, FR3, FR4, NFR6 (the set of checks), ER1–ER5; R1. *Risks addressed:* R1 directly ("appears to work under tested conditions but has an undetected flaw exploitable under adversarial conditions").
*Required checks (per Product Definition):* identity binding (FR1), scope integrity (FR2 + NFR6), expiry (FR3), revocation observability (FR4 within R), signature/tamper (NFR6).
*Metric:* for each required check i, force check i to a failing state while holding all other checks passing; record whether the verifier rejects.
*Threshold:* rejection occurs for every single-check failure (fraction = 1 over checks).
*Test locus:* relying-party verification boundary, adversarial input construction.

**SO6 — Scope-subset enforcement at issuance.** The system shall not issue a delegation whose permission scope is not a strict subset of the principal's own permissions; attempts to create an over-scoped delegation shall fail.
*Traces to:* FR2, ER2. *Risks addressed:* R1 (over-scoped issuance is a silent privilege-escalation flaw).
*Metric:* across attempts to issue delegations with stated scope ⊄ the principal's permission set, the fraction refused at issuance.
*Threshold:* refusal fraction = 1.
*Test locus:* issuance boundary (principal's permission set to delegation scope).

**SO7 — Non-replacement interoperability.** Adoption of delegation-verification shall not require a relying party to remove or wholesale-replace the identity-verification mechanisms already in use in its environment; delegation-verification shall operate through an interface compatible with that baseline.
*Traces to:* NFR5, ER9. *Risks addressed:* R2 ("has not solved the stated problem — it has produced another proprietary identity system"), R3 (standards conflict).
*Metric:* (a) binary — is the relying party's pre-existing identity-verification interface still present and unmodified after delegation-verification is adopted; (b) binary — is delegation-verification exercised through an interface compatible with that baseline rather than a wholesale replacement.
*Threshold:* (a) yes; (b) yes.
*Test locus:* relying-party environment, existing identity-verification interface (a given of the environment, not an artifact proposed here).

**SO8 — Independent reviewability.** The system's security claims shall be stated precisely enough, and its acceptance tests reproducible enough, that an independent reviewer given only this specification package and a system build — and no privileged or non-disclosed information — can reproduce each SO and invariant verdict and reach the same result.
*Traces to:* R6 ("a security primitive's correctness is difficult to self-certify … Without an independent security review, confidence cannot be distinguished from actual correctness"). *Risks addressed:* R6 directly.
*Metric:* binary — for each SO in {SO1..SO7} and each invariant in `03_SYSTEM_INVARIANTS.md`, can an independent reviewer, with no privileged information, reproduce the verdict.
*Threshold:* all SO and invariant verdicts independently reproduced.
*Test locus:* independent reviewer with spec + build only.

## Coverage — security-relevant Product Definition items

| Product item | Security objective |
|---|---|
| FR2 | SO6 |
| FR4 | SO1, SO5 |
| FR5 | SO2 |
| FR6 | SO3 |
| NFR2 | SO2 |
| NFR3 `[HYP]` | SO4 `[HYP]` |
| NFR5 | SO7 |
| NFR6 | SO3, SO5 |
| C1, C2 | SO7 (non-replacement of the standard) |
| C6 | SO2 |
| R1 (silent trust failure) | SO1, SO3, SO4, SO5, SO6 |
| R2 (interop failure) | SO2, SO7 |
| R6 (no external validation) | SO8 |

## Non-objectives — security properties deliberately not asserted

The following are not security objectives because they are not traceable to the Product Definition. Asserting them would invent scope.

- **Replay resistance within the validity window.** The Product Definition specifies expiry (FR3) and tamper-evidence (NFR6) but not resistance to replay of a captured, still-valid delegation. Post-expiry replay is rejected by FR3; within-window replay is a recorded non-objective — see `04_FAILURE_MODEL.md` FM8.
- **Resistance to issuer signing-key compromise.** The Product Definition gives no key-compromise or key-rotation contract. A compromised issuer key can mint delegations that pass signature verification; the product warrants no resistance. See `04_FAILURE_MODEL.md` FM5 (recorded as unmitigated within current scope).
- **Continuous posture re-attestation.** `[HYPOTHESIS]` (FR10/ER14), deferred (`DEFERRED.md` D2). Not a security objective until V1 confirms it is achievable.
- **Cross-protocol interop security with non-adopting relying parties.** `[HYPOTHESIS]` (FR9/ER13), deferred (`DEFERRED.md` D1). Not a security objective until V1 confirms.
- **Security guarantees for ≥3 trust domains or non-SPIFFE environments.** Out of V1 (`DEFERRED.md` D3–D4). SOs above cover exactly the two-domain SPIFFE-coexisting scenario.

## Provenance

- **Primary source:** the Phase 7 Product Definition package (security-relevant subset), frozen by founder statement.
- **Secondary:**
`ASSUMPTIONS_AND_RISKS.md` (R1, R2, R6); `LEVEL0_1_FEASIBILITY_GATE.md` for S1/S4 (observability scope).
- **Confidence:** SO1–SO3, SO5–SO7, SO8 — High (sourced to evidenced Product Definition items). SO4 — inherits NFR3's hypothesis status. Metrics and thresholds are engineering translations of the source; confidence in their *measurability* is High, in their *achievement* depends on V1.
- **Change-condition:** only an amendment to the frozen Product Definition (per `CONTRIBUTING.md` §4) authorizes adding or strengthening an SO. Resolving S1–S5 adjusts R and the observability scope, not the SO set.
