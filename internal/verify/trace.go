package verify

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"time"
)

// The decision trace is the Verification Core's only output channel besides
// the verdict, emitted unconditionally — Accepts included (AP13, SO5, SO8).
// It records which checks fired, in what order, with what inputs (digested,
// never leaking material) and outcome, the injected time and freshness
// readings, and the policy in force — enough for an independent reviewer to
// reproduce the verdict from the specification and a build, and nothing a
// reviewer may not hold (AP8). Observation never alters a verdict.

// CheckName names a pipeline stage.
type CheckName string

const (
	CheckIdentityBinding CheckName = "identity_binding"
	CheckIntegrity       CheckName = "integrity"
	CheckExpiry          CheckName = "expiry"
	CheckScope           CheckName = "scope_integrity"
	CheckRevocation      CheckName = "revocation_status"
)

// CheckOutcome is the per-stage result. NotEvaluated is honest, not a gap: a
// stage that depends on an authenticated read cannot run when integrity did
// not pass, and the trace says so rather than silently omitting it.
type CheckOutcome int

const (
	OutcomePass CheckOutcome = iota
	OutcomeFailDefinitive
	OutcomeInconclusive
	OutcomeNotEvaluated
)

// String renders the outcome for diagnostics.
func (o CheckOutcome) String() string {
	switch o {
	case OutcomePass:
		return "Pass"
	case OutcomeFailDefinitive:
		return "FailDefinitive"
	case OutcomeInconclusive:
		return "Inconclusive"
	case OutcomeNotEvaluated:
		return "NotEvaluated"
	default:
		return "CheckOutcome(unknown)"
	}
}

// CheckEntry is one stage's trace record.
type CheckEntry struct {
	Check        CheckName
	Outcome      CheckOutcome
	Cause        Cause  // CauseNone when the stage passed or was not evaluated
	Detail       string // human-readable, no secrets
	InputsDigest string // digest of the stage's inputs (reproducibility, no material)
}

// PolicySummary is the trace-visible view of the policy in force.
type PolicySummary struct {
	R             time.Duration
	SkewTolerance time.Duration
}

// DecisionTrace is the complete, reproducible account of one verification.
type DecisionTrace struct {
	// Entries are the five stages in canonical label order (binding,
	// integrity, expiry, scope, revocation). Execution runs integrity
	// first as the gate (see verifier.go); the trace presents canonical
	// order for stable reading.
	Entries                []CheckEntry
	TimeReading            time.Time
	RevocationState        string
	RevocationObservedAsOf time.Time
	Policy                 PolicySummary
	Verdict                Verdict
}

// Find returns the entry for a named check, and whether it was present.
func (t DecisionTrace) Find(name CheckName) (CheckEntry, bool) {
	for _, e := range t.Entries {
		if e.Check == name {
			return e, true
		}
	}
	return CheckEntry{}, false
}

// digest produces a stable hex digest of a stage's inputs for reproducibility.
// Inputs are joined with a separator that cannot appear in the parts' own
// structure so the concatenation is unambiguous.
func digest(parts ...string) string {
	h := sha256.Sum256([]byte(strings.Join(parts, "\x1e")))
	return hex.EncodeToString(h[:])
}

// pass, failDefinitive, inconclusive, notEvaluated build entries uniformly so
// every stage records the same shape.
func pass(check CheckName, detail, dig string) CheckEntry {
	return CheckEntry{Check: check, Outcome: OutcomePass, Cause: CauseNone, Detail: detail, InputsDigest: dig}
}

func failDefinitive(check CheckName, cause Cause, detail, dig string) CheckEntry {
	return CheckEntry{Check: check, Outcome: OutcomeFailDefinitive, Cause: cause, Detail: detail, InputsDigest: dig}
}

func inconclusive(check CheckName, cause Cause, detail, dig string) CheckEntry {
	return CheckEntry{Check: check, Outcome: OutcomeInconclusive, Cause: cause, Detail: detail, InputsDigest: dig}
}

func notEvaluated(check CheckName, detail string) CheckEntry {
	return CheckEntry{Check: check, Outcome: OutcomeNotEvaluated, Cause: CauseNone, Detail: detail}
}
