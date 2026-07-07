# Phase A — Primitive Discovery: what is Atlas, actually?

**Status:** Research design review (ADR-class). Additive; no code, no frozen
doc, no scope change *applied*. Answers the founder's Phase A charter.
**Method:** grounded in the frozen Phase 7 Product Definition + Phase 8
engineering package + RFC-000/001/002 + the OMEGA-01..04 findings + the
tagged external research (`EXTERNAL-RESEARCH-2026-07-06.md`).
**Discipline:** the charter says "prove novelty, be brutally skeptical." I
took that literally. The honest conclusions are deflationary; I did not
inflate to match the ambition of the framing.

---

## 0. The one finding that governs all others

The Phase A charter and the tagged research imagine Atlas as a candidate
**new infrastructure primitive** — a trust platform with a query language, a
VM, a graph engine, 1M-delegation scale. Frozen Atlas is the opposite: a
**self-contained, offline-verifiable, single-hop delegation token bound to
SPIFFE workload identity, across two trust domains, with no live authority
and no central store.** These are different design quadrants:

| | Frozen Atlas | The charter's / compendium's vision |
|---|---|---|
| Topology | decentralized, self-contained token | centralized graph / store |
| Time | offline (no live authority) | online query engine |
| Unit | one record, single hop | a graph of 1M delegations |
| Verification | check one token | traverse / query a graph |
| Scope authority | frozen Product Definition | to-be-amended |

**Discovery, stated plainly: Atlas is not a new primitive. It is a specific,
well-chosen point in the existing design space of offline, public-key,
attenuable capability tokens — closest to Biscuit and UCAN — distinguished by
being SPIFFE-native and by an unusually honest treatment of the offline-
revocation-freshness limit.** Turning it into the charter's platform is a
scope amendment and a multi-year research program, not a discovery this
document can assert. Everything below is grounded in that honesty.

## 1. Primitive objects (primitive / derived / removable)

From the frozen model (RFC-002) plus OMEGA:

| Object | Status | Why |
|---|---|---|
| **Principal identity, Delegate identity** | **primitive** (external) | consumed from SPIFFE, never issued (INV10/C1). Atlas does not own identity. |
| **Delegation Record** | **primitive** (the artifact) | the one thing that exists and crosses the boundary (RFC-003 M1). |
| **Scope** | **primitive** | strict-subset attenuation (INV2); the capability content. |
| **Expiration** | **primitive** | the only freshness Atlas verifies without a channel (INV3). |
| **Trust material** | **primitive** (external, RP-local) | the discriminating channel's key (OMEGA-04). |
| **Verification verdict** | **primitive** (the operation) | the 1-node kernel of the trust-composition calculus (OMEGA-02/04). |
| **Revocation (instance-targeted)** | **primitive** | one-way, terminal (INV4). |
| **Revocation-observability / freshness** | **primitive, and the hard one** | the conserved quantity; S4-bounded (OMEGA-02/04). |
| Proof, Policy, Constraint, Trust-edge-graph, Authority-server, Context, Decision-engine | **derived or absent** | Atlas has no policy engine, no graph, no online authority. "Proof" = the signed record. "Decision" = the verdict. A delegation *graph* does not exist at single-hop. |
| Trust score, reputation, capability-algebra beyond strict subset | **absent / out of scope** | not forced by any frozen requirement (AP11). |

**Irreducible set:** {external identity} + {Delegation Record carrying (principal, delegate, scope, expiration, instance-id, integrity)} + {locally-held trust material} + {verdict operation} + {instance-targeted revocation with a freshness bound}. Everything else in the compendium (graph, policy language, scores) is either derived, external, or a superset-only object.

## 2. Mathematical model (choose only what fits; reject the fashionable)

Established in OMEGA-02/04, not re-derived:

- **The verdict** lives in a bounded meet-semilattice `V = {Reject ⊑ Inconclusive ⊑ Accept}`; composition (if it ever existed, i.e. multi-hop) is the meet. Freshness composes by `min` (conserved, non-improvable for *verifiable* freshness). *Fits: lattice/semilattice algebra.*
- **The trust decision** is the discriminating-observation invariant: rational acceptance iff a grounded, adversarially-robust observation with `I(C;E) > 0`; confidence bounded by `I(C;E)` (OMEGA-04, forced by Cox's theorem). *Fits: Bayesian/Shannon.*
- **Scope attenuation** is the subset order on a permission set (a Boolean lattice). *Fits: set algebra / lattice.*

**Rejected as unnecessary (the charter invited "fashionable mathematics if
unnecessary — reject it"):** category theory (no functorial structure is
load-bearing at single-hop); linear logic / linear types (interesting for a
future capability-VM, §8, but nothing in frozen scope consumes exactly-once
typing — the record is copyable, and replay is an *unmitigated* mode by
choice, FM8); rewriting systems / operational semantics of a delegation
calculus (only relevant multi-hop, deferred). Adopting these now would be
mathematics in search of a requirement.

## 3. Execution model

Frozen Atlas's execution is **proof verification of a single self-contained
record**, not graph traversal, query planning, or constraint solving:

1. Parse the presented record (one token).
2. Verify the signature against locally-held trust material (the discriminating channel; crypto).
3. Check expiry against an injected clock within a skew tolerance.
4. Check scope integrity (covered by the signature).
5. Consult a *local* revocation-observation and judge its freshness against R.
6. Route the verdict (definitive dominates inconclusive dominates accept).

This is closest to **proof-carrying authorization (Appel & Felten)**: the
verifier checks a self-contained proof, offline, with no query to an
authority. It is explicitly **not** the compendium's model (graph traversal /
query planning / trust VM), which presupposes a store Atlas does not have.
The honest correction to the charter: Atlas is a *verifier of one proof*, not
a *query engine over a trust graph*.

## 4. Core invariants (formally-verifiable targets)

The frozen INV1–INV12 + C-INV1 are the invariant set; OMEGA sharpened them.
The ones a formal effort should target first (they are the security core):

- **No scope escalation** (INV2): a delegation's scope is a strict subset of the principal's.
- **Revocation monotonicity / terminality** (INV4): revoked never returns to valid.
- **Revocation independence** (INV5/INV6): revoking one never affects another or the underlying identity.
- **No live call on the verification path** (INV7): offline, structural (import lint).
- **Tamper-evidence** (INV8): content integrity (the conformance fuzzer showed this is *content*, not *byte*, integrity — base64 malleability).
- **Observability bound** (INV12): no in-partition freshness claim (the S4 conservation law).
- **Single-check rollback** (SO5): no silent acceptance (the conformance kit already exercises this executably).

These are the natural targets for a TLA+/Coq effort (§8, deferred). The
conformance kit (`tests/conformance/`) is the executable down-payment.

## 5. Architecture derivation

Already done and honest: RFC-003 derived six modules from the frozen
requirements, each traced, none unforced (AP11); E1–E5 implemented them. This
document adds nothing to it and proposes no new layer. Architecture is not the
open question; it emerged and is built. (This is itself a finding: Phase A
does not reopen the architecture.)

## 6. Research gap analysis (the genuinely useful part)

Honest positioning against the named systems:

| System | Solves | Cannot / does not | Relation to Atlas |
|---|---|---|---|
| **Zanzibar / OpenFGA** | centralized relationship authz at massive scale, online | offline; no self-contained token; needs the service | **different quadrant** — Atlas is not this |
| **OPA / Cedar** | policy-as-code decision engines, (mostly) online/local-eval | not a delegation *token*; policy, not attenuable bearer capability | different problem |
| **SPIFFE / SPIRE** | workload **identity** (authN) at scale | no delegation, scope, or authorization semantics (deliberately, C2) | **Atlas is a companion on top** (C1/C2) |
| **OAuth / RFC 8693** | token exchange, `act` delegation claim | online authorization server; not designed for offline cross-domain verify | Atlas composes 8693-style `act` for **offline** verify |
| **Macaroons** | offline attenuable bearer tokens with caveats | **HMAC** — verifier needs the shared root secret; weak cross-trust-domain | Atlas is the **public-key** analogue (verify with a public bundle across domains) |
| **Biscuit** | public-key, offline, attenuable capability tokens + Datalog | not SPIFFE-native; offline revocation unsolved (as for everyone) | **closest cousin** |
| **UCAN** | public-key, offline, DID-based delegation chains | DID/web-native, not workload/SPIFFE-native | **closest cousin** |
| **Proof-carrying authorization** | the academic framing Atlas instantiates | a framework, not a deployable SPIFFE artifact | Atlas's lineage |
| **CHERI / capability OSes** | hardware/OS capabilities | not a cross-domain network delegation token | orthogonal |

**Why Atlas exists (honestly):** there is a real, narrow gap — *a SPIFFE-native, public-key, offline-verifiable, attenuable delegation token that coexists with SPIRE and is honest about the offline-revocation-freshness limit.* Biscuit/UCAN occupy the adjacent space but are not SPIFFE-native; Macaroons are HMAC (not cross-domain); OAuth/Zanzibar/OPA are online. That gap is genuine but **incremental**, not a new primitive.

## 7. Novel contributions — brutally skeptical

