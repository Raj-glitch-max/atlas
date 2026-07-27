# Deploying the Atlas demo

> **Currently live:**
> frontend <https://atlas-dh1.pages.dev> (Cloudflare Pages) ·
> backend <https://atlas-production-c457.up.railway.app> (Railway).
> Verified end to end in a real browser: the console badge reads
> `live · domain-a.test`, issuance/verify/revoke all hit the real engine, and
> two browser sessions cannot see or revoke each other's capabilities.
>
> **The live instance runs fully in-memory, on purpose.** See "Durability vs.
> non-root" below — it is a real trade, and the demo takes the non-root side.

Two pieces: the **backend** (`atlas-server`, a small Go binary in a distroless
container) and the **frontend** (`ui/`, a static Vite build). They are wired
together by one build-time variable and one CORS setting.

Everything here has been verified locally against the real production container.
The only steps that cannot be done without your credentials are the two `login`
commands.

---

## 1. Backend → Railway

Railway builds `deploy/Dockerfile` and injects a `PORT` at runtime.
`railway.json` already points at the right Dockerfile and health check.

```sh
npm i -g @railway/cli
railway login                      # opens a browser — needs you
railway init                       # create/link the project
railway up                         # build + deploy
railway domain                     # mint the public URL
```

Then set the environment. **Do this before pointing the UI at it:**

```sh
railway variables --set ATLAS_ALLOW_ORIGIN=https://<your-pages-domain>
railway variables --set ATLAS_TRUST_DOMAIN=domain-a.test
railway variables --set ATLAS_RATE_LIMIT=120
railway variables --set ATLAS_TRUST_PROXY=true
```

**`ATLAS_TRUST_PROXY` is not optional on a PaaS.** Every request reaches the
container from the platform's proxy, so the server sees one identical source
address for every visitor on earth. Without this flag, `ATLAS_RATE_LIMIT=120`
stops being "120 per visitor" and becomes "120 for the entire internet" — one
enthusiastic visitor or a single crawler locks everyone out, and the demo looks
broken. With it, the client is taken from `X-Forwarded-For`.

Do **not** set it if you ever expose the server directly without a proxy in
front: `X-Forwarded-For` is client-supplied, so trusting it there lets anyone
evade the limit by rotating a header. Rate limiting here is an availability
backstop, never a security control.

### Durability vs. non-root — pick one, knowingly

This is a genuine conflict, and it bites on every PaaS, not just Railway.

The image runs as **uid 65532 (non-root)** — a property the README advertises.
Railway (like most platforms) mounts volumes owned by **root**. So the moment
you attach a volume at `/data` and pass `-key /data/authority.key`, the server
cannot write it and refuses to start:

```
atlas-server: write key file /data/authority.key: permission denied
```

That refusal is deliberate, not a bug. `-key` is a request for durability, and
silently falling back to an ephemeral key would mean every record issued before
a restart quietly stops verifying. The server fails loudly instead.

You have three options:

| Option | Keeps non-root | Durable | Use when |
|---|---|---|---|
| **In-memory (default here)** | ✅ | ❌ | A public demo |
| Volume + `RAILWAY_RUN_UID=0` | ❌ runs as root | ✅ | You need persistence more than the hardening |
| Volume you can chown | ✅ | ✅ | Self-hosted Docker/K8s, where you control the mount |

**The live demo takes option 1**, and the reasoning is worth stating because it
is not obvious: visitor state is *already* ephemeral by design — sessions are
in-memory, capped, and evicted after 30 minutes idle — so durability would buy
almost nothing here. Meanwhile this repository tells people to run a non-root,
distroless container. The flagship public instance should embody the posture the
project advocates rather than quietly contradict it.

On self-hosted Docker or Kubernetes you can have both, because you control the
mount: `chown -R 65532:65532` the host directory, or set `fsGroup: 65532` in the
pod security context. `deploy/docker-compose.yml` already does the right thing.

If you *do* want durability on Railway, set `RAILWAY_RUN_UID=0` and restore the
flags — accepting that the container then runs as root:

```sh
railway variables --set RAILWAY_RUN_UID=0
# and add back:  -store /data/state.json -key /data/authority.key
```

