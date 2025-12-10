package blockchain

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/xmtp/xmtpd/pkg/abi/rateregistry"
	"github.com/xmtp/xmtpd/pkg/config"
	"github.com/xmtp/xmtpd/pkg/fees"
	"github.com/xmtp/xmtpd/pkg/utils"
	"go.uber.org/zap"
)

type IRatesAdmin interface {
	AddRates(ctx context.Context, rates fees.Rates) ProtocolError
}

type RatesAdmin struct {
	logger        *zap.Logger
	client        *ethclient.Client
	signer        TransactionSigner
	paramAdmin    IParameterAdmin
	ratesContract *rateregistry.RateRegistry
}

var _ IRatesAdmin = &RatesAdmin{}

func NewRatesAdmin(
	logger *zap.Logger,
	client *ethclient.Client,
	signer TransactionSigner,
	paramAdmin IParameterAdmin,
	contractsOptions *config.ContractsOptions,
) (IRatesAdmin, error) {
	rateContract, err := rateregistry.NewRateRegistry(
		common.HexToAddress(contractsOptions.SettlementChain.RateRegistryAddress),
		client,
	)
	if err != nil {
		return nil, err
	}

	rateRegistryAdminLogger := logger.Named(utils.RateRegistryAdminLoggerName).With(
		utils.SettlementChainChainIDField(contractsOptions.SettlementChain.ChainID),
	)

	return &RatesAdmin{
		logger:        rateRegistryAdminLogger,
		paramAdmin:    paramAdmin,
		ratesContract: rateContract,
		client:        client,
		signer:        signer,
	}, nil
}

// AddRates adds a new rate to the rates manager.
// The new rate must have a later start time than the last rate in the contract.
func (r *RatesAdmin) AddRates(ctx context.Context, rates fees.Rates) ProtocolError {
	// validations
	if rates.MessageFee < 0 {
		return NewBlockchainError(
			fmt.Errorf("%s must be positive", RateRegistryMessageFeeKey),
		)
	}
	if rates.StorageFee < 0 {
		return NewBlockchainError(
			fmt.Errorf("%s must be positive", RateRegistryStorageFeeKey),
		)
	}
	if rates.CongestionFee < 0 {
		return NewBlockchainError(
			fmt.Errorf("%s must be positive", RateRegistryCongestionFeeKey),
		)
	}

	params := []Uint64Param{
		{Name: RateRegistryMessageFeeKey, Value: uint64(rates.MessageFee)},
		{Name: RateRegistryStorageFeeKey, Value: uint64(rates.StorageFee)},
		{Name: RateRegistryCongestionFeeKey, Value: uint64(rates.CongestionFee)},
		{Name: RateRegistryTargetRatePerMinuteKey, Value: rates.TargetRatePerMinute},
	}
	if err := r.paramAdmin.SetManyUint64Parameters(ctx, params); err != nil {
		return err
	}

	err := ExecuteTransaction(
		ctx,
		r.signer,
		r.logger,
		r.client,
		func(opts *bind.TransactOpts) (*types.Transaction, error) {
			return r.ratesContract.UpdateRates(opts)
		},
		func(log *types.Log) (any, error) {
			return r.ratesContract.ParseRatesUpdated(*log)
		},
		func(event any) {
			e, ok := event.(*rateregistry.RateRegistryRatesUpdated)
			if !ok {
				r.logger.Error("unexpected event type, not of type RateRegistryRatesUpdated")
				return
			}
			r.logger.Info("rates updated",
				zap.Uint64(RateRegistryMessageFeeKey, e.MessageFee),
				zap.Uint64(RateRegistryStorageFeeKey, e.StorageFee),
				zap.Uint64(RateRegistryCongestionFeeKey, e.CongestionFee),
				zap.Uint64(RateRegistryTargetRatePerMinuteKey, e.TargetRatePerMinute),
			)
		},
	)
	if err != nil {
		if err.IsNoChange() {
			r.logger.Info("no update needed",
				zap.Uint64(RateRegistryMessageFeeKey, uint64(rates.MessageFee)),
				zap.Uint64(RateRegistryStorageFeeKey, uint64(rates.StorageFee)),
				zap.Uint64(RateRegistryCongestionFeeKey, uint64(rates.CongestionFee)),
				zap.Uint64(RateRegistryTargetRatePerMinuteKey, rates.TargetRatePerMinute),
			)
			return nil
		}
		r.logger.Error("protocol error", zap.String("error", err.Error()))
		return err
	}

	return nil
}
