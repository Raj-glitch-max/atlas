# V1 Operation Runbook (reference implementation)

**Status:** Operational runbook for the V1 reference implementation
(IMPLEMENTATION_MASTER_PLAN.md E5-T5). Minimal by intent: V1 is a validated
reference, not a production-hardened system (C4). Not frozen.

This runbook covers the two out-of-band operator acts the architecture leaves
to operations (P3/P4 seams) and how to run the boundary drivers.

## Composition roots (drivers)

Three drivers wire the modules at their boundaries. In V1 each runs a
self-contained demonstration of its boundary (production request/response I/O
is deferred until a serialization format is chosen — TP6):

| Driver | Boundary | Demonstrates |
|---|---|---|
| `go run ./cmd/atlas-issue` | issuance (domain A) | proper-subset issuance (Issued) and over-scope refusal (Refused, nothing created) |
| `go run ./cmd/atlas-verify` | verification (RP, domain B) | Accept on a fresh not-revoked observation; InconclusiveRejected on the degenerate provider (fail-closed); latency measured (AT26, no threshold asserted) |
| `go run ./cmd/atlas-revoke` | revocation (domain A) | append-only, terminal register; repeated revocation is a no-op |

## Out-of-band act 1 — trust-material provisioning (M4 / P4)

The relying party verifies offline using trust material it already holds; the
system never fetches it (FM9). Provisioning is a manual operator act (the gate
C1 bundle-exchange pattern):

1. Obtain the issuing domain's public verification key(s) (P-256/ES256, per
   AD-012) and their key IDs, out of band.
2. Build `record.TrustMaterial` for the issuing trust domain from the keyID →
   public-key map.
3. `truststore.Store.Provision(material, now)` on the relying party. Absent
   material is a final answer: until provisioned, verification of that
   domain's delegations fails closed (Inconclusive), never fetches.

Rotation is a re-provision (a new operator act), recorded in the append-only
provisioning log.

## Out-of-band act 2 — permission source (M2 / P3)

The Issuance Authority checks a requested scope against the principal's
permission set through the `PermissionSource` port. Its realization is
deployment-specific and deferred; V1 wires a concrete source at the issuance
composition root. An unavailable source causes issuance to refuse
(`PermissionsUnavailable`) — it never guesses.

## What V1 does NOT operate

- **Revocation propagation** (author register → relying-party view): the
  channel (push/pull/cached-pull) is the spike-selected composition and is
  deferred (E7). Until then the relying party's revocation provider is the
  degenerate realization (answers Indeterminate → verification fails closed),
  or a test/observed provider in acceptance runs.
- **The two-domain SPIRE substrate**: stood up for EXP-001 and the
  substrate-blocked acceptance tests (E6); not part of day-to-day V1 wiring.
- **Any production hardening, HA, or scale** (C4).

## Verification gates

`make ci` runs the full local gate chain: pre-commit hooks, markdown lint,
frozen-doc integrity, secret scan, `go build`/`vet`, the dependency-rule
import lint (`scripts/check-imports.sh`), and `go test`. Every commit must
leave it green.
