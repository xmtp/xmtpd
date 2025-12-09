package db_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/testutils"
)

func TestAdvisoryTryLockWithKey(t *testing.T) {
	ctx := context.Background()
	db, _ := testutils.NewDB(t, ctx)

	// NOTE: We need two transactions in order to compete for the lock.
	tx1, err := db.BeginTx(ctx, nil)
	require.NoError(t, err)
	defer func() {
		_ = tx1.Rollback()
	}()

	tx2, err := db.BeginTx(ctx, nil)
	require.NoError(t, err)
	defer func() {
		_ = tx2.Rollback()
	}()

	// Use a random value for the key that's unlikely to interfere with other tests that may involve locking.
	key := testutils.RandomInt64()

	// Lock should succeed.
	locked, err := queries.New(tx1).AdvisoryTryLockWithKey(ctx, key)
	require.NoError(t, err)
	require.True(t, locked)

	// Another transaction already hold the lock, so this should return but without a lock.
	locked, err = queries.New(tx2).AdvisoryTryLockWithKey(ctx, key)
	require.NoError(t, err)
	require.False(t, locked)

	// Perhaps codifying bad practice but - trying to get a lock while already owning it should work.
	locked, err = queries.New(tx1).AdvisoryTryLockWithKey(ctx, key)
	require.NoError(t, err)
	require.True(t, locked)
}
