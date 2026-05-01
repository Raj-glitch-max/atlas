# The pitch — three lengths, several audiences

Same truth, compressed differently. Pick the one that fits the room.

---

## 30 seconds (anyone)

> Atlas is the software you use when one AI agent needs to safely let another
> agent act with its authority. Instead of sharing an API key or depending on
> a live auth server, an agent hands off a signed, scoped, expiring, revocable
> capability — and the other side can verify it in microseconds, even offline.
> Kill the network and verification still works; reconnect and revocation is
> enforced instantly.

## 2 minutes (technical peer / hiring manager)

> Multi-agent systems have a delegation problem: when agent A asks agent B to
> act on its behalf, today you either share a credential — now B *is* A forever
> — or you route every action through an authorization server that has to be
> up. Atlas replaces both. A delegation is a signed record: *B may do exactly
> this scope, until this time, revocable via this handle.* It's bound to
> SPIFFE workload identity, the scope is strictly attenuated at issuance, and —
> the key part — a relying party verifies it from locally-held trust material
> with **no live call to the issuer**. Revocation rides an independent signed
> channel with *verifiable freshness*: the verifier chooses its own staleness
> budget and fails **closed** past it, never open.
>
> It's a real product, not a paper: a server, a CLI (`atlas delegate`), agent
> tools over MCP so Claude can call it directly, a Python SDK, and a
> 28-vector conformance suite — 18 adversarial — plus fuzzing and published
> benchmarks (~94µs verify). The whole thesis runs live in one script: issue,
> kill the server, verify offline, revoke, reject offline, refuse a tampered
> bundle. And I documented what it *doesn't* do as carefully as what it does.

## 10 minutes (principal engineer / partner)

Structure it as: **problem → why existing tools don't fit → the design → the
proof → the honest limits → why it compounds.**

1. **Problem.** Agents delegating to agents is becoming the default
   architecture. Every hand-off is an authorization decision made with shared
   secrets or bespoke glue — unauditable, credential-leaking, and coupled to a
   live authority. (Concrete pain list: `WHY.md` Q2.)
2. **Why not the obvious tools.** OAuth (live server, client↔server), JWT
   (identity of signer, not continued trustworthiness), Biscuit/Macaroons
   (great offline attenuation, but no identity binding, revocation-freshness,
   or agent runtime). Atlas composes with the first two and credits the third.
   (`WHY.md` Q3–Q5, `docs/product/OBJECTIONS.md`.)
3. **The design.** Single-hop, two-domain, SPIFFE-bound, attenuated-at-issuance
   delegation; a five-check verifier (identity binding, signature, freshness,
   scope, revocation) that emits an unconditional decision trace; revocation as
   signed snapshots with verifiable freshness and a fail-closed staleness
   bound. Deliberately boring crypto (ES256 JWS).
4. **The proof.** `make ci` green; conformance + adversarial vectors; a
   coverage-guided fuzzer with zero silent acceptances; every threat-model
   claim mapped to a passing test (`THREAT_MODEL.md`); the offline+revoke+tamper
   story runnable in one script.
5. **The honest limits.** Single-hop, within-window replay, manual key
   rotation, dev-grade deployment, unaudited crypto, the packet-level offline
   proof still pending real infra. All published (`LIMITATIONS.md`). Drawing
   this boundary is the credibility.
6. **Why it compounds.** The moat is the accumulated
   protocol+verifier+revocation+conformance+ecosystem work and the interop
   trust it earns — the Kubernetes/Git pattern: become the neutral layer
   everyone tests against, not a vendor feature.

---

## The résumé line

> Designed and built Atlas, an open-source runtime for **secure delegation
> between AI agents** — offline-verifiable capability tokens bound to SPIFFE
> identity, with an independent verifiable-freshness revocation channel,
> adversarial conformance testing (28 language-neutral vectors), coverage-guided
> fuzzing, a five-check verification engine, HTTP/CLI/MCP/SDK surfaces, and
> published microsecond-scale benchmarks — with a formally-stated threat model
> and published limitations.

## The one-liner (README/social)

> **Atlas — delegation that survives the network, and can't be lied to.**
