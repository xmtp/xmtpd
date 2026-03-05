package db_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/db"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/testutils"
)

func TestFindOrCreatePayer(t *testing.T) {
	ctx := context.Background()
	rawDB, _ := testutils.NewRawDB(t, ctx)

	querier := queries.New(rawDB)

	address1 := testutils.RandomString(42)
	address2 := testutils.RandomString(42)

	id1, err := querier.FindOrCreatePayer(ctx, address1)
	require.NoError(t, err)

	id2, err := querier.FindOrCreatePayer(ctx, address2)
	require.NoError(t, err)

	require.NotEqual(t, id1, id2)

	reinsertID, err := querier.FindOrCreatePayer(ctx, address1)
	require.NoError(t, err)
	require.Equal(t, id1, reinsertID)
}

func TestFindOrCreatePayerWithRetry(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		ctx, querier := setupTest(t)
		address := testutils.RandomString(42)

		// First call creates the payer
		id1, err := db.FindOrCreatePayerWithRetry(ctx, querier, address, 3)
		require.NoError(t, err)
		require.NotZero(t, id1)

		// Second call finds the existing payer
		id2, err := db.FindOrCreatePayerWithRetry(ctx, querier, address, 3)
		require.NoError(t, err)
		require.Equal(t, id1, id2)
	})

	t.Run("race condition", func(t *testing.T) {
		ctx := context.Background()
		rawDB, _ := testutils.NewRawDB(t, ctx)
		address := testutils.RandomString(42)

		// Start transaction T1 and insert the payer (holds row lock, uncommitted)
		tx1, err := rawDB.BeginTx(ctx, nil)
		require.NoError(t, err)
		defer func() { _ = tx1.Rollback() }()

		_, err = tx1.ExecContext(ctx, "INSERT INTO payers(address) VALUES ($1)", address)
		require.NoError(t, err)

		// Commit T1 after a short delay so the retry can succeed
		go func() {
			time.Sleep(5 * time.Millisecond)
			_ = tx1.Commit()
		}()

		// On a separate connection, the raw FindOrCreatePayer gets sql.ErrNoRows
		// because the CTE INSERT conflicts (T1 holds the lock) and the SELECT
		// uses the pre-commit snapshot.
		poolQuerier := queries.New(rawDB)

		// FindOrCreatePayerWithRetry should succeed after T1 commits
		id, err := db.FindOrCreatePayerWithRetry(ctx, poolQuerier, address, 3)
		require.NoError(t, err)
		require.NotZero(t, id)
	})

	t.Run("context cancellation stops retries", func(t *testing.T) {
		ctx := context.Background()
		rawDB, _ := testutils.NewRawDB(t, ctx)
		address := testutils.RandomString(42)

		// Start transaction T1 and insert the payer (holds row lock, never commits)
		tx1, err := rawDB.BeginTx(ctx, nil)
		require.NoError(t, err)
		defer func() { _ = tx1.Rollback() }()

		_, err = tx1.ExecContext(ctx, "INSERT INTO payers(address) VALUES ($1)", address)
		require.NoError(t, err)

		// Use a context that cancels quickly
		cancelCtx, cancel := context.WithTimeout(ctx, 10*time.Millisecond)
		defer cancel()

		poolQuerier := queries.New(rawDB)
		_, err = db.FindOrCreatePayerWithRetry(cancelCtx, poolQuerier, address, 100)
		assert.Error(t, err)
	})
}
