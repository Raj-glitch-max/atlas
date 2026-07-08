# Atlas — working backlog

Running tracker of open/deferred work. Updated 2026-07-06.
Status legend: ⬜ not started · 🟡 in progress · ✅ done (kept briefly for context) · ⛔ blocked (needs something absent)

---

## A. Site / UI  (`ui/` — Vite + Three.js + GSAP)

**Done this thread:** award-winning site — WebGL constellation hero, verify +
revocation-under-partition cinematics, Lenis smooth-scroll, GSAP camera beat,
scroll-progress rail, active-nav, count-up, sparkline draw-in, cursor
spotlight, filmic grain, real Geist fonts, mobile responsive, reduced-motion
fallback, boot failsafe. Runs at `http://localhost:5173/` (`cd ui && npm run dev`).

- ⬜ **Your visual review + tuning** — motion speed / density / camera-beat
  intensity; anything that feels off at desktop **and** narrow widths.
- ⬜ **Commit `ui/` to git** — currently fully untracked (see §C).
- ⬜ **Wire the real Go verifier** — *explicitly deferred; the named next step.*
  Small local bridge (e.g. a tiny Go HTTP endpoint wrapping `internal/verify`)
  so the pipeline/verdict fire on **real** verdicts instead of scripted data.
- ⬜ (optional) Deeper 3D: WebGL bloom pass · force-directed cluster separation ·
  procedural "ice/crystal" nodes (Igloo-style).
- ⬜ (optional) Third cinematic: key rotation / trust-bundle refresh.
- ⬜ (optional) **Full product app** — the original "Mission Control" scope now
  only sketched as marketing sections: Delegations table, JWS/JSON/binary
  inspectors, signature debugger, packet-capture viewer, conformance replay,
  property-test + fuzzing dashboards, audit-log explorer.
- ⬜ (housekeeping) Decide fate of the hosted Claude Design artifact
  (`claude.ai/code/artifact/5c35e41f…`) — superseded by the local Vite site.

---

## E. Agent-workflow story  (the "software for agents" thesis, sharpened)

- ✅ **`--require-scope`** — `atlas verify` is now an authorization gate
  (valid AND grants the action), not just authentication. The primitive a tool
  adapter uses.
- ✅ **`-grant` flag** — the server's delegatable permission set is configurable
  (real tool scopes like `github:push:acme/landing`, not a hardcoded demo set).
- ✅ **`examples/ship-a-landing-page.sh`** — the multi-tool workflow made safe:
  scoped capabilities → offline per-tool authorization → blast-radius
  containment → instant revocation. Real engine, mock service adapters.
- ✅ **`docs/product/AGENT_WORKFLOWS.md`** — honest positioning (Atlas is the
  trust fabric, not the orchestrator) + every senior-engineer question answered.
- ✅ **UI `#workflows` section** — the landing-page scenario, outcome-first.
- ⬜ **Real service adapters + capability→token broker** — the bridge to a live
  stack (ROADMAP N1). GitHub/Vercel/Slack reference adapters; MCP-server or
  sidecar form.

## D. Product backend  (`cmd/atlas-server` — Layer 3)

**Done:** real HTTP JSON API over the real engine — `POST /issue /verify /revoke`,
`GET /health /version /delegations /audit /graph /metrics`. In-memory store +
audit log, live freshness refresher, CORS for the UI, Prometheus `/metrics`
(also closes the atlas-lab telemetry gap). Composes the frozen kernel via public
APIs only (import-lint 0 violations). Tested end-to-end + live smoke: issue →
verify Accept (333µs) → revoke → verify Reject (`RevokedObservable`). `make ci`
green. Run: `go run ./cmd/atlas-server` (→ 127.0.0.1:8080).

- ✅ **Wire the UI to the server** — the verify console runs live against the
  real API (real verdict/trace/latency), scripted fallback when offline.
- ✅ **Operator console** (`ui/console.html`) — live delegations table (with
  revoke), issue form, audit log, trust graph, metrics; auto-refresh 3s; served
  as a Vite multi-page route, linked from the site nav. New `/stats` JSON
  endpoint on the server backs it.
- ✅ **MCP server** (`cmd/atlas-mcp`) — Atlas as agent tools over MCP stdio
  (`atlas_issue/verify/revoke/delegations/graph/audit`). Hand-rolled protocol,
  zero new deps, 0 import violations. Thin adapter over atlas-server → agents &
  the UI share one trust state. Proven: full agent session (init → tools/list →
  issue → verify accept → revoke → verify reject/RevokedObservable → graph)
  against the real engine. `make ci` green. Register: `claude mcp add atlas --
  ./atlas-mcp`. Next: agent-to-agent delegation demo; MCP `resources` for
  read-only graph/audit; optional embed-mode (no server dependency).
