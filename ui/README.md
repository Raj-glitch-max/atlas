# Atlas — site (Vite + Three.js + GSAP)

Premium dark product site for the Atlas delegation primitive. Raycast
restraint + Linear chrome discipline + Igloo/Lusion-style depth: a **WebGL
hero** (two trust-domain constellations joined by a verifiable signature path)
with **GSAP ScrollTrigger** "camera-beat" scrolling, over a crisp, accessible
HTML overlay. Real **Geist** type (OFL).

## Run

```bash
cd ui
npm install        # first time only (three, gsap, vite)
npm run dev        # → http://localhost:5173/
```

Build / preview the production bundle:

```bash
npm run build      # → dist/
npm run preview    # serves dist/ at http://localhost:4173/
```

## Architecture (reusable, code-split)

```
src/
  main.js      orchestrator — light path first; lazy-loads 3D only when allowed
  styles.css   design system (Geist scale, four-step surface ladder, reveals)
  data.js      shared content (metrics, lab, stages)
  verify.js    2D verify console + accept / partition cinematics (no 3D dep)
  hero.js      Three.js WebGL hero  ← lazy-loaded chunk
  scroll.js    GSAP ScrollTrigger camera beats + reveals  ← lazy-loaded chunk
  palette.js   ⌘K command palette
public/fonts/  Geist · Geist Mono · Geist Pixel (.woff2, OFL)
```

**Performance / safety**
- Initial JS ≈ 7 kB gzip; `three` (~117 kB gz) and `gsap` (~28 kB gz) are
  separate chunks loaded only when motion is allowed **and** WebGL is present.
- WebGL glow uses additive point sprites — **no post-processing pass** — so it
  holds 60fps on mid-range laptops.
- `prefers-reduced-motion`: no 3D, no scroll animation, one static frame,
  all content revealed immediately.
- Boot overlay has a hard failsafe (3.2s) and the whole enhancement path is
  `try/catch` — a WebGL/JS failure degrades to the static site, never a blank
  or trapped page.

## Live mode (wired to the real backend)

The verify console runs against the **real Atlas Server** when it's reachable,
and falls back to the scripted demo when it isn't. Run both:

```bash
# terminal 1 — the backend (real engine)
go run ./cmd/atlas-server            # → 127.0.0.1:8087

# terminal 2 — the site
cd ui && npm run dev                 # → localhost:5173
```

Open <http://localhost:5173/>. The console badge flips to **live · domain-a.test**
and "Run verification" / "Revoke under partition" now issue, verify, and revoke
against the server — showing the **real** decision, the real five-check trace,
and the measured latency (µs). Point at a different server with
`localStorage.atlasApi = "http://host:port"` in the browser console.

When the server isn't running, the badge reads **demo · scripted** and the
cinematics play with representative data — the site is always presentable.

## Operator console

`/console.html` (linked from the nav) is a **live operator surface**: metrics,
a delegations table with one-click revoke, an issue form, the audit log, and
the trust graph — all from the real server, auto-refreshing every 3s. When the
server is down it shows a "start the server" panel. For a guarded server, set
`localStorage.atlasApiKey` in the browser console.

## Data

Demo-mode figures mirror the real implementation (94µs verify, 403-byte proof,
R = 2s, the five checks, `RevokedObservable` / `RevocationKnowledgeStale`).
Live mode shows actual server results.

Additive; not part of the frozen Atlas primitive.
