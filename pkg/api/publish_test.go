package api_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/envelopes"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api"
	apiTestUtils "github.com/xmtp/xmtpd/pkg/testutils/api"
	envelopeTestUtils "github.com/xmtp/xmtpd/pkg/testutils/envelopes"
	"github.com/xmtp/xmtpd/pkg/topic"
	"google.golang.org/protobuf/proto"
)

func TestPublishEnvelope(t *testing.T) {
	api, db, cleanup := apiTestUtils.NewTestAPIClient(t)
	defer cleanup()

	payerEnvelope := envelopeTestUtils.CreatePayerEnvelope(t)

	resp, err := api.PublishPayerEnvelopes(
		context.Background(),
		&message_api.PublishPayerEnvelopesRequest{
			PayerEnvelopes: []*envelopes.PayerEnvelope{payerEnvelope},
		},
	)
	require.NoError(t, err)
	require.NotNil(t, resp)

	unsignedEnv := &envelopes.UnsignedOriginatorEnvelope{}
	require.NoError(
		t,
		proto.Unmarshal(
			resp.GetOriginatorEnvelopes()[0].GetUnsignedOriginatorEnvelope(),
			unsignedEnv,
		),
	)
	clientEnv := &envelopes.ClientEnvelope{}
	require.NoError(
		t,
		proto.Unmarshal(unsignedEnv.GetPayerEnvelope().GetUnsignedClientEnvelope(), clientEnv),
	)

	_, err = topic.ParseTopic(clientEnv.Aad.GetTargetTopic())
	require.NoError(t, err)

	// Check that the envelope was published to the database after a delay
	require.Eventually(t, func() bool {
		envs, err := queries.New(db).
			SelectGatewayEnvelopes(context.Background(), queries.SelectGatewayEnvelopesParams{})
		require.NoError(t, err)

		if len(envs) != 1 {
			return false
		}

		originatorEnv := &envelopes.OriginatorEnvelope{}
		require.NoError(t, proto.Unmarshal(envs[0].OriginatorEnvelope, originatorEnv))
		return proto.Equal(originatorEnv, resp.GetOriginatorEnvelopes()[0])
	}, 500*time.Millisecond, 50*time.Millisecond)
}

func TestUnmarshalErrorOnPublish(t *testing.T) {
	api, _, cleanup := apiTestUtils.NewTestAPIClient(t)
	defer cleanup()

	envelope := envelopeTestUtils.CreatePayerEnvelope(t)
	envelope.UnsignedClientEnvelope = []byte("invalidbytes")
	_, err := api.PublishPayerEnvelopes(
		context.Background(),
		&message_api.PublishPayerEnvelopesRequest{
			PayerEnvelopes: []*envelopes.PayerEnvelope{envelope},
		},
	)
	require.ErrorContains(t, err, "invalid wire-format data")
}

func TestMismatchingOriginatorOnPublish(t *testing.T) {
	api, _, cleanup := apiTestUtils.NewTestAPIClient(t)
	defer cleanup()

	clientEnv := envelopeTestUtils.CreateClientEnvelope()
	clientEnv.Aad.TargetOriginator = 2
	_, err := api.PublishPayerEnvelopes(
		context.Background(),
		&message_api.PublishPayerEnvelopesRequest{
			PayerEnvelopes: []*envelopes.PayerEnvelope{
				envelopeTestUtils.CreatePayerEnvelope(t, clientEnv),
			},
		},
	)
	require.ErrorContains(t, err, "originator")
}

func TestMissingTopicOnPublish(t *testing.T) {
	api, _, cleanup := apiTestUtils.NewTestAPIClient(t)
	defer cleanup()

	clientEnv := envelopeTestUtils.CreateClientEnvelope()
	clientEnv.Aad.TargetTopic = nil
	_, err := api.PublishPayerEnvelopes(
		context.Background(),
		&message_api.PublishPayerEnvelopesRequest{
			PayerEnvelopes: []*envelopes.PayerEnvelope{
				envelopeTestUtils.CreatePayerEnvelope(t, clientEnv),
			},
		},
	)
	require.ErrorContains(t, err, "topic")
}
