package testutils

import (
	"context"
	"database/sql"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/api"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/mocks/blockchain"
	mocks "github.com/xmtp/xmtpd/pkg/mocks/registry"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api"
	"github.com/xmtp/xmtpd/pkg/registrant"
	"github.com/xmtp/xmtpd/pkg/registry"
	"github.com/xmtp/xmtpd/pkg/testutils"
)

func NewTestAPIServer(t *testing.T) (*api.ApiServer, *sql.DB, func()) {
	ctx, cancel := context.WithCancel(context.Background())
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
	mockMessagePublsiher := blockchain.NewMockIMessagePublisher(t)

	svr, err := api.NewAPIServer(
		ctx,
		db,
		log,
		0, /*port*/
		registrant,
		true, /*enableReflection*/
		mockMessagePublsiher,
	)
	require.NoError(t, err)

	return svr, db, func() {
		cancel()
		svr.Close()
		dbCleanup()
	}
}

func NewTestAPIClient(t *testing.T) (message_api.ReplicationApiClient, *sql.DB, func()) {
	svc, db, cleanup := NewTestAPIServer(t)
	conn, err := svc.DialGRPCTest(context.Background())
	require.NoError(t, err)
	client := message_api.NewReplicationApiClient(conn)

	return client, db, func() {
		conn.Close()
		require.NoError(t, err)
		cleanup()
	}
}