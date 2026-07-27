# 05 — RED TEAM — RAW FINDINGS

Phase 2 (Security), after 04. I do not propose designs; I try to produce a record that verifies
and shouldn't, or a legal-operation sequence that composes into an illegal capability. I attacked
the assumptions 03 and 04 handed me, against the real code. Two lists, as always.

## CONFIRMED (exploits that work)

**None.** I could not produce a single false-ACCEPT or a widening composition against this
codebase. For a project this size that is unusual and I am reporting it straight.

## REFUTED (attacks tried, Atlas survived — this list is the security résumé)

EXPLOIT-1  envelope · `alg:none`
  Hypothesis: submit `alg:none` with no signature → accept.
  Result:     REFUTED. `ValidateIntegrity` rejects `header.Algorithm != ES256` before any key
              logic; neg-vector `alg-none` is CI-gated. Repro rejected as Altered.

EXPLOIT-2  envelope · HS256 confusion (public key as HMAC secret)
  Hypothesis: sign HS256 using the issuer's public-key coordinates as the MAC key → forge.
  Result:     REFUTED — structurally. No HMAC path exists; the verify key is an `*ecdsa.PublicKey`.
              Neg-vector present. Cannot even be expressed.

EXPLOIT-3  envelope · `kid` swap / header key injection (`jku`/`x5u`/`jwk`)
  Hypothesis: point `kid` at another key, or smuggle a `jwk` header, to choose the verifier's key.
  Result:     REFUTED. `kid` only *selects among operator-held trusted keys*; go-jose verifies
              against the passed key, never a header-embedded/remote one. No SSRF, no attacker key.

EXPLOIT-4  envelope · `typ` transplant
  Hypothesis: splice a signed revocation snapshot into a delegation slot.
  Result:     REFUTED. `typ` pinned to `atlas-record+jws`; snapshot is a different signed object.

EXPLOIT-5  attenuation · widening by re-ordering / confused deputy / degenerate caveat
  Hypothesis: compose caveats to widen scope across a chain.
  Result:     REFUTED / N/A. There is no multi-hop composition operator — Atlas is single-hop by
              design. Single-hop issuance enforces a *strict* proper subset (`isProperSupersetOf`
              requires `len(distinct) < set.Len()`), so a delegate cannot even equal its principal's
              set, let alone exceed it. No surface to widen.

EXPLOIT-6  attenuation · duplicate `scope` keys
  Hypothesis: `{"scope":[...read],"scope":[...admin]}` — verifier and policy read different values.
  Result:     REFUTED as an exploit. Decode runs on *authenticated* bytes; I cannot inject a
              duplicate key without a valid signature, and the issuer emits canonical single-key
              JSON. Handed to 03/07 as a *differential* (cross-impl), not a forge — see CRYPTO-9.

EXPLOIT-7  chain · cycle / depth bomb / sibling transplant
  Hypothesis: A→B→A loop, 100k-link chain OOM, splice a sub-record across chains.
  Result:     REFUTED / N/A. No chain walker exists; records are independent single-hop artifacts,
              each verified against its issuer's key. Nothing to loop, bomb, or transplant. The
              single-hop scope *removes* this entire attack class rather than defending it.

EXPLOIT-8  partition · stale-snapshot acceptance
  Hypothesis: present a credential to a verifier whose snapshot predates its revocation.
  Result:     REFUTED as a bug — it is by design and *bounded*: accepted only while `now - asOf ≤
              R`; past R → `RevocationKnowledgeStale` → inconclusive-reject. The window is documented
              (LIMITATIONS §6). Finding is only that the *cost* of the chosen R is unmeasured (→ 04
              THREAT-1, S7).

EXPLOIT-9  partition · snapshot rollback to un-revoke
  Hypothesis: feed an older signed snapshot to erase a revocation.
  Result:     REFUTED. `Ingest` adopts only `s.AsOf.After(current.AsOf)`; older-or-equal ignored.
              A tampered set fails the `listID‖asOf‖set` signature and is refused (C8).

EXPLOIT-10 dos · fleet-halt via distribution partition
  Hypothesis: block snapshots → fleet stops.
  Result:     CONFIRMED-BY-DESIGN, not a bypass. It halts (fail-closed), never false-accepts. This
              is a property, not a defect; the open item is the unmeasured availability cost (S7,
              OQ-4). Reported to 04/06, not as a Tier-0 exploit.

EXPLOIT-11 clock · advance/rewind verifier clock
  Hypothesis: rewind to accept an expired record.
  Result:     PARTIAL/owned. Within ±skew, no effect; beyond it, a *compromised local clock* can
              extend expiry — but that is compromised-runtime territory (out of scope). Handed to
              04 as THREAT-3: the boundary should say so explicitly.

EXPLOIT-12 downgrade · version downgrade
  Hypothesis: present an old-version record to skip a newer check.
  Result:     N/A today (no version field, single format) — which is *why* PE-1/CRYPTO-4 matter
              before a v2 exists. There is no downgrade surface now; there will be one the moment a
              second version ships without a version gate. Pre-emptive, not present.

## MANDATORY CLOSING — TWO LISTS
- **CONFIRMED:** (none).
- **REFUTED:** EXPLOIT-1..9, 11 survived; 10 is by-design; 12 is pre-emptive. This survived-attack
  list is a genuine asset — it should go, verbatim, into the TAG-Security self-assessment (→ 12).

Note to 07: nothing new to freeze — the existing 18 adversarial vectors already encode
EXPLOIT-1..4/6. Add the duplicate-key differential (S9) and, when a version field lands, a
downgrade vector.
