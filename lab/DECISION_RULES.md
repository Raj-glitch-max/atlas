# Decision Rules

**Purpose.** This document defines, for every possible experimental outcome the lab may produce, what conclusion is *allowed* from the observed evidence and what conclusion is *forbidden*. It is the lab's anti-bias core. Its single load-bearing goal: **prevent confirmation bias and over-reading of experimental results.**

The lab's rule: a conclusion does not follow from the evidence in general — it follows from the evidence *under a rule frozen before the evidence was collected*. A conclusion reached by a rule not present in pre-registration is void, regardless of how strongly the evidence seems to support it.

---

## 1. Source of authority

Pre-registration primacy. Only criteria, thresholds, and outcome mappings fixed in the pre-registration *before* the run count as evidence; anything chosen after seeing the result does not.

- The pre-registration for EXP-001 (the C4 spike) is `LEVEL0_1_FEASIBILITY_GATE.md` §"Spike question (C4-only)" plus the candidate compositions enumerated there, hash-pinned at pre-flight (`EXPERIMENT_CHECKLIST.md` P-2).
- The applicable outcome table for EXP-001 is §3 of this document, hash-pinned at pre-flight (P-4).
- The Standards Editor's first audit question is always: *were the success and failure thresholds cited in the conclusion present in the pre-registration before the run?* If not, the conclusion is refused (post-run halt H-3).

---

## 2. Universal rules (apply to every experiment in this lab)

### U1 — Burden asymmetry

- To claim a **positive** ("X works", "X is solvable"), *every* pre-registered success criterion for the claim must hold AND no pre-registered failure criterion may trip, AND (for a conclusive verdict) the result must be reproduced (two-run rule, `LAB_README.md` §9).
- To claim a **negative** ("X fails", "X is not solvable"), *any one* pre-registered failure criterion tripping is sufficient, once the failure is reproduced and not attributable to substrate or instrumentation error.
- To claim **inconclusive**, the pre-registered criteria must be undecidable on the evidence AND the change-conditions for resolving the inconclusiveness must be stated.

The asymmetry is deliberate and load-bearing: positives carry the burden-of-proof, negatives carry the lighter burden-of-disproof. This is the opposite of what enthusiasm produces by default; that is the point.

### U2 — Conclusion must name the criterion

Every conclusion cites the specific success-criterion or failure-criterion ID it rests on (e.g., "F4 tripped: sniffer recorded one outbound packet from the RP during T04"). A conclusion that does not name the criterion ID it rests on is malformed and refused. "The run basically failed" is not a conclusion.

### U3 — Confidence label required

Every conclusion carries a project-confidence label (High / Medium / Low / None) with its evidence basis and a stated change-condition, per `agents/GOVERNANCE.md` §5 and the Glossary. A conclusion without a confidence label is refused.

### U4 — Negative-result legitimacy

A clean negative — a run where a failure criterion tripped, was reproduced, and the failure is a capability failure not an instrumentation failure — is a **complete result**, not a failure of the engineer. Judging a clean negative as "the experiment failed" treats the lab as a confirmation engine, which it is forbidden to be.

### U5 — Garden-of-forging-paths guard

The lab may not run many variants and report only the passing one. Every attempted run is logged (with its own run ID), every attempted composition is enumerated, and the count of attempts is part of the evidence. A positive reached after N failed attempts is reported as "1 of N attempts succeeded; the failing N-1 are logged as RUN-02..RUN-(N+1)," not as "the experiment succeeded."

This is the most actively-enforced anti-bias rule. A single passing run selected from many is the canonical shape of confirmation bias; the lab refuses to produce it.

### U6 — Multiple-comparison / false-positive budget

When an experiment tests many criteria or many test cases, the false-positive budget is pre-registered (e.g., "of K criteria, the expected number of spurious passes under the null hypothesis is K·p"). A number of passes within the false-positive budget is *not* evidence of a positive; it is evidence of noise. The budget is set in pre-registration, not after seeing how many passed.

