# Glossary

Project-specific definitions for terms that are precise and load-bearing in `agents/`. Loaded by future sessions to keep terms stable across the team. **Not** an orientation document — see `README.md` — and **not** a rules document — see `GOVERNANCE.md`.

Each entry: project meaning, common confusion to avoid, and the file(s) where the term is enforced.

---

## A

### agents

The current brand name of this framework at `/home/raj/Videos/projects/agents/`. Renamed from "FounderOS" on 2026-06-19. The old name appears in narrative history and changelog-style mentions only — never as the active brand.

**Common confusion:** "agents" here does not mean *autonomous runtime agents*. The system is a set of viewpoint files you load into Claude Code. There is one conversation; named viewpoints are invoked by the founder. Do not conflate with multi-agent runtimes.

---

### anchoring (any X)

The frame that pins a project state item: a transcript index, a file path, or a context-file section. Used so future sessions regenerate the same state without diff. Anchors are not sources of truth themselves; they are **pointers** to sources of truth.

---

## C

### Cartographer

One of five council seats; frame + claim drift discipline. Inspects whether evidence patterns are being read from the wrong angle, whether a problem is scored against the wrong dimension, or whether a vocabulary drift has silently changed meaning mid-discussion. Promotion/demotion of working specialists is Cartographer's call after a quorum of journal entries.

---

### classification (skill)

A label applied to an external skill in the long-term `REFERENCES.md` index. Four values: `Core`, `Temporary`, `Reference`, `Remove`. Promotion requires journal entries demonstrating load-bearing use. Demotion is fast and ceremonial.

**Common confusion:** "Temporary" does not mean "shoddy." It means *opt-in per session, consulted by name*.

---

### confidence label

A calibrated marker expressed by every council and anchor output: `High` (directly evidenced or formally proved within scope), `Medium` (supported by analogy, no direct precedent), `Low` (speculative, with the change condition stated), `None` (refusing to estimate). **Confidence without evidence is forbidden**. Confidence without a stated change-condition is forbidden.

Enforced in `GOVERNANCE.md` §5.

---

### consensus smoothing

The anti-pattern of reshaping a multi-agent review so it converges in apparent agreement without the underlying objections being addressed. Banned explicitly. Disagreement stays visible. "Everyone agrees too quickly" is itself a signal that important objections are missing.

Enforced in `GOVERNANCE.md` §1.

---

### Core

A skill classification (see *classification*): integrated permanently into `agents/`. Currently **no** items classified Core. Promotion requires founder opt-in + ≥ 10 journal entries citing the skill as load-bearing.

---

### council

The five permanent epistemologies — `Empiricist`, `Cartographer`, `Red Team`, `Economist`, `Operator` — which together produce adversarial review of founder artifacts. Housed in `council/`. Council members are *epistemologies*, not roles.

**Common confusion:** council members do **not** vote. They speak in order, in their own voice, then the founder decides.

---

### cross-examination

Opt-in mode, founder-triggered, for high-stakes decisions. Any council member may challenge another council member's *evidence basis* by name. The challenged member must answer in the same session, citing the prior their analogy rests on. If they cannot, the claim is struck from the artifact with a deletion reason recorded.

**Common confusion:** not debate theater. Default off. Turn on when the decision is hard to reverse.

Enforced in `GOVERNANCE.md` §1.

---

## D

### domain anchor

A persistent knowledge territory — `Distributed Systems`, `AI / ML Systems`, `Product Engineering`, `Market & Buyer Strategist` — queryable but not required to speak during review. Housed in `domain/`. Anchors are *facts*, not opinions.

**Common confusion:** anchors aren't "more important" than council seats. Different roles. Anchors answer knowledge questions council epistemologies cannot.

---

### dwell on evidence

The discipline of stating, when citing evidence, what exactly was observed, by whom, and what would change the conclusion. Single citations of vendor-marketing material do not earn `Confidence: High`. Two independent observations earn `= 2`. Three dimensionally-different sources earn `= 3`.

---

## E

### engineer-voiced problem

