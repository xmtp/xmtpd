package blockchain

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// Construct a raw blockchain listener that can be used to listen for events across many contract event types
type ChainStreamerBuilder interface {
	ListenForContractEvent(fromBlock uint64, contractAddress common.Address, topic common.Hash) <-chan types.Log
	Build() ChainStreamer
}

type ChainStreamer interface {
	Start(ctx context.Context) error
	Stop() error
}
