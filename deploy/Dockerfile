# Atlas product image — the server + CLI, statically built, on distroless.
# Multi-stage: build with the Go toolchain, ship a minimal static runtime.
FROM golang:1.22 AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags "-s -w" -o /out/atlas-server ./cmd/atlas-server && \
    CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags "-s -w" -o /out/atlas       ./cmd/atlas

# distroless/static: no shell, no package manager — small attack surface.
FROM gcr.io/distroless/static-debian12
COPY --from=build /out/atlas-server /usr/local/bin/atlas-server
COPY --from=build /out/atlas        /usr/local/bin/atlas
EXPOSE 8087
VOLUME ["/data"]
# Durable by default inside the container: state + key live on the /data volume.
ENTRYPOINT ["atlas-server"]
CMD ["-addr", "0.0.0.0:8087", "-store", "/data/state.json", "-key", "/data/authority.key"]
