# Architecture Readiness Review — Atlas

**Role basis:** Principal Software Engineer + Distributed Systems Architect + Security Engineer.
**Status:** Review artifact. Not an RFC. Not frozen. Binds nothing; every recommendation here becomes binding only through the RFC path defined in RFC-000 (Architecture Review Process) and a founder act.
**Source authority (read, not regenerated):** frozen Phase 7 Product Definition package (`docs/product/`), frozen Phase 8 engineering package (`docs/engineering/` — ER1–ER17, SO1–SO8, INV1–INV12 + C-INV1, FM1–FM11, AT1–AT30), frozen `LEVEL0_1_FEASIBILITY_GATE.md` (S1–S5, spike outcomes α/β/γ/δ), `P5_FALSIFICATION_EXPERIMENT.md`, RFC-000/001/002 (Draft), `lab/EXP-001-EXECUTION-PLAN.md`, `lab/EXPERIMENT_LOG.md` (E-0001 only; no runs), `agents/journal/` (two entries; no scope-act entry exists).
**Discipline:** every claim below traces to one of the above; nothing invents scope; unmitigated modes (FM5, FM8) and the S4 limit are carried forward, not solved; `[HYPOTHESIS]` items stay hypotheses (DR7). Deliberately absent, per instruction and per RFC-000 DR9: languages, frameworks, databases, cloud providers, APIs, code.

---

## 1. Current engineering maturity assessment

| Area | Artifact(s) | State | Maturity |
|---|---|---|---|
| Product definition | `docs/product/` (10 docs) | Frozen, hash-pinned | **Mature** |
| Engineering requirements | ER1–ER17 | Frozen; full FR/NFR/C coverage table; no untraced ER | **Mature** |
| Security objectives | SO1–SO8 | Frozen; each with metric, threshold, test locus | **Mature** |
| Invariants | INV1–INV12, C-INV1 | Frozen; spike-outcome-independent | **Mature** |
| Failure model | FM1–FM11 | Frozen; distributed + adversarial; honest about FM5/FM8 | **Mature** |
| Acceptance tests | AT1–AT30 | Frozen; every ER/SO/INV/FM covered | **Mature** |
| Architecture principles | RFC-000 | **Draft — not Accepted, not Frozen** | Drafted, unratified |
| System boundary | RFC-001 | **Draft** | Drafted, unratified |
| Domain model | RFC-002 | **Draft** | Drafted, unratified |
| Scope parameters | S1–S4 | **Unresolved — no journal scope-act entry exists** | **Absent** |
| C4 feasibility spike | EXP-001 | Pre-registered, execution-planned, **not run** (log has only E-0001) | Blocked on S1–S4 |
| Mechanism architecture | — | None exists (correctly — DR9/TP6) | **Absent by design** |
| Code | — | None (`tests/` empty of product tests; no source tree) | **Absent by design** |

**Verdict:** This is a specification-complete, decision-blocked project. The specification stack is unusually rigorous — full bidirectional traceability from FR/NFR/C through ER/SO/INV/FM to AT, hash-pinned, with honest negative findings (FM5, FM8, S4). The bottleneck is **not** documentation and no further documentation of the existing kind will advance the project. The bottleneck is a chain of three acts, only the last of which is engineering:

1. **Founder scope acts** resolving S1–S4 (blocks everything; `EXP-001-EXECUTION-PLAN.md` §1).
2. **EXP-001 spike execution** (54.5–62.5 eng-h planned; resolves the revocation composition, the only unsolved component — gate C4).
3. **RFC ratification** — RFC-000/001/002 advanced from Draft per the RFC-000 state machine, and the mechanism RFCs they anticipate authorized.

One governance defect blocks item 3 and is flagged here as a finding, not fixed: `DEVELOPMENT_RULES.md` §RFC Policy states *"Do not create new RFCs under current repository governance. Existing files under `rfc/` are audit/archive material until explicitly reclassified"* — while RFC-000 §Architecture Review Process, RFC-001 §18, and RFC-002 §17 all defer decisions to *future numbered RFCs*, and `context/08_AI_HANDOFF.md` §3 treats the RFC track as advanceable. These two rules cannot both be operative. A founder act must reconcile them before any mechanism RFC can exist. (Per the working rule "do not rewrite existing documents," this review records the conflict; it does not resolve it.)

