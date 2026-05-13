# User Model

**Governing caveat, stated once here and inherited by every persona below:** `FOUNDER_DECISION_BRIEF.md` (P5, item 3) explicitly concluded "we do not know" who the buyer or primary beneficiary is — security team, platform team, or an open-source infrastructure component with no direct buyer at all. No prior document resolved this. Every persona below is therefore a **hypothesis** unless otherwise marked, derived from *who would technically interact with the system* (evidenced by `SYSTEM_CONTEXT.md` and `ECOSYSTEM_THESIS.md`'s natural-integrator analysis) rather than *who would adopt or pay for it* (unresolved).

## Personas

### 1. Platform/Infrastructure Engineer — **evidenced**
Operates trust domains, manages workload-identity infrastructure. Directly evidenced by the Operator council role's note that this problem "requires security-process integration" (`PROJECT_HISTORY.md`, R1 cycle) and by `ECOSYSTEM_THESIS.md`'s identification of SPIRE operators as natural integrators.

### 2. Security Engineer / Security Architect — **evidenced**
Defines and reviews delegation policy and trust boundaries. Directly evidenced by the problem's OWASP/DBIR/NIST evidence base (`PROJECT_HISTORY.md`, R1 cycle) and by `FOUNDER_PROBLEM_FIT.md`'s conclusion that security-engineering judgment is a required capability to build this — the same judgment is required to operate and trust it.

### 3. Delegate workload or agent — **evidenced**
The non-human actor presenting proof of delegation. This is not a persona in the conventional sense but is included because it is the entity whose behavior the core hypothesis is directly about (`TECHNICAL_VALIDATION.md`, P5, item 1).

### 4. Relying-party service owner — **evidenced**
A team whose service must verify inbound delegated identity. Directly evidenced by `SYSTEM_CONTEXT.md`'s "relying party" role and by the Technical Validation's success/failure criteria, which are defined from the relying party's point of view.

### 5. Auditor / compliance reviewer — **hypothesis**
An entity that inspects delegation chains after the fact. This role is suggested only by `PRODUCT_THESIS.md`'s observation that structural audit trails "become possible" (P5, item 4) — a capability statement, not a validated user need. No prior document establishes that an auditor persona exists, wants this, or would use it. Treated as a hypothesis to be tested, not a requirement source.

## What this document does not claim

It does not claim any of the above personas has been interviewed, would pay, or has confirmed a need. `ASSUMPTIONS_AND_RISKS.md` carries the unresolved-buyer risk forward explicitly so it is not silently dropped.

<!-- checkpoint: chore(record): clean revstatus snapshot retrieval -->
