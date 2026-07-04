# EXP-001 Execution Plan — C4 Feasibility Spike

**Type.** The final planning artifact before engineering begins. A concrete task breakdown for executing EXP-001 (the C4 Feasibility Spike). It is not a pre-registration, not a lab process document, not code, and not protocol design. It names work, effort, dependencies, tools, risks, and exit criteria for each task.

**What this plan is NOT (binding).** No implementation code. No protocol design. No architecture. No product discussion. Where a task necessarily involves *writing* something (a config file, a manifest entry, a thin revocation layer atop stock libraries), the plan names the task and its constraints; it does not contain the artifact. The implementation may add new code *atop* stock libraries; it may not modify them (per `P5_FALSIFICATION_EXPERIMENT.md` S7/F7).

---

## 0. Source of authority (frozen; do not modify)

This plan executes two frozen upstream documents and operates inside five lab process documents. Nothing in this plan re-scopes, redesigns, or re-interprets them.

| Source | Role | Path |
|---|---|---|
| Pre-registration (spike question, candidate compositions) | Defines **what** EXP-001 tests | `LEVEL0_1_FEASIBILITY_GATE.md` §"Spike question (C4-only)" |
| Level 2 protocol (test cases, criteria, substrate locks) | Source of T-case definitions + S/F criteria referenced by the spike | `P5_FALSIFICATION_EXPERIMENT.md` (frozen) |
| Go/no-go gate + stop conditions | Defines **how** a run is gated and halted | `lab/EXPERIMENT_CHECKLIST.md` |
| Evidence registry | Defines **where** artifacts live and how they are hashed | `lab/EVIDENCE_INDEX.md` |
| Decision rules | Defines **what conclusion** the evidence permits (outcomes α/β/γ/δ) | `lab/DECISION_RULES.md` §3 |
| Lab root + lifecycle | Defines the run/reproduction/journal flow | `lab/LAB_README.md` |
| Observation log | Where observations are appended, never conclusions | `lab/EXPERIMENT_LOG.md` |

**Verbatim spike question (the pre-reg; not paraphrased):**

> On a SPIFFE substrate with RFC 8693 tokens, is there *any* composition of existing standardized primitives (stock JWT verifier + a thin, non-patching revocation layer) that delivers revocation observability within R seconds of revocation — for R as chosen in S1 — with no live call from the RP to domain A and no third-party decision-making broker, *under the S2/S3/S4 read of H1*?

**Candidate compositions (pre-registered, cheapest-first; order is binding):**
1. OAuth Status List with explicit refresh cadence set to R — admissible only if the S2 scope-act admits cached pulls; otherwise excluded by construction.
2. Push-revocation: domain A signs-and-pushes a small revocation event to the RP's inbound listener; the RP treats absence-of-push as "still valid."
3. Accumulator-based revocation: pre-loaded dynamic accumulator updated periodically; the RP verifies membership/non-membership locally.

---

## 1. The blocking prerequisite — S1–S5 scope-act gate (founder, not engineering)

**This is the gate that precedes pre-flight.** The spike question is parameterized by R (S1) and conditioned on the S2/S3/S4 reading of H1. Until those are resolved as semantic scope acts, the spike's success criterion is undefined and `EXPERIMENT_CHECKLIST.md` P-5 ("no threshold is TBD") cannot pass. This is not engineering work and is not in this plan's effort budget; it is named here because **nothing below may start until it is closed.**

