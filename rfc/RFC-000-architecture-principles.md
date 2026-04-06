# RFC-000 — Architecture Principles

**Status:** Founding RFC of the Architecture phase. Enduring.
**Source authority:** the frozen Phase 7 Product Definition package and the frozen Phase 8 Engineering Requirements / Security Objectives / System Invariants / Failure Model / Acceptance Test Plan (26 documents, hash-pinned in `FROZEN.sha256`). RFC-000 operates *downstream* of that package — it bounds how architecture may be proposed, never what the product is.
**What this RFC is:** the principles, trade-off philosophy, and decision rules that every future architectural decision (every subsequent RFC) must satisfy. It is the constitution of the RFC series.
**What this RFC is not — explicitly, and bindingly:** it contains **no repository structure, no technology choice, no API, no protocol, no storage design, no deployment model, no implementation**. Not as examples, not as defaults, not as "candidates for later." Anything in those categories belongs in a future numbered RFC that traces to a requirement forcing it — never here. RFC-000 is the decree that those categories are deferred, not the place to begin them.

---

## Discipline (how RFC-000 was written, and how every future RFC must be)

1. **Trace, don't invent.** Every principle below traces to a frozen Product Definition item, an Engineering Requirement (`01_ENGINEERING_REQUIREMENTS.md` ER), a Security Objective (`02_SECURITY_OBJECTIVES.md` SO), an Invariant (`03_SYSTEM_INVARIANTS.md` INV), or a Failure Mode (`04_FAILURE_MODEL.md` FM). A principle that cannot be traced is, by construction, not allowed here. RFC-000 strengthens nothing; it constrains the *space of architectures* to those that satisfy the frozen requirements.
2. **Spike-outcome-independent.** The C4 feasibility spike outcome is unknown (α composition-works / β technology-gap / γ logical-impossibility / δ unresolvable). Every principle here holds whichever outcome the spike delivers. No principle assumes the spike succeeds; future RFCs must satisfy requirements for **all** values the founder sets for the scope parameters (R, S2/S3, S4, S5), not merely the hoped-for one.
3. **Adversarial and distributed, by construction.** Every principle is reasoned under both distributed-systems conditions (partition, replication, skew, partial failure) and adversarial conditions (forgery, over-issuance, tampering, replay, rollback). An architecture that holds only under cooperative, well-connected conditions satisfies none of these principles.
4. **Honest about limits.** Where the frozen package warrants no response — issuer key compromise (`04_FAILURE_MODEL.md` FM5), within-window replay (FM8) — the principle records that absence and forbids any future RFC from silently "solving" it by burying a mechanism in an architecture. Confidence without evidence is forbidden (`agents/GOVERNANCE.md` §5; `02_SECURITY_OBJECTIVES.md` Non-objectives).
5. **Principles operate on proposals, not on the system.** A *requirement* (ER) states what the system shall do. An *architectural principle* (here) states a property that every proposed architecture must have, and declares any proposal violating it invalid. RFC-000 is a constraint space, not a restatement of requirements.

---

## Architectural principles

Each principle names the frozen items it traces to and the constraint it places on future RFCs.

**AP1 — Offline-by-construction on the core verification path.** No architectural proposal shall require a network call to a shared authority at the moment of verification, for the core determination of validity, scope, expiration, and revocation status. The verification path's correctness must not depend on the liveness of any authority not already held by the relying party.
*Traces to:* FR5, NFR2, C6 (Phase 7); ER7, ER8 (Phase 8); INV7, INV12; SO2; FM9.
*Constraint on future RFCs:* any proposal that places a live call on the verification path is invalid on its face. The burden is on the proposal to show the path is offline; the absence of a live call is not a property to be tested for later, it is a precondition for admission of the RFC.

**AP2 — Companion, not replacement; extend, do not amend.** No architectural proposal shall require modifying the published definition of the workload-identity standard it coexists with, nor shall it assume responsibility for base workload-identity issuance. The architecture operates on already-issued identity material, as a companion capability.
*Traces to:* C1, C2, NFR5 (Phase 7); ER9, ER15, ER16 (Phase 8); INV10, INV11; SO7.
*Constraint on future RFCs:* a proposal that implies "adopting delegation-verification requires replacing the RP's existing identity-verification mechanism" or "the SPIFFE standard must be amended" is invalid. The architecture must be satisfiable against the standard as published.

