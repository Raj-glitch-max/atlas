# Experiment Checklist

**Type.** Go/no-go gate executed before every experiment run, plus the stop conditions monitored during the run. This document is not advisory; a run that has not passed every applicable item below has not started. Items marked **[G]** are universal; items marked **[C4]** are specific to EXP-001 and are marked as such so future experiments can append their own.

Each item has a **witness** — the evidence that the item passed. "Witnessed" means a path, a hash, a command output, or a signature is recorded in the run manifest (`EVIDENCE_INDEX.md` §4). An item with no witness is an item that did not pass.

---

## 1. Pre-registration gate  [G]

| ID | Item | Witness |
|---|---|---|
| P-1 | The experiment's pre-registration document exists and is frozen (version-stamped, not mid-edit). | Path to pre-reg doc + version line |
| P-2 | Pre-reg's SHA-256 hash recorded before any run-specific work begins. | `path#sha256:…` in manifest |
| P-3 | Success criteria, failure criteria, and thresholds are present in the pre-reg and named by ID (e.g., S1, F4). | Pre-reg section reference list |
| P-4 | The applicable outcome table in `DECISION_RULES.md` is identified and hash-pinned to this run. | `DECISION_RULES.md` version + section |
| P-5 | No threshold in the pre-reg is a placeholder or "TBD." Thresholds marked post-run are inadmissible. | Standards-Editor sighted-initials on pre-reg |
| P-6 | The hypothesis tested is copied verbatim into the run's first log entry, not paraphrased. | `EXPERIMENT_LOG.md` E-NNNN reference |

A failure on any P-item is a no-go. The pre-registration is incomplete; do not start.

---

## 2. Environment verification  [G] / [C4]

| ID | Item | Witness |
|---|---|---|
| E-1 [G] | Substrate configuration matches the pre-reg's stated environment section. | Config snapshot artifact path |
| E-2 [G] | All software versions pinned and recorded; no `latest`, no floating tags. | Version manifest artifact path |
| E-3 [G] | System clocks NTP-disciplined; inter-host clock skew measured and recorded. | Skew measurement artifact (e.g., chrony tracking output) |
| E-4 [G] | Hosts reachable per the pre-reg's network topology; isolation rules inactive *at setup time*. | Network reachability probe output |
| E-5 [G] | Free disk capacity at the evidence storage root exceeds the pre-reg's stated maximum capture size with 2× margin. | `df -h lab/evidence/` output |
| E-6 [G] | Out-of-band capture instruments (sniffer, latency probe) are armed and tested on a known-good packet before substrate-specific cases run. | Sniffer self-test artifact |
| E-7 [C4] | Two independent SPIRE deployments present, distinct servers, distinct databases, distinct signing keys. | Two server config snapshots; key fingerprints differ |
| E-8 [C4] | SPIFFE federation disabled in both deployments at the configuration level. | Disabled-federation config lines cited |
| E-9 [C4] | Relying party host's outbound firewall ruleset defined and *inactive* at setup (rules activated per-test only). | Ruleset artifact |
| E-10 [C4] | Stock JWT/JWS / COSE verifier versions pinned; no source modifications to verifier, SPIRE, SPIFFE libs, or RFC 8693 client. | `diff` against upstream tags returned empty (artifact path) |
| E-11 [C4] | TTL configuration locked: delegation TTL = 30 s, Y-SVID TTL = 600 s, per pre-reg. | Config snapshot embedding these values |

A failure on any E-item is a no-go unless the pre-reg explicitly waives that item (and the waiver is recorded in the manifest).

---

## 3. Reproducibility checks  [G]

| ID | Item | Witness |
|---|---|---|
| R-1 | Config snapshot complete enough to reconstruct the substrate on fresh hosts. | Snapshot artifact + Standards-Editor "reconstructable" sighted-initials |
| R-2 | Any randomness seeds or randomized parameters recorded (not left in `Math.random()`). | Seed/parameter log |
| R-3 | A second engineer is available for the reproduction run (named, contactable, adversary-blind where applicable). | Engineer identifier + contact |
| R-4 | Fresh-substrate capability confirmed: a clean reconstruction path exists and was dry-run at setup. | Dry-run log artifact |
| R-5 | Manifest template initialized for this run before the first measurement. | Manifest artifact timestamp precedes first measurement |

A failure on R-1, R-3, or R-4 means the run may proceed *only* as a preliminary run; its result cannot exceed "preliminary" regardless of outcome (see `LAB_README.md` §9 two-run rule).

---

## 4. Evidence readiness  [G]

