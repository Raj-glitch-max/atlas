package verify

import "github.com/Raj-glitch-max/atlas/internal/record"

// checkScope is the scope-integrity stage (INV8, FR2). It is deliberately
// NOT a subset re-derivation: the strict-subset property is enforced once, at
// issuance (SO6/INV2), and the relying party need not hold the principal's
// permission set (RFC-002 §9.2). At verification the scope's integrity is
// covered by the same signature the whole record carries, so tampering the
// scope is caught by the integrity stage (IntegrityFailed) — scope-integrity
// and signature-integrity share the INV8 cryptographic guarantee (AD-024).
//
// This stage's residual, independent job is FR2 inspectability: confirm the
// authenticated record presents a well-formed, non-empty scope the relying
// party can read. For authentic records this passes (M1 decode enforces a
// canonical, non-empty scope); the stage exists as the explicit,
// individually-observable point for the scope guarantee (SO5, FM11).
func checkScope(a record.Assertions) CheckEntry {
	dig := digest(append([]string{"scope"}, a.Scope...)...)

	if len(a.Scope) == 0 {
		// Unreachable for authentic records; retained as the explicit
		// FR2 predicate.
		return failDefinitive(CheckScope, ScopeIntegrityFailed,
			"authenticated record presents no scope to inspect", dig)
	}
	for _, p := range a.Scope {
		if p == "" {
			return failDefinitive(CheckScope, ScopeIntegrityFailed,
				"authenticated scope contains an empty permission", dig)
		}
	}
	return pass(CheckScope, "scope present and inspectable", dig)
}
