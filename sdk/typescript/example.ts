/**
 * Minimal usage example for the Atlas TypeScript SDK.
 * Run:  go run ./cmd/atlas-server   # in one terminal
 *       node --experimental-strip-types sdk/typescript/example.ts
 */
import { AtlasClient } from "./atlas.ts";

const atlas = new AtlasClient("http://127.0.0.1:8087");

// Issue a scoped, expiring, revocable capability — not a shared secret.
const grant = await atlas.issue({
  principal: "spiffe://domain-a.test/workload/payments-api",
  delegate: "spiffe://domain-b.test/agent/booking-worker",
  scope: ["read:orders", "write:audit"],
  ttlSeconds: 900,
});
console.log("issued:", grant.instance, "→", grant.scope.join(", "));

// The delegate side verifies it — in microseconds, with a full decision trace.
const result = await atlas.verify(grant.record);
console.log("verdict:", result.decision, `(${result.latencyMicros}µs)`);
for (const step of result.trace) {
  console.log(`  ${step.check.padEnd(18)} ${step.outcome}`);
}

// Revoke it; verification now fails closed.
await atlas.revoke(grant.instance);
console.log("after revoke:", (await atlas.verify(grant.record)).decision);
