package blockchain

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/xmtp/xmtpd/pkg/abi/erc20"
	ft "github.com/xmtp/xmtpd/pkg/abi/feeToken"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/xmtp/xmtpd/pkg/config"
	"go.uber.org/zap"
)

type IFundsAdmin interface {
	Balances(ctx context.Context, address common.Address) error
}

type fundsAdmin struct {
	logger     *zap.Logger
	app        FundsAdminAppOpts
	settlement FundsAdminSettlementOpts
	feeToken   *ft.FeeToken
	underlying *erc20.ERC20
}

var _ IFundsAdmin = &fundsAdmin{}

type FundsAdminSettlementOpts struct {
	Client *ethclient.Client
	Signer TransactionSigner
}

type FundsAdminAppOpts struct {
	Client *ethclient.Client
	Signer TransactionSigner
}

type FundsAdminOpts struct {
	Logger          *zap.Logger
	ContractOptions config.ContractsOptions
	Settlement      FundsAdminSettlementOpts
	App             FundsAdminAppOpts
}

func NewFundsAdmin(
	opts FundsAdminOpts,
) (IFundsAdmin, error) {
	feeToken, err := ft.NewFeeToken(
		common.HexToAddress(opts.ContractOptions.SettlementChain.FeeToken),
		opts.Settlement.Client,
	)
	if err != nil {
		return nil, err
	}

	underlying, err := erc20.NewERC20(
		common.HexToAddress(opts.ContractOptions.SettlementChain.UnderlyingFeeToken),
		opts.Settlement.Client,
	)
	if err != nil {
		return nil, err
	}

	return &fundsAdmin{
		logger:     opts.Logger.Named("FundsAdmin"),
		app:        opts.App,
		settlement: opts.Settlement,
		feeToken:   feeToken,
		underlying: underlying,
	}, nil
}

func (f *fundsAdmin) Balances(ctx context.Context, address common.Address) error {
	ethBalance, err := f.settlement.Client.BalanceAt(ctx, address, nil)
	if err != nil {
		f.logger.Error("failed to get ETH balance", zap.Error(err))
	} else {
		f.logger.Info(
			"ETH balance of",
			zap.String("address", address.Hex()),
			zap.String("balance", FromWei(ethBalance, 18)),
		)
	}

	feeTokenBalance, err := f.feeToken.BalanceOf(&bind.CallOpts{Context: ctx}, address)
	if err != nil {
		f.logger.Error("failed to get xUSD balance", zap.Error(err))
	} else {
		f.logger.Info(
			"xUSD balance of", zap.String("address", address.Hex()),
			zap.String("balance", FromWei(feeTokenBalance, 6)),
		)
	}

	underlyingTokenBalance, err := f.underlying.BalanceOf(&bind.CallOpts{Context: ctx}, address)
	if err != nil {
		f.logger.Error("failed to get underlying USDC balance", zap.Error(err))
	} else {
		f.logger.Info(
			"USDC balance of", zap.String("address", address.Hex()),
			zap.String("balance", FromWei(underlyingTokenBalance, 6)),
		)
	}

	appGasBalance, err := f.app.Client.BalanceAt(ctx, address, nil)
	if err != nil {
		f.logger.Error("failed to get XMTP balance", zap.Error(err))
	} else {
		f.logger.Info(
			"XMTP balance of",
			zap.String("address", address.Hex()),
			zap.String("balance", FromWei(appGasBalance, 18)),
		)
	}

	return nil
}

// FromWei converts a wei value into a decimal string with the given decimals.
// For ETH, use decimals = 18.
// For an ERC20, use its `decimals()` value.
func FromWei(wei *big.Int, decimals int) string {
	// divisor = 10^decimals
	divisor := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil)

	intPart := new(big.Int).Div(wei, divisor)
	fracPart := new(big.Int).Mod(wei, divisor)

	// Left-pad fractional part with zeros
	fracStr := fracPart.String()
	for len(fracStr) < decimals {
		fracStr = "0" + fracStr
	}

	// Trim trailing zeros for nicer output
	fracStr = trimTrailingZeros(fracStr)

	if fracStr == "" {
		return intPart.String()
	}
	return intPart.String() + "." + fracStr
}

func trimTrailingZeros(s string) string {
	for len(s) > 0 && s[len(s)-1] == '0' {
		s = s[:len(s)-1]
	}
	return s
}
