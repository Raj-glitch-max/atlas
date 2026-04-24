# Ecosystem Thesis

Scope: ecosystem analysis only — not market, not competition, not implementation. The governing question for each candidate primitive: does it have gravity, meaning does its existence pull other independent systems toward depending on it, the way SPIFFE, OpenTelemetry, and OCI images do today.

One framing note, carried consistently from the Product Thesis: neither primitive currently exists as a built, adopted thing. Questions 4–9 below are answered conditionally — "if built as described, what would naturally connect to it" — because that's how they're phrased in the brief. Question 10 and the final verdict are answered honestly about *today's* state, because gravity is a property of realized adoption, not of potential. Conflating the two would inflate something unbuilt, which this project has consistently avoided doing.

---

## P5 — Standardized binding of workload identity to attenuable delegation capabilities

### 1. Existing standards it extends
SPIFFE/SPIRE (CNCF-graduated workload identity), X.509, RFC 8693 (OAuth 2.0 Token Exchange), and the Biscuit token specification. Notably, it would also extend an *active, live* standards conversation rather than starting a new one: four IETF drafts on agent/non-human identity appeared in early 2026, meaning there is a real, ongoing venue this could feed into.

### 2. Existing standards it conflicts with
SPIFFE's own design philosophy is a real conflict, not a minor friction point: SPIFFE maintainers have deliberately scoped the project to authentication only, explicitly excluding authorization/capability semantics. A delegation-capability layer pushes directly into territory SPIFFE chose to stay out of. Separately, it sits in tension with the plain OAuth bearer-token model, which assumes centrally-issued, coarse-grained scopes rather than independently attenuable, offline-verifiable ones.

### 3. Existing standards it could eventually influence
The 2026 IETF agent-identity drafts directly; a possible SPIFFE companion spec (rather than a SPIFFE core change, given point 2); OAuth's Rich Authorization Requests (RFC 9396) work, which is already moving toward finer-grained scope expression; potentially a W3C Verifiable Credentials delegation profile, since VC's proof-chain model overlaps conceptually.

### 4. Existing open-source projects that would naturally integrate it
SPIRE itself (as a plugin, not a core change); Open Policy Agent / OPAL for evaluating capability scope at enforcement points; Envoy (as an ext_authz filter); Istio and Linkerd (mesh-level enforcement); cert-manager (issuance tooling).

### 5. Existing infrastructure products that would naturally consume it
API gateways (Kong, Envoy-based gateways, Apigee); service meshes (Istio, Linkerd, Cilium); secrets managers (HashiCorp Vault, as a bridge to dynamic, capability-scoped secrets); cloud IAM interop points (AWS IAM Roles Anywhere, GCP Workload Identity Federation) as bridges rather than replacements; CI/CD identity issuance (GitHub Actions OIDC tokens).

### 6. Existing developer tools that would naturally expose it
A capability-token inspector (the equivalent of jwt.io for this format); kubectl plugins for issuing/inspecting delegation chains; Terraform providers for policy-as-code issuance; audit/visualization tooling for delegation chains, distinct from generic log viewers.

### 7. Which layer of the cloud-native stack it belongs to
The security and identity layer — a horizontal, cross-cutting layer adjacent to service mesh, not a vertical layer of its own.

