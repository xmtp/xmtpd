# @xmtp/gateway

Run the XMTP gateway as a sidecar process in your Node.js agent.

## Install

```bash
npm install @xmtp/gateway
```

Platform binaries (macOS/Linux, x64/arm64) are included. The right one is picked at runtime.

## Prerequisites

- **Redis** for nonce management
- **Payer wallet** - a funded private key for signing transactions
- **RPC endpoints** for the App Chain and Settlement Chain

See [SETUP.md](./SETUP.md) for how to get these.

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

Or set the env var before creating the agent:

```typescript
process.env.XMTP_GATEWAY_HOST = gateway.url;
const agent = await Agent.createFromEnv();
```

## Configuration

| Option | Required | Default | Description |
|--------|----------|---------|-------------|
| `payerPrivateKey` | yes | - | Private key for signing transactions |
| `redisUrl` | yes | - | Redis connection URL |
| `appChainRpcUrl` | yes | - | App Chain RPC endpoint |
| `appChainWssUrl` | yes | - | App Chain WebSocket endpoint |
| `settlementChainRpcUrl` | yes | - | Settlement Chain RPC endpoint |
| `settlementChainWssUrl` | yes | - | Settlement Chain WebSocket endpoint |
| `contractsEnvironment` | * | - | `"testnet"` or `"mainnet"` |
| `contractsConfigJson` | * | - | Inline JSON contracts config |
| `contractsConfigFilePath` | * | - | Path to JSON contracts config file |
| `port` | no | auto (5050+) | gRPC port |
| `logLevel` | no | `"info"` | `debug`, `info`, `warn`, `error` |
| `nodeSelectorStrategy` | no | `"stable"` | `stable`, `random`, `ordered`, `closest`, `manual` |
| `healthCheckTimeout` | no | `30000` | Startup health check timeout (ms) |

\* One of `contractsEnvironment`, `contractsConfigJson`, or `contractsConfigFilePath` is required.

## GatewayHandle

`startGateway()` returns:

| Property | Type | Description |
|----------|------|-------------|
| `url` | `string` | e.g. `http://localhost:5050` |
| `port` | `number` | Listening port |
| `process` | `ChildProcess` | The child process |
| `stop()` | `() => Promise<void>` | Shut down the gateway |
| `stats()` | `() => GatewayStats` | Publish/error/request counts |

## Logs and monitoring

Gateway logs are forwarded to the console with a `[gateway]` prefix. Errors go to stderr, the rest to stdout. Control verbosity with `logLevel`:

```typescript
const gateway = await startGateway({
  // ...
  logLevel: "debug",
});
```

You can also check stats programmatically:

```typescript
const s = gateway.stats();
console.log(`${s.publishes} publishes, ${s.errors} errors`);
```

## Custom binary

```bash
export XMTP_GATEWAY_BINARY_PATH=/path/to/xmtp-gateway
```

## Supported platforms

- macOS arm64 (Apple Silicon)
- macOS x64 (Intel)
- Linux arm64
- Linux x64
