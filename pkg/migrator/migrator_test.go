package migrator_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/config"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/migrator"
	"github.com/xmtp/xmtpd/pkg/migrator/testdata"
	"github.com/xmtp/xmtpd/pkg/mlsvalidate"
	mlsvalidateMock "github.com/xmtp/xmtpd/pkg/mocks/mlsvalidate"
	"github.com/xmtp/xmtpd/pkg/proto/identity/associations"
	"github.com/xmtp/xmtpd/pkg/testutils"
	"github.com/xmtp/xmtpd/pkg/utils"
)

const (
	// Note that there must be 18 address_log entries in total.
	// In the testdata set there are 2 identity updates, updating the same inboxID.
	// Check rows 6 and 7 in testdata/inbox_log.csv
	addressLogAmount     int64 = 18
	groupMessageAmount   int64 = 19
	groupMessageLastID   int64 = 19
	inboxLogAmount       int64 = 19
	inboxLogLastID       int64 = 19
	welcomeMessageAmount int64 = 19
	welcomeMessageLastID int64 = 19
	installationAmount   int64 = 19
	installationLastID   int64 = 1717490371754970003
)

type migratorTest struct {
	ctx      context.Context
	cleanup  func()
	migrator *migrator.Migrator
	db       *sql.DB
}

func newMigratorTest(t *testing.T) *migratorTest {
	var (
		ctx                  = t.Context()
		writerDB, _          = testutils.NewDB(t, ctx)
		_, dsn, cleanup      = testdata.NewMigratorTestDB(t, ctx)
		mlsValidationService = mlsvalidateMock.NewMockMLSValidationService(
			t,
		)
		payerPrivateKey = testutils.RandomPrivateKey(t)
		nodePrivateKey  = testutils.RandomPrivateKey(t)
	)

	mlsValidationService.EXPECT().
		GetAssociationStateFromEnvelopes(mock.Anything, mock.Anything, mock.Anything).
		Return(&mlsvalidate.AssociationStateResult{
			StateDiff: &associations.AssociationStateDiff{
				NewMembers: []*associations.MemberIdentifier{{
					Kind: &associations.MemberIdentifier_EthereumAddress{
						EthereumAddress: "0x12345",
					},
				}},
			},
		}, nil)

	migrator, err := migrator.NewMigrationService(
		migrator.WithContext(ctx),
		migrator.WithLogger(testutils.NewLog(t)),
		migrator.WithDestinationDB(writerDB),
		migrator.WithMLSValidationService(mlsValidationService),
		migrator.WithMigratorConfig(&config.MigratorOptions{
			Enable:                 true,
			PayerPrivateKey:        utils.EcdsaPrivateKeyToString(payerPrivateKey),
			NodeSigningKey:         utils.EcdsaPrivateKeyToString(nodePrivateKey),
			ReaderConnectionString: dsn,
			ReadTimeout:            1 * time.Second,
			WaitForDB:              5 * time.Second,
			BatchSize:              1000,
			PollInterval:           500 * time.Millisecond,
		}),
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

	<-time.After(1 * time.Second)

	checkMigrationTrackerState(t, test.ctx, test.db)
	checkGatewayEnvelopesLastID(t, test.ctx, test.db)
	checkGatewayEnvelopesMigratedAmount(t, test.ctx, test.db)
	checkGatewayEnvelopesAreUnique(t, test.ctx, test.db)
	checkAddressLogAmount(t, test.ctx, test.db)

	require.NoError(t, test.migrator.Stop())
}

func checkMigrationTrackerState(t *testing.T, ctx context.Context, db *sql.DB) {
	rows, err := db.QueryContext(ctx, "SELECT * FROM migration_tracker")
	require.NoError(t, err)

	defer func() {
		err := rows.Close()
		require.NoError(t, err)
	}()

	state := make(map[string]int64)

	for rows.Next() {
		var (
			tableName      string
			lastMigratedID int64
			createdAt      time.Time
			updatedAt      time.Time
		)

		err := rows.Scan(&tableName, &lastMigratedID, &createdAt, &updatedAt)
		require.NoError(t, err)

		state[tableName] = lastMigratedID
	}

	require.NoError(t, rows.Err())

	require.Equal(t, groupMessageLastID, state["group_messages"])
	require.Equal(t, welcomeMessageLastID, state["welcome_messages"])
	require.Equal(t, inboxLogLastID, state["inbox_log"])
	require.Equal(t, installationLastID, state["installations"])
}

func checkGatewayEnvelopesLastID(t *testing.T, ctx context.Context, db *sql.DB) {
	require.Equal(t, groupMessageLastID, getGatewayEnvelopesLastSequenceID(
		t,
		ctx,
		db,
		int32(migrator.GroupMessageOriginatorID),
	))

	require.Equal(t, inboxLogLastID, getGatewayEnvelopesLastSequenceID(
		t,
		ctx,
		db,
		int32(migrator.InboxLogOriginatorID),
	))

	require.Equal(t, welcomeMessageLastID, getGatewayEnvelopesLastSequenceID(
		t,
		ctx,
		db,
		int32(migrator.WelcomeMessageOriginatorID),
	))

	require.Equal(t, installationLastID, getGatewayEnvelopesLastSequenceID(
		t,
		ctx,
		db,
		int32(migrator.InstallationOriginatorID),
	))
}

func checkGatewayEnvelopesMigratedAmount(t *testing.T, ctx context.Context, db *sql.DB) {
	require.Equal(t, groupMessageAmount, getGatewayEnvelopesAmount(
		t,
		ctx,
		db,
		int32(migrator.GroupMessageOriginatorID),
	))

	require.Equal(t, inboxLogAmount, getGatewayEnvelopesAmount(
		t,
		ctx,
		db,
		int32(migrator.InboxLogOriginatorID),
	))

	require.Equal(t, welcomeMessageAmount, getGatewayEnvelopesAmount(
		t,
		ctx,
		db,
		int32(migrator.WelcomeMessageOriginatorID),
	))

	require.Equal(t, installationAmount, getGatewayEnvelopesAmount(
		t,
		ctx,
		db,
		int32(migrator.InstallationOriginatorID),
	))
}

func checkGatewayEnvelopesAreUnique(t *testing.T, ctx context.Context, db *sql.DB) {
	require.Equal(t, groupMessageAmount, getGatewayEnvelopesUniqueAmount(
		t,
		ctx,
		db,
		int32(migrator.GroupMessageOriginatorID),
	))

	require.Equal(t, inboxLogAmount, getGatewayEnvelopesUniqueAmount(
		t,
		ctx,
		db,
		int32(migrator.InboxLogOriginatorID),
	))

	require.Equal(t, welcomeMessageAmount, getGatewayEnvelopesUniqueAmount(
		t,
		ctx,
		db,
		int32(migrator.WelcomeMessageOriginatorID),
	))

	require.Equal(t, installationAmount, getGatewayEnvelopesUniqueAmount(
		t,
		ctx,
		db,
		int32(migrator.InstallationOriginatorID),
	))
}

func checkAddressLogAmount(t *testing.T, ctx context.Context, db *sql.DB) {
	var (
		count int64
		query = `SELECT COUNT(*)::BIGINT FROM address_log`
	)

	row := db.QueryRowContext(ctx, query)
	require.NoError(t, row.Scan(&count))

	require.Equal(t, addressLogAmount, count)
}

func getGatewayEnvelopesLastSequenceID(
	t *testing.T,
	ctx context.Context,
	db *sql.DB,
	originatorNodeID int32,
) int64 {
	querier := queries.New(db)

	lastSequenceID, err := querier.GetLatestSequenceId(
		ctx,
		originatorNodeID,
	)
	require.NoError(t, err)

	return lastSequenceID
}

func getGatewayEnvelopesAmount(
	t *testing.T,
	ctx context.Context,
	db *sql.DB,
	originatorNodeID int32,
) int64 {
	var (
		count              int64
		getEnvelopesAmount = `SELECT COUNT(*)::BIGINT
FROM gateway_envelopes
WHERE originator_node_id = $1`
	)

	row := db.QueryRowContext(ctx, getEnvelopesAmount, originatorNodeID)
	require.NoError(t, row.Scan(&count))

	return count
}

func getGatewayEnvelopesUniqueAmount(
	t *testing.T,
	ctx context.Context,
	db *sql.DB,
	originatorNodeID int32,
) int64 {
	var (
		count              int64
		getEnvelopesAmount = `SELECT COUNT(DISTINCT originator_sequence_id)::BIGINT
FROM gateway_envelopes 
WHERE originator_node_id = $1`
	)

	row := db.QueryRowContext(ctx, getEnvelopesAmount, originatorNodeID)
	require.NoError(t, row.Scan(&count))

	return count
}
