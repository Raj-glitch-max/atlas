#!/usr/bin/env bash
# Network fault injection for Atlas Lab, applied inside a node container
# (needs NET_ADMIN). Subcommands model the distributed conditions the frozen
# failure model (04_FAILURE_MODEL.md) and EXP-001 require:
#
#   partition   <iface>            sever the link (iptables DROP) — the S4 case
#   heal        <iface>            restore the link
#   latency     <iface> <ms>       add fixed delay (tc netem)
#   packet-loss <iface> <percent>  drop a fraction of packets
#   clock-skew  <seconds>          offset the container clock (needs SYS_TIME)
#   dns-failure                    blackhole DNS (port 53) — resolution failure
#   clear       <iface>            remove all tc qdiscs
#
# Usage (from the host):
#   docker compose exec atlas-relying-party \
#     /atlas-lab/network/faults.sh partition eth-crossdomain
set -euo pipefail

cmd="${1:-}"; shift || true

case "$cmd" in
  partition)
    iface="${1:?iface}"
    # Drop all traffic on the cross-domain link, in both directions.
    iptables -A OUTPUT -o "$iface" -j DROP
    iptables -A INPUT  -i "$iface" -j DROP
    echo "partitioned: $iface (all traffic dropped)"
    ;;
  heal)
    iface="${1:?iface}"
    iptables -D OUTPUT -o "$iface" -j DROP 2>/dev/null || true
    iptables -D INPUT  -i "$iface" -j DROP 2>/dev/null || true
    echo "healed: $iface"
    ;;
  latency)
    iface="${1:?iface}"; ms="${2:?ms}"
    tc qdisc replace dev "$iface" root netem delay "${ms}ms"
    echo "latency: ${ms}ms on $iface"
    ;;
  packet-loss)
    iface="${1:?iface}"; pct="${2:?percent}"
    tc qdisc replace dev "$iface" root netem loss "${pct}%"
    echo "packet-loss: ${pct}% on $iface"
    ;;
  clock-skew)
    secs="${1:?seconds}"
    date -s "@$(( $(date +%s) + secs ))"
    echo "clock skewed by ${secs}s"
    ;;
  dns-failure)
    iptables -A OUTPUT -p udp --dport 53 -j DROP
    iptables -A OUTPUT -p tcp --dport 53 -j DROP
    echo "dns blackholed"
    ;;
  clear)
    iface="${1:?iface}"
    tc qdisc del dev "$iface" root 2>/dev/null || true
    echo "cleared tc on $iface"
    ;;
  *)
    echo "usage: faults.sh {partition|heal|latency|packet-loss|clock-skew|dns-failure|clear} ..." >&2
    exit 2
    ;;
esac
