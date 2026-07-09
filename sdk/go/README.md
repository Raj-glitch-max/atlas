# Atlas Go SDK

A zero-dependency Go client for the [Atlas Server](../../cmd/atlas-server). It
uses only the standard library, so importing it adds no third-party
dependencies. The API mirrors the [Python](../python) and
[TypeScript](../typescript) SDKs.

## Install

```bash
go get github.com/Raj-glitch-max/atlas/sdk/go
```

## Use

```go
import atlas "github.com/Raj-glitch-max/atlas/sdk/go"

c := atlas.New("http://127.0.0.1:8087") // atlas.WithAPIKey("…") if writes are guarded
ctx := context.Background()

grant, err := c.Issue(ctx, atlas.IssueParams{
    Principal: "spiffe://domain-a.test/workload/payments-api",
    Delegate:  "spiffe://domain-b.test/agent/booking-worker",
    Scope:     []string{"read:orders", "write:audit"},
    TTL:       time.Hour,
})

res, err := c.Verify(ctx, grant.Record) // res.Decision == atlas.Accept
err = c.Revoke(ctx, grant.Instance)     // now Verify -> atlas.Reject
```

Errors surface as `*atlas.Error` (with a `.Status` for HTTP errors).
`Health` and `Ready` return a bool and never error.

## Verify it

```bash
go test ./sdk/go/...          # httptest-based, no server needed
go run ./cmd/atlas-server &   # then, against the real engine:
go run ./sdk/go/example
```

For an offline authorization gate built on the same engine, see
[`examples/atlas-gate`](../../examples/atlas-gate).
