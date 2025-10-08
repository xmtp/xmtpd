package payerreport_test

import (
	"context"
	"math"
	"sync"
	"testing"
	"time"

	"github.com/xmtp/xmtpd/pkg/payerreport"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	envsWrapper "github.com/xmtp/xmtpd/pkg/envelopes"
	"github.com/xmtp/xmtpd/pkg/proto/identity/associations"
	envelopesProto "github.com/xmtp/xmtpd/pkg/proto/xmtpv4/envelopes"
	"github.com/xmtp/xmtpd/pkg/testutils"
	envTestUtils "github.com/xmtp/xmtpd/pkg/testutils/envelopes"
	"github.com/xmtp/xmtpd/pkg/topic"
)

func createTestStore(t *testing.T) *payerreport.Store {
	log := testutils.NewLog(t)
	db, _ := testutils.NewDB(t, context.Background())

	return payerreport.NewStore(db, log)
}

func insertRandomReport(
	t *testing.T,
	store *payerreport.Store,
) *payerreport.PayerReportWithStatus {
	startID := testutils.RandomInt64()
	reportID := payerreport.ReportID(randomBytes32())
	numRows, err := store.StoreReport(t.Context(), &payerreport.PayerReport{
		ID:               reportID,
		OriginatorNodeID: uint32(testutils.RandomInt32()),
		StartSequenceID:  uint64(startID),
		EndSequenceID:    uint64(startID + 10),
		PayersMerkleRoot: [32]byte(testutils.RandomBytes(32)),
		ActiveNodeIDs:    []uint32{uint32(testutils.RandomInt32())},
	})
	require.NoError(t, err)
	require.Equal(t, int64(1), numRows)
	require.NotNil(t, reportID)

	returnedVal, err := store.FetchReport(t.Context(), reportID)
	require.NoError(t, err)
	return returnedVal
}

// Helper to create a ClientEnvelope containing a PayerReport payload
func createPayerReportClientEnvelope(
	report *payerreport.PayerReport,
) *envelopesProto.ClientEnvelope {
	protoReport := report.ToProto()
	return &envelopesProto.ClientEnvelope{
		Aad: &envelopesProto.AuthenticatedData{
			TargetTopic: topic.NewTopic(topic.TopicKindGroupMessagesV1, testutils.RandomBytes(3)).
				Bytes(),
		},
		Payload: &envelopesProto.ClientEnvelope_PayerReport{
			PayerReport: protoReport,
		},
	}
}

// Helper to create a ClientEnvelope containing a PayerReportAttestation payload
func createPayerReportAttestationClientEnvelope(
	reportID payerreport.ReportID,
	nodeID uint32,
	sig []byte,
) *envelopesProto.ClientEnvelope {
	return &envelopesProto.ClientEnvelope{
		Aad: &envelopesProto.AuthenticatedData{
			TargetTopic: topic.NewTopic(topic.TopicKindGroupMessagesV1, testutils.RandomBytes(3)).
				Bytes(),
		},
		Payload: &envelopesProto.ClientEnvelope_PayerReportAttestation{
			PayerReportAttestation: &envelopesProto.PayerReportAttestation{
				ReportId: reportID[:],
				Signature: &envelopesProto.NodeSignature{
					NodeId:    nodeID,
					Signature: &associations.RecoverableEcdsaSignature{Bytes: sig},
				},
			},
		},
	}
}

func TestStoreAndRetrieve(t *testing.T) {
	cases := []struct {
		name      string
		report    payerreport.PayerReport
		expectErr bool
	}{
		{
			name: "valid report",
			report: payerreport.PayerReport{
				ID:                  payerreport.ReportID(randomBytes32()),
				OriginatorNodeID:    1,
				StartSequenceID:     0,
				EndSequenceID:       2,
				EndMinuteSinceEpoch: 1,
				PayersMerkleRoot:    randomBytes32(),
				ActiveNodeIDs:       []uint32{1},
			},
			expectErr: false,
		},
		{
			name: "invalid node ID",
			report: payerreport.PayerReport{
				ID:                  payerreport.ReportID(randomBytes32()),
				OriginatorNodeID:    uint32(math.MaxInt32) + 10,
				StartSequenceID:     0,
				EndSequenceID:       2,
				EndMinuteSinceEpoch: 1,
				PayersMerkleRoot:    randomBytes32(),
				ActiveNodeIDs:       []uint32{1},
			},
			expectErr: true,
		},
		{
			name: "missing node IDs",
			report: payerreport.PayerReport{
				ID:                  payerreport.ReportID(randomBytes32()),
				OriginatorNodeID:    1,
				StartSequenceID:     0,
				EndSequenceID:       2,
				EndMinuteSinceEpoch: 1,
			},
			expectErr: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			store := createTestStore(t)
			numRows, err := store.StoreReport(context.Background(), &c.report)
			if c.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, int64(1), numRows)
				storedReport, err := store.FetchReport(context.Background(), c.report.ID)
				require.NoError(t, err)
				require.Equal(t, c.report, storedReport.PayerReport)
			}
		})
	}
}

