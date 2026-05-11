# Module Specification — Atlas

**Status:** Canonical module specification (architecture-phase closure set). Consolidates and supersedes `PROJECT_MODULE_SPECIFICATION.md` (Sprint 1), incorporating deltas AD-013…AD-017 from `ENGINEERING_DECISION_RECORD.md`. The superseded document remains as history.
**Companions:** contracts in `INTERFACE_SPECIFICATION.md`; graphs, flows, and cross-cutting rules in `SYSTEM_ARCHITECTURE.md`; tasks in `IMPLEMENTATION_MASTER_PLAN.md`.
**Rule of this document:** each module carries exactly: responsibilities, public contracts (summary — the interface spec is binding), private responsibilities, dependencies, forbidden dependencies, internal state, ownership, extension rules, lifecycle, configuration, testing. Nothing here may widen an interface-spec contract.

Package mapping (per `AI_BOOTSTRAP.md` §4 stack ruling): M1 `internal/record`, M2 `internal/issuance`, M3 `internal/verify`, M4 `internal/truststore`, M5 `internal/revstatus` (+ `contracttest`), M6 `internal/revorigin`; composition roots in `cmd/`; harness in `tests/`.

---

## M1 — Record Model (`internal/record`) — the stable surface

- **Responsibilities:** define the Delegation Record (principal, delegate, scope, expiration, issuance time, opaque instance identity, opaque optional revocation binding [AD-015], integrity envelope); integrity validation `Intact | Altered` runnable by any holder; reconstruction reading of intact records. (ER1, ER4, INV1, INV8, INV9, SO3.)
- **Public contracts:** `ValidateIntegrity`, `Read`, `InstanceID` (opaque, equality-only), `RevBinding` (opaque, carried-not-interpreted), restricted creation surface (M2-only). Binding text: `INTERFACE_SPECIFICATION.md` §1.
- **Private responsibilities:** envelope construction/parsing (assumption A3: JWS per gate C2/C3/C7 + vetted `go-jose/v3`); canonicalization of protected content; algorithm pinning at validation (SR-2 mitigation — no attacker-selected downgrade).
- **Dependencies:** none in-system (rule R1). Vetted external: `go-jose/v3`; `go-spiffe/v2` identity-type vocabulary.
- **Forbidden dependencies:** every other `internal/` package; `net/*`; any I/O, clock, filesystem; `spire-api-sdk`.
- **Internal state:** none. Pure.
- **Ownership:** stable-surface owner (single); every change requires an `ENGINEERING_DECISION_RECORD.md` entry with frozen trace.
- **Extension rules:** exactly two sanctioned evolutions — (a) instance-identity semantics deepening via the FM5 amendment (O2): changes the minting/comparison rules, never the opaque carriage; (b) the revocation-binding element gains an interpretation *by M5 realizations only* (P2) — M1 continues to carry it opaquely. Any other change is a stable-surface amendment event.
- **Lifecycle:** stateless library; exists from E2; retired only by frozen-package amendment.
- **Configuration:** none, by decision (a configurable record model environment-dependent-izes the stable surface).
- **Testing:** pure-unit boundary — round-trip create/read; mutation corpus (bit-flips, substitutions, truncation, reordering; detection fraction = 1 → AT20); reconstruction-alone sufficiency (AT19/AT21); instance-ID and binding-element opacity (nothing interprets them); algorithm-pinning negative tests.

## M2 — Issuance Authority (`internal/issuance`)

- **Responsibilities:** sole creator of records; strict-subset scope guard via `PermissionSource` port with refusal creating nothing (SO6, FM6); bind all record content including issuance time from `TimePort` [AD-014], fresh `InstanceID` from the minter seam [AD-013], and revocation binding from `RevBindingSource` (empty pre-spike) [AD-015]; ephemeral-delegate support (ER10); `IssuanceTrace` per request, both outcomes.
- **Public contracts:** `Issue → Record | Refused(cause)`; consumed ports `PermissionSource`, `TimePort`, `RevBindingSource`. Binding text: interface spec §2.
- **Private responsibilities:** proper-subset computation; default instance-ID minter (unique per issuance — A4); trace assembly.
- **Dependencies:** `internal/record`; stdlib; `go-spiffe/v2` (consuming already-issued identity — ER15).
- **Forbidden dependencies:** `internal/verify`, `internal/truststore`, `internal/revstatus`, `internal/revorigin`; `net/*` (realizations behind ports may not be reached through this package).
- **Internal state:** none between requests.
- **Ownership:** track A; domain-A operator at runtime.
- **Extension rules:** new refusal causes only via interface-spec amendment; the minter seam accepts the future FM5-amendment semantics without shape change; no other extension.
- **Lifecycle:** per-deployment instantiation with injected ports; stateless; retirement leaves all issued records valid (self-sufficiency).
- **Configuration:** injected ports only. No default TTLs, no policy — expiration comes from the request; judging it is the verifier's job.
- **Testing:** injected-unit boundary — stub `PermissionSource` (subset/over-scope/unavailable → AT4), controllable clock, deterministic test minter; ephemeral issuance (AT18 logic); nothing-created-on-refusal; trace completeness.

