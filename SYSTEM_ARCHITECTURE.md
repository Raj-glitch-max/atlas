# System Architecture — Atlas

**Status:** The complete engineering architecture. Closes the architecture phase. Produced by joint-board review (architecture, distributed systems, security, platform, reliability, interface discipline, OSS maintenance) of the entire corpus reviewed *together*: context package, Phase 7 Product Definition, Phase 8 engineering package (ER/SO/INV/FM/AT), lab governance, RFC-000/001/002/003, and the Sprint 1 engineering set.
**Canonical set:** this document + `MODULE_SPECIFICATION.md` + `INTERFACE_SPECIFICATION.md` + `ENGINEERING_DECISION_RECORD.md` + `IMPLEMENTATION_MASTER_PLAN.md` + `TECHNICAL_DEBT_REGISTER.md` + `RISK_REGISTER.md` + `ENGINEERING_READINESS_CERTIFICATE.md`. Where this set refines the Sprint 1 documents (`PROJECT_MODULE_SPECIFICATION.md`, `MODULE_INTERFACE_SPECIFICATION.md`, `IMPLEMENTATION_ORDER.md`, `ENGINEERING_SPRINT_1_PLAN.md`), **this set is canonical**; the Sprint 1 documents remain as history and agree except where `ENGINEERING_DECISION_RECORD.md` records a delta (AD-014…AD-017).
**Bounds honored:** no cloud provider, database, framework, wire protocol, or code. The stack ruling (Go 1.21+, vetted libraries) is `AI_BOOTSTRAP.md` §4's, not this document's. All frozen limits carried: FM5 and FM8 unmitigated, S4 information-theoretic bound, two domains, single hop, `[HYPOTHESIS]` items unpromoted.

---

## 1. Architecture in one paragraph

Two trust domains, six modules, one artifact. Domain A issues: an **Issuance Authority** (M2) checks scope-strict-subset against an external permission source and mints a self-sufficient, tamper-evident **Delegation Record** (M1 — the stable surface and the *only* thing that crosses the domain boundary on the verification path). Domain B verifies: a pure **Verification Core** (M3) runs five separately-observable checks against injected inputs — local **Trust Material Store** (M4), local **Revocation Status Provider** (M5 — the sole spike-volatile region, permanently isolated behind a three-answer contract), and an injected clock — and returns Accept / Reject / Inconclusive→Reject `[HYPOTHESIS]` plus an unconditional decision trace. Domain A revokes: an append-only **Revocation Origin** register (M6) whose knowledge reaches M5 through a propagation channel that is deliberately undefined until the EXP-001 spike selects it. Everything the frozen requirements fix is in M1–M4/M6 and cannot be perturbed by the spike; everything the spike decides is confined to M5's realization, the propagation channel, and one opaque optional record element reserved for it.

## 2. Module graph and dependency graph

```
                         ┌─────────────────────────────┐
                         │   M1  RECORD MODEL (pure)   │◄─ nothing depends on M1's
                         │   stable surface · artifact │   dependencies: it has none
                         └──▲────────▲────────▲────▲───┘
            uses vocabulary│         │        │    │
        ┌──────────────────┘         │        │    └──────────────────┐
   ┌────┴────────┐          ┌────────┴─────┐  └────────┐      ┌───────┴──────┐
   │ M2 ISSUANCE │          │ M3 VERIFI-   │           │      │ M6 REVOCATION│
   │ AUTHORITY   │          │ CATION CORE  │           │      │ ORIGIN       │
   │ (domain A)  │          │ (domain B,   │           │      │ (domain A,   │
   └──┬───┬───┬──┘          │  pure)       │           │      │  append-only)│
      │   │   │             └─┬───┬───┬────┘           │      └──────┬───────┘
      ▼   ▼   ▼               ▼   ▼   ▼                │             │ read by
 [Permission [Time  [RevBinding [Trust  [RevStatus  [Time            ▼
  Source     port]   Source      Material port]      port]   (propagation channel —
  port]              port·E1]    port]                        DEFERRED, spike-selected,
                                   ▲        ▲                 S2/S3-bounded)
                            fulfils│        │fulfils                 │
                          ┌────────┴─┐  ┌───┴───────┐                │
                          │ M4 TRUST │  │ M5 REVOC. │◄───────────────┘
                          │ MATERIAL │  │ STATUS    │   (feeds M5's view;
                          │ STORE    │  │ PROVIDER  │    never touches M3)
                          └──────────┘  │ (volatile)│
                                        └───────────┘
   cmd/ drivers = composition roots: the only place ports meet providers.
```

