package blockchain

import (
	"context"
	"math/big"

	dm "github.com/xmtp/xmtpd/pkg/abi/distributionmanager"
	pr "github.com/xmtp/xmtpd/pkg/abi/payerregistry"
	prm "github.com/xmtp/xmtpd/pkg/abi/payerreportmanager"
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
	GetDistributionManagerPauseStatus(ctx context.Context) (bool, error)
	SetDistributionManagerPauseStatus(ctx context.Context, paused bool) error

	GetDistributionManagerProtocolFeesRecipient(ctx context.Context) (common.Address, error)
	SetDistributionManagerProtocolFeesRecipient(ctx context.Context, addr common.Address) error
	GetNodeRegistryAdmin(ctx context.Context) (common.Address, error)
	SetNodeRegistryAdmin(ctx context.Context, addr common.Address) error

	GetPayerRegistryMinimumDeposit(ctx context.Context) (*big.Int, error)
	SetPayerRegistryMinimumDeposit(ctx context.Context, v *big.Int) error
	GetPayerRegistryWithdrawLockPeriod(ctx context.Context) (uint32, error)
	SetPayerRegistryWithdrawLockPeriod(ctx context.Context, v uint32) error

	GetPayerReportManagerProtocolFeeRate(ctx context.Context) (uint16, error)
	SetPayerReportManagerProtocolFeeRate(ctx context.Context, v uint16) error

	GetRateRegistryMigrator(ctx context.Context) (common.Address, error)
	SetRateRegistryMigrator(ctx context.Context, addr common.Address) error

	GetRawParameter(ctx context.Context, key string) ([32]byte, error)
	SetRawParameter(ctx context.Context, key string, value [32]byte) error

	BridgeParameters(ctx context.Context, keys []string) error
}

