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

func buildParameterAdmin(t *testing.T) *blockchain.ParameterAdmin {
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

	admin, err := blockchain.NewParameterAdmin(logger, client, signer, contractsOptions)
	require.NoError(t, err)

	return admin
}

func TestSetUint8ParameterAndReadBack(t *testing.T) {
	admin := buildParameterAdmin(t)
	ctx := context.Background()

	// Use the uint8-style key the code exposes.
	const key = blockchain.NODE_REGISTRY_MAX_CANONICAL_NODES_KEY
	const want uint8 = 7

	err := admin.SetUint8Parameter(ctx, key, want)
	require.NoError(t, err)

	got, err := admin.GetParameterUint8(ctx, key)
	require.NoError(t, err)

	// Value should be right-aligned (big-endian), last byte holds our uint8.
	require.Equal(t, want, got, "expected last byte to equal the uint8 value")
}

func TestSetUint8ParameterCanOverwrite(t *testing.T) {
	admin := buildParameterAdmin(t)
	ctx := context.Background()

	const key = blockchain.NODE_REGISTRY_MAX_CANONICAL_NODES_KEY

	err := admin.SetUint8Parameter(ctx, key, 1)
	require.NoError(t, err)

	err = admin.SetUint8Parameter(ctx, key, 255)
	require.NoError(t, err)

	got, err := admin.GetParameterUint8(ctx, key)
	require.NoError(t, err)
	require.EqualValues(t, 255, got)
}

func TestSetAddressParameterAndReadBack(t *testing.T) {
	admin := buildParameterAdmin(t)
	ctx := context.Background()

	const key = blockchain.IDENTITY_UPDATE_PAYLOAD_BOOTSTRAPPER_KEY

	// Any address is fine for storage; use a memorable sentinel.
	wantAddr := common.HexToAddress("0x000000000000000000000000000000000000BEEF")

	err := admin.SetAddressParameter(ctx, key, wantAddr)
	require.NoError(t, err)

	addr, err := admin.GetParameterAddress(ctx, key)
	require.NoError(t, err)
	require.Equal(t, wantAddr, addr)
}

func TestSetAddressParameterCanOverwrite(t *testing.T) {
	admin := buildParameterAdmin(t)
	ctx := context.Background()

	const key = blockchain.IDENTITY_UPDATE_PAYLOAD_BOOTSTRAPPER_KEY

	first := common.HexToAddress("0x000000000000000000000000000000000000CAFE")
	second := common.HexToAddress("0x000000000000000000000000000000000000BEEF")

	err := admin.SetAddressParameter(ctx, key, first)
	require.NoError(t, err)

	err = admin.SetAddressParameter(ctx, key, second)
	require.NoError(t, err)

	addr, err := admin.GetParameterAddress(ctx, key)
	require.NoError(t, err)
	require.Equal(t, second, addr)
}

func TestGetBoolParameter_Unset_ReturnsFalse(t *testing.T) {
	admin := buildParameterAdmin(t)
	ctx := context.Background()

	got, err := admin.GetParameterBool(ctx, blockchain.IDENTITY_UPDATE_PAUSED_KEY)
	require.NoError(t, err)
	require.False(t, got)
}

func TestSetBoolParameter_True_RoundTrip(t *testing.T) {
	admin := buildParameterAdmin(t)
	ctx := context.Background()

	err := admin.SetBoolParameter(ctx, blockchain.IDENTITY_UPDATE_PAUSED_KEY, true)
	require.NoError(t, err)

	got, err := admin.GetParameterBool(ctx, blockchain.IDENTITY_UPDATE_PAUSED_KEY)
	require.NoError(t, err)
	require.True(t, got)
}

func TestSetBoolParameter_False_RoundTrip(t *testing.T) {
	admin := buildParameterAdmin(t)
	ctx := context.Background()

	err := admin.SetBoolParameter(ctx, blockchain.IDENTITY_UPDATE_PAUSED_KEY, false)
	require.NoError(t, err)

	got, err := admin.GetParameterBool(ctx, blockchain.IDENTITY_UPDATE_PAUSED_KEY)
	require.NoError(t, err)
	require.False(t, got)
}

func TestSetBoolParameter_True_NoOp(t *testing.T) {
	admin := buildParameterAdmin(t)
	ctx := context.Background()

	err := admin.SetBoolParameter(ctx, blockchain.IDENTITY_UPDATE_PAUSED_KEY, true)
	require.NoError(t, err)

	err = admin.SetBoolParameter(ctx, blockchain.IDENTITY_UPDATE_PAUSED_KEY, true)
	require.NoError(t, err)
}

func TestSetBoolParameter_False_NoOp(t *testing.T) {
	admin := buildParameterAdmin(t)
	ctx := context.Background()

	err := admin.SetBoolParameter(ctx, blockchain.IDENTITY_UPDATE_PAUSED_KEY, false)
	require.NoError(t, err)

	err = admin.SetBoolParameter(ctx, blockchain.IDENTITY_UPDATE_PAUSED_KEY, false)
	require.NoError(t, err)
}

func TestSetBoolParameter_Unset(t *testing.T) {
	admin := buildParameterAdmin(t)
	ctx := context.Background()

	err := admin.SetBoolParameter(ctx, blockchain.IDENTITY_UPDATE_PAUSED_KEY, true)
	require.NoError(t, err)

	err = admin.SetBoolParameter(ctx, blockchain.IDENTITY_UPDATE_PAUSED_KEY, false)
	require.NoError(t, err)

	got, err := admin.GetParameterBool(ctx, blockchain.IDENTITY_UPDATE_PAUSED_KEY)
	require.NoError(t, err)
	require.False(t, got)
}
