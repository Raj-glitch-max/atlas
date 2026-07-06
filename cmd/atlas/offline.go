package main

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/Raj-glitch-max/atlas/internal/record"
	"github.com/Raj-glitch-max/atlas/internal/revstatus"
	"github.com/Raj-glitch-max/atlas/internal/truststore"
	"github.com/Raj-glitch-max/atlas/internal/verify"
	"github.com/spiffe/go-spiffe/v2/spiffeid"
)

// Offline verification: `atlas verify --offline --bundle trust.json <record|->`
// composes the REAL Verification Core locally — no server, no network. This is
// the product's core claim made physical: kill the network and verification
// still answers; it fails closed once the bundle's revocation snapshot is
// older than your own staleness budget (--max-staleness = the relying party's
// policy R).
//
// The bundle is tamper-evident: the revocation snapshot's ECDSA signature is
// re-verified on import (a doctored bundle is refused), and record signatures
// are checked against the bundled trust material as always.

// bundleFile mirrors cmd/atlas-server's BundleDTO wire shape.
type bundleFile struct {
	Version     int               `json:"version"`
	TrustDomain string            `json:"trustDomain"`
	Keys        map[string]string `json:"keys"`
	Revocation  struct {
		ListID  string    `json:"listId"`
		AsOf    time.Time `json:"asOf"`
		Revoked []string  `json:"revoked"`
		Sig     []byte    `json:"sig"`
	} `json:"revocation"`
	ExportedAt time.Time `json:"exportedAt"`
}

// offlineRevAdapter bridges revstatus.Provider onto verify's port — the same
// composition glue every Atlas relying party wires (AD-020).
type offlineRevAdapter struct{ p revstatus.Provider }

