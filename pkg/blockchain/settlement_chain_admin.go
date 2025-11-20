package blockchain

import (
	"context"
	"errors"
	"math/big"

	settleReg "github.com/xmtp/xmtpd/pkg/abi/settlementchainparameterregistry"

	dm "github.com/xmtp/xmtpd/pkg/abi/distributionmanager"
	nr "github.com/xmtp/xmtpd/pkg/abi/noderegistry"
	pr "github.com/xmtp/xmtpd/pkg/abi/payerregistry"
	prm "github.com/xmtp/xmtpd/pkg/abi/payerreportmanager"
	rr "github.com/xmtp/xmtpd/pkg/abi/rateregistry"
	scg "github.com/xmtp/xmtpd/pkg/abi/settlementchaingateway"

	"github.com/xmtp/xmtpd/pkg/utils"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/xmtp/xmtpd/pkg/config"
	"go.uber.org/zap"
)

type ISettlementChainAdmin interface {
	GetSettlementChainGatewayPauseStatus(ctx context.Context) (bool, error)
	UpdateSettlementChainGatewayPauseStatus(ctx context.Context) error
	GetPayerRegistryPauseStatus(ctx context.Context) (bool, error)
	UpdatePayerRegistryPauseStatus(ctx context.Context) error
	GetDistributionManagerPauseStatus(ctx context.Context) (bool, error)
	UpdateDistributionManagerPauseStatus(ctx context.Context) error

	GetNodeRegistryAdmin(ctx context.Context) (common.Address, error)
	UpdateNodeRegistryAdmin(ctx context.Context) error

	GetDistributionManagerProtocolFeesRecipient(ctx context.Context) (common.Address, error)
	UpdateDistributionManagerProtocolFeesRecipient(ctx context.Context) error

	GetPayerRegistryMinimumDeposit(ctx context.Context) (*big.Int, error)
	UpdatePayerRegistryMinimumDeposit(ctx context.Context) error
	GetPayerRegistryWithdrawLockPeriod(ctx context.Context) (uint32, error)
	UpdatePayerRegistryWithdrawLockPeriod(ctx context.Context) error

	GetPayerReportManagerProtocolFeeRate(ctx context.Context) (uint16, error)
	UpdatePayerReportManagerProtocolFeeRate(ctx context.Context) error

	GetRateRegistryMigrator(ctx context.Context) (common.Address, error)
	UpdateRateRegistryMigrator(ctx context.Context) error

	GetRawParameter(ctx context.Context, key string) ([32]byte, error)
	SetRawParameter(ctx context.Context, key string, value [32]byte) error

	BridgeParameters(ctx context.Context, keys []string) error

	GetSettlementChainGatewayVersion(ctx context.Context) (string, error)
	GetSettlementParameterRegistryVersion(ctx context.Context) (string, error)
	GetPayerReportManagerVersion(ctx context.Context) (string, error)
	GetRateRegistryVersion(ctx context.Context) (string, error)
	GetPayerRegistryVersion(ctx context.Context) (string, error)
	GetNodeRegistryVersion(ctx context.Context) (string, error)
	GetDistributionManagerVersion(ctx context.Context) (string, error)
}

const (
	noUpdateNeededMessage = "no update needed"

	pausedField               = "paused"
	protocolFeeRateField      = "protocol_fee_rate"
	minimumDepositField       = "minimum_deposit"
	withdrawLockPeriodField   = "withdraw_lock_period"
	protocolFeeRecipientField = "protocol_fee_recipient"
)

type settlementChainAdmin struct {
	client                      *ethclient.Client
	signer                      TransactionSigner
	logger                      *zap.Logger
	parameterAdmin              IParameterAdmin
	settlementChainGateway      *scg.SettlementChainGateway
	payerRegistry               *pr.PayerRegistry
	distributionManager         *dm.DistributionManager
	payerReportManager          *prm.PayerReportManager
	nodeRegistry                *nr.NodeRegistry
	rateRegistry                *rr.RateRegistry
	settlementParameterRegistry *settleReg.SettlementChainParameterRegistry
}

