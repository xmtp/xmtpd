package migrator_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/config"
	"github.com/xmtp/xmtpd/pkg/db"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/envelopes"
	"github.com/xmtp/xmtpd/pkg/migrator"
	"github.com/xmtp/xmtpd/pkg/migrator/testdata"
	"github.com/xmtp/xmtpd/pkg/testutils"
	"github.com/xmtp/xmtpd/pkg/testutils/fees"
	"github.com/xmtp/xmtpd/pkg/utils"
)

// TODO: Check 19 identity updates are published to the blockchain.
// TODO: Check 9 group messages are published to the blockchain.

const (
	// Source: 19 total group_messages rows in CSV.
	// - 7 commit messages (ids 1-7, is_commit=true) -> blockchain.
	// - 12 non-commit messages (ids 8-19, is_commit=false) -> database.
	groupMessageAmount int64 = 12
	groupMessageLastID int64 = 19

	// Commit messages go to blockchain, not database.
	commitMessageLastID int64 = 7

	// Identity updates go to blockchain, not database.
	inboxLogAmount       int64 = 0
	inboxLogLastID       int64 = 19
	welcomeMessageAmount int64 = 19
	welcomeMessageLastID int64 = 19
	keyPackageAmount     int64 = 19
	keyPackageLastID     int64 = 19
)

type migratorTest struct {
	ctx      context.Context
	cleanup  func()
	migrator *migrator.Migrator
	db       *db.Handler
	sourceDB *sql.DB
}

func newMigratorTest(t *testing.T) *migratorTest {
	var (
		ctx                    = t.Context()
		writerDB, _            = testutils.NewDB(t, ctx)
		sourceDB, dsn, cleanup = testdata.NewMigratorTestDB(t, ctx)
		chainConfig            = testdata.NewMigratorBlockchain(t)
		nodePrivateKey         = testutils.RandomPrivateKey(t)
		feeCalculator          = fees.NewTestFeeCalculator()
	)

	payerPrivateKey, err := crypto.HexToECDSA(testdata.PayerPrivateKeyString)
	require.NoError(t, err)

	migrator, err := migrator.NewMigrationService(
		migrator.WithContext(ctx),
		migrator.WithLogger(testutils.NewLog(t)),
		migrator.WithDestinationDB(writerDB),
		migrator.WithMigratorConfig(&config.MigrationServerOptions{
			Enable:                 true,
			PayerPrivateKey:        utils.EcdsaPrivateKeyToString(payerPrivateKey),
			NodeSigningKey:         utils.EcdsaPrivateKeyToString(nodePrivateKey),
			ReaderConnectionString: dsn,
			ReaderTimeout:          1 * time.Second,
			WaitForDB:              5 * time.Second,
			BatchSize:              1000,
			PollInterval:           500 * time.Millisecond,
			StartDate:              startDate,
		}),
		migrator.WithContractsOptions(chainConfig),
		migrator.WithFeeCalculator(feeCalculator),
	)
	require.NoError(t, err)

	return &migratorTest{
		ctx:      ctx,
		cleanup:  cleanup,
		migrator: migrator,
		db:       writerDB,
		sourceDB: sourceDB,
	}
}

func TestMigrator(t *testing.T) {
	test := newMigratorTest(t)
	defer test.cleanup()

	require.NoError(t, test.migrator.Start())

	var (
		destDB   = test.db.DB()
		sourceDB = test.sourceDB
		ctx      = test.ctx
	)

	// Wait for migration to complete and stop the migrator.
	require.Eventually(t, func() bool {
		return checkMigrationTrackerState(ctx, destDB)
	}, 20*time.Second, 50*time.Millisecond)

	require.NoError(t, test.migrator.Stop())

	t.Run("gateway_envelopes_last_id", func(t *testing.T) {
		checkGatewayEnvelopesLastID(t, ctx, destDB)
	})

	t.Run("gateway_envelopes_migrated_amount", func(t *testing.T) {
		checkGatewayEnvelopesMigratedAmount(t, ctx, destDB)
	})

	t.Run("gateway_envelopes_are_unique", func(t *testing.T) {
		checkGatewayEnvelopesAreUnique(t, ctx, destDB)
	})

	// Verify group messages data integrity (non-commit only).
	t.Run("group_messages", func(t *testing.T) {
		verifyGroupMessagesIntegrity(t, ctx, sourceDB, destDB)
	})

	// Verify welcome messages data integrity.
	t.Run("welcome_messages", func(t *testing.T) {
		verifyWelcomeMessagesIntegrity(t, ctx, sourceDB, destDB)
	})

	// Verify key packages data integrity.
	t.Run("key_packages", func(t *testing.T) {
		verifyKeyPackagesIntegrity(t, ctx, sourceDB, destDB)
	})
}

