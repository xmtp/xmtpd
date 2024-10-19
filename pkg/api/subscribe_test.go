package api_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/api"
	"github.com/xmtp/xmtpd/pkg/db"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api"
	"github.com/xmtp/xmtpd/pkg/testutils"
	testUtilsApi "github.com/xmtp/xmtpd/pkg/testutils/api"
	envelopeTestUtils "github.com/xmtp/xmtpd/pkg/testutils/envelopes"
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
				envelopeTestUtils.CreateOriginatorEnvelope(t, 1, 1),
			),
		},
		{
			OriginatorNodeID:     2,
			OriginatorSequenceID: 1,
			Topic:                []byte("topicA"),
			OriginatorEnvelope: testutils.Marshal(
				t,
				envelopeTestUtils.CreateOriginatorEnvelope(t, 2, 1),
			),
		},
		// Later rows
		{
			OriginatorNodeID:     1,
			OriginatorSequenceID: 2,
			Topic:                []byte("topicB"),
			OriginatorEnvelope: testutils.Marshal(
				t,
				envelopeTestUtils.CreateOriginatorEnvelope(t, 1, 2),
			),
		},
		{
			OriginatorNodeID:     2,
			OriginatorSequenceID: 2,
			Topic:                []byte("topicB"),
			OriginatorEnvelope: testutils.Marshal(
				t,
				envelopeTestUtils.CreateOriginatorEnvelope(t, 2, 2),
			),
		},
		{
			OriginatorNodeID:     1,
			OriginatorSequenceID: 3,
			Topic:                []byte("topicA"),
			OriginatorEnvelope: testutils.Marshal(
				t,
				envelopeTestUtils.CreateOriginatorEnvelope(t, 1, 3),
			),
		},
	}

	return testUtilsApi.NewTestAPIClient(t)
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
	stream message_api.ReplicationApi_SubscribeEnvelopesClient,
	expectedIndices []int,
) {
	for i := 0; i < len(expectedIndices); {
		envs, err := stream.Recv()
		require.NoError(t, err)
		for _, env := range envs.Envelopes {
			expected := allRows[expectedIndices[i]]
			actual := envelopeTestUtils.UnmarshalUnsignedOriginatorEnvelope(
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
	stream, err := client.SubscribeEnvelopes(
		ctx,
		&message_api.SubscribeEnvelopesRequest{
			Query: &message_api.EnvelopesQuery{
				LastSeen: nil,
			},
		},
	)
	require.NoError(t, err)

	insertAdditionalRows(t, db)
	validateUpdates(t, stream, []int{2, 3, 4})
}

func TestSubscribeEnvelopesByTopic(t *testing.T) {
	client, store, cleanup := setupTest(t)
	defer cleanup()
	insertInitialRows(t, store)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	stream, err := client.SubscribeEnvelopes(
		ctx,
		&message_api.SubscribeEnvelopesRequest{
			Query: &message_api.EnvelopesQuery{
				Topics:   []db.Topic{db.Topic("topicA"), []byte("topicC")},
				LastSeen: nil,
			},
		},
	)
	require.NoError(t, err)

	insertAdditionalRows(t, store)
	validateUpdates(t, stream, []int{4})
}

func TestSubscribeEnvelopesByOriginator(t *testing.T) {
	client, db, cleanup := setupTest(t)
	defer cleanup()
	insertInitialRows(t, db)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	stream, err := client.SubscribeEnvelopes(
		ctx,
		&message_api.SubscribeEnvelopesRequest{
			Query: &message_api.EnvelopesQuery{
				OriginatorNodeIds: []uint32{1, 3},
				LastSeen:          nil,
			},
		},
	)
	require.NoError(t, err)

	insertAdditionalRows(t, db)
	validateUpdates(t, stream, []int{2, 4})
}

func TestSimultaneousSubscriptions(t *testing.T) {
	client, store, cleanup := setupTest(t)
	defer cleanup()
	insertInitialRows(t, store)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	stream1, err := client.SubscribeEnvelopes(
		ctx,
		&message_api.SubscribeEnvelopesRequest{
			Query: &message_api.EnvelopesQuery{},
		},
	)
	require.NoError(t, err)

	stream2, err := client.SubscribeEnvelopes(
		ctx,
		&message_api.SubscribeEnvelopesRequest{
			Query: &message_api.EnvelopesQuery{
				Topics:   []db.Topic{db.Topic("topicB")},
				LastSeen: nil,
			},
		},
	)
	require.NoError(t, err)

	stream3, err := client.SubscribeEnvelopes(
		ctx,
		&message_api.SubscribeEnvelopesRequest{
			Query: &message_api.EnvelopesQuery{
				OriginatorNodeIds: []uint32{2},
				LastSeen:          nil,
			},
		},
	)
	require.NoError(t, err)

	insertAdditionalRows(t, store)
	validateUpdates(t, stream1, []int{2, 3, 4})
	validateUpdates(t, stream2, []int{2, 3})
	validateUpdates(t, stream3, []int{3})
}

func TestSubscribeEnvelopesFromCursor(t *testing.T) {
	client, store, cleanup := setupTest(t)
	defer cleanup()
	insertInitialRows(t, store)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	stream, err := client.SubscribeEnvelopes(
		ctx,
		&message_api.SubscribeEnvelopesRequest{
			Query: &message_api.EnvelopesQuery{
				Topics:   []db.Topic{db.Topic("topicA"), []byte("topicC")},
				LastSeen: &message_api.VectorClock{NodeIdToSequenceId: map[uint32]uint64{1: 1}},
			},
		},
	)
	require.NoError(t, err)

	insertAdditionalRows(t, store)
	validateUpdates(t, stream, []int{1, 4})
}

func TestSubscribeEnvelopesFromEmptyCursor(t *testing.T) {
	client, store, cleanup := setupTest(t)
	defer cleanup()
	insertInitialRows(t, store)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	stream, err := client.SubscribeEnvelopes(
		ctx,
		&message_api.SubscribeEnvelopesRequest{
			Query: &message_api.EnvelopesQuery{
				Topics:   []db.Topic{db.Topic("topicA"), []byte("topicC")},
				LastSeen: &message_api.VectorClock{NodeIdToSequenceId: map[uint32]uint64{}},
			},
		},
	)
	require.NoError(t, err)

	insertAdditionalRows(t, store)
	validateUpdates(t, stream, []int{0, 1, 4})
}

func TestSubscribeEnvelopesInvalidRequest(t *testing.T) {
	client, _, cleanup := setupTest(t)
	defer cleanup()

	stream, err := client.SubscribeEnvelopes(
		context.Background(),
		&message_api.SubscribeEnvelopesRequest{
			Query: &message_api.EnvelopesQuery{
				Topics:            []db.Topic{db.Topic("topicA")},
				OriginatorNodeIds: []uint32{1},
				LastSeen:          nil,
			},
		},
	)
	require.NoError(t, err)
	_, err = stream.Recv()
	require.Error(t, err)
}