var _ ISettlementChainAdmin = (*settlementChainAdmin)(nil)

func NewSettlementChainAdmin(
	logger *zap.Logger,
	client *ethclient.Client,
	signer TransactionSigner,
	contractsOptions config.ContractsOptions,
	parameterAdmin IParameterAdmin,
) (ISettlementChainAdmin, error) {
	acGateway, err := scg.NewSettlementChainGateway(
		common.HexToAddress(contractsOptions.SettlementChain.GatewayAddress),
		client,
	)
	if err != nil {
		return nil, err
	}

	payerRegistry, err := pr.NewPayerRegistry(
		common.HexToAddress(contractsOptions.SettlementChain.PayerRegistryAddress),
		client,
	)
	if err != nil {
		return nil, err
	}

	distributionManager, err := dm.NewDistributionManager(
		common.HexToAddress(contractsOptions.SettlementChain.DistributionManagerAddress),
		client,
	)
	if err != nil {
		return nil, err
	}

	payerReportManager, err := prm.NewPayerReportManager(
		common.HexToAddress(contractsOptions.SettlementChain.PayerReportManagerAddress),
		client,
	)
	if err != nil {
		return nil, err
	}

	nodeContract, err := nr.NewNodeRegistry(
		common.HexToAddress(contractsOptions.SettlementChain.NodeRegistryAddress),
		client,
	)
	if err != nil {
		return nil, err
	}

	rateRegistry, err := rr.NewRateRegistry(
		common.HexToAddress(contractsOptions.SettlementChain.RateRegistryAddress),
		client,
	)
	if err != nil {
		return nil, err
	}

	paramRegistry, err := settleReg.NewSettlementChainParameterRegistry(
		common.HexToAddress(contractsOptions.SettlementChain.ParameterRegistryAddress),
		client,
	)

	settlementChainAdminLogger := logger.Named(utils.SettlementChainAdminLoggerName).With(
		utils.SettlementChainChainIDField(contractsOptions.SettlementChain.ChainID),
	)

	return &settlementChainAdmin{
		client:                      client,
		signer:                      signer,
		logger:                      settlementChainAdminLogger,
		parameterAdmin:              parameterAdmin,
		settlementChainGateway:      acGateway,
		payerRegistry:               payerRegistry,
		distributionManager:         distributionManager,
		payerReportManager:          payerReportManager,
		nodeRegistry:                nodeContract,
		rateRegistry:                rateRegistry,
		settlementParameterRegistry: paramRegistry,
	}, nil
}

func (s settlementChainAdmin) GetSettlementChainGatewayPauseStatus(
	ctx context.Context,
) (bool, error) {
	return s.settlementChainGateway.Paused(&bind.CallOpts{Context: ctx})
}

func (s settlementChainAdmin) UpdateSettlementChainGatewayPauseStatus(
	ctx context.Context,
) error {
	err := ExecuteTransaction(
		ctx,
		s.signer,
		s.logger,
		s.client,
		func(opts *bind.TransactOpts) (*types.Transaction, error) {
			return s.settlementChainGateway.UpdatePauseStatus(opts)
		},
		func(log *types.Log) (any, error) {
			return s.settlementChainGateway.ParsePauseStatusUpdated(*log)
		},
		func(event any) {
			ev, ok := event.(*scg.SettlementChainGatewayPauseStatusUpdated)
			if !ok {
				s.logger.Error(
					"unexpected event type, not of type SettlementChainGatewayPauseStatusUpdated",
				)
				return
			}
			s.logger.Info(
				"settlement-chain gateway pause status updated",
				zap.Bool(pausedField, ev.Paused),
			)
		},
	)
	if err != nil {
		if err.IsNoChange() {
			s.logger.Info(noUpdateNeededMessage)
			return nil
		}
		return err
	}
	return nil
}

func (s settlementChainAdmin) GetPayerRegistryPauseStatus(ctx context.Context) (bool, error) {
	return s.payerRegistry.Paused(&bind.CallOpts{Context: ctx})
}

