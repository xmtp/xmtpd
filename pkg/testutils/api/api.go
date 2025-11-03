// Package apiutils implements the api test utils.
package apiutils

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/xmtp/xmtpd/pkg/config"
	"github.com/xmtp/xmtpd/pkg/interceptors/server"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/api"
	"github.com/xmtp/xmtpd/pkg/api/message"
	"github.com/xmtp/xmtpd/pkg/api/metadata"
	"github.com/xmtp/xmtpd/pkg/api/payer"
	"github.com/xmtp/xmtpd/pkg/authn"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/mocks/blockchain"

	mlsvalidateMocks "github.com/xmtp/xmtpd/pkg/mocks/mlsvalidate"
	mocks "github.com/xmtp/xmtpd/pkg/mocks/registry"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/metadata_api"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/payer_api"
	"github.com/xmtp/xmtpd/pkg/registrant"
	"github.com/xmtp/xmtpd/pkg/registry"
	"github.com/xmtp/xmtpd/pkg/testutils"
	"github.com/xmtp/xmtpd/pkg/testutils/fees"
	"github.com/xmtp/xmtpd/pkg/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewReplicationAPIClient(
	t *testing.T,
	addr string,
) message_api.ReplicationApiClient {
	// https://github.com/grpc/grpc/blob/master/doc/naming.md
	dialAddr := fmt.Sprintf("passthrough://localhost/%s", addr)
	conn, err := grpc.NewClient(
		dialAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(),
	)
	require.NoError(t, err)
	client := message_api.NewReplicationApiClient(conn)
	t.Cleanup(func() {
		require.NoError(t, conn.Close())
	})
	return client
}

func NewPayerAPIClient(
	t *testing.T,
	addr string,
) payer_api.PayerApiClient {
	dialAddr := fmt.Sprintf("passthrough://localhost/%s", addr)
	conn, err := grpc.NewClient(
		dialAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(),
	)
	require.NoError(t, err)
	client := payer_api.NewPayerApiClient(conn)
	t.Cleanup(func() {
		require.NoError(t, conn.Close())
	})
	return client
}

func NewMetadataAPIClient(
	t *testing.T,
	addr string,
) metadata_api.MetadataApiClient {
	dialAddr := fmt.Sprintf("passthrough://localhost/%s", addr)
	conn, err := grpc.NewClient(
		dialAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(),
	)
	require.NoError(t, err)
	client := metadata_api.NewMetadataApiClient(conn)
	t.Cleanup(func() {
		require.NoError(t, conn.Close())
	})
	return client
}

type APIServerMocks struct {
	MockRegistry          *mocks.MockNodeRegistry
	MockValidationService *mlsvalidateMocks.MockMLSValidationService
	MockMessagePublisher  *blockchain.MockIBlockchainPublisher
}

func NewTestAPIServer(t *testing.T) (*api.APIServer, *sql.DB, APIServerMocks) {
	ctx, cancel := context.WithCancel(context.Background())
	log := testutils.NewLog(t)
	db, _ := testutils.NewDB(t, ctx)
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

	jwtVerifier, err := authn.NewRegistryVerifier(
		log,
		mockRegistry,
		registrant.NodeID(),
		testutils.GetLatestVersion(t),
	)
	require.NoError(t, err)

	serviceRegistrationFunc := func(grpcServer *grpc.Server) error {
		replicationService, err := message.NewReplicationAPIService(
			ctx,
			log,
			registrant,
			db,
			mockValidationService,
			metadata.NewCursorUpdater(ctx, log, db),
			fees.NewTestFeeCalculator(),
			config.APIOptions{
				SendKeepAliveInterval: 30 * time.Second,
			},
			false,
		)
		require.NoError(t, err)
		message_api.RegisterReplicationApiServer(grpcServer, replicationService)

		payerService, err := payer.NewPayerAPIService(
			ctx,
			log,
			mockRegistry,
			privKey,
			mockMessagePublisher,
			nil,
			0,
		)
		require.NoError(t, err)
		payer_api.RegisterPayerApiServer(grpcServer, payerService)

		metadataService, err := metadata.NewMetadataAPIService(
			ctx,
			log,
			metadata.NewCursorUpdater(ctx, log, db),
			testutils.GetLatestVersion(t),
			metadata.NewPayerInfoFetcher(db),
		)
		require.NoError(t, err)
		metadata_api.RegisterMetadataApiServer(grpcServer, metadataService)

		return nil
	}

	grpcListener, err := net.Listen("tcp", "localhost:0")
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = grpcListener.Close()
	})

	apiOpts := []api.APIServerOption{
		api.WithContext(ctx),
		api.WithLogger(log),
		api.WithGRPCListener(grpcListener),
		api.WithRegistrationFunc(serviceRegistrationFunc),
		api.WithReflection(true),
		api.WithPrometheusRegistry(prometheus.NewRegistry()),
	}

	// Add auth interceptors
	authInterceptor := server.NewAuthInterceptor(jwtVerifier, log)
	apiOpts = append(apiOpts,
		api.WithUnaryInterceptors(authInterceptor.Unary()),
		api.WithStreamInterceptors(authInterceptor.Stream()),
	)

	svr, err := api.NewAPIServer(apiOpts...)
	require.NoError(t, err)

	allMocks := APIServerMocks{
		MockRegistry:          mockRegistry,
		MockValidationService: mockValidationService,
		MockMessagePublisher:  mockMessagePublisher,
	}

	t.Cleanup(func() {
		cancel()
		svr.Close(0)
	})

	return svr, db, allMocks
}

func NewTestReplicationAPIClient(
	t *testing.T,
) (message_api.ReplicationApiClient, *sql.DB, APIServerMocks) {
	svc, db, allMocks := NewTestAPIServer(t)
	client := NewReplicationAPIClient(t, svc.Addr().String())
	return client, db, allMocks
}

func NewTestMetadataAPIClient(
	t *testing.T,
) (metadata_api.MetadataApiClient, *sql.DB, APIServerMocks) {
	svc, db, allMocks := NewTestAPIServer(t)
	client := NewMetadataAPIClient(t, svc.Addr().String())
	return client, db, allMocks
}
