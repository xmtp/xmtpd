package db_test

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/stretchr/testify/require"
	db2 "github.com/xmtp/xmtpd/pkg/db"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/testutils"
	"sync"
	"testing"
	"time"
)

func TestFillRows(t *testing.T) {
	ctx := context.Background()
	db, _, cleanup := testutils.NewDB(t, ctx)
	defer cleanup()

	querier := queries.New(db)

	err := querier.FillPayerSequence(ctx, queries.FillPayerSequenceParams{
		PendingNonce: 0,
		NumElements:  100,
	})
	require.NoError(t, err)
}

func TestEmptyRows(t *testing.T) {
	ctx := context.Background()
	db, _, cleanup := testutils.NewDB(t, ctx)
	defer cleanup()

	querier := queries.New(db)
	seq, err := querier.GetNextAvailablePayerSequence(ctx)
	require.Error(t, err)

	err = querier.FillPayerSequence(ctx, queries.FillPayerSequenceParams{
		PendingNonce: 0,
		NumElements:  100,
	})
	require.NoError(t, err)

	seq, err = querier.GetNextAvailablePayerSequence(ctx)
	require.NoError(t, err)
	require.EqualValues(t, 0, seq)
}

func getNextPayerSequence(t *testing.T, ctx context.Context, db *sql.DB) (int64, error) {
	return db2.RunInTxWithResult(ctx, db, &sql.TxOptions{},
		func(ctx context.Context, querier *queries.Queries) (int64, error) {
			seq, err := querier.GetNextAvailablePayerSequence(ctx)
			if err != nil {
				return 0, err
			}
			t.Log("Acquired sequence ID: ", seq)
			time.Sleep(10 * time.Millisecond)

			_, err = querier.DeleteAvailablePayerSequence(ctx, seq)
			if err != nil {
				return 0, err
			}

			return int64(seq), nil
		},
	)
}

func failNextPayerSequence(t *testing.T, ctx context.Context, db *sql.DB) (int64, error) {
	return db2.RunInTxWithResult(ctx, db, &sql.TxOptions{},
		func(ctx context.Context, querier *queries.Queries) (int64, error) {
			seq, err := querier.GetNextAvailablePayerSequence(ctx)
			if err != nil {
				return 0, err
			}
			t.Log("Acquired sequence ID: ", seq)
			time.Sleep(10 * time.Millisecond)

			return 0, fmt.Errorf("failed to acquire sequence")
		},
	)
}

func TestConcurrentReads(t *testing.T) {
	ctx := context.Background()
	db, _, cleanup := testutils.NewDB(t, ctx)
	defer cleanup()

	querier := queries.New(db)

	err := querier.FillPayerSequence(ctx, queries.FillPayerSequenceParams{
		PendingNonce: 0,
		NumElements:  100,
	})
	require.NoError(t, err)

	var wg sync.WaitGroup
	numClients := 20
	results := make(chan int64, numClients)

	for i := 0; i < numClients; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			seqID, err := getNextPayerSequence(t, ctx, db)
			if err != nil {
				t.Errorf("Error acquiring sequence: %v", err)
			} else {
				results <- seqID
			}
		}()
	}

	// Wait for all goroutines to complete
	wg.Wait()
	close(results)

	// Collect all allocated sequence IDs
	allocatedIDs := make(map[int64]bool)

	for seqID := range results {
		if allocatedIDs[seqID] {
			t.Errorf("Duplicate sequence ID detected: %d", seqID)
		}
		allocatedIDs[seqID] = true
	}

	for i := 0; i < numClients; i++ {
		if !allocatedIDs[int64(i)] {
			t.Errorf("Missing sequence ID: %d", i)
		}
	}

	t.Log("✅ All 20 sequences were allocated uniquely with no gaps!")
}

