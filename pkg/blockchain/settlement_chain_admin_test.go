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
) (blockchain.ISettlementChainAdmin, blockchain.IParameterAdmin) {
	t.Helper()

	ctx := context.Background()
	logger := testutils.NewLog(t)
	wsURL, rpcURL := anvil.StartAnvil(t, false)
	contractsOptions := testutils.NewContractsOptions(t, rpcURL, wsURL)

	signer, err := blockchain.NewPrivateKeySigner(
		testutils.TestPrivateKey,
		contractsOptions.SettlementChain.ChainID,
	)
	require.NoError(t, err)

	client, err := blockchain.NewRPCClient(ctx, contractsOptions.SettlementChain.RPCURL)
	require.NoError(t, err)

	paramAdmin, err := blockchain.NewSettlementParameterAdmin(
		logger,
		client,
		signer,
		contractsOptions,
	)
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
		name   string
		key    string
		update func(ctx context.Context) error
		get    func(ctx context.Context) (bool, error)
	}

	cases := []pauseCase{
		{
			name:   "settlement-chain-gateway",
			key:    blockchain.SettlementChainGatewayPausedKey,
			update: settlementChainAdmin.UpdateSettlementChainGatewayPauseStatus,
			get:    settlementChainAdmin.GetSettlementChainGatewayPauseStatus,
		},
		{
			name:   "payer-registry",
			key:    blockchain.PayerRegistryPausedKey,
			update: settlementChainAdmin.UpdatePayerRegistryPauseStatus,
			get:    settlementChainAdmin.GetPayerRegistryPauseStatus,
		},
		{
			name:   "distribution-manager",
			key:    blockchain.DistributionManagerPausedKey,
			update: settlementChainAdmin.UpdateDistributionManagerPauseStatus,
			get:    settlementChainAdmin.GetDistributionManagerPauseStatus,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Run(tc.name+"/toggle_true_false", func(t *testing.T) {
				var err error
				require.NoError(t, paramAdmin.SetBoolParameter(ctx, tc.key, true))
				require.NoError(t, tc.update(ctx))
				b, err := paramAdmin.GetParameterBool(ctx, tc.key)
				require.NoError(t, err)
				require.True(t, b)

				got, err := tc.get(ctx)
				require.NoError(t, err)
				require.True(t, got)

				require.NoError(t, paramAdmin.SetBoolParameter(ctx, tc.key, false))
				b, err = paramAdmin.GetParameterBool(ctx, tc.key)
				require.NoError(t, err)
				require.False(t, b)

				require.NoError(t, tc.update(ctx))
				got, err = tc.get(ctx)
				require.NoError(t, err)
				require.False(t, got)
			})

			t.Run(tc.name+"/idempotent_repeat_true", func(t *testing.T) {
				require.NoError(t, paramAdmin.SetBoolParameter(ctx, tc.key, true))
				require.NoError(t, tc.update(ctx))
				require.NoError(t, tc.update(ctx))

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
		blockchain.DistributionManagerProtocolFeesRecipientKey,
	)
	require.NoError(t, err)
	require.Equal(t, common.Address{}, rawDefault)

	// update a value
	want := common.HexToAddress("0x000000000000000000000000000000000000BEEF")
	require.NoError(
		t,
		paramAdmin.SetAddressParameter(
			ctx,
			blockchain.DistributionManagerProtocolFeesRecipientKey,
			want,
		),
	)
	require.NoError(t, scAdmin.UpdateDistributionManagerProtocolFeesRecipient(ctx))

	// read back
	got, err := scAdmin.GetDistributionManagerProtocolFeesRecipient(ctx)
	require.NoError(t, err)
	require.Equal(t, want, got)

	raw, err := paramAdmin.GetParameterAddress(
		ctx,
		blockchain.DistributionManagerProtocolFeesRecipientKey,
	)
	require.NoError(t, err)
	require.Equal(t, want, raw)

	// overwrite + idempotent
	other := common.HexToAddress("0x000000000000000000000000000000000000CAFE")

	require.NoError(
		t,
		paramAdmin.SetAddressParameter(
			ctx,
			blockchain.DistributionManagerProtocolFeesRecipientKey,
			other,
		),
	)
	require.NoError(t, scAdmin.UpdateDistributionManagerProtocolFeesRecipient(ctx))

	// TODO: Reentrancy is broken
	// require.NoError(t, scAdmin.UpdateDistributionManagerProtocolFeesRecipient(ctx))

	got2, err := scAdmin.GetDistributionManagerProtocolFeesRecipient(ctx)
	require.NoError(t, err)
	require.Equal(t, other, got2)

	// zero round-trip
	var zero common.Address

	require.NoError(
		t,
		paramAdmin.SetAddressParameter(
			ctx,
			blockchain.DistributionManagerProtocolFeesRecipientKey,
			zero,
		),
	)
	require.NoError(t, scAdmin.UpdateDistributionManagerProtocolFeesRecipient(ctx))
	gotZero, err := scAdmin.GetDistributionManagerProtocolFeesRecipient(ctx)
	require.NoError(t, err)
	require.Equal(t, zero, gotZero)
}

func TestNodeRegistry_Admin_ReadDefault_WriteThenRead(t *testing.T) {
	scAdmin, paramAdmin := buildSettlementChainAdmin(t)
	ctx := context.Background()

	// update + read
	want := common.HexToAddress("0x000000000000000000000000000000000000DEAD")
	require.NoError(
		t,
		paramAdmin.SetAddressParameter(ctx, blockchain.NodeRegistryAdminKey, want),
	)
	require.NoError(t, scAdmin.UpdateNodeRegistryAdmin(ctx))

	got, err := scAdmin.GetNodeRegistryAdmin(ctx)
	require.NoError(t, err)
	require.Equal(t, want, got)

	raw, err := paramAdmin.GetParameterAddress(ctx, blockchain.NodeRegistryAdminKey)
	require.NoError(t, err)
	require.Equal(t, want, raw)

	// overwrite + zero
	other := common.HexToAddress("0x000000000000000000000000000000000000BEEF")
	require.NoError(
		t,
		paramAdmin.SetAddressParameter(ctx, blockchain.NodeRegistryAdminKey, other),
	)
	require.NoError(t, scAdmin.UpdateNodeRegistryAdmin(ctx))

	var zero common.Address
	require.NoError(
		t,
		paramAdmin.SetAddressParameter(ctx, blockchain.NodeRegistryAdminKey, zero),
	)
	require.NoError(t, scAdmin.UpdateNodeRegistryAdmin(ctx))
	gotZero, err := scAdmin.GetNodeRegistryAdmin(ctx)
	require.NoError(t, err)
	require.Equal(t, zero, gotZero)
}

func TestPayerRegistry_Uint32Params_ReadDefault_WriteThenRead(t *testing.T) {
	scAdmin, paramAdmin := buildSettlementChainAdmin(t)
	ctx := context.Background()

	type u32Case struct {
		name   string
		key    string
		update func(context.Context) error
		get    func(context.Context) (uint32, error)
	}
	cases := []u32Case{
		{
			name:   "withdrawLockPeriod",
			key:    blockchain.PayerRegistryWithdrawLockPeriodKey,
			update: scAdmin.UpdatePayerRegistryWithdrawLockPeriod,
			get:    scAdmin.GetPayerRegistryWithdrawLockPeriod,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Run(tc.name+"/read_default", func(t *testing.T) {
				t.Skip(
					"Some defaults are pre-update https://github.com/xmtp/smart-contracts/issues/126",
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

				require.NoError(t, paramAdmin.SetUint64Parameter(ctx, tc.key, v1))
				require.NoError(t, tc.update(ctx))

				gotV1, err := tc.get(ctx)
				require.NoError(t, err)
				require.EqualValues(t, v1, gotV1)

				rawV1, err := paramAdmin.GetParameterUint64(ctx, tc.key)
				require.NoError(t, err)
				require.EqualValues(t, v1, rawV1)
			})

			t.Run(tc.name+"/write_idempotent", func(t *testing.T) {
				const v1 = 1024
				require.NoError(t, paramAdmin.SetUint64Parameter(ctx, tc.key, v1))
				require.NoError(t, tc.update(ctx))
				require.NoError(t, tc.update(ctx))
			})

			t.Run(tc.name+"/write_back_to_zero", func(t *testing.T) {
				const v1 = 1024
				require.NoError(t, paramAdmin.SetUint64Parameter(ctx, tc.key, v1))
				require.NoError(t, tc.update(ctx))

				const v2 = 0
				require.NoError(t, paramAdmin.SetUint64Parameter(ctx, tc.key, v2))
				err := tc.update(ctx)

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
		name   string
		key    string
		update func(context.Context) error
		get    func(context.Context) (*big.Int, error)
	}
	cases := []u96Case{
		{
			name:   "minimumDeposit",
			key:    blockchain.PayerRegistryMinimumDepositKey,
			update: scAdmin.UpdatePayerRegistryMinimumDeposit,
			get:    scAdmin.GetPayerRegistryMinimumDeposit,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Run(tc.name+"/read_default", func(t *testing.T) {
				t.Skip(
					"Some defaults are pre-update https://github.com/xmtp/smart-contracts/issues/126",
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
				require.NoError(t, paramAdmin.SetUint96Parameter(ctx, tc.key, v1))
				require.NoError(t, tc.update(ctx))

				gotV1, err := tc.get(ctx)
				require.NoError(t, err)
				require.Equal(t, v1, gotV1)

				rawV1, err := paramAdmin.GetParameterUint96(ctx, tc.key)
				require.NoError(t, err)
				require.Equal(t, v1, rawV1)
			})

			t.Run(tc.name+"/write_idempotent", func(t *testing.T) {
				v1 := big.NewInt(1024)
				require.NoError(t, paramAdmin.SetUint96Parameter(ctx, tc.key, v1))
				require.NoError(t, tc.update(ctx))
				require.NoError(t, tc.update(ctx))
			})

			t.Run(tc.name+"/write_back_to_zero", func(t *testing.T) {
				v1 := big.NewInt(1024)
				require.NoError(t, paramAdmin.SetUint96Parameter(ctx, tc.key, v1))
				require.NoError(t, tc.update(ctx))

				v2 := big.NewInt(0)
				require.NoError(t, paramAdmin.SetUint96Parameter(ctx, tc.key, v2))
				err := tc.update(ctx)

				switch tc.name {
				case "minimumDeposit":
					require.ErrorContains(t, err, "ZeroMinimumDeposit()")
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
		blockchain.PayerReportManagerProtocolFeeRateKey,
	)
	require.NoError(t, err)
	require.EqualValues(t, 0, rawDef)

	const v1 = 9999

	require.NoError(
		t,
		paramAdmin.SetUint16Parameter(
			ctx,
			blockchain.PayerReportManagerProtocolFeeRateKey,
			v1,
		),
	)
	require.NoError(t, scAdmin.UpdatePayerReportManagerProtocolFeeRate(ctx))

	got, err := scAdmin.GetPayerReportManagerProtocolFeeRate(ctx)
	require.NoError(t, err)
	require.EqualValues(t, v1, got)

	raw, err := paramAdmin.GetParameterUint64(
		ctx,
		blockchain.PayerReportManagerProtocolFeeRateKey,
	)
	require.NoError(t, err)
	require.EqualValues(t, v1, raw)

	// overwrite + zero
	const v2 = 0
	require.NoError(
		t,
		paramAdmin.SetUint16Parameter(
			ctx,
			blockchain.PayerReportManagerProtocolFeeRateKey,
			v2,
		),
	)
	require.NoError(t, scAdmin.UpdatePayerReportManagerProtocolFeeRate(ctx))

	gotZero, err := scAdmin.GetPayerReportManagerProtocolFeeRate(ctx)
	require.NoError(t, err)
	require.EqualValues(t, 0, gotZero)
}

func TestPayerReportManager_ProtocolFeeAbove100Perc(t *testing.T) {
	scAdmin, paramAdmin := buildSettlementChainAdmin(t)
	ctx := context.Background()

	const v1 = 10001
	require.NoError(
		t,
		paramAdmin.SetUint16Parameter(
			ctx,
			blockchain.PayerReportManagerProtocolFeeRateKey,
			v1,
		),
	)
	err := scAdmin.UpdatePayerReportManagerProtocolFeeRate(ctx)
	require.ErrorContains(t, err, "InvalidProtocolFeeRate()")
}
