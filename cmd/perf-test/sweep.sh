#!/usr/bin/env bash
#
# sweep.sh — Concurrency sweep runner for D14N perf tests
#
# Runs each unary API test at multiple concurrency levels,
# saving JSON snapshots at every step. Then runs streaming
# catch-up tests at multiple stream counts.
#
# Usage:
#   ./sweep.sh [perf-test-binary] [target-addr]
#
# Output:
#   snapshots/<timestamp>/ directory with per-step JSON files
#   snapshots/<timestamp>/summary.json — combined results

set -euo pipefail

BINARY="${1:-/tmp/perf-test}"
ADDR="${2:-grpc.testnet-staging.xmtp.network:443}"
DUR="${DUR:-15s}"
NODE_ID="${NODE_ID:-100}"

TIMESTAMP="$(date -u +%Y%m%dT%H%M%SZ)"
SNAP_DIR="snapshots/${TIMESTAMP}"
mkdir -p "${SNAP_DIR}"

# Concurrency levels for unary tests
CONCURRENCY_LEVELS=(1 4 8 16 32 64 128 256)

# Stream counts for catch-up tests
STREAM_LEVELS=(4 8 16 32 64 128)

# Messages per catch-up test (total pre-published)
CATCHUP_MSGS=500

# Unary tests
UNARY_TESTS=(
  "QueryEnvelopes"
  "GetInboxIds"
  "GetNewestEnvelope"
  "Welcome"
  "GroupMessage-256B"
  "GroupMessage-512B"
  "GroupMessage-1KB"
  "GroupMessage-5KB"
)

log() {
  echo "[$(date -u +%H:%M:%S)] $*"
}

# Write metadata
cat > "${SNAP_DIR}/meta.json" << EOF
{
  "timestamp": "${TIMESTAMP}",
  "addr": "${ADDR}",
  "node_id": ${NODE_ID},
  "duration_per_test": "${DUR}",
  "concurrency_levels": [$(IFS=,; echo "${CONCURRENCY_LEVELS[*]}")],
  "stream_levels": [$(IFS=,; echo "${STREAM_LEVELS[*]}")],
  "catchup_msgs": ${CATCHUP_MSGS},
  "unary_tests": $(printf '%s\n' "${UNARY_TESTS[@]}" | jq -R . | jq -s .)
}
EOF

log "Starting sweep → ${SNAP_DIR}"
log "Target: ${ADDR} (node ${NODE_ID}), duration: ${DUR}"
log "Unary tests: ${#UNARY_TESTS[@]} × ${#CONCURRENCY_LEVELS[@]} concurrency levels"
log "Streaming catch-up: ${#STREAM_LEVELS[@]} stream levels"
echo ""

ALL_RESULTS=()

# Phase 1: Unary tests
for test_name in "${UNARY_TESTS[@]}"; do
  for conc in "${CONCURRENCY_LEVELS[@]}"; do
    # Use half-concurrency for connections (min 1)
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
      2>&1 | tee -a "${SNAP_DIR}/full.log" | grep -E "Requests/sec|ERROR|Err%"; then
      ALL_RESULTS+=("${out_file}")
    else
      log "  ⚠ ${test_name} @ c=${conc} failed"
    fi

    echo ""
  done
done

# Phase 2: Streaming catch-up tests
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
    2>&1 | tee -a "${SNAP_DIR}/full.log" | grep -E "Throughput|Published|Received|ERROR|Err%"; then
    ALL_RESULTS+=("${out_file}")
  else
    log "  ⚠ SubscribeTopics-Catchup @ s=${streams} failed"
  fi

  echo ""
done

# Phase 3: Merge all results into summary.json
log "Merging ${#ALL_RESULTS[@]} result files..."
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
        # data is a list of testResult objects
        for r in data:
            # Add concurrency info from filename
            base = os.path.basename(f).replace('.json', '')
            r['snapshot_file'] = base
            results.append(r)
    except Exception as e:
        print(f'  skip {f}: {e}', file=sys.stderr)

with open(os.path.join(snap_dir, 'summary.json'), 'w') as fh:
    json.dump(results, fh, indent=2)

print(f'  {len(results)} results merged into summary.json')
"

log "Sweep complete! Results in ${SNAP_DIR}/"
log "Run: python3 cmd/perf-test/report.py ${SNAP_DIR} to generate report"
