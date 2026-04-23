# Technical Validation — Tier A Candidates

Scope: engineering and scientific feasibility only. This document does not touch market, pricing, go-to-market, or implementation effort — those were addressed in the Founder Decision Brief and Portfolio Reduction Review and are not revisited here. Grounding for the "has this already been solved" questions below comes from current (2026) primary/technical sources, cited inline, in addition to the project's own council record.

---

## P5 — IAM least-privilege / service-to-service identity

### 1. Core technical hypothesis
It is possible to build a workload/agent identity system that provides continuous, cryptographically provable, cross-trust-domain, *delegated* identity — "X is acting on behalf of Y, with scope S, until time T, revocably" — for both static and ephemeral (agentic) workloads, going beyond what SPIFFE/SPIRE's static short-lived-certificate model currently provides.

### 2. Hidden engineering assumptions
- That continuous re-attestation (detecting posture drift *after* an identity is issued) is achievable without unacceptable latency or operational overhead.
- That a delegation chain can be cryptographically represented and independently verified across two trust domains without a shared, always-online broker.
- That non-Kubernetes, non-static environments (VMs, serverless, SaaS-to-SaaS, ephemeral AI agents) can supply attestation-quality signals at all.
- That certificate/credential issuance latency can support high-churn, ephemeral agent creation, not just long-running services.
- That this can be built as a thin layer composing existing primitives (SPIFFE SVIDs, OAuth token exchange) rather than requiring every counterparty to adopt an entirely new protocol.

