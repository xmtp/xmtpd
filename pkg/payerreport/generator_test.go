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
	"github.com/xmtp/xmtpd/pkg/migrator"
	"github.com/xmtp/xmtpd/pkg/registry"
	"github.com/xmtp/xmtpd/pkg/testutils"
	envelopeTestUtils "github.com/xmtp/xmtpd/pkg/testutils/envelopes"
	registryTestUtils "github.com/xmtp/xmtpd/pkg/testutils/registry"
	"github.com/xmtp/xmtpd/pkg/utils"
	"google.golang.org/protobuf/proto"
)

var domainSeparator = common.BytesToHash(testutils.RandomBytes(32))

func setupGenerator(t *testing.T) (*dbHelpers.Handler, *payerreport.PayerReportGenerator) {
	db, _ := testutils.NewDB(t, context.Background())

	registry := registryTestUtils.CreateMockRegistry(t, []registry.Node{
		registryTestUtils.CreateNode(100, 100, testutils.RandomPrivateKey(t)),
		registryTestUtils.CreateNode(200, 101, testutils.RandomPrivateKey(t)),
	})
	generator := payerreport.NewPayerReportGenerator(
		testutils.NewLog(t),
		db.Query(),
		registry,
		domainSeparator,
	)

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
			SequenceID:        sequenceID,
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
	db, generator := setupGenerator(t)

	payerAddress := testutils.RandomAddress()
	originatorID := int32(100)

	// Two envelopes in the first minute since the epoch
	addEnvelope(t, db.DB(), originatorID, 1, payerAddress, getMinute(1))
	addEnvelope(t, db.DB(), originatorID, 2, payerAddress, getMinute(1))

	// One envelope in the second minute since the epoch
	addEnvelope(t, db.DB(), originatorID, 3, payerAddress, getMinute(2))

	report, err := generator.GenerateReport(
		context.Background(),
		payerreport.PayerReportGenerationParams{
			OriginatorID:            uint32(originatorID),
			LastReportEndSequenceID: 0,
		},
	)
	require.NoError(t, err)

	require.Equal(t, uint32(originatorID), report.OriginatorNodeID)
	// Make sure the end sequence ID is the last sequence ID from the previous minute
	require.Equal(t, uint64(2), report.EndSequenceID)
	require.Equal(t, currency.PicoDollar(200), report.Payers[payerAddress])
}

func TestReportWithMultiplePayers(t *testing.T) {
	db, generator := setupGenerator(t)

	payerAddress1 := testutils.RandomAddress()
	payerAddress2 := testutils.RandomAddress()
	originatorID := int32(100)

	addEnvelope(t, db.DB(), originatorID, 1, payerAddress1, getMinute(1))
	addEnvelope(t, db.DB(), originatorID, 2, payerAddress2, getMinute(1))
	addEnvelope(t, db.DB(), originatorID, 3, payerAddress1, getMinute(2))

	report, err := generator.GenerateReport(
		context.Background(),
		payerreport.PayerReportGenerationParams{
			OriginatorID:            uint32(originatorID),
			LastReportEndSequenceID: 0,
		},
	)
	require.NoError(t, err)

	require.Equal(t, uint32(originatorID), report.OriginatorNodeID)
	require.Equal(t, uint64(2), report.EndSequenceID)
	require.Equal(t, currency.PicoDollar(100), report.Payers[payerAddress1])
	require.Equal(t, currency.PicoDollar(100), report.Payers[payerAddress2])
}

func TestReportWithNoMessages(t *testing.T) {
	_, generator := setupGenerator(t)

	originatorID := int32(100)

	report, err := generator.GenerateReport(
		context.Background(),
		payerreport.PayerReportGenerationParams{
			OriginatorID:            uint32(originatorID),
			LastReportEndSequenceID: 0,
		},
	)
	require.NoError(t, err)
	require.Nil(t, report)
}

func TestSecondReport(t *testing.T) {
	db, generator := setupGenerator(t)

	originatorID := int32(100)
	payerAddress := testutils.RandomAddress()

	addEnvelope(t, db.DB(), originatorID, 1, payerAddress, getMinute(1))
	addEnvelope(t, db.DB(), originatorID, 2, payerAddress, getMinute(1))
	addEnvelope(t, db.DB(), originatorID, 3, payerAddress, getMinute(2))

	report, err := generator.GenerateReport(
		context.Background(),
		payerreport.PayerReportGenerationParams{
			OriginatorID:            uint32(originatorID),
			LastReportEndSequenceID: 0,
		},
	)
	require.NoError(t, err)

	require.Equal(t, uint32(originatorID), report.OriginatorNodeID)
	require.Equal(t, uint64(0), report.StartSequenceID)
	require.Equal(t, uint64(2), report.EndSequenceID)
	require.Equal(t, currency.PicoDollar(200), report.Payers[payerAddress])

	addEnvelope(t, db.DB(), originatorID, 4, payerAddress, getMinute(3))
	addEnvelope(t, db.DB(), originatorID, 5, payerAddress, getMinute(4))

	report, err = generator.GenerateReport(
		context.Background(),
		payerreport.PayerReportGenerationParams{
			OriginatorID:            uint32(originatorID),
			LastReportEndSequenceID: report.EndSequenceID,
		},
	)
	require.NoError(t, err)

	require.Equal(t, uint32(originatorID), report.OriginatorNodeID)
	require.Equal(t, uint64(4), report.EndSequenceID)
	require.Equal(t, currency.PicoDollar(200), report.Payers[payerAddress])
}

