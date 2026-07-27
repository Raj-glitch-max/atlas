# 12 — ENGINEERING MANAGER (CNCF / MATURITY) — RAW FINDINGS

Phase 4 (Polish). Two questions only: would Atlas pass CNCF Sandbox and reach Incubation, and can a
stranger read Staff/Principal-track judgment off this repo without meeting its author? Every
missing-artifact finding carries an hour estimate — reframing paperwork as a cheap win is the job.

FINDING-EM-1
  Type:        sandbox-blocker
  Severity:    P0
  Claim:       No LICENSE. `git ls-files` shows no LICENSE/COPYING/NOTICE tracked. The README does
               not name a license either.
  Missing artifact: LICENSE (Apache-2.0)
  Effort:      1h
  Why it matters: To a TOC reviewer AND any adopter: unlicensed code is legally un-adoptable and
               Sandbox-disqualifying. This is THE ONE THING — highest value/hour in the review.
  Owner:       12  Co-owner: —

FINDING-EM-2
  Type:        governance / bus-factor
  Severity:    P0 (sandbox) / P1 (incubation)
  Claim:       No CODE_OF_CONDUCT.md, MAINTAINERS.md, CODEOWNERS, or repo-level GOVERNANCE.md.
               (CONTRIBUTING.md and an internal agents/GOVERNANCE.md exist; the CNCF-shaped files do
               not.)
  Missing artifact: CODE_OF_CONDUCT.md (CNCF CoC), MAINTAINERS.md, CODEOWNERS, GOVERNANCE.md
  Effort:      3h
  Why it matters: To a TOC reviewer: "who's your second maintainer?" is asked in the first five
               minutes. The files declare intent-to-scale even with one name. Bus factor is 1.
  Owner:       12  Co-owner: —

FINDING-EM-3
  Type:        sandbox-blocker / supply-chain (process)
  Severity:    P0
  Claim:       SECURITY.md exists but is 5 lines — no disclosure SLA, no supported-versions table,
               no CVE process. For a cryptographic project this is a credibility gap.
  Missing artifact: SECURITY.md disclosure policy (SLA + supported versions + embargo per OQ-2)
  Effort:      1h
  Why it matters: To an adopter and a researcher who finds a bug: there is a contact but no terms.
  Owner:       12  Co-owner: 04

FINDING-EM-7
  Type:        supply-chain
  Severity:    P1
  Claim:       No signed releases / SLSA provenance; go-jose not govulncheck-gated. An unsigned
               release of a verification library is a contradiction.
  Missing artifact: cosign signatures + SLSA provenance + govulncheck CI job
  Effort:      1–2d (co-own 06/03, S10)
  Why it matters: OpenSSF Scorecard and TOC both grep for this.
  Owner:       12  Co-owner: 06, 03

FINDING-EM-9
  Type:        incubation-blocker
  Severity:    P1
  Claim:       Zero adopters and (until PM-1) no plan to get one. Sandbox tolerates this; Incubation
               does not.
  Missing artifact: an adopter, via 10's beachhead interviews (S8)
  Effort:      1–2wk (10 owns)
  Why it matters: Adopters are the literal Incubation bar. 10's product work is the incubation
               critical path, not a side quest.
  Owner:       12  Co-owner: 10

FINDING-EM-10
  Type:        legibility / career-evidence
  Severity:    P2
  Claim:       The claim surface over/under-states relative to LIMITATIONS (see TW-2). The repo also
               blends project framing and portfolio framing; a TOC reviewer should see a project.
  Missing artifact: reconciled claims + a project-framed README (career narrative lives outside repo)
  Effort:      2h
  Owner:       12  Co-owner: 11

## REFUTED red flags (the persona predicted these; the repo disproves them — credit, don't file)
- "No ADRs — judgment invisible" → **REFUTED.** AD-001…017+ in ENGINEERING_DECISION_RECORD.md +
  13 journal entries. This was predicted as the single largest unforced error; it is instead a
  standout strength. (See C-3 in the conflict log.)
- "Owned weaknesses not published" → **REFUTED.** LIMITATIONS.md + THREAT_MODEL.md publish them
  well; they are most of a TAG-Security self-assessment already.
- "README reads production-ready" → **REFUTED.** README says v0.1-dev, reference implementation.

## MANDATORY CLOSING — TWO VERDICTS
1. **CNCF verdict.** *Sandbox: BLOCKED on paperwork only — LICENSE (1h), CoC/MAINTAINERS/CODEOWNERS/
   GOVERNANCE (3h), SECURITY.md terms (1h); ~1 day total, then ADMISSIBLE on technical merit.
   Incubation: BLOCKED on adopters (zero) and on the unrun two-domain substrate proof; critical
   path is S8 (a user) + S7 (the number).*
2. **Maturity verdict.** *This repository proves its author can pick a real unsolved infra problem,
   make and defend cryptographic decisions, hold a fail-closed line, record the hard calls as ADRs,
   and — rarest — state its own weaknesses unprompted. That is Staff-track. What it does not yet
   prove is influence beyond self: bus factor 1, no governance files, no adopters. The next ~40
   hours on LICENSE + governance + SECURITY.md + the differentiation doc + three interviews move
   Atlas's standing further than the next 400 hours of code. The code is not the constraint; the
   legibility and transferability of the judgment is — and half of that (the ADRs) is already done.*
