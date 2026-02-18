# SQL Query Benchmarking Suite — Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Benchmark 16 hot-path sqlc queries across 5 groups with multi-tier data scales, producing benchstat-compatible output.

**Architecture:** A `pkg/db/bench/` test package with per-group benchmark files. `TestMain` creates isolated PostgreSQL databases per group, seeds them before any benchmark runs. Envelope benchmarks support 100K/1M/10M tiers controlled by `BENCH_TIER` env var; other groups use fixed scales. Results output in standard `go test -bench` format.

**Tech Stack:** Go `testing.B`, pgx v5, sqlc-generated `queries` package, `benchstat`

**Design Doc:** `docs/plans/2026-02-17-query-benchmarking-suite-design.md`

---

## Task 1: Create `pkg/db/bench/main_test.go`

**Files:**
- Create: `pkg/db/bench/main_test.go`

**What this does:** Foundation for all benchmarks. Creates isolated databases per group in `TestMain`, provides shared helpers. Since `TestMain` receives `*testing.M` (not `*testing.T`), we cannot use `testutils.NewRawDB`. Instead, write direct DB creation using pgx + migrations, with `log.Fatalf` for errors.

**Step 1: Create the file**

```go
package bench

import (
	"context"
	"crypto/rand"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/xmtp/xmtpd/pkg/migrations"
	"github.com/xmtp/xmtpd/pkg/testutils"
)

const (
	localDSNPrefix = "postgres://postgres:xmtp@localhost:8765"
	localDSNSuffix = "?sslmode=disable"
)

// envelopeTier holds a pre-seeded database and metadata for one data scale.
type envelopeTier struct {
	name        string   // sub-benchmark name, e.g. "rows=100K"
	db          *sql.DB
	count       int      // total envelope rows seeded
	topics      [][]byte // topics seeded into this DB
	originators []int32  // originator node IDs seeded
	payerIDs    []int32  // payer IDs created for batch insert benchmarks
}

// Package-level handles, populated by TestMain.
var (
	benchCtx      = context.Background()
	envelopeTiers []*envelopeTier
	congestionDB  *sql.DB
	ledgerDB      *sql.DB
	indexerDB     *sql.DB
	usageDB       *sql.DB

	// Seeded data references for non-envelope groups.
	congestionOriginators []int32
	congestionMaxMinute   int32
	ledgerPayerIDs        []int32
	indexerContracts      []string
	usagePayerIDs         []int32
	usageOriginators      []int32
	usageMaxMinute        int32
)

var cleanups []func()

func TestMain(m *testing.M) {
	ctlDB, err := connectDB(localDSNPrefix + localDSNSuffix)
	if err != nil {
		log.Fatalf("failed to connect to control DB: %v", err)
	}
	defer ctlDB.Close()

	// --- Envelope tiers ---
	tiers := []struct {
		name  string
		count int
	}{
		{"100K", 100_000},
		{"1M", 1_000_000},
		{"10M", 10_000_000},
	}

	benchTier := os.Getenv("BENCH_TIER")
	for _, t := range tiers {
		if benchTier != "" && !strings.EqualFold(benchTier, t.name) {
			continue
		}
		db, cleanup := createBenchDB(ctlDB, "env_"+strings.ToLower(t.name))
		cleanups = append(cleanups, cleanup)

		tier := &envelopeTier{
			name:  "rows=" + t.name,
			db:    db,
			count: t.count,
		}
		seedEnvelopes(benchCtx, tier)
		envelopeTiers = append(envelopeTiers, tier)
	}

	// --- Non-envelope groups ---
	congestionDB, cleanup := createBenchDB(ctlDB, "congestion")
	cleanups = append(cleanups, cleanup)
	seedCongestion(benchCtx, congestionDB)

	ledgerDB, cleanup = createBenchDB(ctlDB, "ledger")
	cleanups = append(cleanups, cleanup)
	seedLedger(benchCtx, ledgerDB)

	indexerDB, cleanup = createBenchDB(ctlDB, "indexer")
	cleanups = append(cleanups, cleanup)
	seedIndexer(benchCtx, indexerDB)

	usageDB, cleanup = createBenchDB(ctlDB, "usage")
	cleanups = append(cleanups, cleanup)
	seedUsage(benchCtx, usageDB)

	code := m.Run()

	for _, fn := range cleanups {
		fn()
	}
	os.Exit(code)
}

// connectDB opens a pgx-backed *sql.DB.
func connectDB(dsn string) (*sql.DB, error) {
	config, err := pgx.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}
	return stdlib.OpenDB(*config), nil
}

// createBenchDB creates an isolated database, runs migrations, returns handle + cleanup.
func createBenchDB(ctlDB *sql.DB, suffix string) (*sql.DB, func()) {
	dbName := "bench_" + suffix + "_" + testutils.RandomStringLower(8)
	log.Printf("creating benchmark database %s...", dbName)

	if _, err := ctlDB.Exec("CREATE DATABASE " + dbName); err != nil {
		log.Fatalf("create database %s: %v", dbName, err)
	}

	dsn := localDSNPrefix + "/" + dbName + localDSNSuffix
	db, err := connectDB(dsn)
	if err != nil {
		log.Fatalf("connect to %s: %v", dbName, err)
	}

	if err := migrations.Migrate(benchCtx, db); err != nil {
		log.Fatalf("migrate %s: %v", dbName, err)
	}

	return db, func() {
		db.Close()
		_, _ = ctlDB.Exec("DROP DATABASE " + dbName)
	}
}

// randomBytes returns n random bytes.
func randomBytes(n int) []byte {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return b
}
```

