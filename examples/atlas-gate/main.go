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
//	# 2) keep the bundle fresh — the gate's staleness budget is measured against
//	#    the snapshot INSIDE it, so a bundle nobody rewrites ages out within R
//	#    and the gate correctly (but uselessly) denies everything:
//	while :; do go run ./cmd/atlas bundle -o /tmp/bundle.json.tmp \
//	        && mv /tmp/bundle.json.tmp /tmp/bundle.json; sleep 1; done &
//
//	# 3) run the gate in front of any upstream, requiring that scope
//	go run ./examples/atlas-gate -bundle /tmp/bundle.json \
//	        -require-scope vercel:deploy:prod -upstream http://127.0.0.1:9000 -addr :8443
//
//	# 4) callers present the capability; no capability => 401, wrong scope => 403
//	curl -H "Atlas-Capability: $REC" http://127.0.0.1:8443/deploy
//
// If you have no refresher (a one-shot demo), raise the budget instead:
// -max-staleness 10m. Do NOT do that in production — the budget IS your
// revocation blind spot.
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
	"strings"
	"sync/atomic"
	"time"

	"github.com/Raj-glitch-max/atlas/internal/record"
	"github.com/Raj-glitch-max/atlas/internal/revstatus"
	"github.com/Raj-glitch-max/atlas/internal/truststore"
	"github.com/Raj-glitch-max/atlas/internal/verify"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spiffe/go-spiffe/v2/spiffeid"
)

// DenyReason categorizes the reason for denial (R6: fail-visible taxonomy).
type DenyReason string

const (
	DenyReasonMissingCapability   DenyReason = "missing_capability"
	DenyReasonInvalidCapability   DenyReason = "invalid_capability"
	DenyReasonRevoked             DenyReason = "revoked"
	DenyReasonStaleSnapshot       DenyReason = "stale_snapshot"
	DenyReasonScopeMismatch       DenyReason = "scope_mismatch"
	DenyReasonUpstreamUnavailable DenyReason = "upstream_unavailable"
	DenyReasonInternalError       DenyReason = "internal_error"
)

// String returns the reason as a string for the Atlas-Reason header.
func (r DenyReason) String() string { return string(r) }

