# BUILD IMAGE --------------------------------------------------------
ARG GO_VERSION=1.22
FROM golang:${GO_VERSION}-bookworm as builder

WORKDIR /app

SHELL ["/bin/bash", "-c"]

RUN apt-get update &&  \
    apt-get install -y curl git jq

RUN curl -L https://foundry.paradigm.xyz | bash && \
    source ~/.bashrc && \
    foundryup &&  \
    cp ~/.foundry/bin/* /usr/local/bin

COPY . .

RUN dev/docker/anvil-background && \
    dev/contracts/deploy-local && \
    dev/register-local-node && \
    dev/register-local-node-2 && \
    sleep 5

# ACTUAL IMAGE -------------------------------------------------------

FROM debian:bookworm

LABEL maintainer="engineering@xmtp.com"
LABEL source="https://github.com/xmtp/xmtpd"
LABEL description="XMTP Node Software"

# color, nocolor, json
ENV GOLOG_LOG_FMT=nocolor

EXPOSE 5050

COPY --from=builder /usr/local/bin/anvil /usr/bin/
COPY --from=builder /app/anvil-baked-state anvil-baked-state

CMD ["anvil", "--state", "anvil-baked-state"]