**Step 2: Verify it compiles**

Run: `go build ./pkg/db/bench/`

This will fail because the `seed*` functions are not yet defined. That is expected — they will be added in subsequent tasks. For now, verify there are no syntax errors by checking with `go vet` after all files exist.

**Step 3: Commit**

Do NOT commit yet — wait until Task 2 so TestMain can compile.

---

## Task 2: Create `pkg/db/bench/indexer_bench_test.go`

**Files:**
- Create: `pkg/db/bench/indexer_bench_test.go`

**Why indexer first:** Simplest group (2 queries, ~100 rows, no partitions). Validates that the entire infrastructure works end-to-end.

**Step 1: Create the file**

```go
package bench

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/db/queries"
)

const numContracts = 100

// seedIndexer populates the latest_block table with contract addresses.
func seedIndexer(ctx context.Context, db *sql.DB) {
	q := queries.New(db)
	indexerContracts = make([]string, numContracts)
	for i := range numContracts {
		addr := fmt.Sprintf("0x%040x", i)
		indexerContracts[i] = addr
		err := q.SetLatestBlock(ctx, queries.SetLatestBlockParams{
			ContractAddress: addr,
			BlockNumber:     int64(1000 + i),
			BlockHash:       randomBytes(32),
		})
		if err != nil {
			log.Fatalf("seed indexer: %v", err)
		}
	}
	log.Printf("seeded indexer: %d contracts", numContracts)
}

func BenchmarkGetLatestBlock(b *testing.B) {
	q := queries.New(indexerDB)
	contract := indexerContracts[0]
	b.ResetTimer()
	for b.Loop() {
		_, err := q.GetLatestBlock(benchCtx, contract)
		require.NoError(b, err)
	}
}

func BenchmarkSetLatestBlock(b *testing.B) {
	q := queries.New(indexerDB)
	var counter atomic.Int64
	counter.Store(100_000)
	b.ResetTimer()
	for b.Loop() {
		blockNum := counter.Add(1)
		err := q.SetLatestBlock(benchCtx, queries.SetLatestBlockParams{
			ContractAddress: indexerContracts[0],
			BlockNumber:     blockNum,
			BlockHash:       randomBytes(32),
		})
		require.NoError(b, err)
	}
}
```

**Step 2: Create stub files so the package compiles**

Create minimal stub files for the other four groups so `TestMain` can call their seed functions. Each stub defines only the `seed*` function with a placeholder body. These will be replaced in subsequent tasks.

Create `pkg/db/bench/envelopes_bench_test.go`:
```go
package bench

import (
	"context"
	"log"
)

func seedEnvelopes(_ context.Context, _ *envelopeTier) {
	log.Printf("seeded envelopes: %d rows (stub)", 0)
}
```

Create `pkg/db/bench/congestion_bench_test.go`:
```go
package bench

import (
	"context"
	"database/sql"
	"log"
)

func seedCongestion(_ context.Context, _ *sql.DB) {
	log.Printf("seeded congestion: 0 rows (stub)")
}
```

