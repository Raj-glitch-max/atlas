# 08 — PRODUCT DESIGNER — RAW FINDINGS

Phase 3 (Ops & Product). I treat the CLI/SDK/errors as UI. Sugar lives in the SDK/CLI layer, never
in `record/` or `verify/` (R7) — and I confirm the core honors that: the ergonomic surface is in
`sdk/` and `cmd/`, the invariant surface in `internal/`. Good separation.

FINDING-PD-1
  Type:        error-as-UI
  Claim:       The verdict trace is rich and machine-honest (five stages, causes, digests), but
               there is no human-facing `atlas explain` that turns "InconclusiveRejected:
               RevocationKnowledgeStale" into "This record is fine, but your revocation snapshot is
               older than your R=5m budget; refresh it or widen R." The raw cause is correct and
               unfriendly.
  Evidence:    `internal/verify/{cause,trace}.go` produce precise causes; no explanation layer maps
               them to operator action.
  Recommendation: LATER (Tier 6, S12) — `atlas explain <record>` reading the existing trace. Cheap
               because the trace already carries everything; it is presentation, not new logic.
  Owner:       08_product_designer

FINDING-PD-2
  Type:        two-channel design
  Claim:       One error surface serves both a trust boundary (the verdict) and a debugging session
               (why). These have different audiences. Today they are the same trace. For the holder
               this is safe (records are public by design, so there is no oracle to leak), but the
               *operator* explanation and the *machine* verdict should be visibly different channels.
  Recommendation: LATER — pair with PD-1: machine verdict stays terse; `explain` is the human channel.
  Owner:       08_product_designer

FINDING-PD-3
  Type:        first-run
  Claim:       The default wiring uses the Degenerate revocation provider (always Indeterminate →
               inconclusive-reject), so a naive first `verify` of a freshly issued record *rejects*
               unless the user wires a real/fresh revocation source. This is correct fail-closed
               behavior but a confusing first-run experience.
  Evidence:    `internal/revstatus/indeterminate.go` (Degenerate default); `examples/` wire a fresh
               source to get an Accept.
  Recommendation: NEXT-ish — the quickstart / `atlas doctor` should surface "you are on the
               degenerate revocation provider; here's why this rejects and how to get an Accept."
  Owner:       08_product_designer

VERDICT: The core/ergonomics boundary is respected and the trace is a strong foundation. The
product-design work is presentation on top of already-honest primitives — appropriately Tier 6,
except the first-run explanation (PD-3) which is cheap and prevents a "why does everything reject?"
bounce.
