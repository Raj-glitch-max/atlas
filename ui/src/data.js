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

// Measured, not estimated. Regenerate with `bash scripts/run-benchmarks.sh`
// and copy from docs/BENCHMARKS.md — these are p50 of the in-process engine
// (Ryzen 5 5600H, Go 1.22.11, linux/amd64, chain depth 1), NOT end-to-end
// server latency, which is ~2-3x and exported live as the
// atlas_verify_latency_seconds Prometheus histogram.
export const METRICS = [
  ["Verify → Accept", "116", "µs", "p50 · ~8.6k/s", [70,66,72,68,64,60,58,62,59,57,55,54]],
  ["Issue (seal)",    "37",  "µs", "p50 · ~27k/s",  [40,38,42,36,35,33,34,32,31,30,30,29]],
  ["Integrity check", "110", "µs", "p50 · ~9.1k/s", [95,90,88,86,84,82,83,81,80,80,79,80]],
  ["Proof size",      "403", "B",  "compact JWS",   [403,403,403,403,403,403,403,403,403,403,403,403]],
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
