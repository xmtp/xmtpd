package message_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"connectrpc.com/connect"
	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/api/message"
	"github.com/xmtp/xmtpd/pkg/db"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/envelopes"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api"
	message_apiconnect "github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api/message_apiconnect"
	"github.com/xmtp/xmtpd/pkg/testutils"
	testUtilsApi "github.com/xmtp/xmtpd/pkg/testutils/api"
	envelopeTestUtils "github.com/xmtp/xmtpd/pkg/testutils/envelopes"
	"github.com/xmtp/xmtpd/pkg/topic"
)

var (
	topicA = topic.NewTopic(topic.TopicKindGroupMessagesV1, []byte("topicA")).Bytes()
	topicB = topic.NewTopic(topic.TopicKindGroupMessagesV1, []byte("topicB")).Bytes()
	topicC = topic.NewTopic(topic.TopicKindGroupMessagesV1, []byte("topicC")).Bytes()
)
var allRows []queries.InsertGatewayEnvelopeParams

func setupTest(
	t *testing.T,
) (message_apiconnect.ReplicationApiClient, *sql.DB, testUtilsApi.APIServerMocks) {
	var (
		client      = testUtilsApi.NewTestGRPCReplicationAPIClient(t, "localhost:0")
		dbHandle, _ = testutils.NewDB(t, t.Context())
		_, _, mocks = testUtilsApi.NewTestFullServer(t)
		payerID     = db.NullInt32(testutils.CreatePayer(t, dbHandle))
	)

	allRows = []queries.InsertGatewayEnvelopeParams{
		// Initial rows
		{
			OriginatorNodeID:     100,
			OriginatorSequenceID: 1,
			Topic:                topicA,
			PayerID:              payerID,
			OriginatorEnvelope: testutils.Marshal(
				t,
				envelopeTestUtils.CreateOriginatorEnvelopeWithTopic(t, 100, 1, topicA),
			),
		},
		{
			OriginatorNodeID:     200,
			OriginatorSequenceID: 1,
			Topic:                topicA,
			PayerID:              payerID,
			OriginatorEnvelope: testutils.Marshal(
				t,
				envelopeTestUtils.CreateOriginatorEnvelopeWithTopic(t, 200, 1, topicA),
			),
		},
		// Later rows
		{
			OriginatorNodeID:     100,
			OriginatorSequenceID: 2,
			Topic:                topicB,
			PayerID:              payerID,
			OriginatorEnvelope: testutils.Marshal(
				t,
				envelopeTestUtils.CreateOriginatorEnvelopeWithTopic(t, 100, 2, topicB),
			),
		},
		{
			OriginatorNodeID:     200,
			OriginatorSequenceID: 2,
			Topic:                topicB,
			PayerID:              payerID,
			OriginatorEnvelope: testutils.Marshal(
				t,
				envelopeTestUtils.CreateOriginatorEnvelopeWithTopic(t, 200, 2, topicB),
			),
		},
		{
			OriginatorNodeID:     100,
			OriginatorSequenceID: 3,
			Topic:                topicA,
			PayerID:              payerID,
			OriginatorEnvelope: testutils.Marshal(
				t,
				envelopeTestUtils.CreateOriginatorEnvelopeWithTopic(t, 100, 3, topicA),
			),
		},
	}

	return client, dbHandle, mocks
}

func insertInitialRows(t *testing.T, store *sql.DB) {
	testutils.InsertGatewayEnvelopes(t, store, []queries.InsertGatewayEnvelopeParams{
		allRows[0], allRows[1],
	})
	time.Sleep(message.SubscribeWorkerPollTime + 100*time.Millisecond)
}

func insertAdditionalRows(t *testing.T, store *sql.DB, notifyChan ...chan bool) {
	testutils.InsertGatewayEnvelopes(t, store, []queries.InsertGatewayEnvelopeParams{
		allRows[2], allRows[3], allRows[4],
	}, notifyChan...)
}

func validateUpdates(
	t *testing.T,
	stream *connect.ServerStreamForClient[message_api.SubscribeEnvelopesResponse],
	expectedIndices []int,
) {
	receivedCount := 0

	for stream.Receive() && receivedCount < len(expectedIndices) {
		envs := stream.Msg()

		for _, env := range envs.Envelopes {
			require.Less(t, receivedCount, len(expectedIndices),
				"received more envelopes than expected")

			expected := allRows[expectedIndices[receivedCount]]
			actual := envelopeTestUtils.UnmarshalUnsignedOriginatorEnvelope(
				t,
				env.UnsignedOriginatorEnvelope,
			)

			require.EqualValues(t, expected.OriginatorNodeID, actual.OriginatorNodeId)
			require.EqualValues(t, expected.OriginatorSequenceID, actual.OriginatorSequenceId)
			require.Equal(t, expected.OriginatorEnvelope, testutils.Marshal(t, env))

			receivedCount++
		}
	}

	// Check for stream errors
	require.NoError(t, stream.Err())
}

