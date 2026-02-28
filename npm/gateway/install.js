// postinstall — validates the gateway binary is available
const os = require("os");
const { existsSync } = require("fs");

const PLATFORM_MAP = { darwin: "darwin", linux: "linux" };
const ARCH_MAP = { arm64: "arm64", x64: "x64" };

const platform = PLATFORM_MAP[process.platform];
const arch = ARCH_MAP[os.arch()];

if (!platform || !arch) {
  console.warn(
    `[@xmtp/gateway] Unsupported platform ${process.platform}-${os.arch()}.`,
  );
  console.warn("Set XMTP_GATEWAY_BINARY_PATH to a custom binary path.");
  process.exit(0);
}

const packageName = `@xmtp/gateway-${platform}-${arch}`;

try {
  const binaryPath = require.resolve(`${packageName}/bin/xmtp-gateway`);
  if (existsSync(binaryPath)) {
    console.log(`[@xmtp/gateway] Binary found: ${packageName}`);
  } else {
    throw new Error("not found");
  }
} catch {
  console.warn(
    `[@xmtp/gateway] Could not find binary for ${platform}-${arch}.`,
  );
  console.warn(
    `If you used --no-optional, set XMTP_GATEWAY_BINARY_PATH instead.`,
  );
}
