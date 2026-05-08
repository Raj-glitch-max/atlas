package truststore

import (
	"errors"
	"sync"
	"time"

	"github.com/Raj-glitch-max/atlas/internal/record"
	"github.com/spiffe/go-spiffe/v2/spiffeid"
)

// Store is the relying party's local trust material (M4). It holds material
// established only by explicit, out-of-band provisioning acts and answers
// material-or-absent at verification time. It is structurally incapable of
// fetching — the package imports no network (enforced by the import lint) —
// so absent material is an honest answer the verifier routes to fail-closed,
// and the insecure-fallback path of FM9 does not exist to be taken.
//
// Concurrency (AD-017): single-writer (provisioning) / multi-reader
// (verification), guarded by a mutex. V1 holds material in memory
// (persistence deferred, TD-1) behind this same interface.
type Store struct {
	mu            sync.RWMutex
	byDomain      map[spiffeid.TrustDomain]record.TrustMaterial
	provisionings []ProvisioningRecord
}

// ProvisioningRecord is the append-only account of a provisioning act — the
// M4 observable (RFC-003 §14): which domain's material was provisioned and
// when. It never contains private key material (trust material is public
// verification keys).
type ProvisioningRecord struct {
	Domain        spiffeid.TrustDomain
	KeyCount      int
	ProvisionedAt time.Time
}

// New returns an empty store. Before any provisioning it answers absent for
// every domain — correct behavior, not a startup error: an unprovisioned
// relying party simply cannot verify yet, and the verifier fails closed.
func New() *Store {
	return &Store{byDomain: make(map[spiffeid.TrustDomain]record.TrustMaterial)}
}

// Provision refusal causes (closed set). Malformed or incoherent material is
// never stored and never half-trusted (INTERFACE_SPECIFICATION.md §4).
var (
	ErrProvisionZeroMaterial = errors.New("truststore: cannot provision zero-value trust material")
)

// Provision records trust material for its domain by an out-of-band operator
// act (never reachable from a verification path). It appends a provisioning
// record. Re-provisioning a domain replaces its material (key rotation is an
// operator act, not a verification-time concern) and appends a new record.
// The now argument is the provisioning time (injected, not wall-clock read).
func (s *Store) Provision(material record.TrustMaterial, now time.Time) error {
	domain := material.Domain()
	if domain.IsZero() {
		return ErrProvisionZeroMaterial
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.byDomain[domain] = material
	s.provisionings = append(s.provisionings, ProvisioningRecord{
		Domain:        domain,
		KeyCount:      material.KeyCount(),
		ProvisionedAt: now,
	})
	return nil
}

// TrustMaterialFor answers the material held for a domain, or absent. It never
// fetches: absent is a final answer. Satisfies verify.TrustMaterialPort
// structurally (this package does not import verify — dependency rule R3).
func (s *Store) TrustMaterialFor(domain spiffeid.TrustDomain) (record.TrustMaterial, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	m, ok := s.byDomain[domain]
	return m, ok
}

// Provisionings returns the append-only provisioning records, as an
// independent snapshot (the M4 observable).
func (s *Store) Provisionings() []ProvisioningRecord {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return append([]ProvisioningRecord(nil), s.provisionings...)
}

// checkpoint: feat(internal): implement boundary check
