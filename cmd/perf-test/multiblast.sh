#!/usr/bin/bash
#
# multiblast.sh — Run N parallel perf-test processes to saturate the server.
# Each process runs blast mode independently; we aggregate results.
#
set -eu

BINARY="${1:-/tmp/perf-test}"
ADDR="${2:-grpc.testnet-staging.xmtp.network:443}"

# Escalation: [num_processes, conc_per_process, conns_per_process]
# Total goroutines = num_processes * conc_per_process
# Total connections = num_processes * conns_per_process
TIERS=(
  "1 100 16"
  "2 100 16"
  "4 100 16"
  "6 100 16"
  "8 100 16"
  "4 200 16"
  "6 200 16"
  "8 200 16"
)

DUR="120s"
COOLDOWN=20

TIMESTAMP="$(date -u +%Y%m%dT%H%M%SZ)"
SNAP_DIR="snapshots/multiblast_${TIMESTAMP}"
mkdir -p "${SNAP_DIR}"

log() {
  echo "[$(date -u +%H:%M:%S)] $*"
}

log "MULTI-PROCESS BLAST — find true server limit"
log "Target: ${ADDR} | Duration per tier: ${DUR}"
log "Output: ${SNAP_DIR}/"
echo ""

cat > "${SNAP_DIR}/meta.json" << EOF
{
  "timestamp": "${TIMESTAMP}",
  "addr": "${ADDR}",
  "duration_per_tier": "${DUR}",
  "type": "multi_process_blast"
}
EOF

STOPPED=""

for tier in "${TIERS[@]}"; do
  read -r NPROCS CONC CONNS <<< "${tier}"
  TOTAL_CONC=$((NPROCS * CONC))
  TOTAL_CONNS=$((NPROCS * CONNS))
  TIER_TAG="${NPROCS}x${CONC}c${CONNS}"
  TIER_DIR="${SNAP_DIR}/${TIER_TAG}"
  mkdir -p "${TIER_DIR}"

  log "================================================================"
  log "  ${NPROCS} procs × ${CONC} goroutines × ${CONNS} conns"
  log "  Total: ${TOTAL_CONC} goroutines, ${TOTAL_CONNS} connections"
  log "================================================================"

  # Launch N processes in parallel
  PIDS=()
  for i in $(seq 1 ${NPROCS}); do
    OUT="${TIER_DIR}/proc_${i}.json"
    "${BINARY}" \
      -addr "${ADDR}" \
      -blast "${CONC}" \
      -conn "${CONNS}" \
      -dur "${DUR}" \
      -out "${OUT}" \
      > "${TIER_DIR}/proc_${i}.log" 2>&1 &
    PIDS+=($!)
  done

  log "  Launched ${NPROCS} processes: PIDs ${PIDS[*]}"

  # Progress monitor — poll total throughput every 15s
  START_TIME=$(date +%s)
  while true; do
    sleep 15
    ELAPSED=$(( $(date +%s) - START_TIME ))

    # Check if any process is still running
    ANY_RUNNING=false
    for pid in "${PIDS[@]}"; do
      if kill -0 "$pid" 2>/dev/null; then
        ANY_RUNNING=true
        break
      fi
    done

    if [ "$ANY_RUNNING" = false ]; then
      break
    fi

    # Sum up OK counts from logs
    TOTAL_OK=0
    for i in $(seq 1 ${NPROCS}); do
      LAST_LINE=$(grep -oP '\d+ OK' "${TIER_DIR}/proc_${i}.log" 2>/dev/null | tail -1 | grep -oP '^\d+' || echo 0)
      TOTAL_OK=$((TOTAL_OK + LAST_LINE))
    done

    if [ "$ELAPSED" -gt 0 ] && [ "$TOTAL_OK" -gt 0 ]; then
      RPS=$((TOTAL_OK / ELAPSED))
      log "  [${ELAPSED}s] ~${TOTAL_OK} OK across ${NPROCS} procs | ~${RPS} req/s"
    fi
  done

  # Wait for all to finish
  for pid in "${PIDS[@]}"; do
    wait "$pid" 2>/dev/null || true
  done

  # Aggregate results
  python3 -c "
import json, glob, os

tier_dir = '${TIER_DIR}'
nprocs = ${NPROCS}
conc = ${CONC}
conns = ${CONNS}

total_ok = 0
total_err = 0
total_rps = 0.0
api_rps = {}
api_p50 = {}
api_p99 = {}
api_avg = {}
api_err = {}

for f in sorted(glob.glob(os.path.join(tier_dir, 'proc_*.json'))):
    try:
        with open(f) as fh:
            d = json.load(fh)
        total_ok += d['aggregate']['total_count'] - d['aggregate']['total_errors']
        total_err += d['aggregate']['total_errors']
        total_rps += d['aggregate']['total_rps']
        for api in d.get('apis', []):
            name = api['name']
            api_rps[name] = api_rps.get(name, 0) + api['rps']
            api_avg.setdefault(name, []).append(api.get('avg_latency_ms', 0))
            api_p50.setdefault(name, []).append(api.get('p50_latency_ms', 0))
            api_p99.setdefault(name, []).append(api.get('p99_latency_ms', 0))
            api_err[name] = api_err.get(name, 0) + api.get('error_pct', 0)
    except Exception as e:
        print(f'  Error reading {f}: {e}')

