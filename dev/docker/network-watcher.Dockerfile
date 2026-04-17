# BUILD IMAGE --------------------------------------------------------
ARG GO_VERSION=1.26
FROM golang:${GO_VERSION}-alpine AS builder

RUN apk add --no-cache build-base

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

ARG VERSION=unknown
RUN go build -ldflags="-X 'main.Version=$VERSION'" -o bin/network-watcher cmd/network-watcher/main.go

# ACTUAL IMAGE -------------------------------------------------------

FROM alpine:3.21

LABEL maintainer="eng@xmtp.com"
LABEL source="https://github.com/xmtp/xmtpd"
LABEL description="XMTP Network Watcher - Global Sync State Observability"

ENV GOLOG_LOG_FMT=nocolor

RUN apk add --no-cache curl ca-certificates

COPY --from=builder /app/bin/network-watcher /usr/bin/

EXPOSE 8009

HEALTHCHECK --interval=30s --timeout=5s --retries=3 \
  CMD curl -f http://localhost:8009/health || exit 1

ENTRYPOINT ["/usr/bin/network-watcher"]
