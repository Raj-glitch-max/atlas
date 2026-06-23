// Thin client for the Atlas Server (cmd/atlas-server). When the server is
// reachable the verify console runs against the real engine; otherwise the
// site falls back to the scripted demo. Override the base URL by setting
// localStorage.atlasApi.
const DEFAULT_BASE = "http://127.0.0.1:8087";

export function createApi(base) {
  base = base || (typeof localStorage !== "undefined" && localStorage.getItem("atlasApi")) || DEFAULT_BASE;

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
  };
}
