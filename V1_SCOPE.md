# V1 Scope

V1 is bounded to exactly what `TECHNICAL_VALIDATION.md`'s minimum falsifying experiment (P5, item 7) covers, because it is the only scope this project has a defined, evidence-based test for. Anything wider is scope not yet earned by evidence.

## In scope for V1

- FR1 — Prove delegated identity, within the two-trust-domain scenario.
- FR2 — Scope a delegation to a subset of the principal's permissions.
- FR3 — Time-bound a delegation with an expiration.
- FR4 — Revoke a specific delegation independent of underlying identity.
- FR5 — Verify a delegation without a live broker call.
- FR6 — Produce a record sufficient to reconstruct a delegation chain.
- FR7 — Issue identity/delegation to a short-lived workload.
- FR8 — Verify across the bounded two-domain scenario specifically defined in the Technical Validation experiment.
- NFR2, NFR4, NFR6 — non-functional properties directly required by the above and directly testable within the same bounded scenario.

## Explicitly not in scope for V1

- FR9 (cross-protocol interoperability with parties that haven't adopted anything new) — hypothesis, unproven, see `ASSUMPTIONS_AND_RISKS.md`.
- FR10 (continuous re-attestation) — hypothesis, unproven.
- NFR1's specific latency number — inherited from a different experiment's success criterion, not yet validated at product scope; V1 should measure and report a number, not commit to one in advance.
- NFR3 (fail-closed behavior) — a reasonable design principle but not yet a validated requirement; V1 should test and document actual behavior rather than assume it's already satisfied.
- Anything involving more than two trust domains.
- Anything involving non-SPIFFE environments (VMs, serverless, SaaS-to-SaaS) — mentioned in `TECHNICAL_VALIDATION.md`'s hidden assumptions as an open question, never validated.
- Any commercial packaging or named buyer — precluded by `CONSTRAINTS.md` C5.

## Definition of done for V1

V1 is complete when the Technical Validation's own success and failure criteria (`TECHNICAL_VALIDATION.md`, P5, items 8–9) have been evaluated against a working implementation of the in-scope requirements above, and the result — success or failure — is documented honestly. A negative result that clearly identifies which requirement failed and why is a valid, complete V1 outcome; it is not a synonym for project failure. See `FOUNDER_PROBLEM_FIT.md` (P5, item 8) for why a rigorous negative result counts as success at this stage.
