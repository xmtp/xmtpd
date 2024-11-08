# BUILD IMAGE --------------------------------------------------------
ARG GO_VERSION=1.22
FROM golang:${GO_VERSION}-alpine AS builder

# Get build tools and required header files
RUN apk add --no-cache build-base

WORKDIR /app
COPY . .

# Build the final node binary
ARG GIT_COMMIT=unknown
RUN go build -ldflags="-X 'main.Commit=$GIT_COMMIT'" -o bin/xmtpd cmd/replication/main.go

# ACTUAL IMAGE -------------------------------------------------------

FROM alpine:3.12

LABEL maintainer="engineering@xmtp.com"
LABEL source="https://github.com/xmtp/xmtpd"
LABEL description="XMTP Node Software"

# color, nocolor, json
ENV GOLOG_LOG_FMT=nocolor

EXPOSE 5050

COPY --from=builder /app/bin/xmtpd /usr/bin/

ENTRYPOINT ["/usr/bin/xmtpd"]
