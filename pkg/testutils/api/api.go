package testutils

import (
	"context"
	"database/sql"
	"fmt"
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
	"github.com/xmtp/xmtpd/pkg/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewAPIClient(t *testing.T, ctx context.Context, addr string) (message_api.ReplicationApiClient, func()) {
	// https://github.com/grpc/grpc/blob/master/doc/naming.md
	dialAddr := fmt.Sprintf("passthrough://localhost/%s", addr)
	conn, err := grpc.NewClient(
		dialAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(),
	)
	require.NoError(t, err)
	client := message_api.NewReplicationApiClient(conn)
	return client, func() {
		err := conn.Close()
		require.NoError(t, err)
	}
}

func NewTestAPIServer(t *testing.T) (*api.ApiServer, *sql.DB, func()) {
	ctx, cancel := context.WithCancel(context.Background())
	log := testutils.NewLog(t)
	db, _, dbCleanup := testutils.NewDB(t, ctx)
	privKey, err := crypto.GenerateKey()
	require.NoError(t, err)
	privKeyStr := "0x" + utils.HexEncode(crypto.FromECDSA(privKey))
	mockRegistry := mocks.NewMockNodeRegistry(t)
	mockRegistry.EXPECT().GetNodes().Return([]registry.Node{
		{NodeID: 1, SigningKey: &privKey.PublicKey},
	}, nil)
	registrant, err := registrant.NewRegistrant(ctx, log, queries.New(db), mockRegistry, privKeyStr)
	require.NoError(t, err)
	mockMessagePublsiher := blockchain.NewMockIBlockchainPublisher(t)

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
	svc, db, svcCleanup := NewTestAPIServer(t)
	client, clientCleanup := NewAPIClient(t, context.Background(), svc.Addr().String())
	return client, db, func() {
		clientCleanup()
		svcCleanup()
	}
}
