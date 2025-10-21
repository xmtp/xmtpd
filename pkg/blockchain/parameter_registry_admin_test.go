package blockchain_test

import (
	"context"
	"math"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/xmtp/xmtpd/pkg/blockchain"
	"github.com/xmtp/xmtpd/pkg/testutils"
	"github.com/xmtp/xmtpd/pkg/testutils/anvil"
)

func buildParameterAdmin(t *testing.T) blockchain.IParameterAdmin {
	t.Helper()

	ctx := context.Background()
	logger := testutils.NewLog(t)
	wsURL, rpcURL := anvil.StartAnvil(t, false)
	contractsOptions := testutils.NewContractsOptions(t, rpcURL, wsURL)

	signer, err := blockchain.NewPrivateKeySigner(
		testutils.TestPrivateKey,
		contractsOptions.AppChain.ChainID,
	)
	require.NoError(t, err)

	client, err := blockchain.NewRPCClient(ctx, contractsOptions.AppChain.RPCURL)
	require.NoError(t, err)

	admin, err := blockchain.NewSettlementParameterAdmin(logger, client, signer, contractsOptions)
	require.NoError(t, err)

	return admin
}

func TestSetUint8ParameterAndReadBack(t *testing.T) {
	admin := buildParameterAdmin(t)
	ctx := context.Background()

	// Use the uint8-style key the code exposes.
	const key = blockchain.NodeRegistryMaxCanonicalNodesKey
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

	const key = blockchain.NodeRegistryMaxCanonicalNodesKey

	err := admin.SetUint8Parameter(ctx, key, 1)
	require.NoError(t, err)

	err = admin.SetUint8Parameter(ctx, key, 255)
	require.NoError(t, err)

	got, err := admin.GetParameterUint8(ctx, key)
	require.NoError(t, err)
	require.EqualValues(t, 255, got)
}

func TestSetUint16ParameterAndReadBack(t *testing.T) {
	admin := buildParameterAdmin(t)
	ctx := context.Background()

	const key = blockchain.PayerReportManagerProtocolFeeRateKey // uint16 param
	const want uint16 = 12345

	err := admin.SetUint16Parameter(ctx, key, want)
	require.NoError(t, err)

	got, err := admin.GetParameterUint16(ctx, key)
	require.NoError(t, err)
	require.EqualValues(t, want, got)
}

func TestSetUint16Parameter_ZeroAndMax(t *testing.T) {
	admin := buildParameterAdmin(t)
	ctx := context.Background()

	const key = blockchain.PayerReportManagerProtocolFeeRateKey

	// zero
	require.NoError(t, admin.SetUint16Parameter(ctx, key, 0))
	got0, err := admin.GetParameterUint16(ctx, key)
	require.NoError(t, err)
	require.EqualValues(t, 0, got0)

	// max
	require.NoError(t, admin.SetUint16Parameter(ctx, key, math.MaxUint16))
	gotMax, err := admin.GetParameterUint16(ctx, key)
	require.NoError(t, err)
	require.EqualValues(t, math.MaxUint16, gotMax)
}

func TestSetUint16ParameterCanOverwrite(t *testing.T) {
	admin := buildParameterAdmin(t)
	ctx := context.Background()

	const key = blockchain.PayerReportManagerProtocolFeeRateKey

	require.NoError(t, admin.SetUint16Parameter(ctx, key, 1))
	require.NoError(t, admin.SetUint16Parameter(ctx, key, 65535))

	got, err := admin.GetParameterUint16(ctx, key)
	require.NoError(t, err)
	require.EqualValues(t, 65535, got)
}

func TestSetUint32ParameterAndReadBack(t *testing.T) {
	admin := buildParameterAdmin(t)
	ctx := context.Background()

	// Use any key; the registry stores bytes32 agnostically.
	// We'll use a rate key in a fresh chain instance for isolation.
	const key = blockchain.RateRegistryCongestionFeeKey
	const want uint32 = 3_000_000_001

	err := admin.SetUint32Parameter(ctx, key, want)
	require.NoError(t, err)

	got, err := admin.GetParameterUint32(ctx, key)
	require.NoError(t, err)
	require.EqualValues(t, want, got)
}

