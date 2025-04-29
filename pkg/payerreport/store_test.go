package payerreport

import (
	"context"
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/testutils"
)

func createTestStore(t *testing.T) *Store {
	log := testutils.NewLog(t)
	db, _, cleanup := testutils.NewDB(t, context.Background())
	t.Cleanup(cleanup)

	return NewStore(queries.New(db), log)
}

func insertRandomReport(
	t *testing.T,
	store *Store,
) *PayerReportWithStatus {
	startID := testutils.RandomInt64()
	insertParams := queries.InsertOrIgnorePayerReportParams{
		ID:               testutils.RandomBytes(32),
		OriginatorNodeID: testutils.RandomInt32(),
		StartSequenceID:  startID,
		EndSequenceID:    startID + 10,
		PayersMerkleRoot: testutils.RandomBytes(32),
		PayersLeafCount:  10,
		NodesHash:        testutils.RandomBytes(32),
		NodesCount:       10,
	}
	require.NoError(t, store.queries.InsertOrIgnorePayerReport(t.Context(), insertParams))

	returnedVal, err := store.FetchReport(t.Context(), insertParams.ID)
	require.NoError(t, err)
	return returnedVal
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
				OriginatorNodeID: 1,
				StartSequenceID:  0,
				EndSequenceID:    2,
				PayersMerkleRoot: randomBytes32(),
				PayersLeafCount:  1,
				NodesHash:        randomBytes32(),
				NodesCount:       1,
			},
			expectErr: false,
		},
		{
			name: "invalid node ID",
			report: PayerReport{
				OriginatorNodeID: uint32(math.MaxInt32) + 1,
				StartSequenceID:  0,
				EndSequenceID:    2,
				PayersMerkleRoot: randomBytes32(),
				PayersLeafCount:  1,
				NodesHash:        randomBytes32(),
				NodesCount:       1,
			},
			expectErr: true,
		},
		{
			name: "invalid nodes count",
			report: PayerReport{
				OriginatorNodeID: 1,
				StartSequenceID:  0,
				EndSequenceID:    2,
				PayersMerkleRoot: randomBytes32(),
				PayersLeafCount:  1,
				NodesHash:        randomBytes32(),
				NodesCount:       uint32(math.MaxInt32) + 1,
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
				require.Equal(t, c.report, *storedReport)
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
		PayersLeafCount:  1,
		NodesHash:        randomBytes32(),
		NodesCount:       1,
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
	report2 := insertRandomReport(t, store)
	// Set the second report's status to Approved
	require.NoError(
		t,
		store.queries.SetReportAttestationStatus(
			t.Context(),
			queries.SetReportAttestationStatusParams{
				NewStatus:  AttestationApproved,
				ReportID:   report2.ID[:],
				PrevStatus: []int16{int16(AttestationPending)},
			},
		),
	)
	report3 := insertRandomReport(t, store)

	cases := []struct {
		name        string
		expectedIDs [][]byte
		query       *FetchReportsQuery
	}{{
		name:        "Get all with created after",
		expectedIDs: [][]byte{report1.ID[:], report2.ID[:], report3.ID[:]},

		query: NewFetchReportsQuery().WithCreatedAfter(report1.CreatedAt.Add(-5 * time.Second)),
	}, {
		name:        "Get newest 2",
		expectedIDs: [][]byte{report2.ID[:], report3.ID[:]},
		query:       NewFetchReportsQuery().WithCreatedAfter(report1.CreatedAt),
	}, {
		name:        "Only approved",
		expectedIDs: [][]byte{report2.ID[:]},
		query: NewFetchReportsQuery().WithCreatedAfter(time.Unix(1, 0)).
			WithAttestationStatus(AttestationApproved),
	}, {
		name:        "Multiple statuses",
		expectedIDs: [][]byte{report2.ID[:]},
		query: NewFetchReportsQuery().
			WithAttestationStatus(AttestationApproved, AttestationRejected),
	}, {
		name:        "No results",
		expectedIDs: [][]byte{},
		query: NewFetchReportsQuery().WithCreatedAfter(time.Unix(1, 0)).
			WithAttestationStatus(AttestationRejected),
	}, {
		name:        "No Params",
		expectedIDs: [][]byte{report1.ID[:], report2.ID[:], report3.ID[:]},
		query:       NewFetchReportsQuery(),
	}, {
		name:        "With start sequence ID",
		expectedIDs: [][]byte{report1.ID[:]},
		query:       NewFetchReportsQuery().WithStartSequenceID(report1.StartSequenceID),
	}, {
		name:        "With end sequence ID",
		expectedIDs: [][]byte{report1.ID[:]},
		query:       NewFetchReportsQuery().WithEndSequenceID(report1.EndSequenceID),
	}, {
		name:        "With start and end sequence ID",
		expectedIDs: [][]byte{report1.ID[:]},
		query: NewFetchReportsQuery().WithStartSequenceID(report1.StartSequenceID).
			WithEndSequenceID(report1.EndSequenceID),
	}}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			results, err := store.FetchReports(t.Context(), c.query)
			require.NoError(t, err)
			require.Len(t, results, len(c.expectedIDs))

			returnedIDs := make([][]byte, len(results))
			for idx, result := range results {
				returnedIDs[idx] = result.ID[:]
			}

			require.ElementsMatch(t, c.expectedIDs, returnedIDs)
		})
	}
}
