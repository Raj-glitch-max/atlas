// Package verify is M3, the Verification Core: the pure decision logic at
// the conformant-verifier boundary, where every guarantee of the system is
// defined (RFC-003 §7 M3; MODULE_SPECIFICATION.md M3;
// INTERFACE_SPECIFICATION.md §3).
//
// Verify runs five separately named, separately observable checks — identity
// binding (INV1), integrity via the record model (INV8), expiry within an
// explicit clock-skew tolerance (INV3, ER3), scope integrity (INV8, not
// subset re-derivation), and revocation status under the freshness policy
// (SO1, FM2/FM4) — and routes to exactly one verdict: Accept, Reject(causes),
// or InconclusiveRejected(causes). The InconclusiveRejected routing is the
// fail-closed posture and is [HYPOTHESIS] (NFR3, ER11, SO4, C-INV1): designed
// for and tested, never documented as warranted, until a V1 confirmation act.
//
// The package is pure: all inputs are injected — trust material, revocation
// status, and time readings arrive through consumer-defined ports
// (TrustMaterialPort, RevocationStatusPort, TimePort) that providers satisfy
// structurally without importing this package. Nothing reachable from a
// verification performs I/O (AP1, INV7, SO2). The freshness policy — R, the
// skew tolerance, and the S4 partition ceiling — lives here and only here
// (rule R8); constructing a verifier with an unset policy refuses.
//
// Every verification emits a DecisionTrace, unconditionally, Accepts
// included (AP13, SO5, SO8). Verifier conformance (FM10, AT30) is defined as
// implementing exactly this package's checks, answer sets, routing, and
// trace obligation.
//
// Traces: ER7, ER3, ER11 [HYP], INV1, INV3, INV7, INV8, SO1, SO2, SO4 [HYP],
// SO5, FM2, FM3, FM4, FM9, FM11, AP1, AP4, AP13.
package verify

// checkpoint: chore(security): optimize fuzzing harness execution
