// Command atlas-verify is the relying-party verification-boundary driver: the
// composition root that constructs the Verification Core (internal/verify)
// with its policy, wires the trust-material, revocation-status, and time
// ports to their providers, and carries the AT26 latency-measurement point —
// around Verify, never inside it. Drivers contain wiring only; a
// delegation-logic check in a driver would be a conformance violation (FD-9).
//
// V1 form: a self-contained demonstration that composes the real modules end
// to end and prints the verdict and decision trace. Production I/O — reading
// a presented record, loading trust material, and policy from configuration —
// is deferred until a serialization/config format is chosen (TP6); the
// composition wiring it would use is exactly what this demonstrates.
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
	"github.com/Raj-glitch-max/atlas/internal/revstatus"
	"github.com/Raj-glitch-max/atlas/internal/truststore"
	"github.com/Raj-glitch-max/atlas/internal/verify"
	"github.com/spiffe/go-spiffe/v2/spiffeid"
)

// systemClock is the production TimePort: the real wall clock. Composition-
// root concern; tests inject a controllable clock instead.
type systemClock struct{}

func (systemClock) Now() time.Time { return time.Now() }

// demoPermissions is an inline permission source for the demonstration.
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

// revocationAdapter bridges a revstatus.Provider onto verify's port (AD-020),
// the composition-root glue a production relying party wires.
type revocationAdapter struct{ p revstatus.Provider }

func (a revocationAdapter) StatusOf(instance record.InstanceID) verify.RevocationStatus {
	ans := a.p.StatusOf(instance)
	var st verify.RevocationState
	switch ans.State {
	case revstatus.NotObservedRevoked:
		st = verify.NotObservedRevoked
	case revstatus.ObservablyRevoked:
		st = verify.ObservablyRevoked
	default:
		st = verify.Indeterminate
	}
	return verify.RevocationStatus{State: st, AsOf: ans.AsOf}
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "atlas-verify:", err)
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
	const keyID = "authority-key-1"

	authority, err := issuance.NewAuthority(
		record.Signer{Key: key, KeyID: keyID},
		demoPermissions{principal: principal, perms: issuance.NewPermissionSet("read:orders", "write:audit", "admin:all")},
		issuance.NoBinding{}, issuance.RandomMinter{}, clock)
	if err != nil {
		return err
	}
	res, err := authority.Issue(issuance.Request{
		Principal: principal, Delegate: delegate,
		Scope: []string{"read:orders", "write:audit"}, Expiration: clock.Now().Add(time.Hour),
	})
	if err != nil {
		return err
	}
	if res.Outcome != issuance.Issued {
		return fmt.Errorf("issuance refused: %s", res.Refusal)
	}

	tm, err := record.NewTrustMaterial(principal.TrustDomain(), map[string]*ecdsa.PublicKey{keyID: &key.PublicKey})
	if err != nil {
		return err
	}
	trust := truststore.New()
	if err := trust.Provision(tm, clock.Now()); err != nil {
		return err
	}
	policy, err := verify.NewPolicy(time.Minute, 30*time.Second)
	if err != nil {
		return err
	}

	// Demonstrate both paths a relying party sees: a fresh not-revoked
	// observation (Accept) and the pre-spike degenerate provider that fails
	// closed (InconclusiveRejected).
	fresh := staticRevocation{state: verify.NotObservedRevoked, clock: clock}
	degenerate := revocationAdapter{p: revstatus.NewDegenerate()}

	for _, scenario := range []struct {
		name string
		port verify.RevocationStatusPort
	}{
		{"fresh not-revoked observation", fresh},
		{"degenerate provider (pre-spike default, fails closed)", degenerate},
	} {
		v, err := verify.NewVerifier(policy, trust, scenario.port, clock)
		if err != nil {
			return err
		}
		start := time.Now()
		verdict, trace := v.Verify(res.Record.Presented()) // AT26 measurement point: around Verify
		elapsed := time.Since(start)

		fmt.Printf("scenario: %s\n", scenario.name)
		fmt.Printf("  verdict: %s\n", verdict.Decision)
		for _, e := range trace.Entries {
			fmt.Printf("    %-18s %-14s %s\n", e.Check, e.Outcome, e.Cause)
		}
		fmt.Printf("  verification latency (measured, no threshold asserted in V1): %s\n\n", elapsed)
	}

	// Scenario 3: the real revocation alpha-path (signed-revoked-set, E7),
	// under the resolved scope act (R = 2s). Domain A signs a timestamped
	// snapshot of revoked instances; the RP verifies it against domain A's
	// public key and answers offline with verifiable freshness. This shows the
	// system actually observing a revocation, not a fake.
	raKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return err
	}
	publisher, err := revstatus.NewPublisher(raKey, "demo-revlist")
	if err != nil {
		return err
	}
	provider := revstatus.NewSignedSetProvider(&raKey.PublicKey, "demo-revlist")
	realPolicy, err := verify.NewPolicy(2*time.Second, 30*time.Second)
	if err != nil {
		return err
	}
	realVerifier, err := verify.NewVerifier(realPolicy, trust, revocationAdapter{p: provider}, clock)
	if err != nil {
		return err
	}
	instance := res.Record.Read().Instance

	fmt.Println("scenario: real revocation alpha-path (signed-revoked-set, R=2s)")
	// (a) publish a fresh empty snapshot -> the delegation verifies Accept
	empty, err := publisher.Publish(nil, clock.Now())
	if err != nil {
		return err
	}
	if _, err := provider.Ingest(empty); err != nil {
		return err
	}
	before, _ := realVerifier.Verify(res.Record.Presented())
	fmt.Printf("  before revocation: %s\n", before.Decision)

	// (b) revoke the instance -> publish a newer signed snapshot -> Reject
	revoked, err := publisher.Publish([]record.InstanceID{instance}, clock.Now().Add(time.Second))
	if err != nil {
		return err
	}
	if _, err := provider.Ingest(revoked); err != nil {
		return err
	}
	after, trace := realVerifier.Verify(res.Record.Presented())
	fmt.Printf("  after revocation:  %s %v\n", after.Decision, after.Causes)
	for _, e := range trace.Entries {
		fmt.Printf("    %-18s %-14s %s\n", e.Check, e.Outcome, e.Cause)
	}
	fmt.Println("  (freshness is verifiable: the as-of is the snapshot's signed timestamp)")
	return nil
}

// staticRevocation reports a fixed knowledge state fresh as of the clock, for
// the demonstration's accept path.
type staticRevocation struct {
	state verify.RevocationState
	clock interface{ Now() time.Time }
}

func (s staticRevocation) StatusOf(record.InstanceID) verify.RevocationStatus {
	return verify.RevocationStatus{State: s.state, AsOf: s.clock.Now()}
}
