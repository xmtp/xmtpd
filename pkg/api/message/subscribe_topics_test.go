package message_test

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
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
	"github.com/xmtp/xmtpd/pkg/registry"
	"github.com/xmtp/xmtpd/pkg/testutils"
	testUtilsApi "github.com/xmtp/xmtpd/pkg/testutils/api"
	envelopeTestUtils "github.com/xmtp/xmtpd/pkg/testutils/envelopes"
	"github.com/xmtp/xmtpd/pkg/topic"
)

// setupTopicTest creates a test API server and returns the client, DB and mocks.
func setupTopicTest(
	t *testing.T,
) (message_apiconnect.ReplicationApiClient, *sql.DB, testUtilsApi.APIServerMocks) {
	nodes := []registry.Node{
		{NodeID: 100, IsCanonical: true},
		{NodeID: 200, IsCanonical: true},
	}
	suite := testUtilsApi.NewTestAPIServer(t, testUtilsApi.WithRegistryNodes(nodes))
	return suite.ClientReplication, suite.DB, suite.APIServerMocks
}

// makeFilter creates a TopicFilter with the given topic and optional LastSeen cursor.
func makeFilter(
	topicBytes []byte, cursor map[uint32]uint64,
) *message_api.SubscribeTopicsRequest_TopicFilter {
	f := &message_api.SubscribeTopicsRequest_TopicFilter{
		Topic: topicBytes,
	}
	if cursor != nil {
		f.LastSeen = &envelopes.Cursor{NodeIdToSequenceId: cursor}
	}
	return f
}

// makeEnvRow constructs an InsertGatewayEnvelopeParams from minimal inputs.
func makeEnvRow(
	t *testing.T, nodeID uint32, seqID uint64, topicBytes []byte, payerID sql.NullInt32,
) queries.InsertGatewayEnvelopeParams {
	t.Helper()
	return queries.InsertGatewayEnvelopeParams{
		OriginatorNodeID:     int32(nodeID),
		OriginatorSequenceID: int64(seqID),
		Topic:                topicBytes,
		PayerID:              payerID,
		OriginatorEnvelope: testutils.Marshal(
			t,
			envelopeTestUtils.CreateOriginatorEnvelopeWithTopic(t, nodeID, seqID, topicBytes),
		),
	}
}

// subscribeTopics opens a SubscribeTopics stream, consumes the initial
// STARTED status message, and returns the ready-to-use stream.
func subscribeTopics(
	t *testing.T,
	client message_apiconnect.ReplicationApiClient,
	ctx context.Context,
	filters []*message_api.SubscribeTopicsRequest_TopicFilter,
) *connect.ServerStreamForClient[message_api.SubscribeTopicsResponse] {
	t.Helper()
	stream, err := client.SubscribeTopics(
		ctx,
		connect.NewRequest(&message_api.SubscribeTopicsRequest{Filters: filters}),
	)
	require.NoError(t, err)
	require.True(t, stream.Receive())
	require.Equal(t,
		message_api.SubscribeTopicsResponse_SUBSCRIPTION_STATUS_STARTED,
		stream.Msg().GetStatusUpdate().GetStatus(),
	)
	return stream
}

// insertAndWait inserts gateway envelopes and waits for the subscribe worker to poll them.
func insertAndWait(t *testing.T, store *sql.DB, rows []queries.InsertGatewayEnvelopeParams) {
	t.Helper()
	testutils.InsertGatewayEnvelopes(t, store, rows)
	time.Sleep(message.SubscribeWorkerPollTime + 100*time.Millisecond)
}

// requireOriginatorOrdering verifies that envelopes are ordered per originator.
func requireOriginatorOrdering(t *testing.T, envs []*envelopes.OriginatorEnvelope) {
	t.Helper()
	lastSeqByNode := make(map[uint32]uint64)
	for _, env := range envs {
		decoded := envelopeTestUtils.UnmarshalUnsignedOriginatorEnvelope(
			t, env.GetUnsignedOriginatorEnvelope(),
		)
		nodeID := decoded.GetOriginatorNodeId()
		seqID := decoded.GetOriginatorSequenceId()
		if last, ok := lastSeqByNode[nodeID]; ok {
			require.Greater(t, seqID, last,
				"out of order for originator %d: got %d after %d", nodeID, seqID, last)
		}
		lastSeqByNode[nodeID] = seqID
	}
}

