// Shared, static content — faithful to the real Atlas implementation.
export const PRINCIPAL = "spiffe://domain-a.test/workload/payments-api";
export const DELEGATE  = "spiffe://domain-b.test/agent/booking-worker";

export const STAGES = [
  ["01", "Identity binding"],
  ["02", "Signature · ES256"],
  ["03", "Freshness"],
  ["04", "Scope"],
  ["05", "Revocation"],
];

export const METRICS = [
  ["Verify → Accept", "94", "µs", "p50 · ~10.6k/s", [70,66,72,68,64,60,58,62,59,57,55,54]],
  ["Issue (seal)",    "30", "µs", "~34k/s",          [40,38,42,36,35,33,34,32,31,30,30,29]],
  ["Integrity check", "80", "µs", "~12.4k/s",        [95,90,88,86,84,82,83,81,80,80,79,80]],
  ["Proof size",      "403","B",  "compact JWS",      [403,403,403,403,403,403,403,403,403,403,403,403]],
];

export const LAB = [
  ["revocation-under-partition", "S4 · AT13 / AT14",        "host"],
  ["zero-egress capture",        "tcpdump · INV7 / SO2",    "host"],
  ["cross-domain verify",        "two SPIRE domains",       "host"],
  ["TrustPerf benchmarks",       "substrate-independent",   "green"],
  ["latency / packet-loss",      "tc netem fault injection","host"],
  ["clock-skew",                 "± tolerance boundary",    "host"],
];

export const HUBS = [
  { n: "ISSUER",   s: "domain-a" },
  { n: "DELEGATE", s: "agent" },
  { n: "VERIFIER", s: "domain-b" },
  { n: "DECISION", s: "verdict" },
];
