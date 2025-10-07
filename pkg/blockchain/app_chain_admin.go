package blockchain

import (
	"context"

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
	UpdateIdentityUpdateBootstrapper(ctx context.Context) error
	GetGroupMessageBootstrapper(ctx context.Context) (common.Address, error)
	UpdateGroupMessageBootstrapper(ctx context.Context) error

	GetGroupMessagePauseStatus(ctx context.Context) (bool, error)
	UpdateGroupMessagePauseStatus(ctx context.Context) error
	GetIdentityUpdatePauseStatus(ctx context.Context) (bool, error)
	UpdateIdentityUpdatePauseStatus(ctx context.Context) error
	GetAppChainGatewayPauseStatus(ctx context.Context) (bool, error)
	UpdateAppChainGatewayPauseStatus(ctx context.Context) error

	GetGroupMessageMaxPayloadSize(ctx context.Context) (uint32, error)
	UpdateGroupMessageMaxPayloadSize(ctx context.Context) error
	GetGroupMessageMinPayloadSize(ctx context.Context) (uint32, error)
	UpdateGroupMessageMinPayloadSize(ctx context.Context) error

	GetIdentityUpdateMaxPayloadSize(ctx context.Context) (uint32, error)
	UpdateIdentityUpdateMaxPayloadSize(ctx context.Context) error
	GetIdentityUpdateMinPayloadSize(ctx context.Context) (uint32, error)
	UpdateIdentityUpdateMinPayloadSize(ctx context.Context) error

	GetRawParameter(ctx context.Context, key string) ([32]byte, error)
}

type appChainAdmin struct {
	client                    *ethclient.Client
	signer                    TransactionSigner
	logger                    *zap.Logger
	parameterAdmin            IParameterAdmin
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
	parameterAdmin IParameterAdmin,
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
	return a.identityUpdateBroadcaster.PayloadBootstrapper(&bind.CallOpts{
		Context: ctx,
	})
}

func (a appChainAdmin) UpdateIdentityUpdateBootstrapper(
	ctx context.Context,
) error {
	err := ExecuteTransaction(
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
		if err.IsNoChange() {
			a.logger.Info("No update needed")
			return nil
		}
		return err
	}

	return nil
}

func (a appChainAdmin) GetGroupMessageBootstrapper(ctx context.Context) (common.Address, error) {
	return a.groupMessageBroadcaster.PayloadBootstrapper(&bind.CallOpts{
		Context: ctx,
	})
}

func (a appChainAdmin) UpdateGroupMessageBootstrapper(
	ctx context.Context,
) error {
	err := ExecuteTransaction(
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
		if err.IsNoChange() {
			a.logger.Info("No update needed")
			return nil
		}
		return err
	}

	return nil
}

func (a appChainAdmin) GetGroupMessagePauseStatus(ctx context.Context) (bool, error) {
	return a.groupMessageBroadcaster.Paused(&bind.CallOpts{
		Context: ctx,
	})
}

func (a appChainAdmin) UpdateGroupMessagePauseStatus(ctx context.Context) error {
	err := ExecuteTransaction(
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
	return a.identityUpdateBroadcaster.Paused(&bind.CallOpts{
		Context: ctx,
	})
}

func (a appChainAdmin) UpdateIdentityUpdatePauseStatus(ctx context.Context) error {
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
	return a.appChainGateway.Paused(&bind.CallOpts{
		Context: ctx,
	})
}

func (a appChainAdmin) UpdateAppChainGatewayPauseStatus(ctx context.Context) error {
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

func (a appChainAdmin) GetGroupMessageMaxPayloadSize(ctx context.Context) (uint32, error) {
	return a.groupMessageBroadcaster.MaxPayloadSize(&bind.CallOpts{Context: ctx})
}

func (a appChainAdmin) UpdateGroupMessageMaxPayloadSize(ctx context.Context) error {
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
			a.logger.Info("No update needed (group-message max payload size)")
			return nil
		}
		return err
	}
	return nil
}

func (a appChainAdmin) GetGroupMessageMinPayloadSize(ctx context.Context) (uint32, error) {
	return a.groupMessageBroadcaster.MinPayloadSize(&bind.CallOpts{Context: ctx})
}

func (a appChainAdmin) UpdateGroupMessageMinPayloadSize(ctx context.Context) error {
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
			a.logger.Info("No update needed (group-message min payload size)")
			return nil
		}
		return err
	}
	return nil
}

func (a appChainAdmin) GetIdentityUpdateMaxPayloadSize(ctx context.Context) (uint32, error) {
	return a.identityUpdateBroadcaster.MaxPayloadSize(&bind.CallOpts{Context: ctx})
}

func (a appChainAdmin) UpdateIdentityUpdateMaxPayloadSize(ctx context.Context) error {
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
			a.logger.Info("No update needed (identity-update max payload size)")
			return nil
		}
		return err
	}
	return nil
}

func (a appChainAdmin) GetIdentityUpdateMinPayloadSize(ctx context.Context) (uint32, error) {
	return a.identityUpdateBroadcaster.MinPayloadSize(&bind.CallOpts{Context: ctx})
}

func (a appChainAdmin) UpdateIdentityUpdateMinPayloadSize(ctx context.Context) error {
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
			a.logger.Info("No update needed (identity-update min payload size)")
			return nil
		}
		return err
	}
	return nil
}

func (a appChainAdmin) GetRawParameter(ctx context.Context, key string) ([32]byte, error) {
	return a.parameterAdmin.GetRawParameter(ctx, key)
}
