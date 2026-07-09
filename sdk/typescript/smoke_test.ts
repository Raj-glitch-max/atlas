/**
 * Live smoke test for the Atlas TypeScript SDK.
 *
 * Requires a running server:  go run ./cmd/atlas-server
 * Run (Node 22.6+, no deps):  node --experimental-strip-types sdk/typescript/smoke_test.ts
 *   or with a base URL:        ATLAS_URL=http://127.0.0.1:8087 node --experimental-strip-types ...
 *
 * Exercises the full alpha-path: issue -> verify(accept) -> revoke ->
 * verify(reject/RevokedObservable), plus /health and /readyz.
 */
import { AtlasClient, Decision } from "./atlas.ts";

const url = process.env.ATLAS_URL ?? "http://127.0.0.1:8087";
const atlas = new AtlasClient(url, { apiKey: process.env.ATLAS_API_KEY });

function assert(cond: unknown, msg: string): void {
  if (!cond) {
    console.error("FAIL:", msg);
    process.exit(1);
  }
}

const PRINCIPAL = "spiffe://domain-a.test/workload/payments-api";
const DELEGATE = "spiffe://domain-b.test/agent/booking-worker";

const healthy = await atlas.health();
if (!healthy) {
  console.error(`atlas-server not reachable at ${url} — start it with:  go run ./cmd/atlas-server`);
  process.exit(2);
}
console.log("health          ok");
console.log("ready           " + (await atlas.ready() ? "ok" : "warming"));

const grant = await atlas.issue({ principal: PRINCIPAL, delegate: DELEGATE, scope: ["read:orders", "write:audit"] });
assert(grant.record && grant.instance, "issue returned a record + instance");
console.log(`issue           ${grant.instance}`);

const v1 = await atlas.verify(grant.record);
assert(v1.decision === Decision.Accept, `fresh record verifies Accept (got ${v1.decision})`);
assert(v1.trace.length === 5, `trace has five checks (got ${v1.trace.length})`);
console.log(`verify          ${v1.decision}  (${v1.latencyMicros}µs, ${v1.trace.length} checks)`);

await atlas.revoke(grant.instance);
console.log("revoke          ok");

const v2 = await atlas.verify(grant.record);
assert(v2.decision === Decision.Reject, `revoked record verifies Reject (got ${v2.decision})`);
assert(v2.causes.some((c) => /Revoked/i.test(c)), `reject cause mentions revocation (got ${JSON.stringify(v2.causes)})`);
console.log(`verify          ${v2.decision}  (${v2.causes.join(", ")})`);

console.log("\nALL PASS — TypeScript SDK works end-to-end against the real engine.");
