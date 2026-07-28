# Atlas conformance vectors

Language-neutral conformance vectors for the Atlas Verification Core (M3).
**Any** verifier implementation — in any language — is conformant iff, for
every vector, it produces the expected `decision` and `causes` from the
vector's inputs. This is the interop backbone: it turns the prose conformance
definition and the Go conformance kit into an artifact a Rust/Python/Zig/…
implementation can run without touching Go, and it is how verifier
differentials (the Frankencerts failure) are caught before they ship.

- **Files:** `verdict-vectors.json` (verdict-space cases),
  `negative-vectors.json` (adversarial/malformed records that MUST be
  rejected), and `revocation-vectors.json` (the revocation LAYER — see below).
  All committed and authoritative.
- **Source of truth:** generated from the Go conformance corpus
  (`tests/conformance`), recording the *reference verifier's actual verdict*.
  Regenerate with `go test ./tests/vectors -run TestVectorsRegenerate -update`
  and commit. The reference implementation is continuously checked against the
  committed file by `TestVectorsReplayAgainstReferenceVerifier`.

## File shape

```jsonc
{
  "schema": 1,                 // bump on any shape change; never repurpose a field
  "description": "...",
  "vectors": [ { /* Vector */ } ]
}
```

## Vector fields

| Field | Type | Meaning |
|---|---|---|
| `name`, `description` | string | human labels |
| `record` | string | the presented delegation record, a **compact JWS** (ASCII) |
| `trust.domain` | string | the trust domain the RP holds material for |
| `trust.keys` | array of JWK | RFC 7517 JWKs (EC P-256, `alg` ES256). **Empty ⇒ the RP holds no material for the domain** (→ `TrustMaterialAbsent`) |
| `revocation.state` | string | the revocation-observation answer supplied to the verifier: `Indeterminate` \| `NotObservedRevoked` \| `ObservablyRevoked` |
| `revocation.as_of` | string \| null | RFC3339 UTC currency of a knowledge answer; `null`/absent for `Indeterminate` |
| `now` | string | RFC3339 UTC — the verifier's clock reading |
| `policy.r_seconds` | number | revocation-observability bound R (seconds) |
| `policy.skew_seconds` | number | clock-skew tolerance (seconds) |
| `expect.decision` | string | `Accept` \| `Reject` \| `InconclusiveRejected` |
| `expect.causes` | array of string | the cause names (order-insensitive); see below |

## The record wire format (enough to consume the vectors)

The `record` is a **compact JWS** (`header.payload.signature`, base64url,
dot-separated) — this is the reference realization (see `internal/record`;
the normative wire format is a mechanism-RFC concern, sketched here to make
the vectors self-contained):

- **Header:** `alg` MUST be exactly `ES256`; `typ` MUST be exactly
  `atlas-record+jws`; `kid` MUST be present. Any other `alg` (including
  `none`), a missing/other `typ`, or a missing `kid` ⇒ the record is not
  authentic (integrity fails).
- **Payload claims:** `sub` (principal SPIFFE ID), `act.sub` (delegate SPIFFE
  ID, RFC 8693 actor form), `scope` (array, canonical: sorted, unique,
  non-empty entries), `exp`/`iat` (RFC 7519 NumericDate seconds, positive),
  `atl_ins` (opaque instance identity, non-empty), `atl_rvb` (optional opaque
  revocation binding, base64url). Unknown claims are tolerated (forward
  compatibility). A payload violating any of these ⇒ integrity fails.
- **Duplicate JSON member names ⇒ reject (MUST).** RFC 8259 permits an object
  to repeat a member name but leaves the result undefined. A verifier that
  resolves duplicates (e.g. last-wins) can disagree with a first-wins verifier
  over the *same signed bytes* — a conformance differential, and for the
  security-bearing `scope`/`sub` a silent-acceptance hazard (last-wins on
  `{"sub":"…/p","sub":"…/evil"}` would accept the wrong principal). A conformant
  verifier MUST reject any authenticated payload containing a duplicate member
  name at any nesting depth, at integrity time. Do **not** rely on your JSON
  library's default — most, including Go's `encoding/json`, silently take
  last-wins.
- **Signature:** ES256 over the signing input, verified with the JWK in
  `trust.keys` whose `kid` matches the header. An unknown `kid` ⇒ integrity
  fails (there is no key-rotation state in scope that makes it legitimate).

Note: the transport encoding is non-canonical (base64url padding-bit
malleability), so two byte-different records may decode identically. Integrity
protects **content**, not transport bytes; do not key anything on record byte
identity.

## Verification algorithm a conformant verifier runs (per vector)

Five checks; produce a verdict by the routing below. See `INTERFACE_SPECIFICATION.md` §3 for the normative contract.

1. **Integrity (gate):** select `trust.keys` for the record's principal
   trust domain. Absent ⇒ inconclusive `TrustMaterialAbsent`. Present but the
   record is not authentic ⇒ definitive `IntegrityFailed`. (Reads that follow
   require an authentic record.)
2. **Identity binding:** principal and delegate both recover and are
   **distinct** ⇒ pass; else definitive `BindingMismatch`.
3. **Expiry (± skew):** `iat > now + skew` ⇒ inconclusive
   `ClockBeyondTolerance`; else `now > exp + skew` ⇒ definitive `Expired`;
   else pass.
