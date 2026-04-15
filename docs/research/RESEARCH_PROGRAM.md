# Research Program

Source of truth: `FOUNDER_DECISION_BRIEF.md` (treated as approved; not regenerated, not rewritten, not summarized). This document plans execution of the eight tasks enumerated in that brief's "Research Plan" section. It is bound by:

- `agents/GOVERNANCE.md`
- `context/00_PROJECT_MISSION.md`
- Sequential council review order: **Empiricist → Cartographer → Red Team → Economist → Operator** (binding, founder-confirmed).

**This document is an execution plan, not the research itself. It does NOT** pick a product, generate interview questions, perform interviews, draft solutions, define a PRD, draw architecture, or discuss implementation. It does NOT move into Product Management, UX, Business Analysis, or any downstream phase. The research program ends at journal-recorded wave completion. PM does not begin inside this program.

---

## 0. Critical-thinking pass on the eight approved tasks

Before sequencing, do any of the eight tasks duplicate each other?

- **Task 1 (founder intake)** vs. **Task 2 (practitioner interviews):** same interview methodology, but distinct evidence targets — founder-self vs. practitioner-observed. Not mergeable.
- **Task 6 (P7 premise re-validation, 3 platform engineers)** vs. **Task 2 (P9/P3/P8 pain, 5–8 per problem):** methodology overlap, but framing is distinct. Task 6 tests a specific Cartographer objection (premise-wrong); Task 2 is general pain corroboration. Merging them would lose the governance traceability that required Task 6's existence. Not merged.
- **Task 7 (vendor re-audit, desk research)** vs. **Task 2 (primary interviews):** different evidence types — Task 7 inspects source independence of **existing cited** evidence; Task 2 generates **new** evidence. Different artifacts.
- **Task 7 (evidence re-audit)** vs. **Task 8 (unit-economics check):** Task 7 is upstream of Task 8 (Task 8 implicitly cites public incumbent pricing — once vendor bias is screened, only independent pricing is admissible). They are sequenced, not redundant.
- **Tasks 3 and 4 (P10 spike, P5 spike):** different problems, different technologies (eBPF/observability vs. SPIFFE/SPIRE identity). Independent.

Verdict: no two tasks collapse into one without losing either governance traceability or evidence-type distinction. The eight are kept. Sequencing is the remaining question.

---

## 1. Research Roadmap

Waves are ordered by **expected reduction in project uncertainty**, not by chronological ordering of effort or by the founder priority ordering in the brief. A wave's job is to make the next wave's findings more grounded, not faster.

| Wave | Title | Tasks | What kind of uncertainty it removes |
|---|---|---|---|
| 1 | Foundational Recalibration | Task 1, Task 7 | Cross-cutting scoring-baseline corruption (A6=0 across ten problems; A1 self-declared "High" with vendor bias flagged) |
| 2 | Long-Shot Retirement | Task 5, Task 6 | P6 buyer existence, P7 premise validity — candidates that already carry stacked cautions |
| 3 | Shortlist Pain Validation | Task 2 | P9/P3/P8 — primary-source corroboration or contradiction of existing secondary-source evidence |
| 4 | Technical Premise Feasibility | Task 3, Task 4 | P10 root-cause accuracy, P5 SPIFFE engineering depth — premises that have never been measured |
| 5 | Unit-Economics Reality Check | Task 8 | ACV / sales-cycle quantification for surviving flagged candidates |

### Wave 1 — Foundational Recalibration (Tasks 1 & 7)

**Objective.** Close the two cross-cutting data gaps that affect every problem's underlying rating **before** any per-problem research effort, so that downstream per-problem findings are interpreted against honest baselines.

**Rationale.** The brief itself flags both gaps at the top: (a) every problem carries A6 = 0 (no founder-advantage evidence), and H3 explicitly requires prior work, not stated interest, to set A6; (b) Domain Research's "Confidence: High / Evidence quality: High" labels were self-declared, pre-council, and the Empiricist pass at R1 quietly noted vendor sourcing on P1, P2, P4 without changing scores. Without recalibrating these, every problem-level finding downstream sits on contested foundations — and the brief already shows this matters (e.g., "evidence still missing" lines on almost every problem).

