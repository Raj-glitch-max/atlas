---
date: 2026-07-06
slug: omega-impact-conformance
artifact: question (founder) — does anything learned in OMEGA-01..04 change Atlas architecturally; if not prove it, if so change it
decision: The KERNEL does not change (proof per finding). The one justified engineering change is an executable conformance kit + property/differential fuzzer for the Verification Core (built), because OMEGA-04 makes verifier-exactness a precondition across implementations. The fuzzer surfaced base64 transport-malleability (benign now; strict-canonicalization proposed-not-applied). Rejected proposals recorded.
agents_consulted: [empiricist, cartographer, red-team, economist, operator]
overrides: false
related_entries: [discriminating-observation-invariant, currency-of-trust-composition, freshness-composition-falsification]
---

# Context

The founder asked the only question that matters after four theory rounds: does OMEGA-01..04 change Atlas — architecturally, mechanically — or not, and to prove whichever is true, with impact (code/tooling), not documents. Standard: would a senior engineer six months out say "this is fundamentally different"; if probably not, delete the proposal. Full authority over everything except frozen governance. Boundaries from the research block still hold (no frozen edits; proposed-not-applied for frozen-touching changes; additive; keep history).

# Decision

**The kernel does not change.** Each OMEGA finding maps to already-realized, outside-the-boundary, or scope-rejected:

- **OMEGA-02 (composition/freshness-min):** not applicable — V1 is single-edge, one discriminating channel; there is nothing to compose. No change.
- **OMEGA-03 (currency indexed to action value):** the action-value decision lives in the relying party's *grant decision*, which is explicitly outside Atlas's boundary (FM10, SYSTEM_CONTEXT). Atlas correctly does not own it; and for a high-stakes primitive above the penalty ceiling, verify-or-exclude is right — so the trichotomy verdict is *validated*, not changed. No change.
- **OMEGA-04 (discriminating-observation invariant):** Atlas is already a faithful realization — grounded (the RP verifies the signature itself), dependent (`I(C;E) > 0`, unforgeable), adversarially robust (fails closed at `I(C;E)=0`, i.e. absent/indeterminate revocation → InconclusiveRejected). No kernel change; OMEGA-04 *explains and justifies* the existing shape.

**One change is justified and was built.** OMEGA-04 makes precise that the verifier IS the discriminating channel, so if competing implementations disagree on an edge case the observation is ill-defined and an adversary selects the accepting verifier — the Frankencerts / verifier-differential failure that has burned X.509, TLS, and JWT. Atlas's conformance was *prose* ("M3 = these five checks"). Prose conformance does not survive the 5-year multi-implementation scenario the founder asked me to design for. So conformance is now **executable**:

- `tests/conformance/` — a reusable kit: a `Verifier` interface, a `Factory` any implementation supplies, an exhaustive `Scenario` corpus over the verdict space (Accept, every definitive cause, every inconclusive cause, definitive-dominates-inconclusive precedence), and `Run(t, factory)` that any implementation calls to prove itself conformant.
- Property / differential fuzzing: invariants over all inputs — content tamper never accepts differing content, garbage never accepts, Accept implies all five stages pass, routing consistency, absent material never accepts, determinism.

This hardens V1 today and turns "a reference implementation" into "a specification with a conformance kit" — the posture of infrastructure meant to have competing implementations (FIDO, WebAuthn, TLS conformance suites).

**The fuzzer earned its keep immediately:** it flagged an "accepted tamper" at a byte that turned out to be the signature segment's final base64 character — a padding-bit flip that decodes to the *same* signature. This is **base64 transport-malleability**: a valid record has multiple byte-encodings decoding to identical content. Content integrity (INV8) holds; byte identity does not. Benign for Atlas (revocation keys on the in-content instance ID, not token bytes; replay is FM8-unmitigated by scope), but recorded. The property was corrected to the honest, stronger form (never accept *differing content*), which passes.

# Evidence cited

- Per-finding mapping above (OMEGA-01..04 discovery docs + their journal entries).
- Conformance corpus + property suite: `tests/conformance/` — 10 corpus scenarios + 6 property/fuzz tests, all green under `make ci`.
- Malleability diagnosis: byte 422 of a 423-byte token is the signature segment's final base64 char (segment-final padding-bit malleability), reproduced and explained; consistent with the E2 mutation-corpus exclusion of segment-final positions (`internal/record/mutation_test.go`).

# Council positions

