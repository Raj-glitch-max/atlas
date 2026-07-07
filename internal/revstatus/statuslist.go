package revstatus

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"sort"
	"sync"
	"time"

	"github.com/Raj-glitch-max/atlas/internal/record"
)

// The signed-revoked-set realization of the Provider contract — the α-path
// composition of the C4 spike (EXP-001), built under the S1–S4 scope act of
// 2026-07-06 (agents/journal/2026-07-06-c4-spike-scope-act.md; S2 admits
// periodic pulls of signed, integrity-protected artifacts; S3 defines a
// passive signed-blob distributor as not-a-broker).
//
// The composition: domain A (the revocation authority) publishes a signed,
// timestamped snapshot of the revoked instance identities as of a time T. A
// relying party ingests the latest such snapshot (out of band, or via a
// passive cache — the S3-compliant channel), verifies its signature against
// domain A's public key it already holds, and answers StatusOf from it. The
// snapshot's signed timestamp is the as-of: freshness is therefore
// **verifiable, not asserted** (OMEGA-04's honest kind — the verifier checks
// a signature over a timestamp rather than trusting a claim). The verifier's
// freshness policy (R) judges whether that as-of is fresh enough; this
// provider applies no policy (dependency rule R8).
//
// This is the reference realization: a signed *set* of revoked IDs, which is
// simple and correct. The OAuth-Status-List bitfield (herd-privacy
// compression, indexed via the record's opaque revocation-binding element,
// AD-015) is a size/privacy optimization deferred with scale (C4 horizon); it
// is semantically this set. The two-domain SPIRE substrate run and the full
// AT13/AT14 measurement remain Epic E6/E7 and are NOT claimed here: this file
// realizes and contract-tests the composition; it does not stand in for the
// substrate-validated spike.

// SignedRevokedSet is a domain-A-signed snapshot of revoked instances as of a
// time. The signature covers the list identity, the as-of, and the (sorted,
// deduplicated) set of revoked identities.
type SignedRevokedSet struct {
	ListID  string
	AsOf    time.Time
	Revoked []record.InstanceID
	Sig     []byte
}

// canonicalDigest is the deterministic hash signed and verified. It sorts and
// deduplicates the revoked set so signer and verifier agree regardless of
// input order or duplicates.
func canonicalDigest(listID string, asOf time.Time, revoked []record.InstanceID) [32]byte {
	sorted := dedupSorted(revoked)
	h := sha256.New()
	h.Write([]byte(listID))
	h.Write([]byte{0})
	var ts [8]byte
	binary.BigEndian.PutUint64(ts[:], uint64(asOf.Unix()))
	h.Write(ts[:])
	h.Write([]byte{0})
	for _, id := range sorted {
		h.Write([]byte(id.String()))
		h.Write([]byte{0})
	}
	var out [32]byte
	copy(out[:], h.Sum(nil))
	return out
}

func dedupSorted(in []record.InstanceID) []record.InstanceID {
	if len(in) == 0 {
		return nil
	}
	cp := append([]record.InstanceID(nil), in...)
	sort.Slice(cp, func(i, j int) bool { return cp[i].String() < cp[j].String() })
	out := cp[:0]
	var prev string
	for i, id := range cp {
		if i > 0 && id.String() == prev {
			continue
		}
		out = append(out, id)
		prev = id.String()
	}
	return out
}

// Publisher is domain A's revocation authority: it signs revoked-set
// snapshots. It holds the private signing key; relying parties verify with
// the corresponding public key.
type Publisher struct {
	key    *ecdsa.PrivateKey
	listID string
}

// Construction refusal causes.
var (
	ErrPublisherNoKey    = errors.New("revstatus: publisher requires a P-256 signing key")
	ErrPublisherNoListID = errors.New("revstatus: publisher requires a non-empty list ID")
)

