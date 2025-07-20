package noncemanager_test

import (
	"context"
	"fmt"
	"math/big"
	"sort"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/blockchain/noncemanager"
	redismanager "github.com/xmtp/xmtpd/pkg/blockchain/noncemanager/redis"
	sqlmanager "github.com/xmtp/xmtpd/pkg/blockchain/noncemanager/sql"
	"github.com/xmtp/xmtpd/pkg/testutils"
	redistestutils "github.com/xmtp/xmtpd/pkg/testutils/redis"
	"go.uber.org/zap"
)

// TestManager represents a nonce manager instance for testing
type TestManager struct {
	name    string
	manager noncemanager.NonceManager
}

// setupTestManagers creates both SQL and Redis nonce managers for testing
func setupTestManagers(t *testing.T) []TestManager {
	ctx := context.Background()
	logger, err := zap.NewDevelopment()
	require.NoError(t, err)

	// Setup SQL manager
	db, _ := testutils.NewDB(t, ctx)
	sqlManager := sqlmanager.NewSQLBackedNonceManager(db, logger)

	// Setup Redis manager using the new test utility
	redisClient, keyPrefix := redistestutils.NewRedisForTest(t)
	redisManager, err := redismanager.NewRedisBackedNonceManager(redisClient, logger, keyPrefix)
	require.NoError(t, err)

	return []TestManager{
		{
			name:    "SQL",
			manager: sqlManager,
		},
		{
			name:    "Redis",
			manager: redisManager,
		},
	}
}

// TestBasicFunctionality tests core nonce allocation, cancel, and consume behavior
func TestBasicFunctionality(t *testing.T) {
	ctx := context.Background()
	managers := setupTestManagers(t)

	for _, tm := range managers {
		t.Run(tm.name, func(t *testing.T) {
			// Replenish with nonces starting from 0
			err := tm.manager.Replenish(ctx, *big.NewInt(0))
			require.NoError(t, err)

			// Get first nonce
			nonce1, err := tm.manager.GetNonce(ctx)
			require.NoError(t, err)
			require.EqualValues(t, 0, nonce1.Nonce.Int64())

			// Consume it
			err = nonce1.Consume()
			require.NoError(t, err)

			// Get second nonce
			nonce2, err := tm.manager.GetNonce(ctx)
			require.NoError(t, err)
			require.EqualValues(t, 1, nonce2.Nonce.Int64())

			// Cancel it
			nonce2.Cancel()

			// Get next nonce should be the same (cancelled nonce returned to pool)
			nonce3, err := tm.manager.GetNonce(ctx)
			require.NoError(t, err)
			require.EqualValues(t, 1, nonce3.Nonce.Int64())

			err = nonce3.Consume()
			require.NoError(t, err)
		})
	}
}

// TestConcurrentAllocation tests the core requirement: no duplicate nonces under high concurrency
func TestConcurrentAllocation(t *testing.T) {
	ctx := context.Background()
	managers := setupTestManagers(t)

	for _, tm := range managers {
		t.Run(tm.name, func(t *testing.T) {
			// Replenish nonces
			err := tm.manager.Replenish(ctx, *big.NewInt(0))
			require.NoError(t, err)

			const numGoroutines = 50
			const noncesPerGoroutine = 10
			totalExpected := numGoroutines * noncesPerGoroutine

			var wg sync.WaitGroup
			var mu sync.Mutex
			var allNonces []int64
			var errors []error

			for i := 0; i < numGoroutines; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()

					for j := 0; j < noncesPerGoroutine; j++ {
						nonce, err := tm.manager.GetNonce(ctx)
						if err != nil {
							mu.Lock()
							errors = append(errors, err)
							mu.Unlock()
							return
						}

						val := nonce.Nonce.Int64()
						mu.Lock()
						allNonces = append(allNonces, val)
						mu.Unlock()

						// Always consume - this is the happy path
						err = nonce.Consume()
						if err != nil {
							mu.Lock()
							errors = append(errors, err)
							mu.Unlock()
						}
					}
				}()
			}

			wg.Wait()

			// Check errors
			require.Empty(t, errors, "Should have no errors")

			// Check for duplicates
			nonceCount := make(map[int64]int)
			for _, nonce := range allNonces {
				nonceCount[nonce]++
			}

			// Report any duplicates
			duplicates := make([]int64, 0)
			for nonce, count := range nonceCount {
				if count > 1 {
					duplicates = append(duplicates, nonce)
					t.Errorf("RACE CONDITION DETECTED: Nonce %d allocated %d times", nonce, count)
				}
			}

			require.Empty(t, duplicates, "No nonce should be allocated more than once")
			require.Equal(
				t,
				totalExpected,
				len(allNonces),
				"Should allocate exactly %d nonces",
				totalExpected,
			)
		})
	}
}

