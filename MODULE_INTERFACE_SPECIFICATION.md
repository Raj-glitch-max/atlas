# Module Interface Specification

**Status:** Engineering document. The binding contract for every public surface named in `PROJECT_MODULE_SPECIFICATION.md`. Not source code: contracts are stated as operations with named inputs, outputs from closed sets, preconditions, postconditions, and trace obligations. Implementation (Sprint 1 W2/W3) must realize exactly these surfaces — no additional public operations, no widened answer sets.
**Authority:** RFC-003 §9 (conceptual interfaces, refined here), frozen items cited per operation.
**Conventions:** `→` separates inputs from the answer set. Every answer set is **closed**: an implementation may never return anything outside it, and "unknown" is always an explicit member where it can occur (FM11). Logical shapes name fields; wire encoding is out of scope (assumption A3 governs the envelope; field-level encoding is Sprint 2).

---

## 1. `internal/record` (M1)

### Types

- **`Record`** — logical content: `principal identity`, `delegate identity`, `scope`, `expiration`, `issuance time`, `instance identity (opaque)`, `integrity envelope` binding all of the above. Immutable after creation.
- **`InstanceID`** — opaque. Supports equality comparison and nothing else. No operation in any package parses, orders, or derives meaning from it (E2 deferral).
- **`Assertions`** — the read result: `principal`, `delegate`, `scope`, `expiration`, `issuance time`, `instance identity`.
- **`TrustMaterial`** — the verification-key material for a trust domain, as held locally by an RP (vocabulary type; population is M4's business).

### Operations

**`ValidateIntegrity(record, trust material) → Intact | Altered`**
- *Pre:* none — any byte sequence claiming to be a record is a legal input.
- *Post:* `Intact` iff the record is exactly what its issuer created, verifiable with the supplied material (INV8). Any post-creation alteration, truncation, reordering of protected content, or unparsable input → `Altered`. Deterministic: same inputs, same answer.
- *Prohibitions:* no I/O; no clock; expired or revoked records still validate `Intact` (validity is M3's question, not M1's).
- *Trace obligation:* none (pure function; the caller traces).
- *AT anchors:* AT5, AT19–AT21 (with §Reconstruction), mutation corpus AT20.

**`Read(record) → Assertions`**
- *Pre:* `ValidateIntegrity` returned `Intact` for this record with appropriate material. Reading an unvalidated record is a caller defect; the operation is defined only on intact records.
- *Post:* the assertions established at issuance, exactly (INV1 determinism — no context-dependent interpretation).
- *AT anchors:* AT1, AT3 (scope inspectability), AT19.

**Reconstruction** (capability, not an extra operation): a third party holding only a record and trust material runs `ValidateIntegrity` + `Read` and obtains delegator, delegate, scope, and time — no verifier runtime state, no privileged input (INV9, SO8). AT19/AT21 exercise precisely this composition.

### Construction

Record construction is **not public**. Only `internal/issuance` creates records (AP10); the constructor is exposed to it through a narrow creation surface that requires every field and refuses partial construction. Test fixtures create records through `issuance` with stub ports, never through a bypass.

---

## 2. `internal/issuance` (M2)

### Types

- **`Request`** — `principal identity`, `delegate identity`, `requested scope`, `requested expiration`.
- **`RefusalCause`** — closed set: `OverScope` | `PermissionsUnavailable` | `MalformedRequest`.
- **`IssuanceTrace`** — per request: request summary, permission-set consultation outcome, subset determination, result (`issued(instanceID)` | `refused(cause)`), issuance time.

### Ports (consumed; defined here)

**`PermissionSource.PermissionsOf(principal) → PermissionSet | Unavailable`**
- Realization deferred (RFC-003 E4). `Unavailable` is an honest answer; the port never blocks indefinitely (the driver owns timeouts and maps them to `Unavailable`).

### Operations

**`Issue(request) → Record | Refused(cause)`**
- *Pre:* none beyond a well-formed request (`MalformedRequest` otherwise).
- *Post (issued):* the record binds exactly the requested principal, delegate, scope, and expiration plus issuance time and a fresh `InstanceID` (assumption A4: unique per issuance); the record is complete and self-sufficient at creation; requested scope ⊂ principal's permission set was affirmatively established (SO6). Strictness: scope equal to the full permission set is refused — the relation is proper subset (ER2 "strict subset").
- *Post (refused):* nothing was created; no partial record, no reserved instance ID, no side effect beyond the trace (FM6).
- *Determinism of refusal:* `PermissionsOf → Unavailable` ⇒ `Refused(PermissionsUnavailable)`, always — no cached prior answer, no retry ladder inside the module.
- *Ephemeral support:* nothing in the operation requires the delegate to pre-exist, persist, or hold a long-lived identity (ER10); the delegate is named, not contacted.
- *Trace obligation:* one `IssuanceTrace` per call, both outcomes (RFC-003 §14).
- *AT anchors:* AT4, AT18.

---

## 3. `internal/verify` (M3)

### Types

- **`Policy`** — `R` (revocation-observability bound), `skew tolerance` (issuer-relative, ER3), `S4 ceiling` (partition-recovery reading parameters). Construction **refuses** unset/zero members: an unparameterized verifier must not exist (FM2/FM4 cannot default; AP7 — values come from founder scope acts at AT time, arbitrary in unit tests).
- **`Verdict`** — closed set: `Accept` | `Reject(causes)` | `InconclusiveRejected(causes)`. The third member *is* the fail-closed routing: an inconclusive determination surfaces as a distinct verdict that **is** a rejection, so the hypothesis behavior (NFR3/ER11/SO4 `[HYPOTHESIS]`) is observable and testable without being conflated with a definitive `Reject` (AT22 needs to tell them apart).
- **`Cause`** — closed set, attributed per check: `BindingMismatch` | `IntegrityFailed` | `Expired` | `ScopeIntegrityFailed` | `RevokedObservable` (definitive) and `TrustMaterialAbsent` | `SignatureUnverifiable` | `ClockBeyondTolerance` | `RevocationStatusIndeterminate` | `RevocationKnowledgeStale` (inconclusive).
- **`DecisionTrace`** — per verification: ordered per-check entries (check name, inputs digest, outcome), the injected time reading and as-of freshness values, policy identification (which parameter values were in force), final verdict + causes. Emitted **unconditionally**, Accepts included (SO5/SO8; RFC-003 §14).

### Ports (consumed; defined here; satisfied structurally — providers never import this package)

- **`TrustMaterialPort.TrustMaterialFor(domain) → material | absent`**
- **`RevocationStatusPort.StatusOf(instanceID) → ObservablyRevoked(asOf) | NotObservedRevoked(asOf) | Indeterminate`**
- **`TimePort.Now() → time reading`**

### Operation

**`Verify(presented record, ports, policy) → Verdict + DecisionTrace`**

Pipeline, in order, every check always evaluated and traced (no short-circuit that hides later-check state from the trace — AT23 needs per-check attribution; the *verdict* may be determined by the first failure, the *trace* records all checks' outcomes):

1. **Identity binding** (INV1): both identities deterministically recoverable and consistent → else `BindingMismatch`.
2. **Integrity** (INV8): `record.ValidateIntegrity` with material from `TrustMaterialFor`. `absent` → inconclusive `TrustMaterialAbsent` (FM9 — never fetch); `Altered` → `IntegrityFailed`; unverifiable-with-held-material → inconclusive `SignatureUnverifiable`.
3. **Expiry** (INV3, ER3): expiration vs `Now()` under `skew tolerance`. Deterministically not-expired → pass; deterministically expired → `Expired`; reading not conclusive within tolerance → inconclusive `ClockBeyondTolerance` (FM3).
4. **Scope integrity** (INV8): the scope field is covered by the intact envelope — integrity check, **not** subset re-derivation (RFC-002 §9.2: the RP need not hold the principal's permission set). Tampered → `ScopeIntegrityFailed`.
5. **Revocation** (SO1, FM2/FM4, INV12): `StatusOf(instanceID)` judged under policy: `ObservablyRevoked(asOf)` → `RevokedObservable`; `NotObservedRevoked(asOf)` with `asOf` within the R-derived freshness bound (and within the S4 ceiling under partition) → pass; `asOf` outside the bound → inconclusive `RevocationKnowledgeStale`; `Indeterminate` → inconclusive `RevocationStatusIndeterminate`.

Verdict routing: all pass → `Accept`; any definitive cause → `Reject`; no definitive cause but ≥1 inconclusive cause → `InconclusiveRejected` (`[HYPOTHESIS]` — designed-for, tested by AT22, never documented as warranted).

- *Prohibitions:* no I/O besides the three port calls; no retry of any port; no fallback between checks; no verdict memory across calls (RFC-002 §9.2); no mutation of any input.
- *Conformance:* a **conformant verifier** (FM10, AT30) is exactly an implementation of this operation: these five checks, these answer sets, this routing, an unconditional trace. Anything less or more is non-conformant.
- *AT anchors:* AT2, AT3, AT6–AT9, AT11, AT13/AT14 (logic), AT15/AT16 (structural: no egress exists to observe), AT22, AT23, AT26 (measured *around* this operation, at the driver).

---

## 4. `internal/truststore` (M4)

**`TrustMaterialFor(domain) → material | absent`**
- *Post:* returns exactly what provisioning established for that domain; `absent` before provisioning, after withdrawal, or for unknown domains. **Never fetches** — the package has no capability to (import-lint enforced).
- Satisfies `verify.TrustMaterialPort` structurally.

**`Provision(material) → accepted | refused(cause)`**
- *Pre:* an out-of-band operator act (never called from a verification path).
- *Post (accepted):* material held and answerable; a provisioning record (what, when, by which act) is appended.
- *Post (refused):* malformed or incoherent material is never stored, never half-trusted; closed cause set: `Malformed` | `DomainMismatch`.

---

## 5. `internal/revstatus` (M5)

### Answer type (the volatile region's fixed contract — every realization, forever)

**`StatusOf(instanceID) → ObservablyRevoked(asOf) | NotObservedRevoked(asOf) | Indeterminate`**
- `ObservablyRevoked(asOf)`: this RP's view contains a revocation of this instance; the view is current as of `asOf`.
- `NotObservedRevoked(asOf)`: this RP's view, current as of `asOf`, contains no revocation of this instance. **Honest-indeterminate rule:** a realization with no view, a corrupted view, or no realization at all must answer `Indeterminate` — ignorance is never expressed as `NotObservedRevoked` (AP5, INV12).
- `Indeterminate`: cannot currently answer. Carries no `asOf` (there is no knowledge to date).
- *Mandatory as-of:* both knowledge answers carry it; during a partition `asOf` stops advancing — the honest, information-theoretically correct signal (S4). The provider **never** applies R; freshness judgment is `verify.Policy`'s alone (RFC-003 R8).
- *Determinism:* answers are a function of (view state, instanceID); two calls against an unchanged view agree.

### `contracttest` suite (exported)

A realization-independent test suite asserting: closed answer set; mandatory as-of on knowledge answers; honest-indeterminate; determinism per view state. **Gate:** no realization — including the Sprint 1 degenerate one and every EXP-001 candidate — is wired into a driver unless it passes this suite.

### Degenerate realization

Answers `Indeterminate` to every query, always. It is the honest realization of spike outcomes β/δ (RFC-003 §6) and the Sprint 1 default wiring: with it, check 5 always yields `RevocationStatusIndeterminate` → `InconclusiveRejected` — the system fails closed rather than pretending revocation knowledge it does not have.

---

## 6. `internal/revorigin` (M6)

**`Revoke(instanceID) → recorded`**
- *Post:* the register contains a revocation entry for this instance; entries are append-only and ordered; revoking an already-revoked instance changes nothing (terminal state, INV4 — idempotent, not an error).
- *Scope of effect:* the register holds only instance IDs — it is structurally incapable of touching identities (INV5) or sibling delegations (INV6). A revocation of a never-issued ID is inert (M3 keys checks to presented records).
- *Failure:* a storage-level inability to append surfaces as a named refusal to the caller; there is no silent drop.

**`View() → ordered revocation entries`**
- *Post:* the complete register, in order, read-only. This is the surface the (deferred) propagation channel and after-the-fact reconstruction read. No operation removes or rewrites entries.

---

## 7. Cross-cutting contract rules

1. **Closed sets are closed.** Adding a member to any answer, cause, or refusal set is a contract amendment: this document changes first (with an RFC-003 trace), then code. CI's import/API surface checks make drive-by widening visible in review.
2. **No public surface beyond this document.** Helper functions, convenience wrappers, and "extra" constructors are private or absent (AP11 at the API level).
3. **Panics never cross module boundaries.** A panic escaping a public operation is a defect class of its own (FM11 — an undefined failure path); every operation's failure behavior is a named member of its answer set.
4. **Hypothesis marking.** Every artifact touching the `InconclusiveRejected` routing — code documentation, trace field names, test names — carries the `[HYPOTHESIS]` label exactly as the AT plan does (DR7). The label is removed only by a V1 confirmation act amending the frozen package.
5. **Trace stability.** `DecisionTrace` and `IssuanceTrace` shapes are reviewer-facing (SO8/AT30): fields may be added with a spec amendment; existing fields are never repurposed.
