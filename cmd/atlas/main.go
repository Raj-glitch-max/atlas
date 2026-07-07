// Command atlas is the Atlas CLI: a human-facing client over the Atlas Server
// (cmd/atlas-server). It issues, verifies, and revokes delegations and inspects
// trust state — the same operations agents get over MCP, on the same shared
// state.
//
//	atlas [--api URL] <command> [args]
//
// Commands: issue · verify · revoke · delegations · graph · audit · doctor · version
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

const cliVersion = "0.1.0-dev"

func main() { os.Exit(run(os.Args[1:], os.Stdout)) }

func run(args []string, out io.Writer) int {
	fs := flag.NewFlagSet("atlas", flag.ContinueOnError)
	fs.SetOutput(out)
	api := fs.String("api", envOr("ATLAS_API", "http://127.0.0.1:8087"), "Atlas Server base URL")
	apiKey := fs.String("api-key", envOr("ATLAS_API_KEY", ""), "bearer token for mutating commands (or $ATLAS_API_KEY)")
	fs.Usage = func() { usage(out) }
	if err := fs.Parse(args); err != nil {
		return exitUsage
	}
	rest := fs.Args()
	if len(rest) == 0 {
		usage(out)
		return exitUsage
	}
	c := newClient(*api, *apiKey)
	cmd, cmdArgs := rest[0], rest[1:]
	switch cmd {
	case "delegate", "issue": // delegate is the canonical verb; issue is the alias
		return cmdIssue(out, c, cmdArgs)
	case "verify":
		return cmdVerify(out, c, cmdArgs)
	case "bundle":
		return cmdBundle(out, c, cmdArgs)
	case "inspect":
		return cmdInspect(out, c, cmdArgs)
	case "revoke":
		return cmdRevoke(out, c, cmdArgs)
	case "delegations", "ls":
		return cmdDelegations(out, c, cmdArgs)
	case "graph":
		return cmdGraph(out, c, cmdArgs)
	case "audit":
		return cmdAudit(out, c, cmdArgs)
	case "doctor":
		return cmdDoctor(out, c, cmdArgs)
	case "version":
		return cmdVersion(out, c, cmdArgs)
	case "help", "-h", "--help":
		usage(out)
		return exitOK
	default:
		fmt.Fprintf(out, "unknown command %q\n\n", cmd)
		usage(out)
		return exitUsage
	}
}

func usage(out io.Writer) {
	fmt.Fprint(out, `atlas — let your agents delegate work safely

usage:
  atlas [--api URL] <command> [args]

commands:
  delegate  --principal ID --delegate ID --scope a,b [--ttl 3600] [-q]
            grant a scoped, expiring, revocable capability (alias: issue)
  verify    <record | ->            verify a capability (exit 3 if not accepted)
            --offline --bundle F [--max-staleness 5m]   verify with NO server, from a trust bundle
            --require-scope ACTION                      authz gate: pass only if valid AND grants ACTION
  bundle    [-o trust-bundle.json]  export the trust bundle for offline verification
  inspect   <record | ->            decode a capability's claims locally (not verified)
  revoke    <instance>              revoke a capability by instance id
  delegations                       list issued capabilities
  graph                             show the delegation graph
  audit     [--limit N]             recent delegate/verify/revoke events
  doctor                            check server reachability + config
  version                           client and server version

env:
  ATLAS_API       server base URL (default http://127.0.0.1:8087)
  ATLAS_API_KEY   bearer token for mutating commands

the loop:
  atlas delegate --principal spiffe://a.test/api --delegate spiffe://b.test/agent \
                 --scope read:orders -q | atlas verify -
  atlas bundle && atlas verify --offline --bundle trust-bundle.json <record>
`)
}

func envOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
