# Deferred

Items explicitly excluded from `V1_SCOPE.md`, listed here so they are not silently lost. Each includes why it's deferred and what would need to be true to bring it back into scope.

**D1 — Cross-protocol interoperability with parties that adopt nothing new (FR9).**
Deferred because: unproven, and is in fact the single condition on which this entire project's primitive-vs-feature status depends (`ECOSYSTEM_THESIS.md`, P5 final verdict). Reconsider when: the V1 experiment produces a result specifically addressing this, per `V1_SCOPE.md`'s definition of done.

**D2 — Continuous re-attestation of delegate posture (FR10).**
Deferred because: classified unproven in `TECHNICAL_VALIDATION.md` (P5, item 3), with no experiment in this package designed to test it. Reconsider when: a dedicated experiment for this specific assumption is defined — it is a different technical question than delegation itself.

**D3 — More than two trust domains / broader federation.**
Deferred because: no evidence exists beyond the two-domain scenario; extrapolating to N domains is unvalidated. Reconsider when: the two-domain result is in hand and specifically suggests (rather than assumes) that broader federation is tractable.

**D4 — Non-SPIFFE environments (VMs, serverless, SaaS-to-SaaS, non-Kubernetes agents).**
Deferred because: `TECHNICAL_VALIDATION.md` names this as a hidden assumption never validated (P5, item 2) — that these environments can supply attestation-quality signals at all. Reconsider when: at least one such environment has been tested directly.

**D5 — Committed latency/performance SLA.**
Deferred because: the only number in this package (NFR1) is inherited from a different experiment's success threshold, not derived from product-scale testing. Reconsider when: V1 produces an actual measured number under realistic conditions.

**D6 — Fail-closed behavior as a formal, tested requirement (NFR3).**
Deferred as a *confirmed* requirement, though retained as a design principle to test. Reconsider when: V1's implementation has been deliberately tested against ambiguous/adversarial verification conditions rather than assumed to behave correctly.

**D7 — Any buyer identification, commercial packaging, pricing, or go-to-market work.**
Deferred because: explicitly out of this project's frozen Strategy Phase scope (`CONSTRAINTS.md` C5) and out of scope for Product Definition entirely. Reconsider when: a separate, later phase is explicitly opened for it — not a matter for this package to anticipate.

**D8 — Auditor persona and any audit-specific tooling or workflow (UC5's downstream use, as distinct from record production itself).**
Deferred because: `USER_MODEL.md` marks the auditor persona as hypothesis, not evidenced. Reconsider when: a specific auditor need is identified and evidenced, not assumed from the fact that a record could theoretically be audited.

**D9 — Trust bridging between organizations with no prior federation relationship (UC7).**
Deferred because: `PRODUCT_THESIS.md` itself offered this only to test generativity, not as a proposal (P5, item 4). Reconsider when: it is raised as an actual requirement by evidence, not inferred from a capability list.