Dependency rules (binding; CI-enforced by import lint — see §14):

| Rule | Statement |
|---|---|
| R1 | M1 depends on nothing in the system |
| R2 | Dependencies point toward stability: M2/M3/M4/M5/M6 → M1 only |
| R3 | M3 and M2 depend on **ports** (named contracts), never on providers; providers satisfy ports structurally and never import their consumers |
| R4 | Records flow; modules never call across the domain boundary — the record is the only cross-domain artifact on the issuance/presentation/verification path |
| R5 | No domain-B module depends on any domain-A module; M5's knowledge arrives via the deferred channel, never via a verification-time dependency |
| R6 | Nothing reachable from a verification invocation may perform egress to a shared authority (structural AP1; AT16) |
| R7 | No cycles; the graph above is the complete permitted edge set |
| R8 | Policy (R, skew tolerance, S4 ceiling) lives in M3's injected policy only — one home, set by founder scope act |

## 3. Ownership graph

| Asset | Design owner | Runtime owner | Change authority |
|---|---|---|---|
| Record definition + integrity semantics (M1) | stable-surface owner (one person) | — (artifact, not a service) | frozen-package trace + `ENGINEERING_DECISION_RECORD.md` amendment |
| Issuance boundary (M2) | track A | domain-A operator | interface spec amendment |
| Verification pipeline + policy + conformance definition (M3) | track A | relying-party operator | interface spec amendment |
| Trust material + provisioning records (M4) | track B | RP operator (out-of-band acts) | interface spec amendment |
| M5 contract + `contracttest` + degenerate realization | track B | RP operator | interface spec amendment |
| M5 *real* realization + propagation channel | unassigned until spike | RP + domain-A operators | EXP-001 outcome + founder acceptance |
| Revocation register (M6) | track B | domain-A operator | interface spec amendment |
| Drivers (`cmd/`) | track A | per-boundary operators | review (no logic allowed in them) |
| Test harness + substrate | track C | lab | lab governance |
| Decision/issuance traces (instances) | — | the invoking boundary's operator | trace shape: interface spec (append-only evolution) |

## 4. Lifecycle

**System lifecycle:** provision M4 (out-of-band trust bundle exchange) → deploy M2/M6 in domain A, M3/M4/M5 wiring in domain B → issue → present → verify (repeatedly) → revoke (occasionally) → reconstruct (any time, by anyone holding a record). No global startup ordering beyond "M4 provisioned before first meaningful verification" — and even that is honest: an unprovisioned M4 answers `absent` and M3 fails closed, which is correct behavior, not a race.

**Module lifecycles** (full table in `MODULE_SPECIFICATION.md`): M1 stateless-always; M2/M3 stateless per-request; M4 provision-then-read; M5 maintenance loop per composition (degenerate: none), read at verification, **as-of freezes during partition as a defined state**; M6 append-only forever, outliving every delegation it names.

**Delegation lifecycle** (RFC-002 §9.1 realized): NotIssued →(M2)→ Issued →(arithmetic in M3)→ Expired ∥ →(M6, observed via M5)→ Revoked. *Expired* is derived, never stored — there is no expiry bookkeeping to fail. *Revoked* is dual: authoritative at M6, observed at M5; the gap between them **is** the revocation-observability state, bounded by R (non-partitioned) and partition recovery (S4), represented honestly rather than papered over.

## 5. Control flow

Exactly three control entry points, one per driver; no background control flow in V1 except M5's realization-defined maintenance (absent in the degenerate realization).

