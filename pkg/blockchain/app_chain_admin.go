package blockchain

import (
	"context"
	"strings"

	acg "github.com/xmtp/xmtpd/pkg/abi/appchaingateway"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	gm "github.com/xmtp/xmtpd/pkg/abi/groupmessagebroadcaster"
	iu "github.com/xmtp/xmtpd/pkg/abi/identityupdatebroadcaster"
	"github.com/xmtp/xmtpd/pkg/config"
	"go.uber.org/zap"
)

type IAppChainAdmin interface {
	GetIdentityUpdateBootstrapper(ctx context.Context) (common.Address, error)
	SetIdentityUpdateBootstrapper(ctx context.Context, address common.Address) error
	GetGroupMessageBootstrapper(ctx context.Context) (common.Address, error)
	SetGroupMessageBootstrapper(ctx context.Context, address common.Address) error

	GetGroupMessagePauseStatus(ctx context.Context) (bool, error)
	SetGroupMessagePauseStatus(ctx context.Context, paused bool) error
	GetIdentityUpdatePauseStatus(ctx context.Context) (bool, error)
	SetIdentityUpdatePauseStatus(ctx context.Context, paused bool) error
	GetAppChainGatewayPauseStatus(ctx context.Context) (bool, error)
	SetAppChainGatewayPauseStatus(ctx context.Context, paused bool) error

	GetGroupMessageMaxPayloadSize(ctx context.Context) (uint64, error)
	SetGroupMessageMaxPayloadSize(ctx context.Context, size uint64) error
	GetGroupMessageMinPayloadSize(ctx context.Context) (uint64, error)
	SetGroupMessageMinPayloadSize(ctx context.Context, size uint64) error

	GetIdentityUpdateMaxPayloadSize(ctx context.Context) (uint64, error)
	SetIdentityUpdateMaxPayloadSize(ctx context.Context, size uint64) error
	GetIdentityUpdateMinPayloadSize(ctx context.Context) (uint64, error)
	SetIdentityUpdateMinPayloadSize(ctx context.Context, size uint64) error
}

type appChainAdmin struct {
	client                    *ethclient.Client
	signer                    TransactionSigner
	logger                    *zap.Logger
	parameterAdmin            *ParameterAdmin
	identityUpdateBroadcaster *iu.IdentityUpdateBroadcaster
	groupMessageBroadcaster   *gm.GroupMessageBroadcaster
	appChainGateway           *acg.AppChainGateway
}

var _ IAppChainAdmin = (*appChainAdmin)(nil)

func NewAppChainAdmin(
	logger *zap.Logger,
	client *ethclient.Client,
	signer TransactionSigner,
	contractsOptions config.ContractsOptions,
	parameterAdmin *ParameterAdmin,
) (IAppChainAdmin, error) {
	iuBroadcaster, err := iu.NewIdentityUpdateBroadcaster(
		common.HexToAddress(contractsOptions.AppChain.IdentityUpdateBroadcasterAddress),
		client,
	)
	if err != nil {
		return nil, err
	}

	gmBroadcaster, err := gm.NewGroupMessageBroadcaster(
		common.HexToAddress(contractsOptions.AppChain.GroupMessageBroadcasterAddress),
		client,
	)
	if err != nil {
		return nil, err
	}

	acGateway, err := acg.NewAppChainGateway(
		common.HexToAddress(contractsOptions.AppChain.GatewayAddress),
		client,
	)
	if err != nil {
		return nil, err
	}

	return &appChainAdmin{
		client:                    client,
		signer:                    signer,
		logger:                    logger.Named("AppChainAdmin"),
		parameterAdmin:            parameterAdmin,
		identityUpdateBroadcaster: iuBroadcaster,
		groupMessageBroadcaster:   gmBroadcaster,
		appChainGateway:           acGateway,
	}, nil
}

func (a appChainAdmin) GetIdentityUpdateBootstrapper(ctx context.Context) (common.Address, error) {
	return a.parameterAdmin.GetParameterAddress(ctx, IDENTITY_UPDATE_PAYLOAD_BOOTSTRAPPER_KEY)
}

