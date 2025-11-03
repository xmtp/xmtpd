package message_test

import (
	"context"
	"math"
	"testing"
	"time"

	"connectrpc.com/connect"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/currency"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	envelopeUtils "github.com/xmtp/xmtpd/pkg/envelopes"
	"github.com/xmtp/xmtpd/pkg/mlsvalidate"
	apiv1 "github.com/xmtp/xmtpd/pkg/proto/mls/api/v1"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/envelopes"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api"
	message_apiconnect "github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api/message_apiconnect"
	"github.com/xmtp/xmtpd/pkg/testutils"
	apiTestUtils "github.com/xmtp/xmtpd/pkg/testutils/api"
	envelopeTestUtils "github.com/xmtp/xmtpd/pkg/testutils/envelopes"
	"github.com/xmtp/xmtpd/pkg/topic"
	"google.golang.org/protobuf/proto"
)

func TestPublishEnvelope(t *testing.T) {
	var (
		client = apiTestUtils.NewTestGRPCReplicationAPIClient(t, "localhost:0")
		db, _  = testutils.NewDB(t, t.Context())
	)

	payerEnvelope := envelopeTestUtils.CreatePayerEnvelope(
		t,
		envelopeTestUtils.DefaultClientEnvelopeNodeID,
	)

	resp, err := client.PublishPayerEnvelopes(
		context.Background(),
		connect.NewRequest(&message_api.PublishPayerEnvelopesRequest{
			PayerEnvelopes: []*envelopes.PayerEnvelope{payerEnvelope},
		}),
	)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, resp.Msg)

	unsignedEnv := &envelopes.UnsignedOriginatorEnvelope{}
	require.NoError(
		t,
		proto.Unmarshal(
			resp.Msg.OriginatorEnvelopes[0].GetUnsignedOriginatorEnvelope(),
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
			SelectGatewayEnvelopesUnfiltered(context.Background(), queries.SelectGatewayEnvelopesUnfilteredParams{})
		require.NoError(t, err)

		if len(envs) != 1 {
			return false
		}

		originatorEnv := &envelopes.OriginatorEnvelope{}
		require.NoError(t, proto.Unmarshal(envs[0].OriginatorEnvelope, originatorEnv))
		return proto.Equal(originatorEnv, resp.Msg.OriginatorEnvelopes[0])
	}, 500*time.Millisecond, 50*time.Millisecond)
}

func TestUnmarshalErrorOnPublish(t *testing.T) {
	client := apiTestUtils.NewTestGRPCReplicationAPIClient(t, "localhost:0")

	envelope := envelopeTestUtils.CreatePayerEnvelope(
		t,
		envelopeTestUtils.DefaultClientEnvelopeNodeID,
	)

	envelope.UnsignedClientEnvelope = []byte("invalidbytes")

	_, err := client.PublishPayerEnvelopes(
		context.Background(),
		connect.NewRequest(&message_api.PublishPayerEnvelopesRequest{
			PayerEnvelopes: []*envelopes.PayerEnvelope{envelope},
		}),
	)
	require.ErrorContains(t, err, "invalid wire-format data")
}

func TestMismatchingOriginatorOnPublish(t *testing.T) {
	client := apiTestUtils.NewTestGRPCReplicationAPIClient(t, "localhost:0")

	nodeID := envelopeTestUtils.DefaultClientEnvelopeNodeID + 100

	clientEnv := envelopeTestUtils.CreateClientEnvelope()
	_, err := client.PublishPayerEnvelopes(
		context.Background(),
		connect.NewRequest(&message_api.PublishPayerEnvelopesRequest{
			PayerEnvelopes: []*envelopes.PayerEnvelope{
				envelopeTestUtils.CreatePayerEnvelope(t, nodeID, clientEnv),
			},
		}),
	)
	require.ErrorContains(t, err, "originator")
}

func TestMissingTopicOnPublish(t *testing.T) {
	var (
		client    = apiTestUtils.NewTestGRPCReplicationAPIClient(t, "localhost:0")
		clientEnv = envelopeTestUtils.CreateClientEnvelope()
	)

	clientEnv.Aad.TargetTopic = nil

	_, err := client.PublishPayerEnvelopes(
		context.Background(),
		connect.NewRequest(&message_api.PublishPayerEnvelopesRequest{
			PayerEnvelopes: []*envelopes.PayerEnvelope{
				envelopeTestUtils.CreatePayerEnvelope(
					t,
					envelopeTestUtils.DefaultClientEnvelopeNodeID,
					clientEnv,
				),
			},
		}),
	)
	require.ErrorContains(t, err, "topic")
}