**Expected deliverables.**

- A founder-profile document with explicit A6 mapping per problem (Task 1).
- An annotated bibliography per P1, P2, P4, separating primary / independent evidence from vendor-originated material (Task 7).
- Re-labeled A1 evidence-quality ratings for P1, P2, P4, tied to each source's independence (Task 7).

**Exit criteria.**

- At least one problem has A6 > 0 grounded in real evidence, not stated interest (Task 1 done).
- Each of P1, P2, P4 has an evidence-quality re-rating tied to specific source independence (Task 7 done).
- A single journal entry records both tasks' outputs, preserving any council dissent on the re-ratings.

**Risks.**

- Founder intake may stall, per H5 ("missing intake is not fabricatable"). *Mitigation:* cap intake at the brief's 1–2 focused sessions; record the failure rather than fabricate values downstream.
- Re-rating may lower net scores on problems currently tied at the top of the shortlist (P4, P2 at net=6). *Mitigation:* treating this as expected output, not a setback — the brief already named it.

**Why this wave exists before the next.** Task 5 and Task 6 (Wave 2) will retire or retain two candidates. Their decisions are easier to interpret against honest evidence and founder-fit than against the current zero-everywhere / maybe-vendor baseline. Wave 1 produces that groundedness for free before any candidate retirement begins.

### Wave 2 — Long-Shot Retirement (Tasks 5 & 6)

**Objective.** Resolve the two lowest-net-score candidates' core uncertainty so they can be formally retired, reframed, or retained — and so subsequent waves are not burdened with candidates that should already be out.

**Tasks and rationale.** P6 (net 1) and P7 (net 0, actively downgraded) carry stacked cautions. The Operating Intelligence doctrine is "killing bad ideas quickly." Council already noted three divergent possible buyers for P6 and concluded the P7 premise is wrong (Cartographer drop). The brief's own research-readiness matrix lists both at "High" additional-research difficulty, but both tasks are bounded (3–5 days each, interview-bound), so they are cheap to resolve even at high difficulty.

**Expected deliverables.**

- A buyer-type memo for P6 (2–3 organizations), with a named budget-holder and described decision process **OR** a documented conclusion that no clear buyer exists.
- A premise-validation memo for P7 from 3 platform engineers, directly testing the Cartographer's "value-promise gap, not missing tooling" objection, with a verdict: premise holds / premise wrong / reframed premise.
- A per-problem re-eligibility decision (retained / retired / reframed) recorded in journal with founder override, if any.

**Exit criteria.**

- P6 has a definitive verdict or a formally documented "no clear buyer" conclusion — not "unresolved."
- P7 has a defensible reframe OR a formal retirement, with the Cartographer's objection explicitly addressed.
- Journal entries preserve any interview dissent verbatim. Sequential review applies.

**Risks.**

- Investigation may uncover a defensible framing (P6 buyer emerges; P7 reframes into a real problem). *Mitigation:* this is a valid Wave 2 outcome — promote the problem into Wave 3 territory, update the candidate set, record it as a "Cartographer objection rebutted" finding.
- Premise validation may produce weak or contradictory evidence. *Mitigation:* "evidence is insufficient" is itself a finding, distinct from "premise holds."

**Why this wave exists before the next.** P9, P3, P8 (Wave 3 candidates) become a more honest shortlist the moment P6 and P7 are resolved. Until then, the candidate set is artificially inflated; running Wave 3 against the inflated set would mix real candidates with candidates that should already be out.

### Wave 3 — Shortlist Pain Validation (Task 2)

**Objective.** Produce independent, primary-source pain validation — not secondary-source — for the three problems the research-readiness matrix flagged as cheapest to de-risk further.

**Tasks.** Task 2 (non-vendor pain validation interviews — P9, P3, P8).

