# Implementation Master Plan — Atlas

**Status:** The executable plan for the entire project, from empty skeleton to V1 verdict. Supersedes `IMPLEMENTATION_ORDER.md` and absorbs `ENGINEERING_SPRINT_1_PLAN.md`'s scope (both remain as history). An engineer can execute this sequentially; parallel tracks are marked.
**Discipline:** one task = one conventional commit; **every task leaves `make ci` green**; every milestone leaves the repository in a working, demonstrable state; substrate- and scope-act-blocked tests skip with named blockers, never TODOs.
**Complexity scale:** S ≈ ≤half a day · M ≈ 1–2 days · L ≈ 3+ days (single engineer). Honesty note: these are planning estimates, not commitments; the only externally priced block is EXP-001 (54.5–62.5 eng-h, from its frozen-adjacent plan).

## 0. Epic map and critical path

```
E1 Foundation ─► E2 Record ─► E3 Verifier ─► E5 Issuance+Drivers+AT(in-process) ─► [SPRINT 1 EXIT]
                     │             ▲                    ▲
                     └── E4 Stores/Register ────────────┘        (E2/E3/E4 partially parallel)
FOUNDER GATE F3 (S1–S4 scope acts) ─► E6 Substrate ─► E7 EXP-001 + revocation mechanism ─► E8 V1 closure
```

Critical path: **E1 → E2 → E3 → E5 → (F3) → E6 → E7 → E8.** E4 is off-path (parallel with E2/E3). Founder gates F1/F2 (envelope + instance-ID confirmations) land at the Sprint 1 review; F3 (S1–S4) gates E6, nothing earlier.

Parallel tracks (disjoint packages, CI-enforced isolation): **Track A** = E2→E3→E5-issuance · **Track B** = E4 · **Track C** = E1 tail + harness/fixtures + E5-acceptance scaffolding.

---

## Epic E1 — Foundation (the architecture enforces itself)

**Milestone E1.M — `make ci` green on a compiling, rule-enforcing skeleton.**

| Task | Description | Deps | Cx | ∥ | Blocking | Test required | Definition of done |
|---|---|---|---|---|---|---|---|
| E1-T1 | Create tree per `REPOSITORY_SKELETON.md` + AD-014/AD-015 additions (`verify` shares its time-port contract with `issuance`; `record/revbinding` element file): `go.mod`, all packages, `doc.go` charters | — | S | — | yes (all) | build | `go build ./...` + `go vet ./...` pass |
| E1-T2 | `scripts/check-imports.sh` — R1–R7 + module forbidden-import tables as a build gate | E1-T1 | S | C | yes | lint self-test (a deliberate violation fails) | violation → non-zero exit, wired locally |
| E1-T3 | Makefile: replace Sprint-0 `test` no-op with build/vet/test/importlint; update `tests/README.md` (scaffold map, skip policy, AT index) | E1-T1 | S | C | yes | `make ci` | all gates green incl. existing frozen/secrets checks |
| E1-T4 | `.github/workflows/go-ci.yml` | E1-T3 | S | C | no | CI run | green on push |

**Working state:** an empty architecture that already rejects rule violations.

## Epic E2 — Record Model (M1, the stable surface)

**Milestone E2.M — tamper-evidence demonstrable: create, validate, mutate-and-catch, reconstruct.**

| Task | Description | Deps | Cx | ∥ | Blocking | Test required | Definition of done |
|---|---|---|---|---|---|---|---|
| E2-T1 | Types: `Record` (all elements incl. opaque `InstanceID` [AD-013] and opaque optional `RevBinding` [AD-015]), `Assertions`, trust-material vocabulary | E1 | S | A | yes (E2) | opacity/shape tests | elements carried; nothing interprets the opaque pair |
| E2-T2 | Envelope (private, JWS per AD-012) + restricted creation surface (issuance-only) + round-trip | E2-T1 | M | A | yes | round-trip | create→parse→read yields issuance-time assertions |
| E2-T3 | `ValidateIntegrity` with **algorithm pinning** (SR-2) | E2-T2 | M | A | yes | altered-detection + wrong-alg negative tests | Intact/Altered per §1 invariants; downgrade attempts → Altered |
| E2-T4 | Mutation corpus (`tests/fixtures/mutations/`) + detection-fraction test | E2-T3 | M | A | no | AT20-class: fraction = 1 | every mutation class detected and rejected |
| E2-T5 | `Read` + reconstruction tests (record-alone, no verifier state) | E2-T3 | S | A | yes | AT19/AT21-class | third-party composition recovers who/whom/scope/time |

