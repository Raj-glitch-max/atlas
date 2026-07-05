# Repository Audit

Date: 2026-07-05
Project: Atlas

Scope: repository hygiene only. No files were deleted during this audit.

## KEEP

| Path | Justification |
|---|---|
| `.github/` | Active repository automation: issue template, PR template, Dependabot, and repo-health workflow. |
| `.gitignore`, `.gitattributes`, `.editorconfig` | Active git/editor hygiene. Naming and scope are consistent with the stack-agnostic foundation. |
| `.pre-commit-config.yaml` | Active local quality gate. Kept after aligning auto-fixer exclusions with `scripts/frozen-docs.list`. |
| `.gitleaks.toml` | Active secret-scanning configuration referenced by `Makefile` and CI. |
| `.devcontainer/`, `Dockerfile`, `.vscode/` | Active contributor-environment support. The container is a dev-toolbox, not product infrastructure. |
| `Makefile` | Active local command surface for setup, lint, docs lint, frozen-doc checks, secrets, and CI-equivalent checks. |
| `README.md`, `CONTRIBUTING.md`, `SECURITY.md`, `CLAUDE.md` | Active repository orientation and contributor guidance. |
| `PROJECT_STRUCTURE.md`, `DEVELOPMENT_RULES.md`, `REPOSITORY_AUDIT.md` | Canonical hygiene outputs from this pass. |
| `FROZEN.sha256`, `scripts/frozen-docs.list`, `scripts/check-frozen-docs.sh` | Active freeze integrity mechanism. |
| Frozen planning documents listed in `scripts/frozen-docs.list` | Hash-pinned planning and lab-governance corpus. Keep unchanged unless the freeze policy is followed. |
| `lab/` | Active engineering research laboratory. `lab/EXPERIMENT_LOG.md` is append-only; static lab process docs are frozen. |
| `context/` | Active project-state context. Keep as durable orientation, not generated scratch. |
| `agents/agents/` | Active reasoning framework. The duplicate nesting is historical and should be normalized only through a separate controlled path migration. |
| `tests/README.md` | Intentional placeholder documenting deferred stack-specific tests. |
| `PROJECT_HISTORY.md` | Historical provenance. Old names and old absolute paths are intentionally preserved there. |
| `GLOSSARY.md` | Active terminology document. Historical references to FounderOS are intentional. |

## ARCHIVE

| Path | Justification |
|---|---|
| `rfc/` | Existing RFC-style planning outputs. This pass does not create or extend RFCs; keep for history unless a later explicit archive move is approved. |
| `CATEGORY_AND_PRIMITIVE_ANALYSIS.md` | Large analysis output not listed in the frozen set. Keep available, but treat as archival/research material until ownership is clarified. |
| `INFRASTRUCTURE_PRIMITIVE_EVALUATION.md` | Research/evaluation output outside the frozen set. Keep available; do not treat as active governance without an explicit decision. |
| `ECOSYSTEM_THESIS.md`, `PRODUCT_THESIS.md`, `PORTFOLIO_REDUCTION_REVIEW.md` | Planning/thesis outputs outside the frozen set. Keep available as archive candidates; do not delete automatically. |
| `claude-skills/` | Third-party acquired reference repository with its own `.git`. Correctly ignored by this repo; keep as local reference-only material or archive outside the repo later. |

## DELETE

| Path | Justification |
|---|---|
| None | No path was confirmed as safe to delete during this pass. Deletion requires explicit follow-up approval. |

## UNKNOWN

| Path | Justification |
|---|---|
| `agents/agents/` flattening | The duplicate directory name is real, but many active and historical paths reference it. Flattening would be a controlled migration, not a hygiene-only edit. |
| `.claude/` | Ignored local agent/session state. It is not a repository artifact; local deletion is optional and was not performed. |
| Untracked research documents outside the frozen list | Several markdown files appear intentionally authored but are not tracked in the current commit. Their long-term home should be decided before moving or deleting them. |

## Verification Notes

- Duplicate folders found: `agents/agents/`; duplicate-looking `claude-skills/research/research/` exists inside ignored third-party reference material.
- Obsolete or temporary folders found: `.claude/` local state; `rfc/` existing RFC-style outputs under the current no-new-RFC constraint.
- Generated/cache folders found outside ignored reference material: none.
- Empty directories found outside `.git/` and ignored reference material: none.
- Nested repositories found: `claude-skills/.git`, already ignored and documented as reference-only.
- Broken symlinks found outside ignored reference material: none.
- Merge conflict markers found outside ignored reference material: none.