4. **Scope integrity:** scope present and well-formed (covered by the
   signature); malformed ⇒ definitive `ScopeIntegrityFailed`. (Subset-ness is
   an issuance-time property, NOT re-derived here.)
5. **Revocation (under policy):** `ObservablyRevoked` ⇒ definitive
   `RevokedObservable`. `NotObservedRevoked(as_of)` with `now - as_of` in
   `[-skew, R]` ⇒ pass; outside ⇒ inconclusive `RevocationKnowledgeStale`.
   `Indeterminate` ⇒ inconclusive `RevocationStatusIndeterminate`.

**Routing (order-independent):** any definitive cause ⇒ `Reject` (with the
definitive causes); else any inconclusive cause ⇒ `InconclusiveRejected` (with
the inconclusive causes) — this is the fail-closed posture, `[HYPOTHESIS]`
pending V1 confirmation; else ⇒ `Accept` (no causes).

## Cause names

Definitive: `BindingMismatch`, `IntegrityFailed`, `Expired`,
`ScopeIntegrityFailed`, `RevokedObservable`.
Inconclusive: `TrustMaterialAbsent`, `ClockBeyondTolerance`,
`RevocationStatusIndeterminate`, `RevocationKnowledgeStale`.
Reserved (not produced in V1): `SignatureUnverifiable`.

## Negative vectors (`negative-vectors.json`)

Adversarial and malformed records that a conformant verifier MUST reject.
Each holds a valid trust key, a fresh `NotObservedRevoked` answer, and an
in-window clock, so the **record's malformation is the sole rejection cause**.
Families covered: `alg=none` (with and without a stale signature),
algorithm substitution (`ES384`) and confusion (`HS256` using the public
key's coordinates as the HMAC secret), missing/wrong `typ`, missing/forged
`kid`, truncation, cross-record signature and payload transplants, non-JWS
garbage, the empty input, **duplicate JSON member names** (a real signature
over a payload that repeats `scope` or `sub` — where a last-wins reader and a
first-wins reader disagree over identical bytes), and **authentic-but-malformed
payloads** (a real signature over a payload that omits the principal or
instance, carries a non-canonical scope, or a non-SPIFFE identity — proving that
signature authenticity is necessary but not sufficient). A verifier that accepts
any of these has a silent-acceptance differential.

## Conformance

An implementation is conformant iff it produces the exact `expect.decision`
and the same set of `expect.causes` for every vector in both files. Compare
causes as a set (order is not significant). A single mismatch — a positive
vector rejected, or (worse) a negative vector accepted — is a verifier
differential and a conformance failure.

## Revocation-layer vectors (`revocation-vectors.json`, schema 1)

The two files above supply `revocation.state` as an **already-computed answer**
(`Indeterminate` | `NotObservedRevoked` | `ObservablyRevoked`). That tests a
verifier's *handling* of an answer and skips its *derivation* entirely: the
snapshot signature is never checked and freshness is never computed.

Since freshness-bounded, fail-closed revocation is the most novel mechanism in
Atlas, that made it the least conformance-tested one — an implementation could
pass all 30 record-level vectors while getting revocation semantics completely
wrong.

`revocation-vectors.json` carries the **signed snapshot itself**. A conformant
implementation must do the real work:

1. verify the snapshot signature against `snapshot_key` (an RFC 7517 JWK),
2. bind the snapshot to `expect_list` — a snapshot from a different revocation
   stream MUST NOT be adopted,
3. apply `ingest[]` **in order**, enforcing *strictly monotonic* `as_of`
   adoption, and match each outcome against `expect_adopt[]`,
4. compute the adopted snapshot's age against `now` and `policy.r_seconds`,
5. derive a revocation state and only then a verdict.

`ingest` is a sequence on purpose: rollback is only expressible as an ordered
pair, and a verifier that adopts snapshots without monotonicity can only be
caught by presenting a newer snapshot followed by an older one.

`expect_adopt` is asserted per-ingest rather than only through the final
verdict, so an implementation that adopts a forged or rolled-back snapshot is
caught at the point of the mistake.

The record-signing key and the revocation-origin key are **deliberately
different** in these vectors, even though the reference server uses one key for
both. That prevents an implementation from passing by accidentally verifying a
snapshot with the record key.

### Coverage

| Vector | Pins |
|---|---|
| `fresh-not-revoked` | the only case that may accept |
| `fresh-revoked` | definitive rejection, `RevokedObservable` |
| `stale-not-revoked-fails-closed` | **the fail-closed property** — not revoked, but unprovable, so refuse |
| `no-snapshot-is-indeterminate` | absence of a record is not evidence of non-revocation |
| `bad-signature-not-adopted` | corrupted signature ⇒ no knowledge, not the snapshot's contents |
| `stripped-entry-not-adopted` | the "hide the revocation" tamper |
| `foreign-key-not-adopted` | correct signature, wrong signer |
| `wrong-list-id-not-adopted` | cross-stream substitution |
| `rollback-rejected` | newer then older, both genuinely signed |
| `replay-same-asof-rejected` | adoption is *strictly* monotonic |
| `newer-snapshot-supersedes` | revocation stays terminal across updates |
| `boundary-exactly-at-R` | age == R is still fresh (an off-by-one here silently changes every deployment) |

Regenerate with `go test ./tests/vectors -run TestRevocationVectorsRegenerate -update`.
