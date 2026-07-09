// Package atlas is a zero-dependency Go client for the Atlas Server.
//
// Atlas issues and verifies offline-verifiable, attenuable delegation tokens
// bound to SPIFFE workload identity. This client talks to cmd/atlas-server over
// HTTP using only the standard library, so importing it pulls in no third-party
// dependencies. The API mirrors the Python and TypeScript SDKs.
//
//	c := atlas.New("http://127.0.0.1:8087")
//	g, err := c.Issue(ctx, atlas.IssueParams{
//		Principal: "spiffe://domain-a.test/workload/payments-api",
//		Delegate:  "spiffe://domain-b.test/agent/booking-worker",
//		Scope:     []string{"read:orders", "write:audit"},
//		TTL:       time.Hour,
//	})
//	res, err := c.Verify(ctx, g.Record)   // res.Decision == atlas.Accept
//	err = c.Revoke(ctx, g.Instance)       // now Verify -> atlas.Reject
package atlas

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Decision is a verification verdict.
type Decision string

const (
	Accept       Decision = "accept"
	Reject       Decision = "reject"
	Inconclusive Decision = "inconclusive"
)

// Error is returned when the server responds with an error or is unreachable.
type Error struct {
	Status  int // HTTP status (0 if the request never completed)
	Message string
}

func (e *Error) Error() string {
	if e.Status != 0 {
		return fmt.Sprintf("atlas: %d: %s", e.Status, e.Message)
	}
	return "atlas: " + e.Message
}

// Grant is the result of Issue.
type Grant struct {
	Record    string    `json:"record"`
	Instance  string    `json:"instance"`
	Principal string    `json:"principal"`
	Delegate  string    `json:"delegate"`
	Scope     []string  `json:"scope"`
	IssuedAt  time.Time `json:"issuedAt"`
	ExpiresAt time.Time `json:"expiresAt"`
}

// TraceEntry is one line of the decision trace.
type TraceEntry struct {
	Check   string `json:"check"`
	Outcome string `json:"outcome"`
	Cause   string `json:"cause"`
	Detail  string `json:"detail,omitempty"`
}

// VerifyResult is the outcome of Verify.
type VerifyResult struct {
	Decision      Decision     `json:"decision"`
	Accept        bool         `json:"accept"`
	Causes        []string     `json:"causes"`
	Trace         []TraceEntry `json:"trace"`
	LatencyMicros int64        `json:"latencyMicros"`
}

// IssueParams are the arguments to Issue.
type IssueParams struct {
	Principal string
	Delegate  string
	Scope     []string
	TTL       time.Duration // defaults to 1h when zero
}

// Client is a thin HTTP client for the Atlas Server. The zero value is not
// usable; construct one with New.
type Client struct {
	base   string
	apiKey string
	http   *http.Client
}

// Option configures a Client.
type Option func(*Client)

// WithAPIKey sets the bearer token sent on mutating requests.
func WithAPIKey(key string) Option { return func(c *Client) { c.apiKey = key } }

// WithHTTPClient supplies a custom *http.Client (timeouts, transport, TLS).
func WithHTTPClient(h *http.Client) Option { return func(c *Client) { c.http = h } }

// New returns a client for the server at baseURL.
func New(baseURL string, opts ...Option) *Client {
	c := &Client{
		base: strings.TrimRight(baseURL, "/"),
		http: &http.Client{Timeout: 8 * time.Second},
	}
	for _, o := range opts {
		o(c)
	}
	return c
}

func (c *Client) call(ctx context.Context, method, path string, body, out any) error {
	var rdr io.Reader
	if body != nil {
		raw, err := json.Marshal(body)
		if err != nil {
			return &Error{Message: "encode request: " + err.Error()}
		}
		rdr = bytes.NewReader(raw)
	}
	req, err := http.NewRequestWithContext(ctx, method, c.base+path, rdr)
	if err != nil {
		return &Error{Message: err.Error()}
	}
	req.Header.Set("Content-Type", "application/json")
	if c.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return &Error{Message: fmt.Sprintf("unreachable at %s (%v)", c.base, err)}
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		msg := string(raw)
		var e struct {
			Error string `json:"error"`
		}
		if json.Unmarshal(raw, &e) == nil && e.Error != "" {
			msg = e.Error
		}
		return &Error{Status: resp.StatusCode, Message: msg}
	}
	if out != nil && len(raw) > 0 {
		if err := json.Unmarshal(raw, out); err != nil {
			return &Error{Message: "decode response: " + err.Error()}
		}
	}
	return nil
}

// Health reports whether the server is reachable and healthy; it never errors.
func (c *Client) Health(ctx context.Context) bool {
	var r struct {
		Status string `json:"status"`
	}
	return c.call(ctx, http.MethodGet, "/health", nil, &r) == nil && r.Status == "ok"
}

// Ready reports orchestrator readiness (a fresh revocation snapshot is held).
func (c *Client) Ready(ctx context.Context) bool {
	var r struct {
		Ready bool `json:"ready"`
	}
	return c.call(ctx, http.MethodGet, "/readyz", nil, &r) == nil && r.Ready
}

// Version returns server + trust-domain metadata.
func (c *Client) Version(ctx context.Context) (map[string]any, error) {
	var r map[string]any
	return r, c.call(ctx, http.MethodGet, "/version", nil, &r)
}

// Issue mints a scoped, expiring, revocable delegation.
func (c *Client) Issue(ctx context.Context, p IssueParams) (*Grant, error) {
	ttl := p.TTL
	if ttl <= 0 {
		ttl = time.Hour
	}
	var g Grant
	err := c.call(ctx, http.MethodPost, "/issue", map[string]any{
		"principal":  p.Principal,
		"delegate":   p.Delegate,
		"scope":      p.Scope,
		"ttlSeconds": int(ttl.Seconds()),
	}, &g)
	if err != nil {
		return nil, err
	}
	return &g, nil
}

// Verify runs the five-check verification against a presented record.
func (c *Client) Verify(ctx context.Context, record string) (*VerifyResult, error) {
	var v VerifyResult
	if err := c.call(ctx, http.MethodPost, "/verify", map[string]any{"record": record}, &v); err != nil {
		return nil, err
	}
	return &v, nil
}

// Revoke revokes a delegation instance.
func (c *Client) Revoke(ctx context.Context, instance string) error {
	return c.call(ctx, http.MethodPost, "/revoke", map[string]any{"instance": instance}, nil)
}

// Delegations lists the known delegations.
func (c *Client) Delegations(ctx context.Context) ([]map[string]any, error) {
	var r struct {
		Delegations []map[string]any `json:"delegations"`
	}
	return r.Delegations, c.call(ctx, http.MethodGet, "/delegations", nil, &r)
}

// Audit returns up to limit recent audit events.
func (c *Client) Audit(ctx context.Context, limit int) ([]map[string]any, error) {
	if limit <= 0 {
		limit = 50
	}
	var r struct {
		Events []map[string]any `json:"events"`
	}
	return r.Events, c.call(ctx, http.MethodGet, "/audit?limit="+strconv.Itoa(limit), nil, &r)
}

// Graph returns the trust graph (edges).
func (c *Client) Graph(ctx context.Context) (map[string]any, error) {
	var r map[string]any
	return r, c.call(ctx, http.MethodGet, "/graph", nil, &r)
}

// Stats returns the JSON metrics view.
func (c *Client) Stats(ctx context.Context) (map[string]any, error) {
	var r map[string]any
	return r, c.call(ctx, http.MethodGet, "/stats", nil, &r)
}
