// Package db implements the database connection and management.
package db

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
	"time"

	"github.com/jackc/pgx/v5/tracelog"
	"github.com/xmtp/xmtpd/pkg/metrics"

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

var allowedNamespaceRe = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)

func waitUntilDBReady(ctx context.Context, db *pgxpool.Pool, waitTime time.Duration) error {
	pingCtx, cancel := context.WithTimeout(ctx, waitTime)
	defer cancel()

	err := db.Ping(pingCtx)
	if err != nil {
		return fmt.Errorf("database is not ready within %s: %w", waitTime, err)
	}
	return nil
}

func parseConfig(dsn string, statementTimeout time.Duration) (*pgxpool.Config, error) {
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

// NewNamespacedDB creates a new database with the given namespace if it doesn't exist and returns the full DSN for the new database.
func NewNamespacedDB(
	ctx context.Context,
	logger *zap.Logger,
	dsn string,
	namespace string,
	waitForDB time.Duration,
	statementTimeout time.Duration,
	prom *prometheus.Registry,
) (*sql.DB, error) {
	// Parse the DSN to get the config
	config, err := parseConfig(dsn, statementTimeout)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", parseDSNErrorMessage, err)
	}

	if err = createNamespace(ctx, config, namespace, waitForDB); err != nil {
		return nil, err
	}

	logger.Info(connectSuccessMessage, zap.String("namespace", namespace))

	config.ConnConfig.Database = namespace

	// enable SQL tracing
	config.ConnConfig.Tracer = &tracelog.TraceLog{
		Logger:   metrics.PromLogger{},
		LogLevel: tracelog.LogLevelTrace,
	}

	db, pool, err := newPGXDB(ctx, config, waitForDB)
	if err != nil {
		return nil, err
	}

	if prom != nil {

		mp, err := bindOTelToProm(prom)
		if err != nil {
			return nil, fmt.Errorf("bind OTel to Prom: %w", err)
		}

		if err := otelpgx.RecordStats(pool, otelpgx.WithStatsMeterProvider(mp)); err != nil {
			return nil, fmt.Errorf("otelpgx.RecordStats: %w", err)
		}
	}

	err = migrations.Migrate(ctx, db)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// ConnectToDB establishes a connection to an existing database using the provided DSN.
// Unlike NewNamespacedDB, this function does not create the database or run migrations.
// If namespace is provided, it overrides the database name in the DSN.
func ConnectToDB(
	ctx context.Context,
	logger *zap.Logger,
	dsn string,
	namespace string,
	waitForDB time.Duration,
	statementTimeout time.Duration,
	prom *prometheus.Registry,
) (*sql.DB, error) {
	config, err := parseConfig(dsn, statementTimeout)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", parseDSNErrorMessage, err)
	}

	if namespace != "" {
		config.ConnConfig.Database = namespace
	}

	// enable SQL tracing
	config.ConnConfig.Tracer = &tracelog.TraceLog{
		Logger:   metrics.PromLogger{},
		LogLevel: tracelog.LogLevelTrace,
	}

	db, pool, err := newPGXDB(ctx, config, waitForDB)
	if err != nil {
		return nil, err
	}

	if prom != nil {

		mp, err := bindOTelToProm(prom)
		if err != nil {
			return nil, fmt.Errorf("bind OTel to Prom: %w", err)
		}

		if err := otelpgx.RecordStats(pool, otelpgx.WithStatsMeterProvider(mp)); err != nil {
			return nil, fmt.Errorf("otelpgx.RecordStats: %w", err)
		}
	}

	logger.Info(connectSuccessMessage, zap.String("database", config.ConnConfig.Database))

	return db, nil
}

func bindOTelToProm(reg *prometheus.Registry) (*sdkmetric.MeterProvider, error) {
	exp, err := exporter.New(
		exporter.WithRegisterer(reg),
	)
	if err != nil {
		return nil, err
	}
	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(exp),
	)
	otel.SetMeterProvider(mp)
	return mp, nil
}
