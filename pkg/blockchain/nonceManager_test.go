package blockchain_test

import (
	"context"
	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/blockchain"
	"github.com/xmtp/xmtpd/pkg/testutils"
	"go.uber.org/zap"
	"sync"
	"testing"
)

type TestNonceManager struct {
	mu     sync.Mutex
	nonce  uint64
	logger *zap.Logger
}

func NewTestNonceManager(logger *zap.Logger) *TestNonceManager {
	return &TestNonceManager{logger: logger}
}

func (tm *TestNonceManager) GetNonce(ctx context.Context) (*blockchain.NonceContext, error) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	nonce := tm.nonce
	tm.nonce++

	tm.logger.Debug("Generated Nonce", zap.Uint64("nonce", nonce))

	return &blockchain.NonceContext{
		Nonce:  nonce,
		Cancel: func() {}, // No-op
		Consume: func() error {
			return nil // No-op
		},
	}, nil
}

func (tm *TestNonceManager) FastForwardNonce(ctx context.Context, nonce uint64) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	tm.nonce = nonce

	return nil
}

func (tm *TestNonceManager) Replenish(ctx context.Context, nonce uint64) error {
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
	err = nonceManager.FillNonces(ctx, 0)
	require.NoError(t, err)

	nonce, err := nonceManager.GetNonce(ctx)
	require.NoError(t, err)
	defer nonce.Cancel()

	require.Equal(t, nonce.Nonce, uint64(0))
}

func TestGetNonce_RevertMany(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	db, _, cleanup := testutils.NewDB(t, ctx)
	defer cleanup()

	logger, err := zap.NewDevelopment()
	require.NoError(t, err)

	nonceManager := blockchain.NewSQLBackedNonceManager(db, logger)
	err = nonceManager.FillNonces(ctx, 0)
	require.NoError(t, err)

	for i := 0; i < 10; i++ {
		nonce, err := nonceManager.GetNonce(ctx)
		require.NoError(t, err)
		require.EqualValues(t, 0, nonce.Nonce)
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
	err = nonceManager.FillNonces(ctx, 0)
	require.NoError(t, err)

	for i := 0; i < 10; i++ {
		nonce, err := nonceManager.GetNonce(ctx)
		require.NoError(t, err)
		require.EqualValues(t, i, nonce.Nonce)
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
	err = nonceManager.FillNonces(ctx, 0)
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