func checkMigrationTrackerState(ctx context.Context, db *sql.DB) bool {
	rows, err := db.QueryContext(ctx, "SELECT * FROM migration_tracker")
	if err != nil {
		return false
	}

	defer func() {
		_ = rows.Close()
	}()

	state := make(map[string]int64)

	for rows.Next() {
		var (
			tableName      string
			lastMigratedID int64
			createdAt      time.Time
			updatedAt      time.Time
		)

		if err := rows.Scan(&tableName, &lastMigratedID, &createdAt, &updatedAt); err != nil {
			return false
		}

		state[tableName] = lastMigratedID
	}

	if rows.Err() != nil {
		return false
	}

	return state["group_messages"] == groupMessageLastID &&
		state["welcome_messages"] == welcomeMessageLastID &&
		state["inbox_log"] == inboxLogLastID &&
		state["key_packages"] == keyPackageLastID &&
		state["commit_messages"] == commitMessageLastID
}

func checkGatewayEnvelopesLastID(t *testing.T, ctx context.Context, db *sql.DB) bool {
	groupMsgSeqID, err := getGatewayEnvelopesLastSequenceID(
		t,
		ctx,
		db,
		int32(migrator.GroupMessageOriginatorID),
	)
	if err != nil || groupMsgSeqID != groupMessageLastID {
		return false
	}

	welcomeMsgSeqID, err := getGatewayEnvelopesLastSequenceID(
		t,
		ctx,
		db,
		int32(migrator.WelcomeMessageOriginatorID),
	)
	if err != nil || welcomeMsgSeqID != welcomeMessageLastID {
		return false
	}

	keyPkgSeqID, err := getGatewayEnvelopesLastSequenceID(
		t,
		ctx,
		db,
		int32(migrator.KeyPackagesOriginatorID),
	)
	if err != nil || keyPkgSeqID != keyPackageLastID {
		return false
	}

	return true
}

func checkGatewayEnvelopesMigratedAmount(t *testing.T, ctx context.Context, db *sql.DB) bool {
	groupMsgAmount, err := getGatewayEnvelopesAmount(
		t,
		ctx,
		db,
		int32(migrator.GroupMessageOriginatorID),
	)
	if err != nil || groupMsgAmount != groupMessageAmount {
		return false
	}

	welcomeMsgAmount, err := getGatewayEnvelopesAmount(
		t,
		ctx,
		db,
		int32(migrator.WelcomeMessageOriginatorID),
	)
	if err != nil || welcomeMsgAmount != welcomeMessageAmount {
		return false
	}

	keyPkgAmount, err := getGatewayEnvelopesAmount(
		t,
		ctx,
		db,
		int32(migrator.KeyPackagesOriginatorID),
	)
	if err != nil || keyPkgAmount != keyPackageAmount {
		return false
	}

	return true
}

func checkGatewayEnvelopesAreUnique(t *testing.T, ctx context.Context, db *sql.DB) bool {
	groupMsgUnique, err := getGatewayEnvelopesUniqueAmount(
		t,
		ctx,
		db,
		int32(migrator.GroupMessageOriginatorID),
	)
	if err != nil || groupMsgUnique != groupMessageAmount {
		return false
	}

	welcomeMsgUnique, err := getGatewayEnvelopesUniqueAmount(
		t,
		ctx,
		db,
		int32(migrator.WelcomeMessageOriginatorID),
	)
	if err != nil || welcomeMsgUnique != welcomeMessageAmount {
		return false
	}

	keyPkgUnique, err := getGatewayEnvelopesUniqueAmount(
		t,
		ctx,
		db,
		int32(migrator.KeyPackagesOriginatorID),
	)
	if err != nil || keyPkgUnique != keyPackageAmount {
		return false
	}

	return true
}