// collectTopicEnvelopes receives envelopes from a SubscribeTopicsResponse stream
// until expectedCount envelopes have been collected or timeout.
// Status update messages are skipped.
func collectTopicEnvelopes(
	t *testing.T,
	stream *connect.ServerStreamForClient[message_api.SubscribeTopicsResponse],
	expectedCount int,
) []*envelopes.OriginatorEnvelope {
	t.Helper()

	var collected []*envelopes.OriginatorEnvelope
	for len(collected) < expectedCount {
		if !stream.Receive() {
			break
		}
		if envMsg := stream.Msg().GetEnvelopes(); envMsg != nil {
			collected = append(collected, envMsg.GetEnvelopes()...)
		}
	}
	return collected
}

// requireTopicStreamError consumes the stream until it closes and checks the error code.
func requireTopicStreamError(
	t *testing.T,
	stream *connect.ServerStreamForClient[message_api.SubscribeTopicsResponse],
	expectedCode connect.Code,
) {
	t.Helper()

	// Consume any messages until the stream closes.
	for stream.Receive() {
		// Keep consuming.
	}
	err := stream.Err()
	require.Error(t, err)
	require.Equal(t, expectedCode, connect.CodeOf(err))
}

// ---- Validation Tests ----

func TestSubscribeTopics_Validation(t *testing.T) {
	client, _, _ := setupTopicTest(t)

	tooManyFilters := make([]*message_api.SubscribeTopicsRequest_TopicFilter, 10001)
	for i := range tooManyFilters {
		tooManyFilters[i] = makeFilter([]byte{byte(i % 256), byte(i / 256), 1}, nil)
	}

	tests := []struct {
		name    string
		filters []*message_api.SubscribeTopicsRequest_TopicFilter
	}{
		{"NilRequest", nil},
		{"EmptyFilters", []*message_api.SubscribeTopicsRequest_TopicFilter{}},
		{
			"EmptyTopic",
			[]*message_api.SubscribeTopicsRequest_TopicFilter{makeFilter([]byte{}, nil)},
		},
		{
			"TopicTooLong",
			[]*message_api.SubscribeTopicsRequest_TopicFilter{makeFilter(make([]byte, 129), nil)},
		},
		{"TooManyFilters", tooManyFilters},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			stream, err := client.SubscribeTopics(
				t.Context(),
				connect.NewRequest(&message_api.SubscribeTopicsRequest{Filters: tc.filters}),
			)
			require.NoError(t, err)
			requireTopicStreamError(t, stream, connect.CodeInvalidArgument)
		})
	}
}

func TestSubscribeTopics_UnknownOriginatorInCursor(t *testing.T) {
	client, store, _ := setupTopicTest(t)
	payerID := db.NullInt32(testutils.CreatePayer(t, store))

	insertAndWait(t, store, []queries.InsertGatewayEnvelopeParams{
		makeEnvRow(t, 100, 1, topicA, payerID),
	})

	// Reference originator 999 which is not known.
	stream, err := client.SubscribeTopics(
		t.Context(),
		connect.NewRequest(&message_api.SubscribeTopicsRequest{
			Filters: []*message_api.SubscribeTopicsRequest_TopicFilter{
				makeFilter(topicA, map[uint32]uint64{999: 0}),
			},
		}),
	)
	require.NoError(t, err)
	requireTopicStreamError(t, stream, connect.CodeInvalidArgument)
}

// ---- Live-Only Tests (nil LastSeen) ----

