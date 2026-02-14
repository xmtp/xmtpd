import { describe, it, expect, beforeEach, afterEach } from "vitest";
import { resolveBinary } from "../src/binary.js";
import path from "node:path";

const localBinaryPath = path.resolve(
  __dirname,
  `../../gateway-${process.platform}-${process.arch}/bin/xmtp-gateway`,
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

  it("should resolve platform package or throw helpful error", () => {
    delete process.env.XMTP_GATEWAY_BINARY_PATH;
    try {
      const resolved = resolveBinary();
      expect(resolved).toContain("xmtp-gateway");
    } catch (err: any) {
      expect(err.message).toContain("Cannot find XMTP gateway binary");
    }
  });
});