| ID | Item | Witness |
|---|---|---|
| V-1 | Evidence directory exists for this run: `lab/evidence/EXP-NNN/RUN-NN/`. | `ls` output |
| V-2 | Subdirectory structure for every artifact type the pre-reg may produce exists. | `ls` output |
| V-3 | Capture instruments write directly to the run directory, not to `/tmp`. | Capture path config |
| V-4 | Hashing procedure (see `EVIDENCE_INDEX.md` §5) tested on a sample artifact. | Sample hash artifact |
| V-5 | Manifest writer ready; the first artifact produced will be appended to the manifest immediately. | Manifest-writer check artifact |

A failure on any V-item is a no-go. Evidence left in `/tmp` is evidence at risk of being lost; do not start.

---

## 5. Personnel and role separation  [G]

| ID | Item | Witness |
|---|---|---|
| H-1 | Engineer role assigned; not also the sole reviewer of this run's conclusion. | Role log |
| H-2 | Adversary role assigned for any adversary-blinded test cases. | Role log |
| H-3 | Adversary has **not** seen the success/failure threshold values in `DECISION_RULES.md` §3 or the pre-reg's S/F criteria. | Adversary sighted-initials confirming blindness |
| H-4 | Standards Editor identified for post-run audit; did not run this experiment. | Role log |
| H-5 | For any test case where blindness is required, the adversary's threshold-blindness is re-confirmed immediately before that case runs. | Per-case blindness re-confirm entry in log |

A failure on H-3 or H-5 invalidates the adversary-blinded case; the case must be re-run with a freshly blinded adversary or recorded as not-adversary-blinded (which lowers its evidentiary tier).

---

## 6. Stop conditions

Stop conditions halt the run. A halt is not a failure of the experiment; it is the lab refusing to produce invalid evidence. A halted run is recorded with the triggering condition's ID and **no conclusion**.

### 6.1 During-run abort triggers (halt immediately)

| ID | Trigger | Action |
|---|---|---|
| A-1 | A substrate property drifts outside pre-reg tolerances during the run (e.g., clock skew exceeds the recorded bound). | Halt; log the drift measurement and the bound it crossed. |
| A-2 | An out-of-band instrument (sniffer, latency probe) fails or returns self-inconsistent results. | Halt; the instrument's validity is the run's validity. |
| A-3 | A test case requests a parameter or threshold not defined in the pre-reg. | Halt; the experiment is attempting to design itself — a scope breach. |
| A-4 | The engineer begins writing product, protocol, or architecture code to make a case pass. | Halt; boundary breach per `LAB_README.md` §2. Record the breach. |
| A-5 | A test case behaves in a way inconsistent with the pre-reg's stated threat model. | Halt; the run is testing something other than the pre-registered hypothesis. |
| A-6 | Evidence capture to the manifest fails (write error, hash mismatch on a freshly written artifact). | Halt; integrity of the evidence stream cannot be guaranteed. |
| A-7 | Hardware MCE, OOM, or host fault affecting any substrate host. | Halt; record the fault. |
| A-8 [C4] | The relying party reaches an address outside domain B during a case that pre-reg requires to be partitioned/isolated. | Halt; this is also a falsifier trigger (F4-class), but the halt prevents the run from continuing on a contaminated substrate. |

### 6.2 Post-run halts (do not begin the next phase)

| ID | Trigger | Action |
|---|---|---|
| H-1 | The manifest is incomplete: an artifact referenced in the log is missing from the manifest, or vice versa. | Do not proceed to Decision; reconcile or re-run. |
| H-2 | The run has not been reproduced and the Engineer wishes to record a conclusive verdict. | Refuse; conclusive verdicts require reproduction (`LAB_README.md` §9). |
| H-3 | The Standards Editor identifies that a threshold cited in a draft conclusion was absent from the pre-reg. | Refuse the conclusion; the run is recorded as inconclusive-on-pre-reg-primacy grounds. |
| H-4 | Two test cases produce contradictory evidence on the same pre-registered criterion. | Do not collapse the contradiction; record both, halt for a re-run or a pre-reg amendment. |

---

## 7. Go / No-Go

The run begins only when a single Go-record is produced:

```
GO RECORD — EXP-NNN/RUN-NN
UTC: <timestamp>
Pre-reg: <path>#sha256:<hash>  (P-1..P-6 passed: Y/Y/Y/Y/Y/Y)
Environment: <snapshot path>#sha256:<hash>  (E-1..E-11 passed: <list>)
Reproducibility: R-1..R-5 passed: <list>  (preliminary-only if any R flagged)
Evidence: V-1..V-5 passed: <list>
Personnel: H-1..H-5 passed: <list>
Stop-conditions armed: A-1..A-8 monitored by <role>
Standards Editor pre-flight sign: <initials>
```

A Go-record with any "N" or unfilled field is not a Go; it is a No-Go and the run does not start.

The Go-record is itself an artifact, hashed and added to the manifest. It is the first artifact of the run.
