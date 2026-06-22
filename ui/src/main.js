// Orchestrator. Loads the light path first (styles, static content, verify
// console, palette, boot). The heavy 3D (three) + GSAP are code-split and
// lazy-loaded ONLY when motion is allowed and WebGL is present — so first
// paint and the reduced-motion path stay cheap.
import "./styles.css";
import { METRICS, LAB } from "./data.js";
import { initVerify, runVerify, runRevocation, enableLive } from "./verify.js";
import { initPalette } from "./palette.js";

const $ = (s, r = document) => r.querySelector(s);
const $$ = (s, r = document) => [...r.querySelectorAll(s)];
const reduce = matchMedia("(prefers-reduced-motion: reduce)").matches;

/* ---------- static content ---------- */
function sparkP(v, w, h) { const mn = Math.min(...v), mx = Math.max(...v), r = (mx - mn) || 1; return v.map((x, i) => `${i ? "L" : "M"} ${(i / (v.length - 1)) * w} ${h - ((x - mn) / r) * (h - 5) - 2}`).join(" "); }
function fillMetrics() {
  $("#metrics").innerHTML = METRICS.map((m) => {
    const flat = Math.min(...m[4]) === Math.max(...m[4]);
    const last = 32 - ((m[4][m[4].length - 1] - Math.min(...m[4])) / ((Math.max(...m[4]) - Math.min(...m[4])) || 1)) * 28 - 2;
    return `<div class="metric"><div class="l">${m[0]}</div><div class="v"><span class="n" data-to="${m[1]}">${m[1]}</span><span class="u"> ${m[2]}</span></div><div class="d">${m[3]}</div>
      <svg class="spk" viewBox="0 0 200 32" preserveAspectRatio="none"><path d="${sparkP(m[4], 200, 28)}" fill="none" stroke="${flat ? "#39404b" : "#7DD3FC"}" stroke-width="1.5"/>${flat ? "" : `<circle cx="200" cy="${last}" r="2.4" fill="#7DD3FC"/>`}</svg></div>`;
  }).join("");
}
function fillLab() {
  $("#labGrid").innerHTML = LAB.map((l) => `<div class="card card-pad feat">
    <div style="color:var(--fg);font-family:var(--mono);font-size:13px">${l[0]}</div>
    <p style="margin:8px 0 16px;font-size:13px">${l[1]}</p>
    ${l[2] === "green" ? '<span class="tag good">● GREEN</span>' : '<span class="tag">○ AWAITING HOST</span>'}</div>`).join("");
}

/* ---------- nav + keycap motif ---------- */
function initNav() { const nav = $("#nav"); const on = () => nav.classList.toggle("stuck", scrollY > 24); on(); addEventListener("scroll", on, { passive: true }); }
function initKeys() {
  const keys = $$("#keys .key"); if (!keys.length || reduce) return;
  let i = 0, iv = null;
  const step = () => { keys.forEach((k) => k.classList.remove("lit")); keys[i % keys.length].classList.add("lit"); i++; };
  // only animate while the keycap row is actually on screen
  const io = new IntersectionObserver(([e]) => {
    if (e.isIntersecting && !iv) iv = setInterval(step, 900);
    else if (!e.isIntersecting && iv) { clearInterval(iv); iv = null; }
  });
  io.observe($("#keys"));
}

/* ---------- mobile nav overlay ---------- */
function initMobileNav() {
  const btn = $("#mbtn"), ov = $("#mnav"); if (!btn || !ov) return;
  const links = $$("a", ov);
  const open = () => { ov.hidden = false; requestAnimationFrame(() => ov.classList.add("open")); btn.classList.add("x"); btn.setAttribute("aria-expanded", "true"); document.documentElement.classList.add("mnav-lock"); };
  const close = () => { if (ov.hidden) return; ov.classList.remove("open"); btn.classList.remove("x"); btn.setAttribute("aria-expanded", "false"); document.documentElement.classList.remove("mnav-lock"); setTimeout(() => { ov.hidden = true; }, 320); };
  btn.addEventListener("click", () => (btn.classList.contains("x") ? close() : open()));
  ov.addEventListener("click", (e) => { if (e.target.closest("a")) close(); });
  addEventListener("keydown", (e) => {
    if (ov.hidden) return;
    if (e.key === "Escape") { close(); btn.focus(); }
    else if (e.key === "Tab") { // keep focus cycling between the button and the overlay links
      const f = [btn, ...links], first = f[0], last = f[f.length - 1];
      if (e.shiftKey && document.activeElement === first) { e.preventDefault(); last.focus(); }
      else if (!e.shiftKey && document.activeElement === last) { e.preventDefault(); first.focus(); }
    }
  });
  addEventListener("resize", () => { if (innerWidth > 820) close(); });
}