func getGatewayEnvelopesLastSequenceID(
	t *testing.T,
	ctx context.Context,
	db *sql.DB,
	originatorNodeID int32,
) (int64, error) {
	querier := queries.New(db)

	lastSequenceID, err := querier.GetLatestSequenceId(ctx, originatorNodeID)
	if err != nil {
		return 0, err
	}

	t.Logf(
		"verified last sequence id %d for originator node id %d",
		lastSequenceID,
		originatorNodeID,
	)

	return lastSequenceID, nil
}

func getGatewayEnvelopesAmount(
	t *testing.T,
	ctx context.Context,
	db *sql.DB,
	originatorNodeID int32,
) (int64, error) {
	var (
		count              int64
		getEnvelopesAmount = `SELECT COUNT(*)::BIGINT
FROM gateway_envelopes_meta
WHERE originator_node_id = $1`
	)

	row := db.QueryRowContext(ctx, getEnvelopesAmount, originatorNodeID)
	if err := row.Scan(&count); err != nil {
		return 0, err
	}

	t.Logf("verified %d envelope amounts for originator node id %d", count, originatorNodeID)

	return count, nil
}

func getGatewayEnvelopesUniqueAmount(
	t *testing.T,
	ctx context.Context,
	db *sql.DB,
	originatorNodeID int32,
) (int64, error) {
	var (
		count              int64
		getEnvelopesAmount = `SELECT COUNT(DISTINCT originator_sequence_id)::BIGINT
FROM gateway_envelopes_meta
WHERE originator_node_id = $1`
	)

	row := db.QueryRowContext(ctx, getEnvelopesAmount, originatorNodeID)
	if err := row.Scan(&count); err != nil {
		return 0, err
	}

	t.Logf("verified %d unique envelope amounts for originator node id %d", count, originatorNodeID)

	return count, nil
}

// verifyGroupMessagesIntegrity checks that non-commit group messages are correctly migrated.
func verifyGroupMessagesIntegrity(t *testing.T, ctx context.Context, sourceDB, destDB *sql.DB) {
	// Query source: only non-commit messages go to database.
	sourceRows, err := sourceDB.QueryContext(ctx, `
		SELECT id, group_id, data 
		FROM group_messages 
		WHERE is_commit = false 
		ORDER BY id
	`)
	require.NoError(t, err)
	defer func() {
		_ = sourceRows.Close()
	}()

	var sourceRecords []struct {
		ID      int64
		GroupID []byte
		Data    []byte
	}

	for sourceRows.Next() {
		var rec struct {
			ID      int64
			GroupID []byte
			Data    []byte
		}
		require.NoError(t, sourceRows.Scan(&rec.ID, &rec.GroupID, &rec.Data))
		sourceRecords = append(sourceRecords, rec)
	}
	require.NoError(t, sourceRows.Err())

	// For each source record, verify it exists in destination and envelope is valid.
	for _, src := range sourceRecords {
		envelope, err := getDestinationEnvelope(
			ctx,
			destDB,
			int32(migrator.GroupMessageOriginatorID),
			src.ID,
		)
		require.NoError(t, err, "failed to get envelope for source id %d", src.ID)
		require.NotNil(t, envelope, "missing envelope for source id %d", src.ID)

		// Verify the envelope can be decoded.
		origEnv, err := envelopes.NewOriginatorEnvelopeFromBytes(envelope)
		require.NoError(t, err, "failed to decode envelope for source id %d", src.ID)

		// Verify sequence ID matches source ID.
		require.Equal(t, uint64(src.ID), origEnv.OriginatorSequenceID(),
			"sequence ID mismatch for source id %d", src.ID)
	}

	t.Logf("verified %d group messages", len(sourceRecords))
}

