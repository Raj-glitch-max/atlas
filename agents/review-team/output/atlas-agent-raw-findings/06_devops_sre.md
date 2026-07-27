# 06 — DEVOPS / SRE — RAW FINDINGS

Phase 3 (Ops & Product). I own SLOs, runbooks, and the availability-coupling arithmetic.

FINDING-SRE-1  (the one that matters)
  Type:        availability-coupling / evidence
  Claim:       `verify_availability ≤ snapshot_availability` is real arithmetic and it is
               UNMEASURED. There is no published curve of freshness-availability vs R vs snapshot
               cadence, and the link-level zero-egress proof (two SPIRE domains, severed cable,
               tcpdump) is authored in `atlas-lab/` but has never run on real infrastructure.
  Evidence:    LIMITATIONS §6–8; `docs/runbooks/V1_OPERATION.md` exists but the number does not.
  Impact:      An operator cannot set R responsibly or write an SLO without this. It is the single
               most valuable measurement in the project. (→ S7, OQ-4)
  Recommendation: NEXT — stand up the substrate, publish the availability curve. Incubation critical path.
  Owner:       06_devops_sre

FINDING-SRE-2
  Type:        supply-chain / release
  Claim:       No signed releases, no SLSA provenance, and go-jose v3 is pinned but not gated by
               govulncheck in CI. An unsigned release of a verification library is a contradiction.
  Evidence:    `.github/workflows/{go-ci,repo-health}.yml` present; secret scan + frozen-doc gates
               exist; provenance/signing absent.
  Recommendation: NEXT — cosign + SLSA provenance + govulncheck job (co-own 03, S10).
  Owner:       06_devops_sre

FINDING-SRE-3
  Type:        observability
  Claim:       The decision trace (five stages, unconditional, emitted on Accept too) is an
               excellent operability primitive and is already built — credit it. What is missing is
               the *aggregate* view: metrics exist on the server (`/metrics`, Prometheus) but there
               is no documented SLO or alert on "verifiers aging past R" (the fleet-halt leading
               indicator).
  Recommendation: NEXT — an SLO + alert on snapshot-age p99 approaching R, so the self-DoS is seen
               before it bites. Pairs with SRE-1's number.
  Owner:       06_devops_sre

FINDING-SRE-4
  Type:        deployment maturity
  Claim:       The server is dev-grade (no TLS termination, no rate-limit hardening beyond per-IP,
               JSON-file store, single-node) — acknowledged honestly (LIMITATIONS §5). The seams
               (Store interface, stateless engine) exist; the hardened deployment does not.
  Recommendation: LATER — harden when a real deployment exists (S11, blocked by S8).
  Owner:       06_devops_sre

VERDICT: Ops posture is honest and the trace/metrics foundation is real. The gap is one
measurement (SRE-1) and the supply-chain paperwork (SRE-2), not a broken operational surface.
