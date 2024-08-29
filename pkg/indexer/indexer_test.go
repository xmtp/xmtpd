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
	"github.com/xmtp/xmtpd/pkg/mocks"
	"github.com/xmtp/xmtpd/pkg/testutils"
)

func TestIndexLogsSuccess(t *testing.T) {
	channel := make(chan types.Log, 10)
	defer close(channel)

	logStorer := mocks.NewMockLogStorer(t)

	event := types.Log{
		Address: common.HexToAddress("0x123"),
	}
	logStorer.EXPECT().StoreLog(mock.Anything, event).Times(1).Return(nil)
	channel <- event

	go indexLogs(context.Background(), channel, testutils.NewLog(t), logStorer)
	time.Sleep(100 * time.Millisecond)
}

func TestIndexLogsRetryableError(t *testing.T) {
	channel := make(chan types.Log, 10)
	defer close(channel)

	logStorer := mocks.NewMockLogStorer(t)

	event := types.Log{
		Address: common.HexToAddress("0x123"),
	}

	// Will fail for the first call with a retryable error and a non-retryable error on the second call
	attemptNumber := 0

	logStorer.EXPECT().
		StoreLog(mock.Anything, event).
		RunAndReturn(func(ctx context.Context, log types.Log) storer.LogStorageError {
			attemptNumber++
			return storer.NewLogStorageError(errors.New("retryable error"), attemptNumber < 2)
		})
	channel <- event

	go indexLogs(context.Background(), channel, testutils.NewLog(t), logStorer)
	time.Sleep(200 * time.Millisecond)

	logStorer.AssertNumberOfCalls(t, "StoreLog", 2)
}
