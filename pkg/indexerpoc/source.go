package indexerpoc

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"go.uber.org/zap"
)

// GethSource is an implementation of the Source interface using go-ethereum.
type GethSource struct {
	client    *ethclient.Client
	rpcClient *rpc.Client
	network   *Network
	log       *zap.Logger

	// Cache for latest block
	cacheMu      sync.Mutex
	latestNumber uint64
	latestHash   common.Hash
	cacheHits    int
	maxCacheHits int
}

// NewGethSource creates a new Ethereum client source.
func NewGethSource(
	network *Network,
	log *zap.Logger,
) (*GethSource, error) {
	if network == nil {
		return nil, fmt.Errorf("network configuration cannot be nil")
	}

	// Create a network-specific named logger
	networkLogger := log.Named("network").With(zap.String("network", network.Name))

	rpcClient, err := rpc.Dial(network.RpcURL)
	if err != nil {
		return nil, fmt.Errorf("connecting to %s node: %w", network.Name, err)
	}

	client := ethclient.NewClient(rpcClient)

	// Verify chain ID matches the expected one
	chainID, err := client.ChainID(context.Background())
	if err != nil {
		return nil, fmt.Errorf("getting chain ID from %s: %w", network.Name, err)
	}

	if chainID.Int64() != network.ChainID {
		return nil, fmt.Errorf("chain ID mismatch for %s: expected %d, got %d",
			network.Name, network.ChainID, chainID.Int64())
	}

	source := &GethSource{
		client:       client,
		rpcClient:    rpcClient,
		network:      network,
		log:          networkLogger,
		maxCacheHits: 20,
	}

	// Start a background goroutine to keep track of the latest block
	go source.pollLatestBlock(context.Background())

	return source, nil
}

// GetNetworkName returns the name of the network
func (s *GethSource) GetNetworkName() string {
	return s.network.Name
}

// GetChainID returns the chain ID of the network
func (s *GethSource) GetChainID() int64 {
	return s.network.ChainID
}

// pollLatestBlock periodically polls for the latest block
func (s *GethSource) pollLatestBlock(ctx context.Context) {
	ticker := time.NewTicker(s.network.PollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			header, err := s.client.HeaderByNumber(ctx, nil) // nil = latest block
			if err != nil {
				s.log.Error("Failed to poll latest block",
					zap.Error(err))
				continue
			}

			s.cacheMu.Lock()
			s.latestNumber = header.Number.Uint64()
			s.latestHash = header.Hash()
			s.cacheHits = 0 // Reset cache hits when we update
			s.cacheMu.Unlock()

			s.log.Debug("Updated latest block",
				zap.Uint64("number", header.Number.Uint64()),
				zap.String("hash", header.Hash().Hex()),
			)
		}
	}
}

// GetLatestBlockNumber retrieves the current latest block number
func (s *GethSource) GetLatestBlockNumber(ctx context.Context) (uint64, error) {
	s.cacheMu.Lock()

	// If we have a cached value and haven't reached max hits, use it
	if s.latestNumber > 0 && s.cacheHits < s.maxCacheHits {
		s.cacheHits++
		latestNumber := s.latestNumber
		s.cacheMu.Unlock()
		return latestNumber, nil
	}

	s.cacheMu.Unlock()

	// Need to fetch from network
	header, err := s.client.HeaderByNumber(ctx, nil) // nil = latest block
	if err != nil {
		return 0, fmt.Errorf("getting latest block header on %s: %w", s.network.Name, err)
	}

	blockNumber := header.Number.Uint64()

	// Update cache
	s.cacheMu.Lock()
	s.latestNumber = blockNumber
	s.latestHash = header.Hash()
	s.cacheHits = 1
	s.cacheMu.Unlock()

	return blockNumber, nil
}

// GetBlockHash retrieves the hash of a specific block
func (s *GethSource) GetBlockHash(ctx context.Context, number uint64) (common.Hash, error) {
	// Check cache first if this is the latest block
	s.cacheMu.Lock()
	if s.latestNumber == number && s.latestHash != (common.Hash{}) {
		hash := s.latestHash
		s.cacheMu.Unlock()
		return hash, nil
	}
	s.cacheMu.Unlock()

	// Otherwise fetch the block header
	header, err := s.client.HeaderByNumber(ctx, big.NewInt(int64(number)))
	if err != nil {
		return common.Hash{}, fmt.Errorf("getting block header %d on %s: %w",
			number, s.network.Name, err)
	}

	return header.Hash(), nil
}

// GetBlockByNumber retrieves a block by its number
func (s *GethSource) GetBlockByNumber(
	ctx context.Context,
	number uint64,
) (*types.Block, error) {
	block, err := s.client.BlockByNumber(ctx, big.NewInt(int64(number)))
	if err != nil {
		return nil, fmt.Errorf("getting block %d on %s: %w",
			number, s.network.Name, err)
	}

	return block, nil
}

// GetLogs retrieves logs based on the filter parameters
func (s *GethSource) GetLogs(
	ctx context.Context,
	startBlock, endBlock uint64,
	filter *Filter,
) ([]types.Log, error) {
	// Create an Ethereum filter query
	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(int64(startBlock)),
		ToBlock:   big.NewInt(int64(endBlock)),
		Addresses: filter.Addresses,
		Topics:    [][]common.Hash{},
	}

	// Add topics if we have them
	if len(filter.Topics) > 0 {
		query.Topics = make([][]common.Hash, len(filter.Topics))
		for i, topicList := range filter.Topics {
			query.Topics[i] = make([]common.Hash, len(topicList))
			copy(query.Topics[i], topicList)
		}
	}

	// Get logs
	logs, err := s.client.FilterLogs(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("filtering logs from %d to %d on %s: %w",
			startBlock, endBlock, s.network.Name, err)
	}

	return logs, nil
}