func (a appChainAdmin) SetIdentityUpdateBootstrapper(
	ctx context.Context,
	address common.Address,
) error {
	err := a.parameterAdmin.SetAddressParameter(
		ctx,
		IDENTITY_UPDATE_PAYLOAD_BOOTSTRAPPER_KEY,
		address,
	)
	if err != nil {
		return err
	}

	err = ExecuteTransaction(
		ctx,
		a.signer,
		a.logger,
		a.client,
		func(opts *bind.TransactOpts) (*types.Transaction, error) {
			return a.identityUpdateBroadcaster.UpdatePayloadBootstrapper(opts)
		},
		func(log *types.Log) (interface{}, error) {
			return a.identityUpdateBroadcaster.ParsePayloadBootstrapperUpdated(*log)
		},
		func(event interface{}) {
			parameterSet, ok := event.(*iu.IdentityUpdateBroadcasterPayloadBootstrapperUpdated)
			if !ok {
				a.logger.Error(
					"unexpected event type, not of type IdentityUpdateBroadcasterPayloadBootstrapperUpdated",
				)
				return
			}
			a.logger.Info("payload bootstrapper updated",
				zap.String("payload_bootstrapper", parameterSet.PayloadBootstrapper.Hex()),
			)
		},
	)
	if err != nil {
		if strings.Contains(err.Error(), "0xa88ee577") {
			a.logger.Info("No update needed",
				zap.String("payload_bootstrapper", address.Hex()),
			)
			return nil
		}
		return err
	}

	return nil
}

func (a appChainAdmin) GetGroupMessageBootstrapper(ctx context.Context) (common.Address, error) {
	return a.parameterAdmin.GetParameterAddress(ctx, GROUP_MESSAGE_PAYLOAD_BOOTSTRAPPER_KEY)
}

func (a appChainAdmin) SetGroupMessageBootstrapper(
	ctx context.Context,
	address common.Address,
) error {
	err := a.parameterAdmin.SetAddressParameter(
		ctx,
		GROUP_MESSAGE_PAYLOAD_BOOTSTRAPPER_KEY,
		address,
	)
	if err != nil {
		return err
	}

	err = ExecuteTransaction(
		ctx,
		a.signer,
		a.logger,
		a.client,
		func(opts *bind.TransactOpts) (*types.Transaction, error) {
			return a.groupMessageBroadcaster.UpdatePayloadBootstrapper(opts)
		},
		func(log *types.Log) (interface{}, error) {
			return a.groupMessageBroadcaster.ParsePayloadBootstrapperUpdated(*log)
		},
		func(event interface{}) {
			parameterSet, ok := event.(*gm.GroupMessageBroadcasterPayloadBootstrapperUpdated)
			if !ok {
				a.logger.Error(
					"unexpected event type, not of type GroupMessageBroadcasterPayloadBootstrapperUpdated",
				)
				return
			}
			a.logger.Info("payload bootstrapper updated",
				zap.String("payload_bootstrapper", parameterSet.PayloadBootstrapper.Hex()),
			)
		},
	)
	if err != nil {
		if strings.Contains(err.Error(), "0xa88ee577") {
			a.logger.Info("No update needed",
				zap.String("payload_bootstrapper", address.Hex()),
			)
			return nil
		}
		return err
	}

	return nil
}

func (a appChainAdmin) GetGroupMessagePauseStatus(ctx context.Context) (bool, error) {
	return a.parameterAdmin.GetParameterBool(ctx, GROUP_MESSAGE_BROADCASTER_PAUSED_KEY)
}

func (a appChainAdmin) SetGroupMessagePauseStatus(ctx context.Context, paused bool) error {
	err := a.parameterAdmin.SetBoolParameter(
		ctx,
		GROUP_MESSAGE_BROADCASTER_PAUSED_KEY,
		paused,
	)
	if err != nil {
		return err
	}

	err = ExecuteTransaction(
		ctx,
		a.signer,
		a.logger,
		a.client,
		func(opts *bind.TransactOpts) (*types.Transaction, error) {
			return a.groupMessageBroadcaster.UpdatePauseStatus(opts)
		},
		func(log *types.Log) (interface{}, error) {
			return a.groupMessageBroadcaster.ParsePauseStatusUpdated(*log)
		},
		func(event interface{}) {
			parameterSet, ok := event.(*gm.GroupMessageBroadcasterPauseStatusUpdated)
			if !ok {
				a.logger.Error(
					"unexpected event type, not of type GroupMessageBroadcasterPauseStatusUpdated",
				)
				return
			}

			a.logger.Info("pause status updated", zap.Bool("paused", parameterSet.Paused))
		},
	)
	if err != nil {
		if err.IsNoChange() {
			a.logger.Info("No update needed")
			return nil
		}
		return err
	}

	return nil
}

func (a appChainAdmin) GetIdentityUpdatePauseStatus(ctx context.Context) (bool, error) {
	return a.parameterAdmin.GetParameterBool(ctx, IDENTITY_UPDATE_BROADCASTER_PAUSED_KEY)
}

