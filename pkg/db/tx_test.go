package db_test

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/db"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/testutils"
)

func TestRunInTxRaw(t *testing.T) {
	var (
		sqlDB, _ = testutils.NewRawDB(t, t.Context())
		// Callback invoked in a transaction that succeeds.
		callbackSuccess = func(ctx context.Context, tx *sql.Tx) error {
			return nil
		}

		// Callback invoked in a transaction that fails, prompting a transaction rollback.
		callbackFailure = func(ctx context.Context, tx *sql.Tx) error {
			return errors.New("failing tx callback")
		}
	)

	t.Run("simple transaction", func(t *testing.T) {
		err := db.RunInTxRaw(t.Context(), sqlDB, nil, callbackSuccess)
		require.NoError(t, err)
	})
	t.Run("on commit hook", func(t *testing.T) {
		var (
			onCommitCalled = false
			onCommitHook   = func() {
				onCommitCalled = true
			}
		)

		err := db.RunInTxRaw(t.Context(), sqlDB, nil, callbackSuccess, db.OnCommit(onCommitHook))
		require.NoError(t, err)

		require.True(t, onCommitCalled)
	})
	t.Run("multiple on commit hooks", func(t *testing.T) {
		var (
			onCommitCounter = 0
			onCommitHook    = func() {
				onCommitCounter += 1
			}
		)

		err := db.RunInTxRaw(t.Context(), sqlDB, nil, callbackSuccess,
			db.OnCommit(onCommitHook),
			db.OnCommit(onCommitHook),
			db.OnCommit(onCommitHook),
		)
		require.NoError(t, err)

		require.Equal(t, 3, onCommitCounter)
	})
	t.Run("on rollback hook", func(t *testing.T) {
		var (
			onRollbackCalled = false
			onRollbackHook   = func() {
				onRollbackCalled = true
			}
		)

		err := db.RunInTxRaw(
			t.Context(),
			sqlDB,
			nil,
			callbackFailure,
			db.OnRollback(onRollbackHook),
		)
		require.Error(t, err)

		require.True(t, onRollbackCalled)
	})
	t.Run("multiple on rollback hooks", func(t *testing.T) {
		var (
			onRollbackCounter = 0
			onRollbackHook    = func() {
				onRollbackCounter += 1
			}
		)

		err := db.RunInTxRaw(t.Context(), sqlDB, nil, callbackFailure,
			db.OnRollback(onRollbackHook),
			db.OnRollback(onRollbackHook),
			db.OnRollback(onRollbackHook),
		)
		require.Error(t, err)

		require.Equal(t, 3, onRollbackCounter)
	})
	t.Run("commit and rollback hooks - transaction success", func(t *testing.T) {
		var (
			onCommitCalled   = false
			onRollbackCalled = false

			onCommitHook = func() {
				onCommitCalled = true
			}
			onRollbackHook = func() {
				onRollbackCalled = true
			}
		)

		err := db.RunInTxRaw(t.Context(), sqlDB, nil, callbackSuccess,
			db.OnCommit(onCommitHook),
			db.OnRollback(onRollbackHook),
		)
		require.NoError(t, err)

		require.True(t, onCommitCalled)
		require.False(t, onRollbackCalled)
	})
	t.Run("commit and rollback hooks - transaction rolled back", func(t *testing.T) {
		var (
			onCommitCalled   = false
			onRollbackCalled = false

			onCommitHook = func() {
				onCommitCalled = true
			}
			onRollbackHook = func() {
				onRollbackCalled = true
			}
		)

		err := db.RunInTxRaw(t.Context(), sqlDB, nil, callbackFailure,
			db.OnCommit(onCommitHook),
			db.OnRollback(onRollbackHook),
		)
		require.Error(t, err)

		require.False(t, onCommitCalled)
		require.True(t, onRollbackCalled)
	})
}

