# OMEGA-04 — The discriminating-observation invariant (trust from first principles)

**Status:** Discovery record (ADR-class). Additive; no code, plan, or frozen
doc changed. Terminal entry of the OMEGA foundational line — it sits beneath
OMEGA-01/02/03 and shows they are instances of it. Presented for a founder
read; nothing built or applied.
**Origin:** founder research prompt (2026-07-06) — invent trust from first
principles; find the minimum information for a rational observer to increase
confidence in another's claim; prove it, destroy it, check whether all trust
systems realize it; don't stop until a more fundamental invariant is found or
shown to probably not exist. Worked by `trust-information-foundations`
(The Epistemologist); full climb log there.
**Doctrine:** no live literature check was possible; novelty is capped **Low**
and every literature attribution is recall-based. The value here is
unification and a usable test, not a novelty claim.

---

## 1. The invariant

**The Discriminating-Observation Invariant.** A rational observer A can
increase confidence in a claim C **iff** it obtains an observation E that is:

1. **Grounded** — access to E terminates, possibly via a chain, in A's own
   perception or computation;
2. **Dependent** — `I(C; E) > 0`: E's distribution differs depending on
   whether C is true; and
3. **Adversarially robust** — the dependence survives optimization by any
   adversary who benefits from A's error (a coupling the adversary controls
   does not count).

Justified confidence change is bounded by `I(C; E)`. Absent such an E, no
protocol, authority, or reasoning can rationally raise confidence — trust is
impossible.

The literal answer to "minimum information for A to trust B's claim": **not
the assertion** (an assertion uncoupled to truth carries zero — cheap talk),
but **a grounded, adversarially-robust observation whose likelihood differs
under the claim's truth vs. falsity.** For a claim mediated by B, that means
the assertion *plus grounded evidence that B's report is counterfactually
coupled to the truth* (B's reliability, or B's cost of lying).

## 2. It is bedrock (why no more fundamental invariant exists)

Clause 2 is not a fact about trust systems; it is forced by the axioms of
rational belief. **Cox's theorem**: any consistent, complete belief calculus
is isomorphic to probability. Probability gives "posterior > prior iff
likelihood ratio > 1 iff C and E are dependent." To go beneath the dependence
requirement one would have to abandon the definition of rational belief — i.e.
leave the question. So the search terminates here: **a more fundamental
invariant probably does not exist**, because this one is a theorem about
rationality itself. Clauses 1 (groundedness) and 3 (adversarial robustness)
are not deeper — they are the two riders that make the core usable by a real,
non-omniscient, adversarially-situated observer.

## 3. It survived every attack (each sharpened it)

Priors (relocate dependence to history, not an increase); deductive claims
(the proof is the discriminating observation — logic is the channel,
unifying verification and testimony); infinite regress (terminates at
first-person observation — all trust grounds out there); Byzantine forgery
(forces clause 3); zero-knowledge (the sharpest realization — isolates the
one C-coupled bit). Full derivation in the specialist's climb log.

## 4. The reframe: trust and verification are one continuum

There is no fundamental line between trust and verification. **Verification**
= A runs the discriminating channel itself (Git rehash, ZK check,
replication) — minimal trust. **Trust** = A relies on a coupling it did not
run but has grounded, robust evidence for (authority tokens, reputation) —
more trust. One continuum, parameterized by how much of the channel A runs
vs. relies on. *Trust is outsourced verification, rational exactly when the
outsourced coupling is grounded and adversarially robust.*

## 5. Every examined trust system is a realization

PKI/TLS, OAuth, Git, blockchains/PoW, transparency logs, capability systems,
Byzantine consensus, reputation, scientific peer review — all nine are
realizations (table in the specialist file), differing only along three axes:

- **Channel**: logical/computational, physical, or incentive-grounded.
- **Robustness mechanism**: cryptographic hardness, physical cost, economic
  stake, or independence assumptions.
- **Temporal placement**: ex ante (checked now) or ex post (detected later).

The invariant also *explains* known facts: cheap talk is worthless (`I=0`);
peer review is weaker than replication (review is reputational, replication
is the dependence channel); ZK is literally minimum-information; correlated
channels (a CA monoculture) overstate `I(C;E)` and breed overconfidence.

Honest scope limit (Empiricist): nine were checked and each fit; "every
existing trust system" is a stronger universal than an enumeration can prove.
The claim is "every system examined is a realization, and the invariant is
forced by rationality, so a non-realization would have to increase confidence
with `I(C;E)=0` — which violates Bayes."

## 6. The OMEGA arc converges (this is the through-line)

The four findings are one principle seen at four depths:

- **OMEGA-02**: freshness composes by `min` and is conserved — you cannot
  compose *more* freshness than the stalest verifiable edge.
- **OMEGA-03**: an uncheckable assertion carries zero epistemic value; only
  cost-of-defection gives it weight, bounded by `P(caught)×penalty`.
- **OMEGA-04**: rational confidence is bounded by `I(C;E)` through a grounded,
  adversarially-robust channel.

All three are the same conservation law: **you cannot get more trust out than
information about C went in, and information requires an adversarially-robust
dependence you can ground.** Freshness, cost-of-defection, and mutual
information are three currencies denominating the same conserved quantity.
The OMEGA line reaches bedrock here.

## 7. Classification and honest limits

- **Classification:** synthesis of Bayesian confirmation theory, Shannon
  mutual information / data-processing, Cox's theorem, and the epistemology of
  testimony. **Not novel** — these are established results; the contribution
  is the unification (one invariant, three axes) and the usable three-clause
  test. Novelty **Low**, no live check.
- **The usable residue (regardless of novelty):** a three-question audit for
  any trust mechanism — *What is the discriminating observation? Does the
  adversary control it? Does my access to it ground out in something I
  observe myself?* — plus the warning that the magnitude bound is broken in
  practice by counting correlated channels as independent.
- **Where the climb stopped:** at Cox's theorem (bedrock — going deeper means
  redefining rationality); and at the novelty question (needs a live check).

## 8. Relation to Atlas (no scope change)

None proposed. The bearing is explanatory: Atlas's verifier is a
discriminating-observation machine (signature = crypto channel, ex ante;
revocation freshness = the conserved currency of OMEGA-02); its fail-closed
default is the correct response to `I(C;E)=0` (no grounded observation →
no rational acceptance). This is a lens on why the kernel is shaped as it is,
not a request to change it.
