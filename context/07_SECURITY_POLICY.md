# Security Policy — Atlas

This document records the security principles, threat model boundaries, and failure model invariants for the Atlas workload delegation system. All claims are derived from frozen requirements.

> [!IMPORTANT]
> This document describes the **security model of the system being designed**, not the security of the Atlas repository itself. Repository-level security (secrets scanning, dependency checks) is handled by `CONTRIBUTING.md`, `.gitleaks.toml`, and `make secrets`.

---

## 1. Security Objectives (Frozen — Phase 8, `02_SECURITY_OBJECTIVES.md`)

| ID | Objective |
|---|---|
| SO1 | A relying party can detect whether a delegation has been revoked, within the declared observability bound R, absent a network partition. |
| SO2 | Verification must be possible offline — no live call to a shared authority is required. |
| SO3 | A delegation's Reconstruction Record is tamper-evident and self-sufficient for independent review. |
| SO4 | The system degrades to an inconclusive verdict (not a silent accept) when a required check cannot be completed. |
| SO5 | A single failed check causes a full rejection (no partial-accept). |
| SO6 | Delegation scope is a strict subset of the principal's permissions; over-scope requests are refused at issuance. |
| SO7 | Delegation verification is additive to, not a replacement for, existing identity verification at the relying party. |
| SO8 | An independent reviewer can reproduce a verification verdict from the Reconstruction Record alone, without access to the original verifier's runtime state. |

---

## 2. System Invariants Relevant to Security (Frozen — Phase 8, `03_SYSTEM_INVARIANTS.md`)

| ID | Invariant |
|---|---|
| INV1 | A delegation that verifies successfully binds exactly one Principal and one Delegate, both deterministically recoverable. |
| INV2 | A delegation's scope is a strict subset of the bound Principal's Permission Set; over-scope is refused at issuance. |
| INV3 | A delegation becomes invalid at its Expiration time; expiry is monotone and terminal. |
| INV4 | Once revoked, a delegation cannot become valid again for that specific instance. |
| INV5 | Revocation does not affect the underlying identities of the Principal or Delegate. |
| INV6 | Revocation targets exactly one delegation instance. |
| INV7 | Core verification requires no live call to a shared authority. |
| INV8 | The Reconstruction Record is tamper-evident; any alteration is detectable. |
| INV9 | The Reconstruction Record is self-sufficient for independent third-party reconstruction. |
| INV10 | The system never issues base workload identity; it consumes identity from external infrastructure. |
| INV11 | The system does not modify the existing workload-identity standard. |
| INV12 | In-partition revocation is not observable to the RP before partition recovery; this bound is honest, not a gap. |

---

## 3. Failure Model — Honest Limits (Frozen — Phase 8, `04_FAILURE_MODEL.md`)

These failures are outside the system's warranted scope. They are documented as honest limits, not gaps to be resolved.

| ID | Failure | Status |
|---|---|---|
| FM5 | Issuer key compromise | **UNMITIGATED** — outside the current-scope boundary. Carried as an honest limit. |
| FM8 | Within-window replay attacks | **UNMITIGATED** — outside the current-scope boundary. Carried as an honest limit. |
| FM1/FM2 | In-partition revocation observability | **BOUNDED HONESTLY** — observability is capped by partition recovery time, not R. Not a bug. |
| FM6 | Over-scope issuance | **REFUSED** at the Issuance boundary; no delegation is created. |
| FM10 | Non-conformant relying party behavior | **OUTSIDE TRUST BOUNDARY** — the system makes no guarantees for non-conformant verifiers. |
| FM3/FM9 | Clock skew / trust material unavailability | **INCONCLUSIVE** — yields Inconclusive verdict; fail-closed hypothesis applies. |

---

## 4. Threat Model Boundaries

### In Scope
- **Scope inflation:** A delegate cannot exceed the authority granted by the delegation.
- **Tamper detection:** Any alteration to a Reconstruction Record is detectable by a conformant verifier.
- **Unauthorized delegation:** The issuance boundary refuses requests not authorized by the Principal.
- **Revocation observability:** The system provides a bounded-latency revocation signal to conformant verifiers.

### Out of Scope
- **Compromised identity infrastructure:** If the underlying workload-identity issuer (e.g., SPIFFE/SPIRE) is compromised, delegation integrity is void. This is a dependency-chain attack, not a scope the system addresses.
- **Malicious relying parties:** Non-conformant RPs are not in the trust boundary. The system cannot control whether an RP chooses to honor its verification result.
- **Within-window replay:** If a valid delegation is intercepted and replayed within its validity window, the system does not detect or prevent this (FM8). This limit is carried openly.
- **Three-or-more-domain federation:** The system is strictly bounded to two trust domains (C3, ER17). Multi-domain scenarios are out of V1 scope (DEFERRED D3–D4).

---

## 5. Repository Security

- **Secret scanning:** Run `make secrets` (gitleaks). CI runs this automatically.
- **Frozen-doc integrity:** Run `make check-frozen` (SHA-256). CI runs this automatically.
- **Pre-commit hooks:** Enforces trailing whitespace, YAML/JSON validity, large-file guard, and private-key detection.
- **No credentials in code:** No secrets, tokens, or credentials may appear in any file in this repository. See `.gitleaks.toml` for the scan configuration.

<!-- checkpoint: docs(threat-model-scenarios): restructure threat model scenarios -->
