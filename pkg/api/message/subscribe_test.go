package message_test

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"slices"
	"sync"
	"sync/atomic"
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
	ctx, cancel := context.WithTimeout(t.Context(), 500*time.Millisecond)
	defer cancel()
	require.NoError(t, suite.MessageService.AwaitCursor(ctx, db.VectorClock{100: 1, 200: 1}))
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

		insertDelay  = 100 * time.Millisecond
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
		received = 0
		streamWG sync.WaitGroup
	)

	streamWG.Go(func() {
		for received < total {
			ok := stream.Receive()
			if !ok {
				break
			}

			n := len(stream.Msg().GetEnvelopes())
			t.Logf("stream produced %v envelopes", n)

			received += n
		}

		cancel()
	})

	// Wait a bit - then start inserting envelopes. Make sure these are streamed.
	time.Sleep(insertDelay)

	for _, env := range envelopeList {
		testutils.InsertGatewayEnvelopes(t, server.DB, []queries.InsertGatewayEnvelopeV3Params{env})
		time.Sleep(insertDelay)
	}

	streamWG.Wait()

	require.Equal(t, total, received)
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

			awaitCtx, awaitCancel := context.WithTimeout(t.Context(), 5*time.Second)
			defer awaitCancel()
			require.NoError(t, server.MessageService.AwaitCursor(
				awaitCtx, db.VectorClock{heavyOriginatorID: uint64(heavyMsgCount)},
			))

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
			awaitCtx, awaitCancel := context.WithTimeout(t.Context(), 5*time.Second)
			defer awaitCancel()
			require.NoError(t, server.MessageService.AwaitCursor(awaitCtx, expectedVC))

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

		insertDelay = 100 * time.Millisecond
	)

	// Envelope data.
	var (
		initialBatchSize = 5
		streamSize       = 5
		totalMessages    = initialBatchSize + streamSize

		sourceEnvelopes = flattenEnvelopeMap(
			generateEnvelopes(
				t,
				nodeIDs,
				totalMessages,
				totalMessages, // Let's get exactly N messages.
				payerID,
				subTopic,
			))

		initialBatch = sourceEnvelopes[:initialBatchSize]
		streamBatch  = sourceEnvelopes[initialBatchSize:]
	)
	defer cancel()

	// Pre-seed envelopes in the DB.
	// These should NOT get picked up by the stream.
	for _, env := range initialBatch {
		testutils.InsertGatewayEnvelopes(t, server.DB, []queries.InsertGatewayEnvelopeV3Params{env})
	}

	// Add a delay so the subscribe worker picks pre-seeded envelopes as known before the streaming started.
	time.Sleep(insertDelay)

	// Start a subscriber stream.
	req := &message_api.SubscribeAllEnvelopesRequest{}
	stream, err := server.ClientNotification.SubscribeAllEnvelopes(ctx, connect.NewRequest(req))
	require.NoError(t, err)

	var (
		received = 0
		streamWG sync.WaitGroup
	)

	streamWG.Go(func() {
		for received < streamSize {
			ok := stream.Receive()
			if !ok {
				break
			}

			n := len(stream.Msg().GetEnvelopes())
			t.Logf("stream produced %v envelopes", n)

			received += n
		}

		cancel()
	})

	// Wait a bit - then start inserting envelopes. These should in fact be streamed.
	time.Sleep(insertDelay)

	for _, env := range streamBatch {
		testutils.InsertGatewayEnvelopes(t, server.DB, []queries.InsertGatewayEnvelopeV3Params{env})
		time.Sleep(insertDelay)
	}

	streamWG.Wait()

	require.Equal(t, streamSize, received)
}

// ---- XIP-83 bidirectional Subscribe (QueryApi) ----

func subMutate(
	mutateID uint64,
	historyOnly bool,
	adds []*message_api.SubscribeRequest_V1_Mutate_Subscription,
	removes [][]byte,
) *message_api.SubscribeRequest {
	return &message_api.SubscribeRequest{
		Version: &message_api.SubscribeRequest_V1_{
			V1: &message_api.SubscribeRequest_V1{
				Request: &message_api.SubscribeRequest_V1_Mutate_{
					Mutate: &message_api.SubscribeRequest_V1_Mutate{
						MutateId:    mutateID,
						HistoryOnly: historyOnly,
						Adds:        adds,
						Removes:     removes,
					},
				},
			},
		},
	}
}

func addSub(
	topicBytes []byte,
	cursor map[uint32]uint64,
) *message_api.SubscribeRequest_V1_Mutate_Subscription {
	sub := &message_api.SubscribeRequest_V1_Mutate_Subscription{Topic: topicBytes}
	if cursor != nil {
		sub.LastSeen = &envelopes.Cursor{NodeIdToSequenceId: cursor}
	}
	return sub
}

// bidiReader drains a Subscribe stream on a background goroutine into a thread-safe buffer, so a
// test can poll for expected frames with require.Eventually. It never answers pings (a test that
// needs liveness reaping relies on that).
type bidiReader struct {
	mu     sync.Mutex
	frames []*message_api.SubscribeResponse
	err    error
}

func newBidiReader(
	stream *connect.BidiStreamForClient[message_api.SubscribeRequest, message_api.SubscribeResponse],
) *bidiReader {
	r := &bidiReader{}
	go func() {
		for {
			resp, err := stream.Receive()
			r.mu.Lock()
			if err != nil {
				r.err = err
				r.mu.Unlock()
				return
			}
			r.frames = append(r.frames, resp)
			r.mu.Unlock()
		}
	}()
	return r
}

func (r *bidiReader) snapshot() ([]*message_api.SubscribeResponse, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := make([]*message_api.SubscribeResponse, len(r.frames))
	copy(out, r.frames)
	return out, r.err
}

// subEnvelopeKeys returns the (originatorNodeID, sequenceID) of every delivered envelope, in order.
func subEnvelopeKeys(t *testing.T, frames []*message_api.SubscribeResponse) [][2]uint64 {
	t.Helper()
	var keys [][2]uint64
	for _, f := range frames {
		env := f.GetV1().GetEnvelopes()
		if env == nil {
			continue
		}
		for _, e := range env.GetEnvelopes() {
			u := envelopeTestUtils.UnmarshalUnsignedOriginatorEnvelope(
				t,
				e.GetUnsignedOriginatorEnvelope(),
			)
			keys = append(
				keys,
				[2]uint64{uint64(u.GetOriginatorNodeId()), u.GetOriginatorSequenceId()},
			)
		}
	}
	return keys
}

