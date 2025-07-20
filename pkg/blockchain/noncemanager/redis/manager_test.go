package redis_test

import (
	"context"
	"math/big"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	redismanager "github.com/xmtp/xmtpd/pkg/blockchain/noncemanager/redis"
	redistestutils "github.com/xmtp/xmtpd/pkg/testutils/redis"
	"go.uber.org/zap"
)

func TestRedisGetNonce_Simple(t *testing.T) {
	client, keyPrefix := redistestutils.NewRedisForTest(t)

	logger, err := zap.NewDevelopment()
	require.NoError(t, err)

	nonceManager := redismanager.NewRedisBackedNonceManager(client, logger, keyPrefix)
	err = nonceManager.Replenish(t.Context(), *big.NewInt(0))
	require.NoError(t, err)

	nonce, err := nonceManager.GetNonce(t.Context())
	require.NoError(t, err)
	defer nonce.Cancel()

	require.EqualValues(t, 0, nonce.Nonce.Int64())
}

func TestRedisGetNonce_RevertMany(t *testing.T) {
	client, keyPrefix := redistestutils.NewRedisForTest(t)

	logger, err := zap.NewDevelopment()
	require.NoError(t, err)

	nonceManager := redismanager.NewRedisBackedNonceManager(client, logger, keyPrefix)
	err = nonceManager.Replenish(t.Context(), *big.NewInt(0))
	require.NoError(t, err)

	for i := 0; i < 10; i++ {
		nonce, err := nonceManager.GetNonce(t.Context())
		require.NoError(t, err)
		require.EqualValues(t, 0, nonce.Nonce.Int64())
		nonce.Cancel()

		// Add a small delay to ensure the Cancel operation completes
		time.Sleep(1 * time.Millisecond)
	}
}

func TestRedisGetNonce_ConsumeMany(t *testing.T) {
	client, keyPrefix := redistestutils.NewRedisForTest(t)

	logger, err := zap.NewDevelopment()
	require.NoError(t, err)

	nonceManager := redismanager.NewRedisBackedNonceManager(client, logger, keyPrefix)
	err = nonceManager.Replenish(t.Context(), *big.NewInt(0))
	require.NoError(t, err)

	for i := 0; i < 10; i++ {
		nonce, err := nonceManager.GetNonce(t.Context())
		require.NoError(t, err)
		require.EqualValues(t, i, nonce.Nonce.Int64())
		err = nonce.Consume()
		require.NoError(t, err)
	}
}

func TestRedisGetNonce_ConsumeManyConcurrent(t *testing.T) {
	client, keyPrefix := redistestutils.NewRedisForTest(t)

	logger, err := zap.NewDevelopment()
	require.NoError(t, err)

	nonceManager := redismanager.NewRedisBackedNonceManager(client, logger, keyPrefix)
	err = nonceManager.Replenish(t.Context(), *big.NewInt(0))
	require.NoError(t, err)

	var wg sync.WaitGroup
	numClients := 20
	errCh := make(chan error, numClients)

	for i := 0; i < numClients; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			nonce, err := nonceManager.GetNonce(t.Context())
			if err != nil {
				errCh <- err
				return
			}
			err = nonce.Consume()
			if err != nil {
				errCh <- err
				return
			}
		}()
	}

	wg.Wait()
	close(errCh)
	for err := range errCh {
		require.NoError(t, err)
	}
}

func TestRedisGetNonce_FastForward(t *testing.T) {
	client, keyPrefix := redistestutils.NewRedisForTest(t)

	logger, err := zap.NewDevelopment()
	require.NoError(t, err)

	nonceManager := redismanager.NewRedisBackedNonceManager(client, logger, keyPrefix)

	// Start with nonces from 0
	err = nonceManager.Replenish(t.Context(), *big.NewInt(0))
	require.NoError(t, err)

	// Fast forward to 100
	err = nonceManager.FastForwardNonce(t.Context(), *big.NewInt(100))
	require.NoError(t, err)

	// Next nonce should be 100
	nonce, err := nonceManager.GetNonce(t.Context())
	require.NoError(t, err)
	require.EqualValues(t, 100, nonce.Nonce.Int64())
	err = nonce.Consume()
	require.NoError(t, err)

	// Next should be 101
	nonce, err = nonceManager.GetNonce(t.Context())
	require.NoError(t, err)
	require.EqualValues(t, 101, nonce.Nonce.Int64())
	err = nonce.Consume()
	require.NoError(t, err)
}

func TestRedisGetNonce_EmptySet(t *testing.T) {
	client, keyPrefix := redistestutils.NewRedisForTest(t)

	logger, err := zap.NewDevelopment()
	require.NoError(t, err)

	nonceManager := redismanager.NewRedisBackedNonceManager(client, logger, keyPrefix)

	// Try to get a nonce without replenishing first
	_, err = nonceManager.GetNonce(t.Context())
	require.Error(t, err)
	require.Contains(t, err.Error(), "no nonces available")
}

func TestRedisGetNonce_ContextCancellation(t *testing.T) {
	client, keyPrefix := redistestutils.NewRedisForTest(t)

	logger, err := zap.NewDevelopment()
	require.NoError(t, err)

	nonceManager := redismanager.NewRedisBackedNonceManager(client, logger, keyPrefix)

	// Test with already cancelled context
	cancelledCtx, cancel := context.WithCancel(t.Context())
	cancel()

	_, err = nonceManager.GetNonce(cancelledCtx)
	require.Error(t, err)
	require.Equal(t, context.Canceled, err)
}

func TestRedisGetNonce_KeyPrefix(t *testing.T) {
	client1, keyPrefix1 := redistestutils.NewRedisForTest(t)
	client2, keyPrefix2 := redistestutils.NewRedisForTest(t)

	logger, err := zap.NewDevelopment()
	require.NoError(t, err)

	// Test different key prefixes don't interfere
	manager1 := redismanager.NewRedisBackedNonceManager(client1, logger, keyPrefix1)
	manager2 := redismanager.NewRedisBackedNonceManager(client2, logger, keyPrefix2)

	// Replenish both with nonces starting from 0
	err = manager1.Replenish(t.Context(), *big.NewInt(0))
	require.NoError(t, err)
	err = manager2.Replenish(t.Context(), *big.NewInt(0))
	require.NoError(t, err)

	// Both should be able to get nonce 0
	nonce1, err := manager1.GetNonce(t.Context())
	require.NoError(t, err)
	require.EqualValues(t, 0, nonce1.Nonce.Int64())

	nonce2, err := manager2.GetNonce(t.Context())
	require.NoError(t, err)
	require.EqualValues(t, 0, nonce2.Nonce.Int64())

	err = nonce1.Consume()
	require.NoError(t, err)
	err = nonce2.Consume()
	require.NoError(t, err)
}

func TestRedisGetNonce_DefaultKeyPrefix(t *testing.T) {
	client, _ := redistestutils.NewRedisForTest(t)

	logger, err := zap.NewDevelopment()
	require.NoError(t, err)

	// Test that empty key prefix gets a default
	nonceManager := redismanager.NewRedisBackedNonceManager(client, logger, "")
	err = nonceManager.Replenish(t.Context(), *big.NewInt(0))
	require.NoError(t, err)

	nonce, err := nonceManager.GetNonce(t.Context())
	require.NoError(t, err)
	require.EqualValues(t, 0, nonce.Nonce.Int64())

	err = nonce.Consume()
	require.NoError(t, err)
}
