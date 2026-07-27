# ATLAS REVIEW — OPEN-QUESTION REGISTER

Conflicts the orchestrator deliberately did **not** resolve. Each is published with both
positions at full strength, the criterion that would settle it, and the cost of each branch.
An honestly escalated conflict is a better output than a confidently wrong ruling.

---

## OQ-1 — Is Atlas a product, or a contribution to an existing standard?

**Positions.**
- **10 (PM):** Existential and unanswered. Atlas is closest in design space to UCAN/Biscuit
  (its own README says so). If the winning move is to contribute the offline-verifiable-for-SPIFFE
  idea to an existing capability standard, then most of the runtime is scaffolding.
- **02 (Architect):** The two-domain, offline, fail-closed topology only pays off at a multi-org
  scale that may never arrive; the differentiator vs SPIFFE/OpenFGA/UCAN must be *written down*.
- **12 (EM/CNCF):** The TOC asks "why are you not SPIFFE / not UCAN / not OpenFGA" in the first
  meeting. Currently that answer lives only in the author's head and in scattered docs.

**What would settle it.** Three beachhead interviews (S8) plus a one-page written differentiation
doc. Until a real team with the delegation pain exists, the product-vs-standard question is being
answered by assumption.

**Cost of each branch.** *Product:* continued runtime investment that a standard might moot.
*Contribution:* surrendering the offline/fail-closed framing that is Atlas's actual novelty, into
a standards process that may not adopt it. Neither is cheap; the interviews are.

---

## OQ-2 — Should the security self-assessment publish the known architectural gaps?

**Positions.**
- **12 (EM/CNCF):** Publishing is the CNCF TAG-Security norm and converts honesty into
  credibility. `LIMITATIONS.md` already does most of this well.
- **04 (Security):** Publishing *unpatched specific exploits* arms attackers.

**Ruling criterion (04 rules per-item, R9/OQ-2 precedent).** *Architectural* limitations (offline
staleness window, within-window replay, single-hop scope, unrun substrate proof) → **publish**
(already largely done). *Specific unpatched exploits* → **embargo** under the SECURITY.md
disclosure process. This is not an unresolved conflict so much as a standing rule; it is listed
here because the SECURITY.md that operationalizes the embargo half does not yet exist (S3).

---

## OQ-3 — Does the wire-format spec ship before the beachhead?

**Positions.**
- **11 (Tech Writer):** For a protocol, the spec *is* the artifact; a reimplementer or a TOC
  reviewer cannot proceed from prose + a Go reference.
- **10 (PM):** Pre-user specification is premature commitment; you will want to change the format
  once a real user appears.

**What would settle it.** Ship `v0.1-draft`, marked explicitly unstable and append-only. Drafting
cost is low (the logical contracts and 28 vectors already pin most of it); the absence cost
— a second-language reimplementer, a TOC reviewer — is high. The "unstable" label buys 10's
flexibility without paying 11's absence cost. (This is S5.)

---

## OQ-4 — How much availability is the fail-closed posture actually costing?

**Positions.**
- **06 (SRE):** The coupling `verify_availability ≤ snapshot_availability` is arithmetic and
  currently unmeasured. Operators cannot set `R` responsibly without the number.
- **04 (Security):** Fail-closed is non-negotiable in the library; availability is bought via
  shorter TTLs and better snapshot distribution, never via a looser verifier.

**These are not in conflict — they are arguing about an unknown.** The resolution is not a
decision; it is a measurement. Run the two-domain substrate proof (S7), publish the availability
curve as a function of `R` and snapshot cadence, and the argument dissolves into a chosen operating
point. Until the number exists, both agents are right and neither can act.

---

### The shape of this register

Three of four open questions resolve to **"go get a number or a user,"** not "make a decision."
That is itself the review's most important structural finding: Atlas's remaining risk is
overwhelmingly *evidentiary*, not architectural. The design decisions have been made and recorded;
what is missing is the measurement and the user that would confirm or falsify them.
