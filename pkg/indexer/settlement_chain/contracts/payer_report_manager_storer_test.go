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
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	p "github.com/xmtp/xmtpd/pkg/abi/payerreportmanager"
	"github.com/xmtp/xmtpd/pkg/blockchain"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	payerreport "github.com/xmtp/xmtpd/pkg/payerreport"
	"github.com/xmtp/xmtpd/pkg/testutils"
	"github.com/xmtp/xmtpd/pkg/testutils/anvil"
	contractsMocks "github.com/xmtp/xmtpd/pkg/testutils/mocks/contracts"
	"github.com/xmtp/xmtpd/pkg/utils"
	re "github.com/xmtp/xmtpd/pkg/utils/retryerrors"
)

type payerReportManagerStorerTester struct {
	abi     *abi.ABI
	storer  *PayerReportManagerStorer
	queries *queries.Queries
}

func TestStorePayerReportManagerErrorNoTopics(t *testing.T) {
	tester := buildPayerReportManagerStorerTester(t)

	err := tester.storer.StoreLog(t.Context(), types.Log{})

	expectedErr := re.NewNonRecoverableError(
		ErrParsePayerReportManagerLog,
		errors.New("no topics"),
	)

	require.Error(t, err)
	require.ErrorAs(t, err, &expectedErr)
}

func TestStorePayerReportManagerErrorUnknownEvent(t *testing.T) {
	tester := buildPayerReportManagerStorerTester(t)

	log := types.Log{
		Topics: []common.Hash{common.HexToHash("UnknownEvent")},
	}

	err := tester.storer.StoreLog(t.Context(), log)

	expectedErr := re.NewNonRecoverableError(
		ErrParsePayerReportManagerLog,
		fmt.Errorf("no event with id: %#x", log.Topics[0].Hex()),
	)

	require.Error(t, err)
	require.ErrorAs(t, err, &expectedErr)
}

func TestStorePayerReportManagerPayerReportSubmitted(t *testing.T) {
	tester := buildPayerReportManagerStorerTester(t)

	originatorNodeID := uint32(1)
	payerReportIndex := uint64(42)

	log := tester.newPayerReportSubmittedLog(t, &payerreport.PayerReport{
		OriginatorNodeID:    originatorNodeID,
		StartSequenceID:     0,
		EndSequenceID:       100,
		EndMinuteSinceEpoch: 200,
		PayersMerkleRoot:    testutils.RandomInboxIDBytes(),
		ActiveNodeIDs:       []uint32{1, 2, 3},
	}, payerReportIndex)

	err := tester.storer.StoreLog(t.Context(), log)
	require.NoError(t, err)

	res, queryErr := tester.queries.FetchPayerReports(t.Context(), queries.FetchPayerReportsParams{
		OriginatorNodeID: utils.NewNullInt32(&originatorNodeID),
	})
	require.NoError(t, queryErr)
	require.Len(t, res, 1)

	require.Equal(t, int32(200), res[0].EndMinuteSinceEpoch)
	require.Equal(t, int64(0), res[0].StartSequenceID)
	require.Equal(t, int64(100), res[0].EndSequenceID)
	require.Equal(t, []int32{1, 2, 3}, res[0].ActiveNodeIds)
	require.True(t, res[0].SubmittedReportIndex.Valid, "SubmittedReportIndex should be set")
	require.Equal(t, int32(payerReportIndex), res[0].SubmittedReportIndex.Int32)
}

func TestStorePayerReportManagerPayerReportSubmittedIdempotency(t *testing.T) {
	tester := buildPayerReportManagerStorerTester(t)

	originatorNodeID := uint32(1)
	payerReportIndex := uint64(123)

	log := tester.newPayerReportSubmittedLog(t, &payerreport.PayerReport{
		OriginatorNodeID: originatorNodeID,
		StartSequenceID:  0,
		EndSequenceID:    100,
		PayersMerkleRoot: testutils.RandomInboxIDBytes(),
		ActiveNodeIDs:    []uint32{1, 2, 3},
	}, payerReportIndex)

	err := tester.storer.StoreLog(t.Context(), log)
	require.NoError(t, err)

	err = tester.storer.StoreLog(t.Context(), log)
	require.NoError(t, err)

	res, queryErr := tester.queries.FetchPayerReports(t.Context(), queries.FetchPayerReportsParams{
		OriginatorNodeID: utils.NewNullInt32(&originatorNodeID),
	})
	require.NoError(t, queryErr)
	require.Len(t, res, 1)
	require.True(t, res[0].SubmittedReportIndex.Valid, "SubmittedReportIndex should be set")
	require.Equal(t, int32(payerReportIndex), res[0].SubmittedReportIndex.Int32)
}

