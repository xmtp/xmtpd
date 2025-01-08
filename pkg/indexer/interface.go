package indexer

import (
	"context"
)

type IBlockTracker interface {
	GetLatestBlockNumber() uint64
	GetLatestBlockHash() []byte
	UpdateLatestBlock(ctx context.Context, block uint64, hash []byte) error
}
