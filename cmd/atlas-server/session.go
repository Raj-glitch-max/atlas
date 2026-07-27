package main

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"sync"
	"time"

	"github.com/Raj-glitch-max/atlas/internal/record"
	"github.com/Raj-glitch-max/atlas/internal/revstatus"
)

// Session is one caller's isolated slice of Atlas state: its own delegation
// store, its own audit log, and — critically — its own revoked set and signed
// revocation snapshot.
//
// Why this exists: a single atlas-server instance otherwise holds ONE global
// delegation graph. That is correct for a single-tenant deployment and wrong
// for a public demo, where every visitor would share one audit log and any
// visitor could revoke another visitor's capability. Sessions make the demo
// honest without pretending the server became multi-tenant: this is per-visitor
// isolation for a sandbox, not a tenancy boundary with an authorization model
// behind it. See LIMITATIONS.md section 5.
//
// The authority key, trust material, and issuance policy are deliberately
// SHARED across sessions — every visitor is issued by the same real authority
// in the same trust domain. Only mutable per-visitor state is isolated.
type Session struct {
	ID       string
	created  time.Time
	lastSeen time.Time

	store    *Store
	provider *revstatus.SignedSetProvider

	mu       sync.RWMutex
	revoked  map[string]record.InstanceID
	lastAsOf time.Time
	lastSet  revstatus.SignedRevokedSet
}

const (
	// sessionHeader is how a browser client claims a session. A header rather
	// than a cookie: the demo is cross-origin (static site → API host), and a
	// header avoids SameSite/credentialed-CORS complexity entirely.
	sessionHeader = "X-Atlas-Session"
	sessionCookie = "atlas_session"

	// sessionIdleTTL evicts a visitor's state after this much inactivity.
	sessionIdleTTL = 30 * time.Minute

	// maxSessions bounds memory on a public instance. At the cap, the least
	// recently seen session is evicted. This is a demo backstop, not a quota
	// system.
	maxSessions = 500
)

// newSessionID returns an unguessable session identifier. Unguessable matters:
// knowing another visitor's id would hand you their isolated state.
func newSessionID() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// validSessionID rejects anything that is not exactly what newSessionID emits.
// Client-supplied ids are accepted (the browser generates its own so it can
// survive a reload), so the format is constrained rather than trusted.
func validSessionID(s string) bool {
	if len(s) != 32 {
		return false
	}
	for i := 0; i < len(s); i++ {
		c := s[i]
		if (c < '0' || c > '9') && (c < 'a' || c > 'f') {
			return false
		}
	}
	return true
}

// newSession builds an in-memory session. Sessions are never persisted: a
// visitor's sandbox state is intentionally ephemeral.
func (a *App) newSession(id string) (*Session, error) {
	s := &Session{
		ID:       id,
		created:  a.clock.Now(),
		lastSeen: a.clock.Now(),
		store:    NewStore(""), // in-memory only
		provider: revstatus.NewSignedSetProvider(a.pubKey, a.listID),
		revoked:  map[string]record.InstanceID{},
	}
	// Seed a signed (empty) snapshot so verification answers immediately with a
	// fresh revocation state instead of Indeterminate.
	if err := a.republish(s); err != nil {
		return nil, err
	}
	return s, nil
}

// sessionFor resolves the session for a request. A request with no valid
// session id gets the default session — which is the durable, single-graph
// behaviour every existing client (CLI, SDKs, MCP, tests) already relies on.
// Only clients that opt in by sending a session id get isolation.
func (a *App) sessionFor(r *http.Request) *Session {
	id := r.Header.Get(sessionHeader)
	if id == "" {
		if c, err := r.Cookie(sessionCookie); err == nil {
			id = c.Value
		}
	}
	if id == "" || !validSessionID(id) {
		return a.def
	}

	a.sessionsMu.Lock()
	defer a.sessionsMu.Unlock()

	if s, ok := a.sessions[id]; ok {
		s.lastSeen = a.clock.Now()
		return s
	}
	a.evictLocked()
	s, err := a.newSession(id)
	if err != nil {
		// Publishing a snapshot failed; fall back to the default session rather
		// than serving a session that cannot answer revocation questions.
		return a.def
	}
	a.sessions[id] = s
	return s
}

// evictLocked drops idle sessions, then enforces the hard cap by evicting the
// least recently seen. Caller holds sessionsMu.
func (a *App) evictLocked() {
	now := a.clock.Now()
	for id, s := range a.sessions {
		if now.Sub(s.lastSeen) > sessionIdleTTL {
			delete(a.sessions, id)
		}
	}
	for len(a.sessions) >= maxSessions {
		var oldestID string
		var oldest time.Time
		for id, s := range a.sessions {
			if oldestID == "" || s.lastSeen.Before(oldest) {
				oldestID, oldest = id, s.lastSeen
			}
		}
		if oldestID == "" {
			return
		}
		delete(a.sessions, oldestID)
	}
}

// SessionCount reports how many isolated sessions are live (the default session
// is not counted).
func (a *App) SessionCount() int {
	a.sessionsMu.Lock()
	defer a.sessionsMu.Unlock()
	return len(a.sessions)
}

// ---- aggregated, anonymized activity ----

// activityKind is the event type recorded in the public aggregate.
type activityKind uint8

const (
	actIssued activityKind = iota
	actRevoked
	actVerified
)

// activityLog is a bounded, append-only ring of (time, kind) pairs used to
// answer "how much has happened recently" ACROSS all sessions.
//
// It records the kind of event and when — and nothing else. No session id, no
// SPIFFE ID, no scope, no instance id, no IP. That is deliberate: this is the
// one surface every visitor can see, so it must not be able to leak anything
// about any individual visitor's activity beyond its existence in a total.
type activityLog struct {
	mu     sync.Mutex
	times  []time.Time
	kinds  []activityKind
	cursor int
	filled bool
}

const activityCapacity = 20000

func newActivityLog() *activityLog {
	return &activityLog{
		times: make([]time.Time, activityCapacity),
		kinds: make([]activityKind, activityCapacity),
	}
}

func (l *activityLog) record(k activityKind, at time.Time) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.times[l.cursor] = at
	l.kinds[l.cursor] = k
	l.cursor++
	if l.cursor == activityCapacity {
		l.cursor = 0
		l.filled = true
	}
}

// countSince returns per-kind counts for events at or after cutoff.
func (l *activityLog) countSince(cutoff time.Time) map[activityKind]int64 {
	l.mu.Lock()
	defer l.mu.Unlock()
	out := map[activityKind]int64{}
	n := l.cursor
	if l.filled {
		n = activityCapacity
	}
	for i := 0; i < n; i++ {
		if !l.times[i].Before(cutoff) {
			out[l.kinds[i]]++
		}
	}
	return out
}
