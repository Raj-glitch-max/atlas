# atlas-mcp — Atlas for agents (MCP server)

Exposes the Atlas delegation engine as **agent tools** over the Model Context
Protocol (stdio JSON-RPC 2.0). An MCP host (Claude, or any MCP client) can then
issue, verify, revoke, and inspect delegations — the "software for agents"
surface.

It is a thin adapter over the Atlas Server (`cmd/atlas-server`), so agents and
the human UI share one trust state: a delegation an agent issues over MCP shows
up in the site's graph, and a revocation there is observed by the agent's next
`atlas_verify`.

## Tools

| Tool | Purpose |
|---|---|
| `atlas_issue` | Issue an attenuable delegation (principal → delegate, scoped). Returns the signed record + instance id. |
| `atlas_verify` | Verify a presented record offline. Returns decision, causes, the five-check trace, and latency. |
| `atlas_revoke` | Revoke by instance id (signed snapshot; next verify → `RevokedObservable`). |
| `atlas_delegations` | List issued delegations. |
| `atlas_graph` | Delegation graph (principal → delegate edges). |
| `atlas_audit` | Recent issue/verify/revoke events. |

## Run

```bash
# 1. the backend must be up (holds the shared trust state)
go run ./cmd/atlas-server            # → 127.0.0.1:8087

# 2. build the MCP server
go build -o atlas-mcp ./cmd/atlas-mcp
```

The server base URL comes from `-api` or `$ATLAS_API` (default
`http://127.0.0.1:8087`). It speaks MCP on stdio — an MCP host launches it.

## Register with Claude Code

```bash
claude mcp add atlas -- /absolute/path/to/atlas-mcp -api http://127.0.0.1:8087
```

Or add to the MCP client config (Claude Desktop `claude_desktop_config.json`):

```json
{
  "mcpServers": {
    "atlas": {
      "command": "/absolute/path/to/atlas-mcp",
      "args": ["-api", "http://127.0.0.1:8087"]
    }
  }
}
```

Then ask the agent to, e.g., *"issue a read:orders delegation from
spiffe://domain-a.test/workload/payments-api to
spiffe://domain-b.test/agent/booking-worker, verify it, then revoke it and
verify again."*

## Design notes

- **No third-party dependency**: the MCP protocol (initialize, tools/list,
  tools/call, ping) is hand-rolled on stdlib — keeps the module supply-chain
  clean and governance-tight (import-lint: 0 violations, no internal imports).
- **Honest errors**: an unreachable server or a refused issuance surfaces as an
  MCP `isError: true` result, not a fabricated success.

<!-- checkpoint: chore(verify): simplify revocation status lookup -->

<!-- checkpoint: chore(security): clean lab environment topology -->