**AP3 — Tamper-evidence of every delegation record.** No architectural proposal shall permit an altered delegation record to verify as the original. The record's integrity must be verifiable by a third party, independently of and after the original verification event, with no access to the original verifier's runtime state.
*Traces to:* NFR6, FR6 (Phase 7); ER4 (Phase 8); INV8, INV9; SO3; FM7.
*Constraint on future RFCs:* any data, transmission, or at-rest representation of a delegation must carry tamper-evidence as a structural property, not as a check performed by trusted infrastructure. Integrity is a property of the record, not of the channel or the store.

**AP4 — Fail-closed under inconclusive verification.** No architectural proposal shall cause a conformant verifier to accept a delegation when verification cannot reach a conclusive determination. When the verdict is undetermined, the architecture rejects.
*Traces to:* NFR3 `[HYPOTHESIS]` (Phase 7); ER11 `[HYPOTHESIS]`, C-INV1 candidate (Phase 8); SO4 `[HYPOTHESIS]`; FM3 (beyond clock tolerance), FM9 (trust material absent).
*Constraint on future RFCs:* the architecture's verification path must have a single, well-defined "inconclusive" state whose mandated outcome is rejection. This principle inherits NFR3's hypothesis status: a proposal may not commit this property as established; it must design for it and let V1 confirm (`DEFERRED.md` D6).

**AP5 — Honesty of claims (no false resistance).** No architectural proposal shall assert a security property the frozen failure model records as unmitigated, nor assert an observability bound the S4 limit forbids. The architecture's claimed guarantees shall be exactly those the frozen package warrants — no more, no fewer.
*Traces to:* FM5 (issuer key compromise, unmitigated), FM8 (within-window replay, unmitigated) (Phase 8); INV12 (S4 observability bound); SO8 (independent reviewability); R1, R6 (Phase 7 `ASSUMPTIONS_AND_RISKS.md`).
*Constraint on future RFCs:* a proposal that implies resistance to issuer key compromise, resistance to within-window replay, or in-partition revocation observability is invalid unless and until an amendment to the frozen Product Definition warrants it. The architecture shall carry the *unmitigated-within-current-scope* labels of FM5/FM8 forward into every claim it makes. An RFC that silently fills a gap the failure model left open is, by this principle, defective.

**AP6 — Two-domain scope discipline.** No architectural proposal shall assert verified behavior for more than two trust domains or for non-SPIFFE-coexisting environments. Architectural claims are bounded to the two-domain minimum experiment.
*Traces to:* C3, FR8 (Phase 7); ER8, ER17 (Phase 8); `DEFERRED.md` D3–D4.
*Constraint on future RFCs:* any proposal whose stated scope exceeds two domains or assumes a non-SPIFFE baseline is out of scope and invalid until the founder's scope act extends V1. Absence of a ≥3-domain result is a scope limit, not a defect to be "fixed" by an RFC.

**AP7 — Parametric satisfaction of the scope parameters.** No architectural proposal shall assume a specific resolution of the scope parameters R (revocation-observability latency), S2/S3 (cached-pull admissibility / broker definition), S4 (partition reading — the information-theoretic limit), or S5 (per-hop authorization, out of V1). The architecture must be statable for any value the founder sets, and must degrade honestly (reject, rather than over-claim) when a bound cannot be met.
*Traces to:* gate S1–S5 (`LEVEL0_1_FEASIBILITY_GATE.md`, cited as evidence only); FM1 (S4 partition), FM2 (R latency), FM4 (composition staleness ≤ R) (Phase 8); INV12.
*Constraint on future RFCs:* a proposal that bakes in a single value of R, or that only works if cached pulls are admissible (S2 relaxation), or that implicitly denies the S4 partition limit, is invalid. The RFC must state the parameter range it satisfies and must surface (not hide) the outcome where no composition meets the bound — that is the β/δ spike-outcome path and it must be designable-for, not precluded.

**AP8 — Independent reviewability of every architectural claim.** No architectural proposal shall rest on privileged or non-disclosed information. Every claim an RFC makes shall be stated precisely enough, and reproducible enough, that an independent reviewer with the specification package and a build — and nothing else — can reach the same verdict.
*Traces to:* SO8 (Phase 8); R6 (Phase 7).
*Constraint on future RFCs:* a proposal whose security argument depends on a secret assumption, an undocumented prior, or reasoning the reviewer cannot replicate from the package is invalid. The architecture is a public artifact; its correctness must be checkable by a stranger.

