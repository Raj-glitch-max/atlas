# RFC-003 — Logical Software Architecture

## 1. Header

| | |
|---|---|
| RFC | 003 |
| Title | Logical Software Architecture |
| Status | Draft |
| Author | founder |
| Date | 2026-07-05 |
| Supersedes | none |
| Superseded by | none |

This RFC defines the logical software architecture of the system: its modules, their responsibilities, the dependency rules among them, their public interfaces stated conceptually, data and state ownership, lifecycles, module interactions, the error-propagation philosophy, observability boundaries, testing boundaries, and extension points. It satisfies RFC-000 (principles), RFC-001 (system boundary), RFC-002 (domain model), and the frozen Phase 8 engineering package, while minimizing coupling and maximizing testability. It names **no** language, framework, protocol, transport, database, cloud provider, or implementation.

## 2. Source authority

- Frozen Phase 7 Product Definition package and frozen Phase 8 engineering package (26 documents, hash-pinned in `FROZEN.sha256`): ER1–ER17, SO1–SO8, INV1–INV12 + C-INV1, FM1–FM11, AT1–AT30.
- Frozen `LEVEL0_1_FEASIBILITY_GATE.md` — scope parameters S1–S5 (all unresolved at this writing; this RFC is parametric over them, AP7) and spike outcomes α/β/γ/δ.
- RFC-000 (principles + Architecture Review Process envelope, followed here), RFC-001 (System Context — this RFC decomposes the inside of that boundary), RFC-002 (Conceptual Domain Model — every module here realizes concepts defined there). All three are unfrozen Drafts; this RFC inherits that status chain and cannot advance past them.
- `ARCHITECTURE_READINESS_REVIEW.md` (2026-07-05) — the candidate decomposition (its §7) and decision analysis (its §5–§8) this RFC formalizes; a review input, not an authority.
- **Founder instruction, 2026-07-05:** produce RFC-003 — Logical Software Architecture — defining exactly the twelve content items in §3, under the constraints in §3, and stopping after. This instruction is the single-RFC authorization required by `context/08_AI_HANDOFF.md` §3. The standing `DEVELOPMENT_RULES.md` §RFC Policy conflict (flagged in `ARCHITECTURE_READINESS_REVIEW.md` §1) is overridden for this one act by the explicit instruction; the policy text itself remains unamended and its reconciliation remains an open founder act (§20).

## 3. Scope of this RFC

This RFC defines:

- **Modules and responsibilities** (§7) — the named modules of the system and what each is answerable for.
- **Dependency rules** (§8) — which module may depend on which, in which direction, and what is forbidden.
- **Public interfaces, conceptually** (§9) — each module's operations stated in prose with named inputs, outputs, and closed answer sets. No signatures, no syntax.
- **Data ownership** (§10.1) and **state ownership** (§10.2).
- **Lifecycle** (§11) — of each module and of the domain lifecycles mapped onto modules.
- **Module interaction diagrams** (§12) — the five boundary flows of RFC-001 §12, realized as module interactions.
- **Error propagation philosophy** (§13).
- **Observability boundaries** (§14).
- **Testing boundaries** (§15).
- **Extension points** (§16).

This RFC deliberately does **not** define (each deferred to a later, separately authorized RFC, or out of scope per the frozen package):

- **Languages, frameworks, runtimes.** Deferred to a Technology Selection RFC.
- **Protocols, transports, wire formats, encodings.** The revocation-propagation channel in particular is the S2/S3-bounded, spike-selected composition — named as a deferred slot (§16 E1), never chosen here.
- **Databases, storage substrates, persistence layouts.** Where §10 assigns *ownership* of data, it says who is answerable for it, not how it is kept.
- **Cloud providers, deployment topologies, runtime placement.**
- **APIs as syntax.** §9 states interface *contracts* in prose; operation names are conceptual labels, not method names.
- **Implementation.** No code, no internal data structures, no algorithms.
- **The record's concrete field encoding.** The Record Model's content is named at the conceptual level (what a record must let a reader determine); its concrete format is the Record & Verifier Contract RFC's job (§20).

## 4. Traces

