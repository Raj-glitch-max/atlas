// Package harness is the acceptance-test instrumentation layer
// (MODULE_SPECIFICATION.md §Non-module code; SYSTEM_ARCHITECTURE.md §11).
//
// It is instrumentation, not product: import rules that bind the internal
// modules do not bind it, and no product package may import it. Across the
// implementation epics it grows the port fakes and controllable clock for
// in-process verification tests (E3-T2) and the substrate-control interfaces
// — two-domain control, link-level partition induction, egress observation —
// realized against the shared EXP-001 substrate (E6-T2). In-process fakes
// and the real substrate implement the same interfaces, so acceptance tests
// run identically at both boundaries.
//
// Honest-negative discipline binds fakes exactly as it binds realizations:
// a fake must never express ignorance as NotObservedRevoked (FD-4,
// TECHNICAL_DEBT_REGISTER.md), and fabricated freshness exists only where a
// test's pre-registered scenario calls for it.
package harness
