package main

import (
	"crypto/subtle"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Router builds the HTTP handler for the API. CORS is permissive by design in
// v1 so the local UI (a different origin) can call it; a real deployment
// pins allowed origins via config.
func (a *App) Router() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", a.handleHealth)
	mux.HandleFunc("/version", a.handleVersion)
	mux.HandleFunc("/issue", a.only(http.MethodPost, a.requireAuth(a.handleIssue)))
	mux.HandleFunc("/verify", a.only(http.MethodPost, a.handleVerify))
	mux.HandleFunc("/revoke", a.only(http.MethodPost, a.requireAuth(a.handleRevoke)))
	mux.HandleFunc("/delegations", a.only(http.MethodGet, a.handleDelegations))
	mux.HandleFunc("/audit", a.only(http.MethodGet, a.handleAudit))
	mux.HandleFunc("/graph", a.only(http.MethodGet, a.handleGraph))
	mux.HandleFunc("/stats", a.only(http.MethodGet, a.handleStats))
	mux.HandleFunc("/bundle", a.only(http.MethodGet, a.handleBundle))
	mux.HandleFunc("/metrics", a.only(http.MethodGet, a.handleMetrics))
	return cors(mux)
}

// ---- middleware ----

func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (a *App) only(method string, h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			writeErr(w, &apiError{Status: 405, Message: "method not allowed; expected " + method})
			return
		}
		h(w, r)
	}
}

// requireAuth guards mutating endpoints with a bearer token when one is
// configured; when no key is set the server is open (dev default). Read-only
// endpoints are never guarded. Constant-time compare avoids token timing leaks.
func (a *App) requireAuth(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if a.apiKey != "" {
			got := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
			if subtle.ConstantTimeCompare([]byte(got), []byte(a.apiKey)) != 1 {
				writeErr(w, &apiError{Status: 401, Message: "unauthorized: a valid bearer token is required for this endpoint"})
				return
			}
		}
		h(w, r)
	}
}

// ---- helpers ----

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeErr(w http.ResponseWriter, e *apiError) { writeJSON(w, e.Status, e) }

func decode(r *http.Request, v any) *apiError {
	if r.Body == nil {
		return badRequest("empty request body")
	}
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(v); err != nil {
		return badRequest("invalid JSON: " + err.Error())
	}
	return nil
}

// ---- handlers ----

func (a *App) handleHealth(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, 200, map[string]any{"status": "ok", "time": a.clock.Now().UTC()})
}

func (a *App) handleVersion(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, 200, map[string]any{
		"service":     "atlas-server",
		"version":     version,
		"trustDomain": a.domain.Name(),
		"algorithm":   "ES256 (P-256)",
		"keyId":       a.keyID,
		"publicKey":   a.pubKeyHex,
		"revocationR": a.revWindow.String(),
		"singleHop":   true,
		"federation":  "off",
	})
}

type issueReq struct {
	Principal  string   `json:"principal"`
	Delegate   string   `json:"delegate"`
	Scope      []string `json:"scope"`
	TTLSeconds int      `json:"ttlSeconds"`
}

func (a *App) handleIssue(w http.ResponseWriter, r *http.Request) {
	var req issueReq
	if e := decode(r, &req); e != nil {
		writeErr(w, e)
		return
	}
	if req.Principal == "" || req.Delegate == "" || len(req.Scope) == 0 {
		writeErr(w, badRequest("principal, delegate, and a non-empty scope are required"))
		return
	}
	res, e := a.Issue(req.Principal, req.Delegate, req.Scope, time.Duration(req.TTLSeconds)*time.Second)
	if e != nil {
		writeErr(w, e)
		return
	}
	writeJSON(w, 200, res)
}

type verifyReq struct {
	Record string `json:"record"`
}

func (a *App) handleVerify(w http.ResponseWriter, r *http.Request) {
	var req verifyReq
	if e := decode(r, &req); e != nil {
		writeErr(w, e)
		return
	}
	if req.Record == "" {
		writeErr(w, badRequest("record is required"))
		return
	}
	writeJSON(w, 200, a.Verify(req.Record))
}

type revokeReq struct {
	Instance string `json:"instance"`
}

func (a *App) handleRevoke(w http.ResponseWriter, r *http.Request) {
	var req revokeReq
	if e := decode(r, &req); e != nil {
		writeErr(w, e)
		return
	}
	if req.Instance == "" {
		writeErr(w, badRequest("instance is required"))
		return
	}
	if e := a.Revoke(req.Instance); e != nil {
		writeErr(w, e)
		return
	}
	writeJSON(w, 200, map[string]any{"revoked": true, "instance": req.Instance, "asOf": a.clock.Now().UTC()})
}

