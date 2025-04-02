#!/bin/bash
set -euo pipefail

# Work always from the root directory
script_dir=$(dirname "$(realpath "$0")")
repo_root=$(realpath "${script_dir}/../../")
cd "${repo_root}"

export RPC_URL=http://localhost:8545

source dev/local.env

echo "Loading anvil state"
anvil --state "${1}" &>/dev/null &
ANVIL_PID=$!

echo "Registering local node 1"
dev/cmd/register-local-node-1

echo "Registering local node 2"
dev/cmd/register-local-node-2

echo "Stopping anvil"
kill $ANVIL_PID && sleep 5