**Rationale.** Per the brief's research-readiness matrix, P9 (on-call/toil), P3 (secrets/config drift), P8 (incident coordination) are the cheapest to de-risk further. They have the strongest non-vendor evidence already (P9 — SRE book/survey/blog; P3 — breach postmortems/standards docs; P8 — incident.io's own named-incumbent material). All prior evidence is secondary (books, surveys, vendor docs); this wave generates primary corroboration or contradiction.

**Expected deliverables.**

- 5–8 structured interview notes per problem, with practitioners who currently experience the stated pain.
- Per-problem verdict: corroborated / contradicted / mixed.
- A follow-on memo: which corroborated candidates proceed to further work (technical spike / founder-go decision), which contradicted candidates are deprioritized or retired, which mixed candidates require a second-pass interview round.

**Exit criteria.**

- All three problems have a verdict with at least the threshold of 5 interviews per problem, transcripts preserved.
- Any contradicted finding triggers a journal-recorded re-evaluation of the problem's standing against the R1 record.
- The "tight ICP" ambiguity on P9 (Red Team framing — focus vs. too-small) is resolved **or** its continued ambiguity is explicitly re-flagged for Wave 5 (Task 8).

**Risks.**

- Interview scheduling slippage, named in the brief as the dominant scheduling risk. *Mitigation:* scope-protection by capping to the bounded count.
- Selection bias. *Mitigation:* active inclusion of teams using the named incumbents (incident.io, PagerDuty, HashiCorp Vault/Terraform) to avoid interviewing only satisfied or only dissatisfied users.
- Cross-wedge disagreement inside one candidate (e.g., P3 "drift is real" alongside "feature, not category"). *Mitigation:* preserve dissent verbatim per governance; do not smooth.

**Why this wave exists before the next.** Wave 4's technical spikes test product premises whose value depends on pain being real. For P9/P3/P8, pain is interview-readable; running Wave 3 first prevents a spike from being run against a problem whose pain turns out to be missing. (For P5 and P10, no interview reads their core premise, so Wave 4 runs unconditionally — those two are "test the technical claim" candidates, not "validate the pain" candidates.)

### Wave 4 — Technical Premise Feasibility (Tasks 3 & 4)

**Objective.** Produce honest measurements of the two technically-deepest candidate problems' core premises, which currently lack any engineered validation.

**Tasks.** Task 3 (root-cause accuracy spike, P10) — Task 4 (SPIFFE issuance spike, P5).

**Rationale.** These are the only problems where the central claim has never been measured. P10's Red Team "research-grade" note is a direct technical challenge that has gone untested. P5 received the only council upgrade in R1 (Cartographer lift citing SPIFFE standardization), but no engineering validation has occurred to confirm whether real differentiation remains after drinking from the existing SPIRE tooling. Without these measurements, both problems' standing in the candidate set rests on narrative, not evidence.

**Expected deliverables.**

- For P10: a precision/recall (or equivalent) figure for root-cause identification against a known dataset, with honestly reported failure cases; no overselling as production-ready.
- For P5: a working minimal SPIFFE-compatible identity issuance flow plus a written account of what remains genuinely unsolved after using existing SPIFFE/SPIRE components.
- Per-problem verdict: viable near-term build target / multi-year research bet / reclassify as OSS contribution rather than product build.

**Exit criteria.**

- For P10, a measured number exists where none did before.
- For P5, an explicit answer exists to "is there real engineering work left here, or does the existing ecosystem already cover it."
- Each verdict is journal-recorded with confidence label and stated change-condition.

**Risks.**

- Both spikes may produce negative verdicts. *Mitigation:* a negative verdict is a valid Wave 4 outcome; it does not invalidate the wave. Logging "research-grade, not productionizable now" or "ecosystem already covers it" closes the most expensive uncertainty at minimum cost.
- Spikes may reveal unanticipated design depth that the problem is larger than solo-feasible. *Mitigation:* log it; route back to project governance via founder-evidence path, not silently.

**Why this wave exists before the next.** Wave 5's unit-economics check is meaningful only for candidates whose pain is corroborated (Wave 3) **and** whose technical premise holds or has been honestly measured (Wave 4). Unit economics on a problem that does not work is a wasted measurement; ACV rating runs last.

