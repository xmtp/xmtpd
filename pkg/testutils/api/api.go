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
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api/message_apiconnect"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/metadata_api"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/metadata_api/metadata_apiconnect"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/payer_api"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/payer_api/payer_apiconnect"
	"github.com/xmtp/xmtpd/pkg/registrant"
	"github.com/xmtp/xmtpd/pkg/registry"
	"github.com/xmtp/xmtpd/pkg/testutils"
	"github.com/xmtp/xmtpd/pkg/testutils/fees"
	"github.com/xmtp/xmtpd/pkg/utils"
	"google.golang.org/grpc"
)

func NewTestReplicationAPIClient(
	t *testing.T,
	addr string,
) message_apiconnect.ReplicationApiClient {
	// https://github.com/grpc/grpc/blob/master/doc/naming.md
	dialAddr := fmt.Sprintf("passthrough://localhost/%s", addr)

	httpClient, err := utils.BuildHTTP2Client(t.Context(), false)
	if err != nil {
		t.Fatalf("failed to build HTTP client: %v", err)
	}

	dialOpts := utils.BuildGRPCDialOptions()

	return message_apiconnect.NewReplicationApiClient(httpClient, dialAddr, dialOpts...)
}

func NewTestPayerAPIClient(
	t *testing.T,
	addr string,
) payer_apiconnect.PayerApiClient {
	dialAddr := fmt.Sprintf("passthrough://localhost/%s", addr)

	httpClient, err := utils.BuildHTTP2Client(t.Context(), false)
	if err != nil {
		t.Fatalf("failed to build HTTP client: %v", err)
	}

	dialOpts := utils.BuildGRPCDialOptions()

	return payer_apiconnect.NewPayerApiClient(httpClient, dialAddr, dialOpts...)
}

func NewTestMetadataAPIClient(
	t *testing.T,
	addr string,
) metadata_apiconnect.MetadataApiClient {
	dialAddr := fmt.Sprintf("passthrough://localhost/%s", addr)

	httpClient, err := utils.BuildHTTP2Client(t.Context(), false)
	if err != nil {
		t.Fatalf("failed to build HTTP client: %v", err)
	}
	dialOpts := utils.BuildGRPCDialOptions()

	return metadata_apiconnect.NewMetadataApiClient(httpClient, dialAddr, dialOpts...)
}

type APIServerMocks struct {
	MockRegistry          *mocks.MockNodeRegistry
	MockValidationService *mlsvalidateMocks.MockMLSValidationService
	MockMessagePublisher  *blockchain.MockIBlockchainPublisher
}

func NewTestAPIServer(
	t *testing.T,
) (mockAPIServer *api.APIServer, mockDB *sql.DB, mockAPIServerMocks APIServerMocks) {
	var (
		ctx, cancel           = context.WithCancel(context.Background())
		log                   = testutils.NewLog(t)
		db, _                 = testutils.NewDB(t, ctx)
		mockRegistry          = mocks.NewMockNodeRegistry(t)
		mockMessagePublisher  = blockchain.NewMockIBlockchainPublisher(t)
		mockValidationService = mlsvalidateMocks.NewMockMLSValidationService(t)
	)

	privKey, err := crypto.GenerateKey()
	require.NoError(t, err)

	privKeyStr := "0x" + utils.HexEncode(crypto.FromECDSA(privKey))

	// Mock registry behavior.
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
			10*time.Millisecond,
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
