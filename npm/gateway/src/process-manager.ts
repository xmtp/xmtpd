import { spawn } from "node:child_process";
import type { ChildProcess } from "node:child_process";
import { createServer } from "node:net";
import { resolveBinary } from "./binary.js";
import type { GatewayConfig, GatewayHandle } from "./types.js";

/**
 * Starts the XMTP gateway as a subprocess.
 *
 * The Go gateway binary is spawned with configuration passed via environment
 * variables. This function waits for the gateway to become healthy before
 * returning.
 */
export async function startGateway(
  config: GatewayConfig,
): Promise<GatewayHandle> {
  const binaryPath = resolveBinary();
  const port = config.port ?? (await findAvailablePort(5050));
  const env = buildEnv(config, port);

  const child = spawn(binaryPath, [], {
    env,
    stdio: ["ignore", "pipe", "pipe"],
  });

  // Forward gateway logs
  setupLogForwarding(child);

  // Handle early exit (e.g. config error)
  const earlyExitPromise = new Promise<never>((_, reject) => {
    child.on("exit", (code, signal) => {
      reject(
        new Error(
          `Gateway exited unexpectedly during startup (code=${code}, signal=${signal})`,
        ),
      );
    });
  });

  // Wait for health check or early exit
  const timeout = config.healthCheckTimeout ?? 30_000;
  try {
    await Promise.race([waitForHealthy(port, timeout), earlyExitPromise]);
  } catch (err) {
    // Clean up on failure
    child.kill("SIGKILL");
    throw err;
  }

  const handle: GatewayHandle = {
    url: `http://localhost:${port}`,
    port,
    process: child,
    stop: () => gracefulShutdown(child),
  };

  // Auto-restart on unexpected exit
  if (config.autoRestart !== false) {
    setupAutoRestart(child, config, handle);
  }

  // Clean up gateway when Node.js exits
  registerCleanup(child);

  return handle;
}

function buildEnv(config: GatewayConfig, port: number): NodeJS.ProcessEnv {
  const env: NodeJS.ProcessEnv = {
    ...process.env,
    // API
    XMTPD_API_PORT: String(port),
    XMTPD_API_ENABLE: "true",
    // Payer
    XMTPD_PAYER_ENABLE: "true",
    XMTPD_PAYER_PRIVATE_KEY: config.payerPrivateKey,
    // Redis
    XMTPD_REDIS_URL: config.redisUrl,
    // App Chain
    XMTPD_APP_CHAIN_RPC_URL: config.appChainRpcUrl,
    XMTPD_APP_CHAIN_WSS_URL: config.appChainWssUrl,
    // Settlement Chain
    XMTPD_SETTLEMENT_CHAIN_RPC_URL: config.settlementChainRpcUrl,
    XMTPD_SETTLEMENT_CHAIN_WSS_URL: config.settlementChainWssUrl,
    // Logging
    XMTPD_LOG_ENCODING: config.logEncoding ?? "json",
    XMTPD_LOG_LEVEL: config.logLevel ?? "info",
  };

  // Contracts config — mutually exclusive options
  if (config.contractsEnvironment) {
    env.XMTPD_CONTRACTS_ENVIRONMENT = config.contractsEnvironment;
  } else if (config.contractsConfigJson) {
    env.XMTPD_CONTRACTS_CONFIG_JSON = config.contractsConfigJson;
  } else if (config.contractsConfigFilePath) {
    env.XMTPD_CONTRACTS_CONFIG_FILE_PATH = config.contractsConfigFilePath;
  }

  // Node selector
  if (config.nodeSelectorStrategy) {
    env.XMTPD_PAYER_NODE_SELECTOR_STRATEGY = config.nodeSelectorStrategy;
  }

  return env;
}

