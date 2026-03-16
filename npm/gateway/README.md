# @xmtp/gateway

Run the XMTP gateway as a sidecar process in your Node.js agent.

## Install

```bash
npm install @xmtp/gateway
```

All platform binaries (macOS/Linux, x64/arm64) are included. The correct one is selected at runtime.

## Prerequisites

- **Redis** — the gateway requires Redis for nonce management
- **Payer wallet** — a funded private key for signing blockchain transactions
- **RPC endpoints** — for the App Chain and Settlement Chain

See [SETUP.md](./SETUP.md) for step-by-step instructions on getting these.

## Quick start

```typescript
import { startGateway } from "@xmtp/gateway";

const gateway = await startGateway({
  payerPrivateKey: process.env.PAYER_PRIVATE_KEY!,
  redisUrl: process.env.REDIS_URL ?? "redis://localhost:6379",
  appChainRpcUrl: process.env.APP_CHAIN_RPC_URL!,
  appChainWssUrl: process.env.APP_CHAIN_WSS_URL!,
  settlementChainRpcUrl: process.env.SETTLEMENT_CHAIN_RPC_URL!,
  settlementChainWssUrl: process.env.SETTLEMENT_CHAIN_WSS_URL!,
  contractsEnvironment: "testnet",
});

console.log(`Gateway running at ${gateway.url}`);

// Stop when done
await gateway.stop();
```

## Connecting your agent

Pass the gateway URL to the XMTP agent SDK:

```typescript
import { startGateway } from "@xmtp/gateway";
import { Agent } from "@xmtp/agent-sdk";

const gateway = await startGateway({
  payerPrivateKey: process.env.PAYER_PRIVATE_KEY!,
  redisUrl: process.env.REDIS_URL ?? "redis://localhost:6379",
  appChainRpcUrl: process.env.APP_CHAIN_RPC_URL!,
  appChainWssUrl: process.env.APP_CHAIN_WSS_URL!,
  settlementChainRpcUrl: process.env.SETTLEMENT_CHAIN_RPC_URL!,
  settlementChainWssUrl: process.env.SETTLEMENT_CHAIN_WSS_URL!,
  contractsEnvironment: "testnet",
});

const agent = await Agent.create(signer, {
  gatewayHost: gateway.url,
  env: "testnet",
});
```

Or set the environment variable before creating the agent:

```typescript
process.env.XMTP_GATEWAY_HOST = gateway.url;
const agent = await Agent.createFromEnv();
```

## Configuration

| Option | Required | Default | Description |
|--------|----------|---------|-------------|
| `payerPrivateKey` | yes | — | Private key for signing blockchain transactions |
| `redisUrl` | yes | — | Redis connection URL |
| `appChainRpcUrl` | yes | — | App Chain RPC endpoint |
| `appChainWssUrl` | yes | — | App Chain WebSocket endpoint |
| `settlementChainRpcUrl` | yes | — | Settlement Chain RPC endpoint |
| `settlementChainWssUrl` | yes | — | Settlement Chain WebSocket endpoint |
| `contractsEnvironment` | * | — | `"testnet"` or `"mainnet"` |
| `contractsConfigJson` | * | — | Inline JSON contracts config |
| `contractsConfigFilePath` | * | — | Path to JSON contracts config file |
| `port` | no | auto (5050+) | gRPC listener port |
| `logLevel` | no | `"info"` | Log level (`debug`, `info`, `warn`, `error`) |
| `nodeSelectorStrategy` | no | `"stable"` | Node selection: `stable`, `random`, `ordered`, `closest`, `manual` |
| `healthCheckTimeout` | no | `30000` | Startup health check timeout (ms) |
| `statusPort` | no | port + 1 | HTTP status dashboard port |

\* Provide one of `contractsEnvironment`, `contractsConfigJson`, or `contractsConfigFilePath`.

## GatewayHandle

`startGateway()` returns a handle with:

| Property | Type | Description |
|----------|------|-------------|
| `url` | `string` | Gateway URL (e.g. `http://localhost:5050`) |
| `port` | `number` | Listening port |
| `statusUrl` | `string` | Status dashboard URL (e.g. `http://localhost:5051`) |
| `process` | `ChildProcess` | Underlying child process |
| `stop()` | `() => Promise<void>` | Gracefully shut down the gateway |
| `stats()` | `() => GatewayStats` | Get current publish/error/request counts |

## Custom binary path

To use a custom-built gateway binary:

```bash
export XMTP_GATEWAY_BINARY_PATH=/path/to/xmtp-gateway
```

## Supported platforms

- macOS arm64 (Apple Silicon)
- macOS x64 (Intel)
- Linux arm64
- Linux x64
