# Decision Log — Atlas

This file records key committed decisions and their outcomes. Every entry traces to a journal file, a frozen document, or a verifiable repository state. Speculation is forbidden.

---

## Decision 1 — Framework Naming: FounderOS → agents

**Date:** 2026-06-19
**Source:** Session 10 — Maintenance Session (`PROJECT_HISTORY.md` §3).
**Decision:** The framework directory was renamed from `/founderos/` to `/agents/`. The brand label "FounderOS" was normalized to "agents" in three documents. The name "FounderOS" is preserved only in historical context files to maintain provenance.
**Rationale:** The path `/home/raj/Videos/projects/founderos/agents/` duplicated the brand unnecessarily.
**Effect:** All active files now reference `agents/` or `Atlas`. No technical content changed.

---

## Decision 2 — External Repository Audit (claude-skills)

**Date:** 2026-06-19
**Source:** Session 8 — External acquisition audit (`PROJECT_HISTORY.md` §3, Session 8).
**Decision:** The `claude-skills/` repository (60+ MB, 345+ skills, 17 domains) was audited. A 7-item subset was classified and indexed in `agents/REFERENCES.md`. No skills were imported into the framework. The folder was deleted from disk in the Atlas repository hygiene session (2026-07-05).
**Rationale:** Zero actual invocations across all sessions. The framework's own doctrine mandates dropping non-used material after the reference index is preserved.
**Effect:** `claude-skills/` removed from disk and from `.gitignore`. The 7-item index remains in `agents/REFERENCES.md`.

---

## Decision 3 — Problem Ranking Framework Established

**Date:** 2026-06-19
**Source:** Session 9 — Problem Ranking Framework (`PROJECT_HISTORY.md` §3, Session 9).
**Decision:** A 10-dimensional opportunity/friction scoring framework was established as a mandatory gate between Domain Research and Product Management.
**Structure:**
- 3 validity-gate axes (Evidence quality ≥ 2, Frequency ≥ 1, Severity ≥ 1).
- 3 opportunity-fitness axes (Tooling under-service, Market under-saturation, Founder advantage). Higher = better.
- 4 friction-penalty axes (Technical friction, Implementation complexity, Buyer friction, Seller friction). Higher = worse.
- Net score = opportunity − friction (range −12 → +18).
**Effect:** Framework is frozen at v1.0. Any scoring cycle must use this framework.

---

## Decision 4 — R1 Problem Ranking Cycle Completed

**Date:** 2026-06-19
**Source:** `agents/journal/2026-06-19-r1-problem-ranking-application.md`.
**Decision:** All 10 Domain Research problems were scored, and sequential council review (Empiricist → Cartographer → Red Team → Economist → Operator) was completed. Founder selection was explicitly deferred.
**Outcome of council review:**
- Cartographer revised P5 (IAM): A5 score 1→2 due to SPIFFE standardization opening competitive room. Net: 4→5.
- Cartographer revised P7 (K8s): A4 score 1→0, because the pain is a value-promise gap, not missing tooling. Net: 1→0.
- No problems were eliminated by any council seat.
- All 10 problems survived with failure-mode and risk annotations applied.
**Final ranked shortlist (post-council top 3):** P9 On-call/toil (6), P4 Cloud cost (6), P2 CI/CD (6).
**Effect:** Cycle R1 is closed. Cycle R2 (Product Management) is blocked until founder makes a selection.

---

## Decision 5 — RFC Architecture Initiation

**Date:** 2026-07-04
**Source:** Founder instruction 2026-07-04; RFC-000, RFC-001, RFC-002.
**Decision:** A formal RFC process was authorized to define the architecture of the workload delegation system identified through the research pipeline (P5 — IAM/service identity). Three RFCs were created:
- `RFC-000`: Architecture Principles (the constitution — architectural properties AP1–AP13, trade-off philosophy TP1–TP6, decision rules DR1–DR10).
- `RFC-001`: System Context (system boundary, trust boundaries, external actors, responsibilities).
- `RFC-002`: Conceptual Domain Model (domain concepts, relationships, state transitions, invariants).
**Status:** All three are in Draft status. RFC policy per `DEVELOPMENT_RULES.md` §RFC Policy prohibits new RFCs until existing ones advance.
**Effect:** `rfc/` directory is now tracked in git.

---

## Decision 6 — Repository Name Normalized to Atlas

**Date:** 2026-07-05
**Source:** Founder instruction, Atlas hygiene session.
**Decision:** The local repository directory was renamed from `projects` to `atlas`. All active repository references now use the canonical project name "Atlas". Historical documents preserve old names where they refer to project history.
**Effect:** Path updated in `CLAUDE.md`, `PROJECT_STRUCTURE.md`, `README.md`, and all context files.

---

## Decision 7 — Research Documents Archived

**Date:** 2026-07-05
**Source:** Atlas hygiene session.
**Decision:** Five loose root-level research documents were moved to `archive/research/`:
- `PRODUCT_THESIS.md`
- `CATEGORY_AND_PRIMITIVE_ANALYSIS.md`
- `INFRASTRUCTURE_PRIMITIVE_EVALUATION.md`
- `ECOSYSTEM_THESIS.md`
- `PORTFOLIO_REDUCTION_REVIEW.md`
**Rationale:** Their content is summarized in `context/00_PROJECT_CONTEXT.md`. The root directory should contain only active governance and orientation files.
**Effect:** Root is clean. Archive files are tracked in git under `archive/research/`.