### U7 — Reproducibility tier

- A conclusion from a single unreproduced run is *preliminary* and may not enter a journal entry.
- A conclusion from a reproduced run by a second engineer on a fresh substrate is *evidence-tier-2* and may enter a journal entry.
- A conclusion whose reproduction contradicts the original is *contradicted*; both runs stand, no conclusion is admissible, and a pre-reg amendment or fresh experiment is required.

### U8 — No scope creep from a result

A lab conclusion is about *what was tested*. A positive on the C4 spike permits a conclusion about C4; it does not permit a conclusion about H1 as a whole, about P5 as a product, about Level 2 being unnecessary, or about the project moving forward. Each of those is a separate decision taken upstream on the journal entry the lab produces.

A conclusion that exceeds what was tested is over-reading and is refused.

### U9 — No re-framing of failed criteria

A failure criterion that trips may not be re-described post-hoc as "not the real criterion" or "implementer's misunderstanding." If the criterion was in the pre-reg, it is the criterion. A change to a criterion is a pre-reg amendment, recorded and dated; it does not retroactively rescue a run.

### U10 — Ockham's razor for divergent outcomes

If the evidence supports two outcome classes, the lab adopts the more conservative (the one permitting fewer downstream conclusions) pending founder decision. It does not pick the favorable reading.

---

## 3. EXP-001 (C4 spike) — outcome table

Four outcome classes are possible for EXP-001, derived from the feasibility gate. The spike tests the C4 component claim only: *offline, fast, partition-tolerant, independent revocation, achieved by composition of existing standardized primitives under the S1–S5 reading of H1.*

For each outcome: the **evidence** that must be observed, the **conclusion** that is allowed, and the **conclusions** that are forbidden.

### Outcome α — A working composition satisfies the C4 spike criterion

**Evidence observed (all required):**
- A candidate composition (named, diff-recorded against stock libraries, with its S1–S5 interpretation stated) satisfies the spike's pre-registered C4 criterion: revocation observable within R seconds (R as set in S1), no live call from the RP to domain A, no third-party decision-making broker, stock JWT/JWS/COSE verifier unmodified.
- The result is reproduced (U7) by a second engineer on a fresh substrate.
- No failure criterion (F1–F10 of the Level 2 protocol) applicable to the spike tripped; the garden-of-forging-paths count (U5) is reported.

**Conclusion allowed:**
- "C4 is solvable, on this substrate, by composition of existing standardized primitive(s): `<named primitives>`, under the S1–S5 reading as `<state which reading>`." (Confidence label required per U3.)
- "The novelty claim of H1 narrows to the revocation layer atop SPIFFE + RFC 8693; an OSS novelty audit is required before any novelty claim for the layer itself."

**Conclusion forbidden:**
- "H1 is proven" or "H1 works" — only C4 was tested; C1/C2/C3 and the multi-hop per-hop authorization (C11) were not.
- "P5 should be built" / "P5 is a viable product" — product and scope decisions are upstream, not lab (U8).
- "The composition is novel" without an OSS novelty audit against existing SPIFFE-interop revocation tools.
- "Level 2 is unnecessary" — that is a founder disposition taken on the journal entry, not a lab conclusion; the lab produces the evidence, the founder retires Level 2 or not.
- Extrapolation to multi-hop (C11) or to a substrate other than the one tested.
- Reporting a single passing run without the count of failed attempts (U5).

### Outcome β — No composition works; the gap is a technology gap (not logical impossibility)

**Evidence observed (all required):**
- Every enumerated candidate composition in the pre-reg was attempted and logged as a separate run; none satisfies the C4 criterion.
- For each failed composition, the failure is identified as a *capability* failure (specific clause: revocation observability, no-live-call, or partition-tolerance) — not an instrumentation failure, not a substrate bug, not an implementation error.
- The failure is reproduced (U7) on a fresh substrate, ruling out setup error.
- A formal argument is recorded that the failure is a *technology* gap, not a *logical* impossibility (the latter is outcome γ). Specifically: there is no information-theoretic obstruction; the obstruction is that no existing standardized primitive supplies the needed channel.