### 8. Layer classification
**Protocol layer**, with a credible but unearned-until-adopted path toward standards layer. It is not application layer (it has no UI or user-facing feature), not platform layer (it doesn't orchestrate compute), and not yet standards layer (no standards body has ratified anything here).

### 9. Could another company realistically build on top of this primitive?
Yes, and concretely: an API gateway vendor could ship "verify delegation chain" as a native ext_authz capability, the way many gateways ship native OIDC verification today. A startup could build an "agent identity broker" the way Auth0 and Okta built businesses directly on top of the OAuth/OIDC standards rather than inventing their own auth model. A cloud provider could adopt it as a cross-cloud trust bridge, since none currently offer one.

### 10. If this primitive disappeared tomorrow (assuming it had been built and adopted as described), what would notice?
Agent-orchestration frameworks that had come to rely on scoped, revocable delegation instead of blanket credentials would regress to broader, less safe credential-sharing patterns. Cross-organizational integrations that had stopped requiring shared federation setups would need to rebuild that shared infrastructure. Audit tooling relying on structural proof would lose that and fall back to log reconstruction.

**Today, honestly: nothing would notice, because it doesn't exist yet.** The above is what the analysis suggests *would* be true if it existed and were adopted — that's a statement about structural potential, not present reality.

---

## P10 — Standardized, portable inferred causal edge for non-call-graph-connected causation

### 1. Existing standards it extends
OpenTelemetry — specifically its existing multi-signal architecture (traces, metrics, logs, and the newer continuous-profiling signal added 2024–2025), which this could be pitched as extending with a fifth signal type. Also extends W3C Trace Context conceptually, though not directly, since Trace Context governs propagation, not causal inference.

### 2. Existing standards it conflicts with
This conflict is more fundamental than P5's. OpenTelemetry has deliberately kept its standard scoped to vendor-neutral, deterministic *collection* — a span records what factually happened, not an inference about why. Every existing OTel signal type is a factual record; a causal edge, by construction, is a probabilistic claim with a confidence score. Embedding inference into a collection standard is a philosophical mismatch with OTel's own stated design principle, not a compatibility detail to be smoothed over.

### 3. Existing standards it could eventually influence
Future OTel semantic conventions specific to causality or AIOps, if the community ever moves in that direction; emerging CNCF conversations at the intersection of AI and observability; academic benchmark standards such as RCAEval, which could adopt a canonical exchange format for causal claims to make cross-method comparison easier than it is today.

### 4. Existing open-source projects that would naturally integrate it
The OpenTelemetry Collector (as a new processor type, if the philosophical conflict in point 2 were resolved); Grafana's Tempo/Loki/Mimir stack and Jaeger for visualization; eBPF-based tools (Pixie, Cilium/Hubble) as a data source specifically for the cross-cutting resource-contention signals that call-graph-based traces structurally can't see.

### 5. Existing infrastructure products that would naturally consume it
APM vendors (Datadog, Honeycomb, New Relic, Grafana Cloud) as a new ingestible signal type; incident-management platforms (PagerDuty, incident.io) for auto-attaching a probable cause to an alert; chaos-engineering tools (Gremlin, LitmusChaos) for validating injected-fault ground truth against inferred edges.

### 6. Existing developer tools that would naturally expose it
Debugger/IDE integrations surfacing causal chains inline; a CLI for "what caused this alert"-style queries; a dedicated causal-graph visualization, distinct from a standard trace waterfall view.

### 7. Which layer of the cloud-native stack it belongs to
The observability layer — the same horizontal layer OpenTelemetry itself occupies.

### 8. Layer classification
**Protocol layer** for the data-interchange format itself (an object other tools could exchange), but the engine that *populates* that format is platform-layer work, not protocol-layer — and that engine is exactly the piece the Technical Validation already flagged as scientifically unresolved. This is a meaningfully different situation than P5, where the protocol-layer object and its populating mechanism are both already proven independently.

### 9. Could another company realistically build on top of this primitive?
Conditionally yes, but the condition is heavier than P5's. An APM vendor could ingest a standardized causal-edge format and build UI/alerting on top of it, the way many vendors built trace explorers on top of standardized OTel trace data. A startup could build a dedicated causal-graph explorer the way Jaeger and Zipkin built trace explorers. But an empty, unreliable standard gives a second company nothing to build on — unlike P5, where the underlying components (SPIFFE, Biscuit) already work regardless of whether anyone adopts the combined standard, P10's underlying inference capability doesn't yet reliably exist at all. A second company can build on a format; it cannot build on a format with no trustworthy producer behind it.

### 10. If this primitive disappeared tomorrow (assuming it had been built and adopted as described), what would notice?
OTel-adjacent tooling and APM vendors that had integrated the signal, and any incident-management platform relying on its confidence scores for alert prioritization, would regress to manual or heuristic correlation.

**Today, honestly: nothing would notice, because it doesn't exist, and — unlike P5 — its non-existence isn't purely an adoption question. The thing that would need to exist first (a reliable populating mechanism) hasn't been demonstrated at all, per the Technical Validation.**

---

## Ecosystem Gravity — Verdicts

### P5 — Verdict: **B, with a structurally credible path to C**
Today, realized gravity is **A (none)** — it doesn't exist yet. Structurally, though, this is shaped like a strong platform primitive: its component parts (SPIFFE, Biscuit/macaroons) already function independently and don't need a scientific breakthrough to work; it would extend a live, currently-forming standards conversation (the 2026 IETF agent-identity drafts) rather than having to create one; and its natural integration points (SPIRE, OPA, Envoy, Istio) are active, healthy OSS projects, not hypothetical ones. What caps it below a confident **C** today is a real, named conflict — SPIFFE's deliberate exclusion of authorization semantics — and a fragmentation risk, since at least one well-funded vendor (Teleport) is already building a proprietary answer to the same problem, which could pull adoption toward a closed solution before an open standard consolidates the space.

### P10 — Verdict: **A today, capped at B even conditionally**
Today, realized gravity is **A (none)** — same as P5, and for the same basic reason. But unlike P5, the ceiling doesn't move to C even in the best-case conditional analysis, for two independent reasons: OpenTelemetry's own standard has a real philosophical objection to embedding probabilistic inference into a collection format, not just a technical integration gap; and the underlying inference capability that would populate the primitive is, per the Technical Validation, an open research question, not an engineering-execution one. Gravity requires other systems to trust and depend on a signal — nothing can credibly depend on a signal whose reliability at production scale hasn't been demonstrated anywhere, by anyone. **C and D are not reachable through ecosystem or standards work alone; they're gated by the same unresolved feasibility question the Technical Validation already identified**, and no amount of standards strategy resolves that.
