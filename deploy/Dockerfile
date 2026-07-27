# Atlas product image — the server + CLI, statically built, on distroless.
# Multi-stage: build with the Go toolchain, ship a minimal static runtime.
FROM golang:1.22 AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags "-s -w" -o /out/atlas-server ./cmd/atlas-server && \
    CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags "-s -w" -o /out/atlas       ./cmd/atlas && \
    mkdir -p /data && chown 65532:65532 /data

# distroless/static:nonroot — no shell, no package manager, and the process
# runs as an unprivileged user (uid 65532) by default: small attack surface,
# and it works under a read-only root filesystem (only /data must be writable,
# which is a mounted volume).
FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/atlas-server /usr/local/bin/atlas-server
COPY --from=build /out/atlas        /usr/local/bin/atlas
# /data is pre-created owned by the nonroot user so durable state + key are
# writable without root, even under a read-only root filesystem.
COPY --from=build --chown=65532:65532 /data /data
USER 65532:65532
EXPOSE 8087
VOLUME ["/data"]

# PORT (not -addr) selects the listen port, so the same image runs unchanged on
# a PaaS that assigns a port at runtime (Railway, Render, Cloud Run all inject
# PORT). The server binds 0.0.0.0 whenever PORT is set; this default keeps local
# `docker run -p 8087:8087` working. Passing -addr explicitly would WIN over
# PORT and silently make the container unreachable on those platforms.
ENV PORT=8087

# Durable by default inside the container: state + key live on the /data volume.
ENTRYPOINT ["atlas-server"]
CMD ["-store", "/data/state.json", "-key", "/data/authority.key"]