**AP9 — Single-choice rollback at the architectural level (falsifiability).** Every architectural choice in a future RFC shall be falsifiable: it must be possible to force that single choice to a wrong state while holding all others correct, and observe the system fail. An architectural choice that cannot be counter-held is a choice that cannot be reviewed.
*Traces to:* SO5 (Phase 8); FM11 (the silent-trust-failure meta-mode); R1.
*Constraint on future RFCs:* a proposal that introduces a check, a boundary, or a guarantee that "cannot fail visibly" carries the silent-failure risk of FM11. The RFC must name the observable consequence of the choice being wrong. No unverifiable cornerstone.

**AP10 — No over-issuance, no scope creep at issuance.** No architectural proposal shall permit a delegation to be issued whose scope is not a strict subset of the principal's permissions, nor permit a delegation's identities to be re-bound after issuance.
*Traces to:* FR1, FR2 (Phase 7); ER1, ER2 (Phase 8); INV1, INV2; SO6; FM6.
*Constraint on future RFCs:* the issuance boundary is a property the architecture must enforce structurally. A proposal that delegates "the scope will be checked downstream" pushes a checkpoint off the issuance boundary and is invalid; issuance refusal is the contract.

**AP11 — Minimize accidental complexity.** No architectural proposal shall introduce a mechanism, layer, or abstraction that is not forced by a frozen requirement or necessary as a substrate of a component that is. Every moving part is a place for a silent failure to hide and a thing the spike outcome can invalidate; the burden is on the RFC to justify each part's existence, and the default is absence.
*Traces to:* `04_FAILURE_MODEL.md` FM11 (the silent-trust-failure meta-mode — complexity is the substrate in which silent flaws live), R1 (Phase 7 `ASSUMPTIONS_AND_RISKS.md`); C4 (no production-readiness claim at the six-month horizon — do not build production machinery for a feasibility horizon); the Discipline §1 "trace, don't invent" doctrine carried through the Phase 8 package.
*Constraint on future RFCs:* a proposal that adds a component, a service, a layer, or a redundant check "just in case" or "for symmetry" is defective unless that component traces (DR1) to a forcing requirement. This principle is the extension of TP6 / DR9 from "do not choose an unjustified technology" to "do not build an unjustified part" — accidental complexity is itself an architectural choice, and it carries the silent-failure risk.

