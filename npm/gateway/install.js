// postinstall script — validates the gateway binary is available
const os = require("os");
const { existsSync } = require("fs");

const PLATFORM_MAP = { darwin: "darwin", linux: "linux", win32: "win32" };
const ARCH_MAP = { arm64: "arm64", x64: "x64" };

const platform = PLATFORM_MAP[process.platform];
const arch = ARCH_MAP[os.arch()];

if (!platform || !arch) {
  console.warn(
    `[@xmtp/gateway] Warning: Unsupported platform ${process.platform}-${os.arch()}.`,
  );
  console.warn(
    `Set XMTP_GATEWAY_BINARY_PATH to a custom binary path to use the gateway.`,
  );
  process.exit(0);
}

const packageName = `@xmtp/gateway-${platform}-${arch}`;
const binaryName = process.platform === "win32" ? "xmtp-gateway.exe" : "xmtp-gateway";

try {
  const binaryPath = require.resolve(`${packageName}/bin/${binaryName}`);
  if (existsSync(binaryPath)) {
    console.log(`[@xmtp/gateway] Binary found: ${packageName}`);
  } else {
    throw new Error("Binary file does not exist");
  }
} catch {
  console.warn(
    `[@xmtp/gateway] Warning: Could not find binary for ${platform}-${arch}.`,
  );
  console.warn(
    `The package ${packageName} may not have been installed.`,
  );
  console.warn(
    `If you used --no-optional, set XMTP_GATEWAY_BINARY_PATH instead.`,
  );
  // Don't fail — this is a warning, not an error
  // The binary might be provided via XMTP_GATEWAY_BINARY_PATH at runtime
}
