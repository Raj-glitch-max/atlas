# Product Definition — Delegated Workload/Agent Identity

Status: Phase 7, Product Definition package. Defines *what* will eventually exist. Contains no architecture, technology choice, API, database design, deployment model, code, or protocol design. Every other document in this package (`SYSTEM_CONTEXT.md`, `USER_MODEL.md`, `USE_CASE_CATALOG.md`, `FUNCTIONAL_REQUIREMENTS.md`, `NON_FUNCTIONAL_REQUIREMENTS.md`, `CONSTRAINTS.md`, `ASSUMPTIONS_AND_RISKS.md`, `V1_SCOPE.md`, `DEFERRED.md`) is cross-referenced from here rather than repeated.

## What this product is

A system that lets one workload or agent prove, to an independent verifier, that it is acting **on behalf of** a specific principal identity — with a bounded scope of permission, a bounded time window, and independent revocability — without requiring a live call to a shared, always-online authority at the moment of verification.

This statement is the same core technical hypothesis validated in `TECHNICAL_VALIDATION.md` (Section: P5, item 1) and refined in `PRODUCT_THESIS.md` (Section: P5, item 3). It is restated here as a product definition, not re-derived.

## Why it matters (evidence-carried, not re-argued)

- The council record (`PROJECT_HISTORY.md`, R1 cycle) is the only case where an independent reviewer *upgraded* the problem mid-review on a specific, checkable basis (SPIFFE-style standardization opening room).
- Existing identity infrastructure (SPIFFE/SPIRE, OAuth bearer tokens) solves *authentication* — "who is this workload" — at production scale, but has no native representation of delegation, scope-narrowing, or independent revocation. This gap is documented in `PRODUCT_THESIS.md` (P5, items 1–2) and is not re-argued here.
- The gap is real and currently contested, not settled: active 2026 IETF drafts and at least one proprietary vendor answer (Teleport, March 2026) confirm the problem is live, per `TECHNICAL_VALIDATION.md` (P5, item 6) and `ECOSYSTEM_THESIS.md` (P5, items 2–3).

## What this product is not

- It is not a replacement for SPIFFE/SPIRE. It is a companion capability that must coexist with the existing workload-identity ecosystem (see `CONSTRAINTS.md`, C1–C2).
- It is not a general-purpose secrets manager, policy engine, or API gateway. Those are natural *consumers* of what this product produces (see `SYSTEM_CONTEXT.md`), not things this product replaces.
- It is not, at this stage, a commercial product with a defined buyer. The buyer question was explicitly left unresolved in `FOUNDER_DECISION_BRIEF.md` (P5, item 3) and is not resolved here — see `CONSTRAINTS.md`, C5.
- It is not validated beyond a bounded prototype. `TECHNICAL_VALIDATION.md` and `FOUNDER_PROBLEM_FIT.md` both concluded that six months of disciplined work realistically produces a validated reference implementation, not a production-hardened system. See `V1_SCOPE.md`.

## How to read this package

Every requirement in `FUNCTIONAL_REQUIREMENTS.md` and `NON_FUNCTIONAL_REQUIREMENTS.md` is either (a) traceable to a specific prior document and finding, or (b) explicitly marked as a hypothesis where no such evidence exists. Nothing in this package should be read as re-opening or re-scoring any prior decision.

<!-- checkpoint: feat(stores): add panic handling middleware (#105) -->

<!-- checkpoint: chore(stores): tweak test assertions -->
