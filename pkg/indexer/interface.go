package indexer

import (
	"context"
)

type IBlockTracker interface {
	GetLatestBlock() (uint64, []byte)
	UpdateLatestBlock(ctx context.Context, block uint64, hash []byte) error
}
