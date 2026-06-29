// The interactive verification console: a 2D-canvas trust graph + the
// five-check pipeline + verdict panel, with two cinematics (accept flow and
// revocation-under-partition). Independent of the WebGL hero, so it works even
// when 3D is unavailable or reduced-motion is on.
import { STAGES } from "./data.js";

const $ = (s, r = document) => r.querySelector(s);
const $$ = (s, r = document) => [...r.querySelectorAll(s)];
const reduce = matchMedia("(prefers-reduced-motion: reduce)").matches;
const rand = (a, b) => a + Math.random() * (b - a);

let gcv, gctx, GW = 0, GH = 0, hubs = [], particles = [], ripples = [], t = 0;
const V = { active: false, t0: 0, dur: reduce ? 1 : 2600, done: false, _r: 0, live: false, realRes: null };
const REV = { active: false, seam: 0, reject: false, t0: 0 };
let _pk = {};

// ---- live mode (wired to cmd/atlas-server) ----
let live = null; // the api client when the server is reachable
const PRINCIPAL = "spiffe://domain-a.test/workload/payments-api";
const DELEGATE = "spiffe://domain-b.test/agent/booking-worker";
const CHECKMAP = { identity_binding: "01", integrity: "02", expiry: "03", scope_integrity: "04", revocation_status: "05" };
const sleep = (ms) => new Promise((r) => setTimeout(r, ms));

// enableLive is called by main.js once /health succeeds. It flips the console
// badge and routes the buttons through the real API.
export function enableLive(api, meta) {
  live = api;
  const el = document.getElementById("liveStatus");
  if (el) el.innerHTML = `<span class="d"></span>live · ${meta && meta.trustDomain ? meta.trustDomain : "atlas-server"}`;
}

function decisionKind(d) { return d === "inconclusive" ? "inconc" : d; }
function outcomeState(o) {
  const s = (o || "").toLowerCase();
  if (s.includes("pass")) return ["pass", "PASS"];
  if (s.includes("reject") || s.includes("fail")) return ["reject", "FAIL"];
  if (s.includes("inconc")) return ["wait", "INCONCL"];
  return ["pass", s.toUpperCase() || "—"];
}
function applyTraceToPipe(trace) {
  (trace || []).forEach((e) => {
    const k = CHECKMAP[e.check]; if (!k) return;
    const [st, fl] = outcomeState(e.outcome);
    setStg(k, st, fl);
  });
}

