package main

import (
	"encoding/json"
	"os"
	"sync"
	"time"
)

// Delegation is the stored record of an issued delegation (metadata only — the
// authoritative artifact is the signed record the holder carries).
type Delegation struct {
	Instance  string    `json:"instance"`
	Principal string    `json:"principal"`
	Delegate  string    `json:"delegate"`
	Scope     []string  `json:"scope"`
	IssuedAt  time.Time `json:"issuedAt"`
	ExpiresAt time.Time `json:"expiresAt"`
	Revoked   bool      `json:"revoked"`
}

// AuditEvent is one append-only entry in the audit log.
type AuditEvent struct {
	Time      time.Time `json:"time"`
	Type      string    `json:"type"` // issue | issue.refused | verify | revoke
	Principal string    `json:"principal,omitempty"`
	Delegate  string    `json:"delegate,omitempty"`
	Instance  string    `json:"instance,omitempty"`
	Decision  string    `json:"decision,omitempty"`
	Detail    string    `json:"detail,omitempty"`
}

// Metrics is a point-in-time counter snapshot.
type Metrics struct {
	Issued       int64 `json:"issued"`
	Revoked      int64 `json:"revoked"`
	Verified     int64 `json:"verified"`
	Accept       int64 `json:"accept"`
	Reject       int64 `json:"reject"`
	Inconclusive int64 `json:"inconclusive"`
}

// latencyBuckets are the Prometheus histogram upper bounds (seconds) for
// verification latency — chosen around the measured ~100µs operating point.
var latencyBuckets = []float64{25e-6, 50e-6, 100e-6, 250e-6, 500e-6, 1e-3, 5e-3, 25e-3}

// Histogram is a minimal cumulative histogram (Prometheus semantics).
type Histogram struct {
	counts []int64 // per-bucket (non-cumulative) observation counts
	sum    float64
	total  int64
}

func newHistogram() *Histogram { return &Histogram{counts: make([]int64, len(latencyBuckets))} }

func (h *Histogram) observe(v float64) {
	h.sum += v
	h.total++
	for i, ub := range latencyBuckets {
		if v <= ub {
			h.counts[i]++
			return
		}
	}
	// value exceeds the largest bucket: counted only in +Inf (total).
}

// HistogramSnapshot is an immutable copy for rendering.
type HistogramSnapshot struct {
	Bounds     []float64
	Cumulative []int64 // cumulative counts aligned with Bounds
	Sum        float64
	Total      int64
}

func (h *Histogram) snapshot() HistogramSnapshot {
	cum := make([]int64, len(h.counts))
	var running int64
	for i, c := range h.counts {
		running += c
		cum[i] = running
	}
	return HistogramSnapshot{Bounds: latencyBuckets, Cumulative: cum, Sum: h.sum, Total: h.total}
}

// persistedState is the on-disk snapshot shape (schema-versioned).
type persistedState struct {
	Version     int           `json:"version"`
	Delegations []*Delegation `json:"delegations"`
	Order       []string      `json:"order"`
	Audit       []AuditEvent  `json:"audit"`
	Metrics     Metrics       `json:"metrics"`
}

const storeSchemaVersion = 1

// Store is the persistence layer: an in-memory index with an optional durable
// backing file. It is deliberately a concrete type with a small surface so a
// real database could replace it without touching the engine or handlers.
//
// Durability model: an atomic JSON snapshot (write-temp + rename). Mutations
// mark the store dirty; the server flushes on a short interval and on
// shutdown. Adequate for the current scale; the interface is the seam where a
// WAL/LSM backend would go.
type Store struct {
	mu      sync.Mutex
	path    string // "" => in-memory only
	dirty   bool
	byInst  map[string]*Delegation
	order   []string
	audit   []AuditEvent
	metrics Metrics
	latency *Histogram
}

func NewStore(path string) *Store {
	return &Store{path: path, byInst: map[string]*Delegation{}, latency: newHistogram()}
}

// ObserveLatency records a verification latency (seconds) into the histogram.
// Not persisted — histograms reset on restart, which is conventional.
func (s *Store) ObserveLatency(seconds float64) {
	s.mu.Lock()
	s.latency.observe(seconds)
	s.mu.Unlock()
}

// LatencySnapshot returns an immutable copy of the latency histogram.
func (s *Store) LatencySnapshot() HistogramSnapshot {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.latency.snapshot()
}

// Load reads the durable snapshot if a path is configured and the file exists.
func (s *Store) Load() error {
	if s.path == "" {
		return nil
	}
	raw, err := os.ReadFile(s.path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // fresh start
		}
		return err
	}
	var st persistedState
	if err := json.Unmarshal(raw, &st); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.byInst = map[string]*Delegation{}
	for _, d := range st.Delegations {
		s.byInst[d.Instance] = d
	}
	s.order = st.Order
	s.audit = st.Audit
	s.metrics = st.Metrics
	return nil
}

