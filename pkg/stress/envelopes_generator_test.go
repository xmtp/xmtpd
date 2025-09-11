package stress

import (
	"context"
	"fmt"
	"net"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/envelopes"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api"
	r "github.com/xmtp/xmtpd/pkg/registry"
	"github.com/xmtp/xmtpd/pkg/testutils"
	"github.com/xmtp/xmtpd/pkg/testutils/anvil"
	apiTestUtils "github.com/xmtp/xmtpd/pkg/testutils/api"
	networkTestUtils "github.com/xmtp/xmtpd/pkg/testutils/network"
	registryTestUtils "github.com/xmtp/xmtpd/pkg/testutils/registry"
	serverTestUtils "github.com/xmtp/xmtpd/pkg/testutils/server"
)

func TestEnvelopesGenerator(t *testing.T) {
	var (
		ctx              = t.Context()
		db, _            = testutils.NewDB(t, ctx)
		grpcPort         = networkTestUtils.OpenFreePort(t)
		httpPort         = networkTestUtils.OpenFreePort(t)
		wsURL, rpcURL    = anvil.StartAnvil(t, false)
		contractsOptions = testutils.NewContractsOptions(t, rpcURL, wsURL)
	)

	privateKey, err := crypto.GenerateKey()
	require.NoError(t, err)

	nodes := []r.Node{
		registryTestUtils.CreateNode(
			100,
			grpcPort.Addr().(*net.TCPAddr).Port,
			privateKey,
		),
	}
	registry := registryTestUtils.CreateMockRegistry(t, nodes)

	server := serverTestUtils.NewTestServer(
		t,
		serverTestUtils.TestServerCfg{
			GRPCListener:     grpcPort,
			HTTPListener:     httpPort,
			Db:               db,
			Registry:         registry,
			PrivateKey:       privateKey,
			ContractsOptions: contractsOptions,
			Services: serverTestUtils.EnabledServices{
				Replication: true,
				Sync:        true,
			},
		},
	)
	defer server.Shutdown(0)

	generator, err := NewEnvelopesGenerator(
		fmt.Sprintf("http://%s", server.Addr().String()),
		testutils.TEST_PRIVATE_KEY,
		100,
	)
	require.NoError(t, err)

	publishResponse, err := generator.PublishWelcomeMessageEnvelopes(context.Background(), 1, 100)
	require.NoError(t, err)
	require.NotNil(t, publishResponse)

	client := apiTestUtils.NewReplicationAPIClient(t, server.Addr().String())
	queryResponse, err := client.QueryEnvelopes(ctx, &message_api.QueryEnvelopesRequest{
		Query: &message_api.EnvelopesQuery{
			OriginatorNodeIds: []uint32{100},
			LastSeen:          &envelopes.Cursor{},
		},
		Limit: 10,
	})
	require.NoError(t, err)
	require.Len(t, queryResponse.Envelopes, 1)
	require.Equal(t, queryResponse.Envelopes[0], publishResponse[0])
}
