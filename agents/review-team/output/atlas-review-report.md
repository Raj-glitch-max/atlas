# ATLAS REVIEW — SYNTHESIS REPORT

*Produced by the 13-agent review roster (`agents/review-team/`), executed against the
real codebase at `github.com/Raj-glitch-max/atlas`, commit state of 2026-07-14.
Phase order per `13_meta_orchestrator.md`: Foundation → Security → Ops&Product →
Polish → Synthesis. Every stack item traces to a named agent finding in
`atlas-agent-raw-findings/`. Arbitration rules R1–R12 applied; conflicts escalated,
not averaged (`atlas-open-questions.md`, `atlas-conflict-log.md`).*

---

## EXECUTIVE SUMMARY

**WHAT IS TRUE.**
The cryptographic core is sound and, unusually, its soundness is *structural* rather than
asserted. `alg` is pinned to ES256 and validated before any key is applied; `typ` is pinned to
`atlas-record+jws`; `kid` must resolve to operator-held material; the curve is enforced at
trust-material construction; only compact single-signature JWS is accepted. The JOSE attack
catalogue (`alg:none`, HS256 confusion, ES384 substitution, `typ`/`kid` tampering, cross-record
transplants) is covered by 18 committed adversarial vectors and an 821k-exec fuzzer, and the
verifier's routing is order-independent and property-tested. Fail-closed is enforced by an
import lint: `internal/verify` cannot perform I/O, so "offline" is a compile-time property, not a
promise. Judgment is *legible*: `ENGINEERING_DECISION_RECORD.md` (AD-001…017+) and 13 dated
journal entries record the hard calls with alternatives and trade-offs; `THREAT_MODEL.md` is a
real RFC-6819-shaped document mapping each claim to a test; `LIMITATIONS.md` publishes the
weaknesses honestly.

**WHAT IS BROKEN.**
Nothing at Tier 0 or Tier 1. No confirmed exploit. No path was found where the verifier accepts
what it must reject. For a cryptographic verification project this empty column is the
significant, rare result — and it is the single most important sentence in this report.

