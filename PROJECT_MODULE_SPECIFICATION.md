# Project Module Specification

**Status:** Engineering document. Realizes RFC-003 (Accepted) as concrete package boundaries. Not frozen.
**Authority:** RFC-003 §7–§16; frozen ER/SO/INV/FM/AT items cited per module; `AI_BOOTSTRAP.md` §4 (vetted tooling: Go 1.21+, `go-spiffe/v2`, `spire-api-sdk`, `go-jose/v3`).
**Reading rule:** each module carries the thirteen specification aspects requested by the founder. The public interfaces named here are summarized; their full contracts (operations, closed sets, pre/postconditions) are in `MODULE_INTERFACE_SPECIFICATION.md`. RFC-003's dependency rules R1–R8 bind everything below; violations are CI failures (forbidden-import lint), not review comments.

Package naming: RFC-003 module → Go package under `internal/`:

| RFC-003 | Package | Role |
|---|---|---|
| M1 Record Model | `internal/record` | stable surface |
| M2 Issuance Authority | `internal/issuance` | domain A |
| M3 Verification Core | `internal/verify` | domain B; conformance definition |
| M4 Trust Material Store | `internal/truststore` | domain B |
| M5 Revocation Status Provider | `internal/revstatus` | domain B; volatile region |
| M6 Revocation Origin | `internal/revorigin` | domain A |

`cmd/` holds the three boundary drivers (composition roots); `tests/` holds harness + acceptance scaffold. Neither is a module; both are specified at the end.

---

## M1 — `internal/record`

- **Responsibilities:** define the Delegation Record type (principal identity, delegate identity, scope, expiration, issuance time, opaque instance identity, integrity envelope); integrity validation (`intact`/`altered`) checkable by any holder; reconstruction reading of intact records. (ER1, ER4, INV1, INV8, INV9, SO3.)
- **Ownership:** the stable surface — one owner; changes require an RFC-003 amendment trace. Sprint 1 owner: track A engineer.
- **Public interface:** `Record` type; `ValidateIntegrity(record, trust material) → Intact | Altered`; `Read(intact record) → assertions`; the opaque `InstanceID` type. Full contract: `MODULE_INTERFACE_SPECIFICATION.md` §1.
- **Private responsibilities:** the JWS envelope construction/parsing (assumption A3); canonicalization of signed content; internal claim mapping. None of these leak into any signature.
- **Dependencies:** none on other atlas packages (RFC-003 R1).
- **Allowed imports:** Go stdlib (no `net/*`); `go-jose/v3` (envelope); `go-spiffe/v2` identity-type vocabulary only (SPIFFE ID parsing/representation).
- **Forbidden imports:** every other `internal/` package; `net`, `net/http`, `os/exec`; `spire-api-sdk`; anything performing I/O. M1 is pure: no clock, no filesystem, no network.
- **Lifecycle:** stateless library; no init, no teardown, no background activity.
- **Configuration:** none. A configurable record model would make the stable surface environment-dependent.
- **Test strategy:** unit + property tests, no fakes needed (pure). Mutation-set tests (bit-flips, field substitution, truncation, reorder) drive `ValidateIntegrity` → AT5, AT19–AT21 in-process; round-trip create/read; instance-ID opacity test (nothing interprets it).
- **Failure boundaries:** returns only members of its closed answer sets; malformed input is `Altered`/a named parse refusal — never a panic across the boundary, never a partial read.
- **Logging boundaries:** none. M1 emits nothing.
- **Metrics boundaries:** none.

## M2 — `internal/issuance`