**Conclusion allowed:**
- "C4 is not solvable, on this substrate, by composition of existing standardized primitives; the gap is a technology gap (named clause), not a logical impossibility." (Confidence label required.)
- "Level 2 is justified and may be narrowed to C4 + C11's per-hop sub-claim, eliminating re-certification of solved components." — *as a recommendation to the founder via the journal entry, not as a lab action.*

**Conclusion forbidden:**
- "H1 is falsified" — only C4 was tested; C1/C2/C3 remain solved.
- "P5 is dead" / "P5 should be retired" — scope and candidate-set decisions are upstream (U8).
- "No delegation primitive is possible" — over-extrapolation; only *these compositions* were tested. A non-standardized or new-primitive approach was not tested by EXP-001.
- Treating a single-run failure as conclusive without reproduction (U7).

### Outcome γ — C4 is a logical impossibility under the S4 combined reading

**Evidence observed:**
- A formal argument (not a run) demonstrates that under the S4 combined reading of T4 ∩ T5 — observe a revocation set *during* a partition that isolated the RP from the issuer at the time of revocation, within a window shorter than partition-recovery — no technology, present or future, can satisfy the criterion. The information has no channel to cross the partition within the window.
- The argument is reviewed by the Standards Editor and recorded as an artifact under `analysis/`.
- (If the spike's pre-reg adopted a non-S4 reading by setting S1 explicitly, this outcome requires the argument that the S4 reading is the *correct* reading of H1, which is a separate analysis.)

**Conclusion allowed:**
- "Under the S4 combined reading, C4 is a logical impossibility; H1 as written (with S4 in scope) is not satisfiable by any technology, present or future. The S4 ambiguity must be resolved by a scope act before any empirical experiment can be informative about H1." (Confidence: High — this is a deductive result, not a measurement.)

**Conclusion forbidden:**
- "H1 was falsified by the experiment" — H1 was falsified by *definition* under the S4 reading, not by the experiment. Labeling a definitional finding as an experimental result misstates the type of the evidence.
- "A technology gap exists" — the finding says nothing about technology; technology is moot under a logical impossibility.
- Running further experiments against the un-corrected H1. The lab may not; the conclusion requires a scope act upstream.
- Continuing to spend experimental effort on the S4 reading as though it were empirical.

### Outcome δ — C4 unsolved AND S1–S5 unresolvable without empirical input

**Evidence observed:**
- Compositional attempts are inconclusive (outcome α and outcome β both fail to meet their evidence requirements).
- The desk audit cannot resolve S1–S5 because the resolution is itself empirical (e.g., the revocation-observability latency bound R cannot be set a priori without first measuring realistic substrate behavior).
- The empirical unresolvability is demonstrated, not assumed: the spike produced measurements whose spread or character shows that the missing parameter cannot be chosen without the very empirical work the desk audit was meant to defer.

**Conclusion allowed:**
- "The desk audit and the C4 spike cannot settle the matter; empirical testing at Level 2 (as frozen) is the only remaining instrument. S1–S5 resolution, where resolvable only empirically, is folded into Level 2's instrumented run." (Confidence label required; this corresponds to feasibility-gate outcome C — Level 2 justified as-is — reachable only after the spike fails to settle.)

**Conclusion forbidden:**
- "H1 works" or "H1 fails" — the spike was inconclusive; no verdict is admissible.
- "Level 2 is unnecessary" — the opposite is the conclusion.
- Skipping the S1–S5 resolution that *is* resolvable by scope act; only the empirically-unresolvable subset is folded into Level 2. The resolvable-by-act subset must still be resolved upstream before Level 2 runs.

---

## 4. Forbidden conclusion patterns (explicit, with examples)

These are over-reading patterns the lab refuses regardless of how strongly the evidence seems to support them. Each is stated as the pattern, an example, and why it is refused.

| Pattern | Example | Why refused |
|---|---|---|
| Adjective-graded verdict | "T04 essentially passed" / "the spike mostly worked" | A criterion passes or fails; "essentially" is a post-hoc relaxation. State which criterion held and which did not. |
| Extrapolation beyond substrate | "C4 holds in general" (from a two-SPIRE-substrate result) | A composition verified on this substrate is verified *on this substrate*. Generalization is a separate claim requiring separate evidence. |
| Treating unrun cases as passed | "T05 would obviously pass given T04 passed" | Unrun cases are unrun; no claim about them is admissible. |
| Reframing a failed criterion | "F4 isn't really the criterion; the real test is latency" | The pre-reg defined F4; U9 forbids redefinition mid-run. |
| Novelty claim without OSS audit | "This composition is genuinely new" | Without an audit of existing SPIFFE-interop revocation tools, novelty is asserted not demonstrated. |
| Selective reporting | "The spike succeeded" (after 7 failed runs and 1 passing run) | U5: count is part of the evidence. |
| Scope creep from a positive | "C4 works, so we should build P5" | U8: a C4 result is about C4, not about P5. |
| Scope creep from a negative | "C4 fails, so P5 is dead" | U8: a C4 result does not retire the candidate set; C1/C2/C3 remain solved and other framings exist. |
| Mislabeling a definitional finding | "The experiment showed H1 is impossible" (when S4 makes it so by definition) | Outcome γ is deductive, not experimental; mislabeling inflates the experiment's standing. |
| Confidence without evidence | "Confidence: High" without a cited witness per claim | Project doctrine: confidence without evidence is forbidden. |
| Conclusion without change-condition | "C4 is solved" with no statement of what would change the conclusion | Project governance: confidence without a change-condition is forbidden. |

---

## 5. Required conclusion structure

A lab conclusion is admitted only in this shape:

```
CONCLUSION — EXP-NNN/RUN-NN — outcome class α|β|γ|δ
UTC: <timestamp>
Pre-reg: <path>#sha256:<hash> (frozen at pre-flight)
Outcome class: <α|β|γ|δ>, per DECISION_RULES.md §3 (<version>/<sha256>)
Evidence: <list of (criterion ID → observed fact → evidence ref)>
Burden: <positive (all S-criteria held, no F tripped, reproduced) | negative (which F tripped, reproduced) | inconclusive (which undecidable, change-condition)>
Confidence: <High|Medium|Low|None> — <evidence basis> — change-condition: <condition>
Reproduction: <second engineer, fresh substrate, run ID, same/narrower/different result>
Forbidden readings checked: <list of §4 patterns explicitly considered and rejected>
Allowed conclusion (verbatim from §3 outcome row): "<text>"
Standards Editor audit: <pre-reg primacy held: Y/N> — <initials>
```

A conclusion in any other shape is refused by the lab. The journal entry the lab produces (per §6) quotes this structure verbatim.

---

## 6. Where the conclusion goes

A conclusion that passes §5 is recorded as a journal entry at `agents/journal/<YYYY-MM-DD>-<slug>.md`, per `agents/GOVERNANCE.md`. The journal entry:

- quotes the §5 conclusion structure verbatim;
- preserves any dissent (e.g., the adversary's divergence from the Engineer's reading) per the preserve-dissent rule;
- is the lab's only channel to the rest of the project.

The lab does not, on the strength of a conclusion, retire the frozen Level 2 protocol, open Product Management, change the candidate set, or recommend funding. Each of those is a founder decision taken *on* the journal entry, not a lab action. The lab's work ends at the journal entry; the founder's begins there.

---

## 7. Change control for this document

A change to §3's outcome table is a pre-registration change for any experiment whose pre-reg post-dates the change but does not alter a sealed run's conclusion. A change to §2's universal rules is a lab-process change logged per `LAB_README.md` §12.

A change that would retroactively alter a sealed run's admissible conclusions is forbidden. Sealed runs are immutable; their conclusions stand as pre-registration at the time of their run allowed.
