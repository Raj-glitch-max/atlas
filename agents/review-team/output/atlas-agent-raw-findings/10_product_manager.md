# 10 — PRODUCT MANAGER — RAW FINDINGS

Phase 3 (Ops & Product). My job is deciding what NOT to build and answering the one question the
engineering roster avoids. I keep career framing out of this and hand it to 12.

FINDING-PM-1
  Type:        beachhead
  Claim:       There is no named beachhead user. "AI agents" is a market, not a customer.
  Evidence:    Owned weakness (context §10, LIMITATIONS §10). Filed here as its consequence (R11),
               not as the bare fact: every downstream claim — DX, sequencing, framework choice,
               revocation-vs-TTL emphasis — is an unvalidated hypothesis until three real teams
               with the delegation pain are interviewed.
  Recommendation: NOW — three beachhead interviews (S8). Worth more than the next 10k lines.
  Un-defer trigger: n/a (this IS the trigger for most of the roadmap).
  Owner:       10_product_manager

FINDING-PM-2
  Type:        cut
  Claim:       The planned scope (5 SDKs, 8 framework integrations, playground, marketplace,
               Terraform/Helm, operator console + marketing site) is a liability, not a roadmap.
  Evidence:    README already ships 3 SDKs, a UI, and an MCP server — arguably *ahead* of a
               validated need. That is scope spent before evidence.
  Recommendation: NOT-NOW for the matrix/marketplace/IaC (S14). Keep the *conformance vectors*
               (cheap, language-neutral, the real differentiator); defer further SDK breadth.
  Un-defer trigger: SDK/framework-specific demand from a named user.
  Owner:       10_product_manager

FINDING-PM-3
  Type:        build-vs-standard
  Claim:       The existential question — product, or contribution to UCAN/Biscuit/MCP? — is
               unanswered in writing.
  Evidence:    OQ-1. README says "closest to Biscuit/UCAN" but does not answer whether Atlas should
               *be* a thing or *feed* a thing.
  Recommendation: NOW — a written position, informed by PM-1's interviews.
  Owner:       10_product_manager

FINDING-PM-4
  Type:        sequencing
  Claim:       Revocation is presented as a headline pillar; short-TTL solves ~90% of cases and is
               the standard answer. This is a sequencing question, not a cut (04 owns the control).
  Evidence:    THREAT_MODEL C4–C7; LIMITATIONS §3 (short TTL as the replay mitigation).
  Recommendation: NEXT — document short-TTL as the *default* trust-refresh path; revocation as
               break-glass. Both ship (R5); the emphasis flips. Coordinate 04.
  Owner:       10_product_manager

FINDING-PM-5
  Type:        metric
  Claim:       No adoption-based success metric for 3/6/12 months.
  Recommendation: NOW — define "working": delegations verified in one real system, not stars.
  Owner:       10_product_manager

## SEQUENCING PROPOSAL (the cut list — becomes 12's public roadmap)
- NOW: 3 interviews (PM-1); differentiation + build-vs-standard doc (PM-3); LICENSE + governance
  (12); success metric (PM-5).
- NEXT: run the substrate proof (S7); short-TTL-as-default docs (PM-4); differential-in-CI (S9).
- LATER: harden the server (S11); pick ONE framework to integrate deeply.
- NOT-NOW: framework matrix, marketplace, Terraform/Helm, more SDKs (S14).

## MANDATORY CLOSING
- **Highest-value next action:** three beachhead interviews (S8). Everything else is built on an
  unexamined premise until they happen.
- **Is the problem real yet?** *Not proven.* Agent-to-agent delegation at scale is a credible bet
  on a future that may be 6 months or 4 years out. That is not fatal — but it must be *said*,
  because it makes "get one real user" outrank "write more runtime."
