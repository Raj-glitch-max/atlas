# Assumptions and Risks

## Assumptions, by evidentiary status (carried verbatim in classification from `TECHNICAL_VALIDATION.md`, P5, item 3 — not re-assessed here)

**Proven:**
- Static workload identity via SPIFFE/SPIRE functions in production at real scale (named deployments: Uber, Stripe, Netflix).

**Partially proven:**
- Cross-trust-domain federation exists in the SPIFFE specification but is operationally heavy in practice, per independent 2026 practitioner sources.

**Unproven:**
- Continuous, post-issuance re-attestation of a delegate's posture (NFR/FR10 in this package are marked hypothesis specifically because of this).
- A cryptographically verifiable, cross-protocol delegation chain for ephemeral agents, interoperable across independently-adopted systems (FR9 in this package is marked hypothesis specifically because of this).

**Probably false:**
- That this problem can be solved through additional SPIRE configuration alone, rather than a new composed primitive. Multiple independent, unaffiliated sources converge on this conclusion.

## Risks (traced to source)

**R1 — Silent trust failure.** A delegation mechanism that appears to work under tested conditions but has an undetected flaw exploitable under adversarial conditions. This is qualitatively different from a conventional software bug: it does not announce itself.
*Source: `FOUNDER_PROBLEM_FIT.md` P5, item 5 (biggest execution risk), item 6.*

**R2 — Interoperability failure.** Building a working delegation mechanism that requires new protocol adoption from every counterparty, at which point it has not solved the stated problem — it has produced another proprietary identity system.
*Source: `TECHNICAL_VALIDATION.md` P5, item 9 (failure criteria); `ECOSYSTEM_THESIS.md` P5, final verdict.*

**R3 — Standards conflict with SPIFFE's own scoping decision.** SPIFFE has deliberately excluded authorization/delegation semantics from its specification. Building against this constraint (see `CONSTRAINTS.md` C2) rather than around it risks rejection by the exact ecosystem this product needs to interoperate with.
*Source: `ECOSYSTEM_THESIS.md` P5, item 2.*

**R4 — Vendor fragmentation.** At least one well-funded vendor (Teleport, March 2026) has already shipped a proprietary partial answer to this problem. Adoption could consolidate around a closed answer before an open one is viable.
*Source: `ECOSYSTEM_THESIS.md` P5, final verdict.*

**R5 — Unresolved buyer/beneficiary.** No prior document has established who adopts or pays for this capability, or whether it has a direct beneficiary at all versus functioning as unattributed open-source infrastructure.
*Source: `FOUNDER_DECISION_BRIEF.md` P5, item 3; carried into `CONSTRAINTS.md` C5 and `USER_MODEL.md`.*

**R6 — No external validation point.** Unlike a product with users or a research result with a benchmark, a security primitive's correctness is difficult to self-certify. Without an independent security review, confidence in correctness cannot be distinguished from actual correctness.
*Source: `FOUNDER_PROBLEM_FIT.md` P5, item 7 (biggest personal risk).*

<!-- checkpoint: chore(lab): harden network partition test -->
