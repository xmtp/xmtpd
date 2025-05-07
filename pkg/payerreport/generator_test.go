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
	"google.golang.org/protobuf/proto"
)

func setup(t *testing.T) (*sql.DB, *PayerReportGenerator) {
	db, _ := testutils.NewDB(t, context.Background())

	generator := NewPayerReportGenerator(testutils.NewLog(t), queries.New(db))

	return db, generator
}

func addEnvelope(
	t *testing.T,
	db *sql.DB,
	originatorID int32,
	sequenceID int64,
	payerAddress common.Address,
	timestamp time.Time,
) {
	payerID := testutils.CreatePayer(t, db, payerAddress.Hex())

	envelope := envelopeTestUtils.CreateOriginatorEnvelopeWithTimestamp(
		t,
		uint32(originatorID),
		uint64(sequenceID),
		timestamp,
	)

	envelopeBytes, err := proto.Marshal(envelope)
	require.NoError(t, err)

	_, err = dbHelpers.InsertGatewayEnvelopeAndIncrementUnsettledUsage(
		context.Background(),
		db,
		queries.InsertGatewayEnvelopeParams{
			OriginatorNodeID:     originatorID,
			OriginatorSequenceID: sequenceID,
			OriginatorEnvelope:   envelopeBytes,
			Topic:                testutils.RandomBytes(32),
			PayerID:              dbHelpers.NullInt32(payerID),
		},
		queries.IncrementUnsettledUsageParams{
			PayerID:           payerID,
			OriginatorID:      originatorID,
			MinutesSinceEpoch: utils.MinutesSinceEpoch(timestamp),
			SpendPicodollars:  100,
		},
	)
	require.NoError(t, err)
}

func getMinute(minutesSinceEpoch int) time.Time {
	return time.Unix(0, 0).Add(time.Duration(minutesSinceEpoch) * time.Minute)
}

func TestFirstReport(t *testing.T) {
	db, generator := setup(t)

	payerAddress := testutils.RandomAddress()
	originatorID := testutils.RandomInt32()

	// Two envelopes in the first minute since the epoch
	addEnvelope(t, db, originatorID, 1, payerAddress, getMinute(1))
	addEnvelope(t, db, originatorID, 2, payerAddress, getMinute(1))

	// One envelope in the second minute since the epoch
	addEnvelope(t, db, originatorID, 3, payerAddress, getMinute(2))

	report, err := generator.GenerateReport(context.Background(), PayerReportGenerationParams{
		OriginatorID:            uint32(originatorID),
		LastReportEndSequenceID: 0,
	})
	require.NoError(t, err)

	require.Equal(t, uint32(originatorID), report.OriginatorNodeID)
	// Make sure the end sequence ID is the last sequence ID from the previous minute
	require.Equal(t, uint64(2), report.EndSequenceID)
	require.Equal(t, currency.PicoDollar(200), report.Payers[payerAddress])
}

func TestReportWithMultiplePayers(t *testing.T) {
	db, generator := setup(t)

	payerAddress1 := testutils.RandomAddress()
	payerAddress2 := testutils.RandomAddress()
	originatorID := testutils.RandomInt32()

	addEnvelope(t, db, originatorID, 1, payerAddress1, getMinute(1))
	addEnvelope(t, db, originatorID, 2, payerAddress2, getMinute(1))
	addEnvelope(t, db, originatorID, 3, payerAddress1, getMinute(2))

	report, err := generator.GenerateReport(context.Background(), PayerReportGenerationParams{
		OriginatorID:            uint32(originatorID),
		LastReportEndSequenceID: 0,
	})
	require.NoError(t, err)

	require.Equal(t, uint32(originatorID), report.OriginatorNodeID)
	require.Equal(t, uint64(2), report.EndSequenceID)
	require.Equal(t, currency.PicoDollar(100), report.Payers[payerAddress1])
	require.Equal(t, currency.PicoDollar(100), report.Payers[payerAddress2])
}