func subTopicsLive(frames []*message_api.SubscribeResponse) [][]byte {
	var out [][]byte
	for _, f := range frames {
		if tl := f.GetV1().GetTopicsLive(); tl != nil {
			out = append(out, tl.GetTopics()...)
		}
	}
	return out
}

func subCatchupCompletes(frames []*message_api.SubscribeResponse) []uint64 {
	var out []uint64
	for _, f := range frames {
		if cc := f.GetV1().GetCatchupComplete(); cc != nil {
			out = append(out, cc.GetMutateId())
		}
	}
	return out
}

func hasEnvKey(keys [][2]uint64, nodeID, seqID uint64) bool {
	for _, k := range keys {
		if k[0] == nodeID && k[1] == seqID {
			return true
		}
	}
	return false
}

func hasTopicBytes(topics [][]byte, want []byte) bool {
	for _, tp := range topics {
		if bytes.Equal(tp, want) {
			return true
		}
	}
	return false
}

func hasMutateID(ids []uint64, want uint64) bool {
	return slices.Contains(ids, want)
}

// TestSubscribe_CatchUpThenLive verifies a subscription delivers a topic's history (catch-up),
// announces the live boundary, then delivers live messages for that topic only — no duplicates
// across the switch, and nothing from an unsubscribed topic.
func TestSubscribe_CatchUpThenLive(t *testing.T) {
	suite := setupTest(t)
	insertInitialRows(t, suite) // topicA: (100,1),(200,1) in the DB; worker polled past them

	stream := suite.ClientQuery.Subscribe(t.Context())
	reader := newBidiReader(stream)

	require.NoError(t, stream.Send(subMutate(
		1, false,
		[]*message_api.SubscribeRequest_V1_Mutate_Subscription{addSub(topicA, nil)},
		nil,
	)))

	require.Eventually(t, func() bool {
		frames, _ := reader.snapshot()
		keys := subEnvelopeKeys(t, frames)
		return hasEnvKey(keys, 100, 1) && hasEnvKey(keys, 200, 1) &&
			hasTopicBytes(subTopicsLive(frames), topicA) &&
			hasMutateID(subCatchupCompletes(frames), 1)
	}, 10*time.Second, 20*time.Millisecond, "expected topicA history + TopicsLive + CatchupComplete(1)")

	// (100,3) is topicA (live); (100,2),(200,2) are topicB (not subscribed).
	insertAdditionalRows(t, suite.DB)

	require.Eventually(t, func() bool {
		frames, _ := reader.snapshot()
		return hasEnvKey(subEnvelopeKeys(t, frames), 100, 3)
	}, 10*time.Second, 20*time.Millisecond, "expected live topicA delivery of (100,3)")

	frames, err := reader.snapshot()
	require.NoError(t, err)
	keys := subEnvelopeKeys(t, frames)
	require.False(t, hasEnvKey(keys, 100, 2), "topicB (100,2) must not be delivered")
	require.False(t, hasEnvKey(keys, 200, 2), "topicB (200,2) must not be delivered")

	seen := make(map[[2]uint64]struct{}, len(keys))
	for _, k := range keys {
		_, dup := seen[k]
		require.False(t, dup, "duplicate envelope across catch-up/live: %v", k)
		seen[k] = struct{}{}
	}
}

// TestSubscribe_MutateRemoveStopsDelivery verifies an in-place remove stops live delivery for a
// topic while an add in the same stream begins it for another — no reconnect.
func TestSubscribe_MutateRemoveStopsDelivery(t *testing.T) {
	suite := setupTest(t)

	stream := suite.ClientQuery.Subscribe(t.Context())
	reader := newBidiReader(stream)

	require.NoError(t, stream.Send(subMutate(
		1, false,
		[]*message_api.SubscribeRequest_V1_Mutate_Subscription{addSub(topicA, nil)},
		nil,
	)))
	require.Eventually(t, func() bool {
		frames, _ := reader.snapshot()
		return hasMutateID(subCatchupCompletes(frames), 1)
	}, 10*time.Second, 20*time.Millisecond)

	// Remove topicA, add topicB, in one mutation.
	require.NoError(t, stream.Send(subMutate(
		2, false,
		[]*message_api.SubscribeRequest_V1_Mutate_Subscription{addSub(topicB, nil)},
		[][]byte{topicA},
	)))
	require.Eventually(t, func() bool {
		frames, _ := reader.snapshot()
		return hasMutateID(subCatchupCompletes(frames), 2)
	}, 10*time.Second, 20*time.Millisecond)

	insertAdditionalRows(t, suite.DB) // topicB (100,2),(200,2); topicA (100,3)

	require.Eventually(t, func() bool {
		keys := subEnvelopeKeys(t, mustFrames(reader))
		return hasEnvKey(keys, 100, 2) && hasEnvKey(keys, 200, 2)
	}, 10*time.Second, 20*time.Millisecond, "topicB must be delivered live after the add")

	frames, err := reader.snapshot()
	require.NoError(t, err)
	require.False(
		t,
		hasEnvKey(subEnvelopeKeys(t, frames), 100, 3),
		"removed topicA must not deliver (100,3)",
	)
}

