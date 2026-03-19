# Development Rules

Atlas reuses the governance already present in `CONTRIBUTING.md`, `.github/`, `Makefile`, `scripts/frozen-docs.list`, and `agents/GOVERNANCE.md`. This file summarizes those active rules without adding product or architecture policy.

## Branch Strategy

- `main` is the current integration branch.
- The repo-health workflow runs on pushes to `main` and on pull requests.
- Keep branches short-lived and scoped to one reviewable change.
- Do not change branches during repository hygiene unless explicitly asked.

## Commit Strategy

- Use Conventional Commits.
- Keep the first line terse and scoped when useful, for example `chore(repo): normalize repository structure and governance`.
- Stage explicit paths only. Do not use broad staging when unrelated or pre-existing work is present.
- Commit frozen-doc baseline changes only with the corresponding allowed frozen-doc amendment.

## Conventional Commits

The `commit-msg` hook enforces Conventional Commits via `.pre-commit-config.yaml`.

Allowed examples from existing governance:

```text
feat(auth): add token verifier scaffold
fix(ci): correct pre-commit exclude
docs: update CONTRIBUTING
chore(sprint0): add dev-toolbox Dockerfile
```

## Review Policy

Use `.github/pull_request_template.md`. A pull request is mergeable when:

- `make ci` passes locally or the failure is explicitly documented;
- `make check-frozen` passes;
- Sprint 0 guardrails are honored;
- required journal entries are linked.

## Freeze Policy

Frozen files are listed in `scripts/frozen-docs.list` and checked against `FROZEN.sha256` by `make check-frozen`.

To change a frozen file, follow `CONTRIBUTING.md` section 4:

1. Record a journal entry at `agents/journal/<YYYY-MM-DD>-<slug>.md`.
2. Add the required dated change note to the changed document.
3. Run `make frozen-baseline`.
4. Commit `FROZEN.sha256` with the amendment.

Editing `FROZEN.sha256` only to silence the integrity check is a violation.

## RFC Policy

- Do not create new RFCs under current repository governance.
- Existing files under `rfc/` are audit/archive material until explicitly reclassified.
- Use the existing journal and freeze process for governed decisions.

## Journal Policy

- Decisions, dissent, overrides, and required frozen-doc changes are recorded under `agents/journal/`.
- Preserve disagreement instead of smoothing it into consensus.
- Founder overrides are allowed, but must be recorded with the reasoning and reconsideration conditions described in `agents/GOVERNANCE.md`.

## Experiment Policy

- `lab/` runs pre-registered experiments and records evidence.
- The lab does not implement product code, design product architecture, choose technologies, or start product-management work.
- Experiment logs are append-only. Static lab process docs are frozen through `scripts/frozen-docs.list`.
- Evidence artifacts follow `lab/LAB_README.md` and `lab/EVIDENCE_INDEX.md`.
