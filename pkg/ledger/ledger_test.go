package ledger_test

import (
	"context"
	"crypto/sha256"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/currency"
	"github.com/xmtp/xmtpd/pkg/db"
	"github.com/xmtp/xmtpd/pkg/ledger"
	"github.com/xmtp/xmtpd/pkg/testutils"
	"go.uber.org/zap"
)

type testFixture struct {
	ctx    context.Context
	db     *db.Handler
	ledger *ledger.Ledger
}

type initialEvent struct {
	payerID   int32
	amount    currency.PicoDollar
	eventID   ledger.EventID
	eventType string // "deposit", "withdrawal", "settlement"
}

type testCase struct {
	name          string
	initialEvents []initialEvent
	action        func(f *testFixture) error
	validate      func(t *testing.T, f *testFixture)
	expectedError error
}

func setupTest(t *testing.T) *testFixture {
	ctx := context.Background()
	db, _ := testutils.NewDB(t, ctx)
	logger := zap.NewNop()
	l := ledger.NewLedger(logger, db)

	return &testFixture{
		ctx:    ctx,
		db:     db,
		ledger: l,
	}
}

func generateEventID(id int) ledger.EventID {
	// Generate a deterministic EventID for testing
	data := fmt.Sprintf("test-event-%d", id)
	return sha256.Sum256([]byte(data))
}

func setupWithInitialState(t *testing.T, events []initialEvent) *testFixture {
	f := setupTest(t)

	for _, event := range events {
		var err error
		switch event.eventType {
		case "deposit":
			err = f.ledger.Deposit(f.ctx, event.payerID, event.amount, event.eventID)
		case "withdrawal":
			err = f.ledger.InitiateWithdrawal(f.ctx, event.payerID, event.amount, event.eventID)
		case "settlement":
			err = f.ledger.SettleUsage(f.ctx, event.payerID, event.amount, event.eventID)
		default:
			t.Fatalf("unknown event type: %s", event.eventType)
		}
		require.NoError(t, err, "failed to setup initial event")
	}

	return f
}

func TestDuplicateEventPrevention(t *testing.T) {
	testCases := []testCase{
		{
			name: "duplicate_deposit_events_ignored",
			initialEvents: []initialEvent{
				{payerID: 1, amount: 1000, eventID: generateEventID(100), eventType: "deposit"},
			},
			action: func(f *testFixture) error {
				// Try to deposit again with same eventID
				return f.ledger.Deposit(f.ctx, 1, 1000, generateEventID(100))
			},
			validate: func(t *testing.T, f *testFixture) {
				// Balance should remain unchanged at 1000
				balance, err := f.ledger.GetBalance(f.ctx, 1)
				require.NoError(t, err)
				require.Equal(t, currency.PicoDollar(1000), balance)
			},
		},
		{
			name: "duplicate_withdrawal_events_ignored",
			initialEvents: []initialEvent{
				{payerID: 2, amount: 5000, eventID: generateEventID(200), eventType: "deposit"},
				{payerID: 2, amount: 1000, eventID: generateEventID(201), eventType: "withdrawal"},
			},
			action: func(f *testFixture) error {
				// Try to withdraw again with same eventID
				return f.ledger.InitiateWithdrawal(f.ctx, 2, 1000, generateEventID(201))
			},
			validate: func(t *testing.T, f *testFixture) {
				// Balance should remain at 4000 (5000 - 1000)
				balance, err := f.ledger.GetBalance(f.ctx, 2)
				require.NoError(t, err)
				require.Equal(t, currency.PicoDollar(4000), balance)
			},
		},
		{
			name: "duplicate_settlement_events_ignored",
			initialEvents: []initialEvent{
				{payerID: 3, amount: 5000, eventID: generateEventID(300), eventType: "deposit"},
				{payerID: 3, amount: 500, eventID: generateEventID(301), eventType: "settlement"},
			},
			action: func(f *testFixture) error {
				// Try to settle again with same eventID
				return f.ledger.SettleUsage(f.ctx, 3, 500, generateEventID(301))
			},
			validate: func(t *testing.T, f *testFixture) {
				// Balance should remain at 4500 (5000 - 500)
				balance, err := f.ledger.GetBalance(f.ctx, 3)
				require.NoError(t, err)
				require.Equal(t, currency.PicoDollar(4500), balance)
			},
		},
		{
			name: "same_event_id_blocks_different_payer",
			initialEvents: []initialEvent{
				{payerID: 4, amount: 1000, eventID: generateEventID(400), eventType: "deposit"},
			},
			action: func(f *testFixture) error {
				// Try to use same eventID for different payer
				return f.ledger.Deposit(f.ctx, 5, 2000, generateEventID(400))
			},
			validate: func(t *testing.T, f *testFixture) {
				// First payer should have their balance
				balance1, err := f.ledger.GetBalance(f.ctx, 4)
				require.NoError(t, err)
				require.Equal(t, currency.PicoDollar(1000), balance1)

				// Second payer should have zero (event was ignored)
				balance2, err := f.ledger.GetBalance(f.ctx, 5)
				require.NoError(t, err)
				require.Equal(t, currency.PicoDollar(0), balance2)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			f := setupWithInitialState(t, tc.initialEvents)

			err := tc.action(f)
			if tc.expectedError != nil {
				require.Error(t, err)
				require.Equal(t, tc.expectedError, err)
			} else {
				require.NoError(t, err)
			}

			tc.validate(t, f)
		})
	}
}

