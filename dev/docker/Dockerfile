# BUILD IMAGE --------------------------------------------------------
ARG GO_VERSION=1.24
FROM golang:${GO_VERSION}-alpine AS builder

# Get build tools and required header files
RUN apk add --no-cache build-base

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

# Build the final node binary
ARG VERSION=unknown
RUN go build -ldflags="-X 'main.Version=$VERSION'" -o bin/xmtpd cmd/replication/main.go

# ACTUAL IMAGE -------------------------------------------------------

FROM alpine:3.21

LABEL maintainer="eng@ephemerahq.com"
LABEL source="https://github.com/xmtp/xmtpd"
LABEL description="XMTP Node Software"

# color, nocolor, json
ENV GOLOG_LOG_FMT=nocolor

EXPOSE 5050

RUN apk add --no-cache curl

COPY --from=builder /app/bin/xmtpd /usr/bin/

ENTRYPOINT ["/usr/bin/xmtpd"]
