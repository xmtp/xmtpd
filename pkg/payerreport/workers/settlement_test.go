package workers

import (
	"context"
	"errors"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/blockchain"
	"github.com/xmtp/xmtpd/pkg/currency"
	"github.com/xmtp/xmtpd/pkg/merkle"
	blockchainMocks "github.com/xmtp/xmtpd/pkg/mocks/blockchain"
	payerreportMocks "github.com/xmtp/xmtpd/pkg/mocks/payerreport"
	"github.com/xmtp/xmtpd/pkg/payerreport"
	"github.com/xmtp/xmtpd/pkg/testutils"
)

func testSettlementWorker(
	t *testing.T,
) (*SettlementWorker, *payerreport.Store, *blockchainMocks.MockPayerReportsManager, *payerreportMocks.MockIPayerReportVerifier) {
	var (
		log            = testutils.NewLog(t)
		ctx            = t.Context()
		db, _          = testutils.NewDB(t, ctx)
		store          = payerreport.NewStore(log, db)
		reportsManager = blockchainMocks.NewMockPayerReportsManager(t)
		verifier       = payerreportMocks.NewMockIPayerReportVerifier(t)
		myNodeID       = uint32(1)
	)

	worker := NewSettlementWorker(
		ctx,
		log,
		store,
		verifier,
		reportsManager,
		myNodeID,
	)

	return worker, store, reportsManager, verifier
}

// TestSettlementStatesAndTransitions covers all possible valid and invalid states for the settlement worker.
//
// Valid transitions:
//
//	(1,X) → (2,X) - SettlementWorker succeeds
//
// Invalid transitions:
//
//	(0,X) - Not submitted yet
//	(2,X) - Already settled
//	(3,X) - Submission rejected
func TestSettlementStatesAndTransitions(t *testing.T) {
	type reportState int
	const (
		stateSubmissionPending reportState = iota
		stateSubmissionSubmittedNotSettled
		stateSubmissionSubmittedAlreadySettled
		stateSubmissionSettled
		stateSubmissionRejected
	)

	prepareReport := func(
		t *testing.T,
		state reportState,
		store payerreport.IPayerReportStore,
		r *payerreport.PayerReportWithStatus,
	) {
		t.Helper()

		switch state {
		case stateSubmissionPending:
			// Nothing to do, report is already in pending state

		case stateSubmissionSubmittedNotSettled:
			require.NoError(t, store.SetReportSubmitted(t.Context(), r.ID, 0))

		case stateSubmissionSubmittedAlreadySettled:
			require.NoError(t, store.SetReportSubmitted(t.Context(), r.ID, 0))

		case stateSubmissionSettled:
			require.NoError(t, store.SetReportSubmitted(t.Context(), r.ID, 0))
			require.NoError(t, store.SetReportSettled(t.Context(), r.ID))

		case stateSubmissionRejected:
			require.NoError(t, store.SetReportSubmissionRejected(t.Context(), r.ID))

		default:
			t.Fatalf("unknown target state: %v", state)
		}
	}

	testCases := []struct {
		name                     string
		state                    reportState
		expectedSubmissionStatus payerreport.SubmissionStatus
		alreadySettledOnChain    bool
	}{
		{
			"don't settle pending report",
			stateSubmissionPending,
			payerreport.SubmissionPending,
			false,
		},
		{
			"settle submitted report",
			stateSubmissionSubmittedNotSettled,
			payerreport.SubmissionSubmitted,
			false,
		},
		{
			"settle submitted report that is already settled on chain",
			stateSubmissionSubmittedAlreadySettled,
			payerreport.SubmissionSubmitted,
			true,
		},
		{
			"don't settle already settled report",
			stateSubmissionSettled,
			payerreport.SubmissionSettled,
			false,
		},
		{
			"don't settle rejected report",
			stateSubmissionRejected,
			payerreport.SubmissionRejected,
			false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			worker, store, reportsManager, verifier := testSettlementWorker(t)

			payers := map[common.Address]currency.PicoDollar{
				common.HexToAddress("0x1"): 100,
				common.HexToAddress("0x2"): 200,
			}

			report, err := payerreport.BuildPayerReport(payerreport.BuildPayerReportParams{
				OriginatorNodeID:    1,
				StartSequenceID:     0,
				EndSequenceID:       10,
				EndMinuteSinceEpoch: 1000,
				DomainSeparator:     domainSeparator,
				NodeIDs:             []uint32{1, 2},
				Payers:              payers,
			})
			require.NoError(t, err)

			stored := storeReport(t, store, &report.PayerReport)
			prepareReport(t, tc.state, store, stored)

			if tc.state == stateSubmissionSubmittedNotSettled ||
				tc.state == stateSubmissionSubmittedAlreadySettled {
				verifier.EXPECT().
					GetPayerMap(mock.Anything, mock.Anything).
					Return(payers, nil)

				if tc.alreadySettledOnChain {
					reportsManager.EXPECT().
						SettlementSummary(mock.Anything, uint32(1), uint64(0)).
						Return(&blockchain.SettlementSummary{
							Offset:    0,
							IsSettled: true,
						}, nil)
				} else {
					reportsManager.EXPECT().
						SettlementSummary(mock.Anything, uint32(1), uint64(0)).
						Return(&blockchain.SettlementSummary{
							Offset:    0,
							IsSettled: false,
						}, nil)

					reportsManager.EXPECT().
						SettleReport(mock.Anything, uint32(1), uint64(0), mock.Anything).
						Return(nil)
				}
			}

			err = worker.SettleReports(t.Context())
			require.NoError(t, err)

			got, err := store.FetchReport(t.Context(), stored.ID)
			require.NoError(t, err)

			expectedFinalStatus := tc.expectedSubmissionStatus
			if tc.state == stateSubmissionSubmittedNotSettled ||
				tc.state == stateSubmissionSubmittedAlreadySettled {
				expectedFinalStatus = payerreport.SubmissionSettled
			}

			require.Equal(t, expectedFinalStatus, got.SubmissionStatus)
		})
	}
}

