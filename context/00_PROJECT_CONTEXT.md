# Project Context — Atlas

## 1. Project Mission & Identity

**Atlas** is a technical venture studio operating in **pre-direction exploration**. It is designed as a persistent, adversarial reasoning workspace to evaluate startup directions and engineering problems before any product, business model, codebase, or stack-specific architecture is generated.

The workspace is built around a persistent technical reasoning framework (`agents/`) that structures adversarial review to prevent emotional attachment to weak ideas, enforce evidence discipline, and kill bad directions quickly.

## 2. Core Thesis

The core thesis of Atlas is that **the problem space must precede the solution space**. We do not build products or startups directly; instead, we systematically discover, rank, and evaluate critical infrastructure and platform engineering pain points.

Our methodology tests candidate problems against:
1. **Evidence Quality** — requiring concrete citations and independent data points, not vendor-marketing claims.
2. **Ecosystem Gravity** — analyzing whether a solved problem would naturally pull independent systems to depend on it (e.g., standard workload delegation vs. isolated features).
3. **Friction Penalties** — evaluating organizational, technical, seller, and buyer constraints that make execution or go-to-market difficult for a small team.

## 3. Pipeline Position

As of the last R1 cycle, the exploration pipeline is at the following stage:

| Stage | Status | Notes |
|---|---|---|
| **Domain Research** | Completed | Identified 10 key pain points in infrastructure/SRE. |
| **External Repository Audit** | Completed | Audited `claude-skills/`, indexed 7 items, deleted raw folder. |
| **Problem Ranking Framework** | Completed | Established 10-dimensional opportunity/friction scoring. |
| **Problem Ranking (R1 Cycle)** | Completed | Ran scoring and council reviews for the 10 pain points. |
| **Product Management (R2 Cycle)** | **Not Started** | Deferred until founder selects a candidate from the R1 shortlist. |

## 4. Candidate Problems & Portfolio Reduction

A comprehensive review of the 10 initial candidate problems resulted in a formal portfolio reduction, narrowing focus to the most viable and defensible engineering spikes:

### Tier A — Research Immediately (Shortlist)
*   **P5: IAM least-privilege / service-to-service identity** (Net: 5, lifted by Cartographer)
    *   *Thesis:* A standardized binding of workload identity to attenuable, offline-verifiable delegation capabilities portable across trust domains.
    *   *Why:* Strong technical novelty, alignment with SPIFFE/SPIRE, and a clear path to two-week spike validation.
*   **P10: Production debugging / root-cause speed** (Net: 4)
    *   *Thesis:* A standardized, portable inferred causal edge for non-call-graph-connected causation.
    *   *Why:* High technical depth; requires immediate validation against a real trace dataset to verify if a reliable populating mechanism is scientifically feasible.

### Tier B — Keep as Backup
*   **P9: On-call burden & toil reduction** (Net: 6)
    *   *Why:* Strongest non-vendor evidence base but faces small-ACV/services-business revenue-model risks.

### Tier C — Reject Before Research
*   **P1: Observability cost & signal-to-noise** (Net: 5) — Rejected due to dependency on uncontrollable storage economics.
*   **P2: CI/CD flakiness & queue contention** (Net: 6) — Rejected due to high incumbent strength and lack of a clear niche wedge.
*   **P3: Secrets & config drift across environments** (Net: 5) — Rejected; drift detection is a feature, not a category, and is being absorbed natively.
*   **P4: Cloud cost attribution & reduction** (Net: 6) — Rejected; requires client tagging discipline and consultative enterprise sales.
*   **P6: Compliance & audit evidence collection** (Net: 1) — Rejected; low technical novelty, high risk of collapsing into a thin AI wrapper.
*   **P7: Kubernetes operational complexity** (Net: 0, dropped by Cartographer) — Rejected; pain is a value-promise gap, not missing tooling.
*   **P8: Incident response coordination overhead** (Net: 3) — Rejected; market saturated, weak vendor-dominated evidence.

## 5. Active Exclusions & Constraints

*   **No Product Generation:** We do not design products, business models, or architectures in this workspace until upstream validity gates are passed.
*   **No Stack Commitment:** Selection of languages, frameworks, test runners, or deployment runtimes is deferred until feasibility spikes are completed.
*   **No Redesign:** The `agents/` framework is frozen at v1.0. No architectural refactoring of the reasoning council is permitted.