func TestIdempotentStore(t *testing.T) {
	store := createTestStore(t)
	report, err := payerreport.BuildPayerReport(payerreport.BuildPayerReportParams{
		OriginatorNodeID: 1,
		StartSequenceID:  0,
		EndSequenceID:    2,
		NodeIDs:          []uint32{1},
		DomainSeparator:  domainSeparator,
	})
	require.NoError(t, err)
	require.Len(t, report.ID, 32)

	numRows1, err := store.StoreReport(context.Background(), &report.PayerReport)
	require.NoError(t, err)
	require.Equal(t, int64(1), numRows1)

	numRows2, err := store.StoreReport(context.Background(), &report.PayerReport)
	require.NoError(t, err)
	require.Equal(t, int64(0), numRows2)

	storedReports, err := store.FetchReports(
		context.Background(),
		payerreport.NewFetchReportsQuery().WithOriginatorNodeID(report.OriginatorNodeID),
	)
	require.NoError(t, err)
	require.Len(t, storedReports, 1)
	require.Equal(t, report.ID, storedReports[0].ID)
}

func TestFetchReport(t *testing.T) {
	store := createTestStore(t)
	report1 := insertRandomReport(t, store)
	time.Sleep(1 * time.Millisecond)
	report2 := insertRandomReport(t, store)
	// Set the second report's status to Approved
	attestation := &payerreport.PayerReportAttestation{
		Report: &report2.PayerReport,
		NodeSignature: payerreport.NodeSignature{
			NodeID:    2,
			Signature: []byte("sig"),
		},
	}
	payerProto := envTestUtils.CreatePayerEnvelope(t, report2.OriginatorNodeID)
	payerEnv, err := envsWrapper.NewPayerEnvelope(payerProto)
	require.NoError(t, err)
	require.NoError(t, store.CreateAttestation(t.Context(), attestation, payerEnv))

	report3 := insertRandomReport(t, store)

	cases := []struct {
		name        string
		expectedIDs []payerreport.ReportID
		query       *payerreport.FetchReportsQuery
	}{{
		name:        "Get all with created after",
		expectedIDs: []payerreport.ReportID{report1.ID, report2.ID, report3.ID},

		query: payerreport.NewFetchReportsQuery().
			WithCreatedAfter(report1.CreatedAt.Add(-5 * time.Second)),
	}, {
		name:        "Get newest 2",
		expectedIDs: []payerreport.ReportID{report2.ID, report3.ID},
		query:       payerreport.NewFetchReportsQuery().WithCreatedAfter(report1.CreatedAt),
	}, {
		name:        "Only approved",
		expectedIDs: []payerreport.ReportID{report2.ID},
		query: payerreport.NewFetchReportsQuery().WithCreatedAfter(time.Unix(1, 0)).
			WithAttestationStatus(payerreport.AttestationApproved),
	}, {
		name:        "Multiple statuses",
		expectedIDs: []payerreport.ReportID{report2.ID},
		query: payerreport.NewFetchReportsQuery().
			WithAttestationStatus(payerreport.AttestationApproved, payerreport.AttestationRejected),
	}, {
		name:        "No results",
		expectedIDs: []payerreport.ReportID{},
		query: payerreport.NewFetchReportsQuery().WithCreatedAfter(time.Unix(1, 0)).
			WithAttestationStatus(payerreport.AttestationRejected),
	}, {
		name:        "No Params",
		expectedIDs: []payerreport.ReportID{report1.ID, report2.ID, report3.ID},
		query:       payerreport.NewFetchReportsQuery(),
	}, {
		name:        "With start sequence ID",
		expectedIDs: []payerreport.ReportID{report1.ID},
		query: payerreport.NewFetchReportsQuery().
			WithStartSequenceID(report1.StartSequenceID),
	}, {
		name:        "With end sequence ID",
		expectedIDs: []payerreport.ReportID{report1.ID},
		query:       payerreport.NewFetchReportsQuery().WithEndSequenceID(report1.EndSequenceID),
	}, {
		name:        "With start and end sequence ID",
		expectedIDs: []payerreport.ReportID{report1.ID},
		query: payerreport.NewFetchReportsQuery().WithStartSequenceID(report1.StartSequenceID).
			WithEndSequenceID(report1.EndSequenceID),
	}, {
		name:        "With originator node ID",
		expectedIDs: []payerreport.ReportID{report1.ID},
		query: payerreport.NewFetchReportsQuery().
			WithOriginatorNodeID(report1.OriginatorNodeID),
	}, {
		name:        "With min attestations",
		expectedIDs: []payerreport.ReportID{report2.ID},
		query:       payerreport.NewFetchReportsQuery().WithMinAttestations(1),
	}}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			results, err := store.FetchReports(t.Context(), c.query)
			require.NoError(t, err)
			require.Len(t, results, len(c.expectedIDs))

			returnedIDs := make([]payerreport.ReportID, len(results))
			for idx, result := range results {
				returnedIDs[idx] = result.ID
			}

			require.ElementsMatch(t, c.expectedIDs, returnedIDs)
		})
	}
}