1. **Issue** (`atlas-issue` → M2): request → permissions-of(port) → subset guard → mint instance-ID → obtain revocation-binding (port; empty pre-spike) → construct record (M1) → return record + issuance trace. Refusal short-circuits before any construction.
2. **Verify** (`atlas-verify` → M3): gather injections (M4 answer, M5 answer, time reading) → run all five checks, each producing a trace entry → route verdict (first-definitive-failure determines Reject; any-inconclusive-without-definitive determines InconclusiveRejected; all-pass determines Accept) → return verdict + full trace. The *trace* always covers all five checks; the *verdict* routing is deterministic and order-independent.
3. **Revoke** (`atlas-revoke` → M6): append(instanceID) → recorded. Idempotent on terminal state.

Control never transfers between modules except through these flows; there are no callbacks, no event subscriptions, no module-to-module notifications in V1.

## 6. Data flow

```
Permission Set ──► M2 ──constructs──► RECORD ──carried by delegate──► M3 reads
(external)          │                  │  ▲                            │
Time reading ───────┘                  │  │ (same artifact)            │
InstanceID (minted) ───────────────────┘  │                            │
RevBinding (opaque; empty pre-spike) ─────┘                            │
                                                                       │
Trust material ── out-of-band ──► M4 ──answer──────────────────────────┤
Revocation act ──► M6 register ──(deferred channel)──► M5 view ──answer┤
Time reading ──────────────────────────────────────────────────────────┤
                                                                       ▼
                                              Verdict + DecisionTrace ──► invoker
RECORD (any copy, any time) ──► third party: ValidateIntegrity + Read (reconstruction)
```

Data-flow properties: the record is write-once (immutable after M2); trust material and revocation knowledge flow only *toward* domain B and only outside verification invocations; nothing flows from M3 back to any store (verification is read-only everywhere); traces flow only to invokers.

## 7. Error flow

Full philosophy in RFC-003 §13; the architecture-level summary:

