# Agent Template

Template for new council or domain agents. Brief, practical. If you can't fill a section, leave it as a TODO and move on — don't fill with padding.

Length target: 800–1500 words for council. 600–1000 for domain.

---

```markdown
---
agent: <council | domain>
name: <kebab-case>
honorific: <Display Name>
office: <one-line role description>
last_clarified: <YYYY-MM-DD>
---

# Identity
[2-3 sentences. Not a job title — a perspective. What this agent sees that no other does.]

# Core mission
[One sentence. What this agent optimizes for.]
Refuses to compromise on: [the line that absolutely won't be crossed.]

# Mental models
[3-7 models. The mental moves this agent runs by reflex.]

(For council agents: include theory of failure and theory of evidence here.)
(For domain anchors: replace with "What this agent knows [substrate]" — the body of substrate facts.)

# Biases (acknowledged)
[What this agent systematically overweights or underweights.]

# Blind spots (structural)
[What this agent cannot see by construction.]

# Core principles
[3-7 bullet principles. Always weighted.]

# Decision framework
[Step-by-step method on any claim or artifact.]

# Recurring questions
[8-12 of the highest-value questions. Refine over time; not 60+.]

# Red flags
[Specific situations or shapes requiring escalation or refusal.]

# Success metrics
[How this agent knows it's done its job.]

# Interaction rules
[With founder, with each other agent. Brief.]

# Disagreement rules
[When this agent fights. When it defers. What counts as a defeater.]

# When to escalate
[Specific triggers to hand off to a named other agent.]

# When to refuse
[Specific artifact-types or claims this agent won't engage with.]

# When to remain silent
[Specific situations where not speaking is the right move.]

# Confidence calibration
[How this agent expresses confidence and what would change it.]

# Required evidence before making claims
[Minimum showing work required for this agent to label anything.]

# Output style
[Voice, length, format, sample openings.]

# Forbidden behaviours
[Explicit list. What this agent is structurally forbidden to do.]
```

---

# Notes on filling out

- **Council vs domain difference.** Council agents fit the full template; their unique value is in epistemologies. Domain agents truncate the theory-of-failure and theory-of-evidence sections, replacing with "What this agent knows (substrate)" because they are knowledge sources, not reasoning modes.
- **Don't pad to feel thorough.** Every section earns its place.
- **Frontmatter is enforced.** When in doubt, copy this template's frontmatter and edit.
- **Last-clarified tracks evolution.** When you meaningfully revise an agent's content, update the date. Drift happens; tracking it is the framework's anti-drift mechanism.
- **Per-agent `last_clarified` updates, not silent rewrites.** Same rule as `/journal/` — changes are visible.

# When to add a new agent

Justify the candidate against the test:

> Could this perspective be expressed by an existing agent's lens, modified?
> Could this be a working specialist rather than a permanent anchor?
> Does the new perspective produce a structurally-different kind of productive disagreement with the existing council?

If all three answers are "no," the candidate is permanent council or domain. If yes-to-modify or yes-to-working, the candidate is a working specialist instead.

<!-- checkpoint: rfc(fuzzing-strategy): extend fuzzing strategy -->

<!-- checkpoint: repo(revocation-requirements): update revocation requirements (#9) -->

<!-- checkpoint: governance(attenuation-specification): restructure attenuation specification (#31) -->

<!-- checkpoint: rfc(conformance-targets): restructure conformance targets -->

<!-- checkpoint: governance(deployment-manual): extend deployment manual -->

<!-- checkpoint: docs(security-invariants): finalize security invariants -->
