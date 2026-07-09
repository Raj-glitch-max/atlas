# Atlas

**Offline-verifiable, attenuable delegation for SPIFFE workload identity.**

When one workload — or one AI agent — needs another to act with its authority,
Atlas lets you hand off a **scoped, expiring, revocable capability** instead of a
shared secret. The receiving side verifies it in microseconds from locally-held
trust material — **no call to the issuer**, even when the network to it is gone —
and a revocation is enforced independently, failing closed when freshness can't
be proven.

Atlas is a **narrow primitive**: single-hop, across two trust domains, a
companion to SPIFFE (never a replacement). It is closest in design space to
Biscuit / UCAN. It is **not** a policy engine, an identity provider, or a
replacement for OAuth/SPIFFE. See [`WHY.md`](WHY.md) and
[`LIMITATIONS.md`](LIMITATIONS.md) for the honest boundaries.

> Status: **v0.1-dev**, a reference implementation. Every claim in this repo is
> either tested or explicitly labelled as deferred/hypothesis — see
> [`LIMITATIONS.md`](LIMITATIONS.md).

## Quickstart

```sh
# build & run the full local gate (build + vet + tests + lints + frozen-doc + import rules)
make ci

# run the server, then verify a capability end-to-end
go run ./cmd/atlas-server                       # → http://127.0.0.1:8087
go run ./cmd/atlas doctor                        # confirm it's up

# the headline demo — offline verify, revoke, tamper-refusal, all live:
bash examples/unforgettable.sh
```

## What's built (all tested; run `make ci`)

- **The engine** (`internal/`) — the six RFC-003 modules: record (M1), issuance
  (M2), verification (M3, the conformance definition), truststore (M4),
  revocation status (M5), revocation origin (M6). Five ordered checks, an
  unconditional decision trace, fail-closed on stale/indeterminate knowledge.
- **Server** (`cmd/atlas-server`) — HTTP JSON API over the real engine
  (`/issue /verify /revoke /delegations /audit /graph /stats /bundle /metrics`,
  plus `/health` and `/readyz`), durable file store, optional bearer auth,
  configurable CORS, per-IP rate limiting, TLS, Prometheus metrics, access logs.
- **CLI** (`cmd/atlas`) — `issue · verify · revoke · delegations · graph · audit
  · doctor · version · inspect · bundle`, including **offline** verification
  (`verify --offline --bundle …`) and the `--require-scope` authorization gate.
- **MCP server** (`cmd/atlas-mcp`) — Atlas as agent tools over MCP stdio.
- **SDKs** — zero-dependency clients in [Python](sdk/python), [TypeScript](sdk/typescript),
  and [Go](sdk/go), each mirroring the same API.
- **Reference gate** (`examples/atlas-gate`) — a deployable reverse-proxy that
  admits a request only if it carries a valid capability granting the required
  scope, verifying offline.
- **Operator console + product site** (`ui/`) — a live operator surface and a
  marketing site (Vite + Three.js + GSAP); see [`ui/README.md`](ui/README.md).
- **Assurance** — 28 conformance vectors (18 adversarial) in `tests/vectors`,
  coverage-guided fuzzing, property tests, published microbenchmarks + a latency
  histogram, an import-boundary lint (dependency rules R1–R7), and frozen-doc
  integrity.

Deploy with the hardened container: `deploy/` (distroless nonroot, read-only
rootfs, `docker compose`).

## Runnable examples

| Script | Shows |
|---|---|
| [`examples/unforgettable.sh`](examples/unforgettable.sh) | offline verify · revoke · staleness fail-closed · tamper-refusal |
| [`examples/ship-a-landing-page.sh`](examples/ship-a-landing-page.sh) | multi-tool agent workflow: least privilege, blast-radius containment |
| [`examples/agent-capability-demo.sh`](examples/agent-capability-demo.sh) | single-hop grant, attenuation, offline verify, revocation |

## Documentation map

The planning corpus is extensive and partly **frozen** (hash-pinned; see below).
Start here:

- **Orientation** — [`context/00_PROJECT_CONTEXT.md`](context/00_PROJECT_CONTEXT.md)
  (mission), [`ROADMAP.md`](ROADMAP.md) (what's done / next), [`WHY.md`](WHY.md),
  [`LIMITATIONS.md`](LIMITATIONS.md).
- **Architecture** — [`SYSTEM_ARCHITECTURE.md`](SYSTEM_ARCHITECTURE.md),
  [`rfc/`](rfc/) (RFC-000…003), `MODULE_SPECIFICATION.md`,
  `INTERFACE_SPECIFICATION.md`.
- **Product / engineering specs (frozen)** — [`docs/product/`](docs/product/),
  [`docs/engineering/`](docs/engineering/).
- **Security** — [`SECURITY.md`](SECURITY.md), [`THREAT_MODEL.md`](THREAT_MODEL.md),
  `docs/engineering/02_SECURITY_OBJECTIVES.md`.
- **Governance** — [`CONTRIBUTING.md`](CONTRIBUTING.md),
  [`DEVELOPMENT_RULES.md`](DEVELOPMENT_RULES.md),
  [`context/01_GOVERNANCE.md`](context/01_GOVERNANCE.md), and the
  reasoning framework in [`agents/`](agents/).
- **Substrate lab** — [`lab/`](lab/) and [`atlas-lab/`](atlas-lab/) (the
  two-domain SPIRE experiment environment).
- **Backlog / debt / risk** — [`BACKLOG.md`](BACKLOG.md),
  [`TECHNICAL_DEBT_REGISTER.md`](TECHNICAL_DEBT_REGISTER.md),
  [`RISK_REGISTER.md`](RISK_REGISTER.md).

## The frozen-planning rule

The planning documents in [`scripts/frozen-docs.list`](scripts/frozen-docs.list)
are frozen: `make check-frozen` verifies their SHA-256 hashes against
`FROZEN.sha256` and is wired into CI. Editing one without a dated amendment and a
journal entry breaks the build — by design. See `CONTRIBUTING.md`
§"Frozen planning documents."

## Toolchain

`go` (module `github.com/Raj-glitch-max/atlas`), plus `python3` + `pre-commit`
and `npx` (Node) for the lint gates; optionally `docker` + `gitleaks`. Run
`make help` for all targets.
