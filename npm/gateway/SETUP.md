# Gateway Setup Guide

Set up and run the XMTP gateway in your Node.js agent.

## Prerequisites

- **Node.js 22+**
- **Redis** — used for nonce management
- **An RPC provider** — [Alchemy](https://www.alchemy.com/) (free tier works) or similar

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

Verify: `redis-cli ping` should return `PONG`.

## 3. Create a payer wallet

The payer wallet signs transactions to publish messages on the XMTP network. Any standard Ethereum wallet works (secp256k1 key).

```bash
node -e "console.log('0x' + require('crypto').randomBytes(32).toString('hex'))"
```

To get the wallet address:

```bash
node -e "
  const { Wallet } = require('ethers');
  console.log(new Wallet('YOUR_PRIVATE_KEY').address);
"
```

Save both the private key and address.

## 4. Fund the payer wallet

Your payer wallet needs funds allocated through the XMTP Funding Portal before the gateway can process messages.

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

Acquire USDC on Base, then allocate through [fund.xmtp.org](https://fund.xmtp.org/).

See the [official funding docs](https://docs.xmtp.org/fund-agents-apps/fund-your-app) for details.

## 5. Get RPC endpoints

Sign up at [Alchemy](https://www.alchemy.com/) and create apps for these two networks:

| Chain | Purpose | Testnet network |
|-------|---------|-----------------|
| **App Chain** (XMTP L3) | Message publishing | XMTP testnet |
| **Settlement Chain** | Contract settlement | Base Sepolia |

You'll get four URLs (HTTP + WebSocket for each chain):

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

Pass the gateway URL to the agent SDK:

```typescript
import { Agent } from "@xmtp/agent-sdk";

const agent = await Agent.create(signer, {
  gatewayHost: gateway.url,
  env: "testnet",
});
```

Or via environment variable:

```typescript
process.env.XMTP_GATEWAY_HOST = gateway.url;
const agent = await Agent.createFromEnv();
```

## 8. Stop the gateway

```typescript
await gateway.stop();
```

This sends SIGTERM to the subprocess and waits up to 5 seconds for graceful shutdown.

---

## Deploy to Render

The easiest way to deploy is using the `render.yaml` blueprint in the xmtpd repo. It sets up a Node.js worker with a managed Redis instance.

1. Fork the [xmtpd repo](https://github.com/xmtp/xmtpd) on GitHub
2. Go to [Render](https://render.com/) and create a new **Blueprint Instance**
3. Connect your fork — Render auto-detects the `render.yaml`
4. Fill in the required env vars when prompted:
   - `PAYER_PRIVATE_KEY` — your funded payer wallet key
   - `APP_CHAIN_RPC_URL` / `APP_CHAIN_WSS_URL` — XMTP chain endpoints
   - `SETTLEMENT_CHAIN_RPC_URL` / `SETTLEMENT_CHAIN_WSS_URL` — Base chain endpoints
5. Redis is provisioned automatically and wired via `REDIS_URL`

The blueprint installs `@xmtp/gateway` from npm (which includes the pre-built Go binary) and runs the bundled start script.

### Deploy without a blueprint

If you prefer to set up manually, create a Render **Background Worker** with:

- **Runtime**: Node
- **Build command**: `npm init -y && npm install @xmtp/gateway`
- **Start command**: `node node_modules/@xmtp/gateway/dist/start.js`

Add the same env vars listed above, plus a Render Redis (Key Value) instance connected as `REDIS_URL`.

---

## Troubleshooting

**Gateway exits immediately**
- Check Redis is running: `redis-cli ping`
- Make sure `contractsEnvironment` is `"testnet"` or `"mainnet"`
- Verify your RPC URLs are correct

**`INSUFFICIENT_PAYER_BALANCE`**
- Allocate more funds at [testnet.fund.xmtp.org](https://testnet.fund.xmtp.org/)

**Health check timeout**
- Increase `healthCheckTimeout` (default is 30s)
- Check that RPC endpoints are reachable

**Binary not found**
- Reinstall with `npm install @xmtp/gateway`
- Or set `XMTP_GATEWAY_BINARY_PATH` to a custom binary path
