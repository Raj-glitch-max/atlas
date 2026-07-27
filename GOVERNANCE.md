# Governance

Atlas is an open-source project under the Apache-2.0 license. This document
describes how decisions are made. It is intentionally lightweight — appropriate
for a project at the CNCF Sandbox stage — and will grow with the contributor base.

## Principles

Atlas is governed by two standing disciplines that predate this file and outrank
convenience:

- **Fail-closed on the security floor.** Changes to `internal/record`,
  `internal/verify`, or `internal/issuance` that could weaken verification are
  held to the highest bar and require an explicit, recorded decision. Ergonomic
  sugar lives in the SDK/CLI layer, never in the verification core.
- **Judgment is recorded, not remembered.** Every architecturally significant or
  irreversible decision produces an entry in
  [`docs/planning/ENGINEERING_DECISION_RECORD.md`](docs/planning/ENGINEERING_DECISION_RECORD.md) (ADRs) and,
  where it concerns the reasoning framework or a freeze act, a dated journal
  entry under [`agents/journal/`](agents/journal/). Planning documents listed in
  [`scripts/frozen-docs.list`](scripts/frozen-docs.list) are hash-pinned; amending
  one follows the process in [`CONTRIBUTING.md`](CONTRIBUTING.md) §4.

## Roles

- **Users** — anyone running Atlas. No obligations; feedback and issues welcome.
- **Contributors** — anyone who opens a pull request or issue. Governed by
  [`CONTRIBUTING.md`](CONTRIBUTING.md) and the [Code of Conduct](CODE_OF_CONDUCT.md).
- **Maintainers** — listed in [`MAINTAINERS.md`](MAINTAINERS.md). Maintainers
  review and merge changes, cut releases, triage security reports, and own the
  technical direction.

## Decision-making

- **Lazy consensus.** Most changes proceed by lazy consensus: a pull request with
  at least one maintainer approval and no unresolved maintainer objection may be
  merged. Trivial changes (docs, tests, typos) need one approval.
- **Verification-core changes** require explicit maintainer approval and an ADR
  entry when they change behavior, an interface, or an invariant.
- **Disagreement.** If maintainers disagree, the change waits. A conflict that
  cannot be resolved by discussion is decided by the lead maintainer, with the
  losing position recorded in the ADR or the PR thread — the same discipline the
  project's own review process uses (`agents/review-team/output/atlas-conflict-log.md`).

## Maintainer lifecycle

- **Adding a maintainer.** A contributor with a sustained track record (guideline:
  three non-trivial, independently reviewed contributions) may be nominated by any
  maintainer. Approval is by consensus of existing maintainers. The change that
  adds them updates `MAINTAINERS.md` and `.github/CODEOWNERS` together.
- **Removing a maintainer.** A maintainer may step down at any time by PR. Inactive
  maintainers (guideline: no activity for six months) may be moved to emeritus by
  consensus of the remaining maintainers.

## Changing this document

Governance changes are proposed by pull request and require consensus of all
current maintainers.

## Code of Conduct

All participation is governed by the [Code of Conduct](CODE_OF_CONDUCT.md).
Violations are handled by the maintainers per the enforcement process therein.