func TestBalanceCalculations(t *testing.T) {
	testCases := []testCase{
		{
			name:          "empty_balance_returns_zero",
			initialEvents: []initialEvent{},
			action: func(f *testFixture) error {
				return nil // No action, just checking balance
			},
			validate: func(t *testing.T, f *testFixture) {
				balance, err := f.ledger.GetBalance(f.ctx, 10)
				require.NoError(t, err)
				require.Equal(t, currency.PicoDollar(0), balance)
			},
		},
		{
			name: "single_deposit_balance",
			initialEvents: []initialEvent{
				{payerID: 11, amount: 5000, eventID: generateEventID(1100), eventType: "deposit"},
			},
			action: func(f *testFixture) error {
				return nil // No action, just checking balance
			},
			validate: func(t *testing.T, f *testFixture) {
				balance, err := f.ledger.GetBalance(f.ctx, 11)
				require.NoError(t, err)
				require.Equal(t, currency.PicoDollar(5000), balance)
			},
		},
		{
			name: "multiple_deposits_sum_correctly",
			initialEvents: []initialEvent{
				{payerID: 12, amount: 1000, eventID: generateEventID(1201), eventType: "deposit"},
				{payerID: 12, amount: 2000, eventID: generateEventID(1202), eventType: "deposit"},
				{payerID: 12, amount: 3000, eventID: generateEventID(1203), eventType: "deposit"},
			},
			action: func(f *testFixture) error {
				return nil // No action, just checking balance
			},
			validate: func(t *testing.T, f *testFixture) {
				balance, err := f.ledger.GetBalance(f.ctx, 12)
				require.NoError(t, err)
				require.Equal(t, currency.PicoDollar(6000), balance)
			},
		},
		{
			name: "mixed_transactions_balance",
			initialEvents: []initialEvent{
				{payerID: 13, amount: 10000, eventID: generateEventID(1301), eventType: "deposit"},
				{
					payerID:   13,
					amount:    2000,
					eventID:   generateEventID(1302),
					eventType: "withdrawal",
				},
				{
					payerID:   13,
					amount:    1000,
					eventID:   generateEventID(1303),
					eventType: "settlement",
				},
				{payerID: 13, amount: 5000, eventID: generateEventID(1304), eventType: "deposit"},
			},
			action: func(f *testFixture) error {
				return nil // No action, just checking balance
			},
			validate: func(t *testing.T, f *testFixture) {
				// Expected: 10000 - 2000 - 1000 + 5000 = 12000
				balance, err := f.ledger.GetBalance(f.ctx, 13)
				require.NoError(t, err)
				require.Equal(t, currency.PicoDollar(12000), balance)
			},
		},
		{
			name: "withdrawal_reduces_balance",
			initialEvents: []initialEvent{
				{payerID: 14, amount: 5000, eventID: generateEventID(1401), eventType: "deposit"},
			},
			action: func(f *testFixture) error {
				return f.ledger.InitiateWithdrawal(f.ctx, 14, 1500, generateEventID(1402))
			},
			validate: func(t *testing.T, f *testFixture) {
				balance, err := f.ledger.GetBalance(f.ctx, 14)
				require.NoError(t, err)
				require.Equal(t, currency.PicoDollar(3500), balance)
			},
		},
		{
			name: "settlement_reduces_balance",
			initialEvents: []initialEvent{
				{payerID: 15, amount: 8000, eventID: generateEventID(1501), eventType: "deposit"},
			},
			action: func(f *testFixture) error {
				return f.ledger.SettleUsage(f.ctx, 15, 2500, generateEventID(1502))
			},
			validate: func(t *testing.T, f *testFixture) {
				balance, err := f.ledger.GetBalance(f.ctx, 15)
				require.NoError(t, err)
				require.Equal(t, currency.PicoDollar(5500), balance)
			},
		},
		{
			name:          "negative_balance_allowed",
			initialEvents: []initialEvent{},
			action: func(f *testFixture) error {
				// Withdraw without deposit
				return f.ledger.InitiateWithdrawal(f.ctx, 16, 1000, generateEventID(1601))
			},
			validate: func(t *testing.T, f *testFixture) {
				balance, err := f.ledger.GetBalance(f.ctx, 16)
				require.NoError(t, err)
				require.Equal(t, currency.PicoDollar(-1000), balance)
			},
		},
		{
			name: "isolated_payer_balances",
			initialEvents: []initialEvent{
				{payerID: 17, amount: 3000, eventID: generateEventID(1701), eventType: "deposit"},
				{payerID: 18, amount: 7000, eventID: generateEventID(1702), eventType: "deposit"},
			},
			action: func(f *testFixture) error {
				return nil // No action, just checking balances
			},
			validate: func(t *testing.T, f *testFixture) {
				balance1, err := f.ledger.GetBalance(f.ctx, 17)
				require.NoError(t, err)
				require.Equal(t, currency.PicoDollar(3000), balance1)

				balance2, err := f.ledger.GetBalance(f.ctx, 18)
				require.NoError(t, err)
				require.Equal(t, currency.PicoDollar(7000), balance2)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			f := setupWithInitialState(t, tc.initialEvents)

			err := tc.action(f)
			if tc.expectedError != nil {
				require.Error(t, err)
				require.Equal(t, tc.expectedError, err)
			} else {
				require.NoError(t, err)
			}

			tc.validate(t, f)
		})
	}
}