func (a appChainAdmin) SetIdentityUpdatePauseStatus(ctx context.Context, paused bool) error {
	if err := a.parameterAdmin.SetBoolParameter(ctx, IDENTITY_UPDATE_BROADCASTER_PAUSED_KEY, paused); err != nil {
		return err
	}

	err := ExecuteTransaction(
		ctx,
		a.signer,
		a.logger,
		a.client,
		func(opts *bind.TransactOpts) (*types.Transaction, error) {
			return a.identityUpdateBroadcaster.UpdatePauseStatus(opts)
		},
		func(log *types.Log) (interface{}, error) {
			return a.identityUpdateBroadcaster.ParsePauseStatusUpdated(*log)
		},
		func(event interface{}) {
			ev, ok := event.(*iu.IdentityUpdateBroadcasterPauseStatusUpdated)
			if !ok {
				a.logger.Error(
					"unexpected event type, not of type IdentityUpdateBroadcasterPauseStatusUpdated",
				)
				return
			}
			a.logger.Info("identity update pause status updated", zap.Bool("paused", ev.Paused))
		},
	)
	if err != nil {
		if err.IsNoChange() {
			a.logger.Info("No update needed")
			return nil
		}
		return err
	}
	return nil
}

func (a appChainAdmin) GetAppChainGatewayPauseStatus(ctx context.Context) (bool, error) {
	return a.parameterAdmin.GetParameterBool(ctx, APP_CHAIN_GATEWAY_PAUSED_KEY)
}

func (a appChainAdmin) SetAppChainGatewayPauseStatus(ctx context.Context, paused bool) error {
	if err := a.parameterAdmin.SetBoolParameter(ctx, APP_CHAIN_GATEWAY_PAUSED_KEY, paused); err != nil {
		return err
	}

	err := ExecuteTransaction(
		ctx,
		a.signer,
		a.logger,
		a.client,
		func(opts *bind.TransactOpts) (*types.Transaction, error) {
			return a.appChainGateway.UpdatePauseStatus(opts)
		},
		func(log *types.Log) (interface{}, error) {
			return a.appChainGateway.ParsePauseStatusUpdated(*log)
		},
		func(event interface{}) {
			ev, ok := event.(*acg.AppChainGatewayPauseStatusUpdated)
			if !ok {
				a.logger.Error(
					"unexpected event type, not of type AppChainGatewayPauseStatusUpdated",
				)
				return
			}
			a.logger.Info("app-chain gateway pause status updated", zap.Bool("paused", ev.Paused))
		},
	)
	if err != nil {
		if err.IsNoChange() {
			a.logger.Info("No update needed")
			return nil
		}
		return err
	}
	return nil
}

func (a appChainAdmin) GetGroupMessageMaxPayloadSize(ctx context.Context) (uint64, error) {
	val, perr := a.parameterAdmin.GetParameterUint64(
		ctx,
		GROUP_MESSAGE_BROADCASTER_MAX_PAYLOAD_SIZE_KEY,
	)
	if perr != nil {
		return 0, perr
	}
	return val, nil
}

func (a appChainAdmin) SetGroupMessageMaxPayloadSize(ctx context.Context, size uint64) error {
	if err := a.parameterAdmin.SetUint64Parameter(ctx, GROUP_MESSAGE_BROADCASTER_MAX_PAYLOAD_SIZE_KEY, size); err != nil {
		return err
	}

	err := ExecuteTransaction(
		ctx,
		a.signer,
		a.logger,
		a.client,
		func(opts *bind.TransactOpts) (*types.Transaction, error) {
			return a.groupMessageBroadcaster.UpdateMaxPayloadSize(opts)
		},
		func(log *types.Log) (interface{}, error) {
			return a.groupMessageBroadcaster.ParseMaxPayloadSizeUpdated(*log)
		},
		func(event interface{}) {
			ev, ok := event.(*gm.GroupMessageBroadcasterMaxPayloadSizeUpdated)
			if !ok {
				a.logger.Error(
					"unexpected event type, not GroupMessageBroadcasterMaxPayloadSizeUpdated",
				)
				return
			}
			a.logger.Info("group-message max payload size updated",
				zap.Uint64("maxPayloadSize", ev.Size.Uint64()))
		},
	)
	if err != nil {
		if err.IsNoChange() {
			a.logger.Info("No update needed (group-message max payload size)",
				zap.Uint64("maxPayloadSize", size))
			return nil
		}
		return err
	}
	return nil
}

func (a appChainAdmin) GetGroupMessageMinPayloadSize(ctx context.Context) (uint64, error) {
	val, perr := a.parameterAdmin.GetParameterUint64(
		ctx,
		GROUP_MESSAGE_BROADCASTER_MIN_PAYLOAD_SIZE_KEY,
	)
	if perr != nil {
		return 0, perr
	}
	return val, nil
}

