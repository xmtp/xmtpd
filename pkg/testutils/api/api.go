// Package apiutils implements the api test utils.
package apiutils

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"net/http"
	"testing"
	"time"

	"connectrpc.com/connect"
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
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api/message_apiconnect"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/metadata_api/metadata_apiconnect"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/payer_api/payer_apiconnect"
	"github.com/xmtp/xmtpd/pkg/registrant"
	"github.com/xmtp/xmtpd/pkg/registry"
	"github.com/xmtp/xmtpd/pkg/testutils"
	"github.com/xmtp/xmtpd/pkg/testutils/fees"
	networkTestUtils "github.com/xmtp/xmtpd/pkg/testutils/network"
	"github.com/xmtp/xmtpd/pkg/utils"
)

// Fix! : Implement a testutils client packages with all types of clients.

/* gRPC clients */

func NewTestGRPCReplicationAPIClient(
	t *testing.T,
	addr string,
	extraDialOpts ...connect.ClientOption,
) message_apiconnect.ReplicationApiClient {
	_, port, err := net.SplitHostPort(addr)
	require.NoError(t, err)

	client, err := utils.NewConnectGRPCReplicationAPIClient(
		t.Context(),
		fmt.Sprintf("http://localhost:%s", port),
		extraDialOpts...,
	)
	if err != nil {
		t.Fatalf("failed to create replication API client: %v", err)
	}

	return client
}

func NewTestGRPCGatewayAPIClient(
	t *testing.T,
	addr string,
	extraDialOpts ...connect.ClientOption,
) payer_apiconnect.PayerApiClient {
	_, port, err := net.SplitHostPort(addr)
	require.NoError(t, err)

	client, err := utils.NewConnectGatewayAPIClient(
		t.Context(),
		fmt.Sprintf("http://localhost:%s", port),
		extraDialOpts...,
	)
	if err != nil {
		t.Fatalf("failed to create gateway API client: %v", err)
	}

	return client
}

func NewTestMetadataAPIClient(
	t *testing.T,
	addr string,
	extraDialOpts ...connect.ClientOption,
) metadata_apiconnect.MetadataApiClient {
	_, port, err := net.SplitHostPort(addr)
	require.NoError(t, err)

	client, err := utils.NewConnectMetadataAPIClient(
		t.Context(),
		fmt.Sprintf("http://localhost:%s", port),
		extraDialOpts...,
	)
	if err != nil {
		t.Fatalf("failed to create metadata API client: %v", err)
	}

	return client
}

type APIServerMocks struct {
	MockRegistry          *mocks.MockNodeRegistry
	MockValidationService *mlsvalidateMocks.MockMLSValidationService
	MockMessagePublisher  *blockchain.MockIBlockchainPublisher
}

type APIServerTestSuite struct {
	APIServer         *api.APIServer
	ClientReplication message_apiconnect.ReplicationApiClient
	ClientPayer       payer_apiconnect.PayerApiClient
	ClientMetadata    metadata_apiconnect.MetadataApiClient
	DB                *sql.DB
	APIServerMocks    APIServerMocks
}

// NewTestAPIServer creates a full API server with all services.
// It creates a mock database, mock registry, mock validation service, mock message publisher,
// and mock API server.
// It returns the mock API server, mock database, and mock API server mocks.
func NewTestAPIServer(
	t *testing.T,
) *APIServerTestSuite {
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

	authInterceptor := server.NewServerAuthInterceptor(jwtVerifier, log)

	serviceRegistrationFunc := func(mux *http.ServeMux, interceptors ...connect.Interceptor) (servicePaths []string, err error) {
		interceptors = append(interceptors, authInterceptor)

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

		replicationPath, replicationHandler := message_apiconnect.NewReplicationApiHandler(
			replicationService,
			connect.WithInterceptors(interceptors...),
		)

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

		payerPath, payerHandler := payer_apiconnect.NewPayerApiHandler(
			payerService,
			connect.WithInterceptors(interceptors...),
		)

		metadataService, err := metadata.NewMetadataAPIService(
			ctx,
			log,
			metadata.NewCursorUpdater(ctx, log, db),
			testutils.GetLatestVersion(t),
			metadata.NewPayerInfoFetcher(db),
		)
		require.NoError(t, err)

		metadataPath, metadataHandler := metadata_apiconnect.NewMetadataApiHandler(
			metadataService,
			connect.WithInterceptors(interceptors...),
		)

		mux.Handle(replicationPath, replicationHandler)
		mux.Handle(payerPath, payerHandler)
		mux.Handle(metadataPath, metadataHandler)

		return []string{
			message_apiconnect.ReplicationApiName,
			payer_apiconnect.PayerApiName,
			metadata_apiconnect.MetadataApiName,
		}, nil
	}

	port := networkTestUtils.OpenFreePort(t)

	apiOpts := []api.APIServerOption{
		api.WithContext(ctx),
		api.WithLogger(log),
		api.WithPort(port),
		api.WithRegistrationFunc(serviceRegistrationFunc),
		api.WithReflection(true),
		api.WithPrometheusRegistry(prometheus.NewRegistry()),
	}

	svr, err := api.NewAPIServer(apiOpts...)
	require.NoError(t, err)

	svr.Start()

	allMocks := APIServerMocks{
		MockRegistry:          mockRegistry,
		MockValidationService: mockValidationService,
		MockMessagePublisher:  mockMessagePublisher,
	}

	t.Cleanup(func() {
		cancel()
		svr.Close()
	})

	clientReplication := NewTestGRPCReplicationAPIClient(t, svr.Addr())
	clientPayer := NewTestGRPCGatewayAPIClient(t, svr.Addr())
	clientMetadata := NewTestMetadataAPIClient(t, svr.Addr())

	return &APIServerTestSuite{
		APIServer:         svr,
		APIServerMocks:    allMocks,
		ClientReplication: clientReplication,
		ClientPayer:       clientPayer,
		ClientMetadata:    clientMetadata,
		DB:                db,
	}
}
