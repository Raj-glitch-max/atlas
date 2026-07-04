# Lab README — Engineering Research Phase

This is the **root document** of the engineering laboratory. It governs how experiments are run, observed, recorded, and decided within this repository. It does not govern what is built, what is shipped, or what is funded — those belong upstream, in `FOUNDER_DECISION_BRIEF.md`, `RESEARCH_PROGRAM.md`, and the project governance in `agents/GOVERNANCE.md`.

The lab exists to convert uncertainty into evidence. It is not a design space, not a product group, and not a decision body. Its output is recorded observation governed by pre-registered decision rules; the founder decides what the observation *means for the project*.

---

## 1. Purpose and non-purpose

**Purpose.** Produce engineering evidence whose validity does not depend on the engineer's enthusiasm for the hypothesis. The lab's entire apparatus — pre-registration, append-only logs, hash-pinned evidence, second-engineer reproduction, frozen decision rules — exists to make a single property true: *a conclusion the lab certifies is one a skeptic would also certify from the artifacts alone.*

**Non-purpose (binding).** The lab does not:
- implement, prototype, or write product code;
- design protocols, systems, or product architecture;
- choose technologies for a product;
- discuss market, business, pricing, go-to-market, or features;
- recommend whether to build, fund, or ship anything;
- retire an upstream experiment (e.g., the frozen Level 2 protocol) or open a downstream phase (e.g., Product Management).

The lab runs pre-registered experiments and records what happened. The scope, the hypothesis, the success/failure criteria, and the disposition of the result are all defined upstream and frozen before any run begins. The lab's job is to execute them faithfully and report.

---

## 2. Scope constraints (inherited, restated for the lab)

From the active Engineering Research instruction set:

- No implementation.
- No protocol design.
- No architecture.
- No product discussion.
- No coding.
- No more strategy documents (planning is closed).

These apply to the lab's *operations* as much as to its *outputs*. An experiment run that begins to design a protocol mid-run has crossed the boundary; the correct action is to stop (see `EXPERIMENT_CHECKLIST.md` §6) and record the boundary breach in the log, not to absorb the design work into the experiment.

---

## 3. Document set

The lab consists of five durable process documents plus the evidence they govern. All live in `/home/raj/Videos/projects/lab/`.

| Document | Role |
|---|---|
| `LAB_README.md` (this file) | Root. Workflow, lifecycle, naming, reproducibility discipline. |
| `EXPERIMENT_LOG.md` | Append-only chronological record of observations. No opinions. |
| `EXPERIMENT_CHECKLIST.md` | Go/no-go pre-flight gate plus run-time stop conditions. |
| `EVIDENCE_INDEX.md` | Registry of every artifact type, its storage path, integrity, and retention. |
| `DECISION_RULES.md` | For each possible outcome: evidence observed → conclusion allowed → conclusion forbidden. Anti-bias core. |

The frozen experiment specifications the lab executes live *upstream*, not in `lab/`:
- `P5_FALSIFICATION_EXPERIMENT.md` — the frozen Level 2 protocol (do not modify).
- `LEVEL0_1_FEASIBILITY_GATE.md` — the approved feasibility gate; its "Spike question (C4-only)" section is the pre-registration for the first experiment (EXP-001).

---

## 4. Roles

The lab recognizes three roles. A single human may hold more than one, except where role separation is explicitly required.

| Role | Responsibility | Separation requirement |
|---|---|---|
| **Engineer** | Sets up substrate, runs the experiment, captures evidence, logs observations. | Must not be the sole reviewer of their own conclusion for any publishable verdict. |
| **Adversary** | Runs adversarial and partition test cases; attempts to break the setup; adversaries pre-registered thresholds. | **Must not see the success/failure threshold values before running** (adversary-blind). See `EXPERIMENT_CHECKLIST.md` §5. |
| **Standards Editor** | Audits that pre-registration primacy, anti-bias rules, and conclusion discipline held; signs off the post-run conclusion. | Must not have run the experiment being audited. |

The project's council review order (Empiricist → Cartographer → Red Team → Economist → Operator) and `agents/GOVERNANCE.md` continue to bind at any point a lab conclusion is elevated to a project-level decision. The lab does not replace governance; it produces the inputs governance reviews.

---

## 5. Experiment lifecycle

Every experiment in this repository follows the same sequence. No stage may be skipped; out-of-order stages are recorded as boundary breaches.

