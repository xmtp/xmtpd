package api_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/api"
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
	time.Sleep(api.SubscribeWorkerPollTime + 100*time.Millisecond)
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

func TestSubscribeEnvelopesAll(t *testing.T) {
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
						LastSeen: nil,
					},
				},
			},
		},
	)
	require.NoError(t, err)

	insertAdditionalRows(t, db)
	validateUpdates(t, stream, []int{2, 3, 4})
}

func TestSubscribeEnvelopesByTopic(t *testing.T) {
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
						Filter:   &message_api.EnvelopesQuery_Topic{Topic: []byte("topicA")},
						LastSeen: nil,
					},
				},
				{
					Query: &message_api.EnvelopesQuery{
						Filter:   &message_api.EnvelopesQuery_Topic{Topic: []byte("topicC")},
						LastSeen: nil,
					},
				},
			},
		},
	)
	require.NoError(t, err)

	insertAdditionalRows(t, db)
	validateUpdates(t, stream, []int{4})
}

func TestSubscribeEnvelopesByOriginator(t *testing.T) {
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
						Filter:   &message_api.EnvelopesQuery_OriginatorNodeId{OriginatorNodeId: 1},
						LastSeen: nil,
					},
				},
				{
					Query: &message_api.EnvelopesQuery{
						Filter:   &message_api.EnvelopesQuery_OriginatorNodeId{OriginatorNodeId: 3},
						LastSeen: nil,
					},
				},
			},
		},
	)
	require.NoError(t, err)

	insertAdditionalRows(t, db)
	validateUpdates(t, stream, []int{2, 4})
}

func TestSimultaneousSubscriptions(t *testing.T) {
	client, db, cleanup := setupTest(t)
	defer cleanup()
	insertInitialRows(t, db)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	stream1, err := client.BatchSubscribeEnvelopes(
		ctx,
		&message_api.BatchSubscribeEnvelopesRequest{
			Requests: []*message_api.BatchSubscribeEnvelopesRequest_SubscribeEnvelopesRequest{
				{
					Query: &message_api.EnvelopesQuery{
						Filter:   nil,
						LastSeen: nil,
					},
				},
			},
		},
	)
	require.NoError(t, err)

	stream2, err := client.BatchSubscribeEnvelopes(
		ctx,
		&message_api.BatchSubscribeEnvelopesRequest{
			Requests: []*message_api.BatchSubscribeEnvelopesRequest_SubscribeEnvelopesRequest{
				{
					Query: &message_api.EnvelopesQuery{
						Filter:   &message_api.EnvelopesQuery_Topic{Topic: []byte("topicB")},
						LastSeen: nil,
					},
				},
			},
		},
	)
	require.NoError(t, err)

	stream3, err := client.BatchSubscribeEnvelopes(
		ctx,
		&message_api.BatchSubscribeEnvelopesRequest{
			Requests: []*message_api.BatchSubscribeEnvelopesRequest_SubscribeEnvelopesRequest{
				{
					Query: &message_api.EnvelopesQuery{
						Filter:   &message_api.EnvelopesQuery_OriginatorNodeId{OriginatorNodeId: 2},
						LastSeen: nil,
					},
				},
			},
		},
	)
	require.NoError(t, err)

	insertAdditionalRows(t, db)
	validateUpdates(t, stream1, []int{2, 3, 4})
	validateUpdates(t, stream2, []int{2, 3})
	validateUpdates(t, stream3, []int{3})
}

func TestSubscribeEnvelopesInvalidRequest(t *testing.T) {
	client, _, cleanup := setupTest(t)
	defer cleanup()

	stream, err := client.BatchSubscribeEnvelopes(
		context.Background(),
		&message_api.BatchSubscribeEnvelopesRequest{
			Requests: []*message_api.BatchSubscribeEnvelopesRequest_SubscribeEnvelopesRequest{
				{
					Query: &message_api.EnvelopesQuery{
						Filter:   &message_api.EnvelopesQuery_OriginatorNodeId{OriginatorNodeId: 1},
						LastSeen: nil,
					},
				},
				{
					Query: &message_api.EnvelopesQuery{
						Filter:   nil,
						LastSeen: nil,
					},
				},
			},
		},
	)
	require.NoError(t, err)
	_, err = stream.Recv()
	require.Error(t, err)

	stream, err = client.BatchSubscribeEnvelopes(
		context.Background(),
		&message_api.BatchSubscribeEnvelopesRequest{
			Requests: []*message_api.BatchSubscribeEnvelopesRequest_SubscribeEnvelopesRequest{
				{
					Query: &message_api.EnvelopesQuery{
						Filter:   nil,
						LastSeen: nil,
					},
				},
				{
					Query: &message_api.EnvelopesQuery{
						Filter:   &message_api.EnvelopesQuery_Topic{Topic: []byte("topicA")},
						LastSeen: nil,
					},
				},
			},
		},
	)
	require.NoError(t, err)
	_, err = stream.Recv()
	require.Error(t, err)

	stream, err = client.BatchSubscribeEnvelopes(
		context.Background(),
		&message_api.BatchSubscribeEnvelopesRequest{
			Requests: []*message_api.BatchSubscribeEnvelopesRequest_SubscribeEnvelopesRequest{
				{
					Query: &message_api.EnvelopesQuery{
						Filter:   &message_api.EnvelopesQuery_OriginatorNodeId{OriginatorNodeId: 1},
						LastSeen: nil,
					},
				},
				{
					Query: &message_api.EnvelopesQuery{
						Filter:   &message_api.EnvelopesQuery_Topic{Topic: []byte("topicA")},
						LastSeen: nil,
					},
				},
			},
		},
	)
	require.NoError(t, err)
	_, err = stream.Recv()
	require.Error(t, err)
}
