# Glossary — Atlas

Domain terminology used throughout this repository. Every term is anchored to a frozen requirement, RFC, or journal entry.

---

## Framework Terms

### Adversarial Council
The five council seats (epistemologies) that review every major decision in the workspace. Designed to surface disagreement rather than produce consensus. Seats are: Empiricist, Cartographer, Red Team, Economist, Operator. Defined in `agents/GOVERNANCE.md`.

### Council Seat (Epistemology)
A permanent, opinionated *method of finding truth* assigned to a single council member. Unlike a role (e.g., "CEO", "CTO"), an epistemology is a way of analyzing — it can be applied to any question without role-collision. Source: `PROJECT_HISTORY.md` §2.

### Domain Anchor
A static knowledge file (`agents/domain/`) providing verified territory expertise in one field (distributed systems, AI/ML, product engineering, market/buyer). Domain anchors are consulted on-demand, do not sit on the council, and do not have opinions. A domain anchor unconsulted for 90 days becomes a `.candidate.md` and must be reactivated.

### Working Specialist
A dynamic agent in `agents/working/` spawned for a specific, narrow task. A specialist is temporary by design and must prove value via journal citations to survive. After 10+ journal citations and founder opt-in, it may be promoted to a domain anchor. After 6 dormant cycles, it is retired.

### Journal Entry
A durable, append-only record at `agents/journal/<YYYY-MM-DD>-<slug>.md` that captures every committed decision, including verbatim dissent and any founder override with reconsideration conditions. Source: `agents/GOVERNANCE.md` §Interaction protocol.

### Founder Override
A founder decision that contradicts the council recommendation. Must be recorded in the journal with: which seat was overruled, the founder's logic, and the specific condition that would trigger reconsideration. Source: `agents/GOVERNANCE.md` §Override protocol.

### Frozen Document
A planning or governance file listed in `scripts/frozen-docs.list` and hash-pinned in `FROZEN.sha256`. Amendments follow the strict process in `CONTRIBUTING.md` §Frozen planning documents and require a journal entry. Editing without following the process breaks CI.

### Confidence Label
A mandatory rating — `High / Medium / Low / None` — attached to every factual assertion in council outputs and journal entries. Confidence without cited evidence is forbidden. Source: `agents/GOVERNANCE.md`.

---

## System / Protocol Terms

### Principal
The workload identity on whose behalf a delegation is made. Holds a Permission Set. Originates delegation requests spanning a subset of its permissions. The system never issues base Principal identity — it consumes it from existing workload-identity infrastructure.
*Forces: FR1, ER1, INV1, C1.*

### Delegate
The workload or agent that receives a delegation and presents it to the relying party. Exercises no authority beyond what the delegation explicitly grants.
*Forces: FR1, ER1, ER10.*

### Delegation
The single presentable unit binding a Principal to a Delegate. Carries a Scope (subset of the Principal's Permission Set), an Expiration, and produces a Reconstruction Record. Terminal once Revoked or Expired.
*Forces: FR1–FR6, ER1–ER6, INV1–INV6.*

### Scope
A permission set that is a **strict subset** of the bound Principal's Permission Set. Enforced at issuance; an over-scoped request is refused at the issuance boundary. At verification time, the scope is checked for integrity (tamper evidence), not re-derived.
*Forces: FR2, ER2, INV2, SO6.*

### Attenuation
The property that a delegation's authority can only be reduced (scoped down), never amplified. A Delegate cannot issue a delegation with more scope than it received. This is a direct consequence of the Scope subset constraint.
*Forces: INV2, SO6.*

### Reconstruction Record
A tamper-evident, self-sufficient record produced for every delegation, from which an independent third party can determine: which identity delegated, to which identity, with what scope, and at what time. Must be reproducible without access to the original verifier's runtime state.
*Forces: FR6, NFR6, ER4, INV8, INV9, SO3.*

### Conformant Relying Party (RP)
A verifier that exercises the documented verification checks against a presented delegation. The system's security guarantees are defined exclusively at the conformant-verifier boundary. A non-conformant RP is outside the system's trust boundary.
*Forces: ER7, FM10, SO2.*

### Trust Domain
An independently-operated trust boundary. The system is bounded to exactly two: Domain A (containing the Principal) and Domain B (containing the Relying Party), with no shared runtime authority between them.
*Forces: C3, FR8, ER8, ER17.*

### Revocation
The one-way, terminal invalidation of exactly one specific delegation instance. Does not affect the underlying workload identities of the Principal or Delegate. The revocation's observability to the Relying Party is bounded by the latency parameter `R` in normal operation and by partition recovery time when the network is partitioned.
*Forces: FR4, ER5, ER6, INV4, INV5, INV6, SO1.*

### Partition (S4)
The network condition in which the Relying Party is isolated from the Revocation-Information Source. Revocations performed during a partition are not observable to the RP before partition recovery. This is an honest, information-theoretic system limit — not a gap.
*Forces: FM1, INV12, S4.*

### Verification Verdict
The outcome of a conformant RP checking a presented delegation. Three possible states:
- **Accept:** all checks pass (identity binding, scope integrity, expiry, revocation-observability).
- **Reject:** any single check fails.
- **Inconclusive [HYPOTHESIS]:** a required check cannot be conclusively determined; under the fail-closed hypothesis, transitions to Reject.
*Forces: SO5, FR3, NFR3, ER11, SO4, C-INV1.*

### Trust Material
Material held locally by the Relying Party, sufficient to verify a delegation offline without a live call to a shared authority.
*Forces: FR5, INV7, NFR2, C6.*

---

## Research Pipeline Terms

### R1 Cycle (Problem Ranking Application)
The first research cycle, completed 2026-06-19. Applied the Problem Ranking Framework v1.0 to all 10 Domain Research problems. Sequential council review completed. Founder selection of a candidate problem is deferred to the founder.

### R2 Cycle (Product Management)
The second research cycle, not yet started. Begins only after the founder selects a candidate problem from the R1 shortlist.

### Net Score
The core ranking metric of the Problem Ranking Framework. Computed as: Opportunity Score (axes A1–A6, max 18) minus Friction Score (axes A7–A10, max 12). Range: −12 to +18.

### Validity Gate
A mandatory filter in the Problem Ranking Framework: a problem must score ≥ 2 on Evidence Quality (A1), ≥ 1 on Frequency (A2), and ≥ 1 on Severity (A3) to proceed to council review. All 10 R1 problems passed.

<!-- checkpoint: feat(stores): implement key derivation -->

<!-- checkpoint: chore(stores): simplify verification controller -->

<!-- checkpoint: fix(scripts): fix Pre-commit validation scripts -->
