# @xmtp/gateway

Run the XMTP gateway as a sidecar process in your Node.js agent.

## Install

```bash
npm install @xmtp/gateway
```

The correct binary for your platform (macOS/Linux, x64/arm64) is installed automatically.

## Prerequisites

- **Redis** — the gateway requires Redis for nonce management
- **Payer wallet** — a funded private key for signing blockchain transactions
- **RPC endpoints** — for the App Chain and Settlement Chain

## Quick start

```typescript
import { startGateway } from "@xmtp/gateway";

const gateway = await startGateway({
  payerPrivateKey: process.env.PAYER_PRIVATE_KEY,
  redisUrl: "redis://localhost:6379",
  appChainRpcUrl: process.env.APP_CHAIN_RPC_URL,
  appChainWssUrl: process.env.APP_CHAIN_WSS_URL,
  settlementChainRpcUrl: process.env.SETTLEMENT_CHAIN_RPC_URL,
  settlementChainWssUrl: process.env.SETTLEMENT_CHAIN_WSS_URL,
  contractsEnvironment: "testnet",
});

console.log(`Gateway running at ${gateway.url}`);

// Stop when done
await gateway.stop();
```

## Connecting your agent

Pass the gateway URL to the XMTP agent SDK via the `XMTP_GATEWAY_HOST` environment variable or the `gatewayHost` option:

### Option 1: Environment variable

```bash
XMTP_GATEWAY_HOST=http://localhost:5050 node my-agent.js
```

### Option 2: Programmatic

```typescript
import { startGateway } from "@xmtp/gateway";
import { Agent } from "@xmtp/agent-sdk";

const gateway = await startGateway({
  payerPrivateKey: process.env.PAYER_PRIVATE_KEY,
  redisUrl: "redis://localhost:6379",
  appChainRpcUrl: process.env.APP_CHAIN_RPC_URL,
  appChainWssUrl: process.env.APP_CHAIN_WSS_URL,
  settlementChainRpcUrl: process.env.SETTLEMENT_CHAIN_RPC_URL,
  settlementChainWssUrl: process.env.SETTLEMENT_CHAIN_WSS_URL,
  contractsEnvironment: "testnet",
});

const agent = await Agent.create(signer, {
  gatewayHost: gateway.url,
  env: "testnet",
});
```

### Option 3: Start gateway, then use `createFromEnv`

```typescript
import { startGateway } from "@xmtp/gateway";
import { Agent } from "@xmtp/agent-sdk";

const gateway = await startGateway({ /* ... */ });

// Set the env var before creating the agent
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
| `logEncoding` | no | `"console"` | Log format (`"json"` or `"console"`) |
| `nodeSelectorStrategy` | no | `"stable"` | Node selection: `stable`, `random`, `ordered`, `closest`, `manual` |
| `healthCheckTimeout` | no | `30000` | Startup health check timeout (ms) |

\* Provide one of `contractsEnvironment`, `contractsConfigJson`, or `contractsConfigFilePath`.

## GatewayHandle

`startGateway()` returns a handle with:

| Property | Type | Description |
|----------|------|-------------|
| `url` | `string` | Gateway URL (e.g. `http://localhost:5050`) |
| `port` | `number` | Listening port |
| `process` | `ChildProcess` | Underlying child process |
| `stop()` | `() => Promise<void>` | Gracefully shut down the gateway |

## Custom binary path

To use a custom-built gateway binary:

```bash
export XMTP_GATEWAY_BINARY_PATH=/path/to/xmtp-gateway
```

## Supported platforms

| OS | Architecture | Package |
|----|-------------|---------|
| macOS | arm64 (Apple Silicon) | `@xmtp/gateway-darwin-arm64` |
| macOS | x64 (Intel) | `@xmtp/gateway-darwin-x64` |
| Linux | arm64 | `@xmtp/gateway-linux-arm64` |
| Linux | x64 | `@xmtp/gateway-linux-x64` |
