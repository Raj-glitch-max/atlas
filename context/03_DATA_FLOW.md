# Data Flows — Atlas

This document describes the high-level flows of data through the system. All flows are defined at the boundary level; no protocol, wire format, API signature, or implementation mechanism is specified here. Every flow traces to a frozen requirement.

---

## 1. Delegation Lifecycle Flow

The delegation lifecycle is the primary data flow of the system.

```
Principal Workload           System Boundary           Relying Party (Domain B)
      │                            │                            │
      │── 1. Delegation Request ──►│                            │
      │   (scope ⊆ permissions)    │                            │
      │                            │ [Issuance Boundary Check]  │
      │                            │  scope ⊂ Permission Set?   │
      │                            │  If not → REFUSE (FM6)     │
      │                            │                            │
      │◄── 2. Delegation Issued ───│                            │
      │    + Reconstruction Record │                            │
      │                            │                            │
      │── 3. Delegation Presented ──────────────────────────►  │
      │                            │                [Conformant Verifier Boundary]
      │                            │                 Checks offline with Trust Material:
      │                            │                 - Identity binding (M-INV1)
      │                            │                 - Scope integrity (M-INV8)
      │                            │                 - Expiry (M-INV3)
      │                            │                 - Revocation-observability (M-INV4)
      │                            │                            │
      │                            │                ◄── 4. Verdict ──►
      │                            │                  Accept / Reject / Inconclusive[HYP]
```

**Forces:** FR1–FR6, ER1–ER7, INV1–INV9, SO1–SO6, AT1–AT30.

---

## 2. Revocation Flow

Revocation targets exactly one delegation instance and has no side-effect on underlying identities.

```
Principal (or Authorized Actor)    Revocation-Information Source    Conformant Relying Party
           │                                │                                │
           │── Revoke(delegation-ref) ─────►│                                │
           │                                │                                │
           │                                │── Observable? ────────────────►│
           │                                │   (within R in non-partitioned │
           │                                │    operation; bounded by        │
           │                                │    partition-recovery if S4)    │
           │                                │                                 │
           │                                │   If not-yet-observable → RP   │
           │                                │   may Accept until Observable  │
```

**Key constraints:**
- Revocation is **one-way and terminal** for that delegation instance (M-INV4, M-INV6).
- Revocation does **not** affect the Principal's or Delegate's underlying workload identity (M-INV5).
- In-partition revocation is **not** observable before partition recovery; this is an honest system limit, not a gap (M-INV12, FM1, S4).

---

## 3. Independent Reconstruction Flow

Any reconstruction is performed by an independent third-party reviewer from the Reconstruction Record alone.

```
Independent Reviewer         Reconstruction Record
       │                             │
       │── Load Record ─────────────►│
       │                             │
       │◄── Determine: Principal, Delegate, Scope, Time ──
       │   (No runtime verifier state. No privileged access.)
       │
       │── Reproduce Verdict ──────────────────────────────►
```

**Forces:** FR6, ER4, INV8–INV9, SO3, SO8, AP3, AP8.

---

## 4. Agent Reasoning Decision Flow (Framework)

The adversarial council is a separate conceptual data flow for workspace decisions.

```
Founder Input (artifact / question)
       │
       ▼
   Empiricist → Cartographer → Red Team → Economist → Operator
   (sequential, one response each, no voting)
       │
       ▼
   Domain Anchors consulted on-demand (by council member request)
       │
       ▼
   Founder Decision (with any override recorded)
       │
       ▼
   Journal Entry written (verbatim dissent preserved)
       │
       ▼
   Context files / frozen docs updated if required
```

**Forces:** `agents/GOVERNANCE.md`, §§ Interaction protocol, Override protocol.

<!-- checkpoint: rfc(architecture-draft): clarify architecture draft -->

<!-- checkpoint: fix(record): fix boundary check -->

<!-- checkpoint: test(internal): test revstatus snapshot retrieval -->

<!-- checkpoint: chore(record): tweak attenuation rule engine (#84) -->

<!-- checkpoint: chore(record): optimize test assertions -->

<!-- checkpoint: chore(lab): tweak integration test runner (#134) -->
