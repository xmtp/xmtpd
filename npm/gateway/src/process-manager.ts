import { spawn } from "node:child_process";
import type { ChildProcess } from "node:child_process";
import { createServer } from "node:net";
import { resolveBinary } from "./binary.js";
import type { GatewayConfig, GatewayHandle } from "./types.js";

/**
 * Starts the XMTP gateway as a subprocess.
 *
 * Spawns the Go gateway binary with configuration passed via environment
 * variables and waits for it to become healthy before returning.
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

  setupLogForwarding(child);

  const earlyExit = new Promise<never>((_, reject) => {
    child.on("exit", (code, signal) => {
      reject(
        new Error(
          `Gateway exited during startup (code=${code}, signal=${signal})`,
        ),
      );
    });
    child.on("error", (err) => {
      reject(new Error(`Gateway failed to spawn: ${err.message}`));
    });
  });

  const timeout = config.healthCheckTimeout ?? 30_000;
  try {
    await Promise.race([waitForHealthy(port, timeout), earlyExit]);
  } catch (err) {
    child.kill("SIGKILL");
    throw err;
  }

  return {
    url: `http://localhost:${port}`,
    port,
    process: child,
    stop: () => gracefulShutdown(child),
  };
}

function buildEnv(config: GatewayConfig, port: number): NodeJS.ProcessEnv {
  const env: NodeJS.ProcessEnv = {
    ...process.env,
    XMTPD_API_PORT: String(port),
    XMTPD_API_ENABLE: "true",
    XMTPD_PAYER_ENABLE: "true",
    XMTPD_PAYER_PRIVATE_KEY: config.payerPrivateKey,
    XMTPD_REDIS_URL: config.redisUrl,
    XMTPD_APP_CHAIN_RPC_URL: config.appChainRpcUrl,
    XMTPD_APP_CHAIN_WSS_URL: config.appChainWssUrl,
    XMTPD_SETTLEMENT_CHAIN_RPC_URL: config.settlementChainRpcUrl,
    XMTPD_SETTLEMENT_CHAIN_WSS_URL: config.settlementChainWssUrl,
    XMTPD_LOG_ENCODING: config.logEncoding ?? "console",
    XMTPD_LOG_LEVEL: config.logLevel ?? "info",
  };

  if (config.contractsEnvironment) {
    env.XMTPD_CONTRACTS_ENVIRONMENT = config.contractsEnvironment;
  } else if (config.contractsConfigJson) {
    env.XMTPD_CONTRACTS_CONFIG_JSON = config.contractsConfigJson;
  } else if (config.contractsConfigFilePath) {
    env.XMTPD_CONTRACTS_CONFIG_FILE_PATH = config.contractsConfigFilePath;
  }

  if (config.nodeSelectorStrategy) {
    env.XMTPD_PAYER_NODE_SELECTOR_STRATEGY = config.nodeSelectorStrategy;
  }

  return env;
}

function setupLogForwarding(child: ChildProcess): void {
  let buffer = "";
  child.stdout?.on("data", (data: Buffer) => {
    buffer += data.toString();
    const parts = buffer.split("\n");
    buffer = parts.pop() ?? "";
    for (const line of parts) {
      if (!line) continue;
      try {
        const log = JSON.parse(line);
        const level = log.level ?? log.L ?? "info";
        const msg = log.msg ?? log.M ?? line;
        if (level === "error" || level === "fatal") {
          console.error(`[gateway] ${msg}`);
        } else if (level === "warn") {
          console.warn(`[gateway] ${msg}`);
        } else if (level !== "debug") {
          console.log(`[gateway] ${msg}`);
        }
      } catch {
        console.log(`[gateway] ${line}`);
      }
    }
  });
  child.stdout?.on("end", () => {
    if (buffer.trim()) {
      console.log(`[gateway] ${buffer.trim()}`);
      buffer = "";
    }
  });

  child.stderr?.on("data", (data: Buffer) => {
    for (const line of data.toString().split("\n")) {
      if (line) console.error(`[gateway:err] ${line}`);
    }
  });
}

async function waitForHealthy(
  port: number,
  timeoutMs: number,
): Promise<void> {
  const start = Date.now();

  while (Date.now() - start < timeoutMs) {
    if (await isPortInUse(port)) {
      await sleep(200);
      return;
    }
    await sleep(500);
  }

  throw new Error(
    `Gateway health check timed out after ${timeoutMs}ms on port ${port}`,
  );
}

function isPortInUse(port: number): Promise<boolean> {
  return new Promise((resolve) => {
    const server = createServer();
    server.once("error", () => resolve(true));
    server.once("listening", () => {
      server.close();
      resolve(false);
    });
    server.listen(port, "127.0.0.1");
  });
}

async function findAvailablePort(startPort: number): Promise<number> {
  for (let port = startPort; port < startPort + 100; port++) {
    if (!(await isPortInUse(port))) return port;
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

function sleep(ms: number): Promise<void> {
  return new Promise((resolve) => setTimeout(resolve, ms));
}