A problem defined in the language of the operator who suffers it (DevOps, platform, SRE). Not a product framing. Used in Domain Research. The voice discipline is part of what holds back premature solutioning.

---

### Empiricist

One of five council seats; falsification + citation discipline. Inspects that confidence labels are earned, that the validity gate is intact, and that single-source vendor claims aren't inflating guard scores. Falsification criterion is required for any "show-stopper objection."

---

### epistemology

A *way of thinking* distinct from a *role persona*. Council seats are epistemologies (Empiricist = falsification; Cartographer = frame drift; Red Team = failure mode; Economist = incentive; Operator = workflow fit). Role-personas conflate domain with stance; epistemologies keep the two axes separate.

**Common confusion:** "Red Team" is not "security person." Red Team forensics failure modes. "Operator" is not "DevOps person." Operator inspects workflow fit. Treat each name as a thinking discipline, not a role ownership.

---

### evidence quality (axis A1)

The first axis of the Problem Ranking Framework. 0–3, with gates at ≥ 2 to enter the eligible list. Mixed-vendor evidence alone is insufficient for `= 3`. Two independent observations are typically the floor.

---

## F

### founder override

When the founder rejects the council (or any single agent), with reasoning recorded: which position was overruled, why, and what would trigger reconsideration. Allowed and expected; the framework is the founder's *adversary*, not the founder's *boss*. Silently throwing away a heavily-dissented decision is forbidden.

Enforced in `GOVERNANCE.md` §6.

---

### friction penalty

The second composite metric in the Problem Ranking Framework, summing axes A7 (technical friction), A8 (implementation complexity), A9 (buyer friction), A10 (seller friction). Higher = worse for founder. Subtracted from opportunity fitness to yield net opportunity score.

---

## G

### gate (validity gate)

The three-axis floor in the Problem Ranking Framework: Evidence quality ≥ 2, Frequency ≥ 1, Severity ≥ 1. A problem that fails any of the three is returned to Domain Research for evidence collection, *not* forwarded to PM.

---

## H

### hidden consensus

The trap where multiple agents arrive at the same conclusion through different lenses (correlation, not independent verification). Cross-examination surfaces this by forcing each agent to state their evidence basis explicitly.

---

## J

### journal

Institutional decision memory at `journal/`. One file per meaningful decision; filename convention `<YYYY-MM-DD>-<verb-object-slug>.md`. Every entry preserves dissent verbatim, cites evidence with source, and records founder overrides. The system's audit trail.

---

### journal-template

A specific scaffold at `templates/journal-template.md` that an entry should follow. Fields include: decision, evidence, dissent (preserved), override (if any), change conditions.

---

## N

### net opportunity score

A score in the Problem Ranking Framework: Opportunity Fitness − Friction Penalty. Range −12 to +18. Higher is better. **Never** taken as the sole basis for founder selection — qualitative council commentary accompanies the numeric rank.

---

## O

### opportunity fitness

Composite metric in the Problem Ranking Framework: sum of axes A1 (evidence quality) + A2 (frequency) + A3 (severity) + A4 (tooling under-service) + A5 (market under-saturation) + A6 (founder advantage). Range 0–18.

---

### Operator

One of five council seats; workflow + adoption discipline. Inspects whether a candidate direction drops into a real workflow or requires a process transformation the founder hasn't signed up for. Multi-stakeholder deployments raise the count.

---

## P

### persona

A *role ownership* of a slice of work (e.g., "DevOps", "CTO", "founder"). **Used deliberately as the contrast to *epistemology*.** Council seats are NOT personas. Using personas as council members creates role collision; epistemologies avoid it.

---

### preserve dissent

The rule that disagreement in any review must remain visible in the journal entry, verbatim. Not summarised at the bottom, not flattened into "the council agreed." Distinct objections remain distinct entries in the `journal/`.

Enforced in `GOVERNANCE.md` §1.

---

### Problem Ranking Framework v1.0

The 10-axis scoring instrument defined for use between Domain Research and Product Management. Three gate axes + three fitness axes + four friction axes. Output is a sorted eligibility list; not a recommendation. Application of this framework to candidate problems is recorded in `journal/2026-06-19-r1-problem-ranking-application.md`.

