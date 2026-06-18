#!/usr/bin/env bash
# The Atlas agent-capability lifecycle — the "OAuth for agents" loop, honestly.
#
# Atlas is SINGLE-HOP by design: a principal grants a scoped capability to one
# delegate, verifiable in another trust domain with no live authority call.
# This is NOT multi-hop re-delegation (A→B→C) — that is deliberately out of the
# primitive.
#
# It demonstrates, end to end against the real engine:
#   1. a principal grants a scoped capability to an agent (single hop)
#   2. attenuation — an over-scope grant is refused by the engine
#   3. a relying party verifies the capability from locally-held trust material
#      (the verification makes no call to the issuer)
#   4. the principal revokes the capability
#   5. the relying party now observes the revocation and rejects
#
# Self-contained: builds the binaries, runs its own server on a temp store,
# tears everything down.
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
BIN="$(mktemp -d)"; DATA="$(mktemp -d)"; PORT="${PORT:-8090}"
SRV=""
cleanup() { [ -n "$SRV" ] && kill "$SRV" 2>/dev/null || true; rm -rf "$BIN" "$DATA"; }
trap cleanup EXIT

echo "building atlas-server + atlas…"
go build -o "$BIN/atlas-server" "$ROOT/cmd/atlas-server"
go build -o "$BIN/atlas" "$ROOT/cmd/atlas"

"$BIN/atlas-server" -addr "127.0.0.1:$PORT" -store "$DATA/state.json" -key "$DATA/authority.key" >"$DATA/server.log" 2>&1 &
SRV=$!
export ATLAS_API="http://127.0.0.1:$PORT"
atlas() { "$BIN/atlas" "$@"; }

# wait for readiness
for _ in $(seq 1 25); do atlas doctor >/dev/null 2>&1 && break; sleep 0.2; done

PRINCIPAL="spiffe://domain-a.test/workload/payments-api"
AGENT="spiffe://domain-b.test/agent/booking-worker"

echo
echo "══ 1. principal grants a SCOPED capability to the agent (single hop) ══"
OUT="$(atlas issue --principal "$PRINCIPAL" --delegate "$AGENT" --scope read:orders,write:audit)"
echo "$OUT" | sed 's/^/   /'
INST="$(echo "$OUT" | awk '/^instance/{print $2}')"
REC="$(echo "$OUT" | tail -1)"

echo
echo "══ 2. attenuation — an over-scope grant is REFUSED by the engine ══"
atlas issue --principal "$PRINCIPAL" --delegate "$AGENT" --scope delete:everything 2>&1 | sed 's/^/   /' || true

echo
echo "══ 3. relying party verifies from locally-held trust material ══"
echo "   (the verification makes no call to the issuer — offline)"
echo "$REC" | atlas verify - | sed 's/^/   /'

echo
echo "══ 4. principal REVOKES the capability ══"
atlas revoke "$INST" | sed 's/^/   /'

echo
echo "══ 5. relying party now observes the revocation and REJECTS ══"
echo "$REC" | atlas verify - | sed 's/^/   /' || true

echo
echo "══ audit trail ══"
atlas audit --limit 6 | sed 's/^/   /'

echo
echo "══ trust graph ══"
atlas graph | sed 's/^/   /'
echo
echo "done — single-hop grant, attenuation, offline verify, independent revocation."
