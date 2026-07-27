# ATLAS REVIEW — CONFLICT LOG

Every arbitration made during synthesis, the rule (R1–R12) that produced it, and the losing
position preserved verbatim-in-substance. A future maintainer who disagrees with a ruling can
find *which rule produced it* and *what the losing argument was*, and overturn it deliberately
rather than rediscover it by accident. The conflict log is not optional.

---

## C-1 — "Attenuation monotonicity is unproven / a bypass" vs the artifact's actual scope

- **Claimants (widen-risk):** 04 (T4 scope widening, ❌), 05 (predicted composition-widening
  break, "high confidence"), 07 (crown-jewel property test demanded).
- **Losing position (preserved):** *"Nothing mechanically proves a child capability cannot exceed
  its parent across a delegation chain; treat capability attenuation as unproven until a shrinking
  property test exists."* This is the correct instinct **for a multi-hop system.**
- **Ruling:** The premise does not hold against this artifact. Atlas is **single-hop by design**
  (README: "a narrow primitive: single-hop"; `LIMITATIONS.md` §1: "There is no A→B→C chain — by
  design, not omission"). The single-hop strict-subset property *is* enforced and tested
  (`isProperSupersetOf` requires a strict proper subset; issuance refuses `OverScope`). There is
  no multi-hop composition operator to widen, so the composition-widening exploit class has no
  surface.
- **Rule applied:** **R2** (a confirmed exploit outranks everything — but 05 produced *no*
  confirmed exploit here; the prediction was REFUTED against real code) and **R6** (a scope
  exclusion — multi-hop — challenged once, then honored). The finding survives **only** re-scoped
  as **S4**: the *review context sheet* `00_atlas_context.md` §6 overstates "Delegation chains ✓",
  contradicting the honest repo. The bug is in the context sheet, not the code.

## C-2 — Crypto (03) "attenuation fields must be `crit`" vs Principal (01) forward-compat

- **Positions:** 03 wants any security-relevant extension field marked `crit` so an old verifier
  rejects rather than silently ignores it; 01 notes `crit` breaks forward-compat (old verifiers
  reject new records).
- **Ruling:** Inside the envelope, **03 wins (R1).** An old verifier that silently ignores a
  restriction is a bypass; rejecting is what fail-closed *means*. But the ruling is **latent, not
  active**: Atlas today defines *no* critical extension fields (the schema is fixed; unknown
  fields are tolerated append-only, which is safe precisely because none of them are
  security-load-bearing yet). So the finding is **S6** (build the `crit` mechanism *before* the
  first security-relevant extension), not a present-tense defect.
- **Rule applied:** **R1** (crypto final inside the envelope), scoped by present reality.

## C-3 — "Judgment is invisible / backfill the ADRs" vs the repository

- **Claimants:** 12 (red flag #4, "single largest unforced error"), 11 (ADR backfill as
  highest-leverage deliverable), 01 (undocumented API contracts).
- **Losing position (preserved):** *"The hardest calls — JWS over COSE, ES256, fail-closed,
  snapshot revocation — exist only in the author's memory; a reviewer reads this as absence of
  judgment."*
- **Ruling:** **REFUTED by evidence.** `ENGINEERING_DECISION_RECORD.md` records AD-001…017+ with
  Alternatives / Reason / Trade-off / Complexity-cost for each, and `agents/journal/` holds 13
  dated decision entries. The judgment is extensively legible. The cluster collapses to its one
  genuine residual: the *wire-format* spec is not yet normative (S5). The predicted P0 is a
  strength, and the report credits it as such.
- **Rule applied:** domain finding treated as an **input to be checked, not a proposal to be
  weighed** (§2 of 13); checking it against the repo reversed it.

## C-4 — Availability (06) vs fail-closed non-negotiability (04)

- **Positions:** 06 — the availability coupling is real and unmeasured, operators need the number;
  04 — fail-closed is non-negotiable in the library, availability is bought via TTL/distribution.
- **Ruling:** **Not averaged, not decided — escalated (OQ-4).** Per **R3/R4**, availability may
  never be answered by weakening verification; but 06's finding is not a request to weaken it, it
  is a request to *measure the cost*. The two positions are arguing about an unknown. Resolution
  is the measurement (S7), not a ruling.
- **Rule applied:** **R3, R4**, then escalation (no rule decides an unmeasured quantity).

## C-5 — PM (10) "revocation is a v2 break-glass feature, deprioritize" vs Security (04)

- **Positions:** 10 — short-TTL solves 90% of cases; revocation is over-weighted as a v1 pillar;
  05 — revocation freshness is a stated guarantee.
- **Ruling:** **Both ship; sequencing differs (R5).** 10 may cut/defer the *feature emphasis*;
  10 may **not** cut the fail-closed-on-stale *control* (04 marks it mandatory, and it is already
  built and tested — C7). Short-TTL becomes the documented default path; revocation is the
  break-glass mechanism. No conflict remains once "feature" and "control" are separated.
- **Rule applied:** **R5** (PM cuts features, not mandatory controls).

## C-6 — Career framing (12) vs the technical/product findings

- **Ruling:** 12's maturity read (Staff-track signals: ADRs, threat model, honest limitations;
  gaps: bus factor, no adopters) is **reported but quarantined**: it altered no technical or
  product finding. The empty Tier 0/1 and the S-item ordering were determined by 03/04/05/07/01
  on the merits, independent of how any of it "looks on a résumé."
- **Rule applied:** **R10** (career framing may never alter a technical or product finding).

## C-7 — "No users" filed as a bare finding

- **Ruling:** Rejected as a standalone finding per **R11**. Re-filed as its *consequence* (S8 /
  OQ-1): every downstream adoption, DX, and sequencing claim is an unvalidated hypothesis. The
  fact is owned (context §10); the consequence is the finding.
- **Rule applied:** **R11**.
