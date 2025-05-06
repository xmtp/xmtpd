# How to onboard a new node in testnet

**⚠️ Experimental:** warning, this file might be out of date!

It is important that both `public key` and `address` are correct as they are immutable and cannot be changed in the future.

## Step 1) Get All Keys

An easy way to generate a new private key is to use our CLI:

```bash
./dev/cmd/cli generate-key | jq
{
  "level": "INFO",
  "time": "2024-10-15T13:21:14.036-0400",
  "message": "generated private key",
  "private-key": "0x7d3cd4989b92c593db9a4b3ac1c2a5d542efad058b2a83e26c3467392b29c6f9",
  "public-key": "0x03da53968d81f4eb3c9dd8b96617575767ec0cccbd28103b2cfd7f1511bb282d30",
  "address": "0x9419db765e6b469edc028ffa72ba2944f2bad169"
}
```

If you already have a private key, you can extract the relevant public details via:

```bash
./dev/cmd/cli get-pub-key --private-key 0xa9b48d687f450ea99a5faaae1be096ddb49487cb28393d3906d7359ede6ea460 | jq
{
  "level": "INFO",
  "time": "2024-10-15T13:21:51.276-0400",
  "message": "parsed private key",
  "pub-key": "0x027a64295b98e48682cb77be1b990d4ecf8f1a86badf051df0af123e6fe3790e3f",
  "address": "0x9419db765e6b469edc028ffa72ba2944f2bad169"
}

```

## Step 2) Share `pub-key` and `address`

Before official testnet launch, only the members of @ephemerahq/backend can register nodes.

## Step 3) Register Node with smart contract

Replace the following environment variables with real values.

```shell
export XMTPD_SETTLEMENT_CHAIN_RPC_URL="http://localhost:7545/"
export XMTPD_SETTLEMENT_CHAIN_CHAIN_ID=31337
export XMTPD_CONTRACTS_NODE_REGISTRY_ADDRESS="0xDEADBEEF"
export ADMIN_PRIVATE_KEY="0xDEADBEEF"
export NODE_HTTP_ADDRESS="https://grpc.example.com"
export NODE_OWNER_ADDRESS="0xDEADBEEF"
export NODE_SIGNING_KEY_PUB="0xDEADBEEF"

dev/cmd/cli register-node \
    --admin.private-key=${ADMIN_PRIVATE_KEY} \
    --http-address=${NODE_HTTP_ADDRESS} \
    --node-owner-address=${NODE_OWNER_ADDRESS} \
    --node-signing-key-pub=${NODE_SIGNING_KEY_PUB}
```

## Step 4) Start the node

This step might differ for every operator. A good starting point is our [Deployment Guide](deploy.md)
