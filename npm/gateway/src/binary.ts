import { existsSync } from "node:fs";
import os from "node:os";

const PLATFORM_MAP: Record<string, string> = {
  darwin: "darwin",
  linux: "linux",
};

const ARCH_MAP: Record<string, string> = {
  arm64: "arm64",
  x64: "x64",
};

/**
 * Resolves the path to the xmtp-gateway binary for the current platform.
 *
 * Fallback order (same pattern as esbuild):
 * 1. XMTP_GATEWAY_BINARY_PATH environment variable
 * 2. Platform-specific optional dependency package
 */
export function resolveBinary(): string {
  const envPath = process.env.XMTP_GATEWAY_BINARY_PATH;
  if (envPath) {
    if (!existsSync(envPath)) {
      throw new Error(
        `XMTP_GATEWAY_BINARY_PATH points to non-existent file: ${envPath}`,
      );
    }
    return envPath;
  }

  const platform = PLATFORM_MAP[process.platform];
  const arch = ARCH_MAP[os.arch()];

  if (!platform || !arch) {
    throw new Error(
      `Unsupported platform: ${process.platform}-${os.arch()}. ` +
        `Supported: darwin-arm64, darwin-x64, linux-arm64, linux-x64`,
    );
  }

  const packageName = `@xmtp/gateway-${platform}-${arch}`;

  try {
    return require.resolve(`${packageName}/bin/xmtp-gateway`);
  } catch {
    throw new Error(
      `Cannot find XMTP gateway binary for ${process.platform}-${os.arch()}.\n` +
        `Expected package: ${packageName}\n\n` +
        `Make sure optional dependencies are installed:\n` +
        `  npm install (without --no-optional)\n\n` +
        `Or set XMTP_GATEWAY_BINARY_PATH to a custom binary path.`,
    );
  }
}
