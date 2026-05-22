# Roadmap

Honest and staged. Dates are intentionally omitted — this lists *order and
intent*, not promises. "Done" items are tested and in `main`; everything else
is open. The live, granular task list is `BACKLOG.md`.

## Done — the core is real

- **Delegation primitive** — single-hop, two-domain, SPIFFE-bound, attenuated,
  offline-verifiable, independently revocable. Conformance-tested.
- **Verification engine** — five ordered checks, unconditional decision trace,
  fail-closed on stale/indeterminate knowledge.
- **Revocation** — signed revoked-set snapshots with verifiable freshness.
- **Product surfaces** — server (HTTP API + durable store + auth + metrics),
  CLI (`atlas delegate/verify/revoke/bundle/inspect/…`, incl. **offline**
  verify), MCP agent tools, Python SDK, operator console.
- **Assurance** — 28 conformance vectors (18 adversarial), coverage-guided
  fuzzing, property tests, published microbenchmarks + latency histogram,
  import-boundary lint, frozen-doc integrity.
- **The unforgettable demo** — `examples/unforgettable.sh` (offline verify +
  revoke + tamper-refusal, live).

## Next — make it trustworthy to depend on

Ordered by leverage, not size (N1–N7).

1. **Reference service adapters + a capability→token broker** — the bridge from
   "compelling demo" to "in your stack": small gateways (or MCP servers /
   sidecars) that hold a real service credential (GitHub/Vercel/Slack/…),
   verify an Atlas capability offline, and exchange it for the real action or a
   short-lived service token. The pattern is proven end-to-end in
   `examples/ship-a-landing-page.sh` with mock adapters and the
   `atlas verify --require-scope` authorization gate; these are the real ones.
   See `docs/product/AGENT_WORKFLOWS.md`.
2. **Proof-of-possession binding** — bind a record to the delegate's own key
   so an intercepted record isn't a bearer token (closes `LIMITATIONS.md` §3,
   the biggest security gap).
3. **Key rotation ceremony** — overlap windows, `kid` lifecycle, key
   revocation; automate what is manual today (§4).
4. **Normative wire-format RFC** — promote `tests/vectors/VECTORS.md` to a
   versioned spec so independent implementations have a fixed target
   (governance-gated).
5. **TypeScript + Go published SDKs** — mirror the Python SDK; the JS client in
   `ui/src/api.js` is the seed.
6. **Durable DB backend** behind the existing `Store` seam (bolt/sqlite) +
   backup/restore.
7. **Deployment hardening** — TLS, non-root/read-only container, rate limits,
   `atlas doctor` preflight for prod.

## Later — prove the hard claims and scale (L1–L5)

1. **Run the two-domain SPIRE substrate** (`atlas-lab/`) on real infra — the
   link-level zero-egress packet proof and partition/propagation measurements
   (`LIMITATIONS.md` §8).
2. **gRPC API** + streaming for internal service-mesh callers.
3. **Third-party security audit** of the crypto/verification path.
4. **Scaling curves** — verification/throughput at 1M–100M delegations;
   revocation-snapshot size/latency behavior.
5. **Independent implementations** — the conformance suite exists precisely so
   a Rust/Zig verifier can appear and provably agree.

## Explicitly out of scope (by design, not neglect)

Multi-hop re-delegation (A→B→C), 3+-domain federation, being an identity
provider, replacing OAuth/SPIFFE. See `LIMITATIONS.md` §§1–2 and `WHY.md` Q3.

## Compatibility & versioning policy

- The **record wire format** and the **conformance vectors** are the
  compatibility contract. The record carries a pinned `typ`
  (`atlas-record+jws`); the vector files carry a `schemaVersion`.
- Changes that would alter a conformant verifier's decision on an existing
  vector are **breaking** and gate on a version bump + a migration note.
- Additive fields must be ignorable by older verifiers without changing a
  verdict.
- The frozen planning corpus (`scripts/frozen-docs.list`) changes only via the
  amendment process in `CONTRIBUTING.md`.

## How to contribute

Setup, commit conventions, and the frozen-doc rules live in `CONTRIBUTING.md`.
This is the *what to work on* for the product surfaces.

**Where things live**

| Area | Path |
|---|---|
| The engine (frozen; compose, don't fork) | `internal/` |
| Server / CLI / MCP / offline verifier | `cmd/` |
| Conformance vectors (the compat contract) | `tests/vectors/` |
| SDKs | `sdk/` |
| Site + operator console | `ui/` |
| Benchmarks | `atlas-lab/bench/` |
| Deploy | `deploy/` |

**Rules of the road**

- The engine (`internal/`) is conformance-frozen. Product code composes it
  through public APIs only — the import-boundary lint (`make ci`) enforces
  this. A change that alters a verifier's decision must move a conformance
  vector, deliberately.
- Everything ships with a test. `make ci` must stay green.

**Good first contributions** (roughly easiest → deepest)

1. A new CLI output format (`--json`) for `verify`/`delegations`.
2. Add a conformance vector for an edge case you can construct (see
   `tests/conformance/`).
3. Port the Python SDK shape to **TypeScript** (`sdk/js`) — the client already
   exists at `ui/src/api.js`; make it standalone + tested.
4. A `Store` implementation over SQLite behind the existing interface
   (`cmd/atlas-server/store.go`).
5. A Grafana dashboard JSON for the real `/metrics` (histogram + counters).
6. Harden the container (non-root, read-only rootfs) in `deploy/`.

Bigger pieces are the numbered roadmap items above — say hi on an issue before
starting one so we don't collide.

## Release cadence

Pre-1.0: released when a coherent capability lands and `make ci` is green — no
fixed calendar. The wire format and conformance vectors are the stability
promise (see the policy above); tooling around them may move faster. A tagged
`v0.x` marks a state where the demo, the tests, and the docs all agree.

<!-- checkpoint: fix(record): fix truststore backend (#49) -->

<!-- checkpoint: test(internal): test truststore backend (#68) -->

<!-- checkpoint: chore(security): clean Docker orchestration config -->
