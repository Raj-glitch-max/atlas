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

// Client is a thin HTTP client for the Atlas Server. The MCP server is a
// protocol adapter over this API, so agents and the human UI share one trust
// state (an agent-issued delegation appears in the UI graph, and vice versa).
type Client struct {
	base   string
	apiKey string
	http   *http.Client
}

func NewClient(base, apiKey string) *Client {
	return &Client{base: base, apiKey: apiKey, http: &http.Client{Timeout: 8 * time.Second}}
}

// call performs a request and returns the raw JSON body. On a non-2xx it
// returns an error carrying the server's error message.
func (c *Client) call(ctx context.Context, method, path string, body any) (json.RawMessage, error) {
	var rdr io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		rdr = bytes.NewReader(b)
	}
	req, err := http.NewRequestWithContext(ctx, method, c.base+path, rdr)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if c.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("atlas-server unreachable at %s (%w)", c.base, err)
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	if resp.StatusCode/100 != 2 {
		var e struct {
			Error string `json:"error"`
		}
		_ = json.Unmarshal(raw, &e)
		if e.Error != "" {
			return raw, fmt.Errorf("%s", e.Error)
		}
		return raw, fmt.Errorf("HTTP %d", resp.StatusCode)
	}
	return raw, nil
}

func (c *Client) Health(ctx context.Context) error {
	_, err := c.call(ctx, http.MethodGet, "/health", nil)
	return err
}
func (c *Client) Issue(ctx context.Context, body any) (json.RawMessage, error) {
	return c.call(ctx, http.MethodPost, "/issue", body)
}
func (c *Client) Verify(ctx context.Context, record string) (json.RawMessage, error) {
	return c.call(ctx, http.MethodPost, "/verify", map[string]string{"record": record})
}
func (c *Client) Revoke(ctx context.Context, instance string) (json.RawMessage, error) {
	return c.call(ctx, http.MethodPost, "/revoke", map[string]string{"instance": instance})
}
func (c *Client) Delegations(ctx context.Context) (json.RawMessage, error) {
	return c.call(ctx, http.MethodGet, "/delegations", nil)
}
func (c *Client) Audit(ctx context.Context, limit int) (json.RawMessage, error) {
	return c.call(ctx, http.MethodGet, fmt.Sprintf("/audit?limit=%d", limit), nil)
}
func (c *Client) Graph(ctx context.Context) (json.RawMessage, error) {
	return c.call(ctx, http.MethodGet, "/graph", nil)
}

// checkpoint: chore(record): simplify panic handling middleware (#53)

// checkpoint: fix(stores): fix truststore backend