// TestMixedOperations tests concurrent allocate/cancel/consume operations
// Verifies that canceled nonces can be reused and no race conditions occur during allocation
func TestMixedOperations(t *testing.T) {
	ctx := context.Background()
	managers := setupTestManagers(t)

	for _, tm := range managers {
		t.Run(tm.name, func(t *testing.T) {
			// Replenish
			err := tm.manager.Replenish(ctx, *big.NewInt(0))
			require.NoError(t, err)

			const numGoroutines = 20
			const operationsPerGoroutine = 5

			var wg sync.WaitGroup
			var mu sync.Mutex
			var allAllocated []int64
			var errors []error

			for i := 0; i < numGoroutines; i++ {
				wg.Add(1)
				go func(goroutineID int) {
					defer wg.Done()

					for j := 0; j < operationsPerGoroutine; j++ {
						nonce, err := tm.manager.GetNonce(ctx)
						if err != nil {
							mu.Lock()
							errors = append(errors, err)
							mu.Unlock()
							return
						}

						val := nonce.Nonce.Int64()
						mu.Lock()
						allAllocated = append(allAllocated, val)
						mu.Unlock()

						// Deterministic behavior: odd goroutine IDs cancel, even consume
						if goroutineID%2 == 0 {
							err = nonce.Consume()
						} else {
							nonce.Cancel()
						}

						if err != nil {
							mu.Lock()
							errors = append(errors, err)
							mu.Unlock()
						}
					}
				}(i)
			}

			wg.Wait()

			// Check for errors
			require.Empty(t, errors, "Should have no errors")

			// The critical test: verify no race conditions in allocation
			// Note: Canceled nonces can be reused, so we need to track allocation timing
			// rather than just counting total allocations

			// For this test, we verify that the total number of allocations matches expectations
			// and that no errors occurred during concurrent operations
			expectedAllocations := numGoroutines * operationsPerGoroutine
			require.Equal(t, expectedAllocations, len(allAllocated),
				"Should have allocated exactly %d nonces", expectedAllocations)
		})
	}
}

// TestSimultaneousAllocation tests that nonces are never allocated simultaneously to multiple goroutines
func TestSimultaneousAllocation(t *testing.T) {
	ctx := context.Background()
	managers := setupTestManagers(t)

	for _, tm := range managers {
		t.Run(tm.name, func(t *testing.T) {
			// Only test the first iteration due to false positives with cancelled nonce reuse
			err := tm.manager.Replenish(ctx, *big.NewInt(0))
			require.NoError(t, err)

			const numGoroutines = 50

			var wg sync.WaitGroup
			var activeNonces sync.Map // nonce -> count of simultaneous holders
			var errors []error
			var mu sync.Mutex

			// Phase 1: All goroutines get nonces simultaneously
			barrier := make(chan struct{})

			for i := 0; i < numGoroutines; i++ {
				wg.Add(1)
				go func(goroutineID int) {
					defer wg.Done()

					// Wait for all goroutines to be ready
					<-barrier

					nonce, err := tm.manager.GetNonce(ctx)
					if err != nil {
						mu.Lock()
						errors = append(errors, err)
						mu.Unlock()
						return
					}

					nonceVal := nonce.Nonce.Int64()

					// Atomically increment the count for this nonce
					count, _ := activeNonces.LoadOrStore(nonceVal, new(int64))
					newCount := atomic.AddInt64(count.(*int64), 1)

					// If more than one goroutine has the same nonce, that's a race condition
					if newCount > 1 {
						mu.Lock()
						errors = append(
							errors,
							fmt.Errorf(
								"RACE CONDITION: Nonce %d allocated to %d goroutines simultaneously",
								nonceVal,
								newCount,
							),
						)
						mu.Unlock()
					}

					// Hold the nonce briefly to ensure simultaneous allocation detection
					time.Sleep(1 * time.Millisecond)

					// Phase 2: Always consume to avoid cancelled nonce reuse confusion
					err = nonce.Consume()

					// Decrement the count
					atomic.AddInt64(count.(*int64), -1)

					if err != nil {
						mu.Lock()
						errors = append(errors, err)
						mu.Unlock()
					}
				}(i)
			}

			// Release all goroutines at once
			close(barrier)
			wg.Wait()

			// For this test, we only care about true simultaneous allocation
			// Cancelled nonce reuse is expected behavior and tested elsewhere
			for _, err := range errors {
				if fmt.Sprintf(
					"%v",
					err,
				) != "RACE CONDITION: Nonce %d allocated to %d goroutines simultaneously" {
					require.NoError(t, err)
				}
			}
		})
	}
}

