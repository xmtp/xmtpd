package testutils

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/api"
	"github.com/xmtp/xmtpd/pkg/api/message"
	"github.com/xmtp/xmtpd/pkg/api/payer"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/mocks/blockchain"
	mocks "github.com/xmtp/xmtpd/pkg/mocks/registry"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/payer_api"
	"github.com/xmtp/xmtpd/pkg/registrant"
	"github.com/xmtp/xmtpd/pkg/registry"
	"github.com/xmtp/xmtpd/pkg/testutils"
	"github.com/xmtp/xmtpd/pkg/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewReplicationAPIClient(
	t *testing.T,
	ctx context.Context,
	addr string,
) (message_api.ReplicationApiClient, func()) {
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

func NewPayerAPIClient(
	t *testing.T,
	ctx context.Context,
	addr string,
) (payer_api.PayerApiClient, func()) {
	dialAddr := fmt.Sprintf("passthrough://localhost/%s", addr)
	conn, err := grpc.NewClient(
		dialAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(),
	)
	require.NoError(t, err)
	client := payer_api.NewPayerApiClient(conn)
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
		{NodeID: 100, SigningKey: &privKey.PublicKey},
	}, nil)
	registrant, err := registrant.NewRegistrant(ctx, log, queries.New(db), mockRegistry, privKeyStr)
	require.NoError(t, err)
	mockMessagePublisher := blockchain.NewMockIBlockchainPublisher(t)

	serviceRegistrationFunc := func(grpcServer *grpc.Server) error {
		replicationService, err := message.NewReplicationApiService(
			ctx,
			log,
			registrant,
			db,
		)
		require.NoError(t, err)
		message_api.RegisterReplicationApiServer(grpcServer, replicationService)

		payerService, err := payer.NewPayerApiService(
			ctx,
			log,
			mockRegistry,
			privKey,
			mockMessagePublisher,
		)
		require.NoError(t, err)
		payer_api.RegisterPayerApiServer(grpcServer, payerService)

		return nil
	}

	svr, err := api.NewAPIServer(
		ctx,
		log,
		"localhost:0", /*listenAddress*/
		true,          /*enableReflection*/
		serviceRegistrationFunc,
	)
	require.NoError(t, err)

	return svr, db, func() {
		cancel()
		svr.Close()
		dbCleanup()
	}
}

func NewTestReplicationAPIClient(t *testing.T) (message_api.ReplicationApiClient, *sql.DB, func()) {
	svc, db, svcCleanup := NewTestAPIServer(t)
	client, clientCleanup := NewReplicationAPIClient(t, context.Background(), svc.Addr().String())
	return client, db, func() {
		clientCleanup()
		svcCleanup()
	}
}
