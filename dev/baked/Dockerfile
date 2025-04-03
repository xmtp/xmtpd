# BUILD IMAGE --------------------------------------------------------
FROM ghcr.io/xmtp/contracts:sha-cff29fc AS builder

WORKDIR /app

SHELL ["/bin/bash", "-c"]

COPY . .

RUN dev/baked/bake.sh anvil-state.json

# ACTUAL IMAGE -------------------------------------------------------

FROM ghcr.io/foundry-rs/foundry

LABEL maintainer="eng@ephemerahq.com"
LABEL source="https://github.com/xmtp/xmtpd"
LABEL description="Foundry with XMTPD Node contracts and registrations"

EXPOSE 8545

COPY --from=builder /app/anvil-state.json anvil-state.json

ENTRYPOINT [ "anvil", "--host", "0.0.0.0", "--mixed-mining", "--block-time", "1", "--load-state", "anvil-state.json" ]
