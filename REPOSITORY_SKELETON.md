# Repository Skeleton

**Status:** Engineering document. The complete target tree for the implementation-ready repository. Folders and filenames only — no implementation, no placeholder logic, no TODO comments exist or will exist in any listed file (Go `doc.go` files carry package documentation, which is specification, not implementation).
**Authority:** RFC-003 (module set, dependency rules), `PROJECT_MODULE_SPECIFICATION.md` (package boundaries), `AI_BOOTSTRAP.md` §4 (Go 1.21+ — file extensions follow).
**Legend:** `(existing)` = already in the repository, unchanged unless marked `(existing, updated)`. Everything else is new. Nothing is created by this document; creation is task T0 in `IMPLEMENTATION_ORDER.md`.

```text
atlas/
├── AI_BOOTSTRAP.md                          (existing)
├── ARCHITECTURE_READINESS_REVIEW.md         (existing)
├── CLAUDE.md                                (existing)
├── CODEGRAPH.md                             (existing)
├── CONTRIBUTING.md                          (existing)
├── DEVELOPMENT_RULES.md                     (existing)
├── Dockerfile                               (existing)
├── ENGINEERING_SPRINT_1_PLAN.md             (existing — this sprint's set)
├── FREEZE_AMENDMENT_PLAN.md                 (existing)
├── FROZEN.sha256                            (existing)
├── IMPLEMENTATION_ORDER.md                  (existing — this sprint's set)
├── LEVEL0_1_FEASIBILITY_GATE.md             (existing, frozen)
├── MODULE_INTERFACE_SPECIFICATION.md        (existing — this sprint's set)
├── Makefile                                 (existing, updated: real test/build/vet/importlint targets)
├── P5_FALSIFICATION_EXPERIMENT.md           (existing, frozen)
├── PROJECT_MODULE_SPECIFICATION.md          (existing — this sprint's set)
├── PROJECT_STRUCTURE.md                     (existing, updated: records this tree)
├── README.md                                (existing)
├── REPOSITORY_AUDIT.md                      (existing)
├── REPOSITORY_FINAL_STATE.md                (existing)
├── REPOSITORY_SKELETON.md                   (existing — this document)
├── SECURITY.md                              (existing)
│
├── go.mod                                   Go module definition (module path, go 1.21)
├── go.sum                                   dependency checksums (created by the toolchain)
│
├── cmd/
│   ├── atlas-issue/
│   │   └── main.go                          issuance boundary driver (wiring only)
│   ├── atlas-verify/
│   │   └── main.go                          RP verification boundary driver (wiring + AT26 measurement point)
│   └── atlas-revoke/
│       └── main.go                          revocation act driver (wiring only)
│
├── internal/
│   ├── record/                              M1 — stable surface (pure)
│   │   ├── doc.go                           package charter + RFC-003 M1 trace
│   │   ├── record.go                        Record type; assertion accessors
│   │   ├── instanceid.go                    opaque InstanceID type (semantics deferred — E2)
│   │   ├── integrity.go                     ValidateIntegrity: Intact | Altered
│   │   ├── envelope.go                      JWS envelope construction/parsing (assumption A3; private)
│   │   ├── reconstruct.go                   Read on intact records (third-party reconstruction)
│   │   ├── record_test.go                   round-trip create/read; opacity of InstanceID
│   │   ├── integrity_test.go                intact/altered basics
│   │   ├── mutation_test.go                 AT20 mutation corpus: flips, substitutions, truncation, reorder
│   │   └── reconstruct_test.go              AT19/AT21 logic: record-alone sufficiency
│   │
│   ├── issuance/                            M2 — sole record creator
│   │   ├── doc.go                           package charter + trace
│   │   ├── authority.go                     Issue: Record | Refused(cause)
│   │   ├── ports.go                         PermissionSource port (consumer-defined)
│   │   ├── refusal.go                       closed refusal-cause set
│   │   ├── trace.go                         IssuanceTrace type + assembly
│   │   ├── authority_test.go                subset guard; ephemeral issuance (AT18 logic)
│   │   └── refusal_test.go                  AT4 logic: over-scope, permissions-unavailable
│   │
│   ├── verify/                              M3 — conformant verifier (pure; ports)
│   │   ├── doc.go                           package charter + conformance statement
│   │   ├── verifier.go                      Verify entry point; pipeline order
│   │   ├── policy.go                        Policy {R, skew tolerance, S4 ceiling}; refuse-unset construction
│   │   ├── ports.go                         TrustMaterialPort, RevocationStatusPort, TimePort
│   │   ├── verdict.go                       Accept | Reject | Inconclusive→Reject [HYPOTHESIS] routing
│   │   ├── cause.go                         closed cause enumeration, per check
│   │   ├── trace.go                         DecisionTrace: per-check entries, unconditional
│   │   ├── check_binding.go                 check 1 — identity binding (INV1)
│   │   ├── check_integrity.go               check 2 — signature/tamper via record.ValidateIntegrity (INV8)
│   │   ├── check_expiry.go                  check 3 — expiry ± stated skew (INV3, ER3)
│   │   ├── check_scope.go                   check 4 — scope integrity (INV8; not subset re-derivation)
│   │   ├── check_revocation.go              check 5 — revocation answer under freshness policy (SO1, FM2/FM4)
│   │   ├── verifier_test.go                 baseline accept/reject paths (AT2/AT3/AT6/AT7 logic)
│   │   ├── rollback_test.go                 AT23: each check forced to fail while others pass
│   │   ├── inconclusive_test.go             AT22 [HYPOTHESIS]: every inconclusive cause → Reject
│   │   ├── expiry_skew_test.go              AT8 logic: within/at/beyond tolerance
│   │   ├── revocation_policy_test.go        AT13/AT14 logic at arbitrary R and S4 ceiling (parametric)
│   │   └── trace_test.go                    unconditional trace, incl. Accepts; per-check attribution
│   │
│   ├── truststore/                          M4 — RP-local trust material
│   │   ├── doc.go                           package charter + trace
│   │   ├── store.go                         TrustMaterialFor: material | absent (never fetch)
│   │   ├── provision.go                     out-of-band Provision + provisioning record
│   │   └── store_test.go                    hit / miss / withdrawn; refusal of malformed material
│   │
│   ├── revstatus/                           M5 — volatile region behind fixed contract
│   │   ├── doc.go                           package charter + E1 deferral statement
│   │   ├── answer.go                        closed answer set + mandatory as-of
│   │   ├── indeterminate.go                 degenerate always-Indeterminate realization (outcome-β honesty)
│   │   ├── indeterminate_test.go            degenerate realization passes the contract
│   │   └── contracttest/
│   │       └── suite.go                     realization-independent contract suite (run by every realization)
│   │
│   └── revorigin/                           M6 — authoritative append-only register
│       ├── doc.go                           package charter + trace
│       ├── register.go                      Revoke (idempotent-terminal) + View (ordered, read-only)
│       └── register_test.go                 append-only; re-revoke no-op; ordering; sibling-independence
│
├── tests/                                   (existing directory; README replaced)
│   ├── README.md                            (existing, updated: scaffold map, AT→file index, skip policy)
│   ├── acceptance/
│   │   ├── identity_test.go                 AT1–AT2
│   │   ├── scope_test.go                    AT3–AT5
│   │   ├── expiry_test.go                   AT6–AT8
│   │   ├── revocation_test.go               AT9–AT14 (AT13/AT14 skip: named blocker S1/substrate)
│   │   ├── offline_test.go                  AT15–AT16 (skip: substrate — egress instrumentation)
│   │   ├── crossdomain_test.go              AT17 (skip: substrate — two SPIRE domains)
│   │   ├── ephemeral_test.go                AT18
│   │   ├── reconstruction_test.go           AT19–AT21
│   │   ├── failclosed_test.go               AT22–AT23 ([HYPOTHESIS]-marked per AT plan)
│   │   ├── interop_test.go                  AT24 (AT25 deferred per D1 — recorded, not scaffolded)
│   │   ├── latency_test.go                  AT26 (skip: substrate + driver measurement point)
│   │   └── boundary_test.go                 AT27–AT29 (AT30 is a process act — indexed in README, not code)
│   ├── harness/
│   │   ├── doc.go                           harness charter: instrumentation, not product
│   │   ├── substrate.go                     two-domain substrate control interface (impl: substrate block)
│   │   ├── partition.go                     partition induction interface
│   │   ├── egress.go                        egress observation interface (AT16)
│   │   ├── clock.go                         controllable TimePort for skew scenarios
│   │   └── fakes.go                         port fakes for in-process ATs
│   └── fixtures/
│       ├── records/                         valid/invalid record corpus (generated by tests, committed forms)
│       ├── trustmaterial/                   bundle fixtures for two synthetic domains
│       └── mutations/                       AT20 mutation corpus definitions
│
├── scripts/
│   ├── check-frozen-docs.sh                 (existing)
│   ├── frozen-docs.list                     (existing)
│   └── check-imports.sh                     dependency-rule lint (RFC-003 R1–R7 as CI gate)
│
├── .github/
│   └── workflows/
│       ├── repo-health.yml                  (existing)
│       └── go-ci.yml                        build + vet + test + import lint
│
├── agents/                                  (existing — unchanged)
├── archive/                                 (existing — unchanged)
├── context/                                 (existing — 08_AI_HANDOFF updated at sprint close per governance)
├── docs/                                    (existing, frozen sets — unchanged)
├── lab/                                     (existing — EXP-001 untouched; substrate work lands there later)
└── rfc/                                     (existing + RFC-003)
```

