---
date: 2026-07-21
slug: open-source-governance-baseline
artifact: maintainer decision — establish the CNCF-Sandbox governance baseline (LICENSE + community files) the review named as the top of the arbitrated stack
decision: Added the open-source governance baseline that was ABSENT and Sandbox-blocking — LICENSE (Apache-2.0), NOTICE, GOVERNANCE.md, MAINTAINERS.md, CODE_OF_CONDUCT.md (Contributor Covenant 2.1), .github/CODEOWNERS — and expanded SECURITY.md from a 5-line note into a real coordinated-disclosure policy (SLA, scope, supported-versions). Reconciled the one honest doc contradiction on the free side (LIMITATIONS.md §10 undercounted the SDKs). No code, no frozen doc, and no security boundary changed.
agents_consulted: [engineering-manager-cncf, security-engineer, technical-writer, principal-engineer]
overrides: false
related_entries: [language-neutral-conformance-vectors]
---

# Context

The 13-agent review (`agents/review-team/output/`) executed against the real
tree found **Tier 0 (confirmed exploits) and Tier 1 (silent-accept bugs) both
empty** — a rare, significant result for a verification project. With nothing
broken, the top of the arbitrated priority stack is **Tier 2: cheap unblockers**,
ranked there only because their value-per-hour is extreme.

The single highest-ratio item ("THE ONE THING"): the repository had **no
LICENSE**. `git ls-files` confirmed no LICENSE/COPYING/NOTICE was tracked, and
the frozen README named no license either. Unlicensed code is legally
un-adoptable and disqualifying at CNCF Sandbox — a ~1-hour fix gating every
adoption and CNCF conversation.

Adjacent Sandbox-blocking absences the review named (agent 12, EM-1..3/6):
no `CODE_OF_CONDUCT.md`, `MAINTAINERS.md`, `CODEOWNERS`, or repo-level
`GOVERNANCE.md`; a `SECURITY.md` that existed but was five lines with no
disclosure SLA, scope, or supported-versions table — a credibility gap for a
cryptographic project (agent 04, SEC-11).

# Decision

Established the governance baseline as new, non-frozen files:

- **LICENSE** — Apache-2.0 (the CNCF norm), copyright "The Atlas Authors" 2026.
- **NOTICE** — Apache attribution + third-party dependency licenses (go-jose,
  go-spiffe: Apache-2.0; x/crypto: BSD-3-Clause).
- **GOVERNANCE.md** — lightweight Sandbox-appropriate governance: roles, lazy
  consensus, a stricter bar for verification-core changes, maintainer lifecycle,
  and the two standing disciplines that predate this file (fail-closed on the
  security floor; judgment recorded via ADRs + journal).
- **MAINTAINERS.md** — one name, stated honestly, with an explicit succession
  path. Declaring the file with a single maintainer is deliberate: it makes the
  bus-factor shape explicit rather than implicit.
- **CODE_OF_CONDUCT.md** — Contributor Covenant 2.1, enforcement routed to the
  maintainer email.
- **.github/CODEOWNERS** — default owner plus explicit ownership of the
  security-critical core and the normative conformance vectors.
- **SECURITY.md** — rewritten into a coordinated-disclosure policy: private
  channels (GitHub Security Advisory + email), a response SLA, in/out-of-scope
  grounded in `THREAT_MODEL.md` and `LIMITATIONS.md`, a supported-versions
  table, and a 90-day coordinated-disclosure default.

Reconciled one honest contradiction on the free side: `LIMITATIONS.md` §10
undercounted the SDKs ("one SDK (Python)") — the repo ships Go/Python/TS client
SDKs. Corrected while preserving the load-bearing distinction the review cared
about: these are thin API clients, **not** independent verifier implementations,
so cross-implementation agreement remains a design goal until a second verifier
passes the conformance suite. The parallel claim on the **frozen** README is
left untouched (README is correct: it already says three SDKs); the frozen
context sheet fix and any frozen-doc reconciliation follow the amendment process
separately.

# Dissent / limits

- This is **plating, not cooking new food.** It changes the legibility of the
  project, not its behavior. The genuine technical bottleneck the review named —
  the unrun two-domain SPIRE substrate proof (E6/E7), which would convert the
  central offline claim from "by construction" to "demonstrated" — is untouched
  here and remains the incubation critical path.
- No frozen document, no `internal/` code, and no security boundary was
  modified. `make ci` (build + vet + tests + lints + frozen-doc integrity +
  import rules) must stay green; verified after this change.
- Sandbox admissibility now turns on paperwork that exists; **incubation** still
  turns on adopters (zero) and the substrate number — neither of which a
  governance file can supply.
