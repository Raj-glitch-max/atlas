// Command atlas-gate is a reference authorization gate: an HTTP reverse proxy
// that sits in front of a protected service and lets a request through only if
// it carries a valid Atlas capability that grants the required scope. It
// verifies the capability OFFLINE — composing the same Verification Core the
// CLI's `verify --offline` uses — so the gate never calls the issuer.
//
// This is the "plug Atlas into your stack" pattern as a runnable, deployable
// service (as opposed to examples/ship-a-landing-page.sh, which demonstrates
// the same idea inline). A real tool adapter is exactly this: verify the
// capability, and only then act with the tool's own credential.
//
//	# 1) run the server, issue a capability, export the trust bundle
//	go run ./cmd/atlas-server -grant 'vercel:deploy:prod' &
//	REC=$(go run ./cmd/atlas -q delegate --principal spiffe://domain-a.test/human/you \
//	        --delegate spiffe://domain-b.test/agent/deployer --scope vercel:deploy:prod)
//	go run ./cmd/atlas bundle -o /tmp/bundle.json
//
//	# 2) run the gate in front of any upstream, requiring that scope
//	go run ./examples/atlas-gate -bundle /tmp/bundle.json \
//	        -require-scope vercel:deploy:prod -upstream http://127.0.0.1:9000 -addr :8443
//
//	# 3) callers present the capability; no capability => 401, wrong scope => 403
//	curl -H "Atlas-Capability: $REC" http://127.0.0.1:8443/deploy
package main

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/Raj-glitch-max/atlas/internal/record"
	"github.com/Raj-glitch-max/atlas/internal/revstatus"
	"github.com/Raj-glitch-max/atlas/internal/truststore"
	"github.com/Raj-glitch-max/atlas/internal/verify"
	"github.com/spiffe/go-spiffe/v2/spiffeid"
)

// bundleFile mirrors the /bundle wire shape (see cmd/atlas/offline.go).
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
}

type revAdapter struct{ p revstatus.Provider }

func (a revAdapter) StatusOf(instance record.InstanceID) verify.RevocationStatus {
	ans := a.p.StatusOf(instance)
	st := verify.Indeterminate
	switch ans.State {
	case revstatus.NotObservedRevoked:
		st = verify.NotObservedRevoked
	case revstatus.ObservablyRevoked:
		st = verify.ObservablyRevoked
	}
	return verify.RevocationStatus{State: st, AsOf: ans.AsOf}
}

type sysClock struct{}

func (sysClock) Now() time.Time { return time.Now() }

func main() {
	addr := flag.String("addr", ":8443", "listen address for the gate")
	bundlePath := flag.String("bundle", "", "path to the Atlas trust bundle (required)")
	requireScope := flag.String("require-scope", "", "scope a capability must grant to pass (required)")
	upstream := flag.String("upstream", "", "upstream URL to proxy authorized requests to (default: reply 200)")
	header := flag.String("header", "Atlas-Capability", "request header carrying the presented record")
	maxStaleness := flag.Duration("max-staleness", 2*time.Second, "relying-party freshness budget (R)")
	flag.Parse()

	if *bundlePath == "" || *requireScope == "" {
		log.Fatal("atlas-gate: -bundle and -require-scope are required")
	}

	verifier, err := buildVerifier(*bundlePath, *maxStaleness)
	if err != nil {
		log.Fatalf("atlas-gate: %v", err)
	}

	// The protected upstream — or a built-in 200 responder for a standalone demo.
	var proxy http.Handler
	if *upstream != "" {
		u, err := url.Parse(*upstream)
		if err != nil {
			log.Fatalf("atlas-gate: bad -upstream: %v", err)
		}
		proxy = httputil.NewSingleHostReverseProxy(u)
	} else {
		proxy = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("authorized: the capability grants " + *requireScope + "\n"))
		})
	}

	gate := func(w http.ResponseWriter, r *http.Request) {
		rec := r.Header.Get(*header)
		if rec == "" {
			deny(w, http.StatusUnauthorized, "no capability presented (set the "+*header+" header)")
			return
		}
		verdict, _ := verifier.Verify([]byte(rec))
		if !verdict.IsAccept() {
			deny(w, http.StatusForbidden, "capability did not verify (not accepted / revoked / stale)")
			return
		}
		if !recordGrantsScope(rec, *requireScope) {
			deny(w, http.StatusForbidden, "capability is valid but does not grant "+*requireScope)
			return
		}
		proxy.ServeHTTP(w, r)
	}

	log.Printf("atlas-gate listening on %s — requiring scope %q, verifying offline against %s",
		*addr, *requireScope, *bundlePath)
	srv := &http.Server{Addr: *addr, Handler: http.HandlerFunc(gate), ReadHeaderTimeout: 5 * time.Second}
	log.Fatal(srv.ListenAndServe())
}

