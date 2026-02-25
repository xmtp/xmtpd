package server_test

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"connectrpc.com/connect"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health/grpc_health_v1"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/constants"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/envelopes"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api"
	r "github.com/xmtp/xmtpd/pkg/registry"
	"github.com/xmtp/xmtpd/pkg/testutils"
	"github.com/xmtp/xmtpd/pkg/testutils/anvil"
	apiTestUtils "github.com/xmtp/xmtpd/pkg/testutils/api"
	envelopeTestUtils "github.com/xmtp/xmtpd/pkg/testutils/envelopes"
	networkTestUtils "github.com/xmtp/xmtpd/pkg/testutils/network"
	registryTestUtils "github.com/xmtp/xmtpd/pkg/testutils/registry"
	serverTestUtils "github.com/xmtp/xmtpd/pkg/testutils/server"
	"github.com/xmtp/xmtpd/pkg/topic"
	"github.com/xmtp/xmtpd/pkg/utils"
)

const (
	server1NodeID = uint32(100)
	server2NodeID = uint32(200)
)

func TestCreateServer(t *testing.T) {
	var (
		ctx = t.Context()
		dbs = testutils.NewDBs(t, ctx, 2)
	)

	privateKey1, err := crypto.GenerateKey()
	require.NoError(t, err)

	privateKey2, err := crypto.GenerateKey()
	require.NoError(t, err)

	port1 := networkTestUtils.OpenFreePort(t)
	port2 := networkTestUtils.OpenFreePort(t)

	nodes := []r.Node{
		registryTestUtils.CreateNode(
			server1NodeID,
			port1,
			privateKey1,
		),
		registryTestUtils.CreateNode(
			server2NodeID,
			port2,
			privateKey2,
		),
	}

	registry := registryTestUtils.CreateMockRegistry(t, nodes)

	wsURL, rpcURL := anvil.StartAnvil(t, false)

	contractsOptions := testutils.NewContractsOptions(t, rpcURL, wsURL)

	server1 := serverTestUtils.NewTestBaseServer(
		t,
		serverTestUtils.TestServerCfg{
			Port:             port1,
			DB:               dbs[0],
			Registry:         registry,
			PrivateKey:       privateKey1,
			ContractsOptions: contractsOptions,
			Services: serverTestUtils.EnabledServices{
				API:     true,
				Reports: true,
				Sync:    true,
			},
		},
	)

	server2 := serverTestUtils.NewTestBaseServer(
		t,
		serverTestUtils.TestServerCfg{
			Port:             port2,
			DB:               dbs[1],
			Registry:         registry,
			PrivateKey:       privateKey2,
			ContractsOptions: contractsOptions,
			Services: serverTestUtils.EnabledServices{
				API:     true,
				Reports: true,
				Sync:    true,
			},
		},
	)
	require.NotEqual(t, server1.Addr(), server2.Addr())

	defer func() {
		server1.Shutdown(0)
		server2.Shutdown(0)
	}()

	client1 := apiTestUtils.NewTestGRPCReplicationAPIClient(t, server1.Addr())
	client2 := apiTestUtils.NewTestGRPCReplicationAPIClient(t, server2.Addr())
	nodeID1 := server1NodeID
	nodeID2 := server2NodeID

	targetTopic := topic.NewTopic(topic.TopicKindGroupMessagesV1, []byte{1, 2, 3}).
		Bytes()

	payerEnvelope1 := envelopeTestUtils.CreatePayerEnvelope(
		t,
		nodeID1,
		envelopeTestUtils.CreateClientEnvelope(
			&envelopeTestUtils.ClientEnvelopeOptions{Aad: &envelopes.AuthenticatedData{
				TargetTopic: targetTopic,
				DependsOn:   &envelopes.Cursor{},
			}},
		),
	)

	p1, err := client1.PublishPayerEnvelopes(
		ctx,
		&connect.Request[message_api.PublishPayerEnvelopesRequest]{
			Msg: &message_api.PublishPayerEnvelopesRequest{
				PayerEnvelopes: []*envelopes.PayerEnvelope{payerEnvelope1},
			},
		},
	)
	require.NoError(t, err)

	payerEnvelope2 := envelopeTestUtils.CreatePayerEnvelope(
		t,
		nodeID2,
		envelopeTestUtils.CreateClientEnvelope(
			&envelopeTestUtils.ClientEnvelopeOptions{Aad: &envelopes.AuthenticatedData{
				TargetTopic: targetTopic,
				DependsOn:   &envelopes.Cursor{},
			}},
		),
	)

	p2, err := client2.PublishPayerEnvelopes(
		ctx,
		&connect.Request[message_api.PublishPayerEnvelopesRequest]{
			Msg: &message_api.PublishPayerEnvelopesRequest{
				PayerEnvelopes: []*envelopes.PayerEnvelope{payerEnvelope2},
			},
		},
	)
	require.NoError(t, err)

	// NOTE: there might be a collection of PayerReports here on top of the actual envelopes

	require.Eventually(t, func() bool {
		q1, err := client1.QueryEnvelopes(ctx, &connect.Request[message_api.QueryEnvelopesRequest]{
			Msg: &message_api.QueryEnvelopesRequest{
				Query: &message_api.EnvelopesQuery{
					OriginatorNodeIds: []uint32{server2NodeID},
					LastSeen:          &envelopes.Cursor{},
				},
				Limit: 10,
			},
		})
		require.NoError(t, err)

		for _, e := range q1.Msg.GetEnvelopes() {
			if reflect.DeepEqual(e, p2.Msg.GetOriginatorEnvelopes()[0]) {
				return true
			}
		}
		return false
	}, 10*time.Second, 200*time.Millisecond)

	require.Eventually(t, func() bool {
		q2, err := client2.QueryEnvelopes(ctx, &connect.Request[message_api.QueryEnvelopesRequest]{
			Msg: &message_api.QueryEnvelopesRequest{
				Query: &message_api.EnvelopesQuery{
					OriginatorNodeIds: []uint32{server1NodeID},
					LastSeen:          &envelopes.Cursor{},
				},
				Limit: 10,
			},
		})
		require.NoError(t, err)

		for _, e := range q2.Msg.GetEnvelopes() {
			if reflect.DeepEqual(e, p1.Msg.GetOriginatorEnvelopes()[0]) {
				return true
			}
		}
		return false
	}, 5000*time.Millisecond, 200*time.Millisecond)
}