- **Responsibilities:** sole creator of records; enforce requested scope ⊂ principal's Permission Set via the Permission Source port; refuse over-scope and permissions-unavailable creating nothing; bind identities/scope/expiry/issuance-time/instance-ID via M1; support ephemeral delegates; emit issuance trace. (ER1–ER3, ER10, INV1–INV3, SO6, FM6, AP10.)
- **Ownership:** track A engineer (after M1); domain-A operational owner at runtime.
- **Public interface:** `Issue(request) → Record | Refused(cause)`; `PermissionSource` port interface (consumer-defined here); `IssuanceTrace` type. Full contract: §2.
- **Private responsibilities:** subset computation; instance-ID minting (assumption A4 — the production rule is private so the FM5 amendment lands in one place); trace assembly.
- **Dependencies:** `internal/record` only.
- **Allowed imports:** stdlib; `internal/record`; `go-spiffe/v2` (consuming already-issued identity material — ER15).
- **Forbidden imports:** `internal/verify`, `internal/truststore`, `internal/revstatus`, `internal/revorigin`; `net/*` in Sprint 1 (the Permission Source realization behind the port may later justify more; the port keeps it out of this package).
- **Lifecycle:** instantiated per domain-A deployment with an injected `PermissionSource`; stateless between requests; retirement leaves issued records valid (self-sufficiency).
- **Configuration:** none in Sprint 1 beyond the injected port. No default TTLs, no policy — expiration comes from the request; validation of it is the verifier's job.
- **Test strategy:** unit tests with a stub `PermissionSource`: over-scope refusal (AT4), permissions-unavailable refusal, ephemeral issuance (AT18 logic), trace completeness. No substrate needed.
- **Failure boundaries:** `Refused(cause)` closed set; refusal creates nothing (no partial record ever exists); port `unavailable` → refusal, never a retry loop or fallback.
- **Logging boundaries:** the issuance trace is the *product-mandated* observable (RFC-003 §14) — a returned/emitted artifact, not a log line. Operational logging only in `cmd/`, never here.
- **Metrics boundaries:** none in-module; issuance counters belong to the driver.

## M3 — `internal/verify`

