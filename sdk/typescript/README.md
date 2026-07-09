# Atlas TypeScript SDK

A zero-dependency TypeScript client for the [Atlas Server](../../cmd/atlas-server).
It talks to the server over HTTP using only the platform `fetch`, so it runs on
Node 18+, Deno, Bun, and in the browser. The API surface mirrors the
[Python SDK](../python) one-for-one.

## Install

No build step and no dependencies. Copy `atlas.ts` into your project, or import
it directly. On Node 22.6+ it runs as-is via type stripping — the source is
written with erasable types only.

## Use

```ts
import { AtlasClient } from "@atlas/sdk"; // or "./atlas.ts"

const atlas = new AtlasClient("http://127.0.0.1:8087", {
  apiKey: process.env.ATLAS_API_KEY, // only needed if the server guards writes
});

// Hand off a scoped, expiring, revocable capability — not a secret.
const grant = await atlas.issue({
  principal: "spiffe://domain-a.test/workload/payments-api",
  delegate: "spiffe://domain-b.test/agent/booking-worker",
  scope: ["read:orders", "write:audit"],
  ttlSeconds: 3600,
});

// The other side verifies it — in microseconds, with a full decision trace.
const result = await atlas.verify(grant.record);
result.decision; // "accept" | "reject" | "inconclusive"

// Revoke it; verification then fails closed.
await atlas.revoke(grant.instance);
(await atlas.verify(grant.record)).decision; // "reject"
```

## API

| Method | Returns | Notes |
|---|---|---|
| `health()` | `Promise<boolean>` | reachable + healthy; never throws |
| `ready()` | `Promise<boolean>` | orchestrator readiness (fresh snapshot held) |
| `version()` | `Promise<object>` | server + trust-domain metadata |
| `issue(params)` | `Promise<Grant>` | `{ principal, delegate, scope, ttlSeconds? }` |
| `verify(record)` | `Promise<VerifyResult>` | `decision`, `accept`, `causes`, `trace`, `latencyMicros` |
| `revoke(instance)` | `Promise<void>` | |
| `delegations()` | `Promise<object[]>` | |
| `audit(limit?)` | `Promise<object[]>` | |
| `graph()` / `stats()` | `Promise<object>` | |

Errors surface as `AtlasError` (with a `.status` for HTTP errors). Reads that
tolerate an offline server (`health`, `ready`) return `false` instead of throwing.

## Verify it works

Against a running server (`go run ./cmd/atlas-server`):

```bash
npm run smoke     # node --experimental-strip-types smoke_test.ts
```

The smoke test drives the full path: issue → verify (accept) → revoke →
verify (reject / `RevokedObservable`), against the real engine.