// TestSubscribe_HalfCloseHistoryOnlyDrains is the bounded catch-up ("sync") flow: history_only +
// half-close, the server finishes the wave (history, TopicsLive, CatchupComplete) then closes the
// stream itself — the client sees a clean io.EOF, not a truncated result.
func TestSubscribe_HalfCloseHistoryOnlyDrains(t *testing.T) {
	suite := setupTest(t)
	insertInitialRows(t, suite)

	stream := suite.ClientQuery.Subscribe(t.Context())
	reader := newBidiReader(stream)

	require.NoError(t, stream.Send(subMutate(
		7, true, // history_only
		[]*message_api.SubscribeRequest_V1_Mutate_Subscription{addSub(topicA, nil)},
		nil,
	)))
	require.NoError(t, stream.CloseRequest())

	require.Eventually(t, func() bool {
		_, err := reader.snapshot()
		return errors.Is(err, io.EOF)
	}, 10*time.Second, 20*time.Millisecond, "stream should close cleanly after the bounded catch-up")

	frames, err := reader.snapshot()
	require.ErrorIs(t, err, io.EOF)
	keys := subEnvelopeKeys(t, frames)
	require.True(t, hasEnvKey(keys, 100, 1) && hasEnvKey(keys, 200, 1), "history must be delivered")
	require.True(t, hasTopicBytes(subTopicsLive(frames), topicA))
	require.True(t, hasMutateID(subCatchupCompletes(frames), 7))

	// Every data frame carries the wave's mutate_id: history_only pages sent with the live
	// tag (0) would satisfy the delivery assertions above but violate XIP-83 requirement 3.
	require.Empty(t, subEnvelopeKeysTagged(t, frames, 0),
		"history_only pages must never ride the live tag")
	for _, f := range frames {
		if env := f.GetV1().GetEnvelopes(); env != nil {
			require.Equal(t, uint64(7), env.GetMutateId(),
				"every history_only data frame carries the wave's mutate_id")
		}
	}
}

// TestSubscribe_HistoryOnlyOnLiveRejected verifies a history_only add targeting a topic already
// live on the same stream is rejected (one cursor floor per topic).
func TestSubscribe_HistoryOnlyOnLiveRejected(t *testing.T) {
	suite := setupTest(t)

	stream := suite.ClientQuery.Subscribe(t.Context())
	reader := newBidiReader(stream)

	require.NoError(t, stream.Send(subMutate(
		1, false,
		[]*message_api.SubscribeRequest_V1_Mutate_Subscription{addSub(topicA, nil)},
		nil,
	)))
	require.Eventually(t, func() bool {
		frames, _ := reader.snapshot()
		return hasMutateID(subCatchupCompletes(frames), 1)
	}, 10*time.Second, 20*time.Millisecond)

	require.NoError(t, stream.Send(subMutate(
		2, true, // history_only on the already-live topicA
		[]*message_api.SubscribeRequest_V1_Mutate_Subscription{addSub(topicA, nil)},
		nil,
	)))

	require.Eventually(t, func() bool {
		_, err := reader.snapshot()
		return err != nil
	}, 10*time.Second, 20*time.Millisecond, "stream should be failed")

	_, err := reader.snapshot()
	require.Equal(t, connect.CodeInvalidArgument, connect.CodeOf(err))
}

// TestSubscribe_NoPongIsReaped verifies an idle client that never answers the server's liveness
// Ping is reaped with DeadlineExceeded.
func TestSubscribe_NoPongIsReaped(t *testing.T) {
	nodes := []registry.Node{
		{NodeID: 100, IsCanonical: true},
		{NodeID: 200, IsCanonical: true},
	}
	suite := testUtilsApi.NewTestAPIServer(
		t,
		testUtilsApi.WithRegistryNodes(nodes),
		testUtilsApi.WithSendKeepAliveInterval(200*time.Millisecond),
	)

	stream := suite.ClientQuery.Subscribe(t.Context())
	reader := newBidiReader(stream) // only reads; never Pongs

	// Subscribe to nothing and stay idle; the server will Ping, get no Pong, and reap.
	require.NoError(t, stream.Send(subMutate(1, false, nil, nil)))

	require.Eventually(t, func() bool {
		_, err := reader.snapshot()
		return err != nil
	}, 5*time.Second, 50*time.Millisecond, "an idle client that never Pongs must be reaped")

	_, err := reader.snapshot()
	require.Equal(t, connect.CodeDeadlineExceeded, connect.CodeOf(err))
}

// TestSubscribe_ReAddLiveTopicDoesNotReplay verifies that re-adding an already-live topic WITHOUT a
// remove (e.g. a herald re-issuing an add with a default/stale cursor) is a no-op: it must not reset
// the live cursor and re-deliver already-sent envelopes (no duplicates, no backwards seqID). Replay
// is only available via remove+re-add.
func TestSubscribe_ReAddLiveTopicDoesNotReplay(t *testing.T) {
	suite := setupTest(t)
	insertInitialRows(t, suite) // topicA: (100,1),(200,1)

	stream := suite.ClientQuery.Subscribe(t.Context())
	reader := newBidiReader(stream)

	// Subscribe to topicA; receive its history and go live.
	require.NoError(t, stream.Send(subMutate(
		1, false,
		[]*message_api.SubscribeRequest_V1_Mutate_Subscription{addSub(topicA, nil)},
		nil,
	)))
	require.Eventually(t, func() bool {
		keys := subEnvelopeKeys(t, mustFrames(reader))
		return hasEnvKey(keys, 100, 1) && hasEnvKey(keys, 200, 1) &&
			hasMutateID(subCatchupCompletes(mustFrames(reader)), 1)
	}, 10*time.Second, 20*time.Millisecond)

	// Deliver a live message so the live cursor advances past the history.
	insertAdditionalRows(t, suite.DB) // topicA (100,3); topicB (100,2),(200,2) not subscribed
	require.Eventually(t, func() bool {
		return hasEnvKey(subEnvelopeKeys(t, mustFrames(reader)), 100, 3)
	}, 10*time.Second, 20*time.Millisecond, "expected live topicA (100,3)")

	// Re-add topicA WITHOUT removing it, with a zero cursor. This must be a no-op — NOT a replay.
	require.NoError(t, stream.Send(subMutate(
		2, false,
		[]*message_api.SubscribeRequest_V1_Mutate_Subscription{addSub(topicA, map[uint32]uint64{})},
		nil,
	)))
	// The no-op mutate is still confirmed (its adds collapsed to none -> removes-only path).
	require.Eventually(t, func() bool {
		return hasMutateID(subCatchupCompletes(mustFrames(reader)), 2)
	}, 10*time.Second, 20*time.Millisecond, "expected CatchupComplete(2) for the no-op re-add")

	// No envelope may be delivered twice: a replay would re-send (100,1),(200,1),(100,3).
	counts := make(map[[2]uint64]int)
	for _, k := range subEnvelopeKeys(t, mustFrames(reader)) {
		counts[k]++
	}
	for k, n := range counts {
		require.Equalf(t, 1, n, "envelope %v delivered %d times: re-add must not replay", k, n)
	}
}

