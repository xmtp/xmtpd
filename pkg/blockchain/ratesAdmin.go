package blockchain

import (
	"context"
	"encoding/binary"
	"errors"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/xmtp/xmtpd/pkg/abi/rateregistry"
	paramReg "github.com/xmtp/xmtpd/pkg/abi/settlementchainparameterregistry"
	"github.com/xmtp/xmtpd/pkg/config"
	"github.com/xmtp/xmtpd/pkg/fees"
	"go.uber.org/zap"
)

const (
	RATE_REGISTRY_MESSAGE_FEE_KEY            = "xmtp.rateRegistry.messageFee"
	RATE_REGISTRY_STORAGE_FEE_KEY            = "xmtp.rateRegistry.storageFee"
	RATE_REGISTRY_CONGESTION_FEE_KEY         = "xmtp.rateRegistry.congestionFee"
	RATE_REGISTRY_TARGET_RATE_PER_MINUTE_KEY = "xmtp.rateRegistry.targetRatePerMinute"
)

/*
*
A RatesAdmin is a struct responsible for calling admin functions on the RatesRegistry contract
*
*/
type RatesAdmin struct {
	client            *ethclient.Client
	signer            TransactionSigner
	parameterContract *paramReg.SettlementChainParameterRegistry
	ratesContract     *rateregistry.RateRegistry
	logger            *zap.Logger
}

func NewRatesAdmin(
	logger *zap.Logger,
	client *ethclient.Client,
	signer TransactionSigner,
	contractsOptions config.ContractsOptions,
) (*RatesAdmin, error) {
	rateContract, err := rateregistry.NewRateRegistry(
		common.HexToAddress(contractsOptions.SettlementChain.RateRegistryAddress),
		client,
	)
	if err != nil {
		return nil, err
	}

	parameterContract, err := paramReg.NewSettlementChainParameterRegistry(
		common.HexToAddress(contractsOptions.SettlementChain.ParameterRegistryAddress),
		client,
	)
	if err != nil {
		return nil, err
	}

	return &RatesAdmin{
		signer:            signer,
		client:            client,
		logger:            logger.Named("RatesAdmin"),
		parameterContract: parameterContract,
		ratesContract:     rateContract,
	}, nil
}

/**
*
* AddRates adds a new rate to the rates manager.
* The new rate must have a later start time than the last rate in the contract.
 */
func (r *RatesAdmin) AddRates(
	ctx context.Context,
	rates fees.Rates,
) error {
	err := ExecuteTransaction(
		ctx,
		r.signer,
		r.logger,
		r.client,
		func(opts *bind.TransactOpts) (*types.Transaction, error) {
			keys := [][]byte{
				[]byte(RATE_REGISTRY_MESSAGE_FEE_KEY),
				[]byte(RATE_REGISTRY_STORAGE_FEE_KEY),
				[]byte(RATE_REGISTRY_CONGESTION_FEE_KEY),
				[]byte(RATE_REGISTRY_TARGET_RATE_PER_MINUTE_KEY),
			}

			if rates.MessageFee < 0 {
				return nil, errors.New("rates.messageFee must be positive")
			}
			if rates.StorageFee < 0 {
				return nil, errors.New("rates.storageFee must be positive")
			}
			if rates.CongestionFee < 0 {
				return nil, errors.New("rates.congestionFee must be positive")
			}

			values := [][32]byte{
				encodeUint64ToBytes32(uint64(rates.MessageFee)),
				encodeUint64ToBytes32(uint64(rates.StorageFee)),
				encodeUint64ToBytes32(uint64(rates.CongestionFee)),
				encodeUint64ToBytes32(rates.TargetRatePerMinute),
			}

			return r.parameterContract.Set(opts, keys, values)
		},
		func(log *types.Log) (interface{}, error) {
			return r.parameterContract.ParseParameterSet(*log)
		},
		func(event interface{}) {
			parameterSet, ok := event.(*paramReg.SettlementChainParameterRegistryParameterSet)
			if !ok {
				r.logger.Error(
					"unexpected event type, not of type SettlementChainParameterRegistryParameterSet",
				)
				return
			}
			r.logger.Info("set parameter",
				zap.String("key", parameterSet.Key.String()),
				zap.Uint64("parameter", decodeBytes32ToUint64(parameterSet.Value)),
			)
		},
	)
	if err != nil {
		return err
	}

	err = ExecuteTransaction(
		ctx,
		r.signer,
		r.logger,
		r.client,
		func(opts *bind.TransactOpts) (*types.Transaction, error) {
			return r.ratesContract.UpdateRates(opts)
		},
		func(log *types.Log) (interface{}, error) {
			return r.ratesContract.ParseRatesUpdated(*log)
		},
		func(event interface{}) {
			rateUpdated, ok := event.(*rateregistry.RateRegistryRatesUpdated)
			if !ok {
				r.logger.Error(
					"unexpected event type, not of type RateRegistryRatesUpdated",
				)
				return
			}
			r.logger.Info("rates updated",
				zap.Uint64("messageFee", rateUpdated.MessageFee),
				zap.Uint64("storageFee", rateUpdated.StorageFee),
				zap.Uint64("congestionFee", rateUpdated.CongestionFee),
				zap.Uint64("targetRatePerMinute", rateUpdated.TargetRatePerMinute),
			)
		},
	)
	if err != nil {
		if strings.Contains(err.Error(), "NoChange") {
			r.logger.Info("No update needed",
				zap.Uint64("messageFee", uint64(rates.MessageFee)),
				zap.Uint64("storageFee", uint64(rates.StorageFee)),
				zap.Uint64("congestionFee", uint64(rates.CongestionFee)),
				zap.Uint64("targetRatePerMinute", rates.TargetRatePerMinute),
			)
			return nil
		}
		return err
	}

	return nil
}

func encodeUint64ToBytes32(v uint64) [32]byte {
	var b [32]byte
	binary.BigEndian.PutUint64(b[24:], v)
	return b
}

func decodeBytes32ToUint64(b [32]byte) uint64 {
	return binary.BigEndian.Uint64(b[24:]) // last 8 bytes
}