| Architectural decision | Forced by |
|---|---|
| A Record Model module: the delegation record as a single, pure, self-sufficient artifact definition with tamper-evidence semantics | ER1, ER4, FR1, FR6, NFR6, INV1, INV8, INV9, SO3, AP3, AP12(a) |
| The presented unit and the reconstruction record are one logical artifact | ER1 ("single presentable unit") + ER4 (record self-sufficiency) + AP11 (a second artifact is unforced); RFC-002 §8 `produces` (1:1) |
| An Issuance Authority module as the sole creator of records, refusing over-scope at its boundary | ER1, ER2, ER3, ER10, INV1, INV2, INV3, SO6, FM6, AP10 |
| A Verification Core module as a per-check pipeline with all inputs injected and zero I/O of its own | ER7, INV7, SO2, SO5, AP1, AP13, FM11, AT15, AT16, AT23 |
| A named Inconclusive verdict routed to Reject `[HYPOTHESIS]` | NFR3, ER11, SO4, C-INV1, FM3, FM9, AP4, RFC-002 §9.2 |
| A Trust Material Store module, RP-local, provisioned out-of-band, read-only at verification | ER7, INV7, NFR2, C6, FM9, gate C1 |
| A Revocation Status Provider module isolated behind a fixed contract; realization deferred to the EXP-001 outcome | ER5, ER6, SO1, FM2, FM4, INV12, AP7, AP12(b), DR3 |
| A Revocation Origin module at the issuing side, owning the authoritative revocation register | ER5, FR4, INV4, INV5, INV6, FM1 (the "revocation source"), RFC-001 §9 |
| Freshness policy (the R-bound and S4 ceiling) lives in the Verification Core, not in the Revocation Status Provider | FM2, FM4 (the bound is a requirement on the *verifier's acceptance*), AP7 (R is a founder parameter, set once, in one place), INV12 |
| Instance identity is carried opaquely; its semantics are deferred | FM5 (open question, verbatim), RFC-002 §7B/§17, DR1 (fixing it here would be untraced invention) |
| Time enters the Verification Core only as an injected reading with a stated skew tolerance | ER3, FM3, AT8 |
| Reconstruction is a capability of the Record Model, not a separate module | ER4, INV9, SO8, AP11 (a dedicated reviewer module is unforced) |
| No module in domain B depends on any module in domain A at verification time; the record is the only cross-domain artifact on the verification path | FR5, C6, ER7, ER8, INV7, AP1, AP6 |

## 5. RFC-000 principle compliance

- **AP1 (offline):** the Verification Core performs no I/O; every input is injected before the verdict is computed (§7 M3, §9.3, §12.2). No module it depends on may reach a shared authority during verification (§8 R6).
- **AP2 (companion):** no module issues base identity or touches the existing standard; the RP's existing identity-verification baseline appears nowhere in the module graph — the system operates alongside it (§7, §8).
- **AP3 (tamper-evidence):** integrity is a property of the Record Model artifact itself, checkable by any holder (§7 M1, §9.1); no module "vouches" for a record.
- **AP4 (fail-closed, `[HYPOTHESIS]`):** one named Inconclusive state, defined in the Verification Core, whose designed outcome is Reject; carried as hypothesis (§13, §18).
- **AP5 (honesty):** §17 carries FM5/FM8/S4 forward; the Revocation Status Provider's answer set contains `indeterminate` precisely so the architecture never manufactures certainty (§9.5).
- **AP6 (two-domain):** the module graph names exactly two domains; nothing is built for a third (§7, §16 E5 is explicitly *not* an extension point).
- **AP7 (parametric):** R, the skew tolerance, and the S4 ceiling are injected policy parameters of the Verification Core, statable for any founder-set value (§9.3, §10.2); no module bakes a value in.
- **AP8 (independent reviewability):** §19; the decision trace plus the record plus this specification suffice.
- **AP9/AP13 (adversarial, observable):** every load-bearing choice has a named observable failure (§13 sub-verdicts, §14 traces, §17 adversarial table in §"Adversarial review").
- **AP10 (no over-issuance):** the Issuance Authority is the only record creator and refuses over-scope; refusal is a first-class outcome, not an error afterthought (§9.2).
- **AP11 (minimize accidental complexity):** six modules; §4 traces each; §7 closes with the argument that no seventh exists and none of the six can merge without violating a principle.
- **AP12 (stable module boundaries):** stable surface = Record Model artifact + the two port contracts of the Verification Core; volatile surface = the realizations of the Revocation Status Provider and Revocation Origin propagation. Resolving S2/S3/R or any spike outcome changes only the volatile side (§6, §8, §16).
- **TP1–TP6:** one trade-off arises and is resolved in §13 (fail-closed vs availability) by TP4, recorded there; TP6 governs every deferral in §3.
- **DR1–DR10:** §4 (DR1); §3 deferrals (DR9); §"Adversarial review" (DR4); §17 (DR5); §"Scope statement" (DR6); §18 (DR7); §19 (DR8); DR2 — no principle is restated or re-ranked; DR3 — §6; DR10 — §21 freeze posture.

## 6. Spike-outcome analysis

The module graph is spike-outcome-independent by construction: the spike decides the *realization* of exactly one contract (the Revocation Status Provider's, and the propagation relationship between it and the Revocation Origin), never the contract itself.

- **α (composition-works):** the Revocation Status Provider is realized as the thin composition the spike selects; no other module changes.
- **β (technology-gap):** no standardized realization exists; the Revocation Status Provider stands as its honest degenerate realization — one that answers `indeterminate` for every query — and the Verification Core's policy turns that into Reject `[HYPOTHESIS]`. The stable region remains valid, testable, and complete; the gap is visible, not hidden (AP5).
- **γ (logical-impossibility at S4):** already carried as the S4 ceiling in the Verification Core's freshness policy; no claim anywhere in this architecture exceeds it.
- **δ (unresolvable):** as β; additionally the founder's re-scope act lands entirely inside the volatile surface.

No outcome forces a change to the Record Model, the Issuance Authority, the Verification Core's pipeline, the Trust Material Store, or any dependency rule. This RFC precludes no outcome.

## 7. Modules and responsibilities

Six modules. Each is named, given a single-sentence charter, and an enumerated responsibility list. "Module" means a logical unit with one owner, one responsibility set, and a public interface — not a process, deliverable, or repository directory (deployment and layout are deferred).

### M1 — Record Model *(stable surface; shared; pure)*

**Charter:** defines what a Delegation Record *is* — its conceptual content, its integrity semantics, and its self-sufficiency — as a pure artifact definition with no I/O, no clock, and no dependency on any other module.

Responsibilities:
- Define the record's conceptual content: principal identity, delegate identity, scope, expiration, issuance time, an **opaque instance identity** (semantics deferred — FM5, §20 Q2), and an integrity envelope binding all of the above.
- Provide **integrity validation**: given a record and trust material, determine `intact` or `altered` — the tamper-evidence check any holder (verifier or third party) can run (INV8).
- Provide **reconstruction reading**: given an intact record, yield who delegated to whom, with what scope, at what time (INV9). Reconstruction by an independent reviewer is this capability exercised by a third party; it needs no other module.
- Guarantee nothing about validity: M1 answers "is this record what its issuer created, and what does it say" — never "should it be accepted." Acceptance is M3's.

### M2 — Issuance Authority *(domain A)*

**Charter:** the sole creator of Delegation Records; the module that realizes the issuance boundary (RFC-001 §10.4).

Responsibilities:
- Accept a delegation request naming a principal, a delegate, a requested scope, and a requested expiration.
- Obtain the principal's Permission Set through the **Permission Source port** (realization deferred — §16 E4, §20 Q4) and enforce requested scope ⊂ Permission Set; **refuse** otherwise, creating nothing (SO6, FM6).
- Bind principal, delegate, scope, expiration, issuance time, and instance identity into a record via M1; the record is complete and self-sufficient at creation — no later enrichment.
- Support issuance to ephemeral delegates: nothing in the issuance path requires the delegate to have a long, statically-provisioned lifetime (ER10).
- Emit an issuance trace (accepted or refused, with reason) — §14.

### M3 — Verification Core *(domain B, RP-side; pure decision logic)*

**Charter:** computes the verification verdict for one presented record from injected inputs only; the module at which the system's guarantees are defined (the conformant-verifier boundary, FM10).

Responsibilities:
- Execute the five required checks as **separately named pipeline stages** (SO5): (1) identity binding [INV1], (2) signature/tamper via M1 integrity validation [INV8], (3) expiry against the injected time reading within the stated skew tolerance [INV3, ER3], (4) scope integrity [INV8; subset-ness is issuance's job, not re-derived — RFC-002 §9.2], (5) revocation status against the injected revocation answer under the freshness policy [SO1, FM2/FM4].
- Own the **freshness policy**: the R-bound (S1 parameter) and the S4 partition ceiling are policy parameters held here, in exactly one place. A revocation answer whose freshness is outside the bound is treated as indeterminate.
- Produce one of three verdicts — **Accept**, **Reject(cause)**, **Inconclusive(cause) → Reject** `[HYPOTHESIS]` — per RFC-002 §9.2, plus a **decision trace** (§14) for every invocation, including Accepts.
- Perform **zero I/O**: trust material, revocation answer, time reading, and policy parameters are injected; M3 can neither fetch nor emit anything but its verdict and trace (AP1 made structural).
- Define **conformance**: a conformant verifier is exactly an implementation of this module's checks and verdict rules (resolves ARR D12; gives FM10 its boundary artifact and AT30 its target).

### M4 — Trust Material Store *(domain B, RP-local)*

**Charter:** holds the trust material the relying party already possesses, and answers read-only queries for it at verification time.

Responsibilities:
- Hold trust material provisioned **out-of-band** (the manual bundle-exchange pattern the gate confirms as standard; the provisioning procedure is operational, §16 E4).
- Answer `material` or `absent` — never fetch on miss (FM9; a miss flows to M3 as an Inconclusive cause).
- Record provisioning events (§14).

### M5 — Revocation Status Provider *(domain B, RP-local; volatile surface)*

**Charter:** answers, for a delegation instance, whether its revocation is observable at this relying party, and how fresh that knowledge is. The only module whose realization the EXP-001 spike decides.

Responsibilities:
- Maintain the RP-local view of revocation knowledge by whatever mechanism the spike-selected composition provides (status-list refresh, push-fed store, accumulator, or — under outcome β/δ — nothing).
- Answer status queries from the closed set in §9.5, always accompanied by an **as-of freshness** disclosure; never convert ignorance into `not-observed-revoked` (the honest-indeterminate rule, AP5).
- Contain **no policy**: it reports knowledge and freshness; whether that satisfies R is M3's judgment. This keeps every candidate realization interchangeable behind one contract (AP12).

### M6 — Revocation Origin *(domain A side)*

**Charter:** owns the authoritative revocation register — the recording of the one-way, terminal revocation of a specific delegation instance (INV4, INV6).

Responsibilities:
- Accept a revocation act targeting exactly one instance identity; record it **append-only** (one-way, terminal — INV4).
- Leave underlying identities untouched (INV5) and sibling delegations unaffected (INV6) — structurally guaranteed because the register knows only instances, not identities.
- Make revocation information available to the propagation channel. The channel itself — push, pull, cached-pull, brokered — is the S2/S3-bounded, spike-selected composition and is **not defined here** (§16 E1).

### Why exactly six

No seventh module is forced: reconstruction is M1 exercised by a third party; the decision trace is M3 output, not a component; the clock is an injected reading, not a module; the RP's existing identity-verification baseline is environment, not system (AP2). No two can merge: M1 into anything sacrifices the pure stable surface (AP12a); M3+M5 puts the spike inside the requirement-fixed region (AP12 violation — the rejected Candidate A of the readiness review); M2+M6 co-locates issuance and revocation authority without a forcing trace and widens the issuance boundary (AP11); M4+M5 merges two stores with different owners, lifecycles, and volatility (§10, §11).

## 8. Dependency rules

```
                    ┌─────────────────────────┐
                    │   M1  Record Model      │  ← depends on nothing
                    └───┬───────┬───────┬─────┘
             depends on │       │       │ depends on
        ┌───────────────┘       │       └───────────────┐
   ┌────┴─────┐          ┌──────┴──────┐          ┌─────┴─────┐
   │ M2       │          │ M3          │          │ M6        │
   │ Issuance │          │ Verification│          │ Revocation│
   │ Authority│          │ Core        │          │ Origin    │
   └────┬─────┘          └──┬───┬───┬──┘          └───────────┘
        │ port              │   │   │ ports (contracts, not modules)
        ▼                   ▼   ▼   ▼
  [Permission         [Trust    [Revocation   [Time
   Source port]        Material  Status        reading +
                       port]     port]         policy params]
                         ▲         ▲
                 fulfils │         │ fulfils
                   ┌─────┴───┐ ┌───┴─────┐
                   │ M4      │ │ M5      │   (M4, M5 also depend on M1
                   │ Trust   │ │ Revoc.  │    for record/instance vocabulary)
                   │ Material│ │ Status  │
                   │ Store   │ │ Provider│
                   └─────────┘ └─────────┘
```

Rules (binding on all later RFCs and implementation):

- **R1 — M1 depends on nothing.** It is the stable surface; anything it depended on could destabilize the record's meaning (AP12a).
- **R2 — Dependencies point toward stability.** M2, M3, M4, M5, M6 may depend on M1. Nothing depends on M2, M4, M5, or M6 except as stated below. Volatile realizations are leaves.
- **R3 — M3 depends on ports, never on providers.** M3 names three contracts — trust-material, revocation-status, time-reading — and knows no realization. M4 and M5 *fulfil* ports; M3 cannot name them. This is what makes the spike outcome invisible to M3 (AP12b) and M3 testable without any substrate (§15).
- **R4 — Records flow; modules do not call across domains.** The only thing that crosses the two-domain boundary on the issuance/presentation path is the record artifact itself. M2 and M3 share no call relationship, only M1's vocabulary.
- **R5 — No module in domain B may depend on any module in domain A.** M5's *knowledge* may originate from M6 via the deferred propagation channel, but no verification-time dependency exists (AP1); under partition, M5 still answers (with degraded freshness), and M3 still verdicts.
- **R6 — No dependency of M3, M4, or M5 may perform egress to a shared authority during a verification.** AT16's zero-egress assertion is a rule of the graph, not a test hope. (M5's *maintenance* activity — e.g., a refresh under an S2-admitting scope act — is not part of any verification invocation and is bounded by the S2/S3 scope acts.)
- **R7 — No cycles.** The graph above is the complete set of permitted edges.
- **R8 — Policy has one home.** R, the skew tolerance, and the S4 ceiling live in M3's injected policy; no other module may hold or interpret them (prevents divergent enforcement, FM2/FM4).

## 9. Public interfaces (conceptual — prose contracts, closed answer sets, no syntax)

Every port and module answer is drawn from a **closed, enumerated set**. No interface may answer outside its set; "unknown" is always an explicit member, never an omission (FM11, §13).

### 9.1 M1 — Record Model

- **validate-integrity**(record, trust material) → `intact` | `altered`. Any alteration after creation yields `altered` (INV8). Callable by anyone holding the record and trust material — verifier or third party alike.
- **read**(intact record) → the record's assertions: principal, delegate, scope, expiration, issuance time, instance identity (opaque). Defined only on intact records.
- Reconstruction (ER4/INV9) = validate-integrity + read, exercised by a third party with no access to the original verifier's state.

### 9.2 M2 — Issuance Authority

- **issue**(principal, delegate, requested scope, requested expiration) → `record` | `refused(cause)`.
  - Guard: requested scope ⊂ principal's Permission Set, as answered by the Permission Source port; otherwise `refused(over-scope)` and **nothing is created** (SO6, FM6).
  - Postcondition on `record`: complete, integrity-enveloped, self-sufficient at creation; carries a fresh instance identity.
- **Permission Source port** (consumed): **permissions-of**(principal) → `permission set` | `unavailable`. On `unavailable`, issuance refuses (`refused(permissions-unavailable)`) — the issuance-side analogue of fail-closed, forced by SO6's refusal contract, not by the NFR3 hypothesis.

### 9.3 M3 — Verification Core

- **verify**(presented record, injected: trust-material answer, revocation-status answer, time reading, policy {R, skew tolerance, S4 ceiling}) → verdict + decision trace.
  - Verdict set: `Accept` | `Reject(cause)` | `Inconclusive(cause) → Reject` `[HYPOTHESIS]`.
  - Causes are named per check (§13's taxonomy); every verdict names the check(s) that produced it.
- **Time-reading port** (consumed): **now**() → a time reading. The verifier's skew tolerance relative to the issuer is an explicit, bounded policy parameter (ER3); a reading that cannot support a deterministic expiry verdict within tolerance is an Inconclusive cause (FM3).

### 9.4 M4 — Trust Material Store (fulfils the trust-material port)

- **trust-material-for**(trust domain) → `material` | `absent`. Never fetches on miss; `absent` is an honest answer that M3 turns into Inconclusive (FM9).
- **provision**(material) — out-of-band operational act, never invoked from a verification (§11, §14).

### 9.5 M5 — Revocation Status Provider (fulfils the revocation-status port)

- **status-of**(instance identity) → exactly one of:
  - `observably-revoked (as-of T)` — this RP's view contains a revocation of this instance, known as of T;
  - `not-observed-revoked (as-of T)` — this RP's view, current as of T, contains no revocation of this instance;
  - `indeterminate` — this RP cannot currently answer (no view, corrupted view, or — under outcome β/δ — no realization exists).
- The as-of disclosure is mandatory: M5 reports knowledge and its age; M3 judges whether that age satisfies R and the S4 ceiling (R8). M5 never claims in-partition observability — during a partition, its as-of simply stops advancing, which is the honest, information-theoretically correct signal (INV12).

### 9.6 M6 — Revocation Origin

- **revoke**(instance identity) → `recorded`. Append-only; recording the same instance again is a no-op on an already-terminal state (INV4). No operation un-revokes.
- **revocations-since / current-view** (conceptual): the register's content is readable by the propagation channel. The channel itself is deferred (§16 E1).

## 10. Data ownership and state ownership

### 10.1 Data ownership

| Data | Owner | Others' access | Notes |
|---|---|---|---|
| Record *definition* (content semantics, integrity semantics) | M1 | all modules use the vocabulary | Changes only by frozen-package amendment + M1 revision (stable surface) |
| Record *instances* | no runtime owner after issuance | any holder | The record is self-sufficient by construction (INV9); it is data that carries its own verifiability, not data a module custodians |
| Principal's Permission Set | **external** (behind the Permission Source port) | M2 reads at issuance | The system never owns it (ARR D10; realization deferred) |
| Trust material (RP's copy) | M4 | M3 reads via port | Provisioned out-of-band; read-only at verification |
| Authoritative revocation register | M6 | propagation channel reads | Append-only (INV4) |
| RP-local revocation view | M5 | M3 reads via port | A *view*, never authoritative; always as-of-stamped |
| Decision traces | the RP (M3's invoker) | independent reviewer reads | §14; input to AT23/AT30 |
| Issuance traces (incl. refusals) | domain A operator (M2's invoker) | reviewer reads | §14 |

### 10.2 State ownership

- **Delegation lifecycle state** (RFC-002 §9.1: Issued / Expired / Revoked) is deliberately **not held by any single module** — this is a load-bearing choice:
  - *Expired* is **derived**, not stored: record's expiration (M1 data) + injected time reading + skew policy (M3). No state transition event exists to miss (INV3 monotonicity is arithmetic, not bookkeeping).
  - *Revoked* is owned **authoritatively by M6** (the register) and **observationally by M5** (the view). The gap between them is exactly the revocation-observability state of RFC-002 §9.3, bounded by R and S4 — the architecture represents the gap honestly instead of pretending a single global state exists (INV12).
- **Verification verdict state** is per-presentation and **ephemeral in M3**: no verdict is remembered; a new presentation is a new verification (RFC-002 §9.2). What persists is the decision trace, owned by the invoker.
- **Policy state** (R, skew tolerance, S4 ceiling) is held by M3's invoker as injected configuration, set by founder scope act, changed only by a new scope act (AP7).
- **M3, M1 are stateless.** All state in the system lives in M4 (trust material), M5 (view + as-of), M6 (register), and the operational traces.

## 11. Lifecycle

### 11.1 Module lifecycles

| Module | Comes into being | Active | Quiescent | Retired |
|---|---|---|---|---|
| M1 | with the Record & Verifier Contract RFC | whenever any module runs | n/a (pure) | only by frozen-package amendment |
| M2 | provisioned in domain A | per issuance request | between requests; holds no session state | when domain A stops issuing; existing records remain valid (self-sufficiency) |
| M3 | instantiated by the RP with injected policy | per presentation | stateless between presentations | replaceable at any time; conformance is the contract, not the instance |
| M4 | first out-of-band provisioning | read at each verification | between provisionings | material withdrawal = subsequent verifications go Inconclusive (fail-closed path, not an error) |
| M5 | realization chosen after EXP-001 | maintenance per its composition; read at each verification | during partition its as-of freezes — a *defined* lifecycle condition, not a fault | swapping realizations is invisible to M3 (R3) |
| M6 | provisioned in domain A | per revocation act | register persists append-only | register outlives any delegation it names (reconstruction, INV9) |

### 11.2 Domain lifecycles mapped to modules

- **Delegation** NotIssued →(M2 issue)→ Issued →(arithmetic in M3)→ Expired, or →(M6 revoke, observed via M5)→ Revoked. Terminal states per INV3/INV4; the record artifact survives both for reconstruction (RFC-002 §9.1).
- **Revocation observability** Not-yet-Observable →(propagation, deferred mechanism)→ Observable at this RP; ceiling = R in non-partitioned operation, partition recovery otherwise (RFC-002 §9.3). Lives entirely in the M6→M5 relationship.
- **Verification** NotPresented →(M3)→ Accept | Reject | Inconclusive→Reject `[HYP]`; terminal per presentation (RFC-002 §9.2).

## 12. Module interaction diagrams

Conceptual sequence flows; no transport, no message format. `──record──►` denotes the artifact moving by whatever means; `─?─` denotes the deferred propagation channel.

### 12.1 Issuance (RFC-001 §12.1)

```
Principal        M2 Issuance          [Permission     M1 Record
(or delegate       Authority           Source port]    Model
 on its behalf)      │                     │             │
   │──request───────►│                     │             │
   │                 │──permissions-of────►│             │
   │                 │◄─permission set─────│             │
   │                 │  (or: unavailable → refused)      │
   │                 │─ scope ⊂ set? ──┐                 │
   │                 │   no → refused(over-scope),       │
   │◄──refused───────│        nothing created            │
   │                 │   yes ──────────┴──create────────►│
   │◄────record──────│◄───────────record─────────────────│
   │                 │  issuance trace emitted (§14)
```

### 12.2 Presentation + verification (RFC-001 §12.2–12.3)

```
Delegate            M3 Verification      M4 (trust     M5 (revocation   [time
   │                   Core               port)         port)            port]
   │──record─────────►│ (RP invokes M3 with injections gathered first:)
   │    (crosses      │──trust-material-for──►│              │             │
   │     two-domain   │◄─material | absent────│              │             │
   │     boundary)    │──status-of(instance)────────────────►│             │
   │                  │◄─observably-revoked | not-observed   │             │
   │                  │   | indeterminate  (as-of T)─────────│             │
   │                  │──now?──────────────────────────────────────────────►
   │                  │◄─time reading──────────────────────────────────────│
   │                  │ pipeline: ①binding ②integrity(M1) ③expiry±skew
   │                  │           ④scope integrity ⑤revocation (R, S4 policy)
   │◄──verdict────────│ Accept | Reject(cause) | Inconclusive(cause)→Reject[HYP]
                      │ decision trace emitted for every invocation (§14)
```

*Zero egress occurs inside the dashed lifetime of a verification; all port answers are local (R5, R6, AT16).*

### 12.3 Revocation + observability (RFC-001 §12.4) and partition (§12.6)

```
Revoking actor     M6 Revocation        (propagation      M5 Status        M3
   │                  Origin             channel —         Provider
   │──revoke(inst)──►│                   DEFERRED)            │
   │◄──recorded──────│ append-only       │                    │
   │                 │═══register═══─?──►│ view updated;      │
   │                 │                   │ as-of advances     │
   │                 │              ┌────┴────┐               │
   │                 │              │PARTITION│ as-of freezes │
   │                 │              └────┬────┘               │
   │                 │                   │  status queries during partition:
   │                 │                   │  honest stale as-of / indeterminate
   │                 │                   │            │──status──►(M3 applies
   │                 │                   │            │            R + S4 policy
   │                 │                   │            │            → Inconclusive
   │                 │                   │            │            → Reject [HYP])
```

### 12.4 Reconstruction (RFC-001 §12.5)

```
Independent reviewer          M1 Record Model
   │──validate-integrity(record, trust material)──►│
   │◄──intact | altered───────────────────────────│
   │──read(intact record)─────────────────────────►│
   │◄──principal, delegate, scope, time───────────│
   (no other module, no verifier runtime state, no privileged input — INV9, SO8)
```

## 13. Error propagation philosophy

Principles, binding on every module and every later RFC:

1. **Verdicts are values, not exceptions.** Every cross-module answer is a member of a closed, enumerated set (§9). No module signals failure by absence, timeout-as-implicit-no, or out-of-band fault channels that bypass the verdict path. A condition the answer set cannot express is a design defect, fixed by amending the contract — never by overloading an existing member.
2. **Three error classes, three destinies.**
   - *Definitive negative* — a check ran and failed (expired, altered, revoked-and-observable, binding mismatch): → `Reject(cause)`, terminal, named check in the trace.
   - *Indeterminate* — a check could not conclude (trust material `absent`, revocation `indeterminate` or too stale for R, time reading beyond skew tolerance): → `Inconclusive(cause) → Reject` `[HYPOTHESIS]`. The **fallback ladder is forbidden**: no retry-with-network, no accept-with-warning, no downgrade to a weaker check (FM9, FM11).
   - *Refusal at issuance* — over-scope or permissions-unavailable: → `refused(cause)`, nothing created, no partial record ever exists (SO6).
3. **Errors carry provenance.** Every cause names the check or port that produced it. An error that cannot be attributed is FM11's dark corner and is forbidden.
4. **No error crosses a boundary as ambiguity.** M5 saying `indeterminate` is not an error *of* M5 — it is an honest answer M3 is required to handle. The only malformed behaviors are answering outside the set and converting ignorance into confidence.
5. **The one trade-off, recorded per RFC-000:** fail-closed (reject on indeterminate) versus availability (never reject a valid delegation presented during a partition). Resolved by **TP4**: the fail-closed path is *present, structural, and testable* (the Inconclusive state exists and routes to Reject) but is carried as `[HYPOTHESIS]` — V1 tests and documents the behavior (AT22) rather than warranting it. This RFC does not silently prefer availability, and does not promote fail-closed to established.

## 14. Observability boundaries

Observability is a structural property of the decision paths (AP13), passive everywhere (observation never alters a verdict), and bounded as follows:

| Boundary | Artifact | Content | Consumers |
|---|---|---|---|
| M3, per verification (including Accepts) | **Decision trace** | which checks fired, in what order, each check's inputs-digest and outcome, the injected as-of freshness and time reading, final verdict + cause | AT23 (single-check rollback), AT30/SO8 (independent reproduction), operators |
| M2, per issuance request | **Issuance trace** | accepted/refused, cause, scope-subset determination | AT4, reviewer |
| M4, per provisioning act | **Provisioning record** | what material, when, by what out-of-band act | FM9 diagnosis, reviewer |
| M5, per answer | **as-of disclosure** (inline in the answer, §9.5) | freshness of the view backing the answer | M3 policy; AT13/AT14 |
| M6, always | **the register itself** (append-only) | every revocation act, ordered | reconstruction, AT9–AT11 |

Rules: the decision trace is emitted unconditionally — an Accept without a trace is non-conformant, because SO5's rollback test and SO8's reproduction need Accepts to be as inspectable as Rejects. No trace contains privileged material a reviewer cannot hold (AP8). Traces observe hypothesis behaviors (the Inconclusive→Reject transition) without promoting them (DR7): the trace records *what happened*, the specification records *what is warranted*.

## 15. Testing boundaries

The dependency rules were drawn to make the acceptance plan executable at the cheapest possible boundary:

| Boundary | What is tested there | ATs |
|---|---|---|
| **M1 alone** (pure) | integrity over the mutation set (bit-flips, substitutions, truncation, reordering); reconstruction reading; third-party self-sufficiency | AT1, AT5, AT19, AT20, AT21 |
| **M3 alone** (pure, all inputs injected) | every check, every verdict, every cause — *without any substrate*: forced single-check failures (rollback), forced indeterminates, skew at/beyond tolerance, staleness vs R, S4-ceiling honoring | AT2, AT3, AT6, AT7, AT8, AT9*, AT11*, AT13*, AT14*, AT22, AT23 (*policy logic at M3 level; end-to-end at system level) |
| **M2 + stub Permission Source port** | over-scope refusal; permissions-unavailable refusal; ephemeral-delegate issuance | AT4, AT18 |
| **M5 contract** (realization-independent test suite) | closed answer set; mandatory as-of; honest-indeterminate; the degenerate always-indeterminate realization passes the *contract* (and drives Inconclusive) — proving β is representable | supports AT9–AT14 |
| **System, two-domain substrate** (shared with EXP-001, plan Phases 2–6) | offline verification with authorities dark; zero-egress instrumentation; cross-domain accept/reject; revocation without restart; partition-at-revocation | AT10, AT12, AT15, AT16, AT17 |
| **Boundary/inspection** | no base-identity issuance; no standard amendment; two-domain discipline; baseline coexistence | AT24, AT27, AT28, AT29 |
| **Measurement + review** | latency measured-not-committed; independent reviewer reproduces every verdict from spec + build + traces | AT26, AT30 |

Consequences of the rules: because M3 is pure (R3), the entire verification logic — including the hypothesis behaviors — is testable before any substrate exists and before the spike runs; the substrate is needed only for the genuinely environmental assertions. Because M5 is contract-bounded, its test suite is written once and reused against every candidate realization the spike attempts — the spike's compositions plug into the same seam production will use.

## 16. Extension points

An extension point is a place where a *future founder act* lands without destabilizing the stable surface. Each is named with the act that opens it. Nothing is pre-built behind any of them (AP11, TP6).

- **E1 — Revocation composition slot.** The M5 realization and the M6→M5 propagation channel. Opened by: EXP-001 outcome + S2/S3 scope acts + a Revocation Mechanism RFC. The M5 contract (§9.5) and R8 are what make this slot safe.
- **E2 — Instance-identity semantics.** M1 carries instance identity opaquely; M5/M6 key to it without interpreting it. Opened by: the frozen-package amendment FM5's open question requires (ARR D7). When resolved, only M1's definition deepens; no contract changes shape.
- **E3 — Pipeline stages.** M3's checks are ordered, separately-named stages; a future check (per-hop authorization if S5 returns multi-hop to scope; posture re-attestation if FR10/ER14 is ever confirmed) enters as a new stage with its own port, cause set, and trace line. Opened by: the corresponding founder scope act; both are out of V1 and nothing anticipatory is built.
- **E4 — Out-of-band provisioning procedures.** The Permission Source port realization (M2 side) and the trust-material provisioning procedure (M4 side) are operational contracts deliberately behind ports. Opened by: the Record & Verifier Contract RFC and the operational runbook, respectively.
- **E5 — explicitly NOT an extension point: ≥3 domains / non-SPIFFE substrates.** The module graph is two-domain by construction (AP6, ER17); generalizing it is a V2 architecture act under a future scope extension, not a slot this RFC leaves open. Declaring this prevents the "natural generalization" drift TP5 forbids.

## Adversarial review (per DR4 — each load-bearing choice counter-held)

1. **M3 pure with injected inputs.** Wrong state: M3 performs its own fetches. Observable consequence: the egress instrumentation of AT16 records outbound traffic during core verification; SO2's egress-count-zero threshold fails; AP1 falsifies on its face.
2. **M5 behind a fixed contract, policy-free.** Wrong state: M3 depends on a concrete realization, or M5 holds the R policy. Observable consequence: a spike outcome or an S2/S3 resolution forces edits inside M3 (the requirement-fixed region) — AP12 falsifies; or two M5 realizations enforce R divergently and AT13 verdicts differ across compositions with identical inputs — R8 falsifies.
3. **Closed answer sets on every port.** Wrong state: a port answers outside its set, or expresses ignorance as `not-observed-revoked`. Observable consequence: AT22 cannot force the Inconclusive state (the state is unreachable because ignorance was laundered), and AT14 shows a partitioned RP accepting a revoked delegation — FM11 and INV12 falsify.
4. **Instance identity opaque in M1.** Wrong state: this RFC fixes instance-identity semantics. Observable consequence: the §4 trace table has an entry with no frozen forcing item — DR1 falsifies by inspection (FM5 explicitly declines the question).
5. **Expiry derived, not stored.** Wrong state: a module stores "expired" as a mutable state. Observable consequence: a stored flag can be absent, stale, or rolled back — AT11's never-returns-to-valid assertion gains a failure path that pure arithmetic on an immutable record cannot have; INV3's monotonicity falsifies exactly there.
6. **The record as the only cross-domain artifact on the verification path.** Wrong state: verification consults domain A. Observable consequence: AT15 (authorities dark) fails — verification cannot complete with the network path disabled; ER7 falsifies.
7. **Unconditional decision trace, including Accepts.** Wrong state: traces only on rejection. Observable consequence: AT23's rollback protocol cannot confirm which checks fired on the accepting side of each forced-failure pair; AT30's reviewer cannot reproduce Accept verdicts — SO5/SO8 falsify.

Each consequence is observable by a reviewer with the package and a build; no choice is dark.

## Honest-claims statement (per DR5)

This architecture carries forward, unchanged in status:

- **FM5 — issuer signing-key compromise: unmitigated within current scope.** No module resists a forged record signed with a compromised key; M1's integrity validation detects alteration of existing records only. The FM5 "re-issue after revocation" sub-case is exactly the E2 open question, carried, not solved.
- **FM8 — within-window replay: unmitigated within current scope.** A captured valid record re-presented within its window passes every M3 check by design; no module claims otherwise.
- **S4 — partition observability limit.** M5's frozen as-of during a partition is the honest signal; M3's ceiling policy never claims in-partition observability (INV12). No claim in this RFC exceeds these limits.

## Scope statement (per DR6)

This RFC asserts structure and behavior only within the two-domain, SPIFFE-coexisting, single-hop scenario: domain A (principal, M2, M6), domain B (relying party, M3, M4, M5), no shared authority, federation disabled. ≥3 domains, non-SPIFFE substrates, and multi-hop delegation are out of V1 (ER17, AP6, S5, DEFERRED D3–D4); §16 E5 explicitly declines to leave a slot for them.

## Hypothesis-preservation statement (per DR7)

Two hypothesis properties are designed-for, not committed: **fail-closed** (the Inconclusive state and its →Reject routing in M3 — NFR3/ER11/SO4/C-INV1; present, structural, testable via AT22, warranted by nothing until V1 confirms) and **latency** (NFR1/ER12 — no latency concept appears in any module contract; AT26 measures and reports only). No `[HYPOTHESIS]` item is promoted anywhere in this RFC.

## Reviewability statement (per DR8)

Every claim here is checkable by direct reading against the frozen package plus RFC-000/001/002: each module and rule traces (§4); each answer set is closed (§9); each wrong state names its observable falsifier (adversarial review). No build is required to review this RFC's *structure*; the module contracts additionally define what a build must expose for the build-complemented review of later RFCs — the decision trace (§14) is the reviewer's instrument, and AT30's reviewer needs the specification package, a build, and the traces, nothing else.

## Open questions (per RFC-000 §12)

| # | Question | Resolved by |
|---|---|---|
| Q1 | S1–S4 values (R, cached-pull admissibility, broker definition, partition reading) — this RFC is parametric over all of them | Founder scope-act journal entry (blocks EXP-001, not this RFC) |
| Q2 | Instance-identity semantics (E2; FM5's open question) | Frozen-package amendment, then reflected into M1 by the Record & Verifier Contract RFC |
| Q3 | The record's concrete content and encoding; the skew-tolerance value (ER3); the conformance checklist derived from M3 | Record & Verifier Contract RFC (next mechanism RFC; spike-independent) |
| Q4 | Permission Source port realization and its operator in the two-domain experiment | Record & Verifier Contract RFC + operational runbook |
| Q5 | M5/M6 realizations and the propagation channel (E1) | EXP-001 outcome + Revocation Mechanism RFC |
| Q6 | Technology selection for all modules | Technology Selection RFC (founder roadmap position: after RFC-003) |
| Q7 | `DEVELOPMENT_RULES.md` §RFC Policy reconciliation (overridden for this act by founder instruction; text still contradicts practice) | Founder governance act + docs amendment |
| Q8 | Ratification of RFC-000/001/002 and of this RFC through the RFC-000 state machine | Founder acts |

## Provenance (per RFC-000 §13)

- **Primary source:** the frozen Phase 7 + Phase 8 packages (hash-pinned in `FROZEN.sha256`); frozen `LEVEL0_1_FEASIBILITY_GATE.md`; RFC-000/001/002 (Drafts); founder instruction of 2026-07-05 authorizing exactly this RFC.
- **Secondary (input, not authority):** `ARCHITECTURE_READINESS_REVIEW.md` §5–§8 (candidate analysis this RFC formalizes); `lab/EXP-001-EXECUTION-PLAN.md` (substrate-sharing premise of §15).
- **Confidence:** High that every module, rule, and contract traces to a frozen forcing item with no invention (§4; adversarial review item 4 is the self-check). High that the six-module decomposition isolates the spike-volatile region per AP12 (§6 shows all four outcomes land in M5/M6 realizations only). Medium that the one-artifact reading (presented unit = reconstruction record) survives the Record & Verifier Contract RFC — no frozen item forces separation, but the concrete format work could surface one; change-condition: that RFC finding a forcing item, which would split M1's artifact into presentable unit + record while preserving the 1:1 `produces` relation. Medium that the M5 answer-set granularity (three members + as-of) is sufficient for every spike composition — change-condition: EXP-001 surfacing a composition whose honest answer cannot be expressed in the set; the contract would then be amended *before* that composition is adopted, per the freeze discipline.
- **Change-condition:** the module set, dependency rules, and port contracts change only by amendment to this RFC (once accepted, per `CONTRIBUTING.md` §4) driven by a frozen-package amendment or a founder scope act. Resolving S1–S4 sets M3's policy parameters; it changes no structure. The EXP-001 outcome selects M5/M6 realizations; it changes no contract.
- **Freeze status of RFC-003:** not frozen. Draft, presented for founder review per the RFC-000 review path (Draft → Adversarially Reviewed → Accepted → Frozen). Eligible for the frozen set (DR10) only by a future, explicitly authorized act.

<!-- checkpoint: planning(threat-model-scenarios): refine threat model scenarios -->