function sizeGraph() {
  if (!gcv) return;
  const dpr = Math.min(devicePixelRatio || 1, 2);
  GW = gcv.clientWidth; GH = gcv.clientHeight;
  gcv.width = GW * dpr; gcv.height = GH * dpr; gctx.setTransform(dpr, 0, 0, dpr, 0, 0);
  const cy = GH * 0.5;
  hubs = [
    { n: "ISSUER", s: "domain-a", dom: 0, x: GW * 0.13, y: cy - GH * 0.14 },
    { n: "DELEGATE", s: "agent", dom: 0, x: GW * 0.38, y: cy + GH * 0.16 },
    { n: "VERIFIER", s: "domain-b", dom: 1, x: GW * 0.64, y: cy - GH * 0.16 },
    { n: "DECISION", s: "verdict", dom: 1, x: GW * 0.88, y: cy + GH * 0.12 },
  ].map((h) => ({ ...h, r: 5.5, lit: 0 }));
}
function pPoint(k) {
  const seg = Math.min(2, Math.floor(k * 3)), l = k * 3 - seg, a = hubs[seg], b = hubs[seg + 1];
  const mx = (a.x + b.x) / 2, my = (a.y + b.y) / 2 - 30, u = 1 - l;
  return { x: u * u * a.x + 2 * u * l * mx + l * l * b.x, y: u * u * a.y + 2 * u * l * my + l * l * b.y };
}
function drawGraph(now) {
  if (!gcv) return;
  t = now / 1000;
  gctx.clearRect(0, 0, GW, GH);
  const segLit = [0, 0, 0]; let pulse = null;
  if (V.active) {
    const k = Math.min(1, (now - V.t0) / V.dur); pipeByK(k);
    let pk = k < 0.55 ? (k / 0.55) * (2 / 3) : (k < 0.70 ? 2 / 3 : ((k - 0.70) / 0.30) * (1 / 3) + 2 / 3);
    pulse = pPoint(Math.min(0.999, pk));
    for (let s = 0; s < 3; s++) segLit[s] = pk * 3 > s ? Math.min(1, pk * 3 - s) : 0;
    if (!reduce && now % 16 < 9) particles.push({ x: pulse.x, y: pulse.y, life: 1, vx: rand(-0.3, 0.3), vy: rand(-0.3, 0.3) });
    if (pk > 0.98 && ripples.length < 3 && now - V._r > 90) { V._r = now; ripples.push({ x: hubs[3].x, y: hubs[3].y, r: 0, life: 1, c: [74, 222, 128] }); }
    if (k >= 1 && !V.done) { V.done = true; if (V.live) finalizeLive(); else onVerified(); }
  }
  for (let s = 0; s < 3; s++) {
    const a = hubs[s], b = hubs[s + 1], mx = (a.x + b.x) / 2, my = (a.y + b.y) / 2 - 30;
    gctx.beginPath(); gctx.moveTo(a.x, a.y); gctx.quadraticCurveTo(mx, my, b.x, b.y);
    gctx.strokeStyle = "rgba(120,150,180,.18)"; gctx.lineWidth = 1.3; gctx.stroke();
    if (segLit[s] > 0) {
      gctx.save(); gctx.globalCompositeOperation = "lighter";
      gctx.beginPath(); gctx.moveTo(a.x, a.y); gctx.quadraticCurveTo(mx, my, b.x, b.y);
      gctx.strokeStyle = `rgba(125,211,252,${0.5 * segLit[s]})`; gctx.lineWidth = 2;
      gctx.shadowColor = "rgba(125,211,252,.8)"; gctx.shadowBlur = 10; gctx.stroke(); gctx.restore();
    }
  }
  if (REV.active || REV.seam > 0) {
    const el = now - REV.t0; REV.seam = el < 1600 ? Math.min(1, el / 280) : Math.max(0, 1 - (el - 1600) / 380);
    if (REV.seam > 0) {
      const sx = GW * 0.5; gctx.save(); const jit = Math.sin(t * 40) * 1.3 + Math.sin(t * 13) * 1;
      const lg = gctx.createLinearGradient(sx, 0, sx, GH);
      lg.addColorStop(0, "rgba(251,113,133,0)"); lg.addColorStop(0.5, `rgba(251,113,133,${0.55 * REV.seam})`); lg.addColorStop(1, "rgba(251,113,133,0)");
      gctx.strokeStyle = lg; gctx.lineWidth = 1.4; gctx.beginPath();
      for (let y = 0; y <= GH; y += 7) { const xx = sx + Math.sin(y * 0.06 + t * 6) * 3 + jit; y === 0 ? gctx.moveTo(xx, y) : gctx.lineTo(xx, y); }
      gctx.stroke(); gctx.restore();
    }
    if (el > 2100) {
      const k = Math.min(1, (el - 2100) / 900);
      for (let s = 0; s < 3; s++) {
        const seg = Math.min(1, Math.max(0, k * 3 - s)); if (seg <= 0) continue;
        const a = hubs[s], b = hubs[s + 1], mx = (a.x + b.x) / 2, my = (a.y + b.y) / 2 - 30;
        gctx.save(); gctx.globalCompositeOperation = "lighter"; gctx.beginPath(); gctx.moveTo(a.x, a.y); gctx.quadraticCurveTo(mx, my, b.x, b.y);
        gctx.strokeStyle = `rgba(251,113,133,${0.5 * seg})`; gctx.lineWidth = 2; gctx.shadowColor = "rgba(251,113,133,.8)"; gctx.shadowBlur = 9; gctx.stroke(); gctx.restore();
      }
      const pp = pPoint(Math.min(0.999, k)); gctx.save(); gctx.globalCompositeOperation = "lighter";
      const gg = gctx.createRadialGradient(pp.x, pp.y, 0, pp.x, pp.y, 13);
      gg.addColorStop(0, "rgba(255,220,225,.9)"); gg.addColorStop(0.45, "rgba(251,113,133,.7)"); gg.addColorStop(1, "rgba(0,0,0,0)");
      gctx.fillStyle = gg; gctx.beginPath(); gctx.arc(pp.x, pp.y, 13, 0, 6.28); gctx.fill(); gctx.restore();
      if (k >= 1) REV.reject = true;
    }
  }
  gctx.globalCompositeOperation = "lighter";
  for (const p of particles) { p.x += p.vx; p.y += p.vy; p.life -= 0.03; if (p.life <= 0) continue; const gg = gctx.createRadialGradient(p.x, p.y, 0, p.x, p.y, 6); gg.addColorStop(0, `rgba(125,211,252,${0.5 * p.life})`); gg.addColorStop(1, "rgba(0,0,0,0)"); gctx.fillStyle = gg; gctx.beginPath(); gctx.arc(p.x, p.y, 6, 0, 6.28); gctx.fill(); }
  particles = particles.filter((p) => p.life > 0);
  for (const r of ripples) { r.r += 3.2; r.life -= 0.02; if (r.life <= 0) continue; gctx.strokeStyle = `rgba(${r.c[0]},${r.c[1]},${r.c[2]},${0.5 * r.life})`; gctx.lineWidth = 1.6; gctx.beginPath(); gctx.arc(r.x, r.y, r.r, 0, 6.28); gctx.stroke(); }
  ripples = ripples.filter((r) => r.life > 0);
  if (pulse) { const gg = gctx.createRadialGradient(pulse.x, pulse.y, 0, pulse.x, pulse.y, 14); gg.addColorStop(0, "rgba(220,245,255,.95)"); gg.addColorStop(0.4, "rgba(125,211,252,.7)"); gg.addColorStop(1, "rgba(0,0,0,0)"); gctx.fillStyle = gg; gctx.beginPath(); gctx.arc(pulse.x, pulse.y, 14, 0, 6.28); gctx.fill(); }
  gctx.globalCompositeOperation = "source-over";
  hubs.forEach((h, i) => {
    const on = ((V.active && ((i === 0) || (i <= 1 && segLit[0] > 0) || (i <= 2 && segLit[1] > 0) || (i === 3 && V.done))) || (i === 3 && REV.reject) || (REV.active && i <= 1)) ? 1 : 0;
    h.lit += (on - h.lit) * 0.12;
    const dc = h.dom ? [174, 185, 204] : [125, 211, 252];
    const gc = (i === 3 && REV.reject) ? [251, 113, 133] : (i === 3 && V.done) ? [74, 222, 128] : (REV.active && i <= 1) ? [251, 113, 133] : dc;
    gctx.save(); gctx.globalCompositeOperation = "lighter";
    const gr = gctx.createRadialGradient(h.x, h.y, 0, h.x, h.y, 32 + h.lit * 24);
    gr.addColorStop(0, `rgba(${gc[0]},${gc[1]},${gc[2]},${0.1 + h.lit * 0.34})`); gr.addColorStop(1, "rgba(0,0,0,0)");
    gctx.fillStyle = gr; gctx.beginPath(); gctx.arc(h.x, h.y, 32 + h.lit * 24, 0, 6.28); gctx.fill(); gctx.restore();
    gctx.beginPath(); gctx.arc(h.x, h.y, h.r + 4, 0, 6.28); gctx.strokeStyle = `rgba(${gc[0]},${gc[1]},${gc[2]},${0.35 + h.lit * 0.5})`; gctx.lineWidth = 1.3; gctx.stroke();
    gctx.beginPath(); gctx.arc(h.x, h.y, h.r, 0, 6.28); gctx.fillStyle = `rgba(${gc[0]},${gc[1]},${gc[2]},${0.55 + h.lit * 0.45})`; gctx.fill();
    gctx.fillStyle = `rgba(245,247,250,${0.5 + h.lit * 0.5})`; gctx.font = '600 9px "Geist Mono",monospace'; gctx.textAlign = "center"; gctx.fillText(h.n, h.x, h.y - 14);
    gctx.fillStyle = "rgba(138,147,163,.7)"; gctx.font = '8px "Geist Mono",monospace'; gctx.fillText(h.s, h.x, h.y + 20);
  });
  requestAnimationFrame(drawGraph);
}

