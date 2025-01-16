package storer

import (
	"context"

	"github.com/ethereum/go-ethereum/core/types"
)

// Takes a log event and stores it, returning either an error that may be retriable, non-retriable, or nil
type LogStorer interface {
	StoreLog(ctx context.Context, event types.Log, appendLog bool) LogStorageError
}
