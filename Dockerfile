# syntax = docker/dockerfile:1.0-experimental

FROM golang:1.12-alpine AS builder
RUN apk add -U git
WORKDIR /work
ENV GO111MODULE=on
COPY cmd cmd
COPY *.go go.* ./
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache \
	CGO_ENABLED=0 go build --ldflags '-extldflags "-static"' -o /buildkit-ninja ./cmd/buildkit-ninja

FROM scratch
COPY --from=builder /buildkit-ninja /buildkit-ninja
ENTRYPOINT ["/buildkit-ninja"]