## Epic E3 — Verification Core (M3, the conformance definition)

**Milestone E3.M — a conformant verifier exists; wired with the degenerate M5 it fails closed rather than pretending revocation knowledge.**

| Task | Description | Deps | Cx | ∥ | Blocking | Test required | Definition of done |
|---|---|---|---|---|---|---|---|
| E3-T1 | `Policy` (refuse-unset), `Verdict`/`Cause` closed sets, `DecisionTrace` | E1 | S | A | yes (E3) | shape + refuse-unset tests | unset policy cannot construct |
| E3-T2 | Ports (`TrustMaterialPort`, `RevocationStatusPort`, `TimePort` — shared contract with M2 per AD-014) + harness fakes + controllable clock | E3-T1 | S | A/C | yes | fakes usable | fakes drive all port answers incl. every honest-negative member |
| E3-T3 | Check 1 — identity binding | E3-T2, E2-T5 | S | A | no | pass/fail unit | per-check trace entry emitted |
| E3-T4 | Check 2 — integrity via M1; `TrustMaterialAbsent`/`SignatureUnverifiable` inconclusive routes | E3-T3 | M | A | no | FM9-path tests | absent material never fetches, routes inconclusive |
| E3-T5 | Check 3 — expiry ± skew | E3-T2 | M | A | no | AT8-logic: within/at/beyond tolerance | deterministic within t; inconclusive beyond |
| E3-T6 | Check 4 — scope integrity (not subset re-derivation) | E3-T4 | S | A | no | tamper-detection unit | tampered scope → definitive reject cause |
| E3-T7 | Check 5 — revocation under freshness policy (R bound + S4 ceiling), parametric | E3-T2 | M | A | no | AT13/AT14-**logic** at arbitrary parameters | stale/indeterminate → inconclusive; observable-revoked → definitive |
| E3-T8 | Pipeline + verdict routing (incl. `InconclusiveRejected` `[HYP]`) + unconditional trace | E3-T3…T7 | M | A | yes | baseline accept/reject (AT2/3/6/7/9/11 logic) | all five checks always traced; routing deterministic |
| E3-T9 | Single-check rollback suite | E3-T8 | M | A | no | **AT23 logic**: each check forced to fail alone → verdict flips | fraction = 1 over checks |
| E3-T10 | Inconclusive-routing suite, `[HYPOTHESIS]`-marked | E3-T8 | S | A | no | **AT22 logic**: every inconclusive cause → InconclusiveRejected | phrased per AT-plan honesty discipline |

## Epic E4 — Stores and Register (M4/M5/M6) — parallel track B

**Milestone E4.M — every port has an honest provider.**

| Task | Description | Deps | Cx | ∥ | Blocking | Test required | Definition of done |
|---|---|---|---|---|---|---|---|
| E4-T1 | `revstatus` answer set (mandatory as-of; honest-indeterminate) | E2-T1 | S | B | yes (E4) | shape tests | closed set exact per §5 |
| E4-T2 | `contracttest` realization-independent suite | E4-T1 | M | B | yes | suite self-checks | asserts all §5 invariants |
| E4-T3 | Degenerate always-Indeterminate realization | E4-T2 | S | B | no | passes `contracttest` | wired as pre-spike default in E5 |
| E4-T4 | `revorigin` append-only register | E2-T1 | S | B | no | append-only, re-revoke no-op, ordering, inert-unknown-ID | §6 invariants demonstrated |
| E4-T5 | `truststore` hold/answer/withdraw + provisioning records + malformed refusal | E2-T1 | M | B | no | hit/miss/withdrawn + refusal | never-fetch structural (importlint) |
| E4-T6 | Store concurrency-posture tests (single-writer/multi-reader per AD-017) | E4-T3/T4/T5 | S | B | no | serialized-mutation + concurrent-read tests | declared posture demonstrated, not assumed |

