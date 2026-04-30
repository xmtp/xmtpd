#!/usr/bin/env python3
"""
report.py — Generate markdown report from perf-test sweep snapshots.

Usage:
    python3 report.py <snapshot_dir>
    python3 report.py snapshots/20260417T120000Z

Outputs:
    <snapshot_dir>/REPORT.md
    Prints to stdout as well.
"""

import json
import os
import sys
from collections import defaultdict

# CB user profile: 100K DAU = 544 avg req/s
CB_DAU_TARGET = 100_000
CB_AVG_RPS = 544
CB_PEAK_RPS = 820
CB_PROFILE = {
    "GetInboxIds": 0.65,       # 65% polling
    "QueryEnvelopes": 0.14,    # 14% query/sync
    "SubscribeTopics": 0.15,   # 15% streaming (reconnects)
    "Writes": 0.03,            # 3% writes (Welcome, GroupMessage)
    # For writes, representative is GroupMessage-256B
}


def load_results(snap_dir):
    """Load all per-step JSON files, extract concurrency from filename."""
    results = []
    for f in sorted(os.listdir(snap_dir)):
        if not f.endswith(".json"):
            continue
        if f in ("meta.json", "summary.json"):
            continue
        path = os.path.join(snap_dir, f)
        try:
            with open(path) as fh:
                data = json.load(fh)
            for r in data:
                # Parse concurrency from filename like GetInboxIds_c32.json
                base = f.replace(".json", "")
                parts = base.rsplit("_", 1)
                if len(parts) == 2 and parts[1].startswith("c"):
                    r["_concurrency"] = int(parts[1][1:])
                elif len(parts) == 2 and parts[1].startswith("s"):
                    r["_concurrency"] = int(parts[1][1:])
                r["_file"] = base
                results.append(r)
        except Exception:
            pass
    return results


def group_by_test(results):
    """Group results by test name."""
    groups = defaultdict(list)
    for r in results:
        groups[r["name"]].append(r)
    for name in groups:
        groups[name].sort(key=lambda x: x.get("_concurrency", 0))
    return groups


def fmt_num(n):
    if n >= 1000:
        return f"{n:,.0f}"
    if n >= 100:
        return f"{n:.0f}"
    return f"{n:.1f}"


def estimate_dau_capacity(groups):
    """
    Estimate max supportable DAU from per-API max RPS.

    CB profile says at 100K DAU:
    - GetInboxIds: 65% of 544 = 354 req/s
    - QueryEnvelopes: 14% of 544 = 76 req/s
    - Writes: 3% of 544 = 16 req/s

    The bottleneck API determines max DAU.
    """
    estimates = {}

    for api, fraction in [
        ("GetInboxIds", 0.65),
        ("QueryEnvelopes", 0.14),
        ("GroupMessage-256B", 0.03),
    ]:
        if api not in groups:
            continue
        # Find max RPS with <5% error rate
        max_rps = 0
        for r in groups[api]:
            if r.get("error_pct", 0) < 5:
                max_rps = max(max_rps, r.get("rps", 0))
        if max_rps > 0:
            # At 100K DAU, this API needs fraction * 544 req/s
            # So max DAU = max_rps / (fraction * 544 / 100_000)
            needed_per_100k = fraction * CB_AVG_RPS
            dau_cap = int(max_rps / needed_per_100k * CB_DAU_TARGET)
            estimates[api] = {
                "max_rps": max_rps,
                "needed_at_100k": needed_per_100k,
                "dau_capacity": dau_cap,
            }

    return estimates


