// Package apiutils implements the api test utils.
package apiutils

import (
	"context"
	"database/sql"
	"net"
	"net/http"
	"testing"
	"time"

	"connectrpc.com/connect"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/xmtp/xmtpd/pkg/config"
	dbPkg "github.com/xmtp/xmtpd/pkg/db"
	"github.com/xmtp/xmtpd/pkg/interceptors/server"
	ledgerPkg "github.com/xmtp/xmtpd/pkg/ledger"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/api"
	"github.com/xmtp/xmtpd/pkg/api/message"
	"github.com/xmtp/xmtpd/pkg/api/metadata"
	"github.com/xmtp/xmtpd/pkg/api/payer"
	"github.com/xmtp/xmtpd/pkg/authn"
	blockchainMocks "github.com/xmtp/xmtpd/pkg/testutils/mocks/blockchain"

	gateway_apiconnect "github.com/xmtp/xmtpd/pkg/proto/xmtpv4/gateway_api/gateway_apiconnect"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api/message_apiconnect"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/metadata_api/metadata_apiconnect"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/payer_api/payer_apiconnect"
	"github.com/xmtp/xmtpd/pkg/registrant"
	"github.com/xmtp/xmtpd/pkg/registry"
	"github.com/xmtp/xmtpd/pkg/testutils"
	"github.com/xmtp/xmtpd/pkg/testutils/fees"
	mlsvalidateMocks "github.com/xmtp/xmtpd/pkg/testutils/mocks/mlsvalidate"
	registryMocks "github.com/xmtp/xmtpd/pkg/testutils/mocks/registry"
	networkTestUtils "github.com/xmtp/xmtpd/pkg/testutils/network"
	"github.com/xmtp/xmtpd/pkg/utils"
)

// TODO: Create Connect-based clients for all APIs.
// TODO: Create gRPC-Web clients for all APIs.

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
		"http://localhost:"+port,
		extraDialOpts...,
	)
	if err != nil {
		t.Fatalf("failed to create replication API client: %v", err)
	}

	return client
}

func NewTestGRPCNotificationAPIClient(
	t *testing.T,
	addr string,
	extraDialOpts ...connect.ClientOption,
) message_apiconnect.NotificationApiClient {
	_, port, err := net.SplitHostPort(addr)
	require.NoError(t, err)

	client, err := utils.NewConnectGRPCNotificationAPIClient(
		t.Context(),
		"http://localhost:"+port,
		extraDialOpts...,
	)
	if err != nil {
		t.Fatalf("failed to create notification API client: %v", err)
	}

	return client
}

func NewTestGRPCQueryAPIClient(
	t *testing.T,
	addr string,
	extraDialOpts ...connect.ClientOption,
) message_apiconnect.QueryApiClient {
	_, port, err := net.SplitHostPort(addr)
	require.NoError(t, err)

	client, err := utils.NewConnectGRPCQueryAPIClient(
		t.Context(),
		"http://localhost:"+port,
		extraDialOpts...,
	)
	if err != nil {
		t.Fatalf("failed to create query API client: %v", err)
	}

	return client
}

func NewTestGRPCPublishAPIClient(
	t *testing.T,
	addr string,
	extraDialOpts ...connect.ClientOption,
) message_apiconnect.PublishApiClient {
	_, port, err := net.SplitHostPort(addr)
	require.NoError(t, err)

	client, err := utils.NewConnectGRPCPublishAPIClient(
		t.Context(),
		"http://localhost:"+port,
		extraDialOpts...,
	)
	if err != nil {
		t.Fatalf("failed to create publish API client: %v", err)
	}

	return client
}

func NewTestGRPCPayerAPIClient(
	t *testing.T,
	addr string,
	extraDialOpts ...connect.ClientOption,
) payer_apiconnect.PayerApiClient {
	_, port, err := net.SplitHostPort(addr)
	require.NoError(t, err)

	client, err := utils.NewConnectGatewayAPIClient(
		t.Context(),
		"http://localhost:"+port,
		extraDialOpts...,
	)
	if err != nil {
		t.Fatalf("failed to create gateway API client: %v", err)
	}

	return client
}

