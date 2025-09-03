package blockchain

import (
	"context"

	pr "github.com/xmtp/xmtpd/pkg/abi/payerregistry"

	scg "github.com/xmtp/xmtpd/pkg/abi/settlementchaingateway"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/xmtp/xmtpd/pkg/config"
	"go.uber.org/zap"
)

type ISettlementChainAdmin interface {
	GetSettlementChainGatewayPauseStatus(ctx context.Context) (bool, error)
	SetSettlementChainGatewayPauseStatus(ctx context.Context, paused bool) error
	GetPayerRegistryPauseStatus(ctx context.Context) (bool, error)
	SetPayerRegistryPauseStatus(ctx context.Context, paused bool) error
}

type settlementChainAdmin struct {
	client                 *ethclient.Client
	signer                 TransactionSigner
	logger                 *zap.Logger
	parameterAdmin         *ParameterAdmin
	settlementChainGateway *scg.SettlementChainGateway
	payerRegistry          *pr.PayerRegistry
}

var _ ISettlementChainAdmin = (*settlementChainAdmin)(nil)

func NewSettlementChainAdmin(
	logger *zap.Logger,
	client *ethclient.Client,
	signer TransactionSigner,
	contractsOptions config.ContractsOptions,
	parameterAdmin *ParameterAdmin,
) (*settlementChainAdmin, error) {
	acGateway, err := scg.NewSettlementChainGateway(
		common.HexToAddress(contractsOptions.SettlementChain.GatewayAddress),
		client,
	)
	if err != nil {
		return nil, err
	}

	payerRegistry, err := pr.NewPayerRegistry(
		common.HexToAddress(contractsOptions.SettlementChain.ParameterRegistryAddress),
		client,
	)
	if err != nil {
		return nil, err
	}

	return &settlementChainAdmin{
		client:                 client,
		signer:                 signer,
		logger:                 logger.Named("AppChainAdmin"),
		parameterAdmin:         parameterAdmin,
		settlementChainGateway: acGateway,
		payerRegistry:          payerRegistry,
	}, nil
}

func (s settlementChainAdmin) GetSettlementChainGatewayPauseStatus(
	ctx context.Context,
) (bool, error) {
	return s.parameterAdmin.GetParameterBool(ctx, SETTLEMENT_CHAIN_GATEWAY_PAUSED_KEY)
}

func (s settlementChainAdmin) SetSettlementChainGatewayPauseStatus(
	ctx context.Context,
	paused bool,
) error {
	if err := s.parameterAdmin.SetBoolParameter(ctx, SETTLEMENT_CHAIN_GATEWAY_PAUSED_KEY, paused); err != nil {
		return err
	}

	err := ExecuteTransaction(
		ctx,
		s.signer,
		s.logger,
		s.client,
		func(opts *bind.TransactOpts) (*types.Transaction, error) {
			return s.settlementChainGateway.UpdatePauseStatus(opts)
		},
		func(log *types.Log) (interface{}, error) {
			return s.settlementChainGateway.ParsePauseStatusUpdated(*log)
		},
		func(event interface{}) {
			ev, ok := event.(*scg.SettlementChainGatewayPauseStatusUpdated)
			if !ok {
				s.logger.Error(
					"unexpected event type, not of type SettlementChainGatewayPauseStatusUpdated",
				)
				return
			}
			s.logger.Info(
				"settlement-chain gateway pause status updated",
				zap.Bool("paused", ev.Paused),
			)
		},
	)
	if err != nil {
		if err.IsNoChange() {
			s.logger.Info("No update needed")
			return nil
		}
		return err
	}
	return nil
}

func (s settlementChainAdmin) GetPayerRegistryPauseStatus(ctx context.Context) (bool, error) {
	return s.parameterAdmin.GetParameterBool(ctx, PAYER_REGISTRY_PAUSED_KEY)
}

func (s settlementChainAdmin) SetPayerRegistryPauseStatus(ctx context.Context, paused bool) error {
	if err := s.parameterAdmin.SetBoolParameter(ctx, PAYER_REGISTRY_PAUSED_KEY, paused); err != nil {
		return err
	}

	err := ExecuteTransaction(
		ctx,
		s.signer,
		s.logger,
		s.client,
		func(opts *bind.TransactOpts) (*types.Transaction, error) {
			return s.payerRegistry.UpdatePauseStatus(opts)
		},
		func(log *types.Log) (interface{}, error) {
			return s.payerRegistry.ParsePauseStatusUpdated(*log)
		},
		func(event interface{}) {
			ev, ok := event.(*pr.PayerRegistryPauseStatusUpdated)
			if !ok {
				s.logger.Error(
					"unexpected event type, not of type PayerRegistryPauseStatusUpdated",
				)
				return
			}
			s.logger.Info(
				"payer registry pause status updated",
				zap.Bool("paused", ev.Paused),
			)
		},
	)
	if err != nil {
		if err.IsNoChange() {
			s.logger.Info("No update needed")
			return nil
		}
		return err
	}
	return nil
}
