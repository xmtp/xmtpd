package redis_test

import (
	"context"
	"math/big"
	"sync"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
	redismanager "github.com/xmtp/xmtpd/pkg/blockchain/noncemanager/redis"
	redistestutils "github.com/xmtp/xmtpd/pkg/testutils/redis"
	"go.uber.org/zap"
)

func TestRedisGetNonce_Simple(t *testing.T) {
	client, keyPrefix := redistestutils.NewRedisForTest(t)

	logger, err := zap.NewDevelopment()
	require.NoError(t, err)

	nonceManager, err := redismanager.NewRedisBackedNonceManager(client, logger, keyPrefix)
	require.NoError(t, err)
	require.NoError(t, nonceManager.Replenish(t.Context(), *big.NewInt(0)))

	nonce, err := nonceManager.GetNonce(t.Context())
	require.NoError(t, err)
	defer nonce.Cancel()

	require.EqualValues(t, 0, nonce.Nonce.Int64())
}

func TestRedisGetNonce_RevertMany(t *testing.T) {
	client, keyPrefix := redistestutils.NewRedisForTest(t)

	logger, err := zap.NewDevelopment()
	require.NoError(t, err)

	nonceManager, err := redismanager.NewRedisBackedNonceManager(client, logger, keyPrefix)
	require.NoError(t, err)
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

	nonceManager, err := redismanager.NewRedisBackedNonceManager(client, logger, keyPrefix)
	require.NoError(t, err)
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

	nonceManager, err := redismanager.NewRedisBackedNonceManager(client, logger, keyPrefix)
	require.NoError(t, err)
	err = nonceManager.Replenish(t.Context(), *big.NewInt(0))
	require.NoError(t, err)

	var wg sync.WaitGroup
	numClients := 20
	errCh := make(chan error, numClients)

	for range numClients {
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

	nonceManager, err := redismanager.NewRedisBackedNonceManager(client, logger, keyPrefix)
	require.NoError(t, err)

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

	nonceManager, err := redismanager.NewRedisBackedNonceManager(client, logger, keyPrefix)
	require.NoError(t, err)

	// Try to get a nonce without replenishing first
	_, err = nonceManager.GetNonce(t.Context())
	require.Error(t, err)
	require.Contains(t, err.Error(), "no nonces available")
}

func TestRedisGetNonce_ContextCancellation(t *testing.T) {
	client, keyPrefix := redistestutils.NewRedisForTest(t)

	logger, err := zap.NewDevelopment()
	require.NoError(t, err)

	nonceManager, err := redismanager.NewRedisBackedNonceManager(client, logger, keyPrefix)
	require.NoError(t, err)

	// Test with already cancelled context
	cancelledCtx, cancel := context.WithCancel(t.Context())
	cancel()

	_, err = nonceManager.GetNonce(cancelledCtx)
	require.Error(t, err)
	require.Contains(t, err.Error(), "context canceled")
}

func TestRedisGetNonce_KeyPrefix(t *testing.T) {
	client1, keyPrefix1 := redistestutils.NewRedisForTest(t)
	// The keyPrefix uses the timestamp as a tiebreak when run from the same test.
	// Ensure we get a distinct prefix
	time.Sleep(1 * time.Millisecond)
	client2, keyPrefix2 := redistestutils.NewRedisForTest(t)

	logger, err := zap.NewDevelopment()
	require.NoError(t, err)

	// Test different key prefixes don't interfere
	manager1, err := redismanager.NewRedisBackedNonceManager(client1, logger, keyPrefix1)
	require.NoError(t, err)
	manager2, err := redismanager.NewRedisBackedNonceManager(client2, logger, keyPrefix2)
	require.NoError(t, err)

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
	nonceManager, err := redismanager.NewRedisBackedNonceManager(client, logger, "")
	require.NoError(t, err)
	err = nonceManager.Replenish(t.Context(), *big.NewInt(0))
	require.NoError(t, err)

	nonce, err := nonceManager.GetNonce(t.Context())
	require.NoError(t, err)
	require.EqualValues(t, 0, nonce.Nonce.Int64())

	err = nonce.Consume()
	require.NoError(t, err)
}

func TestRedisGetNonce_StaleNonceCleanup(t *testing.T) {
	client, keyPrefix := redistestutils.NewRedisForTest(t)

	logger, err := zap.NewDevelopment()
	require.NoError(t, err)

	nonceManager, err := redismanager.NewRedisBackedNonceManager(client, logger, keyPrefix)
	require.NoError(t, err)
	err = nonceManager.Replenish(t.Context(), *big.NewInt(0))
	require.NoError(t, err)

	// Get a nonce but don't consume or cancel it
	staleNonce, err := nonceManager.GetNonce(t.Context())
	require.NoError(t, err)
	require.EqualValues(t, 0, staleNonce.Nonce.Int64())

	// Manually add timestamp to make it stale (31 seconds ago)
	staleTime := float64(time.Now().Add(-31 * time.Second).Unix())
	err = client.ZAdd(t.Context(), keyPrefix+"reserved", []redis.Z{
		{Score: staleTime, Member: 0},
	}...).Err()
	require.NoError(t, err)

	// Get another nonce - this should trigger cleanup
	nonce2, err := nonceManager.GetNonce(t.Context())
	require.NoError(t, err)
	defer nonce2.Cancel()

	// The cleanup should have returned nonce 0 to the available pool
	// and we should get it again
	require.EqualValues(t, 0, nonce2.Nonce.Int64())

	// Verify the stale nonce is no longer in the reserved set
	count, err := client.ZCard(t.Context(), keyPrefix+"reserved").Result()
	require.NoError(t, err)
	require.EqualValues(t, 1, count) // Only nonce2 should be reserved
}

func TestNewRedisBackedNonceManager_NilChecks(t *testing.T) {
	client, keyPrefix := redistestutils.NewRedisForTest(t)
	logger, err := zap.NewDevelopment()
	require.NoError(t, err)

	// Test nil client
	_, err = redismanager.NewRedisBackedNonceManager(nil, logger, keyPrefix)
	require.Error(t, err)
	require.Contains(t, err.Error(), "redis client cannot be nil")

	// Test nil logger
	_, err = redismanager.NewRedisBackedNonceManager(client, nil, keyPrefix)
	require.Error(t, err)
	require.Contains(t, err.Error(), "logger cannot be nil")

	// Test both nil
	_, err = redismanager.NewRedisBackedNonceManager(nil, nil, keyPrefix)
	require.Error(t, err)
	// Should return first error (client)
	require.Contains(t, err.Error(), "redis client cannot be nil")
}
