# Evidence Index

**Type.** Registry of every artifact class the lab may produce, where each is stored, how it is integrity-protected, and how long it is retained. This document defines the manifest schema and the storage layout. It is how a third party — a reviewer, an adversary, a future engineer, a founder — locates and verifies any claim the lab has made.

The lab's rule on evidence: **a conclusion is no stronger than the artifacts it rests on, and the artifacts are no stronger than their hashes.**

---

## 1. Naming convention (full)

| Entity | Pattern | Example |
|---|---|---|
| Experiment | `EXP-<NNN>` (zero-padded, 3 digits) | `EXP-001` |
| Run | `RUN-<NN>` (zero-padded, 2 digits) | `RUN-01` |
| Test case | `T<NN>` (per pre-registration; if pre-reg uses other labels, mirror them) | `T04` |
| Artifact file | `<T-NN-or-slug>-<descriptor>.<ext>` | `T04-replay-post-ttl.pcap` |
| Manifest | `MANIFEST.json` (one per run) | `EXP-001/RUN-01/MANIFEST.json` |
| Go-record | `GO.json` (one per run) | `EXP-001/RUN-01/GO.json` |
| Log entry | `E-<NNNN>` (in `EXPERIMENT_LOG.md`, not on disk) | `E-0007` |

**Slug rules.** Lowercase, ASCII, hyphen-separated, no spaces, no timestamps in the slug (the run directory carries the time context). A slug describes *what* the artifact is, not *when* it was made.

---

## 2. Directory layout

```
lab/
├── LAB_README.md
├── EXPERIMENT_LOG.md
├── EXPERIMENT_CHECKLIST.md
├── EVIDENCE_INDEX.md
├── DECISION_RULES.md
└── evidence/
    └── EXP-NNN/
        └── RUN-NN/
            ├── GO.json
            ├── MANIFEST.json
            ├── config/        # substrate configuration snapshots
            ├── pcap/          # packet captures
            ├── log/           # application and system logs captured during the run
            ├── bench/         # benchmark outputs (latency, throughput, counts)
            ├── trace/         # execution / distributed traces
            ├── token/         # issued and verified token samples (with secrets stripped or redacted)
            ├── screenshot/    # screenshots, where a UI is relevant (rare in this lab)
            ├── recording/     # terminal/screen recordings
            ├── raw/           # raw outputs that do not fit another category
            └── analysis/      # any derived analysis artifacts (never the conclusion itself)
```

Directories are created per-run, not pre-emptively for the whole lab. A run that produces no `screenshot/` artifacts has no `screenshot/` directory; the absence is meaningful and is not an error.

---

## 3. Artifact type registry

| Type | Produced when | Filename pattern | Integrity | Retention | Recorder |
|---|---|---|---|---|---|
| `config` | At setup (substrate snapshot) and at any config change during the run | `<component>-config.<fmt>` | SHA-256 + GPG/age signature by Engineer | Permanent | Engineer |
| `pcap` | During any case where network behavior is observed | `<T-NN>-<descriptor>.pcap` | SHA-256 | Permanent | Engineer |
| `log` | Continuously during the run; per-instrument | `<component>-<T-NN>.log` | SHA-256 | Permanent | Engineer |
| `bench` | When a measurement is taken (latency, throughput, count, time-series) | `<T-NN>-<metric>.<ext>` (e.g., `.csv`, `.json`) | SHA-256 | Permanent | Engineer |
| `trace` | When execution tracing is enabled for a case | `<T-NN>-trace.<fmt>` | SHA-256 | Permanent | Engineer |
| `token` | When a token is issued or verified; secrets redacted before storage | `<T-NN>-<role>-token.json` | SHA-256 + redaction certificate | Permanent | Engineer |
| `screenshot` | Only when a pre-reg case requires a visual artifact | `<T-NN>-<what>.png` | SHA-256 | 365 days | Engineer |
| `recording` | Only when a pre-reg case requires temporal replay | `<T-NN>-<what>.<fmt>` | SHA-256 | 730 days | Engineer |
| `raw` | Outputs not fitting another category | `<T-NN>-<descriptor>.<ext>` | SHA-256 | Permanent | Engineer |
| `analysis` | Derived artifacts (aggregations, plots, post-processing) — never the conclusion | `<T-NN>-<derived>.<ext>` | SHA-256 | Permanent | Engineer |
| `GO.json` | At pre-flight pass | `GO.json` | SHA-256 | Permanent | Standards Editor (sign) |
| `MANIFEST.json` | Started at pre-flight; sealed at run close | `MANIFEST.json` | SHA-256 of the sealed file recorded | Permanent | Engineer |

