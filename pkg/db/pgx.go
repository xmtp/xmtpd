// Package db implements the database connection and management.
package db

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/tracelog"
	"github.com/xmtp/xmtpd/pkg/metrics"
	"github.com/xmtp/xmtpd/pkg/tracing"

	"github.com/exaring/otelpgx"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	exporter "go.opentelemetry.io/otel/exporters/prometheus"

	"go.uber.org/zap"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/xmtp/xmtpd/pkg/migrations"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

const (
	maxNamespaceLength = 32

	connectSuccessMessage = "successfully connected to database"
	parseDSNErrorMessage  = "failed to parse DSN"
)

var (
	bindOTelOnce sync.Once
	bindOTELErr  error
	boundMP      *sdkmetric.MeterProvider
)

var allowedNamespaceRe = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)

// apmQueryTracer implements pgx.QueryTracer to create Datadog APM spans for queries.
// This enables query-level visibility in flame graphs.
type apmQueryTracer struct {
	serviceName string
	role        string // "reader" or "writer" - critical for debugging replica issues
}

// TraceQueryStart creates a span when a query begins.
// Uses StartSpanFromContext to make DB queries children of the active span.
func (t *apmQueryTracer) TraceQueryStart(
	ctx context.Context,
	conn *pgx.Conn,
	data pgx.TraceQueryStartData,
) context.Context {
	// Use StartSpanFromContext so queries appear as children in flame graphs
	span, ctx := tracing.StartSpanFromContext(ctx, tracing.SpanDBQuery)
	tracing.SpanTag(span, tracing.TagDBSystem, "postgresql")
	tracing.SpanTag(span, tracing.TagDBService, t.serviceName)
	tracing.SpanTag(span, tracing.TagDBRole, t.role)
	tracing.SpanTag(span, tracing.TagDBStatement, data.SQL)
	tracing.SpanType(span, "sql")
	tracing.SpanResource(span, data.SQL)

	// Store span in context for TraceQueryEnd
	return context.WithValue(ctx, apmSpanKey{}, span)
}

// TraceQueryEnd finishes the span when query completes.
func (t *apmQueryTracer) TraceQueryEnd(
	ctx context.Context,
	conn *pgx.Conn,
	data pgx.TraceQueryEndData,
) {
	span, ok := ctx.Value(apmSpanKey{}).(tracing.Span)
	if !ok || span == nil {
		return
	}

	if data.Err != nil {
		span.Finish(tracing.WithError(data.Err))
	} else {
		tracing.SpanTag(span, tracing.TagDBRowsAffected, data.CommandTag.RowsAffected())
		span.Finish()
	}
}

type apmSpanKey struct{}

// compositeTracer combines multiple pgx tracers (logging + APM).
type compositeTracer struct {
	logTracer *tracelog.TraceLog
	apmTracer *apmQueryTracer
}

func (c *compositeTracer) TraceQueryStart(
	ctx context.Context,
	conn *pgx.Conn,
	data pgx.TraceQueryStartData,
) context.Context {
	// Call both tracers
	ctx = c.logTracer.TraceQueryStart(ctx, conn, data)
	ctx = c.apmTracer.TraceQueryStart(ctx, conn, data)
	return ctx
}

func (c *compositeTracer) TraceQueryEnd(
	ctx context.Context,
	conn *pgx.Conn,
	data pgx.TraceQueryEndData,
) {
	// Call both tracers
	c.logTracer.TraceQueryEnd(ctx, conn, data)
	c.apmTracer.TraceQueryEnd(ctx, conn, data)
}

func waitUntilDBReady(ctx context.Context, db *pgxpool.Pool, waitTime time.Duration) error {
	pingCtx, cancel := context.WithTimeout(ctx, waitTime)
	defer cancel()

	err := db.Ping(pingCtx)
	if err != nil {
		return fmt.Errorf("database is not ready within %s: %w", waitTime, err)
	}
	return nil
}

func parsePgxPoolConfig(dsn string, statementTimeout time.Duration) (*pgxpool.Config, error) {
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	config.ConnConfig.RuntimeParams["statement_timeout"] = fmt.Sprint(
		statementTimeout.Milliseconds(),
	)

	return config, nil
}

