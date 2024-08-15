package api

import (
	"context"
	"database/sql"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/mocks"
	"github.com/xmtp/xmtpd/pkg/proto/identity/associations"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api"
	"github.com/xmtp/xmtpd/pkg/registrant"
	"github.com/xmtp/xmtpd/pkg/registry"
	test "github.com/xmtp/xmtpd/pkg/testing"
	"google.golang.org/protobuf/proto"
)

func newTestService(t *testing.T) (*Service, *sql.DB, func()) {
	ctx := context.Background()
	log := test.NewLog(t)
	db, _, dbCleanup := test.NewDB(t, ctx)
	privKey, err := crypto.GenerateKey()
	require.NoError(t, err)
	privKeyStr := "0x" + test.HexEncode(crypto.FromECDSA(privKey))
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

func TestSimplePublish(t *testing.T) {
	svc, _, cleanup := newTestService(t)
	defer cleanup()

	resp, err := svc.PublishEnvelope(
		context.Background(),
		&message_api.PublishEnvelopeRequest{
			PayerEnvelope: &message_api.PayerEnvelope{
				UnsignedClientEnvelope: []byte{0x5},
				PayerSignature:         &associations.RecoverableEcdsaSignature{},
			},
		},
	)
	require.NoError(t, err)
	require.NotNil(t, resp)

	unsignedEnv := &message_api.UnsignedOriginatorEnvelope{}
	require.NoError(
		t,
		proto.Unmarshal(resp.GetOriginatorEnvelope().GetUnsignedOriginatorEnvelope(), unsignedEnv),
	)
	require.Equal(t, uint8(0x5), unsignedEnv.GetPayerEnvelope().GetUnsignedClientEnvelope()[0])

	// TODO(rich) Test that the published envelope is retrievable via the query API
}
