// Command example is a minimal end-to-end demo of the Atlas Go SDK.
//
//	go run ./cmd/atlas-server            # terminal 1
//	go run ./sdk/go/example              # terminal 2
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	atlas "github.com/Raj-glitch-max/atlas/sdk/go"
)

func main() {
	c := atlas.New("http://127.0.0.1:8087")
	ctx := context.Background()

	if !c.Health(ctx) {
		log.Fatal("atlas-server not reachable — start it with: go run ./cmd/atlas-server")
	}

	g, err := c.Issue(ctx, atlas.IssueParams{
		Principal: "spiffe://domain-a.test/workload/payments-api",
		Delegate:  "spiffe://domain-b.test/agent/booking-worker",
		Scope:     []string{"read:orders", "write:audit"},
		TTL:       15 * time.Minute,
	})
	if err != nil {
		log.Fatalf("issue: %v", err)
	}
	fmt.Printf("issued   %s  scope=%v\n", g.Instance, g.Scope)

	v, err := c.Verify(ctx, g.Record)
	if err != nil {
		log.Fatalf("verify: %v", err)
	}
	fmt.Printf("verify   %s  (%dµs, %d checks)\n", v.Decision, v.LatencyMicros, len(v.Trace))

	if err := c.Revoke(ctx, g.Instance); err != nil {
		log.Fatalf("revoke: %v", err)
	}
	v2, _ := c.Verify(ctx, g.Record)
	fmt.Printf("revoked  verify -> %s  %v\n", v2.Decision, v2.Causes)
}
