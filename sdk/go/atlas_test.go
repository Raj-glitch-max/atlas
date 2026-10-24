package atlas

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// A stub server lets the client be tested without the real backend — the wire
// contract (paths, request/response JSON) is what we assert.
func TestClient_IssueVerifyRevoke(t *testing.T) {
	revoked := false
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/issue":
			var req map[string]any
			_ = json.NewDecoder(r.Body).Decode(&req)
			if req["principal"] == "" || req["ttlSeconds"] == nil {
				w.WriteHeader(400)
				_ = json.NewEncoder(w).Encode(map[string]string{"error": "bad request"})
				return
			}
			_ = json.NewEncoder(w).Encode(map[string]any{
				"record": "rec-abc", "instance": "inst-1",
				"principal": req["principal"], "delegate": req["delegate"], "scope": req["scope"],
			})
		case "/verify":
			dec := "accept"
			causes := []string{}
			if revoked {
				dec, causes = "reject", []string{"RevokedObservable"}
			}
			_ = json.NewEncoder(w).Encode(map[string]any{
				"decision": dec, "accept": dec == "accept", "causes": causes,
				"latencyMicros": 91,
				"trace": []map[string]string{
					{"check": "identity_binding", "outcome": "Pass"}, {"check": "integrity", "outcome": "Pass"},
					{"check": "expiry", "outcome": "Pass"}, {"check": "scope_integrity", "outcome": "Pass"},
					{"check": "revocation_status", "outcome": "Pass"},
				},
			})
		case "/revoke":
			revoked = true
			_ = json.NewEncoder(w).Encode(map[string]any{"revoked": true})
		default:
			w.WriteHeader(404)
		}
	}))
	defer srv.Close()

	c := New(srv.URL)
	ctx := context.Background()

	g, err := c.Issue(ctx, IssueParams{Principal: "spiffe://a/w", Delegate: "spiffe://b/a", Scope: []string{"read:orders"}})
	if err != nil {
		t.Fatalf("issue: %v", err)
	}
	if g.Record != "rec-abc" || g.Instance != "inst-1" {
		t.Fatalf("unexpected grant: %+v", g)
	}

	v, err := c.Verify(ctx, g.Record)
	if err != nil {
		t.Fatalf("verify: %v", err)
	}
	if v.Decision != Accept || !v.Accept || len(v.Trace) != 5 {
		t.Fatalf("expected accept with 5-check trace, got %+v", v)
	}

	if err := c.Revoke(ctx, g.Instance); err != nil {
		t.Fatalf("revoke: %v", err)
	}

	v2, err := c.Verify(ctx, g.Record)
	if err != nil {
		t.Fatalf("verify after revoke: %v", err)
	}
	if v2.Decision != Reject {
		t.Fatalf("expected reject after revoke, got %s", v2.Decision)
	}
}

func TestClient_HTTPErrorSurfaces(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(401)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "unauthorized"})
	}))
	defer srv.Close()

	err := New(srv.URL).Revoke(context.Background(), "inst-1")
	var ae *Error
	if err == nil {
		t.Fatal("expected an error")
	}
	if e, ok := err.(*Error); !ok || e.Status != 401 {
		t.Fatalf("expected *Error status 401, got %T %v", err, err)
	}
	_ = ae
}

// DEFECT D6: Python and TypeScript SDK verify() and the MCP verify tool make network calls;
// no offline verifier exists outside Go.
// The Go SDK Verify() method does HTTP POST to /verify.
// This test asserts the CURRENT WRONG behavior: Verify requires network call to server.
// After fix, this test should be renamed TestD6_OfflineVerifierExists and inverted.
func TestDefectD6_NoOfflineVerifier(t *testing.T) {
	c := New("http://127.0.0.1:8087") // Non-existent server

	ctx := context.Background()
	_, err := c.Verify(ctx, "some-record")

	// CURRENT WRONG BEHAVIOR: Verify returns network error (unreachable)
	if err == nil {
		t.Fatalf("DEFECT D6 NOT REPRODUCED: Verify succeeded without server (offline verifier exists)")
	}
	atlasErr, ok := err.(*Error)
	if !ok || atlasErr.Status != 0 {
		t.Fatalf("DEFECT D6 NOT REPRODUCED: expected network error (Status=0), got %T %v", err, err)
	}
	if !strings.Contains(atlasErr.Message, "unreachable") {
		t.Fatalf("DEFECT D6 NOT REPRODUCED: expected unreachable message, got %s", atlasErr.Message)
	}
	t.Logf("DEFECT D6 REPRODUCED: SDK Verify() requires network call, no offline path")
}
