# BUILD IMAGE --------------------------------------------------------
ARG GO_VERSION=1.26
FROM golang:${GO_VERSION}-alpine AS builder

# Get build tools and required header files
RUN apk add --no-cache build-base

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

ARG VERSION=unknown
RUN go build -ldflags="-X 'main.Version=$VERSION'" -o bin/xmtpd-e2e cmd/xmtpd-e2e/main.go

# ACTUAL IMAGE -------------------------------------------------------

FROM alpine:3.21

LABEL maintainer="eng@xmtp.com"
LABEL source="https://github.com/xmtp/xmtpd"
LABEL description="XMTP E2E Test Runner"

# color, nocolor, json
ENV GOLOG_LOG_FMT=nocolor

RUN apk add --no-cache curl

COPY --from=builder /app/bin/xmtpd-e2e /usr/bin/

ENTRYPOINT ["/usr/bin/xmtpd-e2e"]
