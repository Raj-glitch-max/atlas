# examples

Runnable demonstrations of Atlas against the real engine.

## `ship-a-landing-page.sh` — the multi-tool agent workflow

The scenario everyone actually has: an agent must act across GitHub, Vercel,
and Slack for you. Instead of handing it every token, you issue each task-agent
a **scoped, expiring, revocable capability**; each tool verifies it **offline**
and acts with its own credential. Shows least privilege, blast-radius
containment (a compromised agent is confined to one scope), and instant
revocation — all offline. See `docs/product/AGENT_WORKFLOWS.md`.

```bash
bash examples/ship-a-landing-page.sh
```

## `unforgettable.sh` — delegation that survives the network

The core thesis in 60 seconds: issue → kill the server → verify offline →
revoke → reject offline → refuse a tampered bundle.

```bash
bash examples/unforgettable.sh
```

## `agent-capability-demo.sh`

The agent-capability lifecycle — the "OAuth for agents" loop, honestly. It is
self-contained: builds the binaries, runs its own server on a temp store, and
tears everything down.

```bash
bash examples/agent-capability-demo.sh
# or on a busy port:  PORT=8092 bash examples/agent-capability-demo.sh
```

It shows, end to end:

1. a **principal grants a scoped capability** to an agent — **single hop**;
2. **attenuation** — an over-scope grant is refused by the engine (`OverScope`);
3. a **relying party verifies** the capability from locally-held trust material
   — the verification makes **no call to the issuer**;
4. the principal **revokes** the capability;
5. the relying party **observes the revocation** and rejects (`RevokedObservable`);
6. the **audit trail** and **trust graph**.

### Honest scope

Atlas is **single-hop by design**. This is a principal → delegate grant
verifiable across two trust domains — *not* multi-hop re-delegation (A→B→C),
which is deliberately outside the primitive. The demo reflects exactly what the
primitive does, no more.