// NewPublisher builds a Publisher for one list identity.
func NewPublisher(key *ecdsa.PrivateKey, listID string) (*Publisher, error) {
	if key == nil {
		return nil, ErrPublisherNoKey
	}
	if listID == "" {
		return nil, ErrPublisherNoListID
	}
	return &Publisher{key: key, listID: listID}, nil
}

// Publish produces a signed snapshot of the revoked set as of asOf. The
// revoked slice is the authoritative revocation set (e.g. read from the
// Revocation Origin register's View at the composition root — the propagation
// channel; revstatus does not import revorigin, dependency rule R5).
func (p *Publisher) Publish(revoked []record.InstanceID, asOf time.Time) (SignedRevokedSet, error) {
	sorted := dedupSorted(revoked)
	digest := canonicalDigest(p.listID, asOf, sorted)
	sig, err := ecdsa.SignASN1(rand.Reader, p.key, digest[:])
	if err != nil {
		return SignedRevokedSet{}, err
	}
	return SignedRevokedSet{ListID: p.listID, AsOf: asOf, Revoked: sorted, Sig: sig}, nil
}

// SignedSetProvider is the relying-party-side Provider realization. It holds
// the latest verified snapshot and answers StatusOf from it. Concurrency
// (AD-017): single-writer (Ingest) / multi-reader (StatusOf).
type SignedSetProvider struct {
	pub    *ecdsa.PublicKey
	listID string

	mu      sync.RWMutex
	current *SignedRevokedSet
	index   map[record.InstanceID]struct{}
}

// Ingest refusal causes (closed set).
var (
	ErrSetBadSignature = errors.New("revstatus: revoked-set signature does not verify")
	ErrSetWrongList    = errors.New("revstatus: revoked-set list ID does not match")
)

// NewSignedSetProvider builds a provider that trusts snapshots for listID
// signed by pub. Before any successful Ingest it answers Indeterminate for
// every instance (honest ignorance — the system fails closed).
func NewSignedSetProvider(pub *ecdsa.PublicKey, listID string) *SignedSetProvider {
	return &SignedSetProvider{pub: pub, listID: listID}
}

// Ingest verifies a snapshot and adopts it if it is authentic, for the right
// list, and newer than the currently held snapshot. A snapshot that fails
// verification is refused and never adopted — a tampered set can neither
// forge a revocation nor erase one; the provider keeps its prior state (or
// stays Indeterminate). Older-or-equal snapshots are ignored (monotone
// freshness). Returns whether the snapshot was adopted.
func (p *SignedSetProvider) Ingest(s SignedRevokedSet) (adopted bool, err error) {
	if s.ListID != p.listID {
		return false, ErrSetWrongList
	}
	digest := canonicalDigest(s.ListID, s.AsOf, s.Revoked)
	if !ecdsa.VerifyASN1(p.pub, digest[:], s.Sig) {
		return false, ErrSetBadSignature
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.current != nil && !s.AsOf.After(p.current.AsOf) {
		return false, nil // not newer; keep current
	}
	idx := make(map[record.InstanceID]struct{}, len(s.Revoked))
	for _, id := range s.Revoked {
		idx[id] = struct{}{}
	}
	snapshot := s
	p.current = &snapshot
	p.index = idx
	return true, nil
}

// StatusOf answers from the held snapshot (Provider contract). With a verified
// snapshot the provider is NOT ignorant about any instance — the snapshot is a
// complete revocation census as of its signed as-of — so an instance absent
// from it is honestly NotObservedRevoked(asOf), and one present is
// ObservablyRevoked(asOf). With no snapshot ingested, every instance is
// Indeterminate (honest ignorance; the honest-indeterminate rule holds).
func (p *SignedSetProvider) StatusOf(instance record.InstanceID) Answer {
	p.mu.RLock()
	defer p.mu.RUnlock()
	if p.current == nil {
		return indeterminate()
	}
	if _, revoked := p.index[instance]; revoked {
		return Answer{State: ObservablyRevoked, AsOf: p.current.AsOf}
	}
	return Answer{State: NotObservedRevoked, AsOf: p.current.AsOf}
}
