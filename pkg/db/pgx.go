package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/xmtp/xmtpd/pkg/migrations"
)

const MAX_NAMESPACE_LENGTH = 32

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
) (*sql.DB, error) {
	dbPool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	if err = waitUntilDBReady(ctx, dbPool, waitForDB); err != nil {
		return nil, err
	}

	db := stdlib.OpenDBFromPool(dbPool)

	return db, nil
}

func isValidNamespace(namespace string) error {
	if len(namespace) == 0 || len(namespace) > MAX_NAMESPACE_LENGTH {
		return fmt.Errorf(
			"namespace length must be between 1 and %d characters",
			MAX_NAMESPACE_LENGTH,
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
		return fmt.Errorf("failed to connect to postgres: %w", err)
	}
	defer adminConn.Close()

	if err = waitUntilDBReady(ctx, adminConn, waitForDB); err != nil {
		return err
	}

	// Create database if it doesn't exist
	_, err = adminConn.Exec(ctx, fmt.Sprintf(`CREATE DATABASE "%s"`, namespace))
	if err != nil {
		// Ignore error if database already exists
		var pgErr *pgconn.PgError
		// Error code 42P04 is for "duplicate database"
		// https://www.postgresql.org/docs/current/errcodes-appendix.html
		if errors.As(err, &pgErr) && pgErr.Code == "42P04" {
			return nil
		}

		return fmt.Errorf("failed to create database: %w", err)
	}

	return nil
}

// Creates a new database with the given namespace if it doesn't exist and returns the full DSN for the new database.
func NewNamespacedDB(
	ctx context.Context,
	dsn string,
	namespace string,
	waitForDB time.Duration, statementTimeout time.Duration,
) (*sql.DB, error) {
	// Parse the DSN to get the config
	config, err := parseConfig(dsn, statementTimeout)
	if err != nil {
		return nil, fmt.Errorf("failed to parse DSN: %w", err)
	}

	if err = createNamespace(ctx, config, namespace, waitForDB); err != nil {
		return nil, err
	}

	config.ConnConfig.Database = namespace

	db, err := newPGXDB(ctx, config, waitForDB)
	if err != nil {
		return nil, err
	}

	err = migrations.Migrate(ctx, db)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func NewDB(
	ctx context.Context,
	dsn string,
	waitForDB, statementTimeout time.Duration,
) (*sql.DB, error) {
	config, err := parseConfig(dsn, statementTimeout)
	if err != nil {
		return nil, err
	}

	db, err := newPGXDB(ctx, config, waitForDB)
	if err != nil {
		return nil, err
	}
	err = migrations.Migrate(ctx, db)
	if err != nil {
		return nil, err
	}
	return db, nil
}