func mustFrames(r *bidiReader) []*message_api.SubscribeResponse {
	frames, _ := r.snapshot()
	return frames
}

// ---- XIP-83 delivery tagging & ordering (server requirements 3 and 4) ----

// envRow builds one insertable gateway envelope row for the tagging/ordering tests.
func envRow(
	t *testing.T,
	payerID sql.NullInt32,
	nodeID int32,
	seqID int64,
	topicBytes []byte,
) queries.InsertGatewayEnvelopeV3Params {
	t.Helper()
	return queries.InsertGatewayEnvelopeV3Params{
		OriginatorNodeID:     nodeID,
		OriginatorSequenceID: seqID,
		Topic:                topicBytes,
		PayerID:              payerID,
		OriginatorEnvelope: testutils.Marshal(
			t,
			envelopeTestUtils.CreateOriginatorEnvelopeWithTopic(
				t,
				uint32(nodeID),
				uint64(seqID),
				topicBytes,
			),
		),
	}
}

// subEnvelopeKeysTagged returns the (originator, sequence) keys of envelopes carried by
// Envelopes frames stamped with the given wave tag, in receive order.
func subEnvelopeKeysTagged(
	t *testing.T,
	frames []*message_api.SubscribeResponse,
	tag uint64,
) [][2]uint64 {
	t.Helper()
	var keys [][2]uint64
	for _, f := range frames {
		env := f.GetV1().GetEnvelopes()
		if env == nil || env.GetMutateId() != tag {
			continue
		}
		for _, e := range env.GetEnvelopes() {
			u := envelopeTestUtils.UnmarshalUnsignedOriginatorEnvelope(
				t,
				e.GetUnsignedOriginatorEnvelope(),
			)
			keys = append(
				keys,
				[2]uint64{uint64(u.GetOriginatorNodeId()), u.GetOriginatorSequenceId()},
			)
		}
	}
	return keys
}

// requirePerOriginatorAscending asserts each originator's sequence ids strictly ascend in
// the given key order — the total-order shape both delivery lanes guarantee per originator.
func requirePerOriginatorAscending(t *testing.T, keys [][2]uint64, desc string) {
	t.Helper()
	last := make(map[uint64]uint64)
	for _, k := range keys {
		if prev, ok := last[k[0]]; ok {
			require.Greater(t, k[1], prev,
				"%s: originator %d sequences must strictly ascend", desc, k[0])
		}
		last[k[0]] = k[1]
	}
}

// requireExactlyOnce asserts no (originator, sequence) key appears twice across the keys.
func requireExactlyOnce(t *testing.T, keys [][2]uint64, desc string) {
	t.Helper()
	seen := make(map[[2]uint64]struct{}, len(keys))
	for _, k := range keys {
		_, dup := seen[k]
		require.Falsef(t, dup, "%s: envelope %v delivered more than once", desc, k)
		seen[k] = struct{}{}
	}
}

// TestSubscribe_ReplayTaggedWithWaveLiveTaggedZero pins delivery tagging under overlapping
// waves: every replay frame carries exactly the mutate_id of the wave that produced it, and
// live frames carry 0.
func TestSubscribe_ReplayTaggedWithWaveLiveTaggedZero(t *testing.T) {
	suite := setupTest(t)
	payerID := db.NullInt32(testutils.CreatePayer(t, suite.DB))

	// History: topicA fed by originator 100, topicB by originator 200.
	testutils.InsertGatewayEnvelopes(t, suite.DB, []queries.InsertGatewayEnvelopeV3Params{
		envRow(t, payerID, 100, 1, topicA),
		envRow(t, payerID, 100, 2, topicA),
		envRow(t, payerID, 100, 3, topicA),
		envRow(t, payerID, 200, 1, topicB),
		envRow(t, payerID, 200, 2, topicB),
	})
	ctx, cancel := context.WithTimeout(t.Context(), 2*time.Second)
	defer cancel()
	require.NoError(t, suite.MessageService.AwaitCursor(ctx, db.VectorClock{100: 3, 200: 2}))

	stream := suite.ClientQuery.Subscribe(t.Context())
	reader := newBidiReader(stream)

	// Two overlapping waves, back to back: their replays race on the same stream.
	require.NoError(t, stream.Send(subMutate(
		7, false,
		[]*message_api.SubscribeRequest_V1_Mutate_Subscription{addSub(topicA, nil)},
		nil,
	)))
	require.NoError(t, stream.Send(subMutate(
		8, false,
		[]*message_api.SubscribeRequest_V1_Mutate_Subscription{addSub(topicB, nil)},
		nil,
	)))
	require.Eventually(t, func() bool {
		cc := subCatchupCompletes(mustFrames(reader))
		return hasMutateID(cc, 7) && hasMutateID(cc, 8)
	}, 10*time.Second, 20*time.Millisecond, "both waves must complete")

	frames := mustFrames(reader)
	wave7 := subEnvelopeKeysTagged(t, frames, 7)
	wave8 := subEnvelopeKeysTagged(t, frames, 8)
	require.Len(t, wave7, 3, "wave 7 delivers exactly topicA's history, tagged 7")
	require.Len(t, wave8, 2, "wave 8 delivers exactly topicB's history, tagged 8")
	for _, k := range wave7 {
		require.Equal(t, uint64(100), k[0], "wave 7 must carry only topicA's originator")
	}
	for _, k := range wave8 {
		require.Equal(t, uint64(200), k[0], "wave 8 must carry only topicB's originator")
	}
	require.Empty(t, subEnvelopeKeysTagged(t, frames, 0), "no replay frame may carry the live tag")

	// Live tail on both topics is tagged 0.
	testutils.InsertGatewayEnvelopes(t, suite.DB, []queries.InsertGatewayEnvelopeV3Params{
		envRow(t, payerID, 100, 4, topicA),
		envRow(t, payerID, 200, 3, topicB),
	})
	require.Eventually(t, func() bool {
		live := subEnvelopeKeysTagged(t, mustFrames(reader), 0)
		return hasEnvKey(live, 100, 4) && hasEnvKey(live, 200, 3)
	}, 10*time.Second, 20*time.Millisecond, "live tail must be tagged 0")

	requireExactlyOnce(t, subEnvelopeKeys(t, mustFrames(reader)), "all lanes")
}

