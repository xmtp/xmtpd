#!/usr/bin/env bash
set -euo pipefail

BINARY=/tmp/perf-test
TARGET="grpc.testnet-staging.xmtp.network:443"
DIR="snapshots/blast_8vcpu_$(date -u +%Y%m%dT%H%M%SZ)"
mkdir -p "$DIR"

DURATION=180
CONNS=32

echo "=== 8-vCPU BLAST SWEEP ==="
echo "Target: $TARGET | Duration: ${DURATION}s | Connections: $CONNS"
echo "Output dir: $DIR"
echo ""

for CONC in 100 200 400 800 1600; do
  TAG="blast_c${CONC}_conn${CONNS}"
  echo "--- $TAG ($(date -u +%H:%M:%S)) ---"
  $BINARY -addr "$TARGET" -blast $CONC -conn $CONNS -dur ${DURATION}s 2>&1 | tee "$DIR/${TAG}.txt"
  echo ""
  sleep 10
done

echo "=== MULTI-PROCESS BLAST ==="
for NPROCS in 2 4 8; do
  TAG="multi_p${NPROCS}_c200_conn16"
  echo "--- $TAG ($(date -u +%H:%M:%S)) ---"
  PIDS=()
  for i in $(seq 1 $NPROCS); do
    $BINARY -addr "$TARGET" -blast 200 -conn 16 -dur ${DURATION}s > "$DIR/${TAG}_proc${i}.txt" 2>&1 &
    PIDS+=($!)
  done
  for pid in "${PIDS[@]}"; do
    wait $pid || true
  done
  TOTAL=0
  for i in $(seq 1 $NPROCS); do
    RPS=$(grep "^Actual:" "$DIR/${TAG}_proc${i}.txt" | grep -oP '\d+' | head -1)
    printf "  proc%d: %s req/s\n" "$i" "${RPS:-0}"
    TOTAL=$((TOTAL + ${RPS:-0}))
  done
  printf "  >>> TOTAL: %d req/s (%d processes)\n\n" "$TOTAL" "$NPROCS"
  sleep 10
done

echo "=== DONE ==="
echo "Results in: $DIR"
