ARG GO_VERSION=none
ARG COMMIT=none
ARG CMD=xmtpd

FROM golang:${GO_VERSION}-alpine as builder
ARG COMMIT
ARG CMD
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -ldflags="-w -s -X 'main.GitCommit=${COMMIT}'" -o bin/${CMD} cmd/${CMD}/main.go

FROM alpine
ARG COMMIT
ARG CMD
ENV CMD=${CMD}
LABEL maintainer="eng@xmtp.com"
LABEL source="https://github.com/xmtp/xmtpd"
LABEL description="XMTP Node Software"
LABEL commit=$COMMIT
COPY --from=builder /app/bin/${CMD} /usr/bin/
ENTRYPOINT /usr/bin/${CMD}
