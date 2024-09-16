package api_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api"
	"github.com/xmtp/xmtpd/pkg/testutils"
)

var allRows []queries.InsertGatewayEnvelopeParams

func setupTest(t *testing.T) (message_api.ReplicationApiClient, *sql.DB, func()) {
	allRows = []queries.InsertGatewayEnvelopeParams{
		// Initial rows
		{
			OriginatorNodeID:     1,
			OriginatorSequenceID: 1,
			Topic:                []byte("topicA"),
			OriginatorEnvelope: testutils.Marshal(
				t,
				testutils.CreateOriginatorEnvelope(t, 1, 1),
			),
		},
		{
			OriginatorNodeID:     2,
			OriginatorSequenceID: 1,
			Topic:                []byte("topicA"),
			OriginatorEnvelope: testutils.Marshal(
				t,
				testutils.CreateOriginatorEnvelope(t, 2, 1),
			),
		},
		// Later rows
		{
			OriginatorNodeID:     1,
			OriginatorSequenceID: 2,
			Topic:                []byte("topicB"),
			OriginatorEnvelope: testutils.Marshal(
				t,
				testutils.CreateOriginatorEnvelope(t, 1, 2),
			),
		},
		{
			OriginatorNodeID:     2,
			OriginatorSequenceID: 2,
			Topic:                []byte("topicB"),
			OriginatorEnvelope: testutils.Marshal(
				t,
				testutils.CreateOriginatorEnvelope(t, 2, 2),
			),
		},
		{
			OriginatorNodeID:     1,
			OriginatorSequenceID: 3,
			Topic:                []byte("topicA"),
			OriginatorEnvelope: testutils.Marshal(
				t,
				testutils.CreateOriginatorEnvelope(t, 1, 3),
			),
		},
	}
	return testutils.NewTestAPIClient(t)
}

func insertInitialRows(t *testing.T, store *sql.DB) {
	testutils.InsertGatewayEnvelopes(t, store, []queries.InsertGatewayEnvelopeParams{
		allRows[0], allRows[1],
	})
	// For the sake of the test, make sure the subscribeWorker has time to receive the rows
	time.Sleep(1 * time.Second)
}

func insertAdditionalRows(t *testing.T, store *sql.DB, notifyChan ...chan bool) {
	testutils.InsertGatewayEnvelopes(t, store, []queries.InsertGatewayEnvelopeParams{
		allRows[2], allRows[3], allRows[4],
	}, notifyChan...)
}

func validateUpdates(
	t *testing.T,
	stream message_api.ReplicationApi_BatchSubscribeEnvelopesClient,
	expectedIndices []int,
) {
	for i := 0; i < len(expectedIndices); {
		envs, err := stream.Recv()
		require.NoError(t, err)
		for _, env := range envs.Envelopes {
			expected := allRows[expectedIndices[i]]
			actual := testutils.UnmarshalUnsignedOriginatorEnvelope(
				t,
				env.UnsignedOriginatorEnvelope,
			)
			require.Equal(t, uint32(expected.OriginatorNodeID), actual.OriginatorNodeId)
			require.Equal(t, uint64(expected.OriginatorSequenceID), actual.OriginatorSequenceId)
			require.Equal(t, expected.OriginatorEnvelope, testutils.Marshal(t, env))
			i++
		}
	}
}

func TestSubscribeAllEnvelopes(t *testing.T) {
	client, db, cleanup := setupTest(t)
	defer cleanup()
	insertInitialRows(t, db)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	stream, err := client.BatchSubscribeEnvelopes(
		ctx,
		&message_api.BatchSubscribeEnvelopesRequest{
			Requests: []*message_api.BatchSubscribeEnvelopesRequest_SubscribeEnvelopesRequest{
				{
					Query: &message_api.EnvelopesQuery{
						Filter:   nil,
						LastSeen: &message_api.VectorClock{},
					},
				},
			},
		},
	)
	require.NoError(t, err)

	insertAdditionalRows(t, db)
	validateUpdates(t, stream, []int{2, 3, 4})
}
