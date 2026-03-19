# BUILD IMAGE --------------------------------------------------------
ARG GO_VERSION=1.26
FROM golang:${GO_VERSION}-alpine AS builder

RUN apk add --no-cache build-base

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go build -o bin/chain-watcher cmd/chain-watcher/main.go

# ACTUAL IMAGE -------------------------------------------------------
FROM alpine:3.21

LABEL maintainer="eng@xmtp.com"
LABEL source="https://github.com/xmtp/xmtpd"
LABEL description="XMTP Chain Watcher - Payer Report Observability"

RUN apk add --no-cache ca-certificates curl

COPY --from=builder /app/bin/chain-watcher /usr/bin/

EXPOSE 8008

HEALTHCHECK --interval=30s --timeout=5s --retries=3 \
  CMD curl -f http://localhost:8008/health || exit 1

ENTRYPOINT ["/usr/bin/chain-watcher"]
