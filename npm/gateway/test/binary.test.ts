import { describe, it, expect, beforeEach, afterEach } from "vitest";
import { resolveBinary } from "../src/binary.js";
import path from "node:path";

const localBinaryPath = path.resolve(
  __dirname,
  `../bin/xmtp-gateway-${process.platform}-${process.arch}`,
);

describe("resolveBinary", () => {
  const originalEnv = process.env;

  beforeEach(() => {
    process.env = { ...originalEnv };
  });

  afterEach(() => {
    process.env = originalEnv;
  });

  it("should use XMTP_GATEWAY_BINARY_PATH when set", () => {
    process.env.XMTP_GATEWAY_BINARY_PATH = localBinaryPath;
    expect(resolveBinary()).toBe(localBinaryPath);
  });

  it("should throw when XMTP_GATEWAY_BINARY_PATH points to non-existent file", () => {
    process.env.XMTP_GATEWAY_BINARY_PATH = "/nonexistent/xmtp-gateway";
    expect(() => resolveBinary()).toThrow("non-existent file");
  });

  it("should resolve binary from bin/ directory", () => {
    delete process.env.XMTP_GATEWAY_BINARY_PATH;
    const resolved = resolveBinary();
    expect(resolved).toContain("xmtp-gateway");
  });
});