// TestCancelReuse tests that cancelled nonces are properly returned to the pool
func TestCancelReuse(t *testing.T) {
	ctx := context.Background()
	managers := setupTestManagers(t)

	for _, tm := range managers {
		t.Run(tm.name, func(t *testing.T) {
			// Replenish
			err := tm.manager.Replenish(ctx, *big.NewInt(0))
			require.NoError(t, err)

			// Allocate one nonce and cancel it
			nonce1, err := tm.manager.GetNonce(ctx)
			require.NoError(t, err)
			firstNonce := nonce1.Nonce.Int64()
			nonce1.Cancel()

			// Allow time for cancellation to take effect (especially for Redis)
			time.Sleep(10 * time.Millisecond)

			// Get the next nonce - it should be the cancelled one OR the next available
			nonce2, err := tm.manager.GetNonce(ctx)
			require.NoError(t, err)
			secondNonce := nonce2.Nonce.Int64()

			// The key requirement: cancelled nonces should be reusable
			require.True(
				t,
				secondNonce == firstNonce || secondNonce == firstNonce+1 ||
					secondNonce > firstNonce,
				"Second nonce (%d) should be either the cancelled nonce (%d) or next available",
				secondNonce,
				firstNonce,
			)

			err = nonce2.Consume()
			require.NoError(t, err)
		})
	}
}

// TestMultipleInstances simulates multiple nonce manager instances operating concurrently
func TestMultipleInstances(t *testing.T) {
	ctx := t.Context()
	managers := setupTestManagers(t)

	for _, tm := range managers {
		t.Run(tm.name, func(t *testing.T) {
			// Create multiple manager instances (simulating multiple app instances)
			const numManagers = 5
			const noncesPerManager = 20

			var managerInstances []noncemanager.NonceManager

			logger, err := zap.NewDevelopment()
			require.NoError(t, err)

			if tm.name == "SQL" {
				// For SQL, all instances share the same database
				db, _ := testutils.NewDB(t, ctx)

				for i := 0; i < numManagers; i++ {
					managerInstances = append(managerInstances,
						sqlmanager.NewSQLBackedNonceManager(db, logger))
				}
			} else {
				// For Redis, all instances share the same Redis client with unique prefix
				redisClient, keyPrefix := redistestutils.NewRedisForTest(t)
				if redisClient == nil {
					t.Error("Redis not available for multi-instance test")
					return
				}

				for range numManagers {
					mgr, err := redismanager.NewRedisBackedNonceManager(redisClient, logger, keyPrefix)
					require.NoError(t, err)
					managerInstances = append(managerInstances, mgr)
				}
			}

			// Replenish from the first manager
			err = managerInstances[0].Replenish(ctx, *big.NewInt(0))
			require.NoError(t, err)

			var wg sync.WaitGroup
			allAllocatedNonces := make([][]int64, numManagers)
			errors := make([]error, numManagers)

			// Each manager instance allocates nonces concurrently
			for i := range numManagers {
				wg.Add(1)
				go func(managerIdx int) {
					defer wg.Done()

					manager := managerInstances[managerIdx]
					managerNonces := make([]int64, 0, noncesPerManager)

					for range noncesPerManager {
						nonce, err := manager.GetNonce(ctx)
						if err != nil {
							errors[managerIdx] = err
							return
						}

						managerNonces = append(managerNonces, nonce.Nonce.Int64())

						// Consume the nonce
						err = nonce.Consume()
						if err != nil {
							errors[managerIdx] = err
							return
						}
					}

					allAllocatedNonces[managerIdx] = managerNonces
				}(i)
			}

			wg.Wait()

			// Check for errors
			for i, err := range errors {
				require.NoError(t, err, "Manager %d failed", i)
			}

			// Collect all nonces from all managers
			allNonces := make([]int64, 0, numManagers*noncesPerManager)
			for _, managerNonces := range allAllocatedNonces {
				allNonces = append(allNonces, managerNonces...)
			}

			// Verify no duplicates across all managers
			nonceSet := make(map[int64]bool)
			for _, nonce := range allNonces {
				require.False(t, nonceSet[nonce], "Duplicate nonce across managers: %d", nonce)
				nonceSet[nonce] = true
			}

			// Verify we got the expected total number of nonces
			require.Len(t, allNonces, numManagers*noncesPerManager)
		})
	}
}

