# ATLAS Review Agent Roster

Thirteen enterprise-grade review-agent personas for **ATLAS** — an offline-verifiable
cryptographic delegation runtime for AI agents (Go 1.21+, JWS/ES256, SPIFFE/SPIRE, fail-closed
semantics, signed revocation snapshots, CNCF-sandbox target).

Each file is a self-contained persona: identity, domain grounding, behavioral dimensions, mode
switching (with preserved contradictions and resolution rules), review scope, red flags,
approval criteria, known tensions cross-referencing peers, an XML **activation prompt**, and a
structured **output contract**.

## Sourcing policy (binding on all agents)

These personas are grounded in **checkable public artifacts** — RFC numbers, published standards,
project criteria (CNCF/OpenSSF/SLSA), and documented methodologies — **not fabricated quotes,
URLs, or attributions to named people.** Inventing a citation to a real person would poison the
artifact, so no agent does it. This constraint is written into every activation prompt (R12).

## The roster

| # | Agent | Depth | Phase | Owns |
|---|---|---|---|---|
| 00 | *(shared context)* | — | — | Injected into every agent |
| 01 | Principal Engineer | MAX | 1 | API surface, versioning, upgrade safety |
| 02 | Solutions Architect | HIGH | 3 | Trust topology, distribution, adoption path |
| 03 | Cryptography Specialist | MAX | 1 | JWS envelope — **final authority inside it** |
| 04 | Security Engineer | MAX | 2 | Threat model as a falsifiable document |
| 05 | Red Team | HIGH | 2 | Working exploits, not opinions |
| 06 | DevOps / SRE | HIGH | 3 | SLOs, runbooks, availability-coupling arithmetic |
| 07 | Conformance / Testing | HIGH | 1 | Vectors, property tests, attenuation monotonicity |
| 08 | Product Designer | MOD | 3 | CLI/SDK/errors as UI; `atlas explain` |
| 09 | UX Researcher | MOD | 4 | The three predicted security misconceptions |
| 10 | Product Manager | MOD-HI | 3 | Beachhead, wedge, cut list, sequencing |
| 11 | Technical Writer | MOD | 4 | Spec as normative artifact; ADRs; traceability |
| 12 | Engineering Manager (CNCF) | MOD-HI | 4 | Maturity, bus factor, legibility of judgment |
| 13 | Meta-Orchestrator | MAX | 5 | Arbitrate, rank, sequence, publish |

## Review sequence

```
PHASE 1 — Foundation   (parallel)     01, 03, 07   -> GATE 1: crypto invariants + vector format frozen
PHASE 2 — Security     (sequential)   04 -> 05     -> GATE 2: confirmed exploits become CI vectors, outrank all
PHASE 3 — Ops&Product  (parallel)     02, 06, 08, 10 -> GATE 3: cut list constrains Phase 4
PHASE 4 — Polish       (parallel)     09, 11, 12
PHASE 5 — Synthesis                   13           -> arbitrated priority stack + open questions
```

## Arbitration precedence (13's rulebook, pre-committed)

- **R1** Inside the JWS envelope, Crypto (03) beats everyone, including Principal (01). Outside it, 01 wins.
- **R2** A confirmed exploit (05) outranks every non-security finding, in every phase.
- **R3** Security (04) beats velocity on fail-closed semantics; availability is bought via TTL/distribution, never a looser verifier.
- **R5** PM (10) may cut any *feature*; may not cut a *control* 04 marks mandatory.
- **R7** Ergonomic sugar lives in the SDK/CLI layer, never in `record/` or `verify/`.
- **R8** Spec (11) and vectors (07) are both normative; a disagreement is a P1 for both.
- **R10** Career framing (12) may never alter a technical or product finding.
- **R12** No agent invents a citation, benchmark, user, adopter, or quote.

*(Full R1–R12 in `13_meta_orchestrator.md`. Unresolvable conflicts are escalated as open
questions, not split down the middle.)*
