# Engineering Requirements

**Phase:** 8 — Engineering Requirements (final specification phase before architecture).
**Source authority:** the frozen Phase 7 Product Definition package — `PRODUCT_DEFINITION.md` and the documents it cross-references (`FUNCTIONAL_REQUIREMENTS.md`, `NON_FUNCTIONAL_REQUIREMENTS.md`, `CONSTRAINTS.md`, `SYSTEM_CONTEXT.md`, `USE_CASE_CATALOG.md`, `USER_MODEL.md`, `V1_SCOPE.md`, `DEFERRED.md`, `ASSUMPTIONS_AND_RISKS.md`). `LEVEL0_1_FEASIBILITY_GATE.md` is cited only as evidence for the scope parameters (S1–S5) and the partition impossibility (S4) that the Product Definition's FR5/C6 already imply; the gate introduces no product requirement here.
**What this document is not:** it contains no architecture, no technology choice, no repository structure, no API, and no protocol. Every requirement below traces to a Product Definition artifact; none invents a product-level requirement absent there.

## Discipline

1. **Translate, don't invent.** Each engineering requirement (ER) carries a `Traces to` line naming the Product Definition requirement(s), use case(s), or constraint(s) it translates. An ER whose `Traces to` line cannot be written is, by construction, not allowed here.
2. **Testable by construction.** Each ER carries a `Testability` line stating how pass/fail is determined at the product boundary (principal, delegate, relying party, presented delegation, trust material, the two-domain scenario) without reference to any internal implementation. `05_ACCEPTANCE_TEST_PLAN.md` operationalizes these.
3. **Hypothesis status is inherited, not laundered.** Items whose source is marked `[HYPOTHESIS]` in the Product Definition are marked `[HYPOTHESIS]` here and are not committed scope (see `V1_SCOPE.md`, `DEFERRED.md`).
4. **Constraints C4 and C5 are claim-discipline, not runtime.** "No production-readiness claim at the six-month horizon" (C4) and "no presumed buyer" (C5) govern the strength of claims, not system behavior; they are enforced through `05_ACCEPTANCE_TEST_PLAN.md` (Definition of Done honesty) and `04_FAILURE_MODEL.md` (no false resistance claims), not as ERs.

## Scope parameters (unresolved, parameterizing not fixing)

These are ambiguities the feasibility gate (`LEVEL0_1_FEASIBILITY_GATE.md`, §"Strengthen") identifies as requiring a founder scope act before the C4 spike runs. They parameterize requirements below; the architecture must satisfy each ER for whatever value the founder sets.