func TestSettleReportWithMultipleBatches(t *testing.T) {
	worker, store, reportsManager, verifier := testSettlementWorker(t)

	// Create a report with many payers to test batching
	payers := make(map[common.Address]currency.PicoDollar)
	for i := range 250 {
		// Create unique addresses
		addr := common.BigToAddress(big.NewInt(int64(i + 1)))
		payers[addr] = currency.PicoDollar(i + 1)
	}

	report, err := payerreport.BuildPayerReport(payerreport.BuildPayerReportParams{
		OriginatorNodeID:    1,
		StartSequenceID:     0,
		EndSequenceID:       10,
		EndMinuteSinceEpoch: 1000,
		DomainSeparator:     domainSeparator,
		NodeIDs:             []uint32{1},
		Payers:              payers,
	})
	require.NoError(t, err)

	stored := storeReport(t, store, &report.PayerReport)
	require.NoError(t, store.SetReportSubmitted(t.Context(), stored.ID, 0))

	verifier.EXPECT().
		GetPayerMap(mock.Anything, mock.Anything).
		Return(payers, nil)

	reportsManager.EXPECT().
		SettlementSummary(mock.Anything, uint32(1), uint64(0)).
		Return(&blockchain.SettlementSummary{
			Offset:    0,
			IsSettled: false,
		}, nil)

	// Expect multiple SettleReport calls due to MAX_PROOF_ELEMENTS limit
	callCount := 0
	reportsManager.EXPECT().
		SettleReport(mock.Anything, uint32(1), uint64(0), mock.Anything).
		Run(func(_ context.Context, _ uint32, _ uint64, _ *merkle.MultiProof) {
			callCount++
		}).
		Return(nil).
		Times(3) // 250 payers / 100 max = 3 batches

	err = worker.SettleReports(t.Context())
	require.NoError(t, err)
	require.Equal(t, 3, callCount)

	got, err := store.FetchReport(t.Context(), stored.ID)
	require.NoError(t, err)
	require.Equal(
		t,
		payerreport.SubmissionStatus(payerreport.SubmissionSettled),
		got.SubmissionStatus,
	)
}

