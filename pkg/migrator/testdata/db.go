package testdata

import (
	"context"
	"database/sql"
	"encoding/csv"
	"encoding/hex"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

const (
	dbName     = "xmtpd"
	dbUser     = "xmtpd"
	dbPassword = "xmtpd"
	dbVersion  = "postgres:16-alpine"
)

func NewMigratorTestDB(t *testing.T, ctx context.Context) (db *sql.DB, dsn string, cleanup func()) {
	postgresContainer, err := postgres.Run(ctx,
		dbVersion,
		postgres.WithInitScripts(filepath.Join("testdata", "schema.sql")),
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		postgres.BasicWaitStrategies(),
	)
	require.NoError(t, err)

	dsn, err = postgresContainer.ConnectionString(
		ctx,
		"sslmode=disable",
	)
	require.NoError(t, err)

	db, err = sql.Open("postgres", dsn)
	require.NoError(t, err)

	insertKeyPackages(t, ctx, db)
	insertGroupMessages(t, ctx, db)
	insertWelcomeMessages(t, ctx, db)
	insertInboxLog(t, ctx, db)

	return db, dsn, func() {
		if err := postgresContainer.Terminate(ctx); err != nil {
			t.Logf("failed to terminate container: %s", err)
		}
	}
}

func decodeBytea(hexStr string) ([]byte, error) {
	hexStr = strings.TrimPrefix(hexStr, `\x`)
	return hex.DecodeString(hexStr)
}

func insertKeyPackages(t *testing.T, ctx context.Context, db *sql.DB) {
	f, err := os.Open(filepath.Join("testdata", "key_packages.csv"))
	require.NoError(t, err)
	defer f.Close()

	reader := csv.NewReader(f)
	rows, err := reader.ReadAll()
	require.NoError(t, err)

	for i, row := range rows {
		if i == 0 {
			continue
		}

		id, err := decodeBytea(row[1])
		require.NoError(t, err)

		keyPackage, err := decodeBytea(row[2])
		require.NoError(t, err)

		_, err = db.ExecContext(
			ctx,
			"INSERT INTO key_packages (sequence_id, installation_id, key_package) VALUES ($1, $2, $3)",
			row[0],
			id,
			keyPackage,
		)
		require.NoError(t, err)
	}
}

func insertGroupMessages(t *testing.T, ctx context.Context, db *sql.DB) {
	f, err := os.Open(filepath.Join("testdata", "group_messages.csv"))
	require.NoError(t, err)
	defer f.Close()

	reader := csv.NewReader(f)
	rows, err := reader.ReadAll()
	require.NoError(t, err)

	for i, row := range rows {
		if i == 0 {
			continue
		}

		groupID, err := decodeBytea(row[2])
		require.NoError(t, err)

		data, err := decodeBytea(row[3])
		require.NoError(t, err)

		hash, err := decodeBytea(row[4])
		require.NoError(t, err)

		senderHMAC, err := decodeBytea(row[6])
		require.NoError(t, err)

		_, err = db.ExecContext(
			ctx,
			"INSERT INTO group_messages (id, created_at, group_id, data, group_id_data_hash, is_commit, sender_hmac, should_push) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
			row[0],
			row[1],
			groupID,
			data,
			hash,
			row[5],
			senderHMAC,
			row[7],
		)
		require.NoError(t, err)
	}
}

func insertWelcomeMessages(t *testing.T, ctx context.Context, db *sql.DB) {
	f, err := os.Open(filepath.Join("testdata", "welcome_messages.csv"))
	require.NoError(t, err)
	defer f.Close()

	reader := csv.NewReader(f)
	rows, err := reader.ReadAll()
	require.NoError(t, err)

	for i, row := range rows {
		if i == 0 {
			continue
		}

		installationKey, err := decodeBytea(row[2])
		require.NoError(t, err)

		data, err := decodeBytea(row[3])
		require.NoError(t, err)

		hpkePublicKey, err := decodeBytea(row[4])
		require.NoError(t, err)

		installationKeyDataHash, err := decodeBytea(row[5])
		require.NoError(t, err)

		welcomeMetadata, err := decodeBytea(row[7])
		require.NoError(t, err)

		_, err = db.ExecContext(
			ctx,
			"INSERT INTO welcome_messages (id, created_at, installation_key, data, hpke_public_key, installation_key_data_hash, wrapper_algorithm, welcome_metadata) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
			row[0],
			row[1],
			installationKey,
			data,
			hpkePublicKey,
			installationKeyDataHash,
			row[6],
			welcomeMetadata,
		)
		require.NoError(t, err)
	}
}

func insertInboxLog(t *testing.T, ctx context.Context, db *sql.DB) {
	f, err := os.Open(filepath.Join("testdata", "inbox_log.csv"))
	require.NoError(t, err)
	defer f.Close()

	reader := csv.NewReader(f)
	rows, err := reader.ReadAll()
	require.NoError(t, err)

	for i, row := range rows {
		if i == 0 {
			continue
		}

		inboxID, err := decodeBytea(row[1])
		require.NoError(t, err)

		identityUpdateProto, err := decodeBytea(row[3])
		require.NoError(t, err)

		_, err = db.ExecContext(
			ctx,
			"INSERT INTO inbox_log (sequence_id, inbox_id, server_timestamp_ns, identity_update_proto) VALUES ($1, $2, $3, $4)",
			row[0],
			inboxID,
			row[2],
			identityUpdateProto,
		)
		require.NoError(t, err)
	}
}
