---
date: 2026-07-21
slug: duplicate-key-rejection
artifact: hardening within the AD-024 doctrine — the verifier rejects authenticated payloads that repeat a JSON member name
decision: decodeClaims now rejects any payload containing a duplicate JSON object member name (at any nesting depth) as Altered, before decoding. This is a strengthening of the existing "authentic signature is necessary but not sufficient" refusal (AD-024), not a new answer-set member — a duplicate-key payload is malformed issuance output, so the record is Altered via the existing IntegrityFailed cause. Strictly fail-closed: it can only reject, never accept. Proven by three cases in the record-package integrity test, two language-neutral adversarial vectors (authentic-duplicate-scope-key, authentic-duplicate-sub-key), and a normative MUST rule in tests/vectors/VECTORS.md.
agents_consulted: [cryptography-specialist, conformance-testing, red-team]
overrides: false
related_entries: [language-neutral-conformance-vectors, open-source-governance-baseline]
---

# Context

The 13-agent review (`agents/review-team/output/`, CRYPTO-9 / TEST-2) named
duplicate-JSON-key handling as an unspecified conformance differential. It is
**not** exploitable in Atlas today: decode runs on signature-verified bytes, so
an attacker cannot inject a duplicate key without the issuer's key, and Seal
emits canonical single-key JSON. But the wire contract left the behavior
undefined, and Go's `encoding/json` silently takes last-wins. A second-language
verifier reading first-wins would then disagree with the reference over the
*same signed bytes* — the Frankencerts failure class. For the security-bearing
claims this is a latent silent-acceptance hazard: last-wins on
`{"sub":"…/p","sub":"…/evil"}` accepts the wrong principal.

# Decision

`decodeClaims` calls `rejectDuplicateKeys(payload)` before unmarshalling. The
check walks the JSON token stream and, per open object, refuses a repeated
member name at any nesting depth (arrays carry no names to collide). A duplicate
returns an error, which `ValidateIntegrity` surfaces as `Altered`.

This is deliberately framed as **within AD-024**, not as a new AD into the
closed architecture set: the answer set is unchanged (no new `Cause`), the
outcome is the existing `Altered`/`IntegrityFailed`, and the refusal is the same
family as the already-decided "authentic-but-malformed payload → Altered" rule
that the suite tests for unsorted scope, missing sub, non-SPIFFE identities,
etc. The `TestValidateIntegrityAuthenticButMalformedPayload` table gains three
cases (duplicate `scope`, duplicate top-level `sub`, duplicate nested `act.sub`)
and the same test asserts the forward-compat "unknown extra field stays Intact"
case still holds — proving the check does not over-reject append-only evolution.

The normative rule is written into `tests/vectors/VECTORS.md` (the wire-format
spec a foreign implementer reads): "Duplicate JSON member names ⇒ reject (MUST)."

# Correction (same day)

An earlier draft of this entry deferred the language-neutral vector on the
belief that the top-level `README.md` was frozen and its "28 vectors (18
adversarial)" count could not be updated without the CONTRIBUTING §4 amendment.
That was wrong: a substring grep had false-matched `lab/LAB_README.md`; the
top-level `README.md` is **not** in `scripts/frozen-docs.list`. On verifying
this, the vector was added properly. Recording the mistake rather than editing
it away, per the project's honesty discipline.

# Dissent / limits

- Two language-neutral adversarial vectors were added to
  `tests/vectors/negative-vectors.json` (now 20 adversarial / 30 total), and the
  (non-frozen) `README.md` count was updated to match. `VECTORS.md` carries the
  normative MUST rule.
- No frozen document was modified. `make ci` (build, vet, tests, lints,
  frozen-doc integrity, import rules) stays green.
- Scope of protection is honest: this closes a *differential/robustness* gap and
  removes a latent hazard; it is not the fix for a live exploit, because none
  existed here.
