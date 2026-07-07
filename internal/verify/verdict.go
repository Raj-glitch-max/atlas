package verify

// Decision is the closed set of verification outcomes (RFC-002 §9.2). The
// zero value is Accept only nominally — a Verdict is never constructed by
// zero value; it is always produced by routing (verifier.go), which yields
// Accept only when every check passed.
type Decision int

const (
	// Accept: every required check passed.
	Accept Decision = iota
	// Reject: at least one check failed definitively.
	Reject
	// InconclusiveRejected is the fail-closed routing [HYPOTHESIS]: no
	// definitive failure, but at least one check could not conclude. It IS
	// a rejection — kept distinct from Reject so the designed fail-closed
	// behavior (NFR3/ER11/SO4/C-INV1) is observable and testable (AT22)
	// without being conflated with a definitive reject. This distinction
	// does not promote the hypothesis (DR7): the verdict records what
	// happened; only a V1 confirmation act warrants it.
	InconclusiveRejected
)

// String renders the decision for traces and diagnostics.
func (d Decision) String() string {
	switch d {
	case Accept:
		return "Accept"
	case Reject:
		return "Reject"
	case InconclusiveRejected:
		return "InconclusiveRejected[HYPOTHESIS]"
	default:
		return "Decision(unknown)"
	}
}

// Verdict is the outcome of a verification: a decision and the causes that
// produced it. Accept carries no causes; Reject carries the definitive
// cause(s); InconclusiveRejected carries the inconclusive cause(s).
type Verdict struct {
	Decision Decision
	Causes   []Cause
}

// IsAccept reports whether the verdict accepted the delegation. Any other
// decision — Reject or InconclusiveRejected — is a non-acceptance; callers
// that must "accept or reject" treat both as reject (fail closed).
func (v Verdict) IsAccept() bool { return v.Decision == Accept }
