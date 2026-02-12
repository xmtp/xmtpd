#!/usr/bin/env bash
set -euo pipefail

# Syncs the version across all gateway npm packages.
# Usage: ./sync-versions.sh <version>
# Example: ./sync-versions.sh 0.2.0
# Example: ./sync-versions.sh 0.2.0-dev.abc1234

VERSION="${1:?Usage: sync-versions.sh <version>}"
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"

# Validate version string to prevent command injection
if [[ ! "$VERSION" =~ ^[0-9]+\.[0-9]+\.[0-9]+(-[a-zA-Z0-9._]+)?$ ]]; then
  echo "ERROR: Invalid version format: ${VERSION}"
  echo "Expected format: X.Y.Z or X.Y.Z-prerelease (e.g. 0.2.0, 0.2.0-dev.abc1234)"
  exit 1
fi

PLATFORM_PACKAGES=(
  "gateway-darwin-arm64"
  "gateway-darwin-x64"
  "gateway-linux-arm64"
  "gateway-linux-x64"
)

echo "Syncing all gateway packages to version ${VERSION}..."

# Update platform packages — pass version via env var to avoid injection
for pkg in "${PLATFORM_PACKAGES[@]}"; do
  PKG_JSON="${SCRIPT_DIR}/${pkg}/package.json"
  if [[ ! -f "$PKG_JSON" ]]; then
    echo "  WARNING: ${PKG_JSON} not found, skipping"
    continue
  fi
  TARGET_VERSION="$VERSION" TARGET_FILE="$PKG_JSON" node -e '
    const fs = require("fs");
    const pkg = JSON.parse(fs.readFileSync(process.env.TARGET_FILE, "utf8"));
    pkg.version = process.env.TARGET_VERSION;
    fs.writeFileSync(process.env.TARGET_FILE, JSON.stringify(pkg, null, 2) + "\n");
  '
  echo "  ${pkg}: ${VERSION}"
done

# Update main package (version + optionalDependencies)
MAIN_PKG="${SCRIPT_DIR}/gateway/package.json"
if [[ ! -f "$MAIN_PKG" ]]; then
  echo "ERROR: ${MAIN_PKG} not found"
  exit 1
fi
TARGET_VERSION="$VERSION" TARGET_FILE="$MAIN_PKG" node -e '
  const fs = require("fs");
  const pkg = JSON.parse(fs.readFileSync(process.env.TARGET_FILE, "utf8"));
  pkg.version = process.env.TARGET_VERSION;
  if (pkg.optionalDependencies) {
    for (const dep of Object.keys(pkg.optionalDependencies)) {
      if (dep.startsWith("@xmtp/gateway-")) {
        pkg.optionalDependencies[dep] = process.env.TARGET_VERSION;
      }
    }
  }
  fs.writeFileSync(process.env.TARGET_FILE, JSON.stringify(pkg, null, 2) + "\n");
'
echo "  gateway (main): ${VERSION}"

echo ""
echo "All packages synced to ${VERSION}"
