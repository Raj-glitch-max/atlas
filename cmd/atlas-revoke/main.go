// Command atlas-revoke is the revocation-act driver: the composition root
// that wires the Revocation Origin register (internal/revorigin). Drivers
// contain wiring only (FD-9).
//
// V1 form: a self-contained demonstration that records revocations and prints
// the append-only register. Production I/O and the propagation channel that
// carries revocations to relying-party views (S2/S3, spike-selected) are
// deferred (E7).
package main

import (
	"fmt"
	"os"
	"time"

	"github.com/Raj-glitch-max/atlas/internal/record"
	"github.com/Raj-glitch-max/atlas/internal/revorigin"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "atlas-revoke:", err)
		os.Exit(1)
	}
}

func run() error {
	register := revorigin.New()
	now := time.Now()

	for _, id := range []string{"inst-alpha", "inst-beta", "inst-alpha"} { // note the repeat
		instance, err := record.InstanceIDFromString(id)
		if err != nil {
			return err
		}
		register.Revoke(instance, now)
	}

	fmt.Printf("revocation register (%d entries, append-only, terminal):\n", register.Len())
	for i, e := range register.View() {
		fmt.Printf("  %d: %s @ %s\n", i, e.Instance, e.RecordedAt.Format(time.RFC3339))
	}
	fmt.Println("  (the repeated revocation of inst-alpha is a no-op: revocation is one-way and terminal)")
	return nil
}
