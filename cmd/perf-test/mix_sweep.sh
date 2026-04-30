#!/usr/bin/env bash
#
# mix_sweep.sh — Escalate mixed CB-user workload through DAU tiers
#
# Runs the CB user simulation at each tier for 3 minutes, saves snapshots,
# stops if rate limited or >10% aggregate error rate.
#
set -euo pipefail

BINARY="${1:-/tmp/perf-test}"
ADDR="${2:-grpc.testnet-staging.xmtp.network:443}"
DUR="180s"  # 3 minutes per tier
COOLDOWN=45 # seconds between tiers

TIMESTAMP="$(date -u +%Y%m%dT%H%M%SZ)"
SNAP_DIR="snapshots/mix_${TIMESTAMP}"
mkdir -p "${SNAP_DIR}"

# DAU tiers to test
DAU_TIERS=(10000 25000 50000 75000 100000 150000 200000 300000 500000)

log() {
  echo "[$(date -u +%H:%M:%S)] $*"
}

log "Mixed CB-user workload escalation test"
log "Target: ${ADDR} | Duration per tier: ${DUR}"
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
  "type": "mixed_cb_user_workload"
}
EOF

STOPPED_TIER=""

for dau in "${DAU_TIERS[@]}"; do
  dau_k=$((dau / 1000))
  out_file="${SNAP_DIR}/mix_${dau_k}k.json"

  log "═══════════════════════════════════════"
  log "Starting ${dau_k}K DAU tier"
  log "═══════════════════════════════════════"

  "${BINARY}" \
    -addr "${ADDR}" \
    -mix "${dau}" \
    -dur "${DUR}" \
    -out "${out_file}" \
    2>&1 | tee -a "${SNAP_DIR}/full.log" || true

  # Check for rate limiting or high error rate
  if [ -f "${out_file}" ]; then
    rate_limited=$(python3 -c "import json; d=json.load(open('${out_file}')); print(d.get('rate_limited', False))" 2>/dev/null || echo "False")
    error_pct=$(python3 -c "import json; d=json.load(open('${out_file}')); print(f'{d[\"aggregate\"][\"error_pct\"]:.1f}')" 2>/dev/null || echo "0.0")

    log "Tier ${dau_k}K: error_pct=${error_pct}%, rate_limited=${rate_limited}"

    if [ "${rate_limited}" = "True" ]; then
      log "⚠ RATE LIMITED at ${dau_k}K DAU — stopping escalation"
      STOPPED_TIER="${dau_k}K (rate limited)"
      break
    fi

    # Check if error rate > 10%
    high_err=$(python3 -c "print('yes' if float('${error_pct}') > 10.0 else 'no')" 2>/dev/null || echo "no")
    if [ "${high_err}" = "yes" ]; then
      log "⚠ HIGH ERROR RATE (${error_pct}%) at ${dau_k}K DAU — stopping escalation"
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
print('╔═══════════╦═══════════╦═══════════╦═══════════╦═══════════╦════════╦═══════════╗')
print('║ DAU       ║ Target    ║ Actual    ║ Total OK  ║ Total Err ║ Err%   ║ Rate Ltd  ║')
print('╠═══════════╬═══════════╬═══════════╬═══════════╬═══════════╬════════╬═══════════╣')
for t in tiers:
    dau = t['dau']
    target = t['target_rps_total']
    actual = t['aggregate']['total_rps']
    ok = t['aggregate']['total_count'] - t['aggregate']['total_errors']
    err = t['aggregate']['total_errors']
    epct = t['aggregate']['error_pct']
    rl = 'YES' if t.get('rate_limited') else 'no'
    dau_str = f'{dau//1000}K'
    print(f'║ {dau_str:<9} ║ {target:>7.0f}/s ║ {actual:>7.0f}/s ║ {ok:>9,} ║ {err:>9,} ║ {epct:>5.1f}% ║ {rl:<9} ║')
print('╚═══════════╩═══════════╩═══════════╩═══════════╩═══════════╩════════╩═══════════╝')

# Per-API breakdown for each tier
print()
print('Per-API Detail:')
for t in tiers:
    dau_k = t['dau'] // 1000
    print(f'\\n--- {dau_k}K DAU ---')
    for api in t.get('apis', []):
        print(f\"  {api['name']:<22} target=N/A  actual={api['rps']:>7.1f}/s  P50={api.get('p50_latency_ms',0):>6.1f}ms  P99={api.get('p99_latency_ms',0):>7.1f}ms  err={api.get('error_pct',0):.1f}%\")

stopped = '${STOPPED_TIER}'
if stopped:
    print(f'\\n⚠ Escalation stopped at: {stopped}')
else:
    print(f'\\nAll tiers completed successfully.')
"

log "Done! Results in ${SNAP_DIR}/"
