package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"
	"time"
)

// Exit codes: 0 ok · 1 error · 2 usage · 3 verify produced a non-accept verdict
// (so `atlas verify $rec && …` gates on Accept).
const (
	exitOK      = 0
	exitError   = 1
	exitUsage   = 2
	exitNonPass = 3
)

func cmdIssue(out io.Writer, c *client, args []string) int {
	fs := flag.NewFlagSet("issue", flag.ContinueOnError)
	fs.SetOutput(out)
	principal := fs.String("principal", "", "SPIFFE ID delegating authority (required)")
	delegate := fs.String("delegate", "", "SPIFFE ID receiving the delegation (required)")
	scope := fs.String("scope", "", "comma-separated capabilities, e.g. read:orders,write:audit (required)")
	ttl := fs.Int("ttl", 3600, "lifetime in seconds")
	quiet := fs.Bool("q", false, "print only the record (for piping to `atlas verify -`)")
	if err := fs.Parse(args); err != nil {
		return exitUsage
	}
	if *principal == "" || *delegate == "" || *scope == "" {
		fmt.Fprintln(out, "issue: --principal, --delegate and --scope are required")
		return exitUsage
	}
	r, err := c.issue(map[string]any{
		"principal": *principal, "delegate": *delegate,
		"scope": splitScope(*scope), "ttlSeconds": *ttl,
	})
	if err != nil {
		fmt.Fprintln(out, "issue failed:", err)
		return exitError
	}
	if *quiet {
		fmt.Fprintln(out, r.Record)
		return exitOK
	}
	fmt.Fprintf(out, "issued  %s → %s\n", r.Principal, r.Delegate)
	fmt.Fprintf(out, "scope   %s\n", strings.Join(r.Scope, ", "))
	fmt.Fprintf(out, "expires %s\n", r.ExpiresAt)
	fmt.Fprintf(out, "instance %s\n", r.Instance)
	fmt.Fprintf(out, "\n%s\n", r.Record)
	return exitOK
}

func cmdVerify(out io.Writer, c *client, args []string) int {
	fs := flag.NewFlagSet("verify", flag.ContinueOnError)
	fs.SetOutput(out)
	offline := fs.Bool("offline", false, "verify locally from a trust bundle — no server, no network")
	bundle := fs.String("bundle", "trust-bundle.json", "trust bundle file (with --offline)")
	maxStale := fs.Duration("max-staleness", 5*time.Minute, "your freshness budget R for the bundled revocation snapshot (with --offline)")
	requireScope := fs.String("require-scope", "", "authorize a specific action: pass only if the capability is valid AND grants this scope (turns verify into an authz gate for a tool adapter)")
	if err := fs.Parse(args); err != nil {
		return exitUsage
	}
	rec, err := readRecord(fs.Args())
	if err != nil {
		fmt.Fprintln(out, "verify:", err)
		return exitUsage
	}
	if *offline {
		return verifyOffline(out, *bundle, rec, *maxStale, *requireScope)
	}
	r, err := c.verify(rec)
	if err != nil {
		fmt.Fprintln(out, "verify failed:", err)
		return exitError
	}
	label := strings.ToUpper(r.Decision)
	fmt.Fprintf(out, "%s  (%dµs)", label, r.LatencyMicros)
	if len(r.Causes) > 0 {
		fmt.Fprintf(out, "  %s", strings.Join(r.Causes, ", "))
	}
	fmt.Fprintln(out)
	tw := tabwriter.NewWriter(out, 0, 2, 2, ' ', 0)
	for _, e := range r.Trace {
		cause := e.Cause
		if cause == "None" || cause == "" {
			cause = "—"
		}
		fmt.Fprintf(tw, "  %s\t%s\t%s\n", e.Check, e.Outcome, cause)
	}
	tw.Flush()
	return scopeGate(out, rec, *requireScope, r.Accept)
}

