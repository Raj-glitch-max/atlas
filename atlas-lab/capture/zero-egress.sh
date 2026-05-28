#!/usr/bin/env bash
# Zero-egress proof (AT16 / INV7 / SO2): capture ALL traffic leaving the
# relying-party node during a verification, out of band, and assert that no
# packet was sent to any shared authority (domain A or a broker) while the
# core verified. This is the *out-of-band* evidence the frozen acceptance
# plan requires — not a flag in the verifier, an actual packet trace.
#
# Run inside the atlas-relying-party container (needs NET_ADMIN):
#   /atlas-lab/capture/zero-egress.sh <domain-a-cidr> <verify-command...>
#
# Exit 0 and print PROOF: ZERO EGRESS iff no packet to the forbidden CIDR was
# observed during the verification window; exit 1 (with the offending packets)
# otherwise.
set -euo pipefail

forbidden_cidr="${1:?domain-a cidr, e.g. 172.20.0.0/16}"; shift
pcap="$(mktemp /tmp/atlas-egress.XXXXXX.pcap)"

# Start capturing egress to the forbidden network before verifying.
tcpdump -n -i any -w "$pcap" "dst net ${forbidden_cidr}" >/dev/null 2>&1 &
tcpdump_pid=$!
sleep 0.5   # let the capture attach

echo "verifying with capture active (forbidden: ${forbidden_cidr})..."
"$@"        # the verification command (e.g. atlas-verify against a presented record)
verify_rc=$?

sleep 0.5
kill "$tcpdump_pid" 2>/dev/null || true
wait "$tcpdump_pid" 2>/dev/null || true

count="$(tcpdump -n -r "$pcap" 2>/dev/null | wc -l | tr -d ' ')"
echo "captured packets to forbidden network during verification: ${count}"

if [ "$count" -eq 0 ]; then
  echo "PROOF: ZERO EGRESS to the shared authority during core verification (INV7/SO2/AT16 satisfied)."
  rm -f "$pcap"
  exit "$verify_rc"
fi
echo "FAIL: egress to a shared authority observed during verification — INV7/SO2 violated:"
tcpdump -n -r "$pcap" 2>/dev/null | head -20
echo "pcap retained at: $pcap"
exit 1
