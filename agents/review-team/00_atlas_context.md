# 00 — ATLAS SHARED CONTEXT (INJECT INTO EVERY AGENT)

> This file is the single source of truth for project context. Every agent file
> references it. If Atlas changes, change it **here**, not in 13 places.

---

## 1. Identity

- **Name**: Atlas
- **Category**: AI Infrastructure / Security Runtime / Agent Runtime / Distributed Systems / Cryptographic Authorization / Developer Infrastructure
- **One-liner**: A runtime and protocol that lets AI agents securely delegate permissions that remain independently verifiable even when the original authorization server is unreachable.
- **Status**: Research Complete → Architecture Complete → Core Runtime Under Development → SDKs Planned → Application Layer Planned

## 2. The Gap Atlas Occupies

| System | Solves | Assumes |
|---|---|---|
| OAuth 2.0 / OIDC | User → app authorization | Live authorization server reachable at decision time |
| SPIFFE / SPIRE | Workload identity | Live SPIRE agent / trust bundle distribution |
| MCP | Tool invocation standardization | Authorization is someone else's problem |
| **Atlas** | **Verifiable delegation across agents** | **Nothing. Verification is offline, deterministic, fail-closed.** |

Atlas does not replace any of the above. It fills the hole they collectively leave:
**cryptographically verifiable, attenuable delegation that survives partition.**

## 3. Tech Stack (Exact — do not hallucinate alternatives)

| Layer | Technology |
|---|---|
| Language | Go 1.21+ |
| Cryptography | JOSE / JWS / ES256 / go-jose v3 |
| Identity | SPIFFE, SPIRE |
| Runtime | Linux, Docker, containers |
| Infra | Docker Compose, Prometheus, Grafana |
| Testing | Go testing, property tests, differential testing, coverage-guided fuzzing, conformance vectors, adversarial vectors |
| CI | GitHub Actions, pre-commit, doc lint, secret scanning, frozen-document integrity |

## 4. Repository Architecture

```
record/      # Delegation record model
issuance/    # Signed delegation creation
verify/      # Verification engine
truststore/  # Trusted issuer management
revstatus/   # Revocation status
revorigin/   # Signed revocation data production
internal/    # Private implementation
context/     # Context management
agents/      # Agent integrations (planned)
lab/         # Distributed experiments, fault injection
tests/       # Test suites
bench/       # Benchmarks
scripts/     # Automation
docs/        # Documentation
```

## 5. Engineering Layers

| Layer | Focus |
|---|---|
| 1. Protocol | Delegation record, cryptographic format, verification rules, revocation semantics |
| 2. Core Runtime | Issue, delegate, verify, reject, revoke |
| 3. Conformance | Language-neutral vectors, property testing, fuzzing, differential verification |
| 4. Performance | Latency, memory, record size, verification cost |
| 5. Research Lab | Partition simulation, fault injection, Docker lab, SPIRE lab |

## 6. Implemented

✓ Offline verification ✓ ES256/JWS signatures ✓ Capability attenuation ✓ Delegation chains
✓ Cross-domain verification ✓ Fail-closed semantics ✓ Revocation snapshots
✓ Differential testing ✓ Property testing ✓ Coverage-guided fuzzing ✓ Benchmarks ✓ Conformance vectors

## 7. Measured Performance (single machine, lab)

| Operation | Cost |
|---|---|
| Verification | ~94 μs |
| Issuance | ~30 μs |
| Integrity validation | ~80 μs |
| Delegation record size | ~403 bytes |

**Every agent must treat these as *unqualified* numbers** — no hardware spec, no percentile,
no chain depth stated. Flagging this is legitimate and expected.

## 8. Security Model

**Guarantees**: tamper-evident records, offline verification, capability attenuation, cryptographic integrity, deterministic verification, fail-closed, replay resistance, independent verification.

**Threats addressed**: signature forgery, algorithm confusion, payload transplantation, malformed records, scope widening, tampering, revocation freshness.

**Explicitly out of scope**: compromised issuer keys, compromised agent runtime, OS/hardware compromise, SPIRE compromise, business policy, human authorization decisions.

## 9. Architecture Principles

Offline First · Cryptographically Verifiable · Deterministic · Fail Closed · Minimal Trust · Portable · Composable · Language Neutral

## 10. Owned Weaknesses (do not let any agent "discover" these as novel)

No hosted playground · No public SDK · No production users · No production deployment ·
No external integrations · No public API · No UI runtime · No ecosystem · No application examples

## 11. Biggest Risks

Adoption · Developer experience · Competing/evolving standards · Protocol evolution ·
Key rotation across partitions · Large-scale revocation distribution · Cross-language compatibility over time

## 12. Scale Targets

| Dimension | Current | Target |
|---|---|---|
| Topology | Single machine, lab | Millions of delegations, thousands of TPS, multiple orgs, multiple trust domains |

## 13. Career / Positioning Context (for agents 10 and 12 only)

Target roles: Staff/Principal-track Infrastructure, Platform, Distributed Systems, Security,
Identity & Access, Developer Infrastructure, Runtime Systems.
Institutional bar to clear: **CNCF Sandbox proposal quality.**

---

## SOURCING POLICY (BINDING ON ALL AGENTS)

Each persona is grounded in **artifacts that can be independently checked** — RFCs, published
books, standards documents, project governance criteria, well-documented engineering practices.

**Agents must never invent quotes, talk titles, timestamps, or URLs attributed to named people.**
If a persona needs authority, it cites the *artifact* (e.g. RFC 7515 §10, "Site Reliability
Engineering" Ch. 4 on SLOs, CNCF Sandbox criteria) rather than a fabricated utterance.
A fabricated citation in a security review is worse than no citation, because it launders
false confidence into a document people act on.