Create `pkg/db/bench/ledger_bench_test.go`:
```go
package bench

import (
	"context"
	"database/sql"
	"log"
)

func seedLedger(_ context.Context, _ *sql.DB) {
	log.Printf("seeded ledger: 0 rows (stub)")
}
```

Create `pkg/db/bench/usage_bench_test.go`:
```go
package bench

import (
	"context"
	"database/sql"
	"log"
)

func seedUsage(_ context.Context, _ *sql.DB) {
	log.Printf("seeded usage: 0 rows (stub)")
}
```

**Step 3: Verify it compiles and runs**

Run: `go test -bench=BenchmarkGetLatestBlock -benchmem -count=1 -timeout=2m ./pkg/db/bench/`

Expected: Benchmark runs successfully against seeded data. Output like:
```
BenchmarkGetLatestBlock-8    NNNN    NNNN ns/op    NNN B/op    N allocs/op
```

**Step 4: Commit**

```
git add pkg/db/bench/
git commit -m "feat: add benchmark infrastructure and indexer benchmarks"
```

---

## Task 3: Create `pkg/db/bench/congestion_bench_test.go`

**Files:**
- Replace: `pkg/db/bench/congestion_bench_test.go` (replace stub)

**Step 1: Write the full file**

```go
package bench

import (
	"context"
	"database/sql"
	"log"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/db/queries"
)

const (
	numCongestionOriginators = 5
	numCongestionMinutes     = 2000 // per originator
)

// seedCongestion populates originator_congestion with time-bucketed message counts.
func seedCongestion(ctx context.Context, db *sql.DB) {
	q := queries.New(db)
	congestionOriginators = make([]int32, numCongestionOriginators)
	for i := range numCongestionOriginators {
		origID := int32(500 + i)
		congestionOriginators[i] = origID
		for minute := int32(0); minute < numCongestionMinutes; minute++ {
			err := q.IncrementOriginatorCongestion(ctx, queries.IncrementOriginatorCongestionParams{
				OriginatorID:      origID,
				MinutesSinceEpoch: minute,
			})
			if err != nil {
				log.Fatalf("seed congestion: %v", err)
			}
		}
	}
	congestionMaxMinute = numCongestionMinutes - 1
	log.Printf("seeded congestion: %d rows", numCongestionOriginators*numCongestionMinutes)
}

func BenchmarkIncrementOriginatorCongestion(b *testing.B) {
	q := queries.New(congestionDB)
	origID := congestionOriginators[0]
	var counter atomic.Int32
	counter.Store(100_000) // start beyond seeded range
	b.ResetTimer()
	for b.Loop() {
		minute := counter.Add(1)
		err := q.IncrementOriginatorCongestion(benchCtx, queries.IncrementOriginatorCongestionParams{
			OriginatorID:      origID,
			MinutesSinceEpoch: minute,
		})
		require.NoError(b, err)
	}
}

func BenchmarkGetRecentOriginatorCongestion(b *testing.B) {
	q := queries.New(congestionDB)
	params := queries.GetRecentOriginatorCongestionParams{
		OriginatorID: congestionOriginators[0],
		EndMinute:    congestionMaxMinute,
		NumMinutes:   60, // last hour
	}
	b.ResetTimer()
	for b.Loop() {
		_, err := q.GetRecentOriginatorCongestion(benchCtx, params)
		require.NoError(b, err)
	}
}

func BenchmarkSumOriginatorCongestion(b *testing.B) {
	q := queries.New(congestionDB)
	params := queries.SumOriginatorCongestionParams{
		OriginatorID:        congestionOriginators[0],
		MinutesSinceEpochGt: 0,
		MinutesSinceEpochLt: int64(congestionMaxMinute),
	}
	b.ResetTimer()
	for b.Loop() {
		_, err := q.SumOriginatorCongestion(benchCtx, params)
		require.NoError(b, err)
	}
}
```

**Step 2: Verify benchmarks run**

Run: `go test -bench=BenchmarkGetRecentOriginatorCongestion -benchmem -count=1 -timeout=2m ./pkg/db/bench/`

Expected: PASS with benchmark output.

**Step 3: Commit**

```
git add pkg/db/bench/congestion_bench_test.go
git commit -m "feat: add congestion query benchmarks"
```

---

