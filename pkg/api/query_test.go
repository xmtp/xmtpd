package api_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/db"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/envelopes"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api"
	"github.com/xmtp/xmtpd/pkg/testutils"
	apiTestUtils "github.com/xmtp/xmtpd/pkg/testutils/api"
	envelopeTestUtils "github.com/xmtp/xmtpd/pkg/testutils/envelopes"
	"github.com/xmtp/xmtpd/pkg/topic"
)

var (
	topicA = topic.NewTopic(topic.TOPIC_KIND_GROUP_MESSAGES_V1, []byte("topicA")).Bytes()
	topicB = topic.NewTopic(topic.TOPIC_KIND_GROUP_MESSAGES_V1, []byte("topicB")).Bytes()
	topicC = topic.NewTopic(topic.TOPIC_KIND_GROUP_MESSAGES_V1, []byte("topicC")).Bytes()
)

func setupQueryTest(t *testing.T, db *sql.DB) []queries.InsertGatewayEnvelopeParams {
	db_rows := []queries.InsertGatewayEnvelopeParams{
		{
			OriginatorNodeID:     1,
			OriginatorSequenceID: 1,
			Topic:                topicA,
			OriginatorEnvelope: testutils.Marshal(
				t,
				envelopeTestUtils.CreateOriginatorEnvelopeWithTopic(t, 1, 1, topicA),
			),
		},
		{
			OriginatorNodeID:     2,
			OriginatorSequenceID: 1,
			Topic:                topicA,
			OriginatorEnvelope: testutils.Marshal(
				t,
				envelopeTestUtils.CreateOriginatorEnvelopeWithTopic(t, 2, 1, topicA),
			),
		},
		{
			OriginatorNodeID:     1,
			OriginatorSequenceID: 2,
			Topic:                topicB,
			OriginatorEnvelope: testutils.Marshal(
				t,
				envelopeTestUtils.CreateOriginatorEnvelopeWithTopic(t, 1, 2, topicB),
			),
		},
		{
			OriginatorNodeID:     2,
			OriginatorSequenceID: 2,
			Topic:                topicB,
			OriginatorEnvelope: testutils.Marshal(
				t,
				envelopeTestUtils.CreateOriginatorEnvelopeWithTopic(t, 2, 2, topicB),
			),
		},
		{
			OriginatorNodeID:     1,
			OriginatorSequenceID: 3,
			Topic:                topicA,
			OriginatorEnvelope: testutils.Marshal(
				t,
				envelopeTestUtils.CreateOriginatorEnvelopeWithTopic(t, 1, 3, topicA),
			),
		},
	}
	testutils.InsertGatewayEnvelopes(t, db, db_rows)
	return db_rows
}

func TestQueryAllEnvelopes(t *testing.T) {
	api, db, _, cleanup := apiTestUtils.NewTestReplicationAPIClient(t)
	defer cleanup()
	db_rows := setupQueryTest(t, db)

	resp, err := api.QueryEnvelopes(
		context.Background(),
		&message_api.QueryEnvelopesRequest{
			Query: &message_api.EnvelopesQuery{},
			Limit: 0,
		},
	)
	require.NoError(t, err)
	checkRowsMatchProtos(t, db_rows, []int{0, 1, 2, 3, 4}, resp.GetEnvelopes())
}

func TestQueryPagedEnvelopes(t *testing.T) {
	api, db, _, cleanup := apiTestUtils.NewTestReplicationAPIClient(t)
	defer cleanup()
	db_rows := setupQueryTest(t, db)

	resp, err := api.QueryEnvelopes(
		context.Background(),
		&message_api.QueryEnvelopesRequest{
			Query: &message_api.EnvelopesQuery{},
			Limit: 2,
		},
	)
	require.NoError(t, err)
	checkRowsMatchProtos(t, db_rows, []int{0, 1}, resp.GetEnvelopes())
}

func TestQueryEnvelopesByOriginator(t *testing.T) {
	api, db, _, cleanup := apiTestUtils.NewTestReplicationAPIClient(t)
	defer cleanup()
	db_rows := setupQueryTest(t, db)

	resp, err := api.QueryEnvelopes(
		context.Background(),
		&message_api.QueryEnvelopesRequest{
			Query: &message_api.EnvelopesQuery{
				OriginatorNodeIds: []uint32{2},
				LastSeen:          nil,
			},
			Limit: 0,
		},
	)
	require.NoError(t, err)
	checkRowsMatchProtos(t, db_rows, []int{1, 3}, resp.GetEnvelopes())
}

func TestQueryEnvelopesByTopic(t *testing.T) {
	api, store, _, cleanup := apiTestUtils.NewTestReplicationAPIClient(t)
	defer cleanup()
	db_rows := setupQueryTest(t, store)

	resp, err := api.QueryEnvelopes(
		context.Background(),
		&message_api.QueryEnvelopesRequest{
			Query: &message_api.EnvelopesQuery{
				Topics:   []db.Topic{topicA},
				LastSeen: nil,
			},
			Limit: 0,
		},
	)
	require.NoError(t, err)
	checkRowsMatchProtos(t, db_rows, []int{0, 1, 4}, resp.GetEnvelopes())
}

