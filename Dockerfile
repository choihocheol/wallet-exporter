FROM golang:1.22.1-bookworm AS build

LABEL org.opencontainers.image.source="https://github.com/choihocheol/wallet-exporter"

RUN apt-get update && \
    apt-get install -y --no-install-recommends make git && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /build

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY . .
RUN --mount=type=cache,target=/root/.cache/go-build \
    make build

FROM debian:bookworm-slim AS deploy

RUN groupadd -r wallet-exporter && \
    useradd -r -g wallet-exporter wallet-exporter
RUN apt-get update && \
    apt-get install -y --no-install-recommends ca-certificates && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=build /build/wallet-exporter /usr/local/bin/wallet-exporter

USER wallet-exporter
ENTRYPOINT ["/usr/local/bin/wallet-exporter"]
