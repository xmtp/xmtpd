package blockchain_test

import (
	"context"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/xmtp/xmtpd/pkg/blockchain"
	"github.com/xmtp/xmtpd/pkg/testutils"
	"github.com/xmtp/xmtpd/pkg/testutils/anvil"
)

func buildAppChainAdmin(t *testing.T) (blockchain.IAppChainAdmin, *blockchain.ParameterAdmin) {
	t.Helper()

	ctx := context.Background()
	logger := testutils.NewLog(t)
	wsURL, rpcURL := anvil.StartAnvil(t, false)
	contractsOptions := testutils.NewContractsOptions(t, rpcURL, wsURL)

	signer, err := blockchain.NewPrivateKeySigner(
		testutils.LOCAL_PRIVATE_KEY,
		contractsOptions.AppChain.ChainID,
	)
	require.NoError(t, err)

	client, err := blockchain.NewRPCClient(ctx, contractsOptions.AppChain.RPCURL)
	require.NoError(t, err)

	paramAdmin, err := blockchain.NewParameterAdmin(logger, client, signer, contractsOptions)
	require.NoError(t, err)

	appAdmin, err := blockchain.NewAppChainAdmin(
		logger,
		client,
		signer,
		contractsOptions,
		paramAdmin,
	)
	require.NoError(t, err)

	return appAdmin, paramAdmin
}

func TestBootstrapperAddress(t *testing.T) {
	appAdmin, paramAdmin := buildAppChainAdmin(t)
	ctx := context.Background()

	type addrCase struct {
		name string
		key  string
		set  func(ctx context.Context, a common.Address) error
		get  func(ctx context.Context) (common.Address, error)
	}

	cases := []addrCase{
		{
			name: "identity",
			key:  blockchain.IDENTITY_UPDATE_PAYLOAD_BOOTSTRAPPER_KEY,
			set:  appAdmin.SetIdentityUpdateBootstrapper,
			get:  appAdmin.GetIdentityUpdateBootstrapper,
		},
		{
			name: "group",
			key:  blockchain.GROUP_MESSAGE_PAYLOAD_BOOTSTRAPPER_KEY,
			set:  appAdmin.SetGroupMessageBootstrapper,
			get:  appAdmin.GetGroupMessageBootstrapper,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name+"/roundtrip", func(t *testing.T) {
			want := common.HexToAddress("0x000000000000000000000000000000000000BEEF")

			require.NoError(t, tc.set(ctx, want))

			// storage sanity
			stored, err := paramAdmin.GetParameterAddress(ctx, tc.key)
			require.NoError(t, err)
			require.Equal(t, want, stored)

			// getter
			got, err := tc.get(ctx)
			require.NoError(t, err)
			require.Equal(t, want, got)
		})

		t.Run(tc.name+"/overwrite", func(t *testing.T) {
			first := common.HexToAddress("0x000000000000000000000000000000000000CAFE")
			second := common.HexToAddress("0x000000000000000000000000000000000000BEEF")

			require.NoError(t, tc.set(ctx, first))
			require.NoError(t, tc.set(ctx, second))

			stored, err := paramAdmin.GetParameterAddress(ctx, tc.key)
			require.NoError(t, err)
			require.Equal(t, second, stored)

			got, err := tc.get(ctx)
			require.NoError(t, err)
			require.Equal(t, second, got)
		})

		t.Run(tc.name+"/unset_returns_zero", func(t *testing.T) {
			newAppAdmin, _ := buildAppChainAdmin(t)
			got, err := func() (common.Address, error) {
				switch tc.name {
				case "identity":
					return newAppAdmin.GetIdentityUpdateBootstrapper(ctx)
				default:
					return newAppAdmin.GetGroupMessageBootstrapper(ctx)
				}
			}()
			require.NoError(t, err)
			require.Equal(t, common.Address{}, got)
		})

		t.Run(tc.name+"/zero_address_roundtrip", func(t *testing.T) {
			var zero common.Address
			require.NoError(t, tc.set(ctx, zero))
			got, err := tc.get(ctx)
			require.NoError(t, err)
			require.Equal(t, zero, got)
		})
	}
}

func TestPauseFlags(t *testing.T) {
	appAdmin, paramAdmin := buildAppChainAdmin(t)
	ctx := context.Background()

	type pauseCase struct {
		name string
		key  string
		set  func(ctx context.Context, paused bool) error
		get  func(ctx context.Context) (bool, error)
	}

	cases := []pauseCase{
		{
			name: "group",
			key:  blockchain.GROUP_MESSAGE_PAUSED_KEY,
			set:  appAdmin.SetGroupMessagePauseStatus,
			get:  appAdmin.GetGroupMessagePauseStatus,
		},
		{
			name: "identity",
			key:  blockchain.IDENTITY_UPDATE_PAUSED_KEY,
			set:  appAdmin.SetIdentityUpdatePauseStatus,
			get:  appAdmin.GetIdentityUpdatePauseStatus,
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name+"/toggle_true_false", func(t *testing.T) {
			require.NoError(t, tc.set(ctx, true))
			b, err := paramAdmin.GetParameterBool(ctx, tc.key)
			require.NoError(t, err)
			require.True(t, b)

			got, err := tc.get(ctx)
			require.NoError(t, err)
			require.True(t, got)

			require.NoError(t, tc.set(ctx, false))
			b, err = paramAdmin.GetParameterBool(ctx, tc.key)
			require.NoError(t, err)
			require.False(t, b)

			got, err = tc.get(ctx)
			require.NoError(t, err)
			require.False(t, got)
		})

		t.Run(tc.name+"/idempotent_repeat_true", func(t *testing.T) {
			require.NoError(t, tc.set(ctx, true))
			require.NoError(t, tc.set(ctx, true))

			got, err := tc.get(ctx)
			require.NoError(t, err)
			require.True(t, got)
		})

		t.Run(tc.name+"/getter_unset_returns_false", func(t *testing.T) {
			newAppAdmin, newParamAdmin := buildAppChainAdmin(t)

			var got bool
			var err error
			switch tc.name {
			case "group":
				got, err = newAppAdmin.GetGroupMessagePauseStatus(ctx)
			default:
				got, err = newAppAdmin.GetIdentityUpdatePauseStatus(ctx)
			}
			require.NoError(t, err)
			require.False(t, got)

			b, err := newParamAdmin.GetParameterBool(ctx, tc.key)
			require.NoError(t, err)
			require.False(t, b)
		})
	}
}
