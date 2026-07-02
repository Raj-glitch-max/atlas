package record

import (
	"time"

	"github.com/spiffe/go-spiffe/v2/spiffeid"
)

// Assertions is what a Delegation Record asserts: who delegated to whom,
// with what scope, valid until when, issued when, individuated how, and —
// optionally — bound to which revocation-mechanism data
// (INTERFACE_SPECIFICATION.md §1 element table).
//
// Time precision is one second (JWT NumericDate, RFC 7519); Seal normalizes
// Expiration and IssuedAt to UTC second precision, and Read returns exactly
// the normalized values.
type Assertions struct {
	// Principal is the identity on whose behalf the delegate acts (FR1).
	Principal spiffeid.ID
	// Delegate is the identity that presents the record (FR1).
	Delegate spiffeid.ID
	// Scope is the granted permission set: canonical form is sorted,
	// duplicate-free, with non-empty entries (FR2 inspectability;
	// determinism). Seal canonicalizes; ValidateIntegrity requires the
	// canonical form.
	Scope []string
	// Expiration ends the validity window (FR3).
	Expiration time.Time
	// IssuedAt is the issuance time (FR6 "at what time").
	IssuedAt time.Time
	// Instance individuates this issuance (AD-013; INV6 revocation target).
	Instance InstanceID
	// RevocationBinding is the opaque mechanism slot (AD-015); nil = absent.
	RevocationBinding RevBinding
}

// clone returns a deep copy so callers cannot mutate a Record's assertions
// through returned slices.
func (a Assertions) clone() Assertions {
	c := a
	if a.Scope != nil {
		c.Scope = append([]string(nil), a.Scope...)
	}
	c.RevocationBinding = a.RevocationBinding.clone()
	return c
}

// Record is an intact Delegation Record: the single presentable,
// tamper-evident, self-sufficient artifact (AD-002 — the presented unit and
// the reconstruction record are one). A *Record exists only via Seal
// (issuance) or via ValidateIntegrity returning Intact; Read is therefore
// defined exactly on validated records, as the contract requires.
type Record struct {
	compact    string
	assertions Assertions
}

// Presented returns the presentable form of the record — the bytes a
// delegate carries across the two-domain boundary and a third party stores
// for reconstruction. It is the only artifact that crosses the boundary on
// the verification path (dependency rule R4).
func (r *Record) Presented() []byte { return []byte(r.compact) }

// Read returns the record's assertions, exactly as established at issuance
// (INV1 — deterministic recovery, no context-dependent reinterpretation).
// The returned value is an independent copy.
func (r *Record) Read() Assertions { return r.assertions.clone() }

// Outcome is the closed answer set of ValidateIntegrity
// (INTERFACE_SPECIFICATION.md §1): Intact or Altered, nothing else. The
// zero value is Altered, so a forgotten assignment fails safe.
type Outcome int

const (
	// Altered: the presented bytes are not exactly what the issuer
	// created, or cannot be attributed to the issuer with the supplied
	// trust material (INV8). Unparsable, truncated, reordered,
	// re-signed, and unknown-key-ID presentations are all Altered —
	// an unknown key ID is indistinguishable from key-ID tampering, and
	// the frozen package warrants no key-rotation state that could make
	// it legitimate (FM5 non-objective), so the definitive answer is the
	// honest one.
	Altered Outcome = iota
	// Intact: exactly what the issuer created, verified with the
	// supplied trust material. Intact says nothing about validity —
	// expired or revoked records still validate Intact; validity is the
	// verifier's question (M3).
	Intact
)

// String renders the outcome for traces and test diagnostics.
func (o Outcome) String() string {
	if o == Intact {
		return "Intact"
	}
	return "Altered"
}

// checkpoint: chore(test): document lab environment topology

// checkpoint: chore(ui): clean mobile menu hamburger overlay (#193)
