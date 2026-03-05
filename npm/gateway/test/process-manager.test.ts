import { describe, it, expect, afterEach } from "vitest";
import { startGateway } from "../src/process-manager.js";
import type { GatewayHandle } from "../src/types.js";
import path from "node:path";

const binaryPath = path.resolve(
  __dirname,
  `../bin/xmtp-gateway-${process.platform}-${process.arch}`,
);

describe("startGateway", () => {
  let handle: GatewayHandle | undefined;
  const savedBinaryPath = process.env.XMTP_GATEWAY_BINARY_PATH;

  afterEach(async () => {
    if (handle) {
      await handle.stop();
      handle = undefined;
    }
    if (savedBinaryPath !== undefined) {
      process.env.XMTP_GATEWAY_BINARY_PATH = savedBinaryPath;
    } else {
      delete process.env.XMTP_GATEWAY_BINARY_PATH;
    }
  });

  it("should fail when gateway exits during startup", async () => {
    process.env.XMTP_GATEWAY_BINARY_PATH = binaryPath;

    await expect(
      startGateway({
        payerPrivateKey: "0x0000000000000000000000000000000000000000000000000000000000000001",
        redisUrl: "redis://localhost:63790",
        appChainRpcUrl: "http://localhost:1111",
        appChainWssUrl: "ws://localhost:1112",
        settlementChainRpcUrl: "http://localhost:1113",
        settlementChainWssUrl: "ws://localhost:1114",
        // Invalid environment to force the gateway to crash on startup
        contractsEnvironment: "invalid_env_for_test",
        healthCheckTimeout: 5_000,
      }),
    ).rejects.toThrow();
  }, 15_000);
});