## The Empiricist
The "kernel does not change" claim is the honest one and it is *provable* by the per-finding mapping — each finding is already-realized / out-of-boundary / scope-rejected, not hand-waved. The conformance kit is real evidence: 10 scenarios + 6 properties, green. One honesty flag: the corpus is exhaustive over the verdict space *as I enumerated it* — completeness is asserted, and a missing verdict path would be a corpus gap; the property tests are the backstop, but "exhaustive" should read "exhaustive over the enumerated verdict space." Confidence the kernel is unchanged: High.

## The Red Team
This is the right instinct — the 5-year killer of a trust primitive is verifier differentials, and an executable oracle is the correct defense. Two things on record. First: the property fuzzer is not *true* differential testing until a second implementation exists; today it is metamorphic/property testing against invariants — valuable, but call it what it is, and the day a second implementation appears, `Run` against both is the real differential test. Second, and sharper: the base64-malleability finding is benign *only because* Atlas keys revocation on the in-content instance ID rather than token bytes. That is load-bearing and easy to lose — if any future consumer (a cache, a dedup, a replay guard) ever keys on token bytes, malleability becomes a live vulnerability. I want strict base64 canonicalization at the verifier **proposed** so the footgun is closed by construction, not left to every downstream consumer's discipline.

## The Operator
The conformance kit is usable by someone who isn't the author: a second implementer imports `conformance`, supplies a `Factory`, calls `Run`, and gets a pass/fail against the exact behavior. That is the difference between "read the spec and hope" and "run the kit." Approve. The property suite also doubles as a regression net for the V1 verifier.

## The Economist
Cost attribution: this bought two things of real value — (1) a *proof* that four rounds of theory do not force a redesign (which prevents a costly, seductive rewrite), and (2) test infrastructure that hardens V1 now AND serves the 5-year scenario. Both are high ROI; the kit is test code, so it adds no product complexity (AP11 is about product mechanism). This is the correct place for the research to have landed. Endorsed. On the proposed base64 hardening: cheap, do it — but as its own reviewed change, not smuggled into the kit.

## The Cartographer
Restate honestly: the headline is "the kernel is validated, not redesigned," and the impact is a conformance kit, not a new architecture. Frame I must hold against the pressure of the prompt: the founder pushed hard for "fundamentally different," and the disciplined answer is that the architecture *should not* fundamentally change because it was well-built and OMEGA validates it — the fundamental difference for a trust primitive's long life is conformance rigor, which is what shipped. Do not let the demand for drama manufacture a teardown. The one genuinely new posture ("specification + conformance kit," not "reference implementation") is real and a senior engineer would notice it; that clears #8. Everything else was correctly rejected (below).

# Rejected proposals (dead ends — preserved per the founder's instruction to show them)

1. **Graded-confidence verifier** (output a confidence scalar bounded by `I(C;E)` instead of the trichotomy). Rejected: unforced generality (AP11/DR9), and OMEGA-03 shows a high-stakes primitive above the penalty ceiling wants verify-or-exclude, not graded admittance. The trichotomy is correct.
2. **Channel-provenance annotation in the trace** (tag each check with its discriminating channel: signature / clock / revocation). Rejected as architecture: cosmetic; #8 fails (a senior engineer would not call it fundamentally different). The independence it would surface is *already* structural — the three channels arrive through three separate injected ports.
3. **Channel-independence static analyzer.** Rejected: redundant with the existing import lint (`scripts/check-imports.sh`) plus port separation, which already forbid coupling between the signature, clock, and revocation channels. Building it would duplicate enforcement.
4. **Currency-indexed verification inside Atlas** (OMEGA-03). Rejected: the currency/action-value decision is the RP's grant decision, outside Atlas's boundary (FM10). Belongs in a future numbered RFC if a scope act ever admits lower-stakes delegation, not in the kernel.
5. **A second, deliberately-independent reference verifier** to enable true differential testing now. Rejected for V1: maintaining two verdict implementations invites drift; property/metamorphic testing gives most of the value without the maintenance burden. Revisit when a real second implementation exists.

# Founder override (if applicable)

None.

# Open questions / proposed-not-applied

- **Strict base64 canonicalization at the verifier** (Red Team + Economist): reject non-canonical JWS encodings so tokens are byte-canonical, closing the malleability footgun by construction. **Proposed, not applied** — it is a change to `internal/record` acceptance semantics and deserves a deliberate founder call; benign under current design (revocation keys on instance ID; replay unmitigated by scope). Recommendation: do it as its own reviewed change.
- **True differential testing** awaits a second implementation; the kit is built to run against it (`Run(t, factory)`).
- **Corpus completeness** is asserted over the enumerated verdict space; the property suite is the backstop.

# Status
- decided: 2026-07-06