func TestSettleReportWithPartialOffset(t *testing.T) {
	worker, store, reportsManager, verifier := testSettlementWorker(t)

	payers := make(map[common.Address]currency.PicoDollar)
	for i := range 150 {
		// Create unique addresses
		addr := common.BigToAddress(big.NewInt(int64(i + 1)))
		payers[addr] = currency.PicoDollar(i + 1)
	}

	report, err := payerreport.BuildPayerReport(payerreport.BuildPayerReportParams{
		OriginatorNodeID:    1,
		StartSequenceID:     0,
		EndSequenceID:       10,
		EndMinuteSinceEpoch: 1000,
		DomainSeparator:     domainSeparator,
		NodeIDs:             []uint32{1},
		Payers:              payers,
	})
	require.NoError(t, err)

	stored := storeReport(t, store, &report.PayerReport)
	require.NoError(t, store.SetReportAttestationApproved(t.Context(), stored.ID))
	require.NoError(t, store.SetReportSubmitted(t.Context(), stored.ID, 0))

	verifier.EXPECT().
		GetPayerMap(mock.Anything, mock.Anything).
		Return(payers, nil)

	// Settlement already partially completed (first 100 elements done)
	reportsManager.EXPECT().
		SettlementSummary(mock.Anything, uint32(1), uint64(0)).
		Return(&blockchain.SettlementSummary{
			Offset:    100,
			IsSettled: false,
		}, nil)

	// Only expect one call for remaining 50 elements
	reportsManager.EXPECT().
		SettleReport(mock.Anything, uint32(1), uint64(0), mock.Anything).
		Return(nil).
		Once()

	err = worker.SettleReports(t.Context())
	require.NoError(t, err)

	got, err := store.FetchReport(t.Context(), stored.ID)
	require.NoError(t, err)
	require.Equal(
		t,
		payerreport.SubmissionStatus(payerreport.SubmissionSettled),
		got.SubmissionStatus,
	)
}

func TestSettleReportErrors(t *testing.T) {
	testCases := []struct {
		name          string
		setupMocks    func(*blockchainMocks.MockPayerReportsManager, *payerreportMocks.MockIPayerReportVerifier)
		expectError   bool
		expectedState payerreport.SubmissionStatus
	}{
		{
			name: "error getting payer map",
			setupMocks: func(rm *blockchainMocks.MockPayerReportsManager, v *payerreportMocks.MockIPayerReportVerifier) {
				v.EXPECT().
					GetPayerMap(mock.Anything, mock.Anything).
					Return(payerreport.PayerMap(nil), errors.New("payer map error"))
			},
			expectError:   true,
			expectedState: payerreport.SubmissionSubmitted,
		},
		{
			name: "error getting settlement summary",
			setupMocks: func(rm *blockchainMocks.MockPayerReportsManager, v *payerreportMocks.MockIPayerReportVerifier) {
				v.EXPECT().
					GetPayerMap(mock.Anything, mock.Anything).
					Return(payerreport.PayerMap{common.HexToAddress("0x1"): 100}, nil)

				rm.EXPECT().
					SettlementSummary(mock.Anything, mock.Anything, mock.Anything).
					Return(nil, errors.New("settlement summary error"))
			},
			expectError:   true,
			expectedState: payerreport.SubmissionSubmitted,
		},
		{
			name: "error settling report",
			setupMocks: func(rm *blockchainMocks.MockPayerReportsManager, v *payerreportMocks.MockIPayerReportVerifier) {
				v.EXPECT().
					GetPayerMap(mock.Anything, mock.Anything).
					Return(payerreport.PayerMap{common.HexToAddress("0x1"): 100}, nil)

				rm.EXPECT().
					SettlementSummary(mock.Anything, mock.Anything, mock.Anything).
					Return(&blockchain.SettlementSummary{Offset: 0, IsSettled: false}, nil)

				rm.EXPECT().
					SettleReport(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return(errors.New("settle error"))
			},
			expectError:   true,
			expectedState: payerreport.SubmissionSubmitted,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			worker, store, reportsManager, verifier := testSettlementWorker(t)

			payers := map[common.Address]currency.PicoDollar{
				common.HexToAddress("0x1"): 100,
			}

			report, err := payerreport.BuildPayerReport(payerreport.BuildPayerReportParams{
				OriginatorNodeID:    1,
				StartSequenceID:     0,
				EndSequenceID:       10,
				EndMinuteSinceEpoch: 1000,
				DomainSeparator:     domainSeparator,
				NodeIDs:             []uint32{1},
				Payers:              payers,
			})
			require.NoError(t, err)

			stored := storeReport(t, store, &report.PayerReport)
			require.NoError(t, store.SetReportAttestationApproved(t.Context(), stored.ID))
			require.NoError(t, store.SetReportSubmitted(t.Context(), stored.ID, 0))

			tc.setupMocks(reportsManager, verifier)

			err = worker.SettleReports(t.Context())
			if tc.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			got, err := store.FetchReport(t.Context(), stored.ID)
			require.NoError(t, err)
			require.Equal(t, tc.expectedState, got.SubmissionStatus)
		})
	}
}

