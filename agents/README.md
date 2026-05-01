# agents v1.0

A persistent technical reasoning framework for Claude Code. (Renamed from "FounderOS" on 2026-06-19; see project journal session of that date for context.)

## What this is

A council of viewpoints and a roster of knowledge anchors that you load into Claude Code sessions when you're exploring startup ideas, evaluating architecture decisions, weighing technical tradeoffs, or testing product strategy. Built so the framework can argue with you — and with itself — for years.

## What this is not

- Not a multi-agent runtime. There is one conversation. The "agents" are viewpoints you invoke by name.
- Not a prompt library. Each agent file is a thinking discipline, not a snippet.
- Not auto-executing. You drive the conversation. The agents refuse to consensus.
- Not jurisdiction over the whole conversation. Load them when you want challenge. Skip them for routine work.

## Layout

```
agents/
├── README.md                       this file
├── GOVERNANCE.md                   interaction rules
│
├── council/                        5 permanent epistemologies
│   ├── empiricist.md
│   ├── red-team.md
│   ├── operator.md
│   ├── economist.md
│   └── cartographer.md
│
├── domain/                         4 permanent knowledge anchors
│   ├── distributed-systems.md
│   ├── ai-ml-systems.md
│   ├── product-engineering.md
│   └── market-buyer.md
│
├── working/                        dynamic specialists (spawn on demand)
│   ├── README.md
│   └── examples/
│
├── templates/                      scaffolding for new files
│   ├── agent-template.md
│   ├── working-specialist-template.md
│   └── journal-template.md
│
└── journal/                        decisions, dissent, evidence
    └── README.md
```

## How to use it

**For a single decision:**
Open Claude Code. Paste or reference the relevant agent file(s) plus the artifact you want reviewed. Address the agent by name: "Council, review this idea" or "Red Team, attack this architecture."

**For a longer exploration:**
Load all five council files plus relevant domain anchors into context. Run the journal protocol — every major decision gets recorded with dissent.

**For narrow technical questions:**
Spawn a working specialist under `working/` using the template. Retire it when the project ends.

**For reviewing past decisions:**
The `journal/` folder is the institutional memory. Read it before re-opening a closed decision.

## House rules (full version in GOVERNANCE.md)

1. Agents do not vote. The founder decides.
2. Disagreement stays visible. No consensus smoothing.
3. High-confidence claims require evidence shown in the response.
4. Founder overrides are recorded with reasoning, not silently applied.
5. Agents escalate to one another by name, not by category.

## Versioning

This is v1.0. Future versions evolve through the journal and the cartographer's quarterly self-review, not through silent file rewrites. See `GOVERNANCE.md`.

<!-- checkpoint: repo(conformance-targets): update conformance targets -->

<!-- checkpoint: repo(conformance-targets): document conformance targets -->

<!-- checkpoint: chore(stores): optimize attenuation rule engine -->
