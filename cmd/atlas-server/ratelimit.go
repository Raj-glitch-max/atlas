package main

import (
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

// rateLimiter is a zero-dependency, per-client token-bucket limiter for the
// mutating endpoints. It protects /issue and /revoke from abuse without adding
// a dependency; read-only and probe endpoints are never limited. A bucket
// holds up to `burst` tokens and refills at `perMinute` tokens per minute.
//
// Buckets are keyed by client IP and swept lazily, so idle clients don't leak
// memory. This is deliberately simple (single mutex): the server is a
// single-domain reference deployment, not a high-fan-out gateway.
type rateLimiter struct {
	perMinute int
	burst     int
	clock     Clock

	mu      sync.Mutex
	buckets map[string]*bucket
}

type bucket struct {
	tokens float64
	last   time.Time
}

func newRateLimiter(perMinute int, clock Clock) *rateLimiter {
	burst := perMinute / 6 // ~10s worth of headroom
	if burst < 1 {
		burst = 1
	}
	return &rateLimiter{perMinute: perMinute, burst: burst, clock: clock, buckets: map[string]*bucket{}}
}

// allow reports whether a request from key may proceed, consuming one token.
func (rl *rateLimiter) allow(key string) bool {
	now := rl.clock.Now()
	rl.mu.Lock()
	defer rl.mu.Unlock()

	b, ok := rl.buckets[key]
	if !ok {
		rl.buckets[key] = &bucket{tokens: float64(rl.burst) - 1, last: now}
		if len(rl.buckets) > 4096 {
			rl.sweep(now)
		}
		return true
	}
	// Refill by elapsed time, cap at burst.
	elapsed := now.Sub(b.last).Minutes()
	b.tokens += elapsed * float64(rl.perMinute)
	if b.tokens > float64(rl.burst) {
		b.tokens = float64(rl.burst)
	}
	b.last = now
	if b.tokens < 1 {
		return false
	}
	b.tokens--
	return true
}

// sweep drops buckets that have fully refilled (idle) to bound memory.
func (rl *rateLimiter) sweep(now time.Time) {
	for k, b := range rl.buckets {
		if now.Sub(b.last).Minutes()*float64(rl.perMinute) >= float64(rl.burst) {
			delete(rl.buckets, k)
		}
	}
}

// clientIP extracts a best-effort client identity for limiting.
//
// trustProxy MUST be false unless the server actually sits behind a proxy that
// sets X-Forwarded-For. The distinction matters in both directions:
//
//   - Behind a PaaS (Railway, Render, Cloud Run, Fly, any CDN) every request
//     arrives from the platform's proxy, so RemoteAddr is IDENTICAL for every
//     visitor on earth. Limiting on it turns a per-IP limit into one global
//     budget: a single busy client locks out everyone.
//   - Exposed directly, X-Forwarded-For is attacker-controlled. Trusting it
//     there lets anyone evade the limit by rotating a header value.
//
// So it is opt-in, and honest about what it buys: with a proxy in front, the
// left-most XFF entry is the originating client by convention, but a client can
// still forge that header before it reaches the proxy. Rate limiting here is an
// availability backstop for a public demo, not a security control — treat it as
// "stops the accidental hammering", never as "stops a determined attacker".
func clientIP(r *http.Request, trustProxy bool) string {
	if trustProxy {
		if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
			if i := strings.IndexByte(xff, ','); i >= 0 {
				xff = xff[:i]
			}
			if ip := strings.TrimSpace(xff); ip != "" {
				return ip
			}
		}
		if ip := strings.TrimSpace(r.Header.Get("X-Real-IP")); ip != "" {
			return ip
		}
	}
	if host, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
		return host
	}
	return r.RemoteAddr
}

// limit wraps a handler with the token bucket when limiting is enabled.
func (a *App) limit(h http.HandlerFunc) http.HandlerFunc {
	if a.limiter == nil {
		return h
	}
	return func(w http.ResponseWriter, r *http.Request) {
		if !a.limiter.allow(clientIP(r, a.trustProxy)) {
			w.Header().Set("Retry-After", "1")
			writeErr(w, &apiError{Status: http.StatusTooManyRequests, Message: "rate limit exceeded; slow down"})
			return
		}
		h(w, r)
	}
}
