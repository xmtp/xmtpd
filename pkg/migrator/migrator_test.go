package migrator_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/config"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/migrator"
	"github.com/xmtp/xmtpd/pkg/migrator/testdata"
	"github.com/xmtp/xmtpd/pkg/testutils"
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
	db       *sql.DB
}

func newMigratorTest(t *testing.T) *migratorTest {
	var (
		ctx             = t.Context()
		writerDB, _     = testutils.NewDB(t, ctx)
		_, dsn, cleanup = testdata.NewMigratorTestDB(t, ctx)
		chainConfig     = testdata.NewMigratorBlockchain(t)
		nodePrivateKey  = testutils.RandomPrivateKey(t)
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
	)
	require.NoError(t, err)

	return &migratorTest{
		ctx:      ctx,
		cleanup:  cleanup,
		migrator: migrator,
		db:       writerDB,
	}
}

func TestMigrator(t *testing.T) {
	test := newMigratorTest(t)
	defer test.cleanup()

	require.NoError(t, test.migrator.Start())

	require.Eventually(t, func() bool {
		return checkMigrationTrackerState(test.ctx, test.db) &&
			checkGatewayEnvelopesLastID(test.ctx, test.db) &&
			checkGatewayEnvelopesMigratedAmount(test.ctx, test.db) &&
			checkGatewayEnvelopesAreUnique(test.ctx, test.db)
	}, 20*time.Second, 50*time.Millisecond)

	require.NoError(t, test.migrator.Stop())
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

func checkGatewayEnvelopesLastID(ctx context.Context, db *sql.DB) bool {
	groupMsgSeqID, err := getGatewayEnvelopesLastSequenceID(
		ctx,
		db,
		int32(migrator.GroupMessageOriginatorID),
	)
	if err != nil || groupMsgSeqID != groupMessageLastID {
		return false
	}

	welcomeMsgSeqID, err := getGatewayEnvelopesLastSequenceID(
		ctx,
		db,
		int32(migrator.WelcomeMessageOriginatorID),
	)
	if err != nil || welcomeMsgSeqID != welcomeMessageLastID {
		return false
	}

	keyPkgSeqID, err := getGatewayEnvelopesLastSequenceID(
		ctx,
		db,
		int32(migrator.KeyPackagesOriginatorID),
	)
	if err != nil || keyPkgSeqID != keyPackageLastID {
		return false
	}

	return true
}

func checkGatewayEnvelopesMigratedAmount(ctx context.Context, db *sql.DB) bool {
	groupMsgAmount, err := getGatewayEnvelopesAmount(
		ctx,
		db,
		int32(migrator.GroupMessageOriginatorID),
	)
	if err != nil || groupMsgAmount != groupMessageAmount {
		return false
	}

	inboxAmount, err := getGatewayEnvelopesAmount(ctx, db, int32(migrator.InboxLogOriginatorID))
	if err != nil || inboxAmount != inboxLogAmount {
		return false
	}

	welcomeMsgAmount, err := getGatewayEnvelopesAmount(
		ctx,
		db,
		int32(migrator.WelcomeMessageOriginatorID),
	)
	if err != nil || welcomeMsgAmount != welcomeMessageAmount {
		return false
	}

	keyPkgAmount, err := getGatewayEnvelopesAmount(ctx, db, int32(migrator.KeyPackagesOriginatorID))
	if err != nil || keyPkgAmount != keyPackageAmount {
		return false
	}

	return true
}

func checkGatewayEnvelopesAreUnique(ctx context.Context, db *sql.DB) bool {
	groupMsgUnique, err := getGatewayEnvelopesUniqueAmount(
		ctx,
		db,
		int32(migrator.GroupMessageOriginatorID),
	)
	if err != nil || groupMsgUnique != groupMessageAmount {
		return false
	}

	inboxUnique, err := getGatewayEnvelopesUniqueAmount(
		ctx,
		db,
		int32(migrator.InboxLogOriginatorID),
	)
	if err != nil || inboxUnique != inboxLogAmount {
		return false
	}

	welcomeMsgUnique, err := getGatewayEnvelopesUniqueAmount(
		ctx,
		db,
		int32(migrator.WelcomeMessageOriginatorID),
	)
	if err != nil || welcomeMsgUnique != welcomeMessageAmount {
		return false
	}

	keyPkgUnique, err := getGatewayEnvelopesUniqueAmount(
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
	ctx context.Context,
	db *sql.DB,
	originatorNodeID int32,
) (int64, error) {
	querier := queries.New(db)

	lastSequenceID, err := querier.GetLatestSequenceId(ctx, originatorNodeID)
	if err != nil {
		return 0, err
	}

	return lastSequenceID, nil
}

func getGatewayEnvelopesAmount(
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

	return count, nil
}

func getGatewayEnvelopesUniqueAmount(
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

	return count, nil
}
