# Technical Debt Register — Atlas

**Status:** Canonical debt register (architecture-phase closure set). Debt here is *chosen and recorded*, never discovered-and-shrugged. Three classes: **expected** (will exist at V1 by frozen decision — not repayable inside V1's scope), **acceptable** (chosen for the C4 feasibility horizon — repayable on a schedule), **forbidden** (never acceptable; incurring it is a defect or a governance breach, not debt).
**Rule:** every acceptable-debt item names its repayment trigger. Debt without a trigger is not acceptable debt; it is a decision that must go to `ENGINEERING_DECISION_RECORD.md`.

---

## 1. Expected debt (frozen-scope; carried with labels, not repaid in V1)

| # | Debt | Source | Carried how |
|---|---|---|---|
| ED-1 | No issuer key-compromise resistance; no key-rotation contract | FM5 (unmitigated within scope) | honest-claims paragraph in every doc; AD-R11 forbids testing it away |
| ED-2 | No within-window replay resistance | FM8 | same |
| ED-3 | No in-partition revocation observability | S4 (information-theoretic) | M5 frozen as-of + M3 ceiling; AT14 asserts the *non-claim* |
| ED-4 | Single-hop, two-domain, SPIFFE-coexisting only | ER17, S5, D3/D4 | non-extension-point ruling (E5-table, `SYSTEM_ARCHITECTURE.md` §10) |
| ED-5 | Fail-closed, latency, cross-protocol interop, posture-binding remain `[HYPOTHESIS]`/deferred | NFR3/NFR1/FR9/FR10; D1/D2/D5/D6 | tested-or-deferred per AT plan; never promoted (AD-R10) |
| ED-6 | Reference implementation only — no production hardening, HA, scale, or ops claims | C4 | every claim bounded; readiness certificate repeats it |

Repayment of any ED item requires a frozen-package amendment or founder scope act — by definition not a V1 engineering activity.

## 2. Acceptable debt (chosen; scheduled)

| # | Debt | Why acceptable now | Repayment trigger | Repayment shape |
|---|---|---|---|---|
| TD-1 | M4/M5/M6 state held in memory behind persistence-agnostic interfaces | no AT requires durability across process restarts except AT12's *workload* continuity (which is the delegate's, not the store's); choosing a substrate now is unforced (AD-R07-adjacent, TP6) | first AT or operational need that requires store durability; else post-V1 | substrate lands behind the unchanged interfaces (AD-D04) |
| TD-2 | Manual, out-of-band trust-bundle and permission-source provisioning; minimal runbooks | matches the gate's standard pattern (manual bundle exchange); automation is unforced | recurring operator error during E6/E8 runs, or post-V1 adoption work | procedure automation behind P3/P4 seams |
| TD-3 | Degenerate always-Indeterminate M5 realization as pre-spike default | the honest representation of not-yet-decided (AD-007); fail-closed, never false-accepting | E7-T2 (spike outcome accepted) | real realization via the same `contracttest`-gated seam |
| TD-4 | Instance identity = opaque unique ID; FM5 semantics unresolved | minimal interim (AD-013); revocation targeting works; deeper semantics unforced until the amendment | AD-D03 (founder amendment act) | new minting/comparison rules in the AD-013 seam; no shape change |
| TD-5 | Drivers one-request-per-invocation; no concurrent-load posture beyond declared store discipline | C4 horizon — production machinery for a feasibility experiment is AP11 violation | post-V1 hardening scope act (AD-D09) | concurrency tests + hardening behind unchanged contracts |
| TD-6 | Substrate is hand-built per the EXP-001 plan, not automated/reproducible-by-script | the frozen plan prices manual build; automation is unforced for one spike + one AT campaign | a second substrate rebuild being needed (repro run, new environment) | scripted provisioning inside `tests/harness`/lab tooling |
| TD-7 | Root-level document sprawl (planning + architecture sets at repository root) | founder-directed filenames; relocation is a founder call (prior IA review exists) | founder editing pass (AD-D10 adjacent) | move into `docs/engineering/` per the IA review's proposal; zero content change |
| TD-8 | Sprint-1 set superseded-but-retained (`PROJECT_MODULE_SPECIFICATION.md`, `MODULE_INTERFACE_SPECIFICATION.md`, `IMPLEMENTATION_ORDER.md`, `ENGINEERING_SPRINT_1_PLAN.md`) | history preservation is repo doctrine; the canonical set names itself | same founder editing pass | archive move, headers already cross-linked |
| TD-9 | Known documentation defects in founder-owned files (AI_BOOTSTRAP 100 ms overclaim; stale `agents/agents/` paths in RFC provenance; DEVELOPMENT_RULES RFC-policy text; CLAUDE.md layout tree) | frozen/authored texts; recording-not-rewriting discipline | AD-D10 (founder editing/amendment pass) | one-line corrections via the amendment path where frozen |

## 3. Forbidden debt (never acceptable; a violation is a defect or breach)

| # | Forbidden item | Why absolute |
|---|---|---|
| FD-1 | Skipping or conditionalizing the decision trace (an Accept without a trace) | breaks conformance, AT23, AT30; FM11's dark corner by construction |
| FD-2 | Widening any closed answer set in code ahead of the spec amendment | interface discipline is the architecture (universal rule 5) |
| FD-3 | Network imports reachable from any verification invocation | AP1/INV7 structural guarantee; the import lint exists to make this impossible, not rare |
| FD-4 | Expressing ignorance as `NotObservedRevoked` (in any M5 realization, fake, or fixture) | the exact silent-trust failure R1 names; AP5/INV12 |
| FD-5 | Test skips without a named blocker; TODO comments as scheduling | rot mechanism; the AT unblock map is the only skip authority |
| FD-6 | Promoting a `[HYPOTHESIS]` property in code docs, trace naming, test phrasing, or claims | DR7; AD-R10 |
| FD-7 | Asserting or implying FM5/FM8 resistance anywhere | AD-R11; doctrine (confidence without evidence) |
| FD-8 | Editing frozen documents or `FROZEN.sha256` outside the amendment process | governance breach per `CONTRIBUTING.md` §4 / `01_GOVERNANCE.md` |
| FD-9 | Business logic in drivers, checks outside M3, record construction outside M2 | boundary erosion — the top of the slope to AD-R01's monolith |
| FD-10 | A "temporary" M5 realization wired without passing `contracttest` | the plugin admission rule is the seam's entire value |

## 4. Repayment schedule (consolidated)

| When | What comes due |
|---|---|
| Sprint 1 review (E5-T6) | F1/F2 confirmations close the A3/A4 provisionality of AD-012/AD-013 |
| F3 + E6 | TD-6 assessed after first substrate build (rebuild needed? → repay) |
| E7-T2 (spike accepted) | TD-3 repaid (real M5) or explicitly converted to V1's honest answer (β/γ/δ) |
| AD-D03 amendment (founder-timed) | TD-4 repaid in the minting seam |
| E8 / V1 closure | AD-D07 resolves ED-5's fail-closed item (promote or redirect); V1 report re-states ED-1…ED-6 verbatim |
| Post-V1 founder pass | TD-1, TD-2, TD-5 (scope acts); TD-7/TD-8/TD-9 (editing pass) |

Unscheduled debt discovered during implementation is not silently absorbed: it gets a register entry and either a trigger or an escalation to the founder — the same rule as `ENGINEERING_DECISION_RECORD.md`'s closing rule.