func TestSubscribeTopics_LiveOnly(t *testing.T) {
	client, store, _ := setupTopicTest(t)
	payerID := db.NullInt32(testutils.CreatePayer(t, store))

	stream := subscribeTopics(
		t,
		client,
		t.Context(),
		[]*message_api.SubscribeTopicsRequest_TopicFilter{
			makeFilter(topicA, nil),
		},
	)

	// Insert envelopes after subscribing.
	testutils.InsertGatewayEnvelopes(t, store, []queries.InsertGatewayEnvelopeParams{
		makeEnvRow(t, 100, 1, topicA, payerID),
	})

	envs := collectTopicEnvelopes(t, stream, 1)
	require.Len(t, envs, 1)

	decoded := envelopeTestUtils.UnmarshalUnsignedOriginatorEnvelope(
		t, envs[0].GetUnsignedOriginatorEnvelope(),
	)
	require.EqualValues(t, 100, decoded.GetOriginatorNodeId())
	require.EqualValues(t, 1, decoded.GetOriginatorSequenceId())
}

func TestSubscribeTopics_LiveOnlyFiltersByTopic(t *testing.T) {
	client, store, _ := setupTopicTest(t)
	payerID := db.NullInt32(testutils.CreatePayer(t, store))

	stream := subscribeTopics(
		t,
		client,
		t.Context(),
		[]*message_api.SubscribeTopicsRequest_TopicFilter{
			makeFilter(topicA, nil),
			makeFilter(topicB, nil),
		},
	)

	// Insert to topicA, topicB, and topicC.
	testutils.InsertGatewayEnvelopes(t, store, []queries.InsertGatewayEnvelopeParams{
		makeEnvRow(t, 100, 1, topicA, payerID),
		makeEnvRow(t, 100, 2, topicB, payerID),
		makeEnvRow(t, 100, 3, topicC, payerID),
	})

	// Should receive only A and B (2 envelopes).
	envs := collectTopicEnvelopes(t, stream, 2)
	require.Len(t, envs, 2)
}

// ---- Catch-Up Tests ----

func TestSubscribeTopics_CatchUpFromEmpty(t *testing.T) {
	client, store, _ := setupTopicTest(t)
	payerID := db.NullInt32(testutils.CreatePayer(t, store))

	insertAndWait(t, store, []queries.InsertGatewayEnvelopeParams{
		makeEnvRow(t, 100, 1, topicA, payerID),
		makeEnvRow(t, 200, 1, topicA, payerID),
	})

	stream := subscribeTopics(
		t,
		client,
		t.Context(),
		[]*message_api.SubscribeTopicsRequest_TopicFilter{
			makeFilter(topicA, map[uint32]uint64{}),
		},
	)

	envs := collectTopicEnvelopes(t, stream, 2)
	require.Len(t, envs, 2)
}

func TestSubscribeTopics_CatchUpFromCursor(t *testing.T) {
	client, store, _ := setupTopicTest(t)
	payerID := db.NullInt32(testutils.CreatePayer(t, store))

	insertAndWait(t, store, []queries.InsertGatewayEnvelopeParams{
		makeEnvRow(t, 100, 1, topicA, payerID),
		makeEnvRow(t, 100, 2, topicA, payerID),
		makeEnvRow(t, 100, 3, topicA, payerID),
	})

	// Cursor says we've seen seq 1 from originator 100, so we should only get 2 and 3.
	stream := subscribeTopics(
		t,
		client,
		t.Context(),
		[]*message_api.SubscribeTopicsRequest_TopicFilter{
			makeFilter(topicA, map[uint32]uint64{100: 1}),
		},
	)

	envs := collectTopicEnvelopes(t, stream, 2)
	require.Len(t, envs, 2)

	for _, env := range envs {
		decoded := envelopeTestUtils.UnmarshalUnsignedOriginatorEnvelope(
			t, env.GetUnsignedOriginatorEnvelope(),
		)
		require.Greater(t, decoded.GetOriginatorSequenceId(), uint64(1))
	}
}