## Task 4: Create `pkg/db/bench/ledger_bench_test.go`

**Files:**
- Replace: `pkg/db/bench/ledger_bench_test.go` (replace stub)

**Step 1: Write the full file**

```go
package bench

import (
	"context"
	"database/sql"
	"log"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/db/queries"
)

const (
	numLedgerPayers        = 50
	numLedgerEventsPerPayer = 100
)

// seedLedger creates payers and populates payer_ledger_events.
func seedLedger(ctx context.Context, db *sql.DB) {
	q := queries.New(db)
	ledgerPayerIDs = make([]int32, numLedgerPayers)

	for i := range numLedgerPayers {
		addr := string(randomBytes(20))
		id, err := q.FindOrCreatePayer(ctx, addr)
		if err != nil {
			log.Fatalf("seed ledger payer: %v", err)
		}
		ledgerPayerIDs[i] = id

		// Insert events: mix of deposits (1), withdrawals (2), settlements (3)
		for j := range numLedgerEventsPerPayer {
			eventType := int16((j % 3) + 1)
			amount := int64(1_000_000) // 1M picodollars
			if eventType == 2 {
				amount = -amount
			}
			err := q.InsertPayerLedgerEvent(ctx, queries.InsertPayerLedgerEventParams{
				EventID:           randomBytes(32),
				PayerID:           id,
				AmountPicodollars: amount,
				EventType:         eventType,
			})
			if err != nil {
				log.Fatalf("seed ledger event: %v", err)
			}
		}
	}
	log.Printf("seeded ledger: %d payers, %d events", numLedgerPayers, numLedgerPayers*numLedgerEventsPerPayer)
}

func BenchmarkInsertPayerLedgerEvent(b *testing.B) {
	q := queries.New(ledgerDB)
	payerID := ledgerPayerIDs[0]
	b.ResetTimer()
	for b.Loop() {
		err := q.InsertPayerLedgerEvent(benchCtx, queries.InsertPayerLedgerEventParams{
			EventID:           randomBytes(32), // unique per iteration
			PayerID:           payerID,
			AmountPicodollars: 1_000_000,
			EventType:         1,
		})
		require.NoError(b, err)
	}
}

func BenchmarkGetPayerBalance(b *testing.B) {
	q := queries.New(ledgerDB)
	payerID := ledgerPayerIDs[0]
	b.ResetTimer()
	for b.Loop() {
		_, err := q.GetPayerBalance(benchCtx, payerID)
		require.NoError(b, err)
	}
}

func BenchmarkGetLastEvent(b *testing.B) {
	q := queries.New(ledgerDB)
	params := queries.GetLastEventParams{
		PayerID:   ledgerPayerIDs[0],
		EventType: 1, // deposits
	}
	b.ResetTimer()
	for b.Loop() {
		_, err := q.GetLastEvent(benchCtx, params)
		require.NoError(b, err)
	}
}
```

**Step 2: Verify benchmarks run**

Run: `go test -bench=BenchmarkGetPayerBalance -benchmem -count=1 -timeout=2m ./pkg/db/bench/`

**Step 3: Commit**

```
git add pkg/db/bench/ledger_bench_test.go
git commit -m "feat: add ledger query benchmarks"
```

---

## Task 5: Create `pkg/db/bench/usage_bench_test.go`

**Files:**
- Replace: `pkg/db/bench/usage_bench_test.go` (replace stub)

**Step 1: Write the full file**

