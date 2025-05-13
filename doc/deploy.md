# How to deploy a local environment

**⚠️ Experimental:** warning, this file might be out of date!

## XMTP Sepolia information

The current testnet environment lives in [XMTP Sepolia Chain](https://xmtp-testnet.explorer.alchemy.com/).

### Deploy XMTPD nodes

Node deployment is currently fully handled by [Ephemera](https://github.com/ephemeraHQ) and only members of `@ephemerahq/backend` have access to it.

The nodes run by Ephemera are:

| DNS Name                           | Location   | Public Key                                                           |
| ---------------------------------- | ---------- | -------------------------------------------------------------------- |
| https://grpc.testnet.xmtp.network  | US-EAST-2  | 0x03e5442c5d1fe2f02b6b9a1a386383a7766860b40a6079a0223994ffa2ce10512c |
| https://grpc2.testnet.xmtp.network | EU-NORTH-1 | 0x02fc261d43a0153539a4c64c29763cb0e7e377c0eac2910c3d4bedb2235ac70371 |

For more info, refer to the infrastructure README.

## Local developer environment

Refer to [XMTP Contracts](https://github.com/xmtp/smart-contracts) for further information regarding the on-chain protocol, deployments and source code.

There are two ways of deploying a local environment:

### Use the XMTP Contracts Image

The [XMTP Contracts Image](https://github.com/xmtp/smart-contracts/blob/main/doc/xmtp-contracts-image.md#using-the-image) can be used to deploy a local environment and test `xmtpd` with it.

The documentation contains the deterministic addresses where all the contracts are deployed, to ease the setup step.

`dev/local.env` contains by default sane values for local deployments. Modify what is necessary.

### Use the dev/up automation

Use the script provided in `dev/up`, which will automatically handle the deployment for you. The blockchain is started at <http://localhost:7545/>

This method automatically pre-register two nodes.

### Register nodes manually

Before nodes can start or peer, they need to be registered with the contract.

To do so, run:

```shell
# Modify environment variables to match your local environment.
export XMTPD_SETTLEMENT_CHAIN_RPC_URL="http://localhost:7545/"
export XMTPD_SETTLEMENT_CHAIN_CHAIN_ID=31337
export XMTPD_SETTLEMENT_CHAIN_NODE_REGISTRY_ADDRESS="0xDEADBEEF"
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

You need to register all (both) nodes with their correct DNS entries and public keys.

### Verify a node registration

To verify registration, use:

```shell
export XMTPD_SETTLEMENT_CHAIN_RPC_URL="http://localhost:7545/"
export XMTPD_SETTLEMENT_CHAIN_CHAIN_ID=31337
export XMTPD_SETTLEMENT_CHAIN_NODE_REGISTRY_ADDRESS="0xDEADBEEF"
export ADMIN_PRIVATE_KEY="0xDEADBEEF"

dev/cmd/get-all-nodes \
    --admin.private-key=$PRIVATE_KEY
```

And you should get something along these lines:

```json
{
	"level": "INFO",
	"time": "2025-05-06T16:39:35.737+0200",
	"message": "got nodes",
	"size": 2,
	"nodes": [
		{
			"node_id": 100,
			"owner_address": "0x70997970C51812dc3A010C7d01b50e0d17dc79C8",
			"signing_key_pub": "0x02ba5734d8f7091719471e7f7ed6b9df170dc70cc661ca05e688601ad984f068b0",
			"http_address": "http://localhost:5050",
			"min_monthly_fee_micro_dollars": 0,
			"in_canonical_network": true
		},
		{
			"node_id": 200,
			"owner_address": "0xe67104BC93003192ab78B797d120DBA6e9Ff4928",
			"signing_key_pub": "0x028f67e68543faafa8540c0f4936435edb66cd5b4f398853914cb066f905e6130f",
			"http_address": "https://grpc.example.com",
			"min_monthly_fee_micro_dollars": 0,
			"in_canonical_network": false
		}
	]
}
```

### Verify deployed nodes

The easiest way is to use [GRPC Health Probe](https://github.com/grpc-ecosystem/grpc-health-probe)

```shell
grpc-health-probe -tls -addr grpc.testnet.xmtp.network:443
status: SERVING
```

```shell
grpc-health-probe -tls -addr grpc2.testnet.xmtp.network:443
status: SERVING
```
