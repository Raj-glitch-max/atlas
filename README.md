# Atlas

**Hand one agent's authority to another without handing over a secret — and
verify it in ~116µs with the issuer offline.**

Atlas issues a **scoped, expiring, revocable capability** instead of a shared
key. The receiving side verifies it from locally-held trust material — **no call
to the issuer, even when the network to it is gone** — and a revocation is
enforced independently, failing closed when freshness can't be proven.

## See it work

```console
$ bash examples/unforgettable.sh

1 · an agent gets a scoped capability
   issued  spiffe://domain-a.test/workload/payments-api → spiffe://domain-b.test/agent/booking-worker
   scope   read:orders, write:audit

3 · THE SERVER DIES — no issuer, no network
   server: UNREACHABLE (connection refused)

4 · verify OFFLINE — still answers, in microseconds
   ACCEPT  (offline · 319µs · snapshot age 38ms / budget 10m0s)
     identity_binding   Pass  ·  integrity  Pass  ·  expiry  Pass
     scope_integrity    Pass  ·  revocation_status  Pass

6 · the revoked capability — rejected, OFFLINE
   REJECT  (offline · 240µs)  RevokedObservable

7 · staleness is YOUR policy — an old snapshot fails CLOSED, never open
   INCONCLUSIVE  (offline · 197µs · age 1.017s / budget 500ms)  RevocationKnowledgeStale

8 · a tampered bundle (revocation stripped) is REFUSED outright
   offline verify: REFUSING bundle — revstatus: revoked-set signature does not verify
```

That is real output from the real engine — no mocks, no scripted text. The
script is self-contained: it builds, runs its own server on a temp store, kills
it mid-demo, and cleans up. **One command, ~60 seconds:**

```sh
bash examples/unforgettable.sh
```

## Quickstart

```sh
make ci                                   # build + vet + tests + lints + frozen-doc + import rules
go run ./cmd/atlas-server                 # → http://127.0.0.1:8087
go run ./cmd/atlas doctor                 # confirm it's up
```

## What Atlas is not

A **narrow primitive**: single-hop, across two trust domains, a companion to
SPIFFE and never a replacement. Closest in design space to Biscuit / UCAN. It is
**not** a policy engine, **not** an identity provider, and **not** a replacement
for OAuth or SPIFFE.

If you have one agent talking to one service, use an API key. If your
permissions are static and long-lived, use your IAM. Atlas earns its complexity
only when delegation actually exists — see
[`LIMITATIONS.md`](LIMITATIONS.md) for the full, deliberately unflattering list.

> **Status: v0.1-dev**, a reference implementation. Every claim here is either
> tested or labelled deferred. No third-party security audit has been done.

## What's built

- **The engine** (`internal/`) — the six RFC-003 modules: record (M1), issuance
  (M2), verification (M3, *the conformance definition*), truststore (M4),
  revocation status (M5), revocation origin (M6). Five ordered checks, an
  unconditional decision trace, fail-closed on stale or indeterminate knowledge.
- **Server** (`cmd/atlas-server`) — HTTP JSON API over the real engine
  (`/issue /verify /revoke /delegations /audit /graph /stats /bundle /metrics`,
  plus `/health` and `/readyz`), durable file store, optional bearer auth,
  pinned CORS, per-IP rate limiting, TLS, Prometheus metrics, access logs.
- **CLI** (`cmd/atlas`) — `delegate · verify · revoke · delegations · graph ·
  audit · doctor · version · inspect · bundle`, including **offline**
  verification (`verify --offline --bundle …`) and the `--require-scope`
  authorization gate.
- **MCP server** (`cmd/atlas-mcp`) — Atlas as agent tools over MCP stdio.
- **SDKs** — zero-dependency clients in [Go](sdk/go), [Python](sdk/python), and
  [TypeScript](sdk/typescript), each mirroring the same API.
- **Reference gate** (`examples/atlas-gate`) — a deployable reverse proxy that
  admits a request only if it carries a capability granting the required scope,
  verified offline.
