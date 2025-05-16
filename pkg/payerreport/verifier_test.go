package payerreport

import (
	"context"
	"database/sql"
	"testing"
	"time"

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
	originatorID uint32,
	payerAddress common.Address,
	timestamp time.Time,
	sequenceID uint64,
) envelopeCreateParams {
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
	}{
		{
			name: "one message per payer in the report",
			messagesToInsert: []envelopeCreateParams{
				newEnvelopeCreateParams(originatorID, payerAddress1, getMinute(1), 0),
				newEnvelopeCreateParams(originatorID, payerAddress2, getMinute(1), 1),
				// This message is in the last minute of the report. Will be ignored by the generator
				newEnvelopeCreateParams(originatorID, payerAddress1, getMinute(2), 2),
			},
		},
		{
			name: "two messages per payer",
			messagesToInsert: []envelopeCreateParams{
				newEnvelopeCreateParams(originatorID, payerAddress1, getMinute(1), 0),
				newEnvelopeCreateParams(originatorID, payerAddress2, getMinute(1), 1),
				newEnvelopeCreateParams(originatorID, payerAddress2, getMinute(2), 2),
				newEnvelopeCreateParams(originatorID, payerAddress2, getMinute(2), 3),
				newEnvelopeCreateParams(originatorID, payerAddress1, getMinute(3), 4),
			},
		},
		{
			name: "messages exist but not enough for a report",
			messagesToInsert: []envelopeCreateParams{
				newEnvelopeCreateParams(originatorID, payerAddress1, getMinute(1), 0),
				newEnvelopeCreateParams(originatorID, payerAddress1, getMinute(1), 1),
			},
		},
		{
			name:             "no messages",
			messagesToInsert: []envelopeCreateParams{},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			db, generator := setupGenerator(t)
			logger := testutils.NewLog(t).With(zap.String("test_case", testCase.name))
			store := NewStore(db, logger)
			verifier := NewPayerReportVerifier(logger, store)

			for _, message := range testCase.messagesToInsert {
				insertEnvelope(t, db, message)
			}

			report, err := generator.GenerateReport(t.Context(), PayerReportGenerationParams{
				OriginatorID:            uint32(originatorID),
				LastReportEndSequenceID: 0,
			})
			generator.log.Info("report", zap.Any("report", report))
			require.NoError(t, err)

			isValid, err := verifier.IsValidReport(t.Context(), nil, &report.PayerReport)
			require.NoError(t, err)
			require.True(t, isValid)
		})
	}
}

func setupVerifier(t *testing.T) (*sql.DB, *PayerReportVerifier) {
	db, _ := testutils.NewDB(t, context.Background())
	logger := testutils.NewLog(t)
	store := NewStore(db, logger)
	verifier := NewPayerReportVerifier(logger, store)

	return db, verifier
}

