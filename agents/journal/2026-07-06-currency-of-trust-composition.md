---
date: 2026-07-06
slug: currency-of-trust-composition
artifact: open question (Red Team, prior entry) — can trust compose from structurally-unverifiable freshness assertions, using a currency other than freshness
decision: Recorded finding (not a build): "verify or exclude" is permanent only for freshness-as-currency; trust CAN compose from unverifiable assertions on a cost-of-defection currency, bounded by P(caught)×credible-penalty per unit of action value, above which verify-or-exclude reasserts. Currency should be action-value-indexed. No V1/kernel/frozen change.
agents_consulted: [empiricist, cartographer, red-team, economist, operator, distributed-systems]
overrides: false
related_entries: [freshness-composition-falsification]
---

# Context

The founder handed a starting point, not instructions, and explicitly asked for a long deep climb: the Red-Team open question from `2026-07-06-freshness-composition-falsification.md` — is there an attestation type that can structurally only ASSERT its freshness (never a verifiable signed observation), and if so, is "verify or exclude" a permanent boundary, or can trust compose from unverifiable assertions using something other than freshness as the currency. Boundaries: no frozen file, no V1 scope, no kernel, additive only, proposed-not-applied for any frozen correction, keep a running log, don't stop at "unsolved." A working specialist (`agents/working/assertion-trust-composition.md`, The Assayer) was spawned to hold the climb.

# Decision

Recorded (this commits no build and changes no scope):

- **The premise is real and load-bearing.** The canonical unverifiable-freshness attestation is *"the signing key is not compromised as of now"* — a traceless non-event, never verifiably observable, only assertable. It is Atlas's **FM5**, and every signature depends on it. **S4** (revocation across a partition) is the same phenomenon, time-boxed. A system excluding all unverifiable freshness would exclude all cryptographic trust — so the ecosystem *already* composes trust on an unverifiable assertion. Existence is not in question; legitimacy and currency are.
- **Information theory forces the currency.** An uncheckable assertion carries *zero epistemic value* on its own (likelihood ratio 1 against a strategic liar); its value can only be manufactured by a **cost on lying** (the cheap-talk / costly-signaling result). Freshness is epistemic and therefore the wrong currency for uncheckable claims by construction; the only currency that works is **cost-of-defection**.
- **Answer:** "verify or exclude" is permanent only *within freshness-as-currency*. Trust composes from unverifiable assertions on cost-of-defection (accountability / stake / reputation / redundancy), which is economic/game-theoretic not epistemic, **bounded by ≈ P(caught) × credible-penalty per unit of action value.** Above that ceiling (irreversible/high-value actions whose lie-gain exceeds any imposable penalty), verify-or-exclude reasserts as permanent. The design lever is **currency selection indexed to action value.**
- **No Atlas change.** V1/kernel/frozen untouched. Bearing: FM5's tolerability-in-the-wild is explained by accountability currency (CT + reputational/legal catastrophe); Atlas's fail-closed default is correct precisely because a security primitive's actions typically sit *above* the ceiling.

# Evidence cited

- Information-theoretic argument (zero epistemic value of an uncheckable assertion; cost-of-lying required) — reasoning path in `assertion-trust-composition.md` Move 3; grounded in cheap-talk/costly-signaling (Crawford–Sobel 1982; Spence 1973) [recall, not live-checked].
- Cost-of-defection currencies and their ceiling — Moves 4–5; crypto-economic slashing, repeated-game reputation, Certificate Transparency / PeerReview accountability [recall].
- Crypto-relocates-not-eliminates (ZK/TEE/accumulators) — Move 6.
- FM5 as the worked example — `04_FAILURE_MODEL.md` [frozen; read only].
- Full synthesis: `docs/discovery/OMEGA-03-the-currency-of-trust.md`.

# Council positions

## The Empiricist
The info-theoretic backbone is sound *as reasoning*: against a strategic liar with no cost, an uncheckable assertion gives no likelihood ratio — that is rigorous, not hand-waving. But the literature grounding is recall, not fetched, so novelty is capped Low-to-Medium and the "cheap talk" attribution needs a live check before it's asserted as prior art. Confidence in the *conclusion* (currency must be cost-of-defection): High as an argument. Confidence in *novelty*: Low. What would shift me: a live check, or a counter-model where an uncheckable assertion carries value with no cost of lying — I bet none exists.

