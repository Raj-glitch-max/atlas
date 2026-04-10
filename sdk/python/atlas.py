"""Atlas Python SDK — a zero-dependency client for the Atlas Server.

Atlas issues and verifies offline-verifiable, attenuable delegation tokens
bound to SPIFFE workload identity. This client talks to `cmd/atlas-server`
over HTTP using only the standard library.

    from atlas import AtlasClient, Decision

    atlas = AtlasClient("http://127.0.0.1:8087", api_key=None)

    grant = atlas.issue(
        principal="spiffe://domain-a.test/workload/payments-api",
        delegate="spiffe://domain-b.test/agent/booking-worker",
        scope=["read:orders", "write:audit"],
        ttl_seconds=3600,
    )
    result = atlas.verify(grant.record)
    assert result.decision == Decision.ACCEPT

    atlas.revoke(grant.instance)
    assert atlas.verify(grant.record).decision == Decision.REJECT
"""
from __future__ import annotations

import json
import urllib.error
import urllib.request
from dataclasses import dataclass, field
from enum import Enum
from typing import Any


class Decision(str, Enum):
    ACCEPT = "accept"
    REJECT = "reject"
    INCONCLUSIVE = "inconclusive"


class AtlasError(Exception):
    """Raised when the server returns an error or is unreachable."""


@dataclass
class Grant:
    record: str
    instance: str
    principal: str
    delegate: str
    scope: list[str]
    expires_at: str

    @classmethod
    def _from(cls, d: dict[str, Any]) -> "Grant":
        return cls(
            record=d["record"], instance=d["instance"],
            principal=d.get("principal", ""), delegate=d.get("delegate", ""),
            scope=d.get("scope", []), expires_at=d.get("expiresAt", ""),
        )


@dataclass
class TraceEntry:
    check: str
    outcome: str
    cause: str


@dataclass
class VerifyResult:
    decision: Decision
    accept: bool
    causes: list[str]
    latency_micros: int
    trace: list[TraceEntry] = field(default_factory=list)

    @classmethod
    def _from(cls, d: dict[str, Any]) -> "VerifyResult":
        return cls(
            decision=Decision(d["decision"]) if d["decision"] in Decision._value2member_map_ else d["decision"],
            accept=bool(d.get("accept")),
            causes=d.get("causes") or [],
            latency_micros=int(d.get("latencyMicros", 0)),
            trace=[TraceEntry(t.get("check", ""), t.get("outcome", ""), t.get("cause", "")) for t in d.get("trace") or []],
        )


class AtlasClient:
    """A thin HTTP client for the Atlas Server."""

    def __init__(self, base_url: str = "http://127.0.0.1:8087",
                 api_key: str | None = None, timeout: float = 8.0) -> None:
        self.base = base_url.rstrip("/")
        self.api_key = api_key
        self.timeout = timeout

    # -- core --

    def _call(self, method: str, path: str, body: dict | None = None) -> Any:
        data = json.dumps(body).encode() if body is not None else None
        req = urllib.request.Request(self.base + path, data=data, method=method)
        req.add_header("Content-Type", "application/json")
        if self.api_key:
            req.add_header("Authorization", "Bearer " + self.api_key)
        try:
            with urllib.request.urlopen(req, timeout=self.timeout) as resp:
                raw = resp.read()
        except urllib.error.HTTPError as e:
            raw = e.read()
            try:
                msg = json.loads(raw).get("error", raw.decode())
            except Exception:
                msg = raw.decode(errors="replace")
            raise AtlasError(f"{e.code}: {msg}") from None
        except urllib.error.URLError as e:
            raise AtlasError(f"atlas-server unreachable at {self.base} ({e.reason})") from None
        return json.loads(raw) if raw else {}

    # -- operations --

    def health(self) -> bool:
        try:
            return self._call("GET", "/health").get("status") == "ok"
        except AtlasError:
            return False

    def version(self) -> dict:
        return self._call("GET", "/version")

    def issue(self, principal: str, delegate: str, scope: list[str],
              ttl_seconds: int = 3600) -> Grant:
        return Grant._from(self._call("POST", "/issue", {
            "principal": principal, "delegate": delegate,
            "scope": scope, "ttlSeconds": ttl_seconds,
        }))

    def verify(self, record: str) -> VerifyResult:
        return VerifyResult._from(self._call("POST", "/verify", {"record": record}))

    def revoke(self, instance: str) -> None:
        self._call("POST", "/revoke", {"instance": instance})

    def delegations(self) -> list[dict]:
        return self._call("GET", "/delegations").get("delegations", [])

    def audit(self, limit: int = 50) -> list[dict]:
        return self._call("GET", f"/audit?limit={limit}").get("events", [])

    def graph(self) -> dict:
        return self._call("GET", "/graph")

    def stats(self) -> dict:
        return self._call("GET", "/stats")


__all__ = ["AtlasClient", "Grant", "VerifyResult", "TraceEntry", "Decision", "AtlasError"]