- **Errors are enumerated answers, not exceptions.** Every port and module answer set is closed; "unknown" is always an explicit member. A panic crossing a module boundary is a defect class (FM11's undefined-failure path), not an error channel.
- **Three classes, three routes:** definitive check failure → `Reject(cause)`; indeterminacy (absent material, unverifiable signature, clock beyond tolerance, stale/indeterminate revocation knowledge) → `InconclusiveRejected(cause)` `[HYPOTHESIS]` with the fallback ladder (retry, refetch, accept-with-warning) **forbidden**; issuance-side problems → `refused(cause)` with nothing created.
- **Error flow is downhill only:** providers surface honest answers upward; M3 converts them to verdicts; drivers convert verdicts to boundary responses. No layer "handles" a lower layer's honesty by hiding it.
- **Every error carries provenance** (which check, which port), and every path — including success — is traced.

## 8. State transitions

Three machines, realized exactly as RFC-002 §9 defines them (no additions):

1. **Delegation:** NotIssued → Issued → {Expired | Revoked}; terminal states terminal (INV3/INV4); record persists past both (INV9).
2. **Verification verdict (per presentation):** NotPresented → {Accept | Reject | Inconclusive→Reject `[HYP]`}; terminal per presentation; no verdict memory (a re-presentation is a fresh machine).
3. **Revocation observability (per revocation, per RP):** Not-yet-Observable → Observable; ≤ R in non-partitioned operation; ceiling = partition recovery under S4; never claimed observable in-partition (INV12).

State-transition ownership: machine 1's transitions are distributed by design (issuance at M2, expiry derived in M3, revocation at M6/M5); machine 2 lives entirely inside one M3 invocation; machine 3 lives in the M6→M5 relationship and is the only machine whose transition *mechanism* awaits the spike.

## 9. Trust boundaries and security boundaries

Realizing RFC-001 §10 (nothing added, nothing weakened):

| Boundary | Architectural realization | Enforced by |
|---|---|---|
| Conformant-verifier boundary | conformance = implementing M3's five checks + verdict rules + unconditional trace, exactly | `INTERFACE_SPECIFICATION.md` §3 (the conformance definition); AT23/AT30 |
| Two-domain boundary | R4/R5: the record is the only crossing artifact; no cross-domain module dependency | import lint + substrate ATs (AT17) |
| Partition boundary (S4) | M5's frozen as-of + M3's ceiling policy; no in-partition observability claim exists to violate | AT14 |
| Issuance boundary | M2 is the only record constructor (M1's creation surface is restricted to it); refusal creates nothing | restricted construction + AT4 |
| **Security non-boundaries (honest):** issuer key compromise (FM5), within-window replay (FM8) | no module claims or implements resistance; labels carried in every doc of this set | AT plan Non-objectives; `TECHNICAL_DEBT_REGISTER.md` (forbidden-debt: pretending otherwise) |

Additional security posture at the implementation seam (new here, traced): the integrity envelope pins its verification algorithm expectations (no attacker-chosen algorithm downgrade — a known JWS-family failure class); this is an implementation obligation recorded as task E2.M1-T3 and risk SR-2, not a new security claim.

## 10. Extension points, dependency-inversion points, plugin boundaries

The three names collapse to the same discipline: a seam is a **port with a closed contract**, and the system has exactly five. Anything else is not extensible by design (AP11).

| Seam | Kind | Opens when |
|---|---|---|
| P1 `RevocationStatusPort` (M5 realization) | plugin boundary — the *only* true plugin surface; every realization must pass `contracttest` before wiring | EXP-001 outcome + founder acceptance |
| P2 Revocation-binding element in the record + `RevBindingSource` port at M2 | dependency-inversion point protecting the stable surface: mechanism-specific per-record data (status-list reference, accumulator witness, or nothing) rides opaquely; issued empty pre-spike; interpreted only by M5 realizations | with P1 (AD-015) |
| P3 `PermissionSource` port (M2) | inversion point; stub in tests, realized per deployment | E5 runbook |
| P4 `TrustMaterialPort` provisioning procedure | operational seam (out-of-band) | E5 runbook |
| P5 `TimePort` (M2 **and** M3) | inversion point for testability (controllable clocks; AT8 skew) | always (harness) |
| Pipeline stages in M3 | extension rule, **closed in V1**: a new check (per-hop S5, posture FR10) enters only by founder scope act + spec amendment | future scope acts |
| ≥3 domains, non-SPIFFE | **explicitly not a seam** — generalizing is a V2 architecture act (TP5) | never in V1 |

## 11. Testing boundaries

Four concentric test boundaries, each cheaper than the next, mapped fully in `IMPLEMENTATION_MASTER_PLAN.md` §AT-map:

1. **Pure-unit** (no substrate, no fakes needed): M1 (mutation corpus AT20-class, reconstruction AT19/21-class), M6 (append-only), M4 (hit/miss/withdrawn).
2. **Injected-unit** (fakes for ports, controllable clock): the whole of M3 — all five checks, rollback (AT23), inconclusive routing (AT22), skew (AT8), freshness policy at arbitrary R/S4 parameters (AT13/14 *logic*); M2 with stub permission source (AT4, AT18).
3. **In-process composition** (drivers + real M1/M2/M3/M4/M6 + degenerate or fake M5): end-to-end issue→verify→revoke flows; the Sprint-1 exit demo; most of AT1–AT11, AT18–AT24, AT27–AT29.
4. **Substrate** (two real SPIRE domains, link-level partition, egress sniffer — shared with EXP-001): AT10, AT12, AT15–AT17, AT13/AT14 end-to-end, AT26 measurement. Lab discipline governs runs (two-run reproducibility, adversary blinding, pre-registration).

Contract tests cut across: `contracttest` gates every M5 realization at boundary 2 regardless of what the spike selects.

## 12. Observability boundaries

Per RFC-003 §14, with ownership resolved (AD-016): observables are **returned values**, not side effects — the module contract is emission-as-value; persistence, shipping, and retention belong to the invoking boundary's operator (driver/harness), keeping M1–M6 free of any logging dependency.

| Observable | Producer | Consumer | Non-negotiable property |
|---|---|---|---|
| DecisionTrace (per verification, **including Accepts**) | M3 | reviewer (AT30), rollback protocol (AT23), operators | unconditional; all five checks present; policy values identified |
| IssuanceTrace (per request, both outcomes) | M2 | AT4, reviewer | refusals as first-class outcomes |
| as-of freshness | M5 (inline in answers) | M3 policy; AT13/14 | mandatory on knowledge answers |
| Revocation register | M6 (the register *is* the observable) | propagation, reconstruction | append-only |
| Provisioning records | M4 | FM9 diagnosis, reviewer | one per provisioning act |
| Latency measurement (AT26) | `atlas-verify` driver, *around* M3 | V1 report | measured, never asserted; never inside M3 |

## 13. Concurrency assumptions

Stated explicitly (they were implicit until this review — AD-017):

- **M1, M3 are pure and immutable-input:** safe under arbitrary concurrent invocation by construction; no test or claim depends on serialization.
- **M4, M5, M6 are single-writer / multi-reader:** provisioning acts (M4), view maintenance (M5), and revocation appends (M6) are serialized per store; reads are concurrent-safe against a consistent view. This is an *assumption the implementation must uphold and test*, not a hope.
- **Drivers process one request per invocation in V1.** Concurrent-load behavior, throughput, and contention are **out of V1 scope** (C4: reference implementation, no production claim) — recorded as accepted debt TD-5, *not* silently assumed to work.
- **Cross-machine ordering:** none assumed anywhere. The only ordering the system relies on is M6's register order (local, single-writer) and monotonic as-of at M5 (local). No global clock beyond ER3's bounded skew tolerance; no distributed consensus exists to get wrong.

## 14. Failure isolation (blast radius)

| Failure | Blast radius | Contained by |
|---|---|---|
| M5 realization broken/absent | verifications go InconclusiveRejected `[HYP]` — availability degrades, **integrity never does** | honest-indeterminate rule + fail-closed routing |
| M4 material absent/corrupt | same containment (FM9); false-accept path is deleted, not handled | never-fetch (structural) + inconclusive routing |
| Partition (RP ↔ revocation source) | as-of freezes → staleness exceeds bound → InconclusiveRejected; valid-token verification (checks 1–4) wholly unaffected | S4 ceiling in M3 policy; gate C5 (valid-token offline verification is solved) |
| M2 down | no new issuance; every already-issued record verifies unaffected (self-sufficiency) | record's independence from its issuer's liveness |
| M6 down | no new revocations recordable; verification unaffected; existing view keeps serving with aging as-of | M6/M5 decoupling via deferred channel |
| Driver failure | one boundary's requests; no shared runtime state exists to corrupt | statelessness of M2/M3 |
| Clock fault beyond tolerance | that verification → InconclusiveRejected `[HYP]` | ER3 explicit tolerance + FM3 routing |
| Silent-failure meta-mode (FM11) | — | per-check traces + AT23 rollback + import lint (no dark path can exist un-traced) |

Enforcement: the import lint (`scripts/check-imports.sh`) is the architecture's teeth — R1–R7 fail the build, not the review.

## 15. Configuration model

- **One composition root per boundary** (the three drivers): all wiring, all configuration, nothing else configures anything.
- **M3 `Policy` {R, skew tolerance, S4 ceiling}:** injected; construction refuses unset values (an unparameterized verifier must not exist); values sourced from the founder S1/S4 scope acts at AT time, arbitrary in unit tests (AP7). Recorded in every DecisionTrace.
- **Module-level configuration is otherwise nil by design:** M1 none (a configurable record model would environment-dependent-ize the stable surface); M2 injected ports only; M4/M5/M6 none in V1 (in-memory; persistence deferred — TD-1). M5's real realization will bring composition-specific configuration with the mechanism decision, confined to the realization.
- **No configuration file format is chosen** — nothing in V1's ATs requires one; drivers may take parameters by the simplest available means. Choosing a format now would be unforced (TP6).

## 16. Cross-review findings (the corpus reviewed together)

Resolved in this set (details in `ENGINEERING_DECISION_RECORD.md`):

- **G1 — M2 had no clock.** RFC-003 gives M2 "issuance time" with no time source. Resolved: M2 consumes the same `TimePort` contract as M3 (AD-014).
- **G2 — Spike-outcome leak into the stable surface.** Status-list-class compositions need per-record status references; accumulator-class may need witness data. Without a reserved slot, the spike outcome would mutate the record — the exact rewrite AP12 forbids. Resolved: opaque optional revocation-binding element + `RevBindingSource` port, empty pre-spike (AD-015).
- **G3 — Trace sinks unowned.** Resolved: traces are returned values; persistence is the invoker's (AD-016).
- **G4 — Concurrency posture unstated.** Resolved: §13 (AD-017).
- **G5 — Instance-ID minting untestable.** Resolved: construction-injectable minter seam, default unique-per-issuance (AD-013).

Left explicitly open (cannot be honestly resolved by architecture):

- **O1 — S1–S4 values.** Founder scope acts. Architecture is parametric; EXP-001 and AT8/13/14 execution wait.
- **O2 — Instance-identity semantics** (FM5's open question). Requires a frozen-package amendment; carried opaquely everywhere until then.
- **O3 — Revocation composition.** EXP-001's to decide; P1/P2 are its landing zone.
- **O4 — Documentation defects** (founder-owned files; recording, not rewriting): `AI_BOOTSTRAP.md` §1 states "under 100 ms" as a core requirement, contradicting frozen NFR1/ER12 `[HYPOTHESIS]`/measure-only — the frozen package wins and the bootstrap line should be corrected; stale `agents/agents/` paths in RFC-000/001/002 provenance (layout since flattened); `DEVELOPMENT_RULES.md` §RFC Policy contradicts the practiced RFC track; `CLAUDE.md` layout tree lags the flattened `agents/`. All four are text defects with zero architectural effect.

## 17. Adversarial self-review — "what forces a rewrite in six months?"

Candidates examined, verdicts:

1. **Spike outcome mutates the record format.** Was the live rewrite risk; **fixed** by G2/AD-015 (opaque binding slot). Residual: a composition needing *non-opaque, verifier-interpreted* record changes would still force a spec amendment — judged unlikely (all three gate candidates fit an opaque slot) and recorded as risk ER-3.
2. **Instance-identity amendment breaks issued-record compatibility.** Mitigated by opacity: semantics can deepen without any interface or record-shape change (only the minting rule and comparison discipline change). Residual risk near zero for V1's bounded lifetime.
3. **Fail-closed hypothesis falsified at V1** (founder resolves toward availability). The Inconclusive state is a *named routing point*: redirecting Inconclusive→(something else) is a one-module, one-spec-line change plus test updates. Designed-for; no rewrite.
4. **The one-artifact reading (record = presented unit) proves wrong.** Splitting would ripple through M1/M2/M3 interfaces. Probability judged low (no frozen item forces separation; ER1+ER4 read naturally as one unit) but this is the **largest un-hedged bet in the architecture**; hedging it (a two-artifact abstraction "just in case") would violate AP11. Documented honestly as risk ER-4 rather than hedged.
5. **Concurrency hardening (post-V1) invalidates store interfaces.** Interfaces are answer-set contracts, not threading contracts; hardening changes internals and adds serialization tests. No rewrite.
6. **Library extraction for adoption (SO7 path)** — `internal/` → public relocation is mechanical by design (no util package, no hidden coupling). No rewrite.

Verdict: after AD-015, no identified failure mode forces a structural rewrite; the honest residuals are ER-3/ER-4 in `RISK_REGISTER.md`.

## 18. Provenance

- **Primary:** frozen Phase 7 + Phase 8 packages (read in full this session, including the complete AT plan and product FR/NFR/DEFERRED/ASSUMPTIONS texts); frozen `LEVEL0_1_FEASIBILITY_GATE.md`; RFC-000/001/002 (Draft) and RFC-003 (Accepted 2026-07-05); lab governance and EXP-001 plan; Sprint 1 engineering set.
- **Confidence:** High for every structural statement (each traces via RFC-003 §4 or `ENGINEERING_DECISION_RECORD.md` to frozen items); Medium for AD-015's sufficiency across all spike compositions (change-condition: EXP-001 surfacing a composition that cannot ride an opaque slot) and for the one-artifact reading (change-condition: a forcing item surfacing during E2 implementation). Predictive confidence in spike outcomes: **None** (doctrine).
- **Change-condition:** this document changes only via `ENGINEERING_DECISION_RECORD.md` entries; decisions there change only with a frozen-package trace or founder act.

<!-- checkpoint: rfc(architecture-draft): extend architecture draft -->

<!-- checkpoint: chore(issuance): simplify CLI flag configuration -->

<!-- checkpoint: chore(record): simplify error wrappers -->
