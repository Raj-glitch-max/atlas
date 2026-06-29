# Repository Audit

Date: 2026-07-05
Project: Atlas
Last updated: 2026-07-05 (hygiene completion pass — canonical memory + archival)

---

## KEEP — Active Governed Assets

| Path | Justification |
|---|---|
| `.github/` | Active repository automation: issue template, PR template, Dependabot, and repo-health workflow. |
| `.gitignore`, `.gitattributes`, `.editorconfig` | Active git/editor hygiene. |
| `.pre-commit-config.yaml` | Active local quality gate. |
| `.gitleaks.toml` | Active secret-scanning configuration referenced by `Makefile` and CI. |
| `.devcontainer/`, `Dockerfile`, `.vscode/` | Active contributor-environment support. |
| `Makefile` | Active local command surface. |
| `README.md`, `CONTRIBUTING.md`, `SECURITY.md`, `CLAUDE.md` | Active repository orientation and contributor guidance. |
| `PROJECT_STRUCTURE.md`, `DEVELOPMENT_RULES.md`, `REPOSITORY_AUDIT.md` | Canonical hygiene outputs. |
| `FROZEN.sha256`, `scripts/frozen-docs.list`, `scripts/check-frozen-docs.sh` | Active freeze integrity mechanism. |
| Frozen planning documents listed in `scripts/frozen-docs.list` | Hash-pinned planning and lab-governance corpus. |
| `lab/` | Active engineering research laboratory. |
| `context/` | Canonical project memory — nine files (00–08). Load before every AI session. |
| `agents/` | Active reasoning framework. |
| `tests/README.md` | Intentional placeholder for deferred stack-specific tests. |
| `docs/project/PROJECT_HISTORY.md` | Historical provenance. Old names and old absolute paths intentionally preserved. |
| `rfc/RFC-000-architecture-principles.md` | Active architecture RFC — Draft status. Tracked in git (2026-07-05). |
| `rfc/RFC-001-system-context.md` | Active architecture RFC — Draft status. Tracked in git (2026-07-05). |
| `rfc/RFC-002-conceptual-domain-model.md` | Active architecture RFC — Draft status. Tracked in git (2026-07-05). |

---

## ARCHIVE — Research Material

| Path | Action taken | Date |
|---|---|---|
| `archive/research/PRODUCT_THESIS.md` | Moved from root to `archive/research/`. Content summarized in `context/00_PROJECT_CONTEXT.md`. | 2026-07-05 |
| `archive/research/CATEGORY_AND_PRIMITIVE_ANALYSIS.md` | Moved from root to `archive/research/`. Content summarized in `context/00_PROJECT_CONTEXT.md`. | 2026-07-05 |
| `archive/research/INFRASTRUCTURE_PRIMITIVE_EVALUATION.md` | Moved from root to `archive/research/`. Content summarized in `context/00_PROJECT_CONTEXT.md`. | 2026-07-05 |
| `archive/research/ECOSYSTEM_THESIS.md` | Moved from root to `archive/research/`. Content summarized in `context/00_PROJECT_CONTEXT.md`. | 2026-07-05 |
| `archive/research/PORTFOLIO_REDUCTION_REVIEW.md` | Moved from root to `archive/research/`. Content summarized in `context/00_PROJECT_CONTEXT.md`. | 2026-07-05 |

Archive files are tracked in git under `archive/research/`. They are read-only; modification requires a journal entry explaining the reason.

---

## DELETED

| Path | Justification | Date |
|---|---|---|
| `claude-skills/` (63 MB) | Third-party reference repository. Audited in Session 8 (2026-06-19). Zero actual invocations. 7-item index preserved in `agents/REFERENCES.md`. Deleted per framework doctrine. | 2026-07-05 |
| `context/00_PROJECT_MISSION.md` | Superseded by `context/00_PROJECT_CONTEXT.md` in canonical memory consolidation. | 2026-07-05 |
| `context/01_CURRENT_STATE.md` | Superseded by `context/08_AI_HANDOFF.md` and `context/05_DECISION_LOG.md`. | 2026-07-05 |
| `context/02_SURVIVING_HYPOTHESES.md` | Superseded by `context/08_AI_HANDOFF.md`. | 2026-07-05 |
| `context/03_RESEARCH_BACKLOG.md` | Superseded by `context/08_AI_HANDOFF.md` and `context/00_PROJECT_CONTEXT.md`. | 2026-07-05 |
| `GLOSSARY.md` (root) | Content migrated to `context/06_GLOSSARY.md`. Root file removed to maintain a clean root directory. | 2026-07-05 |

---

## UNKNOWN / DEFERRED

| Path | Status |
|---|---|
| `agents/` nesting | The duplicate directory structure is historical. Flattening requires a controlled migration with journal entry — not a hygiene-only edit. Deferred. |
| `.claude/` | Ignored local agent/session state. Not a repository artifact. |

---

## Verification Notes (2026-07-05)

- `claude-skills/` deleted from disk and removed from `.gitignore`.
- Five root-level research documents moved to `archive/research/` and tracked in git.
- Four old `context/` files deleted; nine new canonical context files created.
- `rfc/` directory now tracked in git with three active RFC drafts.
- `CLAUDE.md` and `PROJECT_STRUCTURE.md` updated to reflect new layout.
- `make check-frozen` passed — no frozen documents were modified.

<!-- checkpoint: context(system-boundary-definition): clarify system boundary definition -->

<!-- checkpoint: feat(sdk): add cache invalidation -->

<!-- checkpoint: feat(issuance): add cache invalidation -->

<!-- checkpoint: chore(client): simplify viewport styling attributes -->
