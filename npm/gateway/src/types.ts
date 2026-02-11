import type { ChildProcess } from "node:child_process";

export interface GatewayConfig {
  // Required — payer identity
  payerPrivateKey: string;

  // Required — Redis for nonce management
  redisUrl: string;

  // Required — App Chain (L3) endpoints
  appChainRpcUrl: string;
  appChainWssUrl: string;

  // Required — Settlement Chain endpoints
  settlementChainRpcUrl: string;
  settlementChainWssUrl: string;

  // Contract configuration — provide ONE of these
  contractsEnvironment?: string; // e.g. 'dev', 'production'
  contractsConfigJson?: string; // inline JSON string
  contractsConfigFilePath?: string; // path to JSON file

  // Optional
  port?: number; // default: auto-select starting from 5050
  logLevel?: string; // default: 'info'
  logEncoding?: "json" | "console"; // default: 'json'
  nodeSelectorStrategy?: "stable" | "random" | "ordered" | "closest" | "manual";
  autoRestart?: boolean; // default: true
  healthCheckTimeout?: number; // milliseconds, default: 30000
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
