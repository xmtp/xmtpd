import type { ChildProcess } from "node:child_process";

export interface GatewayConfig {
  /** Private key used to sign blockchain transactions */
  payerPrivateKey: string;
  /** Redis URL for nonce management */
  redisUrl: string;
  /** App Chain (L3) RPC endpoint */
  appChainRpcUrl: string;
  /** App Chain (L3) WebSocket endpoint */
  appChainWssUrl: string;
  /** Settlement Chain RPC endpoint */
  settlementChainRpcUrl: string;
  /** Settlement Chain WebSocket endpoint */
  settlementChainWssUrl: string;
  /** Deployed environment: 'testnet' or 'mainnet' */
  contractsEnvironment?: string;
  /** Inline JSON contracts config (alternative to contractsEnvironment) */
  contractsConfigJson?: string;
  /** Path to JSON contracts config file (alternative to contractsEnvironment) */
  contractsConfigFilePath?: string;
  /** gRPC listener port (default: auto-select starting from 5050) */
  port?: number;
  /** Log level (default: 'info') */
  logLevel?: string;
  /** Log format (default: 'console') */
  logEncoding?: "json" | "console";
  /** Node selection strategy (default: 'stable') */
  nodeSelectorStrategy?: "stable" | "random" | "ordered" | "closest" | "manual";
  /** Health check timeout in milliseconds (default: 30000) */
  healthCheckTimeout?: number;
}

export interface GatewayHandle {
  /** The URL to pass to the agent SDK as gatewayUrl */
  url: string;
  /** The port the gateway is listening on */
  port: number;
  /** The underlying child process */
  process: ChildProcess;
  /** Gracefully stop the gateway */
  stop: () => Promise<void>;
}