## Epic E5 — Issuance, drivers, in-process acceptance — **Sprint 1 exit**

**Milestone E5.M — full flow in one process: issue → verify (Accept) → tamper (Reject) → expire (Reject) → revoke+observable-fake (Reject) → reconstruct — all traced.**

| Task | Description | Deps | Cx | ∥ | Blocking | Test required | Definition of done |
|---|---|---|---|---|---|---|---|
| E5-T1 | `issuance` ports (`PermissionSource`, time, `RevBindingSource` [AD-015, answers `none` pre-spike]), refusal causes, trace, injectable minter [AD-013] | E2 | S | A | yes (E5) | shape tests | ports stubbed in harness |
| E5-T2 | `Issue` — proper-subset guard, nothing-created-on-refusal, ephemeral support | E5-T1 | M | A | yes | AT4/AT18 logic | §2 invariants demonstrated |
| E5-T3 | Drivers `atlas-issue`/`atlas-verify`/`atlas-revoke` — wiring, policy loading, trace persistence (AD-016), AT26 measurement point | E3, E4, E5-T2 | M | A | yes | driver smoke-run in CI | zero logic in drivers (review gate) |
| E5-T4 | In-process acceptance subset (AT1–AT9, AT11, AT18–AT24, AT27–AT29 loci); blocked ATs skip with named blockers | E5-T3 | L | A/C | yes | the subset green | AT→file index in `tests/README.md` accurate |
| E5-T5 | Runbooks: permission-source stub operation + trust-material provisioning procedure (AD-D05/D06, minimal V1 form) | E5-T3 | S | C | no | — | an operator can reproduce E5.M from the runbook |
| E5-T6 | **Sprint 1 review with founder:** confirm F1 (AD-012 envelope) + F2 (AD-013 instance-ID); present F3 ask (S1–S4) | E5.M | — | — | gates E6 | — | journal entry if any decision taken |

## FOUNDER GATE F3 — S1–S4 scope acts (no engineering; blocks E6/E7 only)

One journal entry: R value (S1), cached-pull admissibility (S2), broker definition (S3), partition reading (S4) — per `lab/EXP-001-EXECUTION-PLAN.md` §1's exit criteria, including the pre-registration addendum and Standards-Editor sighting.

## Epic E6 — Two-domain substrate (shared: EXP-001 + acceptance; AD-018)

**Milestone E6.M — pre-flight Go: two independent SPIRE domains, federation disabled, link-level partition control, egress observation, NTP-disciplined clocks, evidence tooling — per EXP-001 plan Phases 1–8.**

| Task | Description | Deps | Cx | ∥ | Blocking | Test required | Definition of done |
|---|---|---|---|---|---|---|---|
| E6-T1 | Execute EXP-001 plan Phases 1–6 (groundwork, environment, stock dependencies, SPIRE ×2, isolation, instrumentation) | F3 | L | — | yes | phase exit criteria (frozen plan) | each phase's evidence logged per lab discipline |
| E6-T2 | Realize `tests/harness` substrate interfaces against E6-T1 (domain control, partition, egress, skew) | E6-T1 | M | C | yes | harness self-tests | in-process fakes and substrate implement the same interfaces |
| E6-T3 | Substrate acceptance runs: AT10, AT12, AT15, AT16, AT17 (+ AT8 end-to-end) | E6-T2, E5.M | M | — | no | those ATs | two-run reproducibility per lab discipline |
| E6-T4 | EXP-001 plan Phases 7–8 (roles/blinding, Go/No-Go) | E6-T1 | S | — | yes (E7) | checklist P-gates | GO.json signed |

## Epic E7 — EXP-001 spike + revocation mechanism (the volatile region lands)

**Milestone E7.M — revocation observability realized (or its impossibility-class honestly recorded), AT13/AT14 executed end-to-end.**

