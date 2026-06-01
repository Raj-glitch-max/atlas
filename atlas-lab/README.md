# Atlas Lab

A real distributed-systems laboratory for the Atlas delegation primitive — not
documentation, not mocks, not simulations. It is the concrete realization of
the frozen `lab/EXP-001-EXECUTION-PLAN.md` substrate (Phases 2–6: two
independent SPIRE deployments, network isolation, instrumentation), plus a
benchmark harness, telemetry, and reproducible experiments.

## What runs where (the honest boundary)

Atlas has reached the phase where it is missing an *environment*, not ideas.
This lab supplies the environment as infrastructure-as-code. Two tiers:

| Tier | Needs | Status |
|---|---|---|
| **Benchmarks** (`bench/`) — verification latency/throughput, issuance, integrity, proof size, revocation-snapshot scaling | nothing (pure Go, real implementation) | **runs anywhere; real numbers today** (`make -C atlas-lab bench`, results in `bench/RESULTS.md`) |
| **Substrate experiments** (`experiments/`) — revocation propagation, partition/S4 behavior, cross-domain verification, zero-egress packet proof | a real Docker daemon + SPIRE images | **authored and runnable on a host; NOT executed in the authoring sandbox** |

**No substrate result in this repository is fabricated.** Every experiment
script refuses to run (aborts, tells you to bring the substrate up) when no
Docker daemon is reachable. The benchmark numbers are real but
environment-specific and labeled as such.

## Layout

```
atlas-lab/
  bench/            TrustPerf harness (real Go benchmarks) + RESULTS.md generator
  docker/           two-domain SPIRE topology (docker-compose) + node image + SPIRE configs
  network/          fault injection: partition, latency, packet-loss, clock-skew, dns-failure
  capture/          zero-egress packet proof (AT16 / INV7 / SO2), out of band via tcpdump
  telemetry/        Prometheus scrape config + Grafana dashboard
  experiments/      reproducible experiments (flagship: revocation-under-partition)
  scripts/          run-all.sh (one command; skips substrate honestly if no daemon)
  Makefile          bench / report / config / up / down / experiment / run-all
```

## Quick start

```bash
# anywhere — real benchmark numbers against the real implementation:
make -C atlas-lab bench
make -C atlas-lab report      # writes bench/RESULTS.md

# validate the substrate topology without a daemon:
make -C atlas-lab config

# on a real Docker host — the full lab:
make -C atlas-lab run-all     # bench + two-domain SPIRE + the flagship experiment
```

## The two trust domains

Two **independent** SPIRE deployments (`domain-a.test`, `domain-b.test`) with
**federation disabled** (frozen EXP-001 assumption U1), on separate Docker
networks so the partition experiment severs the relying party's cross-domain
link at the link level. Domain A issues delegations and publishes signed
revocation snapshots; domain B's relying party verifies **offline** — the lab
proves, by packet capture, that it makes no live call to domain A.

## Governance

This lab realizes the frozen `lab/` protocol; it does not modify it. Substrate
runs inherit the lab discipline (`lab/LAB_README.md`,
`lab/EXPERIMENT_CHECKLIST.md`, `lab/DECISION_RULES.md`): two-run
reproducibility, adversary-blinded roles, pre-registered scope parameters
(R and the S4 reading from `agents/journal/2026-07-06-c4-spike-scope-act.md`),
and append-only evidence. `atlas-lab/` is additive; the frozen package is
untouched.

## What the lab will let Atlas finally say

Instead of "I think offline verification makes no live call," the lab lets an
operator *demonstrate* it: start two SPIRE domains, issue a delegation,
disconnect the network, revoke, wait, verify, capture packets, prove zero
egress, collect metrics, repeat, generate a report. That is publishable
engineering evidence — and it is the milestone that turns Atlas from a
rigorous repository into an engineering artifact.
