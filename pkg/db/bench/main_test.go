//go:build bench

package bench

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/xmtp/xmtpd/pkg/db/migrations"
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
	count       int      // total envelope rows seeded
	topics      [][]byte // topics seeded into this DB
	originators []int32  // originator node IDs seeded
	payerIDs    []int32  // payer IDs created for batch insert benchmarks
}

// Package-level handles, populated by TestMain.
// testing.M has no Context() or Cleanup() methods, so we manage lifecycle manually.
var (
	benchCtx      = context.Background()
	envelopeTiers []*envelopeTier
	congestionDB  *sql.DB
	ledgerDB      *sql.DB
	indexerDB     *sql.DB
	usageDB       *sql.DB
	hotPathDB     *sql.DB

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
	ctlDB, err := connectDB(localDSNPrefix + localDSNSuffix)
	if err != nil {
		log.Fatalf("failed to connect to control DB: %v", err)
	}
	defer func() { _ = ctlDB.Close() }()

	// Track cleanups so they run even if a seed function calls log.Fatalf
	// (log.Fatalf calls os.Exit, so defers in this function won't run either —
	// but at least we get cleanup for the normal exit path).
	var cleanups []func()
	defer func() {
		for _, fn := range cleanups {
			fn()
		}
	}()

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
	var cleanup func()
	congestionDB, cleanup = createBenchDB(ctlDB, "congestion")
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

	hotPathDB, cleanup = createBenchDB(ctlDB, "hot_path")
	cleanups = append(cleanups, cleanup)
	seedHotPath(benchCtx, hotPathDB)

	os.Exit(m.Run())
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
		_ = db.Close()
		_, _ = ctlDB.Exec("DROP DATABASE " + dbName + " WITH (FORCE)")
	}
}
