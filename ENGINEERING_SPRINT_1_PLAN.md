# Engineering Sprint 1 Plan — Foundation

**Status:** Sprint plan. Engineering document, not a product document. Not frozen.
**Authority:** RFC-003 (Accepted, 2026-07-05, founder statement), the frozen Phase 8 package (ER/SO/INV/FM/AT), `AI_BOOTSTRAP.md` §4 (vetted tooling), founder instruction of 2026-07-05 ("Engineering Sprint 1 — create the engineering foundation; do not implement business logic; stop after planning").
**Companion documents (this sprint's full deliverable set):** `PROJECT_MODULE_SPECIFICATION.md`, `REPOSITORY_SKELETON.md`, `MODULE_INTERFACE_SPECIFICATION.md`, `IMPLEMENTATION_ORDER.md`.

---

## 1. Sprint objective

Transform RFC-003's logical architecture into an implementation-ready software project: a compiling repository skeleton, fully specified module boundaries and interfaces, a test scaffold that can express the acceptance plan, and CI that enforces the dependency rules — such that implementation work (Sprint 2+) is a sequence of small, independently verifiable tasks with no open structural questions.

Sprint 1 produces **structure, contracts, and scaffolding**. It does not implement delegation semantics, and it does not touch anything blocked on the S1–S4 scope acts or the EXP-001 spike.

## 2. Stated assumptions (per working principles — explicit, checkable)

| # | Assumption | Basis | If wrong |
|---|---|---|---|
| A1 | The implementation stack is **Go 1.21+** with `go-spiffe/v2`, `spire-api-sdk`, `go-jose/v3` and standard libraries only | `AI_BOOTSTRAP.md` §4 "Non-Negotiable Development Rules" (2026-07-05 — newer than, and superseding, the Sprint-0 "stack deferred" placeholders in `Makefile` and `tests/README.md`); founder's prohibition of RFC-004 is consistent with the stack being settled | Only file extensions and toolchain tasks change; module structure is stack-independent |
| A2 | Product code lives in this repository, under `internal/` + `cmd/`, alongside the governance docs | No frozen item forbids it; the repo already carries `Makefile`, `Dockerfile`, `.github/`, `tests/` as engineering foundation | Skeleton re-roots under a subdirectory; nothing else changes |
| A3 | The record's integrity envelope is a JWS-signed token carrying the delegation claims (RFC 8693 `act`-style delegation expression, standard `exp`/`scope` claims) | Gate C2/C3/C7 record these as **solved** components with High confidence; `go-jose/v3` is the vetted library for exactly this; H1's own hypothesis text names this composition | The Record Model's *interface* is unchanged (RFC-003 M1 is envelope-agnostic); only `internal/record` internals change |
| A4 | Instance identity is realized, for Sprint 1 scaffolding, as an opaque, unique-per-issuance identifier element in the record | RFC-003 E2 requires the element to exist opaquely; M5/M6 contracts key on it | Deeper semantics (the FM5 amendment) land in the element's *production and comparison rules*, not in any interface shape |

A3 and A4 are flagged for **founder confirmation at sprint review** (§8). Neither blocks Sprint 1 tasks; both are recorded so no silent choice occurs.

## 3. In scope for Sprint 1

| Workstream | Content | Exit criterion |
|---|---|---|
| W1 — Repository skeleton | The tree in `REPOSITORY_SKELETON.md`: Go module, six `internal/` packages, three `cmd/` drivers, test scaffold under `tests/` | `go build ./...` and `go vet ./...` pass; every package compiles |
| W2 — Contracts as code | Every public type, interface, and closed answer set from `MODULE_INTERFACE_SPECIFICATION.md` exists and compiles: verdicts, causes, port interfaces, trace types, answer sets | Contract types compile; `revstatus/contracttest` suite runs green against the degenerate realization |
| W3 — Pure-logic foundations | M1 (record: integrity, reconstruction, mutation-detection) and M3 (verifier: pipeline, checks, verdict routing, decision trace) implemented against injected fakes — this is *engineering foundation*, not business logic beyond what the frozen checks define | Unit tests green for the AT-derived cases runnable in-process (see `IMPLEMENTATION_ORDER.md` milestone table) |
| W4 — Stores and register scaffolds | M4 (trust store: hold/answer/never-fetch), M5 (answer set + degenerate always-indeterminate realization + contract suite), M6 (append-only register) | Unit + contract tests green |
| W5 — CI and Makefile integration | Replace the Sprint-0 `test` no-op; add Go build/vet/test to `make ci` and `.github/workflows/`; add the dependency-rule lint (forbidden-import check) | `make ci` green locally; forbidden-import violations fail the build |
| W6 — Acceptance scaffold | `tests/acceptance/` files per AT family with the in-process subset expressed; `tests/harness/` interfaces for substrate control (implementations deferred to the substrate sprint) | In-process AT subset green; substrate-dependent ATs compile as skipped-with-cause (named blocker, not TODO) |

## 4. Explicitly out of scope for Sprint 1

- **Anything requiring S1–S4 values:** AT13/AT14 execution (need R and the S4 reading), the EXP-001 spike itself, any M5 realization beyond the degenerate one.
- **The two-domain SPIRE substrate** (EXP-001 plan Phases 2–6; shared with acceptance testing) — a separate, subsequent block of work.
- **The revocation propagation channel** (S2/S3-bounded, spike-selected).
- **Wire-format finalization beyond A3's envelope** — field-level encoding decisions are Sprint 2, gated on A3/A4 confirmation.
- **Any deployment, database, cloud, packaging, or performance work** (C4 horizon discipline; founder instruction).
- **New planning documents, new RFCs** (founder instruction).

## 5. Entry criteria (all satisfied)

- RFC-003 Accepted (founder statement, 2026-07-05). ✔
- Frozen package intact (`make check-frozen`). ✔
- Stack fixed (A1). ✔
- S1–S4 **not** required: every Sprint 1 item is parametric over them (RFC-003 AP7 posture) — policy values are injected, and tests exercise the logic at arbitrary parameter values.

## 6. Sprint risks

| Risk | Mitigation |
|---|---|
| Scope creep from "foundation" into delegation semantics beyond the frozen checks | The task list in `IMPLEMENTATION_ORDER.md` is closed; anything not on it needs a founder add |
| A3/A4 reversal after scaffolding | Both are isolated: A3 inside `internal/record` internals; A4 inside one element's production rule. Interfaces unaffected by construction |
| Dependency rules erode under implementation pressure | W5 makes the forbidden-import table (from `PROJECT_MODULE_SPECIFICATION.md`) a CI gate from the first commit, not a review convention |
| Substrate-dependent ATs silently rot as "skipped" | Every skip names its blocker (S1-scope-act / substrate / spike) and `IMPLEMENTATION_ORDER.md` maps each to the milestone that unblocks it |
| Hypothesis laundering in test names/assertions (fail-closed asserted as guaranteed) | AT22-derived tests carry the `[HYPOTHESIS]` marker in their documentation; they assert *observed designed behavior*, phrased per the AT plan's honesty discipline |

## 7. Definition of done for Sprint 1

1. `make ci` passes and now includes: build, vet, unit tests, contract tests, forbidden-import lint, plus the existing gates (lint, docs-lint, check-frozen, secrets).
2. Every module of RFC-003 exists as a package matching `PROJECT_MODULE_SPECIFICATION.md`, with its public surface exactly `MODULE_INTERFACE_SPECIFICATION.md` — no more, no less.
3. The in-process acceptance subset (identified per-AT in `IMPLEMENTATION_ORDER.md` §5) is green; every non-runnable AT compiles and skips with a named blocker.
4. The decision trace is emitted on every verification, including Accepts, and the single-check rollback protocol (AT23 logic) passes in-process.
5. No frozen document changed; no journal-requiring decision taken without an entry; commits conventional and small (one task, one commit).
6. Sprint review presents A3/A4 for founder confirmation and the Sprint 2 candidate scope (substrate block vs. issuance/driver completion — sequencing decision).

## 8. Founder decisions requested at sprint review (not before — none blocks Sprint 1)

| # | Decision | Smallest form |
|---|---|---|
| F1 | Confirm A3 (JWS envelope per gate C2/C3/C7 + vetted tooling) | yes/no at review |
| F2 | Confirm A4 (opaque unique-per-issuance instance identifier) — and whether to open the FM5 amendment now or defer to the revocation sprint | yes/no + timing |
| F3 | S1–S4 scope acts — unchanged ask from `ARCHITECTURE_READINESS_REVIEW.md` §11; needed before the substrate/spike block, which Sprint 2 planning will sequence | one journal entry |