// TestSubscribe_WaveReplayPerOriginatorOrderAcrossTopics pins wave order: one wave covering
// two topics whose envelopes interleave per originator must replay each originator's
// envelopes in ascending sequence order across BOTH topics — one merged cursor-ordered
// pass, not one topic's burst then the other's.
func TestSubscribe_WaveReplayPerOriginatorOrderAcrossTopics(t *testing.T) {
	suite := setupTest(t)
	payerID := db.NullInt32(testutils.CreatePayer(t, suite.DB))

	// Originator 100 alternates between the two topics; originator 200 interleaves too.
	testutils.InsertGatewayEnvelopes(t, suite.DB, []queries.InsertGatewayEnvelopeV3Params{
		envRow(t, payerID, 100, 1, topicA),
		envRow(t, payerID, 100, 2, topicB),
		envRow(t, payerID, 100, 3, topicA),
		envRow(t, payerID, 100, 4, topicB),
		envRow(t, payerID, 100, 5, topicA),
		envRow(t, payerID, 100, 6, topicB),
		envRow(t, payerID, 200, 1, topicB),
		envRow(t, payerID, 200, 2, topicA),
	})
	ctx, cancel := context.WithTimeout(t.Context(), 2*time.Second)
	defer cancel()
	require.NoError(t, suite.MessageService.AwaitCursor(ctx, db.VectorClock{100: 6, 200: 2}))

	stream := suite.ClientQuery.Subscribe(t.Context())
	reader := newBidiReader(stream)

	require.NoError(t, stream.Send(subMutate(
		3, false,
		[]*message_api.SubscribeRequest_V1_Mutate_Subscription{
			addSub(topicA, nil),
			addSub(topicB, nil),
		},
		nil,
	)))
	require.Eventually(t, func() bool {
		return hasMutateID(subCatchupCompletes(mustFrames(reader)), 3)
	}, 10*time.Second, 20*time.Millisecond)

	replay := subEnvelopeKeysTagged(t, mustFrames(reader), 3)
	require.Len(t, replay, 8, "the wave delivers both topics' history, tagged 3")
	requirePerOriginatorAscending(t, replay, "wave replay across interleaved topics")
	requireExactlyOnce(t, replay, "wave replay")
}

// TestSubscribe_LivePerOriginatorOrderAcrossTopics pins live order: live (mutate_id 0)
// envelopes across all live topics arrive in ascending sequence order per originator.
func TestSubscribe_LivePerOriginatorOrderAcrossTopics(t *testing.T) {
	suite := setupTest(t)
	payerID := db.NullInt32(testutils.CreatePayer(t, suite.DB))

	stream := suite.ClientQuery.Subscribe(t.Context())
	reader := newBidiReader(stream)

	require.NoError(t, stream.Send(subMutate(
		1, false,
		[]*message_api.SubscribeRequest_V1_Mutate_Subscription{
			addSub(topicA, nil),
			addSub(topicB, nil),
		},
		nil,
	)))
	require.Eventually(t, func() bool {
		return hasMutateID(subCatchupCompletes(mustFrames(reader)), 1)
	}, 10*time.Second, 20*time.Millisecond)

	// Each originator alternates between the two live topics.
	testutils.InsertGatewayEnvelopes(t, suite.DB, []queries.InsertGatewayEnvelopeV3Params{
		envRow(t, payerID, 100, 1, topicA),
		envRow(t, payerID, 100, 2, topicB),
		envRow(t, payerID, 100, 3, topicA),
		envRow(t, payerID, 200, 1, topicB),
		envRow(t, payerID, 200, 2, topicA),
	})
	require.Eventually(t, func() bool {
		live := subEnvelopeKeysTagged(t, mustFrames(reader), 0)
		return hasEnvKey(live, 100, 3) && hasEnvKey(live, 200, 2)
	}, 10*time.Second, 20*time.Millisecond)

	live := subEnvelopeKeysTagged(t, mustFrames(reader), 0)
	require.Len(t, live, 5)
	requirePerOriginatorAscending(t, live, "live lane across topics")
	requireExactlyOnce(t, live, "live lane")
}

