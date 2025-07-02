package migrator_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/config"
	"github.com/xmtp/xmtpd/pkg/migrator"
	"github.com/xmtp/xmtpd/pkg/migrator/testdata"
	"github.com/xmtp/xmtpd/pkg/mlsvalidate"
	mlsvalidateMock "github.com/xmtp/xmtpd/pkg/mocks/mlsvalidate"
	"github.com/xmtp/xmtpd/pkg/proto/identity/associations"
	"github.com/xmtp/xmtpd/pkg/testutils"
	"github.com/xmtp/xmtpd/pkg/utils"
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
			PollInterval:           1 * time.Second,
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
	//defer test.cleanup()

	err := test.migrator.Start()
	require.NoError(t, err)

	time.Sleep(5 * time.Second)

	checkMigrationState(t, test.db)

	err = test.migrator.Stop()
	require.NoError(t, err)
}

func checkMigrationState(t *testing.T, db *sql.DB) {
	rows, err := db.QueryContext(t.Context(), "SELECT * FROM migration_tracker")
	require.NoError(t, err)

	defer rows.Close()

	for rows.Next() {
		var (
			tableName      string
			lastMigratedID int64
			createdAt      time.Time
			updatedAt      time.Time
		)

		err := rows.Scan(&tableName, &lastMigratedID, &createdAt, &updatedAt)
		require.NoError(t, err)

		t.Logf("tableName: %s, lastMigratedID: %d", tableName, lastMigratedID)
	}

	require.NoError(t, rows.Err())
}
