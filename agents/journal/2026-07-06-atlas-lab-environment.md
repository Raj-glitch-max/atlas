---
date: 2026-07-06
slug: atlas-lab-environment
artifact: founder direction — Atlas's phase transition is from "missing ideas" to "missing an environment"; build atlas-lab
decision: Built atlas-lab/ — a real distributed-systems laboratory as infrastructure-as-code: a runnable TrustPerf benchmark harness (real numbers now), the two-domain SPIRE topology, network fault injection, zero-egress packet capture, telemetry, and the flagship revocation-under-partition experiment. The substrate-independent tier runs here; the SPIRE tier is authored, runnable on a host, and never fabricates results.
agents_consulted: [empiricist, red-team, operator, economist, cartographer]
overrides: false
related_entries: [c4-spike-scope-act, e7-alpha-signed-revoked-set, language-neutral-conformance-vectors]
---

# Context

The founder named Atlas's phase transition precisely: it is no longer missing
ideas, it is missing an *environment*. Directed to stop spending on
repository-local philosophy and build `atlas-lab/` — real infrastructure: a
two-domain SPIRE topology, network fault injection, packet capture, a
benchmark harness (TrustPerf), telemetry, experiment automation, reproducible
reports. This is the concrete realization of the frozen
`lab/EXP-001-EXECUTION-PLAN.md` substrate (Phases 2-6).

# Decision

Built `atlas-lab/` in two honest tiers:

- **Substrate-independent (runs anywhere, real numbers today):** `bench/` — a
  Go benchmark harness measuring, against the real implementation:
  verification latency (~94µs this env — empirically confirming gate C9's
  sub-100ms finding), throughput (~10.6k verify/s), issuance (~30µs),
  integrity validation (~80µs), record proof size (403 bytes), and the signed
  revoked-set snapshot size vs. revocations (O(revoked)). `RESULTS.md` is
  generated from a measured run, labeled environment-specific/indicative, with
  an explicit "NOT measured here" section for the substrate-dependent metrics.
- **Substrate (authored IaC, runnable on a real Docker host, never
  fabricated):** `docker/` — two independent SPIRE deployments (domain-a.test,
  domain-b.test), federation disabled (frozen EXP-001 U1), on separate
  networks so the partition experiment severs the RP's cross-domain link;
  `network/faults.sh` (partition/latency/loss/clock-skew/dns); `capture/
  zero-egress.sh` (the AT16/INV7/SO2 packet proof, out of band via tcpdump);
  `telemetry/` (Prometheus + Grafana); `experiments/revocation-under-partition`
  (the flagship S4/AT13/AT14 run at the resolved R = 2s), plus
  `scripts/run-all.sh` and `Makefile`.

**Every substrate script refuses to fabricate:** it aborts with a clear
"bring the substrate up on a real host" message when no Docker daemon is
reachable. No substrate result exists in the repo until a real host produces
one.

Governance: `atlas-lab/` is additive and realizes the frozen `lab/` protocol
(it does not modify it); substrate runs inherit the lab discipline (two-run
reproducibility, adversary-blinding, pre-registered R/S4). Validated: bench
runs green; `docker compose config` valid; scripts pass `bash -n`; `make ci`
green (15 packages).

# Evidence cited

- `atlas-lab/bench/RESULTS.md` — a real measured run (this environment).
- `docker compose config` — topology syntactically valid (client-side, no
  daemon).
- `bash -n` on all four scripts — clean.
- Frozen `lab/EXP-001-EXECUTION-PLAN.md` Phases 2-6 (the substrate this
  realizes); `agents/journal/2026-07-06-c4-spike-scope-act.md` (R = 2s).

# Council positions

## The Empiricist
The bench numbers are real and reproducible; the sub-100ms verification result
is now empirical, not asserted (gate C9 confirmed). Crucial discipline held:
the RESULTS.md "NOT measured here" section and the abort-on-no-daemon guards
mean nothing substrate-dependent is claimed. Confidence in the measured
metrics: High (this env); they are labeled environment-specific.

## The Red Team
The single most important property of this lab is that it CANNOT lie: the
experiment scripts abort rather than emit a fabricated pass. That is the right
default for a security lab. Caveat on record: the node `/metrics` surface the
telemetry configs scrape does not exist yet (the drivers are demos) — the
dashboards bind to lab-convention metric names; wiring a real /metrics
endpoint is the next node increment, and until then Grafana shows no data.
Named honestly in prometheus.yml and the dashboard, not hidden.

## The Operator
`make -C atlas-lab bench` gives an engineer real numbers in one command
anywhere; `make config` validates the topology without a daemon;
`run-all.sh` is the one-command full lab on a host. This is the environment
that lets Atlas say "watch" instead of "I think."

## The Economist
Correct spend: the substrate-independent bench delivers real value now (the
first empirical performance evidence), and the IaC is the reusable asset that
unblocks all future substrate evidence. The expensive part (running SPIRE at
scale, 1000× partition runs) is deferred to a real host — appropriately, since
this sandbox cannot run it and fabricating would be worthless.

## The Cartographer
Restate: atlas-lab is the environment, realized as IaC, in two tiers split by
what this sandbox can honestly execute. It realizes the frozen EXP-001
substrate; it does not expand Atlas's scope (the primitive is unchanged) and
does not modify the frozen package. The bench is the runnable evidence; the
SPIRE experiments are the authored apparatus awaiting a host.

# Dissent preserved

No dissent. Red Team's caveat (the node /metrics endpoint the telemetry
assumes is not yet built; dashboards are empty until it is) recorded as the
next honest increment, not a disagreement.

# Founder override (if applicable)

None; built under the technical-lead grant.

# Open questions / next

- **A node `/metrics` endpoint** (Prometheus surface on the drivers) so the
  telemetry/dashboards show real data — the next node increment.
- **Running the substrate on a real Docker host** — the only way to produce
  the partition/cross-domain/zero-egress evidence; needs infrastructure absent
  here.
- **Atlas Bench (TrustPerf) cross-system comparison** (Atlas vs Biscuit/UCAN/
  Macaroons on one harness) — the founder's named next frontier; the
  substrate-independent bench is the seed. Deferred; a real research program.

# Status
- decided: 2026-07-06

<!-- checkpoint: repo(founder-profile-feedback): document founder profile feedback (#25) -->

<!-- checkpoint: rfc(glossary-definitions): document glossary definitions -->

<!-- checkpoint: docs(deployment-manual): restructure deployment manual -->

<!-- checkpoint: repo(conformance-targets): refine conformance targets -->

<!-- checkpoint: planning(conformance-targets): refine conformance targets -->

<!-- checkpoint: feat(sdk): implement truststore backend -->
