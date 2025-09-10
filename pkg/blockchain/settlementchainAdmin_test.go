package blockchain_test

import (
	"context"
	"errors"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"

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

	settlementChainAdmin, err := blockchain.NewSettlementChainAdmin(
		logger,
		client,
		signer,
		contractsOptions,
		paramAdmin,
	)
	require.NoError(t, err)

	return settlementChainAdmin, paramAdmin
}

func TestPauseFlagsSettlement(t *testing.T) {
	settlementChainAdmin, paramAdmin := buildSettlementChainAdmin(t)
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
			set:  settlementChainAdmin.SetSettlementChainGatewayPauseStatus,
			get:  settlementChainAdmin.GetSettlementChainGatewayPauseStatus,
		},
		{
			name: "payer-registry",
			key:  blockchain.PAYER_REGISTRY_PAUSED_KEY,
			set:  settlementChainAdmin.SetPayerRegistryPauseStatus,
			get:  settlementChainAdmin.GetPayerRegistryPauseStatus,
		},
		{
			name: "distribution-manager",
			key:  blockchain.DISTRIBUTION_MANAGER_PAUSED_KEY,
			set:  settlementChainAdmin.SetDistributionManagerPauseStatus,
			get:  settlementChainAdmin.GetDistributionManagerPauseStatus,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
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

			t.Run(tc.name+"/getter_unset_returns_false", func(t *testing.T) {
				newSettlementAdmin, newParamAdmin := buildSettlementChainAdmin(t)

				var got bool
				var err error
				switch tc.name {
				case "settlement-chain-gateway":
					got, err = newSettlementAdmin.GetSettlementChainGatewayPauseStatus(ctx)
				case "payer-registry":
					got, err = newSettlementAdmin.GetPayerRegistryPauseStatus(ctx)
				case "distribution-manager":
					got, err = newSettlementAdmin.GetDistributionManagerPauseStatus(ctx)
				default:
					got, err = false, errors.New("invalid option")
				}
				require.NoError(t, err)
				require.False(t, got)

				b, err := newParamAdmin.GetParameterBool(ctx, tc.key)
				require.NoError(t, err)
				require.False(t, b)
			})
		})
	}
}

func TestDistributionManager_ProtocolFeesRecipient_ReadDefault_WriteThenRead(t *testing.T) {
	scAdmin, paramAdmin := buildSettlementChainAdmin(t)
	ctx := context.Background()

	// default is zero address
	gotDefault, err := scAdmin.GetDistributionManagerProtocolFeesRecipient(ctx)
	require.NoError(t, err)
	require.Equal(t, common.Address{}, gotDefault)

	rawDefault, err := paramAdmin.GetParameterAddress(
		ctx,
		blockchain.DISTRIBUTION_MANAGER_PROTOCOL_FEES_RECIPIENT_KEY,
	)
	require.NoError(t, err)
	require.Equal(t, common.Address{}, rawDefault)

	// set a value
	want := common.HexToAddress("0x000000000000000000000000000000000000BEEF")
	require.NoError(t, scAdmin.SetDistributionManagerProtocolFeesRecipient(ctx, want))

	// read back
	got, err := scAdmin.GetDistributionManagerProtocolFeesRecipient(ctx)
	require.NoError(t, err)
	require.Equal(t, want, got)

	raw, err := paramAdmin.GetParameterAddress(
		ctx,
		blockchain.DISTRIBUTION_MANAGER_PROTOCOL_FEES_RECIPIENT_KEY,
	)
	require.NoError(t, err)
	require.Equal(t, want, raw)

	// overwrite + idempotent
	other := common.HexToAddress("0x000000000000000000000000000000000000CAFE")
	require.NoError(t, scAdmin.SetDistributionManagerProtocolFeesRecipient(ctx, other))
	require.NoError(t, scAdmin.SetDistributionManagerProtocolFeesRecipient(ctx, other))

	got2, err := scAdmin.GetDistributionManagerProtocolFeesRecipient(ctx)
	require.NoError(t, err)
	require.Equal(t, other, got2)

	// zero round-trip
	var zero common.Address
	require.NoError(t, scAdmin.SetDistributionManagerProtocolFeesRecipient(ctx, zero))
	gotZero, err := scAdmin.GetDistributionManagerProtocolFeesRecipient(ctx)
	require.NoError(t, err)
	require.Equal(t, zero, gotZero)
}

