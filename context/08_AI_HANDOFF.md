# AI Handoff — Atlas

This file is loaded by every new AI session. It tells you the current task state, what to do next, what not to do, and what you must preserve. Read this file before any other.

---

## 1. Where We Are Right Now

**Pipeline position:** The system is in a **dormant, governed state** after completing the R1 Problem Ranking Cycle.

| Stage | Status |
|---|---|
| Domain Research (10 problems identified) | ✅ Complete |
| External Repo Audit (claude-skills classified) | ✅ Complete |
| Problem Ranking Framework v1.0 established | ✅ Complete |
| R1 Problem Ranking Cycle (all 10 problems scored + council reviewed) | ✅ Complete |
| RFC Architecture (RFC-000, RFC-001, RFC-002 drafted) | ✅ In Draft (frozen pending founder review) |
| Repository Hygiene (atlas naming, canonical memory) | ✅ Complete (2026-07-05) |
| R2 Product Management Cycle | 🔴 NOT STARTED — blocked on founder selection |

**The next action is a founder decision**, not an autonomous one.

---

## 2. The Surviving Hypotheses (H1–H8)

These hypotheses are formally tracked and have not been invalidated. They constrain reasoning in all future sessions.

| # | Hypothesis | Status |
|---|---|---|
| H1 | External tools are classified, not auto-imported | Surviving — `claude-skills` was audited, index preserved, folder deleted |
| H2 | A problem-ranking gate must sit between Domain Research and Product Management | Surviving — gate was used in R1 |
| H3 | Founder advantage (A6) requires prior work evidence, not stated interest | Surviving — A6 = 0 uniformly in R1 due to no founder-history data |
| H4 | Evidence without citation cannot earn Evidence Quality ≥ 2 | Surviving — integrated as a system rule |
| H5 | Missing intake data is not fabricatable | Surviving — the Founder Profile session returned sparse intake with explicit gaps |
| H6 | System naming is content-neutral | Carried out — resolved via FounderOS → agents rename |
| H7 | Sequential review catches mis-scored dimensions before they calcify | Added 2026-06-19 in R1 — confirmed: Cartographer caught 2 numeric errors in R1 |
| H8 | A flat net-score rank over-states the gap between problems at the top and bottom | Added 2026-06-19 in R1 — remains as an open design question for R2 |

---

## 3. What the Next Session Should Do

The founder has not yet selected a candidate problem from the R1 shortlist. Until that selection is made:

### If the founder returns with a selection:
1. The founder picks one problem from the post-council R1 shortlist (source: `agents/journal/2026-06-19-r1-problem-ranking-application.md` §6).
2. Cycle R2 (Product Management) begins.
3. Write a journal entry for the R2 opening.

### If the founder wants a second-pass on R1:
- Backfill the Founder Profile (H5 surviving) to give A6 (Founder Advantage) a non-zero score in any re-ranking.

### If the founder wants to advance the RFC track:
- Review RFC-000, RFC-001, RFC-002 in `rfc/`.
- Only one new RFC may be authorized per founder instruction. RFC policy is in `DEVELOPMENT_RULES.md` §RFC Policy.
- The current three RFCs are in Draft status and require founder review to advance to Accepted/Frozen.

### If nothing is decided:
- Do not start Product Management work.
- Do not design products, startups, architectures, or features.
- Do not redesign the `agents/` framework.
- Maintain the current dormant state.

---

## 4. What You Must Never Do

These rules are always on, regardless of founder instruction:

- **Do not invent information.** Every factual claim must trace to a frozen RFC, journal entry, or existing repository file.
- **Do not speculate on product direction.** Product ideas, business models, and startup architectures are forbidden until R2 is explicitly opened.
- **Do not redesign `agents/`.** The framework is frozen at v1.0. Architecture is settled.
- **Do not import from `claude-skills/`.** The folder has been deleted. The 7-item index in `agents/REFERENCES.md` is the only reference.
- **Do not smooth dissent.** Journal entries preserve disagreement verbatim. Do not summarize conflict away.
- **Do not amend frozen docs without the amendment process.** Run `make check-frozen` to verify integrity.
- **Do not break working CI.** Run `make ci` locally before committing significant changes.

---

## 5. Files to Load at Session Start

In this order:

1. `context/00_PROJECT_CONTEXT.md` — what we are building and why.
2. `context/01_GOVERNANCE.md` — how decisions are made.
3. `context/08_AI_HANDOFF.md` (this file) — current state and next actions.
4. `agents/GOVERNANCE.md` — binding agent rules.
5. `agents/journal/2026-06-19-r1-problem-ranking-application.md` — R1 closure and shortlist.
6. `rfc/RFC-000-architecture-principles.md` — if doing RFC work.

---

## 6. Key File Locations

| Purpose | Path |
|---|---|
| Canonical project context | `context/00_PROJECT_CONTEXT.md` |
| Governance | `context/01_GOVERNANCE.md` |
| System architecture overview | `context/02_SYSTEM_ARCHITECTURE.md` |
| Data flows | `context/03_DATA_FLOW.md` |
| Make commands & runbooks | `context/04_OPERATION_MANUAL.md` |
| Key historical decisions | `context/05_DECISION_LOG.md` |
| Domain terminology | `context/06_GLOSSARY.md` |
| Security policy | `context/07_SECURITY_POLICY.md` |
| Agent governance rules | `agents/GOVERNANCE.md` |
| R1 journal entry | `agents/journal/2026-06-19-r1-problem-ranking-application.md` |
| Architecture principles | `rfc/RFC-000-architecture-principles.md` |
| System context RFC | `rfc/RFC-001-system-context.md` |
| Domain model RFC | `rfc/RFC-002-conceptual-domain-model.md` |
| Frozen doc list | `scripts/frozen-docs.list` |
| Frozen doc hashes | `FROZEN.sha256` |
| Project history | `PROJECT_HISTORY.md` |
| Archived research | `archive/research/` |

<!-- checkpoint: rfc(deployment-manual): document deployment manual -->

<!-- checkpoint: refactor(verify): refactor revocation status lookup -->

<!-- checkpoint: chore(internal): clean test assertions (#94) -->