// Flush atomically writes the snapshot when dirty and a path is configured.
func (s *Store) Flush() error {
	s.mu.Lock()
	if s.path == "" || !s.dirty {
		s.mu.Unlock()
		return nil
	}
	dels := make([]*Delegation, 0, len(s.order))
	for _, inst := range s.order {
		dels = append(dels, s.byInst[inst])
	}
	st := persistedState{
		Version: storeSchemaVersion, Delegations: dels, Order: s.order,
		Audit: s.audit, Metrics: s.metrics,
	}
	s.dirty = false
	s.mu.Unlock()

	raw, err := json.MarshalIndent(st, "", " ")
	if err != nil {
		return err
	}
	tmp := s.path + ".tmp"
	if err := os.WriteFile(tmp, raw, 0o600); err != nil {
		return err
	}
	return os.Rename(tmp, s.path)
}

func (s *Store) AddDelegation(d *Delegation) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.byInst[d.Instance]; !ok {
		s.order = append(s.order, d.Instance)
	}
	s.byInst[d.Instance] = d
	s.metrics.Issued++
	s.dirty = true
}

func (s *Store) MarkRevoked(instance string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if d, ok := s.byInst[instance]; ok {
		d.Revoked = true
	}
	s.metrics.Revoked++
	s.dirty = true
}

func (s *Store) RecordVerdict(decision string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.metrics.Verified++
	switch decision {
	case "accept":
		s.metrics.Accept++
	case "reject":
		s.metrics.Reject++
	case "inconclusive":
		s.metrics.Inconclusive++
	}
	s.dirty = true
}

// Audit appends an event, keeping the log bounded to the most recent entries.
func (s *Store) Audit(e AuditEvent) {
	s.mu.Lock()
	defer s.mu.Unlock()
	const capN = 10000
	s.audit = append(s.audit, e)
	if len(s.audit) > capN {
		s.audit = s.audit[len(s.audit)-capN:]
	}
	s.dirty = true
}

// RevokedInstances returns the instance ids of delegations marked revoked —
// used at startup to rebuild the signed revoked set.
func (s *Store) RevokedInstances() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := []string{}
	for _, inst := range s.order {
		if d := s.byInst[inst]; d != nil && d.Revoked {
			out = append(out, inst)
		}
	}
	return out
}

// Delegations returns a copy of all delegations, newest first.
func (s *Store) Delegations() []*Delegation {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]*Delegation, 0, len(s.order))
	for i := len(s.order) - 1; i >= 0; i-- {
		out = append(out, s.byInst[s.order[i]])
	}
	return out
}

// AuditLog returns the most recent audit events, newest first, up to limit.
func (s *Store) AuditLog(limit int) []AuditEvent {
	s.mu.Lock()
	defer s.mu.Unlock()
	if limit <= 0 || limit > len(s.audit) {
		limit = len(s.audit)
	}
	out := make([]AuditEvent, 0, limit)
	for i := len(s.audit) - 1; i >= len(s.audit)-limit; i-- {
		out = append(out, s.audit[i])
	}
	return out
}

// Snapshot returns the current metric counters.
func (s *Store) Snapshot() Metrics {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.metrics
}

// Graph returns the delegation graph (principal → delegate edges) for the UI.
func (s *Store) Graph() GraphDTO {
	s.mu.Lock()
	defer s.mu.Unlock()
	nodes := map[string]bool{}
	edges := make([]GraphEdge, 0, len(s.order))
	for _, inst := range s.order {
		d := s.byInst[inst]
		nodes[d.Principal] = true
		nodes[d.Delegate] = true
		edges = append(edges, GraphEdge{From: d.Principal, To: d.Delegate, Instance: d.Instance, Scope: d.Scope, Revoked: d.Revoked})
	}
	ns := make([]GraphNode, 0, len(nodes))
	for id := range nodes {
		ns = append(ns, GraphNode{ID: id})
	}
	return GraphDTO{Nodes: ns, Edges: edges}
}

type GraphDTO struct {
	Nodes []GraphNode `json:"nodes"`
	Edges []GraphEdge `json:"edges"`
}
type GraphNode struct {
	ID string `json:"id"`
}
type GraphEdge struct {
	From     string   `json:"from"`
	To       string   `json:"to"`
	Instance string   `json:"instance"`
	Scope    []string `json:"scope"`
	Revoked  bool     `json:"revoked"`
}
