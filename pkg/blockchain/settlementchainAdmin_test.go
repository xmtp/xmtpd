package blockchain_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/xmtp/xmtpd/pkg/blockchain"
	"github.com/xmtp/xmtpd/pkg/testutils"
	"github.com/xmtp/xmtpd/pkg/testutils/anvil"
)

func buildSettlementChainAdmin(
	t *testing.T,
) (blockchain.ISettlementChainAdmin, *blockchain.ParameterAdmin) {
	t.Helper()

	ctx := context.Background()
	logger := testutils.NewLog(t)
	wsURL, rpcURL := anvil.StartAnvil(t, false)
	contractsOptions := testutils.NewContractsOptions(t, rpcURL, wsURL)

	signer, err := blockchain.NewPrivateKeySigner(
		testutils.LOCAL_PRIVATE_KEY,
		contractsOptions.SettlementChain.ChainID,
	)
	require.NoError(t, err)

	client, err := blockchain.NewRPCClient(ctx, contractsOptions.SettlementChain.RPCURL)
	require.NoError(t, err)

	paramAdmin, err := blockchain.NewParameterAdmin(logger, client, signer, contractsOptions)
	require.NoError(t, err)

	appAdmin, err := blockchain.NewSettlementChainAdmin(
		logger,
		client,
		signer,
		contractsOptions,
		paramAdmin,
	)
	require.NoError(t, err)

	return appAdmin, paramAdmin
}

func TestPauseFlagsSettlement(t *testing.T) {
	appAdmin, paramAdmin := buildSettlementChainAdmin(t)
	ctx := context.Background()

	type pauseCase struct {
		name string
		key  string
		set  func(ctx context.Context, paused bool) error
		get  func(ctx context.Context) (bool, error)
	}

	cases := []pauseCase{
		{
			name: "settlement-chain-gateway",
			key:  blockchain.SETTLEMENT_CHAIN_GATEWAY_PAUSED_KEY,
			set:  appAdmin.SetSettlementChainGatewayPauseStatus,
			get:  appAdmin.GetSettlementChainGatewayPauseStatus,
		},
		{
			name: "payer-registry",
			key:  blockchain.PAYER_REGISTRY_PAUSED_KEY,
			set:  appAdmin.SetPayerRegistryPauseStatus,
			get:  appAdmin.GetPayerRegistryPauseStatus,
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name+"/toggle_true_false", func(t *testing.T) {
			var err error
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
		//
		//t.Run(tc.name+"/getter_unset_returns_false", func(t *testing.T) {
		//	newAppAdmin, newParamAdmin := buildSettlementChainAdmin(t)
		//
		//	var got bool
		//	var err error
		//	switch tc.name {
		//	case "group":
		//		got, err = newAppAdmin.GetGroupMessagePauseStatus(ctx)
		//	default:
		//		got, err = newAppAdmin.GetIdentityUpdatePauseStatus(ctx)
		//	}
		//	require.NoError(t, err)
		//	require.False(t, got)
		//
		//	b, err := newParamAdmin.GetParameterBool(ctx, tc.key)
		//	require.NoError(t, err)
		//	require.False(t, b)
		//})
	}
}
