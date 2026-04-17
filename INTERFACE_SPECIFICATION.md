# Interface Specification — Atlas

**Status:** Canonical, binding contract set (architecture-phase closure). Language-free: contracts only — inputs, outputs, ownership, invariants, failure semantics. No syntax, no signatures. Consolidates and supersedes `MODULE_INTERFACE_SPECIFICATION.md` (Sprint 1), incorporating AD-013…AD-017. An implementation realizes exactly these surfaces; nothing more is public.
**Universal rules (bind every interface below):**
1. **Closed answer sets.** Every output is a member of an enumerated set; "unknown/cannot-answer" is always an explicit member where it can occur; nothing may answer outside its set (FM11).
2. **Answers are values.** No exception, timeout-as-implicit-answer, or out-of-band fault channel substitutes for a set member. An abnormal termination crossing an interface is a defect, not an error path.
3. **Immutable inputs.** No operation mutates what it is given.
4. **Provenance.** Every negative or inconclusive answer names its cause; every cause is attributable to one check, guard, or port.
5. **Set evolution.** Adding a member to any set is a specification amendment (this document first, with a frozen trace via `ENGINEERING_DECISION_RECORD.md`, then code). Existing members are never repurposed.
6. **Hypothesis marking.** Everything touching the Inconclusive→Reject routing carries `[HYPOTHESIS]` (NFR3/ER11/SO4/C-INV1) until a V1 confirmation act.

---

## §1 The Delegation Record (M1) — artifact and operations

### The artifact

One logical unit serving both presentation (ER1) and reconstruction (ER4) — AD-002. Content, at the conceptual level:

| Element | Meaning | Interpretation rights |
|---|---|---|
| Principal identity | on whose behalf | everyone |
| Delegate identity | who acts | everyone |
| Scope | permission subset granted | everyone (inspectable, FR2) |
| Expiration | end of validity window | everyone |
| Issuance time | when created | everyone |
| Instance identity | opaque, unique per issuance (AD-013; semantics pending the FM5 amendment, O2) | **equality comparison only** — no party parses, orders, or derives meaning |
| Revocation binding | opaque, optional, empty pre-spike (AD-015) | **M5 realizations only**; M1 carries, M3 ignores, M2 obtains via port |
| Integrity envelope | binds all of the above | verifiers via ValidateIntegrity |

*Ownership:* the record is owned by no runtime after creation — it is self-sufficient (INV9); any holder holds a full-fidelity instance. *Invariants:* immutable after creation; any post-creation change is detectable (INV8); both identities deterministically recoverable with trust material (INV1).

### ValidateIntegrity

- **Input:** a claimed record (any bytes); trust material.
- **Output:** `Intact | Altered`.
- **Invariants:** deterministic; `Intact` iff exactly what the issuer created, verifiable with the supplied material; unparsable/truncated/reordered/substituted input → `Altered`. Expired or revoked records still validate `Intact` — validity is §3's question.
- **Failure semantics:** there is no failure distinct from `Altered`; the operation cannot be inconclusive (it holds everything it needs by construction). Verification-key absence is the *caller's* prior condition (§3 check 2 handles it as inconclusive before this operation runs).
- **Ownership:** callable by anyone — verifier, third party, reviewer.

### Read

- **Input:** a record for which `ValidateIntegrity` answered `Intact` with appropriate material.
- **Output:** the assertions (all table elements above).
- **Invariants:** returns exactly what issuance established (INV1 — no context-dependent reinterpretation). Defined only on intact records; reading unvalidated bytes is a caller defect, not a supported path.
- **Reconstruction** (ER4/INV9/SO8) = ValidateIntegrity + Read exercised by a third party holding only the record and trust material — no verifier runtime state, no privileged input. This composition is the system's audit answer; no separate interface exists.

### Creation (restricted)

- Only the Issuance Authority constructs records (AP10). The creation surface requires every element and refuses partial construction; test fixtures create records through §2 with stub ports, never through a bypass.

---

## §2 Issuance (M2)

### Issue

- **Input:** a request — principal identity, delegate identity, requested scope, requested expiration.
- **Output:** `Record | Refused(cause)`; cause set: `OverScope | PermissionsUnavailable | MalformedRequest`.
- **Invariants (issued):** requested scope is a **proper** subset of the principal's permission set, affirmatively established via PermissionsOf before construction (ER2 "strict subset" — equality refuses); the record carries a fresh instance identity, the issuance-time reading from the time port (AD-014), and the revocation-binding element from its port (empty pre-spike); the record is complete and self-sufficient at creation; nothing about the delegate's lifetime is assumed (ER10 — the delegate is named, not contacted).
- **Invariants (refused):** nothing exists afterward that did not exist before — no partial record, no reserved instance identity, no side effect beyond the trace (FM6). `PermissionsUnavailable` from the port maps to refusal always: no cached prior answer, no retry within the module.
- **Failure semantics:** the refusal set is total — every non-issued outcome is one of the three causes.
- **Trace obligation:** one issuance trace per request, both outcomes (§7).
- **Ownership:** domain-A operator invokes; the record returned belongs to the requester.

