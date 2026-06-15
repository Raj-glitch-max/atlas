package harness

import (
	"sync"
	"time"
)

// Clock is a controllable TimePort for tests: it makes verification
// deterministic under an injected clock (AT8 clock-skew scenarios) instead of
// depending on wall time. It satisfies both the verify and issuance time-port
// contracts structurally (a single Now() time.Time method; AD-014).
type Clock struct {
	mu  sync.Mutex
	now time.Time
}

// NewClock returns a clock reading exactly t.
func NewClock(t time.Time) *Clock { return &Clock{now: t} }

// Now returns the current reading.
func (c *Clock) Now() time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.now
}

// Set moves the clock to t.
func (c *Clock) Set(t time.Time) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.now = t
}

// Advance moves the clock forward by d.
func (c *Clock) Advance(d time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.now = c.now.Add(d)
}
