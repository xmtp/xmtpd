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
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/api/message"
	"github.com/xmtp/xmtpd/pkg/constants"
	"github.com/xmtp/xmtpd/pkg/db"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/migrator"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/envelopes"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api"
	message_apiconnect "github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api/message_apiconnect"
	"github.com/xmtp/xmtpd/pkg/registry"
	"github.com/xmtp/xmtpd/pkg/testutils"
	testUtilsApi "github.com/xmtp/xmtpd/pkg/testutils/api"
	envelopeTestUtils "github.com/xmtp/xmtpd/pkg/testutils/envelopes"
	testregistry "github.com/xmtp/xmtpd/pkg/testutils/registry"
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
		nodes = []registry.Node{
			{NodeID: 100, IsCanonical: true},
			{NodeID: 200, IsCanonical: true},
		}
		suite   = testUtilsApi.NewTestAPIServer(t, testUtilsApi.WithRegistryNodes(nodes))
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
		for _, env := range msg.GetEnvelopes() {
			actual := envelopeTestUtils.UnmarshalUnsignedOriginatorEnvelope(
				t,
				env.GetUnsignedOriginatorEnvelope(),
			)

			k := key{
				nodeID: int32(actual.GetOriginatorNodeId()),
				seqID:  int64(actual.GetOriginatorSequenceId()),
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
			require.EqualValues(t, expRow.OriginatorNodeID, actual.GetOriginatorNodeId())
			require.EqualValues(t, expRow.OriginatorSequenceID, actual.GetOriginatorSequenceId())
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
	require.Equal(t, 1, receivedMessages, "should receive exactly one keepalive message")

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
	nodeIDs []uint32,
	low int,
	high int,
	payerID int32,
	topic *topic.Topic,
) map[int32][]queries.InsertGatewayEnvelopeParams {
	t.Helper()

	out := make(map[int32][]queries.InsertGatewayEnvelopeParams)

	for _, id := range nodeIDs {

		n := low + rand.Intn(high-low)

		envs := make([]queries.InsertGatewayEnvelopeParams, n)
		for i := range n {
			// Sequence IDs start at 1.
			seqID := int64(i + 1)

			oe := testutils.Marshal(
				t,
				envelopeTestUtils.CreateOriginatorEnvelopeWithTopic(
					t,
					uint32(id),
					uint64(seqID),
					topic.Bytes(),
				),
			)

			envs[i] = queries.InsertGatewayEnvelopeParams{
				OriginatorNodeID:     int32(id),
				OriginatorSequenceID: seqID,
				Topic:                topic.Bytes(),
				PayerID:              db.NullInt32(payerID),
				OriginatorEnvelope:   oe,
			}
		}

		out[int32(id)] = envs
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

func generateNodes(t *testing.T, n int) []registry.Node {
	out := make([]registry.Node, n)
	for i := range n {
		key, err := crypto.GenerateKey()
		require.NoError(t, err)

		// NOTE: Start node IDs from 1000 so we do not collide with nodes created elsewhere (starting at 100).
		node := testregistry.GetHealthyNode(uint32(1000 + 100*i))
		node.SigningKey = &key.PublicKey

		out[i] = node
	}

	return out
}

func nodeIDs(nodes []registry.Node) []uint32 {
	out := make([]uint32, len(nodes))
	for i, node := range nodes {
		out[i] = node.NodeID
	}
	return out
}

func TestSubscribeVariableEnvelopesPerOriginator(t *testing.T) {
	var (
		nodes       = generateNodes(t, 4)
		server      = testUtilsApi.NewTestAPIServer(t, testUtilsApi.WithRegistryNodes(nodes))
		ctx, cancel = context.WithCancel(t.Context())
		payerID     = testutils.CreatePayer(t, server.DB)
		subTopic    = topic.NewTopic(
			topic.TopicKindGroupMessagesV1,
			fmt.Appendf(nil, "generic-topic-%v", rand.Int()),
		)

		// Include messages not coming from our nodes.
		reservedOriginatorIDs = []uint32{
			constants.GroupMessageOriginatorID,
			constants.IdentityUpdateOriginatorID,
			migrator.GroupMessageOriginatorID,
			migrator.WelcomeMessageOriginatorID,
			migrator.KeyPackagesOriginatorID,
		}
		ids             = append(nodeIDs(nodes), reservedOriginatorIDs...)
		sourceEnvelopes = generateEnvelopes(t, ids, 50, 100, payerID, subTopic)

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

	// Subscribe to envelopes from all nodes.
	req := &message_api.SubscribeEnvelopesRequest{
		Query: &message_api.EnvelopesQuery{
			LastSeen: nil,
			Topics: [][]byte{
				subTopic.Bytes(),
			},
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

		for _, env := range msg.GetEnvelopes() {
			received_count += 1

			decoded := envelopeTestUtils.UnmarshalUnsignedOriginatorEnvelope(
				t,
				env.GetUnsignedOriginatorEnvelope(),
			)

			received[keyID(int32(decoded.GetOriginatorNodeId()), int64(decoded.GetOriginatorSequenceId()))] = struct{}{}
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

	t.Logf("processed %v envelopes", received_count)

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

// TestSubscribeCatchUpSkewedOriginators mimics a migrator client that subscribes
// to the migrator server plus the three migration originators (10, 11, 13).
// The subscription has to make sure that all messages from the migrator server are delivered.
func TestSubscribeCatchUpSkewedOriginators(t *testing.T) {
	var (
		// Use just above maxRequestedRows to test the case where the query returns less than maxRequestedRows.
		heavyMsgCount = 1001
		server        = testUtilsApi.NewTestAPIServer(t)
		payerID       = testutils.CreatePayer(t, server.DB)
		subTopic      = topic.NewTopic(
			topic.TopicKindGroupMessagesV1,
			fmt.Appendf(nil, "skewed-topic-%v", rand.Int()),
		)

		// Mimics a migrator client: own nodeID (100) + migration originators.
		originatorIDs = []uint32{
			100,
			migrator.GroupMessageOriginatorID,
			migrator.WelcomeMessageOriginatorID,
			migrator.KeyPackagesOriginatorID,
		}
		heavyOriginatorID = migrator.GroupMessageOriginatorID
	)

	// All messages go to originator 10 (group messages), the heaviest in practice.
	// The old query would have rows_per_originator = max(1000/4, 50) = 250, so the LATERAL subquery
	// returns at most 250 of 500 rows, total < 1000 → catchUp breaks.
	sourceEnvelopes := generateEnvelopes(
		t, []uint32{heavyOriginatorID}, heavyMsgCount, heavyMsgCount+1, payerID, subTopic,
	)

	// Populate the database.
	saveEnvelopes(t, server.DB, sourceEnvelopes)

	// Let the subscribeWorker's catch up.
	time.Sleep(4 * message.SubscribeWorkerPollTime)

	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()

	stream, err := server.ClientReplication.SubscribeEnvelopes(
		ctx,
		connect.NewRequest(&message_api.SubscribeEnvelopesRequest{
			Query: &message_api.EnvelopesQuery{
				OriginatorNodeIds: originatorIDs,
				LastSeen: &envelopes.Cursor{
					NodeIdToSequenceId: map[uint32]uint64{},
				},
			},
		}),
	)
	require.NoError(t, err)

	var (
		total    = len(sourceEnvelopes[int32(heavyOriginatorID)])
		received = make(map[int64]struct{}, total)
	)

	for len(received) < total {
		// If the stream is closed, means the subscribeWorker has caught up.
		// This shouldn't happen after we've saved all envelopes.
		if !stream.Receive() {
			break
		}

		for _, env := range stream.Msg().GetEnvelopes() {
			decoded := envelopeTestUtils.UnmarshalUnsignedOriginatorEnvelope(
				t,
				env.GetUnsignedOriginatorEnvelope(),
			)
			require.Equal(t, heavyOriginatorID, decoded.GetOriginatorNodeId())
			received[int64(decoded.GetOriginatorSequenceId())] = struct{}{}
		}
	}

	cancel()

	err = stream.Err()
	require.Truef(
		t,
		err == nil || errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded),
		"unexpected stream error: %s, received %v/%v envelopes",
		err,
		len(received),
		total,
	)

	require.Lenf(
		t,
		received,
		total,
		"catch-up must deliver all envelopes; LATERAL per-originator cap (%d for %d originators) causes premature pagination termination",
		max(1000/len(originatorIDs), 50),
		len(originatorIDs),
	)
}
