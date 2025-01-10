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
	"github.com/xmtp/xmtpd/pkg/authn"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/mocks/blockchain"
	mlsvalidateMocks "github.com/xmtp/xmtpd/pkg/mocks/mlsvalidate"
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

type ApiServerMocks struct {
	MockRegistry          *mocks.MockNodeRegistry
	MockValidationService *mlsvalidateMocks.MockMLSValidationService
	MockMessagePublisher  *blockchain.MockIBlockchainPublisher
}

func NewTestAPIServer(t *testing.T) (*api.ApiServer, *sql.DB, ApiServerMocks, func()) {
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
	registrant, err := registrant.NewRegistrant(
		ctx,
		log,
		queries.New(db),
		mockRegistry,
		privKeyStr,
		nil,
	)
	require.NoError(t, err)
	mockMessagePublisher := blockchain.NewMockIBlockchainPublisher(t)
	mockValidationService := mlsvalidateMocks.NewMockMLSValidationService(t)

	jwtVerifier := authn.NewRegistryVerifier(mockRegistry, registrant.NodeID())

	serviceRegistrationFunc := func(grpcServer *grpc.Server) error {
		replicationService, err := message.NewReplicationApiService(
			ctx,
			log,
			registrant,
			db,
			mockValidationService,
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
		jwtVerifier,
	)
	require.NoError(t, err)

	allMocks := ApiServerMocks{
		MockRegistry:          mockRegistry,
		MockValidationService: mockValidationService,
		MockMessagePublisher:  mockMessagePublisher,
	}

	return svr, db, allMocks, func() {
		cancel()
		svr.Close()
		dbCleanup()
	}
}

func NewTestReplicationAPIClient(
	t *testing.T,
) (message_api.ReplicationApiClient, *sql.DB, ApiServerMocks, func()) {
	svc, db, allMocks, svcCleanup := NewTestAPIServer(t)
	client, clientCleanup := NewReplicationAPIClient(t, context.Background(), svc.Addr().String())
	return client, db, allMocks, func() {
		clientCleanup()
		svcCleanup()
	}
}
