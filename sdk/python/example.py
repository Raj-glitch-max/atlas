"""Runnable example: the Atlas capability lifecycle via the Python SDK.

    go run ./cmd/atlas-server        # in one terminal
    python3 sdk/python/example.py    # in another
"""
from atlas import AtlasClient, AtlasError, Decision


def main() -> None:
    atlas = AtlasClient("http://127.0.0.1:8087")
    if not atlas.health():
        raise SystemExit("atlas-server not reachable — run: go run ./cmd/atlas-server")

    principal = "spiffe://domain-a.test/workload/payments-api"
    delegate = "spiffe://domain-b.test/agent/booking-worker"

    print("1. issue a scoped capability (single hop)")
    grant = atlas.issue(principal, delegate, ["read:orders", "write:audit"])
    print(f"   instance={grant.instance}")

    print("2. attenuation — over-scope is refused")
    try:
        atlas.issue(principal, delegate, ["delete:everything"])
        print("   ERROR: over-scope was NOT refused")
    except AtlasError as e:
        print(f"   refused as expected: {e}")

    print("3. verify offline")
    r = atlas.verify(grant.record)
    print(f"   {r.decision.value.upper()} ({r.latency_micros}µs)")
    for t in r.trace:
        print(f"     {t.check:<18} {t.outcome} {t.cause if t.cause != 'None' else ''}")

    print("4. revoke")
    atlas.revoke(grant.instance)

    print("5. verify again — now revoked")
    r = atlas.verify(grant.record)
    print(f"   {r.decision.value.upper()} ({r.latency_micros}µs) {r.causes}")

    assert r.decision == Decision.REJECT, "expected rejection after revocation"
    print("\nOK — single-hop grant, attenuation, offline verify, independent revocation.")


if __name__ == "__main__":
    main()
