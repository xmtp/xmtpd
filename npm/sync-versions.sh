#!/usr/bin/env bash
set -euo pipefail

# Syncs the version for the gateway npm package.
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

MAIN_PKG="${SCRIPT_DIR}/gateway/package.json"
if [[ ! -f "$MAIN_PKG" ]]; then
  echo "ERROR: ${MAIN_PKG} not found"
  exit 1
fi
TARGET_VERSION="$VERSION" TARGET_FILE="$MAIN_PKG" node -e '
  const fs = require("fs");
  const pkg = JSON.parse(fs.readFileSync(process.env.TARGET_FILE, "utf8"));
  pkg.version = process.env.TARGET_VERSION;
  fs.writeFileSync(process.env.TARGET_FILE, JSON.stringify(pkg, null, 2) + "\n");
'
echo "gateway: ${VERSION}"
