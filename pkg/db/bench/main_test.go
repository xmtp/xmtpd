//go:build bench

package bench

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/xmtp/xmtpd/pkg/db/migrations"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/testutils"
)

const (
	localDSNPrefix = "postgres://postgres:xmtp@localhost:8765"
	localDSNSuffix = "?sslmode=disable"
)

// envelopeTier holds a pre-seeded database and metadata for one data scale.
type envelopeTier struct {
	name        string // sub-benchmark name, e.g. "rows=100K"
	db          *sql.DB
	queries     *queries.Queries
	count       int      // total envelope rows seeded
	topics      [][]byte // topics seeded into this DB
	originators []int32  // originator node IDs seeded
	payerIDs    []int32  // payer IDs created for batch insert benchmarks
}

// Package-level handles, populated by TestMain.
// testing.M has no Context() or Cleanup() methods, so we manage lifecycle manually.
var (
	benchCtx          = context.Background()
	usePreparedPlans  = parseBenchUsePrepared()
	envelopeTiers     []*envelopeTier
	congestionDB      *sql.DB
	congestionQueries *queries.Queries
	ledgerDB          *sql.DB
	ledgerQueries     *queries.Queries
	indexerDB         *sql.DB
	indexerQueries    *queries.Queries
	usageDB           *sql.DB
	usageQueries      *queries.Queries
	hotPathDB         *sql.DB
	hotPathQueries    *queries.Queries

	// Seeded data references for non-envelope groups.
	congestionOriginators []int32
	congestionMaxMinute   int32
	ledgerPayerIDs        []int32
	indexerContracts      []string
	usagePayerIDs         []int32
	usageOriginators      []int32
	usageMaxMinute        int32
	hotPathPayerIDs       []int32
)

func TestMain(m *testing.M) {
	log.Printf("benchmark prepared queries enabled: %t", usePreparedPlans)

	ctlDB, err := connectDB(localDSNPrefix + localDSNSuffix)
	if err != nil {
		log.Fatalf("failed to connect to control DB: %v", err)
	}

	// Track cleanups so they run even if a seed function calls log.Fatalf
	// (log.Fatalf calls os.Exit, so defers in this function won't run either —
	// but at least we get cleanup for the normal exit path).
	var cleanups []func()

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
		db, querySet, cleanup := createBenchDB(ctlDB, "env_"+strings.ToLower(t.name))
		cleanups = append(cleanups, cleanup)

		tier := &envelopeTier{
			name:    "rows=" + t.name,
			db:      db,
			queries: querySet,
			count:   t.count,
		}
		seedEnvelopes(benchCtx, tier)
		envelopeTiers = append(envelopeTiers, tier)
	}

	// --- Non-envelope groups ---
	var cleanup func()
	congestionDB, congestionQueries, cleanup = createBenchDB(ctlDB, "congestion")
	cleanups = append(cleanups, cleanup)
	seedCongestion(benchCtx)

	ledgerDB, ledgerQueries, cleanup = createBenchDB(ctlDB, "ledger")
	cleanups = append(cleanups, cleanup)
	seedLedger(benchCtx)

	indexerDB, indexerQueries, cleanup = createBenchDB(ctlDB, "indexer")
	cleanups = append(cleanups, cleanup)
	seedIndexer(benchCtx)

	usageDB, usageQueries, cleanup = createBenchDB(ctlDB, "usage")
	cleanups = append(cleanups, cleanup)
	seedUsage(benchCtx)

	hotPathDB, hotPathQueries, cleanup = createBenchDB(ctlDB, "hot_path")
	cleanups = append(cleanups, cleanup)
	seedHotPath(benchCtx)

	code := m.Run()
	for _, fn := range cleanups {
		fn()
	}
	_ = ctlDB.Close()
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
func createBenchDB(
	ctlDB *sql.DB,
	suffix string,
) (db *sql.DB, q *queries.Queries, cleanup func()) {
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

	if !usePreparedPlans {
		querySet := queries.New(db)
		return db, querySet, func() {
			_ = db.Close()
			_, _ = ctlDB.Exec("DROP DATABASE " + dbName + " WITH (FORCE)")
		}
	}

	querySet, err := queries.Prepare(benchCtx, db)
	if err != nil {
		log.Fatalf("prepare queries %s: %v", dbName, err)
	}

	return db, querySet, func() {
		_ = querySet.Close()
		_ = db.Close()
		_, _ = ctlDB.Exec("DROP DATABASE " + dbName + " WITH (FORCE)")
	}
}

func parseBenchUsePrepared() bool {
	raw := os.Getenv("BENCH_USE_PREPARED")
	if raw == "" {
		return true
	}

	enabled, err := strconv.ParseBool(raw)
	if err != nil {
		log.Fatalf("invalid BENCH_USE_PREPARED=%q: %v", raw, err)
	}

	return enabled
}
