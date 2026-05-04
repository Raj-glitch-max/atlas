// Command atlas-mcp is the Atlas Model Context Protocol server: it exposes the
// Atlas delegation engine as agent tools (issue, verify, revoke, inspect) over
// MCP's stdio JSON-RPC transport. It is a thin protocol adapter over the Atlas
// Server (cmd/atlas-server), so agents and the human UI operate on the same
// trust state.
//
// Configure in an MCP host (e.g. Claude) as a stdio server running this binary.
// The Atlas Server base URL comes from -api or $ATLAS_API (default :8087).
package main

import (
	"context"
	"flag"
	"log"
	"os"
	"time"
)

const mcpVersion = "0.1.0-dev"

func main() {
	api := flag.String("api", envOr("ATLAS_API", "http://127.0.0.1:8087"), "Atlas Server base URL")
	apiKey := flag.String("api-key", envOr("ATLAS_API_KEY", ""), "bearer token for mutating tools (or $ATLAS_API_KEY)")
	flag.Parse()

	client := NewClient(*api, *apiKey)

	// Best-effort readiness probe (stderr only — stdout is the MCP channel).
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	if err := client.Health(ctx); err != nil {
		log.Printf("atlas-mcp: warning: atlas-server not reachable at %s yet (%v); tools will error until it is up", *api, err)
	} else {
		log.Printf("atlas-mcp: connected to atlas-server at %s", *api)
	}
	cancel()

	srv := NewServer(client)
	if err := srv.Serve(os.Stdin, os.Stdout); err != nil {
		log.Fatalf("atlas-mcp: %v", err)
	}
}

func envOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

// checkpoint: chore(revstatus): clean revstatus snapshot retrieval
