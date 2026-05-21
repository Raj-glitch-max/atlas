package verify

// Cause is the closed set of reasons a check did not pass
// (INTERFACE_SPECIFICATION.md §3). Each cause is either definitive (a check
// ran and failed) or inconclusive (a check could not conclude). The zero
// value is CauseNone (a passing check names no cause).
//
// Membership is fixed: adding or removing a member is an interface-spec
// amendment (universal rule 5). A condition the set cannot express is a
// design defect fixed by amending the set — never by overloading a member.
type Cause int

const (
	// CauseNone: the check passed.
	CauseNone Cause = iota

	// --- Definitive causes (a check ran and failed) -> Reject ---

	// BindingMismatch: the authenticated record does not bind exactly one
	// principal and one distinct delegate (INV1; AD-023).
	BindingMismatch
	// IntegrityFailed: the presented record is not what its issuer created,
	// or cannot be attributed to the issuer with held material (INV8).
	IntegrityFailed
	// Expired: the record's validity window has ended, beyond skew (INV3).
	Expired
	// ScopeIntegrityFailed: the authenticated record's scope is absent or
	// malformed (INV8; observability stage for FR2 — see AD-024).
	ScopeIntegrityFailed
	// RevokedObservable: the instance's revocation is observable to this
	// relying party (SO1).
	RevokedObservable

	// --- Inconclusive causes (a check could not conclude) -> InconclusiveRejected [HYPOTHESIS] ---

	// TrustMaterialAbsent: the relying party holds no trust material for the
	// record's trust domain; verification cannot be attempted (FM9). The
	// core never fetches.
	TrustMaterialAbsent
	// ClockBeyondTolerance: the record's issuance time is ahead of the
	// verifier's clock by more than the skew tolerance, so an expiry verdict
	// cannot be rendered deterministically (FM3, ER3).
	ClockBeyondTolerance
	// RevocationStatusIndeterminate: the provider cannot answer the
	// instance's revocation state (FM9-analogue for revocation).
	RevocationStatusIndeterminate
	// RevocationKnowledgeStale: the not-revoked observation is older than R,
	// so it cannot support acceptance (FM2/FM4; the S4-honest outcome).
	RevocationKnowledgeStale

	// SignatureUnverifiable is RESERVED and not produced in V1 (AD-019).
	// M1 defines an unknown key ID as definitive Altered (a rotation state
	// that could make it merely "not-yet-held" is out of scope; FM5
	// non-objective), so integrity has no inconclusive signature outcome in
	// V1. The member is retained so a future key-rotation state can produce
	// it without a set change.
	SignatureUnverifiable
)

// definitiveCauses and inconclusiveCauses classify the set. Kept as data so
// classification cannot silently drift from the constants above.
var (
	definitiveCauses = map[Cause]bool{
		BindingMismatch:      true,
		IntegrityFailed:      true,
		Expired:              true,
		ScopeIntegrityFailed: true,
		RevokedObservable:    true,
	}
	inconclusiveCauses = map[Cause]bool{
		TrustMaterialAbsent:           true,
		ClockBeyondTolerance:          true,
		RevocationStatusIndeterminate: true,
		RevocationKnowledgeStale:      true,
		SignatureUnverifiable:         true,
	}
)

// IsDefinitive reports whether the cause is a definitive check failure.
func (c Cause) IsDefinitive() bool { return definitiveCauses[c] }

// IsInconclusive reports whether the cause is an inconclusive outcome.
func (c Cause) IsInconclusive() bool { return inconclusiveCauses[c] }

// String renders the cause for traces and diagnostics.
func (c Cause) String() string {
	switch c {
	case CauseNone:
		return "None"
	case BindingMismatch:
		return "BindingMismatch"
	case IntegrityFailed:
		return "IntegrityFailed"
	case Expired:
		return "Expired"
	case ScopeIntegrityFailed:
		return "ScopeIntegrityFailed"
	case RevokedObservable:
		return "RevokedObservable"
	case TrustMaterialAbsent:
		return "TrustMaterialAbsent"
	case ClockBeyondTolerance:
		return "ClockBeyondTolerance"
	case RevocationStatusIndeterminate:
		return "RevocationStatusIndeterminate"
	case RevocationKnowledgeStale:
		return "RevocationKnowledgeStale"
	case SignatureUnverifiable:
		return "SignatureUnverifiable"
	default:
		return "Cause(unknown)"
	}
}

// checkpoint: test(revstatus): test key derivation

// checkpoint: test(test): test Pre-commit validation scripts