## M3 — Verification Core (`internal/verify`) — the conformance definition

- **Responsibilities:** the five checks as separately named, separately observable pipeline stages — identity binding (INV1), integrity via M1 (INV8), expiry ± explicit skew (INV3/ER3), scope integrity (INV8, not subset re-derivation), revocation status under freshness policy (SO1/FM2/FM4) — with verdict routing `Accept | Reject(causes) | InconclusiveRejected(causes)` `[HYPOTHESIS]` and an **unconditional** `DecisionTrace`; sole home of policy {R, skew tolerance, S4 ceiling} (rule R8); definition of verifier conformance (FM10/AT30).
- **Public contracts:** `Verify`; `Policy` (refuses unset construction); verdict/cause/trace types; consumed ports `TrustMaterialPort`, `RevocationStatusPort`, `TimePort`. Binding text: interface spec §3.
- **Private responsibilities:** check ordering and trace assembly (private so no caller can run a check subset and call it verification); freshness and skew arithmetic.
- **Dependencies:** `internal/record`; stdlib.
- **Forbidden dependencies:** `net`, `net/http` (structural AP1); all provider packages (`truststore`, `revstatus`, `revorigin`), `issuance`; `go-jose` directly (integrity goes through M1); clock reading outside `TimePort`.
- **Internal state:** none; no verdict memory (RFC-002 §9.2).
- **Ownership:** track A; the RP operator runs instances; the public surface changes only with interface-spec amendment (it *is* the conformance definition).
- **Extension rules:** pipeline closed in V1; a sixth check enters only by founder scope act (S5 multi-hop, FR10 posture) + spec amendment, as a named stage with its own causes and trace line. Verdict/cause sets closed; the Inconclusive→Reject routing is a single named point, redirectable only by the V1 hypothesis-resolution act.
- **Lifecycle:** constructed with `Policy`; stateless per presentation; replaceable at will (conformance is the contract, not the instance).
- **Configuration:** `Policy`, injected; scope-act values at AT time, arbitrary in unit tests (AP7); recorded in every trace.
- **Testing:** injected-unit boundary, the richest surface — per-check pass/fail; single-check rollback (AT23: each check forced to fail while others pass; verdict flips every time); every inconclusive cause → InconclusiveRejected (AT22, `[HYPOTHESIS]`-marked); skew within/at/beyond tolerance (AT8 logic); freshness vs R and the S4 ceiling parametrically (AT13/AT14 logic); trace unconditionality including Accepts.

## M4 — Trust Material Store (`internal/truststore`)

- **Responsibilities:** hold out-of-band-provisioned trust material; answer `material | absent`; **never fetch** (the FM9 insecure-fallback path is deleted structurally); provisioning records per act.
- **Public contracts:** `TrustMaterialFor`, `Provision → accepted | refused(cause)`. Satisfies `TrustMaterialPort` structurally. Interface spec §4.
- **Private responsibilities:** holding structure; provisioning-record assembly; material coherence checks at `Provision`.
- **Dependencies:** `internal/record` vocabulary; stdlib; `go-spiffe/v2` bundle types.
- **Forbidden dependencies:** `net/*` (the package must be *incapable* of fetching); `internal/verify`; `internal/revstatus`.
- **Internal state:** the provisioned material and provisioning records. Single-writer (provisioning acts) / multi-reader (verifications) — §13 of `SYSTEM_ARCHITECTURE.md`.
- **Ownership:** track B; RP operator.
- **Extension rules:** provisioning *procedure* is an operational seam (P4, runbook-defined); persistence substrate deferred (TD-1) — arrives behind the same interface.
- **Lifecycle:** empty at construction (answers `absent` — correct, not an error); populated/withdrawn only by explicit provisioning acts.
- **Configuration:** none in V1 (in-memory).
- **Testing:** pure-unit — hit/miss/withdrawn; malformed-material refusal; provisioning-record completeness. Never-fetch enforced by import lint (structural, not asserted by a test).

## M5 — Revocation Status Provider (`internal/revstatus`) — the volatile region

