package main

import (
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
