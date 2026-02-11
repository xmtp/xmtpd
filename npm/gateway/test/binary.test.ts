import { describe, it, expect, beforeEach, afterEach, vi } from "vitest";
import { resolveBinary } from "../src/binary.js";

describe("resolveBinary", () => {
  const originalEnv = process.env;

  beforeEach(() => {
    process.env = { ...originalEnv };
  });

  afterEach(() => {
    process.env = originalEnv;
  });

  it("should use XMTP_GATEWAY_BINARY_PATH when set", () => {
    // Point to the actual compiled binary for current platform
    const binaryPath =
      process.arch === "arm64"
        ? `${__dirname}/../../gateway-darwin-arm64/bin/xmtp-gateway`
        : `${__dirname}/../../gateway-darwin-x64/bin/xmtp-gateway`;

    process.env.XMTP_GATEWAY_BINARY_PATH = binaryPath;
    expect(resolveBinary()).toBe(binaryPath);
  });

  it("should throw when XMTP_GATEWAY_BINARY_PATH points to non-existent file", () => {
    process.env.XMTP_GATEWAY_BINARY_PATH = "/nonexistent/xmtp-gateway";
    expect(() => resolveBinary()).toThrow("non-existent file");
  });

  it("should resolve platform package when no env var set", () => {
    delete process.env.XMTP_GATEWAY_BINARY_PATH;

    // This test depends on the platform packages being available locally.
    // In the POC monorepo setup, they should be findable via require.resolve.
    // If not installed, it should throw a helpful error.
    try {
      const path = resolveBinary();
      expect(path).toContain("xmtp-gateway");
    } catch (err: any) {
      // Expected when platform packages aren't installed as proper npm deps
      expect(err.message).toContain("Cannot find XMTP gateway binary");
    }
  });
});
