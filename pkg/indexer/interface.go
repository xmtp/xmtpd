package indexer

import "context"

type IBlockTracker interface {
	GetLatestBlock() uint64
	UpdateLatestBlock(ctx context.Context, block uint64) error
}