func TestStoreAttestation(t *testing.T) {
	store := createTestStore(t)
	ctx := context.Background()
	reportID := payerreport.ReportID(randomBytes32())

	report := payerreport.PayerReport{
		ID:               reportID,
		OriginatorNodeID: 1,
		StartSequenceID:  0,
		EndSequenceID:    10,
		PayersMerkleRoot: randomBytes32(),
		ActiveNodeIDs:    []uint32{1},
	}

	// First, store the report so that the attestation can reference it
	numRows, err := store.StoreReport(ctx, &report)
	require.NoError(t, err)
	require.Equal(t, int64(1), numRows)

	attestation := &payerreport.PayerReportAttestation{
		Report: &report,
		NodeSignature: payerreport.NodeSignature{
			NodeID:    2,
			Signature: []byte("sig"),
		},
	}

	require.NoError(t, store.StoreAttestation(ctx, attestation))

	// Verify we can fetch the attestation we just stored
	fetchedReport, err := store.FetchReport(ctx, reportID)
	require.NoError(t, err)
	require.Len(t, fetchedReport.AttestationSignatures, 1)
	require.Equal(
		t,
		fetchedReport.AttestationSignatures[0].NodeID,
		fetchedReport.AttestationSignatures[0].NodeID,
	)
}

func TestStoreAttestationInvalidNodeID(t *testing.T) {
	store := createTestStore(t)
	ctx := context.Background()
	reportID := payerreport.ReportID(randomBytes32())

	report := payerreport.PayerReport{
		ID:               reportID,
		OriginatorNodeID: 1,
		StartSequenceID:  0,
		EndSequenceID:    10,
		PayersMerkleRoot: randomBytes32(),
		ActiveNodeIDs:    []uint32{1},
	}

	numRows, err := store.StoreReport(ctx, &report)
	require.NoError(t, err)
	require.Equal(t, int64(1), numRows)

	attestation := &payerreport.PayerReportAttestation{
		Report: &report,
		NodeSignature: payerreport.NodeSignature{
			NodeID:    uint32(math.MaxInt32) + 1,
			Signature: []byte("sig"),
		},
	}

	err = store.StoreAttestation(ctx, attestation)
	require.ErrorIs(t, err, payerreport.ErrOriginatorNodeIDTooLarge)
}

func TestSetReportAttestationStatus(t *testing.T) {
	store := createTestStore(t)
	ctx := context.Background()

	desiredStatuses := []payerreport.AttestationStatus{
		payerreport.AttestationApproved,
		payerreport.AttestationRejected,
	}

	for _, newStatus := range desiredStatuses {
		reportID := payerreport.ReportID(randomBytes32())

		report := payerreport.PayerReport{
			ID:               reportID,
			OriginatorNodeID: 1,
			StartSequenceID:  0,
			EndSequenceID:    10,
			PayersMerkleRoot: randomBytes32(),
			ActiveNodeIDs:    []uint32{1},
		}

		numRows, err := store.StoreReport(ctx, &report)
		require.NoError(t, err)
		require.Equal(t, int64(1), numRows)

		if newStatus == payerreport.AttestationApproved {
			require.NoError(t, store.SetReportAttestationApproved(ctx, reportID))
		} else {
			require.NoError(t, store.SetReportAttestationRejected(ctx, reportID))
		}

		fetched, err := store.FetchReport(ctx, reportID)
		require.NoError(t, err)
		require.Equal(t, newStatus, fetched.AttestationStatus)
	}
}

func TestInvalidStateTransition(t *testing.T) {
	store := createTestStore(t)
	ctx := context.Background()

	reportID := payerreport.ReportID(randomBytes32())
	report := payerreport.PayerReport{
		ID:               reportID,
		OriginatorNodeID: 1,
		StartSequenceID:  0,
		EndSequenceID:    10,
		PayersMerkleRoot: randomBytes32(),
		ActiveNodeIDs:    []uint32{1},
	}

	numRows, err := store.StoreReport(ctx, &report)
	require.NoError(t, err)
	require.Equal(t, int64(1), numRows)

	err = store.SetReportAttestationApproved(ctx, reportID)
	require.NoError(t, err)

	require.NoError(t, store.SetReportAttestationRejected(ctx, reportID))

	fetched, err := store.FetchReport(ctx, reportID)
	require.NoError(t, err)
	require.Equal(
		t,
		payerreport.AttestationStatus(payerreport.AttestationApproved),
		fetched.AttestationStatus,
	)
}

