package revorigin

import (
	"sync"
	"time"

	"github.com/Raj-glitch-max/atlas/internal/record"
)

// Register is the authoritative, append-only revocation register (M6). It
// records the one-way, terminal revocation of specific delegation instances
// (INV4) and exposes its ordered content read-only for the deferred
// propagation channel and after-the-fact reconstruction.
//
// It holds only opaque instance identities and is therefore structurally
// incapable of affecting underlying identities (INV5) or sibling delegations
// (INV6): a revocation names one instance and touches nothing else. Revoking
// a never-issued instance is inert — verification keys its checks to
// presented records, so a spurious entry matches nothing.
//
// Concurrency (AD-017): single-writer semantics for appends, multi-reader for
// View, guarded by a mutex; the register serializes writes and returns
// independent snapshots to readers.
//
// V1 holds the register in memory (persistence deferred, TD-1) behind this
// same interface; the read/write surface does not change when a durable
// substrate lands.
type Register struct {
	mu      sync.RWMutex
	order   []record.InstanceID
	revoked map[record.InstanceID]Entry
}

// Entry is one revocation record: which instance, and when it was first
// recorded (append order is authoritative; the timestamp is informational
// for the propagation channel and reconstruction).
type Entry struct {
	Instance   record.InstanceID
	RecordedAt time.Time
}

// New returns an empty register.
func New() *Register {
	return &Register{revoked: make(map[record.InstanceID]Entry)}
}

// Revoke records the terminal revocation of one instance. It is idempotent on
// an already-revoked instance (a no-op, not an error): revocation is one-way
// and terminal (INV4), so a repeat records nothing new and the original
// RecordedAt stands. The now argument is the recording time (injected, not
// read from a wall clock, for determinism and testability).
func (r *Register) Revoke(instance record.InstanceID, now time.Time) {
	if instance.IsZero() {
		// A zero instance identifies nothing; recording it would create an
		// entry that can never match a presented record. Reject silently by
		// no-op — there is nothing to revoke.
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.revoked[instance]; exists {
		return // terminal; no-op
	}
	r.revoked[instance] = Entry{Instance: instance, RecordedAt: now}
	r.order = append(r.order, instance)
}

// IsRevoked reports whether the register holds a revocation for the instance.
// This is the authoritative fact on the issuing side; the relying-party view
// (M5) is a possibly-stale observation of it, and the gap between them is the
// revocation-observability state (bounded by R and S4).
func (r *Register) IsRevoked(instance record.InstanceID) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, ok := r.revoked[instance]
	return ok
}

// View returns the complete register in append order, as an independent
// snapshot: no operation removes or rewrites entries, and a caller cannot
// mutate the register through the returned slice. This is the surface the
// deferred propagation channel (S2/S3, spike-selected) and reviewers read.
func (r *Register) View() []Entry {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]Entry, len(r.order))
	for i, inst := range r.order {
		out[i] = r.revoked[inst]
	}
	return out
}

// Len returns the number of revocations recorded.
func (r *Register) Len() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.order)
}

// checkpoint: feat(issuance): implement boundary check (#96)
