# P5 Falsification Experiment Specification

**Role basis.** Principal Distributed Systems Engineer (clean partition semantics, observable measurements, reproducible substrate) + Principal Security Engineer (threat-model-aware, adversarial-first) + Research Engineer (pre-registered thresholds, controls explicit, no post-hoc goalpost moves).

**Source of truth (approved, not regenerated):** `TECHNICAL_VALIDATION.md` §P5 (hypothesis, hidden assumptions, classification) and §"Strict Engineering Comparison" (composability reasoning); `FOUNDER_PROBLEM_FIT.md` §P5 §5–7 (silent-failure pattern, interoperability gate, personal-risk asymmetry).

**This document does NOT** build the experiment, write code, choose product, choose stack, or discuss business. The output is a specification only.

---

## 0. Threat model the experiment must honor (precondition)

The hypothesis makes a specific composition claim. We model an adversary only to the extent required to test it; we do not model a full cryptanalytic adversary. The threat model under test:

- **T1 — Online honest-but-curious cross-domain operator.** The relying party does not trust domain A's runtime; it only trusts domain A's published keys and the delegation token it received.
- **T2 — Replay / re-use of an issued delegation after revocation.** Standard token-replay attacker, no key compromise.
- **T3 — Token tampering.** Scope, identity, or TTL fields are edited by a man-in-the-middle or compromised daemon.
- **T4 — Trust-bundle lag.** Domain A rotates its signing keys. Verification must remain possible with whatever keys the relying party last received (no implicit reliance on a live call).
- **T5 — Partition during verification.** The relying party cannot reach any service in domain A or a third-party broker during the verification window.
- **OUT OF SCOPE:** Key compromise of the relying party, side-channels inside primitives (SPIFFE/SPIRE, RFC 8693), and quantum-grade signature concerns — these are delegated to the underlying primitives' threat models.

An experiment that does not isolate T1–T5 from the OUT-OF-SCOPE concerns is invalid; success or failure outside that isolation tells us nothing about the hypothesis.

---

## 1. Exact hypothesis being tested

> **H₁.** Given two independent SPIRE deployments that share no servers, no clocks they trust beyond what each publishes, and no third-party broker, it is possible to compose an existing SPIFFE SVID with an existing RFC 8693 token-exchange flow to produce a *delegation* token — "X in domain A is acting on behalf of principal Y in domain B, with scope S, valid until T, revocable independently of Y's or X's SVID" — such that a relying party *operating inside domain B and forbidden from making a live call to any service in domain A or to any third-party broker* can verify, reject replay, reject tampering, and observe revocation **within a verification window strictly bounded above by the TTL of the underlying Y SVID**, with verification latency below 100 ms end-to-end, using only unmodified off-the-shelf verification libraries.

The strictly-bounded verification window and the unmodified-libraries clause are the load-bearing qualifiers. Without them, "delegation" reduces trivially to "a short-lived cert signed by a key the relying party already has," which is the null hypothesis that SPIFFE/SPIRE already satisfies and would tell us nothing.

This is the single central claim. Everything else in the experiment is a control.

---

## 2. Why this hypothesis is the critical one

**Why it's central.** The Founder-Problem Fit document (§5, §6) explicitly names two distinct risks for P5: (a) *silent failure* — a "looks-finished" delegation primitive the builder believes is sound but isn't; (b) *interoperability failure* — building something that works but requires new protocol adoption from every counterparty. Both risks converge on the same falsifiable claim: **delegation that composes from existing primitives without new protocol, with no shared broker, and with semantics richer than a static short-lived certificate.** H₁ is the single claim whose truth or falsity decides whether P5 is (i) a primitive that didn't previously exist in open form, or (ii) a wrapper around existing SPIRE deployment that does not justify the project's framing.

**Why the other sub-claims are not the target here.** The Technical Validation lists five hidden assumptions (continuous re-attestation, cross-domain delegation without broker, non-K8s attestation, ephemeral-issuance latency, composability). Of these:

- *Continuous (post-issuance) re-attestation* is a separate hypothesis whose falsification requires its own experiment involving timing-side-channel reasoning; bundling it here would muddle what's being tested.
- *Non-K8s attestation* is a substrate-variety test, not a falsifier of the central delegation claim; it cannot fail this experiment meaningfully.
- *High-churn ephemeral issuance latency* is a performance test, not a correctness claim; it can't falsify H₁ — it can only posteriorly bound it.

