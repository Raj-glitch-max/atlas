# Product Thesis — Fundamental Primitive Analysis

Scope: this document ignores UI, APIs, features, pricing, and business entirely. The question is narrower and harder than "is this a good product": does this direction introduce something other systems could be built *on top of*, the way a filesystem, a container image, a Git commit, or a trace span are built on top of — or is it, honestly, a well-executed feature sitting on primitives that already exist.

---

## P5 — Delegated, cross-domain workload/agent identity

### 1. Existing primitive used by industry today
**The static credential** — a certificate (X.509/SPIFFE SVID) or bearer token (OAuth access token, API key) that asserts "I am X" at the moment it's checked. This is the load-bearing primitive underneath essentially all current identity infrastructure, including SPIFFE/SPIRE itself.

### 2. Why it is fundamentally insufficient
A static credential conflates *authentication* (who are you) with *authorization scope* (what can you do, on whose behalf, until when) into a single opaque object. It has no native representation of a delegation chain — "agent A, invoked by service B, invoked by human C" has to be reconstructed after the fact from logs, not verified structurally at the time of the request. Revocation is coarse: revoking the credential kills everything downstream of it, not one specific delegated capability. This primitive was designed for relatively static, long-lived identities, not high-churn, ephemeral, cross-organizational, non-human agents.

### 3. What new primitive could exist
Here's where the analysis has to get more honest than the previous framing allowed. **The attenuable, offline-verifiable delegation capability already exists as a primitive** — Google's macaroons (2014) and the actively maintained Biscuit token format are real, deployed instances of exactly this idea: a cryptographic object that can be narrowed (attenuated), chained, and verified without a broker call. That primitive is not new. What doesn't exist is a **standardized binding between workload identity (SPIFFE) and attenuable capability tokens, portable across independently-operated trust domains with no shared authority** — i.e., not a new primitive in isolation, but a new *load-bearing combination* of two existing primitives that, combined and standardized, would function as a new unit other systems build against.

This is worth being precise about, because it's exactly the pattern behind two of the reference examples in this thesis: Docker didn't invent chroot, cgroups, or union filesystems — it combined them into a portable, addressable image format that became the primitive the entire container ecosystem builds on. Git didn't invent content-addressed Merkle trees — it combined them into a commit graph that became the primitive version control tooling builds on. The test for P5, by the same standard, is not "did you invent delegation" (you didn't, and shouldn't claim to) — it's "does the combined, standardized object become something other systems adopt as a substrate."

### 4. What becomes possible because of that primitive
If the combination becomes a real, adopted standard: genuinely least-privilege AI agent systems, where an agent carries a narrowly-scoped, time-bound, independently-revocable capability instead of a blanket API key or a human-equivalent credential. Cross-organizational workflows that don't require a shared identity provider or federation setup. Structural, cryptographic audit trails instead of log-reconstructed ones. Failure containment that revokes one delegated capability without touching the underlying identity or restarting a workload.

### 5. What products could eventually be built on top of it
Agent-to-agent authorization platforms. Zero-trust API gateways that natively understand delegation chains instead of only bearer tokens. Audit and compliance tooling that consumes structural proof instead of reconstructed intent. Cross-organization B2B integration tooling that skips federation setup entirely. None of these are being proposed here — they're offered only to test whether the primitive is generative, which is the actual question this thesis has to answer.

### 6. Verdict — primitive or feature?
**Conditionally a primitive, and the condition is specific and falsifiable.** It becomes a real primitive if the standardized identity-plus-capability binding is adopted and built against by parties other than its creator — the same bar Docker's image format and Git's object format both cleared. It remains merely a well-executed feature if it ends up as "SPIRE plus a delegation extension only one system understands," which is the exact failure mode the Technical Validation already flagged as the interoperability risk. This thesis does not get to assume the condition is met — that's an empirical question the Technical Validation's own falsification experiment (Task 4) is positioned to answer, not this document.

---

## P10 — Causal root-cause identification in distributed systems

### 1. Existing primitive used by industry today
**The trace span and the metric time series** — OpenTelemetry's standardized units for "this unit of work happened, took this long, had this parent" and "this value was measured at this time." Root-cause analysis today is built as an application layer on top of these — dashboards, correlation heuristics, and increasingly LLM narrative summarization — not as a new primitive.