function buildPipe() { $("#pipe").innerHTML = STAGES.map(([n, nm]) => `<div class="stg" data-k="${n}"><div class="top"><span class="ix">${n}</span><span class="nm">${nm}</span><span class="fl">IDLE</span></div><div class="bar"><i></i></div></div>`).join(""); }
function setStg(n, st, fl) { const r = $(`.stg[data-k="${n}"]`); if (!r) return; r.className = "stg " + st; $(".fl", r).textContent = fl; $(".bar i", r).style.width = st === "idle" ? "0" : (st === "wait" ? "55%" : "100%"); }
function resetPipe() { STAGES.forEach(([n]) => setStg(n, "idle", "IDLE")); }
function pipeByK(k) {
  const seq = [["01", 0.05, 0.20], ["02", 0.22, 0.52], ["03", 0.54, 0.62], ["04", 0.62, 0.70]];
  seq.forEach(([n, s, e]) => { if (k >= e && _pk[n] !== "pass") { setStg(n, "pass", "PASS"); _pk[n] = "pass"; } else if (k >= s && !_pk[n]) { setStg(n, "run", "CHECK"); _pk[n] = "run"; } });
  if (k >= 0.72 && !_pk["05"]) { setStg("05", "wait", "WAIT"); _pk["05"] = "wait"; }
  if (k >= 0.95 && _pk["05"] !== "pass") { setStg("05", "pass", "PASS"); _pk["05"] = "pass"; }
}
function verdict(kind, sub, word) {
  const map = { idle: ["IDLE", '<circle cx="12" cy="12" r="9"/>'], run: ["WORKING", '<circle cx="12" cy="12" r="4"/><path d="M12 3v3M12 18v3M3 12h3M18 12h3" opacity=".55"/>'], accept: ["ACCEPT", '<path d="m5 12 5 5L20 7"/>'], reject: ["REJECT", '<path d="M18 6 6 18M6 6l12 12"/>'], inconc: ["INCONCLUSIVE", '<circle cx="12" cy="12" r="9"/><path d="M12 8v4.5M12 16h.01"/>'] };
  $("#verdict").className = "verdict " + kind; $("#vw").childNodes[0].nodeValue = word || map[kind][0]; $("#vs").textContent = sub || "";
  $("#vg").innerHTML = `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2">${map[kind][1]}</svg>`;
}
export function runVerify() {
  if (V.active || REV.active) return;
  if (live) return void runVerifyLive();
  V.active = true; V.live = false; V.done = false; V.t0 = performance.now(); _pk = {}; REV.reject = false; REV.seam = 0;
  $$("#runBtn,#revBtn").forEach((b) => (b.disabled = true)); resetPipe();
  verdict("run", "Signature traversing the trust path…", "VERIFYING"); $("#mLat").textContent = "—"; $("#mSnap").textContent = "—";
  if (reduce) { V.done = true; onVerified(); V.active = false; }
}