// verifyWelcomeMessagesIntegrity checks that welcome messages are correctly migrated.
func verifyWelcomeMessagesIntegrity(t *testing.T, ctx context.Context, sourceDB, destDB *sql.DB) {
	sourceRows, err := sourceDB.QueryContext(ctx, `
		SELECT id, installation_key, data 
		FROM welcome_messages 
		ORDER BY id
	`)
	require.NoError(t, err)
	defer func() {
		_ = sourceRows.Close()
	}()

	var sourceRecords []struct {
		ID              int64
		InstallationKey []byte
		Data            []byte
	}

	for sourceRows.Next() {
		var rec struct {
			ID              int64
			InstallationKey []byte
			Data            []byte
		}
		require.NoError(t, sourceRows.Scan(&rec.ID, &rec.InstallationKey, &rec.Data))
		sourceRecords = append(sourceRecords, rec)
	}
	require.NoError(t, sourceRows.Err())

	for _, src := range sourceRecords {
		envelope, err := getDestinationEnvelope(
			ctx,
			destDB,
			int32(migrator.WelcomeMessageOriginatorID),
			src.ID,
		)
		require.NoError(t, err, "failed to get envelope for source id %d", src.ID)
		require.NotNil(t, envelope, "missing envelope for source id %d", src.ID)

		origEnv, err := envelopes.NewOriginatorEnvelopeFromBytes(envelope)
		require.NoError(t, err, "failed to decode envelope for source id %d", src.ID)

		require.Equal(t, uint64(src.ID), origEnv.OriginatorSequenceID(),
			"sequence ID mismatch for source id %d", src.ID)
	}

	t.Logf("verified %d welcome messages", len(sourceRecords))
}

// verifyKeyPackagesIntegrity checks that key packages are correctly migrated.
func verifyKeyPackagesIntegrity(t *testing.T, ctx context.Context, sourceDB, destDB *sql.DB) {
	sourceRows, err := sourceDB.QueryContext(ctx, `
		SELECT sequence_id, installation_id, key_package 
		FROM key_packages 
		ORDER BY sequence_id
	`)
	require.NoError(t, err)
	defer func() {
		_ = sourceRows.Close()
	}()

	var sourceRecords []struct {
		SequenceID     int64
		InstallationID []byte
		KeyPackage     []byte
	}

	for sourceRows.Next() {
		var rec struct {
			SequenceID     int64
			InstallationID []byte
			KeyPackage     []byte
		}
		require.NoError(t, sourceRows.Scan(&rec.SequenceID, &rec.InstallationID, &rec.KeyPackage))
		sourceRecords = append(sourceRecords, rec)
	}
	require.NoError(t, sourceRows.Err())

	for _, src := range sourceRecords {
		envelope, err := getDestinationEnvelope(
			ctx,
			destDB,
			int32(migrator.KeyPackagesOriginatorID),
			src.SequenceID,
		)
		require.NoError(t, err, "failed to get envelope for source sequence_id %d", src.SequenceID)
		require.NotNil(t, envelope, "missing envelope for source sequence_id %d", src.SequenceID)

		origEnv, err := envelopes.NewOriginatorEnvelopeFromBytes(envelope)
		require.NoError(
			t,
			err,
			"failed to decode envelope for source sequence_id %d",
			src.SequenceID,
		)

		require.Equal(t, uint64(src.SequenceID), origEnv.OriginatorSequenceID(),
			"sequence ID mismatch for source sequence_id %d", src.SequenceID)
	}

	t.Logf("verified %d key packages", len(sourceRecords))
}

// getDestinationEnvelope retrieves the originator envelope from the destination database.
func getDestinationEnvelope(
	ctx context.Context,
	db *sql.DB,
	originatorNodeID int32,
	sequenceID int64,
) ([]byte, error) {
	var envelope []byte
	query := `
		SELECT b.originator_envelope 
		FROM gateway_envelopes_meta m
		JOIN gateway_envelope_blobs b 
			ON m.originator_node_id = b.originator_node_id 
			AND m.originator_sequence_id = b.originator_sequence_id
		WHERE m.originator_node_id = $1 
			AND m.originator_sequence_id = $2
	`
	err := db.QueryRowContext(ctx, query, originatorNodeID, sequenceID).Scan(&envelope)
	if err != nil {
		return nil, err
	}
	return envelope, nil
}