// TestSubscribe_SeamLiveWaitsForCatchupComplete pins the seam: while a wave replays a
// topic, the topic speaks live (mutate_id 0) only after the wave's CatchupComplete.
// Envelopes published mid-wave arrive exactly once — either folded into the wave (tagged)
// or live after its CatchupComplete — and never as a live frame before it. A pre-live
// sentinel topic publishes throughout: the gate is per-topic, so live delivery for other
// subscriptions keeps flowing while the wave's topic is gated.
func TestSubscribe_SeamLiveWaitsForCatchupComplete(t *testing.T) {
	suite := setupTest(t)
	payerID := db.NullInt32(testutils.CreatePayer(t, suite.DB))
	sentinelTopic := topic.NewTopic(topic.TopicKindGroupMessagesV1, []byte("seam-sentinel")).
		Bytes()

	testutils.InsertGatewayEnvelopes(t, suite.DB, []queries.InsertGatewayEnvelopeV3Params{
		envRow(t, payerID, 100, 1, topicA),
		envRow(t, payerID, 100, 2, topicA),
		envRow(t, payerID, 100, 3, topicA),
	})
	ctx, cancel := context.WithTimeout(t.Context(), 2*time.Second)
	defer cancel()
	require.NoError(t, suite.MessageService.AwaitCursor(ctx, db.VectorClock{100: 3}))

	stream := suite.ClientQuery.Subscribe(t.Context())
	reader := newBidiReader(stream)

	// Sentinel first: subscribed, live, and demonstrably delivering before the wave starts.
	require.NoError(t, stream.Send(subMutate(
		1, false,
		[]*message_api.SubscribeRequest_V1_Mutate_Subscription{addSub(sentinelTopic, nil)},
		nil,
	)))
	require.Eventually(t, func() bool {
		return hasMutateID(subCatchupCompletes(mustFrames(reader)), 1)
	}, 10*time.Second, 20*time.Millisecond)
	testutils.InsertGatewayEnvelopes(t, suite.DB, []queries.InsertGatewayEnvelopeV3Params{
		envRow(t, payerID, 200, 1, sentinelTopic),
	})
	require.Eventually(t, func() bool {
		return hasEnvKey(subEnvelopeKeysTagged(t, mustFrames(reader), 0), 200, 1)
	}, 10*time.Second, 20*time.Millisecond, "sentinel must be live before the wave starts")

	// Start the wave and immediately publish into it, racing the replay — and keep the
	// sentinel publishing during the wave.
	require.NoError(t, stream.Send(subMutate(
		9, false,
		[]*message_api.SubscribeRequest_V1_Mutate_Subscription{addSub(topicA, nil)},
		nil,
	)))
	testutils.InsertGatewayEnvelopes(t, suite.DB, []queries.InsertGatewayEnvelopeV3Params{
		envRow(t, payerID, 100, 4, topicA),
		envRow(t, payerID, 100, 5, topicA),
		envRow(t, payerID, 200, 2, sentinelTopic),
	})

	require.Eventually(t, func() bool {
		frames := mustFrames(reader)
		keys := subEnvelopeKeys(t, frames)
		return hasMutateID(subCatchupCompletes(frames), 9) &&
			hasEnvKey(keys, 100, 4) && hasEnvKey(keys, 100, 5)
	}, 10*time.Second, 20*time.Millisecond, "wave complete + racers delivered")

	// After the wave: the sentinel's live lane must still be flowing.
	testutils.InsertGatewayEnvelopes(t, suite.DB, []queries.InsertGatewayEnvelopeV3Params{
		envRow(t, payerID, 200, 3, sentinelTopic),
	})
	require.Eventually(t, func() bool {
		live := subEnvelopeKeysTagged(t, mustFrames(reader), 0)
		return hasEnvKey(live, 200, 2) && hasEnvKey(live, 200, 3)
	}, 10*time.Second, 20*time.Millisecond,
		"sentinel tag-0 delivery must keep flowing across the wave")

	frames, err := reader.snapshot()
	require.NoError(t, err)
	requireExactlyOnce(t, subEnvelopeKeys(t, frames), "wave + live")
	requirePerOriginatorAscending(t, subEnvelopeKeysTagged(t, frames, 9), "wave 9 replay")
	requirePerOriginatorAscending(t, subEnvelopeKeysTagged(t, frames, 0), "live lane")

	// The sentinel's envelopes (originator 200) travel the live lane only — never the wave's.
	for _, k := range subEnvelopeKeysTagged(t, frames, 9) {
		require.Equal(t, uint64(100), k[0], "wave 9 must not capture the sentinel's envelopes")
	}

	// The seam, scoped to the wave's topic (originator 100 publishes only to topicA here): no
	// live frame for it may precede CatchupComplete(9), and every wave-tagged frame must
	// precede it. Sentinel (originator 200) tag-0 frames may land on either side.
	ccIdx := -1
	for i, f := range frames {
		if cc := f.GetV1().GetCatchupComplete(); cc != nil && cc.GetMutateId() == 9 {
			ccIdx = i
		}
	}
	require.GreaterOrEqual(t, ccIdx, 0)
	for i, f := range frames {
		env := f.GetV1().GetEnvelopes()
		if env == nil || len(env.GetEnvelopes()) == 0 {
			continue
		}
		switch env.GetMutateId() {
		case 0:
			for _, e := range env.GetEnvelopes() {
				u := envelopeTestUtils.UnmarshalUnsignedOriginatorEnvelope(
					t,
					e.GetUnsignedOriginatorEnvelope(),
				)
				if u.GetOriginatorNodeId() == 100 {
					require.Greater(t, i, ccIdx,
						"a live frame for the wave's topic must follow its CatchupComplete")
				}
			}
		case 9:
			require.Less(t, i, ccIdx,
				"a wave replay frame must precede its CatchupComplete")
		default:
			t.Fatalf("frame %d carries unexpected tag %d", i, env.GetMutateId())
		}
	}
}

// TestSubscribe_ResetMidWaveReplaysUnderNewTag covers remove+re-add (reset) racing an
// envelope-bearing first wave: the reset topic is owned by the newer wave, whose replay is
// stamped with the new mutate_id; the stale wave's remaining pages are dropped, so nothing
// is delivered twice and only the owning wave announces the topic.
func TestSubscribe_ResetMidWaveReplaysUnderNewTag(t *testing.T) {
	suite := setupTest(t)
	payerID := db.NullInt32(testutils.CreatePayer(t, suite.DB))
	resetTopic := topic.NewTopic(topic.TopicKindGroupMessagesV1, []byte("reset-mid-wave")).
		Bytes()

	// More than one wave-scan page (topicPageLimit rows) so the reset can land between the
	// stale wave's pages.
	const total = 600
	rows := make([]queries.InsertGatewayEnvelopeV3Params, 0, total)
	for i := int64(1); i <= total; i++ {
		rows = append(rows, envRow(t, payerID, 100, i, resetTopic))
	}
	testutils.InsertGatewayEnvelopes(t, suite.DB, rows)
	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	defer cancel()
	require.NoError(t, suite.MessageService.AwaitCursor(ctx, db.VectorClock{100: total}))

	stream := suite.ClientQuery.Subscribe(t.Context())
	reader := newBidiReader(stream)

	// Wave 1 starts the replay; the reset (remove + re-add from an empty cursor) is sent
	// immediately, racing wave 1's pages.
	require.NoError(t, stream.Send(subMutate(
		1, false,
		[]*message_api.SubscribeRequest_V1_Mutate_Subscription{addSub(resetTopic, nil)},
		nil,
	)))
	require.NoError(t, stream.Send(subMutate(
		2, false,
		[]*message_api.SubscribeRequest_V1_Mutate_Subscription{addSub(resetTopic, nil)},
		[][]byte{resetTopic},
	)))

	require.Eventually(t, func() bool {
		cc := subCatchupCompletes(mustFrames(reader))
		return hasMutateID(cc, 1) && hasMutateID(cc, 2)
	}, 20*time.Second, 20*time.Millisecond, "both waves must complete")

	frames, err := reader.snapshot()
	require.NoError(t, err)
	wave1 := subEnvelopeKeysTagged(t, frames, 1)
	wave2 := subEnvelopeKeysTagged(t, frames, 2)
	requireExactlyOnce(t, wave1, "stale wave replay")
	requireExactlyOnce(t, wave2, "reset wave replay")
	requirePerOriginatorAscending(t, wave1, "stale wave replay")
	requirePerOriginatorAscending(t, wave2, "reset wave replay")

	// Between them the two waves cover every seeded envelope exactly once: the reset wave owns
	// everything the stale wave had not yet delivered when the reset applied. (Do not assert
	// how the split falls — the reset races the stale wave's pages.)
	seen := make(map[[2]uint64]struct{}, total)
	for _, k := range wave1 {
		seen[k] = struct{}{}
	}
	for _, k := range wave2 {
		_, both := seen[k]
		require.Falsef(t, both, "envelope %v delivered under both tag 1 and tag 2", k)
		seen[k] = struct{}{}
	}
	require.Len(t, seen, total, "the two waves together must cover every envelope exactly once")
	for i := int64(1); i <= total; i++ {
		require.Contains(t, seen, [2]uint64{100, uint64(i)})
	}
	require.True(t, hasEnvKey(wave2, 100, total),
		"the reset wave must deliver at least the tail of the history")

	// Only the owning (reset) wave announces the topic, and nothing rides the live tag: the
	// history predates the subscription, so the live lane has nothing to say before CC(2).
	announced := 0
	for _, tl := range subTopicsLive(frames) {
		if bytes.Equal(tl, resetTopic) {
			announced++
		}
	}
	require.Equal(t, 1, announced, "exactly one TopicsLive may announce the reset topic")
	require.Empty(t, subEnvelopeKeysTagged(t, frames, 0), "no envelope may ride the live tag")
}

