# START HERE — Understanding Atlas from Zero

You have never seen this project and know nothing about it. This document is the
single, ordered path from *"what problem does this even solve?"* to *"I can read
the verification core and know why every line is there."* Follow it top to
bottom. Each stage says **what to read, in what order, and the one thing to take
away** — so you are never staring at a 99-file repo wondering where to start.

Budget: ~15 minutes to see it run, ~2 hours for a solid conceptual grasp, ~1 day
to genuinely understand the code.

---

## 0. The 60-second mental model

**The problem.** When software agent A needs agent B to do something on its
behalf, A usually hands B a shared secret or a bearer token — and to check that
token, B normally has to *phone home* to an authorization server. If that server
is unreachable (network partition, another org, an air-gapped edge), B is stuck:
fail open (insecure) or fail shut (unavailable).

**What Atlas is.** A tiny, signed, self-describing "delegation record" that says:

> *Principal `P` grants delegate `D` the scope `S`, valid until time `T`,
> revocable via instance `I`.*

Anyone holding `P`'s public key can verify that record **completely offline** —
no call to any server — in microseconds, and get a deterministic yes/no that
**fails closed** when it can't prove freshness.

**The gap it fills** (it replaces none of these — it fills the hole they leave):

| System | Solves | Assumes |
|---|---|---|
| OAuth / OIDC | user → app authorization | a live auth server at decision time |
| SPIFFE / SPIRE | workload identity | a live SPIRE agent / trust bundle |
| MCP | tool invocation | authorization is someone else's problem |
| **Atlas** | **verifiable delegation across agents** | **nothing — verification is offline, deterministic, fail-closed** |