- **Operator console + site** (`ui/`) — see [`ui/README.md`](ui/README.md).

Deploy with the hardened container: `deploy/` (distroless nonroot, read-only
rootfs, `docker compose`).

## Can you verify the claims?

That's the point — every number below is something you can re-derive.

| Claim | Check it yourself |
|---|---|
| The engine passes its conformance suite | `go test ./tests/...` — 30 vectors, 20 adversarial (`alg:none`, HS256 confusion, signature/payload transplant, duplicate JSON keys) |
| The verifier never accepts garbage | `go test ./internal/verify -run=XXX -fuzz=FuzzVerify -fuzztime=60s` — runs in CI on every PR; ~1.9M executions, zero silent acceptances |
| Verify is ~116µs (p50) | `bash scripts/run-benchmarks.sh` — prints the CPU, Go version and iteration counts *alongside* the numbers ([`docs/BENCHMARKS.md`](docs/BENCHMARKS.md)) |
| Module boundaries hold | `make importlint` — dependency rules R1–R7, 21 packages, 0 violations |
| It works end to end | `bash examples/unforgettable.sh` |

Measured on a Ryzen 5 5600H, Go 1.22.11, chain depth 1 — in-process engine, not
end-to-end server latency (that's ~2-3x, and exported live as the
`atlas_verify_latency_seconds` Prometheus histogram).

## Runnable examples

| Script | Shows |
|---|---|
| [`examples/unforgettable.sh`](examples/unforgettable.sh) | offline verify · revoke · staleness fail-closed · tamper-refusal |
| [`examples/ship-a-landing-page.sh`](examples/ship-a-landing-page.sh) | multi-tool agent workflow: least privilege, blast-radius containment |
| [`examples/agent-capability-demo.sh`](examples/agent-capability-demo.sh) | single-hop grant, attenuation, offline verify, revocation |

## Documentation

**New here?** [`docs/guides/START_HERE.md`](docs/guides/START_HERE.md) is a
zero-prior-knowledge reading path from "what problem does this solve" to reading
the verification core, in order.

| | |
|---|---|
| Why it exists, and why not OAuth/JWT/Biscuit | [`docs/product/WHY.md`](docs/product/WHY.md) |
| What it deliberately does **not** do | [`LIMITATIONS.md`](LIMITATIONS.md) |
| What it defends, and the test proving each claim | [`docs/architecture/THREAT_MODEL.md`](docs/architecture/THREAT_MODEL.md) |
| The engine's design | [`docs/architecture/SYSTEM_ARCHITECTURE.md`](docs/architecture/SYSTEM_ARCHITECTURE.md) · [`rfc/`](rfc/) |
| The wire format, for other implementations | [`tests/vectors/VECTORS.md`](tests/vectors/VECTORS.md) |
| The hard questions, pre-answered | [`docs/product/OBJECTIONS.md`](docs/product/OBJECTIONS.md) |
| Where it's going | [`ROADMAP.md`](ROADMAP.md) |
| Performance, with methodology | [`docs/BENCHMARKS.md`](docs/BENCHMARKS.md) |

## Contributing

Read [`CONTRIBUTING.md`](CONTRIBUTING.md) — it's the single process document
(setup, layout, branches, commits, review, and the frozen-planning rule).
Governance is in [`.github/GOVERNANCE.md`](.github/GOVERNANCE.md), maintainers in
[`.github/MAINTAINERS.md`](.github/MAINTAINERS.md), conduct in
[`.github/CODE_OF_CONDUCT.md`](.github/CODE_OF_CONDUCT.md).

Some planning documents are hash-pinned and verified by `make check-frozen`.
Editing one without the amendment process breaks the build, by design —
`CONTRIBUTING.md` §6.

## Security

Found a vulnerability? **Do not open a public issue.** Follow the private
disclosure process in [`SECURITY.md`](SECURITY.md).

## License

[Apache-2.0](LICENSE). See [`NOTICE`](NOTICE) for attribution and third-party
dependency licenses.
