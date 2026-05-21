package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// fixedClock lets the test control time deterministically.
type fixedClock struct{ t time.Time }

func (c *fixedClock) Now() time.Time { return c.t }

func newTestServer(t *testing.T) (*App, *httptest.Server, *fixedClock) {
	t.Helper()
	clk := &fixedClock{t: time.Unix(1_800_000_000, 0).UTC()}
	app, err := NewApp(Config{Domain: "domain-a.test"}, clk)
	if err != nil {
		t.Fatalf("NewApp: %v", err)
	}
	ts := httptest.NewServer(app.Router())
	t.Cleanup(ts.Close)
	return app, ts, clk
}

func post(t *testing.T, url string, body any) (*http.Response, map[string]any) {
	t.Helper()
	b, _ := json.Marshal(body)
	resp, err := http.Post(url, "application/json", bytes.NewReader(b))
	if err != nil {
		t.Fatalf("POST %s: %v", url, err)
	}
	var out map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&out)
	resp.Body.Close()
	return resp, out
}

func TestHealthAndVersion(t *testing.T) {
	_, ts, _ := newTestServer(t)
	for _, p := range []string{"/health", "/version"} {
		resp, err := http.Get(ts.URL + p)
		if err != nil || resp.StatusCode != 200 {
			t.Fatalf("GET %s: status=%v err=%v", p, resp.StatusCode, err)
		}
		resp.Body.Close()
	}
}

func TestIssueThenVerifyAccept(t *testing.T) {
	_, ts, _ := newTestServer(t)
	_, issued := post(t, ts.URL+"/issue", issueReq{
		Principal: "spiffe://domain-a.test/workload/payments-api",
		Delegate:  "spiffe://domain-b.test/agent/booking-worker",
		Scope:     []string{"read:orders", "write:audit"}, TTLSeconds: 3600,
	})
	rec, _ := issued["record"].(string)
	if rec == "" {
		t.Fatalf("issue returned no record: %v", issued)
	}
	if _, ok := issued["instance"].(string); !ok {
		t.Fatalf("issue returned no instance: %v", issued)
	}
	_, verified := post(t, ts.URL+"/verify", verifyReq{Record: rec})
	if verified["decision"] != "accept" || verified["accept"] != true {
		t.Fatalf("want accept, got %v", verified)
	}
	if _, ok := verified["latencyMicros"]; !ok {
		t.Fatalf("verify missing latency: %v", verified)
	}
}

func TestRevokeThenVerifyReject(t *testing.T) {
	_, ts, clk := newTestServer(t)
	_, issued := post(t, ts.URL+"/issue", issueReq{
		Principal: "spiffe://domain-a.test/workload/payments-api",
		Delegate:  "spiffe://domain-b.test/agent/booking-worker",
		Scope:     []string{"read:orders"}, TTLSeconds: 3600,
	})
	rec := issued["record"].(string)
	inst := issued["instance"].(string)

	// revoke, then advance the clock a hair so the new snapshot's as-of is
	// strictly newer and adopted.
	if resp, out := post(t, ts.URL+"/revoke", revokeReq{Instance: inst}); resp.StatusCode != 200 {
		t.Fatalf("revoke failed: %v", out)
	}
	clk.t = clk.t.Add(10 * time.Millisecond)
	if err := post200Refresh(t, ts); err != nil {
		t.Fatal(err)
	}

	_, verified := post(t, ts.URL+"/verify", verifyReq{Record: rec})
	if verified["decision"] != "reject" {
		t.Fatalf("want reject after revocation, got %v", verified)
	}
}

// post200Refresh nudges the app to re-publish (the background loop is not
// running in tests) so the just-revoked instance is observable.
func post200Refresh(t *testing.T, ts *httptest.Server) error {
	t.Helper()
	// revoke path already republished; nothing to do — kept for clarity.
	return nil
}

func TestAttenuationRefused(t *testing.T) {
	_, ts, _ := newTestServer(t)
	resp, out := post(t, ts.URL+"/issue", issueReq{
		Principal: "spiffe://domain-a.test/workload/payments-api",
		Delegate:  "spiffe://domain-b.test/agent/x",
		Scope:     []string{"delete:everything"}, TTLSeconds: 60,
	})
	if resp.StatusCode != 422 {
		t.Fatalf("over-scope should be refused (422), got %d: %v", resp.StatusCode, out)
	}
	if out["refused"] != true {
		t.Fatalf("expected refused=true, got %v", out)
	}
}

func TestPrincipalDomainEnforced(t *testing.T) {
	_, ts, _ := newTestServer(t)
	resp, _ := post(t, ts.URL+"/issue", issueReq{
		Principal: "spiffe://other-domain.test/x", Delegate: "spiffe://domain-b.test/y",
		Scope: []string{"read:orders"}, TTLSeconds: 60,
	})
	if resp.StatusCode != 400 {
		t.Fatalf("principal outside trust domain should be 400, got %d", resp.StatusCode)
	}
}

func TestMetricsPrometheusFormat(t *testing.T) {
	_, ts, _ := newTestServer(t)
	post(t, ts.URL+"/issue", issueReq{
		Principal: "spiffe://domain-a.test/p", Delegate: "spiffe://domain-b.test/d",
		Scope: []string{"read:orders"}, TTLSeconds: 60,
	})
	resp, err := http.Get(ts.URL + "/metrics")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	body := buf.String()
	for _, want := range []string{"atlas_verify_verdicts_total", "atlas_issued_total 1", "atlas_revocation_snapshot_age_seconds"} {
		if !strings.Contains(body, want) {
			t.Fatalf("metrics missing %q; got:\n%s", want, body)
		}
	}
}

func TestMetricsLatencyHistogram(t *testing.T) {
	_, ts, _ := newTestServer(t)
	_, issued := post(t, ts.URL+"/issue", issueReq{
		Principal: "spiffe://domain-a.test/p", Delegate: "spiffe://domain-b.test/d",
		Scope: []string{"read:orders"}, TTLSeconds: 60,
	})
	post(t, ts.URL+"/verify", verifyReq{Record: issued["record"].(string)})

	resp, err := http.Get(ts.URL + "/metrics")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	body := buf.String()
	for _, want := range []string{
		"atlas_verify_latency_seconds histogram",
		`atlas_verify_latency_seconds_bucket{le="+Inf"} 1`,
		"atlas_verify_latency_seconds_count 1",
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("metrics missing %q; got:\n%s", want, body)
		}
	}
}

func TestCORSPreflight(t *testing.T) {
	_, ts, _ := newTestServer(t)
	req, _ := http.NewRequest(http.MethodOptions, ts.URL+"/verify", nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 204 || resp.Header.Get("Access-Control-Allow-Origin") != "*" {
		t.Fatalf("CORS preflight failed: status=%d origin=%q", resp.StatusCode, resp.Header.Get("Access-Control-Allow-Origin"))
	}
}

// checkpoint: chore(fuzz): optimize Docker orchestration config (#141)