func newPGXDB(
	ctx context.Context,
	config *pgxpool.Config,
	waitForDB time.Duration,
) (*sql.DB, *pgxpool.Pool, error) {
	dbPool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, nil, err
	}
	if err = waitUntilDBReady(ctx, dbPool, waitForDB); err != nil {
		return nil, nil, err
	}
	db := stdlib.OpenDBFromPool(dbPool)
	return db, dbPool, nil
}

func isValidNamespace(namespace string) error {
	if len(namespace) == 0 || len(namespace) > maxNamespaceLength {
		return fmt.Errorf(
			"namespace length must be between 1 and %d characters",
			maxNamespaceLength,
		)
	}
	// PostgreSQL identifiers must start with a letter or underscore
	if !allowedNamespaceRe.MatchString(namespace) {
		return fmt.Errorf(
			"namespace must start with a letter or underscore and contain only letters, numbers, and underscores. Instead is %s",
			namespace,
		)
	}
	return nil
}

// Creates a new database with the given namespace if it doesn't exist
func createNamespace(
	ctx context.Context,
	config *pgxpool.Config,
	namespace string,
	waitForDB time.Duration,
) error {
	if err := isValidNamespace(namespace); err != nil {
		return err
	}

	// Make a copy of the config so we don't dirty it
	config = config.Copy()
	// Change the database to postgres so we are able to create new DBs
	config.ConnConfig.Database = "postgres"

	// Create a temporary connection to the postgres DB
	adminConn, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer adminConn.Close()

	if err = waitUntilDBReady(ctx, adminConn, waitForDB); err != nil {
		return err
	}

	var exists bool
	err = adminConn.QueryRow(ctx,
		"SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1)",
		namespace).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check if database exists: %w", err)
	}

	if !exists {
		_, err = adminConn.Exec(ctx, fmt.Sprintf(`CREATE DATABASE "%s"`, namespace))
		if err != nil {
			return fmt.Errorf("failed to create database: %w", err)
		}
	}

	return nil
}

type dbConnConfig struct {
	pingTimeout        time.Duration
	statementTimeout   time.Duration
	prometheusRegistry *prometheus.Registry
	createNamespace    bool
	runMigrations      bool
	role               string // "reader" or "writer" for APM tagging
}

type dbOptionFunc func(*dbConnConfig)

func dbPingTimeout(d time.Duration) dbOptionFunc {
	return func(cfg *dbConnConfig) {
		cfg.pingTimeout = d
	}
}

func dbStatementTimeout(d time.Duration) dbOptionFunc {
	return func(cfg *dbConnConfig) {
		cfg.statementTimeout = d
	}
}

// doCreateNamespace will create the namespace if it does not exist.
func doCreateNamespace(b bool) dbOptionFunc {
	return func(cfg *dbConnConfig) {
		cfg.createNamespace = b
	}
}

func runMigrations(b bool) dbOptionFunc {
	return func(cfg *dbConnConfig) {
		cfg.runMigrations = b
	}
}

func prometheusRegistry(p *prometheus.Registry) dbOptionFunc {
	return func(cfg *dbConnConfig) {
		cfg.prometheusRegistry = p
	}
}

// dbRole sets the role tag for APM spans (reader or writer).
// This is critical for debugging read-replica issues.
func dbRole(role string) dbOptionFunc {
	return func(cfg *dbConnConfig) {
		cfg.role = role
	}
}

