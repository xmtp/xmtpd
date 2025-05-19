package common_test

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/mock"
	errorMocks "github.com/xmtp/xmtpd/pkg/errors"
	c "github.com/xmtp/xmtpd/pkg/indexer/common"
	blockchainMocks "github.com/xmtp/xmtpd/pkg/mocks/blockchain"
	indexerMocks "github.com/xmtp/xmtpd/pkg/mocks/common"
	"go.uber.org/zap"
)

// TODO: Add more test coverage.

func setup() (chan types.Log, chan uint64, context.Context, context.CancelFunc, types.Log, uint64, common.Hash) {
	channel := make(chan types.Log, 10)
	reorgChannel := make(chan uint64, 1)
	ctx, cancel := context.WithCancel(context.Background())

	newBlockNumber := uint64(10)
	newBlockHash := common.HexToHash(
		"0x0000000000000000000000000000000000000000000000000000000000000000",
	)

	event := types.Log{
		Address:     common.HexToAddress("0x123"),
		BlockNumber: newBlockNumber,
		BlockHash:   newBlockHash,
	}

	return channel, reorgChannel, ctx, cancel, event, newBlockNumber, newBlockHash
}

func TestIndexLogsSuccess(t *testing.T) {
	channel, reorgChannel, ctx, cancel, event, newBlockNumber, newBlockHash := setup()
	defer func() {
		cancel()
		close(channel)
		close(reorgChannel)
	}()

	channel <- event

	mockClient := blockchainMocks.NewMockChainClient(t)

	var wg sync.WaitGroup
	wg.Add(2) // Expecting two calls: StoreLog and UpdateLatestBlock

	contract := indexerMocks.NewMockIContract(t)

	contract.EXPECT().
		Logger().
		Return(zap.NewNop())

	contract.EXPECT().
		Address().
		Return(common.HexToAddress("0x123"))

	contract.EXPECT().
		UpdateLatestBlock(mock.Anything, newBlockNumber, newBlockHash.Bytes()).
		Run(func(ctx context.Context, blockNum uint64, blockHash []byte) {
			wg.Done()
		}).
		Return(nil)

	contract.EXPECT().
		StoreLog(mock.Anything, event).
		Run(func(ctx context.Context, log types.Log) {
			wg.Done()
		}).
		Return(nil)

	go c.IndexLogs(
		ctx,
		mockClient,
		channel,
		reorgChannel,
		contract,
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
	channel, reorgChannel, ctx, cancel, event, _, _ := setup()
	defer func() {
		cancel()
		close(channel)
		close(reorgChannel)
	}()

	var wg sync.WaitGroup
	wg.Add(2) // Expecting two calls: StoreLog and UpdateLatestBlock

	// Will fail for the first call with a retryable error and a non-retryable error on the second call
	attemptNumber := 0

	mockClient := blockchainMocks.NewMockChainClient(t)

	contract := indexerMocks.NewMockIContract(t)

	contract.EXPECT().
		Logger().
		Return(zap.NewNop())

	contract.EXPECT().
		Address().
		Return(common.HexToAddress("0x123"))

	contract.EXPECT().
		StoreLog(mock.Anything, event).
		RunAndReturn(func(ctx context.Context, log types.Log) errorMocks.RetryableError {
			wg.Done()
			attemptNumber++
			if attemptNumber < 2 {
				return errorMocks.NewRecoverableError("retryable error",
					errors.New("retryable error"),
				)
			} else {
				return errorMocks.NewNonRecoverableError("non-retryable error",
					errors.New("non-retryable error"),
				)
			}
		})

	channel <- event

	go c.IndexLogs(
		ctx,
		mockClient,
		channel,
		reorgChannel,
		contract,
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

	contract.AssertNumberOfCalls(t, "StoreLog", 2)
}