type settlementChainAdmin struct {
	client                 *ethclient.Client
	signer                 TransactionSigner
	logger                 *zap.Logger
	parameterAdmin         IParameterAdmin
	settlementChainGateway *scg.SettlementChainGateway
	payerRegistry          *pr.PayerRegistry
	distributionManager    *dm.DistributionManager
	payerReportManager     *prm.PayerReportManager
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

	return &settlementChainAdmin{
		client:                 client,
		signer:                 signer,
		logger:                 logger.Named("SettlementChainAdmin"),
		parameterAdmin:         parameterAdmin,
		settlementChainGateway: acGateway,
		payerRegistry:          payerRegistry,
		distributionManager:    distributionManager,
		payerReportManager:     payerReportManager,
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

func (s settlementChainAdmin) GetDistributionManagerPauseStatus(ctx context.Context) (bool, error) {
	return s.parameterAdmin.GetParameterBool(ctx, DISTRIBUTION_MANAGER_PAUSED_KEY)
}

func (s settlementChainAdmin) SetDistributionManagerPauseStatus(
	ctx context.Context,
	paused bool,
) error {
	if err := s.parameterAdmin.SetBoolParameter(ctx, DISTRIBUTION_MANAGER_PAUSED_KEY, paused); err != nil {
		return err
	}

	err := ExecuteTransaction(
		ctx,
		s.signer,
		s.logger,
		s.client,
		func(opts *bind.TransactOpts) (*types.Transaction, error) {
			return s.distributionManager.UpdatePauseStatus(opts)
		},
		func(log *types.Log) (interface{}, error) {
			return s.distributionManager.ParsePauseStatusUpdated(*log)
		},
		func(event interface{}) {
			ev, ok := event.(*dm.DistributionManagerPauseStatusUpdated)
			if !ok {
				s.logger.Error(
					"unexpected event type, not of type DistributionManagerPauseStatusUpdated",
				)
				return
			}
			s.logger.Info(
				"distribution manager pause status updated",
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

func (s settlementChainAdmin) GetDistributionManagerProtocolFeesRecipient(
	ctx context.Context,
) (common.Address, error) {
	return s.parameterAdmin.GetParameterAddress(
		ctx,
		DISTRIBUTION_MANAGER_PROTOCOL_FEES_RECIPIENT_KEY,
	)
}

func (s settlementChainAdmin) SetDistributionManagerProtocolFeesRecipient(
	ctx context.Context,
	addr common.Address,
) error {
	if err := s.parameterAdmin.SetAddressParameter(ctx, DISTRIBUTION_MANAGER_PROTOCOL_FEES_RECIPIENT_KEY, addr); err != nil {
		return err
	}

	// Apply on-chain (adjust ABI names if they differ in your bindings)
	err := ExecuteTransaction(
		ctx,
		s.signer, s.logger, s.client,
		func(opts *bind.TransactOpts) (*types.Transaction, error) {
			return s.distributionManager.UpdateProtocolFeesRecipient(opts)
		},
		func(log *types.Log) (interface{}, error) {
			return s.distributionManager.ParseProtocolFeesRecipientUpdated(*log)
		},
		func(event interface{}) {
			ev, ok := event.(*dm.DistributionManagerProtocolFeesRecipientUpdated)
			if !ok {
				s.logger.Error(
					"unexpected event type, not DistributionManagerProtocolFeesRecipientUpdated",
				)
				return
			}
			s.logger.Info("distribution manager protocol fees recipient updated",
				zap.String("protocolFeesRecipient", ev.ProtocolFeesRecipient.Hex()))
		},
	)
	if err != nil {
		if err.IsNoChange() {
			s.logger.Info("No update needed (distribution manager protocol fees recipient)",
				zap.String("protocolFeesRecipient", addr.Hex()))
			return nil
		}
		return err
	}
	return nil
}

func (s settlementChainAdmin) GetNodeRegistryAdmin(ctx context.Context) (common.Address, error) {
	return s.parameterAdmin.GetParameterAddress(ctx, NODE_REGISTRY_ADMIN_KEY)
}

func (s settlementChainAdmin) SetNodeRegistryAdmin(ctx context.Context, addr common.Address) error {
	return s.parameterAdmin.SetAddressParameter(ctx, NODE_REGISTRY_ADMIN_KEY, addr)
}

func (s settlementChainAdmin) GetPayerRegistryMinimumDeposit(
	ctx context.Context,
) (*big.Int, error) {
	return s.parameterAdmin.GetParameterUint96(ctx, PAYER_REGISTRY_MINIMUM_DEPOSIT_KEY)
}

func (s settlementChainAdmin) SetPayerRegistryMinimumDeposit(
	ctx context.Context,
	v *big.Int,
) error {
	if err := s.parameterAdmin.SetUint96Parameter(ctx, PAYER_REGISTRY_MINIMUM_DEPOSIT_KEY, v); err != nil {
		return err
	}

	err := ExecuteTransaction(
		ctx,
		s.signer, s.logger, s.client,
		func(opts *bind.TransactOpts) (*types.Transaction, error) {
			return s.payerRegistry.UpdateMinimumDeposit(opts)
		},
		func(log *types.Log) (interface{}, error) {
			return s.payerRegistry.ParseMinimumDepositUpdated(*log)
		},
		func(event interface{}) {
			ev, ok := event.(*pr.PayerRegistryMinimumDepositUpdated)
			if !ok {
				s.logger.Error("unexpected event type, not PayerRegistryMinimumDepositUpdated")
				return
			}
			s.logger.Info("payer registry minimum deposit updated",
				zap.String("minimumDeposit", ev.MinimumDeposit.String()))
		},
	)
	if err != nil {
		if err.IsNoChange() {
			s.logger.Info("No update needed (payer registry minimum deposit)",
				zap.String("minimumDeposit", v.String()))
			return nil
		}
		return err
	}
	return nil
}

func (s settlementChainAdmin) GetPayerRegistryWithdrawLockPeriod(
	ctx context.Context,
) (uint32, error) {
	return s.parameterAdmin.GetParameterUint32(ctx, PAYER_REGISTRY_WITHDRAW_LOCK_PERIOD_KEY)
}

func (s settlementChainAdmin) SetPayerRegistryWithdrawLockPeriod(
	ctx context.Context,
	v uint32,
) error {
	if err := s.parameterAdmin.SetUint32Parameter(ctx, PAYER_REGISTRY_WITHDRAW_LOCK_PERIOD_KEY, v); err != nil {
		return err
	}

	err := ExecuteTransaction(
		ctx,
		s.signer, s.logger, s.client,
		func(opts *bind.TransactOpts) (*types.Transaction, error) {
			return s.payerRegistry.UpdateWithdrawLockPeriod(opts)
		},
		func(log *types.Log) (interface{}, error) {
			return s.payerRegistry.ParseWithdrawLockPeriodUpdated(*log)
		},
		func(event interface{}) {
			ev, ok := event.(*pr.PayerRegistryWithdrawLockPeriodUpdated)
			if !ok {
				s.logger.Error("unexpected event type, not PayerRegistryWithdrawLockPeriodUpdated")
				return
			}
			s.logger.Info("payer registry withdraw lock period updated",
				zap.Uint32("withdrawLockPeriod", ev.WithdrawLockPeriod))
		},
	)
	if err != nil {
		if err.IsNoChange() {
			s.logger.Info("No update needed (payer registry withdraw lock period)",
				zap.Uint32("withdrawLockPeriod", v))
			return nil
		}
		return err
	}
	return nil
}

func (s settlementChainAdmin) GetPayerReportManagerProtocolFeeRate(
	ctx context.Context,
) (uint16, error) {
	return s.parameterAdmin.GetParameterUint16(ctx, PAYER_REPORT_MANAGER_PROTOCOL_FEE_RATE_KEY)
}

func (s settlementChainAdmin) SetPayerReportManagerProtocolFeeRate(
	ctx context.Context,
	v uint16,
) error {
	if err := s.parameterAdmin.SetUint16Parameter(ctx, PAYER_REPORT_MANAGER_PROTOCOL_FEE_RATE_KEY, v); err != nil {
		return err
	}

	err := ExecuteTransaction(
		ctx,
		s.signer, s.logger, s.client,
		func(opts *bind.TransactOpts) (*types.Transaction, error) {
			return s.payerReportManager.UpdateProtocolFeeRate(opts)
		},
		func(log *types.Log) (interface{}, error) {
			return s.payerReportManager.ParseProtocolFeeRateUpdated(*log)
		},
		func(event interface{}) {
			ev, ok := event.(*prm.PayerReportManagerProtocolFeeRateUpdated)
			if !ok {
				s.logger.Error(
					"unexpected event type, not PayerReportManagerProtocolFeeRateUpdated",
				)
				return
			}
			s.logger.Info("payer report manager protocol fee updated",
				zap.Uint16("protocolFeeRate", ev.ProtocolFeeRate))
		},
	)
	if err != nil {
		if err.IsNoChange() {
			s.logger.Info("No update needed (payer report manager protocol fee)",
				zap.Uint16("protocolFeeRate", v))
			return nil
		}
		return err
	}
	return nil
}

func (s settlementChainAdmin) GetRateRegistryMigrator(ctx context.Context) (common.Address, error) {
	return s.parameterAdmin.GetParameterAddress(ctx, RATE_REGISTRY_MIGRATOR_KEY)
}

func (s settlementChainAdmin) SetRateRegistryMigrator(
	ctx context.Context,
	addr common.Address,
) error {
	return s.parameterAdmin.SetAddressParameter(ctx, RATE_REGISTRY_MIGRATOR_KEY, addr)
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
		func(l *types.Log) (interface{}, error) {
			return s.settlementChainGateway.ParseParametersSent(*l)
		},
		func(ev interface{}) {
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
