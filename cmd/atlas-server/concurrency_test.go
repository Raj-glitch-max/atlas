package main

import (
	"fmt"
	"net/http"
	"sync"
	"testing"
	"time"
)

// TestConcurrentSessionsAndRefresh is the test the suite was missing.
//
// `go test -race` only reports races on code paths that ACTUALLY execute
// concurrently during the run, and every other test in this package is
// sequential — so a clean race run proved almost nothing. This drives the
// genuinely shared structures at once: the session map (create/evict/lookup),
// each session's revoked set and snapshot provider, the cross-session activity
// ring buffer, the rate limiter's bucket map, and the background refresher
// republishing every live session's snapshot underneath all of it.
func TestConcurrentSessionsAndRefresh(t *testing.T) {
	app, err := NewApp(Config{Domain: "domain-a.test", RateLimitRPM: 100000}, systemClock{})
	if err != nil {
		t.Fatal(err)
	}
	h := app.Router()

	// The background refresher, as main() runs it: republishes every live
	// session's snapshot while requests are mutating the session set.
	stop := make(chan struct{})
	var refresher sync.WaitGroup
	refresher.Add(1)
	go func() {
		defer refresher.Done()
		tick := time.NewTicker(time.Millisecond)
		defer tick.Stop()
		for {
			select {
			case <-stop:
				return
			case <-tick.C:
				_ = app.Refresh()
			}
		}
	}()

	const workers, perWorker = 24, 40
	var wg sync.WaitGroup
	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func(w int) {
			defer wg.Done()
			// Distinct session per worker, plus deliberate collisions on a
			// shared id so concurrent lookup-or-create races on the same key.
			own := fmt.Sprintf("%032x", w)
			shared := "cccccccccccccccccccccccccccccccc"
			for i := 0; i < perWorker; i++ {
				sess := own
				if i%3 == 0 {
					sess = shared
				}
				code, body := do(t, h, http.MethodPost, "/issue", sess,
					`{"principal":"spiffe://domain-a.test/workload/payments-api",`+
						`"delegate":"spiffe://domain-b.test/agent/booking-worker",`+
						`"scope":["read:orders"],"ttlSeconds":3600}`)
				if code != 200 {
					continue // rate limited or evicted mid-flight; not a race
				}
				rec, _ := body["record"].(string)
				inst, _ := body["instance"].(string)

				do(t, h, http.MethodPost, "/verify", sess, `{"record":"`+rec+`"}`)
				do(t, h, http.MethodPost, "/revoke", sess, `{"instance":"`+inst+`"}`)
				do(t, h, http.MethodPost, "/verify", sess, `{"record":"`+rec+`"}`)
				do(t, h, http.MethodGet, "/activity", sess, "")
				do(t, h, http.MethodGet, "/delegations", sess, "")
				do(t, h, http.MethodGet, "/bundle", sess, "")
			}
		}(w)
	}
	wg.Wait()
	close(stop)
	refresher.Wait()

	// The instance must still be answerable after all that churn.
	if n := app.SessionCount(); n > maxSessions {
		t.Fatalf("session cap breached under concurrency: %d > %d", n, maxSessions)
	}
	code, _ := do(t, h, http.MethodGet, "/activity", "", "")
	if code != 200 {
		t.Fatalf("/activity broken after concurrent load: %d", code)
	}
}