func TestKeyPackageValidationSuccess(t *testing.T) {
	var (
		client         = apiTestUtils.NewTestGRPCReplicationAPIClient(t, "localhost:0")
		_, _, apiMocks = apiTestUtils.NewTestFullServer(t)
	)

	clientEnv := envelopeTestUtils.CreateClientEnvelope(
		&envelopeTestUtils.ClientEnvelopeOptions{Aad: &envelopes.AuthenticatedData{
			TargetTopic: topic.NewTopic(topic.TopicKindKeyPackagesV1, []byte{1, 2, 3}).Bytes(),
			DependsOn:   &envelopes.Cursor{},
		}},
	)

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

	resp, err := client.PublishPayerEnvelopes(
		context.Background(),
		connect.NewRequest(&message_api.PublishPayerEnvelopesRequest{
			PayerEnvelopes: []*envelopes.PayerEnvelope{
				envelopeTestUtils.CreatePayerEnvelope(
					t,
					envelopeTestUtils.DefaultClientEnvelopeNodeID,
					clientEnv,
				),
			},
		}),
	)
	require.Nil(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, resp.Msg)
}

func TestKeyPackageValidationFail(t *testing.T) {
	var (
		client         = apiTestUtils.NewTestGRPCReplicationAPIClient(t, "localhost:0")
		_, _, apiMocks = apiTestUtils.NewTestFullServer(t)
		nodeID         = envelopeTestUtils.DefaultClientEnvelopeNodeID
	)

	clientEnv := envelopeTestUtils.CreateClientEnvelope(
		&envelopeTestUtils.ClientEnvelopeOptions{Aad: &envelopes.AuthenticatedData{
			TargetTopic: topic.NewTopic(topic.TopicKindKeyPackagesV1, []byte{1, 2, 3}).Bytes(),
			DependsOn:   &envelopes.Cursor{},
		}},
	)

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

	_, err := client.PublishPayerEnvelopes(
		context.Background(),
		connect.NewRequest(&message_api.PublishPayerEnvelopesRequest{
			PayerEnvelopes: []*envelopes.PayerEnvelope{
				envelopeTestUtils.CreatePayerEnvelope(t, nodeID, clientEnv),
			},
		}),
	)
	require.Error(t, err)
}