func TestSubscribeEnvelopesAll(t *testing.T) {
	client, db, _ := setupTest(t)

	insertInitialRows(t, db)

	ctx := t.Context()
	stream, err := client.SubscribeEnvelopes(
		ctx,
		connect.NewRequest(&message_api.SubscribeEnvelopesRequest{
			Query: &message_api.EnvelopesQuery{
				LastSeen: nil,
			},
		}),
	)
	require.NoError(t, err)

	insertAdditionalRows(t, db)
	validateUpdates(t, stream, []int{2, 3, 4})
}

func TestSubscribeEnvelopesByTopic(t *testing.T) {
	client, store, _ := setupTest(t)
	insertInitialRows(t, store)

	ctx := t.Context()
	stream, err := client.SubscribeEnvelopes(
		ctx,
		connect.NewRequest(&message_api.SubscribeEnvelopesRequest{
			Query: &message_api.EnvelopesQuery{
				Topics:   []db.Topic{topicA, topicC},
				LastSeen: nil,
			},
		}),
	)
	require.NoError(t, err)

	insertAdditionalRows(t, store)
	validateUpdates(t, stream, []int{4})
}

func TestSubscribeEnvelopesByOriginator(t *testing.T) {
	client, db, _ := setupTest(t)
	insertInitialRows(t, db)

	ctx := t.Context()
	stream, err := client.SubscribeEnvelopes(
		ctx,
		connect.NewRequest(&message_api.SubscribeEnvelopesRequest{
			Query: &message_api.EnvelopesQuery{
				OriginatorNodeIds: []uint32{100, 300},
				LastSeen:          nil,
			},
		}),
	)
	require.NoError(t, err)

	insertAdditionalRows(t, db)
	validateUpdates(t, stream, []int{2, 4})
}

func TestSimultaneousSubscriptions(t *testing.T) {
	client, store, _ := setupTest(t)
	insertInitialRows(t, store)

	ctx := t.Context()
	stream1, err := client.SubscribeEnvelopes(
		ctx,
		connect.NewRequest(&message_api.SubscribeEnvelopesRequest{
			Query: &message_api.EnvelopesQuery{},
		}),
	)
	require.NoError(t, err)

	stream2, err := client.SubscribeEnvelopes(
		ctx,
		connect.NewRequest(&message_api.SubscribeEnvelopesRequest{
			Query: &message_api.EnvelopesQuery{
				Topics:   []db.Topic{topicB},
				LastSeen: nil,
			},
		}),
	)
	require.NoError(t, err)

	stream3, err := client.SubscribeEnvelopes(
		ctx,
		connect.NewRequest(&message_api.SubscribeEnvelopesRequest{
			Query: &message_api.EnvelopesQuery{
				OriginatorNodeIds: []uint32{200},
				LastSeen:          nil,
			},
		}),
	)
	require.NoError(t, err)

	insertAdditionalRows(t, store)
	validateUpdates(t, stream1, []int{2, 3, 4})
	validateUpdates(t, stream2, []int{2, 3})
	validateUpdates(t, stream3, []int{3})
}

func TestSubscribeEnvelopesFromCursor(t *testing.T) {
	client, store, _ := setupTest(t)
	insertInitialRows(t, store)

	ctx := t.Context()
	stream, err := client.SubscribeEnvelopes(
		ctx,
		connect.NewRequest(&message_api.SubscribeEnvelopesRequest{
			Query: &message_api.EnvelopesQuery{
				Topics:   []db.Topic{topicA, topicC},
				LastSeen: &envelopes.Cursor{NodeIdToSequenceId: map[uint32]uint64{100: 1}},
			},
		}),
	)
	require.NoError(t, err)

	insertAdditionalRows(t, store)
	validateUpdates(t, stream, []int{1, 4})
}

func TestSubscribeEnvelopesFromEmptyCursor(t *testing.T) {
	client, store, _ := setupTest(t)
	insertInitialRows(t, store)

	ctx := t.Context()
	stream, err := client.SubscribeEnvelopes(
		ctx,
		connect.NewRequest(&message_api.SubscribeEnvelopesRequest{
			Query: &message_api.EnvelopesQuery{
				Topics:   []db.Topic{topicA, topicC},
				LastSeen: &envelopes.Cursor{NodeIdToSequenceId: map[uint32]uint64{}},
			},
		}),
	)
	require.NoError(t, err)

	insertAdditionalRows(t, store)
	validateUpdates(t, stream, []int{0, 1, 4})
}

func TestSubscribeEnvelopesInvalidRequest(t *testing.T) {
	client, _, _ := setupTest(t)

	stream, err := client.SubscribeEnvelopes(
		context.Background(),
		connect.NewRequest(&message_api.SubscribeEnvelopesRequest{
			Query: &message_api.EnvelopesQuery{
				Topics:            []db.Topic{topicA},
				OriginatorNodeIds: []uint32{1},
				LastSeen:          nil,
			},
		}),
	)
	require.NoError(t, err)

	shouldReceive := stream.Receive()
	require.False(t, shouldReceive)

	msg := stream.Msg()
	require.Nil(t, msg)

	err = stream.Err()
	require.Error(t, err)
}