func TestAmountValidations(t *testing.T) {
	testCases := []struct {
		name          string
		action        func(f *testFixture) error
		expectedError error
	}{
		{
			name: "deposit_zero_amount_rejected",
			action: func(f *testFixture) error {
				return f.ledger.Deposit(f.ctx, 20, 0, generateEventID(2001))
			},
			expectedError: ledger.ErrInvalidAmount,
		},
		{
			name: "deposit_negative_amount_rejected",
			action: func(f *testFixture) error {
				return f.ledger.Deposit(f.ctx, 21, -1000, generateEventID(2101))
			},
			expectedError: ledger.ErrInvalidAmount,
		},
		{
			name: "withdrawal_zero_amount_rejected",
			action: func(f *testFixture) error {
				return f.ledger.InitiateWithdrawal(f.ctx, 22, 0, generateEventID(2201))
			},
			expectedError: ledger.ErrInvalidAmount,
		},
		{
			name: "withdrawal_negative_amount_rejected",
			action: func(f *testFixture) error {
				return f.ledger.InitiateWithdrawal(f.ctx, 23, -500, generateEventID(2301))
			},
			expectedError: ledger.ErrInvalidAmount,
		},
		{
			name: "settlement_zero_amount_rejected",
			action: func(f *testFixture) error {
				return f.ledger.SettleUsage(f.ctx, 24, 0, generateEventID(2401))
			},
			expectedError: ledger.ErrInvalidAmount,
		},
		{
			name: "settlement_negative_amount_rejected",
			action: func(f *testFixture) error {
				return f.ledger.SettleUsage(f.ctx, 25, -100, generateEventID(2501))
			},
			expectedError: ledger.ErrInvalidAmount,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			f := setupTest(t)
			err := tc.action(f)
			require.Error(t, err)
			require.Equal(t, tc.expectedError, err)
		})
	}

	t.Run("valid_positive_amounts_accepted", func(t *testing.T) {
		f := setupTest(t)
		payerID := int32(26)
		amounts := []currency.PicoDollar{1, 100, 1000, 1000000}

		for i, amount := range amounts {
			baseEventID := 2600 + i*3

			// Test deposit
			err := f.ledger.Deposit(f.ctx, payerID, amount, generateEventID(baseEventID))
			require.NoError(t, err)

			// Test withdrawal
			err = f.ledger.InitiateWithdrawal(
				f.ctx,
				payerID,
				amount,
				generateEventID(baseEventID+1),
			)
			require.NoError(t, err)

			// Test settlement
			err = f.ledger.SettleUsage(f.ctx, payerID, amount, generateEventID(baseEventID+2))
			require.NoError(t, err)
		}
	})
}

