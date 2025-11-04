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
	topicA = topic.NewTopic(topic.TopicKindGroupMessagesV1, []byte("topicA")).Bytes()
	topicB = topic.NewTopic(topic.TopicKindGroupMessagesV1, []byte("topicB")).Bytes()
	topicC = topic.NewTopic(topic.TopicKindGroupMessagesV1, []byte("topicC")).Bytes()
)

func setupQueryTest(t *testing.T, dbHandle *sql.DB) []queries.InsertGatewayEnvelopeParams {
	payerId := db.NullInt32(testutils.CreatePayer(t, dbHandle))
	dbRows := []queries.InsertGatewayEnvelopeParams{
		{
			OriginatorNodeID:     100,
			OriginatorSequenceID: 1,
			Topic:                topicA,
			PayerID:              payerId,
			OriginatorEnvelope: testutils.Marshal(
				t,
				envelopeTestUtils.CreateOriginatorEnvelopeWithTopic(t, 100, 1, topicA),
			),
		},
		{
			OriginatorNodeID:     200,
			OriginatorSequenceID: 1,
			Topic:                topicA,
			PayerID:              payerId,
			OriginatorEnvelope: testutils.Marshal(
				t,
				envelopeTestUtils.CreateOriginatorEnvelopeWithTopic(t, 200, 1, topicA),
			),
		},
		{
			OriginatorNodeID:     100,
			OriginatorSequenceID: 2,
			Topic:                topicB,
			PayerID:              payerId,
			OriginatorEnvelope: testutils.Marshal(
				t,
				envelopeTestUtils.CreateOriginatorEnvelopeWithTopic(t, 100, 2, topicB),
			),
		},
		{
			OriginatorNodeID:     200,
			OriginatorSequenceID: 2,
			Topic:                topicB,
			PayerID:              payerId,
			OriginatorEnvelope: testutils.Marshal(
				t,
				envelopeTestUtils.CreateOriginatorEnvelopeWithTopic(t, 200, 2, topicB),
			),
		},
		{
			OriginatorNodeID:     100,
			OriginatorSequenceID: 3,
			Topic:                topicA,
			PayerID:              payerId,
			OriginatorEnvelope: testutils.Marshal(
				t,
				envelopeTestUtils.CreateOriginatorEnvelopeWithTopic(t, 100, 3, topicA),
			),
		},
	}
	testutils.InsertGatewayEnvelopes(t, dbHandle, dbRows)
	return dbRows
}

func TestQueryAllEnvelopes(t *testing.T) {
	api, dbHandle, _ := apiTestUtils.NewTestReplicationAPIClient(t)
	dbRows := setupQueryTest(t, dbHandle)

	resp, err := api.QueryEnvelopes(
		context.Background(),
		&message_api.QueryEnvelopesRequest{
			Query: &message_api.EnvelopesQuery{},
			Limit: 0,
		},
	)
	require.NoError(t, err)
	checkRowsMatchProtos(t, dbRows, []int{0, 1, 2, 3, 4}, resp.GetEnvelopes())
}

func TestQueryPagedEnvelopes(t *testing.T) {
	api, dbHandle, _ := apiTestUtils.NewTestReplicationAPIClient(t)
	dbRows := setupQueryTest(t, dbHandle)

	resp, err := api.QueryEnvelopes(
		context.Background(),
		&message_api.QueryEnvelopesRequest{
			Query: &message_api.EnvelopesQuery{},
			Limit: 2,
		},
	)
	require.NoError(t, err)
	checkRowsMatchProtos(t, dbRows, []int{0, 1}, resp.GetEnvelopes())
}

func TestQueryEnvelopesByOriginator(t *testing.T) {
	api, dbHandle, _ := apiTestUtils.NewTestReplicationAPIClient(t)
	dbRows := setupQueryTest(t, dbHandle)

	resp, err := api.QueryEnvelopes(
		context.Background(),
		&message_api.QueryEnvelopesRequest{
			Query: &message_api.EnvelopesQuery{
				OriginatorNodeIds: []uint32{200},
				LastSeen:          nil,
			},
			Limit: 0,
		},
	)
	require.NoError(t, err)
	checkRowsMatchProtos(t, dbRows, []int{1, 3}, resp.GetEnvelopes())
}

