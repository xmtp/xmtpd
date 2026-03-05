# Local Network (Containerized Multi-Node)

Fully containerized 2-node XMTP network using pre-built GHCR images.
No Go toolchain required — useful for testing, QA, and running tools like `xdbg`.

> This is complementary to `dev/docker/` + `dev/run` which runs nodes natively from source.

## Quick Start

```bash
dev/local-network/start    # Brings up full network (~30s)
dev/local-network/stop     # Tears everything down
```

## Services

| Port  | Service             | Notes |
|-------|---------------------|-------|
| 8545  | Anvil (local chain) | RPC + WS |
| 60051 | MLS Validation      | gRPC |
| 5050  | xmtpd-1 (node 100) | gRPC API |
| 5051  | xmtpd-2 (node 200) | gRPC API |
| 5052  | Gateway             | gRPC API |
| 8765  | PostgreSQL (node 1) | |
| 8766  | PostgreSQL (node 2) | |
| 6379  | Redis               | |

> **Port conflict note:** The existing `dev/docker/` setup uses Anvil on port 7545.
> These two environments cannot run simultaneously.

## Verify

```bash
# Service status
docker compose -p xmtp-network ps

# Chain
curl -s localhost:8545 -X POST -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}'

# Nodes + Gateway (gRPC reflection)
grpcurl -plaintext localhost:5050 list
grpcurl -plaintext localhost:5051 list
grpcurl -plaintext localhost:5052 list

# Logs
docker compose -p xmtp-network logs --tail=50 xmtpd-1 xmtpd-2 gateway
```

## Updating Images

All images track `main`. Since `pull_policy: if_not_present` is set, pull manually to update:

```bash
docker pull ghcr.io/xmtp/xmtpd:main
docker pull ghcr.io/xmtp/xmtpd-gateway:main
docker pull ghcr.io/xmtp/xmtpd-cli:main
docker pull ghcr.io/xmtp/contracts:main
dev/local-network/stop
dev/local-network/start
```

## Origin

Imported from [fbac/xmtp-network-debug](https://github.com/fbac/xmtp-network-debug) with fixes:
- `--rpc-url` → `--settlement-rpc-url` (CLI flag rename)
- `pull_policy: always` → `if_not_present` (avoids GHCR token timeouts)
- Pinned SHA image tags → `main` (original pinned to Nov 2025 images)