func TestCreatePayerReport(t *testing.T) {
	store := createTestStore(t)
	ctx := context.Background()

	report, err := payerreport.BuildPayerReport(payerreport.BuildPayerReportParams{
		OriginatorNodeID: 3,
		StartSequenceID:  0,
		EndSequenceID:    10,
		NodeIDs:          []uint32{3},
		DomainSeparator:  domainSeparator,
	})
	require.NoError(t, err)
	// Build a minimal payer envelope (group message payload is fine for this path)
	payerProto := envTestUtils.CreatePayerEnvelope(t, report.OriginatorNodeID)
	payerEnv, err := envsWrapper.NewPayerEnvelope(payerProto)
	require.NoError(t, err)

	reportID, err := store.CreatePayerReport(ctx, &report.PayerReport, payerEnv)
	require.NoError(t, err)
	require.NotNil(t, reportID)
	require.Equal(t, *reportID, report.ID)

	fetched, err := store.FetchReport(ctx, *reportID)
	require.NoError(t, err)
	require.Equal(t, report.OriginatorNodeID, fetched.OriginatorNodeID)
	require.Equal(t, *reportID, fetched.ID)

	// Ensure a staged originator envelope was created
	staged, err := store.Queries().
		SelectStagedOriginatorEnvelopes(ctx, queries.SelectStagedOriginatorEnvelopesParams{LastSeenID: 0, NumRows: 10})
	require.NoError(t, err)
	require.Len(t, staged, 1)
}

func TestCreateAttestation(t *testing.T) {
	store := createTestStore(t)
	ctx := context.Background()
	reportID := payerreport.ReportID(randomBytes32())

	report := payerreport.PayerReport{
		ID:               reportID,
		OriginatorNodeID: 4,
		StartSequenceID:  0,
		EndSequenceID:    10,
		PayersMerkleRoot: randomBytes32(),
		ActiveNodeIDs:    []uint32{4},
	}

	numRows, err := store.StoreReport(ctx, &report)
	require.NoError(t, err)
	require.Equal(t, int64(1), numRows)

	attestation := &payerreport.PayerReportAttestation{
		Report: &report,
		NodeSignature: payerreport.NodeSignature{
			NodeID:    5,
			Signature: []byte("sig"),
		},
	}

	payerProto := envTestUtils.CreatePayerEnvelope(t, report.OriginatorNodeID)
	payerEnv, err := envsWrapper.NewPayerEnvelope(payerProto)
	require.NoError(t, err)

	require.NoError(t, store.CreateAttestation(ctx, attestation, payerEnv))

	fetchedReport, err := store.FetchReport(ctx, reportID)
	require.NoError(t, err)
	require.Equal(
		t,
		payerreport.AttestationStatus(payerreport.AttestationApproved),
		fetchedReport.AttestationStatus,
	)
	require.Len(t, fetchedReport.AttestationSignatures, 1)
	require.Equal(
		t,
		attestation.NodeSignature.NodeID,
		fetchedReport.AttestationSignatures[0].NodeID,
	)

	staged, err := store.Queries().
		SelectStagedOriginatorEnvelopes(ctx, queries.SelectStagedOriginatorEnvelopesParams{LastSeenID: 0, NumRows: 10})
	require.NoError(t, err)
	require.Len(t, staged, 1)
}

func TestStoreSyncedReport(t *testing.T) {
	store := createTestStore(t)
	ctx := context.Background()

	report, err := payerreport.BuildPayerReport(payerreport.BuildPayerReportParams{
		OriginatorNodeID:    7,
		StartSequenceID:     0,
		EndSequenceID:       10,
		EndMinuteSinceEpoch: uint32(time.Now().Unix() / 60),
		NodeIDs:             []uint32{7},
		DomainSeparator:     domainSeparator,
	})
	require.NoError(t, err)
	t.Logf("report: %+v", report.ID)

	// Build the originator envelope containing the payer report payload
	clientEnv := createPayerReportClientEnvelope(&report.PayerReport)
	payerProto := envTestUtils.CreatePayerEnvelope(t, report.OriginatorNodeID, clientEnv)
	originatorProto := envTestUtils.CreateOriginatorEnvelope(
		t,
		report.OriginatorNodeID,
		1,
		payerProto,
	)

	originatorEnv, err := envsWrapper.NewOriginatorEnvelope(originatorProto)
	require.NoError(t, err)

	payerID, err := store.Queries().FindOrCreatePayer(ctx, testutils.RandomAddress().Hex())
	require.NoError(t, err)

	require.NoError(t, store.StoreSyncedReport(ctx, originatorEnv, payerID, domainSeparator))

	fetched, err := store.FetchReport(ctx, report.ID)
	require.NoError(t, err)
	require.Equal(t, report.OriginatorNodeID, fetched.OriginatorNodeID)
}