// runVerifyLive issues + verifies against the real server, then plays the
// cinematic and applies the REAL decision/trace/latency at the end.
async function runVerifyLive() {
  $$("#runBtn,#revBtn").forEach((b) => (b.disabled = true)); resetPipe();
  REV.reject = false; REV.seam = 0;
  verdict("run", "Issuing + verifying against atlas-server…", "LIVE");
  $("#mLat").textContent = "—"; $("#mSnap").textContent = "—";
  try {
    const iss = await live.issue({ principal: PRINCIPAL, delegate: DELEGATE, scope: ["read:orders", "write:audit"], ttlSeconds: 3600 });
    const res = await live.verify(iss.record);
    V.realRes = res; V.live = true; V.done = false; V.active = true; V.t0 = performance.now(); _pk = {};
    if (reduce) { V.done = true; finalizeLive(); V.active = false; }
  } catch (e) {
    verdict("reject", "atlas-server unreachable — " + e.message, "OFFLINE");
    $$("#runBtn,#revBtn").forEach((b) => (b.disabled = false));
  }
}
function finalizeLive() {
  const res = V.realRes || { decision: "inconclusive", causes: [], trace: [], latencyMicros: 0 };
  applyTraceToPipe(res.trace);
  const kind = decisionKind(res.decision);
  const sub = kind === "accept" ? "Verified by atlas-server · real ES256 record."
    : (res.causes && res.causes[0] ? res.causes[0] : "See trace.");
  verdict(kind, sub);
  $("#mLat").textContent = (res.latencyMicros || 0) + "µs";
  $("#mSnap").textContent = "2s old";
  $("#pipeSt") && ($("#pipeSt").textContent = kind === "accept" ? "5 / 5 pass" : "rejected");
  V.live = false;
  setTimeout(() => { V.active = false; $$("#runBtn,#revBtn").forEach((b) => (b.disabled = false)); }, reduce ? 10 : 700);
}
function onVerified() {
  verdict("accept", "All five checks passed · single‑hop · two‑domain."); $("#mSnap").textContent = "2s old";
  const el = $("#mLat"), t0 = performance.now();
  (function tk(n) { let k = Math.min(1, (n - t0) / 560); el.textContent = Math.round(94 * k) + "µs"; if (k < 1 && !reduce) requestAnimationFrame(tk); else el.textContent = "94µs"; })(performance.now());
  setTimeout(() => { V.active = false; $$("#runBtn,#revBtn").forEach((b) => (b.disabled = false)); }, reduce ? 10 : 700);
}
export function runRevocation() {
  if (V.active || REV.active) return;
  if (live) return void runRevocationLive();
  REV.active = true; REV.reject = false; REV.t0 = performance.now();
  $$("#runBtn,#revBtn").forEach((b) => (b.disabled = true)); resetPipe();
  if (reduce) { ["01", "02", "03", "04"].forEach((n) => setStg(n, "pass", "PASS")); setStg("05", "reject", "REVOKED"); verdict("reject", "Revocation observable after recovery · RevokedObservable."); $("#mLat").textContent = "91µs"; $("#mSnap").textContent = "fresh"; REV.reject = true; REV.active = false; $$("#runBtn,#revBtn").forEach((b) => (b.disabled = false)); return; }
  verdict("run", "Relying party severed from issuer — partition active…", "PARTITIONED"); $("#mLat").textContent = "—"; $("#mSnap").textContent = "—";
  ["01", "02", "03", "04"].forEach((n, i) => setTimeout(() => setStg(n, "pass", "PASS"), 250 + i * 170));
  setTimeout(() => setStg("05", "wait", "WAIT"), 1060);
  setTimeout(() => { verdict("inconc", "Revocation performed during the partition — not observable (S4). Failing closed."); $("#mLat").textContent = "88µs"; $("#mSnap").textContent = "stale"; }, 1600);
  setTimeout(() => verdict("run", "Partition recovered · fresh snapshot propagating…", "RECOVERING"), 2150);
  setTimeout(() => { setStg("05", "reject", "REVOKED"); verdict("reject", "Revocation now observable · RevokedObservable."); $("#mLat").textContent = "91µs"; $("#mSnap").textContent = "fresh"; }, 3100);
  setTimeout(() => { REV.active = false; $$("#runBtn,#revBtn").forEach((b) => (b.disabled = false)); }, 3900);
}