func TestEventIDValidations(t *testing.T) {
	testCases := []struct {
		name          string
		action        func(f *testFixture) error
		expectedError error
	}{
		{
			name: "deposit_zero_event_id_rejected",
			action: func(f *testFixture) error {
				return f.ledger.Deposit(f.ctx, 30, 1000, ledger.EventID{})
			},
			expectedError: ledger.ErrInvalidEventID,
		},
		{
			name: "deposit_negative_event_id_rejected",
			action: func(f *testFixture) error {
				return f.ledger.Deposit(f.ctx, 31, 1000, ledger.EventID{})
			},
			expectedError: ledger.ErrInvalidEventID,
		},
		{
			name: "withdrawal_zero_event_id_rejected",
			action: func(f *testFixture) error {
				return f.ledger.InitiateWithdrawal(f.ctx, 32, 1000, ledger.EventID{})
			},
			expectedError: ledger.ErrInvalidEventID,
		},
		{
			name: "withdrawal_negative_event_id_rejected",
			action: func(f *testFixture) error {
				return f.ledger.InitiateWithdrawal(f.ctx, 33, 1000, ledger.EventID{})
			},
			expectedError: ledger.ErrInvalidEventID,
		},
		{
			name: "settlement_zero_event_id_rejected",
			action: func(f *testFixture) error {
				return f.ledger.SettleUsage(f.ctx, 34, 1000, ledger.EventID{})
			},
			expectedError: ledger.ErrInvalidEventID,
		},
		{
			name: "settlement_negative_event_id_rejected",
			action: func(f *testFixture) error {
				return f.ledger.SettleUsage(f.ctx, 35, 1000, ledger.EventID{})
			},
			expectedError: ledger.ErrInvalidEventID,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			f := setupTest(t)
			err := tc.action(f)
			require.Error(t, err)
			require.Equal(t, tc.expectedError, err)
		})
	}

	t.Run("valid_positive_event_ids_accepted", func(t *testing.T) {
		f := setupTest(t)
		payerID := int32(36)
		amount := currency.PicoDollar(1000)
		eventIDs := []int{1, 100, 1000, 999999} // Test with various valid IDs

		for _, eventID := range eventIDs {
			err := f.ledger.Deposit(f.ctx, payerID, amount, generateEventID(eventID))
			require.NoError(t, err)
		}
	})
}

func TestEventIDGeneration(t *testing.T) {
	t.Run("generate_unique_event_ids", func(t *testing.T) {
		// Test that different inputs generate different EventIDs
		id1 := generateEventID(1)
		id2 := generateEventID(2)
		require.NotEqual(t, id1, id2)

		// Test that same input generates same EventID
		id3 := generateEventID(1)
		require.Equal(t, id1, id3)

		// Test that zero EventID is different from generated ones
		zeroID := ledger.EventID{}
		require.NotEqual(t, zeroID, id1)
	})
}

