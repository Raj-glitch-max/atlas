package main

import (
	"context"
	"encoding/json"
	"fmt"
)

// atlasTools defines the MCP tool surface. Descriptions are written for an
// agent deciding when to call each one.
func atlasTools() []tool {
	return []tool{
		{
			Name:        "atlas_issue",
			Description: "Issue an offline-verifiable delegation from a principal to a delegate, scoped to a set of capabilities. The delegate may only be granted a subset of what the principal holds (attenuation); over-scope requests are refused. Returns the signed record (a compact JWS the delegate presents) and its instance id (used to revoke it).",
			Schema: json.RawMessage(`{"type":"object","properties":{
				"principal":{"type":"string","description":"SPIFFE ID delegating authority, e.g. spiffe://domain-a.test/workload/payments-api"},
				"delegate":{"type":"string","description":"SPIFFE ID receiving the delegation, e.g. spiffe://domain-b.test/agent/booking-worker"},
				"scope":{"type":"array","items":{"type":"string"},"description":"capabilities to grant, e.g. [\"read:orders\",\"write:audit\"]"},
				"ttlSeconds":{"type":"integer","description":"lifetime in seconds (default 3600)"}
			},"required":["principal","delegate","scope"]}`),
			Handler: func(ctx context.Context, c *Client, args json.RawMessage) (string, bool) {
				var a struct {
					Principal  string   `json:"principal"`
					Delegate   string   `json:"delegate"`
					Scope      []string `json:"scope"`
					TTLSeconds int      `json:"ttlSeconds"`
				}
				if err := json.Unmarshal(args, &a); err != nil {
					return errText("issue", err)
				}
				raw, err := c.Issue(ctx, a)
				if err != nil {
					return errText("issue", err)
				}
				var r struct {
					Instance string `json:"instance"`
				}
				_ = json.Unmarshal(raw, &r)
				return fmt.Sprintf("Issued delegation %s → %s.\ninstance: %s (use with atlas_revoke)\n\n%s",
					a.Principal, a.Delegate, r.Instance, pretty(raw)), false
			},
		},
		{
			Name:        "atlas_verify",
			Description: "Verify a presented delegation record offline (no live authority call). Returns the decision (accept | reject | inconclusive), the causes, the five-check trace (identity binding, signature, freshness, scope, revocation), and the measured latency. Use this to decide whether to honor a capability an agent presents.",
			Schema: json.RawMessage(`{"type":"object","properties":{
				"record":{"type":"string","description":"the compact JWS delegation record to verify (from atlas_issue)"}
			},"required":["record"]}`),
			Handler: func(ctx context.Context, c *Client, args json.RawMessage) (string, bool) {
				var a struct {
					Record string `json:"record"`
				}
				if err := json.Unmarshal(args, &a); err != nil {
					return errText("verify", err)
				}
				raw, err := c.Verify(ctx, a.Record)
				if err != nil {
					return errText("verify", err)
				}
				var r struct {
					Decision      string   `json:"decision"`
					Causes        []string `json:"causes"`
					LatencyMicros int64    `json:"latencyMicros"`
				}
				_ = json.Unmarshal(raw, &r)
				head := fmt.Sprintf("Verdict: %s (%dµs)", r.Decision, r.LatencyMicros)
				if len(r.Causes) > 0 {
					head += " — " + fmt.Sprint(r.Causes)
				}
				return head + "\n\n" + pretty(raw), false
			},
		},
		{
			Name:        "atlas_revoke",
			Description: "Revoke a delegation by its instance id. The revocation is published as a signed snapshot; subsequent verifications of that record return reject (RevokedObservable) once the snapshot propagates, bounded by the freshness window R.",
			Schema: json.RawMessage(`{"type":"object","properties":{
				"instance":{"type":"string","description":"the instance id returned by atlas_issue"}
			},"required":["instance"]}`),
			Handler: func(ctx context.Context, c *Client, args json.RawMessage) (string, bool) {
				var a struct {
					Instance string `json:"instance"`
				}
				if err := json.Unmarshal(args, &a); err != nil {
					return errText("revoke", err)
				}
				raw, err := c.Revoke(ctx, a.Instance)
				if err != nil {
					return errText("revoke", err)
				}
				return "Revoked " + a.Instance + ".\n\n" + pretty(raw), false
			},
		},
		{
			Name:        "atlas_delegations",
			Description: "List the delegations issued so far (metadata: principal, delegate, scope, expiry, revoked). Use to inspect current trust relationships.",
			Schema:      json.RawMessage(`{"type":"object","properties":{}}`),
			Handler: func(ctx context.Context, c *Client, _ json.RawMessage) (string, bool) {
				raw, err := c.Delegations(ctx)
				if err != nil {
					return errText("delegations", err)
				}
				return pretty(raw), false
			},
		},
		{
			Name:        "atlas_graph",
			Description: "Return the delegation graph (principal → delegate edges with scope and revocation state) for the current trust state.",
			Schema:      json.RawMessage(`{"type":"object","properties":{}}`),
			Handler: func(ctx context.Context, c *Client, _ json.RawMessage) (string, bool) {
				raw, err := c.Graph(ctx)
				if err != nil {
					return errText("graph", err)
				}
				return pretty(raw), false
			},
		},
		{
			Name:        "atlas_audit",
			Description: "Return the recent audit log (issue / verify / revoke events), newest first. Use to inspect what has happened.",
			Schema: json.RawMessage(`{"type":"object","properties":{
				"limit":{"type":"integer","description":"max events to return (default 50)"}
			}}`),
			Handler: func(ctx context.Context, c *Client, args json.RawMessage) (string, bool) {
				var a struct {
					Limit int `json:"limit"`
				}
				_ = json.Unmarshal(args, &a)
				if a.Limit <= 0 {
					a.Limit = 50
				}
				raw, err := c.Audit(ctx, a.Limit)
				if err != nil {
					return errText("audit", err)
				}
				return pretty(raw), false
			},
		},
	}
}