// runRevocationLive drives the real revocation alpha-path against the server:
// issue → verify (accept) → revoke → verify (reject/RevokedObservable). The
// partition seam plays as a visual metaphor; every verdict shown is real.
async function runRevocationLive() {
  $$("#runBtn,#revBtn").forEach((b) => (b.disabled = true)); resetPipe();
  try {
    verdict("run", "Issuing + verifying against atlas-server…", "LIVE");
    const iss = await live.issue({ principal: PRINCIPAL, delegate: DELEGATE, scope: ["read:orders"], ttlSeconds: 3600 });
    const v1 = await live.verify(iss.record);
    applyTraceToPipe(v1.trace);
    verdict(decisionKind(v1.decision), "Issued & verified — now revoking…");
    $("#mLat").textContent = (v1.latencyMicros || 0) + "µs"; $("#mSnap").textContent = "2s old";
    await sleep(reduce ? 0 : 650);
    // partition seam visual while the revocation propagates
    REV.active = true; REV.reject = false; REV.t0 = performance.now();
    verdict("run", "Revoking · publishing a fresh signed snapshot…", "REVOKING");
    await live.revoke(iss.instance);
    const v2 = await live.verify(iss.record);
    await sleep(reduce ? 0 : 900);
    applyTraceToPipe(v2.trace);
    const kind = decisionKind(v2.decision);
    verdict(kind, kind === "reject" ? ((v2.causes && v2.causes[0]) || "RevokedObservable") + " · real snapshot" : "See trace.");
    $("#mLat").textContent = (v2.latencyMicros || 0) + "µs"; $("#mSnap").textContent = "fresh";
    if (kind === "reject") REV.reject = true;
    await sleep(reduce ? 0 : 500);
  } catch (e) {
    verdict("reject", "atlas-server error — " + e.message, "OFFLINE");
  } finally {
    REV.active = false;
    $$("#runBtn,#revBtn").forEach((b) => (b.disabled = false));
  }
}

export function initVerify() {
  gcv = $("#graph"); gctx = gcv && gcv.getContext("2d");
  buildPipe(); verdict("idle", "Awaiting a presented delegation.");
  if (gcv) { sizeGraph(); requestAnimationFrame(drawGraph); addEventListener("resize", sizeGraph); }
  $("#runBtn").addEventListener("click", runVerify);
  $("#revBtn").addEventListener("click", runRevocation);
}
