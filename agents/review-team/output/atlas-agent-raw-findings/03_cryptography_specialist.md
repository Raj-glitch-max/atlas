# 03 — CRYPTOGRAPHY SPECIALIST — RAW FINDINGS

Phase 1 (Foundation). Final authority inside the JWS envelope. I read the actual bytes and the
actual go-jose call sites, not the prose. Severity is encoded in tone; a bypass would be one line.
There are no bypasses. That is a real finding, and I do not say it reflexively.

The envelope (`internal/record/{envelope,integrity,seal,trustmaterial}.go`):
- `signatureAlgorithm = jose.ES256`, pinned as a constant. `ValidateIntegrity` checks
  `header.Algorithm != string(signatureAlgorithm)` → Altered, **before** any key is applied.
- `alg` is never used for dispatch: the key handed to `jws.Verify(key)` is the operator-held
  `*ecdsa.PublicKey` resolved by `kid`. `jku`/`x5u`/`jwk` header key-injection cannot fetch or
  select a key — go-jose verifies against the key you pass. RFC 7515 §10.7 respected.
- HS256 confusion is *structurally* impossible: there is no HMAC path, and the verification key is
  an `*ecdsa.PublicKey`, not `[]byte`. Curve is enforced at `NewTrustMaterial` (`key.Curve !=
  elliptic.P256()` → refused), so algorithm confusion is unrepresentable in held material, not
  merely rejected at runtime.
- `typ` pinned to `atlas-record+jws`; cross-type transplantation of a signed revocation snapshot
  into a delegation slot fails on `typ`.
- Compact-only shape guard (`isCompactJWS`) + `len(jws.Signatures) != 1` guard shrink the parse
  surface before crypto.
- `negative-vectors.json` (18 vectors) exercises `alg:none` (with/without stale sig), ES384,
  HS256-with-public-key-as-secret, missing/wrong `typ`, missing/forged `kid`, truncation,
  cross-record signature/payload transplants, and authentic-but-malformed payloads. Generation
  *fails* if any is accepted. `FuzzVerify` reports 821k execs with no silent acceptance.

---

FINDING-CRYPTO-4
  Class:        KEY-LIFECYCLE (forward-compat)
  Severity:     MEDIUM
  RFC ref:      RFC 7515 §4.1.11 (crit)
  Location:     internal/record/envelope.go:decodeClaims
  Attack:       Not exploitable today. Future: a v2 issuer adds a security-relevant claim; a v0.1
                verifier tolerates unknown fields and ignores it, accepting a delegation whose new
                restriction it cannot see.
  Precondition: A future critical extension field exists. None does yet.
  Impact:       Silent acceptance of an under-understood restriction (a bypass) once such a field
                ships.
  Test vector:  (future) a record carrying an unrecognized critical marker MUST reject.
  Fix:          Define the critical-field mechanism now (a `crit`-equivalent or the version gate
                in PE-1). An old verifier meeting an unrecognized critical marker rejects. This is
                what fail-closed means; I win this over 01's forward-compat concern (R1).
  Owner:        03_cryptography_specialist

FINDING-CRYPTO-9
  Class:        HARDENING (differential)
  Severity:     MEDIUM
  RFC ref:      RFC 8725 (JSON handling); RFC 7515
  Location:     internal/record/envelope.go:decodeClaims (json.Unmarshal)
  Attack:       Not exploitable in Atlas (decode runs on *authenticated* bytes; an attacker cannot
                inject duplicate keys without a valid signature, and the issuer emits canonical
                JSON). BUT: Go's `encoding/json` is last-wins on duplicate keys and this is
                unspecified in the wire contract. A second-language SDK with first-wins semantics
                diverges — a conformance differential, not a live forge.
  Precondition: A second implementation reads duplicate keys differently.
  Impact:       Cross-implementation verdict divergence (07's Frankencerts failure class).
  Test vector:  A record whose payload JSON contains a duplicate `scope` key → all conformant
                verifiers MUST agree (reject is the clean rule).
  Fix:          Spec rule: duplicate object keys → reject. Add the conformance vector (→ 07 / S9).
  Owner:        03_cryptography_specialist

FINDING-CRYPTO-12
  Class:        SUPPLY-CHAIN
  Severity:     MEDIUM
  RFC ref:      SLSA v1.0
  Location:     .github/workflows/, release process
  Attack:       A backdoored build of a *verification* library verifies anything. No signed
                provenance means an adopter cannot distinguish an authentic release from a tampered
                one.
  Precondition: An adopter consuming a binary/module release.
  Impact:       Full trust compromise of every downstream verifier.
  Fix:          Reproducible build + signed provenance (SLSA ≥ 2); pin go-jose v3 and gate it with
                govulncheck in CI. Co-own with 06 (S10).
  Owner:        03_cryptography_specialist

FINDING-CRYPTO-14
  Class:        HARDENING
  Severity:     LOW
  RFC ref:      FIPS 186-4 / SEC1 (point validation); RFC 7638 (thumbprint)
  Location:     internal/record/trustmaterial.go:NewTrustMaterial; integrity.go:keyFor
  Attack:       (i) `kid` selects among held keys by label rather than by a recomputed thumbprint
                binding; safe here because all held keys are operator-trusted, but a thumbprint
                binding is stronger. (ii) `NewTrustMaterial` checks `Curve == P256` but relies on
                Go's `ecdsa.Verify` for on-curve/in-range point validation of the coordinates.
  Precondition: Operator provisions a malformed or attacker-influenced public key.
  Impact:       None demonstrated (trust material is a trusted, out-of-band input); defense in
                depth only.
  Fix:          Optional: bind `kid` to an RFC 7638 thumbprint; add explicit on-curve/non-identity
                validation at import. Hardening, not a gate.
  Owner:        03_cryptography_specialist

**Minor precision note (not a finding):** `envelope.go`'s comment claims the compact-shape guard
keeps "the bytes presented and the bytes signed over in one-to-one correspondence." Property test
P1 correctly documents that base64url padding-bit malleability lets two byte-different records
decode to identical content. The signature protects *content*, not transport bytes — which is
correct and safe — but the comment overstates it. Reword to match P1.

---

## CRYPTOGRAPHIC ASSUMPTIONS HANDED TO AGENT 05 (attack these directly)

1. ES256/P-256 via Go stdlib + go-jose v3 is correct; no novel construction exists.
2. `alg` never dispatches; the key is chosen out-of-band by `kid` from trusted material.
3. No symmetric/HMAC path is reachable from `verify/`.
4. `typ` cryptographically separates delegation records from revocation snapshots.
5. The revoked-set snapshot signature covers `listID‖asOf‖sorted-revoked-set`; monotone adoption
   (`AsOf` strictly increasing) prevents rollback-to-un-revoke.
6. Signatures are never used as identifiers/cache keys, so ECDSA malleability is inert.
7. Held trust material is trusted and operator-provisioned; `kid` label-selection among it is safe.

VERDICT: **APPROVE.** The envelope resists the JOSE catalogue by construction. Address CRYPTO-4
before the first extension and CRYPTO-9/12 for cross-impl + supply-chain credibility.
