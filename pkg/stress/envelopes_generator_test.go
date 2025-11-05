package stress

import (
	"context"
	"fmt"
	"testing"

	"connectrpc.com/connect"
	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/envelopes"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api"
	"github.com/xmtp/xmtpd/pkg/testutils"
	apiTestUtils "github.com/xmtp/xmtpd/pkg/testutils/api"
)

func TestEnvelopesGenerator(t *testing.T) {
	var (
		suite = apiTestUtils.NewTestAPIServer(t)
		ctx   = t.Context()
	)

	// Create the envelopes generator.
	generator, err := NewEnvelopesGenerator(
		fmt.Sprintf("http://%s", suite.APIServer.Addr()),
		testutils.TestPrivateKey,
		100,
		ProtocolConnectGRPC,
	)
	require.NoError(t, err)

	// Publish the welcome message envelopes.
	publishResponse, err := generator.PublishWelcomeMessageEnvelopes(context.Background(), 1, 100)
	require.NoError(t, err)
	require.NotNil(t, publishResponse)

	// Query the envelopes.
	client := apiTestUtils.NewTestGRPCReplicationAPIClient(t, suite.APIServer.Addr())
	queryResponse, err := client.QueryEnvelopes(
		ctx,
		connect.NewRequest(&message_api.QueryEnvelopesRequest{
			Query: &message_api.EnvelopesQuery{
				OriginatorNodeIds: []uint32{100},
				LastSeen:          &envelopes.Cursor{},
			},
			Limit: 10,
		}),
	)
	require.NoError(t, err)
	require.NotNil(t, queryResponse)
	require.Len(t, queryResponse.Msg.Envelopes, 1)
	require.Equal(t, queryResponse.Msg.Envelopes[0], publishResponse[0])
}