Do **not** set `ATLAS_API_KEY` for a public demo — it would make `/issue` and
`/revoke` require a bearer token, and the whole point is that a visitor can
drive it. `-rate-limit 120` is the protection instead.

### Why `PORT`, not `-addr`

The image sets `ENV PORT=8087` and the start command deliberately does **not**
pass `-addr`. The server binds `0.0.0.0:$PORT` whenever `PORT` is set. Passing
`-addr` explicitly wins over `PORT`, which on Railway silently binds the wrong
port and fails the health check with no useful error. This was a real bug, fixed
in the Dockerfile — don't reintroduce it.

### Fallback: Render free tier

Render's free web services spin down after ~15 minutes idle, so a cold visitor
waits ~30s for a first byte — bad for a "trust it in under a minute" demo. If
you use Render anyway, keep it warm:

```yaml
# .github/workflows/keepalive.yml
name: keepalive
on:
  schedule: [{ cron: "*/10 * * * *" }]
jobs:
  ping:
    runs-on: ubuntu-latest
    steps:
      - run: curl -fsS https://<your-render-host>/health
```

`/health` is a cheap liveness probe and is excluded from access logs by default.

---

## 2. Frontend → Cloudflare Pages

```sh
cd ui
npm ci
VITE_ATLAS_API=https://<your-railway-domain> npm run build
npx wrangler pages deploy dist --project-name atlas   # `wrangler login` first
```

Or connect the repo in the Cloudflare dashboard with:

| Setting | Value |
|---|---|
| Build command | `npm ci && npm run build` |
| Build output directory | `ui/dist` |
| Root directory | `ui` |
| Environment variable | `VITE_ATLAS_API` = your Railway URL |

**`VITE_ATLAS_API` is the whole wiring.** It is baked into the bundle at build
time. Without it the deployed page falls back to `http://127.0.0.1:8087` —
i.e. every visitor's browser tries to reach *their own* laptop, finds nothing,
and silently shows the scripted demo instead of the real engine. Verify it
landed:

```sh
grep -o 'https://[^"]*railway[^"]*' ui/dist/assets/api-*.js
```

---

## 3. The CORS handshake

The two settings must agree exactly, including scheme and no trailing slash:

```text
Railway:            ATLAS_ALLOW_ORIGIN = https://atlas.pages.dev
Cloudflare Pages:   VITE_ATLAS_API     = https://atlas.up.railway.app
```

If `ATLAS_ALLOW_ORIGIN` is left unset it defaults to `*`, which works but is
worth pinning once the Pages domain is stable. The server also allows the
`X-Atlas-Session` header, which the UI needs for per-visitor isolation.

---

## 4. Verify like a stranger

Deployment is not done when it deploys; it is done when a cold visitor succeeds.
In a **fresh incognito window**, with no atlas-server running locally:

```sh
curl -fsS https://<backend>/health      # {"status":"ok",...}
curl -fsS https://<backend>/readyz      # {"ready":true,...}
curl -fsS https://<backend>/activity    # aggregated counts
```

Then in the browser, confirm every one of these:

- [ ] The console badge reads **`live · <trust domain>`**, not `demo · scripted`.
      If it says scripted, `VITE_ATLAS_API` did not make it into the build.
- [ ] "See it verify" produces an **ACCEPT** with a five-row decision trace and
      a real µs latency (a real one varies run to run; the scripted demo always
      shows the same number).
- [ ] The revocation flow ends in **REJECT · RevokedObservable**.
- [ ] A tampered record is refused.
- [ ] Open a second browser profile: its delegation list and audit log are
      **empty**, not showing the first visitor's activity.
- [ ] Revoking in one profile does **not** reject the other profile's record.
- [ ] `/activity` counts go up across both, exposing no identifiers.

## 5. When the backend is down

The failure is designed to be quiet rather than embarrassing. `createApi()`
aborts after 4s, `/health` fails, `enableLive()` is never called, and the site
stays in scripted mode with the badge honestly reading `demo · scripted`. No
error dialog, no broken layout, no fabricated verdict presented as live —
the numbers shown in that mode come from `docs/BENCHMARKS.md` and are labelled.
