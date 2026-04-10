# OMEGA-02 — The trust-composition calculus (extends OMEGA-01)

**Status:** Discovery record (ADR-class). Additive; changes no code, no plan,
no frozen doc. Extends and sharpens `OMEGA-01-the-larger-category.md`
(neither supersedes nor deletes it — repo doctrine preserves history).
Presented for a founder decision.
**Origin:** the second run of the Project OMEGA first-principles exercise
(2026-07-06), run strictly — including step 4 (ecosystem design) that OMEGA-01
under-delivered.
**Doctrine held:** confidence needs evidence; hypotheses stay hypotheses; the
frozen package is untouched; nothing here is built (building it now would
fail the four tests — §5).

---

## 1. The surviving primitive (sharper than OMEGA-01)

OMEGA-01 concluded "a portable attestation envelope; delegation is one claim
type." Attacked harder — *signed, chained claims are ancient (X.509, JWT,
PGP); you have named nothing new* — the envelope alone does not survive. What
survives is one level up:

> **A trust-composition calculus.** The offline, cross-domain verification of
> a directed acyclic graph of heterogeneous attestations — authority,
> provenance, posture, computation — producing not a boolean but a **verdict
> carrying an explicit freshness bound**, degrading monotonically and never
> claiming confidence it cannot have.

Why PKI path validation (the closest existing thing) is not it: path
validation is single-root-hierarchy, *online* for revocation (OCSP/CRL — the
live call the primitive forbids), boolean, and identity-only. The primitive
generalizes it to (a) a DAG across mutually-distrusting roots, (b) offline
with a per-edge staleness bound, (c) a graded verdict, (d) arbitrary claim
types. That generalization is genuinely unsolved by any standardized
composition.

## 2. Atlas is a subset — and the subset is the *verdict*, not the record

OMEGA-01 located Atlas's subset in the record (one claim schema). Sharper: the
subset is the **Verification Core's verdict algebra**.

- M3 produces `Accept | Reject | InconclusiveRejected` + a decision trace + a
  freshness policy (R, skew). That is the **1-node kernel** of the calculus:
  a graph of one edge (principal → delegate), one claim type (authority).
- The composition operator is the **meet** on the lattice
  `Reject ⊑ InconclusiveRejected ⊑ Accept`: a graph verifies Accept iff every
  load-bearing edge does; any edge Reject makes the graph Reject; otherwise
  any inconclusive edge makes the graph InconclusiveRejected. Adding an edge
  can only lower the verdict — **honest, monotone degradation**, exactly the
  posture M3 already enforces at one node.

The record (M1) is one claim schema; the verifier's verdict (M3) is the
kernel of the whole calculus. Generalizing Atlas is generalizing the graph
from one edge to N — the verdict discipline is already correct.

## 3. The conservation law (the deepest finding)

Under trust composition, **freshness is a conserved, non-improvable
quantity**:

- The graph's effective freshness is `min` over its load-bearing edges'
  as-ofs. Composition can only *lower* it.
- No mechanism — caching, mirroring, or a proof-carrying verdict (§4) — can
  manufacture freshness. A cached verdict inherits the as-of of the subgraph
  it summarizes; composing it is valid only within the freshness policy.

Atlas's `LEVEL0_1_FEASIBILITY_GATE.md` S4 result — a revocation performed
across a partition you cannot cross is information-theoretically
unobservable — is the **conservation law of the entire calculus, restricted
to one node**. The project's hardest problem (C4) is not a wart to engineer
around; it is the calculus's deepest theorem, discovered early. That is the
finding that should make an infrastructure engineer stop.

## 4. The ecosystem around the primitive (step 4 — minimal necessary set)

Introduced only where necessary; each item earns its place or is rejected.

**Necessary:**

- **Typed attestations.** One envelope, one schema per claim type (authority,
  provenance, posture, computation-result). Atlas's record is the authority
  schema. *Necessary:* heterogeneous claims are the point.
