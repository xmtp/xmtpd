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
	indexerMocks "github.com/xmtp/xmtpd/pkg/mocks/common"
	"go.uber.org/zap"
)

// TODO: Add more test coverage.

type logHandlerTest struct {
	eventChannel   chan types.Log
	ctx            context.Context
	cancel         context.CancelFunc
	event          types.Log
	newBlockNumber uint64
	newBlockHash   common.Hash
}

func setup(t *testing.T) logHandlerTest {
	ctx, cancel := context.WithCancel(t.Context())

	blockNumber := uint64(10)
	blockHash := common.HexToHash(
		"0x0000000000000000000000000000000000000000000000000000000000000000",
	)

	event := types.Log{
		Address:     common.HexToAddress("0x123"),
		BlockNumber: blockNumber,
		BlockHash:   blockHash,
	}

	return logHandlerTest{
		eventChannel:   make(chan types.Log, 10),
		ctx:            ctx,
		cancel:         cancel,
		event:          event,
		newBlockNumber: blockNumber,
		newBlockHash:   blockHash,
	}
}

func TestIndexLogsSuccess(t *testing.T) {
	cfg := setup(t)
	defer func() {
		cfg.cancel()
		close(cfg.eventChannel)
	}()

	cfg.eventChannel <- cfg.event

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
		UpdateLatestBlock(mock.Anything, cfg.newBlockNumber, cfg.newBlockHash.Bytes()).
		Run(func(ctx context.Context, blockNum uint64, blockHash []byte) {
			wg.Done()
		}).
		Return(nil)

	contract.EXPECT().
		StoreLog(mock.Anything, cfg.event).
		Run(func(ctx context.Context, log types.Log) {
			wg.Done()
		}).
		Return(nil)

	go c.IndexLogs(
		cfg.ctx,
		cfg.eventChannel,
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
	cfg := setup(t)
	defer func() {
		cfg.cancel()
		close(cfg.eventChannel)
	}()

	var wg sync.WaitGroup
	wg.Add(2) // Expecting two calls: StoreLog and UpdateLatestBlock

	// Will fail for the first call with a retryable error and a non-retryable error on the second call
	attemptNumber := 0

	contract := indexerMocks.NewMockIContract(t)

	contract.EXPECT().
		Logger().
		Return(zap.NewNop())

	contract.EXPECT().
		Address().
		Return(common.HexToAddress("0x123"))

	contract.EXPECT().
		StoreLog(mock.Anything, cfg.event).
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

	cfg.eventChannel <- cfg.event

	go c.IndexLogs(
		cfg.ctx,
		cfg.eventChannel,
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
