package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// TestPersistenceRoundTrip proves durable state survives a restart: issue two
// delegations, revoke one, flush; then build a fresh App from the same file and
// confirm the delegations, the revoked flag, and (crucially) the revocation
// verdict are all restored.
func TestPersistenceRoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "atlas-state.json")
	keyPath := filepath.Join(dir, "authority.key")
	cfg := Config{Domain: "domain-a.test", StorePath: path, KeyPath: keyPath}
	clk := &fixedClock{t: time.Unix(1_800_000_000, 0).UTC()}

	// --- session 1: issue two, revoke one, flush ---
	app1, err := NewApp(cfg, clk)
	if err != nil {
		t.Fatalf("NewApp session1: %v", err)
	}
	ts1 := httptest.NewServer(app1.Router())
	r1, e := app1.Issue("spiffe://domain-a.test/workload/a", "spiffe://domain-b.test/agent/x", []string{"read:orders"}, time.Hour)
	if e != nil {
		t.Fatalf("issue1: %v", e.Message)
	}
	r2, e := app1.Issue("spiffe://domain-a.test/workload/b", "spiffe://domain-b.test/agent/y", []string{"read:orders"}, time.Hour)
	if e != nil {
		t.Fatalf("issue2: %v", e.Message)
	}
	_ = r2
	if e := app1.Revoke(r1.Instance); e != nil {
		t.Fatalf("revoke: %v", e.Message)
	}
	if err := app1.Flush(); err != nil {
		t.Fatalf("flush: %v", err)
	}
	ts1.Close()

	// the file must exist and be valid JSON with our schema version
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("state file not written: %v", err)
	}
	var st persistedState
	if err := json.Unmarshal(raw, &st); err != nil || st.Version != storeSchemaVersion {
		t.Fatalf("bad state file (version %d): %v", st.Version, err)
	}

	// --- session 2: fresh App from the same file (same key + store) ---
	clk2 := &fixedClock{t: clk.t.Add(time.Second)} // newer clock so snapshots adopt
	app2, err := NewApp(cfg, clk2)
	if err != nil {
		t.Fatalf("NewApp session2: %v", err)
	}
	dels := app2.store.Delegations()
	if len(dels) != 2 {
		t.Fatalf("want 2 delegations restored, got %d", len(dels))
	}
	// the revoked one must still read revoked, and the verifier must reject its
	// record (the revoked set was rebuilt and republished on load).
	v := app2.Verify(r1.Record)
	if v.Decision != "reject" {
		t.Fatalf("restored revocation not observed: verdict=%s causes=%v", v.Decision, v.Causes)
	}
	// the other one still verifies accept
	v2 := app2.Verify(r2.Record)
	if v2.Decision != "accept" {
		t.Fatalf("non-revoked delegation should still accept, got %s %v", v2.Decision, v2.Causes)
	}
	// metrics carried over (2 issued from session 1)
	if m := app2.store.Snapshot(); m.Issued < 2 {
		t.Fatalf("metrics not restored: issued=%d", m.Issued)
	}
}

// TestAPIKeyGuardsMutations confirms that when an API key is set, mutating
// endpoints require the bearer token while read endpoints stay open.
func TestAPIKeyGuardsMutations(t *testing.T) {
	clk := &fixedClock{t: time.Unix(1_800_000_000, 0).UTC()}
	app, err := NewApp(Config{Domain: "domain-a.test", APIKey: "s3cret"}, clk)
	if err != nil {
		t.Fatal(err)
	}
	ts := httptest.NewServer(app.Router())
	t.Cleanup(ts.Close)

	body := `{"principal":"spiffe://domain-a.test/p","delegate":"spiffe://domain-b.test/d","scope":["read:orders"]}`

	// no token → 401
	resp, _ := http.Post(ts.URL+"/issue", "application/json", strings.NewReader(body))
	if resp.StatusCode != 401 {
		t.Fatalf("issue without token should be 401, got %d", resp.StatusCode)
	}
	// wrong token → 401
	req, _ := http.NewRequest("POST", ts.URL+"/issue", strings.NewReader(body))
	req.Header.Set("Authorization", "Bearer nope")
	resp, _ = http.DefaultClient.Do(req)
	if resp.StatusCode != 401 {
		t.Fatalf("issue with wrong token should be 401, got %d", resp.StatusCode)
	}
	// correct token → 200
	req, _ = http.NewRequest("POST", ts.URL+"/issue", strings.NewReader(body))
	req.Header.Set("Authorization", "Bearer s3cret")
	resp, _ = http.DefaultClient.Do(req)
	if resp.StatusCode != 200 {
		t.Fatalf("issue with correct token should be 200, got %d", resp.StatusCode)
	}
	// read endpoint stays open (no token)
	resp, _ = http.Get(ts.URL + "/delegations")
	if resp.StatusCode != 200 {
		t.Fatalf("read endpoint should stay open, got %d", resp.StatusCode)
	}
}

// TestInMemoryWhenNoPath confirms an empty path means no file is written.
func TestInMemoryWhenNoPath(t *testing.T) {
	clk := &fixedClock{t: time.Unix(1_800_000_000, 0).UTC()}
	app, err := NewApp(Config{Domain: "domain-a.test"}, clk)
	if err != nil {
		t.Fatal(err)
	}
	if _, e := app.Issue("spiffe://domain-a.test/p", "spiffe://domain-b.test/d", []string{"read:orders"}, time.Hour); e != nil {
		t.Fatalf("issue: %v", e.Message)
	}
	if err := app.Flush(); err != nil {
		t.Fatalf("flush (in-memory) should be a no-op, got: %v", err)
	}
}

// checkpoint: chore(sdk): tweak revocation status lookup
