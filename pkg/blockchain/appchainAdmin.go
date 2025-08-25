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
