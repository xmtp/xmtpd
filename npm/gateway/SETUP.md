# Gateway Setup Guide

How to set up and run the XMTP gateway in your Node.js agent.

## Prerequisites

- **Node.js 22+**
- **Redis** for nonce management
- **An RPC provider** like [Alchemy](https://www.alchemy.com/) (free tier works)

## 1. Install

```bash
npm install @xmtp/gateway
```

## 2. Start Redis

```bash
# Docker
docker run -d -p 6379:6379 redis

# Or Homebrew (macOS)
brew install redis && brew services start redis
```

Check it's working: `redis-cli ping` should return `PONG`.

## 3. Create a payer wallet

The payer wallet signs transactions on the XMTP network. Any Ethereum wallet works.

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

## 4. Fund the payer wallet

The payer needs funds allocated via the XMTP Funding Portal.

### Testnet

1. Go to [testnet.fund.xmtp.org](https://testnet.fund.xmtp.org/)
2. Connect your wallet or enter the payer address
3. Register it under "Manage payers"
4. Mint test tokens:

```bash
docker run --rm ghcr.io/xmtp/xmtpd-cli:main \
  funds mint --amount 1000 --to YOUR_WALLET_ADDRESS \
  --private-key YOUR_PRIVATE_KEY \
  --app-rpc-url https://xmtp-testnet.g.alchemy.com/v2/YOUR_KEY \
  --settlement-rpc-url https://base-sepolia.g.alchemy.com/v2/YOUR_KEY \
  --config-file https://github.com/xmtp/smart-contracts/releases/download/v0.5.5/testnet.json
```

5. Allocate funds in the Portal dashboard

You can also get testnet USDC from the [Circle faucet](https://faucet.circle.com/) (10 USDC/hour).

### Mainnet

Get USDC on Base, then allocate through [fund.xmtp.org](https://fund.xmtp.org/).

More info in the [funding docs](https://docs.xmtp.org/fund-agents-apps/fund-your-app).

## 5. Get RPC endpoints

Sign up at [Alchemy](https://www.alchemy.com/) and create apps for:

| Chain | Purpose | Testnet network |
|-------|---------|-----------------|
| **App Chain** (XMTP L3) | Message publishing | XMTP testnet |
| **Settlement Chain** | Contract settlement | Base Sepolia |

You'll get four URLs (HTTP + WebSocket for each):

```
# App Chain
https://xmtp-testnet.g.alchemy.com/v2/YOUR_KEY
wss://xmtp-testnet.g.alchemy.com/v2/YOUR_KEY

# Settlement Chain
https://base-sepolia.g.alchemy.com/v2/YOUR_KEY
wss://base-sepolia.g.alchemy.com/v2/YOUR_KEY
```

## 6. Run the gateway

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
```

## 7. Connect your agent

```typescript
import { Agent } from "@xmtp/agent-sdk";

const agent = await Agent.create(signer, {
  gatewayHost: gateway.url,
  env: "testnet",
});
```

Or via env var:

```typescript
process.env.XMTP_GATEWAY_HOST = gateway.url;
const agent = await Agent.createFromEnv();
```

## 8. Stop the gateway

```typescript
await gateway.stop();
```

Sends SIGTERM and waits up to 5s for graceful shutdown.

---

## Deploy to Render

The `render.yaml` in the xmtpd repo sets up a Node.js worker with managed Redis.

1. Fork the [xmtpd repo](https://github.com/xmtp/xmtpd)
2. On [Render](https://render.com/), create a new **Blueprint Instance**
3. Connect your fork - Render picks up the `render.yaml`
4. Fill in the env vars:
   - `PAYER_PRIVATE_KEY` - your funded payer key
   - `APP_CHAIN_RPC_URL` / `APP_CHAIN_WSS_URL`
   - `SETTLEMENT_CHAIN_RPC_URL` / `SETTLEMENT_CHAIN_WSS_URL`
5. Redis is provisioned automatically as `REDIS_URL`

### Manual setup (no blueprint)

Create a Render **Background Worker**:

- **Runtime**: Node
- **Build command**: `npm init -y && npm install @xmtp/gateway`
- **Start command**: `node node_modules/@xmtp/gateway/dist/start.js`

Add the env vars above, plus a Redis instance connected as `REDIS_URL`.

---

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