---

### promotion

When a working specialist moves up. Requires 10+ journal entries citing it as load-bearing AND founder opt-in. Recorded as a journal entry; not silent.

Enforced in `GOVERNANCE.md` §8.

---

## Q

### quiet disagreement

The trap of an agent that has objections but does not raise them, often because the artifact is well-presented or the founder has signalled strong preference. Cross-examination surfaces these.

---

## R

### Red Team

One of five council seats; catastrophe + attack-surface discipline. Inspects how each candidate direction gets killed in the first 18 months. Failure modes are concrete mechanisms, not vague objections.

---

### Reference

A skill classification (see *classification*): on disk; consulted by reading the file, not by invoking the skill. No auto-load.

---

### Removed

A skill classification (see *classification*): not consulted. Not even read. Most of `claude-skills/` is in this bucket.

---

### Representative (council member)

Disambiguation: a council member is **not** a "representative" of a constituency. Council members vote individually on epistemology, not on anyone's behalf. There is no constituency system.

---

## S

### session

A continuous Claude Code interaction, beginning at message 0 and ending when the founder signals "stop" or moves to a new artifact topic. Sessions may be reopened by re-loading journal entries.

---

### show-stopper objection

A council claim that a direction is fatally flawed. **Requires** a specific falsification criterion. "This will fail" without a mechanism is *vague objection*, not show-stopper, and may be ignored.

---

### surviving hypotheses ledger

The numbered hypothesis ledger previously lived at `context/02_SURVIVING_HYPOTHESES.md`. It was consolidated into `context/08_AI_HANDOFF.md` (surviving hypotheses) and `context/00_PROJECT_CONTEXT.md` (active research questions) during the Atlas repository hygiene session (2026-07-05). Each hypothesis carries a source session, claim, and current status. Distinct from the journal — the handoff file tracks *claims that drive the project*; the journal tracks *decisions already taken*. Hypotheses are about uncertainty to manage, not decisions already taken.

**Common confusion:** "Surviving" does not mean "proven." It means "still being tested without contradicting evidence."

---

## T

### Temporary

A skill classification (see *classification*): invoked by name per session; consults class *if and when* the founder names it. Retires (drops to `Reference`) after six months of non-use.

---

### tolerate merton (no)

The framework **does not** tolerate Normative Drift, where the system drifts toward "the way problems are usually solved." Each decision must be re-grounded in this founder's needs, not industry default.

---

## V

### valid scoring card

A scoring card whose 10-axis values are each cited with evidence and whose validity gate (≥2 / ≥1 / ≥1 for A1–A3) is met. Allowed to enter the eligibility sort. Invalidation passes the problem back to Domain Research.

---

## W

### working specialist

A spawned, task-specific expertise used to fill gaps that a council epistemic reviewer cannot resolve. Housed in `working/`. Lifecycle: active → `.candidate` (90 days untouched) → `examples/` (6 sessions unused). After 10+ journal citations it may be proposed for promotion to a domain anchor.

---

## Census of load-bearing terms

If a future file uses any of the following terms without a definition, this glossary is its anchor: **council, domain anchor, working specialist, validity gate, opportunity fitness, friction penalty, net opportunity score, epistemology, Engineer-voiced problem, surviving hypothesis, journal entry, Temporary, Reference, consensus smoothing, founder override, cross-examination, preserve dissent, show-stopper objection**.

---

## Provenance

Definitions here are anchored to the conversation transcript and to:

- `agents/README.md`
- `agents/GOVERNANCE.md`
- `agents/REFERENCES.md`
- `agents/journal/README.md`
- `agents/journal/2026-06-19-r1-problem-ranking-application.md`
- `context/00_PROJECT_CONTEXT.md`
- `context/08_AI_HANDOFF.md`
- `docs/project/PROJECT_HISTORY.md`

No new terms invented; every glossary entry corresponds to a term used at least once in the captured surfaces. Where a term has a project-specific weight, this glossary names the trap or common confusion directly.
