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
	"sync/atomic"
	"testing"
	"time"
)

func TestFillRows(t *testing.T) {
	ctx := context.Background()
	db, _, cleanup := testutils.NewDB(t, ctx)
	defer cleanup()

	querier := queries.New(db)

	err := querier.FillNonceSequence(ctx, queries.FillNonceSequenceParams{
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
	_, err := querier.GetNextAvailableNonce(ctx)
	require.Error(t, err)

	err = querier.FillNonceSequence(ctx, queries.FillNonceSequenceParams{
		PendingNonce: 0,
		NumElements:  100,
	})
	require.NoError(t, err)

	seq, err := querier.GetNextAvailableNonce(ctx)
	require.NoError(t, err)
	require.EqualValues(t, 0, seq)
}

func getNextPayerSequence(t *testing.T, ctx context.Context, db *sql.DB) (int64, error) {
	return db2.RunInTxWithResult(ctx, db, &sql.TxOptions{},
		func(ctx context.Context, querier *queries.Queries) (int64, error) {
			seq, err := querier.GetNextAvailableNonce(ctx)
			if err != nil {
				return 0, err
			}
			t.Log("Acquired sequence ID: ", seq)
			time.Sleep(10 * time.Millisecond)

			_, err = querier.DeleteAvailableNonce(ctx, seq)
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
			seq, err := querier.GetNextAvailableNonce(ctx)
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

	err := querier.FillNonceSequence(ctx, queries.FillNonceSequenceParams{
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

	err := querier.FillNonceSequence(ctx, queries.FillNonceSequenceParams{
		PendingNonce: 0,
		NumElements:  100,
	})
	require.NoError(t, err)

	seq, err := querier.GetNextAvailableNonce(ctx)
	require.NoError(t, err)
	require.EqualValues(t, 0, seq)

	seq, err = querier.GetNextAvailableNonce(ctx)
	require.NoError(t, err)
	require.EqualValues(t, 0, seq)

	seq, err = querier.GetNextAvailableNonce(ctx)
	require.NoError(t, err)
	require.EqualValues(t, 0, seq)

}
func TestRequestsUsed(t *testing.T) {
	ctx := context.Background()
	db, _, cleanup := testutils.NewDB(t, ctx)
	defer cleanup()

	querier := queries.New(db)

	err := querier.FillNonceSequence(ctx, queries.FillNonceSequenceParams{
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

	err := querier.FillNonceSequence(ctx, queries.FillNonceSequenceParams{
		PendingNonce: 0,
		NumElements:  100,
	})
	require.NoError(t, err)

	seq, err := getNextPayerSequence(t, ctx, db)
	require.NoError(t, err)
	require.EqualValues(t, 0, seq)

	_, err = failNextPayerSequence(t, ctx, db)
	require.Error(t, err)

	_, err = failNextPayerSequence(t, ctx, db)
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

	err := querier.FillNonceSequence(ctx, queries.FillNonceSequenceParams{
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

	_, err = txQuerier.GetNextAvailableNonce(ctx)
	require.NoError(t, err)
	_, err = txQuerier.GetNextAvailableNonce(ctx)
	require.NoError(t, err)
	_, err = txQuerier.GetNextAvailableNonce(ctx)
	require.NoError(t, err)

	err = querier.FillNonceSequence(ctx, queries.FillNonceSequenceParams{
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

	err := querier.FillNonceSequence(ctx, queries.FillNonceSequenceParams{
		PendingNonce: 0,
		NumElements:  10,
	})
	require.NoError(t, err)

	err = querier.FillNonceSequence(ctx, queries.FillNonceSequenceParams{
		PendingNonce: 0,
		NumElements:  30,
	})
	require.NoError(t, err)
}

func TestAbandonNonces(t *testing.T) {
	ctx := context.Background()
	db, _, cleanup := testutils.NewDB(t, ctx)
	defer cleanup()

	querier := queries.New(db)

	err := querier.FillNonceSequence(ctx, queries.FillNonceSequenceParams{
		PendingNonce: 0,
		NumElements:  10,
	})
	require.NoError(t, err)

	_, err = querier.DeleteObsoleteNonces(ctx, 5)
	require.NoError(t, err)

	nonce, err := querier.GetNextAvailableNonce(ctx)
	require.NoError(t, err)

	require.EqualValues(t, 5, nonce)

}

func TestAbandonCanProceedWithOpenTxn(t *testing.T) {
	ctx := context.Background()
	db, _, cleanup := testutils.NewDB(t, ctx)
	defer cleanup()

	querier := queries.New(db)

	err := querier.FillNonceSequence(ctx, queries.FillNonceSequenceParams{
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

	_, err = txQuerier.GetNextAvailableNonce(ctx)
	require.NoError(t, err)
	_, err = txQuerier.GetNextAvailableNonce(ctx)
	require.NoError(t, err)
	_, err = txQuerier.GetNextAvailableNonce(ctx)
	require.NoError(t, err)

	_, err = querier.DeleteObsoleteNonces(ctx, 5)
	require.NoError(t, err)

	nonce, err := querier.GetNextAvailableNonce(ctx)
	require.NoError(t, err)

	require.EqualValues(t, 5, nonce)

}

func TestAbandonSkipsOpenTxn(t *testing.T) {
	ctx := context.Background()
	db, _, cleanup := testutils.NewDB(t, ctx)
	defer cleanup()

	querier := queries.New(db)

	err := querier.FillNonceSequence(ctx, queries.FillNonceSequenceParams{
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

	_, err = txQuerier.GetNextAvailableNonce(ctx)
	require.NoError(t, err)
	_, err = txQuerier.GetNextAvailableNonce(ctx)
	require.NoError(t, err)
	_, err = txQuerier.GetNextAvailableNonce(ctx)
	require.NoError(t, err)

	_, err = querier.DeleteObsoleteNonces(ctx, 5)
	require.NoError(t, err)

	err = tx.Rollback()
	require.NoError(t, err)

	// the nonce manager has at least once semantics
	// it might return a nonce that is too low
	nonce, err := querier.GetNextAvailableNonce(ctx)
	require.NoError(t, err)

	require.EqualValues(t, 0, nonce)

}

func TestAbandonConcurrently(t *testing.T) {
	ctx := context.Background()
	db, _, cleanup := testutils.NewDB(t, ctx)
	defer cleanup()

	querier := queries.New(db)

	err := querier.FillNonceSequence(ctx, queries.FillNonceSequenceParams{
		PendingNonce: 0,
		NumElements:  500,
	})
	require.NoError(t, err)

	var wg sync.WaitGroup
	numClients := 20
	numDeletions := int64(0)

	for i := 1; i <= numClients; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			numrows, err := querier.DeleteObsoleteNonces(ctx, int64(10*i))
			if err != nil {
				t.Errorf("Error deleting nonces: %v", err)
			}
			atomic.AddInt64(&numDeletions, numrows)
		}()
	}

	// Wait for all goroutines to complete
	wg.Wait()
	require.EqualValues(t, 200, numDeletions)
}

func TestAbandonConcurrentlyWithOpenTransaction(t *testing.T) {
	ctx := context.Background()
	db, _, cleanup := testutils.NewDB(t, ctx)
	defer cleanup()

	querier := queries.New(db)

	err := querier.FillNonceSequence(ctx, queries.FillNonceSequenceParams{
		PendingNonce: 0,
		NumElements:  500,
	})
	require.NoError(t, err)

	var wg sync.WaitGroup
	numClients := 20
	numDeletions := int64(0)

	tx, err := db.BeginTx(ctx, &sql.TxOptions{})
	require.NoError(t, err)

	defer func() {
		_ = tx.Rollback()
	}()

	// hold this TX open
	txQuerier := queries.New(db).WithTx(tx)

	_, err = txQuerier.GetNextAvailableNonce(ctx)
	require.NoError(t, err)

	for i := 1; i <= numClients; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			numrows, err := querier.DeleteObsoleteNonces(ctx, int64(10*i))
			if err != nil {
				t.Errorf("Error deleting nonces: %v", err)
			}
			atomic.AddInt64(&numDeletions, numrows)
		}()
	}

	wg.Wait()
	require.EqualValues(t, 199, numDeletions)
}
