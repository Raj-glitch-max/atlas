# Journal — Freeze Phase 7 Product Definition & Phase 8 Engineering Requirements

| | |
|---|---|
| Date | 2026-07-04 |
| Mode | Direct (founder instruction); no council review requested |
| Decision | Add the Phase 7 Product Definition package (10 docs) and the Phase 8 Engineering Requirements set (5 docs) to the hash-pinned frozen set (`scripts/frozen-docs.list` + `FROZEN.sha256`) and commit them as the initial git history. |
| Source authority | Founder instruction, 2026-07-04: "Freeze and commit the approved Phase 7 and Phase 8 documents according to the project's governance." |
| Procedure followed | `CONTRIBUTING.md` §4 (journal entry → re-baseline → commit `FROZEN.sha256`). Step 2 of §4 (dated change note *in the amended document*) addresses amendments of already-frozen docs; this act *adds* docs to the frozen set at their current approved content, so no per-document change note is required (see §5). |
| Active milestone | End of phase-based planning; transition to Architecture RFCs. |
| Stop boundary | After commit of the frozen set and this journal entry. No Phase 9. No architecture, technology, repository, API, or protocol in this act. |

---

## 1. What is being frozen

### 1.1 Phase 7 — Product Definition package (10 documents)

The authoritative scope of "the Product Definition package" is `PRODUCT_DEFINITION.md` line 3, which enumerates the package members verbatim: itself plus the nine documents cross-referenced from it. These were approved and "frozen by founder statement" during Phase 7 (the provenance phrase carried in every Phase 8 doc). This act formalizes that freeze through the integrity-guard mechanism (`scripts/check-frozen-docs.sh` against `FROZEN.sha256`), making the freeze machine-enforced rather than merely declared.

| # | Path |
|---|---|
| 1 | `PRODUCT_DEFINITION.md` |
| 2 | `SYSTEM_CONTEXT.md` |
| 3 | `USER_MODEL.md` |
| 4 | `USE_CASE_CATALOG.md` |
| 5 | `FUNCTIONAL_REQUIREMENTS.md` |
| 6 | `NON_FUNCTIONAL_REQUIREMENTS.md` |
| 7 | `CONSTRAINTS.md` |
| 8 | `ASSUMPTIONS_AND_RISKS.md` |
| 9 | `V1_SCOPE.md` |
| 10 | `DEFERRED.md` |

### 1.2 Phase 8 — Engineering Requirements (5 documents)

| # | Path | Contents |
|---|---|---|
| 11 | `01_ENGINEERING_REQUIREMENTS.md` | ER1–ER17 (ER11–ER14 `[HYPOTHESIS]`) |
| 12 | `02_SECURITY_OBJECTIVES.md` | SO1–SO8 (SO4 `[HYPOTHESIS]`) |
| 13 | `03_SYSTEM_INVARIANTS.md` | INV1–INV12 + C-INV1 (candidate) |
| 14 | `04_FAILURE_MODEL.md` | FM1–FM11 |
| 15 | `05_ACCEPTANCE_TEST_PLAN.md` | AT1–AT30 |

New total in the frozen set: 11 (prior) + 15 = **26** documents.

---

## 2. What is deliberately NOT being frozen (and why)

The following root documents exist on disk and were considered, then deliberately excluded. Each exclusion is traceable, not silent.

- **`PRODUCT_THESIS.md`, `ECOSYSTEM_THESIS.md`, `GLOSSARY.md`, `PROJECT_HISTORY.md`, `PORTFOLIO_REDUCTION_REVIEW.md`** — referenced *by* `PRODUCT_DEFINITION.md` as prior-phase source material (theses, history, glossary), not enumerated *as members of* the Product Definition package. `PRODUCT_DEFINITION.md` line 3 lists exactly the nine (plus itself). Freezing the theses would be a Phase 5/6 retrospective freeze — a separate founder decision, not implied by "freeze Phase 7." Left unfrozen; left untracked in this commit.
- **`CLAUDE.md`, `CONTRIBUTING.md`, `README.md`, `SECURITY.md`** — workspace infrastructure and governance documents, not planning documents. The guard's own header declares its scope as "Planning documents" (`scripts/frozen-docs.list` line 1). Out of scope by type.
- **`context/` (00–03), `agents/` framework, `agents/agents/journal/`** — live working state and the reasoning framework, explicitly intended to evolve. Not freezable planning artifacts.
- **`lab/EXPERIMENT_LOG.md`** — a live experiment record (append-only), not a static plan/rule/index. The other five lab docs (plans, checklists, indexes, decision rules) are already frozen and stay frozen; the log is correctly excluded.

---

## 3. Integrity guard state (before → after)

- **Before:** `make check-frozen` → `OK`, 11 files. `FROZEN.sha256` held 11 baselines.
- **After (target):** rewrite `scripts/frozen-docs.list` to 26 paths; run `make frozen-baseline` to regenerate `FROZEN.sha256` from those 26 files' actual bytes; run `make check-frozen` → must report `OK`, 26 files.

The re-baseline is the *authorized* mechanism for changing the frozen set (`CONTRIBUTING.md` §4; Makefile `frozen-baseline` target). Editing `FROZEN.sha256` by hand to silence the guard is the violation being guarded against; this act does not do that — it changes the *list* (the authorized input) and lets the make target regenerate the baseline from the listed files' actual bytes.

