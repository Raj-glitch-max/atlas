# AI Bootstrap Guide — Atlas

If you are an AI agent starting a new session on Atlas, read this file first.
It contains everything you need to understand the project in under a minute.

---

## 1. The Mission: Cross-Domain SPIFFE Delegation

Atlas enables secure, delegable authorization tokens for SPIFFE workloads
across independent, non-federated trust domains.

* **Core Requirement:** Token verification **must** be offline, local, and
  happen in under 100 ms with **zero** runtime network egress calls.
* **Core Hypothesis (H1):** Offline Relying Parties can verify nested `act`
  (actor) delegation chains using local SPIFFE trust bundles without broker
  communication.

---

## 2. Active Technical Blockers & Work in Progress

* **C4 Offline Revocation:** The primary technical risk is checking token
  revocation under network partitions.
* **Milestone 1 (Targeted Spike):** Evaluate three revocation candidates under
  `lab/spikes/`:
  1. OAuth Status List (`draft-ietf-oauth-status-list-10`).
  2. Signed Push-Revocation events.
  3. Cryptographic Accumulators.

---

## 3. Directory Layout & Where Things Live

```text
atlas/
├── AI_BOOTSTRAP.md         # You are here (AI entry point)
├── CODEGRAPH.md            # CodeGraph index documentation
├── LEVEL0_1_FEASIBILITY_GATE.md  # Root Anchor: pre-registration gate
├── P5_FALSIFICATION_EXPERIMENT.md # Root Anchor: Level 2 experiment spec
├── agents/                 # Persistent reasoning framework
│   ├── GOVERNANCE.md       # Rules binding agent interactions
│   └── journal/            # Journal entries recording all decisions
├── context/                # Canonical AI memory (9 files)
├── docs/                   # Planning specifications (frozen)
│   ├── product/            # Product requirements
│   └── engineering/        # Engineering requirements
├── lab/                    # Isolated local experiment laboratory
└── rfc/                    # Active draft architecture RFCs
```

---

## 4. Non-Negotiable Development Rules

1. **Frozen Document Guard:** All files in `docs/product/`, `docs/engineering/`,
   `docs/research/`, and the two root anchors are **frozen**.
   * Run `make check-frozen` before any commit to verify hashes.
   * Modifying frozen docs requires a governance cycle (documented in a
     journal entry under `agents/journal/` and updating `FROZEN.sha256`).
2. **Offline Verification Invariant:** No code paths may perform live runtime
   network calls to check issuer trust or revocation status during token
   verification.
3. **Traceability:** Every line of code must map to a requirement in
   `docs/engineering/01_ENGINEERING_REQUIREMENTS.md`.
4. **Vetted Tooling Only:** Use only standard Go (1.21+) libraries:
   `go-spiffe/v2`, `spire-api-sdk`, and `go-jose/v3`. No custom crypto forks.

---

## 5. Quickstart Commands

* Verify frozen documents integrity: `make check-frozen`
* Run style & lint checks: `make lint`
* Verify secret scan status: `make secrets`
* Rebuild CodeGraph index: `codegraph index`

<!-- checkpoint: governance(architecture-draft): refine architecture draft -->
