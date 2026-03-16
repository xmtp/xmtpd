import type { ChildProcess } from "node:child_process";

export interface GatewayConfig {
  payerPrivateKey: string;
  redisUrl: string;
  appChainRpcUrl: string;
  appChainWssUrl: string;
  settlementChainRpcUrl: string;
  settlementChainWssUrl: string;
  /** 'testnet' or 'mainnet' */
  contractsEnvironment?: string;
  /** alternative to contractsEnvironment */
  contractsConfigJson?: string;
  /** alternative to contractsEnvironment */
  contractsConfigFilePath?: string;
  port?: number;
  logLevel?: string;
  nodeSelectorStrategy?: "stable" | "random" | "ordered" | "closest" | "manual";
  /** ms, default 30000 */
  healthCheckTimeout?: number;
}

export interface GatewayHandle {
  url: string;
  port: number;
  process: ChildProcess;
  stop: () => Promise<void>;
  stats: () => GatewayStats;
}

export interface GatewayStats {
  online: boolean;
  uptimeSeconds: number;
  gatewayPort: number;
  publishes: number;
  errors: number;
  requests: number;
}
