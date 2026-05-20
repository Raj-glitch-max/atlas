package issuance

import (
	"errors"
	"fmt"
	"time"

	"github.com/Raj-glitch-max/atlas/internal/record"
	"github.com/spiffe/go-spiffe/v2/spiffeid"
)

// Authority is the Issuance Authority (M2): the sole creator of delegation
// records and the realization of the issuance boundary (RFC-001 §10.4). It
// enforces the strict-subset scope guard and refuses otherwise, creating
// nothing (SO6, FM6). It holds only injected configuration; it keeps no
// per-request state.

// RefusalCause is the closed set of reasons issuance refused
// (INTERFACE_SPECIFICATION.md §2). The zero value is a non-refusal sentinel.
type RefusalCause int

const (
	// NoRefusal is the sentinel for an issued (non-refused) result.
	NoRefusal RefusalCause = iota
	// MalformedRequest: the request is not well-formed (missing identity,
	// empty scope, empty scope entry, or missing expiration).
	MalformedRequest
	// PermissionsUnavailable: the permission source could not answer the
	// principal's permissions; issuance fails closed (does not guess).
	PermissionsUnavailable
	// OverScope: the requested scope is not a strict subset of the
	// principal's permissions (reaches outside them, or equals them).
	OverScope
)

// String renders the refusal cause.
func (c RefusalCause) String() string {
	switch c {
	case NoRefusal:
		return "NoRefusal"
	case MalformedRequest:
		return "MalformedRequest"
	case PermissionsUnavailable:
		return "PermissionsUnavailable"
	case OverScope:
		return "OverScope"
	default:
		return "RefusalCause(unknown)"
	}
}

// Request is a delegation issuance request.
type Request struct {
	Principal  spiffeid.ID
	Delegate   spiffeid.ID
	Scope      []string
	Expiration time.Time
}

// Outcome is whether a request was issued or refused.
type Outcome int

const (
	Issued Outcome = iota
	Refused
)

// String renders the outcome.
func (o Outcome) String() string {
	if o == Issued {
		return "Issued"
	}
	return "Refused"
}

// Result is the outcome of an issuance attempt. On Issued, Record is set and
// Refusal is NoRefusal; on Refused, Record is nil and Refusal names the cause.
// Trace is always present (both outcomes; RFC-003 §14).
type Result struct {
	Outcome Outcome
	Record  *record.Record
	Refusal RefusalCause
	Trace   IssuanceTrace
}

// Construction refusal causes (closed set).
var (
	ErrNoSigner           = errors.New("issuance: signer key and key ID are required")
	ErrNoPermissionSource = errors.New("issuance: permission source is nil")
	ErrNoRevBindingSource = errors.New("issuance: revocation-binding source is nil")
	ErrNoMinter           = errors.New("issuance: minter is nil")
	ErrNoClock            = errors.New("issuance: clock is nil")
)

// Authority issues delegation records.
type Authority struct {
	signer  record.Signer
	perms   PermissionSource
	binding RevBindingSource
	minter  Minter
	clock   Clock
}

// NewAuthority constructs an Authority, refusing an incomplete signer or any
// nil port: an authority that cannot sign or cannot consult its inputs must
// not exist.
func NewAuthority(signer record.Signer, perms PermissionSource, binding RevBindingSource, minter Minter, clock Clock) (*Authority, error) {
	if signer.Key == nil || signer.KeyID == "" {
		return nil, ErrNoSigner
	}
	if perms == nil {
		return nil, ErrNoPermissionSource
	}
	if binding == nil {
		return nil, ErrNoRevBindingSource
	}
	if minter == nil {
		return nil, ErrNoMinter
	}
	if clock == nil {
		return nil, ErrNoClock
	}
	return &Authority{signer: signer, perms: perms, binding: binding, minter: minter, clock: clock}, nil
}

// Issue attempts to create a delegation record for the request. It returns a
// Result (Issued or Refused, always with a trace) and a non-nil error only on
// a genuine internal fault (minting or sealing failure) — distinct from a
// policy refusal, which is a Result, not an error. On refusal, nothing is
// created: no record, no reserved instance identity, no side effect beyond
// the returned trace (FM6).
func (a *Authority) Issue(req Request) (Result, error) {
	now := a.clock.Now()
	tr := IssuanceTrace{
		Principal:      req.Principal.String(),
		Delegate:       req.Delegate.String(),
		RequestedScope: append([]string(nil), req.Scope...),
		At:             now,
	}

	if cause := validateRequest(req); cause != NoRefusal {
		return refused(cause, tr), nil
	}

	perms, ok := a.perms.PermissionsOf(req.Principal)
	if !ok {
		tr.PermissionsConsulted = false
		return refused(PermissionsUnavailable, tr), nil
	}
	tr.PermissionsConsulted = true

	if !perms.isProperSupersetOf(req.Scope) {
		tr.SubsetSatisfied = false
		return refused(OverScope, tr), nil
	}
	tr.SubsetSatisfied = true

	instance, err := a.minter.Mint()
	if err != nil {
		return Result{}, fmt.Errorf("issuance: minting instance identity: %w", err)
	}

	rec, err := record.Seal(record.Assertions{
		Principal:         req.Principal,
		Delegate:          req.Delegate,
		Scope:             req.Scope,
		Expiration:        req.Expiration,
		IssuedAt:          now,
		Instance:          instance,
		RevocationBinding: a.binding.RevocationBindingFor(instance),
	}, a.signer)
	if err != nil {
		// All inputs were validated, so a seal failure is an internal/config
		// fault, not a policy refusal. Surface it honestly.
		return Result{}, fmt.Errorf("issuance: sealing record: %w", err)
	}

	tr.Outcome = Issued
	tr.Instance = instance.String()
	return Result{Outcome: Issued, Record: rec, Refusal: NoRefusal, Trace: tr}, nil
}

// validateRequest checks well-formedness before any permission consultation.
func validateRequest(req Request) RefusalCause {
	if req.Principal.IsZero() || req.Delegate.IsZero() {
		return MalformedRequest
	}
	if len(req.Scope) == 0 {
		return MalformedRequest
	}
	for _, p := range req.Scope {
		if p == "" {
			return MalformedRequest
		}
	}
	if req.Expiration.IsZero() {
		return MalformedRequest
	}
	return NoRefusal
}

// refused builds a refused Result, completing the trace.
func refused(cause RefusalCause, tr IssuanceTrace) Result {
	tr.Outcome = Refused
	tr.Refusal = cause
	return Result{Outcome: Refused, Record: nil, Refusal: cause, Trace: tr}
}

// checkpoint: fix(revstatus): fix CLI flag configuration (#72)

// checkpoint: feat(sdk): implement signature validation

// checkpoint: chore(test): improve secrets scanner config
