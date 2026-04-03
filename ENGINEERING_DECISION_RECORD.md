# Engineering Decision Record — Atlas

**Status:** Canonical register of every architectural decision: accepted, rejected, deferred — with alternatives, reasons, trade-offs, and complexity cost. Architecture-phase closure artifact. New entries require a frozen-package trace or a founder act; entries are append-only (superseding, never editing).
**Format:** AD-NNN · decision · status · alternatives considered · reason · trade-off accepted · complexity cost (Low/Med/High — the ongoing cost of living with the decision).

---

## Accepted

**AD-001 — Six-module boundary-aligned decomposition with the revocation region isolated.**
*Alternatives:* (a) single component per domain; (b) delegation checks embedded as an extension of the RP's existing verifier. *Reason:* (a) fails AP12 on its face — the spike-volatile region would live inside everything, so a spike outcome or S2/S3 resolution forces edits across the requirement-fixed region; (b) cannot be specified without naming a baseline implementation (DR9) and endangers SO7's own pass metric ("baseline present and unmodified"). *Trade-off:* one more boundary than a monolith. *Complexity cost:* Low — every part traces; no part is unforced. (RFC-003; ARR §5–6.)

**AD-002 — The presented unit and the reconstruction record are one artifact.**
*Alternatives:* separate presentation token + audit record (1:1). *Reason:* ER1 ("single presentable unit") + ER4 (self-sufficient record) read naturally as one; a second artifact is unforced (AP11) and doubles the tamper-evidence surface. *Trade-off:* if a forcing item surfaces during implementation, splitting ripples through §1–§3 interfaces — the largest un-hedged bet in the architecture, carried openly as risk ER-4 rather than hedged with an abstraction layer. *Complexity cost:* Low now; High only if reversed.

**AD-003 — Verification Core is pure: all inputs injected, zero I/O, consumer-defined ports, structural satisfaction.**
*Alternatives:* verifier that consults its stores directly; provider-defined interfaces. *Reason:* makes AP1 (offline) structural rather than tested-for; makes the entire verification logic testable with no substrate; makes the spike invisible to M3 (R3); AT16's zero-egress becomes an import-lint property. *Trade-off:* composition roots carry wiring boilerplate. *Complexity cost:* Low.

**AD-004 — Freshness/skew/ceiling policy lives only in M3 (rule R8).**
*Alternatives:* each M5 realization enforces R itself. *Reason:* R is a founder parameter (S1); one home means one change point and no divergent enforcement across realizations (FM2/FM4 bound the *verifier's acceptance*, not the provider's reporting). *Trade-off:* M5 answers must carry as-of freshness. *Complexity cost:* Low.

**AD-005 — Expiry is derived arithmetic, never stored state.**
*Alternatives:* lifecycle store tracking expiry transitions. *Reason:* a stored flag can be absent, stale, or rolled back; arithmetic on an immutable record + injected time cannot (INV3 monotonicity by construction). *Trade-off:* correctness of the time port + skew tolerance becomes load-bearing (already required by ER3). *Complexity cost:* Low.

**AD-006 — Revocation state is deliberately dual: authoritative at M6, observed at M5; the gap is the observability state.**
*Alternatives:* single "revocation status" store pretending a global state. *Reason:* the gap between register and view **is** RFC-002 §9.3's state machine, bounded by R and S4; pretending a global state exists is exactly the over-claim INV12 forbids. *Trade-off:* two stores instead of one. *Complexity cost:* Low — the duality is the honest model.

**AD-007 — Degenerate always-Indeterminate M5 realization ships first and is the pre-spike default wiring.**
*Alternatives:* stub returning `NotObservedRevoked`; no default realization. *Reason:* the stub would launder ignorance into knowledge (AP5 violation, and the exact FM11 silent-failure path); no-default leaves outcome β/δ unrepresentable. With the degenerate realization the system fails closed rather than pretending revocation knowledge. *Trade-off:* pre-spike end-to-end Accepts require a test fake with a fresh view (acceptable: that is precisely what "revocation knowledge exists" means). *Complexity cost:* Low.

