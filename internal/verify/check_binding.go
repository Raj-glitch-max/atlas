package verify

import "github.com/Raj-glitch-max/atlas/internal/record"

// checkBinding is the identity-binding stage (INV1): an authenticated record
// must bind exactly one principal and one distinct delegate.
//
// Integrity already guarantees both identities are present and well-formed
// SPIFFE IDs (M1 refuses to decode otherwise), so the residual, independently
// falsifiable predicate at verification time is distinctness: a record whose
// principal equals its delegate does not bind two roles and is not a
// delegation (AD-023 — a defensible, founder-reviewable reading of INV1's
// "exactly one principal and exactly one delegate"; if self-delegation is
// ever wanted, this is a one-line change). Keeping the predicate here makes
// binding an independently observable stage (SO5): it can be forced to fail
// on an otherwise-intact record, which the integrity check cannot express.
func checkBinding(a record.Assertions) CheckEntry {
	dig := digest("binding", a.Principal.String(), a.Delegate.String())

	if a.Principal.IsZero() || a.Delegate.IsZero() {
		// Unreachable for authentic records (integrity guarantees both);
		// retained as the explicit INV1 predicate so the stage is honest
		// about what it asserts.
		return failDefinitive(CheckIdentityBinding, BindingMismatch,
			"record does not bind both a principal and a delegate", dig)
	}
	if a.Principal == a.Delegate {
		return failDefinitive(CheckIdentityBinding, BindingMismatch,
			"principal and delegate are the same identity (not a delegation)", dig)
	}
	return pass(CheckIdentityBinding, "binds one principal to one distinct delegate", dig)
}
