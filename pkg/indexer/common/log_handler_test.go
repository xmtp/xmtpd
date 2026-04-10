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
	c "github.com/xmtp/xmtpd/pkg/indexer/common"
	indexerMocks "github.com/xmtp/xmtpd/pkg/testutils/mocks/common"
	errorMocks "github.com/xmtp/xmtpd/pkg/utils/retryerrors"
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

	waitForWaitGroup(t, &wg, 10*time.Second,
		"timed out waiting for StoreLog and UpdateLatestBlock")
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

	waitForWaitGroup(t, &wg, 10*time.Second,
		"timed out waiting for StoreLog and UpdateLatestBlock")

	contract.AssertNumberOfCalls(t, "StoreLog", 2)
}

// waitForWaitGroup blocks until wg signals done or the timeout elapses.
// The helper goroutine is tracked with t.Cleanup so it can never outlive the
// test function, even if the timeout path is taken.
func waitForWaitGroup(t *testing.T, wg *sync.WaitGroup, timeout time.Duration, failMsg string) {
	t.Helper()

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	t.Cleanup(func() {
		// Drain the helper goroutine before the test returns. If wg is still
		// outstanding we can't cleanly stop it, but at least the test already
		// failed via t.Fatal below so the process will exit shortly.
		select {
		case <-done:
		default:
		}
	})

	select {
	case <-done:
		// success
	case <-time.After(timeout):
		t.Fatal(failMsg)
	}
}
