package contracts

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
	pr "github.com/xmtp/xmtpd/pkg/abi/payerregistry"
	"github.com/xmtp/xmtpd/pkg/blockchain"
	"github.com/xmtp/xmtpd/pkg/currency"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/ledger"
	"github.com/xmtp/xmtpd/pkg/testutils"
	"github.com/xmtp/xmtpd/pkg/testutils/anvil"
	re "github.com/xmtp/xmtpd/pkg/utils/retryerrors"
)

var address = common.HexToAddress("0x1")

type payerRegistryStorerTester struct {
	ctx    context.Context
	abi    *abi.ABI
	storer *PayerRegistryStorer
	ledger ledger.ILedger
}

type testCase struct {
	name          string
	initialLogs   []types.Log
	actionLogs    []types.Log
	validate      func(t *testing.T, tester *payerRegistryStorerTester)
	expectedError re.RetryableError
}

// getBalance is a helper function to get balance for an address with errors already checked
func (st *payerRegistryStorerTester) getBalance(
	t *testing.T,
	address common.Address,
) currency.PicoDollar {
	payerID, err := st.ledger.FindOrCreatePayer(st.ctx, address)
	require.NoError(t, err)

	balance, err := st.ledger.GetBalance(st.ctx, payerID)
	require.NoError(t, err)

	return balance
}

// setupWithInitialLogs sets up a tester with initial logs already processed
func setupWithInitialLogs(t *testing.T, logs []types.Log) *payerRegistryStorerTester {
	tester := buildPayerRegistryStorerTester(t)

	for _, log := range logs {
		err := tester.storer.StoreLog(tester.ctx, log)
		require.NoError(t, err, "failed to setup initial log")
	}

	return tester
}

// runTestCases executes parameterized test cases
func runTestCases(t *testing.T, testCases []testCase) {
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tester := setupWithInitialLogs(t, tc.initialLogs)

			// Execute action logs
			var lastErr error
			for _, log := range tc.actionLogs {
				lastErr = tester.storer.StoreLog(tester.ctx, log)
				// If we expect an error, we only care about the last one
				if tc.expectedError != nil {
					continue
				}
				require.NoError(t, lastErr)
			}

			// Check expected error
			if tc.expectedError != nil {
				require.Error(t, lastErr)
				require.ErrorAs(t, lastErr, &tc.expectedError)

				// Check retryable behavior matches expected
				if retryableErr, ok := lastErr.(re.RetryableError); ok {
					expectedRetryable := tc.expectedError.ShouldRetry()
					actualRetryable := retryableErr.ShouldRetry()
					require.Equal(t, expectedRetryable, actualRetryable,
						"expected ShouldRetry()=%v, got %v", expectedRetryable, actualRetryable)
				} else {
					t.Fatal("expected error to implement RetryableError interface")
				}
			}

			if tc.validate != nil {
				tc.validate(t, tester)
			}
		})
	}
}

