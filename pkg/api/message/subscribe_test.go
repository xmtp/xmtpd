package message_test

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math/rand"
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
	topicA  = topic.NewTopic(topic.TopicKindGroupMessagesV1, []byte("topicA")).Bytes()
	topicB  = topic.NewTopic(topic.TopicKindGroupMessagesV1, []byte("topicB")).Bytes()
	topicC  = topic.NewTopic(topic.TopicKindGroupMessagesV1, []byte("topicC")).Bytes()
	allRows = make([]queries.InsertGatewayEnvelopeParams, 0)
)

func setupTest(
	t *testing.T,
) (message_apiconnect.ReplicationApiClient, *sql.DB, testUtilsApi.APIServerMocks) {
	var (
		suite   = testUtilsApi.NewTestAPIServer(t)
		payerID = db.NullInt32(testutils.CreatePayer(t, suite.DB))
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

	return suite.ClientReplication, suite.DB, suite.APIServerMocks
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
	type key struct {
		nodeID int32
		seqID  int64
	}

	// Build the set of expected (nodeID, seqID) we must observe.
	expected := make(map[key]queries.InsertGatewayEnvelopeParams, len(expectedIndices))
	for _, idx := range expectedIndices {
		r := allRows[idx]
		expected[key{
			nodeID: r.OriginatorNodeID,
			seqID:  r.OriginatorSequenceID,
		}] = r
	}

	seen := make(map[key]struct{}, len(expectedIndices))
	lastSeqByNode := make(map[int32]int64)

	for len(seen) < len(expected) {
		if !stream.Receive() {
			break
		}

		msg := stream.Msg()
		for _, env := range msg.Envelopes {
			actual := envelopeTestUtils.UnmarshalUnsignedOriginatorEnvelope(
				t,
				env.UnsignedOriginatorEnvelope,
			)

			k := key{
				nodeID: int32(actual.OriginatorNodeId),
				seqID:  int64(actual.OriginatorSequenceId),
			}

			// Per-originator ordering must be strictly increasing in the *received stream*.
			if last, ok := lastSeqByNode[k.nodeID]; ok {
				require.Greater(
					t,
					k.seqID,
					last,
					"sequenceID must be strictly increasing for originator nodeID=%s", k.nodeID,
				)
			}
			lastSeqByNode[k.nodeID] = k.seqID

			// Envelope must be one we expected (order across originators doesn't matter).
			expRow, ok := expected[k]
			require.True(t, ok, "received unexpected update: nodeID=%s seqID=%d", k.nodeID, k.seqID)

			// Must not receive duplicates for the expected set.
			_, dup := seen[k]
			require.False(
				t,
				dup,
				"received duplicate update: nodeID=%s seqID=%d",
				k.nodeID,
				k.seqID,
			)

			// Validate contents match expected.
			require.EqualValues(t, expRow.OriginatorNodeID, actual.OriginatorNodeId)
			require.EqualValues(t, expRow.OriginatorSequenceID, actual.OriginatorSequenceId)
			require.Equal(t, expRow.OriginatorEnvelope, testutils.Marshal(t, env))

			seen[k] = struct{}{}

			if len(seen) == len(expected) {
				break
			}
		}
	}

	require.NoError(t, stream.Err())
	require.Len(t, seen, len(expected), "did not receive all expected updates")
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
	validateUpdates(t, stream, []int{})
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
	validateUpdates(t, stream1, []int{})
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

	// Note that streams don't return an error on establishing the connection.
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

	// Consume keepalive messages until stream closes.
	receivedMessages := 0
	for stream.Receive() {
		msg := stream.Msg()
		require.NotNil(t, msg, "keepalive message should not be nil")
		receivedMessages++
	}

	// Verify we received exactly one keepalive message.
	require.Equal(t, receivedMessages, 1, "should receive exactly one keepalive message")

	// Verify the stream closed with InvalidArgument error.
	err = stream.Err()
	require.Error(t, err)
	require.Equal(
		t,
		connect.CodeInvalidArgument,
		connect.CodeOf(err),
		"stream should close with InvalidArgument error",
	)
}

func generateEnvelopes(
	t *testing.T,
	originatorCount int,
	low int,
	high int,
	payerID int32,
) map[int32][]queries.InsertGatewayEnvelopeParams {
	t.Helper()

	const (
		topicCount = 5
	)

	out := make(map[int32][]queries.InsertGatewayEnvelopeParams)

	for i := range originatorCount {
		var (
			id    = int32(100 * (i + 1))
			topic = topic.NewTopic(
				topic.TopicKindGroupMessagesV1,
				[]byte(fmt.Sprintf("topic-%v", rand.Intn(topicCount))),
			)
		)

		n := low + rand.Intn(high-low)

		envs := make([]queries.InsertGatewayEnvelopeParams, n)
		for j := range n {
			// Sequence ID start at 1.
			seqID := int64(j + 1)

			oe := testutils.Marshal(
				t,
				envelopeTestUtils.CreateOriginatorEnvelopeWithTopic(
					t,
					uint32(id),
					uint64(seqID),
					topic.Bytes(),
				),
			)

			envs[j] = queries.InsertGatewayEnvelopeParams{
				OriginatorNodeID:     id,
				OriginatorSequenceID: seqID,
				Topic:                topic.Bytes(),
				PayerID:              db.NullInt32(payerID),
				OriginatorEnvelope:   oe,
			}
		}

		out[id] = envs
	}

	return out
}

func saveEnvelopes(
	t *testing.T,
	store *sql.DB,
	envelopes map[int32][]queries.InsertGatewayEnvelopeParams,
) {
	t.Helper()

	for _, nodeEnvelopes := range envelopes {
		testutils.InsertGatewayEnvelopes(t, store, nodeEnvelopes)
	}
}

func TestSubscribeVariableEnvelopesPerOriginator(t *testing.T) {
	var (
		server      = testUtilsApi.NewTestAPIServer(t)
		ctx, cancel = context.WithCancel(t.Context())
		payerID     = testutils.CreatePayer(t, server.DB)

		sourceEnvelopes = generateEnvelopes(t, 4, 50, 100, payerID)

		// For easier envelope lookup, use "<node-id>-<seq-id>" key.
		keyID = func(nodeID int32, seqID int64) string {
			return fmt.Sprintf("%v-%v", nodeID, seqID)
		}
	)
	defer cancel()

	// Check how many envelopes we have so we know how many to expect back.
	total := 0
	for id, env := range sourceEnvelopes {
		t.Logf("generated %v envelopes for nodeID %v", len(env), id)
		total += len(env)
	}

	// Subscribe to envelopes the node.
	req := &message_api.SubscribeEnvelopesRequest{
		Query: &message_api.EnvelopesQuery{
			LastSeen: nil,
		},
	}

	stream, err := server.ClientReplication.SubscribeEnvelopes(ctx, connect.NewRequest(req))
	require.NoError(t, err)

	// Insert envelopes which will be streamed.
	saveEnvelopes(t, server.DB, sourceEnvelopes)

	// Receive messages and do accounting.
	var (
		received_count = 0
		received       = make(map[string]struct{})
	)
	for received_count < total {

		ok := stream.Receive()
		if !ok {
			break
		}

		msg := stream.Msg()
		for _, env := range msg.Envelopes {
			received_count += 1

			decoded := envelopeTestUtils.UnmarshalUnsignedOriginatorEnvelope(
				t,
				env.UnsignedOriginatorEnvelope,
			)

			received[keyID(int32(decoded.OriginatorNodeId), int64(decoded.OriginatorSequenceId))] = struct{}{}
		}
	}

	cancel()

	err = stream.Err()
	require.Truef(
		t,
		err == nil || errors.Is(err, context.Canceled),
		"unexpected stream error: %s, received %v/%v envelopes",
		err,
		received_count,
		total,
	)

	require.Equal(t, total, received_count)

	// Accounting - verify that query returned everything.
	// Confirm simply that we got back all envelopes based on nodeID and seqID.
	sent := make(map[string]struct{})
	for _, envs := range sourceEnvelopes {
		for _, env := range envs {
			sent[keyID(env.OriginatorNodeID, env.OriginatorSequenceID)] = struct{}{}
		}
	}

	require.Equal(t, sent, received)
}
