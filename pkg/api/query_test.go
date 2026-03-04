package api_test

import (
	"context"
	"database/sql"
	"testing"

	"connectrpc.com/connect"
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
	payerID := db.NullInt32(testutils.CreatePayer(t, dbHandle))
	dbRows := []queries.InsertGatewayEnvelopeParams{
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
	testutils.InsertGatewayEnvelopes(t, dbHandle, dbRows)
	return dbRows
}

func TestQueryAllEnvelopes(t *testing.T) {
	suite := apiTestUtils.NewTestAPIServer(t)
	dbRows := setupQueryTest(t, suite.DB)
	resp, err := suite.ClientReplication.QueryEnvelopes(
		context.Background(),
		connect.NewRequest(&message_api.QueryEnvelopesRequest{
			Query: &message_api.EnvelopesQuery{Topics: [][]byte{topicA, topicB, topicC}},
			Limit: 0,
		}),
	)
	require.NoError(t, err)

	checkRowsMatchProtos(t, dbRows, []int{0, 2, 4, 1, 3}, resp.Msg.GetEnvelopes())
}

func TestQueryPagedEnvelopes(t *testing.T) {
	suite := apiTestUtils.NewTestAPIServer(t)
	dbRows := setupQueryTest(t, suite.DB)

	resp, err := suite.ClientReplication.QueryEnvelopes(
		context.Background(),
		connect.NewRequest(&message_api.QueryEnvelopesRequest{
			Query: &message_api.EnvelopesQuery{Topics: [][]byte{topicA}},
			Limit: 2,
		}),
	)
	require.NoError(t, err)

	checkRowsMatchProtos(t, dbRows, []int{0, 4}, resp.Msg.GetEnvelopes())
}

func TestQueryEnvelopesByOriginator(t *testing.T) {
	suite := apiTestUtils.NewTestAPIServer(t)
	dbRows := setupQueryTest(t, suite.DB)

	resp, err := suite.ClientReplication.QueryEnvelopes(
		context.Background(),
		connect.NewRequest(&message_api.QueryEnvelopesRequest{
			Query: &message_api.EnvelopesQuery{
				OriginatorNodeIds: []uint32{200},
				LastSeen:          nil,
			},
			Limit: 0,
		}),
	)
	require.NoError(t, err)

	checkRowsMatchProtos(t, dbRows, []int{1, 3}, resp.Msg.GetEnvelopes())
}

func TestQueryEnvelopesByTopic(t *testing.T) {
	suite := apiTestUtils.NewTestAPIServer(t)
	dbRows := setupQueryTest(t, suite.DB)

	resp, err := suite.ClientReplication.QueryEnvelopes(
		context.Background(),
		connect.NewRequest(&message_api.QueryEnvelopesRequest{
			Query: &message_api.EnvelopesQuery{
				Topics:   []db.Topic{topicA},
				LastSeen: nil,
			},
			Limit: 0,
		}),
	)
	require.NoError(t, err)

	checkRowsMatchProtos(t, dbRows, []int{0, 4, 1}, resp.Msg.GetEnvelopes())
}

func TestQueryEnvelopesFromLastSeen(t *testing.T) {
	suite := apiTestUtils.NewTestAPIServer(t)
	dbRows := setupQueryTest(t, suite.DB)

	resp, err := suite.ClientReplication.QueryEnvelopes(
		context.Background(),
		connect.NewRequest(&message_api.QueryEnvelopesRequest{
			Query: &message_api.EnvelopesQuery{
				Topics:   [][]byte{topicA, topicB, topicC},
				LastSeen: &envelopes.Cursor{NodeIdToSequenceId: map[uint32]uint64{100: 2}},
			},
			Limit: 0,
		}),
	)
	require.NoError(t, err)

	checkRowsMatchProtos(t, dbRows, []int{4, 1, 3}, resp.Msg.GetEnvelopes())
}

func TestQueryTopicFromLastSeen(t *testing.T) {
	suite := apiTestUtils.NewTestAPIServer(t)
	dbRows := setupQueryTest(t, suite.DB)

	resp, err := suite.ClientReplication.QueryEnvelopes(
		context.Background(),
		connect.NewRequest(&message_api.QueryEnvelopesRequest{
			Query: &message_api.EnvelopesQuery{
				Topics: []db.Topic{topicA},
				LastSeen: &envelopes.Cursor{
					NodeIdToSequenceId: map[uint32]uint64{100: 2, 200: 1},
				},
			},
			Limit: 0,
		}),
	)
	require.NoError(t, err)

	checkRowsMatchProtos(t, dbRows, []int{4}, resp.Msg.GetEnvelopes())
}

func TestQueryMultipleTopicsFromLastSeen(t *testing.T) {
	suite := apiTestUtils.NewTestAPIServer(t)
	dbRows := setupQueryTest(t, suite.DB)

	resp, err := suite.ClientReplication.QueryEnvelopes(
		context.Background(),
		connect.NewRequest(&message_api.QueryEnvelopesRequest{
			Query: &message_api.EnvelopesQuery{
				Topics: []db.Topic{topicA, topicB},
				LastSeen: &envelopes.Cursor{
					NodeIdToSequenceId: map[uint32]uint64{100: 2, 200: 1},
				},
			},
			Limit: 0,
		}),
	)
	require.NoError(t, err)

	checkRowsMatchProtos(t, dbRows, []int{4, 3}, resp.Msg.GetEnvelopes())
}

func TestQueryMultipleOriginatorsFromLastSeen(t *testing.T) {
	suite := apiTestUtils.NewTestAPIServer(t)
	dbRows := setupQueryTest(t, suite.DB)

	resp, err := suite.ClientReplication.QueryEnvelopes(
		context.Background(),
		connect.NewRequest(&message_api.QueryEnvelopesRequest{
			Query: &message_api.EnvelopesQuery{
				OriginatorNodeIds: []uint32{100, 200},
				LastSeen: &envelopes.Cursor{
					NodeIdToSequenceId: map[uint32]uint64{100: 1, 200: 1},
				},
			},
			Limit: 0,
		}),
	)
	require.NoError(t, err)

	checkRowsMatchProtos(t, dbRows, []int{2, 4, 3}, resp.Msg.GetEnvelopes())
}

func TestQueryEnvelopesWithEmptyResult(t *testing.T) {
	suite := apiTestUtils.NewTestAPIServer(t)
	dbRows := setupQueryTest(t, suite.DB)

	resp, err := suite.ClientReplication.QueryEnvelopes(
		context.Background(),
		connect.NewRequest(&message_api.QueryEnvelopesRequest{
			Query: &message_api.EnvelopesQuery{
				Topics: []db.Topic{topicC},
			},
			Limit: 0,
		}),
	)
	require.NoError(t, err)
	checkRowsMatchProtos(t, dbRows, []int{}, resp.Msg.GetEnvelopes())
}

func TestInvalidQuery(t *testing.T) {
	suite := apiTestUtils.NewTestAPIServer(t)

	_, err := suite.ClientReplication.QueryEnvelopes(
		context.Background(),
		connect.NewRequest(&message_api.QueryEnvelopesRequest{
			Query: &message_api.EnvelopesQuery{
				Topics:            []db.Topic{topicA},
				OriginatorNodeIds: []uint32{100},
			},
			Limit: 0,
		}),
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
