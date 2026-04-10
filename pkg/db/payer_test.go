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

		// Start transaction T1 and insert the payer (holds row lock, uncommitted).
		tx1, err := rawDB.BeginTx(ctx, nil)
		require.NoError(t, err)
		defer func() { _ = tx1.Rollback() }()

		_, err = tx1.ExecContext(ctx, "INSERT INTO payers(address) VALUES ($1)", address)
		require.NoError(t, err)

		poolQuerier := queries.New(rawDB)

		// Run FindOrCreatePayerWithRetry concurrently with T1 commit. The
		// retry call will block inside INSERT ... ON CONFLICT on T1's unique-
		// index lock until T1 resolves. Instead of sleeping a fixed duration
		// to guarantee the retry is blocked before we commit, poll
		// pg_stat_activity for a session waiting on a Lock event — a
		// deterministic signal that the contending query is actually stalled.
		type result struct {
			id  int32
			err error
		}
		resCh := make(chan result, 1)
		go func() {
			id, err := db.FindOrCreatePayerWithRetry(ctx, poolQuerier, address, 10)
			resCh <- result{id: id, err: err}
		}()

		require.Eventually(t, func() bool {
			var count int
			err := rawDB.QueryRowContext(ctx, `
				SELECT count(*)
				FROM pg_stat_activity
				WHERE state = 'active'
				  AND wait_event_type = 'Lock'
				  AND query ILIKE '%payers%'
			`).Scan(&count)
			return err == nil && count >= 1
		}, 5*time.Second, 10*time.Millisecond, "retry call should be blocked on T1's lock")

		require.NoError(t, tx1.Commit())

		select {
		case r := <-resCh:
			require.NoError(t, r.err)
			require.NotZero(t, r.id)
		case <-time.After(5 * time.Second):
			t.Fatal("FindOrCreatePayerWithRetry did not return after T1 commit")
		}
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