1. **Pre-registration.** The experiment specification is frozen and hash-pinned. Success criteria, failure criteria, thresholds, and the applicable outcome table (in `DECISION_RULES.md`) are fixed before any run. Pre-registration is upstream of the lab; the lab does not write it.
2. **Pre-flight.** `EXPERIMENT_CHECKLIST.md` is executed in full. Every item's witness is recorded. The Go/No-Go gate (§7 of the checklist) decides whether the run begins.
3. **Run.** The experiment is executed. Observations are logged in `EXPERIMENT_LOG.md` as they occur, not reconstructed afterward. Evidence is captured to the paths defined in `EVIDENCE_INDEX.md`. The adversary runs adversary-blinded cases. Stop conditions are monitored continuously (§6 of the checklist).
4. **Capture close.** Per-run manifest is finalized, every artifact hashed, every artifact's witness recorded. The run is sealed: no new artifacts may be added to a sealed run without a new run ID.
5. **Decision.** The applicable outcome class in `DECISION_RULES.md` is selected from the evidence. The Standards Editor audits that pre-registration primacy held. The conclusion is written with the structure required by `DECISION_RULES.md` §5.
6. **Reproduction (for any conclusive verdict).** A second engineer re-runs from the manifest on a fresh substrate. A conclusive verdict from a single unreproduced run is *preliminary*; a conclusive verdict from a reproduced run is *evidence-tier-2*. See §9.
7. **Journal entry.** A conclusive verdict, once reproduced and audited, is recorded in `agents/journal/<YYYY-MM-DD>-<slug>.md` per project governance. The lab does not unilaterally promote a verdict to a project decision; the journal entry is the channel by which the founder receives it.

An experiment is **closed** when its journal entry exists and its outcome class is recorded. An experiment is **aborted** if a stop condition tripped before decision; an aborted run is recorded with the stop-condition ID that triggered it and no conclusion.

---

## 6. Artifact lifecycle

```
created  →  hashed  →  manifested  →  sealed  →  retained  →  archived or destroyed
```

- **Created.** At the moment the instrument produces it. Logged in `EXPERIMENT_LOG.md` with a reference to its storage path.
- **Hashed.** Every artifact gets a content hash (SHA-256) recorded in the run manifest at creation time. The hash is the artifact's identity; the path is its location.
- **Manifested.** Added to the run's `MANIFEST.json` (see `EVIDENCE_INDEX.md` §4). A sealed manifest is append-only for corrections, never edited in place.
- **Sealed.** At run close, the manifest is itself hashed and the seal recorded in the log. Sealed artifacts are immutable.
- **Retained.** For the duration specified in `EVIDENCE_INDEX.md` §6. No artifact is deleted before its retention window closes.
- **Archived or destroyed.** Per retention rules. Destruction is logged with a reason and the prior hash, so a destroyed artifact remains verifiable-as-having-existed even after its bytes are gone.

---

## 7. Naming convention

All identifiers are deterministic, sortable, and human-readable. No dates in filenames are inserted by hand — they come from the run manifest's UTC timestamp.

| Entity | Pattern | Example |
|---|---|---|
| Experiment | `EXP-<NNN>` | `EXP-001` |
| Run | `EXP-<NNN>/RUN-<NN>` | `EXP-001/RUN-01` |
| Test case | `T<NN>` (from pre-registration) | `T04` |
| Artifact | `EXP-<NNN>/RUN-<NN>/<type>/<slug>.<ext>` | `EXP-001/RUN-01/pcap/T04-replay-post-ttl.pcap` |
| Manifest | `EXP-<NNN>/RUN-<NN>/MANIFEST.json` | `EXP-001/RUN-01/MANIFEST.json` |
| Log entry | `E-<NNNN>` (sequential, never reused) | `E-0007` |

The first experiment in this lab is **EXP-001**, the C4 Feasibility Spike, pre-registered by `LEVEL0_1_FEASIBILITY_GATE.md` §"Spike question (C4-only)."

---

## 8. Evidence handling and integrity

- **Storage root.** `lab/evidence/`. All artifacts live under `lab/evidence/<experiment-id>/<run-id>/<type>/`.
- **No artifact in context.** Raw artifacts (pcaps, logs, traces) are stored on disk and *referenced* from logs and manifests, never pasted inline. A log entry points to a path plus a hash; it does not contain the bytes.
- **Hash primacy.** An artifact's SHA-256 hash is its identity. If a path and a hash disagree, the hash is authoritative and the path is recorded as corrupted. If two artifacts share a hash, they are the same artifact regardless of filename.
- **No silent edits.** A sealed artifact that needs correction is *not* overwritten; a corrected copy is produced with a new filename and a new hash, and the manifest notes the correction. The original is retained.
- **Witness requirement.** Every check, every artifact, every measurement has a witness — a file, a command output, or a signature. Observations without witnesses are not evidence.

---

## 9. Reproducibility discipline