// scopeGate applies the optional --require-scope authorization check on top of
// an authentication verdict. Authorization = the capability is ACCEPTed AND its
// scope grants the requested action. This is the exact decision a tool adapter
// makes before acting with its own credential.
func scopeGate(out io.Writer, rec, requireScope string, accepted bool) int {
	if requireScope == "" {
		if !accepted {
			return exitNonPass
		}
		return exitOK
	}
	granted := accepted && recordGrantsScope(rec, requireScope)
	if granted {
		fmt.Fprintf(out, "authorized: YES  (grants %q)\n", requireScope)
		return exitOK
	}
	reason := "capability not accepted"
	if accepted {
		reason = fmt.Sprintf("valid, but does not grant %q", requireScope)
	}
	fmt.Fprintf(out, "authorized: NO   (%s)\n", reason)
	return exitNonPass
}

// recordGrantsScope decodes the record's scope claim (offline, no verify) and
// reports whether it contains the required scope. Only meaningful once the
// record has been authenticated — scopeGate calls it only after ACCEPT.
func recordGrantsScope(rec, required string) bool {
	parts := strings.Split(strings.TrimSpace(rec), ".")
	if len(parts) != 3 {
		return false
	}
	payload, err := b64urlJSON(parts[1])
	if err != nil {
		return false
	}
	var claims struct {
		Scope []string `json:"scope"`
	}
	if err := json.Unmarshal(payload, &claims); err != nil {
		return false
	}
	for _, s := range claims.Scope {
		if s == required {
			return true
		}
	}
	return false
}

// cmdInspect decodes a record's header and claims locally — no server, offline.
// It does NOT verify the signature; that is what `atlas verify` does.
func cmdInspect(out io.Writer, _ *client, args []string) int {
	rec, err := readRecord(args)
	if err != nil {
		fmt.Fprintln(out, "inspect:", err)
		return exitUsage
	}
	parts := strings.Split(strings.TrimSpace(rec), ".")
	if len(parts) != 3 {
		fmt.Fprintln(out, "inspect: not a compact JWS (expected three dot-separated segments)")
		return exitError
	}
	hdr, err := b64urlJSON(parts[0])
	if err != nil {
		fmt.Fprintln(out, "inspect: bad header:", err)
		return exitError
	}
	payload, err := b64urlJSON(parts[1])
	if err != nil {
		fmt.Fprintln(out, "inspect: bad payload:", err)
		return exitError
	}
	fmt.Fprintln(out, "header")
	fmt.Fprintln(out, indent(hdr))
	fmt.Fprintln(out, "claims")
	fmt.Fprintln(out, indent(payload))
	fmt.Fprintf(out, "signature  %d bytes (base64url) — decoded, NOT verified; use `atlas verify` to check\n", len(parts[2]))
	return exitOK
}

func b64urlJSON(seg string) ([]byte, error) {
	raw, err := base64.RawURLEncoding.DecodeString(seg)
	if err != nil {
		return nil, err
	}
	if !json.Valid(raw) {
		return nil, fmt.Errorf("segment is not JSON")
	}
	return raw, nil
}

func indent(raw []byte) string {
	var buf bytes.Buffer
	if err := json.Indent(&buf, raw, "  ", "  "); err != nil {
		return "  " + string(raw)
	}
	return "  " + buf.String()
}

func cmdRevoke(out io.Writer, c *client, args []string) int {
	if len(args) != 1 {
		fmt.Fprintln(out, "usage: atlas revoke <instance>")
		return exitUsage
	}
	if err := c.revoke(args[0]); err != nil {
		fmt.Fprintln(out, "revoke failed:", err)
		return exitError
	}
	fmt.Fprintf(out, "revoked %s\n", args[0])
	return exitOK
}

func cmdDelegations(out io.Writer, c *client, _ []string) int {
	ds, err := c.delegations()
	if err != nil {
		fmt.Fprintln(out, "delegations failed:", err)
		return exitError
	}
	if len(ds) == 0 {
		fmt.Fprintln(out, "no delegations issued yet")
		return exitOK
	}
	tw := tabwriter.NewWriter(out, 0, 2, 2, ' ', 0)
	fmt.Fprintln(tw, "INSTANCE\tPRINCIPAL\tDELEGATE\tSCOPE\tSTATUS")
	for _, d := range ds {
		status := "active"
		if d.Revoked {
			status = "revoked"
		}
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\n", short(d.Instance), leaf(d.Principal), leaf(d.Delegate), strings.Join(d.Scope, ","), status)
	}
	tw.Flush()
	return exitOK
}

