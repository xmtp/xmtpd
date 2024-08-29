package api

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/mocks"
	"github.com/xmtp/xmtpd/pkg/proto/identity/associations"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api"
	"github.com/xmtp/xmtpd/pkg/registrant"
	"github.com/xmtp/xmtpd/pkg/registry"
	"github.com/xmtp/xmtpd/pkg/testutils"
	"google.golang.org/protobuf/proto"
)

func newTestService(t *testing.T) (*Service, *sql.DB, func()) {
	ctx := context.Background()
	log := testutils.NewLog(t)
	db, _, dbCleanup := testutils.NewDB(t, ctx)
	privKey, err := crypto.GenerateKey()
	require.NoError(t, err)
	privKeyStr := "0x" + testutils.HexEncode(crypto.FromECDSA(privKey))
	mockRegistry := mocks.NewMockNodeRegistry(t)
	mockRegistry.EXPECT().GetNodes().Return([]registry.Node{
		{NodeID: 1, SigningKey: &privKey.PublicKey},
	}, nil)
	registrant, err := registrant.NewRegistrant(ctx, queries.New(db), mockRegistry, privKeyStr)
	require.NoError(t, err)

	svc, err := NewReplicationApiService(ctx, log, registrant, db)
	require.NoError(t, err)

	return svc, db, func() {
		svc.Close()
		dbCleanup()
	}
}

func createClientEnvelope() *message_api.ClientEnvelope {
	return &message_api.ClientEnvelope{
		Payload: nil,
		Aad: &message_api.AuthenticatedData{
			TargetOriginator:   1,
			TargetTopic:        []byte{0x5},
			LastOriginatorSids: []uint64{},
		},
	}
}

func createPayerEnvelope(
	t *testing.T,
	clientEnv ...*message_api.ClientEnvelope,
) *message_api.PayerEnvelope {
	if len(clientEnv) == 0 {
		clientEnv = append(clientEnv, createClientEnvelope())
	}
	clientEnvBytes, err := proto.Marshal(clientEnv[0])
	require.NoError(t, err)

	return &message_api.PayerEnvelope{
		UnsignedClientEnvelope: clientEnvBytes,
		PayerSignature:         &associations.RecoverableEcdsaSignature{},
	}
}

func TestSimplePublish(t *testing.T) {
	svc, db, cleanup := newTestService(t)
	defer cleanup()

	resp, err := svc.PublishEnvelope(
		context.Background(),
		&message_api.PublishEnvelopeRequest{
			PayerEnvelope: createPayerEnvelope(t),
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

func TestUnmarshalError(t *testing.T) {
	svc, _, cleanup := newTestService(t)
	defer cleanup()

	envelope := createPayerEnvelope(t)
	envelope.UnsignedClientEnvelope = []byte("invalidbytes")
	_, err := svc.PublishEnvelope(
		context.Background(),
		&message_api.PublishEnvelopeRequest{
			PayerEnvelope: envelope,
		},
	)
	require.ErrorContains(t, err, "unmarshal")
}

func TestMismatchingOriginator(t *testing.T) {
	svc, _, cleanup := newTestService(t)
	defer cleanup()

	clientEnv := createClientEnvelope()
	clientEnv.Aad.TargetOriginator = 2
	_, err := svc.PublishEnvelope(
		context.Background(),
		&message_api.PublishEnvelopeRequest{
			PayerEnvelope: createPayerEnvelope(t, clientEnv),
		},
	)
	require.ErrorContains(t, err, "originator")
}

func TestMissingTopic(t *testing.T) {
	svc, _, cleanup := newTestService(t)
	defer cleanup()

	clientEnv := createClientEnvelope()
	clientEnv.Aad.TargetTopic = nil
	_, err := svc.PublishEnvelope(
		context.Background(),
		&message_api.PublishEnvelopeRequest{
			PayerEnvelope: createPayerEnvelope(t, clientEnv),
		},
	)
	require.ErrorContains(t, err, "topic")
}