func TestPayerRegistryStorer(t *testing.T) {
	tester := buildPayerRegistryStorerTester(t)

	testCases := []testCase{
		// Error cases
		{
			name:       "error_no_topics",
			actionLogs: []types.Log{{}}, // Empty log with no topics
			expectedError: re.NewNonRecoverableError(
				ErrParsePayerRegistryLog,
				errors.New("no topics"),
			),
		},
		{
			name: "error_unknown_event",
			actionLogs: []types.Log{{
				Topics: []common.Hash{common.HexToHash("UnknownEvent")},
			}},
			expectedError: re.NewNonRecoverableError(
				ErrParsePayerRegistryLog,
				fmt.Errorf("no event with id: %#x", common.HexToHash("UnknownEvent").Hex()),
			),
		},

		// Basic event processing
		{
			name: "usage_settled",
			actionLogs: []types.Log{
				tester.newUsageSettledLog(t, address, 10),
			},
			validate: func(t *testing.T, tester *payerRegistryStorerTester) {
				balance := tester.getBalance(t, address)
				require.Equal(t, currency.FromMicrodollars(10)*-1, balance)
			},
		},
		{
			name: "deposit_single",
			actionLogs: []types.Log{
				tester.newDepositLog(t, address, 20),
			},
			validate: func(t *testing.T, tester *payerRegistryStorerTester) {
				balance := tester.getBalance(t, address)
				require.Equal(t, currency.FromMicrodollars(20), balance)
			},
		},
		{
			name: "deposit_multiple",
			actionLogs: []types.Log{
				tester.newDepositLog(t, address, 10),
				tester.newDepositLog(t, address, 20),
			},
			validate: func(t *testing.T, tester *payerRegistryStorerTester) {
				balance := tester.getBalance(t, address)
				require.Equal(t, currency.FromMicrodollars(30), balance)
			},
		},
		{
			name: "withdrawal_requested",
			actionLogs: []types.Log{
				tester.newWithdrawalRequestedLog(t, address, 10, 100),
			},
			validate: func(t *testing.T, tester *payerRegistryStorerTester) {
				balance := tester.getBalance(t, address)
				require.Equal(t, currency.FromMicrodollars(10)*-1, balance)
			},
		},

		// Withdrawal cancellation scenarios
		{
			name: "withdrawal_cancelled_success",
			initialLogs: []types.Log{
				tester.newDepositLog(t, address, 50),
				tester.newWithdrawalRequestedLog(t, address, 20, 100),
			},
			actionLogs: []types.Log{
				tester.newWithdrawalCancelledLog(t, address),
			},
			validate: func(t *testing.T, tester *payerRegistryStorerTester) {
				balance := tester.getBalance(t, address)
				// 50 (deposit) - 20 (withdrawal) + 20 (cancellation) = 50
				require.Equal(t, currency.FromMicrodollars(50), balance)
			},
		},
		{
			name: "withdrawal_cancelled_without_prior_withdrawal",
			actionLogs: []types.Log{
				tester.newWithdrawalCancelledLog(t, address),
			},
			expectedError: re.NewRecoverableError("error", errors.New("no withdrawal to cancel")),
		},
		{
			name: "multiple_withdrawal_operations",
			initialLogs: []types.Log{
				tester.newDepositLog(t, address, 100),
				tester.newWithdrawalRequestedLog(t, address, 30, 100),
				tester.newWithdrawalRequestedLog(t, address, 20, 200),
			},
			actionLogs: []types.Log{
				tester.newWithdrawalCancelledLog(t, address),
			},
			validate: func(t *testing.T, tester *payerRegistryStorerTester) {
				balance := tester.getBalance(t, address)
				// 100 (deposit) - 30 (first withdrawal) - 20 (second withdrawal) + 20 (cancellation) = 70
				require.Equal(t, currency.FromMicrodollars(70), balance)
			},
		},
	}

	runTestCases(t, testCases)
}

func TestStoreLogIdempotency(t *testing.T) {
	tester := buildPayerRegistryStorerTester(t)

	// Create base logs with identical block/tx/index for true duplicates
	depositLog := tester.newDepositLog(t, address, 50)
	withdrawalLog := tester.newWithdrawalRequestedLog(t, address, 30, 100)
	usageLog := tester.newUsageSettledLog(t, address, 25)
	cancellationLog := tester.newWithdrawalCancelledLog(t, address)

	testCases := []testCase{
		{
			name: "duplicate_deposit_logs",
			actionLogs: []types.Log{
				depositLog,
				depositLog, // Exact same log
			},
			validate: func(t *testing.T, tester *payerRegistryStorerTester) {
				// Balance should only reflect one deposit
				balance := tester.getBalance(t, address)
				require.Equal(t, currency.FromMicrodollars(50), balance)
			},
		},
		{
			name: "duplicate_withdrawal_requested_logs",
			initialLogs: []types.Log{
				tester.newDepositLog(t, address, 100),
			},
			actionLogs: []types.Log{
				withdrawalLog,
				withdrawalLog, // Exact same log
			},
			validate: func(t *testing.T, tester *payerRegistryStorerTester) {
				// Balance should only reflect one withdrawal
				balance := tester.getBalance(t, address)
				require.Equal(t, currency.FromMicrodollars(70), balance) // 100 - 30
			},
		},
		{
			name: "duplicate_usage_settled_logs",
			initialLogs: []types.Log{
				tester.newDepositLog(t, address, 100),
			},
			actionLogs: []types.Log{
				usageLog,
				usageLog, // Exact same log
			},
			validate: func(t *testing.T, tester *payerRegistryStorerTester) {
				// Balance should only reflect one settlement
				balance := tester.getBalance(t, address)
				require.Equal(t, currency.FromMicrodollars(75), balance) // 100 - 25
			},
		},
		{
			name: "duplicate_withdrawal_cancelled_logs",
			initialLogs: []types.Log{
				tester.newDepositLog(t, address, 100),
				tester.newWithdrawalRequestedLog(t, address, 40, 100),
			},
			actionLogs: []types.Log{
				cancellationLog,
				cancellationLog, // Exact same log
			},
			validate: func(t *testing.T, tester *payerRegistryStorerTester) {
				// Balance should be fully restored (only one cancellation applied)
				balance := tester.getBalance(t, address)
				require.Equal(t, currency.FromMicrodollars(100), balance) // 100 - 40 + 40
			},
		},
	}

	runTestCases(t, testCases)
}