function setupLogForwarding(child: ChildProcess): void {
  child.stdout?.on("data", (data: Buffer) => {
    const lines = data.toString().split("\n").filter(Boolean);
    for (const line of lines) {
      try {
        const log = JSON.parse(line);
        const level = log.level ?? log.L ?? "info";
        const msg = log.msg ?? log.M ?? line;
        if (level === "error" || level === "fatal") {
          console.error(`[gateway] ${msg}`);
        } else if (level === "warn") {
          console.warn(`[gateway] ${msg}`);
        } else if (level === "debug") {
          // skip debug in forwarding
        } else {
          console.log(`[gateway] ${msg}`);
        }
      } catch {
        // Not JSON, forward raw
        console.log(`[gateway] ${line}`);
      }
    }
  });

  child.stderr?.on("data", (data: Buffer) => {
    const lines = data.toString().split("\n").filter(Boolean);
    for (const line of lines) {
      console.error(`[gateway:err] ${line}`);
    }
  });
}

/**
 * Poll the gateway's gRPC health endpoint until it responds.
 * Uses a simple TCP connection check since the health endpoint is HTTP/2.
 */
async function waitForHealthy(
  port: number,
  timeoutMs: number,
): Promise<void> {
  const start = Date.now();
  const interval = 500;

  while (Date.now() - start < timeoutMs) {
    try {
      const isOpen = await checkPort(port);
      if (isOpen) {
        // Port is open — give the server a moment to finish initialization
        await sleep(200);
        return;
      }
    } catch {
      // Connection refused, keep trying
    }
    await sleep(interval);
  }

  throw new Error(
    `Gateway health check timed out after ${timeoutMs}ms on port ${port}`,
  );
}

function checkPort(port: number): Promise<boolean> {
  return new Promise((resolve) => {
    const socket = createServer();
    socket.once("error", () => {
      socket.close();
      // Port is in use (gateway is listening) — this is what we want
      resolve(true);
    });
    socket.once("listening", () => {
      socket.close();
      // Port is free — gateway hasn't started listening yet
      resolve(false);
    });
    socket.listen(port, "127.0.0.1");
  });
}

async function findAvailablePort(startPort: number): Promise<number> {
  for (let port = startPort; port < startPort + 100; port++) {
    const isAvailable = await new Promise<boolean>((resolve) => {
      const server = createServer();
      server.once("error", () => resolve(false));
      server.once("listening", () => {
        server.close();
        resolve(true);
      });
      server.listen(port, "127.0.0.1");
    });
    if (isAvailable) return port;
  }
  throw new Error(
    `No available port found in range ${startPort}-${startPort + 99}`,
  );
}

function gracefulShutdown(child: ChildProcess): Promise<void> {
  return new Promise((resolve) => {
    if (child.killed || child.exitCode !== null) {
      resolve();
      return;
    }

    const forceKillTimer = setTimeout(() => {
      child.kill("SIGKILL");
    }, 5_000);

    child.once("exit", () => {
      clearTimeout(forceKillTimer);
      resolve();
    });

    child.kill("SIGTERM");
  });
}

function setupAutoRestart(
  child: ChildProcess,
  config: GatewayConfig,
  handle: GatewayHandle,
): void {
  child.on("exit", (code, signal) => {
    // Don't restart if it was a graceful shutdown (SIGTERM)
    if (signal === "SIGTERM" || signal === "SIGKILL") return;

    console.warn(
      `[gateway] Process exited unexpectedly (code=${code}). Restarting...`,
    );

    // Restart after a brief delay
    setTimeout(async () => {
      try {
        const newHandle = await startGateway({
          ...config,
          port: handle.port,
          autoRestart: true,
        });
        // Update the handle in-place
        handle.process = newHandle.process;
        handle.url = newHandle.url;
        handle.port = newHandle.port;
        handle.stop = newHandle.stop;
      } catch (err) {
        console.error(`[gateway] Failed to restart:`, err);
      }
    }, 1_000);
  });
}

function registerCleanup(child: ChildProcess): void {
  const cleanup = () => {
    if (!child.killed && child.exitCode === null) {
      child.kill("SIGTERM");
    }
  };

  process.on("exit", cleanup);
  process.on("SIGINT", cleanup);
  process.on("SIGTERM", cleanup);
}

function sleep(ms: number): Promise<void> {
  return new Promise((resolve) => setTimeout(resolve, ms));
}