func TestQueryEnvelopesByTopic(t *testing.T) {
	api, dbHandle, _ := apiTestUtils.NewTestReplicationAPIClient(t)
	dbRows := setupQueryTest(t, dbHandle)

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
	checkRowsMatchProtos(t, dbRows, []int{0, 1, 4}, resp.GetEnvelopes())
}

func TestQueryEnvelopesFromLastSeen(t *testing.T) {
	api, dbHandle, _ := apiTestUtils.NewTestReplicationAPIClient(t)
	dbRows := setupQueryTest(t, dbHandle)

	resp, err := api.QueryEnvelopes(
		context.Background(),
		&message_api.QueryEnvelopesRequest{
			Query: &message_api.EnvelopesQuery{
				LastSeen: &envelopes.Cursor{NodeIdToSequenceId: map[uint32]uint64{100: 2}},
			},
			Limit: 0,
		},
	)
	require.NoError(t, err)
	checkRowsMatchProtos(t, dbRows, []int{1, 3, 4}, resp.GetEnvelopes())
}

func TestQueryTopicFromLastSeen(t *testing.T) {
	api, store, _ := apiTestUtils.NewTestReplicationAPIClient(t)
	dbRows := setupQueryTest(t, store)

	resp, err := api.QueryEnvelopes(
		context.Background(),
		&message_api.QueryEnvelopesRequest{
			Query: &message_api.EnvelopesQuery{
				Topics: []db.Topic{topicA},
				LastSeen: &envelopes.Cursor{
					NodeIdToSequenceId: map[uint32]uint64{100: 2, 200: 1},
				},
			},
			Limit: 0,
		},
	)
	require.NoError(t, err)
	checkRowsMatchProtos(t, dbRows, []int{4}, resp.GetEnvelopes())
}

func TestQueryMultipleTopicsFromLastSeen(t *testing.T) {
	api, dbHandle, _ := apiTestUtils.NewTestReplicationAPIClient(t)
	dbRows := setupQueryTest(t, dbHandle)

	resp, err := api.QueryEnvelopes(
		context.Background(),
		&message_api.QueryEnvelopesRequest{
			Query: &message_api.EnvelopesQuery{
				Topics: []db.Topic{topicA, topicB},
				LastSeen: &envelopes.Cursor{
					NodeIdToSequenceId: map[uint32]uint64{100: 2, 200: 1},
				},
			},
			Limit: 0,
		},
	)
	require.NoError(t, err)
	checkRowsMatchProtos(t, dbRows, []int{3, 4}, resp.GetEnvelopes())
}

func TestQueryMultipleOriginatorsFromLastSeen(t *testing.T) {
	api, dbHandle, _ := apiTestUtils.NewTestReplicationAPIClient(t)
	dbRows := setupQueryTest(t, dbHandle)

	resp, err := api.QueryEnvelopes(
		context.Background(),
		&message_api.QueryEnvelopesRequest{
			Query: &message_api.EnvelopesQuery{
				OriginatorNodeIds: []uint32{100, 200},
				LastSeen: &envelopes.Cursor{
					NodeIdToSequenceId: map[uint32]uint64{100: 1, 200: 1},
				},
			},
			Limit: 0,
		},
	)
	require.NoError(t, err)
	checkRowsMatchProtos(t, dbRows, []int{2, 3, 4}, resp.GetEnvelopes())
}

func TestQueryEnvelopesWithEmptyResult(t *testing.T) {
	api, dbHandle, _ := apiTestUtils.NewTestReplicationAPIClient(t)
	dbRows := setupQueryTest(t, dbHandle)

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
	checkRowsMatchProtos(t, dbRows, []int{}, resp.GetEnvelopes())
}

func TestInvalidQuery(t *testing.T) {
	api, dbHandle, _ := apiTestUtils.NewTestReplicationAPIClient(t)
	setupQueryTest(t, dbHandle)

	_, err := api.QueryEnvelopes(
		context.Background(),
		&message_api.QueryEnvelopesRequest{
			Query: &message_api.EnvelopesQuery{
				Topics:            []db.Topic{topicA},
				OriginatorNodeIds: []uint32{100},
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
