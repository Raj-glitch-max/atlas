package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func mockServer(t *testing.T, revoked *bool) *httptest.Server {
	t.Helper()
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{"status": "ok"})
	})
	mux.HandleFunc("/version", func(w http.ResponseWriter, _ *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{"version": "test", "trustDomain": "domain-a.test", "algorithm": "ES256 (P-256)", "revocationR": "2s"})
	})
	mux.HandleFunc("/issue", func(w http.ResponseWriter, _ *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{"record": "REC.J.WS", "instance": "inst-xyz", "principal": "spiffe://domain-a.test/p", "delegate": "spiffe://domain-b.test/d", "scope": []string{"read:orders"}, "expiresAt": "2027-01-01T00:00:00Z"})
	})
	mux.HandleFunc("/verify", func(w http.ResponseWriter, _ *http.Request) {
		if *revoked {
			json.NewEncoder(w).Encode(map[string]any{"decision": "reject", "accept": false, "causes": []string{"RevokedObservable"}, "latencyMicros": 140, "trace": []traceEntry{{Check: "revocation_status", Outcome: "Reject", Cause: "RevokedObservable"}}})
			return
		}
		json.NewEncoder(w).Encode(map[string]any{"decision": "accept", "accept": true, "causes": []string{}, "latencyMicros": 200, "trace": []traceEntry{{Check: "identity_binding", Outcome: "Pass", Cause: "None"}}})
	})
	mux.HandleFunc("/revoke", func(w http.ResponseWriter, _ *http.Request) {
		*revoked = true
		json.NewEncoder(w).Encode(map[string]any{"revoked": true})
	})
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)
	return ts
}

func TestIssueVerifyRevokeFlow(t *testing.T) {
	revoked := false
	ts := mockServer(t, &revoked)
	var out bytes.Buffer

	// issue -q prints only the record
	if code := run([]string{"--api", ts.URL, "issue", "-q", "--principal", "spiffe://domain-a.test/p", "--delegate", "spiffe://domain-b.test/d", "--scope", "read:orders"}, &out); code != exitOK {
		t.Fatalf("issue exit=%d out=%s", code, out.String())
	}
	rec := strings.TrimSpace(out.String())
	if rec != "REC.J.WS" {
		t.Fatalf("issue -q should print only the record, got %q", rec)
	}

	// verify accepts → exit 0
	out.Reset()
	if code := run([]string{"--api", ts.URL, "verify", rec}, &out); code != exitOK {
		t.Fatalf("verify(accept) exit=%d out=%s", code, out.String())
	}
	if !strings.Contains(out.String(), "ACCEPT") {
		t.Fatalf("verify output missing ACCEPT: %s", out.String())
	}

	// revoke, then verify rejects → exit 3
	out.Reset()
	if code := run([]string{"--api", ts.URL, "revoke", "inst-xyz"}, &out); code != exitOK {
		t.Fatalf("revoke exit=%d", code)
	}
	out.Reset()
	code := run([]string{"--api", ts.URL, "verify", rec}, &out)
	if code != exitNonPass {
		t.Fatalf("verify(reject) should exit %d, got %d", exitNonPass, code)
	}
	if !strings.Contains(out.String(), "REJECT") || !strings.Contains(out.String(), "RevokedObservable") {
		t.Fatalf("verify(reject) output wrong: %s", out.String())
	}
}

func TestDoctorReachable(t *testing.T) {
	revoked := false
	ts := mockServer(t, &revoked)
	var out bytes.Buffer
	if code := run([]string{"--api", ts.URL, "doctor"}, &out); code != exitOK {
		t.Fatalf("doctor exit=%d out=%s", code, out.String())
	}
	if !strings.Contains(out.String(), "reachable") || !strings.Contains(out.String(), "domain-a.test") {
		t.Fatalf("doctor output: %s", out.String())
	}
}

func TestDoctorUnreachable(t *testing.T) {
	var out bytes.Buffer
	if code := run([]string{"--api", "http://127.0.0.1:1", "doctor"}, &out); code != exitError {
		t.Fatalf("doctor(unreachable) should exit %d, got %d", exitError, code)
	}
	if !strings.Contains(out.String(), "UNREACHABLE") {
		t.Fatalf("doctor(unreachable) output: %s", out.String())
	}
}

func TestInspectDecodesOffline(t *testing.T) {
	// a compact JWS whose payload carries Atlas claims (offline decode; the
	// signature segment is arbitrary — inspect never verifies it).
	hdr := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"ES256","kid":"authority-key-1","typ":"atlas-record+jws"}`))
	pl := base64.RawURLEncoding.EncodeToString([]byte(`{"sub":"spiffe://domain-a.test/p","act":{"sub":"spiffe://domain-b.test/d"},"scope":["read:orders"]}`))
	rec := hdr + "." + pl + ".c2ln"

	var out bytes.Buffer
	// no --api needed; inspect is offline
	if code := run([]string{"inspect", rec}, &out); code != exitOK {
		t.Fatalf("inspect exit=%d out=%s", code, out.String())
	}
	s := out.String()
	for _, want := range []string{"atlas-record+jws", "spiffe://domain-a.test/p", "read:orders", "NOT verified"} {
		if !strings.Contains(s, want) {
			t.Fatalf("inspect output missing %q:\n%s", want, s)
		}
	}
}

func TestInspectRejectsNonJWS(t *testing.T) {
	var out bytes.Buffer
	if code := run([]string{"inspect", "not-a-jws"}, &out); code != exitError {
		t.Fatalf("inspect of non-JWS should error, got %d", code)
	}
}

func TestUsageAndUnknown(t *testing.T) {
	var out bytes.Buffer
	if code := run([]string{}, &out); code != exitUsage {
		t.Fatalf("no args should be usage exit, got %d", code)
	}
	out.Reset()
	if code := run([]string{"frobnicate"}, &out); code != exitUsage {
		t.Fatalf("unknown command should be usage exit, got %d", code)
	}
}

// checkpoint: refactor(issuance): refactor revstatus snapshot retrieval
