package testdata

import (
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/require"
	gm "github.com/xmtp/xmtpd/pkg/abi/groupmessagebroadcaster"
	iu "github.com/xmtp/xmtpd/pkg/abi/identityupdatebroadcaster"
	paramReg "github.com/xmtp/xmtpd/pkg/abi/settlementchainparameterregistry"
	"github.com/xmtp/xmtpd/pkg/blockchain"
	"github.com/xmtp/xmtpd/pkg/config"
	"github.com/xmtp/xmtpd/pkg/testutils"
	"github.com/xmtp/xmtpd/pkg/testutils/anvil"
	"github.com/xmtp/xmtpd/pkg/utils"
	"go.uber.org/zap"
)

const (
	PayerAddress          = "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"
	PayerPrivateKeyString = "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
)

func NewMigratorBlockchain(t *testing.T) *config.ContractsOptions {
	wsURL, rpcURL := anvil.StartAnvil(t, false)
	cfg := testutils.NewContractsOptions(t, wsURL, rpcURL)

	signer, err := blockchain.NewPrivateKeySigner(
		PayerPrivateKeyString,
		cfg.AppChain.ChainID,
	)
	require.NoError(t, err)

	client, err := ethclient.Dial(wsURL)
	require.NoError(t, err)

	setBroadcasterKeys(t, client, signer, cfg.SettlementChain.ParameterRegistryAddress)
	triggerGroupMessageBroadcasterUpdates(
		t,
		client,
		signer,
		cfg.AppChain.GroupMessageBroadcasterAddress,
	)
	triggerIdentityUpdateBroadcasterUpdates(
		t,
		client,
		signer,
		cfg.AppChain.IdentityUpdateBroadcasterAddress,
	)

	return cfg
}

// TODO: All the following logic should be abstracted to specific ParameterRegistry methods.
// Those methods should be called from the CLI for an operator to be able to pause the contracts,
// and configure different keys.
func setBroadcasterKeys(
	t *testing.T,
	client *ethclient.Client,
	signer blockchain.TransactionSigner,
	parameterRegistryAddress string,
) {
	parameterContract, err := paramReg.NewSettlementChainParameterRegistry(
		common.HexToAddress(parameterRegistryAddress),
		client,
	)
	require.NoError(t, err)

	err = blockchain.ExecuteTransaction(
		t.Context(),
		signer,
		zap.NewNop(),
		client,
		func(opts *bind.TransactOpts) (*types.Transaction, error) {
			keys := []string{
				"xmtp.groupMessageBroadcaster.paused",
				"xmtp.identityUpdateBroadcaster.paused",
				"xmtp.groupMessageBroadcaster.payloadBootstrapper",
				"xmtp.identityUpdateBroadcaster.payloadBootstrapper",
			}

			payerAddressBytes := utils.AddressTo32Slice(common.HexToAddress(PayerAddress))

			values := [][32]byte{
				utils.EncodeUint64ToBytes32(1),
				utils.EncodeUint64ToBytes32(1),
				payerAddressBytes,
				payerAddressBytes,
			}

			return parameterContract.Set0(
				opts,
				keys,
				values,
			)
		},
		func(log *types.Log) (interface{}, error) {
			return parameterContract.ParseParameterSet(*log)
		},
		func(event interface{}) {
			_, ok := event.(*paramReg.SettlementChainParameterRegistryParameterSet)
			require.True(t, ok)
		},
	)
	require.NoError(t, err)
}

func triggerGroupMessageBroadcasterUpdates(
	t *testing.T,
	client *ethclient.Client,
	signer blockchain.TransactionSigner,
	broadcasterAddress string,
) {
	broadcasterContract, err := gm.NewGroupMessageBroadcaster(
		common.HexToAddress(broadcasterAddress),
		client,
	)
	require.NoError(t, err)

	err = blockchain.ExecuteTransaction(
		t.Context(),
		signer,
		zap.NewNop(),
		client,
		func(opts *bind.TransactOpts) (*types.Transaction, error) {
			return broadcasterContract.UpdatePauseStatus(opts)
		},
		func(log *types.Log) (interface{}, error) {
			return broadcasterContract.ParsePauseStatusUpdated(*log)
		},
		func(event interface{}) {
			_, ok := event.(*gm.GroupMessageBroadcasterPauseStatusUpdated)
			require.True(t, ok)
		},
	)
	require.NoError(t, err)

	err = blockchain.ExecuteTransaction(
		t.Context(),
		signer,
		zap.NewNop(),
		client,
		func(opts *bind.TransactOpts) (*types.Transaction, error) {
			return broadcasterContract.UpdatePayloadBootstrapper(opts)
		},
		func(log *types.Log) (interface{}, error) {
			return broadcasterContract.ParsePayloadBootstrapperUpdated(*log)
		},
		func(event interface{}) {
			_, ok := event.(*gm.GroupMessageBroadcasterPayloadBootstrapperUpdated)
			require.True(t, ok)
		},
	)
	require.NoError(t, err)
}

func triggerIdentityUpdateBroadcasterUpdates(
	t *testing.T,
	client *ethclient.Client,
	signer blockchain.TransactionSigner,
	broadcasterAddress string,
) {
	broadcasterContract, err := iu.NewIdentityUpdateBroadcaster(
		common.HexToAddress(broadcasterAddress),
		client,
	)
	require.NoError(t, err)

	err = blockchain.ExecuteTransaction(
		t.Context(),
		signer,
		zap.NewNop(),
		client,
		func(opts *bind.TransactOpts) (*types.Transaction, error) {
			return broadcasterContract.UpdatePauseStatus(opts)
		},
		func(log *types.Log) (interface{}, error) {
			return broadcasterContract.ParsePauseStatusUpdated(*log)
		},
		func(event interface{}) {
			_, ok := event.(*iu.IdentityUpdateBroadcasterPauseStatusUpdated)
			require.True(t, ok)
		},
	)
	require.NoError(t, err)

	err = blockchain.ExecuteTransaction(
		t.Context(),
		signer,
		zap.NewNop(),
		client,
		func(opts *bind.TransactOpts) (*types.Transaction, error) {
			return broadcasterContract.UpdatePayloadBootstrapper(opts)
		},
		func(log *types.Log) (interface{}, error) {
			return broadcasterContract.ParsePayloadBootstrapperUpdated(*log)
		},
		func(event interface{}) {
			_, ok := event.(*iu.IdentityUpdateBroadcasterPayloadBootstrapperUpdated)
			require.True(t, ok)
		},
	)
	require.NoError(t, err)
}
