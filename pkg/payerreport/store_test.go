package payerreport

import (
	"context"
	"math"
	"testing"

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