func NewTestGRPCGatewayAPIClient(
	t *testing.T,
	addr string,
	extraDialOpts ...connect.ClientOption,
) gateway_apiconnect.GatewayApiClient {
	_, port, err := net.SplitHostPort(addr)
	require.NoError(t, err)

	client, err := utils.NewConnectGRPCGatewayAPIClient(
		t.Context(),
		"http://localhost:"+port,
		extraDialOpts...,
	)
	if err != nil {
		t.Fatalf("failed to create gateway API client: %v", err)
	}

	return client
}

func NewTestGRPCMetadataAPIClient(
	t *testing.T,
	addr string,
	extraDialOpts ...connect.ClientOption,
) metadata_apiconnect.MetadataApiClient {
	_, port, err := net.SplitHostPort(addr)
	require.NoError(t, err)

	options := append([]connect.ClientOption{
		connect.WithGRPC(),
	}, extraDialOpts...)

	client, err := utils.NewConnectMetadataAPIClient(
		t.Context(),
		"http://localhost:"+port,
		options...,
	)
	if err != nil {
		t.Fatalf("failed to create metadata API client: %v", err)
	}

	return client
}

type APIServerMocks struct {
	MockRegistry          *registryMocks.MockNodeRegistry
	MockValidationService *mlsvalidateMocks.MockMLSValidationService
	MockMessagePublisher  *blockchainMocks.MockIBlockchainPublisher
}

type APIServerTestSuite struct {
	APIServer          *api.APIServer
	ClientReplication  message_apiconnect.ReplicationApiClient
	ClientNotification message_apiconnect.NotificationApiClient
	ClientQuery        message_apiconnect.QueryApiClient
	ClientPublish      message_apiconnect.PublishApiClient
	ClientPayer        payer_apiconnect.PayerApiClient
	ClientGateway      gateway_apiconnect.GatewayApiClient
	ClientMetadata     metadata_apiconnect.MetadataApiClient
	DB                 *sql.DB
	APIServerMocks     APIServerMocks
	MessageService     *message.Service
}

// APIServerTestConfig allows explicitly setting some components used for tests.
type APIServerTestConfig struct {
	registryNodes               []registry.Node
	requirePayerPositiveBalance bool
}

type TestAPIOption func(*APIServerTestConfig)

func WithRegistryNodes(nodes []registry.Node) TestAPIOption {
	return func(cfg *APIServerTestConfig) {
		cfg.registryNodes = nodes
	}
}

func WithRequirePayerPositiveBalance(enabled bool) TestAPIOption {
	return func(cfg *APIServerTestConfig) {
		cfg.requirePayerPositiveBalance = enabled
	}
}

func createMockRegistry(t *testing.T, nodes []registry.Node) *registryMocks.MockNodeRegistry {
	reg := registryMocks.NewMockNodeRegistry(t)

	reg.EXPECT().GetNodes().Return(nodes, nil)

	// Return a channel for new nodes.
	ch := make(chan []registry.Node)
	reg.EXPECT().OnNewNodes().Return(ch).Maybe()

	return reg
}

