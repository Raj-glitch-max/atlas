# Implementation Order

**Status:** Engineering document. The ordered task list for turning the skeleton into the Sprint 1 exit state, and the forward order beyond it. Not frozen.
**Authority:** RFC-003 dependency rules; `PROJECT_MODULE_SPECIFICATION.md`; `MODULE_INTERFACE_SPECIFICATION.md`; `ENGINEERING_SPRINT_1_PLAN.md` scope.
**Discipline:** one task = one conventional commit; **every task leaves `make ci` green** (each compiles independently and adds only passing tests); every milestone leaves the repository in a working, demonstrable state. Optimized for small commits, high testability, low coupling, and parallel development along the three tracks below.

---

## 1. Tracks (parallelizable by construction)

The dependency rules make three tracks independent after T0 — they touch disjoint packages and meet only at milestone M4:

- **Track A — stable surface and verifier:** `record` → `verify` (the critical mass of frozen-requirement logic; pure, no substrate).
- **Track B — stores and register:** `truststore`, `revstatus` (+ `contracttest`), `revorigin` (small, independent packages).
- **Track C — foundation and scaffold:** CI, import lint, Makefile, harness fakes, fixtures, acceptance scaffold.

A single engineer runs them interleaved; two or three engineers run them concurrently with no coordination beyond T0 and M4.

## 2. Task list

### Phase 0 — repository foundation (Track C; blocks everything; ~half a day)

| ID | Task | Compiles/green because | Commit |
|---|---|---|---|
| T0.1 | Create the tree per `REPOSITORY_SKELETON.md`: `go.mod`, all package directories, `doc.go` per package (charter + trace text only) | empty packages with doc.go compile | `chore(skeleton): create module tree per RFC-003` |
| T0.2 | `scripts/check-imports.sh` — the forbidden-import lint from `PROJECT_MODULE_SPECIFICATION.md` tables | script runs against the empty tree | `build(ci): add dependency-rule import lint` |
| T0.3 | Makefile: real `test`/`build`/`vet`/`importlint` targets replacing the Sprint-0 no-op; wire into `ci` | `make ci` green end-to-end | `build(make): replace stack-agnostic placeholders with Go targets` |
| T0.4 | `.github/workflows/go-ci.yml` (build, vet, test, importlint); update `tests/README.md` to the scaffold map + skip policy | CI green on push | `build(ci): add Go CI workflow` |

**Milestone M0 — the empty architecture enforces itself:** all packages exist, all dependency rules are CI-enforced, `make ci` green. *(Working state: yes — a compiling, linted, rule-enforcing skeleton.)*

### Phase 1 — Track A: `internal/record` (M1)

| ID | Task | Green because | Commit |
|---|---|---|---|
| TA1.1 | Types: `Record` (logical fields), `InstanceID` (opaque, equality only), `Assertions`, `TrustMaterial` vocabulary + `record_test.go` opacity/shape tests | types + tests only | `feat(record): record types and opaque instance identity` |
| TA1.2 | Envelope (private, per assumption A3) + creation surface (restricted to issuance) + round-trip test | create→parse round-trip green | `feat(record): integrity envelope and restricted construction` |
| TA1.3 | `ValidateIntegrity` (`Intact`/`Altered`) + basic altered-detection tests | AT5 seed cases green | `feat(record): integrity validation` |
| TA1.4 | Mutation corpus (`tests/fixtures/mutations/` + `mutation_test.go`): flips, substitutions, truncation, reorder — detection fraction = 1 | AT20 logic green | `test(record): mutation corpus for tamper-evidence` |
| TA1.5 | `Read` + `reconstruct_test.go`: record-alone sufficiency, no verifier state | AT19/AT21 logic green | `feat(record): reconstruction reading` |

**Milestone M1 — the stable surface is real:** a record can be created (via the restricted surface), validated, mutated-and-caught, and reconstructed by a third party. *(Working state: yes — demonstrable tamper-evidence.)*

### Phase 2 — Track B: stores and register (independent of Phase 1 except TA1.1's vocabulary)

| ID | Task | Green because | Commit |
|---|---|---|---|
| TB2.1 | `revstatus`: answer types + honest-indeterminate doc contract | types compile | `feat(revstatus): closed answer set with mandatory as-of` |
| TB2.2 | `revstatus/contracttest`: the realization-independent suite | suite compiles; no realization yet wired | `test(revstatus): realization-independent contract suite` |
| TB2.3 | `revstatus`: degenerate always-Indeterminate realization + pass the suite | contract suite green | `feat(revstatus): degenerate fail-closed realization` |
| TB2.4 | `revorigin`: append-only register + tests (idempotent terminal re-revoke, ordering, view stability) | AT11-origin logic green | `feat(revorigin): append-only revocation register` |
| TB2.5 | `truststore`: hold/answer/withdraw + provisioning record + refusal of malformed material | hit/miss/withdrawn tests green | `feat(truststore): local trust material store` |

**Milestone M2 — every port has an honest provider:** trust material answerable, revocation status answerable (honestly indeterminate), revocations recordable. *(Working state: yes.)*

### Phase 3 — Track A: `internal/verify` (M3) — the sprint's center of gravity

