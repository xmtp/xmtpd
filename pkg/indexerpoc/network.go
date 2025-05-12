package indexerpoc

import (
	"time"
)

const DefaultPollInterval = 1 * time.Second

// Network represents an EVM-compatible blockchain network.
type Network struct {
	Name         string        // Network name (e.g., "ethereum", "arbitrum", "optimism")
	ChainID      int64         // Chain ID of the network
	RpcURL       string        // RPC endpoint URL
	PollInterval time.Duration // Poll interval for this network (optional, defaults to global)
}

// NewNetwork creates a new network configuration.
func NewNetwork(name string, chainID int64, rpcURL string, pollInterval time.Duration) *Network {
	if pollInterval == 0 {
		pollInterval = DefaultPollInterval
	}

	return &Network{
		Name:         name,
		ChainID:      chainID,
		RpcURL:       rpcURL,
		PollInterval: pollInterval,
	}
}
