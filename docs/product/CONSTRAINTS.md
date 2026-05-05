# Constraints

Externally-imposed boundaries this product must operate within. These are not design choices — they are facts about the environment this product must exist inside, established in prior phases.

**C1 — Must coexist with, not replace, existing workload-identity infrastructure.**
The product depends on base workload identity already being issued by SPIFFE/SPIRE-based systems; it does not take on that responsibility itself.
*Source: `ECOSYSTEM_THESIS.md` P5, items 1 and 4; `SYSTEM_CONTEXT.md`.*

**C2 — Must not depend on a change to SPIFFE's core specification.**
SPIFFE maintainers have deliberately scoped the specification to authentication only, excluding authorization/delegation semantics. This product must function as an extension or companion capability, not as a proposed change to SPIFFE itself.
*Source: `ECOSYSTEM_THESIS.md` P5, item 2 (conflict).*

**C3 — Initial technical validation is bounded to a two-trust-domain scenario.**
The only experiment this project has defined and committed to (`TECHNICAL_VALIDATION.md` P5, item 7) covers two independent trust domains. Any claim of validity beyond that scope (more domains, non-SPIFFE environments, production-scale load) is unvalidated by definition and must be labeled as such wherever it appears in this package.
*Source: `TECHNICAL_VALIDATION.md` P5, item 7.*

**C4 — No claim of production-readiness at the six-month horizon.**
Both `TECHNICAL_VALIDATION.md` (P5, item 10) and `FOUNDER_PROBLEM_FIT.md` (P5, item 8) independently concluded that six months of disciplined work realistically yields a validated reference implementation, not a hardened production system.
*Source: `TECHNICAL_VALIDATION.md` P5 item 10; `FOUNDER_PROBLEM_FIT.md` P5 item 8.*

**C5 — No presumed buyer or commercial packaging.**
`FOUNDER_DECISION_BRIEF.md` (P5, item 3) explicitly left the buyer/beneficiary question unresolved. This Product Definition package does not resolve it and must not be read as implying a specific buyer exists.
*Source: `FOUNDER_DECISION_BRIEF.md` P5, item 3; carried into `USER_MODEL.md`.*

**C6 — No dependency on a shared, always-online broker for core verification.**
This is listed as both a functional requirement (FR5) and a constraint because it is not negotiable within scope: a design that requires one would fail the falsification experiment by definition, per `TECHNICAL_VALIDATION.md`'s own failure criteria.
*Source: `TECHNICAL_VALIDATION.md` P5, item 9.*
