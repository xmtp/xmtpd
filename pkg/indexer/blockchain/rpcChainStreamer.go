package blockchain

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

/*
*
A RpcChainStreamer is a naive implementation of the ChainStreamer interface.
It queries a remote blockchain node for log events to backfill history, and then streams new events,
to get a complete history of events on a chain.
*
*/
type RpcChainStreamer struct {
	rpcUrl    string
	listeners []contractEventListener
}

func (a *RpcChainStreamer) Start(ctx context.Context) error {
	return nil
}

func (a *RpcChainStreamer) Stop() error {
	return nil
}

type contractEventListener struct {
	contractAddress common.Address
	topic           common.Hash
	channel         <-chan types.Log
}

// The builder that allows you to configure contract events to listen for
type RpcChainStreamerBuilder struct {
	// All the listeners
	listeners []contractEventListener
	rpcUrl    string
}

func NewRpcChainStreamBuilder(rpcUrl string) *RpcChainStreamerBuilder {
	return &RpcChainStreamerBuilder{rpcUrl: rpcUrl}
}

func (c *RpcChainStreamerBuilder) ListenForContractEvent(fromBlock uint64, contractAddress common.Address, topic common.Hash) <-chan types.Log {
	eventChannel := make(chan types.Log)
	c.listeners = append(c.listeners, contractEventListener{contractAddress, topic, eventChannel})
	return eventChannel
}

func (c *RpcChainStreamerBuilder) Build() *RpcChainStreamer {
	return &RpcChainStreamer{}
}