func TestQueryEnvelopesFromLastSeen(t *testing.T) {
	api, db, _, cleanup := apiTestUtils.NewTestReplicationAPIClient(t)
	defer cleanup()
	db_rows := setupQueryTest(t, db)

	resp, err := api.QueryEnvelopes(
		context.Background(),
		&message_api.QueryEnvelopesRequest{
			Query: &message_api.EnvelopesQuery{
				LastSeen: &envelopes.Cursor{NodeIdToSequenceId: map[uint32]uint64{1: 2}},
			},
			Limit: 0,
		},
	)
	require.NoError(t, err)
	checkRowsMatchProtos(t, db_rows, []int{1, 3, 4}, resp.GetEnvelopes())
}

func TestQueryTopicFromLastSeen(t *testing.T) {
	api, store, _, cleanup := apiTestUtils.NewTestReplicationAPIClient(t)
	defer cleanup()
	db_rows := setupQueryTest(t, store)

	resp, err := api.QueryEnvelopes(
		context.Background(),
		&message_api.QueryEnvelopesRequest{
			Query: &message_api.EnvelopesQuery{
				Topics: []db.Topic{topicA},
				LastSeen: &envelopes.Cursor{
					NodeIdToSequenceId: map[uint32]uint64{1: 2, 2: 1},
				},
			},
			Limit: 0,
		},
	)
	require.NoError(t, err)
	checkRowsMatchProtos(t, db_rows, []int{4}, resp.GetEnvelopes())
}

func TestQueryMultipleTopicsFromLastSeen(t *testing.T) {
	api, store, _, cleanup := apiTestUtils.NewTestReplicationAPIClient(t)
	defer cleanup()
	db_rows := setupQueryTest(t, store)

	resp, err := api.QueryEnvelopes(
		context.Background(),
		&message_api.QueryEnvelopesRequest{
			Query: &message_api.EnvelopesQuery{
				Topics: []db.Topic{topicA, topicB},
				LastSeen: &envelopes.Cursor{
					NodeIdToSequenceId: map[uint32]uint64{1: 2, 2: 1},
				},
			},
			Limit: 0,
		},
	)
	require.NoError(t, err)
	checkRowsMatchProtos(t, db_rows, []int{3, 4}, resp.GetEnvelopes())
}

func TestQueryMultipleOriginatorsFromLastSeen(t *testing.T) {
	api, store, _, cleanup := apiTestUtils.NewTestReplicationAPIClient(t)
	defer cleanup()
	db_rows := setupQueryTest(t, store)

	resp, err := api.QueryEnvelopes(
		context.Background(),
		&message_api.QueryEnvelopesRequest{
			Query: &message_api.EnvelopesQuery{
				OriginatorNodeIds: []uint32{1, 2},
				LastSeen: &envelopes.Cursor{
					NodeIdToSequenceId: map[uint32]uint64{1: 1, 2: 1},
				},
			},
			Limit: 0,
		},
	)
	require.NoError(t, err)
	checkRowsMatchProtos(t, db_rows, []int{2, 3, 4}, resp.GetEnvelopes())
}

func TestQueryEnvelopesWithEmptyResult(t *testing.T) {
	api, store, _, cleanup := apiTestUtils.NewTestReplicationAPIClient(t)
	defer cleanup()
	db_rows := setupQueryTest(t, store)

	resp, err := api.QueryEnvelopes(
		context.Background(),
		&message_api.QueryEnvelopesRequest{
			Query: &message_api.EnvelopesQuery{
				Topics: []db.Topic{topicC},
			},
			Limit: 0,
		},
	)
	require.NoError(t, err)
	checkRowsMatchProtos(t, db_rows, []int{}, resp.GetEnvelopes())
}

func TestInvalidQuery(t *testing.T) {
	api, store, _, cleanup := apiTestUtils.NewTestReplicationAPIClient(t)
	defer cleanup()
	_ = setupQueryTest(t, store)

	_, err := api.QueryEnvelopes(
		context.Background(),
		&message_api.QueryEnvelopesRequest{
			Query: &message_api.EnvelopesQuery{
				Topics:            []db.Topic{topicA},
				OriginatorNodeIds: []uint32{1},
			},
			Limit: 0,
		},
	)
	require.Error(t, err)
}

func checkRowsMatchProtos(
	t *testing.T,
	allRows []queries.InsertGatewayEnvelopeParams,
	matchingIndices []int,
	protos []*envelopes.OriginatorEnvelope,
) {
	require.Len(t, protos, len(matchingIndices))
	for i, p := range protos {
		row := allRows[matchingIndices[i]]
		require.Equal(t, row.OriginatorEnvelope, testutils.Marshal(t, p))
	}
}