func TestNodeRegistry_Admin_ReadDefault_WriteThenRead(t *testing.T) {
	scAdmin, paramAdmin := buildSettlementChainAdmin(t)
	ctx := context.Background()

	// set + read
	want := common.HexToAddress("0x000000000000000000000000000000000000DEAD")
	require.NoError(t, scAdmin.SetNodeRegistryAdmin(ctx, want))

	got, err := scAdmin.GetNodeRegistryAdmin(ctx)
	require.NoError(t, err)
	require.Equal(t, want, got)

	raw, err := paramAdmin.GetParameterAddress(ctx, blockchain.NODE_REGISTRY_ADMIN_KEY)
	require.NoError(t, err)
	require.Equal(t, want, raw)

	// overwrite + zero
	other := common.HexToAddress("0x000000000000000000000000000000000000BEEF")
	require.NoError(t, scAdmin.SetNodeRegistryAdmin(ctx, other))
	var zero common.Address
	require.NoError(t, scAdmin.SetNodeRegistryAdmin(ctx, zero))
	gotZero, err := scAdmin.GetNodeRegistryAdmin(ctx)
	require.NoError(t, err)
	require.Equal(t, zero, gotZero)
}

func TestPayerRegistry_Uint32Params_ReadDefault_WriteThenRead(t *testing.T) {
	scAdmin, paramAdmin := buildSettlementChainAdmin(t)
	ctx := context.Background()

	type u32Case struct {
		name string
		key  string
		set  func(context.Context, uint32) error
		get  func(context.Context) (uint32, error)
	}
	cases := []u32Case{
		{
			name: "withdrawLockPeriod",
			key:  blockchain.PAYER_REGISTRY_WITHDRAW_LOCK_PERIOD_KEY,
			set:  scAdmin.SetPayerRegistryWithdrawLockPeriod,
			get:  scAdmin.GetPayerRegistryWithdrawLockPeriod,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Run(tc.name+"/read_default", func(t *testing.T) {
				t.Skip(
					"Some defaults are pre-set https://github.com/xmtp/smart-contracts/issues/126",
				)
				gotDefault, err := tc.get(ctx)
				require.NoError(t, err)
				require.EqualValues(t, 0, gotDefault)

				rawDefault, err := paramAdmin.GetParameterUint64(ctx, tc.key)
				require.NoError(t, err)
				require.EqualValues(t, 0, rawDefault)
			})

			t.Run(tc.name+"/write_read_back", func(t *testing.T) {
				const v1 = 1024
				require.NoError(t, tc.set(ctx, v1))

				gotV1, err := tc.get(ctx)
				require.NoError(t, err)
				require.EqualValues(t, v1, gotV1)

				rawV1, err := paramAdmin.GetParameterUint64(ctx, tc.key)
				require.NoError(t, err)
				require.EqualValues(t, v1, rawV1)
			})

			t.Run(tc.name+"/write_idempotent", func(t *testing.T) {
				const v1 = 1024
				require.NoError(t, tc.set(ctx, v1))

				require.NoError(t, tc.set(ctx, v1))
			})

			t.Run(tc.name+"/write_back_to_zero", func(t *testing.T) {
				const v1 = 1024
				require.NoError(t, tc.set(ctx, v1))

				const v2 = 0
				err := tc.set(ctx, v2)

				switch tc.name {
				case "withdrawLockPeriod":
					require.NoError(t, err)
				}
			})
		})
	}
}

func TestPayerRegistry_Uint96Params_ReadDefault_WriteThenRead(t *testing.T) {
	scAdmin, paramAdmin := buildSettlementChainAdmin(t)
	ctx := context.Background()

	type u96Case struct {
		name string
		key  string
		set  func(context.Context, *big.Int) error
		get  func(context.Context) (*big.Int, error)
	}
	cases := []u96Case{
		{
			name: "minimumDeposit",
			key:  blockchain.PAYER_REGISTRY_MINIMUM_DEPOSIT_KEY,
			set:  scAdmin.SetPayerRegistryMinimumDeposit,
			get:  scAdmin.GetPayerRegistryMinimumDeposit,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Run(tc.name+"/read_default", func(t *testing.T) {
				t.Skip(
					"Some defaults are pre-set https://github.com/xmtp/smart-contracts/issues/126",
				)
				gotDefault, err := tc.get(ctx)
				require.NoError(t, err)
				require.EqualValues(t, 0, gotDefault)

				rawDefault, err := paramAdmin.GetParameterUint64(ctx, tc.key)
				require.NoError(t, err)
				require.EqualValues(t, 0, rawDefault)
			})

			t.Run(tc.name+"/write_read_back", func(t *testing.T) {
				v1 := big.NewInt(1024)
				require.NoError(t, tc.set(ctx, v1))

				gotV1, err := tc.get(ctx)
				require.NoError(t, err)
				require.EqualValues(t, v1, gotV1)

				rawV1, err := paramAdmin.GetParameterUint96(ctx, tc.key)
				require.NoError(t, err)
				require.EqualValues(t, v1, rawV1)
			})

			t.Run(tc.name+"/write_idempotent", func(t *testing.T) {
				v1 := big.NewInt(1024)
				require.NoError(t, tc.set(ctx, v1))

				require.NoError(t, tc.set(ctx, v1))
			})

			t.Run(tc.name+"/write_back_to_zero", func(t *testing.T) {
				v1 := big.NewInt(1024)
				require.NoError(t, tc.set(ctx, v1))

				v2 := big.NewInt(0)
				err := tc.set(ctx, v2)

				switch tc.name {
				case "minimumDeposit":
					require.ErrorContains(t, err, "0x5bc1c4a0") // invalid min
				}
			})
		})
	}
}