**AD-008 — `contracttest` suite gates every M5 realization (the plugin admission rule).**
*Alternatives:* per-realization ad-hoc tests. *Reason:* the spike attempts up to three compositions; one suite at the production seam means every candidate is judged by identical criteria (lab pre-registration discipline extended to code). *Complexity cost:* Low.

**AD-009 — All modules under `internal/`; drivers as the only composition roots; no shared util package.**
*Alternatives:* public library packages now; a common/util package. *Reason:* V1 is a reference implementation (C4) — external importability is a future founder act, mechanical by relocation because coupling is lint-forbidden; a util package is an untraced module (AP11) and a coupling magnet. *Trade-off:* small duplications tolerated over shared helpers. *Complexity cost:* Low.

**AD-010 — Dependency rules R1–R7 enforced by a build-failing import lint.**
*Alternatives:* review-convention enforcement. *Reason:* the architecture's load-bearing properties (offline, isolation, no cross-domain calls) must not depend on reviewer vigilance (FM11). *Complexity cost:* Low (one script, table-driven).

**AD-011 — Stack: Go 1.21+ with `go-spiffe/v2`, `spire-api-sdk`, `go-jose/v3` — acknowledged from `AI_BOOTSTRAP.md` §4, not chosen here.**
*Alternatives:* none evaluated — the bootstrap fixes it as a Non-Negotiable Development Rule and the founder prohibited a technology-selection RFC, consistent with the ruling standing. The Sprint-0 "stack deferred post-spike" placeholders (Makefile, tests/README) are superseded by the newer bootstrap text and are updated in epic E1. *Complexity cost:* n/a (recorded ruling).

**AD-012 — Record integrity envelope: JWS-signed token carrying the delegation claims (assumption A3).**
*Alternatives:* bespoke signature scheme; detached-signature record. *Reason:* gate C2/C3/C7 record token format, scope/expiry claims, and JWS tamper-rejection as **solved** components (High confidence, fetched primary sources); `go-jose/v3` is the vetted library for exactly this; a bespoke scheme violates the vetted-tooling rule and TP6. *Trade-off:* inherits the JWS family's known failure class (algorithm confusion) → mitigated by algorithm pinning at validation (task E2, risk SR-2). *Pending:* founder confirmation F1 (flagged, non-blocking). *Complexity cost:* Low.

