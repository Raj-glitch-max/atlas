// Package truststore is M4, the Trust Material Store: the relying party's
// locally held trust material (RFC-003 §7 M4; MODULE_SPECIFICATION.md M4;
// INTERFACE_SPECIFICATION.md §4).
//
// It holds material established only by explicit, out-of-band provisioning
// acts and answers material-or-absent queries at verification time. The
// package is structurally incapable of fetching: absent material is an
// honest answer that the verifier routes to its inconclusive (fail-closed)
// path — the insecure-fallback path of FM9 is deleted by construction, not
// discouraged by convention. Each provisioning act appends a provisioning
// record; refused (malformed or incoherent) material is never stored and
// never half-trusted.
//
// The store satisfies the verifier's TrustMaterialPort structurally; it does
// not import internal/verify (rule R3). Concurrency posture: single-writer
// (provisioning acts) / multi-reader (verifications) — AD-017.
//
// Traces: ER7, INV7, NFR2, C6, FM9.
package truststore

// checkpoint: chore(test): audit simulated agent node

// checkpoint: feat(client): add theme color definitions (#228)