### Consumed ports (defined by this module; realized elsewhere; structural satisfaction)

- **PermissionsOf:** principal identity → `PermissionSet | Unavailable`. Honest-unavailable; never blocks indefinitely (the composition root owns timeouts and maps them to `Unavailable`).
- **Time reading** (shared contract with §3): → a time reading. No further semantics here; skew discipline is the verifier's.
- **RevocationBindingFor:** (principal, delegate, scope, expiration, instance identity) → `binding | none`. Pre-spike realization answers `none` always. The binding is opaque to M2; it is minted by the mechanism side (P2) and consumed only by M5 realizations.

---

## §3 Verification (M3) — the conformance definition

### Verify

- **Input:** a presented record; three port answers (gathered by the composition root or taken through injected ports): trust-material answer, revocation-status answer, time reading; a policy (below).
- **Output:** a verdict — `Accept | Reject(causes) | InconclusiveRejected(causes)` `[HYPOTHESIS]` — plus a decision trace (§7), unconditionally.
- **The five checks** (each a named stage; all five always evaluated and traced; verdict routing below):

| # | Check | Pass condition | Definitive-failure cause | Inconclusive cause |
|---|---|---|---|---|
| 1 | Identity binding (INV1) | both identities recoverable and consistent from record + material | `BindingMismatch` | — |
| 2 | Integrity (INV8, via §1 ValidateIntegrity) | material present and record `Intact` | `IntegrityFailed` | `TrustMaterialAbsent`, `SignatureUnverifiable` |
| 3 | Expiry (INV3, ER3) | deterministically un-expired at the time reading within the stated skew tolerance | `Expired` | `ClockBeyondTolerance` |
| 4 | Scope integrity (INV8) | scope covered by the intact envelope — integrity, **not** subset re-derivation (the RP need not hold the principal's permission set; subset-ness was issuance's guard) | `ScopeIntegrityFailed` | — |
| 5 | Revocation (SO1, FM2/FM4, INV12) | status answer is `NotObservedRevoked(asOf)` with asOf within the policy's R-derived freshness bound and S4 ceiling | `RevokedObservable` | `RevocationStatusIndeterminate`, `RevocationKnowledgeStale` |

- **Verdict routing (deterministic, order-independent):** any definitive cause → `Reject(causes)`; no definitive cause but ≥1 inconclusive cause → `InconclusiveRejected(causes)` `[HYPOTHESIS]` — a *distinct* verdict that **is** a rejection, kept distinguishable so AT22 can separate designed fail-closed behavior from definitive rejection; all pass → `Accept`.
- **Invariants:** no I/O beyond the three port consultations; no retry of any port; no fallback between checks; no verdict memory (a re-presentation is a fresh verification); the record's revocation-binding element is ignored (it is M5-realization food, never verifier input).
- **Failure semantics:** the verdict *is* the failure semantics — no other channel exists.
- **Conformance:** a **conformant verifier** (FM10, AT30) is exactly: these five checks, these sets, this routing, an unconditional trace. Fewer, more, or reordered-with-different-outcomes is non-conformant.
- **Ownership:** the RP invokes and owns the verdict + trace instance; the grant decision remains the RP's, outside the system (RFC-001).

### Policy

- **Members:** R (revocation-observability bound), skew tolerance (issuer-relative, ER3), S4 ceiling (partition-recovery reading).
- **Invariants:** construction with any member unset **refuses** — an unparameterized verifier must not exist (FM2/FM4 have no defaults; AP7). Values come from founder scope acts at acceptance time, arbitrary values in unit tests. The policy in force is identified in every trace.
- **Ownership:** held by the verifying boundary's operator; changed only by a new scope act.

### Consumed ports

- **TrustMaterialFor:** trust domain → `material | absent` (§4 realizes).
- **StatusOf:** instance identity → the §5 answer set (§5 realizes).
- **Time reading:** shared contract with §2; the harness's controllable clock realizes it in tests (AT8).

---

## §4 Trust Material Store (M4)

