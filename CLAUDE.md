# Atlas Workspace

This workspace contains the canonical tooling for Atlas.

## Top-level layout

```
/home/raj/Videos/atlas/
├── CLAUDE.md                       this file
├── agents/                         persistent reasoning framework (renamed from FounderOS, 2026-06-19)
│   └── agents/                     current framework root
│       ├── README.md               orientation
│       ├── GOVERNANCE.md           binding rules (load first)
│       ├── council/                5 epistemologies
│       ├── domain/                 4 knowledge anchors
│       ├── working/                dynamic specialists
│       ├── templates/              scaffolds for new files
│       └── journal/                decision memory
├── claude-skills/                  externally acquired third-party skills (reference-only; do not import)
└── context/                        canonical project state
    ├── 00_PROJECT_MISSION.md
    ├── 01_CURRENT_STATE.md
    ├── 02_SURVIVING_HYPOTHESES.md
    └── 03_RESEARCH_BACKLOG.md
```

## Reading order

1. `context/00_PROJECT_MISSION.md` — what this project is and is not.
2. `context/01_CURRENT_STATE.md` — pipeline position; what stage we are in.
3. `agents/agents/GOVERNANCE.md` — how the agents system behaves.
4. `agents/agents/README.md` — overview of the agents layout.
5. `context/03_RESEARCH_BACKLOG.md` — what cycle is next; what cycles follow.

## Constraints (always on)

- Do not redesign `agents/` without founder evidence requiring it.
- Do not generate products, startups, or architectures in research sessions.
- Do not auto-import content from `claude-skills/`; consult only by name.
- Every committed decision produces a journal entry at `agents/agents/journal/<YYYY-MM-DD>-<slug>.md`.

## Working conventions

- The current canonical name of the framework is **`agents`**. "FounderOS" is historical and should not be reintroduced.
- Confidence labels: High / Medium / Low / None. Confidence without evidence is forbidden.
- Files are changed explicitly, not silently. Edits always pass through the relevant governance rule.
