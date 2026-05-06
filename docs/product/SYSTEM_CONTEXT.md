# System Context

Defines the external boundary of the product: what surrounds it and interacts with it, without describing internal structure. See `PRODUCT_DEFINITION.md` for what the product is.

## External entities

- **Principal identity** — the identity (human, service, or organization) on whose behalf a delegation is made. Sourced from `TECHNICAL_VALIDATION.md` (P5, item 1, "acting on behalf of Y").
- **Delegate workload/agent** — the entity presenting proof of delegation. May be long-running or ephemeral/high-churn; the ephemeral case is a distinguishing requirement, not an edge case (`TECHNICAL_VALIDATION.md`, P5, item 2).
- **Relying party** — the entity verifying a presented delegation. May sit in the same trust domain as the principal, or in an independently-operated trust domain (`TECHNICAL_VALIDATION.md`, P5, item 7, minimum experiment).
- **Existing workload-identity infrastructure** — SPIFFE/SPIRE-based systems already deployed in the environment. This product must coexist with, not replace, this infrastructure (`ECOSYSTEM_THESIS.md`, P5, items 1, 4).
- **Existing standards bodies and drafts** — the 2026 IETF agent-identity drafts, the SPIFFE specification, and related OAuth extensions (RFC 8693, RFC 9396) are an external context this product exists alongside, not a component of it (`ECOSYSTEM_THESIS.md`, P5, items 1, 3).
- **Auditor / after-the-fact reviewer** (hypothesis — not evidenced as a validated user, only as a plausible consumer of output; see `USER_MODEL.md`) — an entity that inspects delegation records after the fact rather than at request time.

## What the product receives from its environment

- A request to establish, present, or verify a delegation (from a delegate workload or a relying party).
- Identity material already issued by existing workload-identity infrastructure (the product depends on this existing, not on issuing base identity itself — see `CONSTRAINTS.md`, C1).

## What the product produces into its environment

- A verifiable delegation record consumable by a relying party (`FUNCTIONAL_REQUIREMENTS.md`, FR1, FR5).
- A record sufficient to reconstruct a delegation chain after the fact (`FUNCTIONAL_REQUIREMENTS.md`, FR6).

## What is explicitly outside this boundary

- Issuance of base workload identity (SPIFFE/SPIRE's role, not this product's — `CONSTRAINTS.md`, C1).
- Any policy-decision logic about *what* a delegate should be allowed to do (a relying party or policy engine's role — this product proves *that* a delegation exists and its stated scope, not whether the scope should be granted).
- Any UI, developer tooling, or integration surface (all deferred per `DEFERRED.md`).

<!-- checkpoint: refactor(record): refactor revstatus cache driver -->

<!-- checkpoint: feat(internal): implement verification controller -->
