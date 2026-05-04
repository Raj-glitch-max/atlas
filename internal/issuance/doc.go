// Package issuance is M2, the Issuance Authority: the sole creator of
// Delegation Records and the realization of the issuance trust boundary
// (RFC-003 §7 M2; MODULE_SPECIFICATION.md M2; INTERFACE_SPECIFICATION.md §2).
//
// Issue enforces the strict-subset scope guard — a requested scope must be a
// proper subset of the principal's permission set, obtained through the
// PermissionSource port — and refuses otherwise, creating nothing (SO6, FM6).
// An issued record carries a fresh opaque instance identity from the minter
// seam (AD-013), the issuance-time reading from the time port (AD-014), and
// the revocation-binding element from the RevBindingSource port, which
// answers none until the revocation mechanism exists (AD-015).
//
// The package consumes ports and internal/record only; it never imports
// verification- or relying-party-side packages, and issuance to ephemeral
// delegates requires no long-lived, statically-provisioned identity (ER10).
//
// Every request — issued or refused — produces an IssuanceTrace, returned to
// the invoker (AD-016).
//
// Traces: ER1, ER2, ER3, ER10, INV1, INV2, INV3, SO6, FM6, AP10.
package issuance