- **Responsibilities:** the conformant verifier: five named checks as separate pipeline stages (binding, integrity, expiry±skew, scope integrity, revocation status); freshness policy (R bound, S4 ceiling) applied in exactly one place; verdicts `Accept | Reject(cause) | Inconclusive(cause)→Reject [HYPOTHESIS]`; unconditional decision trace. Defines the three ports (trust material, revocation status, time). (ER7, ER3, ER11[HYP], INV1, INV3, INV7, INV8, SO1, SO2, SO4[HYP], SO5, FM2–FM4, FM9, FM11, AP1, AP4, AP13.)
- **Ownership:** track A engineer; this package **is** the conformance definition (RFC-003 M3) — its public surface changes only with an interface-spec revision.
- **Public interface:** `Verify(presented record, injections, policy) → Verdict + DecisionTrace`; port interfaces `TrustMaterialPort`, `RevocationStatusPort`, `TimePort`; `Policy` type (R, skew tolerance, S4 ceiling); verdict/cause/trace types. Full contract: §3.
- **Private responsibilities:** check ordering; per-check trace-entry assembly; freshness arithmetic; skew arithmetic. Private by design so no caller can invoke a subset of checks and call it verification.
- **Dependencies:** `internal/record` only. Ports are satisfied structurally — providers never import this package, this package never imports providers (RFC-003 R3).
- **Allowed imports:** stdlib (no `net/*`, no `time` for *reading* the clock — `time` types are permitted, clock *reading* comes only through `TimePort`); `internal/record`.
- **Forbidden imports:** `net`, `net/http` (structural AP1 — the zero-egress property is made unviolatable by import policy + CI lint); `internal/truststore`, `internal/revstatus`, `internal/revorigin`, `internal/issuance`; `go-jose` directly (integrity goes through M1); `spire-api-sdk`.
- **Lifecycle:** constructed with `Policy`; stateless per presentation; no verdict memory (RFC-002 §9.2); replaceable at will — conformance is the contract.
- **Configuration:** `Policy` — injected, founder-scope-act-sourced values at AT time, arbitrary values in unit tests (AP7 parametricity). A zero/unset policy refuses to construct: an unparameterized verifier must not exist (FM2/FM4's bound cannot default).
- **Test strategy:** the richest unit surface in the system, entirely in-process with fakes for all three ports: each check individually forced to fail while others pass (AT23 single-check rollback); each Inconclusive cause forced (AT22 — `[HYPOTHESIS]`-marked); skew at/beyond tolerance (AT8 logic); staleness vs R and the S4 ceiling (AT13/AT14 *logic* — end-to-end execution stays blocked on S1); baseline accept/reject (AT2, AT3, AT6, AT7, AT9, AT11 logic).
- **Failure boundaries:** verdicts are values; no error return path exists apart from the verdict causes; a condition outside the cause enumeration is a contract-amendment event, not a new ad-hoc error. No panic across the boundary; no retry; no fallback ladder (RFC-003 §13).
- **Logging boundaries:** the decision trace is the only output channel, emitted unconditionally (Accepts included). No log statements — a check that "logs a warning" instead of tracing is FM11's dark corner.
- **Metrics boundaries:** none in-module. Verification latency (AT26) is measured at the driver boundary around `Verify`, never inside it — measurement must not perturb the decision path.

## M4 — `internal/truststore`

- **Responsibilities:** hold out-of-band-provisioned trust material; answer `material | absent` per trust domain; never fetch; record provisioning events. (ER7, INV7, NFR2, C6, FM9.)
- **Ownership:** track B engineer; RP operator at runtime.
- **Public interface:** `TrustMaterialFor(domain) → material | absent`; `Provision(material) ` (out-of-band path); provisioning record type. Satisfies `verify.TrustMaterialPort` structurally. Full contract: §4.
- **Private responsibilities:** in-memory holding structure; provisioning-record assembly.
- **Dependencies:** `internal/record` (trust-material vocabulary), stdlib.
- **Allowed imports:** stdlib (no `net/*`); `internal/record`; `go-spiffe/v2` bundle types.
- **Forbidden imports:** `net`, `net/http` — *this package must be incapable of fetching* (FM9's insecure-fallback path is deleted by construction); `internal/verify`; `internal/revstatus`.
- **Lifecycle:** empty at construction; populated only by explicit `Provision` acts; `absent` answers before provisioning are correct behavior, not startup errors; material withdrawal is a provisioning act too.
- **Configuration:** none. Storage substrate questions (persistence) are deferred; Sprint 1 holds material in memory behind the same interface.
- **Test strategy:** unit: hit/miss/withdrawn; the never-fetch property is enforced by the import lint (not testable by absence, so it is made structural); provisioning-record completeness.
- **Failure boundaries:** `absent` is an answer, not an error; malformed provisioned material is refused at `Provision` time with a named cause — never stored, never half-trusted.
- **Logging boundaries:** provisioning records only (RFC-003 §14).
- **Metrics boundaries:** none.

## M5 — `internal/revstatus`

- **Responsibilities:** the volatile region's fixed contract: `StatusOf(instanceID) → ObservablyRevoked(asOf) | NotObservedRevoked(asOf) | Indeterminate`; mandatory as-of on knowledge answers; the honest-indeterminate rule (ignorance is never `NotObservedRevoked`); the degenerate always-indeterminate realization (outcome-β representability, RFC-003 §6); the realization-independent contract test suite. (ER5, SO1, FM2, FM4, INV12, AP5, AP7, AP12.)
- **Ownership:** track B engineer for contract + degenerate realization. The spike-selected realization (post-EXP-001) gets its own owner and its own allowed-imports ruling under the Revocation Mechanism decision — not specified here (RFC-003 E1).
- **Public interface:** the answer types; `StatusOf`; sub-package `contracttest` exporting the suite any realization must pass. Satisfies `verify.RevocationStatusPort` structurally. Full contract: §5.
- **Private responsibilities:** none in Sprint 1 beyond the degenerate realization (which has no internals by definition).
- **Dependencies:** `internal/record` (InstanceID vocabulary), stdlib.
- **Allowed imports (contract + degenerate + contracttest):** stdlib (no `net/*`); `internal/record`.
- **Forbidden imports:** `net/*` in everything Sprint 1 ships (a future realization's maintenance path may be granted network imports by the Revocation Mechanism decision — that grant will be explicit, file-scoped, and never reachable from `StatusOf`'s answer path); `internal/verify`; `internal/truststore`; `internal/revorigin` (no direct coupling — the propagation channel is deferred, RFC-003 R5).
- **Lifecycle:** degenerate realization: construct, answer `Indeterminate`, forever. Real realizations: maintenance lifecycle defined by the spike outcome; as-of freezes during partition as a *defined* condition (RFC-003 §11).
- **Configuration:** none in the contract. Realization configuration (refresh cadence etc.) arrives with the spike outcome and is confined to the realization.
- **Test strategy:** `contracttest` suite: closed answer set, mandatory as-of, honest-indeterminate, determinism per view-state; run against the degenerate realization in Sprint 1 and against every candidate realization the spike attempts later — same seam, same suite.
- **Failure boundaries:** `Indeterminate` *is* the failure boundary — every internal problem surfaces as it; no error type escapes this package.
- **Logging boundaries:** the as-of disclosure inline in answers (RFC-003 §14); realization-internal operational logging is a later, realization-scoped ruling.
- **Metrics boundaries:** none in the contract; freshness-lag metrics belong to the future realization + harness.

## M6 — `internal/revorigin`

- **Responsibilities:** the authoritative, append-only revocation register: `Revoke(instanceID) → recorded` (idempotent on terminal state — INV4); expose the register's content for the (deferred) propagation channel; touch nothing but instance IDs (INV5/INV6 structural: the register cannot affect identities or siblings because it does not know them). (FR4, ER5, INV4–INV6, FM1.)
- **Ownership:** track B engineer; domain-A operational owner at runtime.
- **Public interface:** `Revoke(instanceID)`; `View() → ordered revocation entries` (read surface for propagation/reconstruction). Full contract: §6.
- **Private responsibilities:** append-only enforcement; ordering.
- **Dependencies:** `internal/record` (InstanceID), stdlib.
- **Allowed imports:** stdlib (no `net/*`); `internal/record`.
- **Forbidden imports:** everything else `internal/`; `net/*` (publication is the propagation channel's business, not the register's).
- **Lifecycle:** register persists append-only; outlives every delegation it names (reconstruction support); no compaction, no deletion.
- **Configuration:** none. Persistence substrate deferred; Sprint 1 in-memory behind the same interface.
- **Test strategy:** unit: append-only property (re-revoke is a no-op on terminal state — AT11 origin-side logic); ordering; view stability; sibling-independence by construction (register holds only IDs).
- **Failure boundaries:** `Revoke` cannot fail semantically (unknown IDs are recordable — the register does not validate existence; a revocation of a never-issued instance is inert by M3's keying); storage failure surfaces as a named refusal, never a silent drop.
- **Logging boundaries:** the register **is** the observable (append-only, RFC-003 §14); no separate log.
- **Metrics boundaries:** none.

---

## Non-module code

### `cmd/` — boundary drivers (composition roots)

`cmd/atlas-issue`, `cmd/atlas-verify`, `cmd/atlas-revoke`. Thin shells that wire ports to providers and expose the product boundaries the acceptance plan's test loci name (issuance boundary, RP verification boundary, revocation act). They own: wiring, policy loading (scope-act values), operational logging, and the AT26 latency measurement around `Verify`. They contain no delegation logic — a check performed in a driver is a conformance violation. Allowed imports: everything `internal/`; forbidden: business logic (enforced by review, not lint). Trace: AT test loci; RFC-003 §12 flows.

### `tests/harness` + `tests/acceptance` + `tests/fixtures`

The AT scaffold. `harness/` defines substrate-control interfaces (domain control, partition induction, egress observation, clock skew) whose implementations arrive with the substrate block (shared with EXP-001, plan Phases 2–6). `acceptance/` holds one file per AT family; in-process ATs run in Sprint 1, substrate ATs compile and skip with a named blocker. `fixtures/` holds records, trust material, and mutation corpora. The harness may import anything (it is instrumentation, not product); acceptance tests exercise product only through `cmd/` boundaries or public package surfaces, mirroring each AT's stated test locus.

## Dependency-rule enforcement (W5)

The forbidden-imports column above is mechanically checked in CI (a small lint over `go list` import graphs) and locally via `make ci`. The check *is* RFC-003 rules R1–R7 in executable form; changing the allowed/forbidden tables requires amending this specification with an RFC-003 trace, then the lint config — in that order.

<!-- checkpoint: rfc(glossary-definitions): audit glossary definitions -->

<!-- checkpoint: repo(trust-anchors): improve trust anchors -->

<!-- checkpoint: context(attenuation-specification): audit attenuation specification -->

<!-- checkpoint: repo(threat-model-scenarios): audit threat model scenarios (#17) -->

<!-- checkpoint: refactor(verify): refactor error wrappers -->

<!-- checkpoint: chore(revstatus): clean key derivation -->

<!-- checkpoint: fix(test): fix integration test runner -->