---

## 4. Commit shape

This is the initial git commit (the repo previously had zero commits, 0 tracked files). Because zero commits existed, the freeze is not enforceable in git until a commit exists that contains the guard, the list, and the baseline. Commit #1 therefore includes:

- The Sprint-0 guard infrastructure: `Makefile`, `scripts/check-frozen-docs.sh`, `scripts/frozen-docs.list` (expanded to 26), `FROZEN.sha256` (regenerated), `.gitattributes`, `.gitignore`;
- The CI/lint config the committed `Makefile` invokes (so its `lint`/`secrets`/`ci` targets are non-dangling in history): `.pre-commit-config.yaml`, `.gitleaks.toml`, `.github/` (incl. `workflows/repo-health.yml`, `dependabot.yml`, `pull_request_template.md`, `ISSUE_TEMPLATE/`);
- The 26 frozen planning documents themselves (so the baselines they are pinned against are present in the same commit);
- `CONTRIBUTING.md` — the human procedure doc the guard's error message and §4 reference; committing it keeps the guard self-consistent in history (no dangling `see CONTRIBUTING.md §4`);
- This journal entry.

Deliberately excluded from the commit (left untracked, on disk, for a follow-up "workspace docs / dev-env" commit when the founder authorizes): `README.md`, `SECURITY.md`, `Dockerfile`, `.devcontainer/`, `.editorconfig`, `.vscode/` (workspace docs and dev-environment, not part of the freeze or its enforcement); `.claude/` (machine-local session state, gitignored); `claude-skills/` (third-party reference repo carrying its own `.git` — now gitignored so it can never become a broken gitlink); and the prior-phase / supporting root docs and live working state listed in §2. This scoping follows the project doctrine ("Files are changed explicitly, not silently") and the rule that only the work the founder authorized is committed.

Conventional Commits message: `chore(docs): freeze phase 7 product definition and phase 8 engineering requirements`.

---

## 5. Confidence (per GOVERNANCE §5)

- **High** that the 10 Phase 7 docs in §1.1 are exactly the Product Definition package — direct verbatim reading of `PRODUCT_DEFINITION.md` line 3, the authoritative self-declaration. *What would change it:* an amendment to `PRODUCT_DEFINITION.md` altering package membership (itself a frozen-doc amendment act).
- **High** that the 5 Phase 8 docs in §1.2 are the complete Phase 8 set — they were the explicit deliverable of the founder's Phase 8 instruction ("Produce ONLY these five documents") and no others were produced. *What would change it:* discovery of a sixth Phase 8 doc not in the set.
- **High** that the exclusions in §2 are correct — each traces to the package's own scope declaration (line 3), the guard's stated type-scope ("planning documents"), or the live/working nature of the artifact. *What would change it:* a founder ruling that a named excluded doc is in fact part of Phase 7/8.
- **Medium** that no per-document "frozen on \<date\>" change note is required for newly-frozen docs. Rationale: `CONTRIBUTING.md` §4 step 2 addresses *amendment* of an already-frozen doc; adding a doc pins its current approved bytes with no content change to annotate. *What would change it:* a stricter reading of §4 requiring an in-doc freeze annotation; if the founder prefers that, a follow-up amendment act adds them.

---

## 6. Open questions (deliberately not resolved here)

- Should `PRODUCT_THESIS.md`, `ECOSYSTEM_THESIS.md`, and other prior-phase root docs be retrospectively frozen? They predate the Sprint-0 frozen set and are referenced as Phase 7 inputs. A separate freeze act could cover them. Not decided here.
- Should a "frozen on \<date\>" annotation be added inside each newly-frozen document (see §5, confidence item 4)? Not decided here.
- Should the `agents/agents/` nesting be reconciled with `CLAUDE.md`? `CLAUDE.md` states the journal path as `agents/journal/`; the actual path is `agents/agents/journal/`. Documentation discrepancy, not a freeze question. Flagged for a future explicit edit; not fixed silently here (`CLAUDE.md`: "Files are changed explicitly, not silently").

---

## 7. Stop boundary

This entry records the freeze act only. It does not authorize Phase 9, architecture, technology, repositories, APIs, or protocols. The phase-based planning workflow ends here per founder instruction; the project transitions to Architecture RFCs, beginning with RFC-000 (recorded in a separate journal entry).

---

## Provenance trail

- Authority: founder instruction, 2026-07-04 (verbatim in the source-authority row above).
- Scope-of-package source: `PRODUCT_DEFINITION.md` line 3 (self-declaration of the Product Definition package membership).
- Procedure: `CONTRIBUTING.md` §4; Makefile `frozen-baseline` / `check-frozen` targets; `scripts/check-frozen-docs.sh`.
- Prior frozen set: `FROZEN.sha256` (11 baselines) and `scripts/frozen-docs.list` as they stood before this act.
- This journal entry is the durable artifact for the freeze act. It is the source of truth for any future audit of why the frozen set expanded from 11 to 26.

<!-- checkpoint: repo(trust-anchors): improve trust anchors (#15) -->
