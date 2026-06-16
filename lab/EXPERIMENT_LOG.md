# Experiment Log

**Type.** Append-only chronological record of observations during experiment runs.

**Rule (binding).** Observations only. No opinions, no interpretation, no conclusions. Interpretation belongs to `DECISION_RULES.md`, not here. An entry that begins "I think the result shows…" is malformed; rewrite it as "T04 returned REJECT-REPLAY at +7.2 s; sniffer recorded 0 outbound packets" and let the decision rules extract the meaning.

A conclusion reached in the log is void; the log records what was observed, the decision rules determine what it means.

---

## Entry schema

Every entry has the following fields. Fields are non-empty except where explicitly marked optional.

| Field | Meaning | Constraint |
|---|---|---|
| `Entry ID` | `E-NNNN`, sequential, never reused | Monotonic; no gaps without explanation |
| `UTC timestamp` | ISO-8601 UTC, second precision | From system clock at observation time; never hand-set |
| `Experiment ID` | `EXP-NNN` | Must reference a pre-registered experiment |
| `Run ID` | `EXP-NNN/RUN-NN` | Must reference a run that passed pre-flight (or be a lab-process entry; see below) |
| `Entry class` | one of: `SETUP`, `RUN`, `OBSERVATION`, `MEASUREMENT`, `UNEXPECTED`, `STOP`, `SEAL`, `CORRECTION`, `LAB-PROCESS` | Single value |
| `Hypothesis tested` | verbatim from pre-registration | Copied, not paraphrased; for `LAB-PROCESS` entries use "N/A — lab process" |
| `Observations` | factual description of what happened | No adjectives of degree ("good", "essentially", "almost"); no causal claims beyond the directly observed |
| `Measurements` | numeric values with units | Field optional only if the entry class has no measurements; numbers always with units and with the instrument that produced them |
| `Evidence refs` | artifact path + SHA-256 | For every artifact referenced; format: `path#sha256:hexdigest` |
| `Unexpected findings` | deviations from pre-registered predictions | "None" is a valid value; absence of this field is not |
| `Result status` | one of: `RUN` (in progress), `RAN-OK`, `RAN-FAIL`, `INCONCLUSIVE`, `ABORTED`, `N/A` | `ABORTED` requires a stop-condition ID in Observations |
| `Next action` | the next planned step, stated as an action | An intention, not an opinion; "Pre-flight RUN-02 on fresh substrate" not "investigate why T04 misbehaved" |
| `Roles present` | Engineer / Adversary / Standards-Editor initials | At least one; adversary-blind runs note "adversary: pre-flight-blinded" |

**Anti-patterns forbidden in the log:**
- Conclusions inline ("this confirms H1", "the spike succeeded", "the issue is…").
- Adjectives that smuggle verdicts ("essentially", "mostly", "effectively", "appears to", "should have").
- Reconstructed entries — observations are written when observed, not recalled.
- Silent corrections — corrections are new `CORRECTION` entries that reference the original entry ID; originals are never edited.

---

## Template for new entries

```
### E-NNNN — <UTC timestamp> — <Entry class>

- Experiment: EXP-NNN
- Run: EXP-NNN/RUN-NN
- Hypothesis tested: <verbatim from pre-registration, or "N/A — lab process">
- Observations: <factual>
- Measurements: <value unit @ instrument>, or "none"
- Evidence refs: <path#sha256:…>
- Unexpected findings: <factual deviations> or "none"
- Result status: <RUN|RAN-OK|RAN-FAIL|INCONCLUSIVE|ABORTED|N/A>
- Next action: <planned step>
- Roles present: <initials / adversary-blinded note>
```

---

## Entries

### E-0001 — 2026-07-04T03:35:59Z — LAB-PROCESS

- Experiment: N/A (lab establishment)
- Run: N/A
- Hypothesis tested: N/A — lab process
- Observations: Engineering lab initialized at `/home/raj/Videos/projects/lab/`. Five durable process documents written: `LAB_README.md`, `EXPERIMENT_LOG.md`, `EXPERIMENT_CHECKLIST.md`, `EVIDENCE_INDEX.md`, `DECISION_RULES.md`. Upstream frozen pre-registrations confirmed present and unmodified: `P5_FALSIFICATION_EXPERIMENT.md` (Level 2 protocol, frozen), `LEVEL0_1_FEASIBILITY_GATE.md` (feasibility gate, Verdict B approved).
- Measurements: none
- Evidence refs:
  - `/home/raj/Videos/projects/lab/LAB_README.md`
  - `/home/raj/Videos/projects/lab/EXPERIMENT_LOG.md`
  - `/home/raj/Videos/projects/lab/EXPERIMENT_CHECKLIST.md`
  - `/home/raj/Videos/projects/lab/EVIDENCE_INDEX.md`
  - `/home/raj/Videos/projects/lab/DECISION_RULES.md`
- Unexpected findings: none
- Result status: N/A
- Next action: Pre-flight `EXPERIMENT_CHECKLIST.md` against EXP-001 (C4 spike) pre-registration to gate RUN-01 Go/No-Go.
- Roles present: lab-founder. No adversary role required for lab-process entries.

<!--
Subsequent entries are appended below this point. Entries are never edited;
corrections are added as new E-NNNN CORRECTION entries referencing the original.
-->

<!-- checkpoint: feat(lab): implement network partition test -->
