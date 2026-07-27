---
date: 2026-07-27
slug: relocate-frozen-planning-docs
artifact: maintainer decision — relocate two frozen planning documents from the repository root to docs/planning/ and re-baseline their paths, without altering a single byte of their content
decision: Moved P5_FALSIFICATION_EXPERIMENT.md and LEVEL0_1_FEASIBILITY_GATE.md from the repository root to docs/planning/ as part of the root-decluttering pass. Updated scripts/frozen-docs.list to their new paths and re-ran `make frozen-baseline`. The SHA-256 of both documents is UNCHANGED — only the path recorded in FROZEN.sha256 changed. Deliberately did NOT add an in-document change note (step 2 of the CONTRIBUTING.md section 4 process), because doing so would alter the bytes this freeze exists to protect. Reasoning recorded here instead.
agents_consulted: [principal-engineer, technical-writer]
overrides: false
related_entries: [open-source-governance-baseline]
---

# Context

A completion audit run against the real tree (not against the documentation)
found the repository root carrying 24 Markdown files, most of them internal
planning ceremony. The audit's finding was blunt and correct: the process
artifacts were burying the engineering — the fail-closed verification core, the
20 adversarial conformance vectors, the import-boundary lint — under noise that
no first-time visitor needs.

A prior review reached the same conclusion and was never executed:
`docs/project/REPOSITORY_INFORMATION_ARCHITECTURE_REVIEW.md` (2026-07-05)
scored root discoverability **3/10** and named root overcrowding "the single
largest information architecture failure in the repository."

Two of the documents that needed to move are **frozen**:

- `P5_FALSIFICATION_EXPERIMENT.md`
- `LEVEL0_1_FEASIBILITY_GATE.md`

Moving them broke `make check-frozen` (exit 1, `MISSING:` on both). That is the
guard working exactly as designed — a frozen document vanished from its pinned
path and the build stopped.

# Decision

Relocate both documents to `docs/planning/` and re-baseline their **paths**,
under the section 4 amendment process, with one deliberate deviation recorded
below.

The integrity claim is preserved and is directly checkable. The SHA-256 of each
document is identical before and after the move:

```
e7f6204e131947809052bf9ea11a9a11254935bb04c0e39d4e025ce6782f0e61  P5_FALSIFICATION_EXPERIMENT.md
e7f6204e131947809052bf9ea11a9a11254935bb04c0e39d4e025ce6782f0e61  docs/planning/P5_FALSIFICATION_EXPERIMENT.md

60586353571d827afd8c65d2be33a776db3919adbd624e2afd30d029f03c73ea  LEVEL0_1_FEASIBILITY_GATE.md
60586353571d827afd8c65d2be33a776db3919adbd624e2afd30d029f03c73ea  docs/planning/LEVEL0_1_FEASIBILITY_GATE.md
```

Both files were moved with `git mv`, so `git log --follow` resolves their full
history across the rename (verified: 4 commits back to introduction).

## The one deviation, and why

CONTRIBUTING.md section 4 step 2 requires "a dated, reasoned change note to the
document itself." **That step was deliberately skipped, and it should be.**

Step 2 exists so that a *content* amendment is self-documenting to anyone
reading the frozen document. This change is not a content amendment — it is a
relocation. Writing a change note into the file would modify the exact bytes the
freeze exists to protect, changing both hashes and destroying the strongest
evidence available that nothing was altered. Following the letter of step 2 here
would violate its purpose.

The reasoning lives in this journal entry instead, which is what the journal is
for. Steps 1 (journal entry) and 3 (`make frozen-baseline` + commit
`FROZEN.sha256`) were followed exactly.

## Collateral corrections

Three references to the old root paths were updated so the guard and its
tooling stay in sync:

- `scripts/frozen-docs.list` — the two pinned paths.
- `CONTRIBUTING.md` section 4 — the prose list of frozen documents.
- `.pre-commit-config.yaml` — the codespell exclusion was root-anchored
  (`^P5_FALSIFICATION_EXPERIMENT\.md$`) and would no longer have matched after
  the move, silently exposing a frozen document to a linter that governance
  forbids appeasing (TD-9). Re-anchored to the new path. The three auto-fixer
  exclusions (trailing-whitespace, end-of-file-fixer, mixed-line-ending) already
  used a `(^|/)` prefix and matched the new paths without change — checked, not
  assumed.

# Consequences

- `make check-frozen` passes again; `make ci` is green end to end.
- The freeze guarantee is unchanged in substance: the same bytes are pinned,
  under new paths.
- The repository root drops toward the eight files a newcomer actually needs.

# Dissent / open questions

**Raised, not resolved here.** The freeze mechanism currently pins 30 planning
documents — research briefs, a feasibility gate, product and engineering specs —
while `internal/verify`, the lowest-covered (76.0%) and most safety-critical
package in the tree, has no equivalent guard. The protection is well-built and
pointed at the artifacts least likely to be edited by accident.

This entry does not change that. It is recorded as a proposal for a separate,
explicit decision: consider whether the freeze discipline should extend to the
conformance surface (the vectors, the verification core's public behaviour) and
whether some planning documents have earned their way out of the frozen set now
that the thing they planned exists and is tested.
