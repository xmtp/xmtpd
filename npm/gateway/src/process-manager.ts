import { spawn } from "node:child_process";
import type { ChildProcess } from "node:child_process";
import { EventEmitter } from "node:events";
import { createConnection, createServer } from "node:net";
import { resolveBinary } from "./binary.js";
import { startStatusServer } from "./status-server.js";
import type { GatewayConfig, GatewayHandle, GatewayStats } from "./types.js";

const LOG_BUFFER_SIZE = 200;
const LOG_LEVELS: Record<string, number> = {
  debug: 0,
  info: 1,
  warn: 2,
  error: 3,
  fatal: 4,
};

/** Starts the XMTP gateway as a subprocess and waits for it to become healthy. */
export async function startGateway(
  config: GatewayConfig,
): Promise<GatewayHandle> {
  const binaryPath = resolveBinary();
  const port = config.port ?? (await findAvailablePort(5050));
  const statusPort = config.statusPort ?? port + 1;
  const env = buildEnv(config, port);

  const child = spawn(binaryPath, [], {
    env,
    stdio: ["ignore", "pipe", "pipe"],
  });

  // Stats tracking
  const startedAt = Date.now();
  const counters = { publishes: 0, errors: 0, requests: 0 };
  const logBuffer: string[] = [];
  const logEmitter = new EventEmitter();

  const getStats = (): GatewayStats => ({
    online: !child.killed && child.exitCode === null,
    uptimeSeconds: Math.floor((Date.now() - startedAt) / 1000),
    gatewayPort: port,
    publishes: counters.publishes,
    errors: counters.errors,
    requests: counters.requests,
  });

  const userLogLevel = config.logLevel ?? "info";
  setupLogForwarding(child, counters, logBuffer, logEmitter, userLogLevel);

  const onExit = (code: number | null, signal: string | null) => {
    earlyExitReject(
      new Error(
        `Gateway exited during startup (code=${code}, signal=${signal})`,
      ),
    );
  };
  const onError = (err: Error) => {
    earlyExitReject(new Error(`Gateway failed to spawn: ${err.message}`));
  };

  let earlyExitReject: (err: Error) => void;
  const earlyExit = new Promise<never>((_, reject) => {
    earlyExitReject = reject;
  });
  child.on("exit", onExit);
  child.on("error", onError);

  const timeout = config.healthCheckTimeout ?? 30_000;
  try {
    await Promise.race([waitForHealthy(port, timeout), earlyExit]);
  } catch (err) {
    child.kill("SIGKILL");
    throw err;
  } finally {
    child.removeListener("exit", onExit);
    child.removeListener("error", onError);
  }

  // Start status server
  const status = startStatusServer({
    port: statusPort,
    gatewayPort: port,
    getStats,
    logEmitter,
    logBuffer,
  });

  return {
    url: `http://localhost:${port}`,
    port,
    statusUrl: `http://localhost:${statusPort}`,
    process: child,
    stop: async () => {
      await status.close();
      await gracefulShutdown(child);
    },
    stats: getStats,
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
    XMTPD_LOG_ENCODING: "json",
    // Always use debug internally so we can count requests/publishes.
    // The log forwarder filters console output to the user's configured level.
    XMTPD_LOG_LEVEL: "debug",
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

function pushLog(
  logBuffer: string[],
  logEmitter: EventEmitter,
  text: string,
): void {
  logBuffer.push(text);
  if (logBuffer.length > LOG_BUFFER_SIZE) logBuffer.shift();
  logEmitter.emit("log", text);
}

function shouldLog(level: string, minLevel: string): boolean {
  return (LOG_LEVELS[level.toLowerCase()] ?? 1) >= (LOG_LEVELS[minLevel.toLowerCase()] ?? 1);
}

function setupLogForwarding(
  child: ChildProcess,
  counters: { publishes: number; errors: number; requests: number },
  logBuffer: string[],
  logEmitter: EventEmitter,
  consoleLogLevel: string,
): void {
  let buffer = "";
  child.stdout?.on("data", (data: Buffer) => {
    buffer += data.toString();
    const parts = buffer.split("\n");
    buffer = parts.pop() ?? "";
    for (const line of parts) {
      if (!line) continue;
      try {
        const log = JSON.parse(line);
        const level = (log.level ?? log.L ?? "info").toLowerCase();
        const msg = log.msg ?? log.M ?? log.message ?? line;

        // Track stats from all levels (including debug)
        if (level === "error" || level === "fatal") counters.errors++;
        if (msg.includes("received request")) counters.requests++;
        if (msg.includes("publishing to originator") || msg.includes("publishing to blockchain")) counters.publishes++;

        const formatted = `[gateway] ${msg}`;
        pushLog(logBuffer, logEmitter, formatted);

        // Only print to console if at or above user's configured level
        if (shouldLog(level, consoleLogLevel)) {
          if (level === "error" || level === "fatal") {
            console.error(formatted);
          } else if (level === "warn") {
            console.warn(formatted);
          } else {
            console.log(formatted);
          }
        }
      } catch {
        const formatted = `[gateway] ${line}`;
        pushLog(logBuffer, logEmitter, formatted);
        console.log(formatted);
      }
    }
  });
  child.stdout?.on("end", () => {
    if (buffer.trim()) {
      const formatted = `[gateway] ${buffer.trim()}`;
      pushLog(logBuffer, logEmitter, formatted);
      console.log(formatted);
      buffer = "";
    }
  });

  child.stderr?.on("data", (data: Buffer) => {
    for (const line of data.toString().split("\n")) {
      if (line) {
        const formatted = `[gateway:err] ${line}`;
        pushLog(logBuffer, logEmitter, formatted);
        console.error(formatted);
      }
    }
  });
}

async function waitForHealthy(
  port: number,
  timeoutMs: number,
): Promise<void> {
  const start = Date.now();

  while (Date.now() - start < timeoutMs) {
    if (await isPortListening(port)) {
      // Brief settle delay — the port may be open before the gRPC server is ready
      await sleep(200);
      return;
    }
    await sleep(500);
  }

  throw new Error(
    `Gateway health check timed out after ${timeoutMs}ms on port ${port}`,
  );
}

function isPortListening(port: number): Promise<boolean> {
  return new Promise((resolve) => {
    const socket = createConnection({ port, host: "127.0.0.1" });
    socket.once("connect", () => {
      socket.destroy();
      resolve(true);
    });
    socket.once("error", () => {
      socket.destroy();
      resolve(false);
    });
  });
}

async function findAvailablePort(startPort: number): Promise<number> {
  for (let port = startPort; port < startPort + 100; port++) {
    if (!(await isPortBound(port))) return port;
  }
  throw new Error(
    `No available port found in range ${startPort}-${startPort + 99}`,
  );
}

function isPortBound(port: number): Promise<boolean> {
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