## The Red Team
This answers my open question, and the answer is more dangerous as a headline than as a rule. "Trust can compose from unverifiable assertions" — dropped next to an engineer without the ceiling — becomes: authorize an irreversible action on a bonded assertion where the attacker's gain exceeds the bond. The **entire safety content is the ceiling and the action-value index**; without them this finding is a foot-gun. Dissent on record: the admit-with-accountability option must NEVER appear in any downstream artifact without `P(caught)×penalty ≥ lie-gain` and the reversibility gate stated in the same breath. For Atlas specifically, a security primitive's actions are usually above the ceiling, so verify-or-exclude stays the correct default — do not let this finding erode that.

## The Operator
Currency-selection-indexed-to-action-value is usable and already exists in practice under other names (risk-based auth, step-up authentication). An implementer could apply it: classify the action's reversibility/value, then demand verifiable freshness or accept accountability-backed assertion accordingly. Approve as a *future-RFC design principle*, not a V1 feature.

## The Economist
This finding is mine by rights — cost-of-defection is mechanism design, and the info-theoretic forcing is the costly-signaling theorem. The framing is sound. Cost attribution: the research bought a reframing of why FM5 is tolerable and a principle for a future scope act; it bought **nothing V1 needs** — V1 sits above the ceiling by design and correctly fails closed. Dissent (consistent with the prior entry): do not operationalize currency selection until a scope act actually introduces lower-stakes delegation. Until then this is option-value, not work.

## The Cartographer
Restated claim: from "is there a way to compose unverifiable trust (assuming freshness currency)" to "yes, on cost-of-defection currency, bounded by a penalty ceiling above which the freshness boundary reasserts." The move that unlocked it was dropping the freshness-as-currency presupposition. Mind-change preserved: The Assayer began expecting "exclude is likely permanent" and pivoted on realizing FM5 is the universal unverifiable assertion we already accept — that pivot is in the log, not smoothed. Frame I must flag: this is **largely imported economics in security vocabulary**; calling it a "discovery" risks overselling. Partial dissent with the framing: label it an import + a possibly-new synthesis, not a new theorem. The honest deliverable is the *ceiling* and the *currency-selection principle*, both of which are real and usable regardless of novelty.

# Domain anchors consulted

**Distributed Systems:** confirmed S4 is the time-boxed special case of unverifiable-present-state freshness, and that redundancy-based composition recurses to an unverifiable independence assumption (correlated compromise / collusion) rather than grounding out. Endorsed "crypto relocates the trusted base, never eliminates it."

# Working specialists consulted

- `assertion-trust-composition` (The Assayer) — spawned this session; holds the climb log and the taxonomy. session_count → 1.
- `trust-composition-theory` (The Freshness Auditor) — referenced as the sibling handling the *verifiable* case; not re-consulted (session_count unchanged).

# Dissent preserved

- **Red Team:** "the admit-with-accountability option must NEVER appear in any downstream artifact without `P(caught)×penalty ≥ lie-gain` and the reversibility gate stated in the same breath … do not let this finding erode [Atlas's verify-or-exclude default]." — accepted as a binding caveat on the finding, not resolved away.
- **Economist:** "do not operationalize currency selection until a scope act actually introduces lower-stakes delegation." — recorded; no counter-decision.
- **Cartographer (framing):** "this is largely imported economics in security vocabulary; label it an import + synthesis, not a new theorem." — accepted; OMEGA-03 §7 carries the Low-to-Medium novelty caveat.

# Founder override (if applicable)

None. Recorded finding, not a founder decision; commits no build, changes no scope.

# Open questions

- **Novelty live-check** not done (Empiricist): confirm the synthesis against costly-signaling, crypto-economics, and CT literature.
- **Ceiling quantification** for a general action needs a domain-specific value model (The Assayer's stopping point).
- **Redundancy recursion:** does independence ever ground out, or is it unverifiable all the way down? (Distributed Systems.) Unresolved.
- **New question surfaced:** for a *security* primitive specifically, is any action ever low-stakes enough to sit below the ceiling, or does the security context always sit above it (making verify-or-exclude effectively permanent *for Atlas* even though not universal)? Leaning "usually above," but not settled.
- **Proposed-not-applied:** OMEGA-03 §8's note (currency selection belongs in a future numbered RFC if a scope act ever admits lower-stakes delegation) is an input to a founder decision, not applied.

# Status
- decided: 2026-07-06

<!-- checkpoint: chore(lab): optimize simulated agent node -->
