package storer

import (
	"context"

	"github.com/ethereum/go-ethereum/core/types"
)

// Takes a log event and stores it, returning either an error that may be retriable, non-retriable, or nil
// appendLog should be true if the log is part of a reorg; invalidates the old lod and appends the new one
type LogStorer interface {
	StoreLog(ctx context.Context, event types.Log, appendLog bool) LogStorageError
}
