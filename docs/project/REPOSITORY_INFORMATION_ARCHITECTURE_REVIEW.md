# Repository Information Architecture Review — Atlas

**Date:** 2026-07-05
**Reviewer role:** Principal Engineer audit
**Scope:** Read-only. No files were moved, renamed, deleted, or modified.

---

## 1. Executive Summary

The Atlas repository has a well-designed engineering foundation (hooks, CI, frozen-doc integrity, secret scanning) and a coherent reasoning framework. However, the **repository root is severely overcrowded** with 29 Markdown files, most of which are frozen planning specifications that have no business living at root. This creates a discovery failure: an engineer opening the repository cannot distinguish entry points from archived planning artifacts without reading every file header.

A second structural problem exists in `context/`: the directory was established as AI session memory, but it has been loaded with technical documentation (architecture, data flows, security policy, operation manual, governance rules, glossary) that belongs in a `docs/` hierarchy. This conflates two distinct audiences — the AI session loader and the human engineer — and will cause `context/` to grow uncontrollably as the project scales.

A third structural anomaly is the `agents/agents/` double-nesting, which is a historical accident acknowledged in the codebase but never corrected. It creates confusion on every directory listing.

None of these problems require code changes. All are documentation and directory organisation. The repository is not broken — it is pre-implementation and the right time to fix this is now, before engineering begins.

---

## 2. Repository Health Score

**6 / 10**

| Dimension | Score | Rationale |
|---|---|---|
| Engineering foundation (CI, hooks, secrets, frozen-doc guard) | 9/10 | Solid. gitleaks, pre-commit, SHA-256 integrity, Conventional Commits all in place. |
| Root discoverability | 3/10 | 29 Markdown files at root. No navigation hierarchy. Entry points buried alongside frozen specs. |
| Documentation currency | 5/10 | README references deleted files (`context/00_PROJECT_MISSION.md`, `claude-skills/`). DEVELOPMENT_RULES.md RFC Policy is stale. |
| Directory structure clarity | 6/10 | `agents/agents/` double-nesting, `context/` overloaded, no `docs/` hierarchy. |
| Separation of concerns | 5/10 | AI memory and human documentation mixed in `context/`. Planning specs and entry-point files mixed at root. |

---

## 3. Information Architecture Score

**4 / 10**

The root has no navigational hierarchy. A new engineer must read 29 files to understand what is a governance entry point, what is a frozen planning spec, and what is a reference document. There is no `docs/` directory. The frozen Phase 7 and Phase 8 specification packages — 15 files — are dumped directly at root, indistinguishable from `CONTRIBUTING.md` and `README.md` at a glance. This is the single largest information architecture failure in the repository.

---

## 4. Documentation Organisation Score

**4 / 10**

Two glossaries exist with overlapping scope (`agents/agents/GLOSSARY.md` and `context/06_GLOSSARY.md`). `DEVELOPMENT_RULES.md` explicitly describes itself as a summary of rules already in `CONTRIBUTING.md` — creating a drift risk between two files claiming the same authority. `context/` contains at least six files that are documentation, not AI memory. The `context/04_OPERATION_MANUAL.md` duplicates the Make-target table already in `CONTRIBUTING.md` §2.

---

## 5. Root File Audit Table