func cmdGraph(out io.Writer, c *client, _ []string) int {
	g, err := c.graph()
	if err != nil {
		fmt.Fprintln(out, "graph failed:", err)
		return exitError
	}
	if len(g.Edges) == 0 {
		fmt.Fprintln(out, "empty graph")
		return exitOK
	}
	for _, e := range g.Edges {
		arrow := "──►"
		if e.Revoked {
			arrow = "──✗"
		}
		fmt.Fprintf(out, "%s %s %s  [%s]\n", leaf(e.From), arrow, leaf(e.To), strings.Join(e.Scope, ","))
	}
	fmt.Fprintf(out, "\n%d identities · %d delegations\n", len(g.Nodes), len(g.Edges))
	return exitOK
}

func cmdAudit(out io.Writer, c *client, args []string) int {
	fs := flag.NewFlagSet("audit", flag.ContinueOnError)
	fs.SetOutput(out)
	limit := fs.Int("limit", 50, "max events")
	if err := fs.Parse(args); err != nil {
		return exitUsage
	}
	a, err := c.audit(*limit)
	if err != nil {
		fmt.Fprintln(out, "audit failed:", err)
		return exitError
	}
	tw := tabwriter.NewWriter(out, 0, 2, 2, ' ', 0)
	for _, e := range a.Events {
		detail := e.Decision
		if detail == "" {
			detail = e.Detail
		}
		subject := leaf(e.Delegate)
		if subject == "" {
			subject = short(e.Instance)
		}
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n", e.Time, e.Type, subject, detail)
	}
	tw.Flush()
	return exitOK
}

func cmdDoctor(out io.Writer, c *client, _ []string) int {
	fmt.Fprintf(out, "atlas doctor\n  api: %s\n", c.base)
	if err := c.health(); err != nil {
		fmt.Fprintf(out, "  server: UNREACHABLE (%s)\n", err)
		fmt.Fprintln(out, "\n  → start it with:  go run ./cmd/atlas-server")
		return exitError
	}
	fmt.Fprintln(out, "  server: reachable ✓")
	v, err := c.version()
	if err == nil {
		fmt.Fprintf(out, "  trust domain: %v\n  algorithm: %v\n  revocation R: %v\n", v["trustDomain"], v["algorithm"], v["revocationR"])
	}
	return exitOK
}

func cmdVersion(out io.Writer, c *client, _ []string) int {
	fmt.Fprintf(out, "atlas %s\n", cliVersion)
	if v, err := c.version(); err == nil {
		fmt.Fprintf(out, "server %v (%v)\n", v["version"], v["trustDomain"])
	}
	return exitOK
}

// ---- helpers ----

func splitScope(s string) []string {
	parts := strings.Split(s, ",")
	out := parts[:0]
	for _, p := range parts {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, p)
		}
	}
	return out
}

// readRecord returns the record from the first arg, or stdin when the arg is
// "-" (so `atlas issue -q … | atlas verify -` works).
func readRecord(args []string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("usage: atlas verify <record|->")
	}
	if args[0] == "-" {
		b, err := io.ReadAll(os.Stdin)
		if err != nil {
			return "", err
		}
		return strings.TrimSpace(string(b)), nil
	}
	return args[0], nil
}

func short(s string) string {
	if len(s) > 12 {
		return s[:12] + "…"
	}
	return s
}

// leaf trims a SPIFFE ID to its trailing path segment for compact display.
func leaf(id string) string {
	if id == "" {
		return ""
	}
	if i := strings.LastIndex(id, "/"); i >= 0 && i < len(id)-1 {
		return id[i+1:]
	}
	return id
}

// checkpoint: chore(issuance): optimize truststore backend

// checkpoint: chore(sdk): simplify conformance validation