### Wave 5 — Unit-Economics Reality Check (Task 8)

**Objective.** Convert the qualitative "small ACV," "long sales cycle," and "tight pricing" flags the Economist raised on five problems without quantifying them, into order-of-magnitude dollar estimates across the surviving candidates.

**Task.** Task 8 (ACV / unit-economics reality check — flagged pool: P9, P4, P2, P8, P10).

**Rationale.** Structural revenue-model risk is named on five of the ten problems but has never been measured. It currently functions as a qualitative caution, not a decision input. Running this LAST means: (a) Wave 5 can incorporate pain-validation findings (Wave 3) and premise-feasibility findings (Wave 4) into its interpretation, and (b) problems retired or reframed in Waves 2–4 drop out of the ACV pool automatically. The task is the cheapest to abandon among the five if earlier waves narrow the candidate set sufficiently.

**Expected deliverables.**

- A comparative ACV table for each **surviving** candidate from the flagged pool, sourced from public incumbent pricing and interview-derived data.
- A per-candidate disqualification assessment: "disqualified on unit economics alone," "acceptable but tight," or "not disqualifying."
- A journal entry updating each surviving candidate's A8 / A9 / A10 dimensions with quantitative grounding.

**Exit criteria.**

- Each surviving candidate from the flagged pool has an order-of-magnitude ACV estimate attached.
- Each estimate cites at least one independent public-pricing data point.
- A foundation memo records the comparative reading across the surviving pool.

**Risks.**

- The flagged pool will likely change after Waves 2–4 — some problems may retire or reframe, changing which candidates need ACV. *Mitigation:* the wording "surviving candidate from the flagged pool" already handles this; no estimation effort is spent on retired problems.
- ACV estimates are unreliable for narrow niches (e.g., framework-specific compliance). *Mitigation:* prioritize comparable-product pricing; flag nicher estimates with lower confidence and broader change-conditions.

**Why this wave exists last.** It is a synthesizing wave. Stacking it late ensures the unit-economics verdict is grounded in the strongest available evidence, not improvised before per-problem research completed. It is also the cheapest wave to abandon if earlier waves narrow the candidate set to a single problem that PM can already evaluate.

---

## 2. Dependency Graph

```
                            ┌─────────────────────┐
                            │  Wave 1             │
                            │  Foundational        │
                            │  Recalibration       │
                            └──────┬───────────────┘
                                   │
        ┌──────────────────────────┼──────────────────────────┐
        │                          ▼                          │
        │           ┌────────────────────────┐                │
        │           │   Wave 2               │                │
        │           │   Long-Shot Retirement │                │
        │           │   (P6 buyer, P7        │                │
        │           │    premise)           │                │
        │           └────────────┬───────────┘                │
        │                        │                            │
        │                        ▼                            │
        │           ┌────────────────────────┐                │
        │           │   Wave 3               │                │
        │           │   Shortlist Pain       │                │
        │           │   Validation           │                │
        │           │   (P9, P3, P8)         │                │
        │           └────────────┬───────────┘                │
        │                        │                            │
        │                        ▼                            │
        │           ┌────────────────────────┐                │
        │  parallel │   Wave 4               │                │
        │  (with 2) │   Technical Premise    │                │
        │           │   Feasibility          │                │
        │           │   (P10, P5)            │                │
        │           └────────────┬───────────┘                │
        │                        │                            │
        │                        ▼                            │
        │           ┌────────────────────────┐                │
        └─────────► │   Wave 5               │ ◄──────────────┘
                    │   Unit-Economics       │  consumes Wave 1's findings,
                    │   Reality Check        │  Wave 3's interview data,
                    │   (survivors from      │  Wave 4's spike verdicts
                    │    flagged pool)        │
                    └─────────────────────────┘
```

### Critical path

- **Task 1 (founder intake)** — blocks A6 calibration across all ten problems; everything downstream benefits.
- **Task 7 (vendor re-audit)** — blocks evidence-quality reliability on P1, P2, P4 net scores; everything downstream benefits.

### Parallelizable

