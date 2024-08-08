package blockchain

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// The builder that allows you to configure contract events to listen for
type RpcChainStreamBuilder struct {
	// All the listeners
	listeners []contractEventListener
	rpcUrl    string
}

func NewRpcChainStreamBuilder(rpcUrl string) *RpcChainStreamBuilder {
	return &RpcChainStreamBuilder{rpcUrl: rpcUrl}
}

func (c *RpcChainStreamBuilder) ListenForContractEvent(fromBlock uint64, contractAddress common.Address, topic common.Hash) <-chan types.Log {
	eventChannel := make(chan types.Log)
	c.listeners = append(c.listeners, contractEventListener{contractAddress, topic, eventChannel})
	return eventChannel
}

func (c *RpcChainStreamBuilder) Build() *RpcChainStreamer {
	return &RpcChainStreamer{}
}
