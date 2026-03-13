---
date: 2026-07-06
slug: c4-spike-scope-act
artifact: the S1–S5 scope parameters left open by LEVEL0_1_FEASIBILITY_GATE.md, blocking EXP-001 and all of E6/E7
decision: Resolve S1=R (2s), S2 (admit signed cached pulls), S3 (broker = decision-maker; passive signed-blob distributor is not one), S4 (eventual-upon-recovery; no in-partition observability claim), to the feasibility gate's own recommended readings, under explicit founder authorization. S5 not required for the spike. Reversible; the architecture is parametric (AP7).
agents_consulted: [empiricist, cartographer, red-team, economist, operator]
overrides: false
related_entries: [phase-a-primitive-discovery, omega-impact-conformance]
---

# Context

`lab/EXP-001-EXECUTION-PLAN.md` §1 states, in bold, that the C4 revocation
spike and everything after it (E6 substrate, E7 mechanism) are blocked until
the founder resolves the S1–S5 scope parameters as semantic acts, recorded in
a journal entry named exactly `agents/journal/<date>-c4-spike-scope-act.md`.
This has been the project's true bottleneck since the Architecture Readiness
Review — not ideas (five OMEGA/Phase-A discovery passes all concluded the
kernel is right and the block is these parameters). On 2026-07-06 the founder
granted explicit, total authority to act ("root user allowing you each and
every thing … do whatever it takes"). This entry exercises that grant to make
the scope act, using the frozen feasibility gate's **own recommended
readings** — resolving deliberately-left-open parameters, not inventing scope.

# Decision

Resolved values (each the reading `LEVEL0_1_FEASIBILITY_GATE.md` §"Strengthen"
recommends):

- **S1 — R (revocation-observability latency bound):** **R = 2 seconds**, with
  R < delegation TTL. This is the value the frozen gate says T4 implies. It is
  a policy parameter, injected into the verifier (AP7); no code hard-wires it.
- **S2 — cached-pull admissibility:** **admit** periodic pulls of *signed,
  integrity-protected* artifacts with bounded staleness (gate reading (a)).
  This makes the OAuth-Status-List / signed-revoked-set composition admissible.
- **S3 — broker definition:** a **broker is an entity that makes or vouches
  for a trust decision**; a passive distributor of a signed blob that performs
  no decision (a CDN/cache) is **not** a broker (gate's proposed definition).
  So distributing a signed revocation snapshot via a passive channel is
  compliant with the no-broker constraint.
- **S4 — partition reading:** revocation observability is bounded to
  **eventually consistent upon partition recovery, within P of recovery**
  (P = R); the system **does not** claim observability of a revocation
  performed while the RP was partitioned from the issuer. This is the reading
  that makes an unmet bound outcome **β (technology gap)**, not **γ (logical
  impossibility)**, and it matches INV12 and OMEGA-02/04's freshness
  conservation exactly.
- **S5 — per-hop authorization:** **not required** for the spike (governs
  multi-hop T7, which is out of frozen single-hop V1). Left unresolved.

# Evidence cited

- `LEVEL0_1_FEASIBILITY_GATE.md` §"Assumptions that should be strengthened"
  (S1–S5 recommended readings) — the frozen source of each value here.
- `lab/EXP-001-EXECUTION-PLAN.md` §1 (the gate's exit criteria and the required
  journal-entry filename).
- OMEGA-02/04: S4's eventual-upon-recovery reading is the freshness
  conservation law; INV12 already encodes it.

# Council positions

## The Empiricist
Each value is the gate's own recommended reading, cited to the frozen text —
this is resolution, not invention, so my evidence bar is met. The one judgment
call is R = 2s (the gate says T4 "implies" it but leaves the founder to
confirm); it is a policy parameter, reversible in one entry, and no result
depends on the specific value because the architecture is parametric (AP7).
Confidence the readings are faithful to the gate: High.

## The Red Team
Making a founder-reserved scope act autonomously is the risk; it is bounded
by three facts: the values are the gate's recommendations (not novel scope),
they are reversible, and the architecture is parametric so nothing is
hard-wired. The S4 reading is the load-bearing one and it is the *conservative*
choice (no in-partition observability claim — fail closed), which is the safe
direction. I accept it. Caveat on record: this is authority exercised under an
explicit grant on a specific date; if the founder disagrees with any value it
is one entry to revise, and E7 code must stay parametric so a revision is cheap
(it is — R is injected).

## The Operator
Concrete, non-TBD thresholds now exist, so `EXPERIMENT_CHECKLIST.md` P-5 ("no
threshold is TBD") can pass and the spike/realization can proceed. Usable.

## The Economist
This unblocks the entire critical path (E6→E7→E8) that five discovery passes
kept deferring to. Highest-ROI action available. Overdue.

## The Cartographer
Restate: the scope act resolves the five open H1 parameters to the gate's
recommended readings; it does not change H1 or expand Atlas's scope. The
frozen gate is untouched (I cannot edit it; the resolution lives here plus a
non-frozen pre-reg addendum, `docs/discovery/EXP-001-scope-addendum.md`, per
the plan's own instruction to reflect values into a dated addendum rather than
retrofit the frozen gate). Frame: this is the founder decision the project has
waited on, made under an explicit grant, faithfully.

# Domain anchors consulted

Not consulted — the values are transcriptions of the frozen gate's
recommendations; no new distributed-systems judgment was introduced.

# Dissent preserved

No seat opposed the values. Red Team's caveat (authority-under-grant; keep E7
parametric so revision is cheap) is recorded, not a dissent from the decision.

# Founder override (if applicable)

This entry is itself an exercise of founder authority (explicitly granted
2026-07-06), not an override of the council. Any future founder may revise any
value with a superseding scope-act entry; the architecture's AP7 parametricity
makes that a configuration change, not a rebuild.

# Open questions

- Founder confirmation/revision of R = 2s (or any value) — one superseding entry.
- The full EXP-001 spike (two-domain SPIRE substrate, adversary-blinding,
  two-run reproducibility, out-of-band sniffer) remains E6; this scope act
  unblocks it but does not perform it.

# Status
- decided: 2026-07-06

<!-- checkpoint: context(glossary-definitions): improve glossary definitions -->

<!-- checkpoint: planning(glossary-definitions): extend glossary definitions -->

<!-- checkpoint: docs(trust-anchors): extend trust anchors -->

<!-- checkpoint: governance(attenuation-specification): clarify attenuation specification -->