// Prometheus metrics (R6: fail-visible metrics)
var (
	metricRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "atlas_gate_requests_total",
			Help: "Total number of requests processed by the gate",
		},
		[]string{"result", "reason"},
	)
	metricRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "atlas_gate_request_duration_seconds",
			Help:    "Request processing duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"result"},
	)
	metricVerifierReloadsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "atlas_gate_verifier_reloads_total",
			Help: "Total number of verifier reloads",
		},
		[]string{"result"},
	)
	metricSnapshotAge = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "atlas_gate_snapshot_age_seconds",
			Help: "Age of the currently held revocation snapshot in seconds",
		},
	)
	metricActiveVerifier = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "atlas_gate_active_verifier",
			Help: "Whether the active verifier is healthy (1) or aged out (0)",
		},
	)
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
	reload := flag.Duration("reload", 0, "how often to re-read the bundle from disk (default: max-staleness/4)")
	flag.Parse()

	if *bundlePath == "" || *requireScope == "" {
		log.Fatal("atlas-gate: -bundle and -require-scope are required")
	}
	if *reload <= 0 {
		*reload = *maxStaleness / 4
	}

	v0, err := buildVerifier(*bundlePath, *maxStaleness)
	if err != nil {
		log.Fatalf("atlas-gate: %v", err)
	}

	// The bundle on disk is a SNAPSHOT, and its freshness is what the gate's
	// staleness budget is measured against. A gate that reads it once ages out
	// of its own budget within R and then denies everything — fail-closed, so
	// safe, but useless: with the default R=2s the gate stops admitting traffic
	// two seconds after it starts.
	//
	// So re-read it. An out-of-band refresher (`atlas bundle -o …` on a timer,
	// a sidecar, a mounted secret) rewrites the file; the gate picks it up and
	// swaps the verifier atomically. If a reload fails the previous verifier is
	// kept and simply ages out — the failure mode stays closed, never open.
	var current atomic.Pointer[verifierBundle]
	current.Store(v0)
	go func() {
		t := time.NewTicker(*reload)
		defer t.Stop()
		for range t.C {
			v, err := buildVerifier(*bundlePath, *maxStaleness)
			if err != nil {
				metricVerifierReloadsTotal.WithLabelValues("failure").Inc()
				log.Printf("atlas-gate: bundle reload failed (keeping previous, which will age out and fail closed): %v", err)
				continue
			}
			metricVerifierReloadsTotal.WithLabelValues("success").Inc()
			current.Store(v)
		}
	}()

	// Background metrics updater: update snapshot age and active verifier status
	go func() {
		t := time.NewTicker(5 * time.Second)
		defer t.Stop()
		for range t.C {
			bundle := current.Load()
			if bundle == nil {
				metricActiveVerifier.Set(0)
				metricSnapshotAge.Set(-1)
				continue
			}
			age := time.Since(bundle.snapAsOf).Seconds()
			metricSnapshotAge.Set(age)
			metricActiveVerifier.Set(1)
		}
	}()

	// The protected upstream — or a built-in 200 responder for a standalone demo.
	var proxy http.Handler
	if *upstream != "" {
		u, err := url.Parse(*upstream)
		if err != nil {
			log.Fatalf("atlas-gate: bad -upstream: %v", err)
		}
		proxy = httputil.NewSingleHostReverseProxy(u)
		// Add upstream error handling
		proxy.(*httputil.ReverseProxy).ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
			recordMetrics(DenyReasonUpstreamUnavailable, http.StatusBadGateway, time.Now())
			denyWithReason(w, http.StatusBadGateway, DenyReasonUpstreamUnavailable, "upstream unavailable: "+err.Error())
		}
	} else {
		proxy = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("authorized: the capability grants " + *requireScope + "\n"))
		})
	}

	gate := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := r.Header.Get(*header)
		if rec == "" {
			recordMetrics(DenyReasonMissingCapability, http.StatusUnauthorized, start)
			denyWithReason(w, http.StatusUnauthorized, DenyReasonMissingCapability, "no capability presented (set the "+*header+" header)")
			return
		}
		bundle := current.Load()
		if bundle == nil || bundle.verifier == nil {
			recordMetrics(DenyReasonInternalError, http.StatusServiceUnavailable, start)
			denyWithReason(w, http.StatusServiceUnavailable, DenyReasonInternalError, "verifier not initialized")
			return
		}
		verdict, trace := bundle.verifier.Verify([]byte(rec), spiffeid.ID{})
		if !verdict.IsAccept() {
			// Map verdict to deny reason
			reason := mapVerdictToReason(verdict, trace)
			code := reasonToHTTPCode(reason)
			recordMetrics(reason, code, start)
			denyWithReason(w, code, reason, verdictString(verdict))
			return
		}
		if !recordGrantsScope(rec, *requireScope) {
			recordMetrics(DenyReasonScopeMismatch, http.StatusForbidden, start)
			denyWithReason(w, http.StatusForbidden, DenyReasonScopeMismatch, "capability is valid but does not grant "+*requireScope)
			return
		}
		// Success
		metricRequestsTotal.WithLabelValues("allow", "").Inc()
		metricRequestDuration.WithLabelValues("allow").Observe(time.Since(start).Seconds())
		proxy.ServeHTTP(w, r)
	}

	// Metrics endpoint
	http.Handle("/metrics", promhttp.Handler())

	log.Printf("atlas-gate listening on %s — requiring scope %q, verifying offline against %s (R=%s, reload every %s)",
		*addr, *requireScope, *bundlePath, *maxStaleness, *reload)
	srv := &http.Server{Addr: *addr, Handler: http.HandlerFunc(gate), ReadHeaderTimeout: 5 * time.Second}
	log.Fatal(srv.ListenAndServe())
}