func TestConcurrentOperations(t *testing.T) {
	t.Run("concurrent_deposits_different_events", func(t *testing.T) {
		f := setupTest(t)
		payerID := int32(40)
		numDeposits := 10
		amount := currency.PicoDollar(100)

		// Run deposits concurrently
		errChan := make(chan error, numDeposits)
		for i := range numDeposits {
			go func(eventNum int) {
				err := f.ledger.Deposit(f.ctx, payerID, amount, generateEventID(4000+eventNum))
				errChan <- err
			}(i)
		}

		// Collect results
		for range numDeposits {
			err := <-errChan
			require.NoError(t, err)
		}

		// Verify final balance
		balance, err := f.ledger.GetBalance(f.ctx, payerID)
		require.NoError(t, err)
		require.Equal(t, currency.PicoDollar(numDeposits*100), balance)
	})

	t.Run("concurrent_mixed_operations", func(t *testing.T) {
		f := setupWithInitialState(t, []initialEvent{
			{payerID: 41, amount: 10000, eventID: generateEventID(4100), eventType: "deposit"},
		})

		// Run mixed operations concurrently
		type operation struct {
			fn func() error
		}

		operations := []operation{
			{fn: func() error { return f.ledger.Deposit(f.ctx, 41, 1000, generateEventID(4101)) }},
			{fn: func() error { return f.ledger.Deposit(f.ctx, 41, 2000, generateEventID(4102)) }},
			{
				fn: func() error { return f.ledger.InitiateWithdrawal(f.ctx, 41, 500, generateEventID(4103)) },
			},
			{
				fn: func() error { return f.ledger.InitiateWithdrawal(f.ctx, 41, 700, generateEventID(4104)) },
			},
			{
				fn: func() error { return f.ledger.SettleUsage(f.ctx, 41, 300, generateEventID(4105)) },
			},
			{
				fn: func() error { return f.ledger.SettleUsage(f.ctx, 41, 400, generateEventID(4106)) },
			},
		}

		errChan := make(chan error, len(operations))
		for _, op := range operations {
			go func(fn func() error) {
				errChan <- fn()
			}(op.fn)
		}

		// Collect results
		for range len(operations) {
			err := <-errChan
			require.NoError(t, err)
		}

		// Verify final balance: 10000 + 1000 + 2000 - 500 - 700 - 300 - 400 = 11100
		balance, err := f.ledger.GetBalance(f.ctx, 41)
		require.NoError(t, err)
		require.Equal(t, currency.PicoDollar(11100), balance)
	})
}

func TestEdgeCases(t *testing.T) {
	t.Run("large_amounts", func(t *testing.T) {
		f := setupTest(t)
		payerID := int32(50)
		// Test with very large amounts (close to max int64)
		largeAmount := currency.PicoDollar(9223372036854775000)

		err := f.ledger.Deposit(f.ctx, payerID, largeAmount, generateEventID(5001))
		require.NoError(t, err)

		balance, err := f.ledger.GetBalance(f.ctx, payerID)
		require.NoError(t, err)
		require.Equal(t, largeAmount, balance)
	})

	t.Run("many_transactions_single_payer", func(t *testing.T) {
		f := setupTest(t)
		payerID := int32(51)
		numTransactions := 100

		// Perform many small transactions
		for i := range numTransactions {
			eventID := generateEventID(5100 + i)
			switch i % 3 {
			case 0:
				err := f.ledger.Deposit(f.ctx, payerID, 10, eventID)
				require.NoError(t, err)
			case 1:
				err := f.ledger.InitiateWithdrawal(f.ctx, payerID, 3, eventID)
				require.NoError(t, err)
			case 2:
				err := f.ledger.SettleUsage(f.ctx, payerID, 2, eventID)
				require.NoError(t, err)
			}
		}

		// Verify balance is calculated correctly
		balance, err := f.ledger.GetBalance(f.ctx, payerID)
		require.NoError(t, err)
		// 34 deposits of 10 = 340
		// 33 withdrawals of 3 = -99
		// 33 settlements of 2 = -66
		// Total: 340 - 99 - 66 = 175
		require.Equal(t, currency.PicoDollar(175), balance)
	})

	t.Run("cancel_withdrawal_without_prior_withdrawal_fails", func(t *testing.T) {
		f := setupTest(t)
		payerID := int32(52)

		err := f.ledger.CancelWithdrawal(f.ctx, payerID, generateEventID(5201))
		require.Error(t, err)
	})
}