```go
package bench

import (
	"context"
	"database/sql"
	"log"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/db/queries"
)

const (
	numUsagePayers      = 50
	numUsageOriginators = 5
	numUsageMinutes     = 40
)

// seedUsage creates payers and populates unsettled_usage.
func seedUsage(ctx context.Context, db *sql.DB) {
	q := queries.New(db)
	usagePayerIDs = make([]int32, numUsagePayers)
	usageOriginators = make([]int32, numUsageOriginators)

	for i := range numUsageOriginators {
		usageOriginators[i] = int32(600 + i)
	}

	for i := range numUsagePayers {
		addr := string(randomBytes(20))
		id, err := q.FindOrCreatePayer(ctx, addr)
		if err != nil {
			log.Fatalf("seed usage payer: %v", err)
		}
		usagePayerIDs[i] = id

		for _, origID := range usageOriginators {
			for minute := int32(0); minute < numUsageMinutes; minute++ {
				err := q.IncrementUnsettledUsage(ctx, queries.IncrementUnsettledUsageParams{
					PayerID:           id,
					OriginatorID:      origID,
					MinutesSinceEpoch: minute,
					SpendPicodollars:  1_000_000,
					SequenceID:        int64(minute),
					MessageCount:      1,
				})
				if err != nil {
					log.Fatalf("seed usage: %v", err)
				}
			}
		}
	}
	usageMaxMinute = numUsageMinutes - 1
	log.Printf("seeded usage: %d rows", numUsagePayers*numUsageOriginators*numUsageMinutes)
}

func BenchmarkIncrementUnsettledUsage(b *testing.B) {
	q := queries.New(usageDB)
	payerID := usagePayerIDs[0]
	origID := usageOriginators[0]
	var counter atomic.Int32
	counter.Store(100_000) // beyond seeded range
	b.ResetTimer()
	for b.Loop() {
		minute := counter.Add(1)
		err := q.IncrementUnsettledUsage(benchCtx, queries.IncrementUnsettledUsageParams{
			PayerID:           payerID,
			OriginatorID:      origID,
			MinutesSinceEpoch: minute,
			SpendPicodollars:  1_000_000,
			SequenceID:        int64(minute),
			MessageCount:      1,
		})
		require.NoError(b, err)
	}
}
```

**Step 2: Verify benchmarks run**

Run: `go test -bench=BenchmarkIncrementUnsettledUsage -benchmem -count=1 -timeout=2m ./pkg/db/bench/`

**Step 3: Commit**

```
git add pkg/db/bench/usage_bench_test.go
git commit -m "feat: add unsettled usage benchmark"
```

---

## Task 6: Create `pkg/db/bench/envelopes_bench_test.go`

**Files:**
- Replace: `pkg/db/bench/envelopes_bench_test.go` (replace stub)

**This is the most complex task.** Needs partition creation, bulk seeding, and 7 benchmarks with tier sub-benchmarks.

**Step 1: Write the full file**

Key seeding approach:
- 3 originators (100, 200, 300), 100 topics, ~500-byte blobs
- Pre-create partitions via `EnsureGatewayParts` for all (originator, seq_id range) combos
- Use `InsertGatewayEnvelopeWithChecksStandalone` from `pkg/db` for seeding (handles partition creation automatically, simpler than raw SQL)
- For 10M tier, seed in batches with progress logging
- Also create a write-benchmark originator (999) with pre-created partitions
- Create 5 payers for batch insert benchmarks