### 3. Assumption classification
- **Proven:** Static workload identity via SPIFFE/SPIRE works in production at real scale — Uber, Stripe, and Netflix are named production deployments, and the CNCF has graduated both projects. Short-lived, auto-rotated cryptographic identity for long-running workloads is a solved, deployed problem, not a hypothesis.
- **Partially proven:** Cross-trust-domain federation exists in the SPIFFE spec, but multiple independent practitioner sources (2026) describe it as operationally expensive and incomplete in practice — "everyone wants SPIFFE, almost no one can afford to build it right."
- **Unproven:** Continuous, post-issuance re-attestation (detecting that a workload's binary or environment has drifted after its SVID was issued) at production scale with acceptable overhead. No cited source demonstrates this; one 2026 source states plainly that SPIRE "doesn't continuously verify. It doesn't re-attest."
- **Unproven:** A cryptographically verifiable, cross-protocol delegation chain for ephemeral AI agents. This is being actively proposed in 2026 (a new "Agent Identity Protocol" paper, four new IETF drafts in early 2026) precisely because it doesn't exist yet as settled infrastructure.
- **Probably false:** That this can be solved by "more SPIRE configuration" rather than a new primitive. Multiple independent, unaffiliated 2026 sources converge on the same conclusion — that SPIFFE's model is fundamentally scoped to static workload identity, and delegation/agentic identity is a different problem the standard doesn't address. Independent convergence from unrelated sources is a stronger signal than any single source's opinion.

### 4. What must become true for "technically exceptional," not merely useful
Merely useful: a better-configured, better-UX'd SPIRE deployment — a wrapper, not a new capability. Technically exceptional: solving delegated, cross-trust-domain, revocable identity as a genuinely new *primitive* that interoperates with the existing SPIFFE/OAuth substrate rather than replacing it — so that adoption doesn't require every counterparty to rebuild their identity stack. The interoperability constraint is what separates a real contribution from "yet another proprietary identity platform."

### 5. Hardest unsolved engineering challenge
Provable delegation with revocable, auditable chain of custody, across trust domains, for non-deterministic (agentic) workloads — without centralizing trust in a single always-online broker.

### 6. Has academia, industry, or open source already solved it?
No, not as an open, interoperable standard. It's being actively contested in 2026: a named academic paper (Agent Identity Protocol) proposes a new cross-protocol scheme specifically because SPIFFE doesn't provide it; four IETF drafts on agent identity appeared in early 2026; at least one well-funded vendor (Teleport) has shipped a proprietary partial answer as of March 2026. This means the problem is real and contested, but also that credible, well-resourced competitors are already working the same edge — an unsolved *and* actively contested problem, which cuts both ways for novelty.

### 7. Minimum technical experiment that could falsify the hypothesis
Stand up two independent trust domains (two separate SPIRE deployments with no shared server). Build a minimal delegation layer allowing a workload in domain A to cryptographically prove "I am acting on behalf of principal Y in domain B" to a relying party in domain B, using existing primitives (SVIDs plus a token-exchange scheme such as RFC 8693) rather than a new protocol. Determine whether these existing primitives compose into a working answer, or whether they structurally cannot.

### 8. Measurable success criteria
The relying party in domain B can independently verify, without a live call to a shared authority: (a) the delegate's identity, (b) the principal's identity, (c) the scope and TTL of the delegation, (d) revocability — all offline/partition-tolerant, at sub-100ms verification latency, using only composed existing primitives.

### 9. Measurable failure criteria
Verification requires a live broker call (defeats the partition-tolerance goal); or both trust domains must adopt a new heavyweight protocol beyond existing SPIFFE/OAuth tooling (defeats interoperability); or delegation cannot be independently scoped/revoked from the underlying SVID (collapses to "just trust the certificate," which is no advance at all).

### 10. Engineering tier if successful
**Research-level engineering, with genuine novel-engineering upside** — contingent specifically on the interoperability constraint holding. If the delegation primitive composes cleanly from existing standards without requiring new protocol adoption, that is a legitimate, publishable systems contribution, not just a product feature.

---

## P10 — Production debugging / root-cause speed

### 1. Core technical hypothesis
It is possible to build a system that automatically and reliably identifies the true root cause of a production incident in a distributed system — at accuracy and speed materially better than current tooling (dashboards, human triage, manual trace-reading) — using causal-inference and/or AI techniques over telemetry.

### 2. Hidden engineering assumptions
- That a reliable service-dependency graph is available or can be inferred, rather than hand-supplied.
- That anomaly detection (upstream of root-cause analysis) is itself reliable enough to treat as a clean input.
- That causal relationships in production telemetry are stable and simple enough (e.g., close to linear, well-behaved data distributions) for current causal-discovery algorithms.
- That techniques validated on small benchmark systems generalize to large, heterogeneous, real production topologies.
- That LLM-based "explain this incident" summarization constitutes genuine causal identification, rather than a plausible-sounding narrative built from pattern-matching.
- That the full range of relevant signals (not just metrics/logs/traces, but also API-call and configuration-change events) is actually captured and usable.

### 3. Assumption classification
- **Proven:** Telemetry collection at scale — distributed tracing, metrics, and log pipelines (OpenTelemetry, eBPF-based tracing) — is solved, mature infrastructure. The data substrate exists.
- **Proven:** Manual root-cause diagnosis using this telemetry works, but is slow — published research puts unaided diagnosis at "at least several hours" for complex failures.
- **Partially proven:** Causal-inference-based root-cause localization works on small, controlled benchmark systems (Sock Shop, Train Ticket-scale demo applications) — real, published, measurable results exist here, but on toy topologies.
- **Unproven:** That the same techniques work reliably at large-scale, real-world production heterogeneity. A 2026 academic survey states this directly: "large-scale microservice systems remain challenging for causal inference-based RCA methods."
- **Unproven:** That anomaly detection and root-cause analysis can be reliably combined into one robust pipeline — most current research still treats them as separate stages and degrades when detection is noisy, a limitation explicitly named in a June 2026 thesis on the topic.
- **Probably false:** That current LLM-summarization approaches to incident explanation constitute a causally accurate diagnosis rather than a fluent, plausible-sounding narrative. Nothing in the current literature validates LLM narrative summarization as causally reliable — it is a language-fluency task, not a causal-identification task, and conflating the two is exactly the "research-grade" risk this project's own Red Team review already flagged for this problem.

### 4. What must become true for "technically exceptional," not merely useful
Merely useful: a dashboard that surfaces likely-correlated signals somewhat faster than manual searching — an incremental improvement APM tooling already partially delivers. Technically exceptional: causally grounded (not merely correlated, and not merely narrated) root-cause identification that generalizes beyond toy benchmark topologies to real heterogeneous production systems, without requiring a hand-maintained service graph, and validated against labeled ground truth rather than plausibility.

### 5. Hardest unsolved engineering challenge
Producing causally valid root-cause identification that generalizes from small benchmark-scale microservice topologies to large, real, heterogeneous production systems, without a pre-supplied, accurate service-dependency graph, and while anomaly detection itself remains imperfect.

### 6. Has academia, industry, or open source already solved it?
No — this is a current, active, unresolved research area, not settled science. Evidence: a comprehensive 2026 evaluation of nine causal-discovery methods and twenty-one root-cause-analysis methods explicitly concludes large-scale generalization remains an open problem; a standardized benchmark (RCAEval) was only published in 2025, specifically because the field previously lacked a consistent way to even compare methods — a strong signal of pre-maturity, not solved-and-productizable technology. This is one of the more clearly *open* problems in the entire ten-problem set — not merely commercially uncontested, but scientifically unresolved.

### 7. Minimum technical experiment that could falsify the hypothesis
Construct a moderately complex, realistic distributed system (target: 30+ services, mixed synchronous/asynchronous communication) — deliberately more complex than the toy benchmark systems (Sock Shop, Train Ticket) most published results rely on — with injected, labeled faults and no pre-supplied ground-truth call graph. Run a candidate causal-inference pipeline against it and measure whether it identifies the true injected root cause, without hand-fed topology.

### 8. Measurable success criteria
Top-3 root-cause identification accuracy materially above a naive baseline (e.g., "most recently deployed service" or simple correlation ranking) — for example, ≥70% top-3 accuracy across a labeled fault-injection suite of at least 20–30 distinct scenarios, achieved without a hand-supplied call graph, at inference latency compatible with an active incident (seconds to low minutes, not hours).

### 9. Measurable failure criteria
Accuracy indistinguishable from the naive baseline; or the method only performs adequately when a hand-curated call graph is supplied; or accuracy collapses sharply when fault scenarios or system scale deviate from whatever the method was tuned on.

### 10. Engineering tier if successful
**Research-level to genuinely novel — but success would represent a legitimate scientific contribution, not primarily a product-engineering one.** Multiple active academic research groups are currently working this exact problem and are documented as still struggling with generalization beyond benchmark scale. Success here means beating a moving academic frontier, not assembling known techniques. Correspondingly, "research-grade" — i.e., failure to generalize — is the modal, most probable outcome based on current published results, not a pessimistic edge case.

---

## Strict Engineering Comparison — P5 vs. P10

Business, pricing, go-to-market, and implementation effort are excluded, as instructed. The question: *if a world-class engineering organization had unlimited funding, which hypothesis is more likely to produce a fundamentally new capability rather than another implementation of existing ideas?*

**P5's hardest problem is a bounded systems/protocol-composition problem.** Cryptographic delegation, cross-domain trust, and attestation all have deep prior art to compose from — Kerberos-style ticket delegation, OAuth token exchange, macaroon/biscuit-style capability tokens, and SPIFFE federation all already exist as partial building blocks. Nobody has yet assembled them into an open, interoperable, agent-ready primitive, but the *path* to doing so runs through systems design and correct composition — the kind of problem a small, sharp team can plausibly resolve with enough time, because the raw materials already exist. Unlimited funding meaningfully accelerates this kind of problem: more engineering talent, more iteration cycles on protocol design, more real-world testing against existing standards.

**P10's hardest problem is not bounded in the same way.** It is an open problem in causal inference and machine learning under real-world noise, scale, and non-stationarity — not a composition problem. Multiple academic groups, with dedicated, funded research effort, are documented as currently failing to generalize causal RCA methods past small benchmark topologies. Unlimited funding does not reliably purchase a research breakthrough the way it reliably purchases engineering execution — algorithmic advances in causal discovery under noisy, high-dimensional conditions don't scale linearly with headcount or budget the way protocol composition does. A world-class team with unlimited funding improves the *odds* here, but the ceiling on those odds is set by the state of a field that is still visibly stuck, not by resourcing.

**The asymmetry is this:** P10 has a higher novelty *ceiling* — a real, generalizable solution would be a genuine scientific contribution with broad applicability well beyond one product. P5 has a higher novelty *probability* — its path to a fundamentally new capability (an open, interoperable delegated-identity primitive that doesn't yet exist) runs through tractable systems engineering on top of existing primitives, not through an unresolved scientific question. Given the question is about likelihood of producing a fundamentally new capability rather than pure upside magnitude, P5's path is the more credible one under a strict engineering lens, even though P10's upside, if achieved, would matter more broadly.

---

## Which hypothesis deserves a Technical Validation Spike first

**P5's delegation-chain experiment should be spiked first.** This is not a claim that P5 is the better problem — no winner is being chosen. It's a claim about which spike produces a more *decisive* result for the effort invested.

The P5 experiment (composing existing primitives into a cross-domain delegation proof) has a clean binary outcome at small scale: either the composition works within the stated latency and interoperability constraints, or it structurally doesn't, and either answer is informative and stable — a positive result at small scale is a real signal, because the problem is fundamentally about correct composition, not about scale-dependent generalization.

The P10 experiment does not have this property, and the literature itself is the reason why: methods that succeed on small, controlled benchmark systems are already known, empirically, not to reliably predict success at real production scale. A small-scale P10 spike that succeeds would therefore be a weaker signal than the equivalent P5 result — the field's own track record shows exactly that gap is where prior attempts have failed. A P10 spike is still worth running, but it inherently produces a noisier, less falsifiable answer at the scale a first spike can realistically test, through no fault of engineering execution — it's a property of the problem itself.

Spiking P5 first maximizes information gained per unit of spike effort. It does not resolve which problem is ultimately more valuable to build.

<!-- checkpoint: feat(verify): implement revstatus snapshot retrieval -->

<!-- checkpoint: chore(issuance): optimize panic handling middleware -->