func (s settlementChainAdmin) UpdatePayerRegistryPauseStatus(ctx context.Context) error {
	err := ExecuteTransaction(
		ctx,
		s.signer,
		s.logger,
		s.client,
		func(opts *bind.TransactOpts) (*types.Transaction, error) {
			return s.payerRegistry.UpdatePauseStatus(opts)
		},
		func(log *types.Log) (any, error) {
			return s.payerRegistry.ParsePauseStatusUpdated(*log)
		},
		func(event any) {
			ev, ok := event.(*pr.PayerRegistryPauseStatusUpdated)
			if !ok {
				s.logger.Error(
					"unexpected event type, not of type PayerRegistryPauseStatusUpdated",
				)
				return
			}
			s.logger.Info(
				"payer registry pause status updated",
				zap.Bool(pausedField, ev.Paused),
			)
		},
	)
	if err != nil {
		if err.IsNoChange() {
			s.logger.Info(noUpdateNeededMessage)
			return nil
		}
		return err
	}
	return nil
}

func (s settlementChainAdmin) GetDistributionManagerPauseStatus(ctx context.Context) (bool, error) {
	return s.distributionManager.Paused(&bind.CallOpts{Context: ctx})
}

func (s settlementChainAdmin) UpdateDistributionManagerPauseStatus(
	ctx context.Context,
) error {
	err := ExecuteTransaction(
		ctx,
		s.signer,
		s.logger,
		s.client,
		func(opts *bind.TransactOpts) (*types.Transaction, error) {
			return s.distributionManager.UpdatePauseStatus(opts)
		},
		func(log *types.Log) (any, error) {
			return s.distributionManager.ParsePauseStatusUpdated(*log)
		},
		func(event any) {
			ev, ok := event.(*dm.DistributionManagerPauseStatusUpdated)
			if !ok {
				s.logger.Error(
					"unexpected event type, not of type DistributionManagerPauseStatusUpdated",
				)
				return
			}
			s.logger.Info(
				"distribution manager pause status updated",
				zap.Bool(pausedField, ev.Paused),
			)
		},
	)
	if err != nil {
		if err.IsNoChange() {
			s.logger.Info(noUpdateNeededMessage)
			return nil
		}
		return err
	}
	return nil
}

func (s settlementChainAdmin) GetDistributionManagerProtocolFeesRecipient(
	ctx context.Context,
) (common.Address, error) {
	return s.distributionManager.ProtocolFeesRecipient(&bind.CallOpts{Context: ctx})
}

func (s settlementChainAdmin) UpdateDistributionManagerProtocolFeesRecipient(
	ctx context.Context,
) error {
	// Apply on-chain (adjust ABI names if they differ in your bindings)
	err := ExecuteTransaction(
		ctx,
		s.signer, s.logger, s.client,
		func(opts *bind.TransactOpts) (*types.Transaction, error) {
			return s.distributionManager.UpdateProtocolFeesRecipient(opts)
		},
		func(log *types.Log) (any, error) {
			return s.distributionManager.ParseProtocolFeesRecipientUpdated(*log)
		},
		func(event any) {
			ev, ok := event.(*dm.DistributionManagerProtocolFeesRecipientUpdated)
			if !ok {
				s.logger.Error(
					"unexpected event type, not DistributionManagerProtocolFeesRecipientUpdated",
				)
				return
			}
			s.logger.Info("distribution manager protocol fees recipient updated",
				zap.String(protocolFeeRecipientField, ev.ProtocolFeesRecipient.Hex()))
		},
	)
	if err != nil {
		if err.IsNoChange() {
			s.logger.Info("no update needed (distribution manager protocol fees recipient)")
		}
		return err
	}
	return nil
}

func (s settlementChainAdmin) GetPayerRegistryMinimumDeposit(
	ctx context.Context,
) (*big.Int, error) {
	return s.payerRegistry.MinimumDeposit(&bind.CallOpts{Context: ctx})
}

