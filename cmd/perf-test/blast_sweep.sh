#!/usr/bin/bash
#
# blast_sweep.sh — Fire as fast as possible, escalate concurrency until it breaks.
#
set -eu

BINARY="${1:-/tmp/perf-test}"
ADDR="${2:-grpc.testnet-staging.xmtp.network:443}"
DUR="120s"
COOLDOWN=20
CONNS=32

TIMESTAMP="$(date -u +%Y%m%dT%H%M%SZ)"
SNAP_DIR="snapshots/blast_${TIMESTAMP}"
mkdir -p "${SNAP_DIR}"

# Escalate concurrency: goroutines that fire as fast as possible
CONC_TIERS=(50 100 200 400 800 1200 1600 2000 3000)

log() {
  echo "[$(date -u +%H:%M:%S)] $*"
}

log "BLAST MODE — max throughput sweep"
log "Target: ${ADDR} | Duration per tier: ${DUR} | Conns: ${CONNS}"
log "Concurrency tiers: ${CONC_TIERS[*]}"
log "Output: ${SNAP_DIR}/"
echo ""

cat > "${SNAP_DIR}/meta.json" << EOF
{
  "timestamp": "${TIMESTAMP}",
  "addr": "${ADDR}",
  "duration_per_tier": "${DUR}",
  "connections": ${CONNS},
  "concurrency_tiers": [$(IFS=,; echo "${CONC_TIERS[*]}")],
  "type": "blast_max_throughput"
}
EOF

STOPPED_TIER=""
PREV_RPS=0

for conc in "${CONC_TIERS[@]}"; do
  out_file="${SNAP_DIR}/blast_${conc}.json"

  log "================================================================"
  log "  ${conc} goroutines × ${CONNS} conns — FULL BLAST"
  log "================================================================"

  "${BINARY}" \
    -addr "${ADDR}" \
    -blast "${conc}" \
    -conn "${CONNS}" \
    -dur "${DUR}" \
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

    log "Conc ${conc}: actual=${actual_rps}/s, err=${error_pct}%, p99=${p99}ms, rate_limited=${rate_limited}"

    if [ "${rate_limited}" = "True" ]; then
      log "RATE LIMITED at conc ${conc} — stopping"
      STOPPED_TIER="${conc} goroutines (rate limited)"
      break
    fi

    high_err=$(python3 -c "print('yes' if float('${error_pct}') > 25.0 else 'no')" 2>/dev/null || echo "no")
    if [ "${high_err}" = "yes" ]; then
      log "HIGH ERROR RATE (${error_pct}%) at conc ${conc} — stopping"
      STOPPED_TIER="${conc} goroutines (${error_pct}% errors)"
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
for f in sorted(glob.glob(os.path.join(snap_dir, 'blast_*.json')), key=lambda x: int(os.path.basename(x).replace('blast_','').replace('.json',''))):
    try:
        with open(f) as fh:
            tiers.append(json.load(fh))
    except:
        pass

with open(os.path.join(snap_dir, 'all_tiers.json'), 'w') as fh:
    json.dump(tiers, fh, indent=2)

print()
print('==============================================================================')
print('  BLAST MODE RESULTS — MAX THROUGHPUT (single Ohio node)')
print('==============================================================================')
print()
print('+-----------+-----------+-----------+-----------+--------+-----------+')
print('| Conc      | Actual/s  | Total OK  | Total Err | Err%   | Max P99ms |')
print('+-----------+-----------+-----------+-----------+--------+-----------+')
peak_rps = 0
peak_conc = 0
for t in tiers:
    conc = t['concurrency']
    actual = t['aggregate']['total_rps']
    ok = t['aggregate']['total_count'] - t['aggregate']['total_errors']
    err = t['aggregate']['total_errors']
    epct = t['aggregate']['error_pct']
    p99s = [a.get('p99_latency_ms',0) for a in t.get('apis',[])]
    max_p99 = max(p99s) if p99s else 0
    if actual > peak_rps:
        peak_rps = actual
        peak_conc = conc
    marker = ' <-- PEAK' if actual == peak_rps and conc == peak_conc else ''
    print(f'| {conc:<9} | {actual:>7.0f}/s | {ok:>9,} | {err:>9,} | {epct:>5.1f}% | {max_p99:>7.0f}ms |{marker}')
print('+-----------+-----------+-----------+-----------+--------+-----------+')

# DAU equivalent
if peak_rps > 0:
    dau_equiv = peak_rps / (470.0 / 86400.0)
    print(f'\nPeak throughput: {peak_rps:,.0f} req/s at {peak_conc} goroutines')
    print(f'CB Wallet DAU equivalent: ~{dau_equiv:,.0f} DAU ({dau_equiv/1000:.0f}K)')
    print(f'(based on 470 API calls/user/day, excluding 15% SubscribeTopics)')

stopped = '${STOPPED_TIER}'
if stopped:
    print(f'\nEscalation stopped at: {stopped}')
else:
    print(f'\nAll tiers completed!')
"

log "Done! Results in ${SNAP_DIR}/"
