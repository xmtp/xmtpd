# Support a new XMTP testnet node operator

New node operators follow the guidance in [xmtpd-infrastructure](https://github.com/xmtp/xmtpd-infrastructure) to deploy a node to XMTP testnet.

This document outlines the steps the xmtpd dev team takes to support this process.

## Register their node

The new node operator will provide the xmtpd dev team with their node's public key and address. Use them to register their node.

Prompt them to ensure that the public key and address are correct as they are immutable and can't be changed in the future.

Only members of `@xmtp/backend` can currently register nodes.

```bash
export XMTPD_SETTLEMENT_CHAIN_WSS_URL="https://xmtp-testnet.g.alchemy.com/v2/<apikey>"
export XMTPD_SETTLEMENT_CHAIN_CHAIN_ID=34498
export XMTPD_SETTLEMENT_CHAIN_NODE_REGISTRY_ADDRESS=<depends>
export PRIVATE_KEY=<secret>

./dev/cmd/cli register-node \
    --http-address=<node DNS> \
    --node-owner-address=<node address> \
    --admin.private-key=$PRIVATE_KEY \
    --node-signing-key-pub=<node pub-key>
```

## Confirm node registration

To verify the node appears in the registry, run:

```bash
dev/cmd/get-all-nodes --admin.private-key=$ADMIN_PRIVATE_KEY
```

Look for the new node in the output.

## Notify the node operator

Once you've confirmed that their node is registered, let the node operator know.

You can provide them with the following links:

- Take the next step in deploying your node to the XMTP testnet: [Step 3: Set up dependencies](https://github.com/xmtp/xmtpd-infrastructure/blob/main/helm/README.md#step-3-set-up-dependencies)
- Optionally, you can use Kubernetes and Prometheus to set up observability: Set up Prometheus service discovery for xmtpd in Kubernetes using Helm
