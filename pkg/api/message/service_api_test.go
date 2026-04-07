package message_test

import (
	"context"
	"testing"

	"connectrpc.com/connect"
	"github.com/stretchr/testify/require"
	dbPkg "github.com/xmtp/xmtpd/pkg/db"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api"
	payer_api "github.com/xmtp/xmtpd/pkg/proto/xmtpv4/payer_api"
	"github.com/xmtp/xmtpd/pkg/testutils"
	apiTestUtils "github.com/xmtp/xmtpd/pkg/testutils/api"
	envelopeTestUtils "github.com/xmtp/xmtpd/pkg/testutils/envelopes"
	"github.com/xmtp/xmtpd/pkg/topic"
)

// TestQueryApi_QueryEnvelopes verifies that QueryEnvelopes is reachable via the QueryApi client.
func TestQueryApi_QueryEnvelopes(t *testing.T) {
	suite := apiTestUtils.NewTestAPIServer(t)

	payerID := dbPkg.NullInt32(testutils.CreatePayer(t, suite.DB))
	topicBytes := topic.NewTopic(topic.TopicKindGroupMessagesV1, []byte("test-topic")).Bytes()

	envBytes := testutils.Marshal(
		t,
		envelopeTestUtils.CreateOriginatorEnvelopeWithTopic(t, 100, 1, topicBytes),
	)
	testutils.InsertGatewayEnvelopes(t, suite.DB, []queries.InsertGatewayEnvelopeV3Params{
		{
			OriginatorNodeID:     100,
			OriginatorSequenceID: 1,
			Topic:                topicBytes,
			PayerID:              payerID,
			OriginatorEnvelope:   envBytes,
		},
	})

	resp, err := suite.ClientQuery.QueryEnvelopes(
		context.Background(),
		connect.NewRequest(&message_api.QueryEnvelopesRequest{
			Query: &message_api.EnvelopesQuery{
				Topics: [][]byte{topicBytes},
			},
		}),
	)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotEmpty(t, resp.Msg.GetEnvelopes())
}

// TestQueryApi_GetInboxIds verifies that GetInboxIds is reachable via the QueryApi client.
// An empty request returns no results (not an "unimplemented" error).
func TestQueryApi_GetInboxIds(t *testing.T) {
	suite := apiTestUtils.NewTestAPIServer(t)

	resp, err := suite.ClientQuery.GetInboxIds(
		context.Background(),
		connect.NewRequest(&message_api.GetInboxIdsRequest{}),
	)
	require.NoError(t, err)
	require.NotNil(t, resp)
}

// TestGatewayApi_GetNodes verifies that GetNodes is reachable via the GatewayApi client
// and returns the same response as the PayerApi client.
func TestGatewayApi_GetNodes(t *testing.T) {
	suite := apiTestUtils.NewTestAPIServer(t)

	ctx := context.Background()

	respGateway, err := suite.ClientGateway.GetNodes(
		ctx,
		connect.NewRequest(&payer_api.GetNodesRequest{}),
	)
	require.NoError(t, err)
	require.NotNil(t, respGateway)

	respPayer, err := suite.ClientPayer.GetNodes(
		ctx,
		connect.NewRequest(&payer_api.GetNodesRequest{}),
	)
	require.NoError(t, err)
	require.NotNil(t, respPayer)

	require.Equal(t, respPayer.Msg.GetNodes(), respGateway.Msg.GetNodes())
}

// TestPublishApi_PublishPayerEnvelopes verifies that PublishPayerEnvelopes is reachable via the
// PublishApi client. An empty request returns a domain error (not an "unimplemented" error).
func TestPublishApi_PublishPayerEnvelopes(t *testing.T) {
	suite := apiTestUtils.NewTestAPIServer(t)

	_, err := suite.ClientPublish.PublishPayerEnvelopes(
		context.Background(),
		connect.NewRequest(&message_api.PublishPayerEnvelopesRequest{}),
	)
	// An empty request must fail with a domain error, NOT an "unimplemented" error.
	require.Error(t, err)
	require.NotEqual(t, connect.CodeUnimplemented, connect.CodeOf(err))
}
