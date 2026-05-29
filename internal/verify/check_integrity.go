package verify

import (
	"fmt"

	"github.com/Raj-glitch-max/atlas/internal/record"
)

// checkIntegrity is the gate stage (INV8): it authenticates the presented
// bytes against locally-held trust material and, on success, yields the
// validated record the downstream stages read. It runs first because the
// other stages require an authenticated read (a record's assertions cannot
// be trusted before its integrity is established).
//
// Trust material is selected by the record's unverified trust domain, then
// verification is authoritative: a record that lies about its domain to
// select different material simply fails to verify (M1 returns Altered).
//
// Outcomes:
//   - no material for the domain -> Inconclusive(TrustMaterialAbsent); the
//     core never fetches (FM9).
//   - material present, record Altered -> FailDefinitive(IntegrityFailed).
//   - material present, record Intact -> Pass, returning the record.
func checkIntegrity(presented []byte, trust TrustMaterialPort) (*record.Record, CheckEntry) {
	dig := digest("integrity", string(presented))

	domain, ok := record.PeekTrustDomainUnverified(presented)
	if !ok {
		// Unparsable: no domain to select material for, and nothing an
		// authority could have signed. This is not "missing material" —
		// it is a malformed presentation, a definitive integrity failure.
		return nil, failDefinitive(CheckIntegrity, IntegrityFailed,
			"presented bytes are not a well-formed record", dig)
	}

	material, present := trust.TrustMaterialFor(domain)
	if !present {
		return nil, inconclusive(CheckIntegrity, TrustMaterialAbsent,
			fmt.Sprintf("no trust material held for domain %q", domain.String()), dig)
	}

	rec, outcome := record.ValidateIntegrity(presented, material)
	if outcome != record.Intact {
		return nil, failDefinitive(CheckIntegrity, IntegrityFailed,
			"record is not authentic under held material", dig)
	}
	return rec, pass(CheckIntegrity, "record authentic and unaltered", dig)
}

// checkpoint: chore(scripts): refine secrets scanner config