func TestReportWithNoMessages(t *testing.T) {
	_, generator := setup(t)

	originatorID := testutils.RandomInt32()

	report, err := generator.GenerateReport(context.Background(), PayerReportGenerationParams{
		OriginatorID:            uint32(originatorID),
		LastReportEndSequenceID: 0,
	})
	require.NoError(t, err)

	require.Equal(t, uint32(originatorID), report.OriginatorNodeID)
	require.Equal(t, uint64(0), report.StartSequenceID)
	require.Equal(t, uint64(0), report.EndSequenceID)
	require.Equal(t, 0, len(report.Payers))
}

func TestSecondReportWithNoMessages(t *testing.T) {
	db, generator := setup(t)

	originatorID := testutils.RandomInt32()
	payerAddress1 := testutils.RandomAddress()

	addEnvelope(t, db, originatorID, 1, payerAddress1, getMinute(1))

	report, err := generator.GenerateReport(context.Background(), PayerReportGenerationParams{
		OriginatorID:            uint32(originatorID),
		LastReportEndSequenceID: 1,
	})
	require.NoError(t, err)

	require.Equal(t, uint32(originatorID), report.OriginatorNodeID)
	require.Equal(t, uint64(1), report.StartSequenceID)
	require.Equal(t, uint64(1), report.EndSequenceID)
}

func TestSecondReport(t *testing.T) {
	db, generator := setup(t)

	originatorID := testutils.RandomInt32()
	payerAddress := testutils.RandomAddress()

	addEnvelope(t, db, originatorID, 1, payerAddress, getMinute(1))
	addEnvelope(t, db, originatorID, 2, payerAddress, getMinute(1))
	addEnvelope(t, db, originatorID, 3, payerAddress, getMinute(2))

	report, err := generator.GenerateReport(context.Background(), PayerReportGenerationParams{
		OriginatorID:            uint32(originatorID),
		LastReportEndSequenceID: 0,
	})
	require.NoError(t, err)

	require.Equal(t, uint32(originatorID), report.OriginatorNodeID)
	require.Equal(t, uint64(0), report.StartSequenceID)
	require.Equal(t, uint64(2), report.EndSequenceID)
	require.Equal(t, currency.PicoDollar(200), report.Payers[payerAddress])

	addEnvelope(t, db, originatorID, 4, payerAddress, getMinute(3))
	addEnvelope(t, db, originatorID, 4, payerAddress, getMinute(4))

	report, err = generator.GenerateReport(context.Background(), PayerReportGenerationParams{
		OriginatorID:            uint32(originatorID),
		LastReportEndSequenceID: report.EndSequenceID,
	})
	require.NoError(t, err)

	require.Equal(t, uint32(originatorID), report.OriginatorNodeID)
	require.Equal(t, uint64(3), report.EndSequenceID)
	require.Equal(t, currency.PicoDollar(100), report.Payers[payerAddress])
}

// Make sure that we don't pick up sequence IDs from other originators in the report
func TestReportWithNoEnvelopesFromOriginator(t *testing.T) {
	db, generator := setup(t)

	originatorID := testutils.RandomInt32()
	otherOriginatorID := testutils.RandomInt32()
	payerAddress := testutils.RandomAddress()

	addEnvelope(t, db, otherOriginatorID, 1, payerAddress, getMinute(1))
	addEnvelope(t, db, otherOriginatorID, 2, payerAddress, getMinute(2))
	addEnvelope(t, db, otherOriginatorID, 3, payerAddress, getMinute(3))

	report, err := generator.GenerateReport(context.Background(), PayerReportGenerationParams{
		OriginatorID:            uint32(originatorID),
		LastReportEndSequenceID: 0,
	})
	require.NoError(t, err)

	require.Equal(t, uint32(originatorID), report.OriginatorNodeID)
	require.Equal(t, uint64(0), report.StartSequenceID)
	require.Equal(t, uint64(0), report.EndSequenceID)
	require.Equal(t, 0, len(report.Payers))
}
