# Repository Governance — Atlas

Governance regulates the interaction protocol between the founder and the agent system, the git workflows, the commit conventions, and the immutability rules for frozen documents.

## 1. Agent System Roles

The framework operates on a layered reasoning model structured to enforce maximum friction and skepticism:

### Five Council Seats (Epistemologies)
Council members represent distinct ways of analyzing a problem, not personal roles. They respond sequentially once per decision:
1.  **Empiricist** (`agents/council/empiricist.md`): Demands direct citations and independent evidence. Refuses unsubstantiated claims.
2.  **Cartographer** (`agents/council/cartographer.md`): Monitors boundaries, lens drift, and scope definitions.
3.  **Red Team** (`agents/council/red-team.md`): Identifies failure modes, attack surfaces, and fatal execution risks.
4.  **Economist** (`agents/council/economist.md`): Evaluates incentives, capture risks, ACV scale, and value capture.
5.  **Operator** (`agents/council/operator.md`): Assesses workflow integration, practitioner adoption friction, and GTM fit.

### Four Domain Anchors (Knowledge bases)
Domain anchors do not sit on the council. They are consulted only when a council member identifies a knowledge gap:
*   **Distributed Systems** (`agents/domain/distributed-systems.md`)
*   **AI/ML Systems** (`agents/domain/ai-ml-systems.md`)
*   **Product Engineering** (`agents/domain/product-engineering.md`)
*   **Market Buyer** (`agents/domain/market-buyer.md`)

### Working Specialists
Spawned for narrow, task-specific spikes. A working specialist that is cited in 10+ journal entries can be promoted to a domain anchor with founder opt-in. A domain anchor unconsulted for 90 days is demoted to a candidate.

## 2. Decision & Review Path

All decisions and technical designs follow a strict pipeline:
1.  **Draft:** Active working files under research or local experimentation.
2.  **Adversarially Reviewed:** Artifact run through sequential review by the Council. Dissent is recorded verbatim.
3.  **Accepted:** Founder makes a final decision. Overrides of council dissent are recorded with explanation.
4.  **Frozen:** Core plans and architectures are added to `scripts/frozen-docs.list` and hash-pinned in `FROZEN.sha256`.

## 3. Commit Standards & Conventions

The repository enforces Conventional Commits via git hooks configuration in `.pre-commit-config.yaml`.
*   **Format:** `<type>(<scope>): <subject>` (e.g., `feat(auth): add workload identity filter`).
*   **Types allowed:** `feat`, `fix`, `docs`, `chore`, `refactor`, `test`, `style`, `perf`, `build`, `ci`.
*   **Git Discipline:** Only stage files relevant to the specific logical commit.

## 4. Freeze & Immutability Policy

To protect planning integrity, 26 core documents are pinned in `scripts/frozen-docs.list` and monitored by `make check-frozen`.
To amend a frozen file:
1.  Record a journal entry in `agents/journal/` detailing the change rationale.
2.  Add a dated, numbered change note block to the target document.
3.  Run `make frozen-baseline` to regenerate `FROZEN.sha256`.
4.  Commit the updated file and the new `FROZEN.sha256` baseline.
5.  Silencing the check without following this process is a critical governance breach.

## 5. Journal & Override Protocol

*   **Verbatim Dissent:** Every decision journal entry must capture disagreements verbatim. No smoothing or flattening of conflict.
*   **Founder Override:** The founder can overrule a council recommendation or dissent. The entry must record:
    1.  Which council seat or anchor is overruled.
    2.  The founder's logic for the override.
    3.  The specific conditions or observations that would trigger a reconsideration of the override.

<!-- checkpoint: context(security-invariants): finalize security invariants -->

<!-- checkpoint: docs(revocation-requirements): update revocation requirements (#41) -->

<!-- checkpoint: repo(CI-testing-gates): finalize CI testing gates -->
