# 04 — SECURITY ENGINEER — RAW FINDINGS

Phase 2 (Security), first. I consume 01's API facts and 03's approved primitives. I treat the
threat model as a document that must be falsifiable, and its *omissions* as the deliverable.

Unusual for this seat: **the threat model already exists and is good.** `THREAT_MODEL.md` is
RFC-6819-shaped — adversaries A1–A6, claims C1–C10 each mapped to a named test, explicit
boundaries, residual risk per adversary. `LIMITATIONS.md` publishes the gaps honestly. So my job
is less "there is no threat model" and more "interrogate its exclusions once, and find what it
still doesn't cover."

## THREAT TABLE (interrogated against the real code + docs)

| # | Asset | Attacker | Threat | Mitigation | State | Status |
|---|---|---|---|---|---|---|
| T1 | Record | Network (Dolev-Yao) | Forgery | ES256, curve enforced in trust material | TESTED (256+ mutations, fuzz) | ✅ |
| T2 | Record | Network | Algorithm confusion | `alg` pinned pre-verify, no HMAC path | TESTED (neg-vectors) | ✅ |
| T3 | Record | Holds valid record for B | Transplantation to C | **No `aud`**; records are bearer assertions | ASSERTED as out-of-scope (FM8) | ⚠️ owned |
| T4 | Scope | Holds valid record | Scope widening | Strict subset at issuance (single-hop) | TESTED (`OverScope`) | ✅ |
| T5 | Snapshot dist. | Partition | Fleet fail-closed (self-DoS) | Staleness → inconclusive-reject; R bounds it | IMPLEMENTED; cost UNMEASURED | ⚠️ |
| T6 | Snapshot | Replay old snapshot | Un-revocation via rollback | Monotone `AsOf` adoption (`Ingest`) | TESTED (C8, rollback_test) | ✅ |
| T7 | Clock | Controls verifier clock | Extend expiry | ±skew tolerance; injected clock | IMPLEMENTED; clock-trust boundary thin | ⚠️ |
| T8 | Verifier CPU | Deep/looping chain | DoS | **N/A — single-hop, no chain walker** | — | ✅ |
| T9 | Delegate runtime | Compromise B | Lateral use within B's scope | Attenuation + short TTL bounds blast radius | Documented (A2/§5) | ⚠️ owned |
| T10 | Cross-domain | Valid record from A | Unauthorized action in B | Verify ≠ authorize; verdict is a value | Boundary stated (§4.5) | ✅ |

---

THREAT-1  (T5 — the finding that orbits everything)
  Asset:            Snapshot distribution → the whole relying fleet's freshness
  Attacker:         A1, can partition snapshot delivery
  Threat:           Fleet-wide fail-closed once every verifier ages past R → availability collapse
  Class:            availability
  Stated?:          YES — `THREAT_MODEL.md` §5 A1 and `LIMITATIONS.md` §6–7 name it honestly
  Mitigation:       R bounds the window; fails closed, never false-accept
  Mitigation state: IMPLEMENTED (the bound) / **UNMEASURED (the cost)**
  Residual risk:    An operator cannot choose R responsibly because the availability curve as a
                    function of R and snapshot cadence has never been measured — the two-domain
                    substrate proof is authored but unrun (LIMITATIONS §8).
  Handoff to 05:    Try to make the fleet halt and measure how fast. The finding is not "it halts"
                    (by design) — it's "the halt-vs-R tradeoff is an unmeasured number." (→ S7, OQ-4)

THREAT-2  (T3 — transplantation / no audience binding)
  Asset:            A delegation record
  Attacker:         A1/A2, intercepts a valid record in flight
  Threat:           Present it to a different relying party until it expires/revokes (replay)
  Class:            authorization
  Stated?:          YES — carried openly as FM8 (`LIMITATIONS.md` §3), PoP roadmapped
  Mitigation:       Short TTL (default 1h; minutes for agent tasks), scope narrowness, revocation
  Mitigation state: ASSERTED (TTL) — no `aud`, no holder binding shipped
  Residual risk:    Within-window replay works. Honestly disclosed, but "replay resistance" in the
                    top-level guarantee list overstates a mitigation that is really "short TTL."
  Handoff to 05:    Confirm replay within TTL works and confirm the disclosure is accurate (it is).

THREAT-3  (T7 — clock trust)
  Asset:            Expiry enforcement
  Attacker:         Controls the verifier's clock (host/VM/container/NTP)
  Threat:           Rewind to accept an expired record, or advance to reject a valid one
  Class:            integrity/availability
  Stated?:          PARTIAL — skew tolerance (±30s) is documented (C4); clock *manipulation* as an
                    attacker capability is not explicitly in/out of scope
  Mitigation:       Injected clock port + bounded skew; future-dated issuance → inconclusive
  Mitigation state: IMPLEMENTED (skew) / clock-compromise UNSTATED
  Residual risk:    A verifier that trusts a compromised local clock extends expiry. This is
                    adjacent to "compromised runtime" (out of scope), but the boundary should say
                    so explicitly rather than leave it implied.
  Handoff to 05:    Advance/rewind the clock; confirm the boundary and that it degrades to the
                    stated exclusion, not to a surprise.

THREAT-4  (governance/process, not a runtime threat)
  Asset:            The disclosure process itself
  Threat:           A researcher who finds a bug has a one-line SECURITY.md and no terms
  Class:            process
  Stated?:          Partially — SECURITY.md exists (5 lines) but no SLA, supported-versions, or CVE
  Mitigation state: ABSENT
  Residual risk:    For a crypto project this is a credibility gap and a CNCF blocker.
  Fix:              Expand SECURITY.md (S3). This is where OQ-2's embargo half lives.

## SCOPE-EXCLUSION INTERROGATION (challenged once, then honored — R6)
- *Compromised issuer key (A6):* honored. Bounded by expiry + `kid` rotation; automated rotation
  not shipped (LIMITATIONS §4). Load-bearing and correctly stated.
- *Compromised delegate runtime (A2/T9):* honored — but the *blast radius* is documented (attacker
  gets exactly S until T or revocation), which is Atlas's best marketing. Keep it visible.
- *SPIFFE/SPIRE compromise:* honored (dependency, not scope). Correct.
- *Business policy / authorization:* honored, and cleanly — the verdict is a *value* separating
  "cryptographically valid" from "authorized," so the org-boundary hole I usually find is closed
  by design (T10). Approve.

## ASSUMPTIONS HANDED TO 05 (falsifiable claims — break exactly these)
1. No input produces a false ACCEPT (fail-closed is total).
2. Stale snapshot beyond R is always inconclusive-reject, never accept.
3. An older snapshot can never un-revoke (monotone adoption).
4. Single-hop scope cannot be widened at issuance (strict subset).
5. Replay within TTL is *possible* and *disclosed* — confirm the disclosure matches reality.
6. The fleet-halt is real by design; the *cost* is unmeasured — measure it.

VERDICT: **APPROVE-WITH-CHANGES.** No mandatory control is missing from the library; the fail-closed
floor is built and tested. The debt is *evidentiary* (T5 cost) and *process* (SECURITY.md), plus
one honesty trim: "replay resistance" should read "short-TTL-bounded replay window."
