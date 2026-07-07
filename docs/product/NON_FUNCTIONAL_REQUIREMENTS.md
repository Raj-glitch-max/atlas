# Non-Functional Requirements

Each item is testable and traced. Items marked **[HYPOTHESIS]** carry a specific numeric or behavioral target inherited from a bounded experiment context, not a settled product requirement — treat the number as a starting point to validate, not a committed SLA.

**NFR1. [HYPOTHESIS]** Delegation verification (`FUNCTIONAL_REQUIREMENTS.md` FR5) should complete within a latency compatible with synchronous request paths.
*Evidence: `TECHNICAL_VALIDATION.md` (P5, item 8) used sub-100ms as the success threshold for its specific falsification experiment, not as a general product commitment. Carried forward here as a candidate target requiring its own validation at product scope, not assumed to transfer directly.*

**NFR2.** Verification of a delegation (FR5) shall not require a live network call to a central broker, and shall be testable by verifying a delegation with the relying party's network connection to any shared authority disabled.
*Evidence: `TECHNICAL_VALIDATION.md` P5 item 8, item 9 (failure criteria — "requires a live broker call" is explicitly listed as a falsifying outcome).*

**NFR3. [HYPOTHESIS]** When verification cannot conclusively determine a delegation's validity, the system shall reject the delegation rather than accept it (fail closed).
*Evidence: no prior document states this as a formal requirement verbatim. Inferred from `FOUNDER_PROBLEM_FIT.md`'s identification of "silent failure" and adversarial-reasoning requirements as the central risk of this problem (P5, items 6–7), but the specific fail-closed behavior itself is not directly evidenced and must be validated as a design requirement, not assumed as already agreed.*

**NFR4.** Revocation of a delegation (FR4) shall take effect without requiring the affected delegate workload to be redeployed or restarted, testable by revoking an active delegation and confirming subsequent rejection without any change to the workload's running state.
*Evidence: `PRODUCT_THESIS.md` P5 item 4 ("failure containment... without touching the underlying identity or restarting a workload").*

**NFR5.** The delegation-verification mechanism shall not require a relying party to adopt a protocol that is incompatible with, or a wholesale replacement for, identity-verification mechanisms already in use in that relying party's environment.
*Evidence: this is the single most consequential non-functional property in this package. Both `TECHNICAL_VALIDATION.md` (P5, item 9, failure criteria) and `ECOSYSTEM_THESIS.md` (P5, final verdict) identify this specific property as the condition that separates a genuine platform primitive from a feature nobody adopts. Failing this NFR does not make the product non-functional in a conventional sense — it makes it, by this project's own prior analysis, not the thing it was intended to be.*

**NFR6.** The delegation record produced under FR6 shall be tamper-evident, such that any alteration of the record after creation is detectable by a verifier.
*Evidence: implied directly by the core technical hypothesis's requirement that the chain of custody be "provable" (`TECHNICAL_VALIDATION.md`, P5, item 1) — a provable record that cannot detect tampering would not satisfy the hypothesis as stated.*