func TestStoreSyncedAttestation(t *testing.T) {
	store := createTestStore(t)
	ctx := context.Background()
	reportID := payerreport.ReportID(randomBytes32())

	// First create and store the base report so that the attestation references a real report ID
	baseReport := payerreport.PayerReport{
		ID:               reportID,
		OriginatorNodeID: 8,
		StartSequenceID:  0,
		EndSequenceID:    10,
		PayersMerkleRoot: randomBytes32(),
		ActiveNodeIDs:    []uint32{8},
	}
	numRows, err := store.StoreReport(ctx, &baseReport)
	require.NoError(t, err)
	require.Equal(t, int64(1), numRows)

	// Build attestation envelope
	sigBytes := []byte("sig")
	clientEnv := createPayerReportAttestationClientEnvelope(reportID, 9, sigBytes)
	payerProto := envTestUtils.CreatePayerEnvelope(t, baseReport.OriginatorNodeID, clientEnv)
	originatorProto := envTestUtils.CreateOriginatorEnvelope(
		t,
		baseReport.OriginatorNodeID,
		2,
		payerProto,
	)

	originatorEnv, err := envsWrapper.NewOriginatorEnvelope(originatorProto)
	require.NoError(t, err)
	payerID, err := store.Queries().FindOrCreatePayer(ctx, testutils.RandomAddress().Hex())
	require.NoError(t, err)
	require.NoError(t, store.StoreSyncedAttestation(ctx, originatorEnv, payerID))

	fetchedReport, err := store.FetchReport(ctx, reportID)
	require.NoError(t, err)
	require.Len(t, fetchedReport.AttestationSignatures, 1)
	require.Equal(t, uint32(9), fetchedReport.AttestationSignatures[0].NodeID)
}

