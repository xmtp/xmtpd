package api_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api"
	"github.com/xmtp/xmtpd/pkg/testutils"
	"google.golang.org/protobuf/proto"
)

func TestPublishEnvelope(t *testing.T) {
	api, db, cleanup := testutils.NewTestAPIClient(t)
	defer cleanup()

	resp, err := api.PublishEnvelope(
		context.Background(),
		&message_api.PublishEnvelopeRequest{
			PayerEnvelope: testutils.CreatePayerEnvelope(t),
		},
	)
	require.NoError(t, err)
	require.NotNil(t, resp)

	unsignedEnv := &message_api.UnsignedOriginatorEnvelope{}
	require.NoError(
		t,
		proto.Unmarshal(resp.GetOriginatorEnvelope().GetUnsignedOriginatorEnvelope(), unsignedEnv),
	)
	clientEnv := &message_api.ClientEnvelope{}
	require.NoError(
		t,
		proto.Unmarshal(unsignedEnv.GetPayerEnvelope().GetUnsignedClientEnvelope(), clientEnv),
	)
	require.Equal(t, uint8(0x5), clientEnv.Aad.GetTargetTopic()[0])

	// Check that the envelope was published to the database after a delay
	require.Eventually(t, func() bool {
		envs, err := queries.New(db).
			SelectGatewayEnvelopes(context.Background(), queries.SelectGatewayEnvelopesParams{})
		require.NoError(t, err)

		if len(envs) != 1 {
			return false
		}

		originatorEnv := &message_api.OriginatorEnvelope{}
		require.NoError(t, proto.Unmarshal(envs[0].OriginatorEnvelope, originatorEnv))
		return proto.Equal(originatorEnv, resp.GetOriginatorEnvelope())
	}, 500*time.Millisecond, 50*time.Millisecond)
}

func TestUnmarshalErrorOnPublish(t *testing.T) {
	api, _, cleanup := testutils.NewTestAPIClient(t)
	defer cleanup()

	envelope := testutils.CreatePayerEnvelope(t)
	envelope.UnsignedClientEnvelope = []byte("invalidbytes")
	_, err := api.PublishEnvelope(
		context.Background(),
		&message_api.PublishEnvelopeRequest{
			PayerEnvelope: envelope,
		},
	)
	require.ErrorContains(t, err, "unmarshal")
}

func TestMismatchingOriginatorOnPublish(t *testing.T) {
	api, _, cleanup := testutils.NewTestAPIClient(t)
	defer cleanup()

	clientEnv := testutils.CreateClientEnvelope()
	clientEnv.Aad.TargetOriginator = 2
	_, err := api.PublishEnvelope(
		context.Background(),
		&message_api.PublishEnvelopeRequest{
			PayerEnvelope: testutils.CreatePayerEnvelope(t, clientEnv),
		},
	)
	require.ErrorContains(t, err, "originator")
}

func TestMissingTopicOnPublish(t *testing.T) {
	api, _, cleanup := testutils.NewTestAPIClient(t)
	defer cleanup()

	clientEnv := testutils.CreateClientEnvelope()
	clientEnv.Aad.TargetTopic = nil
	_, err := api.PublishEnvelope(
		context.Background(),
		&message_api.PublishEnvelopeRequest{
			PayerEnvelope: testutils.CreatePayerEnvelope(t, clientEnv),
		},
	)
	require.ErrorContains(t, err, "topic")
}
