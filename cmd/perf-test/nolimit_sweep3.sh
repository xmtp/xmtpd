#!/usr/bin/bash
#
# nolimit_sweep3.sh — Multi-connection sweep. FIND THE BREAKING POINT.
#
set -eu

BINARY="${1:-/tmp/perf-test}"
ADDR="${2:-grpc.testnet-staging.xmtp.network:443}"
DUR="180s"
COOLDOWN=30
CONNS=16  # 16 gRPC connections to avoid client-side HTTP/2 bottleneck

TIMESTAMP="$(date -u +%Y%m%dT%H%M%SZ)"
SNAP_DIR="snapshots/nolimit3_${TIMESTAMP}"
mkdir -p "${SNAP_DIR}"

# Start at 200K, jump in big steps
DAU_TIERS=(200000 500000 1000000 1500000 2000000 3000000 5000000)

log() {
  echo "[$(date -u +%H:%M:%S)] $*"
}

log "MULTI-CONN NO-LIMIT SWEEP — ${CONNS} gRPC connections"
log "Target: ${ADDR} | Duration per tier: ${DUR}"
log "Config: WAF=10M/5min, app RL off, single Ohio node, ${CONNS} conns"
log "Tiers: ${DAU_TIERS[*]}"
log "Output: ${SNAP_DIR}/"
echo ""

cat > "${SNAP_DIR}/meta.json" << EOF
{
  "timestamp": "${TIMESTAMP}",
  "addr": "${ADDR}",
  "duration_per_tier": "${DUR}",
  "cooldown_between_tiers": ${COOLDOWN},
  "connections": ${CONNS},
  "dau_tiers": [$(IFS=,; echo "${DAU_TIERS[*]}")],
  "type": "multi_conn_nolimit",
  "config": "waf_10M_5min, app_rate_limits_off, single_ohio_node, ${CONNS}_conns"
}
EOF

STOPPED_TIER=""

for dau in "${DAU_TIERS[@]}"; do
  dau_k=$((dau / 1000))
  out_file="${SNAP_DIR}/mix_${dau_k}k.json"

  target_rps=$(python3 -c "print(f'{${dau}*470/86400:.0f}')")
  log "================================================================"
  log "  ${dau_k}K DAU  |  target ${target_rps} req/s  |  ${CONNS} conns"
  log "================================================================"

  "${BINARY}" \
    -addr "${ADDR}" \
    -mix "${dau}" \
    -dur "${DUR}" \
    -conn "${CONNS}" \
    -out "${out_file}" \
    2>&1 | tee -a "${SNAP_DIR}/full.log"

  if [ -f "${out_file}" ]; then
    rate_limited=$(python3 -c "import json; d=json.load(open('${out_file}')); print(d.get('rate_limited', False))" 2>/dev/null || echo "False")
    error_pct=$(python3 -c "import json; d=json.load(open('${out_file}')); print(f'{d[\"aggregate\"][\"error_pct\"]:.1f}')" 2>/dev/null || echo "0.0")
    actual_rps=$(python3 -c "import json; d=json.load(open('${out_file}')); print(f'{d[\"aggregate\"][\"total_rps\"]:.0f}')" 2>/dev/null || echo "0")
    p99=$(python3 -c "
import json
d=json.load(open('${out_file}'))
lats = [a.get('p99_latency_ms',0) for a in d.get('apis',[])]
print(f'{max(lats):.0f}' if lats else '0')
" 2>/dev/null || echo "0")

    log "Tier ${dau_k}K: actual=${actual_rps}/s, err=${error_pct}%, p99=${p99}ms, rate_limited=${rate_limited}"

    if [ "${rate_limited}" = "True" ]; then
      log "RATE LIMITED at ${dau_k}K DAU — stopping"
      STOPPED_TIER="${dau_k}K (rate limited)"
      break
    fi

    high_err=$(python3 -c "print('yes' if float('${error_pct}') > 20.0 else 'no')" 2>/dev/null || echo "no")
    if [ "${high_err}" = "yes" ]; then
      log "HIGH ERROR RATE (${error_pct}%) at ${dau_k}K DAU — stopping"
      STOPPED_TIER="${dau_k}K (${error_pct}% errors)"
      break
    fi
  fi

  log "Cooling down ${COOLDOWN}s..."
  sleep ${COOLDOWN}
  echo ""
done

log "Generating summary..."
python3 -c "
import json, glob, os

snap_dir = '${SNAP_DIR}'
tiers = []
for f in sorted(glob.glob(os.path.join(snap_dir, 'mix_*.json')), key=lambda x: int(os.path.basename(x).split('_')[1].replace('k.json',''))):
    try:
        with open(f) as fh:
            tiers.append(json.load(fh))
    except:
        pass

with open(os.path.join(snap_dir, 'all_tiers.json'), 'w') as fh:
    json.dump(tiers, fh, indent=2)

print()
print('==============================================================================')
print('  MULTI-CONN NO-LIMIT RESULTS (${CONNS} conns, single Ohio node)')
print('==============================================================================')
print()
print('+-----------+-----------+-----------+-----------+-----------+--------+-----------+')
print('| DAU       | Target/s  | Actual/s  | Total OK  | Total Err | Err%   | Max P99ms |')
print('+-----------+-----------+-----------+-----------+-----------+--------+-----------+')
for t in tiers:
    dau = t['dau']
    target = t['target_rps_total']
    actual = t['aggregate']['total_rps']
    ok = t['aggregate']['total_count'] - t['aggregate']['total_errors']
    err = t['aggregate']['total_errors']
    epct = t['aggregate']['error_pct']
    p99s = [a.get('p99_latency_ms',0) for a in t.get('apis',[])]
    max_p99 = max(p99s) if p99s else 0
    dau_str = f'{dau//1000}K'
    print(f'| {dau_str:<9} | {target:>7.0f}/s | {actual:>7.0f}/s | {ok:>9,} | {err:>9,} | {epct:>5.1f}% | {max_p99:>7.0f}ms |')
print('+-----------+-----------+-----------+-----------+-----------+--------+-----------+')

# Per-API for last 3 tiers
for t in tiers[-3:]:
    dau_k = t['dau'] // 1000
    print(f'\n--- {dau_k}K DAU per-API ---')
    for api in t.get('apis', []):
        print(f'  {api[\"name\"]:<22} {api[\"rps\"]:>7.1f}/s  avg={api.get(\"avg_latency_ms\",0):>6.1f}ms  p50={api.get(\"p50_latency_ms\",0):>6.1f}ms  p99={api.get(\"p99_latency_ms\",0):>7.1f}ms  err={api.get(\"error_pct\",0):.1f}%')

stopped = '${STOPPED_TIER}'
if stopped:
    print(f'\nEscalation stopped at: {stopped}')
else:
    print(f'\nAll tiers completed!')
"

log "Done! Results in ${SNAP_DIR}/"
