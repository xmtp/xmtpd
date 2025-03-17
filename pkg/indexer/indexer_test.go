package indexer

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/mock"
	"github.com/xmtp/xmtpd/pkg/indexer/storer"
	blockchainMocks "github.com/xmtp/xmtpd/pkg/mocks/blockchain"
	indexerMocks "github.com/xmtp/xmtpd/pkg/mocks/indexer"
	storerMocks "github.com/xmtp/xmtpd/pkg/mocks/storer"
	"github.com/xmtp/xmtpd/pkg/testutils"
)

func TestIndexLogsSuccess(t *testing.T) {
	channel := make(chan types.Log, 10)
	reorgChannel := make(chan uint64, 1)
	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		cancel()
		close(channel)
		close(reorgChannel)
	}()

	newBlockNumber := uint64(10)
	newBlockHash := common.HexToHash(
		"0x0000000000000000000000000000000000000000000000000000000000000000",
	)

	event := types.Log{
		Address:     common.HexToAddress("0x123"),
		BlockNumber: newBlockNumber,
		BlockHash:   newBlockHash,
	}

	channel <- event

	mockClient := blockchainMocks.NewMockChainClient(t)

	var wg sync.WaitGroup
	wg.Add(2) // Expecting two calls: StoreLog and UpdateLatestBlock

	blockTracker := indexerMocks.NewMockIBlockTracker(t)
	blockTracker.EXPECT().
		UpdateLatestBlock(mock.Anything, newBlockNumber, newBlockHash.Bytes()).
		Run(func(ctx context.Context, blockNum uint64, blockHash []byte) {
			wg.Done()
		}).
		Return(nil)

	reorgHandler := indexerMocks.NewMockChainReorgHandler(t)

	logStorer := storerMocks.NewMockLogStorer(t)
	logStorer.EXPECT().
		StoreLog(mock.Anything, event).
		Run(func(ctx context.Context, log types.Log) {
			wg.Done()
		}).
		Return(nil)

	go indexLogs(
		ctx,
		mockClient,
		channel,
		reorgChannel,
		testutils.NewLog(t),
		logStorer,
		blockTracker,
		reorgHandler,
		"testContract",
	)

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		// Test passed
	case <-time.After(1 * time.Second):
		t.Fatal("Test timed out waiting for StoreLog and UpdateLatestBlock")
	}
}

func TestIndexLogsRetryableError(t *testing.T) {
	channel := make(chan types.Log, 10)
	reorgChannel := make(chan uint64, 1)
	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		cancel()
		close(channel)
		close(reorgChannel)
	}()

	newBlockNumber := uint64(10)
	newBlockHash := common.HexToHash(
		"0x0000000000000000000000000000000000000000000000000000000000000000",
	)

	event := types.Log{
		Address:     common.HexToAddress("0x123"),
		BlockNumber: newBlockNumber,
		BlockHash:   newBlockHash,
	}

	mockClient := blockchainMocks.NewMockChainClient(t)
	logStorer := storerMocks.NewMockLogStorer(t)
	blockTracker := indexerMocks.NewMockIBlockTracker(t)
	reorgHandler := indexerMocks.NewMockChainReorgHandler(t)

	var wg sync.WaitGroup
	wg.Add(2) // Expecting two calls: StoreLog and UpdateLatestBlock

	// Will fail for the first call with a retryable error and a non-retryable error on the second call
	attemptNumber := 0

	logStorer.EXPECT().
		StoreLog(mock.Anything, event).
		RunAndReturn(func(ctx context.Context, log types.Log) storer.LogStorageError {
			wg.Done()
			attemptNumber++
			if attemptNumber < 2 {
				return storer.NewRetryableLogStorageError(errors.New("retryable error"))
			} else {
				return storer.NewUnrecoverableLogStorageError(errors.New("non-retryable error"))
			}
		})

	channel <- event

	go indexLogs(
		ctx,
		mockClient,
		channel,
		reorgChannel,
		testutils.NewLog(t),
		logStorer,
		blockTracker,
		reorgHandler,
		"testContract",
	)

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		// Test passed
	case <-time.After(1 * time.Second):
		t.Fatal("Test timed out waiting for StoreLog and UpdateLatestBlock")
	}

	logStorer.AssertNumberOfCalls(t, "StoreLog", 2)
}