**AD-013 — Instance identity: opaque element, unique per issuance, minted through a construction-injectable seam.**
*Alternatives:* semantic instance identity (fixing FM5's open question ourselves); no instance identity (revocation keys to content tuple). *Reason:* fixing semantics is untraced invention (FM5 explicitly declines); keying revocation to the content tuple makes distinct issuances indistinguishable, breaking INV6 targeting. Opacity + uniqueness is the minimal interim; the injectable minter makes tests deterministic and gives the future FM5 amendment a single landing point. *Pending:* founder confirmation F2. *Complexity cost:* Low.

**AD-014 — M2 consumes the same time-reading port contract as M3.** *(Cross-review gap G1 — new in this set.)*
*Alternatives:* M2 reads the system clock directly. *Reason:* issuance time is record content (ER4 "at what time"); an uninjected clock makes issuance untestable under controlled time and smuggles I/O into an otherwise port-disciplined module. *Complexity cost:* Low.

**AD-015 — The record reserves an opaque, optional revocation-binding element; M2 obtains it via a `RevBindingSource` port (empty pre-spike); only M5 realizations interpret it.** *(Cross-review gap G2 — the six-month-rewrite fix.)*
*Alternatives:* (a) no slot — add mechanism fields when the spike decides; (b) model each candidate's fields now. *Reason:* (a) is the identified rewrite: status-list-class compositions need a per-record status reference and accumulator-class may need witness data, so the spike outcome would mutate the stable surface — precisely what AP12 exists to prevent; (b) builds three mechanisms' worth of unforced structure (AP11, DR9). An opaque slot is the minimal structure that keeps the stable surface stable across all gate candidates. *Traces:* FM2/FM4 (a composition must be realizable), AP12 (volatility isolation), gate C4 candidate list. *Trade-off:* one element rides the record that V1's pre-spike phase never populates; the slot's opacity must be defended in review (nothing but M5 realizations may ever interpret it). *Complexity cost:* Low; the alternative's cost was High.

**AD-016 — Traces are returned values; persistence and shipping belong to the invoking boundary's operator.** *(Cross-review gap G3.)*
*Alternatives:* modules write traces to a sink/log dependency. *Reason:* keeps M1–M6 free of any logging/IO dependency (purity, R6), makes traces first-class for AT23/AT30, and puts retention where operational ownership already is (drivers/harness). *Complexity cost:* Low.

**AD-017 — Concurrency posture declared: M1/M3 pure-safe; M4/M5/M6 single-writer/multi-reader with serialized mutation; drivers one-request-per-invocation in V1.** *(Cross-review gap G4.)*
*Alternatives:* silence (the prior state); full concurrent hardening now. *Reason:* silence is an unstated assumption waiting to be violated; hardening now is production machinery for a feasibility horizon (C4/AP11). Declared posture is testable and honest; concurrent-load behavior is recorded debt (TD-5), not implied capability. *Complexity cost:* Low.

**AD-018 — Substrate is shared between EXP-001 and acceptance testing.**
*Alternatives:* separate clean-room AT substrate. *Reason:* identical topology requirements (two SPIRE domains, link-level partition, egress observation, NTP discipline — EXP-001 plan Phases 2–6; AT execution discipline); building it twice doubles the dominant infrastructure cost. *Trade-off:* lab discipline (pre-registration, evidence hygiene) governs AT runs too — acceptable, it is stricter, not looser. *Change-condition:* if lab governance demands isolation, the sharing dissolves and the cost returns (recorded in Sprint plan risk table). *Complexity cost:* Low.

## Accepted — surfaced during E3 implementation (Verification Core)

Per the EDR closing rule, decisions forced by implementation are recorded here, not resolved silently in code. AD-019…AD-024 arose building `internal/verify`; none contradicts the frozen package, and each resolves a tension between the two unfrozen architecture-set documents (`INTERFACE_SPECIFICATION.md` §1 vs §3) or fills an under-specified mechanism point.

**AD-019 — `SignatureUnverifiable` is a reserved, not-produced cause in V1.**
*Context:* the interface spec §3 lists `SignatureUnverifiable` as an inconclusive integrity outcome, but M1's `ValidateIntegrity` (built in E2) is binary Intact|Altered and folds unknown-key-ID into definitive Altered (FM5 non-objective: no warranted key-rotation state). *Decision:* integrity produces `TrustMaterialAbsent` (inconclusive, no material for the domain) or `IntegrityFailed` (definitive) in V1; `SignatureUnverifiable` is retained in the closed cause set but never produced, reserved for a future key-rotation state that could distinguish "key not yet held" from "wrong key." *Alternatives:* remove the member (a set change, and premature — the future need is real); make it reachable now (requires inventing rotation state — DR9 violation). *Trade-off:* one enumerated member is currently unreachable, documented in `cause.go`. *Complexity cost:* Low.

**AD-020 — Revocation-answer vocabulary is owned by the verifier (consumer), not the provider.**
*Context:* dependency rule R3 forbids `verify` importing `revstatus` and vice-versa; the shared `RevocationStatus` type cannot live in either without a cycle. *Decision:* `verify` defines `RevocationState`/`RevocationStatus` as its input vocabulary (interfaces defined where consumed). At E4/E5, `revstatus` either produces these via a `cmd`-level adapter or the type is promoted to `record`; deferred to the E4 integration point. *Alternatives:* put the type in `record` now (bloats the stable surface with an M5 concept prematurely). *Trade-off:* the E4↔M3 wiring form is an open integration decision (small). *Complexity cost:* Low.

**AD-021 — Policy is `{R, SkewTolerance}`; the S4 ceiling is satisfied-by-construction, not a third parameter.**
*Context:* the architecture named a three-member policy `{R, skew, S4 ceiling}`. *Decision:* V1 policy carries R and skew only. The revocation answer set carries no partition discriminator, so a partition manifests solely as an aging as-of that R already governs; failing closed on staleness > R **is** the S4 guarantee (no observability claimed that cannot be had, INV12). A separate ceiling would require a partition signal the honest answer set deliberately withholds. *Alternatives:* keep a third numeric knob (inoperative — nothing to apply it to). *Trade-off:* amends the policy definition in the interface/module specs (dated note). *Complexity cost:* Low.

**AD-022 — M1 exposes `PeekTrustDomainUnverified` for trust-material selection.**
*Context:* selecting which trust material to verify with needs the record's trust domain, but reading is defined only on authenticated records (chicken-and-egg). *Decision:* `record` exposes an explicitly-unauthenticated peek of the principal's trust domain, usable ONLY for material selection; `ValidateIntegrity` remains sole authority (a lie about the domain selects wrong material and fails verification — it can never wrongly succeed). Standard JWT/SPIFFE practice (unverified header peek for key selection). *Alternatives:* scatter envelope parsing into `verify` (leaks format out of M1) or into every port provider (worse). *Trade-off:* one additive, clearly-labeled function on the stable surface. *Complexity cost:* Low.

**AD-023 — The identity-binding check enforces principal ≠ delegate.**
*Context:* integrity already guarantees both identities are present and well-formed, leaving the binding stage with no independently-falsifiable predicate — which would make SO5 single-check rollback untestable for binding. *Decision:* the binding check's verification-time predicate is distinctness (a record binding an identity to itself is not a delegation) — a defensible reading of INV1's "exactly one principal and exactly one delegate." *Founder-reviewable:* if self-delegation is ever wanted, this is a one-line change. *Trade-off:* a reading, not a frozen requirement — flagged for founder awareness. *Complexity cost:* Low.

**AD-024 — Scope-integrity and signature-integrity share the INV8 guarantee; scope's rollback manifests through integrity.**
*Context:* the scope field is under the record signature, so it cannot be tampered while the signature still verifies — "scope integrity fails while signature passes" is cryptographically impossible. *Decision:* the scope stage's independent job is FR2 inspectability (scope present, well-formed) — always passing for authentic records; scope tampering is caught by the integrity stage (`IntegrityFailed`), which is the honest cryptographic truth. SO5's "no single scope-level failure yields acceptance" is satisfied via that path (tested in `rollback_test.go`). *Alternatives:* invent a separate scope MAC (unforced duplication, AP11). *Trade-off:* two stages share one cryptographic guarantee; documented. *Complexity cost:* Low.

**AD-025 — Execution order runs integrity first; the trace presents canonical label order.**
*Context:* checks are labelled 1–5 with identity-binding first, but binding/expiry/scope/revocation all require an authenticated read. *Decision:* execute integrity first as the gate; when it does not pass, read-dependent stages are recorded `NotEvaluated` (honest, not omitted); the trace lists all five in canonical label order. Verdict routing is order-independent (definitive dominates inconclusive dominates pass), so the label/execution divergence has no effect on the verdict. *Complexity cost:* Low.

## Rejected

**AD-R01 — Monolithic per-domain component.** Rejected: AP12 violation (see AD-001).
**AD-R02 — Embedding delegation checks inside the RP's existing verifier.** Rejected: DR9 violation + SO7 self-endangerment (see AD-001).
**AD-R03 — Live or cached fetch fallback when trust material is absent.** Rejected: FM9's insecure-fallback path; M4 is made structurally incapable of fetching. The fallback ladder (retry → refetch → accept-with-warning) is rejected wholesale (RFC-003 §13).
**AD-R04 — Verdict caching / verification memoization.** Rejected: RFC-002 §9.2 — a re-presentation is a new verification; a cache is a place where a revocation goes unnoticed (FM2 through the back door) and an unforced mechanism (AP11).
**AD-R05 — A "record registry" service tracking issued delegations.** Rejected: the record is self-sufficient (INV9); a registry adds a liveness dependency AP1 forbids on the verification path, and an unforced component elsewhere.
**AD-R06 — Stored delegation-lifecycle state machine.** Rejected in favor of derived expiry + dual revocation state (AD-005/AD-006).
**AD-R07 — Configuration file format selection.** Rejected as unforced (TP6): nothing in V1's ATs requires one; drivers take parameters by the simplest available means.
**AD-R08 — Shared util/common package.** Rejected: untraced module, coupling magnet (AD-009).
**AD-R09 — Pre-building multi-domain (≥3) or multi-hop structure "while we're at it."** Rejected: TP5/AP6/S5; explicitly a non-extension-point (P-table, `SYSTEM_ARCHITECTURE.md` §10).
**AD-R10 — Asserting fail-closed as a guarantee** (in docs, trace naming, or test phrasing). Rejected: DR7 — the Inconclusive→Reject routing is designed-for, `[HYPOTHESIS]`-marked, resolved only by a V1 act.
**AD-R11 — An AT for FM5 (key compromise) or FM8 (within-window replay) resistance.** Rejected: the AT plan itself forbids it — asserting resistance the Product Definition does not warrant is a doctrine violation. The gaps stay visible (honest-claims paragraphs everywhere), not tested-away.

## Deferred (with owner and trigger)

| # | Decision | Owner | Trigger |
|---|---|---|---|
| AD-D01 | S1–S4 parameter values (R; cached-pull; broker; partition reading) | Founder | scope-act journal entry — gates E6/E7 |
| AD-D02 | Revocation composition + M5 real realization + propagation channel + RevBinding production rule | Founder acceptance of EXP-001 outcome | spike completion |
| AD-D03 | Instance-identity semantics (FM5 open question) | Founder (frozen-package amendment) | amendment act; lands in the AD-013 seam |
| AD-D04 | Persistence substrates for M4/M5/M6 | Engineering, post-V1 or when an AT forces it | TD-1 repayment trigger |
| AD-D05 | Permission Source realization + operator (two-domain experiment) | Engineering + runbook | epic E5 |
| AD-D06 | Trust-material provisioning procedure (runbook) | Engineering + RP operator | epic E5/E6 |
| AD-D07 | Fail-closed hypothesis resolution (promote/redirect the Inconclusive routing) | Founder, on V1 evidence (AT22 results) | V1 closure |
| AD-D08 | Library extraction for external adoption (SO7 path) | Founder | post-V1 scope act |
| AD-D09 | Concurrency hardening beyond declared posture | Engineering | post-V1 (C4 lifts) |
| AD-D10 | Documentation-defect corrections (bootstrap 100 ms line; stale `agents/agents/` paths in RFC provenance; `DEVELOPMENT_RULES.md` RFC-policy text; `CLAUDE.md` layout tree) | Founder (owned files / amendment path) | founder editing pass — zero architectural effect meanwhile |

## Decision-making rule going forward

The architecture phase is closed. A new AD entry is legitimate only when (a) a frozen requirement or founder act forces it, or (b) implementation surfaces a contradiction with the frozen package — in which case the entry records the contradiction and the founder resolves it (never silently in code). Absence of an AD entry means the existing set governs.