**Permanent retention** means: deleted only under a documented destruction event (see §6), with the prior hash preserved so the artifact remains verifiable-as-having-existed.

---

## 4. Per-run manifest specification

`MANIFEST.json` is the authoritative list of artifacts in a run. It is started at pre-flight and sealed at run close. Its schema:

```json
{
  "experiment": "EXP-001",
  "run": "EXP-001/RUN-01",
  "preRegistration": {
    "path": "../LEVEL0_1_FEASIBILITY_GATE.md",
    "sha256": "<hex>",
    "section": "Spike question (C4-only)"
  },
  "decisionRulesVersion": {
    "path": "DECISION_RULES.md",
    "sha256": "<hex>",
    "outcomeTableSection": "§3"
  },
  "goRecord": { "path": "GO.json", "sha256": "<hex>" },
  "configSnapshot": { "path": "config/<component>-config.<fmt>", "sha256": "<hex>" },
  "artifacts": [
    {
      "id": "ART-0001",
      "type": "pcap",
      "testCase": "T04",
      "path": "pcap/T04-replay-post-ttl.pcap",
      "sha256": "<hex>",
      "sizeBytes": 0,
      "witness": "sniffer self-test ART-0000",
      "producedAtUtc": "<ISO-8601>",
      "producedBy": "<role initial>",
      "redacted": false
    }
  ],
  "sealed": {
    "sealedAtUtc": "<ISO-8601>",
    "sealingHash": "<sha256 of MANIFEST.json at seal time>",
    "sealedBy": "<Standards Editor initial>"
  }
}
```

Rules:

- An artifact not in the manifest does not exist as evidence. The manifest is the set.
- A sealed manifest is immutable. Corrections append a `MANIFEST-SUPP-NN.json` supplement; the original is never edited.
- The `sealed.sealingHash` is the hash *after* `sealed` is filled; it is recomputed and matched at audit. Any mismatch is a manifest-integrity failure (`EXPERIMENT_CHECKLIST.md` A-6).
- Artifacts with `redacted: true` carry a redaction certificate (a separate signed file listing what was redacted and why, with the pre-redaction hash attested by the Engineer).

---

## 5. Integrity (hashing) procedure

1. The artifact is written to its final path under `lab/evidence/EXP-NNN/RUN-NN/<type>/`.
2. `sha256sum <path>` is computed. The hex digest is the artifact's identity.
3. The manifest entry is appended *before* the artifact is used in any further step; the artifact must be hashed at creation, not at end-of-run.
4. If the artifact is later re-read, the re-computed hash is matched to the manifest hash. A mismatch is recorded as `EXPERIMENT_CHECKLIST.md` A-6 and the run halts.
5. Cryptographic-signed artifacts (`config`, `GO.json`) are additionally signed; the signature is stored alongside the artifact as `<filename>.sig` and is itself hashed into the manifest.

Hashes are the lab's only identity for artifacts. Filenames are for humans; hashes are for proof.

---

## 6. Retention rules

| Class | Retention | At expiry |
|---|---|---|
| `config`, `pcap`, `log`, `bench`, `trace`, `token` (redacted), `raw`, `analysis`, `GO.json`, `MANIFEST.json` | Permanent | Archived (not destroyed); if destruction is ever required by a separate legal/contractual event, the destruction is logged with the prior hash, a reason, and an authority. |
| `screenshot` | 365 days | Destroyed; manifest entry retained with prior hash; if the screenshot was load-bearing for a conclusion, the conclusion is re-flagged as "evidence destroyed, conclusion rests on manifest reference only." |
| `recording` | 730 days | Same as `screenshot`. |
| `token` (un-redacted, if any were ever stored with secrets) | Never stored un-redacted; if this rule is violated, the artifact is destroyed immediately and the violation is logged as a process breach. | Immediate destruction + breach log entry. |

A conclusion whose supporting evidence has been destroyed (under retention rules) is not invalidated — the manifest's prior hash and the log's observation entry remain — but the conclusion is downgraded to "non-reproducible from primary artifacts" and any future challenge to it cannot be settled by re-examination.

---

## 7. What the index is not

- It is not a place to record conclusions. Conclusions live in the journal entry the lab produces per `DECISION_RULES.md` §6.
- It is not a place to record observations. Observations live in `EXPERIMENT_LOG.md`.
- It is not a place to record interpretation. Interpretation is what `DECISION_RULES.md` permits, bounded by what pre-registration froze.

The index is purely a registry and integrity spine. It tells a reader *what exists, where it is, and whether it is intact* — nothing more.

<!-- checkpoint: chore(fuzz): harden fuzzing harness execution -->