func TestRequestsUnused(t *testing.T) {
	ctx := context.Background()
	db, _, cleanup := testutils.NewDB(t, ctx)
	defer cleanup()

	querier := queries.New(db)

	err := querier.FillPayerSequence(ctx, queries.FillPayerSequenceParams{
		PendingNonce: 0,
		NumElements:  100,
	})
	require.NoError(t, err)

	seq, err := querier.GetNextAvailablePayerSequence(ctx)
	require.NoError(t, err)
	require.EqualValues(t, 0, seq)

	seq, err = querier.GetNextAvailablePayerSequence(ctx)
	require.NoError(t, err)
	require.EqualValues(t, 0, seq)

	seq, err = querier.GetNextAvailablePayerSequence(ctx)
	require.NoError(t, err)
	require.EqualValues(t, 0, seq)

}
func TestRequestsUsed(t *testing.T) {
	ctx := context.Background()
	db, _, cleanup := testutils.NewDB(t, ctx)
	defer cleanup()

	querier := queries.New(db)

	err := querier.FillPayerSequence(ctx, queries.FillPayerSequenceParams{
		PendingNonce: 0,
		NumElements:  100,
	})
	require.NoError(t, err)

	seq, err := getNextPayerSequence(t, ctx, db)
	require.NoError(t, err)
	require.EqualValues(t, 0, seq)

	seq, err = getNextPayerSequence(t, ctx, db)
	require.NoError(t, err)
	require.EqualValues(t, 1, seq)

	seq, err = getNextPayerSequence(t, ctx, db)
	require.NoError(t, err)
	require.EqualValues(t, 2, seq)

}

func TestRequestsFailed(t *testing.T) {
	ctx := context.Background()
	db, _, cleanup := testutils.NewDB(t, ctx)
	defer cleanup()

	querier := queries.New(db)

	err := querier.FillPayerSequence(ctx, queries.FillPayerSequenceParams{
		PendingNonce: 0,
		NumElements:  100,
	})
	require.NoError(t, err)

	seq, err := getNextPayerSequence(t, ctx, db)
	require.NoError(t, err)
	require.EqualValues(t, 0, seq)

	seq, err = failNextPayerSequence(t, ctx, db)
	require.Error(t, err)

	seq, err = failNextPayerSequence(t, ctx, db)
	require.Error(t, err)

	seq, err = getNextPayerSequence(t, ctx, db)
	require.NoError(t, err)
	require.EqualValues(t, 1, seq)

}

func TestFillerCanProceedWithOpenTxn(t *testing.T) {
	ctx := context.Background()
	db, _, cleanup := testutils.NewDB(t, ctx)
	defer cleanup()

	querier := queries.New(db)

	err := querier.FillPayerSequence(ctx, queries.FillPayerSequenceParams{
		PendingNonce: 0,
		NumElements:  10,
	})
	require.NoError(t, err)

	tx, err := db.BeginTx(ctx, &sql.TxOptions{})
	require.NoError(t, err)

	defer func() {
		_ = tx.Rollback()
	}()

	// hold this TX open
	txQuerier := queries.New(db).WithTx(tx)

	_, err = txQuerier.GetNextAvailablePayerSequence(ctx)
	require.NoError(t, err)
	_, err = txQuerier.GetNextAvailablePayerSequence(ctx)
	require.NoError(t, err)
	_, err = txQuerier.GetNextAvailablePayerSequence(ctx)
	require.NoError(t, err)

	err = querier.FillPayerSequence(ctx, queries.FillPayerSequenceParams{
		PendingNonce: 0,
		NumElements:  30,
	})
	require.NoError(t, err)

}

func TestFillerRerun(t *testing.T) {
	ctx := context.Background()
	db, _, cleanup := testutils.NewDB(t, ctx)
	defer cleanup()

	querier := queries.New(db)

	err := querier.FillPayerSequence(ctx, queries.FillPayerSequenceParams{
		PendingNonce: 0,
		NumElements:  10,
	})
	require.NoError(t, err)

	err = querier.FillPayerSequence(ctx, queries.FillPayerSequenceParams{
		PendingNonce: 0,
		NumElements:  30,
	})
	require.NoError(t, err)
}
