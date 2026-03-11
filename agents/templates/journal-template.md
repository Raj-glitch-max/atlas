# Journal Entry Template

Template for recording decisions, dissent, and evidence in `/journal/`.

File naming: `<YYYY-MM-DD>-<short-slug>.md`. The slug is the decision's verb-object — e.g., `2026-06-19-pick-pricing-model.md`.

Length target: variable. Most entries fit on 1–2 pages.

---

```markdown
---
date: <YYYY-MM-DD>
slug: <kebab-case>
artifact: <input artifact — claim, decision, proposal, document, code reference>
decision: <one sentence: what was decided>
agents_consulted: <list of agent names>
overrides: <true | false>
related_entries: <list of related slugs, optional>
---

# Context
[What the founder brought forward. 1–3 paragraphs. The actual claim, not the surrounding anxiety.]

# Decision
[The committed call. One paragraph. What is now true after this entry, what is no longer true.]

# Evidence cited
[Each substantive piece of evidence used. Named. With source or reasoning path.]

# Council positions

## The Empiricist
[Position in their voice. Calibrated confidence. What would shift it.]

## The Red Team
[Position. Mechanism-not-vibe. Failure scenarios specified or refused.]

## The Operator
[Position. Adoption or workflow concerns with named mechanisms.]

## The Economist
[Position. Cost attribution. Actor named.]

## The Cartographer
[Position. Restated claim. Frame surfaced.]

# Domain anchors consulted
[Each domain anchor that spoke. Their contribution, not their reasoning.]

# Working specialists consulted
[Each working specialist that spoke. Brief.]

# Dissent preserved
[Where agents disagreed. Both positions. Don't smooth.]

Quote disagreements verbatim, not paraphrased. The point of the journal is that future-you can read the disagreement as it stood, not as it eventually resolved.

# Founder override (if applicable)
[What was overruled, why, what would have to be observed for reconsideration.]

# Open questions
[What this entry did NOT answer. What it explicitly defers.]

# Status
- decided: <YYYY-MM-DD>
- revisited: <list of dates this entry was reopened, optional>
- superseded_by: <slug if this entry was replaced, optional>
```

---

# Notes on filling out

- **Frontmatter is binding.** When grepping the journal later, the metadata fields drive every search.
- **Sections can be empty.** If a domain anchor wasn't consulted, write "not consulted." Don't invent a position for them. The journal is honest recording, not performed thoroughness.
- **Dissent is mandatory when present.** If two agents disagreed, that disagreement is in this file. No "we eventually aligned" footnotes that erase the conflict.
- **Founder override is its own section.** It is not buried in the decision. The framework's whole point is keeping override visible.
- **Status field tracks the entry's life.** `decided` always populated; `revisited` populated on each re-reading that produced a new sub-decision; `superseded_by` once only.
- **`related_entries` is your best friend when chaining decisions.** A pricing decision that ties to an earlier ICP decision should reference both ways.