func TestSubscribeTopics_DifferentCursorsPerTopic(t *testing.T) {
	client, store, _ := setupTopicTest(t)
	payerID := db.NullInt32(testutils.CreatePayer(t, store))

	// topicA: seq 1, 2, 3 from node 100
	// topicB: seq 4 from node 100
	insertAndWait(t, store, []queries.InsertGatewayEnvelopeParams{
		makeEnvRow(t, 100, 1, topicA, payerID),
		makeEnvRow(t, 100, 2, topicA, payerID),
		makeEnvRow(t, 100, 3, topicA, payerID),
		makeEnvRow(t, 100, 4, topicB, payerID),
	})

	// topicA cursor at seq=2 (should get seq 3), topicB cursor at seq=0 (should get seq 4)
	stream := subscribeTopics(
		t,
		client,
		t.Context(),
		[]*message_api.SubscribeTopicsRequest_TopicFilter{
			makeFilter(topicA, map[uint32]uint64{100: 2}),
			makeFilter(topicB, map[uint32]uint64{100: 0}),
		},
	)

	envs := collectTopicEnvelopes(t, stream, 2)
	require.Len(t, envs, 2)

	seqIDs := make(map[uint64]struct{})
	for _, env := range envs {
		decoded := envelopeTestUtils.UnmarshalUnsignedOriginatorEnvelope(
			t, env.GetUnsignedOriginatorEnvelope(),
		)
		seqIDs[decoded.GetOriginatorSequenceId()] = struct{}{}
	}
	require.Contains(t, seqIDs, uint64(3)) // topicA catch-up
	require.Contains(t, seqIDs, uint64(4)) // topicB catch-up
}

func TestSubscribeTopics_CatchUpThenLive(t *testing.T) {
	client, store, _ := setupTopicTest(t)
	payerID := db.NullInt32(testutils.CreatePayer(t, store))

	insertAndWait(t, store, []queries.InsertGatewayEnvelopeParams{
		makeEnvRow(t, 100, 1, topicA, payerID),
	})

	stream := subscribeTopics(
		t,
		client,
		t.Context(),
		[]*message_api.SubscribeTopicsRequest_TopicFilter{
			makeFilter(topicA, map[uint32]uint64{}),
		},
	)

	// Receive catch-up envelope.
	envs := collectTopicEnvelopes(t, stream, 1)
	require.Len(t, envs, 1)
	decoded := envelopeTestUtils.UnmarshalUnsignedOriginatorEnvelope(
		t,
		envs[0].GetUnsignedOriginatorEnvelope(),
	)
	require.EqualValues(t, 1, decoded.GetOriginatorSequenceId())

	// Now insert a new one for live delivery.
	testutils.InsertGatewayEnvelopes(t, store, []queries.InsertGatewayEnvelopeParams{
		makeEnvRow(t, 100, 2, topicA, payerID),
	})

	liveEnvs := collectTopicEnvelopes(t, stream, 1)
	require.Len(t, liveEnvs, 1)
	decoded = envelopeTestUtils.UnmarshalUnsignedOriginatorEnvelope(
		t,
		liveEnvs[0].GetUnsignedOriginatorEnvelope(),
	)
	require.EqualValues(t, 2, decoded.GetOriginatorSequenceId())
}

func TestSubscribeTopics_NoDuplicatesBetweenCatchUpAndLive(t *testing.T) {
	client, store, _ := setupTopicTest(t)
	payerID := db.NullInt32(testutils.CreatePayer(t, store))

	insertAndWait(t, store, []queries.InsertGatewayEnvelopeParams{
		makeEnvRow(t, 100, 1, topicA, payerID),
		makeEnvRow(t, 100, 2, topicA, payerID),
	})

	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()

	stream, err := client.SubscribeTopics(
		ctx,
		connect.NewRequest(&message_api.SubscribeTopicsRequest{
			Filters: []*message_api.SubscribeTopicsRequest_TopicFilter{
				makeFilter(topicA, map[uint32]uint64{}),
			},
		}),
	)
	require.NoError(t, err)

	// Collect all envelopes for a bit and check for duplicates.
	seen := make(map[uint64]struct{})
	deadline := time.After(5 * time.Second)

	for {
		select {
		case <-deadline:
			// Enough time for catch-up + any live duplicate.
			require.GreaterOrEqual(t, len(seen), 2)
			return

		default:
			if !stream.Receive() {
				cancel()
				return
			}
			if envMsg := stream.Msg().GetEnvelopes(); envMsg != nil {
				for _, env := range envMsg.GetEnvelopes() {
					decoded := envelopeTestUtils.UnmarshalUnsignedOriginatorEnvelope(
						t, env.GetUnsignedOriginatorEnvelope(),
					)
					_, dup := seen[decoded.GetOriginatorSequenceId()]
					require.False(
						t,
						dup,
						"received duplicate seqID=%d",
						decoded.GetOriginatorSequenceId(),
					)
					seen[decoded.GetOriginatorSequenceId()] = struct{}{}
				}
			}
			if len(seen) >= 2 {
				return
			}
		}
	}
}

