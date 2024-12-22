# BUILD IMAGE --------------------------------------------------------
ARG GO_VERSION=1.23
ARG FOUNDRY_VERSION=nightly
FROM golang:${GO_VERSION}-bookworm AS builder

WORKDIR /app

SHELL ["/bin/bash", "-c"]

RUN apt-get update &&  \
    apt-get install -y curl git jq && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

RUN curl -L https://foundry.paradigm.xyz | bash && \
    source ~/.bashrc && \
    foundryup -i "${FOUNDRY_VERSION}" &&  \
    cp ~/.foundry/bin/* /usr/local/bin

COPY . .

# It seems that anvil flushes the file to disk on shutdown and it takes a few ms to be persisted
# That gives us the pkill+sleep requirement
RUN dev/docker/anvil-background && \
    dev/contracts/deploy-local && \
    dev/register-local-node && \
    dev/register-local-node-2 && \
    pkill -f anvil && \
    sleep 5

RUN echo "export XMTPD_CONTRACTS_NODES_ADDRESS="$(jq -r '.deployedTo' build/Nodes.json)"" >> contracts.env && \
    echo "export XMTPD_CONTRACTS_MESSAGES_ADDRESS="$(jq -r '.addresses.groupMessagesProxy' contracts/config/anvil_localnet/GroupMessages.json)""  >> contracts.env && \
    echo "export XMTPD_CONTRACTS_IDENTITY_UPDATES_ADDRESS="$(jq -r '.addresses.identityUpdatesProxy' contracts/config/anvil_localnet/IdentityUpdates.json)""  >> contracts.env

# ACTUAL IMAGE -------------------------------------------------------

FROM ghcr.io/foundry-rs/foundry

LABEL maintainer="engineering@xmtp.com"
LABEL source="https://github.com/xmtp/xmtpd"
LABEL description="Foundry with XMTPD Node contracts and registrations"

EXPOSE 8545

COPY --from=builder /app/anvil-baked-state anvil-baked-state
COPY --from=builder /app/contracts.env contracts.env

ENTRYPOINT ["anvil", "--state", "anvil-baked-state"]