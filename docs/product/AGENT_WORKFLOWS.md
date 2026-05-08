# Atlas for multi-tool agent workflows

*The honest answer to: "I hop across Claude → Cursor → GitHub → Vercel →
Slack → … all day. Does Atlas fix that?"*

## Short answer

**No — and it shouldn't try.** The hopping is an *orchestration* problem
(one thing doing the whole task), which is the job of an agent runtime, not a
trust primitive. If Atlas became an orchestrator it would be a worse version of
tools that already exist, and it would abandon the one thing it's uniquely
good at.

**But the dangerous part of that workflow is Atlas's job exactly.** The moment
you let an agent act across GitHub, Vercel, and Slack *for you*, you face a
security problem that is real and, today, mostly unsolved:

> To let the agent do the work, do you hand it your GitHub token, your Vercel
> token, and your Slack token?

If yes, you just gave a probabilistic, promptable, possibly-compromised process
**standing, unscoped, unrevocable, unauditable** access to everything you own.
That's the pain under the pain. Atlas removes it.

## What Atlas actually is here: the capability broker

Atlas is the **trust fabric** for the agent workflow, not the workflow:

1. You (or a coordinator) issue each task-agent a **scoped, expiring,
   revocable capability** — *"may `github:push:acme/landing` for 15 minutes"* —
   instead of a credential.
2. Each tool sits behind a thin **adapter** that **verifies the capability
   offline** and, only if it authorizes *this exact action*, performs it using
   **the tool's own credential** — which the agent never sees.
3. A compromised or confused agent is confined to its one scope. One `revoke`
   kills it everywhere, instantly, even offline.

The agent holds capabilities; the adapters hold secrets. That inversion is the
whole game.

```
  you ──issue scoped capability──▶ agent ──presents──▶ tool adapter
                                                           │ verify offline
                                                           │ require-scope?
                                                           ▼
                                              acts with the TOOL's own token
```

**Run it:** `bash examples/ship-a-landing-page.sh` — three agents, three tools,
least privilege, blast-radius containment, instant revocation, all offline,
all on the real engine.

## The senior-engineer questions, answered

**"To let the agent push, deploy, and post — do I give it all three tokens?"**
No. You issue three separate capabilities, each granting exactly one action.
The tokens stay with the adapters. The agent's most powerful move is presenting
a capability it already holds.

**"An agent will get prompt-injected or go rogue. What's the blast radius?"**
Exactly its scope, for exactly its TTL. The `coder` capability grants
`github:push:acme/landing` and *nothing else* — the Vercel and Slack adapters
deny it (`atlas verify --require-scope` returns "authorized: NO — valid, but
does not grant …"). Demonstrated live in the script, step 3.

**"Something's wrong mid-run. How fast can I contain it?"**
One `atlas revoke <instance>`. On the adapters' next capability check the agent
is denied — even for the action it *was* allowed — bounded by the freshness
window R, and it works offline (the adapter checks a signed snapshot, not a
live call). Step 4.

**"My deploy runner / edge function can't reach my auth server."**
It doesn't need to. Verification is offline: the adapter holds a trust bundle
(public key + signed revocation snapshot) and answers in microseconds with no
network call to the issuer. That's Atlas's core property, not a bolt-on.

**"How do I audit what happened, with which authority?"**
Every issue, verify, and revoke is a logged event (`atlas audit`), and the
delegation graph (`atlas graph`) shows who could do what. The capability itself
records principal, delegate, scope, expiry, and instance — inspectable with
`atlas inspect`.

**"Can a coordinator agent delegate to worker agents?"**
Single hop: the authority (you) delegates to each worker directly, with
attenuated scope. Atlas is deliberately **not** a re-delegation chain
(A→B→C) — see `LIMITATIONS.md` §1. If you need the coordinator to mint
sub-capabilities, it does so by asking the authority, not by re-signing.

**"Isn't this just OAuth scopes / Vault / SPIFFE?"**
- *OAuth* needs a live authorization server and is client↔server; Atlas is
  agent↔agent and verifies offline (`WHY.md` Q3).
- *Vault/secrets brokers* hand out the real secret; Atlas hands out a
  capability and keeps the secret at the adapter — different blast radius.
- *SPIFFE* issues the identities Atlas binds to; Atlas rides on top to
  authorize actions *between* them (INV10/INV11).

## The honest gap (and the bridge)

The demo's adapters are **mocks** — they verify a real Atlas capability, then
*simulate* the GitHub/Vercel/Slack call. To be usable in your real workflow,
Atlas needs **real service adapters**: small gateways (or MCP servers, or
sidecars) that hold a service credential, verify an Atlas capability, and
exchange it for the real action or a short-lived service token.

That bridge — a **capability→token broker** and a set of reference adapters —
is the roadmapped path from "compelling demo" to "in your stack" (`ROADMAP.md`,
next tier). It's genuinely additive: the trust primitive underneath is done and
tested; adapters are integration work, not new cryptography.

## What this is not

Atlas is not Devin, not Zapier, not an MCP hub, not your IDE. It doesn't reduce
the number of tools you touch. It makes it *safe to let an agent touch them for
you* — which is the only version of "solve the tool hop" that a senior engineer
will actually trust.
