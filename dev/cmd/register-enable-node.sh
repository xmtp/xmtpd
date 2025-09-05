#!/usr/bin/env bash
set -euo pipefail

OWNER_ADDR="$1"
SIGNER_PUB="$2"
HTTP_ADDR="$3"
CFG="$4"

: "${PRIVATE_KEY:?PRIVATE_KEY missing}"
: "${RPC_URL:?RPC_URL missing}"

tmp="$(mktemp -t xmtpd-register-XXXX.jsonl)"
trap 'rm -f "$tmp"' EXIT

set +e
dev/cmd/cli \
  --private-key "${PRIVATE_KEY}" \
  --rpc-url "${RPC_URL}" \
  --log-encoding json \
  --config-file "${CFG}" \
  nodes register \
    --owner-address "${OWNER_ADDR}" \
    --signing-key-pub "${SIGNER_PUB}" \
    --http-address "${HTTP_ADDR}" \
| tee "${tmp}"
status=$?
set -e

NODE_ID=""

if [[ $status -ne 0 ]]; then
  if jq -e 'select(.level=="FATAL" and .message=="signing key public key already registered")' "${tmp}" >/dev/null; then
    echo "⚠️  Signing key already registered; resolving node id…"
    NODE_ID="$(
      dev/cmd/cli \
        --rpc-url "${RPC_URL}" \
        --log-encoding json \
        --config-file "${CFG}" \
        nodes get --all \
      | jq -r '
          select(.message=="Getting all nodes")
          | .nodes[]?
          | select(.signing_key_pub=="'"${SIGNER_PUB}"'")
          | .node_id
        ' \
      | head -n1
    )"
    [[ -n "${NODE_ID}" ]] || { echo "ERROR: Could not resolve node id"; exit 1; }
  else
    echo "ERROR: register failed"; exit $status
  fi
else
  NODE_ID="$(jq -r 'select(.message=="Node registered") | .["node-id"]' "${tmp}")"
fi

echo "==> Node id: ${NODE_ID}"

dev/cmd/cli \
  --private-key "${PRIVATE_KEY}" \
  --rpc-url "${RPC_URL}" \
  --log-encoding json \
  --config-file "${CFG}" \
  nodes canonical-network --add --node-id "${NODE_ID}" \
|| echo "ℹ️  Node ${NODE_ID} already enabled"

echo -e "\033[32m✔\033[0m Node ${NODE_ID} registered+enabled"
