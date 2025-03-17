# Node Registry Migration

## Set xmtpd cli branch

The migrator has been tested in commit `8d6f9252956aa0d7113d56e1deeca887c0d1c888`.

The cli can be used as a docker image using `ghcr.io/xmtp/xmtpd-cli:sha-8d6f925`

## Export V1 Nodes

```shell
go run cmd/cli/main.go get-all-nodes \
    --contracts.rpc-url <RPC_URL> \
    --contracts.nodes-address <NODES_REGISTRY_V1_CONTRACT_ADDRESS> \
    --out-file nodes.json
```

## Check the output

It should be similar to the following:

```json
[
	{
		"node_id": 100,
		"owner_address": "0x70997970C51812dc3A010C7d01b50e0d17dc79C8",
		"signing_key_pub": "0x02ba5734d8f7091719471e7f7ed6b9df170dc70cc661ca05e688601ad984f068b0",
		"http_address": "http://localhost:5050",
		"is_healthy": true
	},
	{
		"node_id": 200,
		"owner_address": "0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC",
		"signing_key_pub": "0x039d9031e97dd78ff8c15aa86939de9b1e791066a0224e331bc962a2099a7b1f04",
		"http_address": "http://localhost:5051",
		"is_healthy": true
	}
]
```

## Import nodes to V2 XMTP Nodes Registry

```shell
go run cmd/cli/main.go migrate-nodes \
    --contracts.rpc-url <RPC_URL> \
    --contracts.nodes-address <NODES_REGISTRY_V2_CONTRACT_ADDRESS> \
    --admin-private-key <PRIVATE_KEY> \
    --in-file nodes.json
```

If everything went well, the output would look like the following.

```text
2025-03-01T00:28:29.111+0100    INFO    NodeRegistryAdminV2     node added to registry V2       {"node_id": 100, "owner": "0x70997970C51812dc3A010C7d01b50e0d17dc79C8", "http_address": "http://localhost:5050", "signing_key_pub": "04ba5734d8f7091719471e7f7ed6b9df170dc70cc661ca05e688601ad984f068b0d67351e5f06073092499336ab0839ef8a521afd334e53807205fa2f08eec74f4", "min_monthly_fee": "0"}
2025-03-01T00:28:29.376+0100    INFO    NodeRegistryAdminV2     node added to registry V2       {"node_id": 200, "owner": "0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC", "http_address": "http://localhost:5051", "signing_key_pub": "049d9031e97dd78ff8c15aa86939de9b1e791066a0224e331bc962a2099a7b1f0464b8bbafe1535f2301c72c2cb3535b172da30b02686ab0393d348614f157fbdb", "min_monthly_fee": "0"}
```

## Activate nodes

In the XMTP Node Registry V2, there is no isHealthy flag.

Instead, the Node Operator is responsible of:

- Enabling API, via `function setIsApiEnabled(uint256 nodeId, bool isApiEnabled)`
- Enabling Replication, via `function setIsReplicationEnabled(uint256 nodeId, bool isReplicationEnabled)`