- Tasks 1 and 7 are parallel: one is reflective intake (sessions), the other is desk research (days).
- Tasks 3 and 4 are parallel: independent technical experiments on different problems.
- Task 2 (Wave 3) and Tasks 3/4 (Wave 4) can overlap on calendar: Task 2 is interview-bound, Tasks 3/4 are engineer-bound.
- The brief explicitly named Task 8 as parallelizable with Task 2. Honored.

### Blocking

- **Wave 2 → Wave 3 (only if P6/P7 get promoted).** If the Wave 2 investigations find a buyer for P6 or a reframed premise for P7, those problems enter Wave 3's interview validation as new shortlist candidates. The brief discourages this outcome but does not exclude it.
- **Wave 3 → Wave 5 (Task 8 consumes interview data).** ACV estimation requires interview-derived pricing data; running Task 8 without Wave 3 produces estimates from public pricing alone, which the brief acknowledges is acceptable but weaker.
- **Wave 4 → Wave 5 (logical).** A problem whose premise is unverified should not have ACV spend on it — that solves a problem that may not exist.

### Anti-pattern avoided

- **Running Wave 5 first.** Premature ACV estimation creates the appearance of confidence before pain or technical premise is validated. Banned.
- **Chronological-by-convenience ordering.** A purely calendar-based schedule (interviews all at once, then spikes all at once, then desk research at the end) would interleave Tasks 1–8 in a way that wastes Wave 1's recalibration on tasks already running on contested baselines. Banned.

---

## 3. Evidence Map

Each task produces a specific artifact with a specific decision it enables. The artifacts are the inputs to PM-handoff; without them, PM would inherit the same uncertainty the R1 cycle warned against.

| Task | Evidence produced | Confidence-typed? | Decision enabled |
|---|---|---|---|
| Task 1 | Founder-profile document; A6 rating per problem | Yes (High / Medium / Low with cited evidence per A6 > 0 claim; H3 standard) | Whether founder fit is a legitimate tiebreaker among research-ready candidates |
| Task 7 | Annotated bibliography per P1, P2, P4; revised A1 ratings | Yes (each source labeled primary / independent / vendor; each re-rating cited) | Whether P1 / P2 / P4 net scores are trusted as-is |
| Task 6 | Premise-validation memo for P7 (3 engineer interviews) | Yes (verbatim interview notes; per-claim confidence label) | Whether P7 is retained, reframed, or retired — and whether the Cartographer's drop holds |
| Task 5 | Buyer-type memo for P6 (2–3 organizations) | Yes (each organization named; budget process described) | Whether P6 is retained or retired |
| Task 2 | 5–8 interview notes per P9, P3, P8 | Yes (verbatim per governance; corroborated / contradicted / mixed verdict per problem) | Whether each problem proceeds to spike, deprioritization, or second-pass interview |
| Task 3 | Accuracy measurement on root-cause identification, against a known dataset | Yes (precision / recall or equivalent; failure cases reported honestly) | Whether P10 is a near-term build target or a multi-year research bet |
| Task 4 | Working minimal SPIFFE issuance flow + unsolved-remaining memo | Yes (each unsolved element explicitly named) | Whether P5 is a differentiated build target or an OSS contribution reclassification |
| Task 8 | Comparative ACV table for surviving flagged candidates | Yes (each row tied to cited incumbent public pricing; per-row confidence) | Whether any surviving candidate is disqualified on unit economics alone |

### Coverage of the brief's research-readiness matrix

| Unquietness axis named in the brief | Closed by |
|---|---|
| Revenue-model risk (P9, P4, P2, P8, P10) | Wave 5 (Task 8) |
| Category viability (P3 — "feature, not category"? ) | Wave 3 (Task 2 interviews) |
| Wedge specificity vs. incumbent (P8) | Wave 3 (Task 2) + Wave 5 (Task 8) |
| Go-to-market realism for solo builder (P4) | Wave 5 (Task 8) + residual signaling from Wave 3 |
| Defensible wedge vs. commodity feature (P2) | Wave 3 (interview differentiation) + Wave 5 (pricing pressure) |
| Structural market dependency (P1 — "storage/cost revolution required") | Wave 1 (Task 7 evidence re-rating) — partial; full resolution may require Wave 6 not in scope |
| Buyer existence and solo-shippability of security-critical software (P5) | Wave 4 (Task 4 spike) — engineering leg; Wave 5 (Task 8) — economic leg |
| Buyer identification (P6) | Wave 2 (Task 5) |
| Whether the central approach works at all (P10) | Wave 4 (Task 3) |
| Whether the problem premise is valid (P7) | Wave 2 (Task 6) |

