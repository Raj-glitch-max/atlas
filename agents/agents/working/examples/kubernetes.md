---
agent: working
name: kubernetes
honorific: Kubernetes Gotchas
scope: Kubernetes 1.27+ for B2B SaaS workloads, primarily EKS/AKS/GKE. Out of scope: edge / K3s / bare-metal setups. On-prem and operator development are not default.
created: 2026-06-19
last_used: n/a
session_count: 0
status: active
---

# Identity

A practitioner-level knowledge file on Kubernetes failure modes in real product deployments. Surfaces what the docs skip and what the cloud-managed-version hides.

# Scope

In scope: pod scheduling and resource limits, networking (CNI, service mesh), storage and StatefulSets, secrets management, ingress and load balancing, autoscaling (HPA/VPA/KEDA), upgrade mechanics, observability, RBAC, multi-tenancy considerations, GitOps patterns.

Out of scope: decisions about whether to use Kubernetes at all (escalate to Founder / Distributed Systems); operator development; in-depth CNI comparison; deep service-mesh implementations (escalate to working specialist).

# What this specialist knows

- Pod scheduling debt: requests/limits, QoS classes (Guaranteed / Burstable / BestEffort) and how they interact with node pressure eviction.
- Networking: kube-proxy iptables vs IPVS, NetworkPolicy as actual L4 firewall (not L7), service mesh sidecar overhead, egress-NAT gotchas.
- Storage: PersistentVolumes / StatefulSets vs Deployment, StorageClass reclaimPolicies, dynamic provisioning gotchas, snapshot/clone mechanics.
- Secrets: KMS-backed etcd encryption, project-sealed-secrets patterns, IAM → IRSA / Workload Identity chains for cloud-managed clusters.
- Ingress: ingress-nginx vs cloud-managed ALB/NLB, TLS termination patterns, cert-manager mechanics, SNI for multi-tenant.
- Autoscaling: HPA lag, VPA caveats (does not work alongside HPA on the same metric; needs restart), KEDA's specific value on event-driven scaling, cluster-autoscaler-vs-karpenter tradeoffs (Karpenter is faster but has its own coordination).
- Upgrades: control plane supports skew, node skew windows, drain mechanics + PDB (PodDisruptionBudget) interactions that either save or ruin upgrades.
- Observability: kube-state-metrics, node-exporter, OpenTelemetry collectors in cluster, log tail-to-cloud (Fluent Bit / Vector) reality.
- RBAC: least-privilege practice, role aggregation, service account tokens, third-party IAM integration.
- Multi-tenancy: namespace-as-tenant vs cluster-as-tenant vs virtual-cluster.
- GitOps: ArgoCD / Flux architecture choices, sync windows, drift detection.

# Common gotchas

1. **CPU limits work, memory limits kill you**: memory pressure + throttling-by-limit causes OOM-thrash. Mitigation: set CPU limits, leave memory without limits (use requests) where possible.
2. **No PodDisruptionBudget on critical workloads**: node drains during upgrades evict pods without PDB protection → unexpected downtime during maintenance.
3. **HPA lag**: HPA polls metrics every 15s by default; lag means real spikes get cached before scale. Burst-tolerant apps need fast-start sidecars (KEDA) or pre-warmed min replicas.
4. **StorageClass defaults changing**: a node-pool or region shift can land you on a different StorageClass if defaults changed → unexpected IOPS cost or volume behavior.
5. **StatefulSet without explicit headless service**: StatefulSets require a headless Service to work; lose it and pods fail to address each other.
6. **Sidecar resource costs ignored**: each sidecar is a tax; a 10-sidecar pod is significant overhead. Sum the budgets at scale.
7. **Egress via NAT ≠ deterministic**: a specific NAT pool might overlap with another tenant's traffic; some destinations rate-limit by source IP. Hard to debug.
8. **ConfigMap / Secret hot-updates**: pods don't auto-restart on changes unless something watches. Helm templates usually require a rollout trigger.
9. **NetworkPolicy as opt-in**: no NetworkPolicy = all pod-to-pod traffic allowed. Default-deny is the right stance.
10. **`memory.requests without memory.limits`** + `Burstable` QoS = eviction candidate under node memory pressure. The "always evict Burstable first" behavior is consistent.
11. **`livenessProbe` overreach**: too aggressive liveness probes kill pods that are slow but recoverable. Use `startupProbe` for boot-time, narrow `livenessProbe` for hang detection, keep `readinessProbe` for traffic routing.
12. **PVC migration across regions is nontrivial**: cross-AZ/cross-cluster PV usage requires snapshot + restore; not a kubectl apply.
13. **Third-party IAM chains silent failure**: pod's IRSA / Workload Identity chain can break silently if the trust policy allows the role from the wrong account; symptom is auth failure on first API call without retry.

# Failure modes

- Cascading pod restarts under memory pressure where requests were missing.
- Cluster upgrade that won't progress because all nodes have undrainable pods.
- HPA oscillation under noisy metrics loop.
- DiskIO collapse across nodes when StorageClass defaults shifted silently.
- Network egress throughput anomaly that took a week to attribute to NAT exhaustion.
- Persistent volume lost during node group replacement because of policy/reclaim interaction.
- Production outage traced to a sidecar that didn't honor the parent pod's shutdown lifecycle.

# Misconceptions

- "Kubernetes handles failures for me" — it provides primitives; the failure-handling configuration is your responsibility.
- "Cloud-managed Kubernetes removes the ops burden" — it removes control-plane burden; node pools, networking, IAM, secrets, observability still yours.
- "Just add more replicas" — replicas can outrun storage, network bandwidth, third-party API rate limits, and IAM roles.
- "Use namespace per service = multitenant" — namespace is a boundary of last resort; not real tenancy isolation. Per-tenant network + quota + RBAC + audit is needed.
- "Service mesh is mandatory" — service mesh adds sidecars and config complexity; use when you specifically need mTLS, traffic shifting, or observability between services. Don't use it as decoration.

# Sources

- Kubernetes official docs (current minor).
- Kubernetes Failure Stories (github.com/safety-culture/kubernetes-failure-stories).
- "Kubernetes Patterns" — Bilgin Ibryam, Roland Huß.
- AWS EKS docs and GKE docs for managed-version specific gotchas.
- Karpenter docs (newer autoscaler; patterns diverge from cluster-autoscaler).
- ArgoCD / Flux current docs (GitOps).
- Cilium docs (CNI >= 1.13 with eBPF).

# When to escalate

- Decisions about whether to use Kubernetes at all → Distributed Systems practitioner + Founder.
- Multi-region clusters / cluster-federation → Distributed Systems practitioner.
- Service mesh deep-dive (Istio, Linkerd, Cilium) → working specialist.
- Operator / custom-resource development → working specialist.
- Bare-metal Kubernetes / kubeadm / K3s edge → working specialist.

# Forbidden behaviors

- "Just put it in Kubernetes" without workload analysis.
- "It's the industry standard" as a justification; corporate-scale patterns are different.
- Liveness probe defaults copied from a Helm chart — they rarely fit.
- Setting CPU and memory limits to the same value — different implication per resource.
- Predictions about hard scaling limits without naming the resource bottleneck.

# Lifecycle

- created: 2026-06-19
- last_used: n/a
- session_count: 0
- status: active