func TestSetReportSettled(t *testing.T) {
	ctx := context.Background()
	store := createTestStore(t)

	// Helper to create unsettled usage for an originator
	// Every increment is 100 picodollars
	createUnsettledUsage := func(originatorID uint32, minutesSinceEpoch int32, payerCount int) {
		for i := 0; i < payerCount; i++ {
			payerID, err := store.Queries().FindOrCreatePayer(ctx, testutils.RandomAddress().Hex())
			require.NoError(t, err)

			err = store.Queries().
				IncrementUnsettledUsage(ctx, queries.IncrementUnsettledUsageParams{
					PayerID:           payerID,
					OriginatorID:      int32(originatorID),
					MinutesSinceEpoch: minutesSinceEpoch,
					SpendPicodollars:  100,
					MessageCount:      1,
					SequenceID:        int64(i),
				})
			require.NoError(t, err)
		}
	}

	// Helper to count unsettled usage for an originator
	countUnsettledUsage := func(originatorID uint32) int {
		rows, err := store.Queries().BuildPayerReport(ctx, queries.BuildPayerReportParams{
			OriginatorID:           int32(originatorID),
			StartMinutesSinceEpoch: 0,
			EndMinutesSinceEpoch:   1000000, // Large value to include all
		})
		require.NoError(t, err)
		count := 0
		for _, row := range rows {
			if row.TotalSpendPicodollars > 0 {
				count += int(row.TotalSpendPicodollars / 100) // Each entry has 100 picodollars
			}
		}
		return count
	}

	t.Run("first settled report clears all unsettled usage", func(t *testing.T) {
		originatorID := uint32(99)

		// Create the first report for this originator
		report := &payerreport.PayerReport{
			ID:                  payerreport.ReportID(randomBytes32()),
			OriginatorNodeID:    originatorID,
			StartSequenceID:     0,
			EndSequenceID:       10,
			EndMinuteSinceEpoch: 500,
			PayersMerkleRoot:    randomBytes32(),
			ActiveNodeIDs:       []uint32{1, 2, 3},
		}
		numRows, err := store.StoreReport(ctx, report)
		require.NoError(t, err)
		require.Equal(t, int64(1), numRows)

		// Create unsettled usage at various times
		createUnsettledUsage(originatorID, 100, 2) // Very old usage
		createUnsettledUsage(originatorID, 300, 2) // Before end minute
		createUnsettledUsage(originatorID, 450, 3) // Before end minute
		createUnsettledUsage(originatorID, 500, 2) // At end minute
		createUnsettledUsage(originatorID, 600, 2) // After end minute (should not be cleared)

		// Verify all unsettled usage exists
		require.Equal(t, 11, countUnsettledUsage(originatorID))

		// Set to submitted then settled
		err = store.SetReportSubmitted(ctx, report.ID)
		require.NoError(t, err)
		err = store.SetReportSettled(ctx, report.ID)
		require.NoError(t, err)

		// Since this is the first settled report, all usage up to and including
		// the end minute should be cleared
		require.Equal(t, 2, countUnsettledUsage(originatorID))

		// Verify report status changed to settled
		fetchedReport, err := store.FetchReport(ctx, report.ID)
		require.NoError(t, err)
		require.Equal(
			t,
			payerreport.SubmissionStatus(payerreport.SubmissionSettled),
			fetchedReport.SubmissionStatus,
		)
	})

	t.Run("clears unsettled usage only for range since last settled report", func(t *testing.T) {
		originatorID := uint32(100)

		// Create first report and settle it
		firstReport := &payerreport.PayerReport{
			ID:                  payerreport.ReportID(randomBytes32()),
			OriginatorNodeID:    originatorID,
			StartSequenceID:     0,
			EndSequenceID:       10,
			EndMinuteSinceEpoch: 300,
			PayersMerkleRoot:    randomBytes32(),
			ActiveNodeIDs:       []uint32{1, 2, 3},
		}
		numRows, err := store.StoreReport(ctx, firstReport)
		require.NoError(t, err)
		require.Equal(t, int64(1), numRows)
		err = store.SetReportSubmitted(ctx, firstReport.ID)
		require.NoError(t, err)
		err = store.SetReportSettled(ctx, firstReport.ID)
		require.NoError(t, err)

		// Create second report
		secondReport := &payerreport.PayerReport{
			ID:                  payerreport.ReportID(randomBytes32()),
			OriginatorNodeID:    originatorID,
			StartSequenceID:     10,
			EndSequenceID:       20,
			EndMinuteSinceEpoch: 500,
			PayersMerkleRoot:    randomBytes32(),
			ActiveNodeIDs:       []uint32{1, 2, 3},
		}
		numRows, err = store.StoreReport(ctx, secondReport)
		require.Equal(t, int64(1), numRows)
		require.NoError(t, err)

		// Create unsettled usage across different time ranges
		createUnsettledUsage(originatorID, 250, 2) // Before first report (should not be cleared)
		createUnsettledUsage(originatorID, 350, 3) // Between reports (should be cleared)
		createUnsettledUsage(originatorID, 450, 2) // Between reports (should be cleared)
		createUnsettledUsage(originatorID, 500, 2) // At second report end (should be cleared)
		createUnsettledUsage(originatorID, 600, 2) // After second report (should not be cleared)

		// Verify all unsettled usage exists
		require.Equal(t, 11, countUnsettledUsage(originatorID))

		// Set second report to submitted then settled
		err = store.SetReportSubmitted(ctx, secondReport.ID)
		require.NoError(t, err)
		err = store.SetReportSettled(ctx, secondReport.ID)
		require.NoError(t, err)

		// Check that only usage outside the range (300, 500] remains
		// Should have 2 from minute 250 and 2 from minute 600
		require.Equal(t, 4, countUnsettledUsage(originatorID))

		// Verify report status changed to settled
		fetchedReport, err := store.FetchReport(ctx, secondReport.ID)
		require.NoError(t, err)
		require.Equal(
			t,
			payerreport.SubmissionStatus(payerreport.SubmissionSettled),
			fetchedReport.SubmissionStatus,
		)
	})

	t.Run("handles multiple reports with submitted and settled states", func(t *testing.T) {
		originatorID := uint32(101)

		// Create and settle first report
		report1 := &payerreport.PayerReport{
			ID:                  payerreport.ReportID(randomBytes32()),
			OriginatorNodeID:    originatorID,
			StartSequenceID:     0,
			EndSequenceID:       10,
			EndMinuteSinceEpoch: 200,
			PayersMerkleRoot:    randomBytes32(),
			ActiveNodeIDs:       []uint32{1, 2},
		}
		numRows, err := store.StoreReport(ctx, report1)
		require.NoError(t, err)
		require.Equal(t, int64(1), numRows)
		err = store.SetReportSubmitted(ctx, report1.ID)
		require.NoError(t, err)
		err = store.SetReportSettled(ctx, report1.ID)
		require.NoError(t, err)

		// Create and submit (but not settle) second report
		report2 := &payerreport.PayerReport{
			ID:                  payerreport.ReportID(randomBytes32()),
			OriginatorNodeID:    originatorID,
			StartSequenceID:     10,
			EndSequenceID:       20,
			EndMinuteSinceEpoch: 400,
			PayersMerkleRoot:    randomBytes32(),
			ActiveNodeIDs:       []uint32{1, 2},
		}
		numRows, err = store.StoreReport(ctx, report2)
		require.NoError(t, err)
		require.Equal(t, int64(1), numRows)
		err = store.SetReportSubmitted(ctx, report2.ID)
		require.NoError(t, err)

		// Create third report
		report3 := &payerreport.PayerReport{
			ID:                  payerreport.ReportID(randomBytes32()),
			OriginatorNodeID:    originatorID,
			StartSequenceID:     20,
			EndSequenceID:       30,
			EndMinuteSinceEpoch: 600,
			PayersMerkleRoot:    randomBytes32(),
			ActiveNodeIDs:       []uint32{1, 2},
		}
		numRows, err = store.StoreReport(ctx, report3)
		require.NoError(t, err)
		require.Equal(t, int64(1), numRows)

		// Create unsettled usage across all time ranges
		createUnsettledUsage(originatorID, 150, 2) // Before first report
		createUnsettledUsage(originatorID, 250, 2) // Between first and second
		createUnsettledUsage(originatorID, 350, 2) // Between first and second
		createUnsettledUsage(originatorID, 450, 2) // Between second and third
		createUnsettledUsage(originatorID, 550, 2) // Between second and third
		createUnsettledUsage(originatorID, 600, 2) // At third report end
		createUnsettledUsage(originatorID, 700, 2) // After third report

		// Verify all unsettled usage exists
		require.Equal(t, 14, countUnsettledUsage(originatorID))

		// Settle the third report
		err = store.SetReportSubmitted(ctx, report3.ID)
		require.NoError(t, err)
		err = store.SetReportSettled(ctx, report3.ID)
		require.NoError(t, err)

		// Should clear usage from last settled/submitted report (report2 at minute 400) to report3 end (minute 600)
		// Remaining: 150 (2), 250 (2), 350 (2), 700 (2) = 8 total
		require.Equal(t, 8, countUnsettledUsage(originatorID))
	})

	t.Run("handles gap when last report is settled vs submitted", func(t *testing.T) {
		originatorID := uint32(102)

		// Create and submit (but not settle) first report
		report1 := &payerreport.PayerReport{
			ID:                  payerreport.ReportID(randomBytes32()),
			OriginatorNodeID:    originatorID,
			StartSequenceID:     0,
			EndSequenceID:       10,
			EndMinuteSinceEpoch: 300,
			PayersMerkleRoot:    randomBytes32(),
			ActiveNodeIDs:       []uint32{1, 2},
		}
		numRows, err := store.StoreReport(ctx, report1)
		require.NoError(t, err)
		require.Equal(t, int64(1), numRows)
		err = store.SetReportSubmitted(ctx, report1.ID)
		require.NoError(t, err)

		// Create second report
		report2 := &payerreport.PayerReport{
			ID:                  payerreport.ReportID(randomBytes32()),
			OriginatorNodeID:    originatorID,
			StartSequenceID:     10,
			EndSequenceID:       20,
			EndMinuteSinceEpoch: 500,
			PayersMerkleRoot:    randomBytes32(),
			ActiveNodeIDs:       []uint32{1, 2},
		}
		numRows, err = store.StoreReport(ctx, report2)
		require.NoError(t, err)
		require.Equal(t, int64(1), numRows)

		// Create unsettled usage
		createUnsettledUsage(originatorID, 200, 2) // Before first report
		createUnsettledUsage(originatorID, 350, 3) // After first report
		createUnsettledUsage(originatorID, 450, 2) // Between reports
		createUnsettledUsage(originatorID, 500, 2) // At second report end
		createUnsettledUsage(originatorID, 600, 2) // After second report

		// Verify all unsettled usage exists
		require.Equal(t, 11, countUnsettledUsage(originatorID))

		// Settle the second report
		err = store.SetReportSubmitted(ctx, report2.ID)
		require.NoError(t, err)
		err = store.SetReportSettled(ctx, report2.ID)
		require.NoError(t, err)

		// Should clear usage from last submitted report (report1 at minute 300) to report2 end (minute 500)
		// Remaining: 200 (2), 600 (2) = 4 total
		require.Equal(t, 4, countUnsettledUsage(originatorID))
	})

	t.Run("does not affect other originators", func(t *testing.T) {
		// Create two reports for different originators
		report1 := &payerreport.PayerReport{
			ID:                  payerreport.ReportID(randomBytes32()),
			OriginatorNodeID:    200,
			StartSequenceID:     0,
			EndSequenceID:       10,
			EndMinuteSinceEpoch: 400,
			PayersMerkleRoot:    randomBytes32(),
			ActiveNodeIDs:       []uint32{1, 2},
		}
		report2 := &payerreport.PayerReport{
			ID:                  payerreport.ReportID(randomBytes32()),
			OriginatorNodeID:    201,
			StartSequenceID:     0,
			EndSequenceID:       10,
			EndMinuteSinceEpoch: 400,
			PayersMerkleRoot:    randomBytes32(),
			ActiveNodeIDs:       []uint32{1, 2},
		}

		numRows, err := store.StoreReport(ctx, report1)
		require.NoError(t, err)
		require.Equal(t, int64(1), numRows)

		numRows, err = store.StoreReport(ctx, report2)
		require.NoError(t, err)
		require.Equal(t, int64(1), numRows)

		// Create unsettled usage for both originators
		createUnsettledUsage(report1.OriginatorNodeID, 350, 3)
		createUnsettledUsage(report2.OriginatorNodeID, 350, 4)

		// Verify both have unsettled usage
		require.Equal(t, 3, countUnsettledUsage(report1.OriginatorNodeID))
		require.Equal(t, 4, countUnsettledUsage(report2.OriginatorNodeID))

		// Set first report to submitted then settled
		err = store.SetReportSubmitted(ctx, report1.ID)
		require.NoError(t, err)
		err = store.SetReportSettled(ctx, report1.ID)
		require.NoError(t, err)

		// Check that only report1's usage was cleared
		require.Equal(t, 0, countUnsettledUsage(report1.OriginatorNodeID))
		require.Equal(t, 4, countUnsettledUsage(report2.OriginatorNodeID))
	})

	t.Run("fails for non-existent report", func(t *testing.T) {
		nonExistentID := payerreport.ReportID(randomBytes32())
		err := store.SetReportSettled(ctx, nonExistentID)
		require.Error(t, err)
	})

	t.Run("transitions from pending to settled", func(t *testing.T) {
		// Create a report in pending state
		report := &payerreport.PayerReport{
			ID:                  payerreport.ReportID(randomBytes32()),
			OriginatorNodeID:    300,
			StartSequenceID:     0,
			EndSequenceID:       10,
			EndMinuteSinceEpoch: 300,
			PayersMerkleRoot:    randomBytes32(),
			ActiveNodeIDs:       []uint32{1},
		}
		numRows, err := store.StoreReport(ctx, report)
		require.NoError(t, err)
		require.Equal(t, int64(1), numRows)

		// Verify initial status is pending
		fetchedReport, err := store.FetchReport(ctx, report.ID)
		require.NoError(t, err)
		require.Equal(
			t,
			payerreport.SubmissionStatus(payerreport.SubmissionPending),
			fetchedReport.SubmissionStatus,
		)

		// Should be able to transition from pending to settled
		err = store.SetReportSettled(ctx, report.ID)
		require.NoError(t, err)

		// Verify status changed
		fetchedReport, err = store.FetchReport(ctx, report.ID)
		require.NoError(t, err)
		require.Equal(
			t,
			payerreport.SubmissionStatus(payerreport.SubmissionSettled),
			fetchedReport.SubmissionStatus,
		)
	})

	t.Run("transitions from submitted to settled", func(t *testing.T) {
		// Create a report and set it to submitted
		report := &payerreport.PayerReport{
			ID:                  payerreport.ReportID(randomBytes32()),
			OriginatorNodeID:    400,
			StartSequenceID:     0,
			EndSequenceID:       10,
			EndMinuteSinceEpoch: 300,
			PayersMerkleRoot:    randomBytes32(),
			ActiveNodeIDs:       []uint32{1},
		}
		numRows, err := store.StoreReport(ctx, report)
		require.NoError(t, err)
		require.Equal(t, int64(1), numRows)

		err = store.SetReportSubmitted(ctx, report.ID)
		require.NoError(t, err)

		// Verify status is submitted
		fetchedReport, err := store.FetchReport(ctx, report.ID)
		require.NoError(t, err)
		require.Equal(
			t,
			payerreport.SubmissionStatus(payerreport.SubmissionSubmitted),
			fetchedReport.SubmissionStatus,
		)

		// Should be able to transition from submitted to settled
		err = store.SetReportSettled(ctx, report.ID)
		require.NoError(t, err)

		// Verify status changed
		fetchedReport, err = store.FetchReport(ctx, report.ID)
		require.NoError(t, err)
		require.Equal(
			t,
			payerreport.SubmissionStatus(payerreport.SubmissionSettled),
			fetchedReport.SubmissionStatus,
		)
	})
}