func TestSubscribeTopics_StatusStartedOnOpen(t *testing.T) {
	client, _, _ := setupTopicTest(t)

	stream, err := client.SubscribeTopics(
		t.Context(),
		connect.NewRequest(&message_api.SubscribeTopicsRequest{
			Filters: []*message_api.SubscribeTopicsRequest_TopicFilter{
				makeFilter(topicA, nil),
			},
		}),
	)
	require.NoError(t, err)

	// First message should be a STARTED status.
	require.True(t, stream.Receive())
	msg := stream.Msg()
	require.Equal(t,
		message_api.SubscribeTopicsResponse_SUBSCRIPTION_STATUS_STARTED,
		msg.GetStatusUpdate().GetStatus(),
	)
}

func TestSubscribeTopics_StatusLifecycle(t *testing.T) {
	client, store, _ := setupTopicTest(t)
	payerID := db.NullInt32(testutils.CreatePayer(t, store))

	// Insert envelopes before subscribing so catch-up is triggered.
	insertAndWait(t, store, []queries.InsertGatewayEnvelopeParams{
		makeEnvRow(t, 100, 1, topicA, payerID),
		makeEnvRow(t, 100, 2, topicA, payerID),
	})

	stream, err := client.SubscribeTopics(
		t.Context(),
		connect.NewRequest(&message_api.SubscribeTopicsRequest{
			Filters: []*message_api.SubscribeTopicsRequest_TopicFilter{
				makeFilter(topicA, map[uint32]uint64{}),
			},
		}),
	)
	require.NoError(t, err)

	// 1. First message: STARTED
	require.True(t, stream.Receive())
	require.Equal(t,
		message_api.SubscribeTopicsResponse_SUBSCRIPTION_STATUS_STARTED,
		stream.Msg().GetStatusUpdate().GetStatus(),
	)

	// 2. Catch-up envelopes
	var catchUpEnvs []*envelopes.OriginatorEnvelope
	for {
		require.True(t, stream.Receive())
		msg := stream.Msg()
		if envMsg := msg.GetEnvelopes(); envMsg != nil {
			catchUpEnvs = append(catchUpEnvs, envMsg.GetEnvelopes()...)
			continue
		}
		// Must be the CATCHUP_COMPLETE status.
		require.Equal(t,
			message_api.SubscribeTopicsResponse_SUBSCRIPTION_STATUS_CATCHUP_COMPLETE,
			msg.GetStatusUpdate().GetStatus(),
		)
		break
	}
	require.Len(t, catchUpEnvs, 2)

	// 3. Insert a new envelope for live delivery.
	testutils.InsertGatewayEnvelopes(t, store, []queries.InsertGatewayEnvelopeParams{
		makeEnvRow(t, 100, 3, topicA, payerID),
	})

	// 4. Next message with envelopes should be the live envelope.
	for {
		require.True(t, stream.Receive())
		msg := stream.Msg()
		if envMsg := msg.GetEnvelopes(); envMsg != nil {
			require.Len(t, envMsg.GetEnvelopes(), 1)
			decoded := envelopeTestUtils.UnmarshalUnsignedOriginatorEnvelope(
				t, envMsg.GetEnvelopes()[0].GetUnsignedOriginatorEnvelope(),
			)
			require.EqualValues(t, 3, decoded.GetOriginatorSequenceId())
			return
		}
		// Skip WAITING keepalives.
		require.Equal(t,
			message_api.SubscribeTopicsResponse_SUBSCRIPTION_STATUS_WAITING,
			msg.GetStatusUpdate().GetStatus(),
		)
	}
}