- **The composition semiring.** Verdicts under `meet` (bounded meet-
  semilattice; Accept is identity, Reject is absorbing) × freshness under
  `min`. The composed trust value of a DAG is the product. *Necessary:* it is
  the calculus itself; it is also small and standardizable.
- **Per-edge offline revocation with a bounded staleness floor.** The C4
  problem, per edge. *Necessary and hard;* the conservation law (§3) bounds
  what any solution can claim.
- **Proof-carrying verdicts.** A verifier may emit its verdict over a subgraph
  as a *new signed attestation* ("V asserts subgraph-verdict = X as-of T"),
  which a downstream verifier composes as one edge without re-walking the
  subgraph. *Necessary:* it is what makes verification itself composable and
  cacheable across mutually-distrusting agents without a central authority —
  and, by §3, it cannot cheat freshness, so it is safe.

**Rejected as unnecessary (they fail the tests in §5):** new OS concepts, new
runtime/compiler concepts, new programming models, a new language. None is
forced by the primitive; each is boil-the-ocean and unadoptable. The calculus
rides on ordinary signed bytes and ordinary verifiers — that is its strength.

## 5. Four-test verdict — building the calculus now is rejected; the kernel is kept

| Test | Full calculus now | Atlas kernel (current) |
|---|---|---|
| 1. Small team? | No — arbitrary claim types + DAG engine + proof-carrying-verdict format is boil-the-ocean | Yes — one edge, one claim type, done |
| 2. Principal respect? | Only if narrow and rigorous; "universal trust calculus" reads as a red flag | Yes — a minimal offline delegation kernel with a proven conservation bound |
| 3. OSS adoption? | No — adoption is through the SPIFFE-coexisting wedge, not a big-bang calculus | Yes — the wedge is the adoption path |
| 4. More long-term value? | The *option* on the calculus is valuable, but captured cheaply by keeping the kernel composable, not by building the graph now | Yes — kernel now, calculus when evidence forces it |

**Verdict: do not rewrite, do not delete, do not build the graph engine.**
Atlas is the correct 1-node kernel of the real primitive. The courageous act
is refusing to prematurely build the graph (AP11/DR9 forbid unforced
generality; test 1 and test 3 reject it), while pre-designing it so the
generalization traces cleanly when evidence forces it. Flipping to a
destructive rewrite to match the exercise's register would itself fail test 4
— it would destroy the traceability spine (SO8) for no long-term value.

## 6. The deferred design (so the generalization is pre-designed, not improvised)

When a founder scope act later admits a second claim type or a multi-edge
graph (post-V1; forced by evidence, per `DEFERRED.md` discipline), this is the
path — no part of it is built now:

1. **Verdict becomes a first-class value carrying its effective freshness.**
   Today M3's decision trace already records the revocation as-of and the
   policy; the only addition is to expose `verdict.effectiveAsOf =
   min(load-bearing edge as-ofs)` — for one node, the revocation as-of. Zero
   new mechanism at one node; it becomes the composition input at N nodes.
2. **Composition function** `compose(v1, v2, …) = (meet(decisions),
   min(asOfs))`, evaluated against one policy (rule R8 — one home for R).
3. **A second claim schema** (provenance is the natural first) validates that
   the record generalizes to typed attestations (OMEGA-01 Path C).
4. **Proof-carrying verdict schema** — a verdict emitted as an authority-style
   attestation, composed as one edge, bounded by its inherited as-of (§3).

Each step traces to the calculus; none is admitted without a forcing scope
act. This section is a map, not a mandate.

## 7. What was deliberately NOT done

No code, plan, or frozen doc changed, despite the exercise granting authority
to delete 90% and rewrite everything. The technically-justified conclusion
(§5) is that destruction fails the founder's own tests; the honest,
higher-value act was to find the *sharper* truth (the calculus, the verdict-
as-kernel, the freshness conservation law) and to pre-design the
generalization. Engineering resumes at `IMPLEMENTATION_MASTER_PLAN.md` (Epics
E1–E5 complete; E6+ gated on the S1–S4 founder scope acts) whenever the
founder returns to it.

<!-- checkpoint: chore(verify): tweak conformance validation -->
