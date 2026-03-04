# Database benchmarks

- [Database benchmarks](#database-benchmarks)
  - [What these benchmarks measure](#what-these-benchmarks-measure)
  - [How results are generated](#how-results-are-generated)
  - [Field-by-field: what each line means](#field-by-field-what-each-line-means)
  - [How to interpret in practice](#how-to-interpret-in-practice)
  - [How to spot problems](#how-to-spot-problems)
  - [Comparing runs correctly](#comparing-runs-correctly)
  - [Common interpretation mistakes](#common-interpretation-mistakes)
  - [Recommended workflow for query tuning](#recommended-workflow-for-query-tuning)

This directory stores benchmark outputs produced by `dev/bench`.

## What these benchmarks measure

The DB benchmarks in `pkg/db/bench` are Go microbenchmarks that execute real SQL against a local Postgres database seeded with test data.

They measure **end-to-end query call cost from Go**, including:

- SQL execution in Postgres
- network hop to local DB
- row scan/decoding in Go (`database/sql` + driver)
- Go-side allocations

They do **not** isolate query planning internals by themselves. Use `EXPLAIN (ANALYZE, BUFFERS)` and `pg_stat_statements` for planner-level diagnosis.

## How results are generated

`dev/bench` runs:

```sh
go test -tags bench -bench=. -benchmem -count=5 -timeout=60m -run='^$' ./pkg/db/bench/
```

Key flags:

- `-bench=.`: run all benchmarks
- `-benchmem`: include memory stats (`B/op`, `allocs/op`)
- `-count=5`: run each benchmark 5 times (helps detect variance)
- `-run='^$'`: skip normal tests, run only benchmarks

Additionally `-benchtime=Xs` can be passed to increase the default bench time (1s).

Output files are written to `benchmarks/results/YYYY-MM-DD.txt`.

## Field-by-field: what each line means

Example line:

```text
BenchmarkSelectGatewayEnvelopesBySingleOriginator/rows=100K-16  3104  402575 ns/op  153105 B/op  1040 allocs/op
```

- `BenchmarkSelectGatewayEnvelopesBySingleOriginator`
  - Benchmark function name.

- `/rows=100K`
  - Sub-benchmark case name (here, seeded data tier).

- `-16`
  - GOMAXPROCS used during run (number of logical threads available to Go runtime).

- `3104`
  - Number of benchmark iterations executed to estimate stable timing.
  - Note that benchmarks in Go run as much times as possible in the benchtime set. Higher is usually better because it means each operation is cheap enough to run many times in the same benchmark window. Compare this field only within the same benchmark configuration.

- `402575 ns/op`
  - Average wall-clock time per operation (one benchmark loop iteration), in nanoseconds.
  - Lower is better (faster). If this increases, the operation got slower.

- `153105 B/op`
  - Average bytes allocated on the Go heap per operation.
  - Lower is generally better (less memory churn, less GC pressure). Higher indicates more allocation work in Go code/scan paths.

- `1040 allocs/op`
  - Average number of heap allocation events per operation.
  - Lower is better. Increases usually indicate extra object creation, less efficient decoding, or changed data-shaping behavior.

Header fields:

- `goos`, `goarch`, `cpu`: machine context for run-to-run comparability.
- `pkg`: benchmark package under test.

Footer fields:

- `PASS`: benchmark binary completed successfully.
- `ok ... 203.147s`: total package runtime, not per-query latency.

## How to interpret in practice

For DB performance decisions, read these in order:

1. `ns/op`: first-order latency signal
2. `allocs/op`: GC pressure and decoding/object churn
3. `B/op`: memory footprint trend (supports allocs diagnosis)
4. variance across repeated lines: stability/reliability of signal

## How to spot problems

Use this checklist:

- **Large variance between repeats**
  - If the same benchmark case differs by ~20%+ between runs, investigate noise before concluding regressions.
  - Common causes: cold caches, background load, WAL/checkpoint spikes, partition creation side effects.
  - Also note that each test is run more than once, so later tests should be run against hot caches.

- **Tier scaling looks wrong**
  - Compare `rows=100K`, `rows=1M`, `rows=10M`.
  - A steep latency jump can indicate index/selectivity issues or plan change.

- **`allocs/op` regresses while `ns/op` is flat**
  - Likely Go decoding/object churn increase; may become latency later under load due to GC.

- **`ns/op` regresses but allocs are flat**
  - More likely DB-side issue (plan, index usage, I/O, contention), not just Go heap behavior.

- **Unexpected query ordering**
  - If a narrower query becomes slower than a broader query unexpectedly, inspect execution plans.

- **Single massive outliers**
  - Example pattern in sample results: one or two much slower lines followed by steady faster lines.
  - Usually warm-up/caching effects; prefer statistical comparison tools over eyeballing.

## Comparing runs correctly

Prefer `benchstat` instead of manual visual diff:

```sh
benchstat benchmarks/results/old.txt benchmarks/results/new.txt
```

Install with:

```sh
go install golang.org/x/perf/cmd/benchstat@latest
```

What to look for:

- Mean/median delta (`+/- %`)
- Statistical significance (`p` values in benchstat output)
- Consistent direction across related benchmarks, not isolated one-off changes

Note that `benchstat` is only precise when comparing benchmarking results generated with the same hardware and settings.

## Common interpretation mistakes

- Treating `B/op` and `allocs/op` as DB engine memory metrics, they're not. They are Go process metrics.
- Reading one run only (`-count=1`) and assuming it is stable.
- Comparing results from different machines/CPU governors/DB settings.
- Ignoring seed shape: synthetic uniform data can hide real-world skew/selectivity issues.

## Recommended workflow for query tuning

1. Reproduce with benchmark (`dev/bench` or focused benchmark).
2. Confirm regression/improvement with `benchstat`.
3. Inspect SQL plan with `EXPLAIN (ANALYZE, BUFFERS)` on same query shape.
4. Validate DB-level behavior via `pg_stat_statements`.
5. Re-run benchmark to confirm end-to-end impact.

This gives both app-level and planner-level truth.
