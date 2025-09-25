package blockchain_test

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/xmtp/xmtpd/pkg/config"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/blockchain"
	"github.com/xmtp/xmtpd/pkg/testutils"
	"github.com/xmtp/xmtpd/pkg/testutils/anvil"

	"github.com/xmtp/xmtpd/pkg/abi/erc20"
)

// --- local helper that also returns client + options so we can assert balances on-chain.
func buildFundsAdminWithDeps(
	t *testing.T,
) (blockchain.IFundsAdmin, context.Context, *ethclient.Client, config.ContractsOptions) {
	t.Helper()

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	logger := testutils.NewLog(t)
	wsURL, rpcURL := anvil.StartAnvil(t, false)
	contractsOptions := testutils.NewContractsOptions(t, rpcURL, wsURL)

	signer, err := blockchain.NewPrivateKeySigner(
		testutils.GetPayerOptions(t).PrivateKey,
		contractsOptions.SettlementChain.ChainID,
	)
	require.NoError(t, err)

	client, err := blockchain.NewRPCClient(ctx, contractsOptions.SettlementChain.RPCURL)
	require.NoError(t, err)

	admin, err := blockchain.NewFundsAdmin(
		blockchain.FundsAdminOpts{
			Logger:          logger,
			ContractOptions: contractsOptions,
			Settlement: blockchain.FundsAdminSettlementOpts{
				Client: client,
				Signer: signer,
			},
			App: blockchain.FundsAdminAppOpts{
				Client: client,
				Signer: signer,
			},
		},
	)
	require.NoError(t, err)

	return admin, ctx, client, contractsOptions
}

// Balances() should not error and should log current balances.
func TestFundsAdmin_Balances_NoError(t *testing.T) {
	admin, ctx, _, _ := buildFundsAdminWithDeps(t)
	require.NoError(t, admin.Balances(ctx))
}

// MintMockUSDC: negative amount should be rejected immediately.
func TestFundsAdmin_MintMockUSDC_RejectsNegative(t *testing.T) {
	admin, ctx, _, _ := buildFundsAdminWithDeps(t)

	addr := testutils.RandomAddress()
	neg := big.NewInt(-1)

	err := admin.MintMockUSDC(ctx, addr, neg)
	require.Error(t, err)
	require.Contains(t, err.Error(), "amount must be positive")
}

// MintMockUSDC: > 10_000 * 1e6 (raw 10000000000) should be rejected.
func TestFundsAdmin_MintMockUSDC_RejectsTooLarge(t *testing.T) {
	admin, ctx, _, _ := buildFundsAdminWithDeps(t)

	addr := testutils.RandomAddress()
	tooLarge := big.NewInt(0).SetUint64(10000000001) // > 10_000e6

	err := admin.MintMockUSDC(ctx, addr, tooLarge)
	require.Error(t, err)
	require.Contains(t, err.Error(), "amount must be less than 10000 mxUSDC")
}

// MintMockUSDC: successful mint should reflect in the ERC20 balance at UnderlyingFeeToken address.
func TestFundsAdmin_MintMockUSDC_SetsBalance(t *testing.T) {
	admin, ctx, client, opts := buildFundsAdminWithDeps(t)

	underlyingAddr := common.HexToAddress(opts.SettlementChain.UnderlyingFeeToken)
	underlying, err := erc20.NewERC20(underlyingAddr, client)
	require.NoError(t, err)

	recipient := testutils.RandomAddress()
	amount := big.NewInt(1_234_567) // 1.234567 USDC, raw with 6 decimals

	require.NoError(t, admin.MintMockUSDC(ctx, recipient, amount))

	// The mock token implements ERC20; verify balance eventually (anvil + tx mining).
	require.Eventually(t, func() bool {
		bal, err := underlying.BalanceOf(&bind.CallOpts{Context: ctx}, recipient)
		if err != nil {
			return false
		}
		return bal.Cmp(amount) == 0
	}, 2*time.Second, 50*time.Millisecond)
}