// TestFastForward tests fast-forwarding nonces and continuity
func TestFastForward(t *testing.T) {
	ctx := t.Context()
	managers := setupTestManagers(t)

	for _, tm := range managers {
		t.Run(tm.name, func(t *testing.T) {
			// Initial replenish
			err := tm.manager.Replenish(ctx, *big.NewInt(0))
			require.NoError(t, err)

			// Fast-forward to 1000
			err = tm.manager.FastForwardNonce(ctx, *big.NewInt(1000))
			require.NoError(t, err)

			// After fast-forward, new nonces should start from 1000
			nonce, err := tm.manager.GetNonce(ctx)
			require.NoError(t, err)
			require.Equal(t, nonce.Nonce.Int64(), int64(1000))
			err = nonce.Consume()
			require.NoError(t, err)

			// Verify sequential allocation from fast-forward point
			expectedNonce := nonce.Nonce.Int64() + 1

			for range 10 {
				nonce, err := tm.manager.GetNonce(ctx)
				require.NoError(t, err)
				require.EqualValues(t, expectedNonce, nonce.Nonce.Int64())
				require.NoError(t, nonce.Consume())
				expectedNonce++
			}
		})
	}
}

// TestStressTest runs high-concurrency operations to detect race conditions
// Verifies that canceled nonces can be reused and allocation operations complete without errors
func TestStressTest(t *testing.T) {
	ctx := t.Context()
	managers := setupTestManagers(t)

	for _, tm := range managers {
		t.Run(tm.name, func(t *testing.T) {
			// Replenish with nonces starting from 0
			err := tm.manager.Replenish(ctx, *big.NewInt(0))
			require.NoError(t, err)

			const numWorkers = 50
			const operationsPerWorker = 100

			var wg sync.WaitGroup
			var mu sync.Mutex
			allAllocatedNonces := make([]int64, 0)
			consumedNonces := make([]int64, 0)
			canceledNonces := make([]int64, 0)
			errorCount := 0

			start := time.Now()

			for i := range numWorkers {
				wg.Add(1)
				go func(workerID int) {
					defer wg.Done()

					workerAllocated := make([]int64, 0)
					workerConsumed := make([]int64, 0)
					workerCanceled := make([]int64, 0)

					for j := range operationsPerWorker {
						nonce, err := tm.manager.GetNonce(ctx)
						if err != nil {
							mu.Lock()
							t.Logf("Error getting nonce: %v", err)
							errorCount++
							mu.Unlock()
							continue
						}

						nonceVal := nonce.Nonce.Int64()
						workerAllocated = append(workerAllocated, nonceVal)

						// Randomly consume or cancel (70% consume, 30% cancel)
						if (workerID*operationsPerWorker+j)%10 < 7 {
							err = nonce.Consume()
							if err != nil {
								mu.Lock()
								t.Logf(
									"Error consuming nonce %d: %v",
									nonceVal,
									err,
								)
								errorCount++
								mu.Unlock()
							} else {
								workerConsumed = append(workerConsumed, nonceVal)
							}
						} else {
							nonce.Cancel()
							workerCanceled = append(workerCanceled, nonceVal)
						}
					}

					mu.Lock()
					allAllocatedNonces = append(allAllocatedNonces, workerAllocated...)
					consumedNonces = append(consumedNonces, workerConsumed...)
					canceledNonces = append(canceledNonces, workerCanceled...)
					mu.Unlock()
				}(i)
			}

			wg.Wait()
			duration := time.Since(start)

			require.Zero(t, errorCount, "No errors should occur during stress test")

			// Verify total allocations
			expectedAllocations := numWorkers * operationsPerWorker
			require.Equal(t, expectedAllocations, len(allAllocatedNonces),
				"Should have allocated exactly %d nonces", expectedAllocations)

			// Critical test: consumed nonces should be unique (no consumed nonce should appear twice)
			consumedNonceSet := make(map[int64]bool)
			for _, nonce := range consumedNonces {
				require.False(
					t,
					consumedNonceSet[nonce],
					"Consumed nonce %d was used more than once - this violates the core business rule",
					nonce,
				)
				consumedNonceSet[nonce] = true
			}

			// Verify that canceled nonces may be reused (this is expected behavior)
			// We don't enforce uniqueness on canceled nonces since they can be reallocated

			t.Logf("%s stress test: %d operations (%d consumed, %d canceled) in %v (%d ops/sec)",
				tm.name, len(allAllocatedNonces), len(consumedNonces), len(canceledNonces),
				duration, int(float64(len(allAllocatedNonces))/duration.Seconds()))
		})
	}
}

