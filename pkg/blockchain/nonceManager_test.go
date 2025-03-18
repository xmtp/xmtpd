package blockchain_test

import (
	"container/heap"
	"context"
	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/blockchain"
	"github.com/xmtp/xmtpd/pkg/testutils"
	"go.uber.org/zap"
	"math/big"
	"os"
	"strconv"
	"sync"
	"testing"
)

type Int64Heap []int64

func (h *Int64Heap) Len() int           { return len(*h) }
func (h *Int64Heap) Less(i, j int) bool { return (*h)[i] < (*h)[j] }
func (h *Int64Heap) Swap(i, j int)      { (*h)[i], (*h)[j] = (*h)[j], (*h)[i] }

func (h *Int64Heap) Push(x interface{}) {
	*h = append(*h, x.(int64))
}

func (h *Int64Heap) Pop() interface{} {
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
	limiter   *blockchain.OpenConnectionsLimiter
}

func NewTestNonceManager(logger *zap.Logger) *TestNonceManager {
	// Read from environment variable, default to 100 if not set or invalid
	limit := 100
	if envLimit, exists := os.LookupEnv("NONCE_MANAGER_LIMIT"); exists {
		if parsedLimit, err := strconv.Atoi(envLimit); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	return &TestNonceManager{logger: logger,
		limiter: blockchain.NewOpenConnectionsLimiter(limit)}
}

func (tm *TestNonceManager) GetNonce(ctx context.Context) (*blockchain.NonceContext, error) {

	tm.limiter.Wg.Add(1)
	select {
	case tm.limiter.Semaphore <- struct{}{}:
	case <-ctx.Done():
		tm.limiter.Wg.Done()
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

	tm.logger.Debug("Generated Nonce", zap.Int64("nonce", nonce))

	return &blockchain.NonceContext{
		Nonce: *new(big.Int).SetInt64(nonce),
		Cancel: func() {
			tm.mu.Lock()
			defer tm.mu.Unlock()
			tm.abandoned.Push(nonce)
			<-tm.limiter.Semaphore
			tm.limiter.Wg.Done()
		}, // No-op
		Consume: func() error {
			<-tm.limiter.Semaphore
			tm.limiter.Wg.Done()
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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	db, _, cleanup := testutils.NewDB(t, ctx)
	defer cleanup()

	logger, err := zap.NewDevelopment()
	require.NoError(t, err)

	nonceManager := blockchain.NewSQLBackedNonceManager(db, logger)
	err = nonceManager.Replenish(ctx, *big.NewInt(0))
	require.NoError(t, err)

	nonce, err := nonceManager.GetNonce(ctx)
	require.NoError(t, err)
	defer nonce.Cancel()

	require.EqualValues(t, 0, nonce.Nonce.Int64())
}

func TestGetNonce_RevertMany(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	db, _, cleanup := testutils.NewDB(t, ctx)
	defer cleanup()

	logger, err := zap.NewDevelopment()
	require.NoError(t, err)

	nonceManager := blockchain.NewSQLBackedNonceManager(db, logger)
	err = nonceManager.Replenish(ctx, *big.NewInt(0))
	require.NoError(t, err)

	for i := 0; i < 10; i++ {
		nonce, err := nonceManager.GetNonce(ctx)
		require.NoError(t, err)
		require.EqualValues(t, 0, nonce.Nonce.Int64())
		nonce.Cancel()
	}
}

func TestGetNonce_ConsumeMany(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	db, _, cleanup := testutils.NewDB(t, ctx)
	defer cleanup()

	logger, err := zap.NewDevelopment()
	require.NoError(t, err)

	nonceManager := blockchain.NewSQLBackedNonceManager(db, logger)
	err = nonceManager.Replenish(ctx, *big.NewInt(0))
	require.NoError(t, err)

	for i := 0; i < 10; i++ {
		nonce, err := nonceManager.GetNonce(ctx)
		require.NoError(t, err)
		require.EqualValues(t, i, nonce.Nonce.Int64())
		err = nonce.Consume()
		require.NoError(t, err)
	}
}

func TestGetNonce_ConsumeManyConcurrent(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	db, _, cleanup := testutils.NewDB(t, ctx)
	defer cleanup()

	logger, err := zap.NewDevelopment()
	require.NoError(t, err)

	nonceManager := blockchain.NewSQLBackedNonceManager(db, logger)
	err = nonceManager.Replenish(ctx, *big.NewInt(0))
	require.NoError(t, err)

	var wg sync.WaitGroup
	numClients := 20

	for i := 0; i < numClients; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			nonce, err := nonceManager.GetNonce(ctx)
			require.NoError(t, err)
			err = nonce.Consume()
			require.NoError(t, err)
		}()
	}

	wg.Wait()
}