// NewTestAPIServer creates a full API server with all services.
// It creates a mock database, mock registry, mock validation service, mock message publisher,
// and mock API server.
// It returns the mock API server, mock database, and mock API server mocks.
func NewTestAPIServer(
	t *testing.T,
	opts ...TestAPIOption,
) *APIServerTestSuite {
	var cfg APIServerTestConfig
	for _, opt := range opts {
		opt(&cfg)
	}

	var (
		ctx, cancel           = context.WithCancel(context.Background())
		log                   = testutils.NewLog(t)
		sqlDB, _              = testutils.NewRawDB(t, ctx)
		db                    = dbPkg.NewDBHandler(sqlDB)
		mockMessagePublisher  = blockchainMocks.NewMockIBlockchainPublisher(t)
		mockValidationService = mlsvalidateMocks.NewMockMLSValidationService(t)
	)

	privKey, err := crypto.GenerateKey()
	require.NoError(t, err)

	privKeyStr := "0x" + utils.HexEncode(crypto.FromECDSA(privKey))

	nodes := append([]registry.Node{
		{NodeID: 100, SigningKey: &privKey.PublicKey, IsCanonical: true},
	}, cfg.registryNodes...)

	mockRegistry := createMockRegistry(t, nodes)

	registrant, err := registrant.NewRegistrant(
		ctx,
		log,
		db.WriteQuery(),
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

	var replicationService *message.Service
	serviceRegistrationFunc := func(mux *http.ServeMux, interceptors ...connect.Interceptor) (servicePaths []string, err error) {
		interceptors = append(interceptors, authInterceptor)

		replicationService, err = message.NewReplicationAPIService(
			ctx,
			log,
			registrant,
			mockRegistry,
			db,
			mockValidationService,
			metadata.NewCursorUpdater(ctx, log, db),
			fees.NewTestFeeCalculator(),
			config.APIOptions{
				SendKeepAliveInterval:       30 * time.Second,
				RequirePayerPositiveBalance: cfg.requirePayerPositiveBalance,
			},
			false,
			10*time.Millisecond,
			dbPkg.NewCachedOriginatorList(db.ReadQuery(), 100*time.Millisecond, log),
			ledgerPkg.NewLedger(log, db),
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
			nil,
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

		notificationPath, notificationHandler := message_apiconnect.NewNotificationApiHandler(
			replicationService,
			connect.WithInterceptors(interceptors...),
		)

		queryPath, queryHandler := message_apiconnect.NewQueryApiHandler(
			replicationService,
			connect.WithInterceptors(interceptors...),
		)
		publishPath, publishHandler := message_apiconnect.NewPublishApiHandler(
			replicationService,
			connect.WithInterceptors(interceptors...),
		)

		gatewayApiPath, gatewayApiHandler := gateway_apiconnect.NewGatewayApiHandler(
			payerService,
			connect.WithInterceptors(interceptors...),
		)

		mux.Handle(replicationPath, replicationHandler)
		mux.Handle(notificationPath, notificationHandler)
		mux.Handle(queryPath, queryHandler)
		mux.Handle(publishPath, publishHandler)
		mux.Handle(payerPath, payerHandler)
		mux.Handle(metadataPath, metadataHandler)
		mux.Handle(gatewayApiPath, gatewayApiHandler)

		return []string{
			message_apiconnect.ReplicationApiName,
			message_apiconnect.NotificationApiName,
			message_apiconnect.QueryApiName,
			message_apiconnect.PublishApiName,
			payer_apiconnect.PayerApiName,
			metadata_apiconnect.MetadataApiName,
			gateway_apiconnect.GatewayApiName,
		}, nil
	}

	ln := networkTestUtils.OpenListener(t)

	apiOpts := []api.APIServerOption{
		api.WithContext(ctx),
		api.WithLogger(log),
		api.WithListener(ln),
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
		svr.Close(10 * time.Second)
	})

	clientReplication := NewTestGRPCReplicationAPIClient(t, svr.Addr())
	clientNotification := NewTestGRPCNotificationAPIClient(t, svr.Addr())
	clientQuery := NewTestGRPCQueryAPIClient(t, svr.Addr())
	clientPublish := NewTestGRPCPublishAPIClient(t, svr.Addr())
	clientPayer := NewTestGRPCPayerAPIClient(t, svr.Addr())
	clientGateway := NewTestGRPCGatewayAPIClient(t, svr.Addr())
	clientMetadata := NewTestGRPCMetadataAPIClient(t, svr.Addr())

	return &APIServerTestSuite{
		APIServer:          svr,
		APIServerMocks:     allMocks,
		ClientReplication:  clientReplication,
		ClientNotification: clientNotification,
		ClientQuery:        clientQuery,
		ClientPublish:      clientPublish,
		ClientPayer:        clientPayer,
		ClientGateway:      clientGateway,
		ClientMetadata:     clientMetadata,
		DB:                 db.DB(),
		MessageService:     replicationService,
	}
}