err_pct = total_err / (total_ok + total_err) * 100 if (total_ok + total_err) > 0 else 0

# Save aggregate
agg = {
    'nprocs': nprocs, 'conc_per_proc': conc, 'conns_per_proc': conns,
    'total_concurrency': nprocs * conc, 'total_connections': nprocs * conns,
    'total_rps': total_rps, 'total_ok': total_ok, 'total_err': total_err,
    'error_pct': err_pct,
    'apis': {name: {'rps': api_rps.get(name, 0),
                     'avg_ms': sum(api_avg.get(name, [0]))/max(len(api_avg.get(name, [0])),1),
                     'p50_ms': max(api_p50.get(name, [0])),
                     'p99_ms': max(api_p99.get(name, [0]))}
             for name in api_rps}
}
with open(os.path.join(tier_dir, 'aggregate.json'), 'w') as fh:
    json.dump(agg, fh, indent=2)

print(f'  {nprocs}×{conc}g×{conns}c = {total_rps:,.0f} req/s | {total_ok:,} OK | {total_err:,} err ({err_pct:.1f}%)')
for name in ['GetInboxIds', 'QueryEnvelopes', 'GetNewestEnvelope', 'GroupMessage-256B']:
    if name in api_rps:
        ap50 = max(api_p50.get(name, [0]))
        ap99 = max(api_p99.get(name, [0]))
        print(f'    {name:<22} {api_rps[name]:>7.0f}/s  p50={ap50:>6.0f}ms  p99={ap99:>6.0f}ms')

# DAU equivalent
if total_rps > 0:
    dau = total_rps / (470.0 / 86400.0)
    print(f'  CB Wallet DAU equivalent: ~{dau:,.0f} ({dau/1000:.0f}K)')
" 2>&1 | tee -a "${SNAP_DIR}/full.log"

  # Check error rate
  ERR_PCT=$(python3 -c "
import json
d=json.load(open('${TIER_DIR}/aggregate.json'))
print(f'{d[\"error_pct\"]:.1f}')
" 2>/dev/null || echo "0.0")

  high_err=$(python3 -c "print('yes' if float('${ERR_PCT}') > 25.0 else 'no')" 2>/dev/null || echo "no")
  if [ "${high_err}" = "yes" ]; then
    log "HIGH ERROR RATE (${ERR_PCT}%) — stopping"
    STOPPED="${TIER_TAG} (${ERR_PCT}% errors)"
    break
  fi

  log "Cooling down ${COOLDOWN}s..."
  sleep ${COOLDOWN}
  echo ""
done

# Final summary
log "Generating summary..."
python3 -c "
import json, glob, os

snap_dir = '${SNAP_DIR}'

print()
print('==============================================================================')
print('  MULTI-PROCESS BLAST — TRUE SERVER CAPACITY')
print('==============================================================================')
print()
print('+------------------+----------+-----------+-----------+-----------+--------+')
print('| Config           | Tot Conc | Actual/s  | Total OK  | Total Err | Err%   |')
print('+------------------+----------+-----------+-----------+-----------+--------+')

peak_rps = 0
peak_cfg = ''

for d in sorted(os.listdir(snap_dir)):
    agg_path = os.path.join(snap_dir, d, 'aggregate.json')
    if not os.path.exists(agg_path):
        continue
    with open(agg_path) as f:
        a = json.load(f)
    cfg = f'{a[\"nprocs\"]}×{a[\"conc_per_proc\"]}g×{a[\"conns_per_proc\"]}c'
    tc = a['total_concurrency']
    rps = a['total_rps']
    ok = a['total_ok']
    err = a['total_err']
    epct = a['error_pct']
    marker = ''
    if rps > peak_rps:
        peak_rps = rps
        peak_cfg = cfg
        marker = ' <-- PEAK'
    print(f'| {cfg:<16} | {tc:>8} | {rps:>7,.0f}/s | {ok:>9,} | {err:>9,} | {epct:>5.1f}% |{marker}')

print('+------------------+----------+-----------+-----------+-----------+--------+')

if peak_rps > 0:
    dau = peak_rps / (470.0 / 86400.0)
    dau_with_sub = dau * 0.85  # account for missing 15% SubscribeTopics
    print(f'\nPeak throughput: {peak_rps:,.0f} req/s ({peak_cfg})')
    print(f'CB Wallet DAU equiv (excl streaming): ~{dau:,.0f} ({dau/1000:.0f}K)')
    print(f'CB Wallet DAU equiv (adj for 15% streaming): ~{dau_with_sub:,.0f} ({dau_with_sub/1000:.0f}K)')
" 2>&1

log "Done! Results in ${SNAP_DIR}/"