func TestCreateAttestationConcurrency(t *testing.T) {
	store := createTestStore(t)
	ctx := context.Background()
	reportID := payerreport.ReportID(randomBytes32())

	report := payerreport.PayerReport{
		ID:               reportID,
		OriginatorNodeID: 4,
		StartSequenceID:  0,
		EndSequenceID:    10,
		PayersMerkleRoot: randomBytes32(),
		ActiveNodeIDs:    []uint32{4},
	}

	err := store.StoreReport(ctx, &report)
	require.NoError(t, err)

	const numWorkers = 10
	var wg sync.WaitGroup
	errCh := make(chan error, numWorkers)

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			attestation := &payerreport.PayerReportAttestation{
				Report: &report,
				NodeSignature: payerreport.NodeSignature{
					NodeID:    5,
					Signature: []byte("sig"),
				},
			}

			payerProto := envTestUtils.CreatePayerEnvelope(t, report.OriginatorNodeID)
			payerEnv, err := envsWrapper.NewPayerEnvelope(payerProto)
			require.NoError(t, err)

			errCh <- store.CreateAttestation(ctx, attestation, payerEnv)
		}()
	}

	wg.Wait()
	close(errCh)

	for err = range errCh {
		require.NoError(t, err)
	}

	fetchedReport, err := store.FetchReport(ctx, reportID)
	require.NoError(t, err)
	require.Equal(
		t,
		payerreport.AttestationStatus(payerreport.AttestationApproved),
		fetchedReport.AttestationStatus,
	)

	staged, err := store.Queries().
		SelectStagedOriginatorEnvelopes(ctx, queries.SelectStagedOriginatorEnvelopesParams{LastSeenID: 0, NumRows: 10})
	require.NoError(t, err)
	require.Len(t, staged, 1)
}
