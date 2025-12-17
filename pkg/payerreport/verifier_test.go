package payerreport_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/xmtp/xmtpd/pkg/payerreport"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/currency"
	dbHelpers "github.com/xmtp/xmtpd/pkg/db"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/testutils"
	envelopeTestUtils "github.com/xmtp/xmtpd/pkg/testutils/envelopes"
	"github.com/xmtp/xmtpd/pkg/utils"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

type envelopeCreateParams struct {
	originatorID uint32
	payerAddress common.Address
	timestamp    time.Time
	sequenceID   uint64
	cost         currency.PicoDollar
}

func newEnvelopeCreateParams(
	t *testing.T,
	originatorID uint32,
	payerAddress common.Address,
	timestamp time.Time,
	sequenceID uint64,
) envelopeCreateParams {
	require.NotEqualValues(t, 0, sequenceID)
	return envelopeCreateParams{
		originatorID: originatorID,
		payerAddress: payerAddress,
		timestamp:    timestamp,
		sequenceID:   sequenceID,
		cost:         currency.PicoDollar(100),
	}
}

func insertEnvelope(t *testing.T, db *sql.DB, params envelopeCreateParams) {
	payerID := testutils.CreatePayer(t, db, params.payerAddress.Hex())

	envelope := envelopeTestUtils.CreateOriginatorEnvelopeWithTimestamp(
		t,
		uint32(params.originatorID),
		params.sequenceID,
		params.timestamp,
	)

	envelopeBytes, err := proto.Marshal(envelope)
	require.NoError(t, err)

	_, err = dbHelpers.InsertGatewayEnvelopeAndIncrementUnsettledUsage(
		context.Background(),
		db,
		queries.InsertGatewayEnvelopeParams{
			OriginatorNodeID:     int32(params.originatorID),
			OriginatorSequenceID: int64(params.sequenceID),
			OriginatorEnvelope:   envelopeBytes,
			Topic:                testutils.RandomBytes(32),
			PayerID:              dbHelpers.NullInt32(payerID),
		},
		queries.IncrementUnsettledUsageParams{
			PayerID:           payerID,
			OriginatorID:      int32(params.originatorID),
			SequenceID:        int64(params.sequenceID),
			MinutesSinceEpoch: utils.MinutesSinceEpoch(params.timestamp),
			SpendPicodollars:  int64(params.cost),
		},
	)
	require.NoError(t, err)
}

func TestValidFirstReport(t *testing.T) {
	payerAddress1 := testutils.RandomAddress()
	payerAddress2 := testutils.RandomAddress()
	originatorID := uint32(100)

	testCases := []struct {
		name             string
		messagesToInsert []envelopeCreateParams
		expectNilReport  bool
	}{
		{
			name: "one message per payer in the report",
			messagesToInsert: []envelopeCreateParams{
				newEnvelopeCreateParams(t, originatorID, payerAddress1, getMinute(1), 1),
				newEnvelopeCreateParams(t, originatorID, payerAddress2, getMinute(1), 2),
				// This message is in the last minute of the report. Will be ignored by the generator
				newEnvelopeCreateParams(t, originatorID, payerAddress1, getMinute(2), 3),
			},
			expectNilReport: false,
		},
		{
			name: "two messages per payer",
			messagesToInsert: []envelopeCreateParams{
				newEnvelopeCreateParams(t, originatorID, payerAddress1, getMinute(1), 1),
				newEnvelopeCreateParams(t, originatorID, payerAddress2, getMinute(1), 2),
				newEnvelopeCreateParams(t, originatorID, payerAddress2, getMinute(2), 3),
				newEnvelopeCreateParams(t, originatorID, payerAddress2, getMinute(2), 4),
				newEnvelopeCreateParams(t, originatorID, payerAddress1, getMinute(3), 5),
			},
			expectNilReport: false,
		},
		{
			name: "messages exist but not enough for a report",
			messagesToInsert: []envelopeCreateParams{
				newEnvelopeCreateParams(t, originatorID, payerAddress1, getMinute(1), 1),
				newEnvelopeCreateParams(t, originatorID, payerAddress1, getMinute(1), 2),
			},
			expectNilReport: true,
		},
		{
			name:             "no messages",
			messagesToInsert: []envelopeCreateParams{},
			expectNilReport:  true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			db, generator := setupGenerator(t)
			logger := testutils.NewLog(t).With(zap.String("test_case", testCase.name))
			store := payerreport.NewStore(logger, db)
			verifier := payerreport.NewPayerReportVerifier(logger, store)

			for _, message := range testCase.messagesToInsert {
				insertEnvelope(t, db.DB(), message)
			}

			report, err := generator.GenerateReport(
				t.Context(),
				payerreport.PayerReportGenerationParams{
					OriginatorID:            uint32(originatorID),
					LastReportEndSequenceID: 0,
				},
			)
			logger.Info("report", zap.Any("report", report))
			require.NoError(t, err)

			if testCase.expectNilReport {
				require.Nil(t, report)
				return
			}

			verifyResult, err := verifier.VerifyReport(t.Context(), nil, &report.PayerReport)
			require.NoError(t, err)
			require.True(t, verifyResult.IsValid)
			require.Equal(t, "valid report", verifyResult.Reason)
		})
	}
}

