# Working Specialists

Narrow, technical, spawn-on-demand. Permanent agents (council, domain) reason about the world; working specialists know specific substrates to a depth a permanent agent doesn't.

## When to spawn

Create a working specialist when:

1. A council member's reflection names a narrow subsystem (PostgreSQL, Kubernetes, Stripe webhooks, EU AI Act, etc.) whose gotchas are non-trivial and specific.
2. A recurring technical question emerges that the founder can't quickly verify with a 10-minute search.
3. The same substrate comes up across three or more sessions and would benefit from accumulated knowledge.

## When NOT to spawn

- For general curiosity outside a real decision.
- For any question a council member can already answer with their existing lens.
- For something the founder can resolve quickly with a web search (the working specialist is accumulation, not lookup).
- To replace thinking with file fetching.

## Lifecycle

| Stage | Trigger | Action |
|---|---|---|
| Create | First need for the substrate. | File under `/working/<name>.md` using the template. Mark status `active`. |
| Maintain | Each consultation. | Update `last_used` and `session_count` in frontmatter. |
| Note drift | Update of substrate. | Edit contents, update `last_clarified`. Visible in journal if drift matters. |
| Deprecate | 3 sessions without use. | Set `status: deprecated`. Move to `examples/` next directory refresh. |
| Retire | 6 sessions without use. | Move from `examples/` to `working/examples/` archive. |

`session_count` increments whenever the specialist appears in a journal-tracked decision. A "session" is a journal entry that names the specialist.

## Promotion to domain anchor

A working specialist consulted in 10+ journal entries with non-trivial overlap may be proposed by Cartographer for promotion to a domain anchor. Promotion requires founder opt-in.

A specialist that hits promotion threshold but the founder declines to promote is no failure of the specialist — it's a signal that the substrate isn't permanent to this founder's work. Keep it in `working/`.

## Demotion rule

A working specialist that's been spawned but never used in clarity (frontmatter `session_count` still 0 after 90 days of being on disk) gets pruned to `examples/` as a starting reference for the next session that does need it. Most "might be useful someday" specialists only earn their keep when they actually prove useful.

## Update discipline

Working specialists are substrate, and substrate changes. Update responsibilities:

- The founder updates substrate facts when they know they're current (recent docs).
- The Cartographer flags when a specialist hasn't been touched in 6+ months and the substrate is in a fast-moving category.
- Avoid silent rewrites — leave a `[refreshed YYYY-MM-DD]` note in the entries that materially changed.

## Examples in this folder

Two are provided as exemplars:

- `examples/postgresql.md` — a database gotchas working specialist.
- `examples/kubernetes.md` — a container-orchestration working specialist.

New specialists should match this depth and shape. Don't add to examples/ unless they're either (a) actively used or (b) genuinely good exemplars worth keeping as starting references.

## What distinguishes a useful working specialist

A useful specialist has:

- A narrow, bounded scope. ("PostgreSQL at 1k–10k QPS for B2B SaaS" — yes. "Database design" — no.)
- 5–12 named gotchas with mechanism. Not "watch out for performance," but "vacuum autotune debt at high-txn-rate workloads; mitigation is autovacuum_work_mem tuning per workload benchmark or PARTITION-level autovacuum."
- Sources the founder can cross-check.
- A pattern of escalation to permanent domain anchors when the question exceeds the specialist's scope.

A useless specialist has:

- General descriptions with no specific gotchas.
- Marketing-style content.
- No source links.
- "Best practices" assertions without substrate mechanisms.