**AP12 — Stable module boundaries.** Module boundaries shall be placed so that (a) the externally-exposed boundary — the delegation record's interface to conformant verifiers and to third-party reviewers — is stable across the internal evolutions the spike outcome and future scope acts will force, and (b) the regions whose design the spike outcome determines are isolated from the regions whose design the requirements fix, so that changing the spike-dependent region does not destabilize the requirement-fixed region.
*Traces to:* ER4, INV9 (Phase 8 — the reconstruction record is a stable, self-sufficient interface a third party can verify without the original verifier's runtime state); INV11 (the external boundary is invariant across internal change); AP7 (parametric satisfaction — isolating the spike-parameter-dependent region is how an architecture becomes statable for any parameter value).
*Constraint on future RFCs:* a proposal whose module boundaries cross the spike-parameter-dependent regions in a way that forces a rewrite when S2/S3 or R are resolved, or that exposes an internal interface to verifiers such that an internal change rebinds a record's meaning, is invalid. The record-facing boundary (AP3, INV8/9) is the stable surface; the spike-facing boundary is the volatile surface; the architecture keeps them separate.

**AP13 — Observability by design.** Every architectural choice shall be introspectable at verification time: a conformant verifier's decision path shall expose which checks fired, in what order, with what outcome, so that the architecture can be exercised in single-check-failure states (SO5), an independent reviewer can reconstruct the verdict from the specification package and a build (SO8), and the silent-trust-failure meta-mode (FM11) has no dark corner to hide in.
*Traces to:* FM11 (the silent-trust-failure meta-mode — observability is its structural mitigation), R1 (Phase 7); SO5 (single-check rollback requires each check to be individually exercisable and observable); SO8 (independent reviewability requires the verdict to be reconstructable); the lab two-run reproducibility discipline (`lab/EXPERIMENT_CHECKLIST.md`, referenced by FM11).
*Constraint on future RFCs:* a proposal that places a load-bearing check on a path that cannot be observed — a check whose firing and outcome are not inspectable by a reviewer with the package and a build — is invalid; it carries the FM11 risk by construction. Observability is a structural property of the verification path, not a logging feature added later. This principle inherits no hypothesis status itself, but where it observes a hypothesis property (e.g., fail-closed behavior, AP4 / ER11), the observation must not promote that property to established (DR7).

---

## Trade-off philosophy

Principles will conflict. When they do, the resolution is not a vote; it is decided by these rules, in priority order. A future RFC that resolves a conflict must record the conflict and the rule it invoked — silent resolution is forbidden.

**TP1 — The information-theoretic limit wins over desirability.** When a stronger guarantee conflicts with the S4 limit (no fresh information across a partition), the S4 limit wins. The architecture shall not claim in-partition revocation observability it cannot have, however useful that claim would be. The trade-off resolves toward an *honest bound*, never toward a stronger-looking-but-false claim. (AP5, AP7, INV12.)

**TP2 — Coexistence wins over elegance.** When a mechanism that requires replacing or amending the existing standard (AP2) is more elegant than one that coexists, coexistence wins. The architecture's value is contingent on adoption; a non-adopting ecosystem makes the perfect architecture worthless. (AP2, SO7.)

**TP3 — Honesty wins over completeness.** When a complete-seeming design would silently "close" a gap the failure model records as unmitigated (FM5 key compromise, FM8 within-window replay), the design is rejected. The trade-off resolves toward recording the gap visibly and claiming less, not toward smuggling a mechanism and claiming more. (AP5, FM5, FM8.)

**TP4 — Fail-closed wins over availability, conditionally on the hypothesis.** When fail-closed (AP4) conflicts with availability (never reject a valid delegation), the conflict is real and the hypothesis status governs: AP4 inherits NFR3's `[HYPOTHESIS]`, so the architecture must make the fail-closed path *present and testable*, not *mandatory-by-claim*. The RFC designs for fail-closed and lets V1 confirm; it does not assert fail-closed as established, nor does it silently prefer availability. (AP4, ER11, `DEFERRED.md` D6.)

**TP5 — Two-domain honesty wins over generalization.** When a design generalizes naturally to N domains (AP6) but the evidence is two-domain, the two-domain bound wins. The trade-off resolves toward the bounded, evidenced claim; generalization is a future RFC's job under a future scope act. (AP6, ER17.)

**TP6 — The requirement forces the mechanism, never the reverse.** When a technology, protocol, API, or storage choice is tempting but no frozen requirement forces it, the choice is not made. The trade-off resolves toward *deferred* — the choice waits for the RFC that a requirement forces it into. (Discipline §1; the entire "no invention" doctrine of the Phase 8 package.)

---

## Decision rules for future RFCs

These are the procedural mandates every subsequent RFC shall satisfy. An RFC that violates a rule is not accepted.

**DR1 — Traceability mandate.** Every RFC shall cite, for each architectural decision it makes, the frozen ER, SO, INV, and/or FM that forces it. A decision with no trace is inventive and is grounds for rejection. (Coupled to Discipline §1.)

**DR2 — No-architecture-of-architecture mandate.** No RFC shall restate or re-rank the principles in RFC-000. RFC-000 is fixed by this founding act; future RFCs operate within it. Changing a principle requires amending RFC-000 via the frozen-document amendment procedure (`CONTRIBUTING.md` §4), not re-deriving it in a numbered RFC.

**DR3 — Spike-outcome-independence mandate.** No RFC shall assume the C4 spike outcome. Each RFC shall state which of the α/β/γ/δ outcomes its proposal supports, and shall either (a) support all four, or (b) state explicitly which it precludes and why that preclusion is acceptable given the founder's scope act. An RFC that only works under α is not a V1 architecture — it is a bet. (AP7, FM2, FM4.)

**DR4 — Adversarial-review mandate.** No RFC shall be accepted without an adversarial pass: for every load-bearing choice, the RFC shall state what observable failure occurs if the choice is wrong (AP9). A choice whose wrongness is invisible fails this mandate and carries the FM11 silent-trust-failure risk.

**DR5 — Honest-claims mandate.** No RFC shall claim a property the failure model records as unmitigated (FM5, FM8), nor an observability bound beyond S4 (INV12). Each RFC shall carry the unmitigated-within-current-scope labels forward into its own claims paragraph. (AP5.)

**DR6 — Two-domain-scope mandate.** No RFC shall assert ≥3-domain or non-SPIFFE behavior. (AP6, ER17.)

**DR7 — Hypothesis-preservation mandate.** No RFC shall convert a `[HYPOTHESIS]` property (ER11–ER14, SO4, NFR1/NFR3-sourced) into a committed architectural guarantee. Hypothesis properties remain testable-but-unconfirmed until V1; an RFC commits only to *designing for* them. (`DEFERRED.md`, `V1_SCOPE.md`.)

**DR8 — Independent-reviewability mandate.** No RFC shall rest its security argument on privileged or undisclosed information. Every claim shall be reproducible from the specification package plus a build. (AP8, SO8.)

**DR9 — No-scope-creep mandate.** No RFC shall introduce a technology, API, protocol, storage design, deployment model, or implementation choice that is not forced by a frozen requirement. The existence of a "natural" or "industry-standard" way to do something is not, by itself, a forcing requirement. (TP6.)

**DR10 — Freeze-eligibility mandate.** Once accepted, an RFC is a candidate for the hash-pinned frozen set (`scripts/frozen-docs.list` + `FROZEN.sha256`), exactly as the Phase 7 and Phase 8 documents are frozen. The freeze discipline of the planning phase continues into the architecture phase: accepted RFCs do not silently drift; they are amended only via a journal entry (`CONTRIBUTING.md` §4) and a re-baseline. RFC-000 itself is a candidate for that freeze in a future, explicitly authorized act — not in this founding act.

---

## What RFC-000 deliberately does not establish

To prevent the founding RFC from becoming an architecture in disguise, the following are *deferred* and may be established only by future numbered RFCs that trace to a requirement forcing them:

- **Repository structure / file layout.** RFC-000 does not prescribe directories, package boundaries, or file naming. Where future RFCs live on disk is an operational concern, not a principle.
- **Technology choices.** No language, library, crypto primitive, storage engine, container runtime, or identity stack component is named or implied.
- **APIs.** No interface, method signature, request/response shape, or schema is defined.
- **Protocols.** No wire format, message ordering, handshake, or transport is defined.
- **Storage.** No data model, persistence boundary, retention policy, or store choice is defined.
- **Deployment.** No topology, runtime placement, or deployment arrangement is defined.
- **Implementation.** No module decomposition, code structure, or build mechanism is defined.

An RFC that begins one of these does so by tracing each element to a forcing requirement (DR1, DR9) and by surviving an adversarial review (DR4). RFC-000's job is to ensure that when those RFCs come, they come honestly.

---

## Provenance

- **Primary source:** the frozen Phase 7 Product Definition package and the frozen Phase 8 Engineering Requirements / Security Objectives / System Invariants / Failure Model / Acceptance Test Plan (26 documents, hash-pinned in `FROZEN.sha256`, frozen by founder act on 2026-07-04 per journal entry `agents/agents/journal/2026-07-04-freeze-phase-7-phase-8.md`).
- **Secondary (evidence, not principle source):** `LEVEL0_1_FEASIBILITY_GATE.md` for the scope parameters S1–S5 and the S4 information-theoretic limit; `agents/agents/GOVERNANCE.md` for the honesty / confidence discipline.
- **Confidence:** AP1–AP3, AP5–AP13 — High, as traced restatements of architectural commitments already evidenced in the frozen package. AP4 — High as a statement of the required architectural posture, but the property it points at (fail-closed) inherits NFR3's `[HYPOTHESIS]` status and is not established. The trade-off philosophy (TP1–TP6) and decision rules (DR1–DR10) — High, as the architectural operationalization of the "no invention / honest claims / spike-outcome-independent" doctrine carried through the frozen Phase 8 documents.
- **Change-condition:** only an amendment to the frozen Product Definition (per `CONTRIBUTING.md` §4 amend procedure: journal entry → dated change note → `make frozen-baseline` → commit `FROZEN.sha256`) authorizes adding, removing, or strengthening a principle in RFC-000. Re-rankings among principles (the trade-off philosophy) are not arithmetic; they are re-decisions and require the same amendment path. Resolving the scope parameters S1–S5 by founder scope act adjusts what architectures must satisfy, not the principles themselves.
- **Freeze status of RFC-000:** not yet frozen. RFC-000 declares itself *eligible* for the frozen set (DR10) but is frozen only by a future, explicitly authorized act — not by this founding. Until then it is an uncommitted, unfrozen document on disk, presented for founder review.

---

## Architecture Review Process

This section defines the required structure every future RFC must follow and the review path it must survive to be accepted. It operationalizes DR1–DR10; it adds no new principles and is binding on RFC-001 onward.

### Required RFC structure

Every RFC shall contain the following sections, in this order. A section with no content is still written, stating "Not applicable to this RFC" and why — silence is not an allowed state for any required section.

1. **Header.** RFC number; title; status (`Draft` → `Adversarially Reviewed` → `Accepted` → `Frozen` → `Superseded`); author; date; `supersedes` / `superseded_by` links (initially none).
2. **Source authority.** The frozen package items (ER / SO / INV / FM / Phase 7 items) the RFC builds on, or the founder scope act that authorizes a new scope. (DR1 scope.)
3. **Scope of this RFC.** What this RFC defines *and* a deferral paragraph naming what it deliberately does not define, inheriting the "What this RFC deliberately does not establish" discipline. (DR9.)
4. **Traces.** A table mapping each architectural decision in the RFC to the frozen ER / SO / INV / FM that forces it. A decision with no forcing trace fails DR1 and is grounds for rejection. (DR1.)
5. **RFC-000 principle compliance.** For each applicable AP, TP, and DR, a one-line statement of how this RFC satisfies it (or why it does not apply). (Operationalizes AP1–AP13, TP1–TP6, DR1–DR10.)
6. **Spike-outcome analysis.** Which of α (composition-works) / β (technology-gap) / γ (logical-impossibility) / δ (unresolvable) outcomes the RFC supports; what it precludes and the founder-scope-act justification for any preclusion. (DR3, AP7.)
7. **Adversarial review.** For each load-bearing architectural choice, the observable consequence if the choice is wrong (what the reviewer would see when the system fails). A choice whose wrongness is invisible fails DR4. (AP9, AP13, DR4, FM11.)
8. **Honest-claims statement.** Which unmitigated-within-current-scope labels (FM5 issuer key compromise; FM8 within-window replay) and which S4 observability limits the RFC carries forward into its claims. (DR5, AP5, INV12.)
9. **Scope statement.** A declaration that the RFC asserts behavior only within the two-domain, SPIFFE-coexisting scenario, or — if it exceeds that — an explicit stop and a statement that it is out of V1. (DR6, AP6, ER17.)
10. **Hypothesis-preservation statement.** Which commitments in the RFC are `[HYPOTHESIS]` and remain testable-but-unconfirmed until V1; a promise that none are promoted to established. (DR7, AP4.)
11. **Reviewability statement.** What an independent reviewer needs — from the specification package and a build, and nothing else — to reproduce each claim's verdict. (DR8, AP8, AP13, SO8.)
12. **Open questions.** What the RFC deliberately does not answer, and which future RFC or founder act resolves each. (Journal-grade honesty.)
13. **Provenance.** Primary source, secondary evidence, confidence (per `agents/agents/GOVERNANCE.md` §5 — High / Medium / Low / None with what would change it), change-condition, and freeze status. (Convention of every document in this workspace.)

### Review path (state machine)

An RFC moves through these states. Only the founder advances a state (`agents/agents/GOVERNANCE.md` §1.5: the founder decides; agents do not vote).

- **Draft.** The RFC exists with all 13 required sections (any "Not applicable" placeholder filled with a reason). Not yet reviewed.
- **Adversarially Reviewed.** A review pass has exercised DR4 on every load-bearing choice: each choice is counter-held to a wrong state and the observable failure named in §7 of the RFC is confirmed to actually occur. A choice whose wrongness produces no observable consequence is struck and reworked. The review appends a dated, attributed review record (dissent preserved verbatim, per `agents/agents/GOVERNANCE.md` §1.6 / §2) to the journal, not to the RFC body.
- **Accepted.** The founder, reading the adversarial-review record and the RFC, accepts the RFC. Acceptance means the RFC's decisions are binding on the architecture. It does not mean they are frozen.
- **Frozen.** The accepted RFC is added to the hash-pinned set (`scripts/frozen-docs.list` → `make frozen-baseline` → commit `FROZEN.sha256`), per DR10, by a future explicitly authorized act — not automatic, not by this process's own momentum.
- **Superseded.** A later RFC replaces this one; `superseded_by` is set one-way (lifecycle per `agents/agents/journal/README.md`). The original is not deleted.

### What the review process does not do

- It does not vote (`GOVERNANCE.md` §1.5). Review records dissent; the founder decides.
- It does not allow silent resolution of principle conflicts — a conflict during review is recorded with the TP rule invoked, per the Trade-off philosophy.
- It does not accept an RFC that fails any of DR1, DR3, DR4, DR5, DR6, DR7, DR8, DR9 — these are gates, not preferences.
- It does not freeze an RFC by its own motion; freeze is a separate founder act (DR10), as it is for every document in this workspace.

<!-- checkpoint: planning(architecture-draft): audit architecture draft (#12) -->
