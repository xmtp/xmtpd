package blockchain_test

import (
	"context"
	"math/big"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/blockchain"
	"github.com/xmtp/xmtpd/pkg/testutils"
	"go.uber.org/zap"
)

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
