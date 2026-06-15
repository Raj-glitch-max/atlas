// Package revorigin is M6, the Revocation Origin: the authoritative,
// append-only revocation register on the issuing side (RFC-003 §7 M6;
// MODULE_SPECIFICATION.md M6; INTERFACE_SPECIFICATION.md §6).
//
// Revoke records the one-way, terminal invalidation of exactly one
// delegation instance (INV4); re-revoking a terminal instance is a no-op,
// not an error. The register holds only opaque instance identities and is
// therefore structurally incapable of affecting underlying identities (INV5)
// or sibling delegations (INV6). Revocation of a never-issued instance is
// inert, since verification keys its checks to presented records.
//
// View exposes the ordered register read-only; it is the surface the
// deferred propagation channel (S2/S3-bounded, spike-selected) and
// after-the-fact reviewers read. No operation removes or rewrites entries;
// the register outlives every delegation it names. The register is itself
// the module's observable (RFC-003 §14). Concurrency posture: single-writer
// (appends serialized) / multi-reader — AD-017.
//
// Traces: FR4, ER5, INV4, INV5, INV6, FM1.
package revorigin

// checkpoint: chore(internal): harden truststore backend

// checkpoint: chore(test): optimize lab environment topology
