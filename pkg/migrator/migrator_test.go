package migrator_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/config"
	"github.com/xmtp/xmtpd/pkg/migrator"
	"github.com/xmtp/xmtpd/pkg/migrator/testdata"
	"github.com/xmtp/xmtpd/pkg/mocks/mlsvalidate"
	"github.com/xmtp/xmtpd/pkg/testutils"
	"github.com/xmtp/xmtpd/pkg/utils"
)

type migratorTest struct {
	ctx      context.Context
	cleanup  func()
	migrator *migrator.Migrator
}

func newMigratorTest(t *testing.T) *migratorTest {
	var (
		ctx                  = t.Context()
		writerDB, _          = testutils.NewDB(t, ctx)
		_, dsn, cleanup      = testdata.NewMigratorTestDB(t, ctx)
		mlsValidationService = mlsvalidate.NewMockMLSValidationService(t)
		payerPrivateKey      = testutils.RandomPrivateKey(t)
		nodePrivateKey       = testutils.RandomPrivateKey(t)
	)

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
			ReadTimeout:            10 * time.Second,
			WaitForDB:              30 * time.Second,
			BatchSize:              1000,
			PollInterval:           10 * time.Second,
		}),
	)
	require.NoError(t, err)

	return &migratorTest{
		ctx:      ctx,
		cleanup:  cleanup,
		migrator: migrator,
	}
}

func TestMigrator(t *testing.T) {
	test := newMigratorTest(t)
	defer test.cleanup()

	err := test.migrator.Start()
	require.NoError(t, err)

	time.Sleep(100 * time.Second)

	err = test.migrator.Stop()
	require.NoError(t, err)
}
