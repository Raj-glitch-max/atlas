"""Self-contained SDK smoke test: builds atlas-server, runs it on a temp store,
exercises the full lifecycle through the SDK, and tears down.

    python3 sdk/python/smoke_test.py
"""
import os
import pathlib
import shutil
import socket
import subprocess
import tempfile
import time

from atlas import AtlasClient, AtlasError, Decision


def free_port() -> int:
    s = socket.socket()
    s.bind(("127.0.0.1", 0))
    port = s.getsockname()[1]
    s.close()
    return port


def main() -> None:
    root = pathlib.Path(__file__).resolve().parents[2]
    workdir = tempfile.mkdtemp(prefix="atlas-sdk-")
    binp = os.path.join(workdir, "atlas-server")
    subprocess.run(["go", "build", "-o", binp, "./cmd/atlas-server"], cwd=root, check=True)

    port = free_port()
    proc = subprocess.Popen(
        [binp, "-addr", f"127.0.0.1:{port}",
         "-store", f"{workdir}/state.json", "-key", f"{workdir}/authority.key"],
        stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL,
    )
    try:
        atlas = AtlasClient(f"http://127.0.0.1:{port}")
        for _ in range(60):
            if atlas.health():
                break
            time.sleep(0.1)
        else:
            raise SystemExit("server did not become healthy")

        p, d = "spiffe://domain-a.test/workload/api", "spiffe://domain-b.test/agent/worker"
        grant = atlas.issue(p, d, ["read:orders", "write:audit"])
        assert atlas.verify(grant.record).decision == Decision.ACCEPT, "fresh grant should accept"

        try:
            atlas.issue(p, d, ["delete:everything"])
            raise AssertionError("over-scope should have been refused")
        except AtlasError:
            pass  # expected

        atlas.revoke(grant.instance)
        r = atlas.verify(grant.record)
        assert r.decision == Decision.REJECT, f"revoked grant should reject, got {r.decision}"
        assert "RevokedObservable" in r.causes, r.causes

        assert atlas.stats()["issued"] >= 1
        assert len(atlas.delegations()) >= 1
        assert atlas.graph()["edges"], "graph should have edges"

        print("PASS: atlas python sdk smoke (issue, attenuation, offline verify, revoke)")
    finally:
        proc.terminate()
        proc.wait()
        shutil.rmtree(workdir, ignore_errors=True)


if __name__ == "__main__":
    main()
