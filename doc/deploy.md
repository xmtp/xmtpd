# How to deploy to Testnet

**⚠️ Experimental:** warning, this file might be out of date!

## Deploy a new Contract

The current environment lives in [Conduit Testnet Staging](https://explorer-testnet-staging-88dqtxdinc.t.conduit.xyz/).
To deploy a new contract you need to run `./contracts/dev/deploy-testnet`
You will need:

- $PRIVATE_KEY which is accessible to all members of @ephemerahq/backend
- $VERIFIER_URL: https://explorer-testnet-staging-88dqtxdinc.t.conduit.xyz/api
- $RPC_URL: https://rpc-testnet-staging-88dqtxdinc.t.conduit.xyz/

If the contract gets deployed correctly, you will get a few addresses.
We definitely need the node contract address (for example `0x7c9A7c92e21E9aC25Ce26C5e724920D84BD5BD9b`)

## Register nodes

Before nodes can start or peer, they need to be registered with the contract.
To do so, run:

```shell
export XMTPD_CONTRACTS_RPC_URL="https://rpc-testnet-staging-88dqtxdinc.t.conduit.xyz/"
export XMTPD_CONTRACTS_CHAIN_ID=34498
export XMTPD_CONTRACTS_NODES_ADDRESS=<from above>
export PRIVATE_KEY=<secret>

dev/cmd/cli register-node \
    --http-address=<node DNS> \
    --node-owner-address=0xd27FDB90A393Ce0E390120aeB58b326AbA910BE0 \
    --admin-private-key=$PRIVATE_KEY \
    --node-signing-key-pub=<node pub key>
```

You need to register all (both) nodes with their correct DNS entries and public keys.

### Verify Registration

To verify registration, use:

```shell
export XMTPD_CONTRACTS_RPC_URL="https://rpc-testnet-staging-88dqtxdinc.t.conduit.xyz/"
export XMTPD_CONTRACTS_CHAIN_ID=34498
export XMTPD_CONTRACTS_NODES_ADDRESS=<from above>
export PRIVATE_KEY=<secret>

dev/cmd/get-all-nodes \
    --admin-private-key=$PRIVATE_KEY
```

And you should get something along these lines:

```json
{
	"size": 2,
	"nodes": [
		{
			"NodeId": 100,
			"Node": {
				"SigningKeyPub": "BOVELF0f4vAra5oaOGODp3ZoYLQKYHmgIjmU/6LOEFEsToqIY97q2FnD1lQKsgJsgvi4k8HFvvbGP0fZ3zOiB9s=",
				"HttpAddress": "https://grpc.testnet.xmtp.network",
				"IsHealthy": true
			}
		},
		{
			"NodeId": 200,
			"Node": {
				"SigningKeyPub": "BPwmHUOgFTU5pMZMKXY8sOfjd8DqwpEMPUvtsiNaxwNxz+fKU3SsqOdYJQDVjLfRL5XsA5XVZIge2WDZ7S0zpx4=",
				"HttpAddress": "https://grpc2.testnet.xmtp.network",
				"IsHealthy": true
			}
		}
	]
}
```

## Deploy XMTPD nodes

Node deployment is currently fully handled by [Ephemera](https://github.com/ephemeraHQ/infrastructure) and only members of @ephemerahq/backend have access to it.

There are currently two nodes running:

| DNS Name                           | Location   | Public Key                                                           |
| ---------------------------------- | ---------- | -------------------------------------------------------------------- |
| https://grpc.testnet.xmtp.network  | US-EAST-2  | 0x03e5442c5d1fe2f02b6b9a1a386383a7766860b40a6079a0223994ffa2ce10512c |
| https://grpc2.testnet.xmtp.network | EU-NORTH-1 | 0x02fc261d43a0153539a4c64c29763cb0e7e377c0eac2910c3d4bedb2235ac70371 |

For more info, refer to the infrastructure README.

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