- **R — revocation-observability latency (gate S1).** The bound within which a revoked delegation must be rejected once the revocation becomes observable, with `R < delegation TTL`. ERs that depend on R are written in terms of R, not a number.
- **Cached-pull admissibility / broker definition (gate S2, S3).** Whether a periodic signed-artifact pull, or a passive signed-blob cache, counts as the forbidden "live call to a shared authority" (C6/FR5). ERs touching a "live call" are written to the **strict** reading; if the founder relaxes S2/S3, the offline/no-live-call requirements **narrow** correspondingly — they never strengthen.
- **Partition reading (gate S4).** Revocation performed *while* the relying party is partitioned from the issuer is an information-theoretic limit (the gate's S4 finding, deducible from FR5/C6: no live call ⇒ no fresh information ⇒ no fresh-observability across the partition). No ER or downstream artifact claims observability stronger than eventual-upon-partition-recovery. See `04_FAILURE_MODEL.md` FM1, `03_SYSTEM_INVARIANTS.md` INV12.
- **Per-hop authorization (gate S5).** Applies only to multi-hop delegation. `FUNCTIONAL_REQUIREMENTS.md` FR1–FR8 treat a single delegation; multi-hop independent per-hop authorization is out of V1 scope. ERs below assume single-hop; S5 is not parameterized further here.

## Requirements

### A. Delegation representation and presentation

**ER1 — Identity binding.** The system shall represent a delegation as a single presentable unit from which a relying party can recover both the principal identity and the delegate identity, such that a conformant verifier determines both identities from the presented unit together with trust material the relying party already holds, without a side channel.
*Traces to:* FR1, UC1. *Evidence:* `TECHNICAL_VALIDATION.md` P5 item 1, item 8.
*Testability:* present a delegation to a conformant verifier; pass = both identities are recovered and match those established at issuance, with no information channel beyond the presented unit and locally-held trust material.

**ER2 — Scope subset and inspectability.** The system shall constrain a delegation's permission scope to be a strict subset of the principal's own permissions, and shall make that scope inspectable by the relying party from the presented delegation.
*Traces to:* FR2, UC2. *Evidence:* `PRODUCT_THESIS.md` P5 item 3.
*Testability:* (a) compare the delegation's stated scope to the principal's permission set — pass = strict subset; (b) confirm the relying party can read the scope from the presented delegation. Issuance of a delegation whose scope is not a strict subset shall fail (see `02_SECURITY_OBJECTIVES.md` SO6).

**ER3 — Expiration and explicit clock tolerance.** The system shall associate an expiration time with each delegation such that a conformant verifier rejects an expired delegation. The verifier's clock-skew tolerance relative to the issuer shall be explicit and bounded, so the expired/not-expired verdict is deterministic within tolerance.
*Traces to:* FR3, UC2, partially NFR3. *Evidence:* `PRODUCT_THESIS.md` P5 item 3; expiry handling is stock-verifier behavior.
*Testability:* (a) present an expired delegation — pass = rejection; present a non-expired delegation (all else equal) — pass = acceptance; (b) set the issuer and verifier clocks apart by the stated tolerance and assert the verdict is deterministic; set them apart beyond tolerance and assert rejection (the fail-closed portion inherits NFR3's hypothesis status).

**ER4 — Reconstruction record and tamper-evidence.** The system shall produce, for each delegation, a record sufficient for a third party to determine which identity delegated to which, with what scope, and at what time; the record shall be tamper-evident (any alteration after creation is detectable by a conformant verifier) and shall be verifiable by a third party independently of, and after, the original verification event, without access to the original verifier's runtime state.
*Traces to:* FR6, NFR6, UC5. *Evidence:* `PRODUCT_THESIS.md` P5 item 4; `TECHNICAL_VALIDATION.md` P5 item 1 ("provable" chain of custody).
*Testability:* (a) hand the unaltered record alone to a third party — pass = delegator, delegate, scope, and time are all recovered; (b) apply a mutation set to the record — pass = every alteration is detected and the altered record is rejected by a conformant verifier; (c) verify the record with no access to the original verifier's state — pass = verification succeeds.

### B. Revocation

**ER5 — Specific revocation, identity unaffected.** The system shall allow a specific delegation to be revoked such that (a) a conformant verifier subsequently rejects it, and (b) the underlying identity of the principal and of the delegate remain valid and unaffected. Revocation targets exactly one delegation; revoking one delegation is not, by itself, a cause for a conformant verifier to reject a different, unrevoked delegation.
*Traces to:* FR4, UC3; the non-cascade clause is the direct reading of FR4's "specific delegation" and `PRODUCT_DEFINITION.md`'s "independent revocability." *Evidence:* `PRODUCT_THESIS.md` P5 item 2, item 4.
*Testability:* revoke delegation d; (a) subsequent verification rejects d; (b) verify a different non-revoked delegation of the same principal — pass = accepted; verify the principal's and delegate's underlying identity — pass = valid.

**ER6 — Revocation without workload restart.** Revocation of a delegation shall take effect without requiring the affected delegate workload to be redeployed or restarted.
*Traces to:* NFR4, UC3. *Evidence:* `PRODUCT_THESIS.md` P5 item 4 ("without … restarting a workload").
*Testability:* revoke an active delegation while the delegate workload is running; confirm subsequent rejection; pass = rejection occurs with no change to the workload's running state (process continuity preserved, no restart/redeploy).

### C. Verification

**ER7 — Offline determination, no live call.** A relying party shall be able to determine the validity, scope, and revocation status of a presented delegation from the presented material together with trust material the relying party already holds, without a network call to a shared authority at the moment of verification.
*Traces to:* FR5, NFR2, C6, UC4. *Evidence:* `TECHNICAL_VALIDATION.md` P5 item 8 (success criteria), item 9 (failure criteria); `CONSTRAINTS.md` C6.
*Testability:* disable the relying party's network path to every shared authority and present a valid, unexpired, unrevoked delegation — pass = validity, scope, and revocation status are all determined and the delegation is accepted; additionally, with the network path enabled, assert out-of-band that no egress to a shared authority occurs during core verification.

**ER8 — Cross-trust-domain verification, bounded to two domains.** The system shall allow a relying party operating in a trust domain independent of the principal's to verify a delegation, within the bounded two-domain scenario defined in `TECHNICAL_VALIDATION.md`'s minimum experiment. No requirement here asserts behavior for more than two trust domains.
*Traces to:* FR8, C3, UC1. *Evidence:* `TECHNICAL_VALIDATION.md` P5 item 7.
*Testability:* place the principal in trust domain A and the relying party in trust domain B with no shared authority and federation disabled; present a valid delegation — pass = accepted; present an invalid one (expired, tampered, or observably-revoked) — pass = rejected. The test covers exactly two domains; absence of a ≥3-domain result is a scope limit, not a failure (see `DEFERRED.md` D3).

**ER9 — Non-replacement interoperability.** The delegation-verification mechanism shall not require a relying party to adopt a protocol that is incompatible with, or a wholesale replacement for, the identity-verification mechanisms already in use in that relying party's environment.
*Traces to:* NFR5, UC4. *Evidence:* `TECHNICAL_VALIDATION.md` P5 item 9; `ECOSYSTEM_THESIS.md` P5 final verdict — the property separating a platform primitive from an unadopted feature.
*Testability:* establish the relying party's existing identity-verification interface as a baseline; pass = delegation-verification operates through a compatible interface and requires neither removal nor wholesale replacement of that baseline.

### D. Ephemeral issuance

**ER10 — Issuance without long-lived static identity.** The system shall allow identity and delegation to be issued to a workload without requiring that workload to have a long, statically-provisioned lifetime; issuance shall succeed for a workload whose lifetime is shorter than the delegation's.
*Traces to:* FR7, UC6. *Evidence:* `TECHNICAL_VALIDATION.md` P5 item 2.
*Testability:* issue a delegation to an ephemeral workload with a bounded lifetime shorter than the delegation's expiry; pass = the delegation is issued and verifiable before the workload's lifetime ends, and issuance did not require a long-lived, statically-provisioned identity for that workload.

### E. Verified-as-hypothesis properties (not committed scope)

**ER11 — Fail-closed on inconclusive verification.** `[HYPOTHESIS]` When verification cannot conclusively determine a delegation's validity, the system shall reject the delegation rather than accept it.
*Traces to:* NFR3. *Evidence:* inferred, not directly stated (`NON_FUNCTIONAL_REQUIREMENTS.md` NFR3 note; `FOUNDER_PROBLEM_FIT.md` P5 items 6–7). *Per `DEFERRED.md` D6 / `V1_SCOPE.md`, V1 tests and documents actual behavior rather than assuming satisfaction.*
*Testability:* induce an inconclusive verification state (e.g., signature unverifiable, or required trust material unavailable) — pass = rejection.

**ER12 — Latency compatible with synchronous paths.** `[HYPOTHESIS]` Delegation verification shall complete within a latency compatible with synchronous request paths. V1 shall measure and report an observed end-to-end verification latency; V1 does not commit to a threshold value in advance.
*Traces to:* NFR1. *Evidence:* `TECHNICAL_VALIDATION.md` P5 item 8 used sub-100 ms as its own experiment's success threshold — a candidate target carried forward, not a committed product SLA (`NON_FUNCTIONAL_REQUIREMENTS.md` NFR1 note). *Per `DEFERRED.md` D5.*
*Testability:* measure end-to-end verification latency at the relying-party verification boundary over the acceptance-test run; pass = the measured value is recorded and documented. No pass/fail threshold is asserted in V1.

**ER13 — Cross-protocol interoperability with non-adopting relying parties.** `[HYPOTHESIS]` Verification shall interoperate with relying parties that have not adopted a new protocol beyond what they currently use for identity verification.
*Traces to:* FR9, UC7 region. *Evidence:* `TECHNICAL_VALIDATION.md` P5 item 9; `PRODUCT_THESIS.md` P5 item 6 — identified as the primitive-vs-feature condition; not confirmed achievable. *Per `DEFERRED.md` D1.*
*Testability:* (deferred) with a relying party that has adopted nothing new, verify a delegation — pass = verification succeeds through the relying party's existing identity-verification mechanism alone. Not a V1-committed test.

**ER14 — Posture-tied continued validity.** `[HYPOTHESIS]` A delegation's continued validity may be tied to the delegate's security posture at verification time, not only at issuance time.
*Traces to:* FR10, UC8. *Evidence:* `TECHNICAL_VALIDATION.md` P5 item 3 classifies this as unproven. *Per `DEFERRED.md` D2.*
*Testability:* (deferred) change the delegate's posture between issuance and verification — pass = the verdict reflects the verification-time posture. Not a V1-committed test.

### F. Boundary and scope

**ER15 — No base-identity issuance.** The system shall operate on workload identity material issued by existing external workload-identity infrastructure and shall not itself perform base workload-identity issuance.
*Traces to:* C1, UC; `SYSTEM_CONTEXT.md` ("outside this boundary"). *Evidence:* `ECOSYSTEM_THESIS.md` P5 items 1, 4.
*Testability:* with the external identity issuer present, the system produces delegations bound to already-issued workload identity; with it absent, the system cannot produce a delegation bound to a workload identity — pass in both cases.

**ER16 — Companion, not a spec change.** The system's operation shall not require any modification to the published definition of the existing workload-identity standard it extends; it shall be satisfiable against the published standard with no outstanding amendment.
*Traces to:* C2, NFR5. *Evidence:* `ECOSYSTEM_THESIS.md` P5 item 2 (the SPIFFE scoping conflict).
*Testability:* inspect the system's dependency on the existing workload-identity standard; pass = the system is satisfiable against the published standard with no pending or required amendment to that standard's definition.

**ER17 — Two-domain scope discipline.** All ERs above assert behavior only within the bounded two-trust-domain scenario. No ER here, and no architecture satisfying them, may assert verified behavior for more than two trust domains or for non-SPIFFE environments.
*Traces to:* C3, FR8, `V1_SCOPE.md`, `DEFERRED.md` D3–D4. *Evidence:* `TECHNICAL_VALIDATION.md` P5 item 7.
*Testability:* by inspection — pass = no ER asserts ≥3-domain or non-SPIFFE behavior, and `05_ACCEPTANCE_TEST_PLAN.md` acceptance tests cover exactly two domains.

## Coverage — every Product Definition item is translated

| Product item | Translated by | Notes |
|---|---|---|
| FR1 | ER1 | |
| FR2 | ER2 | + ER4 (scope tamper-evidence) |
| FR3 | ER3 | clock tolerance folded in |
| FR4 | ER5 | incl. non-cascade reading of "specific" |
| FR5 | ER7 | with NFR2, C6 |
| FR6 | ER4 | with NFR6 |
| FR7 | ER10 | |
| FR8 | ER8 | two-domain bound |
| FR9 `[HYP]` | ER13 `[HYP]` | not committed |
| FR10 `[HYP]` | ER14 `[HYP]` | not committed |
| NFR1 `[HYP]` | ER12 `[HYP]` | measure, don't commit |
| NFR2 | ER7 | |
| NFR3 `[HYP]` | ER11 `[HYP]` | not committed |
| NFR4 | ER6 | |
| NFR5 | ER9 | primitive-vs-feature |
| NFR6 | ER4 | |
| C1 | ER15 | boundary |
| C2 | ER16 | boundary |
| C3 | ER8, ER17 | scope bound |
| C4 | (claim-discipline) | enforced in `05_ACCEPTANCE_TEST_PLAN.md`, `04_FAILURE_MODEL.md` |
| C5 | (claim-discipline) | enforced in `05_ACCEPTANCE_TEST_PLAN.md` |
| C6 | ER7 | non-negotiable; inherits FR5 |

No Product Definition item is dropped; no ER lacks a `Traces to` line.

## Requirements explicitly not added (not invented)

The following look engineering-relevant but are **not** traceable to the Product Definition and are therefore deliberately absent as ERs:

- **Replay resistance within the validity window.** `FUNCTIONAL_REQUIREMENTS.md`/`NON_FUNCTIONAL_REQUIREMENTS.md` specify expiry (FR3 → post-expiry rejection) and tamper-evidence (NFR6) but do not specify resistance to replay of a captured, still-valid delegation. Recording it here would import a requirement from `LEVEL0_1_FEASIBILITY_GATE.md`'s H1 (a different artifact). It is a recorded non-objective in `02_SECURITY_OBJECTIVES.md` and a known, unmitigated mode in `04_FAILURE_MODEL.md` FM8.
- **Key-compromise / key-rotation contract.** The Product Definition gives no response to compromise of an issuer signing key. Stating one would invent scope; `04_FAILURE_MODEL.md` FM5 records it as unmitigated within current product scope, and the discipline forbids claiming resistance not warranted by the source.
- **Multi-hop independent per-hop authorization/revocation.** Outside V1 (gate S5; `FUNCTIONAL_REQUIREMENTS.md` is single-hop). A multi-hop ER would re-open the C4×N problem the gate explicitly and deliberately leaves to the spike.
- **≥3-trust-domain or non-SPIFFE behavior.** Out of V1 (`DEFERRED.md` D3, D4; `V1_SCOPE.md`).

## Provenance

- **Primary source:** the Phase 7 Product Definition package (frozen by founder statement).
- **Secondary (evidence, not requirement source):** `LEVEL0_1_FEASIBILITY_GATE.md` for the S1–S5 scope parameters and the S4 impossibility implied by FR5/C6; `TECHNICAL_VALIDATION.md`, `PRODUCT_THESIS.md`, `ECOSYSTEM_THESIS.md`, `FOUNDER_PROBLEM_FIT.md` at the cited items.
- **Confidence:** each ER's confidence equals its source's. Sources traced to evidenced Product Definition items: High. ERs traced to `[HYPOTHESIS]` sources (ER11–ER14): inherit hypothesis status — not established.
- **Change-condition:** a change to the frozen Product Definition (per `CONTRIBUTING.md` §4 amend procedure) is the only authorized way to add, remove, or re-strengthen an ER here. Resolving S1–S5 by founder scope act adjusts parameters, not the ER set.

<!-- checkpoint: feat(stores): add revstatus snapshot retrieval -->