func buildPayerRegistryStorerTester(t *testing.T) *payerRegistryStorerTester {
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	// Dependencies.
	db, _ := testutils.NewDB(t, ctx)
	wsURL, rpcURL := anvil.StartAnvil(t, false)
	config := testutils.NewContractsOptions(t, rpcURL, wsURL)

	// Chain client.
	client, err := blockchain.NewRPCClient(
		ctx,
		config.AppChain.RPCURL,
	)
	require.NoError(t, err)

	// Contract.
	contract, err := pr.NewPayerRegistry(
		common.HexToAddress(config.SettlementChain.PayerRegistryAddress),
		client,
	)
	require.NoError(t, err)

	payerLedger := ledger.NewLedger(testutils.NewLog(t), queries.New(db))

	// Storer and ABI.
	storer, err := NewPayerRegistryStorer(testutils.NewLog(t), contract, payerLedger)
	require.NoError(t, err)

	abi, err := pr.PayerRegistryMetaData.GetAbi()
	require.NoError(t, err)

	return &payerRegistryStorerTester{
		ctx:    ctx,
		abi:    abi,
		storer: storer,
		ledger: payerLedger,
	}
}

func setLogFields(log *types.Log, blocknumber int, logIndex int, txHash common.Hash) {
	log.BlockNumber = uint64(blocknumber)
	log.Index = uint(logIndex)
	log.TxHash = txHash
}

func (st *payerRegistryStorerTester) newUsageSettledLog(
	t *testing.T,
	payerAddress common.Address,
	amount int64,
) types.Log {
	baseLog := testutils.BuildPayerRegistryUsageSettledLog(
		t,
		payerAddress,
		big.NewInt(int64(amount)),
	)
	setLogFields(&baseLog, 1, 0, testutils.RandomBlockHash())

	return baseLog
}

func (st *payerRegistryStorerTester) newDepositLog(
	t *testing.T,
	payerAddress common.Address,
	amount int64,
) types.Log {
	baseLog := testutils.BuildPayerRegistryDepositLog(t, payerAddress, big.NewInt(int64(amount)))
	setLogFields(&baseLog, 1, 0, testutils.RandomBlockHash())

	return baseLog
}

func (st *payerRegistryStorerTester) newWithdrawalRequestedLog(
	t *testing.T,
	payerAddress common.Address,
	amount int64,
	withdrawableTimestamp uint32,
) types.Log {
	baseLog := testutils.BuildPayerRegistryWithdrawalRequestedLog(
		t,
		payerAddress,
		big.NewInt(int64(amount)),
		withdrawableTimestamp,
	)
	setLogFields(&baseLog, 1, 0, testutils.RandomBlockHash())

	return baseLog
}

func (st *payerRegistryStorerTester) newWithdrawalCancelledLog(
	t *testing.T,
	payerAddress common.Address,
) types.Log {
	baseLog := testutils.BuildPayerRegistryWithdrawalCancelledLog(t, payerAddress)
	setLogFields(&baseLog, 1, 0, testutils.RandomBlockHash())

	return baseLog
}
