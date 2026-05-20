// Package revstatus is M5, the Revocation Status Provider: the system's one
// volatile region, permanently isolated behind a fixed contract (RFC-003 §7
// M5; MODULE_SPECIFICATION.md M5; INTERFACE_SPECIFICATION.md §5).
//
// StatusOf answers, for a delegation instance, exactly one of:
// ObservablyRevoked(asOf), NotObservedRevoked(asOf), or Indeterminate. The
// as-of freshness disclosure is mandatory on knowledge answers, and the
// honest-indeterminate rule is absolute: ignorance is never expressed as
// NotObservedRevoked (AP5, INV12). The provider applies no policy — whether
// an answer's freshness satisfies R or the S4 partition ceiling is judged
// solely by the verifier (rule R8). During a partition the as-of stops
// advancing; that frozen freshness is the honest, information-theoretically
// correct signal (S4).
//
// The realization behind the contract is decided by the EXP-001 spike
// (AD-D02): OAuth Status List, push-fed store, cryptographic accumulator, or
// none. Until then this package ships the degenerate realization, which
// answers Indeterminate to every query — the honest representation of spike
// outcomes β/δ, under which the system fails closed rather than pretending
// revocation knowledge it does not have (AD-007). Every realization —
// degenerate, spike candidate, or production — must pass the
// realization-independent suite in the contracttest subpackage before being
// wired into a composition root (AD-008). Mechanism-specific per-record data
// arrives only through the record's opaque revocation-binding element
// (AD-015).
//
// Traces: ER5, SO1, FM2, FM4, INV12, AP5, AP7, AP12.
package revstatus