## Design notes (why the tree is shaped this way)

1. **`internal/` for all six modules.** Nothing outside this repository may import them in V1 — the reference-implementation posture (C4: no production claim) made structural. Lifting `verify` + `record` to an importable library is a *future founder act* (the SO7 companion-library path), done by relocation, not redesign.
2. **One package per RFC-003 module, no shared `util`/`common` package.** A utility package is an untraced module (AP11) and a coupling magnet; anything two modules both need belongs in `record` (vocabulary) or nowhere.
3. **Consumer-defined ports; structural satisfaction.** `verify/ports.go` and `issuance/ports.go` define what their packages *need*; `truststore`/`revstatus` satisfy them without importing `verify` (RFC-003 R3). The wiring point is `cmd/` — the only place realizations and consumers meet.
4. **Tests live with their module; acceptance tests live at the boundary.** Unit tests exercise package surfaces; `tests/acceptance/` exercises the AT plan's stated test loci through drivers/fakes and never reaches into internals — mirroring SO8's reviewer, who has spec + build, not source-level privilege.
5. **Skips carry blockers, not TODOs.** A substrate- or scope-act-blocked AT compiles and skips with the named blocker (`S1-scope-act`, `substrate`, `spike-outcome`) so `IMPLEMENTATION_ORDER.md`'s unblock map stays executable and nothing rots silently.
6. **`scripts/check-imports.sh` is the architecture's teeth.** The forbidden-import tables in `PROJECT_MODULE_SPECIFICATION.md` fail the build when violated; the dependency rules do not depend on reviewer vigilance.
