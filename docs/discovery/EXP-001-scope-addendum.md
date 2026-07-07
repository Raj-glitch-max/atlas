# EXP-001 pre-registration addendum — resolved scope parameters

**Status:** Non-frozen addendum to the frozen `lab/EXP-001-EXECUTION-PLAN.md`.
Created per that plan's §1 exit criteria, which require the resolved S1–S5
values to be "reflected into a run-specific addendum to the spike's pre-reg
section (a pre-reg amendment, dated, not a retrofit of the frozen gate)."
The frozen gate and plan are unmodified; this addendum records the concrete,
non-TBD thresholds the run witnesses.
**Authority:** `agents/journal/2026-07-06-c4-spike-scope-act.md`.

## Resolved thresholds (pre-registered)

| Param | Resolved value | Meaning for the run |
|---|---|---|
| **S1 — R** | **2 s** (R < delegation TTL) | A `NotObservedRevoked` observation older than 2 s is treated as stale → inconclusive → fail closed. Injected as `verify.Policy.R`. |
| **S2 — cached pulls** | **admitted** for signed, integrity-protected artifacts with bounded staleness | The signed-revoked-set / OAuth-Status-List composition is admissible. |
| **S3 — broker** | decision-maker only; a **passive signed-blob distributor is not a broker** | Distributing a signed revocation snapshot via a passive cache/CDN is compliant. |
| **S4 — partition** | **eventual-upon-recovery, within P = R of recovery; no in-partition observability claim** | A revocation performed while the RP is partitioned is not claimed observable until recovery; the run must not assert otherwise (INV12). Distinguishes outcome β from γ. |
| **S5 — per-hop authz** | **not required** (multi-hop out of V1) | Not exercised. |

## Pre-registered spike criterion (now non-TBD)

> On a SPIFFE substrate with signed delegation records, the signed-revoked-set
> composition (S2-admissible) delivers revocation observability within R = 2 s
> of a revocation becoming published, with no live RP→domain-A call and no
> decision-making broker (S3-compliant), under the S4 eventual-upon-recovery
> reading.

## Disposition of the composition attempt (this session)

The signed-revoked-set realization (`internal/revstatus/statuslist.go`) was
**built and contract-tested in-process** against the resolved parameters
(`agents/journal/2026-07-06-e7-alpha-signed-revoked-set.md`), demonstrating the
α-path logic: verifiable-freshness snapshots, revoked → `RevokedObservable`,
staleness > R → fail closed. This is the composition-attempt in code. The
**substrate-validated spike** (two independent SPIRE deployments, link-level
partition, out-of-band egress sniffer, adversary-blinding, two-run
reproducibility — EXP-001 plan Phases 2–12) is **not** performed here and its
result is **not** claimed; that remains Epic E6/E7 and requires infrastructure
absent from this environment.