func TestPayerReportManager_ProtocolFeeRate_ReadDefault_WriteThenRead(t *testing.T) {
	scAdmin, paramAdmin := buildSettlementChainAdmin(t)
	ctx := context.Background()

	// default
	def, err := scAdmin.GetPayerReportManagerProtocolFeeRate(ctx)
	require.NoError(t, err)
	require.EqualValues(t, 0, def)

	rawDef, err := paramAdmin.GetParameterUint64(
		ctx,
		blockchain.PAYER_REPORT_MANAGER_PROTOCOL_FEE_RATE_KEY,
	)
	require.NoError(t, err)
	require.EqualValues(t, 0, rawDef)

	const v1 = 9999
	require.NoError(t, scAdmin.SetPayerReportManagerProtocolFeeRate(ctx, v1))

	got, err := scAdmin.GetPayerReportManagerProtocolFeeRate(ctx)
	require.NoError(t, err)
	require.EqualValues(t, v1, got)

	raw, err := paramAdmin.GetParameterUint64(
		ctx,
		blockchain.PAYER_REPORT_MANAGER_PROTOCOL_FEE_RATE_KEY,
	)
	require.NoError(t, err)
	require.EqualValues(t, v1, raw)

	// overwrite + zero
	const v2 = 0
	require.NoError(t, scAdmin.SetPayerReportManagerProtocolFeeRate(ctx, v2))

	gotZero, err := scAdmin.GetPayerReportManagerProtocolFeeRate(ctx)
	require.NoError(t, err)
	require.EqualValues(t, 0, gotZero)
}

func TestPayerReportManager_ProtocolFeeAbove100Perc(t *testing.T) {
	scAdmin, _ := buildSettlementChainAdmin(t)
	ctx := context.Background()

	const v1 = 10001
	err := scAdmin.SetPayerReportManagerProtocolFeeRate(ctx, v1)
	require.ErrorContains(t, err, "0x82eeb3b2")
}

func TestRateRegistry_Migrator_ReadDefault_WriteThenRead(t *testing.T) {
	scAdmin, paramAdmin := buildSettlementChainAdmin(t)
	ctx := context.Background()

	def, err := scAdmin.GetRateRegistryMigrator(ctx)
	require.NoError(t, err)
	require.Equal(t, common.Address{}, def)

	rawDef, err := paramAdmin.GetParameterAddress(ctx, blockchain.RATE_REGISTRY_MIGRATOR_KEY)
	require.NoError(t, err)
	require.Equal(t, common.Address{}, rawDef)

	want := common.HexToAddress("0x000000000000000000000000000000000000BEEF")
	require.NoError(t, scAdmin.SetRateRegistryMigrator(ctx, want))

	got, err := scAdmin.GetRateRegistryMigrator(ctx)
	require.NoError(t, err)
	require.Equal(t, want, got)

	raw, err := paramAdmin.GetParameterAddress(ctx, blockchain.RATE_REGISTRY_MIGRATOR_KEY)
	require.NoError(t, err)
	require.Equal(t, want, raw)

	other := common.HexToAddress("0x000000000000000000000000000000000000CAFE")
	require.NoError(t, scAdmin.SetRateRegistryMigrator(ctx, other))

	var zero common.Address
	require.NoError(t, scAdmin.SetRateRegistryMigrator(ctx, zero))
	gotZero, err := scAdmin.GetRateRegistryMigrator(ctx)
	require.NoError(t, err)
	require.Equal(t, zero, gotZero)
}