// TestMigratorOriginatorReport simulates the first report for each migrator originator.
func TestMigratorOriginatorReport(t *testing.T) {
	testCases := []struct {
		name         string
		originatorID uint32
	}{
		{"group_messages", migrator.GroupMessageOriginatorID},
		{"welcome_messages", migrator.WelcomeMessageOriginatorID},
		{"key_packages", migrator.KeyPackagesOriginatorID},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var (
				db, generator = setupGenerator(t)
				payerAddress  = testutils.RandomAddress()
				originatorID  = int32(tc.originatorID)
			)

			addEnvelope(t, db.DB(), originatorID, 1, payerAddress, getMinute(1))
			addEnvelope(t, db.DB(), originatorID, 2, payerAddress, getMinute(1))
			addEnvelope(t, db.DB(), originatorID, 3, payerAddress, getMinute(2))

			report, err := generator.GenerateReport(
				t.Context(),
				payerreport.PayerReportGenerationParams{
					OriginatorID:            tc.originatorID,
					LastReportEndSequenceID: 0,
				},
			)
			require.NoError(t, err)
			require.NotNil(t, report)

			require.Equal(t, tc.originatorID, report.OriginatorNodeID)
			require.Equal(t, uint64(2), report.EndSequenceID)
			require.Equal(t, currency.PicoDollar(200), report.Payers[payerAddress])
		})
	}
}

// TestMigratorOriginatorSequentialReports verifies that sequential reports
// chain correctly for migrator originator IDs, with each report's start
// picking up where the previous one ended.
func TestMigratorOriginatorSequentialReports(t *testing.T) {
	var (
		db, generator = setupGenerator(t)
		originatorID  = int32(migrator.GroupMessageOriginatorID)
		payerAddress  = testutils.RandomAddress()
	)

	addEnvelope(t, db.DB(), originatorID, 1, payerAddress, getMinute(1))
	addEnvelope(t, db.DB(), originatorID, 2, payerAddress, getMinute(1))
	addEnvelope(t, db.DB(), originatorID, 3, payerAddress, getMinute(2))

	// First report.
	report1, err := generator.GenerateReport(
		t.Context(),
		payerreport.PayerReportGenerationParams{
			OriginatorID:            migrator.GroupMessageOriginatorID,
			LastReportEndSequenceID: 0,
		},
	)
	require.NoError(t, err)
	require.NotNil(t, report1)

	require.Equal(t, migrator.GroupMessageOriginatorID, report1.OriginatorNodeID)
	require.Equal(t, uint64(0), report1.StartSequenceID)
	require.Equal(t, uint64(2), report1.EndSequenceID)
	require.Equal(t, currency.PicoDollar(200), report1.Payers[payerAddress])

	addEnvelope(t, db.DB(), originatorID, 4, payerAddress, getMinute(3))
	addEnvelope(t, db.DB(), originatorID, 5, payerAddress, getMinute(4))

	// Second report.
	report2, err := generator.GenerateReport(
		t.Context(),
		payerreport.PayerReportGenerationParams{
			OriginatorID:            migrator.GroupMessageOriginatorID,
			LastReportEndSequenceID: report1.EndSequenceID,
		},
	)
	require.NoError(t, err)
	require.NotNil(t, report2)

	require.Equal(t, migrator.GroupMessageOriginatorID, report2.OriginatorNodeID)
	require.Equal(t, uint64(4), report2.EndSequenceID)
	require.Equal(t, currency.PicoDollar(200), report2.Payers[payerAddress])
}

// TestMigratorReportIDDeterminism verifies that the report ID is deterministic
// for migrator originators — same input always produces the same ID.
func TestMigratorReportIDDeterminism(t *testing.T) {
	db, generator := setupGenerator(t)

	originatorID := int32(migrator.GroupMessageOriginatorID)
	payerAddress := testutils.RandomAddress()

	addEnvelope(t, db.DB(), originatorID, 1, payerAddress, getMinute(1))
	addEnvelope(t, db.DB(), originatorID, 2, payerAddress, getMinute(2))

	report1, err := generator.GenerateReport(
		t.Context(),
		payerreport.PayerReportGenerationParams{
			OriginatorID:            migrator.GroupMessageOriginatorID,
			LastReportEndSequenceID: 0,
		},
	)
	require.NoError(t, err)
	require.NotNil(t, report1)

	report2, err := generator.GenerateReport(
		t.Context(),
		payerreport.PayerReportGenerationParams{
			OriginatorID:            migrator.GroupMessageOriginatorID,
			LastReportEndSequenceID: 0,
		},
	)
	require.NoError(t, err)
	require.NotNil(t, report2)

	require.Equal(
		t,
		report1.ID,
		report2.ID,
		"report IDs must be deterministic for cross-node attestation",
	)
	require.Equal(t, report1.PayersMerkleRoot, report2.PayersMerkleRoot)
}
