# Level 0 / Level 1 Feasibility Gate — P5 / H1

**Status.** A feasibility audit of H1, performed before the frozen Level 2 falsification experiment (`P5_FALSIFICATION_EXPERIMENT.md`) is allowed to run. The Level 2 protocol is **frozen**; this document does not modify it, does not redesign it, does not write code, does not draw architecture, and does not generate implementation. The single objective is to determine whether the expensive Level 2 experiment is warranted.

**Role basis.** Principal Distributed Systems Engineer + Principal Security Architect + Standards Editor.

**Hypothesis under audit (quoted verbatim from the frozen `P5_FALSIFICATION_EXPERIMENT.md` §1 — not paraphrased, not regenerated):**

> **H₁.** Given two independent SPIRE deployments that share no servers, no clocks they trust beyond what each publishes, and no third-party broker, it is possible to compose an existing SPIFFE SVID with an existing RFC 8693 token-exchange flow to produce a *delegation* token — "X in domain A is acting on behalf of principal Y in domain B, with scope S, valid until T, revocable independently of Y's or X's SVID" — such that a relying party *operating inside domain B and forbidden from making a live call to any service in domain A or to any third-party broker* can verify, reject replay, reject tampering, and observe revocation **within a verification window strictly bounded above by the TTL of the underlying Y SVID**, with verification latency below 100 ms end-to-end, using only unmodified off-the-shelf verification libraries.

---

## Evidence basis (provenance)

This audit's claims are grounded in fetched canonical sources, not recall, where load-bearing. Each fetched source is named at the claim that rests on it.

