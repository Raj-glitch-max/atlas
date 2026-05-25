# OMEGA-03 — The currency of trust (composing from unverifiable assertions)

**Status:** Discovery record (ADR-class). Additive; changes no code, no plan,
no frozen doc. Extends the OMEGA series (`OMEGA-01`, `OMEGA-02`) and answers
the Red-Team open question logged in
`agents/journal/2026-07-06-freshness-composition-falsification.md`. Presented
for a founder decision; nothing here is built or applied.
**Origin:** founder research prompt (2026-07-06) — "is there an attestation
type that can structurally only ASSERT its freshness, and can trust compose
from unverifiable assertions using something other than freshness as the
currency." Worked by the `assertion-trust-composition` specialist (The
Assayer); the full climb log lives there.
**Doctrine:** confidence needs cited, fetched evidence — no live literature
check was possible this session, so novelty is capped Low-to-Medium and every
literature claim is marked recall-based.

---

## 1. The question partly dissolved — and that is the finding

The prompt presupposes freshness is the currency of trust composition. The
central discovery is that it is not the only one, and that the presupposition
is what makes the question look like an unclimbed mountain. Once dropped, the
base of the mountain turns out to be well-mapped terrain in economics and
accountability-systems research — with one possibly-new ridge line.

## 2. The premise is real and load-bearing, not hypothetical

An attestation type whose *present-state* freshness is structurally
unverifiable is not an edge case. The most fundamental freshness claim in any
cryptographic system — **"the signing key has not been compromised as of
now"** — is exactly this: compromise is a traceless non-event; non-compromise
can never be verifiably observed, only asserted (or detected via misuse after
the fact). Every signature's meaning rests on it. In Atlas this is **FM5**,
and it is not peripheral — it is the unverifiable freshness assertion every
other attestation silently stands on. Atlas's **S4** (revocation across a
partition) is the same phenomenon, time-boxed to the partition.

Consequence: a system that excluded all unverifiable freshness would exclude
key-non-compromise — i.e., exclude *all cryptographic trust*. So the
ecosystem already, universally, composes trust on top of an unverifiable
freshness assertion. "Can we?" is answered by existence. The real question is
*what makes it legitimate, and what is the currency.*

## 3. Information theory forces the currency to be cost-of-defection

An assertion you cannot check carries **zero epistemic value on its own**: a
liar who always asserts "X" produces the observation "asserts X" identically
whether or not X holds — likelihood ratio 1, no Bayesian update. The
assertion's value can only be *manufactured* by making "asserts X" correlate
with "X true," which requires a **cost imposed on lying**. This is the
**cheap-talk / costly-signaling** result (Crawford–Sobel 1982; Spence 1973)
imported into trust composition: freshness is an epistemic currency and is the
wrong tool for uncheckable claims *by construction*; the only currency that
can give them weight is cost-of-defection.

## 4. The non-freshness currencies (all reduce to cost-of-defection)

| Currency | Mechanism | Grounds out? |
|---|---|---|
| Retrospective attributable accountability | violation becomes evident + attributable later (Certificate Transparency; PeerReview) | yes — on catastrophe-if-caught |
| Economic stake / slashing | bonded value lost on lying (crypto-economic security) | yes — bond is itself verifiable |
| Reputation / repeated game | future value depends on not being caught (folk theorem) | yes — on discounted future |
| Independent redundancy | K-of-N failure-independent witnesses (wrong only if all lie) | **no — recurses to unverifiable independence/collusion** |

## 5. The boundary does not vanish — it relocates and keeps a ceiling

- **Irreversibility limit:** accountability is deferred blame; for an
  irreversible high-value action, catching the liar later does not undo the
  harm. Admit-with-accountability works only where harm is reversible/bounded
  or the ex-ante threat deters.
- **The ceiling (quantifiable):** a rational asserter lies iff
  `gain(lie) > P(caught)·penalty + lost-future-value`. So trust extractable
  from an unverifiable assertion is bounded by roughly
  `P(caught) × credible-penalty` per unit of action value. **Above the
  ceiling — action value exceeds any imposable penalty — no incentive
  suffices and "verify or exclude" reasserts as permanent.**
- **Cryptography relocates, never eliminates:** ZK proves computational
  statements, not present-state world-facts (it cannot prove "my key hasn't
  leaked"); TEE attestation converts "not compromised" into
  verifiable-freshness-of-a-proxy at the cost of trusting a smaller enclave
  root; accumulators compress revocation but still need a fresh signed state.
  Each shrinks the trusted base to a minimal root; the unverifiable assertion
  survives there.

## 6. The answer, stated plainly

**"Verify or exclude" is permanent only within freshness-as-currency** (the
information-theoretic S4 bound). Trust *can* be composed from unverifiable
assertions — but on a different currency, **cost-of-defection**, which is
economic/game-theoretic rather than epistemic, and bounded by
`P(caught) × credible-penalty` per unit of action value. The right design
object is therefore **currency selection indexed to the action's
reversibility/value**: verifiable freshness (or exclude) for
irreversible/high-value actions; unverifiable-assertion-plus-sufficient-
cost-of-defection for reversible/bounded ones. Most PKI/authz is binary
valid/invalid regardless of stakes; this says it should not be.

## 7. Honest classification and where I ran out of angles

- **Classification:** the components are known (cheap talk / costly
  signaling; crypto-economic slashing; repeated-game reputation; Certificate
  Transparency and accountability systems; risk-based/step-up auth). The
  synthesis — freshness as one currency among several, the info-theoretic
  forcing to cost-of-defection, and action-value-indexed currency selection
  with a penalty ceiling — is a framing I have not seen assembled exactly this
  way, but I have **not** confirmed against prior art via live search.
  Novelty: Low-to-Medium, unverified.
- **Angles tried and exhausted:** freshness-as-currency (permanent boundary);
  accountability (works, reversible only); stake (works, to the bond);
  reputation (works, repeated game); redundancy (recurses to independence);
  ZK/TEE/accumulators (relocate, don't eliminate); information theory (the
  backbone — cheap talk). **Where I stopped:** (a) quantifying the exact
  penalty ceiling for a general action needs a domain-specific value model I
  cannot produce generically; (b) confirming the synthesis's novelty needs a
  live literature check I could not run this session. Neither is a dead end in
  the theory — both are honest limits of this session.

## 8. Relation to Atlas (no scope change proposed)

This does not change V1, the kernel, or any frozen doc. Its bearing on Atlas:
FM5 is correctly marked unmitigated *within freshness-currency*; OMEGA-03
explains why the ecosystem tolerates that (accountability currency: CT +
reputational/legal catastrophe) and why Atlas's conservative "fail closed on
unverifiable revocation freshness" is the right default *for a security
primitive whose actions may be irreversible* — it sits above the penalty
ceiling by design. A **proposed, not applied** note for the founder: if a
future scope act ever admits lower-stakes delegation, currency selection
(§6) is the lever, and it belongs in a numbered RFC, not in the kernel. This
is an input to a decision, not a decision.

<!-- checkpoint: chore(test): harden Fuzz Verification core target (#158) -->
