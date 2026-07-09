/**
 * Atlas TypeScript SDK — a zero-dependency client for the Atlas Server.
 *
 * Atlas issues and verifies offline-verifiable, attenuable delegation tokens
 * bound to SPIFFE workload identity. This client talks to `cmd/atlas-server`
 * over HTTP using only the platform `fetch` (Node 18+, Deno, Bun, browsers).
 *
 * ```ts
 * import { AtlasClient } from "./atlas.ts";
 *
 * const atlas = new AtlasClient("http://127.0.0.1:8087");
 * const grant = await atlas.issue({
 *   principal: "spiffe://domain-a.test/workload/payments-api",
 *   delegate:  "spiffe://domain-b.test/agent/booking-worker",
 *   scope:     ["read:orders", "write:audit"],
 *   ttlSeconds: 3600,
 * });
 * const result = await atlas.verify(grant.record);
 * console.log(result.decision); // "accept"
 *
 * await atlas.revoke(grant.instance);
 * (await atlas.verify(grant.record)).decision; // "reject"
 * ```
 *
 * The API surface mirrors the Python SDK (`sdk/python/atlas.py`) one-for-one.
 */

/** A verification verdict. Deterministic: reconstructable from the record + trace. */
export type Decision = "accept" | "reject" | "inconclusive";
export const Decision = {
  Accept: "accept",
  Reject: "reject",
  Inconclusive: "inconclusive",
} as const;

/** Raised when the server returns an error or is unreachable. */
export class AtlasError extends Error {
  status: number;
  constructor(message: string, status = 0) {
    super(message);
    this.name = "AtlasError";
    this.status = status;
  }
}

export interface Grant {
  record: string;
  instance: string;
  principal: string;
  delegate: string;
  scope: string[];
  issuedAt: string;
  expiresAt: string;
}

export interface TraceEntry {
  check: string;
  outcome: string;
  cause: string;
  detail?: string;
}

export interface VerifyResult {
  decision: Decision;
  accept: boolean;
  causes: string[];
  latencyMicros: number;
  trace: TraceEntry[];
}

export interface IssueParams {
  principal: string;
  delegate: string;
  scope: string[];
  ttlSeconds?: number;
}

export interface AtlasClientOptions {
  apiKey?: string;
  /** Request timeout in milliseconds (default 8000). */
  timeoutMs?: number;
}

/** A thin, promise-based HTTP client for the Atlas Server. */
export class AtlasClient {
  private base: string;
  private apiKey: string | undefined;
  private timeoutMs: number;

  constructor(baseUrl = "http://127.0.0.1:8087", opts: AtlasClientOptions = {}) {
    this.base = baseUrl.replace(/\/+$/, "");
    this.apiKey = opts.apiKey;
    this.timeoutMs = opts.timeoutMs ?? 8000;
  }

  // -- core --

  private async call<T>(method: string, path: string, body?: unknown): Promise<T> {
    const headers: Record<string, string> = { "Content-Type": "application/json" };
    if (this.apiKey) headers["Authorization"] = "Bearer " + this.apiKey;

    const ctrl = new AbortController();
    const timer = setTimeout(() => ctrl.abort(), this.timeoutMs);
    let resp: Response;
    try {
      resp = await fetch(this.base + path, {
        method,
        headers,
        body: body === undefined ? undefined : JSON.stringify(body),
        signal: ctrl.signal,
      });
    } catch (e) {
      const reason = e instanceof Error ? e.message : String(e);
      throw new AtlasError(`atlas-server unreachable at ${this.base} (${reason})`);
    } finally {
      clearTimeout(timer);
    }

    const raw = await resp.text();
    if (!resp.ok) {
      let msg = raw;
      try {
        msg = (JSON.parse(raw) as { error?: string }).error ?? raw;
      } catch {
        /* non-JSON body: use raw text */
      }
      throw new AtlasError(`${resp.status}: ${msg}`, resp.status);
    }
    return (raw ? JSON.parse(raw) : {}) as T;
  }

  // -- operations (mirror sdk/python/atlas.py) --

  /** True when the server is reachable and healthy; never throws. */
  async health(): Promise<boolean> {
    try {
      const r = await this.call<{ status?: string }>("GET", "/health");
      return r.status === "ok";
    } catch {
      return false;
    }
  }

  /** Orchestrator readiness: true only when a fresh revocation snapshot is held. */
  async ready(): Promise<boolean> {
    try {
      const r = await this.call<{ ready?: boolean }>("GET", "/readyz");
      return r.ready === true;
    } catch {
      return false;
    }
  }

  version(): Promise<Record<string, unknown>> {
    return this.call("GET", "/version");
  }

  issue(params: IssueParams): Promise<Grant> {
    return this.call<Grant>("POST", "/issue", {
      principal: params.principal,
      delegate: params.delegate,
      scope: params.scope,
      ttlSeconds: params.ttlSeconds ?? 3600,
    });
  }

  verify(record: string): Promise<VerifyResult> {
    return this.call<VerifyResult>("POST", "/verify", { record });
  }

  async revoke(instance: string): Promise<void> {
    await this.call("POST", "/revoke", { instance });
  }

  async delegations(): Promise<Record<string, unknown>[]> {
    const r = await this.call<{ delegations?: Record<string, unknown>[] }>("GET", "/delegations");
    return r.delegations ?? [];
  }

  async audit(limit = 50): Promise<Record<string, unknown>[]> {
    const r = await this.call<{ events?: Record<string, unknown>[] }>("GET", `/audit?limit=${limit}`);
    return r.events ?? [];
  }

  graph(): Promise<Record<string, unknown>> {
    return this.call("GET", "/graph");
  }

  stats(): Promise<Record<string, unknown>> {
    return this.call("GET", "/stats");
  }
}
