# External References Index

Skills outside `agents/` that may be consulted case-by-case. Owner: founder.
Default mode: dormant; opt-in per session. Never auto-loaded.

This file is the lightweight long-term structure for external-skill acquisition. It captures the policy every external skill is governed by and lists currently-classified items.

## Classification key

- **Core** — integrated into `agents/` permanently. *(no items currently at this level)*
- **Temporary** — invoke by name; produce journal entry per invocation; retire when absent 6 months or never load-bearing.
- **Reference** — on disk for reading; never auto-invoke.
- **Remove** — ignored; do not consult. The vast majority of any acquisition ends here.

Promotion rules follow `agents/GOVERNANCE.md` §8: a Temporary becomes Core only after 10+ journal entries cite it as load-bearing AND founder decides. Otherwise retire to Reference or Remove on demotion.

---

## Currently classified items

Long-term subset carried forward from the External Repository Audit session.

| Path (relative to repo root) | Class | Trigger phrases | Note |
|---|---|---|---|
| `claude-skills/productivity/andreessen/` | Temporary | "market-first take", "pressure-test this venture", "anti-sycophancy pass", "is there a real market here" | Anti-hedge operator; useful third opinion when Cartographer + Economist are insufficient. |
| `claude-skills/productivity/handoff/` | Temporary | "hand off this thread", "save session context" | Adjacent to journal; useful templated format. |
| `claude-skills/research-ops/skills/market-research/` | Temporary | "size this market", "TAM SAM SOM", "plan a survey", "segment candidates" | Three stdlib Python tools. Activate only when PM names a candidate market. |
| `claude-skills/research-ops/skills/research-ops-skills/` | Reference | (no trigger — read only) | Routing-pattern reference for orchestrator design. |
| `claude-skills/engineering/skill-security-auditor/` | Reference | "audit this external skill" | Run before any further external acquisition. Self-referential but correct discipline. |
| `claude-skills/engineering/{slo-architect,kubernetes-operator,chaos-engineering,feature-flags-architect}/` | Reference | (no trigger — read on demand) | Reliability substrate. Read when an infrastructure problem surfaces — not auto-loaded. |
| `claude-skills/agents/personas/{solo-founder,startup-cto}.md` | Reference | (no trigger — read only) | Craft exemplars for writing voice-rich agent files. |

Items not on this list are **Remove**. Promoting anything from Remove requires founder override + a journal entry naming (a) the trigger phrase, (b) the load-bearing role, (c) the boundary in `agents/` it would replace or extend.

---

## Stable policy (binding for all future acquisitions)

1. Never auto-import. Every external tool is reviewed against Classification key before activation.
2. A skill earns Code by 10+ journal entries of demonstrated load-bearing use. None currently qualifies.
3. A Temporary goes dormant after 6 months of non-use. Drop to Reference without ceremony.
4. A Temporary that's been used but never produced a load-bearing observation goes to Remove after 12 months.
5. The Maintainer (Cartographer per `GOVERNANCE.md` §7) re-audits this index quarterly.
6. The folder's existence on disk does not constitute inclusion in `agents/`. Only items in this index are considered.

<!-- checkpoint: chore(revstatus): harden test assertions -->

<!-- checkpoint: chore(fuzz): tweak integration test runner (#133) -->