// denyWithReason writes a denial response with the Atlas-Reason header and structured JSON.
func denyWithReason(w http.ResponseWriter, code int, reason DenyReason, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Atlas-Reason", reason.String())
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"error":  msg,
		"reason": reason.String(),
	})
}

// recordMetrics records Prometheus metrics for a request.
func recordMetrics(reason DenyReason, code int, start time.Time) {
	result := "deny"
	if code >= 200 && code < 300 {
		result = "allow"
	}
	metricRequestsTotal.WithLabelValues(result, reason.String()).Inc()
	metricRequestDuration.WithLabelValues(result).Observe(time.Since(start).Seconds())
}

// mapVerdictToReason maps a verification verdict to a deny reason.
func mapVerdictToReason(verdict verify.Verdict, trace verify.DecisionTrace) DenyReason {
	for _, entry := range trace.Entries {
		if entry.Outcome == verify.OutcomeFailDefinitive {
			switch entry.Cause {
			case verify.RevokedObservable:
				return DenyReasonRevoked
			case verify.TrustMaterialAbsent, verify.IntegrityFailed, verify.BindingMismatch,
				verify.ScopeIntegrityFailed, verify.PresenterAbsent, verify.PresenterMismatch:
				return DenyReasonInvalidCapability
			}
		}
		if entry.Outcome == verify.OutcomeInconclusive {
			switch entry.Cause {
			case verify.RevocationKnowledgeStale:
				return DenyReasonStaleSnapshot
			case verify.ClockBeyondTolerance, verify.RevocationStatusIndeterminate:
				return DenyReasonInvalidCapability
			}
		}
	}
	return DenyReasonInvalidCapability
}

// reasonToHTTPCode maps a deny reason to an HTTP status code.
func reasonToHTTPCode(reason DenyReason) int {
	switch reason {
	case DenyReasonMissingCapability:
		return http.StatusUnauthorized // 401
	case DenyReasonInvalidCapability:
		return http.StatusForbidden // 403
	case DenyReasonRevoked:
		return http.StatusForbidden // 403
	case DenyReasonStaleSnapshot:
		return http.StatusServiceUnavailable // 503
	case DenyReasonScopeMismatch:
		return http.StatusForbidden // 403
	case DenyReasonUpstreamUnavailable:
		return http.StatusBadGateway // 502
	case DenyReasonInternalError:
		return http.StatusServiceUnavailable // 503
	default:
		return http.StatusInternalServerError // 500
	}
}

// verdictString returns a human-readable string for the verdict.
func verdictString(verdict verify.Verdict) string {
	var causes []string
	for _, c := range verdict.Causes {
		causes = append(causes, c.String())
	}
	if len(causes) > 0 {
		return "verification failed: " + strings.Join(causes, ", ")
	}
	return "verification failed"
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

// verifierBundle holds a verifier and its associated snapshot AsOf time for metrics.
type verifierBundle struct {
	verifier *verify.Verifier
	snapAsOf time.Time
}

// buildVerifier composes the real Verification Core from a trust bundle — the
// same wiring cmd/atlas/offline.go uses. The revocation snapshot's signature is
// re-verified on ingest, so a tampered bundle is refused here.
func buildVerifier(path string, maxStaleness time.Duration) (*verifierBundle, error) {
	raw, err := record.ReadBundleFile(path)
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
	policy, err := verify.NewPolicy(maxStaleness, 30*time.Second, false)
	if err != nil {
		return nil, err
	}
	v, err := verify.NewVerifier(policy, trust, revAdapter{p: provider}, sysClock{})
	if err != nil {
		return nil, err
	}
	return &verifierBundle{
		verifier: v,
		snapAsOf: b.Revocation.AsOf,
	}, nil
}

var errRefuseBundle = &gateError{"refusing bundle — revocation snapshot signature did not verify"}

type gateError struct{ s string }

func (e *gateError) Error() string { return e.s }