/* ---------- boot — progress tracks real milestones (fonts ready, motion
   modules loaded), compressed to a quick sweep when they finish fast.
   Hard failsafe: can never trap the page. ---------- */
function removeBoot() { const b = $("#boot"); document.body.classList.add("ready"); if (!b || b.classList.contains("gone")) return; b.classList.add("gone"); setTimeout(() => b.remove(), 750); }
const boot = (() => {
  let target = 12, shown = 0, raf = 0, ended = false;
  const LABELS = [[0, "loading type"], [45, "loading the field"], [96, "field online"]];
  function tick() {
    const bar = $("#bootB"), txt = $("#bootT");
    shown += (target - shown) * 0.14; if (target >= 100 && shown > 99.4) shown = 100;
    if (bar) bar.style.width = shown + "%";
    const l = LABELS.filter((x) => shown >= x[0]).pop(); if (txt && l) txt.textContent = l[1];
    if (shown >= 100) { end(); return; }
    raf = requestAnimationFrame(tick);
  }
  function end() { if (ended) return; ended = true; cancelAnimationFrame(raf); setTimeout(removeBoot, 180); }
  return {
    start() {
      if (reduce) { removeBoot(); return; }
      raf = requestAnimationFrame(tick);
      if (document.fonts && document.fonts.ready) document.fonts.ready.then(() => { target = Math.max(target, 55); });
      setTimeout(() => { end(); removeBoot(); }, 3200); // failsafe
    },
    done() { target = 100; },
  };
})();

/* ---------- WebGL support probe ---------- */
function webglOK() { try { const c = document.createElement("canvas"); return !!(window.WebGLRenderingContext && (c.getContext("webgl2") || c.getContext("webgl"))); } catch { return false; } }

/* ---------- init ---------- */
async function init() {
  const palette = [
    { label: "Run verification", key: "↵", icon: '<path d="m5 12 5 5L20 7"/>', run: () => { location.hash = "#verify"; setTimeout(runVerify, 500); } },
    { label: "Simulate revocation under partition", icon: '<circle cx="12" cy="12" r="9"/><path d="m5.6 5.6 12.8 12.8"/>', run: () => { location.hash = "#verify"; setTimeout(runRevocation, 500); } },
    { label: "Jump to TrustPerf", icon: '<path d="M3 3v18h18"/><path d="M7 14l3-4 3 3 5-7"/>', run: () => (location.hash = "#perf") },
    { label: "Jump to Conformance", icon: '<path d="M9 11l3 3L22 4"/>', run: () => (location.hash = "#conformance") },
    { label: "Jump to Atlas Lab", icon: '<path d="M9 3v5.5L4.5 17A2 2 0 0 0 6.3 20h11.4a2 2 0 0 0 1.8-3L15 8.5V3"/>', run: () => (location.hash = "#lab") },
    { label: "Back to top", icon: '<path d="m6 15 6-6 6 6"/>', run: () => (location.hash = "#top") },
  ];
  try {
    fillMetrics(); fillLab(); initVerify(); initNav(); initKeys(); initMobileNav(); initPalette(palette);

    // Detect the Atlas Server; if reachable, route the verify console through
    // the real backend (live verdicts). Silent no-op when it isn't running.
    (async () => {
      try {
        const { createApi } = await import("./api.js");
        const api = createApi();
        const h = await api.health();
        if (h.ok) { const v = await api.version().catch(() => null); enableLive(api, v); }
      } catch { /* stays in scripted demo mode */ }
    })();

    if (!reduce) {
      document.documentElement.classList.add("anim");
      // GSAP smooth-scroll enhancements load whenever motion is allowed; the
      // heavy WebGL hero is loaded only when WebGL is actually present. Any
      // failure degrades to the static, fully-visible site.
      import("./scroll.js")
        .then(async (scrollMod) => {
          let hero = null;
          if (webglOK()) {
            try { const h = await import("./hero.js"); hero = await h.initHero($("#field")); }
            catch (e) { console.warn("[atlas] hero failed, static hero:", e); }
          }
          scrollMod.initScroll(hero);
          boot.done();
        })
        .catch((e) => { console.warn("[atlas] scroll enhancement failed:", e); document.documentElement.classList.remove("anim"); $$(".reveal").forEach((el) => el.classList.add("in")); boot.done(); });
    } else {
      $$(".reveal").forEach((el) => el.classList.add("in"));
    }
  } catch (err) {
    console.error("[atlas] init error:", err);
    $$(".reveal").forEach((el) => el.classList.add("in"));
  } finally {
    boot.start();
  }
}
if (document.readyState === "loading") document.addEventListener("DOMContentLoaded", init); else init();
