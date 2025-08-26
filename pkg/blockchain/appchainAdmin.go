package blockchain

import (
	"context"
	"strings"

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
}

type appChainAdmin struct {
	client                    *ethclient.Client
	signer                    TransactionSigner
	logger                    *zap.Logger
	parameterAdmin            *ParameterAdmin
	identityUpdateBroadcaster *iu.IdentityUpdateBroadcaster
	groupMessageBroadcaster   *gm.GroupMessageBroadcaster
}

var _ IAppChainAdmin = (*appChainAdmin)(nil)

func NewAppChainAdmin(
	logger *zap.Logger,
	client *ethclient.Client,
	signer TransactionSigner,
	contractsOptions config.ContractsOptions,
	parameterAdmin *ParameterAdmin,
) (*appChainAdmin, error) {
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

	return &appChainAdmin{
		client:                    client,
		signer:                    signer,
		logger:                    logger.Named("AppChainAdmin"),
		parameterAdmin:            parameterAdmin,
		identityUpdateBroadcaster: iuBroadcaster,
		groupMessageBroadcaster:   gmBroadcaster,
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
	return a.parameterAdmin.GetParameterBool(ctx, GROUP_MESSAGE_PAUSED_KEY)
}

func (a appChainAdmin) SetGroupMessagePauseStatus(ctx context.Context, paused bool) error {
	err := a.parameterAdmin.SetBoolParameter(
		ctx,
		GROUP_MESSAGE_PAUSED_KEY,
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
		if strings.Contains(err.Error(), "0xa88ee577") {
			a.logger.Info("No update needed")
			return nil
		}
		return err
	}

	return nil
}

func (a appChainAdmin) GetIdentityUpdatePauseStatus(ctx context.Context) (bool, error) {
	return a.parameterAdmin.GetParameterBool(ctx, IDENTITY_UPDATE_PAUSED_KEY)
}

func (a appChainAdmin) SetIdentityUpdatePauseStatus(ctx context.Context, paused bool) error {
	if err := a.parameterAdmin.SetBoolParameter(ctx, IDENTITY_UPDATE_PAUSED_KEY, paused); err != nil {
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
				a.logger.Error("unexpected event type, not of type IdentityUpdateBroadcasterPauseStatusUpdated")
				return
			}
			a.logger.Info("identity update pause status updated", zap.Bool("paused", ev.Paused))
		},
	)
	if err != nil {
		if strings.Contains(err.Error(), "0xa88ee577") {
			a.logger.Info("No update needed")
			return nil
		}
		return err
	}
	return nil
}