func setupVerifier(t *testing.T) (*sql.DB, *payerreport.PayerReportVerifier) {
	db, _ := testutils.NewDB(t, context.Background())
	logger := testutils.NewLog(t)
	store := payerreport.NewStore(logger, db)
	verifier := payerreport.NewPayerReportVerifier(logger, store)

	return db.DB(), verifier
}

func TestValidateReportTransition(t *testing.T) {
	// Test cases for report transition validation
	testCases := []struct {
		name          string
		prevReport    *payerreport.PayerReport
		newReport     *payerreport.PayerReport
		expectedValid bool
		expectedError error
	}{
		{
			name:       "valid first report",
			prevReport: nil,
			newReport: &payerreport.PayerReport{
				OriginatorNodeID: 1,
				StartSequenceID:  0,
				EndSequenceID:    10,
				PayersMerkleRoot: randomBytes32(),
				ActiveNodeIDs:    []uint32{1},
			},
			expectedValid: true,
		},
		{
			name:       "invalid first report - non-zero start",
			prevReport: nil,
			newReport: &payerreport.PayerReport{
				OriginatorNodeID: 1,
				StartSequenceID:  1, // Should be 0 for first report
				EndSequenceID:    10,
				PayersMerkleRoot: randomBytes32(),
				ActiveNodeIDs:    []uint32{1},
			},
			expectedValid: false,
			expectedError: payerreport.ErrInvalidReportStart,
		},
		{
			name: "valid subsequent report",
			prevReport: &payerreport.PayerReport{
				OriginatorNodeID: 1,
				StartSequenceID:  0,
				EndSequenceID:    10,
				PayersMerkleRoot: randomBytes32(),
				ActiveNodeIDs:    []uint32{1},
			},
			newReport: &payerreport.PayerReport{
				OriginatorNodeID: 1,
				StartSequenceID:  10, // Matches prev report's end
				EndSequenceID:    20,
				PayersMerkleRoot: randomBytes32(),
				ActiveNodeIDs:    []uint32{1},
			},
			expectedValid: true,
		},
		{
			name: "invalid subsequent report - mismatched originator",
			prevReport: &payerreport.PayerReport{
				OriginatorNodeID: 1,
				StartSequenceID:  0,
				EndSequenceID:    10,
				PayersMerkleRoot: randomBytes32(),
				ActiveNodeIDs:    []uint32{1},
			},
			newReport: &payerreport.PayerReport{
				OriginatorNodeID: 2, // Different originator
				StartSequenceID:  10,
				EndSequenceID:    20,
				PayersMerkleRoot: randomBytes32(),
				ActiveNodeIDs:    []uint32{1},
			},
			expectedValid: false,
			expectedError: payerreport.ErrMismatchOriginator,
		},
		{
			name: "invalid subsequent report - gap in sequence",
			prevReport: &payerreport.PayerReport{
				OriginatorNodeID: 1,
				StartSequenceID:  0,
				EndSequenceID:    10,
				PayersMerkleRoot: randomBytes32(),
				ActiveNodeIDs:    []uint32{1},
			},
			newReport: &payerreport.PayerReport{
				OriginatorNodeID: 1,
				StartSequenceID:  11, // Gap from prev report's end
				EndSequenceID:    20,
				PayersMerkleRoot: randomBytes32(),
				ActiveNodeIDs:    []uint32{1},
			},
			expectedValid: false,
			expectedError: payerreport.ErrInvalidReportStart,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := payerreport.ValidateReportTransitionTestBinding(tc.prevReport, tc.newReport)
			if tc.expectedError != nil {
				require.Error(t, err)
				require.Equal(t, tc.expectedError, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidateReportStructure(t *testing.T) {
	testCases := []struct {
		name          string
		report        *payerreport.PayerReport
		expectedValid bool
		expectedError error
	}{
		{
			name: "valid report structure",
			report: &payerreport.PayerReport{
				OriginatorNodeID: 1,
				StartSequenceID:  0,
				EndSequenceID:    10,
				PayersMerkleRoot: randomBytes32(),
				ActiveNodeIDs:    []uint32{1},
			},
			expectedValid: true,
		},
		{
			name: "invalid originator ID",
			report: &payerreport.PayerReport{
				OriginatorNodeID: 0, // Invalid originator ID
				StartSequenceID:  0,
				EndSequenceID:    10,
				PayersMerkleRoot: randomBytes32(),
				ActiveNodeIDs:    []uint32{1},
			},
			expectedValid: false,
			expectedError: payerreport.ErrInvalidOriginatorID,
		},
		{
			name: "no active nodes",
			report: &payerreport.PayerReport{
				OriginatorNodeID: 1,
				StartSequenceID:  0,
				EndSequenceID:    10,
				PayersMerkleRoot: randomBytes32(),
				ActiveNodeIDs:    []uint32{}, // Empty node list
			},
			expectedValid: false,
			expectedError: payerreport.ErrNoNodes,
		},
		{
			name: "invalid merkle root length",
			report: &payerreport.PayerReport{
				OriginatorNodeID: 1,
				StartSequenceID:  0,
				EndSequenceID:    10,
				PayersMerkleRoot: [32]byte{}, // Zero merkle root
				ActiveNodeIDs:    []uint32{1},
			},
			expectedValid: true, // Empty merkle root is valid
		},
		{
			name: "invalid sequence ID order",
			report: &payerreport.PayerReport{
				OriginatorNodeID: 1,
				StartSequenceID:  20, // Start > End
				EndSequenceID:    10,
				PayersMerkleRoot: randomBytes32(),
				ActiveNodeIDs:    []uint32{1},
			},
			expectedValid: false,
			expectedError: payerreport.ErrInvalidReportStart,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := payerreport.ValidateReportStructureTestBinding(tc.report)
			if tc.expectedError != nil {
				require.Error(t, err)
				require.Equal(t, tc.expectedError, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidateMerkleRoot(t *testing.T) {
	db, verifier := setupVerifier(t)

	// Create test data
	originatorID := uint32(1)
	payerAddress := testutils.RandomAddress()
	now := time.Now().UTC()
	minute1 := time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		now.Hour(),
		now.Minute(),
		0,
		0,
		time.UTC,
	)
	minute2 := minute1.Add(time.Minute)

	// Insert test envelopes
	insertEnvelope(t, db, envelopeCreateParams{
		originatorID: originatorID,
		sequenceID:   0,
		payerAddress: payerAddress,
		timestamp:    minute1,
		cost:         currency.PicoDollar(100),
	})

	insertEnvelope(t, db, envelopeCreateParams{
		originatorID: originatorID,
		sequenceID:   1,
		payerAddress: payerAddress,
		timestamp:    minute2,
		cost:         currency.PicoDollar(100),
	})

	validMerkleTree200, err := payerreport.GenerateMerkleTree(payerreport.PayerMap{
		payerAddress: currency.PicoDollar(200),
	})
	require.NoError(t, err)
	validMerkleTree0, err := payerreport.GenerateMerkleTree(payerreport.PayerMap{})
	require.NoError(t, err)
	invalidAmountTree, err := payerreport.GenerateMerkleTree(payerreport.PayerMap{
		payerAddress: currency.PicoDollar(400),
	})
	require.NoError(t, err)
	invalidPayerTree, err := payerreport.GenerateMerkleTree(payerreport.PayerMap{
		testutils.RandomAddress(): currency.PicoDollar(100),
	})
	require.NoError(t, err)

	testCases := []struct {
		name          string
		report        *payerreport.PayerReport
		expectedValid bool
		expectedError error
	}{
		{
			name: "empty report",
			report: &payerreport.PayerReport{
				OriginatorNodeID: originatorID,
				StartSequenceID:  0,
				EndSequenceID:    0,
				PayersMerkleRoot: common.BytesToHash(validMerkleTree0.Root()),
				ActiveNodeIDs:    []uint32{1},
			},
			expectedValid: true,
		},
		{
			name: "valid merkle root",
			report: &payerreport.PayerReport{
				OriginatorNodeID: originatorID,
				StartSequenceID:  0,
				EndSequenceID:    1,
				PayersMerkleRoot: common.BytesToHash(validMerkleTree200.Root()),
				ActiveNodeIDs:    []uint32{1},
			},
			expectedValid: true,
		},
		{
			name: "invalid merkle root - wrong amount",
			report: &payerreport.PayerReport{
				OriginatorNodeID: originatorID,
				StartSequenceID:  0,
				EndSequenceID:    1,
				PayersMerkleRoot: common.BytesToHash(invalidAmountTree.Root()),
				ActiveNodeIDs:    []uint32{1},
			},
			expectedValid: false,
		},
		{
			name: "invalid merkle root - wrong payer",
			report: &payerreport.PayerReport{
				OriginatorNodeID: originatorID,
				StartSequenceID:  0,
				EndSequenceID:    1,
				PayersMerkleRoot: common.BytesToHash(invalidPayerTree.Root()),
				ActiveNodeIDs:    []uint32{1},
			},
			expectedValid: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			verifyResult, err := verifier.VerifyReport(context.Background(), nil, tc.report)
			if tc.expectedError != nil {
				require.Error(t, err)
				require.Equal(t, tc.expectedError, err)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tc.expectedValid, verifyResult.IsValid)
		})
	}
}

func TestValidateMinuteBoundaries(t *testing.T) {
	db, verifier := setupVerifier(t)

	// Create test data
	originatorID := uint32(1)
	payerAddress := testutils.RandomAddress()
	now := time.Now().UTC()
	minute1 := time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		now.Hour(),
		now.Minute(),
		0,
		0,
		time.UTC,
	)
	minute2 := minute1.Add(time.Minute)
	minute3 := minute2.Add(time.Minute)

	// Insert test envelopes
	insertEnvelope(t, db, envelopeCreateParams{
		originatorID: originatorID,
		sequenceID:   0,
		payerAddress: payerAddress,
		timestamp:    minute1,
		cost:         currency.PicoDollar(100),
	})

	insertEnvelope(t, db, envelopeCreateParams{
		originatorID: originatorID,
		sequenceID:   1,
		payerAddress: payerAddress,
		timestamp:    minute2,
		cost:         currency.PicoDollar(100),
	})

	insertEnvelope(t, db, envelopeCreateParams{
		originatorID: originatorID,
		sequenceID:   2,
		payerAddress: payerAddress,
		timestamp:    minute3,
		cost:         currency.PicoDollar(100),
	})

	insertEnvelope(t, db, envelopeCreateParams{
		originatorID: originatorID,
		sequenceID:   3,
		payerAddress: payerAddress,
		timestamp:    minute3,
		cost:         currency.PicoDollar(100),
	})

	validMerkleTree, err := payerreport.GenerateMerkleTree(payerreport.PayerMap{
		payerAddress: currency.PicoDollar(200),
	})
	require.NoError(t, err)

	testCases := []struct {
		name          string
		report        *payerreport.PayerReport
		expectedValid bool
		expectedError error
	}{
		{
			name: "valid minute boundaries",
			report: &payerreport.PayerReport{
				OriginatorNodeID: originatorID,
				StartSequenceID:  0,
				EndSequenceID:    1, // Last message of minute2
				PayersMerkleRoot: common.BytesToHash(validMerkleTree.Root()),
				ActiveNodeIDs:    []uint32{1},
			},
			expectedValid: true,
		},
		{
			name: "invalid minute boundary - not last message",
			report: &payerreport.PayerReport{
				OriginatorNodeID: originatorID,
				StartSequenceID:  0,
				EndSequenceID:    2, // Not the last message of minute3
				PayersMerkleRoot: common.BytesToHash(validMerkleTree.Root()),
				ActiveNodeIDs:    []uint32{1},
			},
			expectedValid: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			verifyResult, err := verifier.VerifyReport(context.Background(), nil, tc.report)
			if tc.expectedError != nil {
				require.Error(t, err)
				require.Equal(t, tc.expectedError, err)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tc.expectedValid, verifyResult.IsValid)
		})
	}
}
