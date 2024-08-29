package blockchain

import (
	"context"

	"github.com/ethereum/go-ethereum/ethclient"
)

func NewClient(ctx context.Context, rpcUrl string) (*ethclient.Client, error) {
	return ethclient.DialContext(ctx, rpcUrl)
}
