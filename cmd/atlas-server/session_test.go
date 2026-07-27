package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// newTestApp builds an in-memory app with an ephemeral key — no disk, no
// leftover state between tests.
func newTestApp(t *testing.T) *App {
	t.Helper()
	app, err := NewApp(Config{Domain: "domain-a.test"}, systemClock{})
	if err != nil {
		t.Fatalf("NewApp: %v", err)
	}
	return app
}

// do issues an HTTP request against the router, optionally carrying a session
// id, and decodes the JSON response.
func do(t *testing.T, h http.Handler, method, path, session, body string) (int, map[string]any) {
	t.Helper()
	var rdr *strings.Reader
	if body == "" {
		rdr = strings.NewReader("")
	} else {
		rdr = strings.NewReader(body)
	}
	req := httptest.NewRequest(method, path, rdr)
	req.Header.Set("Content-Type", "application/json")
	if session != "" {
		req.Header.Set(sessionHeader, session)
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	var out map[string]any
	_ = json.Unmarshal(rec.Body.Bytes(), &out)
	return rec.Code, out
}

func issueIn(t *testing.T, h http.Handler, session string) (record, instance string) {
	t.Helper()
	code, body := do(t, h, http.MethodPost, "/issue", session,
		`{"principal":"spiffe://domain-a.test/workload/payments-api",`+
			`"delegate":"spiffe://domain-b.test/agent/booking-worker",`+
			`"scope":["read:orders"],"ttlSeconds":3600}`)
	if code != 200 {
		t.Fatalf("issue in session %q: status %d body %v", session, code, body)
	}
	rec, _ := body["record"].(string)
	inst, _ := body["instance"].(string)
	if rec == "" || inst == "" {
		t.Fatalf("issue returned empty record/instance: %v", body)
	}
	return rec, inst
}

func decisionFor(t *testing.T, h http.Handler, session, rec string) string {
	t.Helper()
	payload, err := json.Marshal(map[string]string{"record": rec})
	if err != nil {
		t.Fatal(err)
	}
	code, body := do(t, h, http.MethodPost, "/verify", session, string(payload))
	if code != 200 {
		t.Fatalf("verify: status %d body %v", code, body)
	}
	dec, _ := body["decision"].(string)
	return dec
}

const (
	sessA = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	sessB = "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
)

// TestSessionsAreIsolated is the property this feature exists for: one
// visitor's sandbox must be invisible and untouchable to another's.
func TestSessionsAreIsolated(t *testing.T) {
	h := newTestApp(t).Router()

	_, instA := issueIn(t, h, sessA)
	_, instB := issueIn(t, h, sessB)

	if instA == instB {
		t.Fatal("two sessions produced the same instance id")
	}

	// A's delegation must not appear in B's graph or delegation list.
	_, delsB := do(t, h, http.MethodGet, "/delegations", sessB, "")
	listB, _ := delsB["delegations"].([]any)
	if len(listB) != 1 {
		t.Fatalf("session B should see exactly its own 1 delegation, saw %d", len(listB))
	}
	for _, d := range listB {
		m, _ := d.(map[string]any)
		if m["instance"] == instA {
			t.Fatal("session B can see session A's delegation — isolation broken")
		}
	}

	// B's audit log must not contain A's events.
	_, audB := do(t, h, http.MethodGet, "/audit", sessB, "")
	evB, _ := audB["events"].([]any)
	for _, e := range evB {
		m, _ := e.(map[string]any)
		if m["instance"] == instA {
			t.Fatal("session A's instance leaked into session B's audit log")
		}
	}
}

// TestSessionCannotRevokeAnotherSessionsCapability is the sharp edge of the
// property: on the shared global store, this revoke WOULD have rejected A's
// record. It must not.
func TestSessionCannotRevokeAnotherSessionsCapability(t *testing.T) {
	h := newTestApp(t).Router()

	recA, instA := issueIn(t, h, sessA)
	if got := decisionFor(t, h, sessA, recA); got != "accept" {
		t.Fatalf("A's fresh record should accept in A, got %q", got)
	}

	// B attempts to revoke A's instance id.
	code, _ := do(t, h, http.MethodPost, "/revoke", sessB,
		`{"instance":"`+instA+`"}`)
	if code != 200 {
		t.Fatalf("revoke call itself should succeed (it just affects B's set), got %d", code)
	}

	// A's record must still verify in A.
	if got := decisionFor(t, h, sessA, recA); got != "accept" {
		t.Fatalf("session B revoked session A's capability — isolation broken (A now %q)", got)
	}

	// And A revoking its own must still work.
	if code, _ := do(t, h, http.MethodPost, "/revoke", sessA, `{"instance":"`+instA+`"}`); code != 200 {
		t.Fatalf("A revoking its own instance failed: %d", code)
	}
	if got := decisionFor(t, h, sessA, recA); got != "reject" {
		t.Fatalf("A's own revocation must be enforced in A, got %q", got)
	}
}

// TestDefaultSessionUnchangedForNonSessionClients guards backward
// compatibility: the CLI, SDKs and MCP server send no session header and must
// keep seeing the single durable graph.
func TestDefaultSessionUnchangedForNonSessionClients(t *testing.T) {
	h := newTestApp(t).Router()

	rec, inst := issueIn(t, h, "") // no session header
	if got := decisionFor(t, h, "", rec); got != "accept" {
		t.Fatalf("default session verify: got %q", got)
	}
	if code, _ := do(t, h, http.MethodPost, "/revoke", "", `{"instance":"`+inst+`"}`); code != 200 {
		t.Fatal("default session revoke failed")
	}
	if got := decisionFor(t, h, "", rec); got != "reject" {
		t.Fatalf("default session revocation not enforced: %q", got)
	}

	// An invalid session id must fall back to the default session, not create a
	// silent empty sandbox that makes state appear to vanish.
	if got := decisionFor(t, h, "not-a-valid-session-id", rec); got != "reject" {
		t.Fatalf("invalid session id should fall back to default, got %q", got)
	}
}

// TestActivityAggregatesWithoutLeaking checks the one cross-session surface:
// it must count across sessions and expose no per-session identifiers.
func TestActivityAggregatesWithoutLeaking(t *testing.T) {
	app := newTestApp(t)
	h := app.Router()

	_, instA := issueIn(t, h, sessA)
	issueIn(t, h, sessB)
	do(t, h, http.MethodPost, "/revoke", sessA, `{"instance":"`+instA+`"}`)

	code, body := do(t, h, http.MethodGet, "/activity", "", "")
	if code != 200 {
		t.Fatalf("/activity status %d", code)
	}
	if got := body["issuedToday"]; got != float64(2) {
		t.Fatalf("issuedToday should aggregate both sessions, got %v", got)
	}
	if got := body["revocationsLastHour"]; got != float64(1) {
		t.Fatalf("revocationsLastHour = %v, want 1", got)
	}
	if got := body["liveSessions"]; got != float64(2) {
		t.Fatalf("liveSessions = %v, want 2", got)
	}

	// The response must not carry any identifier that ties activity to a
	// session or a workload.
	raw, _ := json.Marshal(body)
	for _, leak := range []string{sessA, sessB, instA, "payments-api", "booking-worker"} {
		if strings.Contains(string(raw), leak) {
			t.Fatalf("/activity leaked %q: %s", leak, raw)
		}
	}
}

// TestSessionEvictionBoundsMemory checks the public-instance backstop.
func TestSessionEvictionBoundsMemory(t *testing.T) {
	app := newTestApp(t)
	app.sessionsMu.Lock()
	for i := 0; i < maxSessions+10; i++ {
		id, err := newSessionID()
		if err != nil {
			t.Fatal(err)
		}
		app.evictLocked()
		s, err := app.newSession(id)
		if err != nil {
			t.Fatal(err)
		}
		app.sessions[id] = s
	}
	n := len(app.sessions)
	app.sessionsMu.Unlock()
	if n > maxSessions {
		t.Fatalf("session count %d exceeded cap %d", n, maxSessions)
	}
}

// TestIdleSessionsAreEvicted verifies the idle TTL actually drops state.
func TestIdleSessionsAreEvicted(t *testing.T) {
	app := newTestApp(t)
	id, _ := newSessionID()
	s, err := app.newSession(id)
	if err != nil {
		t.Fatal(err)
	}
	s.lastSeen = time.Now().Add(-2 * sessionIdleTTL)
	app.sessionsMu.Lock()
	app.sessions[id] = s
	app.evictLocked()
	_, still := app.sessions[id]
	app.sessionsMu.Unlock()
	if still {
		t.Fatal("idle session past its TTL was not evicted")
	}
}
