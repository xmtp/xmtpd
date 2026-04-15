package message_test

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"connectrpc.com/connect"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/constants"
	"github.com/xmtp/xmtpd/pkg/db"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/migrator"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/envelopes"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api/message_apiconnect"
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
	allRows = make([]queries.InsertGatewayEnvelopeV3Params, 0)
)

func setupTest(t *testing.T) *testUtilsApi.APIServerTestSuite {
	var (
		nodes = []registry.Node{
			{NodeID: 100, IsCanonical: true},
			{NodeID: 200, IsCanonical: true},
		}
		suite   = testUtilsApi.NewTestAPIServer(t, testUtilsApi.WithRegistryNodes(nodes))
		payerID = db.NullInt32(testutils.CreatePayer(t, suite.DB))
	)

	allRows = []queries.InsertGatewayEnvelopeV3Params{
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

	return suite
}

func insertInitialRows(t *testing.T, suite *testUtilsApi.APIServerTestSuite) {
	t.Helper()
	testutils.InsertGatewayEnvelopes(t, suite.DB, []queries.InsertGatewayEnvelopeV3Params{
		allRows[0], allRows[1],
	})
	// Wait until the subscribe worker has polled past the inserted rows so that
	// a subsequent subscription with LastSeen=nil won't see them.
	require.Eventually(t, func() bool {
		return suite.MessageService.DispatchedMet(db.VectorClock{100: 1, 200: 1})
	}, 500*time.Millisecond, 5*time.Millisecond)
}

func insertAdditionalRows(t *testing.T, store *sql.DB, notifyChan ...chan bool) {
	testutils.InsertGatewayEnvelopes(t, store, []queries.InsertGatewayEnvelopeV3Params{
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
	expected := make(map[key]queries.InsertGatewayEnvelopeV3Params, len(expectedIndices))
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

func TestSubscribeEnvelopesByTopic(t *testing.T) {
	suite := setupTest(t)
	insertInitialRows(t, suite)

	ctx := t.Context()
	stream, err := suite.ClientReplication.SubscribeEnvelopes(
		ctx,
		connect.NewRequest(&message_api.SubscribeEnvelopesRequest{
			Query: &message_api.EnvelopesQuery{
				Topics:   []db.Topic{topicA, topicC},
				LastSeen: nil,
			},
		}),
	)
	require.NoError(t, err)

	insertAdditionalRows(t, suite.DB)
	validateUpdates(t, stream, []int{4})
}

func TestSubscribeEnvelopesByOriginator(t *testing.T) {
	suite := setupTest(t)
	insertInitialRows(t, suite)

	ctx := t.Context()
	stream, err := suite.ClientReplication.SubscribeEnvelopes(
		ctx,
		connect.NewRequest(&message_api.SubscribeEnvelopesRequest{
			Query: &message_api.EnvelopesQuery{
				OriginatorNodeIds: []uint32{100, 300},
				LastSeen:          nil,
			},
		}),
	)
	require.NoError(t, err)

	insertAdditionalRows(t, suite.DB)
	validateUpdates(t, stream, []int{2, 4})
}

func TestSimultaneousSubscriptions(t *testing.T) {
	suite := setupTest(t)
	insertInitialRows(t, suite)

	ctx := t.Context()
	stream1, err := suite.ClientReplication.SubscribeEnvelopes(
		ctx,
		connect.NewRequest(&message_api.SubscribeEnvelopesRequest{
			Query: &message_api.EnvelopesQuery{},
		}),
	)
	require.NoError(t, err)

	stream2, err := suite.ClientReplication.SubscribeEnvelopes(
		ctx,
		connect.NewRequest(&message_api.SubscribeEnvelopesRequest{
			Query: &message_api.EnvelopesQuery{
				Topics:   []db.Topic{topicB},
				LastSeen: nil,
			},
		}),
	)
	require.NoError(t, err)

	stream3, err := suite.ClientReplication.SubscribeEnvelopes(
		ctx,
		connect.NewRequest(&message_api.SubscribeEnvelopesRequest{
			Query: &message_api.EnvelopesQuery{
				OriginatorNodeIds: []uint32{200},
				LastSeen:          nil,
			},
		}),
	)
	require.NoError(t, err)

	insertAdditionalRows(t, suite.DB)
	validateUpdates(t, stream1, []int{})
	validateUpdates(t, stream2, []int{2, 3})
	validateUpdates(t, stream3, []int{3})
}

func TestSubscribeEnvelopesFromCursor(t *testing.T) {
	suite := setupTest(t)
	insertInitialRows(t, suite)

	ctx := t.Context()
	stream, err := suite.ClientReplication.SubscribeEnvelopes(
		ctx,
		connect.NewRequest(&message_api.SubscribeEnvelopesRequest{
			Query: &message_api.EnvelopesQuery{
				Topics:   []db.Topic{topicA, topicC},
				LastSeen: &envelopes.Cursor{NodeIdToSequenceId: map[uint32]uint64{100: 1}},
			},
		}),
	)
	require.NoError(t, err)

	insertAdditionalRows(t, suite.DB)
	validateUpdates(t, stream, []int{1, 4})
}

func TestSubscribeEnvelopesFromEmptyCursor(t *testing.T) {
	suite := setupTest(t)
	insertInitialRows(t, suite)

	ctx := t.Context()
	stream, err := suite.ClientReplication.SubscribeEnvelopes(
		ctx,
		connect.NewRequest(&message_api.SubscribeEnvelopesRequest{
			Query: &message_api.EnvelopesQuery{
				Topics:   []db.Topic{topicA, topicC},
				LastSeen: &envelopes.Cursor{NodeIdToSequenceId: map[uint32]uint64{}},
			},
		}),
	)
	require.NoError(t, err)

	insertAdditionalRows(t, suite.DB)
	validateUpdates(t, stream, []int{0, 1, 4})
}

func TestSubscribeEnvelopesInvalidRequest(t *testing.T) {
	var (
		suite = setupTest(t)
		ctx   = t.Context()
		tests = []struct {
			name  string
			query *message_api.EnvelopesQuery
		}{
			{
				name: "no filter",
				// Neither topics nor originators set.
				query: &message_api.EnvelopesQuery{
					LastSeen: nil,
				},
			},
			{
				name: "incompatible filters",
				// Both topics and originators set.
				query: &message_api.EnvelopesQuery{
					Topics:            []db.Topic{topicA},
					OriginatorNodeIds: []uint32{1},
					LastSeen:          nil,
				},
			},
		}
	)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			stream, err := suite.ClientReplication.SubscribeEnvelopes(
				ctx,
				connect.NewRequest(&message_api.SubscribeEnvelopesRequest{
					Query: test.query,
				}),
			)
			require.NoError(t, err)

			_ = stream.Receive()

			streamErr := stream.Err()
			require.Equal(t, connect.CodeInvalidArgument, connect.CodeOf(streamErr))
		})
	}
}

func generateEnvelopes(
	t *testing.T,
	nodeIDs []uint32,
	low int,
	high int,
	payerID int32,
	topic *topic.Topic,
) map[int32][]queries.InsertGatewayEnvelopeV3Params {
	t.Helper()

	out := make(map[int32][]queries.InsertGatewayEnvelopeV3Params)

	for _, id := range nodeIDs {

		n := low
		// Add variance if requested.
		if high > low {
			n += rand.Intn(high - low)
		}

		envs := make([]queries.InsertGatewayEnvelopeV3Params, n)
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

			envs[i] = queries.InsertGatewayEnvelopeV3Params{
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

func flattenEnvelopeMap(
	m map[int32][]queries.InsertGatewayEnvelopeV3Params,
) []queries.InsertGatewayEnvelopeV3Params {
	var out []queries.InsertGatewayEnvelopeV3Params
	for _, list := range m {
		out = append(out, list...)
	}

	return out
}

func saveEnvelopes(
	t *testing.T,
	store *sql.DB,
	envelopes map[int32][]queries.InsertGatewayEnvelopeV3Params,
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
		sourceEnvelopes = generateEnvelopes(t, ids, 10, 30, payerID, subTopic)

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
		receivedCount atomic.Int64
		received      = make(map[string]struct{})
		streamWG      sync.WaitGroup
	)

	streamWG.Go(func() {
		for receivedCount.Load() < int64(total) {
			if ok := stream.Receive(); !ok {
				break
			}

			msg := stream.Msg()

			for _, env := range msg.GetEnvelopes() {
				receivedCount.Add(1)

				decoded := envelopeTestUtils.UnmarshalUnsignedOriginatorEnvelope(
					t,
					env.GetUnsignedOriginatorEnvelope(),
				)

				received[keyID(int32(decoded.GetOriginatorNodeId()), int64(decoded.GetOriginatorSequenceId()))] = struct{}{}
			}
		}

		cancel()
	})

	require.Eventually(t, func() bool {
		return receivedCount.Load() >= int64(total)
	}, 10*time.Second, 100*time.Millisecond, "not all envelopes received")

	streamWG.Wait()

	require.Equal(t, int64(total), receivedCount.Load())

	err = stream.Err()
	require.Truef(
		t,
		err == nil || errors.Is(err, context.Canceled),
		"unexpected stream error: %s, received %v/%v envelopes",
		err,
		receivedCount.Load(),
		total,
	)

	t.Logf("processed %v envelopes", receivedCount.Load())

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

	// Block until the subscribeWorker has polled past the last inserted sequence ID.
	require.Eventually(t, func() bool {
		return server.MessageService.DispatchedMet(
			db.VectorClock{heavyOriginatorID: uint64(heavyMsgCount)},
		)
	}, 5*time.Second, 5*time.Millisecond)

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

func TestSubscribeAll(t *testing.T) {
	var (
		nodes       = generateNodes(t, 2)
		nodeIDs     = nodeIDs(nodes)
		server      = testUtilsApi.NewTestAPIServer(t, testUtilsApi.WithRegistryNodes(nodes))
		ctx, cancel = context.WithCancel(t.Context())
		payerID     = testutils.CreatePayer(t, server.DB)
		subTopic    = topic.NewTopic(
			topic.TopicKindGroupMessagesV1,
			fmt.Appendf(nil, "generic-topic-%v", rand.Int()),
		)

		minEnvelopes = 10
		maxEnvelopes = 20

		envelopeList = flattenEnvelopeMap(
			generateEnvelopes(
				t,
				nodeIDs,
				minEnvelopes,
				maxEnvelopes,
				payerID,
				subTopic,
			))
		total = len(envelopeList)
	)
	defer cancel()

	t.Logf("generated total %v envelopes from %v nodes", total, len(nodeIDs))

	// Start a subscriber stream.
	req := &message_api.SubscribeAllEnvelopesRequest{}
	stream, err := server.ClientNotification.SubscribeAllEnvelopes(ctx, connect.NewRequest(req))
	require.NoError(t, err)

	var (
		received atomic.Int64
		streamWG sync.WaitGroup
	)

	streamWG.Go(func() {
		for received.Load() < int64(total) {
			ok := stream.Receive()
			if !ok {
				break
			}

			n := len(stream.Msg().GetEnvelopes())
			t.Logf("stream produced %v envelopes", n)

			received.Add(int64(n))
		}

		cancel()
	})

	// Wait until the server has registered the listener before inserting, so
	// no inserts race the listener registration.
	require.Eventually(t, func() bool {
		return server.MessageService.GlobalListenerCount() >= 1
	}, 5*time.Second, 5*time.Millisecond)

	for _, env := range envelopeList {
		testutils.InsertGatewayEnvelopes(t, server.DB, []queries.InsertGatewayEnvelopeV3Params{env})
	}

	streamWG.Wait()

	require.Equal(t, int64(total), received.Load())
}

func readOriginatorsStream(
	t *testing.T,
	stream *connect.ServerStreamForClient[message_api.SubscribeOriginatorsResponse],
	n int,
) []*envelopes.OriginatorEnvelope {
	t.Helper()
	var received []*envelopes.OriginatorEnvelope
	for len(received) < n {
		if !stream.Receive() {
			break
		}
		received = append(received, stream.Msg().GetEnvelopes().GetEnvelopes()...)
	}
	require.NoError(t, stream.Err())
	return received
}

func TestSubscribeOriginators(t *testing.T) {
	suite := setupTest(t)
	insertInitialRows(t, suite)

	ctx := t.Context()
	stream, err := suite.ClientReplication.SubscribeOriginators(
		ctx,
		connect.NewRequest(&message_api.SubscribeOriginatorsRequest{
			Filter: &message_api.SubscribeOriginatorsRequest_OriginatorFilter{
				OriginatorNodeIds: []uint32{100},
				LastSeen:          &envelopes.Cursor{NodeIdToSequenceId: map[uint32]uint64{}},
			},
		}),
	)
	require.NoError(t, err)

	// Expect initial keepalive.
	ok := stream.Receive()
	require.True(t, ok, "expected initial keepalive")
	require.Empty(t, stream.Msg().GetEnvelopes().GetEnvelopes())

	// Catch-up: the initial rows for originator 100 (seq 1) should arrive.
	catchUp := readOriginatorsStream(t, stream, 1)
	require.Len(t, catchUp, 1)

	insertAdditionalRows(t, suite.DB)

	// Live: new rows for originator 100 (seq 2, seq 3) should arrive.
	received := readOriginatorsStream(t, stream, 2)
	require.Len(t, received, 2)
	for _, env := range received {
		decoded := envelopeTestUtils.UnmarshalUnsignedOriginatorEnvelope(
			t,
			env.GetUnsignedOriginatorEnvelope(),
		)
		require.EqualValues(t, 100, decoded.GetOriginatorNodeId())
	}
}

func TestSubscribeOriginators_NilLastSeen(t *testing.T) {
	suite := setupTest(t)
	ctx := t.Context()

	stream, err := suite.ClientReplication.SubscribeOriginators(
		ctx,
		connect.NewRequest(&message_api.SubscribeOriginatorsRequest{
			Filter: &message_api.SubscribeOriginatorsRequest_OriginatorFilter{
				OriginatorNodeIds: []uint32{100},
				LastSeen:          nil,
			},
		}),
	)
	require.NoError(t, err)

	_ = stream.Receive()
	require.Equal(t, connect.CodeInvalidArgument, connect.CodeOf(stream.Err()))
}

func TestSubscribeOriginators_NilFilter(t *testing.T) {
	suite := setupTest(t)
	ctx := t.Context()

	stream, err := suite.ClientReplication.SubscribeOriginators(
		ctx,
		connect.NewRequest(&message_api.SubscribeOriginatorsRequest{
			Filter: nil,
		}),
	)
	require.NoError(t, err)

	_ = stream.Receive()
	require.Equal(t, connect.CodeInvalidArgument, connect.CodeOf(stream.Err()))
}

func TestSubscribeOriginators_EmptyOriginatorIds(t *testing.T) {
	suite := setupTest(t)
	ctx := t.Context()

	stream, err := suite.ClientReplication.SubscribeOriginators(
		ctx,
		connect.NewRequest(&message_api.SubscribeOriginatorsRequest{
			Filter: &message_api.SubscribeOriginatorsRequest_OriginatorFilter{
				OriginatorNodeIds: []uint32{},
			},
		}),
	)
	require.NoError(t, err)

	_ = stream.Receive()
	require.Equal(t, connect.CodeInvalidArgument, connect.CodeOf(stream.Err()))
}

func TestSubscribeOriginators_FromCursor(t *testing.T) {
	suite := setupTest(t)
	insertInitialRows(t, suite)

	ctx := t.Context()
	stream, err := suite.ClientReplication.SubscribeOriginators(
		ctx,
		connect.NewRequest(&message_api.SubscribeOriginatorsRequest{
			Filter: &message_api.SubscribeOriginatorsRequest_OriginatorFilter{
				OriginatorNodeIds: []uint32{100},
				LastSeen: &envelopes.Cursor{
					NodeIdToSequenceId: map[uint32]uint64{100: 1},
				},
			},
		}),
	)
	require.NoError(t, err)

	// Expect initial keepalive.
	ok := stream.Receive()
	require.True(t, ok, "expected initial keepalive")

	insertAdditionalRows(t, suite.DB)

	// allRows[2] = originator 100 seq 2, allRows[4] = originator 100 seq 3 (both past cursor 1).
	received := readOriginatorsStream(t, stream, 2)
	require.Len(t, received, 2)
	for _, env := range received {
		decoded := envelopeTestUtils.UnmarshalUnsignedOriginatorEnvelope(
			t,
			env.GetUnsignedOriginatorEnvelope(),
		)
		require.EqualValues(t, 100, decoded.GetOriginatorNodeId())
		require.Greater(t, decoded.GetOriginatorSequenceId(), uint64(1))
	}
}

func TestSubscribeOriginators_MultipleOriginators(t *testing.T) {
	suite := setupTest(t)
	insertInitialRows(t, suite)

	ctx := t.Context()
	stream, err := suite.ClientReplication.SubscribeOriginators(
		ctx,
		connect.NewRequest(&message_api.SubscribeOriginatorsRequest{
			Filter: &message_api.SubscribeOriginatorsRequest_OriginatorFilter{
				OriginatorNodeIds: []uint32{100, 200},
				LastSeen:          &envelopes.Cursor{NodeIdToSequenceId: map[uint32]uint64{}},
			},
		}),
	)
	require.NoError(t, err)

	// Expect initial keepalive.
	ok := stream.Receive()
	require.True(t, ok, "expected initial keepalive")

	// Catch-up: initial rows for both originators (100 seq 1, 200 seq 1).
	catchUp := readOriginatorsStream(t, stream, 2)
	require.Len(t, catchUp, 2)

	seenOriginators := make(map[uint32]bool)
	for _, env := range catchUp {
		decoded := envelopeTestUtils.UnmarshalUnsignedOriginatorEnvelope(
			t,
			env.GetUnsignedOriginatorEnvelope(),
		)
		seenOriginators[decoded.GetOriginatorNodeId()] = true
	}
	require.True(t, seenOriginators[100], "expected envelope from originator 100")
	require.True(t, seenOriginators[200], "expected envelope from originator 200")

	insertAdditionalRows(t, suite.DB)

	// Live: new rows for both originators (100 seq 2+3, 200 seq 2) = 3 envelopes.
	received := readOriginatorsStream(t, stream, 3)
	require.Len(t, received, 3)

	seenOriginators = make(map[uint32]bool)
	for _, env := range received {
		decoded := envelopeTestUtils.UnmarshalUnsignedOriginatorEnvelope(
			t,
			env.GetUnsignedOriginatorEnvelope(),
		)
		seenOriginators[decoded.GetOriginatorNodeId()] = true
	}
	require.True(t, seenOriginators[100], "expected live envelope from originator 100")
	require.True(t, seenOriginators[200], "expected live envelope from originator 200")
}

// ---------------------------------------------------------------------------
// Parity tests: run identical scenarios against both SubscribeEnvelopes
// (originator filter) and SubscribeOriginators to verify behavioral equivalence.
// ---------------------------------------------------------------------------

// envelopeStreamAdapter abstracts over SubscribeEnvelopes and SubscribeOriginators
// stream types so the same test logic can drive both.
type envelopeStreamAdapter struct {
	receiveFn   func() bool
	envelopesFn func() []*envelopes.OriginatorEnvelope
	errFn       func() error
}

func (a *envelopeStreamAdapter) Receive() bool { return a.receiveFn() }

func (a *envelopeStreamAdapter) Envelopes() []*envelopes.OriginatorEnvelope { return a.envelopesFn() }
func (a *envelopeStreamAdapter) Err() error                                 { return a.errFn() }

// subscribeByOriginatorsFn creates a subscription stream filtered by originator IDs
// with a non-nil cursor (empty map = catch up from the beginning).
type subscribeByOriginatorsFn func(
	ctx context.Context,
	client message_apiconnect.ReplicationApiClient,
	originatorIDs []uint32,
	cursor map[uint32]uint64,
) (*envelopeStreamAdapter, error)

func subscribeEnvelopesAdapter(
	ctx context.Context,
	client message_apiconnect.ReplicationApiClient,
	originatorIDs []uint32,
	cursor map[uint32]uint64,
) (*envelopeStreamAdapter, error) {
	stream, err := client.SubscribeEnvelopes(
		ctx,
		connect.NewRequest(&message_api.SubscribeEnvelopesRequest{
			Query: &message_api.EnvelopesQuery{
				OriginatorNodeIds: originatorIDs,
				LastSeen:          &envelopes.Cursor{NodeIdToSequenceId: cursor},
			},
		}),
	)
	if err != nil {
		return nil, err
	}
	return &envelopeStreamAdapter{
		receiveFn:   stream.Receive,
		envelopesFn: func() []*envelopes.OriginatorEnvelope { return stream.Msg().GetEnvelopes() },
		errFn:       stream.Err,
	}, nil
}

func subscribeOriginatorsAdapter(
	ctx context.Context,
	client message_apiconnect.ReplicationApiClient,
	originatorIDs []uint32,
	cursor map[uint32]uint64,
) (*envelopeStreamAdapter, error) {
	stream, err := client.SubscribeOriginators(
		ctx,
		connect.NewRequest(&message_api.SubscribeOriginatorsRequest{
			Filter: &message_api.SubscribeOriginatorsRequest_OriginatorFilter{
				OriginatorNodeIds: originatorIDs,
				LastSeen:          &envelopes.Cursor{NodeIdToSequenceId: cursor},
			},
		}),
	)
	if err != nil {
		return nil, err
	}
	return &envelopeStreamAdapter{
		receiveFn:   stream.Receive,
		envelopesFn: func() []*envelopes.OriginatorEnvelope { return stream.Msg().GetEnvelopes().GetEnvelopes() },
		errFn:       stream.Err,
	}, nil
}

var originatorSubscribeBackends = []struct {
	name      string
	subscribe subscribeByOriginatorsFn
}{
	{"SubscribeEnvelopes", subscribeEnvelopesAdapter},
	{"SubscribeOriginators", subscribeOriginatorsAdapter},
}

func readAdaptedStream(
	t *testing.T,
	stream *envelopeStreamAdapter,
	n int,
) []*envelopes.OriginatorEnvelope {
	t.Helper()
	var received []*envelopes.OriginatorEnvelope
	for len(received) < n {
		if !stream.Receive() {
			break
		}
		received = append(received, stream.Envelopes()...)
	}
	require.NoError(t, stream.Err())
	return received
}

// consumeKeepalive reads and discards the initial keepalive message that both
// SubscribeEnvelopes and SubscribeOriginators send on stream open.
func consumeKeepalive(t *testing.T, stream *envelopeStreamAdapter) {
	t.Helper()
	ok := stream.Receive()
	require.True(t, ok, "expected initial keepalive")
}

// Test 1: Catch-up pagination with >maxRequestedRows envelopes on a single
// originator. Verifies the pagination loop delivers every envelope.
func TestOriginatorParity_SkewedPagination(t *testing.T) {
	for _, backend := range originatorSubscribeBackends {
		t.Run(backend.name, func(t *testing.T) {
			var (
				heavyMsgCount     = 1001 // just above maxRequestedRows (1000)
				heavyOriginatorID = uint32(100)
				server            = testUtilsApi.NewTestAPIServer(t)
				payerID           = testutils.CreatePayer(t, server.DB)
				subTopic          = topic.NewTopic(
					topic.TopicKindGroupMessagesV1,
					fmt.Appendf(nil, "skewed-parity-%v", rand.Int()),
				)
			)

			sourceEnvelopes := generateEnvelopes(
				t, []uint32{heavyOriginatorID},
				heavyMsgCount, heavyMsgCount+1, payerID, subTopic,
			)
			saveEnvelopes(t, server.DB, sourceEnvelopes)

			require.Eventually(t, func() bool {
				return server.MessageService.DispatchedMet(
					db.VectorClock{heavyOriginatorID: uint64(heavyMsgCount)},
				)
			}, 5*time.Second, 5*time.Millisecond)

			ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
			defer cancel()

			stream, err := backend.subscribe(
				ctx, server.ClientReplication,
				[]uint32{heavyOriginatorID},
				map[uint32]uint64{},
			)
			require.NoError(t, err)
			consumeKeepalive(t, stream)

			total := len(sourceEnvelopes[int32(heavyOriginatorID)])
			received := make(map[uint64]struct{}, total)

			for len(received) < total {
				if !stream.Receive() {
					break
				}
				for _, env := range stream.Envelopes() {
					decoded := envelopeTestUtils.UnmarshalUnsignedOriginatorEnvelope(
						t, env.GetUnsignedOriginatorEnvelope(),
					)
					require.Equal(t, heavyOriginatorID, decoded.GetOriginatorNodeId())
					received[decoded.GetOriginatorSequenceId()] = struct{}{}
				}
			}

			cancel()

			err = stream.Err()
			require.Truef(
				t,
				err == nil || errors.Is(err, context.Canceled) ||
					errors.Is(err, context.DeadlineExceeded),
				"unexpected stream error: %s, received %v/%v envelopes",
				err,
				len(received),
				total,
			)
			require.Lenf(t, received, total,
				"catch-up must deliver all %d envelopes", total,
			)
		})
	}
}

// Test 2: Per-originator ordering. Sequence IDs within each originator must
// arrive in strictly increasing order across both catch-up and live phases.
func TestOriginatorParity_Ordering(t *testing.T) {
	for _, backend := range originatorSubscribeBackends {
		t.Run(backend.name, func(t *testing.T) {
			suite := setupTest(t)
			insertInitialRows(t, suite)

			ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
			defer cancel()

			stream, err := backend.subscribe(
				ctx, suite.ClientReplication,
				[]uint32{100, 200},
				map[uint32]uint64{},
			)
			require.NoError(t, err)
			consumeKeepalive(t, stream)

			insertAdditionalRows(t, suite.DB)

			// Expect all 5 rows: catch-up (100:1, 200:1) + live (100:2, 200:2, 100:3).
			const total = 5
			lastSeqByNode := make(map[uint32]uint64)
			seen := make(map[string]struct{}, total)

			for len(seen) < total {
				if !stream.Receive() {
					break
				}
				for _, env := range stream.Envelopes() {
					decoded := envelopeTestUtils.UnmarshalUnsignedOriginatorEnvelope(
						t, env.GetUnsignedOriginatorEnvelope(),
					)
					origID := decoded.GetOriginatorNodeId()
					seqID := decoded.GetOriginatorSequenceId()

					if last, ok := lastSeqByNode[origID]; ok {
						require.Greater(t, seqID, last,
							"sequence ID must be strictly increasing for originator %d", origID,
						)
					}
					lastSeqByNode[origID] = seqID

					key := fmt.Sprintf("%d-%d", origID, seqID)
					require.NotContains(t, seen, key, "duplicate envelope %s", key)
					seen[key] = struct{}{}
				}
			}

			require.NoError(t, stream.Err())
			require.Len(t, seen, total)
		})
	}
}

// Test 3: Two concurrent streams with different originator filters must be
// fully isolated — each receives only its filtered originator's envelopes.
func TestOriginatorParity_SimultaneousStreams(t *testing.T) {
	for _, backend := range originatorSubscribeBackends {
		t.Run(backend.name, func(t *testing.T) {
			suite := setupTest(t)
			insertInitialRows(t, suite)

			ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
			defer cancel()

			stream100, err := backend.subscribe(
				ctx, suite.ClientReplication,
				[]uint32{100},
				map[uint32]uint64{},
			)
			require.NoError(t, err)
			consumeKeepalive(t, stream100)

			stream200, err := backend.subscribe(
				ctx, suite.ClientReplication,
				[]uint32{200},
				map[uint32]uint64{},
			)
			require.NoError(t, err)
			consumeKeepalive(t, stream200)

			insertAdditionalRows(t, suite.DB)

			// stream100: catch-up (100:1) + live (100:2, 100:3) = 3 envelopes.
			received100 := readAdaptedStream(t, stream100, 3)
			require.Len(t, received100, 3)
			for _, env := range received100 {
				decoded := envelopeTestUtils.UnmarshalUnsignedOriginatorEnvelope(
					t, env.GetUnsignedOriginatorEnvelope(),
				)
				require.EqualValues(t, 100, decoded.GetOriginatorNodeId(),
					"stream100 must only contain originator 100",
				)
			}

			// stream200: catch-up (200:1) + live (200:2) = 2 envelopes.
			received200 := readAdaptedStream(t, stream200, 2)
			require.Len(t, received200, 2)
			for _, env := range received200 {
				decoded := envelopeTestUtils.UnmarshalUnsignedOriginatorEnvelope(
					t, env.GetUnsignedOriginatorEnvelope(),
				)
				require.EqualValues(t, 200, decoded.GetOriginatorNodeId(),
					"stream200 must only contain originator 200",
				)
			}
		})
	}
}

// Test 4: Multiple originators with random envelope counts (10-30 each).
// Verifies every envelope is delivered exactly once.
func TestOriginatorParity_VariableVolume(t *testing.T) {
	for _, backend := range originatorSubscribeBackends {
		t.Run(backend.name, func(t *testing.T) {
			var (
				nodes   = generateNodes(t, 4)
				ids     = nodeIDs(nodes)
				server  = testUtilsApi.NewTestAPIServer(t, testUtilsApi.WithRegistryNodes(nodes))
				payerID = testutils.CreatePayer(t, server.DB)
				sub     = topic.NewTopic(
					topic.TopicKindGroupMessagesV1,
					fmt.Appendf(nil, "variable-parity-%v", rand.Int()),
				)
				sourceEnvelopes = generateEnvelopes(t, ids, 10, 30, payerID, sub)
			)

			total := 0
			for id, envs := range sourceEnvelopes {
				t.Logf("generated %d envelopes for originator %d", len(envs), id)
				total += len(envs)
			}

			saveEnvelopes(t, server.DB, sourceEnvelopes)

			expectedVC := make(db.VectorClock, len(sourceEnvelopes))
			for nodeID, envs := range sourceEnvelopes {
				expectedVC[uint32(nodeID)] = uint64(len(envs))
			}
			require.Eventually(t, func() bool {
				return server.MessageService.DispatchedMet(expectedVC)
			}, 5*time.Second, 5*time.Millisecond)

			ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
			defer cancel()

			stream, err := backend.subscribe(
				ctx, server.ClientReplication, ids, map[uint32]uint64{},
			)
			require.NoError(t, err)
			consumeKeepalive(t, stream)

			keyID := func(nodeID uint32, seqID uint64) string {
				return fmt.Sprintf("%d-%d", nodeID, seqID)
			}

			received := make(map[string]struct{}, total)
			for len(received) < total {
				if !stream.Receive() {
					break
				}
				for _, env := range stream.Envelopes() {
					decoded := envelopeTestUtils.UnmarshalUnsignedOriginatorEnvelope(
						t, env.GetUnsignedOriginatorEnvelope(),
					)
					k := keyID(decoded.GetOriginatorNodeId(), decoded.GetOriginatorSequenceId())
					require.NotContains(t, received, k, "duplicate envelope %s", k)
					received[k] = struct{}{}
				}
			}

			cancel()

			err = stream.Err()
			require.Truef(t,
				err == nil || errors.Is(err, context.Canceled),
				"unexpected stream error: %s, received %v/%v", err, len(received), total,
			)

			sent := make(map[string]struct{}, total)
			for _, envs := range sourceEnvelopes {
				for _, env := range envs {
					sent[keyID(uint32(env.OriginatorNodeID), uint64(env.OriginatorSequenceID))] = struct{}{}
				}
			}
			require.Equal(t, sent, received)
		})
	}
}

func TestSubscribeAll_StreamsOnlyNewMessages(t *testing.T) {
	var (
		nodes       = generateNodes(t, 2)
		nodeIDs     = nodeIDs(nodes)
		server      = testUtilsApi.NewTestAPIServer(t, testUtilsApi.WithRegistryNodes(nodes))
		ctx, cancel = context.WithCancel(t.Context())
		payerID     = testutils.CreatePayer(t, server.DB)
		subTopic    = topic.NewTopic(
			topic.TopicKindGroupMessagesV1,
			fmt.Appendf(nil, "generic-topic-%v", rand.Int()),
		)
	)

	// Envelope data.
	//
	// generateEnvelopes returns `perNode` envelopes *per* originator, so with
	// 2 nodes we get 2 * perNode total envelopes. We only insert the first
	// `initialBatchSize` as pre-seeds and the next `streamSize` as the new
	// batch the stream should observe; the remainder are generated but never
	// inserted. Slicing by exact bounds (rather than [initialBatchSize:])
	// keeps the post-subscribe insert count to exactly streamSize — otherwise
	// the stream races past streamSize and the final equality check fails.
	var (
		initialBatchSize = 5
		streamSize       = 5
		perNode          = initialBatchSize + streamSize

		sourceEnvelopes = flattenEnvelopeMap(
			generateEnvelopes(
				t,
				nodeIDs,
				perNode,
				perNode, // Exactly perNode per node.
				payerID,
				subTopic,
			))

		initialBatch = sourceEnvelopes[:initialBatchSize]
		streamBatch  = sourceEnvelopes[initialBatchSize : initialBatchSize+streamSize]
	)
	defer cancel()

	// Pre-seed envelopes in the DB. These should NOT get picked up by the
	// stream because the subscribe worker marks them as known before the
	// stream subscription is registered.
	for _, env := range initialBatch {
		testutils.InsertGatewayEnvelopes(t, server.DB, []queries.InsertGatewayEnvelopeV3Params{env})
	}

	// Block until the subscribe worker has polled past every pre-seeded row.
	preSeedVC := make(db.VectorClock)
	for _, env := range initialBatch {
		nodeID := uint32(env.OriginatorNodeID)
		seq := uint64(env.OriginatorSequenceID)
		if cur, ok := preSeedVC[nodeID]; !ok || seq > cur {
			preSeedVC[nodeID] = seq
		}
	}
	require.Eventually(t, func() bool {
		return server.MessageService.DispatchedMet(preSeedVC)
	}, 5*time.Second, 5*time.Millisecond)

	// Start a subscriber stream.
	req := &message_api.SubscribeAllEnvelopesRequest{}
	stream, err := server.ClientNotification.SubscribeAllEnvelopes(ctx, connect.NewRequest(req))
	require.NoError(t, err)

	var (
		received atomic.Int64
		streamWG sync.WaitGroup
	)

	streamWG.Go(func() {
		for received.Load() < int64(streamSize) {
			ok := stream.Receive()
			if !ok {
				break
			}

			n := len(stream.Msg().GetEnvelopes())
			t.Logf("stream produced %v envelopes", n)

			received.Add(int64(n))
		}

		cancel()
	})

	// Wait until the server has registered the listener before inserting, so
	// no inserts race the listener registration.
	require.Eventually(t, func() bool {
		return server.MessageService.GlobalListenerCount() >= 1
	}, 5*time.Second, 5*time.Millisecond)

	for _, env := range streamBatch {
		testutils.InsertGatewayEnvelopes(t, server.DB, []queries.InsertGatewayEnvelopeV3Params{env})
	}

	streamWG.Wait()

	require.Equal(t, int64(streamSize), received.Load())
}
