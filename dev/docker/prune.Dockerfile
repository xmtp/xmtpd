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
RUN go build -ldflags="-X 'main.Version=$VERSION'" -o bin/xmtpd-prune cmd/prune/main.go

# ACTUAL IMAGE -------------------------------------------------------

FROM alpine:3.21

LABEL maintainer="eng@ephemerahq.com"
LABEL source="https://github.com/xmtp/xmtpd"
LABEL description="XMTPD Prune"

# color, nocolor, json
ENV GOLOG_LOG_FMT=nocolor

ENV XMTPD_LOG_ENCODING=json

COPY --from=builder /app/bin/xmtpd-prune /usr/bin/

ENTRYPOINT ["/usr/bin/xmtpd-prune"]
