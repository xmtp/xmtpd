import { startGateway } from "./index.js";

const required = (name: string): string => {
  const val = process.env[name];
  if (!val) throw new Error(`Missing required env var: ${name}`);
  return val;
};

const gateway = await startGateway({
  payerPrivateKey: required("PAYER_PRIVATE_KEY"),
  redisUrl: required("REDIS_URL"),
  appChainRpcUrl: required("APP_CHAIN_RPC_URL"),
  appChainWssUrl: required("APP_CHAIN_WSS_URL"),
  settlementChainRpcUrl: required("SETTLEMENT_CHAIN_RPC_URL"),
  settlementChainWssUrl: required("SETTLEMENT_CHAIN_WSS_URL"),
  contractsEnvironment: process.env.CONTRACTS_ENVIRONMENT ?? "testnet",
  logLevel: process.env.LOG_LEVEL ?? "info",
});

console.log(`Gateway running at ${gateway.url}`);

process.on("SIGINT", async () => {
  await gateway.stop();
  process.exit(0);
});

process.on("SIGTERM", async () => {
  await gateway.stop();
  process.exit(0);
});
