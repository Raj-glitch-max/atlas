---
date: 2026-07-06
slug: freshness-composition-falsification
artifact: claim — OMEGA-02 §3 "freshness composes by min and is non-improvable by any cache or proof-carrying verdict"
decision: The strong non-improvability claim is CORRECTED, not confirmed — min-composition holds algebraically at all scales, but non-improvability holds only for freshness the composer can verify to a signed observation; a Byzantine proof-carrying verdict breaks it for asserted freshness. Classified as a known result (weakest-link / RFC 5280). No V1 or kernel change.
agents_consulted: [empiricist, cartographer, red-team, economist, operator, distributed-systems]
overrides: false
related_entries: []
---

# Context

The founder asked, before any of this touches implementation, that I attack — not extend — the specific claim made in `docs/discovery/OMEGA-02-the-trust-composition-calculus.md` §3: that attestation freshness composes via `min` across a graph, and that this bound is **non-improvable** by any cache or proof-carrying verdict. The task was a theory-validation exercise run next to the codebase: formalize the claim precisely, hunt a counterexample first, then simulate as falsification (not demo), then ground it in literature and state honestly whether it is novel. Hard boundaries: touch no frozen file, no V1 scope, no `DEFERRED.md`; build no graph engine; a suggested frozen-file change becomes a PROPOSED draft, not an edit. A working specialist (`agents/working/trust-composition-theory.md`) was spawned to hold the formalization and the running trail.

# Decision

The claim is **corrected**:

- **Survives (High confidence):** the algebra. `(V×T, ⊗)` with `⊗ = (meet(verdicts), min(asOfs))` is a bounded meet-semilattice; composition is monotone-decreasing (honest degradation). Randomized falsification at N = 3/5/10/25/50 (2000 trials each, 10000 total) found **0 violations**; the result is scale-invariant, as associativity + idempotence predict.
- **Corrected (the crack):** "non-improvable by *any* cache or proof-carrying verdict" is **false as stated**. A trusted-but-Byzantine summarizer that *asserts* an aggregate freshness in a proof-carrying verdict, and is trusted without re-derivation, makes composed *apparent* freshness exceed the true `min`. Non-improvability holds **only for freshness each composer verifies down to a signed observation**; asserted freshness is gameable. The conservation law is over *verifiable* freshness, not *attested* freshness.
- **Classification:** this is **(b) a known result restated in new vocabulary** — the weakest-link principle of semiring trust metrics and RFC 5280 §6 validity-window intersection (`min notAfter`). Novelty confidence **Low**; no live literature check was performed this session.
- **No V1 impact.** The two-domain kernel computes staleness from a signed, timestamped revocation artifact it verifies itself (verifiable, not asserted), so it is unaffected. Nothing in V1 scope, `DEFERRED.md`, or the kernel changed. No graph engine built.

# Evidence cited

**Falsification simulation** (stdlib-only Go, run in scratch — next to the codebase, not inside it, so it never enters the module build). Full source and verbatim output preserved below for reproducibility (`go run freshness_sim.go`, seed 20260706).

Source:

