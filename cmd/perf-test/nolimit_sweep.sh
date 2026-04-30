#!/usr/bin/env bash
#
# nolimit_sweep.sh — No-ratelimit CB user workload escalation
#
# Rate limits OFF, single node (Ohio only). Push until it breaks.
#
set -euo pipefail

BINARY="${1:-/tmp/perf-test}"
ADDR="${2:-grpc.testnet-staging.xmtp.network:443}"
DUR="180s"  # 3 minutes per tier
COOLDOWN=30 # seconds between tiers

TIMESTAMP="$(date -u +%Y%m%dT%H%M%SZ)"
SNAP_DIR="snapshots/nolimit_${TIMESTAMP}"
mkdir -p "${SNAP_DIR}"

# Aggressive DAU tiers — go until it breaks
DAU_TIERS=(10000 50000 100000 200000 300000 500000 750000 1000000 1500000 2000000)

log() {
  echo "[$(date -u +%H:%M:%S)] $*"
}

log "NO-RATELIMIT CB user workload escalation"
log "Target: ${ADDR} | Duration per tier: ${DUR}"
log "Config: Rate limits OFF, single Ohio node"
log "Tiers: ${DAU_TIERS[*]}"
log "Output: ${SNAP_DIR}/"
echo ""

# Save metadata
cat > "${SNAP_DIR}/meta.json" << EOF
{
  "timestamp": "${TIMESTAMP}",
  "addr": "${ADDR}",
  "duration_per_tier": "${DUR}",
  "cooldown_between_tiers": ${COOLDOWN},
  "dau_tiers": [$(IFS=,; echo "${DAU_TIERS[*]}")],
  "type": "nolimit_cb_user_workload",
  "config": "rate_limits_off, single_ohio_node"
}
EOF

STOPPED_TIER=""

for dau in "${DAU_TIERS[@]}"; do
  dau_k=$((dau / 1000))
  out_file="${SNAP_DIR}/mix_${dau_k}k.json"

  log "================================================================"
  log "Starting ${dau_k}K DAU tier (target $(python3 -c "print(f'{${dau}*470/86400:.0f}')")  req/s)"
  log "================================================================"

  "${BINARY}" \
    -addr "${ADDR}" \
    -mix "${dau}" \
    -dur "${DUR}" \
    -out "${out_file}" \
    2>&1 | tee -a "${SNAP_DIR}/full.log"

  # Check for rate limiting or high error rate
  if [ -f "${out_file}" ]; then
    rate_limited=$(python3 -c "import json; d=json.load(open('${out_file}')); print(d.get('rate_limited', False))" 2>/dev/null || echo "False")
    error_pct=$(python3 -c "import json; d=json.load(open('${out_file}')); print(f'{d[\"aggregate\"][\"error_pct\"]:.1f}')" 2>/dev/null || echo "0.0")
    actual_rps=$(python3 -c "import json; d=json.load(open('${out_file}')); print(f'{d[\"aggregate\"][\"total_rps\"]:.0f}')" 2>/dev/null || echo "0")

    log "Tier ${dau_k}K: actual=${actual_rps} req/s, error_pct=${error_pct}%, rate_limited=${rate_limited}"

    if [ "${rate_limited}" = "True" ]; then
      log "RATE LIMITED at ${dau_k}K DAU — stopping escalation"
      STOPPED_TIER="${dau_k}K (rate limited)"
      break
    fi

    # Stop on >15% error rate (slightly more tolerant for high load)
    high_err=$(python3 -c "print('yes' if float('${error_pct}') > 15.0 else 'no')" 2>/dev/null || echo "no")
    if [ "${high_err}" = "yes" ]; then
      log "HIGH ERROR RATE (${error_pct}%) at ${dau_k}K DAU — stopping escalation"
      STOPPED_TIER="${dau_k}K (${error_pct}% errors)"
      break
    fi
  fi

  log "Cooling down ${COOLDOWN}s..."
  sleep ${COOLDOWN}
  echo ""
done

# Generate combined summary
log "Generating combined summary..."
python3 -c "
import json, glob, os

snap_dir = '${SNAP_DIR}'
tiers = []
for f in sorted(glob.glob(os.path.join(snap_dir, 'mix_*.json'))):
    try:
        with open(f) as fh:
            tiers.append(json.load(fh))
    except:
        pass

with open(os.path.join(snap_dir, 'all_tiers.json'), 'w') as fh:
    json.dump(tiers, fh, indent=2)

print()
print('==============================================================================')
print('  NO-RATELIMIT CB USER WORKLOAD RESULTS (single Ohio node)')
print('==============================================================================')
print()
print('+-----------+-----------+-----------+-----------+-----------+--------+---------+')
print('| DAU       | Target    | Actual    | Total OK  | Total Err | Err%   | RateLtd |')
print('+-----------+-----------+-----------+-----------+-----------+--------+---------+')
for t in tiers:
    dau = t['dau']
    target = t['target_rps_total']
    actual = t['aggregate']['total_rps']
    ok = t['aggregate']['total_count'] - t['aggregate']['total_errors']
    err = t['aggregate']['total_errors']
    epct = t['aggregate']['error_pct']
    rl = 'YES' if t.get('rate_limited') else 'no'
    dau_str = f'{dau//1000}K'
    print(f'| {dau_str:<9} | {target:>7.0f}/s | {actual:>7.0f}/s | {ok:>9,} | {err:>9,} | {epct:>5.1f}% | {rl:<7} |')
print('+-----------+-----------+-----------+-----------+-----------+--------+---------+')

# Per-API latency summary for highest completed tier
if tiers:
    t = tiers[-1]
    dau_k = t['dau'] // 1000
    print(f'\nPer-API Detail at peak tier ({dau_k}K DAU):')
    print('+------------------------+-----------+----------+----------+----------+--------+')
    print('| API                    | Actual/s  | Avg(ms)  | P50(ms)  | P99(ms)  | Err%   |')
    print('+------------------------+-----------+----------+----------+----------+--------+')
    for api in t.get('apis', []):
        print(f'| {api[\"name\"]:<22} | {api[\"rps\"]:>7.1f}/s | {api.get(\"avg_latency_ms\",0):>8.2f} | {api.get(\"p50_latency_ms\",0):>8.2f} | {api.get(\"p99_latency_ms\",0):>8.2f} | {api.get(\"error_pct\",0):>5.1f}% |')
    print('+------------------------+-----------+----------+----------+----------+--------+')

stopped = '${STOPPED_TIER}'
if stopped:
    print(f'\nEscalation stopped at: {stopped}')
else:
    print(f'\nAll tiers completed successfully!')
"

log "Done! Results in ${SNAP_DIR}/"