func (a appChainAdmin) SetGroupMessageMinPayloadSize(ctx context.Context, size uint64) error {
	if err := a.parameterAdmin.SetUint64Parameter(ctx, GROUP_MESSAGE_BROADCASTER_MIN_PAYLOAD_SIZE_KEY, size); err != nil {
		return err
	}

	err := ExecuteTransaction(
		ctx,
		a.signer,
		a.logger,
		a.client,
		func(opts *bind.TransactOpts) (*types.Transaction, error) {
			return a.groupMessageBroadcaster.UpdateMinPayloadSize(opts)
		},
		func(log *types.Log) (interface{}, error) {
			return a.groupMessageBroadcaster.ParseMinPayloadSizeUpdated(*log)
		},
		func(event interface{}) {
			ev, ok := event.(*gm.GroupMessageBroadcasterMinPayloadSizeUpdated)
			if !ok {
				a.logger.Error(
					"unexpected event type, not GroupMessageBroadcasterMinPayloadSizeUpdated",
				)
				return
			}
			a.logger.Info("group-message min payload size updated",
				zap.Uint64("minPayloadSize", ev.Size.Uint64()))
		},
	)
	if err != nil {
		if err.IsNoChange() {
			a.logger.Info("No update needed (group-message min payload size)",
				zap.Uint64("minPayloadSize", size))
			return nil
		}
		return err
	}
	return nil
}

func (a appChainAdmin) GetIdentityUpdateMaxPayloadSize(ctx context.Context) (uint64, error) {
	val, perr := a.parameterAdmin.GetParameterUint64(
		ctx,
		IDENTITY_UPDATE_BROADCASTER_MAX_PAYLOAD_SIZE_KEY,
	)
	if perr != nil {
		return 0, perr
	}
	return val, nil
}

func (a appChainAdmin) SetIdentityUpdateMaxPayloadSize(ctx context.Context, size uint64) error {
	if err := a.parameterAdmin.SetUint64Parameter(ctx, IDENTITY_UPDATE_BROADCASTER_MAX_PAYLOAD_SIZE_KEY, size); err != nil {
		return err
	}

	err := ExecuteTransaction(
		ctx,
		a.signer,
		a.logger,
		a.client,
		func(opts *bind.TransactOpts) (*types.Transaction, error) {
			return a.identityUpdateBroadcaster.UpdateMaxPayloadSize(opts)
		},
		func(log *types.Log) (interface{}, error) {
			return a.identityUpdateBroadcaster.ParseMaxPayloadSizeUpdated(*log)
		},
		func(event interface{}) {
			ev, ok := event.(*iu.IdentityUpdateBroadcasterMaxPayloadSizeUpdated)
			if !ok {
				a.logger.Error(
					"unexpected event type, not IdentityUpdateBroadcasterMaxPayloadSizeUpdated",
				)
				return
			}
			a.logger.Info("identity-update max payload size updated",
				zap.Uint64("maxPayloadSize", ev.Size.Uint64()))
		},
	)
	if err != nil {
		if err.IsNoChange() {
			a.logger.Info("No update needed (identity-update max payload size)",
				zap.Uint64("maxPayloadSize", size))
			return nil
		}
		return err
	}
	return nil
}

func (a appChainAdmin) GetIdentityUpdateMinPayloadSize(ctx context.Context) (uint64, error) {
	val, perr := a.parameterAdmin.GetParameterUint64(
		ctx,
		IDENTITY_UPDATE_BROADCASTER_MIN_PAYLOAD_SIZE_KEY,
	)
	if perr != nil {
		return 0, perr
	}
	return val, nil
}

func (a appChainAdmin) SetIdentityUpdateMinPayloadSize(ctx context.Context, size uint64) error {
	if err := a.parameterAdmin.SetUint64Parameter(ctx, IDENTITY_UPDATE_BROADCASTER_MIN_PAYLOAD_SIZE_KEY, size); err != nil {
		return err
	}

	err := ExecuteTransaction(
		ctx,
		a.signer,
		a.logger,
		a.client,
		func(opts *bind.TransactOpts) (*types.Transaction, error) {
			return a.identityUpdateBroadcaster.UpdateMinPayloadSize(opts)
		},
		func(log *types.Log) (interface{}, error) {
			return a.identityUpdateBroadcaster.ParseMinPayloadSizeUpdated(*log)
		},
		func(event interface{}) {
			ev, ok := event.(*iu.IdentityUpdateBroadcasterMinPayloadSizeUpdated)
			if !ok {
				a.logger.Error(
					"unexpected event type, not IdentityUpdateBroadcasterMinPayloadSizeUpdated",
				)
				return
			}
			a.logger.Info("identity-update min payload size updated",
				zap.Uint64("minPayloadSize", ev.Size.Uint64()))
		},
	)
	if err != nil {
		if err.IsNoChange() {
			a.logger.Info("No update needed (identity-update min payload size)",
				zap.Uint64("minPayloadSize", size))
			return nil
		}
		return err
	}
	return nil
}