---

## 2. Missing architectural decisions

Ordered by how much downstream work each blocks. "Owner" is who can legitimately decide.

| # | Missing decision | Blocks | Owner | Smallest sufficient decision |
|---|---|---|---|---|
| D1 | **S1 — value of R** (revocation-observability latency, R < 30 s delegation TTL) | Spike pass/fail threshold; FM2/FM4 bounds; AT13 | Founder scope act | Confirm the T4-implied 2 s or set another value; journal entry per EXP-001 plan §1 |
| D2 | **S2 — cached-pull admissibility** | Whether composition 1 (Status List) is admissible; the spike's candidate set | Founder scope act | Choose (a) admit periodic signed-artifact pulls with bounded staleness, or (b) push-only |
| D3 | **S3 — broker definition** | Composition 1's compliance decidability | Founder scope act | Adopt the gate's proposed definition: broker = entity making/vouching trust decisions; passive signed-blob caches excluded |
| D4 | **S4 — partition reading** | Distinguishes spike outcome β from γ; FM1 required response; INV12 bound | Founder scope act | Bound observability to "eventual upon partition recovery, within P of recovery" (the gate's recommended strengthening) |
| D5 | **RFC-policy reconciliation** (§1 above) | All mechanism RFCs; ratification of RFC-000/001/002 | Founder governance act | Reclassify `rfc/` as the active architecture track; amend `DEVELOPMENT_RULES.md` §RFC Policy accordingly |
| D6 | **Revocation composition** (Status List / push / accumulator / none) | The Revocation Status Provider region (§7); FM2/FM4 staleness floor | **EXP-001 spike outcome** + founder acceptance | Run the spike as planned; the outcome table in the gate already maps result → disposition |
| D7 | **Delegation-instance identity** — what individuates two issuances to the same (principal, delegate, scope, expiration). FM5 records this as open; RFC-002 §17 defers it; revocation (INV4/INV6) keys to an *instance* | Any revocation mechanism; the record format; AT9–AT13 test construction | Frozen-package amendment (RFC-002 §18 change-condition) — founder | A V1 ruling that each issuance is a distinct instance with a unique issuance-time identity carried in the record. Must go through the amendment process; the current package deliberately does not fix it |
| D8 | **Clock-skew tolerance value** — ER3 requires the tolerance to be "explicit and bounded"; no document sets it | AT8; FM3 required response; deterministic expiry verdicts | Mechanism RFC (record/verifier contract) | State one bounded tolerance in the verifier contract; measure against it in AT8 |
| D9 | **Record/token relationship** — whether the presented unit (FR1/ER1) and the reconstruction record (FR6/ER4) are one artifact or two. RFC-002 relates them 1:1 (`produces`) but does not identify them | Record-format RFC; Verification Core and Reviewer interfaces | Mechanism RFC | Decide one artifact serving both roles unless a frozen requirement forces separation (AP11 favors one) |
| D10 | **Permission Set authority at issuance** — ER2/SO6 require the issuance boundary to check scope ⊆ principal's permissions; nothing states where the principal's permission set authoritatively resides or how the Issuance Authority reads it | Issuance Authority design; AT4 | Mechanism RFC | Name the permission-set source as an input the issuance boundary must hold locally at issuance time (mirrors the RP's locally-held trust material pattern) |
| D11 | **Module decomposition** — AP12 anticipates a module-boundary RFC; none exists | All implementation | Mechanism RFC (this review's §7 is the candidate input to it) | Ratify a decomposition isolating the spike-volatile region (§7) |
| D12 | **Conformance definition artifact** — FM10 defines guarantees at the "conformant verifier" boundary; SO5 lists the five required checks; no single artifact yet *defines* conformance for an implementer or reviewer (needed for AT30/SO8) | Independent review; any third-party verifier claim | Mechanism RFC | The verifier contract enumerates the required checks, their inputs, and the verdict rules; conformance = implementing exactly those |

D1–D5 require zero engineering and block everything. D6 requires ~8 eng-days. D7 is the only one requiring a frozen-package amendment and is on the critical path of the revocation region only.

---

## 3. Engineering risks

| # | Risk | Consequence if ignored | Mitigation (structural, already available) |
|---|---|---|---|
| R-A | **Spike-outcome coupling.** Any mechanism committed before EXP-001 concludes may be invalidated by outcome β/γ/δ | Rework of everything touching revocation | AP12: isolate the revocation region behind a stable internal boundary; build the requirement-fixed region first (§9). Never let the revocation mechanism's shape leak into the record format beyond an opaque instance-identity + status-check contract |
| R-B | **Single-decider bottleneck.** S1–S4, RFC ratification, D5, and D7 are all founder acts; engineering throughput is zero until they land | Indefinite stall in "dormant, governed state" | §11 packages the founder's queue as one sitting's worth of smallest-possible decisions |
| R-C | **Instance-identity leakage (D7).** The temptation is to "just add an ID" to the record. Doing so without the amendment process silently resolves what FM5 explicitly leaves open — a DR1/AP5 violation | Governance breach; an untraced load-bearing design element at the heart of revocation | Treat D7 as a blocking amendment before the record-format RFC; the interim fail-closed stub (§9 step 4) keeps the Verification Core buildable meanwhile |
| R-D | **Hypothesis laundering.** Fail-closed (ER11/SO4/C-INV1) is `[HYPOTHESIS]`. An implementation that hard-codes reject-on-inconclusive and *documents it as a guarantee* promotes the hypothesis (DR7 violation); one that silently accepts is FM11 | Either a false claim or a silent-trust failure | The Verification Core carries an explicit, named Inconclusive state (RFC-002 §9.2) whose reject outcome is implemented, tested (AT22), and *labeled* as behavior-under-test, not warranted |
| R-E | **Test-apparatus underestimation.** AT13/AT14 need partition induction; AT16 needs egress instrumentation; SO8/AT30 needs the whole verdict set reproducible by a stranger; the lab requires two-run reproducibility and adversary blinding | Acceptance testing becomes the schedule-dominant item, discovered late | The EXP-001 substrate (two SPIRE domains, firewall isolation, sniffer, latency probe — plan Phases 2–6) is deliberately reusable as the AT substrate; budget it as shared infrastructure, not spike-only |
| R-F | **Verification-path retrofit cost.** SO5 (single-check rollback) and AP13 (per-check observability) are structural: a verifier built as an opaque pass/fail cannot be retrofitted cheaply | FM11 has a dark corner to live in; AT23 unimplementable | The check pipeline is decomposed per-check with an inspectable decision trace from the first commit (§7, Verification Core) |
| R-G | **Claims drift.** FM5 (issuer key compromise) and FM8 (within-window replay) are unmitigated *by scope*; future artifacts (docs, tests, even variable names) may accidentally imply resistance | AP5/DR5 violation; a reviewer finds the project claiming what it cannot | Every engineering artifact carries the honest-claims paragraph forward, as RFC-001 §14 already does; AT plan already refuses tests for unwarranted properties |
| R-H | **Freeze-friction on iteration.** 26 frozen docs + amendment ceremony is correct for planning but will meet high-frequency engineering change | Either governance violations under deadline pressure, or paralysis | Keep the frozen set at *specification* altitude; mechanism RFCs enter the freeze only at Accepted (DR10), and code is never frozen — verified instead by `make ci` + AT harness |

---

## 4. Dependency graph

```
FOUNDER ACTS (no engineering)                     ENGINEERING (unblocked today)
─────────────────────────────                     ─────────────────────────────
 D5 RFC-policy reconciliation ──┐
 RFC-000/001/002 ratification ──┤
                                │                  AT-harness substrate design
 S1 (R) ─┐                      │                  (shared with EXP-001 Phases 2–6:
 S2      ├─ scope-act journal   │                   2 SPIRE domains, isolation,
 S3      │  entry + pre-reg     │                   sniffer, latency probe)
 S4 ─────┘  addendum            │
     │                          │
     ▼                          ▼
 EXP-001 spike (Phases 1–12) ── requires ratified governance to journal its outcome
     │
     ├─ outcome α  → revocation collapses to thin layer over standardized primitive
     ├─ outcome β  → technology gap confirmed → Level 2 narrowed to C4 (+C11 per-hop)
     ├─ outcome γ  → S4 re-scope returns to founder (already pre-empted if D4 adopts
     │               the eventual-upon-recovery reading)
     └─ outcome δ  → Level 2 as frozen
     │
     ▼
 D7 instance-identity amendment (frozen-package amendment; founder)
     │
     ▼
 MECHANISM RFC TRACK (sequential, each traced per DR1, each adversarially reviewed per DR4)
     RFC-00X module boundaries (AP12)         ← §7 of this review is the candidate input
     RFC-00X record & verifier contract       ← needs D7, D8, D9, D12; spike-INDEPENDENT
     RFC-00X revocation mechanism             ← needs D6 (spike outcome); spike-DEPENDENT
     │
     ▼
 IMPLEMENTATION (§9)
     stable region first (Record, Verification Core minus revocation, Issuance Authority)
     volatile region last (Revocation Status Provider)
     AT1–AT30 harness woven through, not appended after
```

Critical path: **S1–S4 → EXP-001 → revocation RFC → Revocation Status Provider**. Everything else — including most of the implementation — is off the critical path and can proceed in parallel once the RFC track is unblocked (D5) and the record/verifier contract RFC exists.

---

## 5. Candidate high-level architectures

Three candidate responsibility allocations, all technology-free. All three assume the RFC-001 boundary (issuance / record / verification / revocation / reconstruction inside; base identity, grant decisions, non-conformant RPs outside) — that boundary is frozen-adjacent and not re-litigated here. A fourth family — any design with a live shared-authority call on the verification path — is invalid on its face (AP1) and is not analyzed.

### Candidate A — Single component per domain

One deliverable in domain A (issues, records, revokes) and one in domain B (verifies, tracks revocation observability, reconstructs). All responsibilities co-located; the revocation mechanism is internal detail.

### Candidate B — Boundary-aligned decomposition, revocation isolated

Five parts, aligned one-to-one with the RFC-001 trust boundaries and RFC-002 concepts:

1. **Issuance Authority** (domain A) — enforces scope ⊆ permission set (refuses over-scope), binds identities and expiration, emits the record. Owns the issuance boundary (RFC-001 §10.4).
2. **Delegation Record** — an *artifact*, not a component: the single presentable, tamper-evident, self-sufficient unit (D9 resolved toward one artifact). The stable surface of AP12(a). Everything else may change; this cannot without a frozen amendment.
3. **Verification Core** (RP-side, domain B) — a decomposed check pipeline: identity binding, signature/tamper, expiry-within-stated-skew, scope integrity, revocation-observability. Each check individually forceable to failure (SO5) and individually observable (AP13) via an emitted decision trace. Carries the named Inconclusive state → reject `[HYPOTHESIS]`. Consults 4 and 5 read-only; performs zero egress.
4. **Revocation Status Provider** (RP-local) — the *only* spike-dependent part. Answers exactly one question for the Verification Core: "is this delegation instance observably revoked, and how fresh is that knowledge." Its realization (status list / push-fed store / accumulator / none) is whatever EXP-001 selects; its boundary does not move. When it cannot answer within the R/S4 bound, it reports inconclusive — it never guesses.
5. **Trust Material Store** (RP-local) — locally-held trust material, provisioned out-of-band (the gate's manual bundle exchange), read-only at verification time. Absent/corrupted material ⇒ Inconclusive path, not fetch.

### Candidate C — Embedded extension of the RP's existing verifier

Delegation checks integrated *into* the relying party's existing identity-verification mechanism as an extension point, rather than a companion component invoked alongside it. Motivated by ER9/SO7 (non-replacement): the tightest possible interoperability story.

---

## 6. Trade-off analysis

| Criterion | A — Single component | B — Boundary-aligned, revocation isolated | C — Embedded extension |
|---|---|---|---|
| Complexity | Lowest part-count, but internal coupling hides the spike-volatile region inside everything | One more boundary than A; every part traces to a requirement (AP11 satisfied — no part is unforced) | Lowest *new*-surface count, but couples to the internals of an external baseline |
| Operational cost | One deliverable per domain | Same runtime footprint as A (3–5 are RP-co-located; decomposition is structural, not topological) | Depends on baseline's extension mechanics — unknown, uncontrolled |
| Maintainability | Spike outcome β/δ or an S2/S3 re-reading forces changes across the whole RP component — **AP12 violation on its face** | Spike outcome perturbs part 4 only; record surface stable | Baseline upgrades can break the extension; maintenance hostage to external release cadence |
| Correctness | SO5 single-check rollback achievable but not structural; checks can silently merge | SO5/AP13 structural: checks are separate by construction | Check ordering/short-circuiting owned by the baseline; SO5 not guaranteed demonstrable |
| Performance | Fine (gate C9: verification clears the budget by orders of magnitude) | Same | Same |
| Security | FM11 risk: monolithic verdict path is where silent flaws live | Smallest dark-corner area; decision trace per check | **SO7 metric (a) endangered**: "baseline present and *unmodified*" — an in-baseline extension risks failing its own interop test; also strains H1's unmodified-stock-verifier clause |
| Debugging | Verdict opaque without extra work | Decision trace is the debugging surface (AP13 for free) | Debugging spans a foreign codebase boundary |
| Testing | AT23 (single-check rollback) needs internal seams that don't structurally exist | AT1–AT30 map onto part boundaries almost one-to-one; the revocation stub makes the stable region testable pre-spike | Cannot test without choosing a specific baseline — imports an unforced technology decision (DR9 violation to even specify) |
| Extensibility (multi-hop S5, ≥3 domains, later scope acts) | Rework | New checks are new pipeline stages; new compositions are new part-4 realizations | Bounded by the host's extension model |
| RFC-000 admission | **Fails AP12** (volatile region not isolated) | Passes AP1–AP13 (shown in §8) | **Fails DR9** (forces a baseline choice now) and risks SO7; also weakens AP12(a) |

**Decision: Candidate B. A and C are rejected.** A is rejected because AP12 is a gate, not a preference — a proposal whose module boundaries force a rewrite when S2/S3/R resolve is invalid, and A is exactly that proposal. C is rejected because it cannot even be *specified* without naming a baseline implementation (DR9), and because it puts SO7's own pass metric at risk. C's legitimate concern — companion-not-replacement — is fully satisfied by B: the RP's existing mechanism stays untouched and delegation verification operates alongside it, which is precisely the coexistence reading RFC-001 §9 permits and AT24 tests.

---

## 7. Recommended architecture

Candidate B, stated precisely enough to seed the module-boundary RFC (D11). Conceptual component diagram (no deployment, no protocol, no technology):

```
        DOMAIN A (principal's)                    DOMAIN B (relying party's)
┌────────────────────────────────┐      ┌──────────────────────────────────────────┐
│  Principal ──request──►        │      │            Delegate ──presents──┐        │
│  ┌──────────────────────────┐  │      │                                 ▼        │
│  │ ISSUANCE AUTHORITY       │  │      │  ┌────────────────────────────────────┐  │
│  │  · scope ⊆ perm-set guard│  │      │  │ VERIFICATION CORE (per-check       │  │
│  │    (refuse over-scope)   │  │      │  │  pipeline, zero egress)            │  │
│  │  · bind P,D,scope,expiry │  │      │  │  1 identity binding      [INV1]    │  │
│  │  · emit Record           │  │      │  │  2 signature / tamper    [INV8]    │  │
│  └──────────┬───────────────┘  │      │  │  3 expiry ± stated skew  [INV3]    │  │
│             │ produces          │      │  │  4 scope integrity       [INV8]    │  │
│             ▼                   │      │  │  5 revocation status     [SO1]     │  │
│   ═══ DELEGATION RECORD ═══════════════►│  verdict: Accept / Reject /         │  │
│   (single tamper-evident,      │      │  │  Inconclusive→Reject [HYPOTHESIS]  │  │
│    self-sufficient artifact —  │      │  │  emits DECISION TRACE   [AP13]     │  │
│    THE STABLE SURFACE, AP12a)  │      │  └───────┬──────────────────┬─────────┘  │
│                                │      │   reads  │           reads  │            │
│  Revocation decision ──────────┼──?──►│  ┌───────▼────────┐ ┌───────▼─────────┐  │
│  (origin side of the           │      │  │ REVOCATION     │ │ TRUST MATERIAL  │  │
│   revocation-information       │      │  │ STATUS PROVIDER│ │ STORE (local,   │  │
│   source; mechanism = spike    │      │  │ THE VOLATILE   │ │ out-of-band     │  │
│   outcome, S2/S3-bounded)      │      │  │ REGION (spike- │ │ provisioned)    │  │
│                                │      │  │ selected impl) │ │                 │  │
└────────────────────────────────┘      │  └────────────────┘ └─────────────────┘  │
                                        └──────────────────────────────────────────┘
   Third party / Independent Reviewer ◄── reconstructs from the RECORD alone [INV9]
```

Boundary properties:

- **Stable surface (AP12a):** the Delegation Record — its content and meaning are what conformant verifiers and third-party reviewers bind to. It changes only by frozen-package amendment.
- **Volatile surface (AP12b):** the Revocation Status Provider's *internals*. Its contract with the Verification Core — *(instance-identity) → {observably-revoked | not-observed-revoked | inconclusive} + freshness* — is fixed regardless of spike outcome; only the realization behind it varies. Under outcome γ/δ for a given parameter set, the honest realization is one that always answers inconclusive beyond the S4 bound — the architecture degrades by rejecting, not by over-claiming (AP7).
- **Conformance (D12):** "conformant verifier" = implements checks 1–5 with the verdict rules above. This gives FM10 its boundary artifact and AT30 its review target.
- **No component exists that is not forced:** issuance boundary (ER1/ER2/SO6/FM6), record (ER4/INV8/INV9), verification (ER7/SO5), revocation status (ER5/SO1/FM2/FM4), local trust material (ER7/INV7/FM9). There is no sixth part (AP11).

---

## 8. Architecture decision rationale

Per-part forcing traces (DR1 form) and principle compliance:

| Decision | Forced by |
|---|---|
| Issuance Authority as the sole creator of delegations, refusing over-scope at its boundary | ER1, ER2, ER3, INV1–INV3, SO6, FM6, AP10 |
| Record as a single artifact serving presentation and reconstruction (D9 proposal) | ER1 (single presentable unit) + ER4/INV9 (self-sufficient record) + AP11 (two artifacts unforced) |
| Verification Core as a decomposed per-check pipeline with a decision trace | SO5 (single-check rollback), AP13 (observability), FM11 (no dark corners), AT23 |
| Named Inconclusive state routed to reject, labeled `[HYPOTHESIS]` | NFR3/ER11/SO4/C-INV1, FM3, FM9, AP4, DR7, AT22 |
| Zero egress during verification; trust material and revocation status read locally | FR5/NFR2/C6, ER7, INV7, SO2, AP1, AT15/AT16 |
| Revocation Status Provider isolated behind a fixed contract; realization deferred to EXP-001 | AP12, AP7, DR3 (spike-outcome-independence), FM2/FM4 (staleness ≤ R), INV12/S4 (honest bound) |
| Trust Material Store provisioned out-of-band; absence ⇒ inconclusive, never fetch | FM9, gate C1 (manual bundle exchange is the standard pattern), ER7 |
| RP's existing identity-verification baseline untouched; core operates alongside | ER9, SO7, AT24, AP2, TP2 |
| Two-domain, single-hop only; no multi-hop or ≥3-domain accommodation built "for later" | ER17, AP6, TP5, RFC-002 §11 |

Honest-claims statement (DR5): this architecture warrants **no** resistance to issuer signing-key compromise (FM5), **no** within-window replay resistance (FM8), and **no** in-partition revocation observability (S4/INV12). The Revocation Status Provider's freshness answer is bounded by R in non-partitioned operation and by partition recovery otherwise; nothing in §7 shrinks those limits, and no implementation artifact may imply otherwise.

Spike-outcome analysis (DR3): §7 supports all four outcomes. α → part 4 realized as the thin composition found; β → part 4 unrealizable with standard parts, Level 2 narrows, stable region remains valid and testable; γ → already pre-empted if D4 adopts the eventual-upon-recovery reading; δ → stable region still stands, part 4 stays a fail-closed stub. No other part changes under any outcome.

---

## 9. Implementation order

Smallest complete steps, each verifiable, ordered by dependency and risk retirement. Steps 0–2 contain zero product code.

| Step | Work | Gate to next |
|---|---|---|
| 0 | **Founder queue (one sitting):** resolve S1–S4 as a journal scope-act entry; reconcile the RFC policy (D5); advance RFC-000/001/002 through the RFC-000 state machine | Journal entries exist; `make check-frozen` clean |
| 1 | **EXP-001 spike** exactly per `lab/EXP-001-EXECUTION-PLAN.md` (Phases 1–12). Build the substrate (two SPIRE domains, isolation, instrumentation) as *shared* infrastructure — it is also the AT harness substrate (R-E) | Outcome class journaled (α/β/γ/δ); disposition per the gate's outcome table |
| 2 | **Mechanism RFC pair:** (a) module-boundary RFC from §7 (D11); (b) record & verifier contract RFC (needs D7 amendment, D8, D9, D12). Both spike-independent; (b) can be drafted during step 1 and accepted after | Both Accepted per RFC-000 review path (adversarial review record in journal) |
| 3 | **Delegation Record + Issuance Authority.** Stable surface first. AT1–AT7, AT18–AT21, AT27–AT29 become runnable here | ATs pass on the shared substrate |
| 4 | **Verification Core** with checks 1–4 live and check 5 (revocation) stubbed to *inconclusive → reject* — an honest, shippable intermediate state, not a mock. AT8, AT15–AT17, AT22–AT24 runnable | ATs pass; decision trace demonstrably supports single-check rollback (AT23) |
| 5 | **Revocation mechanism RFC** (post-spike, spike-selected composition), then **Revocation Status Provider** realization. AT9–AT14 runnable | ATs pass including partition cases (AT14 honors S4) |
| 6 | **Latency measurement** (AT26 — measured and reported, no threshold asserted) and **independent-review package** (SO8/AT30: spec + build + decision traces sufficient for a stranger to reproduce every verdict) | AT30 executed by someone who built none of it |
| 7 | **V1 verdict:** evaluate `TECHNICAL_VALIDATION.md`'s success/failure criteria against the implementation and document the result honestly — negative included (`V1_SCOPE.md` definition of done) | Journal entry; V1 closed |

Rationale for the order: step 4 before step 5 retires the largest structural risks (R-F, FM11) while the spike-dependent region is still undecided, and produces a system whose only unimplemented check fails closed — never a system that silently skips revocation.

---

## 10. Technical debt expected at V1

Deliberate, recorded debt — none of it accidental, all of it traceable:

1. **Scope debt (by frozen decision, not shortcut):** no issuer key rotation/compromise response (FM5); no within-window replay resistance (FM8); single-hop only (S5); two domains only (ER17); SPIFFE-coexisting environments only (D4 deferral). Each carries its unmitigated/deferred label into V1 docs.
2. **Hypothesis debt:** fail-closed behavior tested and documented but not warrantable (D6 deferral, `DEFERRED.md`); latency measured, not committed (D5 deferral); cross-protocol interop and posture-binding untouched (D1, D2).
3. **Instance-identity debt:** D7 will be resolved minimally for V1 (enough for revocation to key to an instance); the full "re-issue after revocation" semantics of FM5's second sub-case remains open pending a future amendment.
4. **Operational debt:** manual, out-of-band trust-bundle and (depending on spike outcome) revocation-artifact provisioning; a hand-built two-domain test substrate rather than repeatable environment automation; single-verifier-implementation conformance evidence (D12 defines conformance, but only one implementation will exist to witness it).
5. **Margin debt:** whatever composition EXP-001 selects will have a staleness floor characterized by the spike; if that floor sits close to R, the operating margin is thin and only re-parameterization (founder act) can widen it.
6. **Reference-implementation debt (C4 constraint, by design):** V1 is a validated reference implementation, not production-hardened — no HA, no scale characterization, no hardening pass. Claiming otherwise is forbidden.

---

## 11. Questions that MUST be answered before code exists

The founder's pre-code queue, smallest sufficient form. Q1–Q6 block all product code; Q7–Q10 block only the artifacts named.

| # | Question | Blocks | Answer vehicle |
|---|---|---|---|
| Q1 | What is R? (S1 — confirm 2 s or set another value < 30 s) | EXP-001; FM2/FM4; AT13 | Scope-act journal entry |
| Q2 | Are periodic signed-artifact cached pulls admissible? (S2) | EXP-001 candidate set | Same entry |
| Q3 | Is a passive signed-blob cache a "broker"? (S3) | Composition 1 decidability | Same entry |
| Q4 | Is revocation-during-partition observability required, or bounded to eventual-upon-recovery? (S4) | β vs γ interpretability of the spike; FM1 response | Same entry |
| Q5 | Is the `rfc/` track active architecture or archive? (D5 — the `DEVELOPMENT_RULES.md` conflict) | Every mechanism RFC | Governance act + `DEVELOPMENT_RULES.md` amendment |
| Q6 | Are RFC-000/001/002 accepted as the constitution/boundary/domain-model they claim to be? | Step 2 of §9 | RFC state advancement per RFC-000 |
| Q7 | What individuates a delegation instance for V1? (D7) | Record format; all revocation ATs | Frozen-package amendment |
| Q8 | Is the presented unit the reconstruction record (one artifact), or two? (D9) | Record & verifier contract RFC | Mechanism RFC decision |
| Q9 | What clock-skew tolerance does the verifier state? (D8) | AT8; FM3 | Mechanism RFC decision |
| Q10 | Where does the Issuance Authority obtain the principal's permission set, and who operates issuance in the two-domain experiment? (D10) | Issuance Authority; AT4 | Mechanism RFC decision |

---

## Provenance

- **Primary sources:** frozen Phase 7 + Phase 8 packages; frozen `LEVEL0_1_FEASIBILITY_GATE.md`; `P5_FALSIFICATION_EXPERIMENT.md`; RFC-000/001/002 (Draft); `lab/EXP-001-EXECUTION-PLAN.md`; `lab/EXPERIMENT_LOG.md`; `agents/journal/`; `DEVELOPMENT_RULES.md`; `context/` package. All read this session; none modified.
- **Assumptions stated:** (i) the EXP-001 substrate can double as the AT substrate — High confidence, both need the identical two-domain isolated topology per the plan's Phases 2–6 and the AT plan's execution discipline; change-condition: if AT execution discipline demands a clean-room substrate separate from the spike's, the shared-infrastructure saving in §3 R-E is lost. (ii) D9's one-artifact reading is compatible with ER1+ER4 — Medium confidence; the record/verifier contract RFC must confirm no frozen item forces separation. (iii) No other repository contains prior implementation work — High; the tree has no source code.
- **Confidence:** §1–§4 findings — High (direct reading of repository state). §5–§8 recommendation — High that A and C fail RFC-000 gates as argued; Medium that §7's five-part decomposition is minimal — what would change it: the module-boundary RFC's adversarial review collapsing parts 3/5 or 3/4 without violating AP12/AP13. §9 order — High on dependency correctness; effort estimates deliberately absent except where the frozen plan supplies them. §10–§11 — High (enumerations of recorded frozen decisions and open items).
- **Change-condition:** any founder resolution of Q1–Q10, any EXP-001 outcome, or any frozen-package amendment supersedes the corresponding section of this review. This document is input to the module-boundary RFC (D11), not a substitute for it.

**Stop.** Awaiting founder decision on §11 Q1–Q6 before any further action.

<!-- checkpoint: feat(internal): implement signature validation -->

<!-- checkpoint: fix(stores): fix revocation status lookup (#83) -->

<!-- checkpoint: chore(revstatus): optimize attenuation rule engine -->

<!-- checkpoint: chore(record): simplify conformance validation (#91) -->

<!-- checkpoint: chore(record): clean boundary check -->

<!-- checkpoint: chore(sdk): tweak conformance validation -->

<!-- checkpoint: chore(verify): clean revstatus snapshot retrieval (#55) -->