### 2. Why it is fundamentally insufficient
Spans and metrics encode temporal and call-structure relationships, not causal ones. A span can tell you A happened before B and B was nested inside A's context — OpenTelemetry's span links even capture *direct* triggering relationships in async messaging. What no current primitive captures is an **inferred causal relationship between two events that have no direct call or messaging relationship at all** — the canonical incident-response case where service A's failure was actually caused by resource contention with unrelated service B, and nothing in the trace graph connects them because they never called each other. Every RCA tool has to reconstruct this from scratch, per system, usually requiring a hand-fed dependency graph, because the substrate itself has no native way to represent it.

### 3. What new primitive could exist
A standardized, portable **inferred causal edge** — a structured object distinct from a span or a trace link, representing "event/state-change X had a measured causal influence on event/state-change Y, with confidence C, based on evidence E," specifically covering relationships that aren't already visible in the call graph. This is a sharper and more defensible claim than "a causal graph" in general, because span links already partially cover directly-connected causation — the genuinely missing primitive is the cross-cutting, non-call-graph-connected case.

The precedent here is OpenTelemetry itself: it didn't invent distributed tracing (Dapper and Zipkin predated it by years) — it standardized something that existed in ad hoc, vendor-specific silos into a portable format every tool could consume. The honest version of this thesis is that the *shape* of a causal-edge primitive is a coherent, precedented idea. What's not established — and this is where this document has to be more skeptical than optimistic — is whether it can actually be **populated with trustworthy data at all**, at real-world scale. The Technical Validation already concluded that reliably inferring this kind of causal structure beyond toy benchmarks is an open research problem, not an engineering-standardization problem. A primitive that can be defined but not reliably filled in is not yet a primitive — it's a hypothesis about one.

### 4. What becomes possible because of that primitive
If it becomes reliably populatable: root-cause tools that consume a standard causal substrate instead of each re-deriving causal inference from raw telemetry — the same leap OpenTelemetry produced for tracing. Genuinely comparable, benchmarkable RCA research, since there'd be a standard object to evaluate against instead of ad hoc per-system reconstructions. Causal debugging across organizational boundaries, if causal edges are emitted and shared as structured data rather than reverse-engineered after an incident.

### 5. What products could eventually be built on top of it
RCA and AIOps tooling that doesn't each re-solve causal inference independently. Causally-aware chaos-engineering and fault-injection tooling. Causally-aware cost-attribution tooling — notably, this would eventually touch the same territory as P4, one of the problems already rejected in the Portfolio Reduction Review, though that is not a reason to revisit that decision here. Again, these are offered only to test generativity, not proposed.

### 6. Verdict — primitive or feature?
**A genuine primitive-shaped idea, gated by an unresolved feasibility question rather than an unresolved standardization question.** This is a meaningfully different verdict than P5's. P5's open question is "will anyone adopt the standard" — an adoption/execution risk. P10's open question is "can this be reliably computed at all" — a scientific risk that sits upstream of adoption entirely. It would be dishonest to call this a confirmed primitive; it would be equally dishonest to call it merely a feature, because the concept, if it worked, would be structurally novel in exactly the way the thesis is testing for. The correct statement is: **this is a candidate primitive whose existence as a real, buildable thing has not yet been established** — which is the same conclusion the Technical Validation reached from a different angle, now confirmed from the primitive-design angle as well.

---

## Cross-cutting observation

Neither candidate clears the bar for "confirmed platform primitive" today. Both clear the bar for "a primitive worth testing," but for different reasons and with different failure modes: P5's combination already has all its components proven and working independently (SPIFFE, macaroons/biscuits both exist and function) — the risk is purely whether the combination gets adopted as a standard by anyone besides its creator. P10's core object doesn't yet have a proven way to be populated with trustworthy data at all — the risk is upstream, in whether the concept is even computable at scale, independent of adoption. That is the same asymmetry the Technical Validation identified — higher probability of success for P5, higher scientific novelty ceiling for P10 — now visible again at the level of primitive design rather than engineering feasibility.