func TestCancelWithdrawal(t *testing.T) {
	testCases := []testCase{
		{
			name: "cancel_withdrawal_restores_balance",
			initialEvents: []initialEvent{
				{payerID: 60, amount: 5000, eventID: generateEventID(6001), eventType: "deposit"},
				{
					payerID:   60,
					amount:    1000,
					eventID:   generateEventID(6002),
					eventType: "withdrawal",
				},
			},
			action: func(f *testFixture) error {
				return f.ledger.CancelWithdrawal(f.ctx, 60, generateEventID(6003))
			},
			validate: func(t *testing.T, f *testFixture) {
				// Balance should be restored to original 5000
				balance, err := f.ledger.GetBalance(f.ctx, 60)
				require.NoError(t, err)
				require.Equal(t, currency.PicoDollar(5000), balance)
			},
		},
		{
			name:          "cancel_withdrawal_without_prior_withdrawal",
			initialEvents: []initialEvent{},
			action: func(f *testFixture) error {
				return f.ledger.CancelWithdrawal(f.ctx, 61, generateEventID(6101))
			},
			validate:      func(t *testing.T, f *testFixture) {},
			expectedError: ledger.ErrWithdrawalNotFound,
		},
		{
			name: "duplicate_cancellation_same_event_id_succeeds",
			initialEvents: []initialEvent{
				{payerID: 62, amount: 3000, eventID: generateEventID(6201), eventType: "deposit"},
				{payerID: 62, amount: 500, eventID: generateEventID(6202), eventType: "withdrawal"},
			},
			action: func(f *testFixture) error {
				// First cancellation
				err := f.ledger.CancelWithdrawal(f.ctx, 62, generateEventID(6203))
				if err != nil {
					return err
				}
				// Second cancellation with same event ID should succeed (idempotent)
				return f.ledger.CancelWithdrawal(f.ctx, 62, generateEventID(6203))
			},
			validate: func(t *testing.T, f *testFixture) {
				balance, err := f.ledger.GetBalance(f.ctx, 62)
				require.NoError(t, err)
				require.Equal(t, currency.PicoDollar(3000), balance)
			},
		},
		{
			name: "duplicate_cancellation_different_event_id_fails",
			initialEvents: []initialEvent{
				{payerID: 63, amount: 4000, eventID: generateEventID(6301), eventType: "deposit"},
				{payerID: 63, amount: 800, eventID: generateEventID(6302), eventType: "withdrawal"},
			},
			action: func(f *testFixture) error {
				// First cancellation
				err := f.ledger.CancelWithdrawal(f.ctx, 63, generateEventID(6303))
				if err != nil {
					return err
				}
				// Second cancellation with different event ID should fail
				return f.ledger.CancelWithdrawal(f.ctx, 63, generateEventID(6304))
			},
			validate:      func(t *testing.T, f *testFixture) {},
			expectedError: ledger.ErrWithdrawalAlreadyCanceled,
		},
		{
			name: "cancel_most_recent_withdrawal_only",
			initialEvents: []initialEvent{
				{payerID: 64, amount: 6000, eventID: generateEventID(6401), eventType: "deposit"},
				{
					payerID:   64,
					amount:    1000,
					eventID:   generateEventID(6402),
					eventType: "withdrawal",
				},
				{payerID: 64, amount: 500, eventID: generateEventID(6403), eventType: "withdrawal"},
			},
			action: func(f *testFixture) error {
				// Should cancel the most recent withdrawal (500)
				return f.ledger.CancelWithdrawal(f.ctx, 64, generateEventID(6404))
			},
			validate: func(t *testing.T, f *testFixture) {
				// Balance should be 6000 - 1000 - 500 + 500 = 5000
				balance, err := f.ledger.GetBalance(f.ctx, 64)
				require.NoError(t, err)
				require.Equal(t, currency.PicoDollar(5000), balance)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			f := setupWithInitialState(t, tc.initialEvents)

			err := tc.action(f)
			if tc.expectedError != nil {
				require.Error(t, err)
				require.Equal(t, tc.expectedError, err)
			} else {
				require.NoError(t, err)
			}

			tc.validate(t, f)
		})
	}
}