func connectToDB(
	ctx context.Context,
	logger *zap.Logger,
	dsn string,
	namespace string,
	opts ...dbOptionFunc,
) (*sql.DB, error) {
	var cfg dbConnConfig
	for _, opt := range opts {
		opt(&cfg)
	}

	poolcfg, err := parsePgxPoolConfig(dsn, cfg.statementTimeout)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", parseDSNErrorMessage, err)
	}

	if cfg.createNamespace {
		err := createNamespace(ctx, poolcfg, namespace, cfg.pingTimeout)
		if err != nil {
			return nil, fmt.Errorf("could not create namespace (name: %v): %w", namespace, err)
		}
	}

	if namespace != "" {
		poolcfg.ConnConfig.Database = namespace
	}

	// Determine service name for APM spans
	serviceName := "xmtpd-db"
	if namespace != "" {
		serviceName = "xmtpd-db-" + namespace
	}

	// Default role to "writer" if not specified
	role := cfg.role
	if role == "" {
		role = "writer"
	}

	// Set up SQL tracing. When APM is enabled, use a composite tracer
	// (Prometheus metrics logging + Datadog APM spans). Otherwise, keep
	// only the existing Prometheus metrics logger to avoid per-query overhead.
	logTracer := &tracelog.TraceLog{
		Logger:   metrics.PromLogger{},
		LogLevel: tracelog.LogLevelTrace,
	}

	if tracing.IsEnabled() {
		poolcfg.ConnConfig.Tracer = &compositeTracer{
			logTracer: logTracer,
			apmTracer: &apmQueryTracer{
				serviceName: serviceName,
				role:        role, // reader or writer - critical for replica debugging
			},
		}
	} else {
		poolcfg.ConnConfig.Tracer = logTracer
	}

	db, pool, err := newPGXDB(ctx, poolcfg, cfg.pingTimeout)
	if err != nil {
		return nil, err
	}

	logger.Info(connectSuccessMessage, zap.String("namespace", namespace))

	if cfg.prometheusRegistry != nil {
		mp, err := bindOTelToProm(cfg.prometheusRegistry)
		if err != nil {
			return nil, fmt.Errorf("bind OTel to Prom: %w", err)
		}

		if err := otelpgx.RecordStats(pool, otelpgx.WithStatsMeterProvider(mp)); err != nil {
			return nil, fmt.Errorf("otelpgx.RecordStats: %w", err)
		}
	}

	if cfg.runMigrations {
		err = migrations.Migrate(ctx, db)
		if err != nil {
			return nil, fmt.Errorf("could not run migrations: %w", err)
		}
	}

	return db, nil
}

// NewNamespacedDB creates a new database with the given namespace if it doesn't exist and returns the full DSN for the new database.
// This is typically used for the writer connection (creates schema, runs migrations).
func NewNamespacedDB(
	ctx context.Context,
	logger *zap.Logger,
	dsn string,
	namespace string,
	waitForDB time.Duration,
	statementTimeout time.Duration,
	prom *prometheus.Registry,
) (*sql.DB, error) {
	return connectToDB(
		ctx,
		logger,
		dsn,
		namespace,
		dbPingTimeout(waitForDB),
		dbStatementTimeout(statementTimeout),
		prometheusRegistry(prom),
		doCreateNamespace(true),
		runMigrations(true),
		dbRole("writer"),
	)
}

// NewNamespacedReaderDB is like NewNamespacedDB but tags the connection as a
// read replica for APM. The "reader" role helps debug read-replica lag issues.
func NewNamespacedReaderDB(
	ctx context.Context,
	logger *zap.Logger,
	dsn string,
	namespace string,
	waitForDB time.Duration,
	statementTimeout time.Duration,
	prom *prometheus.Registry,
) (*sql.DB, error) {
	return connectToDB(
		ctx,
		logger,
		dsn,
		namespace,
		dbPingTimeout(waitForDB),
		dbStatementTimeout(statementTimeout),
		prometheusRegistry(prom),
		doCreateNamespace(false),
		runMigrations(false),
		dbRole("reader"),
	)
}

// ConnectToDB establishes a connection to an existing database using the provided DSN.
// Unlike NewNamespacedDB, this function does not create the database or run migrations.
// If namespace is provided, it overrides the database name in the DSN.
func ConnectToDB(ctx context.Context,
	logger *zap.Logger,
	dsn string,
	namespace string,
	waitForDB time.Duration,
	statementTimeout time.Duration,
	prom *prometheus.Registry,
) (*sql.DB, error) {
	return connectToDB(ctx, logger, dsn, namespace,
		dbPingTimeout(waitForDB),
		dbStatementTimeout(statementTimeout),
		prometheusRegistry(prom),
		dbRole("writer"),
		// Not creating namespace.
		// Not running migrations.
	)
}

func bindOTelToProm(reg *prometheus.Registry) (*sdkmetric.MeterProvider, error) {
	bindOTelOnce.Do(func() {
		exp, err := exporter.New(exporter.WithRegisterer(reg))
		if err != nil {
			bindOTELErr = err
			return
		}
		mp := sdkmetric.NewMeterProvider(sdkmetric.WithReader(exp))
		otel.SetMeterProvider(mp)
		boundMP = mp
	})
	return boundMP, bindOTELErr
}