```go
package bench

import (
	"context"
	"database/sql"
	"log"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/db"
	"github.com/xmtp/xmtpd/pkg/db/queries"
)

const (
	numOriginators     = 3
	numTopics          = 100
	blobSize           = 500
	writeOriginatorID  = int32(999) // dedicated originator for write benchmarks
	numBenchPayers     = 5
)

var envelopeOriginators = []int32{100, 200, 300}

// seedEnvelopes populates gateway_envelopes_meta and gateway_envelope_blobs.
func seedEnvelopes(ctx context.Context, tier *envelopeTier) {
	q := queries.New(tier.db)
	tier.originators = envelopeOriginators

	// Generate topics
	tier.topics = make([][]byte, numTopics)
	for i := range numTopics {
		tier.topics[i] = randomBytes(32)
	}

	// Create payers for batch insert benchmarks
	tier.payerIDs = make([]int32, numBenchPayers)
	for i := range numBenchPayers {
		id, err := q.FindOrCreatePayer(ctx, string(randomBytes(20)))
		if err != nil {
			log.Fatalf("seed envelope payer: %v", err)
		}
		tier.payerIDs[i] = id
	}

	// Pre-create partitions for write benchmark originator
	perOriginator := tier.count / numOriginators
	for seqID := int64(0); seqID < int64(perOriginator)+db.GatewayEnvelopeBandWidth; seqID += db.GatewayEnvelopeBandWidth {
		for _, origID := range tier.originators {
			_ = q.EnsureGatewayParts(ctx, queries.EnsureGatewayPartsParams{
				OriginatorNodeID:     origID,
				OriginatorSequenceID: seqID,
				BandWidth:            db.GatewayEnvelopeBandWidth,
			})
		}
	}
	// Partitions for write benchmark originator
	for seqID := int64(0); seqID < 10*db.GatewayEnvelopeBandWidth; seqID += db.GatewayEnvelopeBandWidth {
		_ = q.EnsureGatewayParts(ctx, queries.EnsureGatewayPartsParams{
			OriginatorNodeID:     writeOriginatorID,
			OriginatorSequenceID: seqID,
			BandWidth:            db.GatewayEnvelopeBandWidth,
		})
	}

	// Seed envelopes distributed across originators and topics
	batchSize := 10_000
	blob := randomBytes(blobSize) // reuse same blob for speed
	seqIDs := make([]int64, numOriginators) // per-originator sequence counters

	for i := range tier.count {
		origIdx := i % numOriginators
		origID := tier.originators[origIdx]
		seqIDs[origIdx]++
		topicIdx := i % numTopics

		_, err := db.InsertGatewayEnvelopeWithChecksStandalone(ctx, q, queries.InsertGatewayEnvelopeParams{
			OriginatorNodeID:     origID,
			OriginatorSequenceID: seqIDs[origIdx],
			Topic:                tier.topics[topicIdx],
			Expiry:               time.Now().Add(24 * time.Hour).Unix(),
			OriginatorEnvelope:   blob,
		})
		if err != nil {
			log.Fatalf("seed envelope %d: %v", i, err)
		}

		if (i+1)%batchSize == 0 {
			log.Printf("seeded %d/%d envelopes for tier %s", i+1, tier.count, tier.name)
		}
	}
	log.Printf("seeded envelopes: %d rows for tier %s", tier.count, tier.name)
}

// --- Read benchmarks ---

func BenchmarkSelectGatewayEnvelopesByTopics(b *testing.B) {
	for _, tier := range envelopeTiers {
		b.Run(tier.name, func(b *testing.B) {
			q := queries.New(tier.db)
			// Query 10 topics, all originators, cursor at 50% through data
			midSeq := int64(tier.count / numOriginators / 2)
			params := queries.SelectGatewayEnvelopesByTopicsParams{
				Topics:            tier.topics[:10],
				RowLimit:          100,
				CursorNodeIds:     tier.originators,
				CursorSequenceIds: []int64{midSeq, midSeq, midSeq},
			}
			b.ResetTimer()
			for b.Loop() {
				_, err := q.SelectGatewayEnvelopesByTopics(benchCtx, params)
				require.NoError(b, err)
			}
		})
	}
}

func BenchmarkSelectGatewayEnvelopesByOriginators(b *testing.B) {
	for _, tier := range envelopeTiers {
		b.Run(tier.name, func(b *testing.B) {
			q := queries.New(tier.db)
			midSeq := int64(tier.count / numOriginators / 2)
			params := queries.SelectGatewayEnvelopesByOriginatorsParams{
				OriginatorNodeIds: tier.originators,
				RowsPerOriginator: 50,
				RowLimit:          100,
				CursorNodeIds:     tier.originators,
				CursorSequenceIds: []int64{midSeq, midSeq, midSeq},
			}
			b.ResetTimer()
			for b.Loop() {
				_, err := q.SelectGatewayEnvelopesByOriginators(benchCtx, params)
				require.NoError(b, err)
			}
		})
	}
}

func BenchmarkSelectGatewayEnvelopesBySingleOriginator(b *testing.B) {
	for _, tier := range envelopeTiers {
		b.Run(tier.name, func(b *testing.B) {
			q := queries.New(tier.db)
			midSeq := int64(tier.count / numOriginators / 2)
			params := queries.SelectGatewayEnvelopesBySingleOriginatorParams{
				OriginatorNodeID: tier.originators[0],
				CursorSequenceID: midSeq,
				RowLimit:         100,
			}
			b.ResetTimer()
			for b.Loop() {
				_, err := q.SelectGatewayEnvelopesBySingleOriginator(benchCtx, params)
				require.NoError(b, err)
			}
		})
	}
}

func BenchmarkSelectGatewayEnvelopesUnfiltered(b *testing.B) {
	for _, tier := range envelopeTiers {
		b.Run(tier.name, func(b *testing.B) {
			q := queries.New(tier.db)
			midSeq := int64(tier.count / numOriginators / 2)
			params := queries.SelectGatewayEnvelopesUnfilteredParams{
				RowLimit:          100,
				CursorNodeIds:     tier.originators,
				CursorSequenceIds: []int64{midSeq, midSeq, midSeq},
			}
			b.ResetTimer()
			for b.Loop() {
				_, err := q.SelectGatewayEnvelopesUnfiltered(benchCtx, params)
				require.NoError(b, err)
			}
		})
	}
}

func BenchmarkSelectNewestFromTopics(b *testing.B) {
	for _, tier := range envelopeTiers {
		b.Run(tier.name, func(b *testing.B) {
			q := queries.New(tier.db)
			topics := tier.topics[:10]
			b.ResetTimer()
			for b.Loop() {
				_, err := q.SelectNewestFromTopics(benchCtx, topics)
				require.NoError(b, err)
			}
		})
	}
}

// --- Write benchmarks ---

func BenchmarkInsertGatewayEnvelope(b *testing.B) {
	for _, tier := range envelopeTiers {
		b.Run(tier.name, func(b *testing.B) {
			q := queries.New(tier.db)
			blob := randomBytes(blobSize)
			topic := tier.topics[0]
			var counter atomic.Int64
			counter.Store(1_000_000) // beyond seeded range for writeOriginatorID
			b.ResetTimer()
			for b.Loop() {
				seqID := counter.Add(1)
				_, err := q.InsertGatewayEnvelope(benchCtx, queries.InsertGatewayEnvelopeParams{
					OriginatorNodeID:     writeOriginatorID,
					OriginatorSequenceID: seqID,
					Topic:                topic,
					Expiry:               time.Now().Add(24 * time.Hour).Unix(),
					OriginatorEnvelope:   blob,
				})
				require.NoError(b, err)
			}
		})
	}
}

func BenchmarkInsertGatewayEnvelopeBatch(b *testing.B) {
	for _, tier := range envelopeTiers {
		b.Run(tier.name, func(b *testing.B) {
			q := queries.New(tier.db)
			blob := randomBytes(blobSize)
			batchLen := 10
			var counter atomic.Int64
			counter.Store(5_000_000)
			b.ResetTimer()
			for b.Loop() {
				baseSeq := counter.Add(int64(batchLen))
				nodeIDs := make([]int32, batchLen)
				seqIDs := make([]int64, batchLen)
				topics := make([][]byte, batchLen)
				payerIDs := make([]int32, batchLen)
				times := make([]time.Time, batchLen)
				expiries := make([]int64, batchLen)
				blobs := make([][]byte, batchLen)
				spends := make([]int64, batchLen)
				now := time.Now()
				exp := now.Add(24 * time.Hour).Unix()

				for j := range batchLen {
					nodeIDs[j] = writeOriginatorID
					seqIDs[j] = baseSeq + int64(j)
					topics[j] = tier.topics[j%numTopics]
					payerIDs[j] = tier.payerIDs[j%numBenchPayers]
					times[j] = now
					expiries[j] = exp
					blobs[j] = blob
					spends[j] = 1_000_000
				}

				_, err := q.InsertGatewayEnvelopeBatchAndIncrementUnsettledUsage(
					benchCtx,
					queries.InsertGatewayEnvelopeBatchAndIncrementUnsettledUsageParams{
						OriginatorNodeIds:     nodeIDs,
						OriginatorSequenceIds: seqIDs,
						Topics:                topics,
						PayerIds:              payerIDs,
						GatewayTimes:          times,
						Expiries:              expiries,
						OriginatorEnvelopes:   blobs,
						SpendPicodollars:      spends,
					},
				)
				require.NoError(b, err)
			}
		})
	}
}
```

