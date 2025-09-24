package server_test

import (
	"fmt"
	"net"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health/grpc_health_v1"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
)

const (
	server1NodeID = uint32(100)
	server2NodeID = uint32(200)
)

func TestCreateServer(t *testing.T) {
	ctx := t.Context()
	dbs := testutils.NewDBs(t, ctx, 2)
	privateKey1, err := crypto.GenerateKey()
	require.NoError(t, err)
	privateKey2, err := crypto.GenerateKey()
	require.NoError(t, err)

	server1Port := networkTestUtils.OpenFreePort(t)
	server2Port := networkTestUtils.OpenFreePort(t)

	nodes := []r.Node{
		registryTestUtils.CreateNode(
			server1NodeID,
			server1Port.Addr().(*net.TCPAddr).Port,
			privateKey1,
		),
		registryTestUtils.CreateNode(
			server2NodeID,
			server2Port.Addr().(*net.TCPAddr).Port,
			privateKey2,
		),
	}

	registry := registryTestUtils.CreateMockRegistry(t, nodes)

	wsURL, rpcURL := anvil.StartAnvil(t, false)

	contractsOptions := testutils.NewContractsOptions(t, rpcURL, wsURL)

	server1 := serverTestUtils.NewTestReplicationServer(
		t,
		serverTestUtils.TestServerCfg{
			GRPCListener:     server1Port,
			DB:               dbs[0],
			Registry:         registry,
			PrivateKey:       privateKey1,
			ContractsOptions: contractsOptions,
			Services: serverTestUtils.EnabledServices{
				API:  true,
				Sync: true,
			},
		},
	)
	server2 := serverTestUtils.NewTestReplicationServer(
		t,
		serverTestUtils.TestServerCfg{
			GRPCListener:     server2Port,
			DB:               dbs[1],
			Registry:         registry,
			PrivateKey:       privateKey2,
			ContractsOptions: contractsOptions,
			Services: serverTestUtils.EnabledServices{
				API:  true,
				Sync: true,
			},
		},
	)

	require.NotEqual(t, server1.Addr(), server2.Addr())

	defer func() {
		server1.Shutdown(0)
		server2.Shutdown(0)
	}()

	client1 := apiTestUtils.NewReplicationAPIClient(t, server1.Addr().String())
	client2 := apiTestUtils.NewReplicationAPIClient(t, server2.Addr().String())
	nodeID1 := server1NodeID
	nodeID2 := server2NodeID

	targetTopic := topic.NewTopic(topic.TopicKindGroupMessagesV1, []byte{1, 2, 3}).
		Bytes()

	p1, err := client1.PublishPayerEnvelopes(
		ctx,
		&message_api.PublishPayerEnvelopesRequest{
			PayerEnvelopes: []*envelopes.PayerEnvelope{envelopeTestUtils.CreatePayerEnvelope(
				t,
				nodeID1,
				envelopeTestUtils.CreateClientEnvelope(
					&envelopeTestUtils.ClientEnvelopeOptions{Aad: &envelopes.AuthenticatedData{
						TargetTopic: targetTopic,
						DependsOn:   &envelopes.Cursor{},
					}},
				),
			)},
		},
	)
	require.NoError(t, err)
	p2, err := client2.PublishPayerEnvelopes(
		ctx,
		&message_api.PublishPayerEnvelopesRequest{
			PayerEnvelopes: []*envelopes.PayerEnvelope{envelopeTestUtils.CreatePayerEnvelope(
				t,
				nodeID2,
				envelopeTestUtils.CreateClientEnvelope(
					&envelopeTestUtils.ClientEnvelopeOptions{Aad: &envelopes.AuthenticatedData{
						TargetTopic: targetTopic,
						DependsOn:   &envelopes.Cursor{},
					}},
				),
			)},
		},
	)
	require.NoError(t, err)

	require.Eventually(t, func() bool {
		q1, err := client1.QueryEnvelopes(ctx, &message_api.QueryEnvelopesRequest{
			Query: &message_api.EnvelopesQuery{
				OriginatorNodeIds: []uint32{server2NodeID},
				LastSeen:          &envelopes.Cursor{},
			},
			Limit: 10,
		})
		require.NoError(t, err)
		if len(q1.Envelopes) != 1 {
			return false
		}
		if !assert.Equal(t, q1.Envelopes[0], p2.OriginatorEnvelopes[0]) {
			return false
		}
		return true
	}, 3000*time.Millisecond, 200*time.Millisecond)

	require.Eventually(t, func() bool {
		q2, err := client2.QueryEnvelopes(ctx, &message_api.QueryEnvelopesRequest{
			Query: &message_api.EnvelopesQuery{
				OriginatorNodeIds: []uint32{server1NodeID},
				LastSeen:          &envelopes.Cursor{},
			},
			Limit: 10,
		})
		require.NoError(t, err)
		if len(q2.Envelopes) != 1 {
			return false
		}
		if !assert.Equal(t, q2.Envelopes[0], p1.OriginatorEnvelopes[0]) {
			return false
		}
		return true
	}, 3000*time.Millisecond, 200*time.Millisecond)
}

