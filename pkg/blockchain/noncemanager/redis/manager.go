// Package redis provides a Redis-backed implementation of the NonceManager interface.
// It uses Redis sorted sets and atomic operations to ensure consistent nonce allocation
// even under high concurrency across multiple instances.
package redis

import (
	"context"
	"errors"
	"fmt"
	"math"
	"math/big"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/xmtp/xmtpd/pkg/blockchain/noncemanager"
	"github.com/xmtp/xmtpd/pkg/metrics"
	"github.com/xmtp/xmtpd/pkg/utils"
	"go.uber.org/zap"
)

const (
	// StaleReservationTimeout is how long before a reservation is considered stale
	StaleReservationTimeout = 30 * time.Second
	// BatchSize is the number of nonces to generate in a single replenish operation
	BatchSize = 10000
	// RedisOperationTimeout bounds how long cancel/consume Redis calls can take
	RedisOperationTimeout = 5 * time.Second
	// maxRetries is the number of retry attempts for cancel/consume operations
	maxRetries = 3
	// initialBackoff is the starting delay between retries (doubles each attempt)
	initialBackoff = 50 * time.Millisecond
)

// Lua scripts for atomic operations
var (
	luaCleanupAndReserveScript = `
		local availableKey = KEYS[1]
		local reservedKey = KEYS[2]
		local staleThreshold = ARGV[1]
		local currentTime = ARGV[2]

		-- First, cleanup stale reservations
		local staleNonces = redis.call('ZRANGEBYSCORE', reservedKey, '-inf', staleThreshold)
		local cleanupCount = 0

		if #staleNonces > 0 then
			-- Move each stale nonce back to available pool
			for i, nonce in ipairs(staleNonces) do
				redis.call('ZADD', availableKey, nonce, nonce)
			end

			-- Remove from reserved set
			redis.call('ZREMRANGEBYSCORE', reservedKey, '-inf', staleThreshold)
			cleanupCount = #staleNonces
		end

		-- Then, reserve the next available nonce
		local result = redis.call('ZPOPMIN', availableKey, 1)
		if #result == 0 then
			return {nil, cleanupCount}
		end

		local nonce = result[1]
		-- Add it to the reserved set with current timestamp as score for cleanup
		redis.call('ZADD', reservedKey, currentTime, nonce)

		return {nonce, cleanupCount}
	`
)

// RedisBackedNonceManager implements NonceManager using Redis for persistence.
// It provides distributed nonce allocation with configurable concurrency limits.
type RedisBackedNonceManager struct {
	client    redis.UniversalClient
	logger    *zap.Logger
	limiter   *noncemanager.OpenConnectionsLimiter
	keyPrefix string
}

var _ noncemanager.NonceManager = &RedisBackedNonceManager{}

// NewRedisBackedNonceManager creates a new Redis-backed nonce manager
func NewRedisBackedNonceManager(
	client redis.UniversalClient,
	logger *zap.Logger,
	keyPrefix string,
) (*RedisBackedNonceManager, error) {
	if client == nil {
		return nil, errors.New("redis client cannot be nil")
	}
	if logger == nil {
		return nil, errors.New("logger cannot be nil")
	}
	if keyPrefix == "" {
		keyPrefix = "xmtpd:nonces:"
	}

	redisBackedNonceManagerLogger := logger.Named(utils.RedisNonceManagerLoggerName)

	return &RedisBackedNonceManager{
		client:    client,
		logger:    redisBackedNonceManagerLogger,
		limiter:   noncemanager.NewOpenConnectionsLimiter(noncemanager.BestGuessConcurrency),
		keyPrefix: keyPrefix,
	}, nil
}

// availableKey returns the Redis key for the available nonces sorted set
func (r *RedisBackedNonceManager) availableKey() string {
	return r.keyPrefix + "available"
}