Only **delegation across trust domains without a shared broker, composable from existing primitives** is the load-bearing falsifier. Everything else in this document is either a control that keeps the experiment honest or an explicitly out-of-scope follow-up.

**Why a falsification-optimized version of this is non-trivial.** A naive experiment that "tries to compose" will confirm trivially if delegation is permitted a live broker call — every delegation scheme in existence works when the relying party has a live trust anchor. The force of H₁ comes from *forbidding the live broker call* and *forbidding library modification*. Those two constraints are what makes the experiment capable of actually failing.

---

## 3. Smallest possible experiment

Five components. Each is as small as the falsifiability constraints allow.

### 3.1 Substrate (build once, run nine test cases)

- Two SPIRE deployments, **completely independent**: distinct servers, distinct databases, distinct signing keys, clocks independently skewed within 250 ms of each other.
- One workload in domain A (delegating principal X).
- One principal (Y) registered in domain B but not present as an active workload.
- A relying party (RP) running inside domain B.
- RP is network-isolated from domain A's SPIRE server and from any plausible third-party broker during every test case.

This is the only plausible substrate. It will not recognize any service in domain A as a verification endpoint; verification is offline-or-nothing.

### 3.2 Configuration locks (pre-registered, do not change during test)

- All SPIFFE federation between the two domains is **disabled** at the SPIRE configuration level. Federation is the closest SPIRE-shipped mechanism that would otherwise smuggle a shared broker back in; we lock it out so any "compose" answer does not silently rely on it.
- The RP's verifier is a **stock JWT/JWS / COSE verifier** with only the public keys of domain A and domain B explicitly loaded. No source modification of SPIRE, SPIFFE libraries, or RFC 8693 client libraries is permitted; if it doesn't work with stock libraries, the experiment is informing us that the hypothesis fails on the composability requirement.
- TTL of delegation token: **30 seconds**. TTL of Y's SVID: **10 minutes**. The narrow delegation window inside the wider Y-SVID window is the test surface — if a "delegation" can only exist for the full SVID lifetime, it has not been demonstrated as a separate primitive.

### 3.3 Test cases (nine; each is binary pass/fail by §5/§6)

| # | Test case | Exercise of H₁ |
|---|---|---|
| T1 | RP verifies identity, scope, and TTL of an untampered delegation token issued at T₀ | happy-path identity check |
| T2 | Replay of the same token post-TTL | replay rejection |
| T3 | Token with one character mutated in the scope field | tamper rejection |
| T4 | Token whose revocation status was set at T₀+5s; RP verifies at T₀+7s using only the token (no live consultation) | independent revocation |
| T5 | RP verification while SP1 in domain A and any plausible broker are network-unreachable for the full 30-second TTL window | offline / partition tolerance |
| T6 | Domain A rotates its key bundle at T₀; the RP holds only the prior bundle and verifies at T₀+1s | trust-bundle lag |
| T7 | A delegation chain two levels deep: X → X' → Y (X' is registered in domain A and has its own SVID; X delegates to X' who delegates to Y) | chain-of-delegation correctness |
| T8 | An ephemeral agent (created at T₀, destroyed at T₀+5s) requests a delegation and the relying party verifies within the verification window, while domain A's SPIRE server is unreachable | ephemeral-workload substrate |
| T9 | Control: the same RP verifies a normal SPIRE SVID directly issued in domain B (no cross-domain composability) — this case must succeed trivially | proves RP setup is honest; not testing H₁ |

Test cases T1, T8 probe happy-path; T2–T6 probe adversarial and operational conditions; T7 probes chain depth; T9 is a control.

### 3.4 What is *not* in the experiment (explicitly)

- Continuous (post-issuance) re-attestation detection — separate hypothesis.
- Non-Kubernetes attestation substrate — separate hypothesis.
- Adversarial cryptography under key compromise — out of scope per §0.
- Performance tuning beyond the 100 ms verification threshold — separate experiment.
- Agent-to-agent signed communication payload format — separate problem.

These are named here so a future reader doesn't infer that "the experiment passed" includes them.

