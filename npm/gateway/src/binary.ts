import { existsSync } from "node:fs";
import os from "node:os";
import path from "node:path";

const PLATFORM_MAP: Record<string, string> = {
  darwin: "darwin",
  linux: "linux",
};

const ARCH_MAP: Record<string, string> = {
  arm64: "arm64",
  x64: "x64",
};

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

  const binaryPath = path.join(
    __dirname,
    "..",
    "bin",
    `xmtp-gateway-${platform}-${arch}`,
  );

  if (!existsSync(binaryPath)) {
    throw new Error(
      `Cannot find XMTP gateway binary at ${binaryPath}.\n` +
        `Run npm/build.sh to build the binaries, ` +
        `or set XMTP_GATEWAY_BINARY_PATH to a custom binary path.`,
    );
  }

  return binaryPath;
}
