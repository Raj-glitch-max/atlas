# 09 — UX RESEARCHER — RAW FINDINGS

Phase 4 (Polish). I own the three predicted security *misconceptions* — the ways a competent
integrator will misread Atlas and build an insecure system while believing they are safe. I have
no real comprehension data (no users — R11), so these are hypotheses to validate in PM-1's
interviews, filed as consequences, not as user claims I cannot support.

FINDING-UXR-1
  Misconception: "Verified means authorized."
  Why predicted: The word "verify" invites the reader to believe a passing verdict grants
                 permission. Atlas is careful here — the verdict is a *value* separating
                 cryptographically-valid from authorized (T10) — but the *docs* must make an
                 integrator do the second step explicitly or they will treat Accept as allow.
  Evidence:      `verify` returns a `Verdict`; the authorization decision is the integrator's. The
                 `--require-scope` gate and `examples/atlas-gate` model this correctly; a reader who
                 skips them may not.
  Severity:      inherits authorization-bypass severity if mis-integrated → treat as MAJOR docs.
  Validate:      Ask 3 interviewees to describe what they'd do after an Accept. If any says "allow
                 the action" without a scope check, the misconception is real.
  Owner:         09_ux_researcher

FINDING-UXR-2
  Misconception: "Inconclusive means try again / probably fine."
  Why predicted: `InconclusiveRejected` is carried `[HYPOTHESIS]` and reads softer than `Reject`.
                 An operator under pressure may treat stale-knowledge inconclusive as "good enough"
                 and widen R until it always passes — silently converting fail-closed into
                 fail-open by configuration.
  Evidence:      `policy.go` lets R be any `>0`; LIMITATIONS §6 warns "configurable ≈ unbounded."
  Severity:      MAJOR docs — this is 04's "an operator under pressure will raise it."
  Validate:      Watch whether interviewees reach for "raise R" as the fix for inconclusive.
  Owner:         09_ux_researcher

FINDING-UXR-3
  Misconception: "Offline verification means no setup."
  Why predicted: "No call to the issuer" can be misread as "no provisioning." In fact the RP must
                 hold trust material (out-of-band) and a reasonably fresh snapshot, or everything
                 fails closed (correctly, but surprisingly — see PD-3).
  Evidence:      `truststore` is empty until provisioned; Degenerate provider default rejects.
  Severity:      MODERATE docs / first-run.
  Validate:      First-run task: can a new user get one Accept without reading the source? Time it.
  Owner:         09_ux_researcher

## MANDATORY CLOSING
There is no real comprehension data yet, and that absence is the finding: the three misconceptions
above are the *specific* things PM-1's interviews must probe. Each, if real, inherits the severity
of the security mistake it causes — not the severity of a typo. The docs that pre-empt them
(especially UXR-1 and UXR-2) are the highest-value documentation in the project.
