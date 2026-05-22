# Infrastructure Primitive Evaluation — Delegated Workload/Agent Identity

Question under evaluation, exactly as posed: if this project succeeds perfectly, what new computing primitive has entered the cloud-native ecosystem? This document does not assume the answer is "yes, a new primitive." It tests that claim against seven primitives that are uncontested members of that category, and follows the evidence wherever it leads.

The candidate object under test, carried forward from `PRODUCT_THESIS.md` (P5, item 3): a standardized binding of workload identity to an attenuable, offline-verifiable delegation capability, portable across independently-operated trust domains.

---

## Comparison 1 — Git Commit

**What made it fundamental:** Content-addressing (the object's identity *is* a hash of its content, so identity and integrity are the same fact), full composability into a DAG, and — critically — **unilateral usefulness**. A commit has value the instant it's created, on a single machine, with no network, no counterparty, and no coordination with anyone else. Value comes before adoption, not after it.

**Does this project share that characteristic:** Partially, and only on the integrity axis. The delegation object aims for tamper-evidence and composability into a chain, which rhymes with a commit graph. But it fails the unilateral-usefulness test outright: a delegation proof has zero value until a second, independent party verifies it. Nothing in this package describes a mode where the object is useful alone.

**Where it falls short:** Git's primitive spreads virally because early adopters get value before anyone else adopts anything. This candidate requires simultaneous or sequential buy-in from a counterparty before it does anything at all — a structurally harder adoption path, closer to how slowly two-sided protocols like DNSSEC spread than how fast a one-sided tool like Git did.

**What would need to become true:** A discovered use case where the delegation object has standalone value — for example, local policy simulation or testing — before any relying party exists to verify it. Nothing in the current documents suggests this exists.

---

## Comparison 2 — Docker Container

**What made it fundamental:** It collapsed three previously separate hard problems (packaging, distribution, execution) into one artifact, and it was demoable in under a minute with zero specialist knowledge — `docker run` produces visceral, immediate proof of value.

**Does this project share that characteristic:** No. The entire validated scope of this project, per `V1_SCOPE.md`, is a two-trust-domain lab experiment with a cryptographic verification property as its success criterion. That is a meaningful engineering result; it is not a demoable "wow" moment for a general audience, and its intended audience (security/platform engineers) is inherently narrower than "every developer," which was Docker's audience.

**Where it falls short:** No trivial, visceral, one-command demonstration of value exists or is planned in this package.

**What would need to become true:** An equivalently simple demonstration — something like proving delegated identity across two unrelated systems in one command, legible to someone with no protocol background — would need to exist and actually surprise people. Nothing in the current scope produces this.

---

## Comparison 3 — Kubernetes Pod

**What made it fundamental:** It gave a durable name to an abstraction operators already needed (the smallest co-scheduled deployment unit), and it rode on top of Kubernetes' own already-massive adoption gravity to reach ubiquity almost immediately.

**Does this project share that characteristic:** Partially, in concept — "a delegation" could become a similarly nameable unit engineers reason about. But it has no equivalent platform to ride on. SPIFFE, the closest available host platform, has real but meaningfully smaller adoption gravity than Kubernetes had by the time Pod became vocabulary, and this package has no defined path for riding on SPIFFE's gravity — only a stated intent to coexist with it (`CONSTRAINTS.md`, C1–C2).

**Where it falls short:** No distribution vehicle with Pod-level gravity exists to carry this vocabulary anywhere.

**What would need to become true:** Either SPIFFE would need to reach Kubernetes-era ubiquity and formally adopt this as a companion vocabulary, or some other already-ubiquitous platform would need to adopt it first. Neither condition is close to being met.

---

## Comparison 4 — OpenTelemetry Trace

**What made it fundamental:** It didn't invent tracing — it ended vendor lock-in on tracing data by giving competing vendors a shared, neutral format, because every vendor was paying the same cost (maintaining N proprietary agents) with zero differentiation benefit from doing so. The economics pushed competitors toward convergence.

**Does this project share that characteristic:** This is the closest structural analogy in the whole set — both projects standardize something that already exists in silos. But the economic incentive runs in the opposite direction. Telemetry collection was never a competitive differentiator for APM vendors; identity and delegation schemes plausibly are. At least one funded vendor (Teleport) has already shipped a proprietary answer rather than converging on an open one — evidence pointing toward fragmentation, not the convergence that made OTel possible.

**Where it falls short:** No evidence of the specific economic pressure that made OTel's convergence rational for competing vendors.

**What would need to become true:** At least two unrelated, competing vendors would need to publicly converge on a shared specification the way OpenTracing and OpenCensus merged into OpenTelemetry — a sign that maintaining separate schemes has become pure cost, not competitive advantage. No such signal currently exists.

---

## Comparison 5 — Terraform Desired State

**What made it fundamental:** It's a different *kind* of primitive than the others in this list — not a new object, but a new control paradigm (declarative reconciliation instead of imperative operation). Its vocabulary won because it changed how people think about operations, not just what format they use.

**Does this project share that characteristic:** No, and it was never attempting to. This candidate is squarely a "noun" (a new kind of object), not a "verb" or control loop. This comparison mostly serves to clarify category, not to find a gap — measuring this candidate against a paradigm-shift bar would be a category error in either direction.

**Where it falls short:** Not applicable in the way the other six comparisons are.

**What would need to become true:** Nothing — this axis doesn't apply, and pretending otherwise would manufacture a gap that isn't real.

---

## Comparison 6 — Vault Secret

**What made it fundamental:** It turned "the secret" into a dynamically-issued, leased, revocable, audit-logged object, replacing static credentials scattered in config files — the closest prior-art match to what this project is attempting.

**Does this project share that characteristic:** Uncomfortably closely, and this is the most aggressive challenge in the whole comparison set. Vault already solved leasing, TTL, and revocation for credentials generally. SPIFFE already solved workload identity. Macaroons and Biscuit tokens already solved attenuable capabilities. Stress-tested honestly: this project's contribution may be closer to *composing three already-solved problems into one coherent cross-domain story* than to inventing something at the level Vault's dynamic secret represented when it was genuinely new in 2015.

**Where it falls short:** Vault won distribution through a single vendor (HashiCorp) shipping a widely-adopted product and getting instant reach through that one channel. This project has no equivalent distribution vehicle and, per `ECOSYSTEM_THESIS.md`, must succeed through open multi-party standardization instead — a structurally harder path than Vault's unilateral rollout.

**What would need to become true:** Either a single already-adopted product would need to ship this as a built-in feature the way Vault did, or the multi-vendor convergence described in the OpenTelemetry comparison would need to happen. Absent either, this risks being seen as a re-composition of existing, already-successful primitives rather than a new one.

---

## Comparison 7 — Kafka Log

**What made it fundamental:** Generality far beyond its original design intent. The append-only log turned out to underlie pub/sub, event sourcing, stream processing, and change-data-capture — patterns nobody was designing for when the log was first built. This unplanned generality is arguably the single strongest test of genuine primitive status of all seven examples.

**Does this project share that characteristic:** This is where the candidate fails most clearly. A delegation capability is narrowly and deliberately scoped to identity and authorization. Every use case in this project's own `USE_CASE_CATALOG.md` sits inside that domain. Nothing in `PRODUCT_THESIS.md`'s own "what becomes possible" list (P5, item 4) gestures at an unrelated use the way the log's applications did. There is no evidence, anywhere in this project's record, of generality beyond designed purpose.

**Where it falls short:** Zero evidence of use outside its intended domain — the sharpest, most objective gap in this entire evaluation.

**What would need to become true:** Someone would need to find a genuinely unrelated use for the delegation-capability object that nobody anticipated at design time. Its absence right now isn't neutral — it's informative, and it currently points away from primitive status.

---

## Verdict

Feature, product, platform, protocol, standard, or genuine infrastructure primitive?

**At best, if this project succeeds perfectly: a protocol, with a realistic path to becoming a narrow, specialist standard — not a genuine ecosystem primitive in the sense the other seven examples represent.**

The evidence supports this conclusion directly, not by default:

- It fails the two tests that most reliably separate a primitive from a well-executed standard: unilateral usefulness before coordination (Git Commit) and generality beyond designed purpose (Kafka Log). Both failures are total, not partial — no evidence anywhere in this project's record points the other way on either axis.
- It fails the two tests about *how* a vocabulary actually spreads: a trivial, visceral adoption moment (Docker) and riding on pre-existing platform gravity (Kubernetes Pod). Neither exists nor is planned for.
- Its closest analogy (Vault Secret) raises, rather than resolves, the concern that this is a recomposition of already-successful primitives (Vault's leasing model, SPIFFE's identity, macaroon-style attenuation) rather than a new one — the same tension `PRODUCT_THESIS.md` flagged when it distinguished "a new primitive" from "a new load-bearing combination of existing ones."
- Its most structurally favorable analogy (OpenTelemetry) depends on an economic convergence among competing vendors that has no supporting evidence yet, and one data point (Teleport's proprietary answer) currently points the opposite direction.

This is not the same as saying the project has no value. A narrow, well-adopted standard for cross-domain workload delegation — something in the reference class of SAML, WebAuthn, or mTLS: real, important, standardized, but understood by a specialist audience rather than reshaping how engineers in general talk about systems — is a legitimate, non-trivial outcome and is not foreclosed by anything in this evaluation. It is a smaller claim than "the next Docker or the next Git," and the evidence currently available does not support the larger claim.

## What evidence is still missing

1. **A demonstrated case of unilateral value** — the delegation object doing something useful before any counterparty exists to verify it. None currently exists.
2. **A demonstrated case of use outside the identity/authorization domain** — an application nobody designed for. None currently exists; this is the single most telling gap, because it's the one that would be hardest to manufacture artificially and most convincing if it appeared naturally.
3. **Evidence of vendor economic convergence rather than fragmentation** — at least two independent, competing parties converging on a shared spec rather than shipping proprietary answers. The one data point available (Teleport) currently points the other way.
4. **A trivial, legible demonstration of value** to an audience broader than security/platform specialists. Nothing in the validated scope produces this, and nothing in `V1_SCOPE.md` is designed to.
5. **A clear differentiation from Vault's existing dynamic-secret model** sufficient to establish this as something new rather than a recombination of Vault, SPIFFE, and macaroons under a different name.
6. **Adoption riding on an already-ubiquitous host platform**, rather than a standalone standardization effort competing for attention on its own.

None of these six gaps are closed by anything in this project's prior documents. This evaluation does not conclude the project should or shouldn't proceed — only that the specific claim "this becomes a new cloud-native primitive" is not yet supported, and the evidence that would support it is identifiable, specific, and currently absent.

<!-- checkpoint: chore(fuzz): improve Fuzz Verification core target -->
