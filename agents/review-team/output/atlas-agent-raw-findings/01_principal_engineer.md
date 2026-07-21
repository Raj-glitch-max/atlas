# 01 — PRINCIPAL ENGINEER — RAW FINDINGS

Phase 1 (Foundation). Reviewed the public surface, package boundaries, error taxonomy, and
protocol evolution posture against the real tree. Verdict at the end.

Context note that reframes half my usual checklist: **there is essentially no public Go API.**
`record/ issuance/ verify/ truststore/ revstatus/ revorigin/` are all under `internal/`, so they
cannot be imported outside the module. The public surface is `sdk/go`, `cmd/*`, and the wire
format. That is a *good* decision and it retires my "exported symbol leaks go-jose" red flag —
`grep go-jose` outside `internal/` and tests returns nothing. Confirmed clean.

---

FINDING-PE-1
  Severity:     MAJOR
  Category:     protocol-evolution
  Location:     internal/record/envelope.go (payloadClaims)
  Claim:        The delegation record carries no protocol version field.
  Failure:      In 18 months a second-language reimplementer adds a security-relevant claim. A
                v0.1 verifier has no version to switch on and no signal that it is reading a
                newer record; it applies old rules to new semantics. There is no flag day
                available because there is no version to gate one on.
  Fix:          Add a required, signed `atl_ver` (or reuse `typ` with a version suffix), verified
                at decode. Document the rule: unrecognized major version → reject.
  Cost of fix:  now = 2h (field + vector) | after v1 = a format break for every SDK.
  Owner:        01_principal_engineer

FINDING-PE-2
  Severity:     MAJOR
  Category:     protocol-evolution
  Location:     internal/record/envelope.go:decodeClaims (unknown-field tolerance)
  Claim:        Forward-compat is "tolerate unknown fields, append-only" with no `crit`-style
                mechanism to make a future field mandatory-to-understand.
  Failure:      The append-only rule is safe *today* because no tolerated field is
                security-load-bearing. The first one that is (a new caveat, an audience) will be
                silently ignored by an old verifier — a bypass. There is no mechanism to prevent
                that.
  Fix:          Define now, before any extension exists, how a critical field is signalled and
                that an unrecognized critical field rejects. Coordinate with 03 (this is their
                envelope call — CRYPTO-4).
  Cost of fix:  now = design note + 1 field | after first extension = a bypass window.
  Owner:        01_principal_engineer

FINDING-PE-3
  Severity:     MINOR
  Category:     honesty
  Location:     00_atlas_context.md §6; context §7; LIMITATIONS.md §10
  Claim:        The claim surface is internally inconsistent (see TW-2 / EM-10). The review's own
                context sheet lists "Delegation chains ✓" that the repo contradicts; the SDK count
                disagrees between README (3) and LIMITATIONS (1); benchmarks are unqualified.
  Failure:      A reviewer who reads two docs finds two Atlases. For a project whose entire pitch
                is *honesty*, an over-claim in the summary undoes the credibility the honest docs
                earned.
  Fix:          Make the top claim surface match `LIMITATIONS.md` (the honest floor). Qualify
                every benchmark with CPU / Go version / p50-p99 / chain-depth=1.
  Cost of fix:  now = 2h | after v1 = a "they over-claimed" first impression.
  Owner:        01_principal_engineer

**Things I checked and approve (REFUTED red flags):**
- `interface{}`/`any` in the verify path: none. Closed types throughout (Outcome, Cause, Verdict).
- `context.Context` + I/O in verify: none. `Verify([]byte)` takes no context; `internal/verify`
  imports no network (import lint enforces it). Offline-first is structural. Approve.
- Errors as strings: no — sentinel errors + closed `Cause`/`RefusalCause` enums, documented.
- Same struct for issuance-input and verify-output: no — `issuance.Request` ≠ `record.Assertions`
  (returned only from an authenticated read). Approve.
- Import graph / layering: `verify` does not import `issuance`; dependency rules R1–R7 are
  lint-enforced. Approve.
- README "production-ready" vs alpha: README says "v0.1-dev, reference implementation." No
  contradiction. Approve.
- MaxChainDepth / cycle detection: N/A — single-hop, no chain walker, so no depth-bomb surface.
  This is a *strength* of the single-hop scope, not a missing rail.

VERDICT: **APPROVE-WITH-CHANGES.** Ship the version field (PE-1) and the `crit` design (PE-2)
before the first extension; reconcile the claims (PE-3). The core API discipline is genuinely
strong — the internal/ boundary and the lint-enforced offline property are Staff-level moves.