| Scope item | To be resolved by founder as a scope act | Effect on the spike if unresolved |
|---|---|---|
| **S1** | Choose R (revocation-observability latency bound, R < delegation TTL = 30 s). The Level 2 protocol's T4 implies R ≈ 2 s; the founder either confirms 2 s or picks another value. | The spike's pass/fail threshold is undefined. P-5 fails; no Go. |
| **S2** | Decide (a) admit periodic cached pulls of signed integrity-protected artifacts (bound staleness), or (b) forbid all RP-initiated fetches from domain A and require push. | Composition 1 (Status List) is either admissible or excluded by construction. Admissibility set of candidates is undefined. |
| **S3** | Define "third-party broker" as an entity that makes or vouches for trust decisions, explicitly excluding passive signed-blob caches that perform no decision. | Whether a CDN/list host counts as a broker is undefined; composition 1's compliance is undecidable. |
| **S4** | Decide whether H1 requires observability of revocations *performed during a partition that isolated the RP from the issuer at revocation time*, or whether revocation observability is bounded to "eventually consistent upon partition recovery, within P seconds of recovery." | Distinguishes outcome β (technology gap) from outcome γ (logical impossibility). The combined T4∩T5 reading under a strict S4 is information-theoretically unsatisfiable; running the spike against it would force outcome γ by construction. |
| **S5** | Not required for the spike (S5 governs T7, multi-hop, which is C4×N and out of the spike's C4-only scope). Resolve only if the spike fails → outcome β → Level 2 is justified and narrows to C4 + C11's per-hop sub-claim. | None for EXP-001. |

**Exit criteria for this gate (all four required before Phase 1):**
- A journal entry at `agents/journal/<date>-c4-spike-scope-act.md` records the chosen values of R, S2-reading, S3-definition, and S4-reading, with confidence labels and change-conditions per `agents/GOVERNANCE.md`.
- Those values are reflected into a run-specific addendum to the spike's pre-reg section (a pre-reg amendment, dated, not a retrofit of the frozen gate) so that P-2/P-3/P-5 can witness concrete non-TBD thresholds.
- The Standards Editor sights the amendment and confirms it does not retroactively alter the frozen gate (it scopes a previously-ambiguous parameter; it does not redefine H1).

**Risk if skipped:** the spike runs against an undefined threshold, every conclusion is refused at H-3 (post-run halt: threshold absent from pre-reg), and the run is void. There is no path to a valid conclusion that bypasses this gate.

---

## 2. Effort summary

Effort is in engineering-hours (eng-h), single Engineer, exclusive of the founder's scope-act work above and exclusive of reproduction (Phases 11) and decision/journal (Phase 12), which are listed separately. Honesty note: `LEVEL0_1_FEASIBILITY_GATE.md` estimated "2–3 days" for the C4 compositional question *assuming a substrate already existed*. This plan includes substrate build, instrumentation, and reproduction per lab discipline, which the audit's estimate did not. The two are not in conflict; they price different scopes.

| Phase | Tasks | Effort (eng-h) | Notes |
|---|---|---|---|
| 1 — Pre-flight groundwork | T1.1–T1.3 | 2 | Manifest + hash-pinning |
| 2 — Environment setup | T2.1–T2.4 | 3 | 3 hosts, NTP, disk, version manifest |
| 3 — Dependency installation | T3.1–T3.6 | 4 | Stock-only; diff-verified |
| 4 — SPIRE deployment | T4.1–T4.6 | 6 | Two independent deployments |
| 5 — Network isolation | T5.1–T5.3 | 3 | Firewall fire-and-reconnect |
| 6 — Instrumentation | T6.1–T6.4 | 4 | Sniffer + latency probe + hashing |
| 7 — Personnel + roles | T7.1–T7.3 | 1 | Blinding, reproduction readiness |
| 8 — Pre-flight Go/No-Go | T8.1–T8.2 | 1.5 | Checklist + GO.json |
| 9 — Composition attempts | T9.0–T9.3 | 18–26 | Cheapest-first; stop at first α |
| 10 — Evidence collection | (continuous) | 2 | Woven through 9; called out |
| 11 — Reproduction | T11.1–T11.2 | 7 | Second engineer, fresh substrate |
| 12 — Decision + journal | T12.1–T12.4 | 3 | Outcome class + audit + journal |
| **Total** | | **54.5–62.5 eng-h** | ~7–8 eng-days single-engineer; +reproduction block |

The dominant cost is Phase 9 (composition attempts) and Phase 11 (reproduction). Phase 9's range reflects: composition 1 may be excluded by the S2 scope-act (cheaper), or composition 3 may have no standardized OSS implementation (excluded, also cheaper); the wide band is accumulator work if all three are attempted.

---

## 3. Dependency graph

```
Phase 0 (scope-act gate) ── blocks ──> everything
        │
        ▼
Phase 1 (pre-flight groundwork: hash-pin pre-reg, init manifest)
        │
        ├─> Phase 2 (environment) ──> Phase 3 (dependencies)
        │                                 │
        │                                 ├─> Phase 4 (SPIRE deploy)
        │                                 │       │
        │                                 │       └─> Phase 5 (network isolation)
        │                                 │               │
        │                                 └───────────────┴─> Phase 6 (instrumentation)
        │                                                         │
        └─> Phase 7 (personnel/roles) ─────────────────────────────┤
                                                                      ▼
                                                              Phase 8 (Go/No-Go)
                                                                      │
                                                                      ▼
                 Phase 9 (composition attempts: RUN-01 → RUN-02 → RUN-03, cheapest-first, stop at α)
                                                                      │
                          (continuous) Phase 10 (evidence collection) ─┤
                                                                      ▼
                                                              Phase 11 (reproduction)
                                                                      │
                                                                      ▼
                                                              Phase 12 (decision + journal)
```

**Critical path:** 0 → 1 → 2 → 3 → 4 → 5 → 6 → 8 → 9 → 11 → 12.
**Parallelizable:** Phase 7 (personnel/blinding) can run alongside Phases 2–6. Phase 3's tool installation (T3.1–T3.6) is parallelizable across hosts once Phase 2 is done.
**Anti-pattern (forbidden):** starting any Phase-9 composition work before Phase 8 GO.json is signed. This is the canonical garden-of-forging-paths entry point (`DECISION_RULES.md` U5); an attempt before Go is an unlogged attempt, which the lab refuses.

---

## 4. Phase-by-phase task breakdown

### Phase 1 — Pre-flight groundwork

#### T1.1 — Hash-pin the frozen pre-reg and decision rules
- **Effort:** 0.5 eng-h
- **Depends on:** Phase 0 closed.
- **Required tools:** `sha256sum`, the run manifest template (`EVIDENCE_INDEX.md` §4).
- **Risks:** Pre-reg drift between hash-pinning and run start (mitigated: pin at the start of Phase 1, re-verify at Phase 8). The frozen gate must not be edited to receive the S1–S5 amendment — the amendment is a separate dated addendum referenced by the manifest.
- **Exit criteria:** Manifest `preRegistration.sha256` and `decisionRulesVersion.sha256` filled (P-2, P-4); both hashes re-computed and matched against the recorded values.

#### T1.2 — Copy H1 + spike question verbatim into the run's first log entry
- **Effort:** 0.5 eng-h
- **Depends on:** T1.1.
- **Required tools:** `EXPERIMENT_LOG.md` template.
- **Risks:** Paraphrase drift (mitigated: literal copy; the Standards Editor sights the verbatim text at pre-flight).
- **Exit criteria:** `EXPERIMENT_LOG.md` has a new SETUP entry (the run's E-NNNN) containing the verbatim H1 from `P5_FALSIFICATION_EXPERIMENT.md` §1 and the verbatim spike question from the gate (P-6).

#### T1.3 — Initialize run directory + manifest template
- **Effort:** 1 eng-h
- **Depends on:** T1.1.
- **Required tools:** `mkdir`, the `MANIFEST.json` schema (`EVIDENCE_INDEX.md` §4).
- **Risks:** Premature over-creation of unused type-directories (e.g., `screenshot/`); per `EVIDENCE_INDEX.md` §2, directories are created per-run and absence is meaningful. Create only `config/`, `pcap/`, `log/`, `bench/`, `token/`, `raw/`, `analysis/` (screenshot/recording not anticipated by the spike).
- **Exit criteria:** `lab/evidence/EXP-001/RUN-01/` exists with a started (unsealed) `MANIFEST.json` whose timestamp precedes any measurement (V-1, V-2, R-5).

### Phase 2 — Environment setup

#### T2.1 — Provision three hosts
- **Effort:** 1 eng-h
- **Depends on:** Phase 1.
- **Required tools:** Whatever provisioning the substrate uses (VMs or containers on a single laptop, per `P5_FALSIFICATION_EXPERIMENT.md` §4 — but network isolation must use **host-level firewall rules, not container network namespaces**, because containers share kernel/trust state).
- **Risks:** Container-network-namespace isolation (mitigated: use host firewall; called out in the protocol). Resource contention if all three on one laptop (mitigated: 4 vCPU/8 GiB minimum per host; the spike is latency-sensitive at the 100 ms budget which is two orders of magnitude of headroom, so contention is a low risk to S1 but a real risk to clean measurement).
- **Exit criteria:** Three hosts up: `host-A` (SPIRE domain A), `host-B` (SPIRE domain B), `host-RP` (relying party). Config snapshot artifact produced (E-1).

#### T2.2 — NTP-discipline clocks; measure inter-host skew
- **Effort:** 1 eng-h
- **Depends on:** T2.1.
- **Required tools:** `chrony` (or `ntpd`), skew measurement (e.g., `chronyc tracking`, pairwise offset).
- **Risks:** Skew drift during the run exceeding the pre-reg's 250 ms bound (stop-condition A-1). This is a low-risk-but-load-bearing item: revocation-observability timing relative to T₀+5s/T₀+7s requires clocks to mean the same thing across hosts.
- **Exit criteria:** All three clocks NTP-disciplined; pairwise skew artifact produced and recorded (E-3); skew within the pre-reg bound. Manifest `bench/` entry for the skew measurement.

#### T2.3 — Disk capacity check at evidence root
- **Effort:** 0.25 eng-h
- **Depends on:** T2.1.
- **Required tools:** `df -h lab/evidence/`.
- **Risks:** Filling the capture root mid-run (stop-condition A-6 via write error). The spike's capture volume is small (no full-run pcaps, only T4/T5 case excerpts), so 2× margin is easy.
- **Exit criteria:** Free space ≥ 2× pre-reg max-capture-size (E-5); `df` output recorded.

#### T2.4 — Build the version manifest
- **Effort:** 0.75 eng-h
- **Depends:**** T2.1.
- **Required tools:** `git`-pinned checkouts or release-tarball hashes for every component installed in Phase 3.
- **Risks:** Floating tags (`latest`) silently pulling a newer (and potentially S7-violating) library version. This is the root of the interoperability gate (S7/F7): a floating tag that upgrades a verifier mid-spike could invalidate a prior pass.
- **Exit criteria:** Version manifest artifact lists every component with a pinned version + upstream hash (E-2); no `latest` anywhere.

### Phase 3 — Dependency installation (stock only)

**Cross-cutting constraint for T3.1–T3.6:** every installation is diff-verified against its upstream tag. Any non-empty diff fails E-10 and is a no-go. The Engineer may add new code *atop* stock; the Engineer may not modify stock source (S7/F7). A stop-condition A-4 (boundary breach) fires if installation drifts into patching stock.

#### T3.1 — Install SPIRE server + agent (both domains)
- **Effort:** 1 eng-h
- **Depends on:** T2.4.
- **Required tools:** SPIRE release artifacts (2026-vintage stable, pinned), `sha256sum`, `diff`/`git verify-tag`.
- **Risks:** A "convenience" build flag that re-enables federation by default (mitigated: explicit federation-disabled config in Phase 4; E-8 verifies the config lines, not just the binary).
- **Exit criteria:** SPIRE installed on `host-A` and `host-B`; versions pinned in the manifest; diff against upstream tag empty (E-10).

#### T3.2 — Install the stock JWT/JWS/COSE verifier
- **Effort:** 0.5 eng-h
- **Depends on:** T2.4.
- **Required tools:** Pinned verifier (e.g., a stock go-jose / nimbus / equivalent JWT-SVID verifier; the choice is an instrumental one the Engineer records, not a protocol decision), `diff`/`git verify-tag`.
- **Risks:** A verifier convenience wrapper that performs an implicit network call for key resolution (a JWKSet fetcher). The RP's verifier must resolve keys only from the pre-loaded domain-A and domain-B bundles — no fetch-on-miss. An implicit fetch is a latent F2/F4 (sniffer hit outside domain B).
- **Exit criteria:** Verifier installed, pinned, diff-empty (E-10); a logged check that its key-resolution path is bundle-only, no network fallback.

#### T3.3 — Install the RFC 8693 token-exchange client library
- **Effort:** 0.5 eng-h
- **Depends on:** T2.4.
- **Required tools:** Pinned RFC 8693 client, `diff`/`git verify-tag`.
- **Risks:** A library convenience that calls a "token introspection" endpoint (RFC 7662) as a verification step — that is a live call and would falsify the spike on the no-live-call clause. The client is used for *issuance* of the delegation token in domain A; verification at the RP uses the stock verifier (T3.2), not the token-exchange client.
- **Exit criteria:** RFC 8693 client installed on domain-A side, pinned, diff-empty (E-10).

#### T3.4 — Install the OAuth Status List library (composition 1) — *conditional*
- **Effort:** 0.75 eng-h
- **Depends on:** T2.4 **and the Phase-0 S2 scope-act.**
- **Required tools:** Pinned `draft-ietf-oauth-status-list` implementation (if S2 admits cached pulls); `diff`.
- **Risks:** (a) S2 forbids cached pulls → skip this task entirely and log it as "excluded by scope-act" (this is honest evidence, not a gap — U5 demands the count of attempts include excluded-by-scope-act candidates). (b) The status-list impl requires a fetch from the issuer's URI even when the list is cached → that is a live call to domain A, F2/F4; the composition would fail the sniffer test, which is the spike's job to surface.
- **Exit criteria if S2 admits:** Library installed, pinned, diff-empty (E-10); or a logged exclusion entry if no standardized impl exists (evidence for outcome β). **Exit criteria if S2 forbids:** A log entry records composition 1 as excluded-by-scope-act; RUN-01 (Phase 9) is skipped.

#### T3.5 — Install capture instruments (sniffer + latency probe)
- **Effort:** 1 eng-h
- **Depends on:** T2.4.
- **Required tools:** Out-of-band sniffer (e.g., `tcpdump`/libpcap on a separate capture interface or `iptables LOG` mirror — *not* in-band with the verifier, so the instrument cannot mask the very calls it is meant to detect); a latency probe instrumentation hook at the RP.
- **Risks:** An in-band sniffer that shares the verifier's network path could drop or reorder the very packets it must witness. The sniffer must be out-of-band and write directly to the run's `pcap/` directory (V-3).
- **Exit criteria:** Sniffer installed on `host-RP`, configured to capture any TCP/UDP from the RP to addresses outside `domain-B.prefix`, writing to `evidence/EXP-001/RUN-NN/pcap/`; latency-probe hook present (E-6 prep).

#### T3.6 — Install hashing + signing tooling
- **Effort:** 0.25 eng-h
- **Depends on:** T2.4.
- **Required tools:** `sha256sum`, `age` or `gpg` (for `config` and `GO.json` signatures per `EVIDENCE_INDEX.md` §3/§5).
- **Risks:** A signature scheme whose signing key is not separable from the run (the Engineer's signing key must be identified so a second engineer reproducing can verify signatures without the original key — mitigate: signatures are by role-initial, recorded in the manifest, and reproduction verifies the signature against a published role key, not a private one).
- **Exit criteria:** Hashing + signing tooling present; sample-hashing test passes (V-4).

### Phase 4 — SPIRE deployment

#### T4.1 — Deploy SPIRE domain A
- **Effort:** 1.5 eng-h
- **Depends on:** T3.1.
- **Required tools:** SPIRE server + agent config (recorded as a `config/` artifact).
- **Risks:** Reusing a default database or signing key across domains (mitigated: distinct DB, distinct keys; E-7 witnesses key fingerprints differing).
- **Exit criteria:** SPIRE domain A running on `host-A` with its own DB and signing key; config snapshot produced and signed (E-1, E-7).

#### T4.2 — Deploy SPIRE domain B
- **Effort:** 1.5 eng-h
- **Depends on:** T3.1, T4.1.
- **Required tools:** SPIRE server + agent config.
- **Risks:** Same as T4.1 plus the trap of "copy domain A's config and tweak" — which can inherit federation defaults. Domain B is built from the same template but with distinct keys (E-7) and federation disabled (E-8).
- **Exit criteria:** SPIRE domain B running on `host-B`; key fingerprints differ from domain A (E-7); federation disabled (E-8).

#### T4.3 — Verify federation disabled in both deployments
- **Effort:** 0.5 eng-h
- **Depends on:** T4.1, T4.2.
- **Required tools:** SPIRE config inspection.
- **Risks:** Federation enabled by a default the Engineer didn't explicitly disable. Federation is the closest SPIRE-shipped mechanism that would smuggle a shared broker back in (`P5_FALSIFICATION_EXPERIMENT.md` §3.2); leaving it on would let a composition "succeed" by cheating.
- **Exit criteria:** Disabled-federation config lines cited in the manifest (E-8); both deployments confirmed federation-off.

#### T4.4 — Register workloads and principals
- **Effort:** 1 eng-h
- **Depends on:** T4.1, T4.2, T4.3.
- **Required tools:** `spire` CLI registrations.
- **Risks:** Registering principal Y in domain B as an *active workload* (the protocol specifies Y is registered in domain B but **not present as an active workload** — registering it active would change what is being tested).
- **Exit criteria:** Workload X registered and running in domain A; principal Y registered in domain B (not active); RP registered and running in domain B; registrations recorded in config snapshot.

#### T4.5 — Manual trust-bundle exchange into the RP verifier
- **Effort:** 0.5 eng-h
- **Depends on:** T3.2, T4.1, T4.2.
- **Required tools:** SPIRE bundle export (`spire bundle show`), the stock verifier's key-loading interface.
- **Risks:** The verifier auto-fetching bundles (a federation-like back-channel) instead of using the manually-loaded bundles. The RP must verify with whatever bundles it last received — no implicit live call (T6/S5 territory, but the bundle *loading* discipline applies from setup).
- **Exit criteria:** Domain A and domain B public bundles exported, loaded into the RP's stock verifier, loading recorded; verifier confirmed to use loaded bundles only (E-1).

#### T4.6 — Lock TTLs
- **Effort:** 0.5 eng-h
- **Depends on:** T4.4.
- **Required tools:** Token-issuance config (delegation TTL), SPIRE SVID TTL config (Y-SVID TTL).
- **Risks:** Mid-spike TTL change to make a case pass — this is exactly F10 (tunable threshold) and A-3 (parameter not in pre-reg). TTLs are frozen at pre-reg, embedded in the config snapshot, and the manifest witnesses the values.
- **Exit criteria:** Delegation TTL = 30 s, Y-SVID TTL = 600 s, both embedded in the config snapshot artifact (E-11).

### Phase 5 — Network isolation

#### T5.1 — Define the RP outbound firewall ruleset (inactive at setup)
- **Effort:** 1 eng-h
- **Depends on:** T4.5.
- **Required tools:** Host firewall (`iptables`/`nftables` on `host-RP`).
- **Risks:** Rules active at setup that break initial bundle exchange or smoke tests (mitigated: rules defined but inactive; E-9 requires inactive-at-setup, activated per-test only).
- **Exit criteria:** Ruleset artifact produced (allow domain B; deny domain A + any plausible broker), inactive at setup (E-9).

#### T5.2 — Verify reachability at setup
- **Effort:** 0.5 eng-h
- **Depends on:** T5.1.
- **Required tools:** Network reachability probe (`ping`/`curl` to each host's expected endpoints).
- **Risks:** A host unreachable at setup that is then blamed on the partition test (the *fire-and-reconnect* workflow tests partition-on-demand; a setup-time unreachable host is a substrate fault, not a test result).
- **Exit criteria:** All hosts reachable at setup per the pre-reg topology; probe output recorded (E-4).

#### T5.3 — Dry-run the firewall fire-and-reconnect workflow
- **Effort:** 1.5 eng-h
- **Depends on:** T5.1, T5.2.
- **Required tools:** Firewall activation/revert scripts.
- **Risks:** Revert failure leaving the RP partitioned for a subsequent case (stop-condition A-8's reverse: a contaminated *next* case). The fire-and-reconnect cycle must be reliable and fast, because T5's "full 30-second partition" is timed and a slow revert corrupts the measurement window.
- **Exit criteria:** Activate → confirm deny → revert → confirm allow, executed cleanly on a non-test target, logged as a dry-run artifact (R-4 evidence; also de-risks Phase 9 timing).

### Phase 6 — Instrumentation

#### T6.1 — Arm and self-test the out-of-band sniffer
- **Effort:** 1.5 eng-h
- **Depends on:** T3.5, T5.1.
- **Required tools:** Sniffer, a known-good test packet (a deliberate outbound packet from the RP to a domain-B address and one to a non-domain-B address).
- **Risks:** The sniffer failing to capture the very packet class it must catch (the falsifier for F2/F4 is "sniffer records an outbound call outside domain B" — a sniffer that misses silently turns a FAIL into a false PASS). Mitigated: self-test on a known-good packet (E-6) where the sniffer is confirmed to record both an in-domain-B and an out-of-domain-B packet, then the out-of-domain-B one and only that one.
- **Exit criteria:** Sniffer self-test artifact (ART-0000-class) produced; sniffer confirmed to record out-of-domain-B traffic and ignore in-domain-B traffic; recorded (E-6).

#### T6.2 — Stand up the latency probe
- **Effort:** 1 eng-h
- **Depends on:** T3.5.
- **Required tools:** ms-precision timing instrumentation at the RP (token-received → verdict-returned).
- **Risks:** Timer resolution coarser than the budget (the budget is 100 ms; the timing must resolve comfortably below that, e.g., sub-ms). Clock-source confusion (using a different clock than the system clock the manifest timestamps).
- **Exit criteria:** Latency probe records token-received and verdict-returned timestamps in ms; a sample measurement recorded to `bench/` (S1/F1 instrument).

#### T6.3 — Test the hashing procedure on a sample artifact
- **Effort:** 0.5 eng-h
- **Depends on:** T3.6, T1.3.
- **Required tools:** `sha256sum`, the manifest writer.
- **Risks:** Hash-at-creation-vs-hash-at-end drift (mitigated: a sample artifact hashed, manifested, re-read, re-hashed, matched — per `EVIDENCE_INDEX.md` §5 step 4).
- **Exit criteria:** Sample artifact created, hashed, manifest entry appended, re-read hash matches (V-4).

#### T6.4 — Manifest-writer readiness check
- **Effort:** 0.5 eng-h
- **Depends on:** T1.3, T6.3.
- **Required tools:** Manifest writer.
- **Risks:** A first artifact produced before the manifest writer is ready → an unlogged artifact → invalid evidence (the artifact "does not exist" per `EVIDENCE_INDEX.md` §4).
- **Exit criteria:** Manifest-writer check artifact recorded confirming the first real artifact will be appended immediately (V-5).

### Phase 7 — Personnel and role separation

#### T7.1 — Assign Engineer, Adversary, Standards Editor
- **Effort:** 0.25 eng-h
- **Depends on:** Phase 1.
- **Required tools:** Role log.
- **Risks:** Same person holding Engineer and sole reviewer (H-1 violation); Standards Editor having run the experiment (H-4 violation). For a lab this small, role separation is the most likely procedural slip.
- **Exit criteria:** Roles assigned and logged (H-1, H-2, H-4); Standards Editor confirmed distinct from Engineer.

#### T7.2 — Adversary threshold-blinding
- **Effort:** 0.25 eng-h
- **Depends on:** T7.1.
- **Required tools:** Sighted-initials confirmation.
- **Risks:** The Adversary seeing the S/F thresholds in `DECISION_RULES.md` §3 or the pre-reg's S/F criteria before running (H-3 violation) — invalidates adversary-blinded cases. Re-confirmation before each case (H-5).
- **Exit criteria:** Adversary sighted-initials confirming blindness; per-case re-confirm plan logged (H-3, H-5).

#### T7.3 — Confirm reproduction readiness
- **Effort:** 0.5 eng-h
- **Depends on:** T7.1, R-1 (the config snapshot from Phase 4).
- **Required tools:** Second-engineer contact + fresh-substrate dry-run (R-4 already dry-run in T5.3; here extended to a fresh host).
- **Risks:** No named second engineer available → the run is preliminary-only regardless of outcome (the two-run rule, `LAB_README.md` §9). This is acceptable for a preliminary verdict but blocks Phase 12's journal entry.
- **Exit criteria:** Second engineer named and contactable (R-3); fresh-substrate reconstruction path dry-run on a clean host (R-4).

### Phase 8 — Pre-flight Go/No-Go

#### T8.1 — Execute the full checklist
- **Effort:** 1 eng-h
- **Depends on:** Phases 2–7 complete.
- **Required tools:** `EXPERIMENT_CHECKLIST.md`.
- **Risks:** An item "essentially" passed (the lab forbids adjective-graded verdicts even at the gate). Every P/E/R/V/H item either has a witness or fails.
- **Exit criteria:** P-1..P-6, E-1..E-11, R-1..R-5, V-1..V-5, H-1..H-5 each witnessed or flagged; any R-flag records the run as preliminary-only.

#### T8.2 — Produce and sign GO.json
- **Effort:** 0.5 eng-h
- **Depends on:** T8.1.
- **Required tools:** `GO.json` template (`EXPERIMENT_CHECKLIST.md` §7), Standards Editor signature.
- **Risks:** A Go-record with an unfilled or "N" field silently treated as a Go. The Go-record is the first artifact of the run; it is hashed and manifested (EVIDENCE_INDEX.md §3).
- **Exit criteria:** `GO.json` produced, every field filled, no "N", Standards Editor pre-flight initial present, hashed into the manifest.

### Phase 9 — Composition attempts (cheapest-first; the load-bearing phase)

**Cross-cutting discipline for T9.0–T9.3:**
- **Run-ID discipline.** Each composition attempt is its own `RUN-NN`. The cheapest-first ordering is binding. The spike **stops at the first composition that satisfies the C4 criterion** (outcome α); the other candidates are *not* run, and the count is reported honestly as "1 of up-to-3 attempted; the others not run because α was reached." Running further candidates after α just to inflate or verify the count is wasted experimental capital — the very thing the audit warned against.
- **Case set (instrumental, faithful to the spike question — Standards Editor confirms at T8.1):** each composition attempt runs {T1 (smoke), T4 (revocation observability within R), T5 (partition + sniffer-clean), T9 (control)}. This case set is the minimal realization of the spike question's clauses: T4 realizes "revocation observability within R seconds," T5 realizes "no live call... no third-party broker" under partition, T1 is the smoke prerequisite (you cannot observe revocation of a token you cannot first verify), T9 is the substrate control (S8 — if T9 fails, no claim is admissible). **The spike does NOT run T2/T3/T6/T7/T8** — those exercise already-solved components (per the feasibility audit's "Why not C" point 1) and running them is misallocated capital. Running them anyway is a boundary breach (A-4 / scope creep, U8).
- **Boundary breach trigger (A-4):** if the Engineer begins writing protocol code to make a composition pass (beyond a *thin, non-patching* revocation layer atop stock libraries), halt. The composition may add new code atop stock; it may not modify stock (S7/F7) and may not invent a new protocol (that is Level 2's territory, outcome β's "new-primitive approach was not tested by EXP-001").
- **Stop-condition A-8 (C4-specific):** if the RP reaches an address outside domain B during a partitioned case, halt — the substrate is contaminated.

#### T9.0 — Composition 0 (smoke before any candidate): T9 control + T1 smoke on the substrate
- **Effort:** 2 eng-h
- **Depends on:** Phase 8.
- **Required tools:** The substrate from Phases 4–6, the stock verifier (T3.2).
- **Risks:** T9 failing — the substrate is broken; **no claim is admissible** (S8/F8). T9 is the gate that says "the RP can verify a normal in-domain-B SVID at all." If T9 fails here, do not proceed to candidate compositions; fix the substrate and re-run T9 (this is `P5_FALSIFICATION_EXPERIMENT.md` F-J territory, methodologically). T1 failing means the basic JWT-SVID + delegation-token issuance/verification path is broken at the substrate level — also a fix-before-proceeding signal, not a C4 verdict.
- **Exit criteria:** T9 verifies-OK (control passes; S8 held); T1 verifies-OK within 100 ms with sniffer clean (smoke passed); both recorded as RUN-01 SETUP/OBSERVATION entries. If either fails, halt and fix the substrate — do not interpret as a C4 result.

#### T9.1 — Composition 1: OAuth Status List (RUN-02, if admissible)
- **Effort:** 4–6 eng-h
- **Depends on:** T9.0, **S2 scope-act admitting cached pulls (Phase 0)**, T3.4.
- **Required tools:** The Status List library (T3.4, if installed), the sniffer (T6.1), the latency probe (T6.2), the firewall (T5.1).
- **Risks:** (a) The Status List's refresh cadence cannot be set ≤ R without defeating the batching purpose (the audit's observability-latency falsification) — this surfaces as F-class failure, logged; the composition is recorded as failing the C4 criterion, not as a bug. (b) The list fetch *is* a live call to domain A (the audit's no-live-call falsification) — the sniffer catches it (F2/F4). (c) If a CDN/passive cache is used to evade the no-live-call clause, the S3 scope-act decides whether that counts as a broker — surface, do not resolve; log the ambiguity for the decision phase.
- **Exit criteria (run result, not "pass"):** T4 attempts a revocation set at T₀+5s observable at T₀+7s; T5 attempts partitioned verification. Record: verdict (OK/REJECT/REVOKED-observable?), latency, sniffer result (zero vs N outbound-outside-B packets), X/Y SVIDs still valid post-T4 (S3). Manifest + log entries sealed as RUN-02. **If the C4 criterion is met here (outcome α): stop, report 1 of up-to-3, go to Phase 11 (reproduce).**

#### T9.2 — Composition 2: Push-revocation (RUN-03)
- **Effort:** 6–8 eng-h
- **Depends on:** T9.1 (run regardless of T9.1 outcome, *unless* α was already reached at T9.1).
- **Required tools:** Stock JWT verifier + a thin inbound-listener layer (new code atop stock, not modifying stock), sniffer, latency probe, firewall.
- **Risks:** (a) The push itself cannot arrive during the T5 partition — under the S4 combined reading (revocation set *during* the partition), this composition fails by information-theoretic limit, which is **the spike surfacing outcome γ**, exactly its job. Under the S4 "eventual post-recovery" reading, the push arrives post-recovery and the composition may satisfy C4 modulo latency. **The S4 scope-act (Phase 0) is what makes this decidable** — without it the result is ambiguous and the run is void. (b) The inbound-listener layer being judged "not a thin revocation layer" but "a new protocol" → boundary breach (A-4); the Standards Editor confirms the listener is non-patching atop stock at audit.
- **Exit criteria (run result):** T4, T5, T9 executed, same record shape as T9.1; ambiguity about whether the listener is "thin" is logged verbatim for the decision phase. **If the C4 criterion is met (α): stop, report count, reproduce.** **If the failure is information-theoretic under S4-combined: that is γ-evidence; flag for Phase 12.**

#### T9.3 — Composition 3: Accumulator-based revocation (RUN-04)
- **Effort:** 8–12 eng-h *or excluded*
- **Depends on:** T9.2 (run only if neither T9.1 nor T9.2 reached α and the failure character is not already γ).
- **Required tools:** A dynamic-accumulator implementation (Camenisch-Lysyanskaya-style or equivalent). **The audit flagged this as "academic, no standardized OSS implementation interoperating with SPIFFE/JWT/SPIRE"** (`LEVEL0_1_FEASIBILITY_GATE.md` C4 analysis) — so the first sub-step is an OSS-existence check.
- **Risks:** (a) No standardized OSS impl exists → the composition is "not tested by EXP-001" (per outcome β's forbidden-conclusion clause: "only *these compositions* were tested"). This is itself evidence — record it; do not build a bespoke accumulator (that is new-protocol work, A-4 boundary breach). (b) An accumulator impl exists but uses heavyweight algebra whose verification latency approaches the 100 ms budget → F1; logged as a capability failure, not a bug. (c) Accumulator membership proof size blowing the measurement record uncomfortably — not a falsifier per se, recorded honestly.
- **Exit criteria (run result):** Either (i) the composition is excluded due to no standardized impl — logged as "not tested, no standardized OSS composition available" (this contributes to outcome β's evidence that the gap is a technology gap); or (ii) the composition is attempted and T4/T5/T9 recorded as in T9.1/T9.2. **If α: stop, report, reproduce.** **If all three attempted/excluded and none satisfied → outcome β or δ (decision phase).**

### Phase 10 — Evidence collection (continuous, woven through Phase 9)

This is not a separate scheduled phase; it is named here because its discipline is load-bearing and must not be improvised.

- **Per-test-case artifact bundle** (`P5_FALSIFICATION_EXPERIMENT.md` §7 + `EVIDENCE_INDEX.md` §3): token (redacted), verdict log line, latency (bench), sniffer pcap excerpt (T4/T5 only), wall-clock start/end, config snapshots, source-diff-against-stock (must be empty for S7), second-engineer commentary recorded within 24 h.
- **Hash-at-creation** (`EVIDENCE_INDEX.md` §5): every artifact hashed *before* it is used in any further step; manifest appended immediately (V-5).
- **No convenience summaries as evidence** ("T1–T8 all passed" is refused; raw artifacts are).
- **Count accounting (U5):** every `RUN-NN` is logged — passing, failing, and excluded-by-scope-act. The garden-of-forging-paths count is part of the evidence and is read out at Phase 12.

**Exit criteria:** Every Phase-9 run has a sealed manifest with every artifact hashed; the log has one OBSERVATION/MEASUREMENT entry per case; no convenience summary replaces raw artifacts.

### Phase 11 — Reproduction

#### T11.1 — Fresh-substrate reconstruction by the second engineer
- **Effort:** 5 eng-h
- **Depends on:** A composition that satisfied the C4 criterion (α) at Phase 9, **or** a conclusive negative (β) requiring reproduction per U1. *Preliminary-only run if R-3/R-4 flagged at Phase 8.*
- **Required tools:** The sealed manifest of the run being reproduced; a fresh host; the same pinned versions (T2.4).
- **Risks:** Latent state in the original substrate being load-bearing for the result (the point of fresh-substrate reproduction — `LAB_README.md` §9). The second engineer must *not* have access to the first run's artifacts beyond the config snapshot (per F-K in the Level 2 protocol).
- **Exit criteria:** RUN-(N+1) on a fresh substrate reconstructed from the manifest; same case set; result recorded as same / narrower / different (U7). **Different → contradicted; no conclusion admissible; pre-reg amendment or fresh experiment required** (this blocks Phase 12).

#### T11.2 — Compare reproduced verdict to original
- **Effort:** 2 eng-h
- **Depends on:** T11.1.
- **Required tools:** The two sealed manifests, the decision rules.
- **Risks:** A "narrower" reproduction being read as confirmation (U10 — Ockham's razor: divergent outcomes adopt the more conservative reading; do not pick the favorable one).
- **Exit criteria:** Reproduction verdict recorded in the log with its relationship to the original (same/narrower/different). The conclusive verdict is now *evidence-tier-2* (U7) and may enter a journal entry.

### Phase 12 — Decision + journal entry

#### T12.1 — Select the outcome class
- **Effort:** 1 eng-h
- **Depends on:** T11.2 (reproduction done) or a recorded preliminary verdict.
- **Required tools:** `DECISION_RULES.md` §3.
- **Risks:** Selecting the favorable outcome class (α) when the evidence supports a more conservative one (β) — U10 forbids this; the Standards Editor audits. Reading a single-run result as conclusive without reproduction — U7 / H-2 forbid this.
- **Exit criteria:** One of α / β / γ / δ selected, with the evidence list (criterion ID → observed fact → evidence ref) per `DECISION_RULES.md` §5.

#### T12.2 — Standards Editor audit
- **Effort:** 0.75 eng-h
- **Depends on:** T12.1.
- **Required tools:** The pre-reg (and its S1–S5 amendment), the decision rules, the sealed manifests.
- **Risks:** A threshold cited in the conclusion that was absent from the pre-reg + amendment (H-3 — refuse the conclusion). A boundary breach during the run not flagged (A-4 — refuse).
- **Exit criteria:** Pre-reg primacy held (Y); Standards Editor audit line filled; forbidden-readings checks (U8/U9/§4 patterns) explicitly considered and rejected.

#### T12.3 — Write the conclusion in the required structure
- **Effort:** 1 eng-h
- **Depends on:** T12.2.
- **Required tools:** The `DECISION_RULES.md` §5 template.
- **Risks:** A conclusion in any other shape is refused (`DECISION_RULES.md` §5).
- **Exit criteria:** Conclusion block produced verbatim in the §5 structure.

#### T12.4 — Record the journal entry
- **Effort:** 0.25 eng-h
- **Depends on:** T12.3.
- **Required tools:** `agents/journal/` per `agents/GOVERNANCE.md`.
- **Risks:** The lab unilaterally promoting the verdict to a project decision (retire Level 2, build P5, etc.) — `LAB_README.md` §13 forbids this; the journal entry *channels* the founder, it does not decide for them. The conclusion names C4 only (U8).
- **Exit criteria:** `agents/journal/<YYYY-MM-DD>-exp-001-c4-spike.md` exists, quoting the §5 conclusion verbatim and preserving any dissent. **The lab's work ends here; the founder's begins.**

---

## 5. Execution order (binding)

1. **Phase 0** (scope-act gate) — founder; blocks all.
2. **Phases 1–8** in dependency order (§3 graph) — single Engineer; Phase 7 parallelizable.
3. **Phase 9** cheapest-first: T9.0 (smoke+control) → T9.1 (Status List) → T9.2 (push) → T9.3 (accumulator), **stopping at the first composition that satisfies the C4 criterion** and reporting the count honestly.
4. **Phase 10** continuous through Phase 9.
5. **Phase 11** reproduction of whichever run reached α, or of a conclusive β.
6. **Phase 12** outcome class → audit → conclusion → journal entry → *stop*.

**No phase may begin out of order.** An out-of-order phase is a boundary breach logged as a `STOP` entry (A-class trigger where applicable) and is itself an item the Standards Editor audits at T12.2.

---

## 6. Evidence → criterion map (verification plan)

This is the verification spine: which artifact witnesses which S/F criterion, mechanized so the conclusion does not depend on the Engineer's reading.

| Criterion (frozen Level 2) | Witnessed by | Mechanization | Spike phase |
|---|---|---|---|
| S1 (≤100 ms verification) | `bench/<T-NN>-latency.json` | latency probe (T6.2); auto-fail F1 if > 100 ms | T9.0–T9.3 |
| S2 / S4 (sniffer clean; reject tamper/replay/revoked locally) | `pcap/<T-NN>-<descriptor>.pcap` + verdict log | sniffer (T6.1); any out-of-domain-B packet on a verified-OK case = F2/F4 | T9.1–T9.3 (T4/T5) |
| S3 (independent revocation; X/Y SVIDs still valid post-T4) | `token/<T-NN>-*-token.json` (redacted) + a post-T4 T9 re-run | T9 must pass after T4; F3 if it doesn't | T9.1–T9.3 |
| S5 (trust-bundle lag; not load-bearing for C4 directly, not run) | — | not exercised by the spike; deferred to Level 2 if β | — |
| S7 (zero source mods to stock) | `analysis/<component>-diff-upstream.txt` | `diff` against upstream tag (E-10); empty required; F7 if non-empty | T3.1–T3.6, audited T12.2 |
| S8 (T9 control passes) | `log/RP-T09.log` + verdict | T9 verifies-OK; F8 (→ methodological fix, not a C4 verdict) if not | T9.0 |
| F9 (reproducible) | RUN-(N+1) manifest vs RUN-NN manifest | fresh-substrate reproduction (T11.1) | Phase 11 |
| F10 (no post-hoc tuning) | manifest + log | every threshold cited is in the pre-reg + S1–S5 amendment; audited at T12.2 (H-3) | Phase 12 |

**Outcome class selection (post-run, `DECISION_RULES.md` §3):**
- **α** — a composition met T4's R-bound + T5's sniffer-clean + S3 + T9 passed, reproduced → C4 solvable by named primitives under the stated S1–S5 reading.
- **β** — all admissible compositions attempted/excluded, none satisfied, failure is a capability gap (not information-theoretic), reproduced → technology gap; Level 2 justified and narrows to C4 + C11 per-hop.
- **γ** — a formal argument (run-independent, recorded under `analysis/`) shows T4∩T5 under the S4-combined reading is information-theoretically unsatisfiable → logical impossibility; scope act required, no further experiments against un-corrected H1.
- **δ** — compositional attempts inconclusive AND S1–S5 unresolvable without empirical input (the empirically-unresolvable subset, e.g., R cannot be set a priori) → Level 2 justified as-is.

---

## 7. Risk register (consolidated, ranked)

| ID | Risk | Likelihood | Impact | Mitigation | Owning phase |
|---|---|---|---|---|---|
| R-1 | Phase 0 scope-act not closed; spike runs against undefined R | Med | **Run void** (H-3) | Hard block in §1; no Phase-1 start until amendment exists | Phase 0 |
| R-2 | Stock verifier silently calls a key-resolution endpoint (JWKSet fetch) | Med | **False α** (F2/F4 sniffer hit) | T3.2 explicit bundle-only check; sniffer self-test (T6.1) | Phase 3/6 |
| R-3 | Federation left enabled; composition "succeeds" via smuggled broker | Low-Med | **False α** (invalidates the no-broker clause) | T4.3 explicit disable + E-8 cite | Phase 4 |
| R-4 | Engineer writes protocol code to make a case pass (boundary breach A-4) | Med | **Run aborted**; or silent if undetected | T9 cross-cutting discipline; Standards Editor diff audit at T12.2 | Phase 9 |
| R-5 | Single-run result treated as conclusive (no reproduction) | Med | **H-2 halt**; conclusion refused | U7 two-run rule; Phase 11 mandatory for α/β | Phase 11/12 |
| R-6 | Garden-of-forging-paths: only the passing run reported (U5 violation) | Med | **Conclusion refused** | One RUN-NN per attempt, including excluded; count read at T12.1 | Phase 9/10 |
| R-7 | S4-ambiguity surfacing as a "technology gap" when it is logical impossibility (γ misread as β) | Med | **Wrong outcome class** | S4 scope-act (Phase 0) decides this *before* the run; T12.1 checks the reading | Phase 0/12 |
| R-8 | Floating version tag upgrades a stock lib mid-spike (S7 violation retroactively) | Low | **F7** | T2.4 version manifest; diff-verification at T3.x | Phase 2/3 |
| R-9 | Sniffer misses the very out-of-domain-B packets it must catch | Low | **False PASS** (the falsifier fails to falsify) | T6.1 self-test on known-good packets | Phase 6 |
| R-10 | T9 (control) fails and is misread as a C4 result instead of a substrate fault | Low-Med | **Wrong verdict type** (F8/F-J is methodological, not H1) | T9.0 explicit gate: fix substrate before candidates | Phase 9 |
| R-11 | Second engineer unavailable → preliminary-only verdict cannot reach journal | Med | **Blocked journal entry** (acceptable, not a fault) | T7.3 confirm early; if unavailable, the run is honestly preliminary | Phase 7/11 |
| R-12 | Accumulator composition drifts into bespoke protocol design (A-4) | Med | **Run aborted** | T9.3: if no standardized impl, exclude and log — do not build | Phase 9 |

---

## 8. Stop conditions (pointer, plus spike-specific)

All during-run aborts (A-1..A-8) and post-run halts (H-1..H-4) of `EXPERIMENT_CHECKLIST.md` §6 apply unchanged. The spike-specific emphases:

- **A-3 (parameter not in pre-reg)** — tripped if the Engineer reaches for an R or a TTL not fixed in the Phase-0 amendment. Halt; the run is designing itself.
- **A-4 (boundary breach)** — tripped if the Engineer writes protocol/architecture/product code to make a case pass. The thin revocation layer atop stock is the allowed envelope; beyond it, halt.
- **A-8 (C4)** — RP reaches an address outside domain B during a partitioned case. Halt; substrate contaminated.
- **H-3 (post-run)** — a threshold cited in the conclusion was absent from the pre-reg + amendment. Refuse the conclusion. This is the S1–S5 gate's enforcement arm.

A halted run is recorded with the trigger ID and **no conclusion** (`LAB_README.md` §5, `EXPERIMENT_CHECKLIST.md` §6).

---

## 9. Exit criteria for the plan as a whole (when is the spike "done")

The spike is **done** when exactly one of the following holds:

1. **Outcome α reached, reproduced, audited, journaled.** A composition satisfied the C4 criterion under the S1–S5 reading; a second engineer reproduced it on a fresh substrate; the Standards Editor confirmed pre-reg primacy; the journal entry at `agents/journal/<date>-exp-001-c4-spike.md` quotes the §5 conclusion. *The lab stops; the founder decides whether to retire Level 2.*
2. **Outcome β reached, reproduced, audited, journaled.** All admissible compositions attempted/excluded; the gap is a technology gap, not logical impossibility; reproduced; audited; journaled. *The lab stops; the founder decides whether to narrow Level 2 to C4 + C11 per-hop.*
3. **Outcome γ reached, audited, journaled.** A formal information-theoretic argument (under the S4-combined reading) is recorded under `analysis/` and Standards-Editor-reviewed; audited; journaled. *The lab stops; the founder resolves the S4 ambiguity by scope act before any further experiment.*
4. **Outcome δ reached, audited, journaled.** Compositional attempts inconclusive and the empirically-unresolvable subset of S1–S5 is demonstrated unresolvable without empirical input; audited; journaled. *The lab stops; Level 2 is justified as-is and the spike folds the empirically-unresolvable subset into Level 2's instrumented run.*
5. **The run is aborted** (a stop condition tripped) with the trigger ID and no conclusion. *The lab stops; the next action is a fix or a pre-reg amendment, recorded as a new LAB-PROCESS or CORRECTION log entry.*

In all five cases, **the lab does not, on the strength of the result, build, fund, retire Level 2 unilaterally, or open Product Management** (`LAB_README.md` §13). The journal entry is the channel; the founder decides on it.

---

## 10. Provenance, confidence, and stop rule

**Provenance.** Pre-registration: `LEVEL0_1_FEASIBILITY_GATE.md` §"Spike question (C4-only)" (frozen, Verdict B). Test-case + criteria definitions: `P5_FALSIFICATION_EXPERIMENT.md` (frozen Level 2 protocol, referenced by the spike, not re-scoped). Lab process: the five documents in `lab/`. No external sources past the feasibility gate's already-fetched canonical URLs are required for *planning*; composition-specific library versions are pinned at Phase 2 (T2.4), not pre-chosen here (doing so would be an implementation decision the plan forbids).

**Confidence (per project doctrine: evidence + change-condition required):**
- **Plan completeness (every Phase-9 task, every criterion, every exit criterion mapped): High.** Evidence: the plan traces each S/F criterion (§6) to a witnessing artifact and a mechanized check; each phase has exit criteria mapping to checklist items; the S1–S5 prerequisite is named as a hard block. Change-condition: if `EXPERIMENT_CHECKLIST.md` gains a new item, the affected phase's exit criteria update.
- **Plan executability (the spike *can* be run as planned on a 3-host substrate): Medium-High.** Evidence: the substrate, instruments, and case set are each individually realizable per the frozen protocol; the open executability risk is composition 3's accumulation step (no standardized impl — handled by exclude-and-log, T9.3). Change-condition: if a candidate composition proves unconstructable atop stock without modification, that composition is excluded (evidence for β), not forced.
- **Effort estimate (54.5–62.5 eng-h single-engineer +reproduction): Medium.** Evidence: bottom-up per-task estimates against the substrate the protocol specifies; the band reflects composition-3 uncertainty. Change-condition: if the substrate is reused from prior work (not the case here — fresh lab), Phase 2–4 effort drops; if composition 3 must be attempted against an emerging impl, the upper bound rises.
- **Predictive confidence about the spike outcome: None.** Per project doctrine (confidence without evidence is forbidden and refused-by-default for forecasts). This plan asserts nothing about which of α/β/γ/δ will occur.

**Stop rule.** This document is the execution plan. It does not write code, design protocol, choose product, or start the run. It does not execute Phase 0; the founder does. The next action belongs to the founder: (a) close the Phase-0 S1–S5 scope-act as a journal entry, (b) approve this plan as written, (c) approve with revisions, or (d) redirect. No engineering begins until (a) and (b)/(c) are complete; the lab `EXPERIMENT_LOG.md` records the approval as the next LAB-PROCESS entry.