**WHAT IS UNPROVEN** (this is where Atlas's real risk lives, and naming it is the point):
- **The availability cost of fail-closed is unmeasured.** `R` bounds staleness and fails closed
  beyond it — correct — but the actual availability number, and the *link-level* zero-egress
  proof (two real SPIRE domains, severed cable, tcpdump), are authored in `atlas-lab/` and **not
  run** (`LIMITATIONS.md` §8). This is the project's genuine technical bottleneck.
- **Cross-implementation agreement is claimed but not differentially enforced in CI.** Three SDKs
  (Go/Python/TS) and 28 language-neutral vectors exist, but no CI job runs all three against the
  vectors and asserts identical verdicts, and duplicate-JSON-key handling is unspecified (Go is
  last-wins) — a real divergence risk the moment a second implementation reads first-wins.
- **Demand is a hypothesis.** Zero users. Every downstream adoption, DX, and sequencing claim is
  unvalidated. (Owned weakness — filed as its *consequence*, per R11, not as a bare fact.)

**THE ONE THING.**
Add a `LICENSE` (Apache-2.0). ~1 hour. It is legally absent today — the code is not adoptable
and no CNCF conversation can begin without it. Owner: 12. *(The highest-leverage **technical**
evidence action is separate and much larger — the E6/E7 two-domain substrate run, S7 — which
converts the central offline claim from "by construction" to "demonstrated." Do the license this
hour; sequence the substrate run as the incubation critical path.)*

**THE HONEST VERDICT.**
Atlas would be **admissible to CNCF Sandbox on technical merit and is blocked only on paperwork**
— a `LICENSE` (1h) and four cheap governance files (~1 day) — which is the most fixable failure
mode a project can have; **Incubation is blocked on adopters (zero) and on the unrun two-domain
substrate proof.** It demonstrates the engineering judgment it claims: the ADRs, the falsifiable
threat model, the published limitations, and fail-closed-by-construction are Staff-track signals;
the remaining gap is *influence beyond self* (bus factor 1) and *evidence* (no user, no
substrate run) — not code quality.

---

## THE ARBITRATED PRIORITY STACK

Machine-readable form in `atlas-priority-stack.csv`. Narrative below. Tier 0 is reserved for
confirmed exploits (R2); it is empty, and severity is not inflated to fill it.

### TIER 0 — CONFIRMED EXPLOITS
**Empty.** Red Team (05) produced zero CONFIRMED breaks. Every envelope attack in the playbook
was REFUTED against the real code (see `05_red_team.md`). This is a result, not an omission.

### TIER 1 — SILENT-ACCEPT CORRECTNESS BUGS / MANDATORY MISSING CONTROLS
**Empty.** No path found where the verifier returns Accept for input it should reject. The
fail-closed posture (absent material → inconclusive-reject; stale snapshot → inconclusive-reject;
Degenerate provider default → inconclusive-reject) holds by construction and under property
tests P1–P6 (`tests/conformance/properties_test.go`). Verified, stated plainly.

### TIER 2 — CHEAP UNBLOCKERS (rank here *only* because the value/hour ratio is extreme)
- **S1 · LICENSE absent (P0).** `owner 12 · 1h.` No `LICENSE`/`COPYING`/`NOTICE` is git-tracked.
  Legally un-adoptable; hard CNCF Sandbox blocker. *(EM-1)*
- **S2 · Repo governance files absent.** `owner 12 · 3h.` No `CODE_OF_CONDUCT.md`, `MAINTAINERS.md`,
  `CODEOWNERS`, or repo-level `GOVERNANCE.md`. Each declares intent-to-scale even with one name.
  *(EM-2, EM-6)*
- **S3 · SECURITY.md is 5 lines.** `owner 04/12 · 1h.` Present but has no disclosure SLA, no
  supported-versions table, no CVE process — a researcher who finds a bug has a channel but no
  terms. *(EM-3, SEC-11)*
- **S4 · Doc-claim contradictions.** `owner 11 · 2h.` The review context sheet
  (`00_atlas_context.md` §6) lists "Delegation chains ✓", which the repo contradicts (single-hop
  by design — README, `LIMITATIONS.md` §1); `LIMITATIONS.md` §10 undercounts SDKs (says one
  Python; repo ships Go+Python+TS); benchmarks are published without CPU/percentile/chain-depth.
  Reconcile so the top claim surface matches the honest lower one. *(PE-6, TW-2, EM-10)*

### TIER 3 — LEGIBILITY OF JUDGMENT
- **S5 · No normative *wire-format* spec.** `owner 11 · 1–2d (07 co-owns).` The logical contracts
  are normative (`INTERFACE_SPECIFICATION.md`, "Canonical, binding") and the vectors are
  machine-normative, but the byte layout is explicitly deferred ("the normative wire format is a
  mechanism-RFC concern", `VECTORS.md`). A reimplementer cannot build from prose alone. *(TW-1)*
- **S6 · No protocol version field, no `crit` mechanism.** `owner 01 · 1d (03 co-owns).` The
  record carries no version; forward-compat is "tolerate unknown fields, append-only", which is
  safe *today* (no critical extensions exist) but has no way to introduce a future
  mandatory-to-understand field without a flag day. Ship the version field now (cheap,
  irreversible-to-omit); design the `crit` story before the first security-relevant extension.
  *(PE-2, PE-3, CRYPTO-4)*
- *(ADRs: **not a finding** — they exist and are strong. Credited, not backfilled. This refutes
  the roster's predicted "judgment is invisible" cluster.)*

### TIER 4 — EVIDENCE
- **S7 · Run the two-domain SPIRE substrate proof (E6/E7).** `owner 06/02 · weeks.` The
  zero-egress link-level demonstration is authored but unrun (`LIMITATIONS.md` §8). This is the
  highest-leverage technical work in the repo and the incubation critical path. *(SRE-1, ARCH-2)*
- **S8 · Three beachhead interviews / one real user.** `owner 10 · 1–2wk.` Resolves OQ-1 and
  un-defers most of the roadmap. Until it lands, every product claim is a hypothesis. *(PM-1,
  UXR-1, EM-9)*

### TIER 5 — HARDENING & OPS
- **S9 · Differential-in-CI + duplicate-key vector.** `owner 07 · 2–3d.` Run all three SDKs
  against the 28 vectors in CI and assert identical verdicts; add a conformance vector + a spec
  rule pinning duplicate-JSON-key handling (reject, or a stated winner). *(TEST-1, TEST-2,
  CRYPTO-9)*
- **S10 · Signed releases + SLSA provenance + govulncheck.** `owner 06/03 · 2–3d.` An unsigned
  release of a cryptographic verification library is a credibility contradiction. *(CRYPTO-12,
  EM-7)*
- **S11 · Server hardening.** `owner 06 · deferred.` TLS/rate-limit/store are dev-grade and
  acknowledged (`LIMITATIONS.md` §5); the seams exist, the hardened deployment does not. *(SRE-4)*

### TIER 6 — ERGONOMICS
- **S12 · `atlas explain` + error taxonomy + first-run DX.** `owner 08 · deferred.` The trace is
  rich; the human-facing explanation of a Reject/Inconclusive is not yet a product surface.
  *(PD-1, UXR-2)*

### TIER 7 — DEFERRED (10's cut list; each carries an un-defer trigger)
- **Multi-hop delegation (A→B→C).** Explicitly deferred (DEFERRED D3–D4). *Un-defer:* a named user
  needs re-delegation and single-hop composition at the integrator is insufficient.
- **8-framework integration matrix, marketplace, Terraform/Helm.** *Un-defer:* per S8 — a real
  deployment/user exists to provision or integrate.
- **Proof-of-possession / holder binding (replay fix, `LIMITATIONS.md` §3).** *Un-defer:*
  within-window replay becomes a live threat for an actual user.

---

## FINDING CLUSTERS (deduplicated; each reported once, all contributors named)

| Cluster | Contributing agents | The one underlying issue | Disposition |
|---|---|---|---|
| **Attenuation monotonicity** | 03, 04(T4), 05(B), 07 | *Is child-⊆-parent mechanically proven?* | **Re-scoped, not broken.** Atlas is single-hop *by design* (README, `LIMITATIONS.md` §1). Single-hop strict-subset is enforced at issuance (`isProperSupersetOf`, strict) and tested. Multi-hop attenuation is deferred, not claimed by the repo. The finding survives only as **S4** (the *review context sheet* overstates "chains") — the artifact is honest. |
| **Offline vs revocation** | 04, 06, 10, 11 | *Verification availability ≤ snapshot availability, unmeasured.* | **UNPROVEN → S7 + OQ-4.** Bound is explicit and fails closed; the *number* is missing and the link-level proof is unrun. |
| **Judgment is invisible** | 11, 12, 01 | *Hard calls made well, recorded nowhere.* | **REFUTED.** ADRs (AD-001…017+) + 13 journal entries + normative interface spec exist. Credited. Residual is only the *wire-format* spec (S5). |
| **No user, therefore no evidence** | 02, 09, 10, 12 | *Every downstream claim is a hypothesis.* | **UNPROVEN → S8 + OQ-1.** Owned; filed as consequence per R11. |
| **Errors leak / errors confuse** | 04, 08, 09, 11 | *One error surface serves a trust boundary and a debugging session.* | **Tier 6 → S12.** The verdict trace is honest and non-leaking to the holder (records are public by design); the *human explanation* is the gap, not an oracle. |

---

## BROKEN / UNPROVEN / ABSENT (the three-column truth, strictly separated)

**BROKEN** (wrong behavior shipping today):
- *Nothing.* No confirmed exploit; no silent-accept path.

**UNPROVEN** (built, plausibly correct, not yet demonstrated — different fix from "broken"):
- Availability cost of fail-closed; the two-domain zero-egress substrate proof (S7).
- Cross-SDK agreement (no differential-in-CI) and duplicate-key determinism (S9).
- Product demand / beachhead (S8).

**ABSENT** (not built; a decision to make, not a bug to fix):
- LICENSE (S1); repo governance files (S2); SECURITY.md disclosure terms (S3).
- Normative wire-format spec (S5); protocol version field + `crit` mechanism (S6).
- Signed releases / SLSA provenance (S10); differential-in-CI (S9); production-grade server (S11).
- Multi-hop delegation, framework matrix, PoP replay binding (Tier 7 — deferred by design).

---

## OPEN QUESTIONS (escalated, not resolved — full text in `atlas-open-questions.md`)

- **OQ-1** Product or contribution-to-a-standard (UCAN/Biscuit/MCP)? → three beachhead interviews.
- **OQ-2** Does the security self-assessment publish the known gaps? → 04 rules per-item
  (architectural limitations publish; specific unpatched exploits embargo).
- **OQ-3** Ship the wire-format spec before the beachhead? → ship `v0.1-draft`, marked unstable.
- **OQ-4** How much availability does fail-closed actually cost? → *measure it* (S7).

Three of four resolve to "go get a number or a user," not "make a decision." That is the
review's most important structural finding.
