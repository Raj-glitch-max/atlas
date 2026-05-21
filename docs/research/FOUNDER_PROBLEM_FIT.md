# Founder-Problem Fit Analysis

**Scope limitation, stated up front:** No founder profile exists anywhere in this project's record. Session 5 (Founder Profile) explicitly returned sparse intake rather than a fabricated profile, and H3 confirmed founder-advantage = 0 across all ten original problems as a direct consequence. This document does not fabricate one. Where a question is genuinely about *you* — your background, your risk tolerance, your available time — it's marked as an open gap, not answered. Where a question is about what the *problem* demands of whoever builds it, that's answerable now, from the Technical Validation already on record, and is answered in full below.

---

## P5 — IAM least-privilege / service-to-service identity

### 1. Why this could be a good fit
This is answerable only conditionally. It would be a strong fit for someone who: enjoys correctness-critical systems work where the cost of a subtle bug is high; has patience for standards-and-protocol-style thinking rather than pure feature-building; is comfortable operating in a space with real prior art (Kerberos, OAuth, SPIFFE) to study and compose rather than invent from nothing; and wants a project with a clean, falsifiable engineering question (Task 4 in the Technical Validation) rather than an open-ended one. **Unknown whether this describes you.**

### 2. Why this could be a bad fit
Also conditional. It would be a poor fit for someone who: is energized primarily by fast, visible iteration (security-critical work rewards slowness and review, not shipping speed); doesn't yet have real comfort with applied cryptography or PKI concepts (getting this wrong isn't just a bug, it's a trust failure); or is looking for a project where "mostly working" has real value (a delegation system that's mostly correct is arguably worse than no system, because it creates false confidence). **Unknown whether any of this describes you.**

### 3. Capabilities needed before beginning
These are demanded by the problem, not by you specifically:
- Working knowledge of distributed-systems fundamentals (trust boundaries, network partitions, replay attacks).
- Applied familiarity with PKI/X.509 and at least one prior exposure to a real identity or auth system, even as a user/integrator, not necessarily a builder.
- Enough security-engineering judgment to reason about adversarial cases, not just happy-path correctness — this is arguably the single hardest prerequisite, and the hardest to self-assess honestly.
- Comfort reading standards documents and RFCs (SPIFFE spec, RFC 8693) as primary sources, since the project explicitly requires composing existing primitives correctly rather than inventing new ones.

### 4. Capabilities that could reasonably be learned during the project
- The specific internals of SPIFFE/SPIRE (attestation flow, workload API, trust bundle federation) — this is documented, learnable technology, not a prerequisite.
- Delegation-token formats and token-exchange mechanics (RFC 8693 specifics).
- Kubernetes-native operational patterns, if not already familiar, since most current tooling in this space assumes a Kubernetes substrate.
- General protocol-design instincts, which sharpen naturally through doing this kind of work carefully.

### 5. Biggest execution risk
Shipping something that *looks* finished but has an invisible trust flaw — a delegation path that works for the tested cases but fails under an adversarial or edge-case scenario nobody thought to test. Unlike most software bugs, this class of failure doesn't announce itself; the system keeps working right up until it's exploited. That makes normal "it works, ship it" instincts actively dangerous here.

### 6. Biggest technical risk
Landing on the wrong side of the interoperability constraint identified in the Technical Validation — building something that technically works but requires new protocol adoption from every counterparty, at which point it hasn't actually solved the stated hard problem. It becomes just another proprietary identity system, indistinguishable in kind from what Teleport or others have already shipped.

### 7. Biggest personal risk
A "silent failure" pattern: months of careful, correctness-focused work that produces something real but has no natural external validation point. Unlike a product with users, or a research result with a benchmark number, a security primitive's correctness is genuinely hard to self-certify. Without a security-focused peer or auditor to review the work, there's a real risk of high effort producing something the builder *believes* is sound without actually knowing whether it is.

### 8. What success would realistically look like after 6 months of disciplined work
A working, documented delegation-chain prototype that passes the specific falsification experiment defined in the Technical Validation (cross-domain proof, sub-100ms offline verification, no new heavyweight protocol required) — likely released as a minimal open-source reference implementation with a technical writeup, not a polished product. Six months is a realistic timeline for *a rigorous prototype*, not for a production-hardened system.

### 9. What failure would most likely look like
Partial implementation that handles the happy path convincingly but was never adversarially tested — no fuzzing, no red-team review, no explicit reasoning about failure modes — combined with an overrun timeline, since correctness-critical work is notoriously hard to time-box honestly. The failure mode here is quiet, not loud: something that "seems done" but isn't trustworthy.

### 10. Project shape
- **Flagship portfolio project:** Strong fit — a working, honestly-tested delegation prototype is a specific, demonstrable, hard-to-fake artifact.
- **Open-source project:** Strong fit, arguably necessary — the interoperability constraint that defines success basically requires openness; a closed delegation protocol undermines its own value proposition.
- **Startup:** Weaker fit on engineering grounds alone (setting business questions aside per scope) — the Technical Validation classified this as research-level *systems engineering*, not open science, meaning a startup would need to compete on execution and trust in a space with well-funded incumbents already working the same edge.
- **Research project:** Weaker fit — the hard problem here is composition and correct engineering of known primitives, not generating new science. Framing it as "research" risks giving it more open-endedness than the problem actually has.

---

## P10 — Production debugging / root-cause speed

