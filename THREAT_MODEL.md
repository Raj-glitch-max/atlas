# Atlas Threat Model

This is the product-facing threat model: who the adversaries are, what they
can and cannot achieve, and **which test proves each claim**. It builds on the
frozen security corpus (`docs/engineering/02_SECURITY_OBJECTIVES.md`,
`03_SYSTEM_INVARIANTS.md`, `04_FAILURE_MODEL.md`, summarized in
`context/07_SECURITY_POLICY.md`) and extends it to the product surfaces
(server, CLI, MCP, trust bundle) that post-date those documents. Nothing here
weakens a frozen boundary; the frozen documents win on any conflict.

Honest companion: `LIMITATIONS.md` — what Atlas does **not** defend, stated
plainly.

---

## 1. What Atlas protects

A **delegation record**: a signed assertion that *principal P grants delegate
D the scope S until time T, revocable via instance I*. The security goal is
that a conformant verifier's decision about that assertion is correct even
when the verifier is offline, and that no party — including the network, the
delegate, or a bundle courier — can silently alter the answer.

## 2. Adversaries

| Adversary | Capability assumed |
|---|---|
| **A1 · Network attacker** | Read/modify/replay anything in transit; partition any link at any time |
| **A2 · Malicious delegate** | Holds a valid record; wants more scope, more time, or survival past revocation |
| **A3 · Record forger** | No key; can craft arbitrary bytes and present them as records |
| **A4 · Bundle courier** | Delivers trust bundles to offline verifiers; may tamper with them |
| **A5 · Rogue caller of the server** | Network access to atlas-server's API |
| **A6 · Compromised issuer key** | Holds the authority's private key |

## 3. Claims, and the test that proves each

| # | Claim | Mechanism | Proven by |
|---|---|---|---|
| C1 | A record altered in any way is rejected (**INV8**) | ES256 over the full payload; alg pinned; `typ` pinned | 256+-mutation corpus (`internal/record`), fuzzer (`FuzzVerify`, 821k execs), adversarial vectors 1–18 |
| C2 | A delegate can never hold more scope than its principal (**INV2**) | Strict proper-subset check at issuance; scope covered by signature | issuance tests; `OverScope` refusal live in every demo |
| C3 | Forged/garbage/confused records never verify | Closed algorithm set (no `alg=none`, no HS256 confusion), kid pinning, structural checks | `negative-vectors.json` (18 adversarial vectors; generation fails if any is accepted) |
| C4 | Expiry is enforced within a declared skew (**INV3**) | Clock port + ±30s tolerance, boundary-tested | expiry boundary tests, conformance vectors |
| C5 | A revoked delegation is rejected once the revocation is observable (**INV4**) | Signed revoked-set snapshots, monotone adoption | E2E tests; `examples/unforgettable.sh` steps 5–6 |
| C6 | Verification needs no live authority (**INV7**) | All inputs local: trust material + snapshot + clock | offline CLI verify with the server dead; (link-level packet proof: atlas-lab, pending real host) |
| C7 | Stale revocation knowledge fails **closed**, never open (**INV12/S4**) | Freshness bound R; staleness → `Inconclusive`, routed to reject | `TestOfflineStaleSnapshotFailsClosed`; hypothesis tests in `internal/verify` |
| C8 | A tampered trust bundle is refused, not silently trusted | Snapshot signature covers listID+asOf+revoked-set (canonical digest) | `TestOfflineTamperedBundleRefused`; demo step 8 |
| C9 | Every decision is independently reconstructable (**INV9**) | Unconditional five-check decision trace, no hidden state | trace emitted on every verdict, incl. accepts |
| C10 | Mutating the server requires authorization (when configured) | Bearer token, constant-time compare | `TestAPIKeyGuardsMutations` |

## 4. Security boundaries (explicit)

1. **The verifier's decision** is the security boundary — not transport.
   Records are designed to cross hostile networks; nothing about their
   confidentiality is assumed (they are capability assertions, not secrets —
   but see A2/replay in LIMITATIONS).
2. **The issuance boundary**: only the authority's key can mint; permission
   attenuation is enforced *here*, so a verifier never needs the principal's
   permission list.
3. **The trust bundle** is integrity-protected (signature), not
   confidentiality-protected. Couriers can read it; they cannot alter it
   undetected — including *removing revocations* (C8).
4. **The server's mutating API** (`/issue`, `/revoke`) is guarded by bearer
   token when configured; read/verify paths are deliberately open (verifying
   a record you hold is not a privilege).
5. **Out of the boundary**: the workload-identity layer beneath Atlas
   (SPIFFE/SPIRE compromise voids delegation integrity — dependency, not
   scope), and non-conformant relying parties (we cannot force anyone to
   honor a verdict).

## 5. What the adversaries actually get (residual risk)

- **A1** can delay or partition. Result: verifiers keep answering from their
  last fresh snapshot **within their own staleness budget R**, then fail
  closed. An attacker who controls the network can cause *unavailability* of
  fresh knowledge, never a false ACCEPT. In-partition revocations are not
  observable before recovery (INV12) — Atlas states this instead of
  pretending otherwise.
- **A2** gets exactly S until exactly T, or until revocation is observable —
  nothing else. Within-window replay of a valid record is **not prevented**
  (FM8, carried openly; see LIMITATIONS §3).
- **A3** gets rejections. All 18 adversarial constructions are refused; the
  fuzzer has never produced a silent acceptance.
- **A4** gets refused (C8) or delivers honestly.
- **A5** can verify records (public by design) and read state; cannot mint or
  revoke without the token (when configured). Availability of the API itself
  is ordinary service hardening (TLS/rate limits: LIMITATIONS §5).
- **A6** is game over for new issuance — as with any signature scheme. The
  damage is bounded by expiry (records are short-lived by default) and by
  rotation: verifiers pin keys by `kid` in the trust bundle, so a rotated
  bundle retires the compromised key. Automated rotation is not shipped yet
  (LIMITATIONS §4).

## 6. Verification of this document

Every "proven by" above is runnable today:

```bash
make ci                          # gates: conformance, adversarial vectors, property tests
go test ./internal/verify -run x -fuzz FuzzVerify -fuzztime 30s
bash examples/unforgettable.sh   # the offline + revocation + tamper story, live
```

<!-- checkpoint: chore(scripts): refine Docker orchestration config (#132) -->
