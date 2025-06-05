# Deploy a local environment for developing with xmtpd

**⚠️ Experimental:** This file might be out of date!

This document describes how to deploy a local environment for developing with xmtpd.

The local environment includes everything you need to run xmtpd, including:

- Required databases
- The MLS validation service
- A blockchain with all the contracts

Use one of the following methods to deploy a local environment:

- Use the XMTP Contracts Image
- Use the `dev/up` automation

### Use the XMTP Contracts Image

You can use the XMTP Contracts Image to deploy a local environment, and test `xmtpd` with it.

The XMTP Contracts Image [documentation](https://github.com/xmtp/smart-contracts/blob/main/doc/xmtp-contracts-image.md#using-the-image) contains the deterministic addresses where all the contracts are deployed, to ease the setup step.

By default, `dev/local.env` contains sane values for local deployments. Modify what's necessary.

For more information about the onchain protocol, deployments, and source code, see the XMTP [smart-contracts](https://github.com/xmtp/smart-contracts) repo.

### Use the dev/up automation

Use the script provided in `dev/up` to automatically handle the deployment for you. The blockchain starts at <http://localhost:7545/>.

This method automatically pre-registers one node.

## Register a node with the local network

Once you've deployed your local environment, manually register your nodes to use it.

Before a node can start or peer on the local network, you must first register it.

To register a node on the local network, run:

```shell
# Modify environment variables to match your local environment.
export XMTPD_SETTLEMENT_CHAIN_WSS_URL="ws://localhost:7545/"
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

### Verify node registration

To verify node registration on the local network, run:

```shell
export XMTPD_SETTLEMENT_CHAIN_WSS_URL="ws://localhost:7545/"
export XMTPD_SETTLEMENT_CHAIN_CHAIN_ID=31337
export XMTPD_SETTLEMENT_CHAIN_NODE_REGISTRY_ADDRESS="0xDEADBEEF"
export ADMIN_PRIVATE_KEY="0xDEADBEEF"

dev/cmd/get-all-nodes \
    --admin.private-key=$PRIVATE_KEY
```

The response should look something like this:

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

### Verify node deployment

The easiest way to verify node deployment on the local network is to use the [gRPC Health Probe](https://github.com/grpc-ecosystem/grpc-health-probe).

For example, you can run:

```shell
grpc-health-probe -tls -addr grpc.testnet.xmtp.network:443
status: SERVING
```

```shell
grpc-health-probe -tls -addr grpc2.testnet.xmtp.network:443
status: SERVING
```
