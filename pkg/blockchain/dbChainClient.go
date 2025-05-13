package blockchain

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
)

type DatabaseChainClient struct {
}

func NewDatabaseChainClient() *DatabaseChainClient {
	return &DatabaseChainClient{}
}

func (d *DatabaseChainClient) BlockNumber(ctx context.Context) (uint64, error) {
	return 0, fmt.Errorf("BlockNumber not implemented")
}

func (d *DatabaseChainClient) FilterLogs(
	ctx context.Context,
	q ethereum.FilterQuery,
) ([]types.Log, error) {
	return nil, fmt.Errorf("FilterLogs not implemented")
}

func (d *DatabaseChainClient) SubscribeFilterLogs(
	ctx context.Context,
	q ethereum.FilterQuery,
	ch chan<- types.Log,
) (ethereum.Subscription, error) {
	return nil, fmt.Errorf("SubscribeFilterLogs not implemented")
}

func (d *DatabaseChainClient) ChainID(ctx context.Context) (*big.Int, error) {
	return nil, fmt.Errorf("ChainID not implemented")
}

func (d *DatabaseChainClient) BlockByNumber(
	ctx context.Context,
	number *big.Int,
) (*types.Block, error) {
	return nil, fmt.Errorf("BlockByNumber not implemented")
}
