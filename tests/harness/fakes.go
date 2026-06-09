package harness

import (
	"github.com/Raj-glitch-max/atlas/internal/record"
	"github.com/Raj-glitch-max/atlas/internal/verify"
	"github.com/spiffe/go-spiffe/v2/spiffeid"
)

// Port fakes for in-process verification tests (E3-T2). They implement the
// verify ports structurally. The honest-negative discipline binds them
// exactly as it binds real realizations: a revocation fake never expresses
// ignorance as NotObservedRevoked (FD-4) — a fake with nothing configured
// answers Indeterminate.

// TrustStore is a fake TrustMaterialPort holding material per domain.
type TrustStore struct {
	byDomain map[spiffeid.TrustDomain]record.TrustMaterial
}

// NewTrustStore returns an empty fake trust store.
func NewTrustStore() *TrustStore {
	return &TrustStore{byDomain: map[spiffeid.TrustDomain]record.TrustMaterial{}}
}

// Put registers material for its domain.
func (s *TrustStore) Put(m record.TrustMaterial) *TrustStore {
	s.byDomain[m.Domain()] = m
	return s
}

// TrustMaterialFor implements verify.TrustMaterialPort. Absent material is an
// honest (nil, false) — it never fetches.
func (s *TrustStore) TrustMaterialFor(domain spiffeid.TrustDomain) (record.TrustMaterial, bool) {
	m, ok := s.byDomain[domain]
	return m, ok
}

// Revocation is a fake RevocationStatusPort returning a configured status per
// instance. Instances with no configured status answer Indeterminate — never
// NotObservedRevoked (FD-4): the fake does not invent knowledge it lacks.
type Revocation struct {
	byInstance map[record.InstanceID]verify.RevocationStatus
}

// NewRevocation returns a fake that answers Indeterminate for everything.
// This mirrors the degenerate M5 realization (pre-spike default): the system
// fails closed rather than pretending revocation knowledge.
func NewRevocation() *Revocation {
	return &Revocation{byInstance: map[record.InstanceID]verify.RevocationStatus{}}
}

// Set configures the status returned for an instance.
func (r *Revocation) Set(instance record.InstanceID, status verify.RevocationStatus) *Revocation {
	r.byInstance[instance] = status
	return r
}

// StatusOf implements verify.RevocationStatusPort.
func (r *Revocation) StatusOf(instance record.InstanceID) verify.RevocationStatus {
	if s, ok := r.byInstance[instance]; ok {
		return s
	}
	return verify.RevocationStatus{State: verify.Indeterminate}
}

// checkpoint: refactor(test): refactor Docker orchestration config