### Coverage of "evidence still missing" lines from the brief

Every "evidence still missing before an experienced engineering leader approves building" line in the brief maps to a specific task. The build-out above is sufficient; nothing in the brief's evidence-gaps list is left uncovered by this program.

**One acknowledged partial coverage.** P1's "evidence still missing — real cost-per-GB benchmark against at least one incumbent" is partially closed by Wave 1 (Task 7 — vendor re-audit) but a true cost-benchmark requires independent measurement, which is a separate engineering task. *Decision:* this is flagged in §6 as an admissible residual uncertainty carried into PM-handoff, not expanded into a new task inside this program. Adding such a benchmark would require building infrastructure, which is out of research phase.

---

## 4. Decision Gates

Each wave has a single named gate. Gates are journal-recorded events, not a re-run of the brief. Sequential council review applies to all gate decisions.

### Wave 1 gate — "Foundations are honest"

- **Trigger.** Both Task 1 founder intake and Task 7 vendor re-audit have produced their documented outputs.
- **Pass condition.** At least one A6 > 0 with cited evidence; P1, P2, P4 each have a re-rated A1 grounded in source independence.
- **Fail handling.** If A6 still equals 0 for all ten after Task 1 (the brief explicitly flags this as plausible), record that finding explicitly and let downstream waves continue to run, but with founder-fit held at zero and explicit confidence labels propagating.
- **Journal.** A strict, council-aware entry combining both tasks' outputs. Empiricist dissent preserved if A1 re-ratings conflict with Domain Research's original self-declared ratings.

### Wave 2 gate — "Long-shots are resolved"

- **Trigger.** Task 5 (P6) and Task 6 (P7) verdicts recorded.
- **Pass condition.** Each of P6 and P7 has a binary or reframe verdict, not "unresolved."
- **Fail handling.** If a verdict is still "unresolved" after this wave, retry once with one additional bounded interview cycle; if unresolved persists, retire the problem on doctrine (H5 — "missing intake is not fabricatable"). Continued unresolved status is itself a finding, not a deferral.
- **Journal.** Two entries (one per problem) preserving interview dissent verbatim. Sequential review applies.

### Wave 3 gate — "Shortlist pain is real or contradicted"

- **Trigger.** Task 2 interview verdicts for P9, P3, P8 recorded.
- **Pass condition.** Each of P9, P3, P8 has at least 5 interviews and a corroborated / contradicted / mixed verdict.
- **Fail handling.** Contradictions highlight scoring mismatches against R1; preserve dissent; allow one additional interview round only.
- **Journal.** One entry per problem with confidence labels and dissent preserved.

### Wave 4 gate — "Premises are measured"

- **Trigger.** Task 3 (P10) and Task 4 (P5) spike outputs recorded.
- **Pass condition.** Each spike has a measured number (P10) or an explicit "what is left to solve" mapping (P5), with caveats.
- **Fail handling.** Negative verdicts are valid outputs — record them. Spike data is research data; do not "spin" toward promotion or retirement based on the spike alone.
- **Journal.** One entry per spike with engineering verdict, confidence label, and stated change-condition.

### Wave 5 gate — "Unit economics are quantified"

- **Trigger.** Task 8 ACV table produced for surviving candidates from the flagged pool.
- **Pass condition.** Each row has at least one independent public-pricing citation.
- **Fail handling.** If a problem cannot get a comparable pricing reference, record it as a low-confidence estimate with the gap explicitly named.
- **Journal.** One entry comparing across the surviving pool with founder override filed, if any.

### Crossing all wave gates — "Decision block is consolidated"