func TestStorePayerReportManagerPayerReportSubsetSettled(t *testing.T) {
	testCases := []struct {
		name             string
		originatorNodeID uint32
		payerReportIndex uint64
		count            uint32
		remaining        uint32
		feesSettled      *big.Int
		shouldSetSettled bool
	}{
		{
			name:             "Report fully settled when remaining is 0",
			originatorNodeID: uint32(1),
			payerReportIndex: 0,
			count:            10,
			remaining:        0,
			feesSettled:      big.NewInt(1000),
			shouldSetSettled: true,
		},
		{
			name:             "Report not settled when remaining is greater than 0",
			originatorNodeID: uint32(1),
			payerReportIndex: 0,
			count:            5,
			remaining:        5,
			feesSettled:      big.NewInt(500),
			shouldSetSettled: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			t.Cleanup(cancel)

			// Dependencies
			db, _ := testutils.NewDB(t, ctx)

			// Create mock contract
			mockContract := contractsMocks.NewMockPayerReportManagerContract(t)

			// Set up DOMAINSEPARATOR expectation
			domainSeparator := testutils.RandomBytes(32)
			var domainSeparatorArray [32]byte
			copy(domainSeparatorArray[:], domainSeparator)
			mockContract.EXPECT().
				DOMAINSEPARATOR(mock.Anything).
				Return(domainSeparatorArray, nil).
				Once()

			// Create storer with mock
			storer, err := NewPayerReportManagerStorer(db, testutils.NewLog(t), mockContract)
			require.NoError(t, err)

			// Create and store a PayerReport
			report := &payerreport.PayerReport{
				OriginatorNodeID:    tc.originatorNodeID,
				StartSequenceID:     0,
				EndSequenceID:       100,
				EndMinuteSinceEpoch: 200,
				PayersMerkleRoot:    testutils.RandomInboxIDBytes(),
				ActiveNodeIDs:       []uint32{1, 2, 3},
			}

			// Create helper tester for log generation
			tester := &payerReportManagerStorerTester{
				storer: storer,
			}

			// Use real ParsePayerReportSubmitted function
			submittedLog := tester.newPayerReportSubmittedLog(t, report, tc.payerReportIndex)

			// Create a real contract instance to use its parsing function
			realContract, err := p.NewPayerReportManager(
				common.HexToAddress("0x0000000000000000000000000000000000000000"),
				nil, // We don't need a client for parsing
			)
			require.NoError(t, err)

			mockContract.EXPECT().
				ParsePayerReportSubmitted(submittedLog).
				RunAndReturn(realContract.ParsePayerReportSubmitted)

			err = storer.StoreLog(ctx, submittedLog)
			require.NoError(t, err)

			// Create the PayerReportSubsetSettled event
			settledLog := testutils.BuildPayerReportSubsetSettledLog(
				t,
				tc.originatorNodeID,
				tc.payerReportIndex,
				tc.count,
				tc.remaining,
				tc.feesSettled,
			)

			mockContract.EXPECT().
				ParsePayerReportSubsetSettled(settledLog).
				RunAndReturn(realContract.ParsePayerReportSubsetSettled)

			// If remaining is 0, mock GetPayerReport
			if tc.remaining == 0 {
				mockReport := p.IPayerReportManagerPayerReport{
					StartSequenceId:     report.StartSequenceID,
					EndSequenceId:       report.EndSequenceID,
					EndMinuteSinceEpoch: uint32(report.EndMinuteSinceEpoch),
					PayersMerkleRoot:    report.PayersMerkleRoot,
					NodeIds:             report.ActiveNodeIDs,
					FeesSettled:         tc.feesSettled,
					IsSettled:           false,
				}
				mockContract.EXPECT().GetPayerReport(
					mock.AnythingOfType("*bind.CallOpts"),
					tc.originatorNodeID,
					mock.MatchedBy(func(index *big.Int) bool {
						return index.Cmp(big.NewInt(int64(tc.payerReportIndex))) == 0
					}),
				).Return(mockReport, nil)
			}

			err = storer.StoreLog(ctx, settledLog)
			require.NoError(t, err)

			// Verify the report status
			q := db.Query()
			res, queryErr := q.FetchPayerReports(ctx, queries.FetchPayerReportsParams{
				OriginatorNodeID: utils.NewNullInt32(&tc.originatorNodeID),
			})
			require.NoError(t, queryErr)
			require.Len(t, res, 1)

			// Check the submission status
			if tc.shouldSetSettled {
				// When remaining is 0, the report should be marked as settled (2)
				require.Equal(t, int16(2), res[0].SubmissionStatus)
			} else {
				// When remaining > 0, the report should remain as submitted (1)
				require.Equal(t, int16(1), res[0].SubmissionStatus)
			}
		})
	}
}

