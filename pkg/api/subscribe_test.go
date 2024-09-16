package api_test

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api"
	"github.com/xmtp/xmtpd/pkg/testutils"
)

var allRows = []queries.InsertGatewayEnvelopeParams{
	// Initial rows
	{
		OriginatorNodeID:     1,
		OriginatorSequenceID: 1,
		Topic:                []byte("topicA"),
		OriginatorEnvelope:   []byte("envelope1"),
	},
	{
		OriginatorNodeID:     2,
		OriginatorSequenceID: 1,
		Topic:                []byte("topicA"),
		OriginatorEnvelope:   []byte("envelope2"),
	},
	// Later rows
	{
		OriginatorNodeID:     1,
		OriginatorSequenceID: 2,
		Topic:                []byte("topicA"),
		OriginatorEnvelope:   []byte("envelope3"),
	},
	{
		OriginatorNodeID:     2,
		OriginatorSequenceID: 2,
		Topic:                []byte("topicA"),
		OriginatorEnvelope:   []byte("envelope4"),
	},
	{
		OriginatorNodeID:     1,
		OriginatorSequenceID: 3,
		Topic:                []byte("topicA"),
		OriginatorEnvelope:   []byte("envelope5"),
	},
}

func insertInitialRows(t *testing.T, store *sql.DB) {
	testutils.InsertGatewayEnvelopes(t, store, []queries.InsertGatewayEnvelopeParams{
		allRows[0], allRows[1],
	})
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
		fmt.Println("waiting for update")
		envs, err := stream.Recv()
		fmt.Printf("got update of length %d\n", len(envs.Envelopes))
		require.NoError(t, err)
		for _, env := range envs.Envelopes {
			expected := allRows[expectedIndices[i]].OriginatorEnvelope
			require.Equal(t, expected, testutils.Marshal(t, env))
			i++
		}
	}
}

func TestQAllEnvelopes(t *testing.T) {
	t.Skip("skipping test")
	client, db, cleanup := testutils.NewTestAPIClient(t)
	defer cleanup()
	insertInitialRows(t, db)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	res, err := client.QueryEnvelopes(
		ctx,
		&message_api.QueryEnvelopesRequest{
			Query: &message_api.EnvelopesQuery{
				Filter:   nil,
				LastSeen: &message_api.VectorClock{},
			},
		},
	)
	require.NoError(t, err)

	insertAdditionalRows(t, db)
	require.Equal(t, 5, len(res.Envelopes))
}

func TestSubscribeAllEnvelopes(t *testing.T) {
	client, db, cleanup := testutils.NewTestAPIClient(t)
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