func deny(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

// recordGrantsScope decodes the record's scope claim (offline) and reports
// whether it contains the required scope. Mirrors cmd/atlas's scope gate; only
// meaningful after the record has been accepted.
func recordGrantsScope(rec, required string) bool {
	parts := strings.Split(strings.TrimSpace(rec), ".")
	if len(parts) != 3 {
		return false
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return false
	}
	var claims struct {
		Scope []string `json:"scope"`
	}
	if json.Unmarshal(payload, &claims) != nil {
		return false
	}
	for _, s := range claims.Scope {
		if s == required {
			return true
		}
	}
	return false
}

// buildVerifier composes the real Verification Core from a trust bundle — the
// same wiring cmd/atlas/offline.go uses. The revocation snapshot's signature is
// re-verified on ingest, so a tampered bundle is refused here.
func buildVerifier(path string, maxStaleness time.Duration) (*verify.Verifier, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var b bundleFile
	if err := json.Unmarshal(raw, &b); err != nil {
		return nil, err
	}
	td, err := spiffeid.TrustDomainFromString(b.TrustDomain)
	if err != nil {
		return nil, err
	}
	keys := map[string]*ecdsa.PublicKey{}
	for kid, b64 := range b.Keys {
		der, err := base64.StdEncoding.DecodeString(b64)
		if err != nil {
			return nil, err
		}
		pub, err := x509.ParsePKIXPublicKey(der)
		if err != nil {
			return nil, err
		}
		ek, ok := pub.(*ecdsa.PublicKey)
		if !ok {
			continue
		}
		keys[kid] = ek
	}
	tm, err := record.NewTrustMaterial(td, keys)
	if err != nil {
		return nil, err
	}
	trust := truststore.New()
	if err := trust.Provision(tm, time.Now()); err != nil {
		return nil, err
	}
	revoked := make([]record.InstanceID, 0, len(b.Revocation.Revoked))
	for _, s := range b.Revocation.Revoked {
		id, err := record.InstanceIDFromString(s)
		if err != nil {
			return nil, err
		}
		revoked = append(revoked, id)
	}
	snap := revstatus.SignedRevokedSet{
		ListID: b.Revocation.ListID, AsOf: b.Revocation.AsOf,
		Revoked: revoked, Sig: b.Revocation.Sig,
	}
	var provider *revstatus.SignedSetProvider
	for _, pub := range keys {
		p := revstatus.NewSignedSetProvider(pub, b.Revocation.ListID)
		if _, err := p.Ingest(snap); err == nil {
			provider = p
			break
		}
	}
	if provider == nil {
		return nil, errRefuseBundle
	}
	policy, err := verify.NewPolicy(maxStaleness, 30*time.Second)
	if err != nil {
		return nil, err
	}
	return verify.NewVerifier(policy, trust, revAdapter{p: provider}, sysClock{})
}

var errRefuseBundle = &gateError{"refusing bundle — revocation snapshot signature did not verify"}

type gateError struct{ s string }

func (e *gateError) Error() string { return e.s }
