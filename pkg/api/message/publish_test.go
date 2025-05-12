package message_test

import (
	"context"
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/currency"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	envelopeUtils "github.com/xmtp/xmtpd/pkg/envelopes"
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

	payerEnvelope := envelopeTestUtils.CreatePayerEnvelope(
		t,
		envelopeTestUtils.DefaultClientEnvelopeNodeId,
	)

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
	payerEnv := &envelopes.PayerEnvelope{}
	require.NoError(
		t,
		proto.Unmarshal(unsignedEnv.GetPayerEnvelopeBytes(), payerEnv),
	)

	clientEnv := &envelopes.ClientEnvelope{}
	require.NoError(
		t,
		proto.Unmarshal(payerEnv.GetUnsignedClientEnvelope(), clientEnv),
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

	envelope := envelopeTestUtils.CreatePayerEnvelope(
		t,
		envelopeTestUtils.DefaultClientEnvelopeNodeId,
	)
	envelope.UnsignedClientEnvelope = []byte("invalidbytes")
	_, err := api.PublishPayerEnvelopes(
		context.Background(),
		&message_api.PublishPayerEnvelopesRequest{
			PayerEnvelopes: []*envelopes.PayerEnvelope{envelope},
		},
	)
	require.ErrorContains(t, err, "invalid wire-format data")
}

func TestMismatchingAADOriginatorOnPublishNoLongerFails(t *testing.T) {
	api, _, _, cleanup := apiTestUtils.NewTestReplicationAPIClient(t)
	defer cleanup()

	nid := envelopeTestUtils.DefaultClientEnvelopeNodeId + 100

	clientEnv := envelopeTestUtils.CreateClientEnvelope()
	// nolint:staticcheck
	clientEnv.Aad.TargetOriginator = &nid
	_, err := api.PublishPayerEnvelopes(
		context.Background(),
		&message_api.PublishPayerEnvelopesRequest{
			PayerEnvelopes: []*envelopes.PayerEnvelope{
				envelopeTestUtils.CreatePayerEnvelope(
					t,
					envelopeTestUtils.DefaultClientEnvelopeNodeId,
					clientEnv,
				),
			},
		},
	)
	require.NoError(t, err)
}

func TestMismatchingOriginatorOnPublish(t *testing.T) {
	api, _, _, cleanup := apiTestUtils.NewTestReplicationAPIClient(t)
	defer cleanup()

	nid := envelopeTestUtils.DefaultClientEnvelopeNodeId + 100

	clientEnv := envelopeTestUtils.CreateClientEnvelope()
	_, err := api.PublishPayerEnvelopes(
		context.Background(),
		&message_api.PublishPayerEnvelopesRequest{
			PayerEnvelopes: []*envelopes.PayerEnvelope{
				envelopeTestUtils.CreatePayerEnvelope(t, nid, clientEnv),
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
				envelopeTestUtils.CreatePayerEnvelope(
					t,
					envelopeTestUtils.DefaultClientEnvelopeNodeId,
					clientEnv,
				),
			},
		},
	)
	require.ErrorContains(t, err, "topic")
}

func TestKeyPackageValidationSuccess(t *testing.T) {
	api, _, apiMocks, cleanup := apiTestUtils.NewTestReplicationAPIClient(t)
	defer cleanup()

	nid := envelopeTestUtils.DefaultClientEnvelopeNodeId

	clientEnv := envelopeTestUtils.CreateClientEnvelope(&envelopes.AuthenticatedData{
		TargetTopic:      topic.NewTopic(topic.TOPIC_KIND_KEY_PACKAGES_V1, []byte{1, 2, 3}).Bytes(),
		TargetOriginator: &nid,
		DependsOn:        &envelopes.Cursor{},
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
				envelopeTestUtils.CreatePayerEnvelope(
					t,
					envelopeTestUtils.DefaultClientEnvelopeNodeId,
					clientEnv,
				),
			},
		},
	)
	require.Nil(t, err)
}

func TestKeyPackageValidationFail(t *testing.T) {
	api, _, apiMocks, cleanup := apiTestUtils.NewTestReplicationAPIClient(t)
	defer cleanup()

	nid := envelopeTestUtils.DefaultClientEnvelopeNodeId

	clientEnv := envelopeTestUtils.CreateClientEnvelope(&envelopes.AuthenticatedData{
		TargetTopic:      topic.NewTopic(topic.TOPIC_KIND_KEY_PACKAGES_V1, []byte{1, 2, 3}).Bytes(),
		TargetOriginator: &nid,
		DependsOn:        &envelopes.Cursor{},
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
				envelopeTestUtils.CreatePayerEnvelope(t, nid, clientEnv),
			},
		},
	)
	require.Error(t, err)
}

func TestPublishEnvelopeBlockchainCursorAhead(t *testing.T) {
	api, _, _, cleanup := apiTestUtils.NewTestReplicationAPIClient(t)
	defer cleanup()

	err := publishPayerEnvelopeWithNodeIDAndCursor(
		t,
		api,
		envelopeTestUtils.DefaultClientEnvelopeNodeId,
		&envelopes.Cursor{
			NodeIdToSequenceId: map[uint32]uint64{
				1: 105,
			},
		},
	)
	require.Error(t, err)
	require.Contains(t, err.Error(), "DependsOn has not been seen by this node")
}

