package blockchain

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	acg "github.com/xmtp/xmtpd/pkg/abi/appchaingateway"
	scg "github.com/xmtp/xmtpd/pkg/abi/settlementchaingateway"

	"github.com/ethereum/go-ethereum/core/types"
	mft "github.com/xmtp/xmtpd/pkg/abi/mockunderlyingfeetoken"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/xmtp/xmtpd/pkg/abi/erc20"
	ft "github.com/xmtp/xmtpd/pkg/abi/feetoken"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/xmtp/xmtpd/pkg/config"
	"go.uber.org/zap"
)

type IFundsAdmin interface {
	MintMockUSDC(ctx context.Context, addr common.Address, amount *big.Int) error
	Balances(ctx context.Context) error
	Deposit(ctx context.Context, amount *big.Int, gasLimit *big.Int, gasPrice *big.Int) error
	Withdraw(ctx context.Context, amount *big.Int) error
	ReceiveWithdrawal(ctx context.Context) error
}

type fundsAdmin struct {
	logger                 *zap.Logger
	app                    FundsAdminAppOpts
	settlement             FundsAdminSettlementOpts
	feeToken               *ft.FeeToken
	underlying             *erc20.ERC20
	settlementGateway      *scg.SettlementChainGateway
	appGateway             *acg.AppChainGateway
	mockUnderlyingFeeToken *mft.MockUnderlyingFeeToken
	spender                common.Address
	appChainId             int64
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

	mockToken, err := mft.NewMockUnderlyingFeeToken(
		common.HexToAddress(opts.ContractOptions.SettlementChain.UnderlyingFeeToken),
		opts.Settlement.Client)
	if err != nil {
		return nil, err
	}

	settlementGateway, err := scg.NewSettlementChainGateway(
		common.HexToAddress(opts.ContractOptions.SettlementChain.GatewayAddress),
		opts.Settlement.Client)
	if err != nil {
		return nil, err
	}

	appGateway, err := acg.NewAppChainGateway(
		common.HexToAddress(opts.ContractOptions.AppChain.GatewayAddress),
		opts.App.Client)
	if err != nil {
		return nil, err
	}

	return &fundsAdmin{
		logger:                 opts.Logger.Named("FundsAdmin"),
		app:                    opts.App,
		settlement:             opts.Settlement,
		feeToken:               feeToken,
		underlying:             underlying,
		mockUnderlyingFeeToken: mockToken,
		spender: common.HexToAddress(
			opts.ContractOptions.SettlementChain.GatewayAddress,
		),
		settlementGateway: settlementGateway,
		appGateway:        appGateway,
		appChainId:        opts.ContractOptions.AppChain.ChainID,
	}, nil
}

