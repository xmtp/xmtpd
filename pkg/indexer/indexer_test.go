package indexer

import (
	"context"
	"errors"
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

	blockTracker := indexerMocks.NewMockIBlockTracker(t)
	blockTracker.EXPECT().
		UpdateLatestBlock(mock.Anything, newBlockNumber, newBlockHash.Bytes()).
		Return(nil)

	reorgHandler := indexerMocks.NewMockChainReorgHandler(t)

	logStorer := storerMocks.NewMockLogStorer(t)
	logStorer.EXPECT().
		StoreLog(mock.Anything, event).
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
	)

	time.Sleep(100 * time.Millisecond)
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

	// Will fail for the first call with a retryable error and a non-retryable error on the second call
	attemptNumber := 0

	logStorer.EXPECT().
		StoreLog(mock.Anything, event).
		RunAndReturn(func(ctx context.Context, log types.Log) storer.LogStorageError {
			attemptNumber++
			return storer.NewLogStorageError(errors.New("retryable error"), attemptNumber < 2)
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
	)

	time.Sleep(200 * time.Millisecond)

	logStorer.AssertNumberOfCalls(t, "StoreLog", 2)
}
