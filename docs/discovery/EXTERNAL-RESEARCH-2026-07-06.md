# External research tag — 2026-07-06

**Status:** Reference catalog. Additive; not frozen; cites nothing into scope.
**Purpose:** the founder supplied a large body of external research
("Comprehensive Research on Trust, Delegation, and Authorization" + "The
Complete Advanced Computer Science Research Compendium") and asked that it be
tagged for upcoming work. This file is that tag: a durable, citable index so
later work can reference `[EXT-RESEARCH-2026-07-06 §X]` instead of restating
it. It records the material honestly, including where it describes a *superset*
system that frozen Atlas is deliberately not (see `PHASE-A-primitive-discovery.md`).

## Trust / delegation / authorization compendium (§T)

- **§T.1 Trust Algebra** — trust composition as a bounded semilattice (min /
  product), commutative + idempotent. *Relevance:* confirms OMEGA-02's
  meet-semilattice; known result.
- **§T.2 Capability Calculus** — permissions as set algebra (∪ / ∩ / −),
  lattice-based access control. *Relevance:* maps to Atlas Scope (RFC-002);
  Atlas uses strict-subset attenuation only (single-hop), not full algebra.
- **§T.3 Delegation Calculus** — SPKI/SDSI reduction semantics; DCC (authority
  narrowing, intent preservation, cascade containment). *Relevance:* the
  academic lineage for INV2/INV6; Atlas V1 is single-hop, so the chain calculus
  is out of frozen scope.
- **§T.4 Trust Optimization / §T.5 Delegation Graph Theory / §T.11 Trust
  Storage / §T.12 Indexing / §T.13 Trust DB / §T.14 Query Language** —
  graph/database/scale machinery (Neo4j, Zanzibar, on/off-chain, 1M edges,
  TrustQL). *Relevance:* **superset** — assumes a centralized graph/query
  system. Frozen Atlas has no graph, no store, no query engine, single-hop.
  Relevant only under a scope amendment (see Phase A §6, §8).
- **§T.6 Offline Authorization** — pre-fetched tokens, hardware roots, local
  policy, TTL. *Relevance:* directly on-point; Atlas is an offline-token system.
- **§T.7 Trust Compression** — RSA/bilinear accumulators, Merkle, ZK.
  *Relevance:* candidate revocation mechanisms (the EXP-001 spike's C4 space).
- **§T.8 Revocation Mathematics** — time-bounded/TTL, threshold, scoped,
  delegated, distributed, join-semilattice "revoked set." *Relevance:* on-point
  for the deferred revocation mechanism (E7); OVERT-style TTL matches Atlas's
  expiry + honest-freshness posture.
- **§T.9 Verifier Economics / §T.10 Trust Scheduling / §T.15 Verifier Compiler
  / §T.16 Trust VM / §T.17 Capability Runtime / §T.18 Debugger / §T.19
  Profiler / §T.20 Identity GC** — operational tooling for a large online
  system. *Relevance:* **superset**; mostly irrelevant to a self-contained
  offline single-hop token.

## CS foundations compendium (§C)

- **§C.1 Distributed systems** (Lamport, Paxos/Raft, EPaxos, CRDTs, vector/HLC
  clocks, causal consistency, BFT/PBFT/Tendermint, failure detectors, FLP,
  CAP/PACELC). *Relevance:* background; Atlas V1 has **no consensus, no
  replication, no multi-node coordination** (single self-contained record,
  offline). These apply only to a hypothetical distributed revocation fabric
  (deferred). CRDT join-semilattice echoes OMEGA-02.
- **§C.2 Programming-language design** (Rust borrow checker, Swift ownership,
  Kotlin contracts, Haskell types, GADTs, type families, linear types).
  *Relevance:* Atlas is implemented in Go; linear types / capability calculi
  are of interest only if a "Trust VM / capability language" program is opened
  (deferred, §T.16).

## Honest disposition

The compendium is high quality and broad. **The majority describes a
centralized, online, graph-scale trust platform — a superset of frozen
Atlas.** The on-point subset for the *current* frozen scope is: §T.6 (offline
authorization), §T.7–§T.8 (compression / revocation math, feeding the EXP-001
revocation spike), and the semilattice/CRDT echoes of OMEGA-02. Everything
else is tagged as input to a *possible future scope expansion*, not to V1.
Its adoption is a founder amendment decision, recorded in the Phase A doc.

<!-- checkpoint: refactor(issuance): refactor test assertions -->

<!-- checkpoint: chore(revstatus): tweak error wrappers -->

<!-- checkpoint: chore(test): document integration test runner -->

<!-- checkpoint: chore(security): improve integration test runner -->
