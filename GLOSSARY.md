# Glossary

> Scope: this document defines terms already established elsewhere in the repository. It does not introduce new concepts, rules, or structure. If a definition here ever conflicts with `GOVERNANCE.md`, `README.md`, or a journal entry, those source documents govern — this glossary should be corrected to match them, not the reverse.
>
> What belongs here: short, source-anchored definitions of terms that recur across files and are prone to being misread by someone new to the project.
> What must never be placed here: decisions, rationale, status, opinions, roadmap, or anything that changes over time. If it can become outdated on its own, it belongs in a context file or the journal, not the glossary.

---

**agents** — The current name of the persistent technical reasoning framework at `/agents/`. Renamed from "FounderOS" on 2026-06-19. Not a multi-agent runtime, not a prompt library.

**FounderOS** — Deprecated former name of the `agents` framework. May still appear in file paths, filenames, or journal entries written before 2026-06-19. Not to be reused going forward.

**Epistemology** — A council seat. An opinionated *way of thinking* (e.g., Empiricist, Red Team), not a role or persona. Five are permanent: Empiricist, Cartographer, Red Team, Economist, Operator.

**Council** — The set of five permanent epistemologies. Each responds once, in order, in their own voice, to a founder-presented artifact. Council members do not vote.

**Domain anchor** — One of four permanent knowledge territories (`distributed-systems`, `ai-ml-systems`, `product-engineering`, `market-buyer`). Consultable but not required to speak. An anchor unconsulted for 90 days is renamed `.candidate.md`.

**Working specialist / working layer** — A task-specific specialist that spawns on demand. May be promoted to a domain anchor after 10+ journal citations and founder opt-in, or retired after 6 dormant cycles.

**Journal** — The institutional decision memory. Every committed decision produces an entry. Dissent is preserved verbatim. Founder overrides are recorded. Confidence labels and change conditions are mandatory fields.

**Founder override** — A recorded instance of the founder overruling the framework's output. Always logged, never silently accepted or erased.

**Confidence label** — A mandatory High / Medium / Low / None rating attached to any claim. Confidence without cited evidence is not permitted.

**Cross-examination** — An opt-in, high-stakes-only step where one agent challenges another agent's evidence basis by name.

**Validity gate** — The first stage of the Problem Ranking Framework. Three axes — Evidence quality, Frequency, Severity — must each clear a minimum threshold (≥2, ≥1, ≥1) before a problem is scored further.

**Opportunity-fitness axes** — Three scored-higher-is-better axes in the Problem Ranking Framework: Tooling under-service, Market under-saturation, Founder advantage.

**Friction-penalty axes** — Four scored-higher-is-worse axes in the Problem Ranking Framework: Technical friction, Implementation complexity, Buyer friction, Seller friction.

**Net opportunity score ("net score")** — Opportunity fitness minus friction penalty. Range −12 to +18. Produces a sorted eligible list, not a recommendation.

**Domain Research** — The pipeline stage that identifies candidate problems without proposing solutions, products, or AI.

**Problem Ranking (R1)** — The pipeline stage that scores identified problems through the validity gate, opportunity/friction axes, and sequential council review.

**R2 / Product Management (PM)** — The next pipeline stage after Problem Ranking. As of the last recorded session, unblocked but not started; does not begin without explicit founder instruction.

**Hypothesis (H1–H8)** — A numbered, tracked claim about the framework or project, carrying a survive / carried / dropped status. Tracked in `context/02_SURVIVING_HYPOTHESES.md`.

**claude-skills** — A third-party acquired repository (~63 MB, 345+ skills) present on disk but not merged or imported into `agents/`. Subject to a classification audit, not automatic adoption.

**Classified item** — An external skill or tool that has been formally categorized (e.g., Temporary, Reference, Reject) rather than auto-imported. Adoption requires explicit, named invocation.