func publishPayerEnvelopeWithNodeIDAndCursor(
	t *testing.T,
	api message_api.ReplicationApiClient, nodeId uint32, cursor *envelopes.Cursor,
) error {
	targetTopic := topic.NewTopic(topic.TOPIC_KIND_GROUP_MESSAGES_V1, []byte{1, 2, 3}).
		Bytes()
	_, err := api.PublishPayerEnvelopes(
		context.Background(),
		&message_api.PublishPayerEnvelopesRequest{
			PayerEnvelopes: []*envelopes.PayerEnvelope{envelopeTestUtils.CreatePayerEnvelope(
				t, nodeId,
				envelopeTestUtils.CreateClientEnvelope(&envelopes.AuthenticatedData{
					TargetOriginator: &nodeId,
					TargetTopic:      targetTopic,
					DependsOn:        cursor,
				}),
			)},
		},
	)

	return err
}

func TestPublishEnvelopeOriginatorUnknown(t *testing.T) {
	api, _, _, cleanup := apiTestUtils.NewTestReplicationAPIClient(t)
	defer cleanup()

	err := publishPayerEnvelopeWithNodeIDAndCursor(
		t,
		api,
		envelopeTestUtils.DefaultClientEnvelopeNodeId,
		&envelopes.Cursor{
			NodeIdToSequenceId: map[uint32]uint64{
				1600: 1,
			},
		},
	)
	require.Error(t, err)
	require.Contains(t, err.Error(), "DependsOn has not been seen by this node")
}

func TestPublishEnvelopeFees(t *testing.T) {
	api, db, _, cleanup := apiTestUtils.NewTestReplicationAPIClient(t)
	defer cleanup()

	payerEnvelope := envelopeTestUtils.CreatePayerEnvelope(
		t,
		envelopeTestUtils.DefaultClientEnvelopeNodeId,
	)

	resp, err := api.PublishPayerEnvelopes(
		context.Background(),
		&message_api.PublishPayerEnvelopesRequest{
			PayerEnvelopes: []*envelopes.PayerEnvelope{payerEnvelope},
		},
	)
	require.NoError(t, err)
	require.NotNil(t, resp)

	returnedEnv, err := envelopeUtils.NewOriginatorEnvelope(resp.GetOriginatorEnvelopes()[0])
	require.NoError(t, err)
	// BaseFee will always be > 0
	require.Greater(t, returnedEnv.UnsignedOriginatorEnvelope.BaseFee(), currency.PicoDollar(0))
	// CongestionFee will be 0 for now.
	// TODO:nm: Set this to the actual congestion fee
	require.Equal(t, returnedEnv.UnsignedOriginatorEnvelope.CongestionFee(), currency.PicoDollar(0))

	envs, err := queries.New(db).
		SelectGatewayEnvelopes(context.Background(), queries.SelectGatewayEnvelopesParams{})
	require.NoError(t, err)
	require.Equal(t, len(envs), 1)

	originatorEnv, err := envelopeUtils.NewOriginatorEnvelopeFromBytes(envs[0].OriginatorEnvelope)
	require.NoError(t, err)
	require.Equal(
		t,
		originatorEnv.UnsignedOriginatorEnvelope.BaseFee(),
		returnedEnv.UnsignedOriginatorEnvelope.BaseFee(),
	)
	require.Equal(
		t,
		originatorEnv.UnsignedOriginatorEnvelope.CongestionFee(),
		returnedEnv.UnsignedOriginatorEnvelope.CongestionFee(),
	)
}

func TestPublishEnvelopeWithVarExpirations(t *testing.T) {
	api, _, _, cleanup := apiTestUtils.NewTestReplicationAPIClient(t)
	defer cleanup()

	tests := []struct {
		name        string
		expiry      uint32
		wantErr     bool
		expectedErr string
	}{
		{
			name:        "0 expiry",
			expiry:      0,
			wantErr:     true,
			expectedErr: "invalid expiry retention days",
		},
		{
			name:        "short expiry",
			expiry:      1,
			wantErr:     true,
			expectedErr: "invalid expiry retention days",
		},
		{
			name:    "minimal expiry",
			expiry:  2,
			wantErr: false,
		},
		{
			name:    "1 week expiry",
			expiry:  7,
			wantErr: false,
		},
		{
			name:    "30 day expiry",
			expiry:  30,
			wantErr: false,
		},
		{
			name:    "90 day expiry",
			expiry:  90,
			wantErr: false,
		},
		{
			name:        "5 year expiry",
			expiry:      5 * 365,
			wantErr:     true,
			expectedErr: "invalid expiry retention days",
		},
		{
			name:    "infinite expiry",
			expiry:  math.MaxUint32,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payerEnvelope := envelopeTestUtils.CreatePayerEnvelopeWithExpiration(
				t,
				envelopeTestUtils.DefaultClientEnvelopeNodeId,
				tt.expiry,
			)

			_, err := api.PublishPayerEnvelopes(
				context.Background(),
				&message_api.PublishPayerEnvelopesRequest{
					PayerEnvelopes: []*envelopes.PayerEnvelope{payerEnvelope},
				},
			)
			if tt.wantErr {
				if tt.expectedErr == "" {
					require.NoError(t, err)
				} else {
					require.ErrorContains(t, err, tt.expectedErr)
				}
			} else {
				require.NoError(t, err)
			}
		})
	}
}
