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

func TestSetAndGetIdentityUpdateBootstrapper_RoundTrip(t *testing.T) {
	appAdmin, paramAdmin := buildAppChainAdmin(t)
	ctx := context.Background()

	// Use a distinctive non-zero address so alignment issues are obvious.
	want := common.HexToAddress("0x000000000000000000000000000000000000BEEF")

	// Set via public API
	err := appAdmin.SetIdentityUpdateBootstrapper(ctx, want)
	require.NoError(t, err)

	_, err = paramAdmin.GetParameterAddress(
		ctx,
		blockchain.IDENTITY_UPDATE_PAYLOAD_BOOTSTRAPPER_KEY,
	)
	require.NoError(t, err)

	// Now read it back through the getter
	got, err := appAdmin.GetIdentityUpdateBootstrapper(ctx)
	require.NoError(t, err)
	require.Equal(t, want, got)
}

func TestSetIdentityUpdateBootstrapper_Overwrite(t *testing.T) {
	appAdmin, paramAdmin := buildAppChainAdmin(t)
	ctx := context.Background()

	first := common.HexToAddress("0x000000000000000000000000000000000000CAFE")
	second := common.HexToAddress("0x000000000000000000000000000000000000BEEF")

	err := appAdmin.SetIdentityUpdateBootstrapper(ctx, first)
	require.NoError(t, err)

	err = appAdmin.SetIdentityUpdateBootstrapper(ctx, second)
	require.NoError(t, err)

	// Raw storage should contain the latest value
	addr, err := paramAdmin.GetParameterAddress(
		ctx,
		blockchain.IDENTITY_UPDATE_PAYLOAD_BOOTSTRAPPER_KEY,
	)
	require.NoError(t, err)
	require.Equal(t, second, addr)

	// Getter should reflect latest value too
	got, err := appAdmin.GetIdentityUpdateBootstrapper(ctx)
	require.NoError(t, err)
	require.Equal(t, second, got)
}

func TestGetIdentityUpdateBootstrapper_Unset_ReturnsZeroAddress(t *testing.T) {
	appAdmin, _ := buildAppChainAdmin(t)
	ctx := context.Background()

	got, err := appAdmin.GetIdentityUpdateBootstrapper(ctx)
	require.NoError(t, err)
	require.Equal(t, (common.Address{}), got, "unset param should decode to 0x000...000")
}

func TestSetIdentityUpdateBootstrapper_ZeroAddress_RoundTrip(t *testing.T) {
	appAdmin, _ := buildAppChainAdmin(t)
	ctx := context.Background()

	var zero common.Address
	err := appAdmin.SetIdentityUpdateBootstrapper(ctx, zero)
	require.NoError(t, err)

	got, err := appAdmin.GetIdentityUpdateBootstrapper(ctx)
	require.NoError(t, err)
	require.Equal(t, zero, got, "zero address should round-trip unless validation is added")
}

func TestSetAndGetGroupMessageBootstrapper_RoundTrip(t *testing.T) {
	appAdmin, paramAdmin := buildAppChainAdmin(t)
	ctx := context.Background()

	// Use a distinctive non-zero address so alignment issues are obvious.
	want := common.HexToAddress("0x000000000000000000000000000000000000BEEF")

	// Set via public API
	err := appAdmin.SetGroupMessageBootstrapper(ctx, want)
	require.NoError(t, err)

	_, err = paramAdmin.GetParameterAddress(
		ctx,
		blockchain.GROUP_MESSAGE_PAYLOAD_BOOTSTRAPPER_KEY,
	)
	require.NoError(t, err)

	// Now read it back through the getter
	got, err := appAdmin.GetGroupMessageBootstrapper(ctx)
	require.NoError(t, err)
	require.Equal(t, want, got)
}

func TestSetGroupMessageBootstrapper_Overwrite(t *testing.T) {
	appAdmin, paramAdmin := buildAppChainAdmin(t)
	ctx := context.Background()

	first := common.HexToAddress("0x000000000000000000000000000000000000CAFE")
	second := common.HexToAddress("0x000000000000000000000000000000000000BEEF")

	err := appAdmin.SetGroupMessageBootstrapper(ctx, first)
	require.NoError(t, err)

	err = appAdmin.SetGroupMessageBootstrapper(ctx, second)
	require.NoError(t, err)

	// Raw storage should contain the latest value
	addr, err := paramAdmin.GetParameterAddress(
		ctx,
		blockchain.GROUP_MESSAGE_PAYLOAD_BOOTSTRAPPER_KEY,
	)
	require.NoError(t, err)
	require.Equal(t, second, addr)

	// Getter should reflect latest value too
	got, err := appAdmin.GetGroupMessageBootstrapper(ctx)
	require.NoError(t, err)
	require.Equal(t, second, got)
}

func TestGetGroupMessageBootstrapper_Unset_ReturnsZeroAddress(t *testing.T) {
	appAdmin, _ := buildAppChainAdmin(t)
	ctx := context.Background()

	got, err := appAdmin.GetGroupMessageBootstrapper(ctx)
	require.NoError(t, err)
	require.Equal(t, (common.Address{}), got, "unset param should decode to 0x000...000")
}

func TestSetGroupMessageBootstrapper_ZeroAddress_RoundTrip(t *testing.T) {
	appAdmin, _ := buildAppChainAdmin(t)
	ctx := context.Background()

	var zero common.Address
	err := appAdmin.SetGroupMessageBootstrapper(ctx, zero)
	require.NoError(t, err)

	got, err := appAdmin.GetGroupMessageBootstrapper(ctx)
	require.NoError(t, err)
	require.Equal(t, zero, got, "zero address should round-trip unless validation is added")
}