func TestSettleReportWithNilIndex(t *testing.T) {
	worker, store, _, _ := testSettlementWorker(t)

	payers := map[common.Address]currency.PicoDollar{
		common.HexToAddress("0x1"): 100,
	}

	report, err := payerreport.BuildPayerReport(payerreport.BuildPayerReportParams{
		OriginatorNodeID:    1,
		StartSequenceID:     0,
		EndSequenceID:       10,
		EndMinuteSinceEpoch: 1000,
		DomainSeparator:     domainSeparator,
		NodeIDs:             []uint32{1},
		Payers:              payers,
	})
	require.NoError(t, err)

	stored := storeReport(t, store, &report.PayerReport)

	// Manually fetch and settle without setting SubmittedReportIndex
	reportWithStatus, err := store.FetchReport(t.Context(), stored.ID)
	require.NoError(t, err)
	require.Nil(t, reportWithStatus.SubmittedReportIndex)

	// Try to settle with nil index
	err = worker.settleReport(t.Context(), reportWithStatus)
	require.Error(t, err)
	require.Contains(t, err.Error(), "report index is nil")
}

func TestSettleMultipleReports(t *testing.T) {
	worker, store, reportsManager, verifier := testSettlementWorker(t)

	// Create multiple reports
	reports := make([]*payerreport.PayerReportWithInputs, 3)
	for i := range 3 {
		payers := map[common.Address]currency.PicoDollar{
			common.HexToAddress("0x1"): currency.PicoDollar(100 * (i + 1)),
		}

		report, err := payerreport.BuildPayerReport(payerreport.BuildPayerReportParams{
			OriginatorNodeID:    1,
			StartSequenceID:     uint64(i * 10),
			EndSequenceID:       uint64((i + 1) * 10),
			EndMinuteSinceEpoch: uint32(1000 + i),
			DomainSeparator:     domainSeparator,
			NodeIDs:             []uint32{1},
			Payers:              payers,
		})
		require.NoError(t, err)

		reports[i] = report
		stored := storeReport(t, store, &report.PayerReport)
		require.NoError(t, store.SetReportAttestationApproved(t.Context(), stored.ID))
		require.NoError(t, store.SetReportSubmitted(t.Context(), stored.ID, int32(i)))

		verifier.EXPECT().
			GetPayerMap(mock.Anything, mock.Anything).
			Return(payers, nil)

		reportsManager.EXPECT().
			SettlementSummary(mock.Anything, uint32(1), uint64(i)).
			Return(&blockchain.SettlementSummary{
				Offset:    0,
				IsSettled: false,
			}, nil)

		reportsManager.EXPECT().
			SettleReport(mock.Anything, uint32(1), uint64(i), mock.Anything).
			Return(nil)
	}

	err := worker.SettleReports(t.Context())
	require.NoError(t, err)

	// Verify all reports are settled
	for _, report := range reports {
		got, err := store.FetchReport(t.Context(), report.ID)
		require.NoError(t, err)
		require.Equal(
			t,
			payerreport.SubmissionStatus(payerreport.SubmissionSettled),
			got.SubmissionStatus,
		)
	}
}