- ✅ **CLI** (`cmd/atlas`) — `issue · verify · revoke · delegations · graph ·
  audit · doctor · version`. Thin client over the API, tabwriter tables, ASCII
  graph (──►/──✗), `issue -q | verify -` pipe, script-friendly exit codes
  (0 accept / 3 reject). Tested + live-verified; shares state with MCP + UI.
  `make ci` green (18 packages).
  `atlas inspect` decodes a record's claims offline (no server).
- ✅ **Persistence** — durable file snapshot (atomic write) + persistent
  authority key; state survives restarts (proven). `-store` / `-key` flags.
- ✅ **AuthN** — optional bearer-token guard on `/issue` + `/revoke`
  (`-api-key` / `$ATLAS_API_KEY`, constant-time); CLI + MCP + SDK send it.
- ✅ **Docker** (`deploy/`) — distroless multi-stage image + compose
  (server + Prometheus + Grafana on real `/metrics`); compose-validated.
- ✅ **Latency histogram** in `/metrics` (`atlas_verify_latency_seconds_*`).
- ✅ **Python SDK** (`sdk/python`) — zero-dep client + example + smoke test.
- ✅ **Product overview** — `PRODUCT.md`.
- ⬜ **Remaining:** gRPC · published TS/Go SDKs · config file (YAML) + hot
  reload · non-root/read-only container hardening · TLS · durable DB backend
  behind the `Store` seam.

## B. Atlas engineering  (from earlier in the session)

- ⛔ **Run `atlas-lab/` on a real Docker host** — produce the substrate evidence
  (revocation-under-partition / cross-domain / zero-egress capture). Authored &
  validated (`docker compose config` clean); **blocked**: needs Docker + SPIRE +
  network isolation, absent in this sandbox. Scripts refuse to fabricate.
- ✅ **Node `/metrics` endpoint** — now served by `cmd/atlas-server` in
  Prometheus format (verdicts_total, issued/revoked/verified, snapshot_age).
  Point the atlas-lab Prometheus scrape at the server to light up Grafana.
  (Histogram buckets for latency still a follow-up.)
- 🟡 **Normative record wire-format RFC** — governance-gated. `tests/vectors/
  VECTORS.md` sketches the wire format; the real RFC touches the flagged
  RFC-policy tension → **founder decision** before writing.
- ⬜ **Atlas Bench / TrustPerf cross-system comparison** — Atlas vs
  Biscuit / UCAN / Macaroons on one harness. The substrate-independent bench is
  the seed; the comparison is a deferred research program.

---

## C. Repo housekeeping  (uncommitted work)

Nothing below is committed yet — decide what to keep.

- ⬜ **Untracked planning corpus** (produced in the architecture phase, never
  committed): `SYSTEM_ARCHITECTURE.md`, `MODULE_SPECIFICATION.md`,
  `INTERFACE_SPECIFICATION.md`, `IMPLEMENTATION_MASTER_PLAN.md`,
  `IMPLEMENTATION_ORDER.md`, `ENGINEERING_READINESS_CERTIFICATE.md`,
  `RISK_REGISTER.md`, `TECHNICAL_DEBT_REGISTER.md`, `ARCHITECTURE_READINESS_REVIEW.md`,
  `MODULE_INTERFACE_SPECIFICATION.md`, `PROJECT_MODULE_SPECIFICATION.md`,
  `ENGINEERING_SPRINT_1_PLAN.md`, `REPOSITORY_SKELETON.md`, `CODEGRAPH.md`,
  `AI_BOOTSTRAP.md`, `rfc/RFC-003-logical-software-architecture.md`.
- ⬜ **Assets/app**: `ui/` (the site). Web fonts ship the three Geist `.woff2`
  weights the build uses from `ui/public/fonts/` (OFL retained there); the full
  upstream Geist distribution is not vendored.
- ⬜ **Modified, uncommitted**: `.gitignore`, `tests/vectors/verdict-vectors.json`
  (both were already modified at session start).

---

## Blockers summary
- **Infrastructure-gated:** running atlas-lab (Docker/SPIRE).
- **Governance-gated:** wire-format RFC (founder RFC-policy decision).
- **Awaiting you:** UI review/tuning; commit decisions; go-ahead to wire the verifier.