func TestSetUint32Parameter_ZeroAndMax(t *testing.T) {
	admin := buildParameterAdmin(t)
	ctx := context.Background()

	const key = blockchain.RateRegistryStorageFeeKey // fresh chain per test

	// zero
	require.NoError(t, admin.SetUint32Parameter(ctx, key, 0))
	got0, err := admin.GetParameterUint32(ctx, key)
	require.NoError(t, err)
	require.EqualValues(t, 0, got0)

	// max
	require.NoError(t, admin.SetUint32Parameter(ctx, key, math.MaxUint32))
	gotMax, err := admin.GetParameterUint32(ctx, key)
	require.NoError(t, err)
	require.EqualValues(t, math.MaxUint32, gotMax)
}

func TestSetUint32ParameterCanOverwrite(t *testing.T) {
	admin := buildParameterAdmin(t)
	ctx := context.Background()

	const key = blockchain.RateRegistryMessageFeeKey

	require.NoError(t, admin.SetUint32Parameter(ctx, key, 42))
	require.NoError(t, admin.SetUint32Parameter(ctx, key, 99_999_999))

	got, err := admin.GetParameterUint32(ctx, key)
	require.NoError(t, err)
	require.EqualValues(t, 99_999_999, got)
}

func TestSetAddressParameterAndReadBack(t *testing.T) {
	admin := buildParameterAdmin(t)
	ctx := context.Background()

	const key = blockchain.IdentityUpdatePayloadBootstrapperKey

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

	const key = blockchain.IdentityUpdatePayloadBootstrapperKey

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

	got, err := admin.GetParameterBool(ctx, blockchain.IdentityUpdateBroadcasterPausedKey)
	require.NoError(t, err)
	require.False(t, got)
}

func TestSetBoolParameter_True_RoundTrip(t *testing.T) {
	admin := buildParameterAdmin(t)
	ctx := context.Background()

	err := admin.SetBoolParameter(ctx, blockchain.IdentityUpdateBroadcasterPausedKey, true)
	require.NoError(t, err)

	got, err := admin.GetParameterBool(ctx, blockchain.IdentityUpdateBroadcasterPausedKey)
	require.NoError(t, err)
	require.True(t, got)
}

func TestSetBoolParameter_False_RoundTrip(t *testing.T) {
	admin := buildParameterAdmin(t)
	ctx := context.Background()

	err := admin.SetBoolParameter(ctx, blockchain.IdentityUpdateBroadcasterPausedKey, false)
	require.NoError(t, err)

	got, err := admin.GetParameterBool(ctx, blockchain.IdentityUpdateBroadcasterPausedKey)
	require.NoError(t, err)
	require.False(t, got)
}

func TestSetBoolParameter_True_NoOp(t *testing.T) {
	admin := buildParameterAdmin(t)
	ctx := context.Background()

	err := admin.SetBoolParameter(ctx, blockchain.IdentityUpdateBroadcasterPausedKey, true)
	require.NoError(t, err)

	err = admin.SetBoolParameter(ctx, blockchain.IdentityUpdateBroadcasterPausedKey, true)
	require.NoError(t, err)
}

func TestSetBoolParameter_False_NoOp(t *testing.T) {
	admin := buildParameterAdmin(t)
	ctx := context.Background()

	err := admin.SetBoolParameter(ctx, blockchain.IdentityUpdateBroadcasterPausedKey, false)
	require.NoError(t, err)

	err = admin.SetBoolParameter(ctx, blockchain.IdentityUpdateBroadcasterPausedKey, false)
	require.NoError(t, err)
}

func TestSetBoolParameter_Unset(t *testing.T) {
	admin := buildParameterAdmin(t)
	ctx := context.Background()

	err := admin.SetBoolParameter(ctx, blockchain.IdentityUpdateBroadcasterPausedKey, true)
	require.NoError(t, err)

	err = admin.SetBoolParameter(ctx, blockchain.IdentityUpdateBroadcasterPausedKey, false)
	require.NoError(t, err)

	got, err := admin.GetParameterBool(ctx, blockchain.IdentityUpdateBroadcasterPausedKey)
	require.NoError(t, err)
	require.False(t, got)
}