- **Trigger.** All five wave gates cleared.
- **Pass condition.** The PM cycle may begin only when:
  1. A6 ratings re-anchored against real input (Wave 1 / Task 1) or explicitly zeroed with rationale recorded.
  2. A1 evidence-quality ratings for P1, P2, P4 re-anchored by source-independence test (Wave 1 / Task 7).
  3. P6 and P7 no longer "unresolved" (Wave 2).
  4. P9, P3, P8 pain corroborated, contradicted, or irrecoverably mixed (Wave 3).
  5. P5 and P10 technical premises measured (Wave 4).
  6. ACV comparable across surviving candidates (Wave 5).
  7. A journal entry per wave captures founder override (if any), preserved dissent, and stated change-conditions.

**Critical:** program completion does **not** mean "all candidates are equally viable." It means the candidate set is honest. PM is permitted to begin only against the honest candidate set, not against a fabricated shortlist.

---

## 5. Exit Criteria for Every Wave

A consolidated list, restating §1's per-wave exit criteria with explicit journal obligations:

### Wave 1 exit

- Founder profile document exists with A6 mapping per problem (Task 1).
- Annotated bibliographies exist for P1, P2, P4 (Task 7).
- A1 evidence-quality ratings re-labeled per problem (Task 7).
- Single journal entry combining both tasks' outputs, with preserved dissent and source-by-source independence labeling.

### Wave 2 exit

- P6 has a buyer-type verdict **or** a documented conclusion that no clear buyer exists.
- P7 has a premise-holds / premise-wrong / reframed-premise verdict.
- Per-problem re-eligibility decision recorded.
- Verbatim interview notes preserved per governance.

### Wave 3 exit

- 5–8 structured interview notes per P9, P3, P8.
- Per-problem verdict: corroborated / contradicted / mixed.
- Follow-on memo: which candidates proceed, which retire, which need a second-pass interview round.
- Verbatim interview notes preserved; per-claim-cluster confidence labels.

### Wave 4 exit

- For P10: a measured root-cause-identification accuracy figure on a known dataset, with explicit failure cases.
- For P5: a working SPIFFE-compatible issuance flow + written remaining-unsolved account.
- Per-problem viability verdict: near-term build / multi-year research / OSS-contribution reclassification.
- Confidence labels and change-conditions recorded on each verdict.

### Wave 5 exit

- Comparative ACV table for **surviving** candidates from the flagged pool.
- Per-candidate disqualification assessment (disqualified / tight / not disqualifying).
- All citations to incumbent public pricing.
- Journal entry comparing across the surviving pool with founder override filed if any.

---

## 6. Criteria Required Before Product Management May Begin

These are program-level requirements, distinct from per-wave exit criteria. The PM stage may not begin without meeting all of them.

1. **Founder-fit is grounded, not zero everywhere.** Either (a) at least one A6 > 0 grounded in real evidence, **or** (b) an explicit, journal-recorded statement that "no founder fit can be drawn from available input, founder-fit cannot be a tiebreaker." If (b), this is a known opaque factor carried visibly into the founder decision step.
2. **Evidence quality is re-rated.** A1 ratings for P1, P2, P4 are explicitly grounded in source-independence tests, not self-declared as High.
3. **P6 and P7 are no longer open questions.** Each has a binary verdict or a reframe-and-keep verdict with dissent preserved.
4. **P9, P3, P8 have interview verdicts.** Corroborated / contradicted / mixed, each with at least the threshold interview count.
5. **P5 and P10 are measured, not narrated.** Spike output exists with confidence labels and stated change-conditions.
6. **ACV table is comparable across surviving flagged candidates.** At least one independent public-pricing citation per row.
7. **Candidate set is reduced to research-ready problems only.** PM operates against a candidate set where uncertain problems are out, not in.
8. **Each research outcome has journal traceability.** Every wave's decision is in `agents/journal/<YYYY-MM-DD>-<slug>.md`, with founder override (if any) recorded.
9. **No new problem has been added.** The R1 problem set plus the brief defined the candidates. The research program does not generate new candidates mid-stream.
10. **No architecture, no product choice, no PRD, no instance of moving into PM.** The research program ends at journal-recorded wave completion. PM does not begin inside this program.

