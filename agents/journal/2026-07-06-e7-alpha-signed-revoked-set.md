---
date: 2026-07-06
slug: e7-alpha-signed-revoked-set
artifact: the missing revocation realization — the system fails closed because M5 is the degenerate stub
decision: Built the signed-revoked-set realization of the M5 Provider contract (EXP-001 composition 1, the α-path), under the 2026-07-06 scope act. The system now observes revocations offline with cryptographically verifiable freshness. In-process realization + contract/property/e2e tests; the substrate-validated spike remains E6/E7 and is not claimed.
agents_consulted: [empiricist, cartographer, red-team, economist, operator, distributed-systems]
overrides: false
related_entries: [c4-spike-scope-act, omega-impact-conformance, discriminating-observation-invariant]
---

# Context

With the S1–S4 scope act resolved (`c4-spike-scope-act`), the revocation
composition space is defined and composition 1 (signed periodic artifact,
S2-admissible) is admissible. Until now the only M5 realization was the
degenerate always-`Indeterminate` stub (AD-007), so every verification failed
closed on revocation — correct but useless for observing real revocations.
Under the founder's total-authority grant of 2026-07-06, this builds the real
α-path realization.

# Decision

Built `internal/revstatus/statuslist.go`: a **signed-revoked-set** realization
of the M5 Provider contract.

- **Domain A (`Publisher`)** signs a timestamped snapshot of the revoked
  instance set (ECDSA P-256 over a canonical digest of listID + as-of + sorted
  revoked IDs).
- **Relying party (`SignedSetProvider`)** ingests snapshots, **verifies the
  signature** against domain A's public key it holds out of band, adopts only
  strictly-newer snapshots (monotone freshness), and answers `StatusOf` from
  the held snapshot: present → `ObservablyRevoked(asOf)`, absent from a
  complete verified census → `NotObservedRevoked(asOf)`, no snapshot →
  `Indeterminate`.
- The **as-of is the snapshot's signed timestamp**, so freshness is
  **verifiable, not asserted** — precisely OMEGA-04's honest kind, and the
  reason a Byzantine distributor cannot forge freshness (it lacks the key).
  This closes the whole OMEGA theory arc into the implementation.
- The **propagation channel** (register → published snapshot → RP) is
  composition-root wiring that reads `revorigin.View()`; `revstatus` never
  imports `revorigin` (dependency rule R5).

Wired end to end (`tests/acceptance/revocation_realization_test.go`): issue →
verify Accept → revoke → propagate signed snapshot → verify
Reject(`RevokedObservable`) → terminal on re-presentation; plus the R-governs
test (a snapshot older than R = 2 s fails closed with
`RevocationKnowledgeStale`). All at the resolved R.

**Honesty bound:** this is the composition **realized and tested in-process**.
The substrate-validated spike (two SPIRE domains, partition injection,
out-of-band egress sniffer, adversary-blinding, two-run reproducibility) is
Epic E6/E7 and is **not** performed or claimed here.

# Evidence cited

- `internal/revstatus/statuslist.go` + `statuslist_test.go` (publish/ingest/
  status, tamper-refusal across add/erase/as-of-move, wrong-key, wrong-list,
  monotone freshness, contract compliance in the ignorance state) — green
  under `-race`.
- `tests/acceptance/revocation_realization_test.go` — the end-to-end α-path
  and the R-governs freshness bound.
- The tamper tests confirm a forged snapshot can neither add nor erase a
  revocation nor move the as-of forward (each fails signature verification).

# Council positions

## The Empiricist
The realization is genuinely offline-verifiable and its freshness is
verifiable (signed timestamp), not asserted — the tamper tests demonstrate a
distributor without the key cannot forge, erase, or re-date a revocation.
Confidence in the in-process correctness: High. What I will not let stand: any
implication that the C4 spike is *settled*. It is not — no substrate ran. This
is composition 1 built; the spike is E6/E7.

## The Red Team
Two real properties the tests pin: (1) monotone freshness — an equal-or-older
snapshot is refused, so a captured old snapshot cannot roll back a revocation
(the e2e test even surfaced this by failing until timestamps advanced — good);
(2) tamper-refusal in both directions (cannot forge a revocation to deny
service, cannot erase one to keep a revoked token alive). Residual honest
limits, carried not hidden: this realization is a signed *set* — it reveals
which instances are revoked and grows with revocations; the OAuth-Status-List
bitfield (herd privacy, size) is the deferred optimization. And the propagation
channel's liveness (does the RP actually receive fresh snapshots) is exactly
where the S4 partition bound bites — handled by fail-closed-on-stale, not by
pretending delivery.

## The Distributed Systems anchor
The signed-snapshot-with-monotone-adoption is a sound CRDT-flavored design: the
revoked set is a grow-only join (a snapshot only ever adds revocations as time
advances), and adoption-by-newer-as-of is last-writer-wins on a signed clock.
Consistent with OMEGA-02's join-semilattice. The eventual-upon-recovery S4
reading is honored: a partitioned RP simply holds a stale snapshot and fails
closed once it ages past R.

## The Operator
The composition root wires it in three lines (read register → publish →
ingest); the RP holds domain A's public key like a trust bundle. Deployable
shape is clear. The bit that needs an operational runbook later: the pull/push
delivery of snapshots (deferred with E6).

## The Economist
This is the first time in many turns that real critical-path code shipped
rather than another discovery doc. The scope act unblocked it; the realization
is small, correct, and tested. Good ROI. The remaining cost (substrate + full
spike) is E6 and needs infrastructure this environment lacks.

## The Cartographer
Restate: composition 1 of the C4 spike is realized and contract-tested,
in-process, under the resolved scope parameters. It does not expand Atlas's
scope (single-hop, two-domain, SPIFFE-companion all intact) and it fills the
M5 volatile region exactly where AP12 placed it — the stable surfaces (record,
verifier) are untouched. The frozen package is untouched. This is E7's α-path,
honestly bounded short of the substrate spike.

# Dissent preserved

No seat dissented from building it. Empiricist's and Red Team's insistence —
"do not imply the spike is settled; carry the privacy/size limit and the
delivery-liveness limit honestly" — is recorded as binding framing on any
external description, not a dissent from the work.

# Founder override (if applicable)

None; built under the explicit authority grant of 2026-07-06.

# Open questions

- Full EXP-001 substrate spike (E6): two SPIRE domains, partition injection,
  egress sniffer, adversary-blinding, two-run reproducibility — not done here.
- Snapshot delivery channel (pull/push) operational runbook — deferred with E6.
- OAuth-Status-List bitfield optimization (herd privacy, size) — deferred with
  scale (C4 horizon); semantically equal to this set.
- Founder confirmation of R and the other scope-act values.

# Status
- decided: 2026-07-06

<!-- checkpoint: repo(CI-testing-gates): refine CI testing gates -->

<!-- checkpoint: docs(threat-model-scenarios): document threat model scenarios (#19) -->

<!-- checkpoint: repo(revocation-requirements): restructure revocation requirements (#8) -->

<!-- checkpoint: planning(API-path-design): clarify API path design -->