func TestStorePayerReportManagerPayerReportSubsetSettledIdempotency(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	// Dependencies
	db, _ := testutils.NewDB(t, ctx)

	// Create mock contract
	mockContract := contractsMocks.NewMockPayerReportManagerContract(t)

	// Set up DOMAINSEPARATOR expectation
	domainSeparator := testutils.RandomBytes(32)
	var domainSeparatorArray [32]byte
	copy(domainSeparatorArray[:], domainSeparator)
	mockContract.EXPECT().DOMAINSEPARATOR(mock.Anything).Return(domainSeparatorArray, nil).Once()

	// Create storer with mock
	storer, err := NewPayerReportManagerStorer(db, testutils.NewLog(t), mockContract)
	require.NoError(t, err)

	originatorNodeID := uint32(1)
	payerReportIndex := uint64(0)

	// Create a PayerReport
	report := &payerreport.PayerReport{
		OriginatorNodeID:    originatorNodeID,
		StartSequenceID:     0,
		EndSequenceID:       100,
		EndMinuteSinceEpoch: 200,
		PayersMerkleRoot:    testutils.RandomInboxIDBytes(),
		ActiveNodeIDs:       []uint32{1, 2, 3},
	}

	// Create helper tester for log generation
	tester := &payerReportManagerStorerTester{
		storer: storer,
	}

	// Use real ParsePayerReportSubmitted function
	submittedLog := tester.newPayerReportSubmittedLog(t, report, payerReportIndex)

	// Create a real contract instance to use its parsing function
	realContract, err := p.NewPayerReportManager(
		common.HexToAddress("0x0000000000000000000000000000000000000000"),
		nil, // We don't need a client for parsing
	)
	require.NoError(t, err)

	mockContract.EXPECT().
		ParsePayerReportSubmitted(submittedLog).
		RunAndReturn(realContract.ParsePayerReportSubmitted).
		Once()

	err = storer.StoreLog(ctx, submittedLog)
	require.NoError(t, err)

	// Create a PayerReportSubsetSettled event with remaining = 0
	settledLog := testutils.BuildPayerReportSubsetSettledLog(
		t,
		originatorNodeID,
		payerReportIndex,
		10,
		0,
		big.NewInt(1000),
	)

	mockContract.EXPECT().
		ParsePayerReportSubsetSettled(settledLog).
		RunAndReturn(realContract.ParsePayerReportSubsetSettled)

	// Mock GetPayerReport (called twice for idempotency)
	mockReport := p.IPayerReportManagerPayerReport{
		StartSequenceId:     report.StartSequenceID,
		EndSequenceId:       report.EndSequenceID,
		EndMinuteSinceEpoch: uint32(report.EndMinuteSinceEpoch),
		PayersMerkleRoot:    report.PayersMerkleRoot,
		NodeIds:             report.ActiveNodeIDs,
		FeesSettled:         big.NewInt(1000),
		IsSettled:           false,
	}
	mockContract.EXPECT().GetPayerReport(
		mock.Anything,
		originatorNodeID,
		mock.MatchedBy(func(index *big.Int) bool {
			return index.Cmp(big.NewInt(int64(payerReportIndex))) == 0
		}),
	).Return(mockReport, nil).Twice()

	// Store the event twice to test idempotency
	err = storer.StoreLog(ctx, settledLog)
	require.NoError(t, err)

	err = storer.StoreLog(ctx, settledLog)
	require.NoError(t, err)

	// Verify the report exists and is marked as settled
	res, queryErr := db.Query().FetchPayerReports(ctx, queries.FetchPayerReportsParams{
		OriginatorNodeID: utils.NewNullInt32(&originatorNodeID),
	})
	require.NoError(t, queryErr)
	require.Len(t, res, 1)
	require.Equal(t, int16(2), res[0].SubmissionStatus) // 2 = settled
}

func buildPayerReportManagerStorerTester(t *testing.T) *payerReportManagerStorerTester {
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
	contract, err := p.NewPayerReportManager(
		common.HexToAddress(config.SettlementChain.PayerReportManagerAddress),
		client,
	)
	require.NoError(t, err)

	// Storer and ABI.
	storer, err := NewPayerReportManagerStorer(db, testutils.NewLog(t), contract)
	require.NoError(t, err)

	abi, err := p.PayerReportManagerMetaData.GetAbi()
	require.NoError(t, err)

	return &payerReportManagerStorerTester{
		abi:     abi,
		storer:  storer,
		queries: db.Query(),
	}
}

func (st *payerReportManagerStorerTester) newPayerReportSubmittedLog(
	t *testing.T,
	report *payerreport.PayerReport,
	payerReportIndex uint64,
) types.Log {
	return testutils.BuildPayerReportSubmittedEvent(
		t,
		report.OriginatorNodeID,
		payerReportIndex,
		report.StartSequenceID,
		report.EndSequenceID,
		uint64(report.EndMinuteSinceEpoch),
		report.PayersMerkleRoot,
		report.ActiveNodeIDs,
	)
}
