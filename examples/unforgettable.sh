#!/usr/bin/env bash
# THE Atlas demo — the unforgettable moment.
#
#   1. An agent gets a scoped capability.        (atlas delegate)
#   2. The relying party exports a trust bundle. (atlas bundle)
#   3. THE SERVER DIES. The network is gone.
#   4. Verification still works — locally, in µs.   ← the point of Atlas
#   5. Reconnect. Revoke. Re-export.
#   6. The revoked capability is rejected — offline.
#   7. A snapshot older than YOUR staleness budget fails CLOSED.
#   8. A tampered bundle (revocation stripped) is REFUSED outright.
#
# Everything here is the real engine — no mocks, no scripted output.
# Self-contained: builds, runs its own server on a temp store, cleans up.
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
BIN="$(mktemp -d)"; DATA="$(mktemp -d)"; PORT="${PORT:-8095}"
SRV=""
cleanup() { [ -n "$SRV" ] && kill "$SRV" 2>/dev/null || true; rm -rf "$BIN" "$DATA"; }
trap cleanup EXIT
say() { printf '\n\033[1m%s\033[0m\n' "$*"; }

echo "building…"
go build -o "$BIN/atlas-server" "$ROOT/cmd/atlas-server"
go build -o "$BIN/atlas" "$ROOT/cmd/atlas"
export ATLAS_API="http://127.0.0.1:$PORT"
atlas() { "$BIN/atlas" "$@"; }
start_server() {
  "$BIN/atlas-server" -addr "127.0.0.1:$PORT" -store "$DATA/state.json" -key "$DATA/authority.key" >>"$DATA/server.log" 2>&1 &
  SRV=$!
  for _ in $(seq 1 25); do atlas doctor >/dev/null 2>&1 && return; sleep 0.2; done
  echo "server failed to start"; exit 1
}
stop_server() { kill -TERM "$SRV" 2>/dev/null; wait "$SRV" 2>/dev/null || true; SRV=""; }

start_server

say "1 · an agent gets a scoped capability"
OUT="$(atlas delegate --principal spiffe://domain-a.test/workload/payments-api \
                      --delegate  spiffe://domain-b.test/agent/booking-worker \
                      --scope read:orders,write:audit)"
echo "$OUT" | sed 's/^/   /' | head -4
INST="$(echo "$OUT" | awk '/^instance/{print $2}')"
REC="$(echo "$OUT" | tail -1)"
# a second, never-revoked capability (used for the staleness beat in step 7)
REC_OK="$(atlas delegate -q --principal spiffe://domain-a.test/workload/payments-api \
                          --delegate spiffe://domain-b.test/agent/audit-reader --scope read:orders)"

say "2 · export the trust bundle (public key + signed revocation snapshot)"
atlas bundle -o "$DATA/bundle.json" | sed 's/^/   /' | head -4

say "3 · THE SERVER DIES — no issuer, no network"
stop_server
atlas doctor 2>&1 | grep 'server:' | sed 's/^/   /' || true

say "4 · verify OFFLINE — still answers, in microseconds"
echo "$REC" | atlas verify --offline --bundle "$DATA/bundle.json" --max-staleness 10m - | sed 's/^/   /'

say "5 · reconnect · revoke · re-export · die again"
start_server
atlas revoke "$INST" | sed 's/^/   /'
atlas bundle -o "$DATA/bundle2.json" >/dev/null
stop_server
echo "   (fresh bundle carries the revocation; server is dead again)"

say "6 · the revoked capability — rejected, OFFLINE"
echo "$REC" | atlas verify --offline --bundle "$DATA/bundle2.json" --max-staleness 10m - | sed 's/^/   /' || true

say "7 · staleness is YOUR policy — an old snapshot fails CLOSED, never open"
sleep 1
echo "   (a VALID capability, but the snapshot is now older than a 500ms budget)"
echo "$REC_OK" | atlas verify --offline --bundle "$DATA/bundle2.json" --max-staleness 500ms - | sed 's/^/   /' | head -2 || true

say "8 · a tampered bundle (revocation stripped) is REFUSED outright"
python3 - "$DATA/bundle2.json" <<'PY'
import json, sys
b = json.load(open(sys.argv[1]))
b["revocation"]["revoked"] = []           # attacker hides the revocation
json.dump(b, open(sys.argv[1].replace('.json', '-tampered.json'), 'w'))
PY
echo "$REC" | atlas verify --offline --bundle "$DATA/bundle2-tampered.json" - 2>&1 | sed 's/^/   /' | head -1 || true

say "that's Atlas: delegation that survives the network — and can't be lied to."
