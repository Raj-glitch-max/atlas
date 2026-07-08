# Atlas

Atlas is a venture-studio technical workspace. Repository layout:

- **`agents/`** — persistent reasoning framework (governance, council, journal). See `agents/README.md`.
- **`context/`** — canonical project state (mission, current state, hypotheses, backlog).
- **`lab/`** — the engineering research laboratory: five process docs (`LAB_README.md`, `EXPERIMENT_LOG.md`, `EXPERIMENT_CHECKLIST.md`, `EVIDENCE_INDEX.md`, `DECISION_RULES.md`), the EXP-001 execution plan, and (future) `evidence/` runs.
- **Frozen documents** — hash-pinned planning and lab-governance files listed in `scripts/frozen-docs.list`; see `CONTRIBUTING.md` section "Frozen planning documents."
- **`docs/`** — planning specifications and project documentation. See `docs/product/`, `docs/engineering/`, `docs/research/`, `docs/project/`.

Agent instructions: `CLAUDE.md`. Project mission: `context/00_PROJECT_CONTEXT.md`. Filesystem rules: `PROJECT_STRUCTURE.md`. Development rules: `DEVELOPMENT_RULES.md`.

## Engineering foundation (Sprint 0)

This repo ships a **stack-agnostic** engineering foundation:

- git hygiene (`.gitignore`, `.gitattributes`, `.editorconfig`)
- pre-commit quality gates (`.pre-commit-config.yaml`) — trailing whitespace, end-of-file, YAML/JSON/TOML validity, large-file guard, private-key detection, markdown-aware, misspellings, Conventional Commits commit-msg enforcement
- CI repo-health pipeline (`.github/workflows/repo-health.yml`)
- secret scanning (`.gitleaks.toml`)
- commit-message conventions (Conventional Commits)
- dev-toolbox container (`Dockerfile`, `.devcontainer/devcontainer.json`)
- frozen-planning-docs integrity guard (`scripts/check-frozen-docs.sh` + `FROZEN.sha256`)
- repository docs and templates

**Stack-specific infrastructure is intentionally deferred** — no language, no test runner, no product container, no build pipeline — because selecting a stack is a decision that may depend on the C4 feasibility spike outcome. See `CONTRIBUTING.md` §"Sprint 0 scope."

### Quickstart

```sh
make init     # installs pre-commit + commit-msg hooks
make help     # list all targets
make ci       # full local CI equivalent
```

Toolchain: `git`, `python3` + `pre-commit`, `npx` (Node) for markdown lint, and optionally `docker` + `gitleaks` for `devshell`/`secrets`. CI provisions its own.

### The frozen-planning rule

The planning documents are frozen. `make check-frozen` verifies their SHA-256 hashes against `FROZEN.sha256` and is wired into CI. Editing a frozen doc without a journal entry and dated amendment breaks CI. See `CONTRIBUTING.md` §"Frozen planning documents."

<!-- checkpoint: test(verify): verify conformance against final test vectors -->

<!-- checkpoint: feat(ui): refine landing page workflows and styles -->
