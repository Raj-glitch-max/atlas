# 02 — SOLUTIONS ARCHITECT — RAW FINDINGS

Phase 3 (Ops & Product). I own trust topology, distribution, cost, and the adoption path.

FINDING-ARCH-1
  Type:        differentiation
  Claim:       The "why Atlas is not SPIFFE/SPIRE, not OpenFGA, not UCAN/Biscuit, not plain JWT"
               answer is scattered across README, WHY.md, and OMEGA docs but is not one page a
               reviewer or an architect can consume in 60 seconds.
  Evidence:    README correctly positions "closest to Biscuit/UCAN, companion to SPIFFE"; the
               fuller argument lives in docs/discovery/OMEGA-*. It needs consolidation.
  Impact:      This is the first question a CNCF TOC and any evaluating architect asks. Diffuse ≈ absent.
  Recommendation: NOW — one-page differentiation doc (feeds 12/EM-3, S4).
  Owner:       02_solutions_architect

FINDING-ARCH-2
  Type:        topology / distribution
  Claim:       The trust topology (one issuing domain, one relying domain, out-of-band snapshot +
               bundle distribution) is coherent and its cost is honestly bounded — but the
               distribution channel that makes it work at scale is unbuilt and unmeasured. The
               two-domain SPIRE substrate proof is authored (`atlas-lab/`) and unrun (LIMITATIONS §8).
  Evidence:    `revstatus` signed-set realization is contract-tested; the "passive signed-blob
               distributor" is defined (journal 2026-07-06). No live two-domain run exists.
  Impact:      The topology's payoff (offline verification across a partition) is asserted by
               construction, not demonstrated at link level. This is the same gap 04/06 orbit.
  Recommendation: NEXT — run the substrate proof (S7); it is the architecture's load-bearing demo.
  Owner:       02_solutions_architect

FINDING-ARCH-3
  Type:        topology
  Claim:       Two-domain, single-hop is a deliberate, defensible scope — but the topology only
               *pays off* at a multi-org scale that may never arrive. That is a strategic bet, not
               a technical flaw; it should be named as one.
  Evidence:    LIMITATIONS §1–2; scale target "millions of delegations, multiple orgs" (context §12).
  Recommendation: LATER — federation productization; un-defer when a second org actually asks (OQ-1).
  Owner:       02_solutions_architect

VERDICT: The topology is sound and honestly scoped. Its adoption path depends on (a) the one-page
differentiation doc and (b) the substrate proof turning a by-construction claim into a demo.