Classifying each candidate as the charter demanded:

| Candidate | Verdict |
|---|---|
| Offline attenuable public-key delegation token | **already known** (Biscuit, UCAN) |
| SPIFFE-native delegation companion | **engineering improvement / integration** — real value, not research novelty |
| Freshness composes by `min`, conserved | **known** (weakest-link trust metrics; RFC 5280 window intersection) — OMEGA-02 |
| Uncheckable assertions need cost-of-defection | **known** (cheap talk / costly signaling) — OMEGA-03 |
| Discriminating-observation invariant | **known** (Bayesian confirmation + Shannon + Cox) — OMEGA-04 |
| Clean statement of the offline-revocation-freshness impossibility + verify-or-exclude currency for SPIFFE delegation | **combination of known work; possibly a workshop-grade systematization** — the most defensible "contribution," and it is modest |
| Conformance-kit-as-specification + property/differential fuzzing for the verifier | **engineering best-practice, well-executed** — not novel; valuable |
| Trust VM / query language / calculus / graph engine (compendium) | **not Atlas** — superset research programs, unproven, deferred (§8) |

**Verdict: Atlas contains no proven publishable primitive.** Its honest
identity is *an integration + a rigorous, honest engineering realization of a
known design point*, plus one modest systematization (the freshness-limit
statement) that is workshop-adjacent at best. Claiming more would fail the
charter's own "prove novelty" bar.

## 8. Research programs (cataloged as DEFERRED; honest novelty)

Each compendium program, with scope status. **All are outside the frozen V1
scope; none is authorized; each would require a founder scope amendment to
pursue.** Cited to `EXTERNAL-RESEARCH-2026-07-06.md`.

| Program | Novelty | Scope status |
|---|---|---|
| Multi-hop Delegation / Capability Calculus (§T.2/T.3) | known (SPKI/DCC) | deferred (S5; frozen single-hop) |
| Offline revocation mechanism (§T.7/T.8) | open engineering problem (the C4 spike) | **in the plan** (E7), gated on S1–S4 |
| Trust Query Language / DB / Indexing / Graph engine (§T.11–T.14) | known (Zanzibar/TrustQL) | superset; deferred; needs redefinition |
| Trust VM / Verifier Compiler / Capability Runtime (§T.15–T.17) | speculative | superset; deferred |
| Trust Debugger / Profiler (§T.18/T.19) | engineering | the decision trace already seeds this; deferred |
| Formal verification of INV1–12 (TLA+/Coq) | engineering + modest research | **worth doing**; deferred; conformance kit is the down-payment |
| Trust compression / accumulators (§T.7) | known crypto | feeds E7 revocation spike |
| ≥3-domain / N-domain composition (the OMEGA calculus) | known (weakest-link) + honest-freshness twist | deferred; the OMEGA-02 pre-design exists |

The single program that is *both* in-scope and genuinely open is **offline
revocation (C4)** — already the EXP-001 spike, already planned, already
blocked on your S1–S4 scope acts.

## 9. Implementation readiness

**Nothing new to implement from Phase A.** The frozen kernel (E1–E5) is built,
green, conformance-tested. Phase A did not surface a missing primitive or a
forced new package; per the charter's own rule ("nothing exists until
mathematically or operationally justified"), no package is justified by this
document. The next *code* is unchanged from the plan: E6 (substrate) → E7
(revocation spike) → E8 (V1 verdict), gated on S1–S4.

## The decision this actually surfaces (yours, not mine)

Phase A converges — for the fifth time across OMEGA-01..04 and now here — on
the same honest answer: **frozen Atlas is a correct, minimal, known-design-
point primitive; the grand vision in the compendium is a different, larger
system.** So the real fork is a governance decision only you can make:

- **Path 1 — Atlas stays the narrow primitive.** Then the compendium is
  tagged reference for later, the research programs stay deferred, and the
  highest-value next action is **resolving S1–S4 to unblock the built kernel's
  spike** — not more discovery.
- **Path 2 — Atlas becomes the trust platform** (query language / VM / graph /
  N-domain). Then this is a deliberate **amendment to the frozen Product
  Definition** (via `CONTRIBUTING.md` §4: journal → dated change note →
  `make frozen-baseline`), re-opening scope with eyes open to a multi-year
  research program that the four-test analysis (OMEGA) has repeatedly judged
  fails small-team / adoption / near-term-value tests.

I will not choose for you by writing the Path-2 architecture as if it were
decided — that would be the scope-expansion-by-document failure. This document
is the input to your decision. If you choose Path 2, say so and I will open the
amendment properly. If Path 1, the next move is S1–S4.
