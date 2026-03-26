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
│       ├── REFERENCES.md           external-skills policy + 7-item index
│       ├── council/                5 epistemologies
│       ├── domain/                 4 knowledge anchors
│       ├── working/                dynamic specialists
│       ├── templates/              scaffolds for new files
│       └── journal/                decision memory
├── archive/
│   └── research/                   archived root-level research documents (2026-07-05)
├── context/                        canonical project state (9 files)
│   ├── 00_PROJECT_CONTEXT.md       project mission, pipeline, candidate ranking
│   ├── 01_GOVERNANCE.md            agent roles, review path, commit standards, freeze policy
│   ├── 02_SYSTEM_ARCHITECTURE.md   system boundary, trust model, framework layout
│   ├── 03_DATA_FLOW.md             delegation, revocation, and reasoning flows
│   ├── 04_OPERATION_MANUAL.md      make commands, runbooks, journal process
│   ├── 05_DECISION_LOG.md          key committed decisions and outcomes
│   ├── 06_GLOSSARY.md              domain terminology
│   ├── 07_SECURITY_POLICY.md       security objectives, invariants, honest limits
│   └── 08_AI_HANDOFF.md            current state, surviving hypotheses, next actions
├── rfc/                            architecture RFC documents
│   ├── RFC-000-architecture-principles.md
│   ├── RFC-001-system-context.md
│   └── RFC-002-conceptual-domain-model.md
└── scripts/
    ├── frozen-docs.list            list of hash-pinned planning documents
    └── check-frozen-docs.sh        integrity verification script
```

## Reading order

1. `context/08_AI_HANDOFF.md` — current pipeline state and next actions. **Read this first.**
2. `context/00_PROJECT_CONTEXT.md` — what this project is and is not.
3. `context/01_GOVERNANCE.md` — how the agents system behaves and how decisions are made.
4. `agents/GOVERNANCE.md` — binding rules for agent interaction.
5. `agents/journal/` — committed decisions with verbatim dissent.

## Constraints (always on)

- Do not redesign `agents/` without founder evidence requiring it.
- Do not generate products, startups, or architectures in research sessions.
- Do not auto-import content from any external skill library; consult only by name via `agents/REFERENCES.md`.
- Every committed decision produces a journal entry at `agents/journal/<YYYY-MM-DD>-<slug>.md`.
- Do not amend frozen docs without following the amendment process in `CONTRIBUTING.md`.

## Working conventions

- The current canonical name of the framework is **`agents`**. "FounderOS" is historical only.
- The current canonical project name is **Atlas**. Previous names are preserved in `docs/project/PROJECT_HISTORY.md`.
- Confidence labels: High / Medium / Low / None. Confidence without cited evidence is forbidden.
- Files are changed explicitly, not silently. Edits always pass through the relevant governance rule.