**Step 2: Verify envelopes benchmarks run (100K tier only for speed)**

Run: `BENCH_TIER=100K go test -bench=BenchmarkSelectGatewayEnvelopesByTopics -benchmem -count=1 -timeout=5m ./pkg/db/bench/`

Expected: Sub-benchmark `rows=100K` runs and produces output.

**Step 3: Commit**

```
git add pkg/db/bench/envelopes_bench_test.go
git commit -m "feat: add envelope query benchmarks with multi-tier seeding"
```

---

## Task 7: Create `dev/bench` script and results directory

**Files:**
- Create: `dev/bench`
- Create: `benchmarks/results/.gitkeep`

**Step 1: Create the runner script**

```bash
#!/bin/bash
set -euo pipefail

RESULTS_DIR="benchmarks/results"
mkdir -p "$RESULTS_DIR"

OUTFILE="$RESULTS_DIR/$(date +%Y-%m-%d).txt"

echo "Running benchmarks..."
echo "Tier filter: ${BENCH_TIER:-all}"
echo "Output: $OUTFILE"
echo

go test -bench=. -benchmem -count=5 -timeout=60m \
    -run='^$' \
    ./pkg/db/bench/ 2>&1 | tee "$OUTFILE"

echo
echo "Results written to $OUTFILE"
echo "Compare with: benchstat benchmarks/results/old.txt $OUTFILE"
```

