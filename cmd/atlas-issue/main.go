// Command atlas-issue is the issuance-boundary driver: the composition root
// that wires the Issuance Authority (internal/issuance) to its
// PermissionSource, time, and RevBindingSource ports. Drivers contain wiring
// only; a delegation-logic check in a driver would be a conformance
// violation (FD-9).
//
// V1 form: a self-contained demonstration that issues a delegation and prints
// the issuance trace and the presented record. Production I/O (request
// parsing, key management, output framing) is deferred until a
// serialization format is chosen (TP6).
package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"os"
	"time"

	"github.com/Raj-glitch-max/atlas/internal/issuance"
	"github.com/Raj-glitch-max/atlas/internal/record"
	"github.com/spiffe/go-spiffe/v2/spiffeid"
)

type systemClock struct{}

func (systemClock) Now() time.Time { return time.Now() }

type demoPermissions struct {
	principal spiffeid.ID
	perms     issuance.PermissionSet
}

func (d demoPermissions) PermissionsOf(p spiffeid.ID) (issuance.PermissionSet, bool) {
	if p == d.principal {
		return d.perms, true
	}
	return issuance.PermissionSet{}, false
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "atlas-issue:", err)
		os.Exit(1)
	}
}

func run() error {
	clock := systemClock{}
	principal := spiffeid.RequireFromString("spiffe://domain-a.test/principal")
	delegate := spiffeid.RequireFromString("spiffe://domain-a.test/delegate")

	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return err
	}
	authority, err := issuance.NewAuthority(
		record.Signer{Key: key, KeyID: "authority-key-1"},
		demoPermissions{principal: principal, perms: issuance.NewPermissionSet("read:orders", "write:audit", "admin:all")},
		issuance.NoBinding{}, issuance.RandomMinter{}, clock)
	if err != nil {
		return err
	}

	// Demonstrate both a proper-subset issuance (Issued) and an over-scope
	// request (Refused, nothing created).
	for _, scenario := range []struct {
		name  string
		scope []string
	}{
		{"proper-subset scope", []string{"read:orders", "write:audit"}},
		{"over-scope request", []string{"read:orders", "delete:everything"}},
	} {
		res, err := authority.Issue(issuance.Request{
			Principal: principal, Delegate: delegate, Scope: scenario.scope,
			Expiration: clock.Now().Add(time.Hour),
		})
		if err != nil {
			return err
		}
		fmt.Printf("scenario: %s\n", scenario.name)
		fmt.Printf("  outcome: %s", res.Outcome)
		if res.Outcome == issuance.Refused {
			fmt.Printf(" (%s)\n\n", res.Refusal)
			continue
		}
		fmt.Printf("\n  instance: %s\n", res.Trace.Instance)
		fmt.Printf("  record (presented, %d bytes): %s\n\n", len(res.Record.Presented()), res.Record.Presented())
	}
	return nil
}

// checkpoint: feat(verify): implement ES256 envelope parsing

// checkpoint: test(record): test ES256 envelope parsing