func TestPublishEnvelopeBlockchainCursorAhead(t *testing.T) {
	client := apiTestUtils.NewTestGRPCReplicationAPIClient(t, "localhost:0")

	err := publishPayerEnvelopeWithNodeIDAndCursor(
		t,
		client,
		envelopeTestUtils.DefaultClientEnvelopeNodeID,
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
	client message_apiconnect.ReplicationApiClient,
	nodeID uint32,
	cursor *envelopes.Cursor,
) error {
	targetTopic := topic.NewTopic(topic.TopicKindGroupMessagesV1, []byte{1, 2, 3}).
		Bytes()

	_, err := client.PublishPayerEnvelopes(
		context.Background(),
		connect.NewRequest(&message_api.PublishPayerEnvelopesRequest{
			PayerEnvelopes: []*envelopes.PayerEnvelope{envelopeTestUtils.CreatePayerEnvelope(
				t,
				nodeID,
				envelopeTestUtils.CreateClientEnvelope(
					&envelopeTestUtils.ClientEnvelopeOptions{Aad: &envelopes.AuthenticatedData{
						TargetTopic: targetTopic,
						DependsOn:   cursor,
					}},
				),
			)},
		}),
	)

	return err
}

func TestPublishEnvelopeOriginatorUnknown(t *testing.T) {
	client := apiTestUtils.NewTestGRPCReplicationAPIClient(t, "localhost:0")

	err := publishPayerEnvelopeWithNodeIDAndCursor(
		t,
		client,
		envelopeTestUtils.DefaultClientEnvelopeNodeID,
		&envelopes.Cursor{
			NodeIdToSequenceId: map[uint32]uint64{
				97: 1,
			},
		},
	)
	require.Error(t, err)
	require.Contains(t, err.Error(), "DependsOn has not been seen by this node")
}

func TestPublishEnvelolopeDependsOnOriginator(t *testing.T) {
	client := apiTestUtils.NewTestGRPCReplicationAPIClient(t, "localhost:0")

	err := publishPayerEnvelopeWithNodeIDAndCursor(
		t,
		client,
		envelopeTestUtils.DefaultClientEnvelopeNodeID,
		&envelopes.Cursor{
			NodeIdToSequenceId: map[uint32]uint64{
				100: 1,
			},
		},
	)
	require.Error(t, err)
	require.Contains(t, err.Error(), "A message can not depend on a non-commit")
}

func TestPublishEnvelopeFees(t *testing.T) {
	client := apiTestUtils.NewTestGRPCReplicationAPIClient(t, "localhost:0")
	db, _ := testutils.NewDB(t, t.Context())

	payerEnvelope := envelopeTestUtils.CreatePayerEnvelope(
		t,
		envelopeTestUtils.DefaultClientEnvelopeNodeID,
	)

	resp, err := client.PublishPayerEnvelopes(
		context.Background(),
		connect.NewRequest(&message_api.PublishPayerEnvelopesRequest{
			PayerEnvelopes: []*envelopes.PayerEnvelope{payerEnvelope},
		}),
	)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, resp.Msg)

	returnedEnv, err := envelopeUtils.NewOriginatorEnvelope(resp.Msg.OriginatorEnvelopes[0])
	require.NoError(t, err)

	// BaseFee will always be > 0
	require.Greater(t, returnedEnv.UnsignedOriginatorEnvelope.BaseFee(), currency.PicoDollar(0))

	// CongestionFee will be 0 for now.
	// TODO:nm: Set this to the actual congestion fee
	require.Equal(t, returnedEnv.UnsignedOriginatorEnvelope.CongestionFee(), currency.PicoDollar(0))

	envs, err := queries.New(db).
		SelectGatewayEnvelopesUnfiltered(context.Background(), queries.SelectGatewayEnvelopesUnfilteredParams{})
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

func TestPublishEnvelopeFeesReservedTopic(t *testing.T) {
	var (
		client  = apiTestUtils.NewTestGRPCReplicationAPIClient(t, "localhost:0")
		db, _   = testutils.NewDB(t, t.Context())
		querier = queries.New(db)
	)

	clientEnv := envelopeTestUtils.CreatePayerReportClientEnvelope(100)

	// Create a payer envelope with a reserved topic (PAYER_REPORTS_V1)
	payerEnvelope := envelopeTestUtils.CreatePayerEnvelope(
		t,
		envelopeTestUtils.DefaultClientEnvelopeNodeID,
		clientEnv,
	)

	// Attempt to publish the envelope through the API
	_, err := client.PublishPayerEnvelopes(
		context.Background(),
		connect.NewRequest(&message_api.PublishPayerEnvelopesRequest{
			PayerEnvelopes: []*envelopes.PayerEnvelope{payerEnvelope},
		}),
	)
	require.Error(t, err)
	require.Contains(t, err.Error(), "reserved topics")

	payerEnvelopeBytes, err := proto.Marshal(payerEnvelope)
	require.NoError(t, err)

	// Write to the DB directly to simulate publishing to a reserved topic
	// since the API will eventually block publishing to reserved topics
	_, err = querier.InsertStagedOriginatorEnvelope(
		context.Background(),
		queries.InsertStagedOriginatorEnvelopeParams{
			Topic:         clientEnv.Aad.TargetTopic,
			PayerEnvelope: payerEnvelopeBytes,
		},
	)
	require.NoError(t, err)

	require.Eventually(t, func() bool {
		envs, err := querier.
			SelectGatewayEnvelopesUnfiltered(
				context.Background(),
				queries.SelectGatewayEnvelopesUnfilteredParams{},
			)
		require.NoError(t, err)
		if len(envs) != 1 {
			return false
		}
		originatorEnv, err := envelopeUtils.NewOriginatorEnvelopeFromBytes(
			envs[0].OriginatorEnvelope,
		)
		require.NoError(t, err)
		return originatorEnv.UnsignedOriginatorEnvelope.BaseFee() == currency.PicoDollar(0) &&
			originatorEnv.UnsignedOriginatorEnvelope.CongestionFee() == currency.PicoDollar(0)
	}, 2*time.Second, 500*time.Millisecond)
}

func TestPublishEnvelopeWithVarExpirations(t *testing.T) {
	client := apiTestUtils.NewTestGRPCReplicationAPIClient(t, "localhost:0")

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
				envelopeTestUtils.DefaultClientEnvelopeNodeID,
				tt.expiry,
			)

			_, err := client.PublishPayerEnvelopes(
				context.Background(),
				connect.NewRequest(&message_api.PublishPayerEnvelopesRequest{
					PayerEnvelopes: []*envelopes.PayerEnvelope{payerEnvelope},
				}),
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

func TestPublishCommitViaNodeGetsRejected(t *testing.T) {
	var (
		client = apiTestUtils.NewTestGRPCReplicationAPIClient(t, "localhost:0")
		nodeID = envelopeTestUtils.DefaultClientEnvelopeNodeID
	)

	clientEnv := envelopeTestUtils.CreateClientEnvelope(
		&envelopeTestUtils.ClientEnvelopeOptions{Aad: &envelopes.AuthenticatedData{
			TargetTopic: topic.NewTopic(topic.TopicKindGroupMessagesV1, []byte{1, 2, 3}).
				Bytes(),
			DependsOn: &envelopes.Cursor{},
		}, IsCommit: true},
	)
	_, err := client.PublishPayerEnvelopes(
		context.Background(),
		connect.NewRequest(&message_api.PublishPayerEnvelopesRequest{
			PayerEnvelopes: []*envelopes.PayerEnvelope{
				envelopeTestUtils.CreatePayerEnvelope(t, nodeID, clientEnv),
			},
		}),
	)
	require.ErrorContains(t, err, "published via the blockchain")
}