### 1. Why this could be a good fit
Conditionally: a strong fit for someone energized by open problems with no guaranteed answer, comfortable reading and building on academic literature (causal inference, ML papers), and genuinely motivated by rigor over shipping — someone for whom a well-documented negative result feels like a real outcome, not a failure. **Unknown whether this describes you.**

### 2. Why this could be a bad fit
A poor fit for someone who needs a working v1 to stay motivated, who hasn't previously sat with an unsolved problem for months without a payoff, or who is uncomfortable with the possibility that six months of disciplined work could honestly conclude "this doesn't generalize" with no product to show for it. **Unknown whether any of this describes you.**

### 3. Capabilities needed before beginning
- Genuine comfort with applied statistics / causal inference concepts (causal discovery, confounding, the specific limitations of methods like PC algorithm or LiNGAM) — not just ML familiarity in general.
- Distributed-systems and observability engineering experience, to build realistic test systems and injected-fault scenarios rather than relying on toy benchmarks.
- Research literacy: the ability to read a paper like the 2026 survey cited in the Technical Validation, understand *why* a method fails to generalize, and avoid repeating the same mistake rather than reinventing it.
- Experimental discipline — the willingness to design a fair falsification test (Task 3, Technical Validation) rather than unconsciously tuning toward a benchmark that flatters the method, which is precisely the failure mode the field's own literature documents repeatedly.

### 4. Capabilities that could reasonably be learned during the project
- Specific causal-discovery toolchains (gCastle and similar) — these are learnable tools, not prerequisites.
- Benchmark methodology, by directly engaging with public benchmarks like RCAEval.
- Fault-injection system design, which improves with iteration.

### 5. Biggest execution risk
Scope creep into an unbounded research loop with no natural stopping point. Because the field itself hasn't converged, there's always a plausible "one more experiment" — a different causal-discovery method, a different benchmark, a different feature set — and the project can absorb unlimited time without ever producing a clear answer. This is the specific trap named in the prompt, and it's not hypothetical: the published literature shows credentialed research teams still iterating on this exact problem for years.

### 6. Biggest technical risk
The technique works on the test system it was built and tuned against, and fails to generalize the moment scale or topology changes — the single most common documented failure mode in this literature, not a rare edge case.

### 7. Biggest personal risk
"Sunk cost drift" — research has an inherently seductive quality where results are always "almost there," making it easy to lose track of calendar time chasing marginal accuracy gains rather than reaching a clear stopping point. Unlike P5's silent-failure risk, this risk is loud and visible in retrospect, but hard to notice from inside it.

### 8. What success would realistically look like after 6 months of disciplined work
An honest, well-documented benchmark result on a moderately realistic system (per the Technical Validation's minimum experiment: 30+ services, labeled fault injection, no hand-fed topology) — either a genuinely promising early signal worth continuing, or a clear, rigorous negative result explaining specifically where and why the approach breaks down. For a problem this open, "a trustworthy answer" is a legitimate six-month success outcome, not a consolation prize — but it is a different kind of success than "a working v1."

### 9. What failure would most likely look like
Months spent narrowing and re-tuning hyperparameters on an ever-smaller, ever-more-favorable benchmark variant, without ever honestly testing generalization — producing a method that looks good on paper (or in a demo) and would fail immediately outside the conditions it was tuned for. This is exactly the pattern the field's own literature warns is common, not a hypothetical worst case.

### 10. Project shape
- **Flagship portfolio project:** Plausible fit specifically *because* it's hard — attempting a genuinely unsolved problem and producing a rigorous, honest result (positive or negative) is a differentiated signal most portfolios don't have. Requires the discipline to stop and document rather than quietly abandon.
- **Open-source project:** Plausible fit if the work produces a reusable contribution — e.g., an extension to or result on a public benchmark like RCAEval — but less structurally necessary to the project's value than it is for P5.
- **Startup:** Weakest fit of the four framings on engineering grounds alone — the Technical Validation concluded success requires beating an active academic frontier, not composing known techniques, and a startup needs a credible v1 on a business timeline that open science doesn't reliably provide.
- **Research project:** Strongest natural fit — this is, by the Technical Validation's own classification, closer to unresolved science than product engineering. Framing it explicitly as research (rather than as a product-in-disguise) is the framing most consistent with what the problem actually is.

---

## What's missing to turn this into real fit analysis

The items below are the specific gaps that, if filled, would let this comparison be about *you* rather than about the problems in the abstract. Consistent with this project's own intake discipline, none of them are guessed at here:

- Prior technical background: languages, systems worked in, depth of distributed-systems or security experience.
- Any prior completed projects (shipped, open-sourced, or otherwise), regardless of scale — H3's own rule is that stated interest doesn't count, only prior behavior does.
- Comfort level specifically with security-adjacent work versus statistics/ML-adjacent work, since P5 and P10 draw on genuinely different skill sets.
- Tolerance for an extended period without a demoable v1 — this cuts differently for P5 (silent-failure risk) than for P10 (open-ended-research risk).
- Available time horizon and whether six months is a hard constraint or a checkpoint.
- Access to any external review capability — a security-literate peer for P5, or a research-literate peer for P10 — since both projects' biggest personal risks are versions of "no one is checking my work."

This document does not recommend one problem over the other, and does not resolve these gaps. It's structured so that once you have answers to the questions above, the fit conclusion should follow relatively directly from what's already written here.

<!-- checkpoint: chore(scripts): audit Docker orchestration config -->
