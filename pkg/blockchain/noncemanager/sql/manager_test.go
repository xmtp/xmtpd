package sql_test

import (
	"container/heap"
	"context"
	"math/big"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/blockchain/noncemanager"
	"github.com/xmtp/xmtpd/pkg/blockchain/noncemanager/sql"
	"github.com/xmtp/xmtpd/pkg/testutils"
	"go.uber.org/zap"
)

type Int64Heap []int64

func (h *Int64Heap) Len() int           { return len(*h) }
func (h *Int64Heap) Less(i, j int) bool { return (*h)[i] < (*h)[j] }
func (h *Int64Heap) Swap(i, j int)      { (*h)[i], (*h)[j] = (*h)[j], (*h)[i] }

func (h *Int64Heap) Push(x any) {
	*h = append(*h, x.(int64))
}

func (h *Int64Heap) Pop() any {
	old := *h
	n := len(old)
	x := old[0] // Get the smallest element
	*h = old[1:n]
	return x
}

func (h *Int64Heap) Peek() int64 {
	if len(*h) == 0 {
		return -1 // Return an invalid value if empty
	}
	return (*h)[0]
}

type TestNonceManager struct {
	mu        sync.Mutex
	nonce     int64
	logger    *zap.Logger
	abandoned Int64Heap
}

func NewTestNonceManager(logger *zap.Logger) *TestNonceManager {
	return &TestNonceManager{logger: logger}
}

func (tm *TestNonceManager) GetNonce(ctx context.Context) (*noncemanager.NonceContext, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	tm.mu.Lock()
	defer tm.mu.Unlock()

	var nonce int64
	if tm.abandoned.Len() > 0 {
		nonce = heap.Pop(&tm.abandoned).(int64)
	} else {
		nonce = tm.nonce
		tm.nonce++
	}

	return &noncemanager.NonceContext{
		Nonce: *new(big.Int).SetInt64(nonce),
		Cancel: func() {
			tm.mu.Lock()
			defer tm.mu.Unlock()
			tm.abandoned.Push(nonce)
		}, // No-op
		Consume: func() error {
			return nil // No-op
		},
	}, nil
}

func (tm *TestNonceManager) FastForwardNonce(ctx context.Context, nonce big.Int) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	tm.nonce = nonce.Int64()

	return nil
}

func (tm *TestNonceManager) Replenish(ctx context.Context, nonce big.Int) error {
	return nil
}

func TestGetNonce_Simple(t *testing.T) {
	ctx := t.Context()
	db, _ := testutils.NewDB(t, ctx)

	logger, err := zap.NewDevelopment()
	require.NoError(t, err)

	nonceManager := sql.NewSQLBackedNonceManager(db, logger)
	err = nonceManager.Replenish(ctx, *big.NewInt(0))
	require.NoError(t, err)

	nonce, err := nonceManager.GetNonce(ctx)
	require.NoError(t, err)
	defer nonce.Cancel()

	require.EqualValues(t, 0, nonce.Nonce.Int64())
}

func TestGetNonce_RevertMany(t *testing.T) {
	ctx := t.Context()
	db, _ := testutils.NewDB(t, ctx)

	logger, err := zap.NewDevelopment()
	require.NoError(t, err)

	nonceManager := sql.NewSQLBackedNonceManager(db, logger)
	err = nonceManager.Replenish(ctx, *big.NewInt(0))
	require.NoError(t, err)

	for range 10 {
		nonce, err := nonceManager.GetNonce(ctx)
		require.NoError(t, err)
		require.EqualValues(t, 0, nonce.Nonce.Int64())
		nonce.Cancel()
	}
}

func TestGetNonce_ConsumeMany(t *testing.T) {
	ctx := t.Context()
	db, _ := testutils.NewDB(t, ctx)

	logger, err := zap.NewDevelopment()
	require.NoError(t, err)

	nonceManager := sql.NewSQLBackedNonceManager(db, logger)
	err = nonceManager.Replenish(ctx, *big.NewInt(0))
	require.NoError(t, err)

	for i := range 10 {
		nonce, err := nonceManager.GetNonce(ctx)
		require.NoError(t, err)
		require.EqualValues(t, i, nonce.Nonce.Int64())
		err = nonce.Consume()
		require.NoError(t, err)
	}
}

func TestGetNonce_ConsumeManyConcurrent(t *testing.T) {
	ctx := t.Context()
	db, _ := testutils.NewDB(t, ctx)

	logger, err := zap.NewDevelopment()
	require.NoError(t, err)

	nonceManager := sql.NewSQLBackedNonceManager(db, logger)
	err = nonceManager.Replenish(ctx, *big.NewInt(0))
	require.NoError(t, err)

	var wg sync.WaitGroup
	numClients := 20
	errCh := make(chan error, numClients)

	for range numClients {
		wg.Go(func() {
			nonce, err := nonceManager.GetNonce(ctx)
			if err != nil {
				errCh <- err
				return
			}
			err = nonce.Consume()
			if err != nil {
				errCh <- err
				return
			}
		})
	}

	wg.Wait()
	close(errCh)
	for err := range errCh {
		require.NoError(t, err)
	}
}
