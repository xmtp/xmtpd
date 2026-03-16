package message_test

import (
	"context"
	"math"
	"testing"
	"time"

	"connectrpc.com/connect"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/currency"
	dbPkg "github.com/xmtp/xmtpd/pkg/db"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	envelopeUtils "github.com/xmtp/xmtpd/pkg/envelopes"
	"github.com/xmtp/xmtpd/pkg/ledger"
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
	suite := apiTestUtils.NewTestAPIServer(t)

	payerEnvelope := envelopeTestUtils.CreatePayerEnvelope(
		t,
		envelopeTestUtils.DefaultClientEnvelopeNodeID,
	)

	resp, err := suite.ClientReplication.PublishPayerEnvelopes(
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
			resp.Msg.GetOriginatorEnvelopes()[0].GetUnsignedOriginatorEnvelope(),
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

	_, err = topic.ParseTopic(clientEnv.GetAad().GetTargetTopic())
	require.NoError(t, err)

	// Check that the envelope was published to the database after a delay
	require.Eventually(t, func() bool {
		envs, err := queries.New(suite.DB).
			SelectGatewayEnvelopesUnfiltered(context.Background(), queries.SelectGatewayEnvelopesUnfilteredParams{})
		require.NoError(t, err)

		if len(envs) != 1 {
			return false
		}

		originatorEnv := &envelopes.OriginatorEnvelope{}
		require.NoError(t, proto.Unmarshal(envs[0].OriginatorEnvelope, originatorEnv))
		return proto.Equal(originatorEnv, resp.Msg.GetOriginatorEnvelopes()[0])
	}, 5*time.Second, 50*time.Millisecond)
}

func TestUnmarshalErrorOnPublish(t *testing.T) {
	suite := apiTestUtils.NewTestAPIServer(t)

	envelope := envelopeTestUtils.CreatePayerEnvelope(
		t,
		envelopeTestUtils.DefaultClientEnvelopeNodeID,
	)

	envelope.UnsignedClientEnvelope = []byte("invalidbytes")

	_, err := suite.ClientReplication.PublishPayerEnvelopes(
		context.Background(),
		connect.NewRequest(&message_api.PublishPayerEnvelopesRequest{
			PayerEnvelopes: []*envelopes.PayerEnvelope{envelope},
		}),
	)
	require.ErrorContains(t, err, "cannot parse invalid wire-format data")
}

func TestMismatchingOriginatorOnPublish(t *testing.T) {
	suite := apiTestUtils.NewTestAPIServer(t)

	nodeID := envelopeTestUtils.DefaultClientEnvelopeNodeID + 100

	clientEnv := envelopeTestUtils.CreateClientEnvelope()
	_, err := suite.ClientReplication.PublishPayerEnvelopes(
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
		suite     = apiTestUtils.NewTestAPIServer(t)
		clientEnv = envelopeTestUtils.CreateClientEnvelope()
	)

	clientEnv.Aad.TargetTopic = nil

	_, err := suite.ClientReplication.PublishPayerEnvelopes(
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
	suite := apiTestUtils.NewTestAPIServer(t)

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

	suite.APIServerMocks.MockValidationService.EXPECT().
		ValidateKeyPackages(mock.Anything, mock.Anything).
		Return(
			[]mlsvalidate.KeyPackageValidationResult{
				{
					IsOk: true,
				},
			},
			nil,
		)

	resp, err := suite.ClientReplication.PublishPayerEnvelopes(
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
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, resp.Msg)
}

func TestKeyPackageValidationFail(t *testing.T) {
	suite := apiTestUtils.NewTestAPIServer(t)
	nodeID := envelopeTestUtils.DefaultClientEnvelopeNodeID

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

	suite.APIServerMocks.MockValidationService.EXPECT().
		ValidateKeyPackages(mock.Anything, mock.Anything).
		Return(
			[]mlsvalidate.KeyPackageValidationResult{
				{
					IsOk: false,
				},
			},
			nil,
		)

	_, err := suite.ClientReplication.PublishPayerEnvelopes(
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
	suite := apiTestUtils.NewTestAPIServer(t)

	err := publishPayerEnvelopeWithNodeIDAndCursor(
		t,
		suite.ClientReplication,
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
	suite := apiTestUtils.NewTestAPIServer(t)

	err := publishPayerEnvelopeWithNodeIDAndCursor(
		t,
		suite.ClientReplication,
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
	suite := apiTestUtils.NewTestAPIServer(t)

	err := publishPayerEnvelopeWithNodeIDAndCursor(
		t,
		suite.ClientReplication,
		envelopeTestUtils.DefaultClientEnvelopeNodeID,
		&envelopes.Cursor{
			NodeIdToSequenceId: map[uint32]uint64{
				100: 1,
			},
		},
	)
	require.Error(t, err)
	require.Contains(t, err.Error(), "a message can not depend on a non-commit")
}

func TestPublishEnvelopeFees(t *testing.T) {
	suite := apiTestUtils.NewTestAPIServer(t)

	payerEnvelope := envelopeTestUtils.CreatePayerEnvelope(
		t,
		envelopeTestUtils.DefaultClientEnvelopeNodeID,
	)

	resp, err := suite.ClientReplication.PublishPayerEnvelopes(
		context.Background(),
		connect.NewRequest(&message_api.PublishPayerEnvelopesRequest{
			PayerEnvelopes: []*envelopes.PayerEnvelope{payerEnvelope},
		}),
	)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, resp.Msg)

	returnedEnv, err := envelopeUtils.NewOriginatorEnvelope(resp.Msg.GetOriginatorEnvelopes()[0])
	require.NoError(t, err)

	// BaseFee will always be > 0
	require.Greater(t, returnedEnv.UnsignedOriginatorEnvelope.BaseFee(), currency.PicoDollar(0))

	// CongestionFee will be 0 for now.
	// TODO:nm: Set this to the actual congestion fee
	require.Equal(t, returnedEnv.UnsignedOriginatorEnvelope.CongestionFee(), currency.PicoDollar(0))

	envs, err := queries.New(suite.DB).
		SelectGatewayEnvelopesUnfiltered(context.Background(), queries.SelectGatewayEnvelopesUnfilteredParams{})
	require.NoError(t, err)
	require.Len(t, envs, 1)

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
		suite   = apiTestUtils.NewTestAPIServer(t)
		querier = queries.New(suite.DB)
	)

	clientEnv := envelopeTestUtils.CreatePayerReportClientEnvelope(100)

	// Create a payer envelope with a reserved topic (PAYER_REPORTS_V1)
	payerEnvelope := envelopeTestUtils.CreatePayerEnvelope(
		t,
		envelopeTestUtils.DefaultClientEnvelopeNodeID,
		clientEnv,
	)

	// Attempt to publish the envelope through the API
	_, err := suite.ClientReplication.PublishPayerEnvelopes(
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
			Topic:         clientEnv.GetAad().GetTargetTopic(),
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
	}, 10*time.Second, 500*time.Millisecond)
}

func TestPublishEnvelopeWithVarExpirations(t *testing.T) {
	suite := apiTestUtils.NewTestAPIServer(t)

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

			_, err := suite.ClientReplication.PublishPayerEnvelopes(
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
		suite  = apiTestUtils.NewTestAPIServer(t)
		nodeID = envelopeTestUtils.DefaultClientEnvelopeNodeID
	)

	clientEnv := envelopeTestUtils.CreateClientEnvelope(
		&envelopeTestUtils.ClientEnvelopeOptions{Aad: &envelopes.AuthenticatedData{
			TargetTopic: topic.NewTopic(topic.TopicKindGroupMessagesV1, []byte{1, 2, 3}).
				Bytes(),
			DependsOn: &envelopes.Cursor{},
		}, IsCommit: true},
	)
	_, err := suite.ClientReplication.PublishPayerEnvelopes(
		context.Background(),
		connect.NewRequest(&message_api.PublishPayerEnvelopesRequest{
			PayerEnvelopes: []*envelopes.PayerEnvelope{
				envelopeTestUtils.CreatePayerEnvelope(t, nodeID, clientEnv),
			},
		}),
	)
	require.ErrorContains(t, err, "published via the blockchain")
}

func TestPublishEnvelopeEmpty(t *testing.T) {
	suite := apiTestUtils.NewTestAPIServer(t)

	resp, err := suite.ClientReplication.PublishPayerEnvelopes(
		context.Background(),
		connect.NewRequest(&message_api.PublishPayerEnvelopesRequest{}),
	)
	require.Error(t, err)
	require.Nil(t, resp)
	require.Contains(t, err.Error(), "missing payer envelope")
}

func TestPublishEnvelopeBatchPublish(t *testing.T) {
	suite := apiTestUtils.NewTestAPIServer(t)

	payerEnvelope := envelopeTestUtils.CreatePayerEnvelope(
		t,
		envelopeTestUtils.DefaultClientEnvelopeNodeID,
	)

	resp, err := suite.ClientReplication.PublishPayerEnvelopes(
		context.Background(),
		connect.NewRequest(&message_api.PublishPayerEnvelopesRequest{
			PayerEnvelopes: []*envelopes.PayerEnvelope{payerEnvelope, payerEnvelope, payerEnvelope},
		}),
	)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, resp.Msg)

	require.Len(t, resp.Msg.GetOriginatorEnvelopes(), 3)

	require.Eventually(t, func() bool {
		envs, err := queries.New(suite.DB).
			SelectGatewayEnvelopesUnfiltered(context.Background(), queries.SelectGatewayEnvelopesUnfilteredParams{})
		require.NoError(t, err)

		return len(envs) == 3
	}, 5*time.Second, 50*time.Millisecond)
}

func TestPublishEnvelopeBatchPublishNoPartialError(t *testing.T) {
	suite := apiTestUtils.NewTestAPIServer(t)

	payerEnvelope := envelopeTestUtils.CreatePayerEnvelope(
		t,
		envelopeTestUtils.DefaultClientEnvelopeNodeID,
	)

	clientEnv := envelopeTestUtils.CreateClientEnvelope(
		&envelopeTestUtils.ClientEnvelopeOptions{Aad: &envelopes.AuthenticatedData{
			TargetTopic: topic.NewTopic(topic.TopicKindGroupMessagesV1, []byte{1, 2, 3}).
				Bytes(),
			DependsOn: &envelopes.Cursor{},
		}, IsCommit: true},
	)
	invalidPayerEnvelope := envelopeTestUtils.CreatePayerEnvelope(
		t,
		envelopeTestUtils.DefaultClientEnvelopeNodeID,
		clientEnv,
	)

	resp, err := suite.ClientReplication.PublishPayerEnvelopes(
		context.Background(),
		connect.NewRequest(&message_api.PublishPayerEnvelopesRequest{
			PayerEnvelopes: []*envelopes.PayerEnvelope{
				payerEnvelope,
				payerEnvelope,
				invalidPayerEnvelope,
			},
		}),
	)
	require.Error(t, err)
	require.Nil(t, resp)
	require.Contains(t, err.Error(), "published via the blockchain")

	// give this some time to process just in case
	time.Sleep(100 * time.Millisecond)

	envs, err := queries.New(suite.DB).
		SelectGatewayEnvelopesUnfiltered(context.Background(), queries.SelectGatewayEnvelopesUnfilteredParams{})
	require.NoError(t, err)
	require.Empty(t, envs)
}

func TestPublishEnvelopeBalanceEnforcement(t *testing.T) {
	tests := []struct {
		name           string
		enforce        bool
		deposit        currency.PicoDollar
		unsettledUsage int64
		wantCode       connect.Code // 0 means expect success
	}{
		{
			name:    "enforcement off, no balance — succeeds",
			enforce: false,
		},
		{
			name:     "enforcement on, no balance — rejected",
			enforce:  true,
			wantCode: connect.CodeFailedPrecondition,
		},
		{
			name:    "enforcement on, sufficient balance — succeeds",
			enforce: true,
			deposit: 1_000_000_000_000, // 1 dollar
		},
		{
			name:           "enforcement on, balance consumed by unsettled usage — rejected",
			enforce:        true,
			deposit:        1_000_000_000_000, // 1 dollar
			unsettledUsage: 999_999_999_999,   // nearly 1 dollar
			wantCode:       connect.CodeFailedPrecondition,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var opts []apiTestUtils.TestAPIOption
			if tc.enforce {
				opts = append(opts, apiTestUtils.WithRequirePayerPositiveBalance(true))
			}
			suite := apiTestUtils.NewTestAPIServer(t, opts...)

			payerEnvelope := envelopeTestUtils.CreatePayerEnvelope(
				t,
				envelopeTestUtils.DefaultClientEnvelopeNodeID,
			)

			if tc.deposit > 0 || tc.unsettledUsage > 0 {
				payerEnv, err := envelopeUtils.NewPayerEnvelope(payerEnvelope)
				require.NoError(t, err)
				payerAddr, err := payerEnv.RecoverSigner()
				require.NoError(t, err)

				payerLedger := ledger.NewLedger(
					testutils.NewLog(t),
					dbPkg.NewDBHandler(suite.DB),
				)
				payerID, err := payerLedger.FindOrCreatePayer(
					context.Background(),
					*payerAddr,
				)
				require.NoError(t, err)

				if tc.deposit > 0 {
					eventID := ledger.EventID{}
					copy(eventID[:], []byte("test-deposit-event-id-00001"))
					err = payerLedger.Deposit(
						context.Background(),
						payerID,
						tc.deposit,
						eventID,
					)
					require.NoError(t, err)
				}

				if tc.unsettledUsage > 0 {
					err = queries.New(suite.DB).IncrementUnsettledUsage(
						context.Background(),
						queries.IncrementUnsettledUsageParams{
							PayerID:           payerID,
							OriginatorID:      int32(envelopeTestUtils.DefaultClientEnvelopeNodeID),
							MinutesSinceEpoch: 1,
							SpendPicodollars:  tc.unsettledUsage,
							SequenceID:        1,
							MessageCount:      1,
						},
					)
					require.NoError(t, err)
				}
			}

			resp, err := suite.ClientReplication.PublishPayerEnvelopes(
				context.Background(),
				connect.NewRequest(&message_api.PublishPayerEnvelopesRequest{
					PayerEnvelopes: []*envelopes.PayerEnvelope{payerEnvelope},
				}),
			)

			if tc.wantCode == 0 {
				require.NoError(t, err)
				require.NotNil(t, resp)
			} else {
				require.Error(t, err)
				require.Equal(t, tc.wantCode, connect.CodeOf(err))
			}
		})
	}
}

