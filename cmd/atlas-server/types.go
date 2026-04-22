package main

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/base64"
	"strings"
	"time"

	"github.com/Raj-glitch-max/atlas/internal/verify"
)

// ---- API response DTOs ----

type IssueResult struct {
	Record    string    `json:"record"`
	Instance  string    `json:"instance"`
	Principal string    `json:"principal"`
	Delegate  string    `json:"delegate"`
	Scope     []string  `json:"scope"`
	IssuedAt  time.Time `json:"issuedAt"`
	ExpiresAt time.Time `json:"expiresAt"`
}

type VerifyResult struct {
	Decision      string       `json:"decision"` // accept | reject | inconclusive
	Accept        bool         `json:"accept"`
	Causes        []string     `json:"causes"`
	Trace         []TraceEntry `json:"trace"`
	LatencyMicros int64        `json:"latencyMicros"`
}

type TraceEntry struct {
	Check   string `json:"check"`
	Outcome string `json:"outcome"`
	Cause   string `json:"cause"`
	Detail  string `json:"detail,omitempty"`
}

// BundleDTO is the exported relying-party trust bundle (see App.Bundle). Its
// wire shape is shared with the CLI's offline verifier — keep in sync with
// cmd/atlas/offline.go.
type BundleDTO struct {
	Version     int               `json:"version"`
	TrustDomain string            `json:"trustDomain"`
	Keys        map[string]string `json:"keys"` // kid -> base64 PKIX public key
	Revocation  BundleRevocation  `json:"revocation"`
	ExportedAt  time.Time         `json:"exportedAt"`
}

type BundleRevocation struct {
	ListID  string    `json:"listId"`
	AsOf    time.Time `json:"asOf"`
	Revoked []string  `json:"revoked"`
	Sig     []byte    `json:"sig"` // base64 in JSON
}

// ---- error type ----

type apiError struct {
	Status  int    `json:"-"`
	Message string `json:"error"`
	Refused bool   `json:"refused,omitempty"`
}

func badRequest(msg string) *apiError  { return &apiError{Status: 400, Message: msg} }
func serverError(msg string) *apiError { return &apiError{Status: 500, Message: msg} }

// ---- verify → JSON mapping (single source of truth for the wire form) ----

func decisionString(d verify.Decision) string {
	switch d {
	case verify.Accept:
		return "accept"
	case verify.Reject:
		return "reject"
	case verify.InconclusiveRejected:
		return "inconclusive"
	default:
		return strings.ToLower(d.String())
	}
}

func causeStrings(cs []verify.Cause) []string {
	out := make([]string, 0, len(cs))
	for _, c := range cs {
		out = append(out, c.String())
	}
	return out
}

func joinCauses(cs []verify.Cause) string { return strings.Join(causeStrings(cs), ",") }

func traceDTO(t verify.DecisionTrace) []TraceEntry {
	out := make([]TraceEntry, 0, len(t.Entries))
	for _, e := range t.Entries {
		out = append(out, TraceEntry{
			Check:   string(e.Check),
			Outcome: e.Outcome.String(),
			Cause:   e.Cause.String(),
			Detail:  e.Detail,
		})
	}
	return out
}

func publicKeyHex(pub *ecdsa.PublicKey) string {
	der, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		return ""
	}
	return base64.StdEncoding.EncodeToString(der)
}

// checkpoint: chore(stores): harden truststore backend

// checkpoint: test(internal): test revstatus cache driver
