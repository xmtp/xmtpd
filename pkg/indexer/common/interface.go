package common

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	re "github.com/xmtp/xmtpd/pkg/errors"
	"go.uber.org/zap"
)

// Takes a log event and stores it, returning either an error that may be retriable, non-retriable, or nil.
type ILogStorer interface {
	StoreLog(ctx context.Context, event types.Log) re.RetryableError
}

// Tracks the latest block number and hash for a contract.
type IBlockTracker interface {
	GetLatestBlock() (uint64, []byte)
	UpdateLatestBlock(ctx context.Context, block uint64, hash []byte) error
}

type IReorgHandler interface {
	HandleLog(ctx context.Context, event types.Log) re.RetryableError
}

// An IContract is a contract that can be indexed.
type IContract interface {
	IBlockTracker
	IReorgHandler
	ILogStorer
	Address() common.Address
	Topics() []common.Hash
	Logger() *zap.Logger
}
