# Experiment: revocation under partition

The flagship substrate experiment — the concrete run of the S4 case (FM1) and
acceptance tests AT13/AT14 against two real, independent SPIRE domains, at the
resolved scope act (R = 2 s; `agents/journal/2026-07-06-c4-spike-scope-act.md`).

**What it proves (on a real host):**
- offline verification with **zero egress** to domain A (INV7/SO2/AT16), via
  packet capture, not a verifier flag;
- that a revocation performed **while the RP is partitioned** from the issuer
  is **not claimed observable** during the partition (INV12/S4 — the honest
  bound), and
- that after recovery + ≤ R, the revoked delegation is **rejected**
  (`RevokedObservable`), with the propagation latency recorded.

**Governance:** this is the executable form of frozen
`lab/EXP-001-EXECUTION-PLAN.md` (Phases 9–12), and it inherits the lab
discipline: run twice for reproducibility, adversary-blinded role separation,
pre-registered R and S4 reading, append-only evidence. It does not modify the
frozen plan; it realizes it.

**Honesty:** `run.sh` refuses to fabricate — if no Docker daemon/topology is
reachable it aborts and tells you to bring the substrate up. No result in this
repository is a substrate result until a real host produces one.

Run:

```
atlas-lab/scripts/run-all.sh            # brings up the topology, runs the experiment
# or, once the topology is up:
atlas-lab/experiments/revocation-under-partition/run.sh
```

<!-- checkpoint: chore(test): clarify simulated agent node -->