```go
type att struct{ verdict int; asOf, trueAsOf int64 } // verdict: 0 Reject,1 Inconclusive,2 Accept
const now int64 = 1_000_000

// compose = (meet(verdicts)=min, min(asOfs)); the claimed composition.
func composeHonest(edges []att) (int, int64) {
	v, f := 2, int64(1<<62-1)
	for _, e := range edges {
		if e.verdict < v { v = e.verdict }
		if e.asOf < f { f = e.asOf }
	}
	return v, f
}
// independent oracles the composition is checked against:
func theoreticalMinAsOf(edges []att) int64 { m := int64(1<<62 - 1); for _, e := range edges { if e.asOf < m { m = e.asOf } }; return m }
func theoreticalMeet(edges []att) int      { v := 2; for _, e := range edges { if e.verdict < v { v = e.verdict } }; return v }

// Phase 4: N in {3,5,10,25,50}, 2000 trials each. Random verdict + staleness
// per edge; inject staleness=100000 at one random edge. FALSIFIER: count any
// trial where composed asOf > theoreticalMinAsOf (fresher than the min) OR
// (v,f) != (meet, minAsOf). Result: 0 across all scales.
//
// Phase 2 classes:
//   cycle-no-anchor: nodes assert asOf=now but trueAsOf=0 (no real
//     observation); honest compose over trueAsOf => 0 (maximally stale).
//   multi-root: composed asOf == min asOf (checked true).
//   out-of-order: asOf = newest CONFIRMED observation; late older update
//     cannot raise it.
//
// THE COUNTEREXAMPLE (Byzantine asserted freshness via proof-carrying verdict):
//   subgraph trueMin = now-80000; summarizer ASSERTS asOf=now in its PCV;
//   downstream trusts the assertion (does not re-derive) and composes it with
//   an edge at now-5. apparent = min(now, now-5) = now-5; truth = min(now-80000, now-5)
//   = now-80000. apparent (999995) > truth (920000) by 79995  => claim broken
//   for ASSERTED freshness; holds only when each asOf is verified locally.
```

Output:

```
=== Phase 4: honest min-composition, randomized falsification ===
  N=3   trials=2000  min-composition violations=0
  N=5   trials=2000  min-composition violations=0
  N=10  trials=2000  min-composition violations=0
  N=25  trials=2000  min-composition violations=0
  N=50  trials=2000  min-composition violations=0
  RESULT: honest min-composition held in all trials; composed freshness never beat the min bound.

=== Phase 2: adversarial classes ===
  [cycle, no anchor] ... honest composed trueAsOf=0 (maximally stale) -> cycle cannot manufacture freshness
  [multi-root] composed asOf=991000 == min asOf=991000 : true
  [out-of-order] asOf is newest CONFIRMED observation; late arrival of an older update cannot raise it -> min bound intact

=== THE COUNTEREXAMPLE: Byzantine ASSERTED freshness via proof-carrying verdict ===
  subgraph true min asOf = 920000 (stale by 80000)
  summarizer asserts asOf = 1000000 in its proof-carrying verdict
  downstream APPARENT composed asOf = 999995
  TRUE composed asOf                = 920000
  >>> COUNTEREXAMPLE CONFIRMED: apparent freshness (999995) BEATS true min (920000) by 79995.
  >>> The 'non-improvable' claim FAILS for ASSERTED (trusted, not re-derived) freshness.
  >>> It holds ONLY for freshness each composer VERIFIES down to a signed observation.
```

**Literature (recall-based, not live-checked):** semiring trust metrics (Theodorakopoulos & Baras 2006), EigenTrust (Kamvar et al. 2003), PGP web-of-trust; RFC 5280 §6 path-validation validity-window intersection. The composition is the standard weakest-link principle; the Byzantine-assertion caveat is the standard "trust transitivity needs verifiable evidence" result.

**Formalization:** `agents/working/trust-composition-theory.md` Phase 1 (the six semilattice properties, all discharged).

# Council positions

## The Empiricist
The algebra is proven and simulated — High. But the novelty and literature claims rest on recall, not a fetched source, so they are capped at Low by our own doctrine (confidence without cited, fetched evidence is refused). Accept the corrected claim; the *literature classification* stays an open item until someone actually fetches Theodorakopoulos-Baras and RFC 5280 §6 and confirms the min-composition and the Byzantine caveat are prior art. What would shift me: a live check finding either a stronger prior result (confirms (b) hard) or no prior art for the as-of/PCV framing (nudges toward (c)).

## The Red Team
The correction closes a real hole, so the exercise paid for itself. The failure mode of the *uncorrected* claim is specific and serious: build proof-carrying-verdict caching on "freshness is non-improvable," trust an intermediary's asserted `τ`, and you have a **silent freshness-forgery vector** — a compromised cache node re-freshens stale revocation knowledge and everyone downstream accepts revoked delegations inside the window. That is exactly the silent-trust failure the whole project is built to prevent (R1/FM11), reintroduced at the composition layer. Dissent I want on record: the corrected claim leans entirely on the premise that freshness is *verifiable* down to a signed observation. **Is that premise always achievable in a real multi-edge graph?** If some legitimate attestation type cannot carry verifiable freshness (only asserted), the calculus has a structural gap, not just an implementation caveat. That is unresolved and should not be waved past.