func (f *fundsAdmin) Balances(ctx context.Context) error {
	address := f.settlement.Signer.FromAddress()

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

func (f *fundsAdmin) MintMockUSDC(
	ctx context.Context,
	addr common.Address,
	amount *big.Int,
) error {
	// sanity check
	if amount.Sign() == -1 {
		return fmt.Errorf("amount must be positive")
	}
	if amount.Cmp(big.NewInt(10000000000)) > 0 {
		return fmt.Errorf("amount must be less than 10000 mxUSDC")
	}

	err := ExecuteTransaction(
		ctx,
		f.settlement.Signer, f.logger, f.settlement.Client,
		func(opts *bind.TransactOpts) (*types.Transaction, error) {
			return f.mockUnderlyingFeeToken.Mint(opts, addr, amount)
		},
		func(log *types.Log) (interface{}, error) {
			return f.mockUnderlyingFeeToken.ParseTransfer(*log)
		},
		func(event interface{}) {
			transfer, ok := event.(*mft.MockUnderlyingFeeTokenTransfer)
			if !ok {
				f.logger.Error("unexpected event type, not MockUnderlyingFeeTokenTransfer")
				return
			}
			f.logger.Info(
				"tokens minted",
				zap.String("from", transfer.From.Hex()),
				zap.String("to", transfer.To.Hex()),
				zap.String("amount", transfer.Value.String()),
			)
		},
	)
	if err != nil {
		if strings.Contains(err.Error(), "FiatToken") {
			return fmt.Errorf("not a XMTP mock USDC token: %s", err.Error())
		}
		return err
	}
	return nil
}

func (f *fundsAdmin) Deposit(
	ctx context.Context,
	amount *big.Int,
	gasLimit *big.Int,
	gasPrice *big.Int,
) error {
	from := f.settlement.Signer.FromAddress()

	feeTokenBalance, err := f.feeToken.BalanceOf(&bind.CallOpts{Context: ctx}, from)
	if err != nil {
		return err
	}
	f.logger.Info("Current balance", zap.String("raw", feeTokenBalance.String()))

	if feeTokenBalance.Cmp(amount) < 0 {
		return fmt.Errorf(
			"insufficient balance: need %s tokens, have %s tokens",
			amount.String(),
			feeTokenBalance.String(),
		)
	}

	f.logger.Info("Executing bridge transaction…")

	xmtpChainId := big.NewInt(f.appChainId)

	err = ExecuteTransaction(
		ctx,
		f.settlement.Signer,
		f.logger,
		f.settlement.Client,
		func(opts *bind.TransactOpts) (*types.Transaction, error) {
			return f.settlementGateway.Deposit(
				opts,
				xmtpChainId,
				from,
				amount,
				gasLimit,
				gasPrice,
			)
		},
		func(log *types.Log) (interface{}, error) {
			return f.settlementGateway.ParseDeposit(*log)
		},
		func(event interface{}) {
			deposit, ok := event.(*scg.SettlementChainGatewayDeposit)
			if !ok {
				f.logger.Error("node added event is not of type SettlementChainGatewayDeposit")
				return
			}
			f.logger.Info("deposited", zap.String("amount", deposit.Amount.String()))
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func (f *fundsAdmin) Withdraw(
	ctx context.Context,
	amount *big.Int,
) error {
	from := f.settlement.Signer.FromAddress()

	appGasBalance, err := f.app.Client.BalanceAt(ctx, from, nil)
	if err != nil {
		return err
	}
	f.logger.Info("Current balance", zap.String("raw", appGasBalance.String()))

	if appGasBalance.Cmp(amount) < 0 {
		return fmt.Errorf(
			"insufficient balance: need %s tokens, have %s tokens",
			amount.String(),
			appGasBalance.String(),
		)
	}

	f.logger.Info("Executing bridge transaction…")

	err = ExecuteTransaction(
		ctx,
		f.app.Signer,
		f.logger,
		f.app.Client,
		func(opts *bind.TransactOpts) (*types.Transaction, error) {
			opts.Value = amount
			return f.appGateway.Withdraw(
				opts,
				from,
			)
		},
		func(log *types.Log) (interface{}, error) {
			return f.appGateway.ParseWithdrawal(*log)
		},
		func(event interface{}) {
			withdrawal, ok := event.(*acg.AppChainGatewayWithdrawal)
			if !ok {
				f.logger.Error("node added event is not of type AppChainGatewayWithdrawal")
				return
			}
			f.logger.Info(
				"withdrawn",
				zap.String("amount", withdrawal.Amount.String()),
				zap.String("recipient", withdrawal.Recipient.Hex()),
			)
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func (f *fundsAdmin) ReceiveWithdrawal(
	ctx context.Context,
) error {
	from := f.settlement.Signer.FromAddress()

	f.logger.Info("Executing bridge transaction…")

	err := ExecuteTransaction(
		ctx,
		f.settlement.Signer,
		f.logger,
		f.settlement.Client,
		func(opts *bind.TransactOpts) (*types.Transaction, error) {
			return f.settlementGateway.ReceiveWithdrawal(
				opts,
				from,
			)
		},
		func(log *types.Log) (interface{}, error) {
			return f.settlementGateway.ParseWithdrawalReceived(*log)
		},
		func(event interface{}) {
			withdrawal, ok := event.(*scg.SettlementChainGatewayWithdrawalReceived)
			if !ok {
				f.logger.Error(
					"node added event is not of type SettlementChainGatewayWithdrawalReceived",
				)
				return
			}
			f.logger.Info(
				"withdrawn",
				zap.String("amount", withdrawal.Amount.String()),
				zap.String("recipient", withdrawal.Recipient.Hex()),
			)
		},
	)
	if err != nil {
		return err
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