func TestReadOwnWritesGuarantee(t *testing.T) {
	ctx := t.Context()
	dbs := testutils.NewDBs(t, ctx, 1)
	privateKey1, err := crypto.GenerateKey()
	require.NoError(t, err)
	server1Port := networkTestUtils.OpenFreePort(t)

	nodeID1 := server1NodeID

	nodes := []r.Node{
		registryTestUtils.CreateNode(
			server1NodeID,
			server1Port.Addr().(*net.TCPAddr).Port,
			privateKey1,
		),
	}
	registry := registryTestUtils.CreateMockRegistry(t, nodes)
	wsURL, rpcURL := anvil.StartAnvil(t, false)

	contractsOptions := testutils.NewContractsOptions(t, rpcURL, wsURL)

	server1 := serverTestUtils.NewTestReplicationServer(
		t,
		serverTestUtils.TestServerCfg{
			GRPCListener:     server1Port,
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

	client1 := apiTestUtils.NewReplicationAPIClient(t, server1.Addr().String())

	targetTopic := topic.NewTopic(topic.TopicKindGroupMessagesV1, []byte{1, 2, 3}).
		Bytes()

	_, err = client1.PublishPayerEnvelopes(
		ctx,
		&message_api.PublishPayerEnvelopesRequest{
			PayerEnvelopes: []*envelopes.PayerEnvelope{envelopeTestUtils.CreatePayerEnvelope(
				t,
				nodeID1,
				envelopeTestUtils.CreateClientEnvelope(
					&envelopeTestUtils.ClientEnvelopeOptions{Aad: &envelopes.AuthenticatedData{
						TargetTopic: targetTopic,
						DependsOn:   &envelopes.Cursor{},
					}},
				),
			)},
		},
	)
	require.NoError(t, err)

	// query the same server immediately after writing
	// the server should return the write on the first attempt

	q1, err := client1.QueryEnvelopes(ctx, &message_api.QueryEnvelopesRequest{
		Query: &message_api.EnvelopesQuery{
			OriginatorNodeIds: []uint32{server1NodeID},
			LastSeen:          &envelopes.Cursor{},
		},
		Limit: 10,
	})
	require.NoError(t, err)
	require.Len(t, q1.Envelopes, 1)
}

func TestGRPCHealthEndpoint(t *testing.T) {
	ctx := t.Context()
	dbs := testutils.NewDBs(t, ctx, 1)
	privateKey, err := crypto.GenerateKey()
	require.NoError(t, err)

	grpcPort := networkTestUtils.OpenFreePort(t)

	nodes := []r.Node{
		registryTestUtils.CreateNode(
			server1NodeID,
			grpcPort.Addr().(*net.TCPAddr).Port,
			privateKey,
		),
	}
	registry := registryTestUtils.CreateMockRegistry(t, nodes)
	wsURL, rpcURL := anvil.StartAnvil(t, false)
	contractsOptions := testutils.NewContractsOptions(t, rpcURL, wsURL)

	server := serverTestUtils.NewTestReplicationServer(t, serverTestUtils.TestServerCfg{
		GRPCListener:     grpcPort,
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
				fmt.Sprintf("dns:///localhost:%d", grpcPort.Addr().(*net.TCPAddr).Port),
				grpc.WithTransportCredentials(insecure.NewCredentials()),
			)
			if err != nil {
				return false
			}
			defer func() { _ = conn.Close() }()

			healthClient := grpc_health_v1.NewHealthClient(conn)
			grpcResp, err = healthClient.Check(ctx, &grpc_health_v1.HealthCheckRequest{})
			return err == nil && grpcResp.GetStatus() == grpc_health_v1.HealthCheckResponse_SERVING
		}, 3*time.Second, 100*time.Millisecond)
	})
}
