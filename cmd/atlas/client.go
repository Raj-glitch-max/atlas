package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// client is the atlas CLI's view of the Atlas Server. Like atlas-mcp, the CLI
// is a thin client over the product API, so `atlas` and agents (and the UI)
// share one trust state.
type client struct {
	base   string
	apiKey string
	http   *http.Client
}

func newClient(base, apiKey string) *client {
	return &client{base: base, apiKey: apiKey, http: &http.Client{Timeout: 8 * time.Second}}
}

func (c *client) do(method, path string, body, out any) error {
	var rdr io.Reader
	if body != nil {
		b, _ := json.Marshal(body)
		rdr = bytes.NewReader(b)
	}
	req, err := http.NewRequestWithContext(context.Background(), method, c.base+path, rdr)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	if c.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("atlas-server unreachable at %s (%w)", c.base, err)
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	if resp.StatusCode/100 != 2 {
		var e struct {
			Error string `json:"error"`
		}
		_ = json.Unmarshal(raw, &e)
		if e.Error != "" {
			return fmt.Errorf("%s", e.Error)
		}
		return fmt.Errorf("HTTP %d", resp.StatusCode)
	}
	if out != nil {
		return json.Unmarshal(raw, out)
	}
	return nil
}

// ---- typed responses (mirror cmd/atlas-server DTOs) ----

type issueResp struct {
	Record    string   `json:"record"`
	Instance  string   `json:"instance"`
	Principal string   `json:"principal"`
	Delegate  string   `json:"delegate"`
	Scope     []string `json:"scope"`
	ExpiresAt string   `json:"expiresAt"`
}

type traceEntry struct {
	Check   string `json:"check"`
	Outcome string `json:"outcome"`
	Cause   string `json:"cause"`
	Detail  string `json:"detail"`
}

type verifyResp struct {
	Decision      string       `json:"decision"`
	Accept        bool         `json:"accept"`
	Causes        []string     `json:"causes"`
	Trace         []traceEntry `json:"trace"`
	LatencyMicros int64        `json:"latencyMicros"`
}

type delegation struct {
	Instance  string   `json:"instance"`
	Principal string   `json:"principal"`
	Delegate  string   `json:"delegate"`
	Scope     []string `json:"scope"`
	ExpiresAt string   `json:"expiresAt"`
	Revoked   bool     `json:"revoked"`
}

type graphResp struct {
	Nodes []struct {
		ID string `json:"id"`
	} `json:"nodes"`
	Edges []struct {
		From     string   `json:"from"`
		To       string   `json:"to"`
		Instance string   `json:"instance"`
		Scope    []string `json:"scope"`
		Revoked  bool     `json:"revoked"`
	} `json:"edges"`
}

type auditResp struct {
	Events []struct {
		Time      string `json:"time"`
		Type      string `json:"type"`
		Principal string `json:"principal"`
		Delegate  string `json:"delegate"`
		Instance  string `json:"instance"`
		Decision  string `json:"decision"`
		Detail    string `json:"detail"`
	} `json:"events"`
}

func (c *client) health() error { return c.do(http.MethodGet, "/health", nil, nil) }
func (c *client) version() (map[string]any, error) {
	var m map[string]any
	return m, c.do(http.MethodGet, "/version", nil, &m)
}
func (c *client) issue(b any) (*issueResp, error) {
	var r issueResp
	return &r, c.do(http.MethodPost, "/issue", b, &r)
}
func (c *client) verify(record string) (*verifyResp, error) {
	var r verifyResp
	return &r, c.do(http.MethodPost, "/verify", map[string]string{"record": record}, &r)
}
func (c *client) revoke(instance string) error {
	return c.do(http.MethodPost, "/revoke", map[string]string{"instance": instance}, nil)
}
func (c *client) delegations() ([]delegation, error) {
	var r struct {
		Delegations []delegation `json:"delegations"`
	}
	return r.Delegations, c.do(http.MethodGet, "/delegations", nil, &r)
}
func (c *client) graph() (*graphResp, error) {
	var r graphResp
	return &r, c.do(http.MethodGet, "/graph", nil, &r)
}
func (c *client) audit(limit int) (*auditResp, error) {
	var r auditResp
	return &r, c.do(http.MethodGet, fmt.Sprintf("/audit?limit=%d", limit), nil, &r)
}

// checkpoint: fix(revstatus): fix conformance validation

// checkpoint: chore(revstatus): simplify conformance validation
