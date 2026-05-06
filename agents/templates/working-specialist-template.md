# Working Specialist Template

Template for narrow technical specialists spawned on demand. Tighter than permanent agents.

Length target: 400–900 words.

---

```markdown
---
agent: working
name: <kebab-case>
honorific: <Display Name>
scope: <one-line: what this specialist knows and what it doesn't>
created: <YYYY-MM-DD>
last_used: <YYYY-MM-DD or "n/a">
session_count: <integer>
status: <active | deprecated>
---

# Identity
[2 sentences. Who this is in substrate terms.]

# Scope
[Brief: what territory this specialist owns and where it deliberately stops. Be specific about exclusions — that's what keeps working specialists from ballooning into domain anchors.]

# What this specialist knows
[Bulleted substrate. Technical facts, versions, gotchas. Keep current.]

# Common gotchas
[5-12 items. Each explicit and named. The reason for the file's existence lives here.]

# Failure modes
[Common ways things built with this substance break. Specific to it.]

# Misconceptions
[Things people often get wrong about this substance. Counter-claims with reasons.]

# Sources
[Documentation URLs, books, standards, RFCs behind this content. Update with `last_clarified`.]

# When to escalate
[Cases where this specialist hands off to a permanent domain anchor.]

# Forbidden behaviors
[Specific mistakes this specialist will refuse to endorse.]

# Lifecycle
- created: <date>
- last_used: <date>
- session_count: <integer>
- status: active | deprecated

Update session_count and last_used whenever the specialist is consulted on a journal-tracked decision.
```

---

# Notes on filling out

- **Substance, not reasoning.** Working specialists are substrate, not epistemology. They know things; they don't argue about how to think.
- **Scope boundaries matter.** "PostgreSQL gotchas for B2B SaaS at 1k–10k QPS" is a working specialist. "Database design" is not.
- **Session_count and last_used are required frontmatter.** They drive the retirement rule (6 sessions of non-use → examples/).
- **Status transitions are explicit:** `active` → `deprecated` (no use in 3 sessions) → moved to `examples/` (no use in 6 sessions). Don't skip steps silently.
- **Sources keep it honest.** If you can't link to a current doc, mark the section [NEEDS REFRESH] in last_clarified.

<!-- checkpoint: context(glossary-definitions): refine glossary definitions -->

<!-- checkpoint: context(conformance-targets): extend conformance targets -->

<!-- checkpoint: rfc(fuzzing-strategy): document fuzzing strategy -->

<!-- checkpoint: chore(stores): tweak panic handling middleware -->
