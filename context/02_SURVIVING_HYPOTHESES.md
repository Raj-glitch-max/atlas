# Surviving Hypotheses

Hypotheses still being tested as of 2026-06-19. Each is attributable to a session decision and may be retired or rejected after the relevant cycle produces evidence.

## H1 — External tools are *classified*, not auto-imported

Source: External Repository Audit session.
Claim: External skill repositories (e.g., `claude-skills/`) should be evaluated against a classification (Core / Temporary / Reference / Reject) before any adoption. Auto-import dilutes the operating system.
Status: **Surviving** — applied to one external repository so far; more acquisitions expected.

## H2 — A problem-ranking gate must sit between Domain Research and PM

Source: Maintenance session prior to Problem Ranking application.
Claim: Without a uniform scoring pass across candidate problems, the Product Manager would inherit every candidate as equally worthy. The 10-dimensional Problem Ranking Framework introduces a validity gate (evidence, frequency, severity) and a net opportunity score (opportunity fitness minus friction penalty).
Status: **Surviving** — Framework built; one application cycle required to validate or refute.

## H3 — Founder advantage requires prior work, not stated interest

Source: Problem Ranking Framework defining rules.
Claim: "I find this interesting" is not evidence of founder advantage. Founder-advantage scoring demands prior behavior: past shipping, IP, network, or distribution in the target community.
Status: **Carried into R1 application.** All 10 scored problems received A6=0 because the Founder Profile session returned no founder history. The hypothesis **held**: scoring upheld the rule; founder-advantage was not inflated by any problem's appeal alone.

## H7 — Sequential review catches mis-scored dimensions before they calcify

Source: Problem Ranking Application R1 cycle (2026-06-19).
Claim: Pre-council scoring produces at least one over-scoring error per cycle; sequential council review surfaces those; pre-council numbers must not be treated as final until reviewed.
Status: **Carried out in R1.** Cartographer revised P7 down (tooling under-service mis-scored as gap when gap was abstraction-promise, not tooling) and P5 up (Market under-saturation understated given SPIFFE-style standardization). Empiricist flagged vendor-bias weighting without changing scores. No eliminations but two numeric revisions.

## H8 — A flat net-score rank over-states the gap between problems at the top and bottom

Source: Problem Ranking Application R1 cycle.
Claim: The 10-dimension net-score compression undersells qualitative differences. Three problems tied at net=6 but had different opportunity-fitness profiles; one problem at net=0 (P7) is qualitatively different from one at net=1 (P6) — vendor-dominance vs cost-of-build. Quantitative rank + qualitative council commentary together form the actual shortlist.
Status: **Surviving** — net scores are a heuristic; preserved alongside council commentary; do not collapse to rank alone for founder selection.

## H4 — Evidence without citation cannot earn evidence-quality ≥ 2

Source: Problem Ranking Framework; Empiricist discipline.
Claim: A score of 2 on evidence quality requires at least two independent observations; a score of 3 requires three dimensionally different sources. Anecdotes, single-source claims, or vendor marketing do not qualify.
Status: **Surviving** — applies to scoring cards as they are produced.

## H5 — Missing intake is not fabricatable

Source: Founder profile (Session #1).
Claim: When intake data is sparse, a founder profile is not to be invented. Outputs are made with explicit confidence labels and explicit gaps named. The next-cycle hypothesis prefers "gather missing intake" over "commit to a direction on thin evidence."
Status: **Surviving** — applied once; integrated as a system rule.

## H6 — System naming is content-neutral

Source: Maintenance session 2026-06-19.
Claim: Renaming "FounderOS" → "agents" is permissible without architectural redesign because the system name was always a brand, not a method. Method is in `council/`, `domain/`, `working/`, `journal/`, `templates/`, and `GOVERNANCE.md`.
Status: **Carried out** — directory now `/agents/`; doc references updated.