- **TrustMaterialFor:** trust domain → `material | absent`. *Invariants:* answers only what provisioning established; `absent` before provisioning, after withdrawal, for unknown domains; **the interface has no fetch semantics and its realization must be incapable of acquiring them** (FM9 — the insecure-fallback path is deleted, not discouraged). *Failure semantics:* `absent` is an answer, not an error.
- **Provision:** material → `accepted | refused(Malformed | DomainMismatch)`. *Invariants:* an out-of-band operator act, never reachable from a verification; refused material is never stored, never half-trusted; each act appends a provisioning record (what, when, which act). *Ownership:* the RP operator owns the store's content and its provisioning records.

## §5 Revocation Status (M5) — the volatile region's permanent contract

- **StatusOf:** instance identity → exactly one of:
  - `ObservablyRevoked(asOf)` — this RP's view contains a revocation of this instance, current as of the stated time;
  - `NotObservedRevoked(asOf)` — this RP's view, current as of the stated time, contains no such revocation;
  - `Indeterminate` — cannot currently answer (no view, corrupted view, no realization). Carries no asOf — there is no knowledge to date.
- **Invariants:** the as-of is mandatory on knowledge answers; **honest-indeterminate** — ignorance is never expressed as `NotObservedRevoked` (AP5/INV12); the provider applies **no policy** — freshness judgment is §3's alone (R8); answers are a deterministic function of (view state, instance identity); during a partition the as-of stops advancing — the honest, information-theoretically correct signal (S4); no answer path may perform egress (R6; a realization's *maintenance* path is separately governed by the mechanism decision).
- **Failure semantics:** `Indeterminate` *is* the failure boundary — every internal problem surfaces as it; nothing else escapes.
- **Admission rule (the plugin boundary, P1):** no realization — degenerate, spike candidate, or production — is wired into a composition root unless it passes the realization-independent contract suite asserting all invariants above. Mechanism-specific per-record data comes only from the record's opaque revocation-binding element (P2).
- **Ownership:** the RP operator owns the view; the *authority* over revocation facts remains M6's — this interface answers about observation, never about truth.

## §6 Revocation Origin (M6)

- **Revoke:** instance identity → `recorded`. *Invariants:* append-only; one-way and terminal per instance (INV4) — re-revocation is a no-op, not an error; the register holds only instance identities and is structurally incapable of touching underlying identities (INV5) or sibling delegations (INV6); revocation of a never-issued identity is inert (verification keys to presented records). *Failure semantics:* a storage-level inability to append surfaces as a named refusal; there is no silent drop.
- **View:** → the ordered, complete register, read-only. *Invariants:* no operation removes or rewrites entries; the register outlives every delegation it names. *Ownership:* domain-A operator owns the register; the propagation channel (deferred, S2/S3, spike-selected) and after-the-fact reviewers read it.

## §7 Traces (cross-cutting observables)

- **DecisionTrace** (produced by §3, per verification, **including Accepts**): ordered entries for all five checks — check name, inputs digest, outcome; the injected time reading and as-of values; the policy in force; the final verdict + causes. *Invariants:* unconditional (an Accept without a trace is non-conformant — AT23 and AT30 need Accepts inspectable); contains nothing a reviewer may not hold (AP8); fields append-only across versions, never repurposed.
- **IssuanceTrace** (produced by §2, per request, both outcomes): request summary, permission-consultation outcome, subset determination, result with instance identity or refusal cause, issuance time.
- **Ownership (AD-016):** traces are **returned values**; persistence, shipping, retention belong to the invoking boundary's operator (composition root / harness). Modules carry no logging dependency; observation never alters a verdict (AP13 passivity).

## §8 Interface-level invariant summary (the reviewer's checklist)

| Invariant | Where it lives |
|---|---|
| INV1 identity binding | §1 Read determinism; §3 check 1 |
| INV2 scope subset | §2 issuance guard (proper subset; refusal total) |
| INV3 expiry monotone | §3 check 3 (derived, nothing stored) |
| INV4 revocation terminal | §6 Revoke idempotent-terminal; §5 Observable answers |
| INV5/INV6 identity/sibling independence | §6 structural (register holds only instance IDs) |
| INV7 offline verification | §3 no-I/O invariant + R6 |
| INV8 tamper-evidence | §1 ValidateIntegrity; §3 checks 2 and 4 |
| INV9 reconstruction self-sufficiency | §1 Reconstruction composition |
| INV10/INV11 companion boundary | no interface here issues identity or touches the standard |
| INV12 observability bound | §5 as-of semantics + §3 check 5 ceiling |
| C-INV1 fail-closed `[HYP]` | §3 InconclusiveRejected routing — designed-for, unpromoted |

<!-- checkpoint: repo(verification-criteria): restructure verification criteria (#9) -->

<!-- checkpoint: planning(API-path-design): document API path design -->

<!-- checkpoint: refactor(issuance): refactor attenuation rule engine -->
