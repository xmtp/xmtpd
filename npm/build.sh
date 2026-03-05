#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
XMTPD_DIR="${XMTPD_DIR:-$(cd "${SCRIPT_DIR}/.." && pwd)}"

echo "Building XMTP gateway binaries from ${XMTPD_DIR}..."

build_target() {
  local goos=$1
  local goarch=$2
  local pkg_arch=$3
  local output="${SCRIPT_DIR}/gateway/bin/xmtp-gateway-${goos}-${pkg_arch}"

  echo "  Building ${goos}/${goarch} -> ${output}"
  mkdir -p "$(dirname "${output}")"
  CGO_ENABLED=0 GOOS="${goos}" GOARCH="${goarch}" \
    go build -C "${XMTPD_DIR}" -ldflags="-s -w" -o "${output}" ./cmd/gateway
  echo "  Done: $(ls -lh "${output}" | awk '{print $5}')"
}

build_target darwin arm64 arm64
build_target darwin amd64 x64
build_target linux arm64 arm64
build_target linux amd64 x64

echo ""
echo "All binaries built:"
ls -lh "${SCRIPT_DIR}/gateway/bin/"
