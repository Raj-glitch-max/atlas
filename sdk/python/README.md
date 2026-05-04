# Atlas Python SDK

A zero-dependency (stdlib-only) client for the Atlas Server — integrate
offline-verifiable delegation into a Python app in a few lines.

## Use

```python
from atlas import AtlasClient, Decision

atlas = AtlasClient("http://127.0.0.1:8087")   # api_key=... if the server is guarded

grant = atlas.issue(
    principal="spiffe://domain-a.test/workload/payments-api",
    delegate="spiffe://domain-b.test/agent/booking-worker",
    scope=["read:orders", "write:audit"],
    ttl_seconds=3600,
)

result = atlas.verify(grant.record)
if result.decision is Decision.ACCEPT:
    ...  # honor the capability

atlas.revoke(grant.instance)
assert atlas.verify(grant.record).decision is Decision.REJECT
```

Methods: `health · version · issue · verify · revoke · delegations · audit ·
graph · stats`. Errors raise `AtlasError` (unreachable server, refused
issuance, etc.).

## Run the example / test

```bash
go run ./cmd/atlas-server            # terminal 1
python3 sdk/python/example.py        # terminal 2

# or fully self-contained (builds + runs its own server):
python3 sdk/python/smoke_test.py
```

## Notes

- Single file, no pip install, no third-party deps (uses `urllib`).
- Atlas is **single-hop** by design — a principal → delegate grant verifiable
  across trust domains, not multi-hop re-delegation.

<!-- checkpoint: feat(sdk): implement conformance validation (#114) -->