// reservedKey returns the Redis key for the reserved nonces sorted set
func (r *RedisBackedNonceManager) reservedKey() string {
	return r.keyPrefix + "reserved"
}

// GetNonce atomically reserves the next available nonce from Redis.
// It moves nonces from available set to reserved set to prevent concurrent allocation,
// similar to how SQL uses SELECT FOR UPDATE SKIP LOCKED.
func (r *RedisBackedNonceManager) GetNonce(
	ctx context.Context,
) (*noncemanager.NonceContext, error) {
	r.limiter.WG.Add(1)

	// Block until there is an available slot in the blockchain rate limiter
	select {
	case r.limiter.Semaphore <- struct{}{}:
	case <-ctx.Done():
		r.limiter.WG.Done()
		return nil, ctx.Err()
	}

	// Clean up stale reservations and get next available nonce in a single Redis call
	nonce, err := r.cleanupAndReserveNonce(ctx)
	if err != nil {
		r.releaseLimiter()
		return nil, err
	}

	metrics.EmitGatewayCurrentNonce(float64(nonce))

	return r.createNonceContext(nonce), nil
}

// Replenish ensures a sufficient number of nonces are available starting from the given nonce.
// It generates up to 10,000 nonces in a single batch operation using ZADD.
func (r *RedisBackedNonceManager) Replenish(ctx context.Context, nonce big.Int) error {
	if !nonce.IsInt64() || nonce.Sign() < 0 {
		return fmt.Errorf(
			"nonce value %s is out of int64 range, cannot replenish",
			nonce.String(),
		)
	}
	startNonce := nonce.Int64()
	if startNonce > math.MaxInt64-int64(BatchSize) {
		return fmt.Errorf(
			"nonce value %d would overflow when generating batch of %d",
			startNonce,
			BatchSize,
		)
	}

	// Prepare the nonces to add
	members := make([]redis.Z, BatchSize)
	for i := range int64(BatchSize) {
		nonceVal := startNonce + i
		members[i] = redis.Z{
			Score:  float64(nonceVal),
			Member: nonceVal,
		}
	}

	// Use a pipeline for efficiency
	pipe := r.client.Pipeline()
	pipe.ZAdd(ctx, r.availableKey(), members...)

	if _, err := pipe.Exec(ctx); err != nil {
		return fmt.Errorf("failed to replenish nonces in Redis: %w", err)
	}

	r.logger.Debug(
		"replenished nonces",
		utils.StartingNonceField(nonce.Uint64()),
		utils.NumNoncesField(BatchSize),
	)

	return nil
}

// FastForwardNonce sets the nonce sequence to start from the given value and removes
// all nonces below it. This is typically used when recovering from blockchain state issues.
func (r *RedisBackedNonceManager) FastForwardNonce(ctx context.Context, nonce big.Int) error {
	if !nonce.IsInt64() || nonce.Sign() < 0 {
		return fmt.Errorf(
			"nonce value %s is out of int64 range, cannot fast forward",
			nonce.String(),
		)
	}

	// First replenish nonces starting from the given value
	err := r.Replenish(ctx, nonce)
	if err != nil {
		return err
	}

	// Remove all obsolete nonces below the given threshold
	_, err = r.client.ZRemRangeByScore(ctx, r.availableKey(), "-inf", fmt.Sprintf("(%d", nonce.Int64())).
		Result()
	if err != nil {
		return fmt.Errorf("failed to remove obsolete nonces from Redis: %w", err)
	}

	return nil
}

