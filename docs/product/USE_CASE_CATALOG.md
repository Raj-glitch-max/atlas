# Use Case Catalog

Each use case references the persona from `USER_MODEL.md` and the evidence source it derives from. Use cases marked hypothesis are not dropped, but are separated from evidenced ones so `V1_SCOPE.md` can exclude them without ambiguity.

## Evidenced use cases

**UC1 — Prove delegated identity across trust domains**
A delegate workload presents proof that it is acting on behalf of a principal to a relying party in an independently-operated trust domain.
Actors: Delegate workload, Relying-party service owner.
Evidence: `TECHNICAL_VALIDATION.md` (P5, item 1, core hypothesis; item 7, minimum experiment).

**UC2 — Scope and time-bound a delegation**
A delegation is issued with a permission scope narrower than the principal's own, and an expiration after which it is no longer valid.
Actors: Security Engineer (defines scope), Delegate workload (presents it).
Evidence: `PRODUCT_THESIS.md` (P5, item 3).

**UC3 — Revoke a specific delegation independently of the underlying identity**
A specific delegation is invalidated without affecting the principal's or delegate's underlying identity, and without requiring the delegate workload to be redeployed or restarted.
Actors: Security Engineer, Platform/Infrastructure Engineer.
Evidence: `PRODUCT_THESIS.md` (P5, item 2, "revocation is coarse" insufficiency; item 4, "failure containment").

**UC4 — Verify a delegation without a live broker call**
A relying party checks the validity, scope, and revocation status of a presented delegation without a network call to a shared authority at verification time.
Actors: Relying-party service owner.
Evidence: `TECHNICAL_VALIDATION.md` (P5, item 8, success criteria — offline/partition-tolerant verification).

**UC5 — Reconstruct a delegation chain after the fact**
A verifiable record is produced that is sufficient to determine which identity delegated to which, with what scope, and when.
Actors: Security Engineer, Auditor (hypothesis persona — record production is evidenced, the auditor's use of it is not).
Evidence: `PRODUCT_THESIS.md` (P5, item 4, structural audit trail).

**UC6 — Issue identity and delegation to an ephemeral, high-churn workload**
A short-lived agent or workload receives identity and a delegation without requiring long-lived, statically-provisioned identity.
Actors: Delegate workload, Platform/Infrastructure Engineer.
Evidence: `TECHNICAL_VALIDATION.md` (P5, item 2, ephemeral/agentic workload assumption).

## Hypothesis use cases (not evidenced, not excluded — see `DEFERRED.md`)

**UC7 — Bridge trust between two organizations with no shared federation setup**
Two independently-operated organizations establish delegation-verifiable trust without a pre-existing shared identity federation.
Evidence status: `PRODUCT_THESIS.md` explicitly offered this only "to test generativity, not proposed" (P5, item 4). Not a validated requirement.

**UC8 — Continuous re-attestation of a delegate's posture after issuance**
A delegation's validity is tied to the delegate's ongoing security posture, not only its state at issuance time.
Evidence status: `TECHNICAL_VALIDATION.md` classifies continuous re-attestation as **unproven** (P5, item 3). Included here only because it appears in the technical hypothesis space, not because it is a validated use case.
