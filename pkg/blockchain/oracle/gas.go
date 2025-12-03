package oracle

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type GasPriceSource int

const (
	gasPriceSourceArbMinimum GasPriceSource = iota
	gasPriceSourceDefault
)

func (g GasPriceSource) String() string {
	switch g {
	case gasPriceSourceArbMinimum:
		return "arbitrum"
	case gasPriceSourceDefault:
		return "default"
	default:
		return "unknown"
	}
}

var (
	// ArbGasInfo precompile address.
	arbGasInfoAddr = common.HexToAddress("0x000000000000000000000000000000000000006C")

	// getMinimumGasPrice() selector = keccak256("getMinimumGasPrice()")[:4].
	// Verified with: `cast sig "getMinimumGasPrice()"`.
	getMinimumGasPriceSelector = common.FromHex("0xf918379a")
)

func isArbChain(ctx context.Context, client *ethclient.Client) bool {
	code, err := client.CodeAt(ctx, arbGasInfoAddr, nil)
	if err != nil {
		return false
	}

	return len(code) > 0
}

func getGasPrice(
	ctx context.Context,
	client *ethclient.Client,
	gasPriceSource GasPriceSource,
) (*big.Int, error) {
	if gasPriceSource == gasPriceSourceArbMinimum {
		return getArbMinimumGasPrice(ctx, client)
	}

	return getDefaultGasPrice(ctx, client)
}

func getArbMinimumGasPrice(ctx context.Context, client *ethclient.Client) (*big.Int, error) {
	result, err := client.CallContract(ctx, ethereum.CallMsg{
		To:   &arbGasInfoAddr,
		Data: getMinimumGasPriceSelector,
	}, nil)
	if err != nil {
		return nil, err
	}

	return new(big.Int).SetBytes(result), nil
}

func getDefaultGasPrice(ctx context.Context, client *ethclient.Client) (*big.Int, error) {
	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		return nil, err
	}

	return gasPrice, nil
}