func (a *App) handleDelegations(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, 200, map[string]any{"delegations": a.store.Delegations()})
}

func (a *App) handleAudit(w http.ResponseWriter, r *http.Request) {
	limit := 100
	if q := r.URL.Query().Get("limit"); q != "" {
		if n, err := strconv.Atoi(q); err == nil {
			limit = n
		}
	}
	writeJSON(w, 200, map[string]any{"events": a.store.AuditLog(limit)})
}

func (a *App) handleGraph(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, 200, a.store.Graph())
}

// handleStats is a JSON metrics view for dashboards (the counters plus derived
// figures), complementing the Prometheus /metrics.
func (a *App) handleStats(w http.ResponseWriter, _ *http.Request) {
	m := a.store.Snapshot()
	writeJSON(w, 200, map[string]any{
		"issued":            m.Issued,
		"revoked":           m.Revoked,
		"verified":          m.Verified,
		"accept":            m.Accept,
		"reject":            m.Reject,
		"inconclusive":      m.Inconclusive,
		"delegations":       len(a.store.Delegations()),
		"snapshotAgeSecond": a.snapshotAge().Seconds(),
		"trustDomain":       a.domain.Name(),
		"revocationR":       a.revWindow.String(),
	})
}

// handleBundle exports the relying-party trust bundle (trust material + the
// latest signed revocation snapshot) for offline verification.
func (a *App) handleBundle(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, 200, a.Bundle())
}

// handleMetrics emits Prometheus text exposition format — the surface the
// atlas-lab telemetry (Prometheus + Grafana) scrapes.
func (a *App) handleMetrics(w http.ResponseWriter, _ *http.Request) {
	m := a.store.Snapshot()
	age := a.snapshotAge().Seconds()
	w.Header().Set("Content-Type", "text/plain; version=0.0.4; charset=utf-8")
	fmt.Fprintf(w, "# HELP atlas_verify_verdicts_total Verification verdicts by decision.\n")
	fmt.Fprintf(w, "# TYPE atlas_verify_verdicts_total counter\n")
	fmt.Fprintf(w, "atlas_verify_verdicts_total{decision=\"accept\"} %d\n", m.Accept)
	fmt.Fprintf(w, "atlas_verify_verdicts_total{decision=\"reject\"} %d\n", m.Reject)
	fmt.Fprintf(w, "atlas_verify_verdicts_total{decision=\"inconclusive\"} %d\n", m.Inconclusive)
	fmt.Fprintf(w, "# HELP atlas_issued_total Delegations issued.\n# TYPE atlas_issued_total counter\natlas_issued_total %d\n", m.Issued)
	fmt.Fprintf(w, "# HELP atlas_revoked_total Revocations recorded.\n# TYPE atlas_revoked_total counter\natlas_revoked_total %d\n", m.Revoked)
	fmt.Fprintf(w, "# HELP atlas_verified_total Verifications performed.\n# TYPE atlas_verified_total counter\natlas_verified_total %d\n", m.Verified)
	fmt.Fprintf(w, "# HELP atlas_revocation_snapshot_age_seconds Age of the held signed revocation snapshot.\n")
	fmt.Fprintf(w, "# TYPE atlas_revocation_snapshot_age_seconds gauge\natlas_revocation_snapshot_age_seconds %.3f\n", age)

	h := a.store.LatencySnapshot()
	fmt.Fprintf(w, "# HELP atlas_verify_latency_seconds Verification latency.\n")
	fmt.Fprintf(w, "# TYPE atlas_verify_latency_seconds histogram\n")
	for i, ub := range h.Bounds {
		fmt.Fprintf(w, "atlas_verify_latency_seconds_bucket{le=\"%g\"} %d\n", ub, h.Cumulative[i])
	}
	fmt.Fprintf(w, "atlas_verify_latency_seconds_bucket{le=\"+Inf\"} %d\n", h.Total)
	fmt.Fprintf(w, "atlas_verify_latency_seconds_sum %g\n", h.Sum)
	fmt.Fprintf(w, "atlas_verify_latency_seconds_count %d\n", h.Total)
}

// checkpoint: refactor(stores): refactor signature validation

// checkpoint: chore(test): refine Fuzz Verification core target (#149)

// checkpoint: chore(fuzz): optimize Fuzz Verification core target
