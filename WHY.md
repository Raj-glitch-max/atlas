# Why Atlas

> **Atlas is the software you use when one AI agent needs to safely authorize
> another AI agent.**

Your agents already call tools, services, and other agents. Today they do it
by sharing API keys, forwarding credentials, or trusting a central
authorization server to be up. Atlas lets them **hand off scoped, expiring,
revocable permission instead of handing over secrets** — and lets the
receiving side check that permission even when the network to the issuer is
gone.

See it in 60 seconds:

```bash
bash examples/unforgettable.sh
```

An agent gets a capability → **the server dies** → verification still answers
in microseconds → the network comes back → revocation is enforced immediately
→ a tampered trust bundle is refused outright. Real engine, no mocks.

---

## Q1 · Why should I care?

Because the painful part of multi-agent systems isn't the AI — it's what
happens when agent A asks agent B to do something with agent A's authority.
Today you either share a key (now B *is* A, forever, everywhere) or you build
a delegation service (now every action depends on it being reachable). Atlas
replaces both with a signed capability: *B may do exactly this, until then,
unless revoked* — checkable anywhere.

## Q2 · What pain disappears?

- **Shared API keys** between agents and services — gone.
- **Permanent credentials** for temporary jobs — gone.
- **Hand-rolled delegation logic** (the code everyone hates writing and no one
  audits) — gone.
- **Authorization that dies when the auth server does** — gone; verification
  is local.
- **"Who could do what, when?" audits that are impossible** — every grant,
  check, and revocation is a signed, replayable event.
- **Every service re-implementing verification differently** — one conformant
  verifier, one test-vector suite, identical decisions everywhere.

## Q3 · Why not OAuth?

Use OAuth. Atlas isn't competing with it. OAuth solves identity,
authentication, and authorization **between clients and servers** with a live
authorization server in the loop. Atlas solves **delegation between
autonomous agents** — where the checker may not be able to reach the issuer,
and where authority must be attenuated (narrowed) at each hand-off. They
compose: your humans log in with OAuth; your agents delegate with Atlas.

## Q4 · Why not just JWT?

A JWT answers *"who signed this?"* Atlas answers *"can this delegated
authority still be trusted — here, now, under these conditions?"* That extra
distance is exactly: scope attenuation enforced at issuance, expiry, an
independent revocation channel with **verifiable freshness**, and a defined
fail-closed answer when knowledge is stale. (Under the hood an Atlas record
*is* a signed JWS — deliberately boring cryptography. The value is the
verification semantics, not a new token format.)

## Q5 · Why not Biscuit or Macaroons?

They're good prior art, and Atlas doesn't pretend to have invented offline
attenuation — they solved that first. What Atlas adds is everything *around*
the token that makes delegation deployable between agents:

- binding to **workload identity** (SPIFFE) instead of bare keys,
- a **revocation channel with verifiable freshness** and honest partition
  semantics (what exactly can be claimed when the network is split — and what
  cannot),
- a **runtime**: server, CLI, agent tools (MCP), SDK, audit, metrics,
- a **conformance suite** — 28 language-neutral test vectors, including 18
  adversarial ones — so independent implementations provably agree,
- fuzzing, property tests, and benchmarks published, not promised.

If Biscuit fits your problem, use Biscuit. Atlas is for the agent-delegation
workflow, end to end.

## Q6 · Why shouldn't I build this myself?

Because "just sign a token" becomes, within a quarter: protocol design, a
verifier with five ordered checks and a fail-closed policy, revocation
distribution, freshness bounds, key rotation, conformance testing, adversarial
vectors (alg-none, key confusion, signature transplants), fuzzing,
cross-language compatibility, SDKs, audit trails, and replay tooling. You
wouldn't build OAuth in-house. This is the same category of decision.

## Q7 · When should I use Atlas?

Use Atlas when **most** of these are true:

- agents (or services) delegate work to other agents,
- the permission should be **temporary** and **narrower** than the grantor's,
- delegation crosses a team, org, or trust-domain boundary,
- the checker can't — or shouldn't — depend on a live call to the issuer,
- you need an audit trail that survives interrogation.

## Q8 · When should I NOT use Atlas?

Honestly: often. Don't use Atlas when —

- you have **one agent** talking to **one service** → use an API key,
- permissions are **static and long-lived** → use your IAM,
- everything lives inside one process or one trust domain with a reliable
  authorization service → use that service,
- nothing is ever delegated → there's nothing for Atlas to do.

Atlas earns its complexity only when delegation actually exists.

## Q9 · Why now?

AI systems are moving from one model answering questions to **many
specialized agents doing work** — calling tools, spawning sub-agents,
crossing service and organization boundaries. Every one of those hand-offs is
a delegation decision made today with shared secrets or ad-hoc glue.
Delegation is becoming infrastructure; infrastructure needs a standard shape.

## Q10 · Why will this become inevitable?

We won't claim "Atlas will be the next OAuth." The defensible version: **if**
multi-agent systems keep growing — and every framework roadmap says they will
— developers will need a standard way for one agent to safely authorize
another without sharing credentials or depending on a live authority. Atlas
is built for exactly that problem, with the conformance surface that lets it
outgrow any single implementation.

---

## Where to go next

| You want to… | Go to |
|---|---|
| **See it save a real workflow** | `docs/product/AGENT_WORKFLOWS.md` · `bash examples/ship-a-landing-page.sh` |
| **Try it in 5 minutes** | `PRODUCT.md` → quickstart, or `bash examples/unforgettable.sh` |
| Read what can go wrong | `THREAT_MODEL.md` · `LIMITATIONS.md` |
| See the hard questions answered | `docs/product/OBJECTIONS.md` |
| Understand the engine | `SYSTEM_ARCHITECTURE.md` · `tests/vectors/VECTORS.md` |
| The research underneath | `docs/discovery/` · `rfc/` |

<!-- checkpoint: repo(trust-anchors): update trust anchors -->

<!-- checkpoint: test(stores): test error wrappers -->