| ID | Task | Green because | Commit |
|---|---|---|---|
| TA3.1 | `Policy` (refuse-unset construction), `Verdict`, `Cause`, `DecisionTrace` types + trace-shape tests | types + tests | `feat(verify): policy, verdict, cause, and trace contracts` |
| TA3.2 | Ports (`TrustMaterialPort`, `RevocationStatusPort`, `TimePort`) + `tests/harness/fakes.go` + `clock.go` controllable time | fakes compile, usable by later tests | `feat(verify): consumer-defined ports and test fakes` |
| TA3.3 | Check 1 — identity binding + tests | per-check green | `feat(verify): identity binding check` |
| TA3.4 | Check 2 — integrity via record + absent/unverifiable inconclusive causes + tests | FM9 paths green | `feat(verify): integrity check with fail-closed material handling` |
| TA3.5 | Check 3 — expiry ± skew + `expiry_skew_test.go` (within/at/beyond tolerance) | AT8 logic green | `feat(verify): expiry check with explicit clock tolerance` |
| TA3.6 | Check 4 — scope integrity + tests | green | `feat(verify): scope integrity check` |
| TA3.7 | Check 5 — revocation under freshness policy + `revocation_policy_test.go` parametric over R and the S4 ceiling | AT13/AT14 *logic* green at arbitrary parameters | `feat(verify): revocation check under parametric freshness policy` |
| TA3.8 | Pipeline assembly + verdict routing (incl. `InconclusiveRejected` `[HYPOTHESIS]`) + unconditional trace + baseline accept/reject tests | AT2/3/6/7/9/11 logic green | `feat(verify): pipeline, verdict routing, unconditional trace` |
| TA3.9 | `rollback_test.go` — every check forced to fail while all others pass; verdict flips each time | AT23 logic green | `test(verify): single-check rollback` |
| TA3.10 | `inconclusive_test.go` — every inconclusive cause → `InconclusiveRejected`, `[HYPOTHESIS]`-marked | AT22 logic green | `test(verify): fail-closed inconclusive routing [HYPOTHESIS]` |

**Milestone M3 — a conformant verifier exists:** every frozen check, every verdict, single-check rollback, and fail-closed routing demonstrable in-process with fakes. Wired with the degenerate `revstatus`, the system already behaves honestly: it rejects rather than claims revocation knowledge. *(Working state: yes — the sprint's headline demo.)*

### Phase 4 — issuance and composition (needs M1 + M3)

| ID | Task | Green because | Commit |
|---|---|---|---|
| TA4.1 | `issuance`: refusal causes, `PermissionSource` port, `IssuanceTrace` + stub-port tests | types + tests | `feat(issuance): ports, refusals, trace` |
| TA4.2 | `Issue` — strict-subset guard, fresh `InstanceID` (assumption A4), nothing-created-on-refusal + AT4/AT18-logic tests | green | `feat(issuance): issuance with strict-subset refusal` |
| TA4.3 | `cmd/atlas-issue`, `cmd/atlas-verify`, `cmd/atlas-revoke` — wiring only; policy loading; AT26 measurement point in `atlas-verify` | drivers build; smoke-run in CI | `feat(cmd): boundary drivers` |
| TA4.4 | `tests/acceptance/` in-process subset — issue → present → verify → revoke → verify again, through public surfaces/drivers with fakes and real stores; substrate-blocked ATs skip with named blockers | in-process AT subset green | `test(acceptance): in-process acceptance subset` |

**Milestone M4 — end-to-end in one process:** issue a delegation, verify it (Accept), tamper (Reject), expire via controlled clock (Reject), revoke at the register + a fake observable status (Reject), all with decision traces — the full RFC-003 §12 flow minus the real substrate. *(Working state: yes — Sprint 1 exit.)*

### Phase 5+ — beyond Sprint 1 (ordered; listed for continuity, not scheduled here)

| Order | Block | Unblocked by |
|---|---|---|
| 5 | Two-domain substrate (shared with EXP-001: SPIRE ×2, isolation, sniffer, latency probe) + substrate ATs (AT10, AT12, AT15–AT17) | founder S1–S4 scope acts (F3) — the substrate serves the spike first (lab discipline: spike before product use) |
| 6 | EXP-001 execution per `lab/EXP-001-EXECUTION-PLAN.md` | substrate + scope acts |
| 7 | Revocation mechanism decision + real `revstatus` realization (passes `contracttest`) + AT13/AT14 end-to-end | spike outcome |
| 8 | AT26 measurement runs; AT30 independent-review package assembly; V1 verdict per `V1_SCOPE.md` | all prior |

## 3. Why this order

- **Stable before volatile (AP12):** record → verifier → stores land before any revocation mechanism exists; the degenerate `revstatus` realization means the volatile region's absence is *represented*, not worked around.
- **The hardest correctness surface gets the longest soak:** `verify` lands mid-sprint with per-check commits, so rollback and inconclusive behavior are exercised by every subsequent task.
- **Nothing waits on the founder:** every Sprint 1 task is parametric over S1–S4 (policy injected) and independent of the spike. The founder's queue (F1–F3) gates Phase 5+, not Sprint 1.
- **Parallel-safe:** Tracks A/B/C touch disjoint packages under CI-enforced dependency rules; merge conflicts are structurally confined to `go.mod` and the Makefile.
- **Each milestone demos:** M0 a self-enforcing skeleton; M1 tamper-evidence; M3 a conformant verifier failing closed; M4 the full flow in-process. No milestone is "plumbing done, nothing visible."

## 4. AT unblock map (every skip names its milestone)

| AT | Runnable at | Blocker until then |
|---|---|---|
| AT1–AT9, AT11, AT18–AT24, AT27–AT29 | **M4 (Sprint 1)** — in-process | — |
| AT10, AT12, AT15–AT17 | Phase 5 | substrate |
| AT13, AT14 | Phase 7 | S1 scope act (R value) + spike-selected realization |
| AT25 | deferred | D1 (frozen deferral — not scheduled) |
| AT26 | Phase 8 | substrate + measurement runs |
| AT30 | Phase 8 | complete build + traces + review package |

<!-- checkpoint: rfc(system-boundary-definition): clarify system boundary definition -->

<!-- checkpoint: governance(API-path-design): extend API path design (#9) -->