func (a offlineRevAdapter) StatusOf(instance record.InstanceID) verify.RevocationStatus {
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

// cmdBundle fetches the trust bundle from the server and writes it to a file.
func cmdBundle(out io.Writer, c *client, args []string) int {
	fs := flag.NewFlagSet("bundle", flag.ContinueOnError)
	fs.SetOutput(out)
	output := fs.String("o", "trust-bundle.json", "output file")
	if err := fs.Parse(args); err != nil {
		return exitUsage
	}
	var raw json.RawMessage
	if err := c.do(http.MethodGet, "/bundle", nil, &raw); err != nil {
		fmt.Fprintln(out, "bundle failed:", err)
		return exitError
	}
	var pretty map[string]any
	_ = json.Unmarshal(raw, &pretty)
	buf, _ := json.MarshalIndent(pretty, "", "  ")
	if err := os.WriteFile(*output, buf, 0o600); err != nil {
		fmt.Fprintln(out, "bundle: write:", err)
		return exitError
	}
	var b bundleFile
	_ = json.Unmarshal(raw, &b)
	fmt.Fprintf(out, "wrote %s\n", *output)
	fmt.Fprintf(out, "  trust domain   %s\n", b.TrustDomain)
	fmt.Fprintf(out, "  keys           %d\n", len(b.Keys))
	fmt.Fprintf(out, "  revoked        %d (snapshot as-of %s)\n", len(b.Revocation.Revoked), b.Revocation.AsOf.Format(time.RFC3339))
	fmt.Fprintln(out, "\nverify offline with:\n  atlas verify --offline --bundle", *output, "<record|->")
	return exitOK
}

// verifyOffline runs the real Verification Core against a record using only
// the bundle — no server, no network. When requireScope is non-empty it also
// applies the authorization gate (valid AND grants the action).
func verifyOffline(out io.Writer, bundlePath, rec string, maxStaleness time.Duration, requireScope string) int {
	raw, err := os.ReadFile(bundlePath)
	if err != nil {
		fmt.Fprintln(out, "offline verify: read bundle:", err)
		return exitError
	}
	var b bundleFile
	if err := json.Unmarshal(raw, &b); err != nil {
		fmt.Fprintln(out, "offline verify: bad bundle JSON:", err)
		return exitError
	}
	if b.Version != 1 {
		fmt.Fprintf(out, "offline verify: unsupported bundle version %d\n", b.Version)
		return exitError
	}

	// 1. Trust material from the bundled public keys.
	td, err := spiffeid.TrustDomainFromString(b.TrustDomain)
	if err != nil {
		fmt.Fprintln(out, "offline verify: bundle trust domain:", err)
		return exitError
	}
	keys := map[string]*ecdsa.PublicKey{}
	for kid, b64 := range b.Keys {
		der, err := base64.StdEncoding.DecodeString(b64)
		if err != nil {
			fmt.Fprintf(out, "offline verify: key %s: %v\n", kid, err)
			return exitError
		}
		pub, err := x509.ParsePKIXPublicKey(der)
		if err != nil {
			fmt.Fprintf(out, "offline verify: key %s: %v\n", kid, err)
			return exitError
		}
		ek, ok := pub.(*ecdsa.PublicKey)
		if !ok {
			fmt.Fprintf(out, "offline verify: key %s is not ECDSA\n", kid)
			return exitError
		}
		keys[kid] = ek
	}
	tm, err := record.NewTrustMaterial(td, keys)
	if err != nil {
		fmt.Fprintln(out, "offline verify: trust material:", err)
		return exitError
	}
	trust := truststore.New()
	if err := trust.Provision(tm, time.Now()); err != nil {
		fmt.Fprintln(out, "offline verify: provision:", err)
		return exitError
	}

	// 2. Revocation snapshot — signature re-verified on Ingest, so a tampered
	//    bundle is refused here, not silently trusted. Any bundled key may be
	//    the snapshot signer (key rotation), so each is tried; refusal only if
	//    none verifies.
	revoked := make([]record.InstanceID, 0, len(b.Revocation.Revoked))
	for _, s := range b.Revocation.Revoked {
		id, err := record.InstanceIDFromString(s)
		if err != nil {
			fmt.Fprintln(out, "offline verify: bundle revoked entry:", err)
			return exitError
		}
		revoked = append(revoked, id)
	}
	snap := revstatus.SignedRevokedSet{
		ListID: b.Revocation.ListID, AsOf: b.Revocation.AsOf,
		Revoked: revoked, Sig: b.Revocation.Sig,
	}
	var provider *revstatus.SignedSetProvider
	var ingestErr error
	for _, pub := range keys {
		p := revstatus.NewSignedSetProvider(pub, b.Revocation.ListID)
		if _, err := p.Ingest(snap); err == nil {
			provider = p
			break
		} else {
			ingestErr = err
		}
	}
	if provider == nil {
		fmt.Fprintln(out, "offline verify: REFUSING bundle —", ingestErr)
		return exitError
	}

	// 3. The real Verification Core, composed locally. --max-staleness is the
	//    relying party's own freshness policy (R).
	policy, err := verify.NewPolicy(maxStaleness, 30*time.Second)
	if err != nil {
		fmt.Fprintln(out, "offline verify: policy:", err)
		return exitError
	}
	v, err := verify.NewVerifier(policy, trust, offlineRevAdapter{p: provider}, sysClock{})
	if err != nil {
		fmt.Fprintln(out, "offline verify:", err)
		return exitError
	}
	start := time.Now()
	verdict, trace := v.Verify([]byte(rec))
	elapsed := time.Since(start)

	label := "INCONCLUSIVE"
	switch {
	case verdict.IsAccept():
		label = "ACCEPT"
	case verdict.Decision == verify.Reject:
		label = "REJECT"
	}
	fmt.Fprintf(out, "%s  (offline · %s · snapshot age %s / budget %s)",
		label, elapsed.Round(time.Microsecond), time.Since(b.Revocation.AsOf).Round(time.Millisecond), maxStaleness)
	if len(verdict.Causes) > 0 {
		cs := make([]string, 0, len(verdict.Causes))
		for _, c := range verdict.Causes {
			cs = append(cs, c.String())
		}
		fmt.Fprintf(out, "  %s", strings.Join(cs, ", "))
	}
	fmt.Fprintln(out)
	tw := tabwriter.NewWriter(out, 0, 2, 2, ' ', 0)
	for _, e := range trace.Entries {
		cause := e.Cause.String()
		if cause == "None" || cause == "" {
			cause = "—"
		}
		fmt.Fprintf(tw, "  %s\t%s\t%s\n", string(e.Check), e.Outcome.String(), cause)
	}
	tw.Flush()
	return scopeGate(out, rec, requireScope, verdict.IsAccept())
}

type sysClock struct{}

func (sysClock) Now() time.Time { return time.Now() }

// checkpoint: chore(style): test scroll animation trigger (#249)
