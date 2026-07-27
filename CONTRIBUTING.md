# Contributing

Contributions are welcome. This is the only process document you need — branch
strategy, commits, review, and the freeze rule all live here.

## 1. First-time setup

```sh
make init     # installs pre-commit + commit-msg hooks
make help     # list all targets
```

Toolchain: `go` (1.22+, see `go.mod`), `git`, `python3` + `pre-commit`, and
`npx` (Node) for the markdown lint. Optionally `docker` and `gitleaks` — CI
installs its own copy of the latter.

## 2. Day-to-day commands

| Target | Does |
|---|---|
| `make build` | `go build ./...` |
| `make test` | `go test ./...` |
| `make vet` | `go vet ./...` |
| `make importlint` | Enforce dependency rules R1–R7 (see §7) |
| `make lint` | pre-commit hooks on all files |
| `make format` | Auto-fixing hooks only |
| `make docs-lint` | Markdown lint (markdownlint-cli2) |
| `make secrets` | Secret scan (gitleaks) |
| `make check-frozen` | Verify frozen planning docs unchanged |
| `make ci` | Everything above — the local equivalent of CI |
| `make upgrade` | `pre-commit autoupdate` |

Two things CI runs that `make ci` does not:

- **Coverage-guided fuzzing** — `go test ./internal/verify -run=XXX -fuzz=FuzzVerify -fuzztime=60s`.
  Run it locally before touching the verification core.
- **Benchmarks** — `bash scripts/run-benchmarks.sh` regenerates
  [`docs/BENCHMARKS.md`](docs/BENCHMARKS.md) with the machine's details
  attached. Re-run it if you change anything on a hot path, and commit the
  result.

## 3. Repository layout

```text
cmd/          binaries: atlas (CLI), atlas-server, atlas-mcp, atlas-{issue,verify,revoke}
internal/     the engine — the six RFC-003 modules, the conformance definition:
              record (M1) · issuance (M2) · verify (M3) · truststore (M4)
              revstatus (M5) · revorigin (M6)
sdk/          zero-dependency clients: go, python, typescript
examples/     runnable demos against the real engine, plus atlas-gate (reverse proxy)
tests/        conformance kit, language-neutral vectors, acceptance tests, harness
ui/           product site + operator console (Vite)
deploy/       hardened container (distroless, nonroot) + compose
docs/         architecture, product, engineering, planning, research, guides
rfc/          architecture RFCs (RFC-000…003)
agents/       reasoning framework + dated decision journal
lab/          pre-registered SPIRE experiments (governed separately — see §8)
scripts/      CI helpers: import lint, frozen-doc guard, benchmark capture
```

**`internal/` is the load-bearing part.** Everything else composes it through
public APIs only, and the import lint enforces that.

## 4. Branches and commits

- `main` is the integration branch. CI runs on pushes to `main` and on all PRs.
- Keep branches short-lived and scoped to one reviewable change.
- Stage explicit paths. Don't sweep unrelated work into a commit.

Conventional Commits 1.0.0, enforced by the `commit-msg` hook:

```text
feat(server): add readiness probe
fix(verify): reject duplicate JSON keys
docs: correct the benchmark methodology
chore(ci): pin the gitleaks version
```

Scope is optional; `!` marks a breaking change. Only the first line is enforced,
so write a real body — explain *why*, not *what*.

## 5. Pull requests

Use the PR template. A PR is mergeable when:

- `make ci` passes locally,
- new behaviour has a test, and a change to the verification core has a
  conformance vector,
- any journal entry required by §6 is linked.

Changes to `internal/verify` carry a higher bar than the rest of the tree: it is
the conformance definition, and a silent behaviour change there desynchronizes
every implementation that trusts the vectors.

## 6. Frozen planning documents

Some planning documents are **frozen** — their bytes are hash-pinned. The list
is `scripts/frozen-docs.list`; the baseline is `FROZEN.sha256`; the guard is
`make check-frozen`, wired into CI.

Frozen today: the four `docs/research/` briefs,
`docs/planning/P5_FALSIFICATION_EXPERIMENT.md`,
`docs/planning/LEVEL0_1_FEASIBILITY_GATE.md`, the ten `docs/product/` specs
anchored by `PRODUCT_DEFINITION.md`, the five `docs/engineering/` specs, and the
five static `lab/` process docs. (`lab/EXPERIMENT_LOG.md` is append-only by
design and deliberately *not* frozen — the log is meant to grow.)

To change a frozen document:

1. Record a **journal entry** at `agents/journal/<YYYY-MM-DD>-<slug>.md` per
   `agents/GOVERNANCE.md`.
2. Add a dated, reasoned change note to the document itself (lab docs: per
   `lab/LAB_README.md` §12).
3. Re-baseline: `make frozen-baseline`, and commit `FROZEN.sha256` alongside
   the amendment.

A frozen doc changed without these steps fails CI. **Editing the hash baseline
to silence the alarm is itself the violation being guarded against.**

*Relocating* a frozen document is not an amendment: move it, update its path in
`scripts/frozen-docs.list`, re-baseline, and record the move in a journal entry.
The hash must come out identical — if it doesn't, you changed the content. See
`agents/journal/2026-07-27-relocate-frozen-planning-docs.md` for a worked
example.

## 7. Dependency rules (R1–R7)

`scripts/check-imports.sh` enforces the module dependency graph from RFC-003 —
which module may import which. Violations fail the build rather than review.
The script self-tests (`--self-test`) so the checker itself is checked.

## 8. Boundaries

- **`lab/`** runs pre-registered experiments and records evidence. It does not
  implement product code or choose technologies. Experiment logs are
  append-only.
- **`rfc/`** is a record of accepted architecture decisions. Don't add new RFCs
  under current governance; use the journal and the freeze process instead.
- **`agents/journal/`** preserves disagreement instead of smoothing it into
  consensus. Overrides are allowed, but recorded with reasoning and the
  conditions under which they'd be reconsidered.

## 9. Security

Do not open a public issue for a vulnerability. Follow the private disclosure
process in [`SECURITY.md`](SECURITY.md).

## 10. License

Atlas is licensed under the [Apache License 2.0](LICENSE). Contributions are
accepted under the same license. See [`NOTICE`](NOTICE) for attribution and
third-party dependency licenses.