func TestSettleReportsWithPartialFailure(t *testing.T) {
	worker, store, reportsManager, verifier := testSettlementWorker(t)

	// Create two reports, first will fail, second will succeed
	payers1 := map[common.Address]currency.PicoDollar{
		common.HexToAddress("0x1"): 100,
	}
	report1, err := payerreport.BuildPayerReport(payerreport.BuildPayerReportParams{
		OriginatorNodeID:    1,
		StartSequenceID:     0,
		EndSequenceID:       10,
		EndMinuteSinceEpoch: 1000,
		DomainSeparator:     domainSeparator,
		NodeIDs:             []uint32{1},
		Payers:              payers1,
	})
	require.NoError(t, err)
	stored1 := storeReport(t, store, &report1.PayerReport)
	require.NoError(t, store.SetReportAttestationApproved(t.Context(), stored1.ID))
	require.NoError(t, store.SetReportSubmitted(t.Context(), stored1.ID, 0))
	// Make sure the created_at is different between the two reports
	time.Sleep(1 * time.Millisecond)

	payers2 := map[common.Address]currency.PicoDollar{
		common.HexToAddress("0x2"): 200,
	}
	report2, err := payerreport.BuildPayerReport(payerreport.BuildPayerReportParams{
		OriginatorNodeID:    2,
		StartSequenceID:     10,
		EndSequenceID:       20,
		EndMinuteSinceEpoch: 2000,
		DomainSeparator:     domainSeparator,
		NodeIDs:             []uint32{1},
		Payers:              payers2,
	})
	require.NoError(t, err)
	stored2 := storeReport(t, store, &report2.PayerReport)
	require.NoError(t, store.SetReportAttestationApproved(t.Context(), stored2.ID))
	require.NoError(t, store.SetReportSubmitted(t.Context(), stored2.ID, 1))

	// Reports are processed in database order, so we need to set up expectations for both
	// First report (index 0) will fail on GetPayerMap
	verifier.EXPECT().
		GetPayerMap(mock.Anything, mock.MatchedBy(func(report *payerreport.PayerReport) bool {
			return report.OriginatorNodeID == 1
		})).
		Return(payerreport.PayerMap(nil), errors.New("verifier error")).
		Once()

	// Second report (index 1) will succeed
	verifier.EXPECT().
		GetPayerMap(mock.Anything, mock.MatchedBy(func(report *payerreport.PayerReport) bool {
			return report.OriginatorNodeID == 2
		})).
		Return(payers2, nil).
		Once()

	// We need to accept either index 0 or 1 for settlement summary since database order is not guaranteed
	reportsManager.EXPECT().
		SettlementSummary(mock.Anything, uint32(2), mock.Anything).
		Return(&blockchain.SettlementSummary{
			Offset:    0,
			IsSettled: false,
		}, nil).
		Maybe()

	reportsManager.EXPECT().
		SettleReport(mock.Anything, uint32(2), mock.Anything, mock.Anything).
		Return(nil).
		Maybe()

	err = worker.SettleReports(t.Context())
	require.Error(t, err) // Should return the latest error
	// First report should still be submitted
	got1, err := store.FetchReport(t.Context(), stored1.ID)
	require.NoError(t, err)
	require.Equal(
		t,
		payerreport.SubmissionStatus(payerreport.SubmissionSubmitted),
		got1.SubmissionStatus,
	)

	// Second report should be settled
	got2, err := store.FetchReport(t.Context(), stored2.ID)
	require.NoError(t, err)
	require.Equal(
		t,
		payerreport.SubmissionStatus(payerreport.SubmissionSettled),
		got2.SubmissionStatus,
	)
}