- **Responsibilities:** the fixed contract every realization honors forever: `StatusOf → ObservablyRevoked(asOf) | NotObservedRevoked(asOf) | Indeterminate`, mandatory as-of on knowledge answers, honest-indeterminate rule (ignorance is never `NotObservedRevoked`), zero policy (freshness judgment is M3's alone); the degenerate always-Indeterminate realization (spike outcome β/δ representability); the exported `contracttest` suite gating every realization.
- **Public contracts:** answer set + `StatusOf`; `contracttest` suite. Satisfies `RevocationStatusPort` structurally. Interface spec §5.
- **Private responsibilities:** none in the degenerate realization; view maintenance in real realizations (spike-defined).
- **Dependencies:** `internal/record` (`InstanceID`, `RevBinding` interpretation once P1/P2 open); stdlib.
- **Forbidden dependencies:** `internal/verify`, `internal/truststore`, `internal/revorigin` (no direct coupling — the channel is deferred); `net/*` in everything shipped pre-spike. A real realization's *maintenance* path may receive an explicit, file-scoped network-import grant with the mechanism decision — never reachable from `StatusOf`'s answer path (rule R6).
- **Internal state:** the RP-local view + its as-of. Single-writer (maintenance) / multi-reader (verifications). As-of freezes during partition — a defined state, not a fault.
- **Ownership:** track B for contract + degenerate + suite; the real realization's owner is assigned with the mechanism decision (O3).
- **Extension rules:** this is the system's one plugin boundary (P1). Admission rule: pass `contracttest`, honor the answer set unchanged, take mechanism-specific per-record data only from the opaque `RevBinding` element (P2), never widen the contract. A composition that cannot satisfy this is surfaced to the founder as a spec-amendment question (risk ER-3), never wedged in.
- **Lifecycle:** degenerate — construct, answer `Indeterminate`, forever. Real — maintenance loop per composition; read at verification.
- **Configuration:** none in the contract; realization-specific configuration arrives with the mechanism decision, confined to the realization.
- **Testing:** `contracttest` (closed set, mandatory as-of, honest-indeterminate, determinism per view state) — run against the degenerate realization now and every spike candidate later, at the same seam production uses.

## M6 — Revocation Origin (`internal/revorigin`)

- **Responsibilities:** the authoritative append-only revocation register: `Revoke(instanceID) → recorded` (idempotent on terminal state, INV4); ordered read-only `View` for the deferred propagation channel and reconstruction; structural incapacity to affect identities (INV5) or siblings (INV6) — it stores only instance IDs.
- **Public contracts:** `Revoke`, `View`. Interface spec §6.
- **Private responsibilities:** append-only enforcement; ordering.
- **Dependencies:** `internal/record` (`InstanceID`); stdlib.
- **Forbidden dependencies:** everything else `internal/`; `net/*` (publication is the channel's business).
- **Internal state:** the register. Single-writer / multi-reader; append serialized.
- **Ownership:** track B; domain-A operator.
- **Extension rules:** the propagation channel reads `View`; its mechanism (push/pull/cached-pull per S2/S3) attaches *outside* this module with the spike outcome. The register's shape does not change for any composition.
- **Lifecycle:** append-only forever; outlives every delegation it names; no compaction, no deletion.
- **Configuration:** none in V1 (in-memory; persistence deferred, TD-1).
- **Testing:** pure-unit — append-only property; re-revoke no-op on terminal state (AT11 origin side); ordering; view stability; revocation of never-issued IDs is inert.

---

## Non-module code (specified for completeness; not modules)

- **`cmd/atlas-issue`, `cmd/atlas-verify`, `cmd/atlas-revoke`** — composition roots: wiring, policy loading, trace persistence (AD-016), operational logging, the AT26 measurement point (around `Verify`, never inside). Zero delegation logic — a check in a driver is a conformance violation. One request per invocation in V1 (§13, `SYSTEM_ARCHITECTURE.md`).
- **`tests/harness`** — port fakes, controllable clock, substrate-control interfaces (two-domain SPIRE control, link-level partition, egress observation) realized in epic E6; instrumentation, unconstrained imports.
- **`tests/acceptance`** — one file per AT family; product exercised only through drivers/public surfaces (mirrors each AT's test locus); substrate-blocked tests skip with named blockers, never TODOs.
- **`scripts/check-imports.sh`** — rules R1–R7 as a build gate. Changing an allowed/forbidden table requires amending this document with a trace, then the lint config — in that order.

<!-- checkpoint: fix(revstatus): fix revstatus snapshot retrieval -->