| File | Current Location | Recommended Location | Decision | Justification |
|---|---|---|---|---|
| `README.md` | Root | Root | **KEEP** | Repository entry point. Update stale references to deleted files. |
| `CONTRIBUTING.md` | Root | Root | **KEEP** | Engineering entry point. Primary governance reference. |
| `SECURITY.md` | Root | Root | **KEEP** | GitHub convention. Required at root. |
| `CLAUDE.md` | Root | Root | **KEEP** | AI session entry point. Must be at root. |
| `Makefile` | Root | Root | **KEEP** | Build entry point. |
| `Dockerfile` | Root | Root | **KEEP** | Dev toolbox. |
| `FROZEN.sha256` | Root | Root | **KEEP** | Referenced by `make check-frozen` by path. Cannot move without updating Makefile and CI. |
| `DEVELOPMENT_RULES.md` | Root | Root → `docs/project/` | **MOVE TO docs/project/** | Self-described as a summary of `CONTRIBUTING.md`. Not an entry point. RFC Policy section is now stale (RFCs are active, not archived). Retaining it at root creates a two-source-of-truth problem for commit and freeze policy. |
| `PROJECT_STRUCTURE.md` | Root | `docs/project/` | **MOVE TO docs/project/** | Useful reference documentation but not a repository entry point. Engineers reach for this after orientation, not before. |
| `PROJECT_HISTORY.md` | Root | `docs/project/` | **MOVE TO docs/project/** | 391-line historical chronicle. Valuable provenance. Not an entry point. A link from README is sufficient. |
| `REPOSITORY_AUDIT.md` | Root | `docs/project/` | **MOVE TO docs/project/** | Meta-document about hygiene actions. Not an entry point. |
| `PRODUCT_DEFINITION.md` | Root | `docs/specs/product/` | **MOVE TO docs/specs/product/** | Frozen Phase 7 planning spec. Not an entry point. Cross-references nine sibling specs. |
| `SYSTEM_CONTEXT.md` | Root | `docs/specs/product/` | **MOVE TO docs/specs/product/** | Frozen Phase 7 spec. Part of the product definition package. |
| `USER_MODEL.md` | Root | `docs/specs/product/` | **MOVE TO docs/specs/product/** | Frozen Phase 7 spec. |
| `USE_CASE_CATALOG.md` | Root | `docs/specs/product/` | **MOVE TO docs/specs/product/** | Frozen Phase 7 spec. |
| `FUNCTIONAL_REQUIREMENTS.md` | Root | `docs/specs/product/` | **MOVE TO docs/specs/product/** | Frozen Phase 7 spec. |
| `NON_FUNCTIONAL_REQUIREMENTS.md` | Root | `docs/specs/product/` | **MOVE TO docs/specs/product/** | Frozen Phase 7 spec. |
| `CONSTRAINTS.md` | Root | `docs/specs/product/` | **MOVE TO docs/specs/product/** | Frozen Phase 7 spec. |
| `ASSUMPTIONS_AND_RISKS.md` | Root | `docs/specs/product/` | **MOVE TO docs/specs/product/** | Frozen Phase 7 spec. |
| `V1_SCOPE.md` | Root | `docs/specs/product/` | **MOVE TO docs/specs/product/** | Frozen Phase 7 spec. |
| `DEFERRED.md` | Root | `docs/specs/product/` | **MOVE TO docs/specs/product/** | Frozen Phase 7 spec. |
| `01_ENGINEERING_REQUIREMENTS.md` | Root | `docs/specs/engineering/` | **MOVE TO docs/specs/engineering/** | Frozen Phase 8 spec. Numeric prefix implies belonging in a numbered package, not root. |
| `02_SECURITY_OBJECTIVES.md` | Root | `docs/specs/engineering/` | **MOVE TO docs/specs/engineering/** | Frozen Phase 8 spec. |
| `03_SYSTEM_INVARIANTS.md` | Root | `docs/specs/engineering/` | **MOVE TO docs/specs/engineering/** | Frozen Phase 8 spec. |
| `04_FAILURE_MODEL.md` | Root | `docs/specs/engineering/` | **MOVE TO docs/specs/engineering/** | Frozen Phase 8 spec. |
| `05_ACCEPTANCE_TEST_PLAN.md` | Root | `docs/specs/engineering/` | **MOVE TO docs/specs/engineering/** | Frozen Phase 8 spec. |
| `FOUNDER_DECISION_BRIEF.md` | Root | `docs/research/` | **MOVE TO docs/research/** | 459-line research analysis. Active reference but not an entry point. Frozen. |
| `RESEARCH_PROGRAM.md` | Root | `docs/research/` | **MOVE TO docs/research/** | Frozen research planning doc. Not an entry point. |
| `TECHNICAL_VALIDATION.md` | Root | `docs/research/` | **MOVE TO docs/research/** | Frozen research output. Not an entry point. |
| `FOUNDER_PROBLEM_FIT.md` | Root | `docs/research/` | **MOVE TO docs/research/** | Frozen research output. Not an entry point. |
| `P5_FALSIFICATION_EXPERIMENT.md` | Root | `docs/research/` | **MOVE TO docs/research/** | Frozen experiment spec. Directly references `lab/`. Should live near `lab/` conceptually but `docs/research/` is the best available location without creating a new top-level directory. |
| `LEVEL0_1_FEASIBILITY_GATE.md` | Root | `docs/research/` | **MOVE TO docs/research/** | Frozen feasibility gate. Same rationale as above. |

---

## 6. Directory Audit Table

| Directory | Current State | Recommendation | Rationale |
|---|---|---|---|
| `agents/` | Contains only `agents/agents/`. Outer directory serves no purpose. | **Rename outer: flatten to single `agents/`** | The double-nesting `agents/agents/` is a confirmed historical accident. Every internal file path already uses the inner directory. The outer `agents/` shell is a navigation dead-end. This is the highest-value structural fix in the repository. |
| `context/` | 9 files. Mix of AI memory and technical documentation. | **Narrow scope: AI memory only.** Migrate documentation to `docs/`. | `context/` should contain only files loaded by AI sessions to understand current project state. Files `01_GOVERNANCE.md`, `02_SYSTEM_ARCHITECTURE.md`, `03_DATA_FLOW.md`, `04_OPERATION_MANUAL.md`, `06_GLOSSARY.md`, and `07_SECURITY_POLICY.md` are documentation, not session memory. They belong in `docs/`. |
| `rfc/` | 3 active RFC drafts (RFC-000, RFC-001, RFC-002). | **Correct. Keep.** | RFCs are correctly isolated. Top-level placement signals their architectural weight. The `DEVELOPMENT_RULES.md` RFC Policy is stale (says "no new RFCs") and must be updated, but the directory itself is correct. |
| `lab/` | 6 process docs + `evidence/` placeholder. | **Correct. Keep.** | Lab is appropriately isolated. Frozen process docs are correctly separated from live log (`EXPERIMENT_LOG.md`). Boundary is clear. No engineering code has leaked in. |
| `archive/` | Contains only `archive/research/`. | **Correct. Keep.** | Archive correctly houses retired research material. Could add `archive/YYYY-QN/` dating for long-term navigability but not a current problem. |
| `scripts/` | 2 files: `check-frozen-docs.sh`, `frozen-docs.list`. | **Correct. Keep.** | Appropriately scoped. `frozen-docs.list` will need updating if frozen files are moved to `docs/specs/`. |
| `tests/` | 1 file: `README.md` placeholder. | **Correct. Keep.** | Correct placeholder. Will be replaced with a real test suite once a stack is selected. |
| `.github/` | Issue template, PR template, Dependabot, CI workflow. | **Correct. Keep.** | Standard placement. CI workflow path references should be verified after any file moves. |
| `docs/` | **Does not exist.** | **Create. Required.** | Without `docs/`, all non-root-entry documentation has no coherent home, forcing it to either clutter root or be buried inside `context/` or `agents/`. |

---

## 7. Proposed Final Repository Tree

Conceptual only. No files are moved by this document.

```text
atlas/
│
├── README.md                        # entry point: orientation + quickstart
├── CONTRIBUTING.md                  # entry point: engineering rules
├── SECURITY.md                      # entry point: security reporting
├── CLAUDE.md                        # entry point: AI session loader
├── Makefile                         # entry point: commands
├── Dockerfile                       # entry point: dev toolbox
├── FROZEN.sha256                    # integrity baseline (referenced by Makefile by path)
│
├── .github/                         # CI, templates, Dependabot
├── .devcontainer/                   # dev container
├── .vscode/                         # editor settings
├── .pre-commit-config.yaml
├── .gitleaks.toml
├── .gitignore / .gitattributes / .editorconfig
│
├── agents/                          # reasoning framework (FLATTENED from agents/agents/)
│   ├── README.md
│   ├── GOVERNANCE.md
│   ├── GLOSSARY.md                  # single authoritative glossary (merge with context/06)
│   ├── REFERENCES.md
│   ├── council/
│   ├── domain/
│   ├── working/
│   ├── templates/
│   └── journal/
│
├── context/                         # AI session memory ONLY
│   ├── 00_PROJECT_CONTEXT.md        # project mission, pipeline position
│   ├── 05_DECISION_LOG.md           # key decisions (borderline; keep for session continuity)
│   └── 08_AI_HANDOFF.md             # current state, next actions, surviving hypotheses
│
├── docs/
│   ├── project/
│   │   ├── PROJECT_HISTORY.md
│   │   ├── PROJECT_STRUCTURE.md
│   │   ├── DEVELOPMENT_RULES.md     # or merged into CONTRIBUTING.md
│   │   └── REPOSITORY_AUDIT.md
│   │
│   ├── research/
│   │   ├── FOUNDER_DECISION_BRIEF.md
│   │   ├── RESEARCH_PROGRAM.md
│   │   ├── TECHNICAL_VALIDATION.md
│   │   ├── FOUNDER_PROBLEM_FIT.md
│   │   ├── P5_FALSIFICATION_EXPERIMENT.md
│   │   └── LEVEL0_1_FEASIBILITY_GATE.md
│   │
│   ├── specs/
│   │   ├── product/                 # Phase 7 product definition package (10 files)
│   │   │   ├── PRODUCT_DEFINITION.md
│   │   │   ├── SYSTEM_CONTEXT.md
│   │   │   ├── USER_MODEL.md
│   │   │   ├── USE_CASE_CATALOG.md
│   │   │   ├── FUNCTIONAL_REQUIREMENTS.md
│   │   │   ├── NON_FUNCTIONAL_REQUIREMENTS.md
│   │   │   ├── CONSTRAINTS.md
│   │   │   ├── ASSUMPTIONS_AND_RISKS.md
│   │   │   ├── V1_SCOPE.md
│   │   │   └── DEFERRED.md
│   │   │
│   │   └── engineering/             # Phase 8 engineering requirements package (5 files)
│   │       ├── 01_ENGINEERING_REQUIREMENTS.md
│   │       ├── 02_SECURITY_OBJECTIVES.md
│   │       ├── 03_SYSTEM_INVARIANTS.md
│   │       ├── 04_FAILURE_MODEL.md
│   │       └── 05_ACCEPTANCE_TEST_PLAN.md
│   │
│   └── engineering/                 # future: architecture, runbooks, API docs
│       ├── SYSTEM_ARCHITECTURE.md   # moved from context/02
│       ├── DATA_FLOW.md             # moved from context/03
│       ├── OPERATION_MANUAL.md      # moved from context/04
│       └── SECURITY_POLICY.md       # moved from context/07
│
├── rfc/                             # architecture RFCs
│   ├── RFC-000-architecture-principles.md
│   ├── RFC-001-system-context.md
│   └── RFC-002-conceptual-domain-model.md
│
├── lab/                             # engineering research laboratory
│   ├── LAB_README.md
│   ├── EXPERIMENT_LOG.md
│   ├── EXPERIMENT_CHECKLIST.md
│   ├── EVIDENCE_INDEX.md
│   ├── DECISION_RULES.md
│   ├── EXP-001-EXECUTION-PLAN.md
│   └── evidence/
│
├── archive/
│   └── research/                    # retired research documents
│
├── scripts/
│   ├── check-frozen-docs.sh
│   └── frozen-docs.list             # must be updated with new paths after any move
│
└── tests/
    └── README.md
```

---

## 8. Migration Plan

Ordered smallest-risk first. Each step is independently completable and independently rollback-safe.

### Step 1 — Fix stale references (zero file moves, zero frozen-doc risk)
Update `README.md` to remove references to `context/00_PROJECT_MISSION.md` (deleted) and `claude-skills/` (deleted). Update `DEVELOPMENT_RULES.md` RFC Policy to reflect that `rfc/` contains active drafts, not archival material. No frozen docs touched. CI passes without changes.

### Step 2 — Resolve the two-glossary problem (one file change)
`agents/agents/GLOSSARY.md` (330 lines, framework-specific) and `context/06_GLOSSARY.md` (created in the last hygiene pass, protocol + framework terms) overlap. Decide which is authoritative and consolidate. The agents-internal glossary has richer framework content and older provenance. The context glossary has protocol terms not in the agents glossary. The correct answer is one file at `agents/agents/GLOSSARY.md` (or `agents/GLOSSARY.md` post-flatten) that covers both domains, with `context/06_GLOSSARY.md` removed. This requires a journal entry per governance.

### Step 3 — Create `docs/` hierarchy and move non-frozen root docs (no frozen-doc amendments)
Move files that are NOT in `scripts/frozen-docs.list`:
- `PROJECT_HISTORY.md` → `docs/project/`
- `PROJECT_STRUCTURE.md` → `docs/project/`
- `DEVELOPMENT_RULES.md` → `docs/project/` (or consolidate into `CONTRIBUTING.md`)
- `REPOSITORY_AUDIT.md` → `docs/project/`

Update `CLAUDE.md` and `README.md` to point to new paths. No `make check-frozen` risk. No frozen docs touched.

### Step 4 — Move frozen research docs to `docs/research/` (requires frozen-doc amendment process per file)
Move the 6 frozen research-planning documents:
`FOUNDER_DECISION_BRIEF.md`, `RESEARCH_PROGRAM.md`, `TECHNICAL_VALIDATION.md`, `FOUNDER_PROBLEM_FIT.md`, `P5_FALSIFICATION_EXPERIMENT.md`, `LEVEL0_1_FEASIBILITY_GATE.md`.

Each requires: journal entry → dated change-note in file → `scripts/frozen-docs.list` path update → `make frozen-baseline` → single commit. All 6 can be bundled into one journal entry covering the migration rationale. The frozen-doc script checks content hashes, not paths, so updating `frozen-docs.list` with the new relative paths is the only mechanical change needed.

### Step 5 — Move frozen Phase 7 product specs to `docs/specs/product/` (same amendment process)
Move the 10-file Phase 7 product definition package. Same process as Step 4. These are heavily cross-referenced by internal paths — update all cross-references before committing. `CONTRIBUTING.md` §4 lists some of these explicitly and must be updated.

### Step 6 — Move frozen Phase 8 engineering specs to `docs/specs/engineering/` (same amendment process)
Move the 5-file Phase 8 engineering requirements package. Same process. Verify that `02_SECURITY_OBJECTIVES.md`, `03_SYSTEM_INVARIANTS.md`, `04_FAILURE_MODEL.md` are not referenced by path in `rfc/` documents before moving.

### Step 7 — Narrow `context/` to AI memory only (no frozen docs involved)
Move the documentation-leak files out of `context/`:
- `context/01_GOVERNANCE.md` → `docs/engineering/` (governance rules already exist in `CONTRIBUTING.md` and `agents/agents/GOVERNANCE.md`; evaluate whether this file adds anything unique before moving vs. removing)
- `context/02_SYSTEM_ARCHITECTURE.md` → `docs/engineering/`
- `context/03_DATA_FLOW.md` → `docs/engineering/`
- `context/04_OPERATION_MANUAL.md` → `docs/engineering/` (or remove; its Make-target table duplicates `CONTRIBUTING.md` §2 exactly)
- `context/06_GLOSSARY.md` → consolidated into `agents/GLOSSARY.md` per Step 2
- `context/07_SECURITY_POLICY.md` → `docs/engineering/` or `docs/security/`

After this step, `context/` contains: `00_PROJECT_CONTEXT.md`, `05_DECISION_LOG.md`, `08_AI_HANDOFF.md`. Three files. Clear purpose. AI session memory only.

Update `CLAUDE.md` reading order to point to the new locations.

### Step 8 — Flatten `agents/agents/` to `agents/` (highest-risk, do last)
This is the highest-coordination step. Every internal cross-reference in `agents/agents/GOVERNANCE.md`, all journal entries, `CONTRIBUTING.md` §4, `CLAUDE.md`, `context/` files, `lab/LAB_README.md`, and `README.md` reference the inner path. A comprehensive path search and replace is required before committing. Do this as one atomic commit. Verify with `git grep agents/agents` before and after. This step should not be attempted while any other step is in-flight.

---

## 9. Risks

| Risk | Affected Steps | Severity | Mitigation |
|---|---|---|---|
| `make check-frozen` fails after moving frozen files if `frozen-docs.list` is not updated simultaneously | Steps 4, 5, 6 | High | Update `frozen-docs.list` paths and run `make frozen-baseline` in the same commit as the file move. Never split these. |
| CI workflow references a moved file by path | Steps 3–7 | Medium | `grep` the `.github/workflows/` directory for any hardcoded paths before each step. |
| `lab/LAB_README.md` references `agents/GOVERNANCE.md` (not double-nested path) — already partially correct | Step 8 | Low | Verify all intra-repo cross-references with `git grep agents/agents` before flattening. |
| `context/` files reference old frozen-doc paths after Steps 4–6 | Steps 4–6 | Medium | Update `context/08_AI_HANDOFF.md` and `context/00_PROJECT_CONTEXT.md` path tables in the same commit as each move. |
| `CONTRIBUTING.md` §4 names specific frozen files by filename | Steps 4–6 | Medium | Update `CONTRIBUTING.md` in the same commit. Not a frozen doc itself, so no baseline re-generation needed. |
| `agents/agents/GLOSSARY.md` provenance references old context paths | Step 2 | Low | Update the Provenance section of the consolidated glossary as part of the merge commit. |
| Two glossaries diverge further if Step 2 is deferred | Step 2 | Medium | Do Step 2 before any work that would add new terms to either file. |
| `DEVELOPMENT_RULES.md` RFC Policy still says "no new RFCs" — misleads future contributors | Step 1 | Low | Fix in Step 1 alongside the README stale-reference cleanup. No file moves required. |

---

## 10. Final Recommendation

**C — Moderate reorganisation recommended before implementation.**

The engineering foundation is sound. The core structural problems are:

1. **29 root Markdown files with no navigation hierarchy.** An engineer cannot find the entry points. Fifteen of these files are frozen planning specifications that belong in `docs/specs/`.

2. **`context/` overloaded with documentation.** Six of nine context files are technical documentation, not AI session memory. This will cause `context/` to become an unmanageable catch-all as the project grows.

3. **`agents/agents/` double-nesting.** A known historical accident. It confuses every directory listing and creates unnecessary path verbosity in all cross-references.

4. **Two glossaries with overlapping scope.** A drift problem waiting to compound.

5. **Stale references in `README.md` and `DEVELOPMENT_RULES.md`.** Minor, but signals the documentation already needs maintenance before engineering has begun.

None of these block the first engineering commit. But implementing features into a repository whose documentation structure is already confusing makes all future maintenance harder. The correct time to fix this is before the first `feat:` commit, not after.

The migration is achievable in one focused session for Steps 1–3 (low risk, no frozen-doc amendments) and two to three sessions for Steps 4–8 (higher coordination, frozen-doc amendment process required per batch).

---

*This document is the output of a read-only audit. No files were moved, renamed, deleted, archived, frozen, committed, or pushed during its production.*

<!-- checkpoint: feat(record): implement revstatus snapshot retrieval -->