---

## 4. Required environment

Minimum to be reproducible and falsifiable.

| Component | Minimum |
|---|---|
| Compute | 3 hosts (Linux x86_64, 4 vCPU / 8 GiB each, dedicated), one each for: SPIRE domain A, SPIRE domain B, the relying party. All clocks NTP-disciplined within drift bounds; clock skew between the two SPIRE servers managed by sleeping for randomized intervals. |
| Network | All hosts reachable initially; per-test isolation enforced via firewall drop rules (T5, T8), not via clean-room networks — the fire-and-reconnect workflow is part of what we're testing. |
| SPIRE | Latest stable SPIRE server and agent (2026 vintage). Identical Helm-free, single-server config in both domains. Federation **disabled** in both. |
| Primitives | SPIFFE workload API for SVID issuance; RFC 8693 token-exchange implementation that issues standard JWT or CWT tokens with claims restricted to registered names; stock JWS / COSE verifier with public-key resolution. **No source modification to any of these.** |
| Time | All clocks visible to the test harness; the harness records ms-precision measurements for every verification and every "live call attempted" sniffer. |
| Telemetry | An out-of-band sniffer on the RP host that records every TCP/UDP attempt from the RP to any address outside `domain-B.prefix`. The sniffer is the *fail-the-experiment* instrument: any verified-OK test case that produced a sniffer hit on an address outside domain B is automatically graded FAIL on the offline/partition-tolerance criterion regardless of any other signal. |
| Personnel / handoffs | One engineer for setup, plus a second engineer for the adversary role (T2–T7). The second engineer does **not** see §5 thresholds before running. |
| Time-of-record | Wall-clock time included with every test artifact; reproducibility inside a 24-hour window is sufficient. |