| Task | Description | Deps | Cx | ∥ | Blocking | Test required | Definition of done |
|---|---|---|---|---|---|---|---|
| E7-T1 | EXP-001 Phases 9–12: composition attempts cheapest-first, evidence, reproduction, decision + journal | E6-T4 | L (18–26h + 7h repro, priced by the frozen plan) | — | yes | pre-registered spike criteria | outcome class α/β/γ/δ journaled |
| E7-T2 | **Founder acceptance** of outcome → AD-D02 resolution (mechanism, RevBinding production rule, M5 realization owner) | E7-T1 | — | — | yes | — | AD entry appended |
| E7-T3 | Real M5 realization (α path): passes `contracttest`; maintenance path per composition; file-scoped network grant if warranted (R6 honored) | E7-T2 | L | B | yes | `contracttest` + realization-specific | wired via the same seam as the degenerate |
| E7-T4 | `RevBindingSource` production realization (α path; per AD-015) | E7-T2 | M | A | no | issuance round-trip with populated binding | binding remains opaque outside M5 |
| E7-T5 | Propagation channel per S2/S3 + spike selection (M6 View → M5 view) | E7-T3 | M | B | no | freshness-advance tests | as-of advances ≤ R non-partitioned |
| E7-T6 | AT13 + AT14 end-to-end at the pre-registered R and S4 reading | E7-T3…T5, E6 | M | — | no | AT13/AT14 | outcomes recorded even if negative (valid V1 result per AT plan DoD) |
| *β/γ/δ path* | If no composition lands: degenerate realization **is** V1's honest revocation answer; AT13/AT14 record the gap; Level 2 disposition returns to the founder per the gate's outcome table | E7-T2 | — | — | — | AT13/14 document behavior | the gap surfaced, never patched over |

## Epic E8 — V1 closure

**Milestone E8.M — V1 verdict documented, honestly, either way.**

| Task | Description | Deps | Cx | ∥ | Blocking | Test required | Definition of done |
|---|---|---|---|---|---|---|---|
| E8-T1 | AT26: latency measured at the driver boundary across the acceptance run; number reported, no threshold asserted | E6.M, E5.M | S | — | no | AT26 | measured value in the V1 report |
| E8-T2 | Full AT sweep at pre-registered parameters; two-run reproducibility | E7.M | M | — | yes | AT1–AT24, AT26–AT29 | every result maps to a named requirement (AT fail-routes) |
| E8-T3 | AT30 package: spec set + build + traces + runbooks for an independent reviewer; reviewer executes | E8-T2 | M | — | yes | AT30 | every SO/INV verdict reproduced by someone who built none of it |
| E8-T4 | V1 verdict vs `TECHNICAL_VALIDATION.md` P5 items 8–9 criteria; journal entry; hypothesis dispositions (AD-D07) to founder | E8-T3 | S | — | yes | — | success **or** failure documented per `V1_SCOPE.md` DoD |

---

## AT unblock map (single source of truth for skips)

| AT | Runnable at | Named blocker until then |
|---|---|---|
| AT1–AT9, AT11, AT18–AT24, AT27–AT29 | E5.M (in-process) | — |
| AT8 end-to-end, AT10, AT12, AT15–AT17 | E6.M | `substrate` |
| AT13, AT14 | E7.M | `S1-scope-act` + `spike-outcome` |
| AT25 | — (deferred, D1) | frozen deferral — not scheduled |
| AT26 | E8 | `substrate` + measurement run |
| AT30 | E8 | complete build + traces + package |
| (ER14/FR10) | — (deferred, D2) | no AT exists in V1 by design |

## Expected output at plan completion

A validated two-domain reference implementation (or a rigorous negative result naming the failed requirement — equally valid per `V1_SCOPE.md`), with: all six modules realized; the degenerate-or-real revocation provider behind an unchanged contract; every AT executed or skip-accounted; decision traces sufficient for independent reproduction; the V1 verdict journaled; and zero unrecorded deviations from the frozen package.

<!-- checkpoint: context(problem-scoring-guidelines): extend problem scoring guidelines -->

<!-- checkpoint: governance(threat-model-scenarios): extend threat model scenarios -->

<!-- checkpoint: repo(founder-profile-feedback): extend founder profile feedback -->

<!-- checkpoint: feat(internal): implement boundary check -->

<!-- checkpoint: chore(record): tweak CLI flag configuration -->