// TestSubscribe_WaveScanPaginatesPastPageLimit drives the wave's merged keyset scan across a
// page boundary that lands mid-originator: the resume must be strictly-after the last row (a
// >= row-value comparison would re-deliver the boundary row) without skipping the next
// originator's rows.
func TestSubscribe_WaveScanPaginatesPastPageLimit(t *testing.T) {
	suite := setupTest(t)
	payerID := db.NullInt32(testutils.CreatePayer(t, suite.DB))
	pageT1 := topic.NewTopic(topic.TopicKindGroupMessagesV1, []byte("wave-page-t1")).Bytes()
	pageT2 := topic.NewTopic(topic.TopicKindGroupMessagesV1, []byte("wave-page-t2")).Bytes()

	// topicPageLimit (500) + 1 rows for originator 100, interleaved across both topics, put
	// the page boundary inside originator 100's run; originator 200's rows follow on page 2.
	const heavy = 501
	const light = 5
	rows := make([]queries.InsertGatewayEnvelopeV3Params, 0, heavy+light)
	for i := int64(1); i <= heavy; i++ {
		tp := pageT1
		if i%2 == 0 {
			tp = pageT2
		}
		rows = append(rows, envRow(t, payerID, 100, i, tp))
	}
	for i := int64(1); i <= light; i++ {
		rows = append(rows, envRow(t, payerID, 200, i, pageT1))
	}
	testutils.InsertGatewayEnvelopes(t, suite.DB, rows)
	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	defer cancel()
	require.NoError(
		t,
		suite.MessageService.AwaitCursor(ctx, db.VectorClock{100: heavy, 200: light}),
	)

	stream := suite.ClientQuery.Subscribe(t.Context())
	reader := newBidiReader(stream)

	require.NoError(t, stream.Send(subMutate(
		5, false,
		[]*message_api.SubscribeRequest_V1_Mutate_Subscription{
			addSub(pageT1, nil),
			addSub(pageT2, nil),
		},
		nil,
	)))
	require.Eventually(t, func() bool {
		return hasMutateID(subCatchupCompletes(mustFrames(reader)), 5)
	}, 20*time.Second, 20*time.Millisecond)

	frames, err := reader.snapshot()
	require.NoError(t, err)
	replay := subEnvelopeKeysTagged(t, frames, 5)
	require.Len(t, replay, heavy+light, "every seeded envelope arrives tagged with the wave")
	requireExactlyOnce(t, replay, "wave replay across pages")
	requirePerOriginatorAscending(t, replay, "wave replay across pages")

	waveFrames := 0
	for _, f := range frames {
		if env := f.GetV1().GetEnvelopes(); env != nil && env.GetMutateId() == 5 {
			waveFrames++
		}
	}
	require.GreaterOrEqual(t, waveFrames, 2, "the replay must span multiple scan pages")
}

// TestSubscribe_CursorNamedUnknownOriginatorReplayed covers a client cursor naming an
// originator the TTL-cached originator list has not seen (its rows exist in
// gateway_envelopes_meta, but the cached list is stale): the wave must still pin a ceiling
// for it and replay its rows, rather than silently dropping them from the catch-up (the
// legacy per-topic path replayed them — see TestSubscribeTopics_AcceptsUnknownOriginatorInCursor).
func TestSubscribe_CursorNamedUnknownOriginatorReplayed(t *testing.T) {
	suite := setupTest(t)
	payerID := db.NullInt32(testutils.CreatePayer(t, suite.DB))

	// topicA history from originator 100 (known everywhere) and originator 300 (not in the
	// registry, so the worker never polls it — nothing can arrive via the live lane).
	testutils.InsertGatewayEnvelopes(t, suite.DB, []queries.InsertGatewayEnvelopeV3Params{
		envRow(t, payerID, 100, 1, topicA),
		envRow(t, payerID, 300, 1, topicA),
		envRow(t, payerID, 300, 2, topicA),
		envRow(t, payerID, 300, 3, topicA),
	})
	ctx, cancel := context.WithTimeout(t.Context(), 2*time.Second)
	defer cancel()
	require.NoError(t, suite.MessageService.AwaitCursor(ctx, db.VectorClock{100: 1}))

	// Fabricate the stale-cache state: drop 300 from gateway_envelopes_latest (the cached
	// originator list's source) while its rows stay in gateway_envelopes_meta. This is what a
	// TTL-stale CachedOriginatorList sees when an originator's first rows land after the
	// cache was filled — deterministic here instead of racing the 100ms test TTL.
	_, err := suite.DB.ExecContext(t.Context(),
		"DELETE FROM gateway_envelopes_latest WHERE originator_node_id = 300")
	require.NoError(t, err)

	stream := suite.ClientQuery.Subscribe(t.Context())
	reader := newBidiReader(stream)

	// The cursor names originator 300 below its newest row: the wave owes (300,2),(300,3)
	// even though the originator list has never heard of 300.
	require.NoError(t, stream.Send(subMutate(
		1, false,
		[]*message_api.SubscribeRequest_V1_Mutate_Subscription{
			addSub(topicA, map[uint32]uint64{300: 1}),
		},
		nil,
	)))

	require.Eventually(t, func() bool {
		frames := mustFrames(reader)
		replay := subEnvelopeKeysTagged(t, frames, 1)
		return hasEnvKey(replay, 100, 1) &&
			hasEnvKey(replay, 300, 2) && hasEnvKey(replay, 300, 3) &&
			hasMutateID(subCatchupCompletes(frames), 1)
	}, 10*time.Second, 20*time.Millisecond,
		"the wave must replay the cursor-named originator 300, then CatchupComplete(1)")

	frames, err := reader.snapshot()
	require.NoError(t, err)
	replay := subEnvelopeKeysTagged(t, frames, 1)
	require.False(t, hasEnvKey(replay, 300, 1), "the cursor floor (300,1) must not be re-delivered")
	requireExactlyOnce(t, subEnvelopeKeys(t, frames), "unknown-originator replay")
	requirePerOriginatorAscending(t, replay, "unknown-originator replay")
}

