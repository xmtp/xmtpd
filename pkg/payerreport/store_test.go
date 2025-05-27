package payerreport

import (
	"context"
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	envsWrapper "github.com/xmtp/xmtpd/pkg/envelopes"
	"github.com/xmtp/xmtpd/pkg/proto/identity/associations"
	envelopesProto "github.com/xmtp/xmtpd/pkg/proto/xmtpv4/envelopes"
	"github.com/xmtp/xmtpd/pkg/testutils"
	envTestUtils "github.com/xmtp/xmtpd/pkg/testutils/envelopes"
	"github.com/xmtp/xmtpd/pkg/topic"
)

func createTestStore(t *testing.T) *Store {
	log := testutils.NewLog(t)
	db, _ := testutils.NewDB(t, context.Background())

	return NewStore(db, log)
}

func insertRandomReport(
	t *testing.T,
	store *Store,
) *PayerReportWithStatus {
	startID := testutils.RandomInt64()
	reportID, err := store.StoreReport(t.Context(), &PayerReport{
		OriginatorNodeID: uint32(testutils.RandomInt32()),
		StartSequenceID:  uint64(startID),
		EndSequenceID:    uint64(startID + 10),
		PayersMerkleRoot: [32]byte(testutils.RandomBytes(32)),
		ActiveNodeIDs:    []uint32{uint32(testutils.RandomInt32())},
	})
	require.NoError(t, err)

	returnedVal, err := store.FetchReport(t.Context(), reportID)
	require.NoError(t, err)
	return returnedVal
}

// Helper to create a ClientEnvelope containing a PayerReport payload
func createPayerReportClientEnvelope(report *PayerReport) *envelopesProto.ClientEnvelope {
	protoReport := report.ToProto()
	return &envelopesProto.ClientEnvelope{
		Aad: &envelopesProto.AuthenticatedData{
			TargetTopic: topic.NewTopic(topic.TOPIC_KIND_GROUP_MESSAGES_V1, testutils.RandomBytes(3)).
				Bytes(),
		},
		Payload: &envelopesProto.ClientEnvelope_PayerReport{
			PayerReport: protoReport,
		},
	}
}

// Helper to create a ClientEnvelope containing a PayerReportAttestation payload
func createPayerReportAttestationClientEnvelope(
	reportID ReportID,
	nodeID uint32,
	sig []byte,
) *envelopesProto.ClientEnvelope {
	return &envelopesProto.ClientEnvelope{
		Aad: &envelopesProto.AuthenticatedData{
			TargetTopic: topic.NewTopic(topic.TOPIC_KIND_GROUP_MESSAGES_V1, testutils.RandomBytes(3)).
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
		report    PayerReport
		expectErr bool
	}{
		{
			name: "valid report",
			report: PayerReport{
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
			report: PayerReport{
				OriginatorNodeID:    uint32(math.MaxInt32) + 1,
				StartSequenceID:     0,
				EndSequenceID:       2,
				EndMinuteSinceEpoch: 1,
				PayersMerkleRoot:    randomBytes32(),
				ActiveNodeIDs:       []uint32{1},
			},
			expectErr: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			store := createTestStore(t)
			id, err := store.StoreReport(context.Background(), &c.report)
			if c.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotNil(t, id)
				require.Len(t, id, 32)
				storedReport, err := store.FetchReport(context.Background(), id)
				require.NoError(t, err)
				require.Equal(t, c.report, storedReport.PayerReport)
			}
		})
	}
}

func TestIdempotentStore(t *testing.T) {
	store := createTestStore(t)
	report := PayerReport{
		OriginatorNodeID: 1,
		StartSequenceID:  0,
		EndSequenceID:    2,
		PayersMerkleRoot: randomBytes32(),
		ActiveNodeIDs:    []uint32{1},
	}
	reportID, err := report.ID()
	require.NoError(t, err)

	returnedID, err := store.StoreReport(context.Background(), &report)
	require.NoError(t, err)
	require.NotNil(t, returnedID)
	require.Len(t, returnedID, 32)
	require.Equal(t, reportID, returnedID)

	newID, err := store.StoreReport(context.Background(), &report)
	require.NoError(t, err)
	require.NotNil(t, newID)
	require.Equal(t, newID, returnedID)
}

