package common_test

import (
	"sync"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	c "github.com/xmtp/xmtpd/pkg/indexer/common"
	"github.com/xmtp/xmtpd/pkg/testutils"
)

var address = common.HexToAddress("0x0000000000000000000000000000000000000000")

func TestInitialize(t *testing.T) {
	ctx := t.Context()
	db, _ := testutils.NewDB(t, ctx)
	querier := queries.New(db)

	tracker, err := c.NewBlockTracker(ctx, address, querier)
	require.NoError(t, err)
	blockNumber, blockHash := tracker.GetLatestBlock()
	require.NoError(t, err)
	require.NotNil(t, tracker)
	require.Equal(t, uint64(0), blockNumber)
	require.Equal(t, common.Hash{}.Bytes(), blockHash)
}

func TestUpdateLatestBlock(t *testing.T) {
	ctx := t.Context()
	db, _ := testutils.NewDB(t, ctx)
	querier := queries.New(db)

	tracker, err := c.NewBlockTracker(ctx, address, querier)
	require.NoError(t, err)

	blockHigh := testutils.Int64ToHash(100).Bytes()
	blockLow := testutils.Int64ToHash(50).Bytes()

	// Test updating to a higher block
	err = tracker.UpdateLatestBlock(ctx, 100, blockHigh)
	blockNumber, blockHash := tracker.GetLatestBlock()
	require.NoError(t, err)
	require.Equal(t, uint64(100), blockNumber)
	require.Equal(t, blockHigh, blockHash)

	// Test updating to a lower block (should not update)
	err = tracker.UpdateLatestBlock(ctx, 50, blockLow)
	require.NoError(t, err)
	blockNumber, blockHash = tracker.GetLatestBlock()
	require.Equal(t, uint64(100), blockNumber)
	require.Equal(t, blockHigh, blockHash)

	// Test updating to the same block (should not update)
	err = tracker.UpdateLatestBlock(ctx, 100, blockHigh)
	require.NoError(t, err)
	blockNumber, blockHash = tracker.GetLatestBlock()
	require.Equal(t, uint64(100), blockNumber)
	require.Equal(t, blockHigh, blockHash)

	// Verify persistence
	newTracker, err := c.NewBlockTracker(ctx, address, querier)
	require.NoError(t, err)
	blockNumber, blockHash = newTracker.GetLatestBlock()
	require.Equal(t, uint64(100), blockNumber)
	require.Equal(t, blockHigh, blockHash)
}

func TestConcurrentUpdates(t *testing.T) {
	ctx := t.Context()
	db, _ := testutils.NewDB(t, ctx)
	querier := queries.New(db)

	tracker, err := c.NewBlockTracker(ctx, address, querier)
	require.NoError(t, err)

	numGoroutines := 10
	updatesPerGoroutine := 100

	errCh := make(chan error, numGoroutines)
	wg := sync.WaitGroup{}

	// Launch multiple goroutines to update the block number
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(startBlock int) {
			defer wg.Done()
			for j := 0; j < updatesPerGoroutine; j++ {
				blockNum := uint64(startBlock + j)
				err := tracker.UpdateLatestBlock(
					ctx,
					blockNum,
					testutils.Int64ToHash(int64(blockNum)).Bytes(),
				)
				if err != nil {
					errCh <- err
					break
				}
			}
			errCh <- nil
		}(i * updatesPerGoroutine)
	}

	wg.Wait()

	close(errCh)
	for err := range errCh {
		require.NoError(t, err)
	}

	// The final block number should be the highest one attempted
	expectedFinalBlock := uint64((numGoroutines-1)*updatesPerGoroutine + (updatesPerGoroutine - 1))
	blockNumber, blockHash := tracker.GetLatestBlock()
	require.Equal(t, expectedFinalBlock, blockNumber)

	expectedFinalHash := testutils.Int64ToHash(int64(expectedFinalBlock)).Bytes()
	require.Equal(t, expectedFinalHash, blockHash)

	// Verify persistence
	newTracker, err := c.NewBlockTracker(ctx, address, querier)
	require.NoError(t, err)
	blockNumber, blockHash = newTracker.GetLatestBlock()
	require.Equal(t, expectedFinalBlock, blockNumber)
	require.Equal(t, expectedFinalHash, blockHash)
}

func TestMultipleContractAddresses(t *testing.T) {
	ctx := t.Context()
	db, _ := testutils.NewDB(t, ctx)
	querier := queries.New(db)

	address1 := common.HexToAddress("0x0000000000000000000000000000000000000001")
	address2 := common.HexToAddress("0x0000000000000000000000000000000000000002")

	tracker1, err := c.NewBlockTracker(ctx, address1, querier)
	require.NoError(t, err)
	tracker2, err := c.NewBlockTracker(ctx, address2, querier)
	require.NoError(t, err)

	blockHash1 := testutils.Int64ToHash(100).Bytes()
	blockHash2 := testutils.Int64ToHash(200).Bytes()

	// Update trackers independently
	err = tracker1.UpdateLatestBlock(ctx, 100, blockHash1)
	require.NoError(t, err)
	err = tracker2.UpdateLatestBlock(ctx, 200, blockHash2)
	require.NoError(t, err)

	// Verify different addresses maintain separate block numbers
	blockNumber, blockHash := tracker1.GetLatestBlock()
	require.Equal(t, uint64(100), blockNumber)
	require.Equal(t, blockHash1, blockHash)
	blockNumber, blockHash = tracker2.GetLatestBlock()
	require.Equal(t, uint64(200), blockNumber)
	require.Equal(t, blockHash2, blockHash)

	// Verify persistence for both addresses
	newTracker1, err := c.NewBlockTracker(ctx, address1, querier)
	require.NoError(t, err)
	newTracker2, err := c.NewBlockTracker(ctx, address2, querier)
	require.NoError(t, err)

	blockNumber, blockHash = newTracker1.GetLatestBlock()
	require.Equal(t, uint64(100), blockNumber)
	require.Equal(t, blockHash1, blockHash)
	blockNumber, blockHash = newTracker2.GetLatestBlock()
	require.Equal(t, uint64(200), blockNumber)
	require.Equal(t, blockHash2, blockHash)
}
