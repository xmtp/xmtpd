#!/bin/bash

set -euo pipefail

error() {
  echo "Error: $1" >&2
  exit 1
}

# Get the directory where the script is located
SCRIPT_DIR=$(dirname "$(realpath "$0")")

TOP_LEVEL_DIR=$(realpath "${SCRIPT_DIR}/../../.." 2>/dev/null) || error "Failed to resolve top-level directory"

[ -d "$TOP_LEVEL_DIR" ] || error "Top level directory not found: $TOP_LEVEL_DIR"

cd "$TOP_LEVEL_DIR" || error "Failed to change to top level directory"

ENV_FILE="./dev/local.env"
[ -f "$ENV_FILE" ] || error "Environment file not found: $ENV_FILE"
[ -r "$ENV_FILE" ] || error "Environment file not readable: $ENV_FILE"
. "$ENV_FILE"

# a subset of all of them
REQUIRED_VARS=(
  "XMTPD_SIGNER_PRIVATE_KEY"
  "XMTPD_DB_WRITER_CONNECTION_STRING"
  "XMTPD_APP_CHAIN_RPC_URL"
  "XMTPD_SETTLEMENT_CHAIN_RPC_URL"
)

for var in "${REQUIRED_VARS[@]}"; do
  [ -n "${!var:-}" ] || error "Required environment variable not set: $var"
done

# Display and validate XMTPD variables
XMTPD_VARS=$(env | grep XMTPD) || error "No XMTPD environment variables found"
echo "$XMTPD_VARS"