func TestFetchReport(t *testing.T) {
	store := createTestStore(t)
	report1 := insertRandomReport(t, store)
	time.Sleep(1 * time.Millisecond)
	report2 := insertRandomReport(t, store)
	// Set the second report's status to Approved
	attestation := &PayerReportAttestation{
		Report: &report2.PayerReport,
		NodeSignature: NodeSignature{
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
		expectedIDs []ReportID
		query       *FetchReportsQuery
	}{{
		name:        "Get all with created after",
		expectedIDs: []ReportID{report1.ID, report2.ID, report3.ID},

		query: NewFetchReportsQuery().WithCreatedAfter(report1.CreatedAt.Add(-5 * time.Second)),
	}, {
		name:        "Get newest 2",
		expectedIDs: []ReportID{report2.ID, report3.ID},
		query:       NewFetchReportsQuery().WithCreatedAfter(report1.CreatedAt),
	}, {
		name:        "Only approved",
		expectedIDs: []ReportID{report2.ID},
		query: NewFetchReportsQuery().WithCreatedAfter(time.Unix(1, 0)).
			WithAttestationStatus(AttestationApproved),
	}, {
		name:        "Multiple statuses",
		expectedIDs: []ReportID{report2.ID},
		query: NewFetchReportsQuery().
			WithAttestationStatus(AttestationApproved, AttestationRejected),
	}, {
		name:        "No results",
		expectedIDs: []ReportID{},
		query: NewFetchReportsQuery().WithCreatedAfter(time.Unix(1, 0)).
			WithAttestationStatus(AttestationRejected),
	}, {
		name:        "No Params",
		expectedIDs: []ReportID{report1.ID, report2.ID, report3.ID},
		query:       NewFetchReportsQuery(),
	}, {
		name:        "With start sequence ID",
		expectedIDs: []ReportID{report1.ID},
		query:       NewFetchReportsQuery().WithStartSequenceID(report1.StartSequenceID),
	}, {
		name:        "With end sequence ID",
		expectedIDs: []ReportID{report1.ID},
		query:       NewFetchReportsQuery().WithEndSequenceID(report1.EndSequenceID),
	}, {
		name:        "With start and end sequence ID",
		expectedIDs: []ReportID{report1.ID},
		query: NewFetchReportsQuery().WithStartSequenceID(report1.StartSequenceID).
			WithEndSequenceID(report1.EndSequenceID),
	}, {
		name:        "With originator node ID",
		expectedIDs: []ReportID{report1.ID},
		query:       NewFetchReportsQuery().WithOriginatorNodeID(report1.OriginatorNodeID),
	}, {
		name:        "With min attestations",
		expectedIDs: []ReportID{report2.ID},
		query:       NewFetchReportsQuery().WithMinAttestations(1),
	}}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			results, err := store.FetchReports(t.Context(), c.query)
			require.NoError(t, err)
			require.Len(t, results, len(c.expectedIDs))

			returnedIDs := make([]ReportID, len(results))
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

	report := PayerReport{
		OriginatorNodeID: 1,
		StartSequenceID:  0,
		EndSequenceID:    10,
		PayersMerkleRoot: randomBytes32(),
		ActiveNodeIDs:    []uint32{1},
	}

	// First, store the report so that the attestation can reference it
	reportID, err := store.StoreReport(ctx, &report)
	require.NoError(t, err)

	attestation := &PayerReportAttestation{
		Report: &report,
		NodeSignature: NodeSignature{
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

	report := PayerReport{
		OriginatorNodeID: 1,
		StartSequenceID:  0,
		EndSequenceID:    10,
		PayersMerkleRoot: randomBytes32(),
		ActiveNodeIDs:    []uint32{1},
	}

	_, err := store.StoreReport(ctx, &report)
	require.NoError(t, err)

	attestation := &PayerReportAttestation{
		Report: &report,
		NodeSignature: NodeSignature{
			NodeID:    uint32(math.MaxInt32) + 1,
			Signature: []byte("sig"),
		},
	}

	err = store.StoreAttestation(ctx, attestation)
	require.ErrorIs(t, err, ErrOriginatorNodeIDTooLarge)
}

func TestSetReportAttestationStatus(t *testing.T) {
	store := createTestStore(t)
	ctx := context.Background()

	report := PayerReport{
		OriginatorNodeID: 1,
		StartSequenceID:  0,
		EndSequenceID:    10,
		PayersMerkleRoot: randomBytes32(),
		ActiveNodeIDs:    []uint32{1},
	}

	reportID, err := store.StoreReport(ctx, &report)
	require.NoError(t, err)

	// Move from Pending -> Approved
	require.NoError(
		t,
		store.SetReportAttestationStatus(
			ctx,
			reportID,
			[]AttestationStatus{AttestationPending},
			AttestationApproved,
		),
	)

	fetched, err := store.FetchReport(ctx, reportID)
	require.NoError(t, err)
	require.Equal(t, AttestationStatus(AttestationApproved), fetched.AttestationStatus)
}

func TestCreatePayerReport(t *testing.T) {
	store := createTestStore(t)
	ctx := context.Background()

	report := PayerReport{
		OriginatorNodeID: 3,
		StartSequenceID:  0,
		EndSequenceID:    10,
		PayersMerkleRoot: randomBytes32(),
		ActiveNodeIDs:    []uint32{3},
	}

	// Build a minimal payer envelope (group message payload is fine for this path)
	payerProto := envTestUtils.CreatePayerEnvelope(t, report.OriginatorNodeID)
	payerEnv, err := envsWrapper.NewPayerEnvelope(payerProto)
	require.NoError(t, err)

	reportID, err := store.CreatePayerReport(ctx, &report, payerEnv)
	require.NoError(t, err)
	require.NotNil(t, reportID)

	fetched, err := store.FetchReport(ctx, reportID)
	require.NoError(t, err)
	require.Equal(t, report.OriginatorNodeID, fetched.OriginatorNodeID)

	// Ensure a staged originator envelope was created
	staged, err := store.Queries().
		SelectStagedOriginatorEnvelopes(ctx, queries.SelectStagedOriginatorEnvelopesParams{LastSeenID: 0, NumRows: 10})
	require.NoError(t, err)
	require.Len(t, staged, 1)
}

func TestCreateAttestation(t *testing.T) {
	store := createTestStore(t)
	ctx := context.Background()

	report := PayerReport{
		OriginatorNodeID: 4,
		StartSequenceID:  0,
		EndSequenceID:    10,
		PayersMerkleRoot: randomBytes32(),
		ActiveNodeIDs:    []uint32{4},
	}

	reportID, err := store.StoreReport(ctx, &report)
	require.NoError(t, err)

	attestation := &PayerReportAttestation{
		Report: &report,
		NodeSignature: NodeSignature{
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
	require.Equal(t, AttestationStatus(AttestationApproved), fetchedReport.AttestationStatus)
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

	report := PayerReport{
		OriginatorNodeID:    7,
		StartSequenceID:     0,
		EndSequenceID:       10,
		EndMinuteSinceEpoch: uint32(time.Now().Unix() / 60),
		PayersMerkleRoot:    randomBytes32(),
		ActiveNodeIDs:       []uint32{7},
	}
	// Build the originator envelope containing the payer report payload
	clientEnv := createPayerReportClientEnvelope(&report)
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

	require.NoError(t, store.StoreSyncedReport(ctx, originatorEnv, payerID))

	reportID, err := report.ID()
	require.NoError(t, err)
	fetched, err := store.FetchReport(ctx, reportID)
	require.NoError(t, err)
	require.Equal(t, report.OriginatorNodeID, fetched.OriginatorNodeID)
}

func TestStoreSyncedAttestation(t *testing.T) {
	store := createTestStore(t)
	ctx := context.Background()

	// First create and store the base report so that the attestation references a real report ID
	baseReport := PayerReport{
		OriginatorNodeID: 8,
		StartSequenceID:  0,
		EndSequenceID:    10,
		PayersMerkleRoot: randomBytes32(),
		ActiveNodeIDs:    []uint32{8},
	}
	reportID, err := store.StoreReport(ctx, &baseReport)
	require.NoError(t, err)

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
