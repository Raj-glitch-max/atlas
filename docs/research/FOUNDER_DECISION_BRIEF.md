# Founder Decision Brief

Prepared as a pre-decision uncertainty review, not a recommendation. Source of truth: `00_PROJECT_MISSION.md`, `01_CURRENT_STATE.md`, `02_SURVIVING_HYPOTHESES.md`, `03_RESEARCH_BACKLOG.md`, and the R1 problem set / council commentary recorded in `PROJECT_HISTORY.md`. All ten problems from the R1 cycle survived the validity gate and sequential council review, so all ten are treated as live candidates here — a low net score is a data point, not a disqualification, and a high net score is not a green light.

Every dimension rating below is my own analytical judgment as reviewer, not a repository fact. Where I state something as fact, it is traceable to a document. Where I state a rating or a risk, it is an assessment you should argue with, not a conclusion to defer to.

**A cross-cutting flag before the individual reviews:** every one of the ten problems currently carries Founder-advantage = 0 (H3, confirmed by the Founder Profile session returning no usable history). That means none of these problems has been evaluated with real founder-fit data — the comparison below is evaluating the *problems*, not "which problem fits you." That gap doesn't go away by picking a problem; it has to be closed separately (see Research Plan, Task 1).

A second flag: the "Confidence: High / Evidence quality: High" labels attached to all ten problems at the Domain Research stage (session 7) were *self-declared, pre-council*. The subsequent Empiricist pass at R1 flagged vendor-originated evidence on multiple problems (P1, P2, and implicitly P4's FinOps/McKinsey sourcing) without changing scores. "High" evidence quality and "vendor-sourced evidence" are not the same claim, and the repo currently treats them as compatible. That's worth being skeptical of going in.

---

## Problem-by-Problem Review

### P1 — Observability: cost & signal-to-noise at scale
**Net score: 5.** Council: storage/cost economics must shift for a wedge to exist (Red Team); evidence is mixed and partly vendor-originated (Empiricist); founder-buyer alignment plausible (Economist); workflow fit good (Operator).

| Dimension | Rating | Note |
|---|---|---|
| Engineering value | M | Real distributed-systems and cost-modeling depth, but the hard part (storage economics) is not yours to solve |
| Technical depth | M | Query/aggregation and sampling techniques are non-trivial; the ingestion pipeline itself is well-trodden ground |
| Resume value | M | "Built an observability cost tool" is a crowded line on a resume |
| Portfolio value | M | Demoable, but hard to make visually distinctive |
| Open-source potential | M–H | OSS observability tooling has a healthy adoption culture (OTel ecosystem) |
| Startup potential | L–M | Datadog/Honeycomb/Grafana pricing pressure is already commoditizing this |
| Defensibility | L | Cost/signal tooling is a feature incumbents ship, not a standalone moat |
| Long-term differentiation | L | Red Team's own framing ("storage/cost revolution required") says the differentiation depends on a market shift you don't control |
| Implementation realism (solo) | M | Ingestion + storage + query at meaningful scale is a lot for one engineer without cloud spend |
| Market maturity | H | Mature, well-served market |
| Incumbent pressure | H | Datadog, Honeycomb, Grafana Cloud, New Relic all compete directly |
| AI leverage | L–M | Anomaly/noise reduction is a plausible ML use, but it's additive, not foundational |
| Infrastructure complexity | H | Time-series storage at cost-competitive scale is genuinely hard |
| Research difficulty | M | The economics research (what actually drives observability spend) is more the gap than the engineering |
| Required domain expertise | H | Needs real production-scale familiarity with cardinality, retention, and query cost |
| Expected learning value | H | Whether or not it ships, the storage-economics learning is transferable |
| P(genuinely novel engineering) | M | Possible in the storage/query layer, unlikely in the product surface |
| P("just another AI wrapper") | L | Not AI-shaped as scoped; risk is "another dashboard," not "another wrapper" |

1. **What we know:** Ten-field Domain Research card exists (session 7); R1 net score 5; council flagged mixed/vendor evidence without eliminating it.
2. **What we only believe:** That "storage/cost revolution" framing is accurate — no cost-model evidence has actually been produced yet.
3. **What we do not know:** Who specifically would buy this pre-built vs. building it themselves on OTel; actual willingness to pay at any price point.
4. **Biggest technical uncertainty:** Whether a solo engineer can build storage/query performance that's cost-competitive with hyperscale vendors, at any scale worth demoing.
5. **Biggest product uncertainty:** Whether the wedge is "cheaper storage" (commodity, low differentiation) or "better signal" (harder, more defensible) — the problem card conflates both.
6. **Biggest market uncertainty:** Whether buyers switch observability vendors for cost reasons alone, given migration cost and lock-in.
7. **Experiment that removes the most uncertainty:** Interview 5–8 engineers who own an observability budget; ask what they've actually tried to reduce cost and why it failed or succeeded. This tests the "cost revolution required" premise directly.
8. **Evidence still missing before an experienced engineering leader approves building:** A real cost-per-GB benchmark against at least one incumbent, and a non-vendor-sourced account of actual buyer pain (not marketing pain).

---

### P2 — CI/CD pipeline flakiness & queue contention
**Net score: 6.** Council: Buildkite/Harness are direct competitors, a niche wedge is required (Red Team); some evidence is vendor-originated (Empiricist); per-seat pricing is tight (Economist); workflow fit good (Operator).

| Dimension | Rating | Note |
|---|---|---|
| Engineering value | M | Scheduling, flakiness-detection, and test-quarantine logic are real problems, but well-studied |
| Technical depth | M | Statistical flakiness detection has genuine depth; queue orchestration is largely solved |
| Resume value | M | Solid but generic "CI tooling" line |
| Portfolio value | M | Easy to demo, hard to make memorable |
| Open-source potential | H | CI tooling has strong OSS precedent (Buildkite agent, Tekton, etc.) |
| Startup potential | L | Explicitly named as directly competitive with funded incumbents by the project's own Red Team pass |
| Defensibility | L | Feature, not category — flakiness detection is being absorbed into CI platforms natively |
| Long-term differentiation | L | Same conclusion the council already reached: a real wedge, not a platform play, is required |
| Implementation realism (solo) | M–H | A flakiness-detection layer alone (not a full CI platform) is buildable solo |
| Market maturity | H | Very mature, heavily funded competitive set |
| Incumbent pressure | H | Buildkite, Harness, GitHub Actions, CircleCI all ship adjacent features |
| AI leverage | M | Flaky-test classification is a legitimate ML use case, not decorative |
| Infrastructure complexity | M | Moderate — depends on scope (plugin vs. platform) |
| Research difficulty | L | The problem is well documented in engineering literature already |
| Required domain expertise | M | CI/CD familiarity, not deep specialization |
| Expected learning value | M | Useful but not frontier learning |
| P(genuinely novel engineering) | L–M | Possible in flakiness statistics, unlikely in orchestration |
| P("just another AI wrapper") | M | Easy to slip into "GPT explains your failed build" territory if scoped loosely |

1. **What we know:** Net score 6 (tied highest); council explicitly named Buildkite and Harness as direct competitors.
2. **What we only believe:** That a "niche wedge" actually exists distinct from what incumbents already ship — this was asserted by Red Team, not evidenced.
3. **What we do not know:** Whether flakiness detection alone is sellable as a standalone product versus something teams expect bundled free.
4. **Biggest technical uncertainty:** Whether flaky-test classification can be made accurate enough (low false-positive rate) to be trusted in a merge-blocking workflow.
5. **Biggest product uncertainty:** Plugin-to-existing-CI vs. standalone platform — these are almost different products with different defensibility profiles.
6. **Biggest market uncertainty:** Per-seat pricing pressure was already flagged by the Economist; unclear if any price point clears both willingness-to-pay and unit economics.
7. **Experiment that removes the most uncertainty:** Build a minimal flakiness-classifier on a public OSS repo's CI history and measure precision/recall — this tests the technical premise cheaply, without building a product.
8. **Evidence still missing:** A concrete definition of "the wedge" that isn't already a checkbox feature in Buildkite/Harness roadmaps.

---

### P3 — Secrets & config drift across environments
**Net score: 5.** Council: drift detection is a feature, not a category (Cartographer); strong evidence base — HashiCorp, AWS, breach postmortems (Empiricist); IaC-layer feature risk (Red Team).

| Dimension | Rating | Note |
|---|---|---|
| Engineering value | M | Real correctness/security problem, but largely a detection-and-diffing problem |
| Technical depth | L–M | State comparison and drift detection are well-understood techniques |
| Resume value | M | Security-adjacent framing helps; drift-detection specifically is less impressive |
| Portfolio value | M | Demoable via before/after drift scenarios |
| Open-source potential | H | Strong precedent: driftctl-style tools, Terraform ecosystem plugins |
| Startup potential | L | Cartographer's own framing — "a feature, not a category" — is a direct warning against building a company here |
| Defensibility | L | HashiCorp, AWS Config, and IaC platforms are absorbing this natively |
| Long-term differentiation | L | Same conclusion twice from two different council members (Cartographer + Red Team) is a strong signal, not a coincidence |
| Implementation realism (solo) | H | Genuinely buildable solo at meaningful scope |
| Market maturity | H | Mature; well-documented via breach postmortems |
| Incumbent pressure | H | HashiCorp (Vault, Terraform), AWS Config, Pulumi all compete |
| AI leverage | L | Little natural AI surface — this is a deterministic diffing problem |
| Infrastructure complexity | M | Multi-cloud state reconciliation adds real complexity |
| Research difficulty | L | Evidence quality is already the strongest of the ten (breach postmortems, vendor docs, standards) |
| Required domain expertise | M | IaC and cloud-provider API familiarity required |
| Expected learning value | M | Solid infra learning, less frontier |
| P(genuinely novel engineering) | L | This is well-trodden engineering ground |
| P("just another AI wrapper") | L | Not AI-shaped; if anything, risks being "just another IaC linter" |

1. **What we know:** Two independent council members (Cartographer and Red Team) both converged on "feature, not category" / "IaC-layer feature risk" — the strongest internal consensus against category-defining potential of any problem in the set.
2. **What we only believe:** That there's meaningful whitespace between what Vault/AWS Config already do and what a standalone tool could add.
3. **What we do not know:** Whether any team would pay for drift detection separately rather than demanding it be bundled into their existing IaC tool.
4. **Biggest technical uncertainty:** Multi-cloud, multi-format state reconciliation (Terraform state, cloud-native config, secrets managers) at acceptable accuracy.
5. **Biggest product uncertainty:** Standalone tool vs. plugin — and if plugin, into which ecosystem, given the fragmentation of IaC tooling.
6. **Biggest market uncertainty:** Whether "detection" alone is valuable without "remediation," and whether remediation is a much bigger, riskier scope.
7. **Experiment that removes the most uncertainty:** Talk to 5 platform engineers about the last drift incident they had — ask specifically what tool they reached for and why it wasn't enough. This directly tests the Cartographer's "feature not category" claim.
8. **Evidence still missing:** A documented case where an existing IaC/security tool's native drift detection genuinely failed a team, not just was absent.

---

### P4 — Cloud cost attribution & reduction
**Net score: 6.** Council: must beat a tag-discipline prerequisite (Red Team); FinOps + McKinsey + case-study evidence (Empiricist); long, consultative sales cycle (Economist); requires FinOps practice maturity (Operator).

| Dimension | Rating | Note |
|---|---|---|
| Engineering value | M | Attribution/allocation logic is real, but the hard part is organizational (tagging discipline), not technical |
| Technical depth | L–M | Mostly aggregation, allocation rules, and reporting |
| Resume value | M–H | "Cloud cost/FinOps" carries real market recognition right now |
| Portfolio value | M | Dashboards demo well but are generic-looking |
| Open-source potential | M | OpenCost and similar OSS projects already exist |
| Startup potential | M | Real budget line exists (FinOps is now an established discipline) but sales motion is enterprise-shaped |
| Defensibility | L–M | Cost-visibility tools converge toward commodity dashboards quickly |
| Long-term differentiation | L | Red Team's "must beat tag-discipline prerequisite" is a structural problem: the product's value is capped by a customer behavior it can't control |
| Implementation realism (solo) | M | Buildable, but the go-to-market (long consultative sales cycle, per Economist) is not solo-friendly |
| Market maturity | H | FinOps Foundation, established vendors (CloudHealth, Kubecost/OpenCost, Vantage) |
| Incumbent pressure | H | Strong incumbents plus native cloud-provider cost tools (AWS Cost Explorer, etc.) |
| AI leverage | L–M | Anomaly detection and forecasting are legitimate but incremental uses |
| Infrastructure complexity | M | Multi-cloud billing API integration is fiddly but not deep |
| Research difficulty | L | Evidence base is the strongest cited (McKinsey, case studies) of any problem |
| Required domain expertise | M | FinOps domain knowledge more than deep systems knowledge |
| Expected learning value | M | Valuable business-context learning, less deep-systems learning |
| P(genuinely novel engineering) | L | Low — this is an integration and reporting problem |
| P("just another AI wrapper") | L–M | Real risk of becoming "ChatGPT summarizes your AWS bill" if scoped lazily |

1. **What we know:** Tied for highest net score (6); has the best-documented evidence base of the ten problems (named sources: FinOps Foundation-style material, McKinsey, case studies).
2. **What we only believe:** That "tag-discipline prerequisite" is solvable by the product itself rather than being an organizational precondition outside the product's control.
3. **What we do not know:** Actual sales-cycle length and CAC for a solo/small team selling into this category — the Economist flagged risk but didn't quantify it.
4. **Biggest technical uncertainty:** Whether accurate cost attribution is even possible without the tagging discipline the Red Team flagged as a prerequisite — i.e., the product may only work for customers who don't need it.
5. **Biggest product uncertainty:** Whether the product is a dashboard (commodity, low differentiation) or a workflow/process tool that changes tagging behavior (harder, more defensible, unproven).
6. **Biggest market uncertainty:** Whether a solo-built product can close FinOps-category enterprise sales cycles at all, given the Economist's own "long, consultative" framing.
7. **Experiment that removes the most uncertainty:** Shadow or interview 3–5 FinOps practitioners through an actual cost-attribution exercise, specifically to test whether tagging discipline is the blocker in practice, not just in theory.
8. **Evidence still missing:** Any data on solo-founder or small-team success selling into FinOps-category enterprise budgets, given the sales-cycle risk already flagged internally.

---

### P5 — IAM least-privilege / service-to-service identity
**Net score: 5 (lifted from 4).** Council: SPIFFE-style standardization opens room (Cartographer, lift); OWASP/DBIR/NIST evidence (Empiricist); cloud-bundle capture risk (Economist); requires security-process integration (Operator).

| Dimension | Rating | Note |
|---|---|---|
| Engineering value | H | Identity/authZ systems are genuinely deep, well-regarded engineering |
| Technical depth | H | Cryptographic identity, attestation, and policy engines are non-trivial |
| Resume value | H | Security-infrastructure work carries strong signal |
| Portfolio value | H | A working SPIFFE/SPIRE-adjacent system is a credible, specific artifact |
| Open-source potential | H | SPIFFE/SPIRE is itself an active OSS/CNCF ecosystem to plug into or extend |
| Startup potential | M | Real category, but cloud providers are absorbing identity natively (Economist's "cloud-bundle capture risk") |
| Defensibility | M | Standards-based work is easier to differentiate on execution than on secrecy |
| Long-term differentiation | M | The Cartographer's lift (SPIFFE standardization "opens room") is the most substantively bullish single note in the entire council record |
| Implementation realism (solo) | L–M | Correctness-critical security systems are unforgiving for a solo engineer to ship safely |
| Market maturity | M | Standard exists (SPIFFE) but adoption is still early-to-mid |
| Incumbent pressure | M | Cloud-native IAM (AWS IAM, GCP Workload Identity) competes, but standards-based portability is a real counter-argument |
| AI leverage | L | Minimal natural AI surface — this is a correctness and cryptography problem |
| Infrastructure complexity | H | Distributed trust, attestation, and rotation at scale are hard |
| Research difficulty | M–H | Requires real security domain depth to do credibly |
| Required domain expertise | H | Security engineering background matters a lot here |
| Expected learning value | H | Among the highest of the ten — deep, transferable, hard-to-fake expertise |
| P(genuinely novel engineering) | M–H | Highest of the ten alongside P1's storage-economics angle |
| P("just another AI wrapper") | L | Essentially zero — this problem has no AI-wrapper shape at all |

1. **What we know:** The only problem in the set that the council actively *upgraded* mid-cycle (Cartographer lift), on the specific stated basis that SPIFFE-style standardization changes the opportunity picture.
2. **What we only believe:** That standards-based portability is enough to counter "cloud-bundle capture" — the Economist flagged the risk but it was never tested against the Cartographer's optimism.
3. **What we do not know:** Actual buyer segment — is this sold to security teams, platform teams, or absorbed as an open-source infrastructure component with no direct buyer at all?
4. **Biggest technical uncertainty:** Whether a solo engineer can build and, critically, be trusted to ship a security-critical identity system — trust, not raw capability, may be the real constraint.
5. **Biggest product uncertainty:** Whether this is a product at all versus an open-source infrastructure contribution with no monetizable surface — those are very different paths from the same technical work.
6. **Biggest market uncertainty:** Whether cloud providers' native workload-identity features close the gap this problem is targeting before a small team could ship.
7. **Experiment that removes the most uncertainty:** A scoped technical spike — implement a minimal SPIFFE-compatible identity issuance flow and see how much of the "hard part" is actually solved by existing SPIRE tooling versus requiring new work. This tests both technical feasibility and differentiation simultaneously.
8. **Evidence still missing:** Any account of who currently struggles with this in production and what they've tried from the existing SPIFFE/SPIRE ecosystem before concluding it's insufficient.

---

### P6 — Compliance & audit evidence collection
**Net score: 1.** Council: framework-specific niche possible, e.g. EU AI Act (Red Team); vendor-bias caution (Empiricist); audit-firm bundling risk (Economist); requires compliance-officer sponsorship (Operator).

| Dimension | Rating | Note |
|---|---|---|
| Engineering value | L | Largely evidence-collection and reporting plumbing, not deep systems work |
| Technical depth | L | Low — mostly integration and document generation |
| Resume value | L–M | Compliance-tooling experience is a narrower signal than infra/security work |
| Portfolio value | L | Hard to make visually or technically compelling for a portfolio |
| Open-source potential | L | Compliance tooling is rarely open-sourced due to liability and customization needs |
| Startup potential | L | Lowest net score in the set for a reason — Operator flagged a sponsorship dependency (compliance officer buy-in) that a solo builder can't easily generate |
| Defensibility | L | Audit-firm bundling risk (Economist) means the natural incumbent isn't even a software vendor — it's the audit relationship itself |
| Long-term differentiation | L | Framework-specific niches (e.g., EU AI Act) are real but narrow and regulation-dependent, i.e., differentiation with an expiration date tied to policy |
| Implementation realism (solo) | L–M | Buildable narrowly, but requires ongoing regulatory-interpretation work that doesn't scale with one engineer |
| Market maturity | M | Established category (Vanta, Drata, etc.) with real budget |
| Incumbent pressure | H | Vanta/Drata-class incumbents plus the audit firms themselves |
| AI leverage | M | Evidence summarization is a plausible legitimate use, but also an easy wrapper trap |
| Infrastructure complexity | L | Low — mostly integrations, not infrastructure |
| Research difficulty | M | Requires real regulatory/compliance domain literacy to avoid building something wrong |
| Required domain expertise | M–H | Compliance-domain knowledge, not engineering depth, is the bottleneck |
| Expected learning value | L–M | Domain learning is narrow and regulatory, less transferable |
| P(genuinely novel engineering) | L | Lowest of the ten |
| P("just another AI wrapper") | M–H | High risk — "AI reads your logs and writes your audit evidence" is close to the median AI-wrapper pitch |

1. **What we know:** Lowest net score of the ten (1); every council note attached to it is a caution (vendor bias, bundling risk, sponsorship dependency).
2. **What we only believe:** That a "framework-specific niche" (e.g., EU AI Act) is a real wedge — this is speculative and was flagged by Red Team as merely "possible."
3. **What we do not know:** Whether any regulatory niche is stable and large enough to build a durable business on, versus a moving target that changes faster than a solo builder can track.
4. **Biggest technical uncertainty:** Minimal — this problem's uncertainty is almost entirely non-technical.
5. **Biggest product uncertainty:** Whether the buyer is the compliance officer, the audit firm, or the engineering org — three very different products.
6. **Biggest market uncertainty:** Whether audit firms themselves become the actual competitive bottleneck (bundling risk), making a standalone tool structurally excluded from the sale.
7. **Experiment that removes the most uncertainty:** Interview 2–3 compliance officers or auditors about how a specific regulatory framework's evidence is currently collected, specifically probing the audit-firm bundling dynamic.
8. **Evidence still missing:** Any concrete answer to "who actually signs the check," which the current problem card doesn't resolve.

---

### P7 — Kubernetes operational complexity
**Net score: 0 (dropped from 1).** Council: pain is a value-promise gap, not missing tooling (Cartographer, drop); managed K8s continues to improve (Red Team); cloud-bundle capture (Economist); requires multi-stakeholder coordination (Operator).

| Dimension | Rating | Note |
|---|---|---|
| Engineering value | L | The council's own conclusion is that the pain isn't a tooling gap at all |
| Technical depth | M | K8s internals are deep, but that's not where this problem's stated pain lives |
| Resume value | L | "Built a Kubernetes tool" is one of the most saturated resume lines in infra |
| Portfolio value | L | Extremely crowded category to stand out in |
| Open-source potential | M | K8s ecosystem has strong OSS culture, but that cuts against differentiation too |
| Startup potential | L | Explicitly the only problem the council actively downgraded (dropped from 1 to 0) |
| Defensibility | L | Managed K8s (EKS/GKE/AKS) keeps absorbing exactly this complexity |
| Long-term differentiation | L | The Cartographer's drop rationale — "value-promise gap, not missing tooling" — is a direct statement that the premise itself is wrong |
| Implementation realism (solo) | L | Requires "multi-stakeholder coordination" per Operator — structurally hard for a solo builder |
| Market maturity | H | Extremely mature, saturated tooling landscape |
| Incumbent pressure | H | Managed K8s providers plus a huge OSS tooling ecosystem |
| AI leverage | L | No clear natural AI surface distinct from general ops-copilot territory |
| Infrastructure complexity | H | Genuinely complex domain, but complexity ≠ opportunity here |
| Research difficulty | L | Well-documented territory; the issue isn't lack of information |
| Required domain expertise | H | Deep K8s expertise required just to have credibility in the space |
| Expected learning value | M | Learning K8s internals deeply is valuable regardless of product outcome |
| P(genuinely novel engineering) | L | Low — this is the most "solved" category in the set |
| P("just another AI wrapper") | M | "AI-powered kubectl" is an extremely well-worn wrapper pattern already |

1. **What we know:** The only problem the council downgraded during review, and the only one with a Cartographer *drop* rather than a lift or hold — the strongest negative signal in the R1 record.
2. **What we only believe:** Nothing further should be believed here without new evidence — the existing council record already argues against the premise.
3. **What we do not know:** Whether there's a genuinely different framing of K8s pain (not "more tooling") that the original Domain Research pass simply didn't surface.
4. **Biggest technical uncertainty:** Effectively moot given the Cartographer's framing — the uncertainty is conceptual (is there a real problem here) more than technical.
5. **Biggest product uncertainty:** Whether any product framing survives the "value-promise gap, not missing tooling" objection.
6. **Biggest market uncertainty:** Whether managed K8s providers close the remaining gap before any new entrant could matter.
7. **Experiment that removes the most uncertainty:** Before any further work, directly test the Cartographer's objection — find 3 platform teams and ask what they wish existed, explicitly screening for "another tool" answers vs. deeper process/expectation problems.
8. **Evidence still missing:** A reframed problem statement that isn't just "K8s is complex" — as currently defined, there isn't enough here for an engineering leader to approve further investment.

---

### P8 — Incident response coordination overhead
**Net score: 3.** Council: postmortem-quality niche (Red Team); incident.io-style documented evidence (Empiricist); per-seat recurring revenue model (Economist); workflow fit good (Operator).

| Dimension | Rating | Note |
|---|---|---|
| Engineering value | M | Coordination/workflow engineering, moderate depth |
| Technical depth | L–M | Mostly state machines, notifications, and integrations |
| Resume value | M | Reasonable, not standout |
| Portfolio value | M | Demoable via a simulated incident scenario |
| Open-source potential | M | Some OSS precedent, less than CI/CD or IaC categories |
| Startup potential | L–M | incident.io is a direct, well-funded incumbent already named in the evidence base |
| Defensibility | L | Workflow tools converge to commodity fast once the core loop is copied |
| Long-term differentiation | L–M | "Postmortem-quality niche" (Red Team) is a specific enough wedge to be worth testing, but narrow |
| Implementation realism (solo) | M | Buildable at MVP scope solo |
| Market maturity | M–H | Established category with a clear incumbent |
| Incumbent pressure | H | incident.io, PagerDuty, Opsgenie all compete directly |
| AI leverage | M | Postmortem drafting/summarization is a legitimate, bounded AI use case |
| Infrastructure complexity | L–M | Moderate — mostly integration-heavy, not deep infra |
| Research difficulty | L | Evidence base already cites a named, credible incumbent's public material |
| Required domain expertise | M | SRE/on-call process familiarity needed |
| Expected learning value | M | Useful but not frontier |
| P(genuinely novel engineering) | L | Low — workflow and integration engineering |
| P("just another AI wrapper") | M | "AI writes your postmortem" is a plausible-but-common pitch |

1. **What we know:** Net score 3; evidence base explicitly references incident.io's public material, meaning the closest incumbent is also the cited evidence source.
2. **What we only believe:** That "postmortem quality" specifically (versus coordination generally) is a defensible enough wedge — Red Team named it narrowly, and narrowness cuts both ways (focus vs. small market).
3. **What we do not know:** Whether teams already using incident.io or PagerDuty perceive a real gap, or whether this is solving a problem the incumbents already solve adequately.
4. **Biggest technical uncertainty:** Low — this is one of the more technically tractable problems in the set; the uncertainty is elsewhere.
5. **Biggest product uncertainty:** Whether "better postmortems" alone is a sellable wedge distinct from full incident-coordination platforms, or whether it's a feature that gets absorbed.
6. **Biggest market uncertainty:** Per-seat recurring pricing (Economist) against an incumbent with strong brand recognition in the exact niche named.
7. **Experiment that removes the most uncertainty:** Interview teams currently using incident.io/PagerDuty specifically about postmortem quality — this directly tests whether the named wedge is real or already served.
8. **Evidence still missing:** A specific, named gap in incident.io's actual postmortem workflow, rather than a general "quality" claim.

---

### P9 — On-call burden & toil reduction
**Net score: 6.** Council: tight ICP / narrow feature (Red Team); SRE book, survey, and blog evidence (Empiricist); small-ACV / services-business risk (Economist); workflow fit good (Operator).

| Dimension | Rating | Note |
|---|---|---|
| Engineering value | M | Toil-reduction automation has real engineering content, moderate depth |
| Technical depth | M | Alert correlation, runbook automation — solid but not frontier |
| Resume value | M | SRE-adjacent framing has decent recognition |
| Portfolio value | M | Demoable via before/after on-call load scenarios |
| Open-source potential | M | Reasonable OSS precedent (runbook automation projects exist) |
| Startup potential | L | Economist explicitly flagged small-ACV / services-business risk — a structural revenue-model warning, not a minor caveat |
| Defensibility | L–M | Tight ICP (Red Team) could mean real focus or could mean "too small to matter" — ambiguous which |
| Long-term differentiation | M | Evidence base (SRE literature, survey data) is genuinely strong and less vendor-tainted than several other problems |
| Implementation realism (solo) | M–H | Scoped narrowly, this is buildable solo |
| Market maturity | M–H | Established category (PagerDuty, Opsgenie, incident.io overlap) |
| Incumbent pressure | M–H | Direct overlap with P8's incumbent set |
| AI leverage | M | Alert triage/correlation is a legitimate, bounded AI use case |
| Infrastructure complexity | M | Moderate |
| Research difficulty | L | Strong non-vendor evidence base already exists (SRE book, survey data) |
| Required domain expertise | M | On-call/SRE process familiarity |
| Expected learning value | M | Solid, not frontier |
| P(genuinely novel engineering) | L–M | Low-to-moderate |
| P("just another AI wrapper") | M | "AI triages your alerts" is a common pitch pattern |

1. **What we know:** Tied for highest net score (6); evidence base is unusually strong and cites a named book, survey data, and blog sources rather than primarily vendor material — the least vendor-tainted evidence of the ties.
2. **What we only believe:** That the "tight ICP" the Red Team flagged is a feature (focus) rather than a bug (too narrow to sustain a business) — this hasn't been tested either way.
3. **What we do not know:** What "small-ACV / services-business risk" actually means in dollar terms — the Economist flagged the shape of the risk, not its magnitude.
4. **Biggest technical uncertainty:** Whether alert-correlation/triage automation can be accurate enough to be trusted without human override, which is the difference between a real toil reduction and a new noise source.
5. **Biggest product uncertainty:** Whether this is a standalone product or a feature that PagerDuty/Opsgenie/incident.io will simply ship natively, given the incumbent overlap with P8.
6. **Biggest market uncertainty:** Whether the flagged small-ACV risk means this is structurally a services business in disguise, which is a very different founder commitment than a software business.
7. **Experiment that removes the most uncertainty:** Get real numbers — interview 5+ on-call engineers on actual toil hours and current spend on tools addressing it, to directly test the small-ACV concern with data instead of inference.
8. **Evidence still missing:** Concrete pricing/willingness-to-pay data; right now this problem has the best pain evidence in the set and the weakest revenue-model evidence.

---

### P10 — Production debugging / root-cause speed
**Net score: 4.** Council: Honeycomb/eBPF-era tooling flattens the wedge; AI-summary approaches are still research-grade (Red Team); small ACV (Economist); workflow fit good (Operator).

| Dimension | Rating | Note |
|---|---|---|
| Engineering value | M–H | Root-cause analysis and distributed tracing correlation are genuinely hard problems |
| Technical depth | H | Among the deepest technical problems in the set if pursued seriously |
| Resume value | M–H | "Built root-cause analysis tooling" is a strong, specific line if it actually works |
| Portfolio value | M | Compelling if demoable; hard to demo convincingly without a real production system |
| Open-source potential | M | Some precedent (eBPF tooling ecosystem), less than infra-config categories |
| Startup potential | L | Small ACV flagged by Economist, on top of Red Team explicitly calling AI-summary approaches "research-grade" — i.e., not yet productizable |
| Defensibility | L–M | eBPF-era tooling (per Red Team) is already flattening the differentiation this would need |
| Long-term differentiation | L | The council's own note is the clearest warning in the set: the core proposed mechanism (AI summarization) isn't mature enough to differentiate on |
| Implementation realism (solo) | L | Genuinely difficult — distributed tracing correlation at production scale is a large undertaking |
| Market maturity | H | Honeycomb, Datadog APM, and eBPF-native tools are all active here |
| Incumbent pressure | H | Well-funded incumbents plus emerging eBPF-native players |
| AI leverage | L (as currently framed) | Red Team's own "research-grade" label is a direct signal this isn't ready |
| Infrastructure complexity | H | High — distributed tracing and correlation at scale |
| Research difficulty | H | The underlying ML/analysis problem is not solved industry-wide, let alone solo |
| Required domain expertise | H | Deep distributed-systems and observability expertise required |
| Expected learning value | H | If pursued for learning rather than shipping, this is one of the highest-value problems in the set |
| P(genuinely novel engineering) | M–H | Real novel engineering is plausible here, specifically because the AI-summary approach is acknowledged as unsolved |
| P("just another AI wrapper") | H | This is the problem most explicitly flagged as at risk of exactly that — the council named it directly |

1. **What we know:** Net score 4; the council's Red Team note is unusually specific and unusually blunt — "AI-summary research-grade" is closer to a warning than a scoring note.
2. **What we only believe:** That root-cause speed is even the right framing — "speed" implies the bottleneck is analysis time, when it may instead be tooling fragmentation or team process.
3. **What we do not know:** Whether there's a viable non-AI-summary approach to this problem that would sidestep the Red Team's specific objection.
4. **Biggest technical uncertainty:** Whether automated root-cause correlation can clear a usable accuracy bar at all with current techniques — this is an open research question, not an engineering execution question.
5. **Biggest product uncertainty:** Whether the product is honest about being "research-grade" (a research tool, sold as such) or oversells itself as production-ready, which the evidence suggests would be premature.
6. **Biggest market uncertainty:** Small ACV combined with strong incumbent presence (Honeycomb, Datadog) — unclear if there's room for a new entrant at any price.
7. **Experiment that removes the most uncertainty:** A technical spike measuring actual root-cause-identification accuracy on a real (or realistic) distributed trace dataset, honestly reported — this directly tests the Red Team's "research-grade" claim rather than assuming it.
8. **Evidence still missing:** Any accuracy benchmark at all. Right now this problem is being evaluated entirely on narrative, with the one specific technical claim in the record being a caution against it.

---

## Comparison Matrix — Research Readiness (not desirability)

This matrix does not rank which problem is best. It ranks which problems could support an informed founder decision **soonest**, based on how much additional uncertainty-reduction each one needs.

| Problem | Net score | Evidence quality (non-vendor) | Core uncertainty type | Additional research needed before decision |
|---|---|---|---|---|
| P9 — On-call toil | 6 | High (book/survey/blog) | Revenue-model risk (small ACV) | **Low–Medium** — pain evidence is strong; the gap is pricing/ACV data, which is fast to gather |
| P3 — Secrets/config drift | 5 | High (breach postmortems, standards docs) | Category viability (two council members said "feature not category") | **Low–Medium** — evidence is strong; the open question (is there a standalone wedge) is answerable via a small number of interviews |
| P8 — Incident coordination | 3 | Medium (named incumbent's own material) | Wedge specificity vs. incumbent | **Medium** — needs a concrete gap identified in incident.io's workflow, not just general interviews |
| P4 — Cloud cost/FinOps | 6 | High (McKinsey, case studies) | Go-to-market realism for a solo builder | **Medium** — pain and evidence are solid; the real unknown is sales-motion feasibility, which takes longer to test than a survey |
| P2 — CI/CD flakiness | 6 | Medium (partly vendor-sourced) | Defensible wedge vs. commodity feature | **Medium** — technical premise is cheaply testable (a classifier spike), but competitive differentiation remains genuinely unclear |
| P1 — Observability cost | 5 | Medium (mixed, partly vendor) | Structural market dependency (needs a "storage/cost revolution") | **Medium–High** — depends on an external condition (cost economics shifting) that isn't founder-controllable to validate quickly |
| P5 — IAM/service identity | 5 | High (OWASP/DBIR/NIST) | Buyer existence and solo-shippability of security-critical software | **Medium–High** — strong evidence, but requires an actual technical spike (not just interviews) to know if it's real |
| P6 — Compliance/audit | 1 | Low–Medium (vendor-bias flagged) | Who the buyer even is | **High** — three different possible buyers were never disambiguated |
| P10 — Root-cause debugging | 4 | Low (no benchmark exists) | Whether the core technical approach works at all | **High** — the central technical claim is currently unverified and explicitly flagged as unproven |
| P7 — Kubernetes complexity | 0 | Medium (well-documented domain, weak premise) | Whether the problem premise is even valid | **High** — the council already concluded the stated premise is likely wrong; research has to start by re-establishing there's a problem at all |

**Read of this matrix:** the problems that are cheapest to de-risk further (P9, P3) are not the same as the problems with the highest net score (P4, P2 tied with P9 at 6) or the ones with the deepest technical ceiling (P5, P10). Net score measured opportunity-minus-friction at a point in time; it did not measure "distance to a confident founder decision." Those are different axes, and conflating them would be a mistake.

---

# Research Plan

This plan investigates uncertainty, not solutions. No task below produces a product decision by itself; each is designed to remove one specific unknown identified above.

### Task 1 — Founder-advantage intake backfill
- **Objective:** Replace the current A6=0 (no founder-advantage evidence) with real data across all ten problems.
- **Why it matters:** Every dimension rating above evaluates the *problem*. None of it accounts for founder fit, which H5/H3 already flagged as unresolved and structurally important.
- **Expected output:** A founder-profile document with actual prior work, network, and domain exposure mapped against each of the ten problems.
- **Success criteria:** At least one problem can be scored with founder-advantage > 0 based on real evidence, not stated interest.
- **Decision it enables:** Whether founder fit should be a tiebreaker among the research-ready problems (P9, P3) or a reason to weight a technically deeper problem (P5) more heavily.
- **Estimated effort:** 1–2 focused sessions; primarily reflective/interview work, not building anything.

### Task 2 — Non-vendor pain validation interviews (P9, P3, P8)
- **Objective:** Independently verify pain and current-workaround behavior for the three problems with the strongest existing evidence bases, via direct interviews rather than published sources.
- **Why it matters:** All prior evidence for these problems is secondary (books, surveys, vendor docs). None of it is primary, founder-collected evidence.
- **Expected output:** 5–8 structured interview notes per problem with practitioners who currently experience the stated pain.
- **Success criteria:** Interviews either corroborate or contradict the existing council notes (e.g., does "tight ICP" on P9 read as focus or as too-small in practice).
- **Decision it enables:** Whether to proceed to a technical spike on any of these three, or deprioritize based on contradicted evidence.
- **Estimated effort:** 1–2 weeks, interview scheduling-bound rather than engineering-bound.

### Task 3 — Technical feasibility spike: root-cause accuracy (P10)
- **Objective:** Produce an honest accuracy benchmark for automated root-cause correlation on a real or realistic distributed trace dataset.
- **Why it matters:** This is the only problem where the central technical premise is currently unverified rather than merely under-evidenced — Red Team's "research-grade" note is a direct challenge that hasn't been tested.
- **Expected output:** A measured precision/recall (or equivalent) figure for root-cause identification against a known dataset, reported honestly including failure cases.
- **Success criteria:** A number exists where none did before, whatever that number is.
- **Decision it enables:** Whether P10 is a viable near-term build target or a multi-year research bet that shouldn't be scoped as a product yet.
- **Estimated effort:** 1–2 weeks of focused technical work; this is the most build-adjacent task in the plan but its output is a measurement, not a product.

### Task 4 — Technical feasibility spike: SPIFFE-based identity issuance (P5)
- **Objective:** Implement a minimal SPIFFE-compatible identity issuance flow to determine how much of the "hard part" existing SPIRE tooling already solves.
- **Why it matters:** P5 has the strongest single positive signal in the entire council record (the Cartographer lift), but no technical validation has occurred.
- **Expected output:** A working minimal spike plus a written account of what remains genuinely unsolved after using existing SPIFFE/SPIRE components.
- **Success criteria:** A clear answer to "is there real engineering work left here, or does the existing ecosystem already cover it."
- **Decision it enables:** Whether P5 is pursued as a differentiated build or reclassified as "contribute to existing OSS" rather than "build a product."
- **Estimated effort:** 1–2 weeks.

### Task 5 — Buyer disambiguation (P6)
- **Objective:** Determine definitively who the buyer is for compliance/audit evidence tooling — compliance officer, audit firm, or engineering org.
- **Why it matters:** This is currently unknown and materially changes what the product would even be; net score alone (1) doesn't capture why it's unresolved.
- **Expected output:** A short interview-based memo identifying the actual budget holder and decision process in 2–3 real organizations.
- **Success criteria:** A named buyer type with a described budget process, or a documented conclusion that no clear buyer exists.
- **Decision it enables:** Whether P6 should remain in consideration at all, or be formally retired given its already-low net score and the council's stacked cautions.
- **Estimated effort:** 3–5 days.

### Task 6 — Premise re-validation (P7)
- **Objective:** Directly test the Cartographer's "value-promise gap, not missing tooling" objection before any further Kubernetes-complexity work is considered.
- **Why it matters:** This is the only problem the council actively downgraded; proceeding without addressing that specific objection would ignore the project's own review discipline.
- **Expected output:** A short memo from 3 platform-engineer interviews, explicitly screening for whether the pain is tooling-shaped or expectation/process-shaped.
- **Success criteria:** A clear answer on whether a reframed, non-tooling problem statement exists, or confirmation that the Cartographer's objection holds.
- **Decision it enables:** Whether P7 stays in the candidate set in any form, or is formally retired.
- **Estimated effort:** 3–5 days.

### Task 7 — Vendor-bias re-audit of existing evidence (P1, P2, P4)
- **Objective:** Re-examine the evidence underlying the three highest-net-score-adjacent problems specifically for vendor-originated bias, which the Empiricist flagged but did not resolve.
- **Why it matters:** "High confidence / high evidence quality" labels from Domain Research were self-declared and pre-council; this task tests whether they hold up.
- **Expected output:** A short annotated bibliography per problem separating primary/independent evidence from vendor-originated material.
- **Success criteria:** Each problem's evidence quality is re-labeled based on actual source independence, not the original self-assessment.
- **Decision it enables:** Whether the net scores for these three problems should be trusted as-is, or whether the Empiricist's evidence-quality axis needs to be revisited before further investment.
- **Estimated effort:** 2–4 days, desk research.

### Task 8 — ACV / unit-economics reality check (P9, P4, P2, P8, P10)
- **Objective:** Attach real numbers to the "small ACV," "long sales cycle," and "per-seat pricing tight" flags that the Economist raised across five of the ten problems without quantifying.
- **Why it matters:** Structural revenue-model risk was named repeatedly but never measured — it currently functions as a qualitative caution, not a decision input.
- **Expected output:** A comparable pricing/ACV estimate for each of the five flagged problems, based on public pricing of named incumbents plus interview data from Tasks 2 and 5.
- **Success criteria:** Each flagged problem has an order-of-magnitude ACV estimate attached, not just a qualitative risk note.
- **Decision it enables:** Whether any of these five problems is disqualified on unit economics alone, independent of technical merit.
- **Estimated effort:** 3–5 days, can run in parallel with Task 2.

---

## Closing note

Nothing in this brief should be read as narrowing the set. If anything, it widens the honest uncertainty: the two highest-net-score problems (P4, P2) both carry structural go-to-market risk that the R1 scoring didn't fully price in, while the problem with the single strongest positive council signal (P5) has the least founder-realism validation of any candidate. The research-readiness matrix and the eight tasks above are designed to close those specific gaps — not to build toward a particular answer.

<!-- checkpoint: fix(issuance): fix truststore backend -->
