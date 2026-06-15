#!/usr/bin/env bash
# One-command Atlas Lab run: bring up the two-domain SPIRE substrate, run the
# substrate-independent benchmarks (these work anywhere), then the flagship
# substrate experiment (needs a real Docker daemon). Refuses to fabricate.
set -euo pipefail
cd "$(dirname "$0")/.."
COMPOSE="docker compose -f docker/docker-compose.yml"

echo "== Atlas Lab =="

echo "[1/3] substrate-independent benchmarks (TrustPerf; runs anywhere) ..."
( cd ../ && go test ./atlas-lab/bench -run TestTrustPerfReport -report ) && \
  echo "     -> atlas-lab/bench/RESULTS.md"

if ! docker info >/dev/null 2>&1; then
  echo "[2/3] SKIP substrate: no Docker daemon reachable."
  echo "[3/3] SKIP experiment: needs the substrate."
  echo
  echo "To run the full lab on a real host:"
  echo "  $COMPOSE up -d --build"
  echo "  ./experiments/revocation-under-partition/run.sh"
  echo "No substrate result is fabricated."
  exit 0
fi

echo "[2/3] bringing up the two-domain SPIRE topology ..."
$COMPOSE up -d --build

echo "[3/3] running the flagship experiment (twice, per lab reproducibility discipline) ..."
./experiments/revocation-under-partition/run.sh
./experiments/revocation-under-partition/run.sh

echo "done. dashboards: http://localhost:3000  metrics: http://localhost:9090"