func TestReadOwnWritesGuarantee(t *testing.T) {
	var (
		ctx     = t.Context()
		dbs     = testutils.NewDBs(t, ctx, 1)
		port    = networkTestUtils.OpenFreePort(t)
		nodeID1 = server1NodeID
	)

	privateKey1, err := crypto.GenerateKey()
	require.NoError(t, err)

	nodes := []r.Node{
		registryTestUtils.CreateNode(
			server1NodeID,
			port,
			privateKey1,
		),
	}
	registry := registryTestUtils.CreateMockRegistry(t, nodes)
	wsURL, rpcURL := anvil.StartAnvil(t, false)

	contractsOptions := testutils.NewContractsOptions(t, rpcURL, wsURL)

	server1 := serverTestUtils.NewTestBaseServer(
		t,
		serverTestUtils.TestServerCfg{
			Port:             port,
			DB:               dbs[0],
			Registry:         registry,
			PrivateKey:       privateKey1,
			ContractsOptions: contractsOptions,
			Services: serverTestUtils.EnabledServices{
				API: true,
			},
		},
	)
	defer func() {
		server1.Shutdown(0)
	}()

	client1 := apiTestUtils.NewTestGRPCReplicationAPIClient(t, server1.Addr())

	targetTopic := topic.NewTopic(topic.TopicKindGroupMessagesV1, []byte{1, 2, 3}).
		Bytes()

	payerEnvelope1 := envelopeTestUtils.CreatePayerEnvelope(
		t,
		nodeID1,
		envelopeTestUtils.CreateClientEnvelope(
			&envelopeTestUtils.ClientEnvelopeOptions{Aad: &envelopes.AuthenticatedData{
				TargetTopic: targetTopic,
				DependsOn:   &envelopes.Cursor{},
			}},
		),
	)

	_, err = client1.PublishPayerEnvelopes(
		ctx,
		&connect.Request[message_api.PublishPayerEnvelopesRequest]{
			Msg: &message_api.PublishPayerEnvelopesRequest{
				PayerEnvelopes: []*envelopes.PayerEnvelope{payerEnvelope1},
			},
		},
	)
	require.NoError(t, err)

	// query the same server immediately after writing
	// the server should return the write on the first attempt

	q1, err := client1.QueryEnvelopes(ctx, &connect.Request[message_api.QueryEnvelopesRequest]{
		Msg: &message_api.QueryEnvelopesRequest{
			Query: &message_api.EnvelopesQuery{
				OriginatorNodeIds: []uint32{server1NodeID},
				LastSeen:          &envelopes.Cursor{},
			},
			Limit: 10,
		},
	})
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(q1.Msg.GetEnvelopes()), 1)
}

func TestGRPCHealthEndpoint(t *testing.T) {
	var (
		ctx  = t.Context()
		dbs  = testutils.NewDBs(t, ctx, 1)
		port = networkTestUtils.OpenFreePort(t)
	)

	privateKey, err := crypto.GenerateKey()
	require.NoError(t, err)

	nodes := []r.Node{
		registryTestUtils.CreateNode(
			server1NodeID,
			port,
			privateKey,
		),
	}
	registry := registryTestUtils.CreateMockRegistry(t, nodes)
	wsURL, rpcURL := anvil.StartAnvil(t, false)
	contractsOptions := testutils.NewContractsOptions(t, rpcURL, wsURL)

	server := serverTestUtils.NewTestBaseServer(t, serverTestUtils.TestServerCfg{
		Port:             port,
		DB:               dbs[0],
		Registry:         registry,
		PrivateKey:       privateKey,
		ContractsOptions: contractsOptions,
		Services: serverTestUtils.EnabledServices{
			API: true,
		},
	})
	defer server.Shutdown(0)

	t.Run("gRPC /v1/health should return SERVING", func(t *testing.T) {
		var grpcResp *grpc_health_v1.HealthCheckResponse

		require.Eventually(t, func() bool {
			conn, err := grpc.NewClient(
				fmt.Sprintf("dns:///localhost:%d", port),
				grpc.WithTransportCredentials(insecure.NewCredentials()),
			)
			if err != nil {
				return false
			}
			defer func() { _ = conn.Close() }()

			healthClient := grpc_health_v1.NewHealthClient(conn)
			grpcResp, err = healthClient.Check(ctx, &grpc_health_v1.HealthCheckRequest{})
			return err == nil && grpcResp.GetStatus() == grpc_health_v1.HealthCheckResponse_SERVING
		}, 10*time.Second, 100*time.Millisecond)
	})
}

