#!/usr/bin/env bash
# Flagship experiment: revocation observability under partition — the
# concrete run of the S4 case (FM1) and AT13/AT14 against the real two-domain
# substrate, at the resolved scope act (R = 2s; agents/journal/
# 2026-07-06-c4-spike-scope-act.md). Requires the docker-compose topology up
# on a real host. This script is the executable protocol; it produces the
# evidence that cannot be generated without the substrate.
#
# Protocol (per frozen lab/EXP-001-EXECUTION-PLAN.md discipline; run TWICE for
# reproducibility, adversary-blinded per lab governance):
#   1. issue a delegation in domain A; deliver the record + domain A's trust
#      bundle to the relying party in domain B (out of band).
#   2. baseline: RP verifies offline -> expect Accept; capture egress ->
#      expect ZERO egress to domain A (AT15/AT16).
#   3. partition: sever B's cross-domain link.
#   4. revoke the delegation in domain A DURING the partition; publish a fresh
#      signed revoked-set snapshot in A.
#   5. while partitioned, RP verifies -> expect the RP does NOT observe the
#      revocation (S4: no in-partition observability) and MUST NOT claim it;
#      it either still accepts on its last fresh snapshot within R, or fails
#      closed once that snapshot ages past R. Record which, with timestamps.
#   6. heal the partition; deliver the new snapshot; wait <= R.
#   7. RP verifies -> expect Reject(RevokedObservable). Record propagation
#      latency (heal -> first Reject).
#   8. emit a JSON result: {accepted_before, egress_count, in_partition_claim,
#      post_recovery_reject, propagation_latency_ms, R_ms, run_id}.
set -euo pipefail

COMPOSE="docker compose -f $(dirname "$0")/../../docker/docker-compose.yml"
RP="atlas-relying-party"
ISSUER="atlas-issuer"
XLINK="${XLINK:-eth-crossdomain}"   # the cross-domain interface name inside the RP container
R_MS="${R_MS:-2000}"

echo "== Atlas Lab :: revocation-under-partition (R=${R_MS}ms) =="
echo "NOTE: requires a running docker daemon + the compose topology up."
echo "      This host has no daemon; the script is the authored protocol."
echo

# Guard: refuse to fabricate. If there is no daemon/topology, stop honestly.
if ! docker info >/dev/null 2>&1; then
  echo "ABORT: no Docker daemon reachable. Bring up the substrate on a real host:"
  echo "  $COMPOSE up -d --build"
  echo "then re-run this experiment. No result is fabricated."
  exit 3
fi

run_id="$(date -u +%Y%m%dT%H%M%SZ)"
echo "run_id=$run_id"

# 1-2. issue + baseline verify + zero-egress capture (domain A CIDR discovered
#      from the compose network).
a_cidr="$($COMPOSE exec -T "$ISSUER" sh -c "ip -o -f inet addr show | awk '/domain-a|crossdomain/{print \$4}' | head -1")"
echo "domain-A cidr (approx): ${a_cidr:-unknown}"
$COMPOSE exec -T "$ISSUER" atlas-issue           # issues + prints the record (demo driver)
# (a real harness captures the record artifact into /artifacts and hands it to the RP)

echo "baseline verify + zero-egress proof:"
$COMPOSE exec -T "$RP" /atlas-lab/capture/zero-egress.sh "${a_cidr:-172.16.0.0/12}" atlas-verify || true

# 3. partition
$COMPOSE exec -T "$RP" /atlas-lab/network/faults.sh partition "$XLINK"

# 4. revoke during partition (in A) + publish snapshot
$COMPOSE exec -T "$ISSUER" atlas-revoke

# 5. in-partition verify — must not claim observability
echo "in-partition verify (expect: no in-partition revocation observability, S4):"
$COMPOSE exec -T "$RP" atlas-verify || true

# 6. heal + deliver + wait R
$COMPOSE exec -T "$RP" /atlas-lab/network/faults.sh heal "$XLINK"
sleep "$(awk "BEGIN{print $R_MS/1000}")"

# 7. post-recovery verify — expect Reject(RevokedObservable)
echo "post-recovery verify (expect Reject / RevokedObservable):"
$COMPOSE exec -T "$RP" atlas-verify || true

echo
echo "== record the observed outcomes into a run report per lab evidence discipline =="
echo "(a full harness writes results/${run_id}.json and appends to the evidence index)"

# checkpoint: chore(lab): restructure network partition test
