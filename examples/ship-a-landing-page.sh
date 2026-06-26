#!/usr/bin/env bash
# "Ship a landing page" — the multi-tool agent workflow, made safe with Atlas.
#
# The pain: one agent (or a team of them) has to act across GitHub, Vercel, and
# Slack on your behalf. Today you'd hand it your GitHub PAT, your Vercel token,
# your Slack token — now a compromised or confused agent can do ANYTHING to
# ALL of them, forever, and you can't tell what it did.
#
# Atlas is NOT the agent and NOT the orchestrator. Atlas is the trust fabric:
# you hand each agent a scoped, expiring, revocable *capability* instead of a
# secret. Each tool ("adapter") verifies the capability OFFLINE and, only if
# it authorizes the exact action, performs it using the tool's OWN credential —
# which the agent never sees.
#
# This script proves, against the real engine:
#   • least privilege        — each agent can touch exactly one tool + action
#   • blast-radius control    — a compromised agent is confined to its scope
#   • instant containment     — one revoke kills a capability everywhere
#   • offline authorization   — adapters never call the issuer
#
# Self-contained: builds, runs its own server, cleans up.
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
BIN="$(mktemp -d)"; DATA="$(mktemp -d)"; PORT="${PORT:-8101}"
SRV=""
cleanup() { [ -n "$SRV" ] && kill "$SRV" 2>/dev/null || true; rm -rf "$BIN" "$DATA"; }
trap cleanup EXIT
say() { printf '\n\033[1m%s\033[0m\n' "$*"; }

echo "building…"
go build -o "$BIN/atlas-server" "$ROOT/cmd/atlas-server"
go build -o "$BIN/atlas" "$ROOT/cmd/atlas"
export ATLAS_API="http://127.0.0.1:$PORT"
atlas() { "$BIN/atlas" "$@"; }

# The authority grants these tool actions to your principal to delegate onward.
"$BIN/atlas-server" -addr "127.0.0.1:$PORT" -store "$DATA/s.json" -key "$DATA/k.pem" \
  -trust-domain acme.dev \
  -grant 'github:push:acme/landing,vercel:deploy:landing-prod,slack:post:launches' >"$DATA/log" 2>&1 &
SRV=$!
for _ in $(seq 1 25); do atlas doctor >/dev/null 2>&1 && break; sleep 0.2; done

YOU="spiffe://acme.dev/human/raj"

# A tool adapter: verify the capability OFFLINE, check it authorizes THIS action,
# and only then "act" with the tool's own credential (which the agent never has).
adapter() {
  local tool="$1" scope="$2" rec="$3" action="$4"
  if echo "$rec" | atlas verify --offline --bundle "$DATA/bundle.json" --require-scope "$scope" - >/dev/null 2>&1; then
    printf '   [%s] \033[32m✓ authorized\033[0m → %s\n' "$tool" "$action"
    printf '        (performed with %s'"'"'s own credential — the agent never saw it)\n' "$tool"
  else
    printf '   [%s] \033[31m✗ DENIED\033[0m — capability does not authorize %s\n' "$tool" "$scope"
  fi
}

say "1 · you hand each task-agent a narrowly-scoped, expiring capability"
CODER=$(atlas delegate -q --principal "$YOU" --delegate spiffe://acme.dev/agent/coder    --scope github:push:acme/landing   --ttl 900)
DEPLOYER=$(atlas delegate -q --principal "$YOU" --delegate spiffe://acme.dev/agent/deployer --scope vercel:deploy:landing-prod --ttl 600)
NOTIFIER=$(atlas delegate -q --principal "$YOU" --delegate spiffe://acme.dev/agent/notifier --scope slack:post:launches      --ttl 300)
CODER_INST=$(echo "$CODER" | atlas inspect - | grep atl_ins | grep -oE '[0-9a-f]{32}')
atlas bundle -o "$DATA/bundle.json" >/dev/null
echo "   coder    → github:push:acme/landing   (15m)"
echo "   deployer → vercel:deploy:landing-prod (10m)"
echo "   notifier → slack:post:launches         (5m)"
echo "   (no agent holds a real GitHub/Vercel/Slack token — only an Atlas capability)"

say "2 · the happy path — each agent acts on exactly its tool, verified offline"
adapter github "github:push:acme/landing"   "$CODER"    "git push origin main → acme/landing"
adapter vercel "vercel:deploy:landing-prod" "$DEPLOYER" "deploy → landing-prod.vercel.app"
adapter slack  "slack:post:launches"        "$NOTIFIER" "post → #launches: 'landing page is live 🚀'"

say "3 · blast radius — the CODER agent is compromised and reaches for more"
echo "   it holds a valid capability, so a naive token-sharing setup would let it"
echo "   deploy and post. With Atlas, its capability only grants github:push:"
adapter vercel "vercel:deploy:landing-prod" "$CODER" "malicious deploy"
adapter slack  "slack:post:launches"        "$CODER" "spam #launches"

say "4 · instant containment — you revoke the coder; re-export the bundle"
atlas revoke "$CODER_INST" >/dev/null
atlas bundle -o "$DATA/bundle.json" >/dev/null
echo "   the coder's own github:push now dies too — offline, everywhere:"
adapter github "github:push:acme/landing" "$CODER" "git push"
echo "   …while the deployer's independent capability is untouched:"
adapter vercel "vercel:deploy:landing-prod" "$DEPLOYER" "redeploy hotfix"

say "Atlas didn't ship the landing page. Atlas is why you could let the agents do it."