func TestCreateServer_AllOptionPermutations(t *testing.T) {
	var (
		ctx          = t.Context()
		serverNodeID = server1NodeID

		// Use a single registry port for all subtests – registry is shared.
		registryPort = networkTestUtils.OpenFreePort(t)
	)

	privateKey, err := crypto.GenerateKey()
	require.NoError(t, err)

	nodes := []r.Node{
		registryTestUtils.CreateNode(
			serverNodeID,
			registryPort,
			privateKey,
		),
	}

	registry := registryTestUtils.CreateMockRegistry(t, nodes)

	wsURL, rpcURL := anvil.StartAnvil(t, false)
	contractsOptions := testutils.NewContractsOptions(t, rpcURL, wsURL)

	for mask := range 16 {
		services := serverTestUtils.EnabledServices{
			API:     mask&1 != 0,
			Reports: mask&2 != 0,
			Sync:    mask&4 != 0,
			Indexer: mask&8 != 0,
		}

		name := fmt.Sprintf(
			"api=%t/reports=%t/sync=%t/indexer=%t",
			services.API,
			services.Reports,
			services.Sync,
			services.Indexer,
		)

		t.Run(name, func(t *testing.T) {
			var (
				port  = networkTestUtils.OpenFreePort(t)
				db, _ = testutils.NewDB(t, ctx)
			)

			server := serverTestUtils.NewTestBaseServer(
				t,
				serverTestUtils.TestServerCfg{
					Port:             port,
					DB:               db.DB(),
					Registry:         registry,
					PrivateKey:       privateKey,
					ContractsOptions: contractsOptions,
					Services:         services,
				},
			)

			require.NotNil(t, server)
			server.Shutdown(0)
		})
	}
}

func TestGRPCPayloadLimit(t *testing.T) {
	var (
		ctx              = t.Context()
		dbs              = testutils.NewDBs(t, ctx, 1)
		port             = networkTestUtils.OpenFreePort(t)
		privateKey       = testutils.RandomPrivateKey(t)
		wsURL, rpcURL    = anvil.StartAnvil(t, false)
		contractsOptions = testutils.NewContractsOptions(t, rpcURL, wsURL)

		nodes = []r.Node{
			registryTestUtils.CreateNode(
				server1NodeID,
				port,
				privateKey,
			),
		}

		registry = registryTestUtils.CreateMockRegistry(t, nodes)

		server = serverTestUtils.NewTestBaseServer(t, serverTestUtils.TestServerCfg{
			Port:             port,
			DB:               dbs[0],
			Registry:         registry,
			PrivateKey:       privateKey,
			ContractsOptions: contractsOptions,
			Services: serverTestUtils.EnabledServices{
				API: true,
			},
		})
	)

	defer server.Shutdown(0)

	largePayload := make([]byte, 500*1024)
	for i := range largePayload {
		largePayload[i] = byte(i % 256)
	}

	totalPayloadSize := 0
	payerEnvelopes := make([]*envelopes.PayerEnvelope, 0)

	for totalPayloadSize < constants.GRPCPayloadLimit {
		clientEnv := envelopeTestUtils.CreateGroupMessageClientEnvelope(
			[16]byte{1, 2, 3},
			largePayload,
		)
		payerEnvelope := envelopeTestUtils.CreatePayerEnvelopeWithSigner(
			t, server1NodeID, privateKey, constants.DefaultStorageDurationDays, clientEnv,
		)
		totalPayloadSize += len(payerEnvelope.GetUnsignedClientEnvelope())
		payerEnvelopes = append(payerEnvelopes, payerEnvelope)
	}

	t.Run("gRPC payload limit should be respected", func(t *testing.T) {
		client, err := utils.NewConnectReplicationAPIClient(
			ctx,
			fmt.Sprintf("http://localhost:%d", port),
		)
		require.NoError(t, err)

		_, err = client.PublishPayerEnvelopes(
			ctx,
			&connect.Request[message_api.PublishPayerEnvelopesRequest]{
				Msg: &message_api.PublishPayerEnvelopesRequest{
					PayerEnvelopes: payerEnvelopes,
				},
			},
		)
		require.Error(t, err)

		code := connect.CodeOf(err)
		require.Equal(t, connect.CodeResourceExhausted, code)
	})
}