func TestValidateReportTransition(t *testing.T) {
	// Test cases for report transition validation
	testCases := []struct {
		name          string
		prevReport    *PayerReport
		newReport     *PayerReport
		expectedValid bool
		expectedError error
	}{
		{
			name:       "valid first report",
			prevReport: nil,
			newReport: &PayerReport{
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
			newReport: &PayerReport{
				OriginatorNodeID: 1,
				StartSequenceID:  1, // Should be 0 for first report
				EndSequenceID:    10,
				PayersMerkleRoot: randomBytes32(),
				ActiveNodeIDs:    []uint32{1},
			},
			expectedValid: false,
			expectedError: ErrInvalidReportStart,
		},
		{
			name: "valid subsequent report",
			prevReport: &PayerReport{
				OriginatorNodeID: 1,
				StartSequenceID:  0,
				EndSequenceID:    10,
				PayersMerkleRoot: randomBytes32(),
				ActiveNodeIDs:    []uint32{1},
			},
			newReport: &PayerReport{
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
			prevReport: &PayerReport{
				OriginatorNodeID: 1,
				StartSequenceID:  0,
				EndSequenceID:    10,
				PayersMerkleRoot: randomBytes32(),
				ActiveNodeIDs:    []uint32{1},
			},
			newReport: &PayerReport{
				OriginatorNodeID: 2, // Different originator
				StartSequenceID:  10,
				EndSequenceID:    20,
				PayersMerkleRoot: randomBytes32(),
				ActiveNodeIDs:    []uint32{1},
			},
			expectedValid: false,
			expectedError: ErrMismatchOriginator,
		},
		{
			name: "invalid subsequent report - gap in sequence",
			prevReport: &PayerReport{
				OriginatorNodeID: 1,
				StartSequenceID:  0,
				EndSequenceID:    10,
				PayersMerkleRoot: randomBytes32(),
				ActiveNodeIDs:    []uint32{1},
			},
			newReport: &PayerReport{
				OriginatorNodeID: 1,
				StartSequenceID:  11, // Gap from prev report's end
				EndSequenceID:    20,
				PayersMerkleRoot: randomBytes32(),
				ActiveNodeIDs:    []uint32{1},
			},
			expectedValid: false,
			expectedError: ErrInvalidReportStart,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateReportTransition(tc.prevReport, tc.newReport)
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
		report        *PayerReport
		expectedValid bool
		expectedError error
	}{
		{
			name: "valid report structure",
			report: &PayerReport{
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
			report: &PayerReport{
				OriginatorNodeID: 0, // Invalid originator ID
				StartSequenceID:  0,
				EndSequenceID:    10,
				PayersMerkleRoot: randomBytes32(),
				ActiveNodeIDs:    []uint32{1},
			},
			expectedValid: false,
			expectedError: ErrInvalidOriginatorID,
		},
		{
			name: "no active nodes",
			report: &PayerReport{
				OriginatorNodeID: 1,
				StartSequenceID:  0,
				EndSequenceID:    10,
				PayersMerkleRoot: randomBytes32(),
				ActiveNodeIDs:    []uint32{}, // Empty node list
			},
			expectedValid: false,
			expectedError: ErrNoNodes,
		},
		{
			name: "invalid merkle root length",
			report: &PayerReport{
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
			report: &PayerReport{
				OriginatorNodeID: 1,
				StartSequenceID:  20, // Start > End
				EndSequenceID:    10,
				PayersMerkleRoot: randomBytes32(),
				ActiveNodeIDs:    []uint32{1},
			},
			expectedValid: false,
			expectedError: ErrInvalidReportStart,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateReportStructure(tc.report)
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

	validMerkleTree200, err := generateMerkleTree(payerMap{
		payerAddress: currency.PicoDollar(200),
	})
	require.NoError(t, err)
	validMerkleTree0, err := generateMerkleTree(payerMap{})
	require.NoError(t, err)
	invalidAmountTree, err := generateMerkleTree(payerMap{
		payerAddress: currency.PicoDollar(400),
	})
	require.NoError(t, err)
	invalidPayerTree, err := generateMerkleTree(payerMap{
		testutils.RandomAddress(): currency.PicoDollar(100),
	})
	require.NoError(t, err)

	testCases := []struct {
		name          string
		report        *PayerReport
		expectedValid bool
		expectedError error
	}{
		{
			name: "empty report",
			report: &PayerReport{
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
			report: &PayerReport{
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
			report: &PayerReport{
				OriginatorNodeID: originatorID,
				StartSequenceID:  0,
				EndSequenceID:    1,
				PayersMerkleRoot: common.BytesToHash(invalidAmountTree.Root()),
				ActiveNodeIDs:    []uint32{1},
			},
			expectedValid: false,
			expectedError: ErrMerkleRootMismatch,
		},
		{
			name: "invalid merkle root - wrong payer",
			report: &PayerReport{
				OriginatorNodeID: originatorID,
				StartSequenceID:  0,
				EndSequenceID:    1,
				PayersMerkleRoot: common.BytesToHash(invalidPayerTree.Root()),
				ActiveNodeIDs:    []uint32{1},
			},
			expectedValid: false,
			expectedError: ErrMerkleRootMismatch,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			isValid, err := verifier.IsValidReport(context.Background(), nil, tc.report)
			if tc.expectedError != nil {
				require.Error(t, err)
				require.Equal(t, tc.expectedError, err)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tc.expectedValid, isValid)
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

	validMerkleTree, err := generateMerkleTree(payerMap{
		payerAddress: currency.PicoDollar(200),
	})
	require.NoError(t, err)

	testCases := []struct {
		name          string
		report        *PayerReport
		expectedValid bool
		expectedError error
	}{
		{
			name: "valid minute boundaries",
			report: &PayerReport{
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
			report: &PayerReport{
				OriginatorNodeID: originatorID,
				StartSequenceID:  0,
				EndSequenceID:    2, // Not the last message of minute3
				PayersMerkleRoot: common.BytesToHash(validMerkleTree.Root()),
				ActiveNodeIDs:    []uint32{1},
			},
			expectedValid: false,
			expectedError: ErrMessageNotAtMinuteEnd,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			isValid, err := verifier.IsValidReport(context.Background(), nil, tc.report)
			if tc.expectedError != nil {
				require.Error(t, err)
				require.Equal(t, tc.expectedError, err)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tc.expectedValid, isValid)
		})
	}
}