func TestPublishEnvelopeMultiEnvelopeBatchBalance(t *testing.T) {
	suite := apiTestUtils.NewTestAPIServer(
		t,
		apiTestUtils.WithRequirePayerPositiveBalance(true),
	)

	// Create a shared signer so all envelopes have the same payer
	signerKey := testutils.RandomPrivateKey(t)

	nodeID := envelopeTestUtils.DefaultClientEnvelopeNodeID

	// Create 3 envelopes from the same payer
	env1 := envelopeTestUtils.CreatePayerEnvelopeWithSigner(
		t, nodeID, signerKey, 30, envelopeTestUtils.CreateClientEnvelope(),
	)
	env2 := envelopeTestUtils.CreatePayerEnvelopeWithSigner(
		t, nodeID, signerKey, 30, envelopeTestUtils.CreateClientEnvelope(),
	)
	env3 := envelopeTestUtils.CreatePayerEnvelopeWithSigner(
		t, nodeID, signerKey, 30, envelopeTestUtils.CreateClientEnvelope(),
	)

	// Deposit enough for 1 envelope but not 3
	payerAddr := crypto.PubkeyToAddress(signerKey.PublicKey)
	payerLedger := ledger.NewLedger(testutils.NewLog(t), dbPkg.NewDBHandler(suite.DB))
	payerID, err := payerLedger.FindOrCreatePayer(context.Background(), payerAddr)
	require.NoError(t, err)

	// Deposit a very small amount — enough for maybe 1 message but not 3
	eventID := ledger.EventID{}
	copy(eventID[:], []byte("test-batch-event-id-0000001"))
	err = payerLedger.Deposit(
		context.Background(),
		payerID,
		currency.PicoDollar(1), // 1 picodollar — nearly nothing
		eventID,
	)
	require.NoError(t, err)

	// Should fail — batch fee total exceeds tiny balance
	_, err = suite.ClientReplication.PublishPayerEnvelopes(
		context.Background(),
		connect.NewRequest(&message_api.PublishPayerEnvelopesRequest{
			PayerEnvelopes: []*envelopes.PayerEnvelope{env1, env2, env3},
		}),
	)
	require.Error(t, err)
	require.Equal(t, connect.CodeFailedPrecondition, connect.CodeOf(err))
}

func TestPublishEnvelopeMixedPayerAddressesRejected(t *testing.T) {
	suite := apiTestUtils.NewTestAPIServer(t)

	nodeID := envelopeTestUtils.DefaultClientEnvelopeNodeID

	// Create two envelopes with different signers (different payer addresses)
	key1 := testutils.RandomPrivateKey(t)
	key2 := testutils.RandomPrivateKey(t)

	env1 := envelopeTestUtils.CreatePayerEnvelopeWithSigner(
		t, nodeID, key1, 30, envelopeTestUtils.CreateClientEnvelope(),
	)
	env2 := envelopeTestUtils.CreatePayerEnvelopeWithSigner(
		t, nodeID, key2, 30, envelopeTestUtils.CreateClientEnvelope(),
	)

	// Should fail with InvalidArgument — mixed payers
	_, err := suite.ClientReplication.PublishPayerEnvelopes(
		context.Background(),
		connect.NewRequest(&message_api.PublishPayerEnvelopesRequest{
			PayerEnvelopes: []*envelopes.PayerEnvelope{env1, env2},
		}),
	)
	require.Error(t, err)
	require.Equal(t, connect.CodeInvalidArgument, connect.CodeOf(err))
	require.Contains(t, err.Error(), "same payer")
}
