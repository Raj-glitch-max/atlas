// Atlas operator console — a live view of trust state (metrics, delegations,
// audit, graph) against the real atlas-server, with issue + revoke actions and
// auto-refresh. Degrades to a "start the server" panel when offline.
import "./styles.css";
import { createApi } from "./api.js";

const $ = (s, r = document) => r.querySelector(s);
const api = createApi();
let timer = null;

function setConn(ok, stats) {
  const el = $("#conn");
  if (ok) {
    el.innerHTML = `<span class="dot ok"></span>live · ${stats ? stats.trustDomain : "atlas-server"}`;
    $("#offline").hidden = true;
    $("#livewrap").hidden = false;
  } else {
    el.innerHTML = `<span class="dot warn"></span>offline`;
    $("#offline").hidden = false;
    $("#livewrap").hidden = true;
  }
}

function renderStats(s) {
  $("#td").textContent = s.trustDomain;
  $("#rr").textContent = s.revocationR;
  const cards = [
    ["Delegations", s.delegations, ""],
    ["Verifications", s.verified, `${s.accept} accept · ${s.reject} reject`],
    ["Revocations", s.revoked, ""],
    ["Snapshot age", (s.snapshotAgeSecond || 0).toFixed(1), "sec · bound by R"],
  ];
  $("#stats").innerHTML = cards.map(([l, v, d]) =>
    `<div class="metric"><div class="l">${l}</div><div class="v">${v}</div><div class="d">${d}</div></div>`).join("");
}

function pill(revoked) {
  return revoked
    ? `<span class="pill reject"><span class="pd"></span>REVOKED</span>`
    : `<span class="pill accept"><span class="pd"></span>ACTIVE</span>`;
}
function leaf(id) { const i = id.lastIndexOf("/"); return i >= 0 ? id.slice(i + 1) : id; }
function short(s) { return s.length > 14 ? s.slice(0, 14) + "…" : s; }

function renderDelegations(ds) {
  if (!ds || !ds.length) { $("#delBody").innerHTML = `<tr><td colspan="6" class="sub" style="padding:18px">no delegations yet — issue one above</td></tr>`; return; }
  $("#delBody").innerHTML = ds.map((d) => `<tr>
    <td class="mono">${short(d.instance)}</td>
    <td class="mono sub">${leaf(d.principal)}</td>
    <td class="mono">${leaf(d.delegate)}</td>
    <td class="sub">${(d.scope || []).join(", ")}</td>
    <td>${pill(d.revoked)}</td>
    <td>${d.revoked ? "" : `<button class="revbtn" data-inst="${d.instance}" title="revoke">✗</button>`}</td>
  </tr>`).join("");
  $$(".revbtn", $("#delBody")).forEach((b) => b.addEventListener("click", async () => {
    b.disabled = true;
    try { await api.revoke(b.dataset.inst); } catch (e) { flash("revoke failed: " + e.message, true); }
    tick();
  }));
}

function renderAudit(evs) {
  $("#auditCount").textContent = evs ? evs.length + " events" : "";
  $("#audit").innerHTML = (evs || []).map((e) => {
    const subject = e.delegate ? leaf(e.delegate) : (e.instance ? short(e.instance) : "");
    const detail = e.decision || e.detail || "";
    const cls = e.decision === "accept" ? "good" : (e.decision === "reject" ? "bad" : "");
    const t = (e.time || "").slice(11, 19);
    return `<div class="arow"><span class="at">${t}</span><span class="atype mono">${e.type}</span><span class="asub mono">${subject}</span><span class="adet ${cls}">${detail}</span></div>`;
  }).join("") || `<div class="sub" style="padding:14px">no events yet</div>`;
}

function renderGraph(g) {
  const svg = $("#cgraph");
  const W = svg.clientWidth || 800, H = 240;
  svg.setAttribute("viewBox", `0 0 ${W} ${H}`);
  const NS = "http://www.w3.org/2000/svg";
  svg.innerHTML = "";
  if (!g.edges || !g.edges.length) { svg.innerHTML = `<text x="${W / 2}" y="${H / 2}" text-anchor="middle" fill="#5A626F" font-family="Geist Mono,monospace" font-size="12">empty graph</text>`; return; }
  const principals = [...new Set(g.edges.map((e) => e.from))];
  const delegates = [...new Set(g.edges.map((e) => e.to))];
  const posL = {}, posR = {};
  principals.forEach((p, i) => posL[p] = { x: W * 0.18, y: (H / (principals.length + 1)) * (i + 1) });
  delegates.forEach((d, i) => posR[d] = { x: W * 0.82, y: (H / (delegates.length + 1)) * (i + 1) });
  const el = (t, a) => { const e = document.createElementNS(NS, t); for (const k in a) e.setAttribute(k, a[k]); return e; };
  g.edges.forEach((e) => {
    const a = posL[e.from], b = posR[e.to];
    const col = e.revoked ? "#FB7185" : "#7DD3FC";
    const path = el("path", { d: `M ${a.x} ${a.y} C ${(a.x + b.x) / 2} ${a.y}, ${(a.x + b.x) / 2} ${b.y}, ${b.x} ${b.y}`, fill: "none", stroke: col, "stroke-width": 1.4, "stroke-opacity": e.revoked ? .5 : .7, "stroke-dasharray": e.revoked ? "4 4" : "none" });
    svg.appendChild(path);
  });
  const node = (p, label, tint) => {
    svg.appendChild(el("circle", { cx: p.x, cy: p.y, r: 5, fill: tint }));
    const t = el("text", { x: p.x, y: p.y - 12, "text-anchor": "middle", fill: "#C3CAD6", "font-family": "Geist Mono,monospace", "font-size": 11 }); t.textContent = label; svg.appendChild(t);
  };
  principals.forEach((p) => node(posL[p], leaf(p), "#7DD3FC"));
  delegates.forEach((d) => node(posR[d], leaf(d), "#AEB9CC"));
}

function flash(msg, err) {
  const el = $("#issueMsg"); el.textContent = msg; el.className = "form-msg" + (err ? " err" : " ok");
  setTimeout(() => { el.textContent = ""; }, 4000);
}

async function tick() {
  try {
    const [stats, dels, aud, g] = await Promise.all([api.stats(), api.delegations(), api.audit(30), api.graph()]);
    setConn(true, stats);
    renderStats(stats); renderDelegations(dels.delegations); renderAudit(aud.events); renderGraph(g);
  } catch {
    setConn(false);
  }
}

function $$(s, r = document) { return [...r.querySelectorAll(s)]; }

$("#issueForm").addEventListener("submit", async (e) => {
  e.preventDefault();
  const f = new FormData(e.target);
  const scope = String(f.get("scope")).split(",").map((s) => s.trim()).filter(Boolean);
  try {
    const r = await api.issue({ principal: f.get("principal"), delegate: f.get("delegate"), scope, ttlSeconds: Number(f.get("ttl")) });
    flash("issued " + r.instance);
  } catch (err) {
    flash(err.message, true);
  }
  tick();
});

tick();
timer = setInterval(tick, 3000);
document.addEventListener("visibilitychange", () => {
  if (document.hidden) clearInterval(timer);
  else { tick(); timer = setInterval(tick, 3000); }
});
