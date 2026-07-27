package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// fakeClock is a controllable clock for deterministic rate-limit tests.
type fakeClock struct{ t time.Time }

func (c *fakeClock) Now() time.Time          { return c.t }
func (c *fakeClock) advance(d time.Duration) { c.t = c.t.Add(d) }

func TestRateLimiter_BurstThenRefuse(t *testing.T) {
	clk := &fakeClock{t: time.Unix(1_700_000_000, 0)}
	rl := newRateLimiter(60, clk) // 60/min => burst of 10

	// Burst of `burst` tokens should pass, then be refused within the same instant.
	allowed := 0
	for i := 0; i < 50; i++ {
		if rl.allow("1.2.3.4") {
			allowed++
		}
	}
	if allowed != rl.burst {
		t.Fatalf("expected exactly burst=%d allowed in an instant, got %d", rl.burst, allowed)
	}
	if rl.allow("1.2.3.4") {
		t.Fatal("expected refusal once the bucket is empty")
	}
}

func TestRateLimiter_RefillsOverTime(t *testing.T) {
	clk := &fakeClock{t: time.Unix(1_700_000_000, 0)}
	rl := newRateLimiter(60, clk) // 1 token/sec

	for rl.allow("ip") { // drain the bucket
	}
	if rl.allow("ip") {
		t.Fatal("bucket should be empty")
	}
	clk.advance(1100 * time.Millisecond) // ~1 token refilled
	if !rl.allow("ip") {
		t.Fatal("expected a refilled token to be available after ~1s")
	}
}

func TestRateLimiter_IsolatesClients(t *testing.T) {
	clk := &fakeClock{t: time.Unix(1_700_000_000, 0)}
	rl := newRateLimiter(6, clk)
	for rl.allow("a") { // exhaust client a
	}
	if !rl.allow("b") {
		t.Fatal("client b must not be limited by client a's usage")
	}
}

// TestRateLimitBehindProxyIdentifiesRealClient pins the PaaS failure mode.
//
// Behind Railway/Render/Cloud Run/a CDN, every request arrives from the
// platform's proxy, so RemoteAddr is identical for all visitors. Without
// -trust-proxy a per-IP limit silently becomes ONE GLOBAL BUDGET: a single busy
// client locks out the entire internet. That is a demo-killing outage, not a
// subtle degradation, so it gets a test.
func TestRateLimitBehindProxyIdentifiesRealClient(t *testing.T) {
	proxied := func(xff string) *http.Request {
		r := httptest.NewRequest(http.MethodPost, "/issue", nil)
		r.RemoteAddr = "10.0.0.1:54321" // the proxy — same for everyone
		r.Header.Set("X-Forwarded-For", xff)
		return r
	}

	// trustProxy=false: two different visitors collapse to one identity.
	a := clientIP(proxied("203.0.113.7"), false)
	b := clientIP(proxied("198.51.100.9"), false)
	if a != b {
		t.Fatalf("precondition: without trust-proxy both should be the proxy IP, got %q and %q", a, b)
	}

	// trustProxy=true: they are distinguished.
	a = clientIP(proxied("203.0.113.7"), true)
	b = clientIP(proxied("198.51.100.9"), true)
	if a == b {
		t.Fatal("with -trust-proxy, distinct clients must not share a rate-limit bucket")
	}
	if a != "203.0.113.7" || b != "198.51.100.9" {
		t.Fatalf("wrong client IPs extracted: %q, %q", a, b)
	}

	// A multi-hop chain: the left-most entry is the originating client.
	if got := clientIP(proxied("203.0.113.7, 70.41.3.18, 150.172.238.178"), true); got != "203.0.113.7" {
		t.Fatalf("multi-hop XFF: got %q, want the left-most entry", got)
	}

	// Without the header, fall back to RemoteAddr rather than an empty bucket
	// key that every caller would share.
	r := httptest.NewRequest(http.MethodPost, "/issue", nil)
	r.RemoteAddr = "192.0.2.44:1234"
	if got := clientIP(r, true); got != "192.0.2.44" {
		t.Fatalf("missing XFF should fall back to RemoteAddr host, got %q", got)
	}
}
