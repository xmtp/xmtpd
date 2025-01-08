package message_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/mlsvalidate"
	apiv1 "github.com/xmtp/xmtpd/pkg/proto/mls/api/v1"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/envelopes"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api"
	apiTestUtils "github.com/xmtp/xmtpd/pkg/testutils/api"
	envelopeTestUtils "github.com/xmtp/xmtpd/pkg/testutils/envelopes"
	"github.com/xmtp/xmtpd/pkg/topic"
	"google.golang.org/protobuf/proto"
)

func TestPublishEnvelope(t *testing.T) {
	api, db, _, cleanup := apiTestUtils.NewTestReplicationAPIClient(t)
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
	api, _, _, cleanup := apiTestUtils.NewTestReplicationAPIClient(t)
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
	api, _, _, cleanup := apiTestUtils.NewTestReplicationAPIClient(t)
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
	api, _, _, cleanup := apiTestUtils.NewTestReplicationAPIClient(t)
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

func TestKeyPackageValidationSuccess(t *testing.T) {
	api, _, apiMocks, cleanup := apiTestUtils.NewTestReplicationAPIClient(t)
	defer cleanup()

	clientEnv := envelopeTestUtils.CreateClientEnvelope(&envelopes.AuthenticatedData{
		TargetTopic:      topic.NewTopic(topic.TOPIC_KIND_KEY_PACKAGES_V1, []byte{1, 2, 3}).Bytes(),
		TargetOriginator: 100,
		LastSeen:         &envelopes.Cursor{},
	})
	clientEnv.Payload = &envelopes.ClientEnvelope_UploadKeyPackage{
		UploadKeyPackage: &apiv1.UploadKeyPackageRequest{
			KeyPackage: &apiv1.KeyPackageUpload{
				KeyPackageTlsSerialized: []byte{1, 2, 3},
			},
		},
	}

	apiMocks.MockValidationService.EXPECT().
		ValidateKeyPackages(mock.Anything, mock.Anything).
		Return(
			[]mlsvalidate.KeyPackageValidationResult{
				{
					IsOk: true,
				},
			},
			nil,
		)

	_, err := api.PublishPayerEnvelopes(
		context.Background(),
		&message_api.PublishPayerEnvelopesRequest{
			PayerEnvelopes: []*envelopes.PayerEnvelope{
				envelopeTestUtils.CreatePayerEnvelope(t, clientEnv),
			},
		},
	)
	require.Nil(t, err)
}

func TestKeyPackageValidationFail(t *testing.T) {
	api, _, apiMocks, cleanup := apiTestUtils.NewTestReplicationAPIClient(t)
	defer cleanup()

	clientEnv := envelopeTestUtils.CreateClientEnvelope(&envelopes.AuthenticatedData{
		TargetTopic:      topic.NewTopic(topic.TOPIC_KIND_KEY_PACKAGES_V1, []byte{1, 2, 3}).Bytes(),
		TargetOriginator: 100,
		LastSeen:         &envelopes.Cursor{},
	})
	clientEnv.Payload = &envelopes.ClientEnvelope_UploadKeyPackage{
		UploadKeyPackage: &apiv1.UploadKeyPackageRequest{
			KeyPackage: &apiv1.KeyPackageUpload{
				KeyPackageTlsSerialized: []byte{1, 2, 3},
			},
		},
	}

	apiMocks.MockValidationService.EXPECT().
		ValidateKeyPackages(mock.Anything, mock.Anything).
		Return(
			[]mlsvalidate.KeyPackageValidationResult{
				{
					IsOk: false,
				},
			},
			nil,
		)

	_, err := api.PublishPayerEnvelopes(
		context.Background(),
		&message_api.PublishPayerEnvelopesRequest{
			PayerEnvelopes: []*envelopes.PayerEnvelope{
				envelopeTestUtils.CreatePayerEnvelope(t, clientEnv),
			},
		},
	)
	require.Error(t, err)
}
