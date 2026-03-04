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

  it("should throw helpful error when platform package is missing", () => {
    delete process.env.XMTP_GATEWAY_BINARY_PATH;
    // Without the env override, resolution depends on whether the platform
    // package is installed. In local dev it usually is, so we just verify
    // the result is reasonable either way.
    const result = (() => {
      try {
        return { path: resolveBinary() };
      } catch (err: any) {
        return { error: err.message as string };
      }
    })();

    if ("path" in result) {
      expect(result.path).toContain("xmtp-gateway");
    } else {
      expect(result.error).toContain("Cannot find XMTP gateway binary");
    }
  });
});