func (s settlementChainAdmin) UpdatePayerRegistryMinimumDeposit(
	ctx context.Context,
) error {
	err := ExecuteTransaction(
		ctx,
		s.signer, s.logger, s.client,
		func(opts *bind.TransactOpts) (*types.Transaction, error) {
			return s.payerRegistry.UpdateMinimumDeposit(opts)
		},
		func(log *types.Log) (any, error) {
			return s.payerRegistry.ParseMinimumDepositUpdated(*log)
		},
		func(event any) {
			ev, ok := event.(*pr.PayerRegistryMinimumDepositUpdated)
			if !ok {
				s.logger.Error("unexpected event type, not PayerRegistryMinimumDepositUpdated")
				return
			}
			s.logger.Info("payer registry minimum deposit updated",
				zap.String(minimumDepositField, ev.MinimumDeposit.String()))
		},
	)
	if err != nil {
		if err.IsNoChange() {
			s.logger.Info(noUpdateNeededMessage)
			return nil
		}
		return err
	}
	return nil
}

func (s settlementChainAdmin) GetPayerRegistryWithdrawLockPeriod(
	ctx context.Context,
) (uint32, error) {
	return s.payerRegistry.WithdrawLockPeriod(&bind.CallOpts{Context: ctx})
}

func (s settlementChainAdmin) UpdatePayerRegistryWithdrawLockPeriod(
	ctx context.Context,
) error {
	err := ExecuteTransaction(
		ctx,
		s.signer, s.logger, s.client,
		func(opts *bind.TransactOpts) (*types.Transaction, error) {
			return s.payerRegistry.UpdateWithdrawLockPeriod(opts)
		},
		func(log *types.Log) (any, error) {
			return s.payerRegistry.ParseWithdrawLockPeriodUpdated(*log)
		},
		func(event any) {
			ev, ok := event.(*pr.PayerRegistryWithdrawLockPeriodUpdated)
			if !ok {
				s.logger.Error("unexpected event type, not PayerRegistryWithdrawLockPeriodUpdated")
				return
			}
			s.logger.Info("payer registry withdraw lock period updated",
				zap.Uint32(withdrawLockPeriodField, ev.WithdrawLockPeriod))
		},
	)
	if err != nil {
		if err.IsNoChange() {
			s.logger.Info(noUpdateNeededMessage)
			return nil
		}
		return err
	}
	return nil
}

func (s settlementChainAdmin) GetPayerReportManagerProtocolFeeRate(
	ctx context.Context,
) (uint16, error) {
	return s.payerReportManager.ProtocolFeeRate(&bind.CallOpts{Context: ctx})
}

func (s settlementChainAdmin) UpdatePayerReportManagerProtocolFeeRate(
	ctx context.Context,
) error {
	err := ExecuteTransaction(
		ctx,
		s.signer, s.logger, s.client,
		func(opts *bind.TransactOpts) (*types.Transaction, error) {
			return s.payerReportManager.UpdateProtocolFeeRate(opts)
		},
		func(log *types.Log) (any, error) {
			return s.payerReportManager.ParseProtocolFeeRateUpdated(*log)
		},
		func(event any) {
			ev, ok := event.(*prm.PayerReportManagerProtocolFeeRateUpdated)
			if !ok {
				s.logger.Error(
					"unexpected event type, not PayerReportManagerProtocolFeeRateUpdated",
				)
				return
			}
			s.logger.Info("payer report manager protocol fee updated",
				zap.Uint16(protocolFeeRateField, ev.ProtocolFeeRate))
		},
	)
	if err != nil {
		if err.IsNoChange() {
			s.logger.Info(noUpdateNeededMessage)
			return nil
		}
		return err
	}
	return nil
}

func (s settlementChainAdmin) GetNodeRegistryAdmin(ctx context.Context) (common.Address, error) {
	return s.nodeRegistry.Admin(&bind.CallOpts{Context: ctx})
}

