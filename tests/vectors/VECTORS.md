# Atlas conformance vectors

Language-neutral conformance vectors for the Atlas Verification Core (M3).
**Any** verifier implementation — in any language — is conformant iff, for
every vector, it produces the expected `decision` and `causes` from the
vector's inputs. This is the interop backbone: it turns the prose conformance
definition and the Go conformance kit into an artifact a Rust/Python/Zig/…
implementation can run without touching Go, and it is how verifier
differentials (the Frankencerts failure) are caught before they ship.

- **Files:** `verdict-vectors.json` (verdict-space cases) and
  `negative-vectors.json` (adversarial/malformed records that MUST be
  rejected). Both committed and authoritative.
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
garbage, the empty input, and **authentic-but-malformed payloads** (a real
signature over a payload that omits the principal or instance, carries a
non-canonical scope, or a non-SPIFFE identity — proving that signature
authenticity is necessary but not sufficient). A verifier that accepts any of
these has a silent-acceptance differential.

## Conformance

An implementation is conformant iff it produces the exact `expect.decision`
and the same set of `expect.causes` for every vector in both files. Compare
causes as a set (order is not significant). A single mismatch — a positive
vector rejected, or (worse) a negative vector accepted — is a verifier
differential and a conformance failure.
