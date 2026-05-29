# Journal

The institutional memory of agents v1.0.

This folder holds every meaningful decision, with its evidence, its dissent, and its overrides preserved.

## What goes in here

- **Decisions**: choices that commit the founder to a path, cost time / money / political capital, or close off alternatives.
- **Dissent preserved verbatim**: where agents disagreed, the disagreement is here. Future-you will need it.
- **Evidence cited**: what justified each decision. Source or reasoning path. Each piece named.
- **Founder overrides**: when the founder rejects the council, the override is here with reasoning and the condition that would trigger reconsideration.
- **Open questions**: anything this entry deliberately did NOT answer.

## What does NOT go in here

- Routine tactical decisions (aconfig tweak, a single-line rename) that don't commit the founder's path. Direct mode decisions can be journalized briefly or skipped.
- Agent chatter that produced no decision. Use journal entries to record resolved items, not deliberations.
- Personal notes unrelated to a decision. The journal is decision-memory.
- AI-generated slop that wasn't reviewed. Audit each entry before file.

## File naming

`<YYYY-MM-DD>-<verb-object-slug>.md`

Examples:
- `2026-06-19-pick-pricing-model.md`
- `2026-07-02-build-v0-without-auth.md`
- `2026-07-14-decline-enterprise-pilot.md`

Slug should let you grep for the decision in 60 seconds. Verb-object, not abstract nouns.

## Lifecycle

| State | Trigger | Action |
|---|---|---|
| Created | Decision moment | File with `decided` populated. |
| Revisited | Subsequent journal entry depends on this one or disagrees with it | Update `revisited:` with date and brief reason. |
| Superseded | This decision later replaced by another | Set `superseded_by:` link. Don't delete. |
| Archived | Routine cleanup at year-end | Move to `<year>/archive/`. Visible but out of the way. |

`superseded_by` is one-way. The original entry stays intact; the new entry references it. Most decisions restate prior decisions rather than supersede them — re-stating isn't overriding for this purpose.

## How to grep

```bash
# All decisions on a topic
ls journal/ | grep -i pricing

# All entries where cartographer dissented
grep -l "Cartographer: dissented" journal/*.md

# All founder overrides in the last 90 days
grep -l "founder override\|override applied" journal/*-$(date +%Y-%m-%d).md journal/*.md | head

# Recurring dissent agent
grep -h "^[A-Za-z].*dissented" journal/*.md | sort | uniq -c | sort -nr
```

These are starting scripts. Build up to your own as your journal grows.

## What this folder is for — long-term

After a year of decisions, this folder answers questions you currently can't:

- What did I commit to before, and is the commitment still alive?
- Who dissented on the heads-of-terms for the deal I now think was wrong?
- What was the evidence base for the architectural decision I keep defending?
- Which Open Questions from a year ago are still open, and which can be closed safely?
- What did I override three times in a row, and was the override pattern productive?

The journal is the framework's way of remembering things you'd otherwise reconstruct from memory.

## Compaction

When entries exceed 50:

- Keep dissent and override information intact always.
- Compress routine decision entries into yearly summary files (`journal/<year>-summary.md`).
- Don't compress entries from the last 90 days.
- Don't compress entries that future searches are likely to land on (specific ICP choices, major pricing decisions, architecturally-load-bearing decisions).

Compression is opt-in by the founder; the journal's bias is to retain.

## Reading the journal

Before re-opening a closed decision, read the original entry. The point is reviewing decisions against original reasoning, not against current reasons. Drift between the two is a Cartographer flag.

## Quarterly self-review

Every 90 days, Cartographer runs `/journal/review-quarterly.md` with:

- Which council agents caught unexpected things this quarter (good signal).
- Which council agents missed things other agents had to catch (improvement signal).
- Which founder overrides have aged well, which have aged poorly (calibration signal).
- Which Open Questions are still open a year later (frame health signal).

Findings write back into agent files and templates, but only as explicit edits with `last_clarified` updates.

## What this folder cannot do

- Cannot defend against rewriting. The founder owns the integrity of this folder.
- Cannot enforce consistency across entries. Each entry is a snapshot at decision time.
- Cannot solve memory fragmentation across multiple working trees / multiple laptops. For that, use a git repo.

<!-- checkpoint: planning(conformance-targets): clarify conformance targets -->

<!-- checkpoint: repo(architecture-draft): finalize architecture draft -->

<!-- checkpoint: chore(record): clean attenuation rule engine -->

<!-- checkpoint: chore(record): harden boundary check -->

<!-- checkpoint: chore(security): audit Fuzz Verification core target -->

<!-- checkpoint: refactor(scripts): refactor lab environment topology -->