### Admissible residual uncertainties

Two gaps may legitimately survive the program and remain as documented unknowns at PM-handoff, rather than being expanded into new research tasks:

- **P1's independent cost-per-GB benchmark.** Partial coverage only (Wave 1 / Task 7 vendor re-audit). True measurement requires building or renting storage infrastructure, which is out of research phase. Carry forward as residual.
- **Cross-spike integration tests** (e.g., does the SPIFFE spike interoperate with a real workload identity chain?). Single-task spikes do not produce system-level assurance; integration is a later-phase concern. Carry forward as residual.

These are recorded as known unknowns at PM-handoff, not silently dropped.

### Stop rules (binding)

- The research program ends when all wave-exit criteria are met. It does not continue into PM even if more questions arise.
- If, after any wave, additional research is identified as needed, the supplement goes through a **new Founder Decision Brief** — not a Research Program extension. This preserves the discipline that the brief is the source of truth, and changes to it are founder-authorized.
- If founder override is invoked to skip a candidate's research task, the journal records the override with reasoning and a stated change-condition for later reconsideration. Silently skipping a heavily-dissented task is forbidden.
- The program may be abandoned, paused, or have its scope reduced only via journal-recorded founder override. No silent termination.

---

## Provenance

- Source of truth for the eight tasks: `FOUNDER_DECISION_BRIEF.md` (treated as approved).
- Project governance: `agents/GOVERNANCE.md`. Council review order: `PROJECT_HISTORY.md` (founder-confirmed in R1 closure).
- Mission: `context/00_PROJECT_MISSION.md`.
- Pipeline position: `context/01_CURRENT_STATE.md`.
- Research backlog predecessor: `context/03_RESEARCH_BACKLOG.md` (R1 closed; R2/R3 cycles gated on program completion).
- Surviving hypotheses referenced: H3 (founder advantage requires prior work, not stated interest — applied in Wave 1 / Task 1), H5 (missing intake is not fabricatable — applied as a stop rule), H7 (sequential review catches mis-scoring — governing journal review for every wave's journal entries).
- R1 precedent: `agents/journal/2026-06-19-r1-problem-ranking-application.md`. / GLOSSARY.md provides definitional anchoring.

## Confidence

- **Overall program confidence: Medium.** Forward-looking sequencing rests on (a) brief acceptance as stated, (b) each wave delivering its stated artifact, (c) interview scheduling and engineer-time completing within the brief's bounds.
- **Wave-by-wave confidence:**
  - Wave 1 — Medium. Task 7 (desk research) is high-feasibility; Task 1 (founder intake) is the lower-confidence element and may stall (H5).
  - Wave 2 — Medium-High. Both are bounded interview tasks.
  - Wave 3 — Medium. Interview scheduling slippage is named as the dominant risk in the brief.
  - Wave 4 — Medium-Low. Technical spikes may produce negative verdicts; the verdicts themselves are the deliverable, so this risk does not negate the wave's value.
  - Wave 5 — Medium. Depends on surviving candidate set being stable enough to estimate.
- **Change conditions that would lower confidence:**
  - Brief itself is challenged or revised → all waves resequence.
  - A wave produces unanticipated negative verdicts → cascading waves may shrink.
  - Founder-intake stall → the program runs with A6 = 0 explicitly named; downstream waves operate against residual uncertainty.
- **Change conditions that would raise confidence:**
  - Wave 1 produces concrete A6 > 0 with cited prior work.
  - Wave 3 produces corroboration on all three shortlist candidates.
  - Wave 4 produces negative verdicts on both P5 and P10 (these are the lowest-uncertainty end-states because they explain the most). Spikes showing positive differentiation would *lower* confidence because they expand the candidate set rather than reducing it.

---

**End of Research Program.** No artifacts follow this document. PM does not begin here. The next journal entry records program completion; the next stage begins after founder instruction.

<!-- checkpoint: feat(stores): implement verification controller (#63) -->
