package blockchain

import (
	"context"
	"errors"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/xmtp/xmtpd/pkg/abi/rateregistry"
	"github.com/xmtp/xmtpd/pkg/config"
	"github.com/xmtp/xmtpd/pkg/fees"
	"go.uber.org/zap"
)

/*
*
A RatesAdmin is a struct responsible for calling admin functions on the RatesRegistry contract
*
*/
type RatesAdmin struct {
	logger        *zap.Logger
	paramAdmin    *ParameterAdmin
	ratesContract *rateregistry.RateRegistry
}

func NewRatesAdmin(
	logger *zap.Logger,
	client *ethclient.Client,
	signer TransactionSigner,
	contractsOptions config.ContractsOptions,
) (*RatesAdmin, error) {
	paramAdmin, err := NewSettlementParameterAdmin(logger, client, signer, contractsOptions)
	if err != nil {
		return nil, err
	}

	rateContract, err := rateregistry.NewRateRegistry(
		common.HexToAddress(contractsOptions.SettlementChain.RateRegistryAddress),
		client,
	)
	if err != nil {
		return nil, err
	}

	return &RatesAdmin{
		logger:        logger.Named("RatesAdmin"),
		paramAdmin:    paramAdmin,
		ratesContract: rateContract,
	}, nil
}

/**
*
* AddRates adds a new rate to the rates manager.
* The new rate must have a later start time than the last rate in the contract.
 */
func (r *RatesAdmin) AddRates(ctx context.Context, rates fees.Rates) ProtocolError {
	// validations
	if rates.MessageFee < 0 {
		return NewBlockchainError(errors.New("rates.messageFee must be positive"))
	}
	if rates.StorageFee < 0 {
		return NewBlockchainError(errors.New("rates.storageFee must be positive"))
	}
	if rates.CongestionFee < 0 {
		return NewBlockchainError(errors.New("rates.congestionFee must be positive"))
	}

	params := []Uint64Param{
		{Name: RATE_REGISTRY_MESSAGE_FEE_KEY, Value: uint64(rates.MessageFee)},
		{Name: RATE_REGISTRY_STORAGE_FEE_KEY, Value: uint64(rates.StorageFee)},
		{Name: RATE_REGISTRY_CONGESTION_FEE_KEY, Value: uint64(rates.CongestionFee)},
		{Name: RATE_REGISTRY_TARGET_RATE_PER_MINUTE_KEY, Value: rates.TargetRatePerMinute},
	}
	if err := r.paramAdmin.SetManyUint64Parameters(ctx, params); err != nil {
		return err
	}

	err := ExecuteTransaction(
		ctx,
		r.paramAdmin.signer,
		r.logger,
		r.paramAdmin.client,
		func(opts *bind.TransactOpts) (*types.Transaction, error) {
			return r.ratesContract.UpdateRates(opts)
		},
		func(log *types.Log) (interface{}, error) {
			return r.ratesContract.ParseRatesUpdated(*log)
		},
		func(event interface{}) {
			e, ok := event.(*rateregistry.RateRegistryRatesUpdated)
			if !ok {
				r.logger.Error("unexpected event type, not of type RateRegistryRatesUpdated")
				return
			}
			r.logger.Info("rates updated",
				zap.Uint64("messageFee", e.MessageFee),
				zap.Uint64("storageFee", e.StorageFee),
				zap.Uint64("congestionFee", e.CongestionFee),
				zap.Uint64("targetRatePerMinute", e.TargetRatePerMinute),
			)
		},
	)
	if err != nil {
		if err.IsNoChange() {
			r.logger.Info("No update needed",
				zap.Uint64("messageFee", uint64(rates.MessageFee)),
				zap.Uint64("storageFee", uint64(rates.StorageFee)),
				zap.Uint64("congestionFee", uint64(rates.CongestionFee)),
				zap.Uint64("targetRatePerMinute", rates.TargetRatePerMinute),
			)
			return nil
		}
		r.logger.Error("Protocol error", zap.String("error", err.Error()))
		return err
	}

	return nil
}

func (r *RatesAdmin) GetParamAdmin() *ParameterAdmin {
	return r.paramAdmin
}