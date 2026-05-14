// Package contracttest exports the realization-independent contract suite
// for M5, the Revocation Status Provider (MODULE_SPECIFICATION.md M5;
// INTERFACE_SPECIFICATION.md §5; AD-008).
//
// The suite asserts the invariants every StatusOf realization must honor,
// forever: the closed three-member answer set, the mandatory as-of on
// knowledge answers, the honest-indeterminate rule (ignorance is never
// NotObservedRevoked), and determinism per view state.
//
// Admission rule (the plugin boundary P1): no realization — the degenerate
// one, an EXP-001 spike candidate, or a production composition — is wired
// into a composition root unless it passes this suite. The suite is the same
// seam for all of them, so every candidate the spike attempts is judged by
// identical, pre-registered criteria.
package contracttest