A single laptop with three VMs or containers satisfies this, provided the network-isolation test cases use host-level firewall rules rather than container network namespaces (containers share kernel and sometimes share trust state — that's a different experiment).

---

## 5. Success criteria

Pre-registered. Defined before any test case runs. The hypothesis is *not* supported unless all are met simultaneously.

| ID | Criterion | Why hard |
|---|---|---|
| S1 | T1, T8 verify-OK within **100 ms** end-to-end at the RP, measured from "token received at RP" to "RP verdict returned." | Hard ceiling on the practical claim; above this it isn't a usable delegation primitive. |
| S2 | T2, T3, T4 each return a verifiable **rejection** at the RP with no spurious network calls to anything outside domain B (sniffer clean). | Tamper and replay must be caught locally; "live call to remake decision" is verification failure. |
| S3 | T4 revocation works **without revoking X's or Y's underlying SVID**. The X and Y SVIDs remain verifiable (as in T9) at the conclusion of T4. | Independent revocation is what makes delegation a separate primitive from short-lived cert, the falsifier's whole point. |
| S4 | T5 verifies-OK within 100 ms despite partition; the sniffer records **zero** outbound TCP/UDP from the RP to anything outside domain B for the full duration of the test. | The live-broker clause. |
| S5 | T6 verifies-OK using only trust material previously received by the RP. RP's verifier accepts the prior key bundle per its own policy and rejects anything else. | Trust rotation without implicit broker. |
| S6 | T7 verifies-OK with chain-of-delegation verifiable under the same conditions as T1; sub-claims (X' is reachable in chain, X' → Y is verified, X → X' is verified) each pass independently. | Multi-hop delegation is a structurally separate failure mode from the single-hop case; it must work or the primitive isn't actually a delegation primitive. |
| S7 | The implementation required **zero source modifications** to SPIRE / SPIFFE libraries and zero source modifications to the RFC 8693 client or stock JWS / COSE verifier. | The interoperability gate. Implementation may add new code atop stock libraries; it may not modify them. |
| S8 | T9 verifies-OK (substrate-control passes). | If T9 fails, the substrate is broken — no claim can be made about the experimental substrate; rerun after fix, do not infer from this run. |

A single criterion missing means the experiment records **FAIL**.

## 6. Failure criteria

Pre-registered. The hypothesis is *falsified* if **any** of the following is true (each criterion is a distinct, mechanized detection):

| ID | Failure condition |
|---|---|
| F1 | S1 violated at any test case: end-to-end verification exceeds 100 ms in T1, T5, T6, or T8 even once. |
| F2 | S2 violated at T2, T3, **or** T4: the RP accepts a tampered/replayed/revoked token, **or** the sniffer records an outbound call outside domain B during the verification. |
| F3 | S3 violated: revocation cannot be effected without also invalidating Y's or X's SVID. After T4, T9 must still pass. |
| F4 | S4 violated: T5 cannot complete verification within the 30-second window, **or** the sniffer records even a single outbound call outside domain B. |
| F5 | S5 violated: T6 verification depends on a callback, live key lookup, or other interaction not reducible to "load trust bundle at startup / on demand." |
| F6 | S6 violated: T7 fails for any sub-claim, even when T1 passes. |
| F7 | S7 violated: implementation modified source of any SPIRE / SPIFFE / RFC 8693 / stock-JWS library; **or** a custom binary token format was required rather than standard JWT / CWT. |
| F8 | S8 violated: T9 fails. In this case, the substrate experiment is broken; it is *not* a falsification of H₁, but the experiment must be re-run cleanly before any H₁ claim can be drawn. |
| F9 | Any single test case is unreproducible across two independent runs (engineer handoff, fresh substrate) without re-engineering. Reproducibility failure is not "luck of the draw" — it is an experimental-design failure. |
| F10 | Any criterion depends on a configurable number the engineer is permitted to tune after viewing test-case results. |

F1–F8 are the structural falsifiers. F9–F10 are the methodology falsifiers — they don't falsify H₁ specifically, they falsify the experiment's status as evidence about H₁.

---

## 7. Evidence produced

Each test case produces a fixed artifact bundle. The bundle is the evidence; it is preserved verbatim and is the only thing engineering claims about this experiment rest on.

| Artifact | Source | Format |
|---|---|---|
| Token issued | Test harness | JWT / CWT, recorded as issued |
| Verification verdict (`OK` / `REJECT-REPLAY` / `REJECT-TAMPER` / `REJECT-REVOKED` / etc.) | RP | terse log line |
| Latency (ms) | Test harness | single float |
| Outbound-call trace from RP host | Out-of-band sniffer | pcap excerpt (only meaningful in T2–T6, T8) |
| Wall-clock timestamp (start, end) | Test harness | ISO-8601 |
| Configuration snapshot of both SPIRE domains | Test harness | YAML or equivalent, hash-pinned |
| Concrete source files used in the implementation | Engineer | version-controlled, with diff against stock — must be empty for S7 |
| Engineer commentary recorded inside 24h of running | Second engineer | short note per test case |

The bundle for **every** test case is preserved, including T9 (the control). Convenience summaries ("T1–T8 all passed") are not accepted as evidence; raw artifacts are.

---

## 8. What uncertainty disappears if it succeeds

If H₁ is supported (all S1–S8 satisfied, no F1–F10 tripped), the following uncertainties collapse to specific, recorded facts:

- **Whether open, interoperable delegated identity composes from existing primitives.** This is now an artifact, not a claim. The Founder-Problem Fit document's §6 risk ("landing on the wrong side of the interoperability constraint") is closed: the experiment either confirms the interoperability constraint holds or rejects it by §S7.
- **Whether the "no shared broker" requirement is achievable in principle for cross-domain delegation.** Per §S4, it is a measured fact for the substrate, not a hope.
- **Whether revocation is a separate primitive from short-lived cert issuance.** Per §S3, it is or it isn't.
- **Whether the Cartographer's lift on P5 in the R1 cycle (the only council upgrade in this set) was right.** A positive result here directly rewards that signal; a negative result directly overrules it. The signal has, for the first time, evidence rather than analogy behind it.

The following are **not** collapsed by success and are explicitly named as remaining:

- Continuous post-issuance re-attestation (separate hypothesis).
- Operational cost of running the substrate at production scale (separate measurement).
- Substrate portability beyond Linux-on-x86 (separate substrate experiment).
- Threat-model coverage beyond T1–T5 (e.g., insider-attack within a domain).
- Multi-trust-domain compositions beyond 2 domains with mixed trust (heuristic; the test hinges on 2).

## 9. What uncertainty disappears if it fails

If H₁ is falsified (any F1–F10 trioped), the following collapse to specific, recorded facts:

- **Delegation that meets H₁ is not achievable with the existing primitives** (on this substrate, with unmodified libraries, at this verification window). This is a measured limit, not a speculation.
- **P5's central technical claim is reframed.** The hypothesis is restricted to weaker forms. Specifically, at least one of the following holds, and the experiment tells us which:
  - Composition requires live broker calls; the primitive reduces to "short-lived cert with online trust anchor."
  - Revocation is not independent of SVID issuance; the primitive reduces to "short-lived cert with TTL."
  - The primitive requires a new protocol beyond what RFC 8693 + standard JWT/COSE provides; the interoperability gate is closed.
- **The Cartographer's lift is overruled as a basis for product framing.** A failure here is direct evidence the council's only upgrade does not survive engineering contact, independent of any market considerations.
- **P5's engineering tier drops from "research-level, novel upside" to "compounding wrapper around existing SPIRE deployment,"** as flagged in §6 of the Founder-Problem Fit document. The hidden assumption that "more SPIRE configuration is sufficient" (an explicit watchpoint in §P5 §3 of the Technical Validation) collapses to a documented fact: the existing SPIRE configuration is sufficient, full stop.
- **A clear delisting-or-reframe decision for P5 becomes available.** The Founder Decision Brief's framing of P5 ("viable near-term build target / multi-year research bet / OSS-contribution reclassification") resolves to: *OSS-contribution overlay or reframe*, with the technical reasoning recorded.

The following are **not** collapsed by failure and remain:

- Continuous post-issuance re-attestation is still open and still potentially meaningful — it's a separate hypothesis this experiment does not test.
- Whether *some other* composition (a different RFC, a different token format, or a different base primitive) could meet H₁ is *not* disproven by this experiment; this experiment only widens the falsification to a specific compositional family. The experiment disconfirms "compose from SPIFFE + RFC 8693 to meet H₁" — not "any open delegation primitive is impossible."
- Market and pricing questions remain wholly out of scope per the brief.

---

## 10. Whether another experiment becomes necessary afterward

Yes — and the specifics depend on the outcome, which is exactly why the experiment has to be falsification-first rather than confirmation-first.

### 10.1 If H₁ is supported (S1–S8 all satisfied)

The follow-on experiments are required to make the result productization-ready:

| # | Follow-on | What it tests | Trades |
|---|---|---|---|
| F-A | Continuous re-attestation (post-issuance) | Does a delegation correctly invalidate when the workload's attestation signals shift materially *after* SVID issuance, *without* invalidating unrelated delegations? | Separate hypothesis; this is currently the cold spot in the brief. |
| F-B | Trust-bundle rotation across the verification window | What happens when domain A rotates its key and the RP holds both old and new bundles in its verification window? Is there a measurable availability or security gap? | Operational reality check; not a yes/no falsifier. |
| F-C | Adversarial cryptography under realistic adversarial models | Move from T1–T5 (discrete attackers) to a T6-class attacker with token-format-shape flexibility | Goes deeper than T1–T5 threat model. |
| F-D | Multi-trust-domain compositions (≥3 domains, ring topology) | Does H₁ hold transitively across multiple domain-to-domain handoffs? | Generalization; outside the substrate tested here. |
| F-E | Performance at threshold (1k delegations / second sustained) | Does the 100 ms verification hold under load? | Operational, not falsification. |

F-A through F-D are falsification experiments; F-E is operational benchmarking.

### 10.2 If H₁ is falsified (any F1–F8 trioped)

Different follow-ons, because the question changes:

| # | Follow-on | What it tests | Why |
|---|---|---|---|
| F-F | Weakest-link identification across F1–F8 | Which of the *eight* structural failures actually tripped? Bounds can tighten: which specific clause of H₁ is not satisfiable as stated? | Tells us exactly which clause of H₁ the experiment killed, so a reframe can target the minimum modification. |
| F-G | Composition with a *different* standard primitive (e.g., DPoP, SD-JWT, macaroon, biscuit, anonymous credentials) | Can H₁'s clauses (especially S3 — independent revocation, and S4 — offline verification) be satisfied with a non-SPIFFE / non-RFC-8693 primitive base? | Reframes the question from "which wrapper around existing primitives" to "which primitive family works at all." |
| F-H | Composition with a *new* (non-stock) protocol and explicit "every counterparty must adopt" semantics | What is the smallest bespoke protocol that meets H₁? | Establishes the cost of H₁'s clauses in protocol-design terms; useful only as a ceiling on cost. |
| F-I | Reframe test: can P5's opportunistic claim be salvaged by *narrowing* the claim to one-clause-at-a-time? (E.g., drop the "no broker" clause and re-test) | Tells us whether there's a strictly weaker P5 claim worth keeping in the candidate set as a fallback. | A "yes, but simpler" reframe is logged as a candidate-set adjustment — **not** as a replacement for H₁. |

F-F is the immediate follow-on if H₁ fails; F-G through F-I are downstream of it.

### 10.3 Substrate-only failure (F8: T9 fails)

This is **not** an H₁ outcome but still requires a follow-on:

- **F-J — Substrate integrity repeatability test.** Build a fresh substrate (new SPIRE instances, new keys, new clocks). Run T9 until it succeeds. Without a passing T9, no H₁ evidence at all is admissible. This follow-on is methodological, not scientific.

### 10.4 Reproducibility failures (F9, F10)

- **F-K — Independent reproducibility test.** A second engineer re-runs the full experiment from scratch (without access to the first run's artifacts beyond the configuration snapshot). If their results agree, methodology holds; if not, the experiment itself is broken and the artifact is not evidence about H₁.

### 10.5 What is **not** owed by H₁

By design, this experiment does not produce evidence about:

- Agent-workload ephemerality beyond T8's specific scenario.
- Operational scale-up of the substrate.
- Whether market demand exists for the primitive (which belongs in Wave 5 of the Research Program — Task 8 — *not* here).
- Whether the founder has the profile to ship this responsibility-grade software (per `FOUNDER_PROBLEM_FIT.md` §"What's missing" — explicitly unfilled intake).

The experiment is engineered. Market, founder, business — out.

---

## Provenance and confidence

**Provenance.** Source of truth: `TECHNICAL_VALIDATION.md` §P5 (approved); secondary anchor: `FOUNDER_PROBLEM_FIT.md` §P5 (approved); Research-program Wave 4 / Task 4 plan: `RESEARCH_PROGRAM.md` §4 Wave 4 gate (approved). This document specifies the experiment; it does not execute it.

**Confidence labels** (per the project's standing doctrine):

- **H₁ as stated:** Medium-High. Direct restatement from the Technical Validation's central technical hypothesis; the qualifiers "without a shared broker" and "without modifying stock libraries" are our explicit additions to render it falsifier-shaped, not the Technical Validation's wording — those additions are *engineering constraints added to make the hypothesis testable,* not revisions to its content. The qualifiers qualify the falsification condition, not the claim.
- **Falsifiability of the experiment as specified:** High. F1–F8 are mechanizable, pre-registered, and structurally prevent post-hoc goalpost movement. F9–F10 honestly define when the methodology itself is broken rather than the hypothesis.
- **Predictive confidence about outcome:** None — and per project doctrine (confidence-without-evidence is forbidden), refusing to estimate is itself the correct label here. Nothing in this document is evidence about whether H₁ is true.
- **Change conditions** that would alter the experiment specification:
  - If Technical Validation §P5 §3's classifications are revised (e.g., a hidden assumption is reclassified), §2's "why this is the critical one" is revised accordingly.
  - If `FOUNDER_PROBLEM_FIT.md` §P5 §6's interoperability-gate risk is reframed, §S7's verification criterion narrows or widens correspondingly.
  - If a new fact emerges (e.g., a published SPIRE federation primitive that eliminates the no-shared-broker requirement), the experiment's substrate in §3.1 widens; the experiment re-runs; the verdict applies to the wider substrate, not the original one.
  - All such revisions documented as journal entries before the experiment runs.

**Stop rule.** This document is the experiment specification. It does not build the experiment. It does not write or run code. It does not choose a protocol. It does not enter PM, BA, architecture, or implementation phases. The next action belongs to the founder and is either (a) approve the specification as written, (b) approve with revisions, or (c) redirect to a different falsifier. No follow-on artifact is produced by this document.

<!-- checkpoint: feat(issuance): implement boundary check (#71) -->