// cleanupAndReserveNonce atomically cleans up stale reservations and reserves the next available nonce
func (r *RedisBackedNonceManager) cleanupAndReserveNonce(ctx context.Context) (int64, error) {
	staleThreshold := float64(time.Now().Add(-StaleReservationTimeout).Unix())
	currentTime := float64(time.Now().Unix())

	result, err := r.client.Eval(ctx, luaCleanupAndReserveScript, []string{r.availableKey(), r.reservedKey()}, staleThreshold, currentTime).
		Result()
	if err != nil {
		return 0, fmt.Errorf("failed to cleanup and reserve nonce from Redis: %w", err)
	}

	resultArray, ok := result.([]any)
	if !ok {
		return 0, errors.New("invalid result format from cleanup and reserve script")
	}

	// Check if a nonce was reserved (first element)
	if len(resultArray) == 0 || resultArray[0] == nil {
		return 0, errors.New("no nonces available in Redis")
	}

	nonceStr, ok := resultArray[0].(string)
	if !ok {
		return 0, fmt.Errorf("unexpected nonce type %T from Redis", resultArray[0])
	}

	nonce, err := strconv.ParseInt(nonceStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse nonce from Redis: %w", err)
	}

	// Log cleanup count if any stale reservations were cleaned up
	if cleanupCount, ok := resultArray[1].(int64); ok && cleanupCount > 0 {
		r.logger.Info("cleaned up stale nonce reservations", utils.CountField(cleanupCount))
	}

	return nonce, nil
}

// releaseLimiter releases the semaphore and decrements the wait group
func (r *RedisBackedNonceManager) releaseLimiter() {
	<-r.limiter.Semaphore
	r.limiter.WG.Done()
}

// createNonceContext creates a NonceContext with Cancel and Consume functions
func (r *RedisBackedNonceManager) createNonceContext(nonce int64) *noncemanager.NonceContext {
	var operationDone atomic.Int32 // 0 = not done, 1 = done

	return &noncemanager.NonceContext{
		Nonce: *new(big.Int).SetInt64(nonce),
		Cancel: func() {
			if !operationDone.CompareAndSwap(0, 1) {
				return // Already cancelled or consumed
			}

			r.releaseLimiter()
			r.cancelNonce(nonce)
		},
		Consume: func() error {
			if !operationDone.CompareAndSwap(0, 1) {
				return fmt.Errorf("nonce %d already consumed or cancelled", nonce)
			}

			r.releaseLimiter()
			r.consumeNonce(nonce)
			return nil
		},
	}
}

// cancelNonce returns a nonce to the available pool with retry on transient failures
func (r *RedisBackedNonceManager) cancelNonce(nonce int64) {
	attempts, err := utils.RetryWithBackoff(func() error {
		ctx, cancel := context.WithTimeout(context.Background(), RedisOperationTimeout)
		defer cancel()

		pipe := r.client.Pipeline()
		pipe.ZRem(ctx, r.reservedKey(), nonce)
		pipe.ZAdd(ctx, r.availableKey(), redis.Z{
			Score:  float64(nonce),
			Member: nonce,
		})
		_, err := pipe.Exec(ctx)
		return err
	}, maxRetries, initialBackoff)
	if err != nil {
		r.logger.Error("failed to return cancelled nonce to Redis",
			utils.NonceField(uint64(nonce)), zap.Error(err))
		return
	}
	if attempts > 1 {
		r.logger.Warn("cancelled nonce returned to Redis after retry",
			utils.NonceField(uint64(nonce)), zap.Int("attempts", attempts))
	}
}

// consumeNonce removes a nonce from the reserved pool with retry on transient failures
func (r *RedisBackedNonceManager) consumeNonce(nonce int64) {
	attempts, err := utils.RetryWithBackoff(func() error {
		ctx, cancel := context.WithTimeout(context.Background(), RedisOperationTimeout)
		defer cancel()

		return r.client.ZRem(ctx, r.reservedKey(), nonce).Err()
	}, maxRetries, initialBackoff)
	if err != nil {
		r.logger.Error("failed to remove consumed nonce from reserved set",
			utils.NonceField(uint64(nonce)), zap.Error(err))
		return
	}
	if attempts > 1 {
		r.logger.Warn("consumed nonce removed from reserved set after retry",
			utils.NonceField(uint64(nonce)), zap.Int("attempts", attempts))
	}
}