func TestRunInTxWithResult(t *testing.T) {
	// Dummy DB record to be returned by the function passed to the transaction.
	type dummyRecord struct {
		str string
	}

	var (
		sqlDB, _ = testutils.NewRawDB(t, t.Context())
		// generator function to generate per-test case payload.
		generate = func() []dummyRecord {
			// Generate 5-10 records
			count := 5 + rand.Intn(5)

			out := make([]dummyRecord, count)
			for i := range count {
				out[i] = dummyRecord{
					str: fmt.Sprint(rand.Int()),
				}
			}
			return out
		}
	)

	t.Run("simple transaction", func(t *testing.T) {
		var (
			input = generate()
			txfn  = func(context.Context, *queries.Queries) ([]dummyRecord, error) {
				return input, nil
			}
		)

		res, err := db.RunInTxWithResult(t.Context(), sqlDB, nil, txfn)
		require.NoError(t, err)
		require.Equal(t, input, res)
	})
	t.Run("on commit hook", func(t *testing.T) {
		var (
			input = generate()
			txfn  = func(context.Context, *queries.Queries) ([]dummyRecord, error) {
				return input, nil
			}
		)

		var (
			onCommitCalled = false
			onCommitHook   = func() {
				onCommitCalled = true
			}
		)

		res, err := db.RunInTxWithResult(t.Context(), sqlDB, nil, txfn,
			db.OnCommit(onCommitHook),
		)
		require.NoError(t, err)
		require.Equal(t, input, res)
		require.True(t, onCommitCalled)
	})
	t.Run("multiple on commit hooks", func(t *testing.T) {
		var (
			input = generate()
			txfn  = func(context.Context, *queries.Queries) ([]dummyRecord, error) {
				return input, nil
			}
		)

		var (
			onCommitCounter = 0
			onCommitHook    = func() {
				onCommitCounter += 1
			}
		)

		res, err := db.RunInTxWithResult(t.Context(), sqlDB, nil, txfn,
			db.OnCommit(onCommitHook),
			db.OnCommit(onCommitHook),
			db.OnCommit(onCommitHook),
		)
		require.NoError(t, err)
		require.Equal(t, input, res)
		require.Equal(t, 3, onCommitCounter)
	})
	t.Run("on rollback hook", func(t *testing.T) {
		var (
			input = generate()
			txfn  = func(context.Context, *queries.Queries) ([]dummyRecord, error) {
				return input, errors.New("transaction function failed")
			}
		)

		var (
			onRollbackCalled = false
			onRollbackHook   = func() {
				onRollbackCalled = true
			}
		)

		_, err := db.RunInTxWithResult(t.Context(), sqlDB, nil, txfn,
			db.OnRollback(onRollbackHook),
		)
		require.Error(t, err)
		require.True(t, onRollbackCalled)
	})
	t.Run("multiple on rollback hooks", func(t *testing.T) {
		var (
			input = generate()
			txfn  = func(context.Context, *queries.Queries) ([]dummyRecord, error) {
				return input, errors.New("transaction function failed")
			}
		)

		var (
			onRollbackCounter = 0
			onRollbackHook    = func() {
				onRollbackCounter += 1
			}
		)

		_, err := db.RunInTxWithResult(t.Context(), sqlDB, nil, txfn,
			db.OnRollback(onRollbackHook),
			db.OnRollback(onRollbackHook),
			db.OnRollback(onRollbackHook),
		)
		require.Error(t, err)
		require.Equal(t, 3, onRollbackCounter)
	})
	t.Run("commit and rollback hooks - transaction success", func(t *testing.T) {
		var (
			input = generate()
			txfn  = func(context.Context, *queries.Queries) ([]dummyRecord, error) {
				return input, nil
			}
		)

		var (
			onCommitCalled   = false
			onRollbackCalled = false

			onCommitHook = func() {
				onCommitCalled = true
			}
			onRollbackHook = func() {
				onRollbackCalled = true
			}
		)

		res, err := db.RunInTxWithResult(t.Context(), sqlDB, nil, txfn,
			db.OnCommit(onCommitHook),
			db.OnRollback(onRollbackHook),
		)
		require.NoError(t, err)
		require.Equal(t, input, res)
		require.True(t, onCommitCalled)

		require.True(t, onCommitCalled)
		require.False(t, onRollbackCalled)
	})
	t.Run("commit and rollback hooks - transaction rolled back", func(t *testing.T) {
		var (
			input = generate()
			txfn  = func(context.Context, *queries.Queries) ([]dummyRecord, error) {
				return input, errors.New("transaction function failed")
			}
		)

		var (
			onCommitCalled   = false
			onRollbackCalled = false

			onCommitHook = func() {
				onCommitCalled = true
			}
			onRollbackHook = func() {
				onRollbackCalled = true
			}
		)

		_, err := db.RunInTxWithResult(t.Context(), sqlDB, nil, txfn,
			db.OnCommit(onCommitHook),
			db.OnRollback(onRollbackHook),
		)
		require.Error(t, err)
		require.False(t, onCommitCalled)
		require.True(t, onRollbackCalled)
	})
}