// ---- Ordering Tests ----

func TestSubscribeTopics_PerOriginatorOrdering(t *testing.T) {
	client, store, _ := setupTopicTest(t)
	payerID := db.NullInt32(testutils.CreatePayer(t, store))

	// Insert multiple envelopes from two originators.
	rows := make([]queries.InsertGatewayEnvelopeParams, 0, 10)
	for seq := uint64(1); seq <= 5; seq++ {
		rows = append(rows,
			makeEnvRow(t, 100, seq, topicA, payerID),
			makeEnvRow(t, 200, seq, topicA, payerID),
		)
	}
	insertAndWait(t, store, rows)

	stream := subscribeTopics(
		t,
		client,
		t.Context(),
		[]*message_api.SubscribeTopicsRequest_TopicFilter{
			makeFilter(topicA, map[uint32]uint64{}),
		},
	)

	envs := collectTopicEnvelopes(t, stream, 10)
	require.Len(t, envs, 10)
	requireOriginatorOrdering(t, envs)
}

func TestSubscribeTopics_MultiOriginatorMultiTopic(t *testing.T) {
	nodes := []registry.Node{
		{NodeID: 100, IsCanonical: true},
		{NodeID: 200, IsCanonical: true},
		{NodeID: 300, IsCanonical: true},
	}
	suite := testUtilsApi.NewTestAPIServer(t, testUtilsApi.WithRegistryNodes(nodes))
	payerID := db.NullInt32(testutils.CreatePayer(t, suite.DB))

	topics := [][]byte{topicA, topicB, topicC}
	nodeIDList := []uint32{100, 200, 300}

	rows := make([]queries.InsertGatewayEnvelopeParams, 0)
	expectedTotal := 0
	seq := uint64(1)
	for _, tp := range topics {
		for _, nid := range nodeIDList {
			for range 3 {
				rows = append(rows, makeEnvRow(t, nid, seq, tp, payerID))
				seq++
				expectedTotal++
			}
		}
	}
	insertAndWait(t, suite.DB, rows)

	// Subscribe to all topics with empty cursors.
	filters := make([]*message_api.SubscribeTopicsRequest_TopicFilter, len(topics))
	for i, tp := range topics {
		filters[i] = makeFilter(tp, map[uint32]uint64{})
	}

	stream := subscribeTopics(t, suite.ClientReplication, t.Context(), filters)

	envs := collectTopicEnvelopes(t, stream, expectedTotal)
	require.Len(t, envs, expectedTotal)
	requireOriginatorOrdering(t, envs)
}

// ---- Scale Tests ----

func TestSubscribeTopics_LargeCatchUpMultiplePages(t *testing.T) {
	client, store, _ := setupTopicTest(t)
	payerID := db.NullInt32(testutils.CreatePayer(t, store))

	// Insert >500 envelopes to trigger pagination.
	total := 600
	rows := make([]queries.InsertGatewayEnvelopeParams, total)
	for i := range total {
		rows[i] = makeEnvRow(t, 100, uint64(i+1), topicA, payerID)
	}
	insertAndWait(t, store, rows)

	stream := subscribeTopics(
		t,
		client,
		t.Context(),
		[]*message_api.SubscribeTopicsRequest_TopicFilter{
			makeFilter(topicA, map[uint32]uint64{}),
		},
	)

	envs := collectTopicEnvelopes(t, stream, total)
	require.Len(t, envs, total)
}