def generate_report(snap_dir, results, meta):
    groups = group_by_test(results)
    lines = []
    a = lines.append

    a(f"# D14N Staging Performance Report")
    a(f"")
    a(f"**Date:** {meta.get('timestamp', 'unknown')}")
    a(f"**Target:** `{meta.get('addr', 'unknown')}` (node {meta.get('node_id', '?')})")
    a(f"**Duration per test:** {meta.get('duration_per_test', '?')}")
    a(f"**Machine:** EC2 (8 vCPU, 30GB RAM)")
    a(f"")

    # DAU Capacity Estimate
    estimates = estimate_dau_capacity(groups)
    if estimates:
        a("## Estimated DAU Capacity (CB Wallet Profile)")
        a("")
        a("Based on CB user traffic profile (544 avg req/s at 100K DAU):")
        a("")
        a("| API | Max Clean RPS | Needed @ 100K DAU | Estimated Max DAU |")
        a("|-----|---------------|-------------------|-------------------|")
        bottleneck_dau = float("inf")
        bottleneck_api = ""
        for api, est in sorted(estimates.items()):
            dau_str = f"{est['dau_capacity']:,}"
            a(f"| {api} | {est['max_rps']:,.0f} | {est['needed_at_100k']:.0f} req/s | **{dau_str}** |")
            if est["dau_capacity"] < bottleneck_dau:
                bottleneck_dau = est["dau_capacity"]
                bottleneck_api = api
        a("")
        a(f"> **Bottleneck:** `{bottleneck_api}` — staging can support ~**{bottleneck_dau:,} DAU** before this API saturates.")
        a(f"> Note: This is single-client, single-machine load generation. Real capacity may differ with distributed load.")
        a("")

    # Per-API Scaling Tables
    a("## Per-API Scaling Curves")
    a("")

    # Categorize tests
    read_tests = ["QueryEnvelopes", "GetInboxIds", "GetNewestEnvelope"]
    write_tests = [
        "Welcome",
        "GroupMessage-256B",
        "GroupMessage-512B",
        "GroupMessage-1KB",
        "GroupMessage-5KB",
    ]
    stream_tests = ["SubscribeTopics-Catchup"]

    for category, test_names, label in [
        ("Read APIs", read_tests, "Workers"),
        ("Write APIs", write_tests, "Workers"),
        ("Streaming (Catch-up)", stream_tests, "Streams"),
    ]:
        found = [t for t in test_names if t in groups]
        if not found:
            continue

        a(f"### {category}")
        a("")

        for test_name in found:
            runs = groups[test_name]
            a(f"**{test_name}**")
            a("")
            a(f"| {label} | Count | RPS | Avg (ms) | P50 (ms) | P95 (ms) | P99 (ms) | StdDev | Err% |")
            a(f"|---------|-------|-----|----------|----------|----------|----------|--------|------|")
            for r in runs:
                conc = r.get("_concurrency", "?")
                a(
                    f"| {conc} "
                    f"| {r.get('count', 0):,} "
                    f"| {r.get('rps', 0):,.1f} "
                    f"| {r.get('avg_latency_ms', 0):.2f} "
                    f"| {r.get('p50_latency_ms', 0):.2f} "
                    f"| {r.get('p95_latency_ms', 0):.2f} "
                    f"| {r.get('p99_latency_ms', 0):.2f} "
                    f"| {r.get('stddev_ms', 0):.2f} "
                    f"| {r.get('error_pct', 0):.1f}% |"
                )
            a("")

            # Find saturation point
            peak_rps = 0
            peak_conc = 0
            saturated_at = None
            for r in runs:
                if r.get("error_pct", 0) < 5 and r.get("rps", 0) > peak_rps:
                    peak_rps = r.get("rps", 0)
                    peak_conc = r.get("_concurrency", 0)
            for r in runs:
                if r.get("error_pct", 0) >= 5 and saturated_at is None:
                    saturated_at = r.get("_concurrency", 0)

            note = f"Peak: **{peak_rps:,.0f} RPS** at {peak_conc} {label.lower()}"
            if saturated_at:
                note += f" | Errors spike at {saturated_at}+ {label.lower()}"
            a(f"> {note}")
            a("")

    # Raw numbers for quick copy-paste
    a("## Quick Reference: Peak RPS per API")
    a("")
    a("| API | Peak RPS | @ Workers | P99 @ Peak | Error% |")
    a("|-----|----------|-----------|------------|--------|")
    for test_name in read_tests + write_tests + stream_tests:
        if test_name not in groups:
            continue
        runs = groups[test_name]
        best = max(
            (r for r in runs if r.get("error_pct", 0) < 5),
            key=lambda r: r.get("rps", 0),
            default=None,
        )
        if best:
            a(
                f"| {test_name} "
                f"| {best.get('rps', 0):,.0f} "
                f"| {best.get('_concurrency', '?')} "
                f"| {best.get('p99_latency_ms', 0):.1f}ms "
                f"| {best.get('error_pct', 0):.1f}% |"
            )
    a("")

    # Streaming bug note
    a("## Notes")
    a("")
    a("- **Live streaming is broken on staging** — SubscribeTopics live delivery returns 0 messages. ")
    a("  Catch-up mode works. This is a server-side subscribe worker bug (funnel race condition).")
    a("- Catch-up results above measure the DB→stream delivery path, not the live pub→sub path.")
    a("- Load was generated from a single EC2 instance (8 vCPU). Client-side saturation may cap ")
    a("  some high-concurrency results — distributed load generation would give higher ceilings.")
    a("")

    report = "\n".join(lines)

    report_path = os.path.join(snap_dir, "REPORT.md")
    with open(report_path, "w") as fh:
        fh.write(report)

    print(report)
    print(f"\n--- Report saved to {report_path} ---")


def main():
    if len(sys.argv) < 2:
        print(f"Usage: {sys.argv[0]} <snapshot_dir>")
        sys.exit(1)

    snap_dir = sys.argv[1]

    meta_path = os.path.join(snap_dir, "meta.json")
    meta = {}
    if os.path.exists(meta_path):
        with open(meta_path) as f:
            meta = json.load(f)

    results = load_results(snap_dir)
    if not results:
        print(f"No results found in {snap_dir}")
        sys.exit(1)

    generate_report(snap_dir, results, meta)


if __name__ == "__main__":
    main()
