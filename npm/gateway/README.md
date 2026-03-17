# @xmtp/gateway

Run the XMTP gateway as a sidecar in your Node.js agent.

## Install

```bash
npm install @xmtp/gateway
```

Platform binaries (macOS/Linux, x64/arm64) are included. The right one is picked at runtime.

## Prerequisites

- **Node.js 22+**
- **Redis** for nonce management
- **Payer wallet** - a funded private key for signing transactions
- **RPC endpoints** for the App Chain and Settlement Chain

## Setup

### 1. Start Redis

```bash
# Docker
docker run -d -p 6379:6379 redis

# Or Homebrew (macOS)
brew install redis && brew services start redis
```

Check it's working: `redis-cli ping` should return `PONG`.

### 2. Create a payer wallet

Any Ethereum wallet works. Generate one:

```bash
node -e "console.log('0x' + require('crypto').randomBytes(32).toString('hex'))"
```

Get the address:

```bash
node -e "
  const { Wallet } = require('ethers');
  console.log(new Wallet('YOUR_PRIVATE_KEY').address);
"
```

Save both the private key and address.

### 3. Fund the payer wallet

The payer needs funds allocated via the XMTP Funding Portal.

**Testnet:**

1. Go to [testnet.fund.xmtp.org](https://testnet.fund.xmtp.org/)
2. Connect your wallet or enter the payer address
3. Register it under "Manage payers"
4. Mint test tokens:

```bash
docker run --rm ghcr.io/xmtp/xmtpd-cli:latest \
  funds mint --amount 1000 --to YOUR_WALLET_ADDRESS \
  --private-key YOUR_PRIVATE_KEY \
  --config-file=config://testnet
```

5. Allocate funds in the Portal dashboard

You can also get testnet USDC from the [Circle faucet](https://faucet.circle.com/) (10 USDC/hour).

**Mainnet:**

Get USDC on Base, then allocate through [fund.xmtp.org](https://fund.xmtp.org/). More info in the [funding docs](https://docs.xmtp.org/fund-agents-apps/fund-your-app).

### 4. Get RPC endpoints

Sign up at [Alchemy](https://www.alchemy.com/) and create apps for:

| Chain | Purpose | Testnet network |
|-------|---------|-----------------|
| **App Chain** (XMTP L3) | Message publishing | XMTP Testnet |
| **Settlement Chain** | Contract settlement | Base Sepolia |

You'll get four URLs (HTTP + WebSocket for each):

```
# App Chain
https://xmtp-ropsten.g.alchemy.com/v2/YOUR_KEY
wss://xmtp-ropsten.g.alchemy.com/v2/YOUR_KEY

# Settlement Chain
https://base-sepolia.g.alchemy.com/v2/YOUR_KEY
wss://base-sepolia.g.alchemy.com/v2/YOUR_KEY
```

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

Check stats programmatically:

```typescript
const s = gateway.stats();
console.log(`${s.publishes} publishes, ${s.errors} errors`);
```

## Deploy to Render

The `render.yaml` in the xmtpd repo sets up a Node.js worker with managed Redis.

1. Fork the [xmtpd repo](https://github.com/xmtp/xmtpd)
2. On [Render](https://render.com/), create a new **Blueprint Instance**
3. Connect your fork - Render picks up the `render.yaml`
4. Fill in the env vars:
   - `PAYER_PRIVATE_KEY`
   - `APP_CHAIN_RPC_URL` / `APP_CHAIN_WSS_URL`
   - `SETTLEMENT_CHAIN_RPC_URL` / `SETTLEMENT_CHAIN_WSS_URL`
5. Redis is provisioned automatically as `REDIS_URL`

**Manual setup (no blueprint):**

Create a Render **Background Worker**:

- **Runtime**: Node
- **Build command**: `npm init -y && npm install @xmtp/gateway`
- **Start command**: `node node_modules/@xmtp/gateway/dist/start.js`

Add the env vars above plus a Redis instance connected as `REDIS_URL`.

## Custom binary

```bash
export XMTP_GATEWAY_BINARY_PATH=/path/to/xmtp-gateway
```

## Supported platforms

- macOS arm64 (Apple Silicon)
- macOS x64 (Intel)
- Linux arm64
- Linux x64

## Troubleshooting

**Gateway exits immediately**
- Check Redis: `redis-cli ping`
- `contractsEnvironment` must be `"testnet"` or `"mainnet"`
- Check your RPC URLs

**`INSUFFICIENT_PAYER_BALANCE`**
- Allocate more funds at [testnet.fund.xmtp.org](https://testnet.fund.xmtp.org/)

**Health check timeout**
- Bump `healthCheckTimeout` (default 30s)
- Check RPC endpoints are reachable

**Binary not found**
- Reinstall: `npm install @xmtp/gateway`
- Or set `XMTP_GATEWAY_BINARY_PATH`