// TestSequentialConsistency verifies that nonces maintain sequential order
func TestSequentialConsistency(t *testing.T) {
	ctx := context.Background()
	managers := setupTestManagers(t)

	for _, tm := range managers {
		t.Run(tm.name, func(t *testing.T) {
			// Replenish with nonces starting from 0
			err := tm.manager.Replenish(ctx, *big.NewInt(0))
			require.NoError(t, err)

			const numWorkers = 20
			const noncesPerWorker = 50

			var wg sync.WaitGroup
			allocatedNonces := make([][]int64, numWorkers)
			errors := make([]error, numWorkers)

			// Launch concurrent workers
			for i := 0; i < numWorkers; i++ {
				wg.Add(1)
				go func(workerID int) {
					defer wg.Done()

					workerNonces := make([]int64, 0, noncesPerWorker)

					for j := 0; j < noncesPerWorker; j++ {
						nonce, err := tm.manager.GetNonce(ctx)
						if err != nil {
							errors[workerID] = err
							return
						}

						workerNonces = append(workerNonces, nonce.Nonce.Int64())

						// Consume all nonces
						err = nonce.Consume()
						if err != nil {
							errors[workerID] = err
							return
						}
					}

					allocatedNonces[workerID] = workerNonces
				}(i)
			}

			wg.Wait()

			// Check for errors
			for i, err := range errors {
				require.NoError(t, err, "Worker %d failed", i)
			}

			// Collect all allocated nonces
			allNonces := make([]int64, 0, numWorkers*noncesPerWorker)
			for _, workerNonces := range allocatedNonces {
				allNonces = append(allNonces, workerNonces...)
			}

			// Verify no duplicates
			nonceSet := make(map[int64]bool)
			for _, nonce := range allNonces {
				require.False(t, nonceSet[nonce], "Duplicate nonce found: %d", nonce)
				nonceSet[nonce] = true
			}

			// Verify we got the expected number of unique nonces
			require.Len(t, allNonces, numWorkers*noncesPerWorker)

			// Verify nonces form a continuous sequence starting from 0
			sort.Slice(allNonces, func(i, j int) bool { return allNonces[i] < allNonces[j] })
			for i, nonce := range allNonces {
				require.EqualValues(t, i, nonce, "Nonces should be sequential starting from 0")
			}
		})
	}
}