// TestSubscribe_DuplicateAddsFirstCursorWins pins handleMutate's add dedup: when one Mutate
// carries two adds for the same topic, the first add's cursor is the floor (the duplicate is
// dropped), so replay starts strictly after it — and the topic is announced once.
func TestSubscribe_DuplicateAddsFirstCursorWins(t *testing.T) {
	suite := setupTest(t)
	payerID := db.NullInt32(testutils.CreatePayer(t, suite.DB))
	dupTopic := topic.NewTopic(topic.TopicKindGroupMessagesV1, []byte("dup-adds")).Bytes()

	testutils.InsertGatewayEnvelopes(t, suite.DB, []queries.InsertGatewayEnvelopeV3Params{
		envRow(t, payerID, 100, 1, dupTopic),
		envRow(t, payerID, 100, 2, dupTopic),
		envRow(t, payerID, 100, 3, dupTopic),
		envRow(t, payerID, 100, 4, dupTopic),
	})
	ctx, cancel := context.WithTimeout(t.Context(), 2*time.Second)
	defer cancel()
	require.NoError(t, suite.MessageService.AwaitCursor(ctx, db.VectorClock{100: 4}))

	stream := suite.ClientQuery.Subscribe(t.Context())
	reader := newBidiReader(stream)

	// First add carries cursor {100: 2}; the duplicate second add carries an empty cursor.
	require.NoError(t, stream.Send(subMutate(
		4, false,
		[]*message_api.SubscribeRequest_V1_Mutate_Subscription{
			addSub(dupTopic, map[uint32]uint64{100: 2}),
			addSub(dupTopic, map[uint32]uint64{}),
		},
		nil,
	)))
	require.Eventually(t, func() bool {
		return hasMutateID(subCatchupCompletes(mustFrames(reader)), 4)
	}, 10*time.Second, 20*time.Millisecond)

	frames, err := reader.snapshot()
	require.NoError(t, err)
	require.Equal(t, [][2]uint64{{100, 3}, {100, 4}}, subEnvelopeKeysTagged(t, frames, 4),
		"replay must start after the FIRST add's cursor, exactly once")
	require.Empty(t, subEnvelopeKeysTagged(t, frames, 0))

	announced := 0
	for _, tl := range subTopicsLive(frames) {
		if bytes.Equal(tl, dupTopic) {
			announced++
		}
	}
	require.Equal(t, 1, announced, "the deduped topic is announced exactly once")
}

// TestSubscribe_AddsRequireNonzeroMutateId pins the request-side rule the tag depends on:
// a Mutate with adds and mutate_id 0 (the live tag) fails the stream with InvalidArgument.
func TestSubscribe_AddsRequireNonzeroMutateId(t *testing.T) {
	suite := setupTest(t)

	stream := suite.ClientQuery.Subscribe(t.Context())
	reader := newBidiReader(stream)

	require.NoError(t, stream.Send(subMutate(
		0, false,
		[]*message_api.SubscribeRequest_V1_Mutate_Subscription{addSub(topicA, nil)},
		nil,
	)))

	require.Eventually(t, func() bool {
		_, err := reader.snapshot()
		return err != nil
	}, 10*time.Second, 20*time.Millisecond, "adds with mutate_id 0 must fail the stream")

	_, err := reader.snapshot()
	require.Equal(t, connect.CodeInvalidArgument, connect.CodeOf(err))
}

// TestSubscribe_EmptyMutateAcked pins the ack rule for the degenerate Mutate shape: a Mutate
// with no adds and no removes is still confirmed with exactly one CatchupComplete echoing its
// mutate_id, and the stream stays healthy afterwards (a subsequent subscribe works end to end).
func TestSubscribe_EmptyMutateAcked(t *testing.T) {
	suite := setupTest(t)

	stream := suite.ClientQuery.Subscribe(t.Context())
	reader := newBidiReader(stream)

	// Empty Mutate: no adds, no removes.
	require.NoError(t, stream.Send(subMutate(3, false, nil, nil)))
	require.Eventually(t, func() bool {
		return hasMutateID(subCatchupCompletes(mustFrames(reader)), 3)
	}, 10*time.Second, 20*time.Millisecond, "an empty Mutate must be acked promptly")

	// The stream is still usable: a subsequent subscribe completes its own wave.
	require.NoError(t, stream.Send(subMutate(
		4, false,
		[]*message_api.SubscribeRequest_V1_Mutate_Subscription{addSub(topicA, nil)},
		nil,
	)))
	require.Eventually(t, func() bool {
		return hasMutateID(subCatchupCompletes(mustFrames(reader)), 4)
	}, 10*time.Second, 20*time.Millisecond, "the stream must stay healthy after an empty Mutate")

	frames, err := reader.snapshot()
	require.NoError(t, err)
	acks := 0
	for _, id := range subCatchupCompletes(frames) {
		if id == 3 {
			acks++
		}
	}
	require.Equal(t, 1, acks, "exactly one CatchupComplete must echo the empty Mutate's id")
}