func (s settlementChainAdmin) UpdateNodeRegistryAdmin(ctx context.Context) error {
	err := ExecuteTransaction(
		ctx,
		s.signer, s.logger, s.client,
		func(opts *bind.TransactOpts) (*types.Transaction, error) {
			return s.nodeRegistry.UpdateAdmin(opts)
		},
		func(log *types.Log) (any, error) {
			return s.nodeRegistry.ParseAdminUpdated(*log)
		},
		func(event any) {
			ev, ok := event.(*nr.NodeRegistryAdminUpdated)
			if !ok {
				s.logger.Error(
					"unexpected event type, not NodeRegistryAdminUpdated",
				)
				return
			}
			s.logger.Info("node registry admin updated",
				zap.String("admin", ev.Admin.Hex()))
		},
	)
	if err != nil {
		if err.IsNoChange() {
			s.logger.Info(noUpdateNeededMessage)
			return nil
		}
		return err
	}
	return nil
}

func (s settlementChainAdmin) GetRateRegistryMigrator(ctx context.Context) (common.Address, error) {
	return common.Address{}, errors.New("not implemented")
}

func (s settlementChainAdmin) UpdateRateRegistryMigrator(ctx context.Context) error {
	return errors.New("not implemented")
}

func (s settlementChainAdmin) GetRawParameter(ctx context.Context, key string) ([32]byte, error) {
	return s.parameterAdmin.GetRawParameter(ctx, key)
}

func (s settlementChainAdmin) SetRawParameter(
	ctx context.Context,
	key string,
	value [32]byte,
) error {
	return s.parameterAdmin.SetRawParameter(ctx, key, value)
}

func (s settlementChainAdmin) BridgeParameters(ctx context.Context, keys []string) error {
	chainIds := make([]*big.Int, 1)
	chainIds[0] = big.NewInt(351243127)

	gasLimit := big.NewInt(3000000)
	maxFeePerGas := big.NewInt(2000000000)
	amountToSend := big.NewInt(100000)

	err := ExecuteTransaction(
		ctx,
		s.signer,
		s.logger,
		s.client,
		func(opts *bind.TransactOpts) (*types.Transaction, error) {
			return s.settlementChainGateway.SendParameters(
				opts,
				chainIds,
				keys,
				gasLimit,
				maxFeePerGas,
				amountToSend,
			)
		},
		func(l *types.Log) (any, error) {
			return s.settlementChainGateway.ParseParametersSent(*l)
		},
		func(ev any) {
			e, ok := ev.(*scg.SettlementChainGatewayParametersSent)
			if ok {
				s.logger.Info("parameters sent",
					zap.Uint64("nonce", e.Nonce.Uint64()),
					zap.Int("keys", len(e.Keys)),
				)
			}
		},
	)

	return err
}

func (s settlementChainAdmin) GetSettlementChainGatewayVersion(
	ctx context.Context,
) (string, error) {
	return s.settlementChainGateway.Version(&bind.CallOpts{Context: ctx})
}

func (s settlementChainAdmin) GetSettlementParameterRegistryVersion(
	ctx context.Context,
) (string, error) {
	return s.settlementParameterRegistry.Version(&bind.CallOpts{Context: ctx})
}

func (s settlementChainAdmin) GetPayerReportManagerVersion(ctx context.Context) (string, error) {
	return s.payerReportManager.Version(&bind.CallOpts{Context: ctx})
}

func (s settlementChainAdmin) GetRateRegistryVersion(ctx context.Context) (string, error) {
	return s.rateRegistry.Version(&bind.CallOpts{Context: ctx})
}

func (s settlementChainAdmin) GetPayerRegistryVersion(ctx context.Context) (string, error) {
	return s.payerRegistry.Version(&bind.CallOpts{Context: ctx})
}

func (s settlementChainAdmin) GetNodeRegistryVersion(ctx context.Context) (string, error) {
	return s.nodeRegistry.Version(&bind.CallOpts{Context: ctx})
}

func (s settlementChainAdmin) GetDistributionManagerVersion(ctx context.Context) (string, error) {
	return s.distributionManager.Version(&bind.CallOpts{Context: ctx})
}
