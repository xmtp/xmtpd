# BUILD IMAGE --------------------------------------------------------
ARG GO_VERSION=1.22
FROM golang:${GO_VERSION}-bookworm as builder

WORKDIR /app

SHELL ["/bin/bash", "-c"]

RUN apt-get update &&  \
    apt-get install -y curl git jq

RUN curl -L https://foundry.paradigm.xyz | bash && \
    source ~/.bashrc && \
    foundryup -v nightly-2044faec64f99a21f0e5f0094458a973612d0712 &&  \
    cp ~/.foundry/bin/* /usr/local/bin

COPY . .

RUN dev/docker/anvil-background && \
    dev/contracts/deploy-local && \
    dev/register-local-node && \
    dev/register-local-node-2 && \
    sleep 5

# ACTUAL IMAGE -------------------------------------------------------

FROM ghcr.io/foundry-rs/foundry:nightly-2044faec64f99a21f0e5f0094458a973612d0712

LABEL maintainer="engineering@xmtp.com"
LABEL source="https://github.com/xmtp/xmtpd"
LABEL description="Foundry with XMTPD Node contracts and registrations"

EXPOSE 8545

COPY --from=builder /app/anvil-baked-state anvil-baked-state

CMD ["anvil", "--state", "anvil-baked-state", "--host", "0.0.0.0"]