func TestSubscribeTopics_ManyTopics(t *testing.T) {
	client, store, _ := setupTopicTest(t)
	payerID := db.NullInt32(testutils.CreatePayer(t, store))

	numTopics := 1000
	topics := make([][]byte, numTopics)
	filters := make([]*message_api.SubscribeTopicsRequest_TopicFilter, numTopics)
	rows := make([]queries.InsertGatewayEnvelopeParams, numTopics)

	for i := range numTopics {
		tp := topic.NewTopic(topic.TopicKindGroupMessagesV1, fmt.Appendf(nil, "many-topic-%d", i)).
			Bytes()
		topics[i] = tp
		filters[i] = makeFilter(tp, map[uint32]uint64{})
		rows[i] = makeEnvRow(t, 100, uint64(i+1), tp, payerID)
	}
	insertAndWait(t, store, rows)

	stream := subscribeTopics(t, client, t.Context(), filters)

	envs := collectTopicEnvelopes(t, stream, numTopics)
	require.Len(t, envs, numTopics)
}

// ---- Error Path Tests ----

func TestSubscribeTopics_ContextCancelledDuringLive(t *testing.T) {
	client, _, _ := setupTopicTest(t)
	ctx, cancel := context.WithCancel(t.Context())

	stream := subscribeTopics(t, client, ctx, []*message_api.SubscribeTopicsRequest_TopicFilter{
		makeFilter(topicA, nil),
	})

	cancel()

	// Stream should close cleanly.
	for stream.Receive() {
		// Drain.
	}

	err := stream.Err()
	if err != nil {
		// Either nil or cancelled is acceptable.
		require.True(
			t,
			errors.Is(err, context.Canceled) || connect.CodeOf(err) == connect.CodeCanceled,
		)
	}
}

func TestSubscribeTopics_ContextCancelledDuringCatchUp(t *testing.T) {
	client, store, _ := setupTopicTest(t)
	payerID := db.NullInt32(testutils.CreatePayer(t, store))

	// Insert a large amount of data to make catch-up take time.
	rows := make([]queries.InsertGatewayEnvelopeParams, 200)
	for i := range 200 {
		rows[i] = makeEnvRow(t, 100, uint64(i+1), topicA, payerID)
	}
	insertAndWait(t, store, rows)

	ctx, cancel := context.WithCancel(t.Context())

	stream, err := client.SubscribeTopics(
		ctx,
		connect.NewRequest(&message_api.SubscribeTopicsRequest{
			Filters: []*message_api.SubscribeTopicsRequest_TopicFilter{
				makeFilter(topicA, map[uint32]uint64{}),
			},
		}),
	)
	require.NoError(t, err)

	// Cancel quickly.
	cancel()

	// Drain stream.
	for stream.Receive() {
	}

	// Stream should close with an error or nil.
	err = stream.Err()
	if err != nil {
		require.True(t,
			errors.Is(err, context.Canceled) ||
				connect.CodeOf(err) == connect.CodeCanceled ||
				connect.CodeOf(err) == connect.CodeInternal)
	}
}

// ---- Concurrent Tests ----

func TestSubscribeTopics_SimultaneousWithSubscribeEnvelopes(t *testing.T) {
	client, store, _ := setupTopicTest(t)
	payerID := db.NullInt32(testutils.CreatePayer(t, store))

	ctx := t.Context()

	topicStream := subscribeTopics(
		t,
		client,
		ctx,
		[]*message_api.SubscribeTopicsRequest_TopicFilter{
			makeFilter(topicA, nil),
		},
	)

	envelopeStream, err := client.SubscribeEnvelopes(
		ctx,
		connect.NewRequest(&message_api.SubscribeEnvelopesRequest{
			Query: &message_api.EnvelopesQuery{
				Topics:   []db.Topic{topicA},
				LastSeen: nil,
			},
		}),
	)
	require.NoError(t, err)

	// Consume keepalive for SubscribeEnvelopes.
	require.True(t, envelopeStream.Receive())

	// Insert an envelope.
	testutils.InsertGatewayEnvelopes(t, store, []queries.InsertGatewayEnvelopeParams{
		makeEnvRow(t, 100, 1, topicA, payerID),
	})

	// Both streams should receive the envelope.
	topicEnvs := collectTopicEnvelopes(t, topicStream, 1)
	require.Len(t, topicEnvs, 1)

	// For the SubscribeEnvelopes stream.
	require.True(t, envelopeStream.Receive())
	msg := envelopeStream.Msg()
	require.NotEmpty(t, msg.GetEnvelopes())
}