Note: `-run='^$'` ensures no regular tests run, only benchmarks.

**Step 2: Make it executable**

```
chmod +x dev/bench
```

**Step 3: Create results directory**

```
mkdir -p benchmarks/results
touch benchmarks/results/.gitkeep
```

**Step 4: Commit**

```
git add dev/bench benchmarks/results/.gitkeep
git commit -m "feat: add dev/bench runner script and results directory"
```

---

## Task 8: Add benchstat to `tools/go.mod`

**Files:**
- Modify: `tools/go.mod`

**Why:** `benchstat` is the standard Go tool for comparing benchmark results across runs. It should be tracked as a project tool dependency like `sqlc`, `mockery`, etc.

**Step 1: Add benchstat to the tool block**

Add `golang.org/x/perf/cmd/benchstat` to the `tool (...)` block in `tools/go.mod`:

```
tool (
	github.com/bufbuild/buf/cmd/buf
	github.com/ethereum/go-ethereum/cmd/abigen
	github.com/golang-migrate/migrate/v4/cmd/migrate
	github.com/sqlc-dev/sqlc/cmd/sqlc
	github.com/vektra/mockery/v2
	golang.org/x/perf/cmd/benchstat
)
```

**Step 2: Run `go mod tidy` in the tools directory**

```
cd tools && go mod tidy
```

This resolves the `golang.org/x/perf` dependency and updates `tools/go.sum`.

**Step 3: Verify benchstat works**

```
go tool benchstat -h
```

Expected: Help output from benchstat.

**Step 4: Commit**

```
git add tools/go.mod tools/go.sum
git commit -m "feat: add benchstat to tools/go.mod"
```

---

## Task 9: Full verification run

**Step 1: Run all non-envelope benchmarks**

Run: `BENCH_TIER=none go test -bench=. -benchmem -count=1 -timeout=5m -run='^$' ./pkg/db/bench/`

Note: `BENCH_TIER=none` skips all envelope tiers (no tier matches "none"). This runs only indexer, congestion, ledger, and usage benchmarks.

Verify: All 10 non-envelope benchmarks produce output. No errors.

**Step 2: Run envelope benchmarks at 100K**

Run: `BENCH_TIER=100K go test -bench=Benchmark -benchmem -count=1 -timeout=5m -run='^$' ./pkg/db/bench/`

Verify: All 7 envelope benchmarks produce `rows=100K` sub-benchmark output plus the 10 non-envelope benchmarks. No errors.

**Step 3: Run the dev/bench script at 100K**

Run: `BENCH_TIER=100K dev/bench`

Verify: Results file created at `benchmarks/results/YYYY-MM-DD.txt`. File contains valid benchstat-format output.

**Step 4: Final commit with any fixes**

If any issues were found and fixed, commit them.

---

## Notes

- **Seeding time for 10M tier:** Using `InsertGatewayEnvelopeWithChecksStandalone` one-at-a-time will be slow (~10-30 minutes for 10M rows). If this is unacceptable, a follow-up task can optimize seeding with raw batch SQL or `pgx.CopyFrom`. The 100K and 1M tiers are practical for development iteration.
- **Write benchmark isolation:** Write benchmarks use originator 999 (separate from seeded originators 100/200/300) so they don't pollute read benchmark data within the same `go test` process. However, across `b.N` iterations, writes accumulate. This is intentional — it benchmarks insert performance into a growing table.
- **`b.Loop()` vs `for i := 0; i < b.N; i++`:** The plan uses `b.Loop()` which requires Go 1.24+. The project uses Go 1.25 so this is safe. `b.Loop()` prevents compiler optimizations from eliding the benchmark body.
