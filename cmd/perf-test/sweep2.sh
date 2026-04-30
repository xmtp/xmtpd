#!/usr/bin/env bash
#
# sweep2.sh — Re-run failed tests with cooldowns between API groups
# Writes into the same snapshot directory as the first sweep.
#
set -eu

BINARY="${1:-/tmp/perf-test}"
ADDR="${2:-grpc.testnet-staging.xmtp.network:443}"
DUR="${DUR:-15s}"
NODE_ID="${NODE_ID:-100}"

SNAP_DIR="snapshots/20260416T103919Z"

CONCURRENCY_LEVELS=(1 4 8 16 32 64 128 256)
COOLDOWN=30  # seconds between API groups

# Tests that got 429'd and need re-running
REMAINING_TESTS=(
  "GetNewestEnvelope"
  "Welcome"
  "GroupMessage-256B"
  "GroupMessage-512B"
  "GroupMessage-1KB"
  "GroupMessage-5KB"
)

# Stream levels for catch-up
STREAM_LEVELS=(4 8 16 32 64 128)
CATCHUP_MSGS=500

log() {
  echo "[$(date -u +%H:%M:%S)] $*"
}

log "Re-running failed tests with ${COOLDOWN}s cooldown between API groups"
log "Snapshot dir: ${SNAP_DIR}"
echo ""

# Phase 1: Re-run failed unary tests with cooldowns
for test_name in "${REMAINING_TESTS[@]}"; do
  log "=== ${test_name} (${#CONCURRENCY_LEVELS[@]} levels) ==="

  for conc in "${CONCURRENCY_LEVELS[@]}"; do
    conn=$(( conc / 2 ))
    [[ $conn -lt 1 ]] && conn=1

    out_file="${SNAP_DIR}/${test_name}_c${conc}.json"
    log "▶ ${test_name} @ c=${conc} conn=${conn}"

    if "${BINARY}" \
      -addr "${ADDR}" \
      -node-id "${NODE_ID}" \
      -tests "${test_name}" \
      -c "${conc}" \
      -conn "${conn}" \
      -dur "${DUR}" \
      -out "${out_file}" \
      2>&1 | tee -a "${SNAP_DIR}/full2.log" | grep -E "Requests/sec|ERROR|Err%"; then
      :
    else
      log "  ⚠ ${test_name} @ c=${conc} failed"
    fi
  done

  # Cooldown between API groups to avoid rate limiting
  log "Cooling down ${COOLDOWN}s..."
  sleep ${COOLDOWN}
  echo ""
done

# Phase 2: Re-run streaming catch-up tests
log "=== Streaming Catch-up ==="
for streams in "${STREAM_LEVELS[@]}"; do
  conn=$(( streams / 2 ))
  [[ $conn -lt 1 ]] && conn=1

  out_file="${SNAP_DIR}/SubscribeTopics-Catchup_s${streams}.json"
  log "▶ SubscribeTopics-Catchup @ streams=${streams} msgs=${CATCHUP_MSGS}"

  if "${BINARY}" \
    -addr "${ADDR}" \
    -node-id "${NODE_ID}" \
    -tests "SubscribeTopics-Catchup" \
    -c "${streams}" \
    -conn "${conn}" \
    -pub-rate "${CATCHUP_MSGS}" \
    -dur 30s \
    -out "${out_file}" \
    2>&1 | tee -a "${SNAP_DIR}/full2.log" | grep -E "Throughput|Published|Received|ERROR|Err%"; then
    :
  else
    log "  ⚠ SubscribeTopics-Catchup @ s=${streams} failed"
  fi
done

# Re-merge
log "Re-merging all results..."
python3 -c "
import json, glob, os, sys
snap_dir = '${SNAP_DIR}'
results = []
for f in sorted(glob.glob(os.path.join(snap_dir, '*.json'))):
    if os.path.basename(f) in ('meta.json', 'summary.json'):
        continue
    try:
        with open(f) as fh:
            data = json.load(fh)
        if data is None:
            continue
        for r in data:
            base = os.path.basename(f).replace('.json', '')
            r['snapshot_file'] = base
            results.append(r)
    except Exception as e:
        print(f'  skip {f}: {e}', file=sys.stderr)
with open(os.path.join(snap_dir, 'summary.json'), 'w') as fh:
    json.dump(results, fh, indent=2)
print(f'  {len(results)} results merged')
"

log "Done! Run: python3 report.py ${SNAP_DIR}"
