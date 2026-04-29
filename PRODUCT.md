# Atlas — the product

> **Let your AI agents delegate work to each other — safely.**
> No shared API keys. No live auth server in the hot path. Verification that
> survives the network, and can't be lied to.

When one agent needs another to act with its authority, Atlas hands off a
**scoped, expiring, revocable capability** instead of a secret — and the
receiving side can verify it in microseconds, even offline. Underneath: a
narrow, conformance-tested delegation primitive (the frozen kernel in
`internal/`) wrapped in a real product — a server, a CLI, an agent interface
(MCP), and SDKs, over one shared trust state.

**New here? Start with [`WHY.md`](WHY.md)** — the outcome, the pain it removes,
and why-not-OAuth/JWT/Biscuit. This file is the *how to run it*.

## See the whole thesis in 60 seconds

```bash
bash examples/unforgettable.sh
```

An agent gets a capability → **the server is killed** → verification still
answers in microseconds → reconnect → revoke → the revoked capability is
rejected *offline* → a tampered trust bundle is refused. Real engine, no mocks.

*(The engineering/governance workspace is described in `README.md`; the
primitive's specs live in `docs/` and `rfc/`.)*

## Four ways in, one engine

```
                    ┌────────────────────────────┐
   humans  ──CLI──▶ │                            │
   humans  ──UI───▶ │   atlas-server (:8087)     │──▶  internal/ (frozen engine)
   agents  ──MCP──▶ │   API · store · audit      │      issue · verify · revoke
   apps    ──SDK──▶ │                            │
                    └────────────────────────────┘
```

| Surface | Path | What it is |
|---|---|---|
| **Server** | `cmd/atlas-server` | HTTP API over the real engine: `POST /issue /verify /revoke`, `GET /health /version /delegations /audit /graph /stats /metrics`. Durable store + authority key, audit log, optional bearer-token auth, Prometheus metrics. |
| **CLI** | `cmd/atlas` | `issue · verify · inspect · revoke · delegations · graph · audit · doctor`. Script-friendly exit codes; `issue -q \| verify -` pipes. |
| **MCP** | `cmd/atlas-mcp` | Atlas as **agent tools** over the Model Context Protocol (`atlas_issue/verify/revoke/…`). Register with Claude: `claude mcp add atlas -- ./atlas-mcp`. |
| **SDK** | `sdk/python` | Zero-dependency Python client. (The UI ships a JS client in `ui/src/api.js`.) |
| **UI** | `ui/` | A site with a live verify console, and an **operator console** (`/console.html`) — delegations, audit, graph, metrics, live. |
| **Deploy** | `deploy/` | Multi-stage Docker image (distroless) + compose (server + Prometheus + Grafana). |

## 60-second quickstart

```bash
# 1. run the backend (durable)
go run ./cmd/atlas-server -store ./atlas-state.json -key ./authority.key   # → :8087

# 2. drive it from the CLI
go run ./cmd/atlas doctor
REC=$(go run ./cmd/atlas issue -q \
  --principal spiffe://domain-a.test/workload/payments-api \
  --delegate  spiffe://domain-b.test/agent/booking-worker \
  --scope read:orders,write:audit)
echo "$REC" | go run ./cmd/atlas verify -        # ACCEPT
echo "$REC" | go run ./cmd/atlas inspect -        # decode claims (offline)

# 3. one-command lifecycle demo
bash examples/agent-capability-demo.sh

# 4. the UI (in ui/)
cd ui && npm install && npm run dev              # → localhost:5173  (+ /console.html)
```

## What it does — and deliberately does not

- **Single hop.** A principal → delegate grant, verifiable across two trust
  domains. **Not** multi-hop re-delegation (A→B→C) — that is outside the
  primitive by design.
- **Attenuation.** A delegate is granted a subset of the principal's
  capabilities; over-scope requests are refused by the engine.
- **Offline verification.** The verifier decides from locally-held trust
  material and a signed revocation snapshot — no call to the issuer. Freshness
  is bounded by R (default 2s); past R it fails closed.
- **Independent revocation.** Signed revoked-set snapshots with verifiable
  freshness; a revoked record verifies to `RevokedObservable`.

## Status (honest)

- The **engine** is the frozen, conformance-tested kernel; the product surfaces
  compose it through public APIs only (import-boundary lint: 0 violations).
- Storage is a durable **file snapshot** (atomic write); a real DB is the next
  backend behind the same `Store` seam.
- The **Docker image** is authored + compose-validated but not built in the
  authoring sandbox (no daemon). Everything else here is tested and
  live-verified.
- Deferred: gRPC + more SDKs, non-root container hardening, TLS, the two-domain
  SPIRE substrate run (`atlas-lab/`, needs real infra). See `BACKLOG.md`.

## The document map

| You want… | Read |
|---|---|
| Why it exists, outcome-first | [`WHY.md`](WHY.md) |
| Does it fix my multi-tool agent workflow? | [`docs/product/AGENT_WORKFLOWS.md`](docs/product/AGENT_WORKFLOWS.md) |
| What can go wrong / what it defends | [`THREAT_MODEL.md`](THREAT_MODEL.md) |
| What it deliberately does **not** do | [`LIMITATIONS.md`](LIMITATIONS.md) |
| The hard questions, pre-answered | [`docs/product/OBJECTIONS.md`](docs/product/OBJECTIONS.md) |
| The 30s / 2m / 10m pitch | [`docs/product/PITCH.md`](docs/product/PITCH.md) |
| Where it's going + how to help | [`ROADMAP.md`](ROADMAP.md) |
| The wire format for other implementations | [`tests/vectors/VECTORS.md`](tests/vectors/VECTORS.md) |

<!-- checkpoint: repo(revocation-requirements): extend revocation requirements (#41) -->

<!-- checkpoint: rfc(glossary-definitions): improve glossary definitions -->

<!-- checkpoint: chore(revstatus): tweak attenuation rule engine -->

<!-- checkpoint: chore(sdk): tweak verification controller -->

<!-- checkpoint: refactor(revstatus): refactor revstatus snapshot retrieval -->
