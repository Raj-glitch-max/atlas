# Operation Manual — Atlas

This document covers the day-to-day commands and runbooks for the Atlas repository. It does not prescribe product or architecture decisions — those are governed by the planning pipeline in `context/00_PROJECT_CONTEXT.md` and `context/01_GOVERNANCE.md`.

---

## 1. Prerequisites

| Tool | Required For | Install |
|---|---|---|
| `git` | All operations | system |
| `python3` + `pre-commit` | Hooks, linting, CI | `pip install pre-commit` |
| `node` + `npx` | Markdown lint | system |
| `gitleaks` | Secret scanning | optional (CI provisions it) |
| `docker` | Dev shell | optional |

---

## 2. First-Time Setup

```sh
# Install pre-commit and commit-msg hooks
make init
```

This installs all hooks defined in `.pre-commit-config.yaml` plus the `commit-msg` hook enforcing Conventional Commits.

---

## 3. Common Make Targets

| Command | What It Does |
|---|---|
| `make help` | List all available targets |
| `make init` | First-time setup — install git hooks |
| `make lint` | Run all pre-commit hooks on all files |
| `make format` | Run auto-fixing hooks (whitespace, line-endings) |
| `make docs-lint` | Markdown lint via `markdownlint-cli2` (requires `npx`) |
| `make secrets` | Scan for accidentally committed secrets via `gitleaks` |
| `make check-frozen` | Verify frozen planning doc hashes against `FROZEN.sha256` |
| `make frozen-baseline` | Regenerate `FROZEN.sha256` **after** a documented amendment |
| `make devshell` | Start the stack-agnostic dev-toolbox container |
| `make test` | Deferred — no stack selected yet (Sprint 0) |
| `make ci` | Full local CI equivalent: lint + docs-lint + check-frozen + secrets + test |
| `make upgrade` | Bump pre-commit hooks to latest versions |

---

## 4. Freeze Amendment Runbook

Use this procedure when a frozen document must change. Do not silence the check without following this process.

1. Create a journal entry at `agents/journal/<YYYY-MM-DD>-<slug>.md` documenting the rationale.
2. Add a dated, numbered change-note block to the document being amended.
3. Run `make frozen-baseline` to regenerate `FROZEN.sha256`.
4. Commit the amended document and the updated `FROZEN.sha256` together in a single commit.
5. Reference the journal entry in the commit message body.

---

## 5. Adding a Working Specialist

1. Copy `agents/templates/working-specialist-template.md` to `agents/working/<name>.md`.
2. Define the specialist's narrow scope, activation conditions, and retirement conditions.
3. Activate by naming the specialist in a founder prompt.
4. After 10+ journal citations with founder opt-in, a working specialist may be promoted to a domain anchor.
5. A domain anchor not consulted in 90 days is renamed `<name>.candidate.md` and must be reactivated by journal entry.

---

## 6. Writing a Journal Entry

Journal entries are the primary record of all committed decisions in Atlas.

```
agents/journal/<YYYY-MM-DD>-<slug>.md
```

Minimum required sections (from `agents/templates/journal-template.md`):
- **Date, Cycle, Context**
- **Decision** — precisely what was decided.
- **Dissent** — verbatim disagreements from any council member. Never smoothed.
- **Override (if applicable)** — what was overruled, by whom, and the exact reconsideration condition.
- **Confidence labels** — High / Medium / Low / None on every assertion.
- **Change condition** — what concrete observation would make this decision obsolete.

---

## 7. CI Pipeline

The CI workflow (`.github/workflows/repo-health.yml`) runs on push to `main` and on pull requests. It mirrors `make ci`:

1. `lint` — pre-commit hooks.
2. `docs-lint` — markdown lint.
3. `check-frozen` — SHA-256 integrity.
4. `secrets` — gitleaks scan.
5. `test` — deferred (no stack yet).
