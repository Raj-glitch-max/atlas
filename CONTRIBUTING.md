# Contributing

Atlas is a venture-studio technical workspace with a **frozen planning system** and an **engineering research lab**. These rules govern how work is added.

## 1. First-time setup

```sh
make init     # installs pre-commit + commit-msg hooks
make help     # list all targets
```

Toolchain (stack-agnostic): `git`, `python3` + `pre-commit`, `npx` (Node) for markdown lint, optionally `docker` and `gitleaks`. CI installs its own.

## 2. Day-to-day commands

| Target | Does |
|---|---|
| `make lint` | Run pre-commit hooks on all files |
| `make format` | Run auto-fixing hooks |
| `make docs-lint` | Markdown lint (markdownlint-cli2) |
| `make secrets` | Secret scan (gitleaks) |
| `make check-frozen` | Verify frozen planning docs unchanged (SHA-256 vs `FROZEN.sha256`) |
| `make test` | Test suite — **DEFERRED**: no stack selected yet (see §5) |
| `make ci` | Full local equivalent of CI |
| `make devshell` | Run the stack-agnostic dev-toolbox container |
| `make upgrade` | `pre-commit autoupdate` |

Stack-specific targets (`build`, a product-container `docker-build`, a real `test` runner, language linters) are added once a stack is selected — see §5.

## 3. Commits

Conventional Commits 1.0.0, enforced by the `commit-msg` hook:

```
feat(auth): add token verifier scaffold
fix(ci): correct pre-commit exclude
docs: update CONTRIBUTING
chore(sprint0): add dev-toolbox Dockerfile
```

Scope is optional; `!` marks a breaking change. Multi-line messages are fine; only the first line is enforced.

## 4. Frozen planning documents

The following are **frozen** and protected by `make check-frozen` (list in `scripts/frozen-docs.list`, baseline in `FROZEN.sha256`):

- `docs/research/FOUNDER_DECISION_BRIEF.md`, `docs/research/RESEARCH_PROGRAM.md`, `docs/research/TECHNICAL_VALIDATION.md`, `docs/research/FOUNDER_PROBLEM_FIT.md`, `P5_FALSIFICATION_EXPERIMENT.md`, `LEVEL0_1_FEASIBILITY_GATE.md`
- `docs/product/PRODUCT_DEFINITION.md` (and the nine cross-referenced specs in `docs/product/`)
- `docs/engineering/01_ENGINEERING_REQUIREMENTS.md` (and the four sibling specs in `docs/engineering/`)
- `lab/LAB_README.md`, `lab/EXPERIMENT_CHECKLIST.md`, `lab/EVIDENCE_INDEX.md`, `lab/DECISION_RULES.md`, `lab/EXP-001-EXECUTION-PLAN.md`

(`lab/EXPERIMENT_LOG.md` is append-only by design and is **not** in the frozen set — the log is meant to grow.)

To change a frozen document:

1. Record a **journal entry** at `agents/journal/<YYYY-MM-DD>-<slug>.md` per `agents/GOVERNANCE.md`.
2. Add a dated, reasoned change note to the document itself (lab docs: per `lab/LAB_README.md` §12).
3. Re-baseline: `make frozen-baseline`, commit `FROZEN.sha256` alongside the amendment.

A frozen doc changed in a PR without these steps fails `make check-frozen` and fails CI. Editing the hash baseline to silence the alarm is itself the violation being guarded against.

## 5. Sprint 0 scope — and what is deferred

Sprint 0 deliberately builds **stack-agnostic** infrastructure only, consistent with the directive not to make architectural decisions that depend on the C4 feasibility spike outcome.

Delivered now: git hygiene, pre-commit gates, CI repo-health pipeline, dev-toolbox container, secret scanning, commit conventions, frozen-docs guard, repo docs, PR/issue templates.

**Deferred until a stack is selected** (post-C4-spike, or on an explicit founder decision):

- a language/stack and its package manifest (`go.mod` / `pyproject.toml` / `package.json`…),
- that stack's test runner + a real `make test` target,
- that stack's linter/formatter (added to `.pre-commit-config.yaml`),
- a **product** container (the existing `Dockerfile` is a dev-toolbox, not a product image),
- build/test jobs in CI for product artifacts.

When the stack is chosen, add these in a single `feat(tooling): adopt <stack>` change and extend this document.

## 6. Lab boundary

Engineering work lives outside `lab/`. `lab/` is governed by `lab/LAB_README.md`; experiment runs follow `lab/EXPERIMENT_CHECKLIST.md` and produce evidence under `lab/evidence/`. Do not modify the lab process documents or frozen experiment specs as part of engineering work — see §4.

## 7. Pull requests

Use the PR template. A PR is mergeable when:

- `make ci` passes locally,
- `make check-frozen` passes,
- the Sprint 0 guardrail is honored (no product functionality, no spike-dependent architecture),
- any journal entry required by §4 is linked.

## 8. License

Not yet set. Adding a `LICENSE` is a recorded decision (proprietary vs OSS, and which license) deferred to the founder. Engineering infrastructure added before that decision carries no licensing claim.

<!-- checkpoint: context(API-path-design): document API path design (#7) -->

<!-- checkpoint: repo(deployment-manual): document deployment manual -->

<!-- checkpoint: docs(founder-profile-feedback): improve founder profile feedback -->

<!-- checkpoint: governance(fuzzing-strategy): audit fuzzing strategy -->

<!-- checkpoint: repo(glossary-definitions): refine glossary definitions (#40) -->

<!-- checkpoint: docs(revocation-requirements): finalize revocation requirements -->