func TestParameterAdmin_SetManyUint64Parameters_RoundTrip(t *testing.T) {
	admin := buildParameterAdmin(t)
	ctx := context.Background()

	items := []blockchain.Uint64Param{
		{Name: blockchain.RateRegistryMessageFeeKey, Value: 123},
		{Name: blockchain.RateRegistryStorageFeeKey, Value: 456},
		{Name: blockchain.RateRegistryCongestionFeeKey, Value: 789},
		{Name: blockchain.RateRegistryTargetRatePerMinuteKey, Value: 42},
	}

	err := admin.SetManyUint64Parameters(ctx, items)
	require.NoError(t, err)

	gotMsg, err := admin.GetParameterUint64(ctx, blockchain.RateRegistryMessageFeeKey)
	require.NoError(t, err)
	require.EqualValues(t, 123, gotMsg)

	gotStore, err := admin.GetParameterUint64(ctx, blockchain.RateRegistryStorageFeeKey)
	require.NoError(t, err)
	require.EqualValues(t, 456, gotStore)

	gotCong, err := admin.GetParameterUint64(ctx, blockchain.RateRegistryCongestionFeeKey)
	require.NoError(t, err)
	require.EqualValues(t, 789, gotCong)

	gotTarget, err := admin.GetParameterUint64(
		ctx,
		blockchain.RateRegistryTargetRatePerMinuteKey,
	)
	require.NoError(t, err)
	require.EqualValues(t, 42, gotTarget)
}

func TestParameterAdmin_SetManyUint64Parameters_Overwrite(t *testing.T) {
	admin := buildParameterAdmin(t)
	ctx := context.Background()

	first := []blockchain.Uint64Param{
		{Name: blockchain.RateRegistryMessageFeeKey, Value: 1},
		{Name: blockchain.RateRegistryStorageFeeKey, Value: 2},
		{Name: blockchain.RateRegistryCongestionFeeKey, Value: 3},
		{Name: blockchain.RateRegistryTargetRatePerMinuteKey, Value: 4},
	}
	require.NoError(t, admin.SetManyUint64Parameters(ctx, first))

	second := []blockchain.Uint64Param{
		{Name: blockchain.RateRegistryMessageFeeKey, Value: 11},
		{Name: blockchain.RateRegistryStorageFeeKey, Value: 22},
		{Name: blockchain.RateRegistryCongestionFeeKey, Value: 33},
		{Name: blockchain.RateRegistryTargetRatePerMinuteKey, Value: 44},
	}
	require.NoError(t, admin.SetManyUint64Parameters(ctx, second))

	gotMsg, err := admin.GetParameterUint64(ctx, blockchain.RateRegistryMessageFeeKey)
	require.NoError(t, err)
	require.EqualValues(t, 11, gotMsg)

	gotStore, err := admin.GetParameterUint64(ctx, blockchain.RateRegistryStorageFeeKey)
	require.NoError(t, err)
	require.EqualValues(t, 22, gotStore)

	gotCong, err := admin.GetParameterUint64(ctx, blockchain.RateRegistryCongestionFeeKey)
	require.NoError(t, err)
	require.EqualValues(t, 33, gotCong)

	gotTarget, err := admin.GetParameterUint64(
		ctx,
		blockchain.RateRegistryTargetRatePerMinuteKey,
	)
	require.NoError(t, err)
	require.EqualValues(t, 44, gotTarget)
}

func TestSetUint64ParameterAndReadBack(t *testing.T) {
	admin := buildParameterAdmin(t)
	ctx := context.Background()

	const key = blockchain.RateRegistryMessageFeeKey
	const want uint64 = 123456789

	err := admin.SetUint64Parameter(ctx, key, want)
	require.NoError(t, err)

	got, err := admin.GetParameterUint64(ctx, key)
	require.NoError(t, err)
	require.EqualValues(t, want, got)
}

func TestSetUint64ParameterCanOverwrite(t *testing.T) {
	admin := buildParameterAdmin(t)
	ctx := context.Background()

	const key = blockchain.RateRegistryStorageFeeKey

	err := admin.SetUint64Parameter(ctx, key, 42)
	require.NoError(t, err)

	err = admin.SetUint64Parameter(ctx, key, 99)
	require.NoError(t, err)

	got, err := admin.GetParameterUint64(ctx, key)
	require.NoError(t, err)
	require.EqualValues(t, 99, got)
}
