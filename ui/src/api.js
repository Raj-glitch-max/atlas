// Thin client for the Atlas Server (cmd/atlas-server). When the server is
// reachable the verify console runs against the real engine; otherwise the
// site falls back to the scripted demo.
//
// Base URL resolution, in priority order:
//   1. an explicit argument
//   2. localStorage.atlasApi          — lets anyone repoint a deployed build
//   3. import.meta.env.VITE_ATLAS_API — baked in at build time (the deployed default)
//   4. http://127.0.0.1:8087          — local development
//
// The build-time variable is what makes a deployed page talk to a deployed
// backend; without it a public visitor's browser would try to reach their own
// localhost and always fall back to the scripted demo.
const BUILD_BASE =
  (typeof import.meta !== "undefined" && import.meta.env && import.meta.env.VITE_ATLAS_API) || "";
const DEFAULT_BASE = BUILD_BASE || "http://127.0.0.1:8087";

// Session isolation: each visitor gets their own capability graph, audit log
// and revocation set, so nobody can see or revoke anyone else's state. The id
// is stored locally so a page reload keeps the same sandbox.
const SESSION_KEY = "atlasSession";
const SESSION_HEADER = "X-Atlas-Session";

function hex32() {
  const b = new Uint8Array(16);
  (globalThis.crypto || {}).getRandomValues
    ? globalThis.crypto.getRandomValues(b)
    : b.forEach((_, i) => (b[i] = Math.floor(Math.random() * 256)));
  return [...b].map((x) => x.toString(16).padStart(2, "0")).join("");
}

function sessionId() {
  if (typeof localStorage === "undefined") return hex32();
  let id = localStorage.getItem(SESSION_KEY);
  if (!id || !/^[0-9a-f]{32}$/.test(id)) {
    id = hex32();
    localStorage.setItem(SESSION_KEY, id);
  }
  return id;
}

export function createApi(base) {
  base =
    base ||
    (typeof localStorage !== "undefined" && localStorage.getItem("atlasApi")) ||
    DEFAULT_BASE;

  const session = sessionId();

  async function j(path, opts) {
    const ctrl = new AbortController();
    const to = setTimeout(() => ctrl.abort(), 4000);
    const key = typeof localStorage !== "undefined" && localStorage.getItem("atlasApiKey");
    try {
      const r = await fetch(base + path, {
        ...opts,
        signal: ctrl.signal,
        headers: {
          "Content-Type": "application/json",
          [SESSION_HEADER]: session,
          ...(key ? { Authorization: "Bearer " + key } : {}),
          ...(opts && opts.headers),
        },
      });
      const data = await r.json().catch(() => ({}));
      if (!r.ok) throw new Error(data.error || "HTTP " + r.status);
      return data;
    } finally {
      clearTimeout(to);
    }
  }

  return {
    base,
    session,
    async health() {
      try {
        const d = await j("/health");
        return { ok: d.status === "ok" };
      } catch {
        return { ok: false };
      }
    },
    version() { return j("/version"); },
    issue(body) { return j("/issue", { method: "POST", body: JSON.stringify(body) }); },
    verify(record) { return j("/verify", { method: "POST", body: JSON.stringify({ record }) }); },
    revoke(instance) { return j("/revoke", { method: "POST", body: JSON.stringify({ instance }) }); },
    delegations() { return j("/delegations"); },
    audit(limit = 40) { return j("/audit?limit=" + limit); },
    graph() { return j("/graph"); },
    stats() { return j("/stats"); },
    // Aggregated, anonymized, cross-session activity — counts only.
    activity() { return j("/activity"); },
  };
}