- **RFC 8693** ([datatracker.ietf.org/doc/html/rfc8693](https://datatracker.ietf.org/doc/html/rfc8693)) — fetched 2026-07-04. Verbatim confirmation: `act` (actor) claim §4.1 defines delegation; nested `act` represents a delegation chain ("a chain of delegation can be expressed by nesting one `act` claim within another"); nested prior actors are "informational only" for access-control decisions (only the outermost actor matters).
- **draft-ietf-oauth-status-list** ([datatracker.ietf.org/doc/html/draft-ietf-oauth-status-list-10](https://datatracker.ietf.org/doc/html/draft-ietf-oauth-status-list-10)) — fetched 2026-07-04. Bitmap-indexed revocation; no per-token live call; explicit design goal "the specification shall enable caching policies and offline support"; revocation observability is **eventual consistency on the order of the list refresh interval**, not real-time; the RP fetches the list from the published URI.
- **RFC 7009** ([datatracker.ietf.org/doc/html/rfc7009](https://datatracker.ietf.org/doc/html/rfc7009)) — fetched 2026-07-04. Exclusively an online revocation endpoint; no offline or partition-tolerant mechanism; explicitly acknowledges propagation delay without defining how to handle it.
- **SPIFFE standards** ([raw.githubusercontent.com/spiffe/spiffe/main/standards/SPIFFE.md](https://raw.githubusercontent.com/spiffe/spiffe/main/standards/SPIFFE.md)) — fetched 2026-07-04. Trust bundles are trust-domain-scoped, distributed via the Workload API, used for federation; JWT-SVID is a separately specified document. Offline cryptographic verification of a JWT-SVID with the issuer's public trust bundle is the canonical SPIFFE design — confidence **Medium-High** (the architecture is confirmed in the fetched text; the JWT-SVID verification prose specifically sits in the referenced sub-spec, not the fetched excerpt).
- **RFC 7519 (JWT) / RFC 7515 (JWS)** — standard `exp`, `scope`, `iss`, `sub` claims and signature verification. Durable canonical references; not re-fetched for this audit.

**Note on tool availability.** Web *search* was unavailable in this transport (sandbox configuration: `web_search`/`web_fetch` not wired to the active upstream without a local handler). Web *fetch* on canonical URLs succeeded. The audit therefore grounds claims in primary fetched sources where load-bearing and labels recall-based claims with their actual (lower) confidence. No claim here is presented as verified when it rests only on recall.

---

## Decomposition of H1

H1 is decomposed into twelve component claims. Each is audited as solved / partially solved / unsolved against existing standards, OSS, or straightforward composition. Verdicts and their confidence labels follow the project doctrine: confidence without evidence is forbidden; confidence without a stated change-condition is forbidden.

### Solved (Confidence: High unless noted)

| ID | Component claim | Why solved | Evidence | Change-condition |
|---|---|---|---|---|
| C1 | Cross-trust-domain JWT-SVID verification, offline, using the issuer's published trust bundle | JWT-SVID is a JWT signed by the trust domain's keys; the RP verifies with the bundle's public keys, no live call to domain A. Two independent SPIRE deployments + manual bundle exchange = standard SPIFFE deployment pattern. | SPIFFE standards (fetched); trust-bundle architecture confirmed. JWT-SVID verification prose confidence Medium-High. | Only if a SPIFFE revision requires online federation for JWT-SVID verification (it does not per current spec). |
| C2 | Delegation token *format* expressing "X acting on behalf of Y, scope S, until T" | RFC 8693 `act` claim is defined precisely for this (fetched verbatim). `scope` (RFC 6749) and `exp` (RFC 7519) are standard JWT claims. | RFC 8693 §4.1 (fetched). High. | Only if `act` claim is deprecated in a successor RFC (not in progress). |
| C3 | Scope and TTL constraints on the delegation | Standard JWT `scope` and `exp` claims; stock JWT libraries handle both. | RFC 7519 standard. High. | — |
| C6 | Replay rejection after TTL | JWT `exp` and `nbf`; clock-skew handling is stock-verifier behavior. | RFC 7519. High. | — |
| C7 | Tamper rejection | JWS signature verification over the token. Stock. | RFC 7515. High. | — |
| C9 | Verification latency < 100 ms end-to-end | JWT signature verification with cached keys is sub-millisecond at ordinary key sizes; the 100 ms budget clears by ~2–3 orders of magnitude. | Standard OSS verifier performance (go-jose, nimbus, etc.). High. | Only if a deliberately exotic signature scheme is mandated. |
| C8 | Verification window bounded above by Y-SVID TTL | This is a temporal configuration constraint (set delegation TTL < Y-SVID TTL), not a primitive. | First principles. High. | — |
| C12 | Ephemeral workload substrate | Ephemeral workloads with short-lived SVIDs are a core SPIFFE/SPIRE use case. T8 is an instance of C1, not a new primitive. | SPIFFE architecture (fetched). High. | — |
| C10 (for the solved parts) | Unmodified off-the-shelf verification libs suffice | Stock JWT/JWS verifiers handle C1, C2, C3, C6, C7, C9, C12. | OSS verifier ecosystems. High for these parts. | — |

### Partially solved (Confidence: Medium-High)

| ID | Component claim | What is solved | What is unsolved | Change-condition |
|---|---|---|---|---|
| C5 | Offline / partition-tolerant verification | Verification of an *already-issued, untampered, un-expired* token is fully offline (C1). | "**observe revocation**" across a partition that isolates the RP from the issuer *at the time of revocation* cannot be done by any push-then-verify composition — the revocation information has no channel to cross the partition within the verification window. This is an information-theoretic limit, not a technology gap. | Only if H1 is re-scoped to "revocation observable *eventually*, post-partition-recovery" — see §"Strengthen." |
| C11 | Multi-hop delegation: X → X' → Y | *Format* is solved — RFC 8693 explicitly supports nested `act` (fetched verbatim). | *Per-hop independent authorization* is **not** solved by RFC 8693: the spec states nested prior actors are "informational only," with only the outermost actor mattering for access-control. T7's sub-claim "X → X' verified, X' → Y verified, X → X' verified, each *independently*" implies per-hop authorization gates, which RFC 8693 does not provide. Independently revoking per-hop is also barred by C4 (below). | Only if a per-hop capability-token chain (e.g., macaroon/biscuit-style attenuation or a new capability scheme) is admitted — that is not a "straightforward composition of existing primitives" and not RFC 8693. |

### Genuinely unsolved (Confidence: Medium-High that no standardized composition solves it)

**C4 — Offline, fast, partition-tolerant, *independent* revocation of a delegation token, with no live call to domain A and no third-party broker.**

This is the single load-bearing unsolved component of H1. Three candidate primitives were considered and falsified against H1's no-live-call/no-broker clauses:

1. **RFC 7009 (OAuth Token Revocation).** Falsified. It is exclusively an *online* endpoint (fetched verbatim). Calling it is a live call to the authorization server — forbidden by H1. It defines no offline, no list, no partition-tolerant mechanism. Out.

2. **OAuth Status List (`draft-ietf-oauth-status-list`).** Falsified on **two** independent grounds:
   - *Observability latency:* revocation is eventual consistency on the order of the list's refresh interval (fetched verbatim), not real-time. T4's test demands revocation set at T₀+5s be observable at T₀+7s — a 2-second requirement. A status list whose refresh cadence is < 2s is operationally absurd at any scale (re-signing and republishing a compressed herd-privacy list every second defeats the batching purpose). The mechanism's whole point is herded, periodic refresh.
   - *No-live-call clause:* the status list is fetched from a URI (the status issuer's published endpoint). If domain A is the status issuer, the periodic fetch *is* a live call to a service in domain A — forbidden. Mirroring to a CDN may evade the "domain A" clause but the CDN is a third party, and H1 forbids "any third-party broker." Whether a passive cache counts as a "broker" is ambiguous (see §"Strengthen," S3); under a strict reading it does. Out, modulo that ambiguity.

3. **Short-TTL-as-revocation (let the delegation expire).** Falsified by H1's own configuration. The Level 2 protocol locks delegation TTL = 30 s. T4 requires a revocation set at T₀+5s to be observable at T₀+7s — well inside the 30 s window. Short-TTL delivers *expiry*, not *revocation*, and gives no observability inside a live window. Out.

What is **not** falsified (because no standardized composition exists to test):
- **Push-revocation** (domain A pushes revocation events to the RP's inbound listener): no standardized protocol; no OSS that interoperates with SPIFFE/RFC 8693 RP code; would require new composition code on the RP, which may or may not satisfy H1's "unmodified off-the-shelf verification libraries" clause S7 (the JWT verifier stays unmodified; the *listener* is new — ambiguous).
- **Cryptographic-accumulator-based revocation** (Camenisch-Lysyanskaya-style, dynamic accumulators): academic, no standardized OSS implementation interoperating with SPIFFE/JWT/SPIRE; uses heavyweight algebra which carries its own verification-latency overhead.
- **Capability-token attenuation-as-revocation** (macaroon/biscuit): yields *narrowing*, not *revocation*; cannot nullify a hold arbitrarily without an external coordination signal.

**Confidence that no standardized composition solves C4: Medium-High.** Evidence basis: fetched RFC 7009 + fetched OAuth Status List spec + first-principles on short-TTL. Push/accumulator remain *in-principle* candidates with no current standardized implementation; whether either composes within H1's library and partition constraints is exactly what an empirical test would settle.

**Change-condition that would re-open C4:** (i) a standardized push-revocation or accumulator-based protocol achieves IETF RFC status and interoperates with SPIFFE; (ii) H1's no-live-call clause is re-interpreted to admit periodic cached pulls (relaxing the no-broker clause for passive caches); (iii) the revocation-observability latency bound is loosened to "eventual consistency within the delegation TTL." Any of (i)/(ii)/(iii) re-opens C4 and re-runs this audit.

---

## Assumptions in H1 that are actually unnecessary

Stated as a feasibility audit — "unnecessary" means *not load-bearing to the falsifier*, not "wrong to test inside Level 2." Level 2 may keep them as rigor; the feasibility question does not depend on them.

| # | Unnecessary assumption | Why not load-bearing to the falsifier |
|---|---|---|
| U1 | The "two independent SPIRE deployments with federation disabled" substrate precondition | This is an experimental-design choice that guarantees no federation back-channel smuggles delegated trust back in. It hardens Level 2's internal validity. It is not a property of *whether the delegation primitive exists in principle*; the feasibility of cross-domain delegation does not turn on federation being disabled. For the audit, this is instrumentation, not hypothesis. |
| U2 | C12 / ephemeral-workload substrate (T8) as a *separate* falsifier | Once C1/C2/C3 are solved, ephemeral workloads are handled natively by SPIFFE. T8 is an instance of C1, not a novel falsifier of H1. Keeping it as a Level 2 test is fine; it carries no independent weight for the *feasibility* question. |
| U3 | C11 multi-hop (T7) as a *separate* falsifier | The multi-hop *format* is solved (RFC 8693 nesting, fetched verbatim). The only sub-claim of T7 not already solved is "per-hop independent authorization/revocation," which is C4 × N — i.e., the C4 problem again, not a second distinct falsifier. T7 inherits its verdict from C4. |
| U4 | The explicit < 100 ms latency threshold (C9) as a falsifier | Stock JWT verification clears 100 ms by two orders of magnitude. This is a performance *budget*, not a feasibility *question*. Retaining it in Level 2 is harmless; it is not the falsifier. |

## Assumptions that should be strengthened

These are ambiguities in H1 whose resolution materially changes C4's (and therefore H1's) feasibility. Resolving them is a *descriptive* scope act, not a Level 2 redesign — it is the feasibility gate defining what H1 actually claims.

| # | Ambiguity in H1 | Why it matters | Strengthen to |
|---|---|---|---|
| S1 | Revocation-observability latency is unspecified as a *number*. H1 says "observe revocation within a verification window strictly bounded above by the TTL of the underlying Y SVID" — that bounds the *window*, not the *latency*. T4 then implies a 2-second latency. These are conflated. | Whether "revocation observability" means **2 s** (T4-strict) or **eventual within the window** decides whether the OAuth Status List primitive is admissible (2 s → inadmissible; within-window → admissible modulo refresh cadence). This is a feasibility *branch point* hiding inside H1, not a minor wording issue. | State the latency bound explicitly: "revocation observable within R seconds, R < delegation TTL," and choose R. |
| S2 | "Forbidden from making a live call to any service in domain A or to any third-party broker" is binary but does not address the periodic-cached-pull case (Status List's whole design). | Decides whether the OAuth Status List — the closest existing primitive — is admissible at all. If the periodic cached pull from domain A's published URI is "a live call to a service in domain A," the Status List is excluded by construction and C4 must be solved by push/accumulator (no standard). If the cached pull is admissible, C4 collapses to a refresh-cadence parameter and the problem shrinks. | Either (a) admit periodic cached pulls of *signed, integrity-protected* artifacts (and bound the staleness), or (b) forbid all RP-initiated fetches from domain A and require push. |
| S3 | "Any third-party broker" does not distinguish an *active trust broker* (making decisions) from a *passive cache* (CDN serving a signed blob). | A CDN serving a Status List is arguably a passive cache; under (a) above, it is the natural delivery channel. Under a strict reading, it may be a "third-party broker." This ambiguity alone can change C4 from solvable-by-composition to unsolved. | Define "broker" as an entity that makes or vouches for trust decisions, explicitly excluding passive signed-blob caches that perform no decision. |
| S4 | The revocation-**during-partition** case (T4 ∩ T5 combined: revocation set on the issuer side *while* the RP is partitioned from the issuer) is unspecified. | This is an information-theoretic limit. No technology — present or future — can let an RP observe a revocation performed on the other side of a partition that isolated the RP from the issuer at the time of revocation, within a window shorter than partition-recovery. Strictly read, T4 ∩ T5 is a *logical impossibility*, not a technology gap; running Level 2 against that combined reading would falsify H1 by construction and tell us nothing. | Either explicitly bound revocation observability to "eventually consistent upon partition recovery, within P seconds of recovery," OR explicitly state H1 does not require observability of revocations performed during a partition that isolated the RP from the issuer at revocation time. |
| S5 | T7's "X → X' verified, X' → Y verified, X → X' verified, **each independently**" vs RFC 8693's "nested prior actors informational only." | RFC 8693 nested `act` is a *provenance trail*, not a *per-hop authorization gate chain*. If T7 requires per-hop independent authorization, RFC 8693 does not provide it and the experiment would require a non-RFC-8693 capability layer (macaroon/biscuit attenuation or new capability scheme) — which is not "straightforward composition of existing primitives." | Decide whether T7 requires per-hop independent *authorization* (capability-token chain, unsolved by standard) or an auditable *provenance trail* (RFC 8693 nesting, solved). Strengthen the wording accordingly. |

---

## Verdict — is Level 2 justified?

### **B. Partially. A lightweight feasibility spike should precede Level 2.**

### Why not A ("no, existing technology satisfies H1")

C4 is genuinely unsolved by any standardized composition under H1's strict no-live-call / no-broker clauses. RFC 7009 is online-only; the OAuth Status List is eventual-consistency-with-periodic-fetch-from-the-issuer and is excluded under a strict reading; short-TTL does not deliver in-window revocation. No existing technology or straightforward composition fully satisfies H1 as written. The audit therefore rejects "A."

### Why not C ("yes, Level 2 immediately")

The audit does establish that a load-bearing sub-claim (C4) is unsolved — which under a coarse reading would justify the experiment. But three points argue against running the *full* Level 2 protocol immediately:

1. **~80% of H1 is already solved by standardized composition.** C1, C2, C3, C6, C7, C8, C9, C11-format, C12, and C10-for-the-solved-parts are all solvable with stock libraries. The full Level 2 protocol — nine test cases, eight success criteria, two-engineer adversary-blind apparatus, out-of-band sniffer, frozen substrate — re-tests all of them. That apparatus is well-designed *for certifying a primitive*, but its load-bearing falsification target collapses to **one** open question, C4. Spending the full protocol to re-confirm the solved ~80% is misallocated experimental capital.

2. **C4 itself contains a cheap-or-expensive branch point.** The sharp audit finding is that *whether C4 is solvable by any near-standard composition* (push-revocation layer over stock JWT verifier; OAuth Status List with S2/S3 relaxed and refresh ≤ delegation TTL; accumulator-based) is exactly what a small, targeted spike would settle — *without* the full substrate, sniffer, adversary-blind apparatus, or nine-case suite. A 2–3 day spike on C4 alone either (a) finds a composition that satisfies T4 → in which case H1 collapses to "SPIFFE + a thin revocation layer," the novelty claim narrows sharply, and Level 2 in its current breadth is not needed, or (b) confirms the gap and **narrows** Level 2 to C4 + the genuinely-novel remainder, making the expensive experiment both better-targeted and cheaper.

3. **The audit surfaced a possible *logical* tension in H1, not just a technology gap.** S4 (T4 ∩ T5: observe revocation set during a partition that isolated the RP from the issuer) is, under a strict combined reading, an information-theoretic impossibility — no technology can satisfy it. Running Level 2 against that combined reading would falsify H1 *by construction*, telling us nothing about P5. The feasibility gate must first resolve S4 before any Level 2 verdict is interpretable. This is a definitional act (defining what H1 actually claims), not a Level 2 redesign.

Therefore the cheap spike is the correct gate: it either dissolves the expensive experiment, narrows it, or exposes a needed scope correction — each outcome strictly better than running Level 2 blind against a hypothesis whose load-bearing clause is partly ambiguous and partly already solvable.

### What the lightweight feasibility spike should target (not a Level 2 redesign — a smaller, prior gate)

The spike answers exactly one question and uses exactly the per-component audit above as its scope contract:

> **Spike question (C4-only).** On a SPIFFE substrate with RFC 8693 tokens, is there *any* composition of existing standardized primitives (stock JWT verifier + a thin, non-patching revocation layer) that delivers revocation observability within R seconds of revocation — for R as chosen in S1 — with no live call from the RP to domain A and no third-party decision-making broker, *under the S2/S3/S4 read of H1*?

Candidate compositions to attempt, cheapest first:
- OAuth Status List with explicit refresh cadence set to R (resolves only if S2/S3 admit cached pulls; otherwise excluded by construction).
- Push-revocation: domain A signs-and-pushes a small revocation event to the RP's inbound listener on a best-effort channel; the RP treats absence-of-push as "still valid" and is partition-tolerant for verification of already-issued tokens.
- Accumulator-based revocation: pre-loaded dynamic accumulator updated periodically; the RP verifies membership/non-membership locally.

**Spike outcome → Level 2 disposition:**

| Spike outcome | Disposition for Level 2 |
|---|---|
| A working composition satisfies the spike's C4 criterion (with R as set in S1) | H1 collapses to "SPIFFE + a thin revocation layer over a standardized primitive." **Level 2 in its frozen breadth is not justified** — most of it re-certifies solved components. Re-scope decision is the founder's (this audit does not redesign Level 2). |
| No composition works, **and** C4 is confirmed to be a *technology gap* (not a logical impossibility) under the S1–S5strengthened reading | **Level 2 is justified**, and can be narrowed to C4 + C11's per-hop sub-claim (which is C4 × N), eliminating the re-certification of solved components. |
| C4 is confirmed to be a *logical impossibility* under the S4 combined reading | **Level 2 is not justified as written** — it would falsify H1 by construction. The gate returns to the founder to re-scope H1 (per S4) before any experiment. |
| C4 is unsolved by composition **and** the S1–S5 scope ambiguities cannot be resolved without empirical input | **Level 2 is justified as-is**, because the desk audit cannot reduce the uncertainty further — the empirical test is the only remaining instrument. (This is the branch that would correspond to verdict "C," reached only *after* the spike fails to settle the matter.) |

### Stop rule

This document is the Level 0/1 feasibility-gate verdict. It does not redesign Level 2, write code, design architecture, or generate implementation. Its single output is verdict **B** and the C4-targeted spike question above. The next action belongs to the founder: approve the spike scope, resolve S1–S5 as semantic acts (defining what H1 actually claims), or redirect.

---

## Confidence (per project doctrine: evidence + change-condition required)

- **Verdict (B): High.** Evidence: the audit demonstrates a specific, cheaper C4-targeted spike that can dissolve, narrow, or correct the Level 2 question before it is run; the audit also demonstrates ~80% of H1 is already solved by standardized composition. Change-condition: if the C4 spike cannot be scoped tightly (e.g., requires the full Level 2 substrate to be meaningful), revisit whether the spike retains its cost advantage over Level 2 — that would shift the verdict toward C.
- **C4 genuinely unsolved by standardized composition: Medium-High.** Evidence: fetched RFC 7009, fetched OAuth Status List spec, first-principles on short-TTL. Change-condition: any of S1–S4 re-opens it (see above).
- **C1, C2, C3, C6, C7, C9, C11-format, C12 solved: High** (C1's specific JWT-SVID verification prose: Medium-High, awaiting the referenced sub-spec).
- **C5 partition-tolerance of verification of valid tokens: High.**
- **C5 observability of revocation performed during a partition (S4 combined reading): High** — that this is an information-theoretic limit (not a technology gap) is a first-principles conclusion, deducible without external sources.
- **Predictive confidence in the spike outcome: None.** Per project doctrine (confidence without evidence is forbidden and refused-by-default for forecasts), no claim is made here about *which* spike outcome will occur.

## Provenance

Source of truth: `P5_FALSIFICATION_EXPERIMENT.md` §1 (H1, frozen; not modified). Secondary: `TECHNICAL_VALIDATION.md` §P5 (approved), `FOUNDER_PROBLEM_FIT.md` §P5 §5–7 (approved), `RESEARCH_PROGRAM.md` §4 Wave 4 / Task 4 (approved). External standards grounded via fetched canonical URLs listed in the Evidence Basis section.

## Sources

- [RFC 8693 — Token Exchange](https://datatracker.ietf.org/doc/html/rfc8693) (fetched 2026-07-04; `act` claim, nested delegation, prior-actors-informational confirmations)
- [draft-ietf-oauth-status-list-10 — Status List Token](https://datatracker.ietf.org/doc/html/draft-ietf-oauth-status-list-10) (fetched 2026-07-04; bitmap revocation, eventual consistency, offline-support design goal)
- [RFC 7009 — OAuth 2.0 Token Revocation](https://datatracker.ietf.org/doc/html/rfc7009) (fetched 2026-07-04; online endpoint only, no offline mechanism)
- [SPIFFE standards](https://raw.githubusercontent.com/spiffe/spiffe/main/standards/SPIFFE.md) (fetched 2026-07-04; trust-bundle/trust-domain architecture confirmed)
- RFC 7519 (JWT) and RFC 7515 (JWS) — referenced as durable canonical standards; not re-fetched for this audit.

<!-- checkpoint: context(threat-model-scenarios): refine threat model scenarios -->

<!-- checkpoint: feat(verify): implement truststore backend -->

<!-- checkpoint: chore(lab): audit lab environment topology (#154) -->
