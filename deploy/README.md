# deploy — the Atlas product stack

Containerized Atlas Server (durable), with Prometheus scraping its real
`/metrics` and Grafana for dashboards. This is the **deployable product** —
distinct from `atlas-lab/` (the SPIRE experiment substrate).

## Run

```bash
# from repo root
docker compose -f deploy/docker-compose.yml up -d --build

# server API
curl localhost:8087/health
# metrics (scraped by Prometheus)
curl localhost:8087/metrics
# Prometheus → http://localhost:9090   Grafana → http://localhost:3000
```

Guard mutating endpoints with a token:

```bash
ATLAS_API_KEY=$(openssl rand -hex 16) docker compose -f deploy/docker-compose.yml up -d
# then: atlas --api-key $ATLAS_API_KEY issue …
```

## Image

`deploy/Dockerfile` is a multi-stage build: the Go toolchain compiles static
`atlas-server` + `atlas` binaries (`-trimpath -ldflags "-s -w"`, `CGO_ENABLED=0`),
shipped on `gcr.io/distroless/static-debian12` — no shell, no package manager,
small attack surface. State and the authority key persist on the `atlas-data`
volume, so **records survive restarts**.

Validate the compose without a daemon:

```bash
docker compose -f deploy/docker-compose.yml config >/dev/null && echo valid
```

## Not included (follow-ups)

- Non-root/read-only-rootfs hardening (distroless `:nonroot` + a writable
  `/data` owned by the runtime user).
- TLS termination (front with a reverse proxy, or add server TLS).
- Latency histogram buckets in `/metrics` (currently counters + snapshot-age
  gauge).
