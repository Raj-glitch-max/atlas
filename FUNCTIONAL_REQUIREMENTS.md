# Functional Requirements

Each requirement is testable (pass/fail determinable by inspection or experiment) and traced to a use case in `USE_CASE_CATALOG.md` and an evidence source. Requirements marked **[HYPOTHESIS]** are not supported by prior evidence and must not be treated as confirmed scope — see `V1_SCOPE.md` and `DEFERRED.md` for disposition.

**FR1.** The system shall allow a delegate workload to present proof that it is acting on behalf of a specified principal identity, such that a relying party can determine both identities from the presented proof alone.
*Traces to: UC1. Evidence: `TECHNICAL_VALIDATION.md` P5 item 1, item 8.*

**FR2.** The system shall allow a delegation to be constrained to a permission scope that is a strict subset of the principal's own permissions, such that the scope is inspectable by the relying party.
*Traces to: UC2. Evidence: `PRODUCT_THESIS.md` P5 item 3.*

**FR3.** The system shall allow a delegation to carry an expiration time after which a relying party must treat it as invalid, testable by presenting an expired delegation and confirming rejection.
*Traces to: UC2. Evidence: `PRODUCT_THESIS.md` P5 item 3.*

**FR4.** The system shall allow a specific delegation to be revoked such that: (a) the delegation is subsequently rejected by a relying party, and (b) the underlying identity of the principal and delegate remain valid and unaffected.
*Traces to: UC3. Evidence: `PRODUCT_THESIS.md` P5 item 2.*

**FR5.** The system shall allow a relying party to determine the validity, scope, and revocation status of a presented delegation without making a network call to a shared authority at the moment of verification.
*Traces to: UC4. Evidence: `TECHNICAL_VALIDATION.md` P5 item 8 (success criteria).*

**FR6.** The system shall produce a record, associated with each delegation, sufficient for a third party to determine which identity delegated to which, with what scope, and at what time — independent of and after the original verification event.
*Traces to: UC5. Evidence: `PRODUCT_THESIS.md` P5 item 4.*

**FR7.** The system shall allow identity and delegation to be issued to a workload without requiring that workload to have a long, statically-provisioned lifetime.
*Traces to: UC6. Evidence: `TECHNICAL_VALIDATION.md` P5 item 2.*

**FR8.** The system shall allow a relying party operating in a trust domain independent of the principal's trust domain to verify a delegation, within the bounded two-domain scenario defined in `TECHNICAL_VALIDATION.md`'s minimum experiment (P5, item 7).
*Traces to: UC1. Evidence: `TECHNICAL_VALIDATION.md` P5 item 7. Note: this requirement's validation is bounded to the two-domain experiment scope; see `V1_SCOPE.md`.*

**FR9. [HYPOTHESIS]** The system shall allow delegation verification to interoperate with relying parties that do not themselves adopt a new protocol beyond what they currently use for identity verification.
*Traces to: none directly — this is the "interoperability constraint" `TECHNICAL_VALIDATION.md` (P5, item 9) and `PRODUCT_THESIS.md` (P5, item 6) both identify as the condition separating a real primitive from a mere feature, but neither document confirms it is achievable. Presented here as a requirement to be tested, not a capability to be assumed.*

**FR10. [HYPOTHESIS]** The system shall tie a delegation's continued validity to the delegate's security posture at verification time, not only at issuance time.
*Traces to: UC8. Evidence status: `TECHNICAL_VALIDATION.md` classifies this assumption as unproven (P5, item 3). Not supported as an achievable requirement by any prior evidence.*
