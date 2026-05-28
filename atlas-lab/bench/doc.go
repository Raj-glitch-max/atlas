// Package bench is the substrate-independent core of Atlas Bench (working
// title: TrustPerf) — a benchmark harness that measures the delegation
// primitive's costs directly against the real implementation, no mocks.
//
// It measures the metrics that do NOT require the two-domain SPIRE substrate:
// verification latency and throughput, issuance latency, integrity-validation
// latency, the record (proof) size, and the signed revoked-set snapshot size
// as revocations scale. The substrate-dependent metrics — revocation
// propagation latency, partition behavior, cross-domain verification, cold
// start — live in atlas-lab/experiments and require a real Docker host; they
// are NOT measured here and no number for them is fabricated.
//
// Run: `go test ./atlas-lab/bench -bench .`  (raw Go benchmarks)
//
//	`go test ./atlas-lab/bench -run TestTrustPerfReport -report`  (writes RESULTS.md)
package bench