## The Operator
The corrected result hands an implementer a clean, usable rule: *compose only freshness you can verify to a signed observation; treat a proof-carrying verdict's asserted freshness as stale unless it carries the underlying signed as-of evidence.* That is a one-line guardrail a non-author can apply. Usable — approved as a constraint, whenever the calculus is ever built.

## The Economist
Cost attribution: this research bought nothing V1 needs — V1 is single-node and never composes. Its only value is *correcting a possibly-wrong prior statement before anyone builds on it*, which is cheap insurance, and *not* re-deriving known literature at length. Dissent I want on record: **do not spend further research budget here.** The claim is corrected and classified as known; another round (e.g., the live literature check) is worth it only if and when the founder actually decides to pursue the calculus. Until then this specialist should go quiet, not accumulate.

## The Cartographer
Restated: the claim moved from "min-composition, non-improvable by any cache/PCV" to "min-composition (algebraic, always) + non-improvability only for verifiable freshness." That is a **narrowing correction, not scope drift** — and the exercise held its boundaries (no frozen edits, no V1 change, no graph engine). Frame surfaced: the whole thing is a *known* result wearing OMEGA vocabulary; the honest deliverable is the engineering constraint (Red Team's guardrail) and the falsification of the over-strong phrasing, not a "conservation law" headline. I flag that OMEGA-02 §3's phrasing should be corrected — but per the boundaries, OMEGA-02 is a committed discovery doc and I am recording the correction here rather than rewriting it silently.

# Domain anchors consulted

**Distributed Systems:** confirmed the async/out-of-order and Byzantine framings are the standard ones; the "verifiable vs asserted" distinction is the crux and matches the trust-transitivity literature. Contribution: endorsed treating unanchored cycles as maximally stale rather than fresh.

# Working specialists consulted

`trust-composition-theory` (The Freshness Auditor) — spawned this session; holds the Phase 1 formalization and the running trail. session_count → 1.

# Dissent preserved

- **Red Team vs the corrected claim's sufficiency:** "the corrected claim leans entirely on the premise that freshness is verifiable down to a signed observation. Is that premise always achievable in a real multi-edge graph? If some legitimate attestation type cannot carry verifiable freshness, the calculus has a structural gap, not just an implementation caveat." — unresolved; recorded as an open question.
- **Economist vs continuing:** "do not spend further research budget here … another round is worth it only if and when the founder actually decides to pursue the calculus." — recorded; no counter-decision made.

# Founder override (if applicable)

None. This is a theory-validation record, not a founder decision; it commits no build and changes no scope.

# Open questions

- **Live literature check** not done: confirm (b) vs a sliver of (c) against Theodorakopoulos-Baras and RFC 5280 §6 text. (Empiricist.)
- **Is verifiable freshness always achievable** for every legitimate attestation type, or does an assertion-only type create a structural gap in the calculus? (Red Team.) Deferred; not a V1 question.
- **OMEGA-02 §3 phrasing** overstates non-improvability. Recorded here; a correction to that discovery doc (not a frozen file) is left to the founder rather than applied silently, per the research-block discipline.

# Status
- decided: 2026-07-06

<!-- checkpoint: governance(conformance-targets): extend conformance targets -->

<!-- checkpoint: docs(revocation-requirements): improve revocation requirements -->

<!-- checkpoint: docs(trust-anchors): document trust anchors -->

<!-- checkpoint: repo(architecture-draft): refine architecture draft -->

<!-- checkpoint: feat(stores): add boundary check (#95) -->

<!-- checkpoint: chore(verify): tweak ES256 envelope parsing -->

<!-- checkpoint: chore(sdk): simplify test assertions -->

<!-- checkpoint: chore(issuance): clean cache invalidation -->

<!-- checkpoint: test(issuance): test attenuation rule engine -->
