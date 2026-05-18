// Package record is M1, the Record Model: the stable surface of the system
// (RFC-003 §7 M1; MODULE_SPECIFICATION.md M1; INTERFACE_SPECIFICATION.md §1).
//
// It defines the Delegation Record — the single, self-sufficient,
// tamper-evident artifact that binds a principal identity, a delegate
// identity, a scope, an expiration, an issuance time, an opaque instance
// identity (AD-013), and an opaque optional revocation binding (AD-015) —
// together with integrity validation (Intact | Altered) runnable by any
// holder, and reconstruction reading of intact records.
//
// The package is pure: it depends on no other package in this module and
// performs no I/O, holds no clock, and touches no network. Records validate
// as Intact or Altered regardless of validity; validity (expiry, revocation,
// scope) is the verifier's question, answered in internal/verify.
//
// Record construction is restricted to the Issuance Authority
// (internal/issuance); no other caller may create records (AP10).
//
// Traces: ER1, ER4, INV1, INV8, INV9, SO3.
package record
