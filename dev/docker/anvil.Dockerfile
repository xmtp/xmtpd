FROM ghcr.io/foundry-rs/foundry

WORKDIR /anvil

ENTRYPOINT anvil --host 0.0.0.0