A run that another engineer cannot reproduce from the manifest is a run whose conclusion cannot exceed *preliminary*. Reproducibility is not a courtesy; it is the threshold between a claim and a fact.

- **Config snapshot.** Every run records a snapshot of substrate configuration (versions, configs, keys in use, clock state, isolation rules) sufficient for a second engineer to reconstruct it. The snapshot is itself an artifact and is hashed.
- **Fresh substrate.** Reproduction uses a *fresh* substrate (new instances, new keys where applicable), not the original, so that latent state in the original is not load-bearing for the result.
- **Two-run rule.** A conclusive verdict (per `DECISION_RULES.md`) requires a reproduced run by a second engineer who was not the original Engineer. Preliminary verdicts may be recorded from a single run but must be labeled as such; they may not enter a journal entry.
- **No re-tuning.** If the first run fails and the engineer changes a parameter and re-runs, *both* runs are logged in full; the second run is a new run with a new ID, not a correction of the first. Reporting only the passing run is forbidden (see `DECISION_RULES.md` §2 "Garden-of-forging-paths guard").

---

## 10. Pre-registration primacy

The single most load-bearing discipline in the lab: **only criteria frozen before the run counts as evidence; thresholds chosen after seeing the result are not evidence.**

- Success criteria, failure criteria, thresholds, and the outcome table that maps evidence to conclusions are all fixed in pre-registration.
- A run that discovers a *better* threshold and applies it is, by definition, not testing the pre-registered hypothesis; that run may be informative for a *future* pre-registration but is inadmissible as evidence for the current conclusion.
- The Standards Editor's first audit question is always: *were the success and failure thresholds cited in the conclusion present in the pre-registration before the run?* If not, the conclusion is refused.

---

## 11. Stop conditions (pointer)

Stop conditions — both during-run aborts and post-run halts — are defined in `EXPERIMENT_CHECKLIST.md` §6. They are not guidelines; triggering a stop condition halts the run and the run is recorded as aborted with the trigger ID, no conclusion.

---

## 12. Change control for these documents

The five lab documents are durable but not immutable. A change to any of them:

1. Must be motivated by a concrete failure of the existing document observed during a real run, not by preference.
2. Must be recorded in `EXPERIMENT_LOG.md` as a lab-process entry (separate class from experimental entries).
3. Must not retroactively alter a pre-registration or a sealed run. A change to `DECISION_RULES.md` does not re-open a sealed run's conclusion; it applies to runs whose pre-registration post-dates the change.
4. Is itself versioned: the affected document gains a version line at the top citing the log entry that motivated it.

Lab-document changes do not require a journal entry (they are process, not project decisions) unless the change alters the boundary between the lab and upstream governance, in which case a journal entry is required.

---

## 13. Relationship to project governance

- The lab operates *inside* `agents/GOVERNANCE.md`. Confidence labels (High / Medium / Low / None), the requirement for evidence and change-conditions, the prohibition on consensus smoothing, and the preserve-dissent rule all apply to lab outputs as much as to any other project artifact.
- A lab conclusion that survives reproduction and Standards-Editor audit is recorded as a journal entry. The journal entry is the lab's only channel to the rest of the project.
- The lab does not, and cannot, open Product Management, retire the frozen Level 2 protocol, or change the candidate set. Those are founder decisions taken on the journal entry, not lab actions.
- If the lab produces evidence that calls an upstream document into question (e.g., the frozen Level 2 protocol contains a logical tension the lab surfaced), the lab records the finding in the log, raises it via a journal entry, and stops. It does not edit the upstream document.

---

## 14. First experiment — EXP-001 (C4 Feasibility Spike)

The lab's first experiment is the C4 Feasibility Spike, authorized by `LEVEL0_1_FEASIBILITY_GATE.md` Verdict B and pre-registered by its "Spike question (C4-only)" section.

- **Pre-registration:** `LEVEL0_1_FEASIBILITY_GATE.md` §"Spike question (C4-only)" plus the candidate compositions enumerated there.
- **Outcome table applicable:** `DECISION_RULES.md` §3 (outcomes α / β / γ / δ).
- **Scope:** the spike tests the C4 component claim only (offline, fast, partition-tolerant, independent revocation via existing-standard composition). Other H1 components (C1, C2, C3, etc.) are not re-tested by this spike.
- **Stakeholders of the result:** per the feasibility gate, the spike outcome disposes of the frozen Level 2 protocol (narrow, justify-as-is, or correct) — but the spike does not make that disposition; the founder does, on the journal entry the spike produces.

The lab is now ready to pre-flight EXP-001/RUN-01. The next action is the pre-flight checklist, executed against the spike's pre-registration.
