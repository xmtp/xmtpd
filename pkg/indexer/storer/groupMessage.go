package storer

import (
	"context"
	"errors"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"go.uber.org/zap"
)

type GroupMessageStorer struct {
	queries *queries.Queries
	logger  *zap.Logger
}

func NewGroupMessageStorer(queries *queries.Queries, logger *zap.Logger) *GroupMessageStorer {
	return &GroupMessageStorer{queries: queries, logger: logger}
}

// Validate and store a group message log event
func (s *GroupMessageStorer) StoreLog(ctx context.Context, event types.Log) LogStorageError {
	return NewLogStorageError(errors.New("not implemented"), true)
}