**What it deliberately is NOT** (read this now so you don't misread it later):

- **Single-hop only.** `P → D`. The delegate cannot re-delegate (`A → B → C`
  chains do not exist). This is by design — see `LIMITATIONS.md` §1.
- **Two trust domains,** not a federation.
- **Not a policy engine.** It tells you a record is *cryptographically valid*;
  deciding whether that valid record is *authorized* to do a thing is your job.
- **Closest in design space to** Biscuit / UCAN — but for SPIFFE identity, and
  offline-first.

If you remember one sentence: **Atlas turns "can this delegation be trusted?"
from a network question into a math question.**

---

## 1. See it run first (15 minutes, before you read any theory)

Understanding lands faster once you've watched it work. From the repo root:

```sh
make ci                              # build + vet + tests + lints + frozen-doc + import rules
bash examples/unforgettable.sh       # the headline story, live
```

`examples/unforgettable.sh` is the whole project in one script: it issues a
delegation, verifies it **offline with the server dead**, revokes it, shows the
revocation bite, shows a stale-snapshot verifier **fail closed**, and shows a
tampered record **refused**. Watch that once and the rest of this document is
just "how each of those steps works."

Two more worth running:

```sh
go run ./cmd/atlas-server            # the HTTP API (http://127.0.0.1:8087)
go run ./cmd/atlas doctor            # confirm it's up
```

**Takeaway:** you have now seen offline verification, revocation, staleness
fail-closed, and tamper-refusal actually happen. Everything below explains them.

---

## 2. Prerequisites — the five concepts (skim, ~20 min)

You do not need to be an expert, but these five ideas recur constantly. If one is
unfamiliar, read a primer before Level 3.

1. **JWS / JWT** — a JSON payload with a cryptographic signature over it, encoded
   as `header.payload.signature` (base64url, dot-separated). Atlas records are a
   pinned, restricted kind of JWS. (RFC 7515.)
2. **ES256** — the one signature algorithm Atlas uses: ECDSA over the P-256 curve
   with SHA-256. "Pinned" = no other algorithm is ever accepted, which structurally
   kills the entire class of "algorithm confusion" attacks. (RFC 7518.)
3. **SPIFFE ID** — a workload identity that looks like `spiffe://domain/path`. The
   principal and delegate in a record are SPIFFE IDs. The part before the first `/`
   is the **trust domain**, which selects which public key verifies the record.
4. **Capability / attenuation** — a delegation grants a *subset* of what the
   grantor holds ("attenuation" = narrowing). Atlas enforces the subset **at
   issuance** (a delegate can never be granted more than its principal has).
5. **Fail-closed** — when the verifier cannot *prove* something is safe (e.g. it
   can't tell whether a record was revoked because its revocation data is too old),
   it says **no**, never "probably fine." This is the project's spine.

**Takeaway:** Atlas is "a pinned JWS over SPIFFE identities, granting an attenuated
capability, verified offline and fail-closed." Every prerequisite maps to one word
in that sentence.

---

## 3. The reading path

Six levels, in order. Levels 0–2 are concepts (read fast). Level 3 is the code
(read carefully). Levels 4–6 are proof, operation, and the "why."

### Level 0 — Orient (read fast, ~15 min)

| Read | Why | Take away |
|---|---|---|
| `README.md` | the map | what's built, what commands exist |
| `WHY.md` | motivation + when *not* to use it | the honest boundaries of the idea |
| `LIMITATIONS.md` | what Atlas does NOT do | single-hop, two-domain, replay window, unaudited crypto — read all 10 |
| `context/00_PROJECT_CONTEXT.md` | mission + pipeline | what the project is and is not |
| `context/06_GLOSSARY.md` | the vocabulary | this repo is acronym-dense; keep this open |

> This repo uses a lot of short codes: **M1–M6** (the six modules), **INV** (invariants),
> **FM** (failure modes), **SO** (security objectives), **ER** (engineering
> requirements), **AP** (architecture principles), **AD-###** (architecture
> decisions), **R** (the revocation-freshness bound). There's a quick glossary at
> the bottom of *this* file (§8), and the full one is `context/06_GLOSSARY.md`.

### Level 1 — The problem & the product (~20 min)

| Read | Take away |
|---|---|
| `docs/product/PITCH.md` | the one-paragraph "why anyone cares" |
| `docs/product/USE_CASE_CATALOG.md` | concrete scenarios the design must serve |
| `docs/product/USER_MODEL.md` + `docs/product/SYSTEM_CONTEXT.md` | who uses it and what's inside vs outside the boundary |
| `docs/product/FUNCTIONAL_REQUIREMENTS.md` (FR#) | the requirements the code traces back to |

**Takeaway:** you now know *what it must do* and *for whom*, before seeing *how*.

### Level 2 — The architecture (~30 min)

Read the RFCs in numerical order — they build on each other:

| Read | Take away |
|---|---|
| `rfc/RFC-000-architecture-principles.md` | the 8 principles (offline-first, fail-closed, minimal-trust, …) that every later decision obeys |
| `rfc/RFC-001-system-context.md` | the system boundary and trust model |
| `rfc/RFC-002-conceptual-domain-model.md` | the domain objects and the revocation state machine |
| `rfc/RFC-003-logical-software-architecture.md` | **the six modules M1–M6** and the dependency rules (R1–R7) between them |
| `INTERFACE_SPECIFICATION.md` | **the normative, binding contract** for every module: inputs, outputs, closed answer sets, failure semantics. This is the spec the code implements. |
| `SYSTEM_ARCHITECTURE.md` | the whole picture in one place |

**The six modules (memorize this map — the code is organized exactly like it):**

```
M1 record      internal/record      the signed artifact + integrity check
M2 issuance    internal/issuance    the ONLY thing that creates records (enforces attenuation)
M3 verify      internal/verify      the verification core — THE heart of the system
M4 truststore  internal/truststore  the relying party's locally-held public keys
M5 revstatus   internal/revstatus   "is this instance revoked, and how fresh is that knowledge?"
M6 revorigin   internal/revorigin   the authoritative revocation register (producer side)
```

**Takeaway:** M1 is the noun; M2 makes them; M3 judges them; M4/M5/M6 are the
inputs M3 consults. Verification (M3) depends on M4/M5 but **never** on M2 —
verifying must not need the signing code.

### Level 3 — The code, in dependency order (~2–4 hours) — the heart

Read the packages **bottom-up**, in this exact order. Each depends only on the
ones above it, so nothing is a forward reference. For each package the files are
listed in reading order with the single thing to grasp.

#### 3.1 `internal/record` (M1) — the artifact

| File | Grasp |
|---|---|
| `doc.go` | the package's job in prose |
| `record.go` | `Assertions` (what a record asserts) and `Record` (the sealed artifact); `Outcome` = `Intact`/`Altered`, zero value is `Altered` (fails safe) |
| `envelope.go` | the wire format: `alg`/`typ`/`kid` pinned, the claim layout, `decodeClaims` (runs only on verified bytes), and duplicate-key rejection |
| `seal.go` | `Seal()` — the *only* constructor of a record; requires the private key; refuses partial construction |
| `integrity.go` | `ValidateIntegrity()` — the 6-step pipeline (shape → parse → pin → resolve key → verify signature → decode). **Every exit is `Altered` except a fully-valid record.** Read this twice. |
| `trustmaterial.go` | `TrustMaterial` — public keys for one domain; curve enforced at construction |
| `peek.go` | `PeekTrustDomainUnverified` — reads the domain from *unverified* bytes only to pick which key to try; never a trust decision |
| `revbinding.go`, `instanceid.go` | the opaque revocation-binding slot and the per-issuance instance identity |

**Takeaway:** a `*Record` can only come into existence two ways — `Seal` (issuance)
or `ValidateIntegrity` returning `Intact`. There is no other door. That is the
whole tamper-evidence story.

#### 3.2 `internal/issuance` (M2) — the only maker of records

| File | Grasp |
|---|---|
| `doc.go`, `ports.go` | the injected inputs: a `PermissionSource`, a clock, a minter, a signer |
| `permissionset.go` | `isProperSupersetOf` — the **strict** subset check: a delegate's scope must be *strictly smaller* than the principal's permissions (this is attenuation, enforced here and nowhere else) |
| `authority.go` | `Issue()` — validate → consult permissions → check subset → seal. On any refusal, **nothing is created** |

**Takeaway:** attenuation (a delegate can't out-scope its principal) is an
**issuance-time** guarantee. The verifier does *not* re-derive it — it trusts
that the issuer's key signed only a valid subset. (This is why single-hop matters:
there's no chain to re-check.)

#### 3.3 `internal/truststore` (M4) — the relying party's keys

| File | Grasp |
|---|---|
| `store.go` | holds public keys per domain; **structurally cannot fetch** (imports no network); absent keys → an honest "absent" that the verifier fails closed on |

#### 3.4 `internal/revstatus` (M5) — revocation observation

| File | Grasp |
|---|---|
| `answer.go` | the closed answer set: `Indeterminate` / `NotObservedRevoked` / `ObservablyRevoked`, each with an `AsOf` freshness; zero value is `Indeterminate` (fails closed) |
| `indeterminate.go` | the **default** provider answers `Indeterminate` to everything (honest ignorance until a real mechanism is wired) — this is why a naive first verify *rejects* |
| `statuslist.go` | the real realization: a signed, timestamped snapshot of revoked instances; the signature covers `listID + asOf + set`; snapshots are adopted only if strictly newer (rollback-proof) |
| `contracttest/` | the contract every provider realization must pass |

**Takeaway:** the provider *reports* knowledge + freshness; it never decides
policy. Whether the knowledge is fresh *enough* is the verifier's call.

#### 3.5 `internal/revorigin` (M6) — the revocation producer

| File | Grasp |
|---|---|
| `register.go` | the authoritative "these instances are revoked" register that snapshots are published from |

#### 3.6 `internal/verify` (M3) — **the verification core, the single most important package**

Read `verifier.go` last and slowest — it is the locus where every guarantee
becomes true.

| File | Grasp |
|---|---|
| `policy.go` | `Policy{ R, skew }` — `R` is the max staleness of a "not revoked" answer the verifier will accept; a verifier can't be built without it |
| `ports.go` | the three injected inputs (trust material, revocation status, clock). **No I/O in this package** — offline is structural |
| `cause.go` | the closed set of rejection reasons, split into **definitive** (→ Reject) and **inconclusive** (→ InconclusiveRejected) |
| `verdict.go` | `Accept` / `Reject` / `InconclusiveRejected` |
| `check_integrity.go` | stage 1 (the gate): authenticate the bytes |
| `check_binding.go` | principal and delegate present and *distinct* |
| `check_expiry.go` | within the validity window, ± skew; future-dated → inconclusive |
| `check_scope.go` | scope present and well-formed (integrity, *not* re-derived subset — see 3.2) |
| `check_revocation.go` | apply `R` to the freshness of the revocation answer; stale → inconclusive |
| `verifier.go` | **`Verify()`** — runs the five checks, collects all causes, and routes: *any definitive cause → Reject; else any inconclusive → InconclusiveRejected; else Accept.* Emits a full decision trace **every time**, even on Accept. |

**Takeaway — the one idea to leave with:** verification is **five independent
checks + one order-independent routing rule**, and the routing is biased so that
*anything unproven becomes a rejection*. Forcing any single check to fail flips
the verdict away from Accept. Read `Verify()` and `route()` until that's obvious.

#### 3.7 The edges — `cmd/` and `sdk/` (skim)

The core above is pure and I/O-free. Everything user-facing wires it up:

| Read | Grasp |
|---|---|
| `cmd/atlas/` | the CLI (`issue`, `verify`, `--offline`, `doctor`, …) |
| `cmd/atlas-server/` | the HTTP API over the same engine |
| `cmd/atlas-mcp/` | Atlas exposed as agent tools over MCP |
| `sdk/go`, `sdk/python`, `sdk/typescript` | thin client SDKs mirroring the API |
| `examples/atlas-gate/` | a reverse proxy that admits a request only if it carries a valid capability |

**Takeaway:** all ergonomics live out here; the security-bearing logic lives only
in `internal/`. That boundary is enforced by a lint (dependency rules R1–R7).

### Level 4 — How it's proven (~45 min)

| Read | Take away |
|---|---|
| `THREAT_MODEL.md` | adversaries A1–A6, claims C1–C10, and **the test that proves each** — read this fully |
| `tests/vectors/VECTORS.md` | the language-neutral conformance spec: the exact wire format + verification algorithm a re-implementer in any language follows |
| `tests/vectors/verdict-vectors.json` + `negative-vectors.json` | 30 vectors (20 adversarial) — the machine-readable spec; the negative ones are attacks that MUST be rejected |
| `tests/conformance/properties_test.go` | property tests P1–P6 (tamper never accepts, garbage never accepts, fail-closed totality, determinism) |
| `internal/verify/fuzz_test.go` | coverage-guided fuzzing of the verifier |

**Takeaway:** the conformance vectors *are* the spec for anyone building a second
implementation. The negative vectors are where lenient verifiers get caught.

### Level 5 — Running & operating (~30 min, optional)

| Read | Take away |
|---|---|
| `examples/unforgettable.sh` | re-read it now that you know the pieces |
| `docs/runbooks/V1_OPERATION.md` | how to actually operate it |
| `deploy/` | the hardened container (distroless, nonroot, read-only rootfs) |
| `Makefile` (`make help`) | every build/test/lint target |

### Level 6 — The "why" behind every hard call (~1 hour, optional but high-value)

This is what makes Atlas unusual: the judgment is written down.

| Read | Take away |
|---|---|
| `ENGINEERING_DECISION_RECORD.md` | **AD-001 … AD-025**, each with the alternatives considered, the reason, and the trade-off accepted. This is where "why JWS not COSE?", "why ES256?", "why fail-closed?" are answered. |
| `agents/journal/` | dated decision memory — the freeze acts, the spike scope, the discoveries |
| `context/01_GOVERNANCE.md` + `GOVERNANCE.md` | how changes are governed (and the frozen-planning rule) |
| `context/07_SECURITY_POLICY.md` | the security objectives and honest limits |
| `agents/` | the reasoning framework the project is built with |

**Takeaway:** if you ever wonder "why did they do it *this* way?", the answer is
almost always already written — in an AD, a journal entry, or a `doc.go` comment.

---

## 4. How a single `Verify` call flows (the whole system in one trace)

Once you've read Level 3, this should read like a summary of everything:

```
presented bytes ──▶ verify.Verify()
   │
   ├─ check_integrity ─▶ record.ValidateIntegrity(bytes, trustMaterial[domain])
   │        │              shape → parse → pin(alg/typ/kid) → resolve key → verify sig → decode
   │        └─ not authentic? ▶ definitive IntegrityFailed ▶ REJECT   (gate: stops here)
   │
   ├─ check_binding     ▶ principal ≠ delegate?
   ├─ check_expiry      ▶ now within [iat, exp] ± skew?
   ├─ check_scope       ▶ scope present & well-formed?
   ├─ check_revocation  ▶ revstatus answer + is it fresher than R?
   │
   └─ route(all causes):
         any DEFINITIVE cause  ▶ REJECT
         else any INCONCLUSIVE ▶ INCONCLUSIVE-REJECTED   (fail-closed)
         else                  ▶ ACCEPT
      (+ a full decision trace, emitted every time)
```

Trust material (M4), the revocation answer (M5), and the clock are all **injected**
— the core reads no files and opens no sockets. That is why "offline" is a
compile-time fact, not a promise.

---

## 5. Reading paths by goal (shortcuts)

- **"I just want to evaluate the security."** → `THREAT_MODEL.md` →
  `LIMITATIONS.md` → `internal/record/integrity.go` → `internal/verify/verifier.go`
  → `agents/review-team/output/atlas-review-report.md` (an adversarial review of
  this exact code).
- **"I want to implement a verifier in another language."** → `tests/vectors/VECTORS.md`
  → `INTERFACE_SPECIFICATION.md` → run your impl against `verdict-vectors.json` +
  `negative-vectors.json` until all 30 pass.
- **"I want to contribute."** → `CONTRIBUTING.md` → `GOVERNANCE.md` →
  `DEVELOPMENT_RULES.md` → pick a package's `doc.go`.
- **"I want the business case."** → `docs/product/PITCH.md` → `WHY.md` →
  `docs/product/USE_CASE_CATALOG.md`.
- **"Why is X the way it is?"** → search `ENGINEERING_DECISION_RECORD.md` and
  `agents/journal/` before assuming it was an accident.

---

## 6. Three things newcomers get wrong (avoid these)

1. **"Verified means authorized."** No. `Accept` means *cryptographically valid*.
   Whether that valid capability is *allowed* to do the action is a separate,
   explicit decision you make. (See the `--require-scope` gate and
   `examples/atlas-gate`.)
2. **"Inconclusive means probably fine, retry."** No. `InconclusiveRejected` is a
   **rejection** — it means the verifier could not prove safety (e.g. revocation
   knowledge older than `R`). Widening `R` to make it pass is turning fail-closed
   into fail-open by configuration. Don't.
3. **"Offline means no setup."** No. The relying party must be *provisioned*
   out-of-band with trust material (M4) and a reasonably fresh revocation snapshot
   (M5). Offline means "no call *at decision time*," not "no configuration ever."
   With nothing provisioned, everything correctly fails closed (which surprises
   people — see the degenerate revocation default in `internal/revstatus/indeterminate.go`).

---

## 7. If you only read five files

In priority order, to understand the *essence*:

1. `internal/verify/verifier.go` — the five checks and the fail-closed routing.
2. `internal/record/integrity.go` — how a record is proven authentic (or not).
3. `THREAT_MODEL.md` — what's defended and how each claim is tested.
4. `LIMITATIONS.md` — what's honestly *not* defended.
5. `INTERFACE_SPECIFICATION.md` — the binding contract the whole thing implements.

---

## 8. Mini-glossary of the codes you'll hit everywhere

| Code | Means |
|---|---|
| **M1–M6** | the six modules: record, issuance, verify, truststore, revstatus, revorigin |
| **AD-###** | an Architecture Decision (in `ENGINEERING_DECISION_RECORD.md`) |
| **INV#** | a System Invariant (`docs/engineering/03_SYSTEM_INVARIANTS.md`) |
| **FM#** | a Failure Mode (`docs/engineering/04_FAILURE_MODEL.md`) |
| **SO#** | a Security Objective (`docs/engineering/02_SECURITY_OBJECTIVES.md`) |
| **ER#** | an Engineering Requirement (`docs/engineering/01_ENGINEERING_REQUIREMENTS.md`) |
| **FR#** | a Functional Requirement (`docs/product/FUNCTIONAL_REQUIREMENTS.md`) |
| **AP#** | an Architecture Principle (`rfc/RFC-000`) |
| **R** | the revocation-freshness bound: max staleness of a "not revoked" answer the verifier accepts |
| **skew** | the clock-skew tolerance applied to expiry |
| **kid** | key id — the JWS header field naming which public key verifies the record |
| **`atl_ins`** | the per-issuance instance identity (the revocation target) |
| **`atl_rvb`** | the opaque, optional revocation-binding slot |
| **attenuation** | narrowing a capability: a delegate gets a strict subset of the principal's scope |
| **fail-closed** | when safety can't be proven, reject — never "probably fine" |

The full glossary is `context/06_GLOSSARY.md`.

---

*You're oriented. Start at Level 0, run the demo in §1, and read down. When a
line of code seems arbitrary, check `ENGINEERING_DECISION_RECORD.md` — the reason
is almost always written down.*
