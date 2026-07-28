# Known Limitations

Published deliberately. Every system has these; most hide them. If a
limitation below blocks your use case, Atlas is the wrong tool today — see
`docs/product/WHY.md` §Q8 for when *not* to use it.

## 1. Single-hop delegation only

A principal delegates to a delegate. **The delegate cannot re-delegate.**
There is no A→B→C chain — by design, not omission: multi-hop attenuation
multiplies revocation and audit complexity, and the two-domain single-hop case
is the one we can make bulletproof. If you need deep delegation chains today,
look at Biscuit. (Multi-hop is a researched, deliberately deferred extension —
`docs/discovery/`, DEFERRED D3–D4.)

## 2. Two trust domains, no federation

Atlas is scoped to *one issuing domain, one relying domain*. Three-plus-domain
federation, brokered trust, and transitive bundles are out of scope (C3,
ER17).

## 3. Within-window replay is not prevented (FM8)

A valid record intercepted in flight can be presented by anyone until it
expires or is revoked — records are bearer assertions. Mitigations available
today: short TTLs (default 1h; set minutes for agent tasks), scope narrowness,
and revocation. Proof-of-possession binding (record usable only by the
delegate's own key) is the roadmapped fix, not a shipped one.

## 4. Key rotation is manual

Trust material supports multiple keys (`kid`-pinned) and bundles carry the key
set, so rotation is *designed for* — but there is no automated rotation
ceremony, overlap-window tooling, or revocation-of-a-key story shipped yet.
The authority key sits in one PEM file (0600). An HSM/KMS integration does not
exist.

## 5. The server is a dev-grade deployment today

TLS (`-tls-cert`/`-tls-key`), per-IP rate limiting (`-rate-limit`), optional
bearer auth on mutating endpoints, pinned CORS, and a hardened container
(distroless, nonroot, read-only rootfs) do now exist. What is still missing is
what makes it *dev-grade*: **no multi-tenant isolation, single-node only, no
HA, no horizontal scale.** The state store is an atomic-write JSON snapshot —
correct, durable, and honest about scale: it is not a database. The seams for
replacing it (a `Store` interface, a stateless engine) exist; a clustered
deployment does not.

One consequence worth stating plainly, and one partial mitigation.

By default a single server instance holds **one trust domain and one global
delegation graph.** Every caller shares the same audit log and the same
revocation set. That is correct for a single-tenant deployment, and it is what
the CLI, the SDKs, and the MCP server all use.

For the public sandbox there is **per-session isolation**: a client that sends
an `X-Atlas-Session` header gets its own delegation store, audit log, and
revoked set, so one visitor cannot see or revoke another's capabilities
(`cmd/atlas-server/session_test.go` proves both directions). Be clear about what
this is and is not:

- It **is** a sandbox boundary, sufficient for a shared public demo.
- It is **not** multi-tenancy. There is no authentication behind a session id,
  no authorization model, no quota, and no durability — sessions are in-memory,
  capped at 500, and evicted after 30 minutes idle. Anyone who learns a session
  id has that session's state.
- The authority key and trust material are **shared** across sessions by
  design: every visitor is issued by the same real authority.

Do not mistake it for a tenancy feature in a production deployment.

## 6. Scope strings are uninterpreted — attenuation is only as sound as your vocabulary

Atlas performs **set inclusion on opaque strings**. It has no opinion on what a
scope *means*, and cannot acquire one without becoming a policy engine.

So this attack needs no forgery at all: obtain a legitimate capability for
`read:orders`; all five checks pass; the resource interprets `read:orders` as
"read *all* orders" rather than "read the orders in this request's context".
Attenuation was cryptographically perfect and semantically meaningless.

**Deployment obligation:** scope strings MUST map one-to-one onto enforcement
boundaries the resource actually implements. A scope coarser than its
enforcement is a vulnerability, not a naming preference. Avoid wildcard scopes
(`deploy:*`) — they grant authority over resources that do not exist yet. And
note that individually-harmless capabilities can aggregate: read-customer-data
plus write-external-bucket plus trigger-job is exfiltration, and Atlas verifies
each independently with no notion of a holder's total authority.

## 7. Single-hop is a property of the model, not of your deployment

The delegate cannot re-delegate *in Atlas*. But when service B verifies a
capability and then calls service C **with B's own credential**, the chain
A→B→C exists in the real system and Atlas cannot see, express, or audit it.

The blast-radius containment demonstrated in `examples/ship-a-landing-page.sh`
holds because every tool there is a **terminal resource**. If any tool were
itself a delegating service, containment would end at that boundary. Atlas
bounds what a *verifier* must reason about, not what a *system* does.

## 8. The audit log is not tamper-evident

The audit trail is a plain JSON file: mutable by anyone with filesystem access,
no hash chain, no signatures. It is useful for operations and **must not be
relied on for forensic non-repudiation**. Signing or chaining it is roadmapped,
not shipped.

## 9. Issuer availability is moved, not eliminated

Atlas removes the issuer from the **hot verification path**. It does not remove
it from the system. Short TTLs — the recommended mitigation for the bearer-token
window and several other gaps — require the issuer reachable once per TTL
period. What you gain is that a partition *during a capability's lifetime* does
not stop work; verification keeps answering. That is narrower than "no issuer
dependency," and the difference matters when sizing TTLs.

## 10. NTP integrity is a security dependency

Both expiry and freshness are denominated in wall-clock time, and there is no
cryptographic way to verify a clock. An attacker who controls a verifier's time
source can resurrect expired capabilities or make stale snapshots appear fresh,
within the skew tolerance and beyond it. Authenticated time is out of scope and
is assumed.

## 11. The trust domain is an index; the key is the boundary

Two organizations both using `domain-a.test` — the default shipped in every
example — are separated by **key material alone**. The trust-domain string
selects which material to use; it does not isolate. Cross-domain verification
fails at the integrity check, so the failure is safe, but do not read the domain
name as a security boundary.

## 12. Freshness costs staleness — physics, not a bug

Offline verification answers from the last snapshot the verifier holds. Your
staleness budget (`--max-staleness`, policy R) is exactly your revocation
blind spot: R = 5 minutes means a revocation may take up to 5 minutes to bite
for that verifier. Atlas's contribution is making this bound *explicit,
chosen, and fail-closed* — not making it disappear. Nothing can make it
disappear; distrust anyone who claims otherwise.

## 13. Your effective revocation latency is the WORST-configured verifier

R is chosen per verifier, by whoever operates that verifier. The security
property, however, is global.

An attacker holding a revoked capability does not present it at your strictest
verifier. They present it at the **loosest** one — the service where a team set
R to 24 hours after an outage, to stop the pager. Your fleet-wide revocation
window is the maximum R in the fleet, not the value you configured.

Today an issuer cannot constrain R, and nothing exposes a verifier's effective
R, so an operator cannot even *find* the outlier. Treat R as a fleet-wide
policy that must be audited, not a per-service knob. An issuer-side ceiling is
tracked as a design change.

## 14. In-partition revocation is invisible until recovery (INV12 / S4)

If the network is split when a revocation is published, a partitioned
verifier cannot observe it until the partition heals + one snapshot delivery.
Within its staleness budget it will keep accepting; past the budget it fails
closed. This is stated, tested, and demonstrated — never papered over.

## 15. The zero-egress packet proof awaits real infrastructure

Offline verification is proven at the process level (server dead, verify
answers — `examples/unforgettable.sh`) and by construction (no network
dependency in `internal/verify`). The *link-level* proof — two real SPIRE
domains, a severed cable, tcpdump showing zero egress — is fully authored in
`atlas-lab/` but requires a Docker/SPIRE host this project's sandbox doesn't
have. No result is claimed until a real host produces one.

## 16. Cryptography is deliberately boring — and unaudited

ES256/P-256 JWS via Go's standard library; no novel constructions anywhere.
That's a feature. But no third-party security audit has been performed, and
until one has, treat Atlas as pre-production for high-stakes use.

## 17. Young ecosystem

One reference implementation of the **verification core** (Go, `internal/`).
Client SDKs ship for **Go, Python, and TypeScript** (`sdk/`), but these are thin
clients that mirror the server API — they are *not* independent re-implementations
of the verifier, and no independent verifier implementation (Rust/Zig/…) exists
yet. That distinction matters: the 30-vector language-neutral conformance suite
(`tests/vectors`) exists precisely so that independent verifiers can appear and
be shown to agree, byte-for-byte, before anyone relies on them. Until a second
verifier passes the suite, cross-implementation agreement is a design goal, not a
demonstrated fact (see the review stack item on differential-in-CI).
