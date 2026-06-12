# tests/

Acceptance-test scaffold for Atlas. Unit tests live beside their packages
under `internal/`; this tree holds everything that exercises the system at
its product boundaries.

## Layout

| Path | Content | Lands at |
|---|---|---|
| `harness/` | instrumentation: port fakes, controllable clock (E3-T2); substrate-control interfaces — two-domain control, link-level partition, egress observation — realized against the shared EXP-001 substrate (E6-T2) | E3/E6 |
| `acceptance/` | one file per acceptance-test family (AT plan: `docs/engineering/05_ACCEPTANCE_TEST_PLAN.md`); tests exercise product only through drivers and public package surfaces, mirroring each AT's stated test locus | E5-T4 |
| `fixtures/` | record corpus, trust-material bundles for the two synthetic domains, mutation-corpus definitions (AT20) | E2-T4 |

Rules: no product package may import this tree (enforced by
`scripts/check-imports.sh`, SR-4). Harness fakes obey the honest-negative
discipline — a fake never expresses ignorance as `NotObservedRevoked` (FD-4).

## Skip policy (binding)

A test that cannot run yet **skips with a named blocker** — never a TODO,
never a silent omission. Valid blockers and their unblock points
(authoritative map: `IMPLEMENTATION_MASTER_PLAN.md` §AT unblock map):

| Blocker | Meaning | Unblocks at |
|---|---|---|
| `S1-scope-act` | needs the founder's R value / S4 reading (gate F3) | F3 journal entry |
| `substrate` | needs the two-domain SPIRE substrate with partition + egress instrumentation | E6.M |
| `spike-outcome` | needs the EXP-001-selected revocation realization | E7.M |

## AT → file index

| Acceptance tests | File (under `acceptance/`) | First runnable |
|---|---|---|
| AT1–AT2 identity binding | `identity_test.go` | E5.M (in-process) |
| AT3–AT5 scope | `scope_test.go` | E5.M |
| AT6–AT8 expiration/clock | `expiry_test.go` | E5.M (AT8 end-to-end: E6.M) |
| AT9–AT14 revocation | `revocation_test.go` | E5.M (AT10/12: `substrate`; AT13/14: `S1-scope-act` + `spike-outcome`) |
| AT15–AT16 offline/no-egress | `offline_test.go` | `substrate` (E6.M) |
| AT17 two-domain | `crossdomain_test.go` | `substrate` (E6.M) |
| AT18 ephemeral issuance | `ephemeral_test.go` | E5.M |
| AT19–AT21 reconstruction | `reconstruction_test.go` | E5.M |
| AT22–AT23 fail-closed `[HYPOTHESIS]` / rollback | `failclosed_test.go` | E5.M |
| AT24 non-replacement (AT25 deferred, D1) | `interop_test.go` | E5.M |
| AT26 latency (measure-only) `[HYPOTHESIS]` | `latency_test.go` | E8 (`substrate`) |
| AT27–AT29 boundary | `boundary_test.go` | E5.M |
| AT30 independent review | process act, not code — package assembled at E8-T3 | E8 |

## Execution discipline

Substrate-boundary runs follow lab governance (`lab/LAB_README.md`,
`lab/EXPERIMENT_CHECKLIST.md`, `lab/DECISION_RULES.md`): two-run
reproducibility, adversary-blinded roles, pre-registered parameters (the R
and S4 readings each run used are recorded with the run). Hypothesis-marked
tests (`AT8` fail-closed portion, `AT22`, `AT26`) document behavior; they are
not pass-gated acceptance for committed scope and their phrasing must not
promote hypotheses (DR7).
