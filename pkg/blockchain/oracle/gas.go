package oracle

import (
	"context"
	"errors"
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

// ArbGasInfo is a Nitro precompile.
// Source: https://github.com/OffchainLabs/nitro-precompile-interfaces/blob/main/ArbGasInfo.sol
var (
	// ArbGasInfo precompile address.
	arbGasInfoAddr = common.HexToAddress("0x000000000000000000000000000000000000006C")

	// getPricesInWei() selector = keccak256("getPricesInWei()")[:4].
	getPricesInWeiSelector = common.FromHex("0x41b247a8")
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
		price, err := getArbGasPrices(ctx, client)
		if err != nil {
			return nil, err
		}
		return price, nil
	}

	return getDefaultGasPrice(ctx, client)
}

func getArbGasPrices(ctx context.Context, client *ethclient.Client) (*big.Int, error) {
	result, err := client.CallContract(ctx, ethereum.CallMsg{
		To:   &arbGasInfoAddr,
		Data: getPricesInWeiSelector,
	}, nil)
	if err != nil {
		return nil, err
	}

	// Response is 6 consecutive uint256 (32 bytes each) = 192 bytes total.
	// Direct byte slicing is faster than abi.Arguments.Unpack() (no reflection).
	const wordSize = 32
	if len(result) < 6*wordSize {
		return nil, errors.New("unexpected response length from getPricesInWei")
	}

	// Last 32 bytes are the total gas price. Index 5 * 32 bytes.
	return new(big.Int).SetBytes(result[5*wordSize : 6*wordSize]), nil
}

func getDefaultGasPrice(ctx context.Context, client *ethclient.Client) (*big.Int, error) {
	return client.SuggestGasPrice(ctx)
}
