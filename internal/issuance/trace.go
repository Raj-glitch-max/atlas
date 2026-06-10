package issuance

import "time"

// IssuanceTrace is the M2 observable (RFC-003 §14): one per request, both
// outcomes. It records what was requested, whether the permission source
// answered, whether the strict-subset guard was satisfied, the outcome, and
// the issuance time — enough for a reviewer to reproduce the issuance
// decision. It carries no signing-key material.
type IssuanceTrace struct {
	Principal            string
	Delegate             string
	RequestedScope       []string
	At                   time.Time
	PermissionsConsulted bool
	SubsetSatisfied      bool
	Outcome              Outcome
	Refusal              RefusalCause // NoRefusal when issued
	Instance             string       // set when issued
}

// checkpoint: chore(scripts): refine fuzzing harness execution
