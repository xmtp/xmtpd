package payer_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/api/payer"
	"github.com/xmtp/xmtpd/pkg/registry"
	"github.com/xmtp/xmtpd/pkg/testutils"
	apiTestUtils "github.com/xmtp/xmtpd/pkg/testutils/api"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func formatAddress(addr string) string {
	chunks := strings.Split(addr, ":")
	return fmt.Sprintf("http://localhost:%s", chunks[len(chunks)-1])
}

func TestClientManager(t *testing.T) {
	server1, _, cleanup1 := apiTestUtils.NewTestAPIServer(t)
	defer cleanup1()
	server2, _, cleanup2 := apiTestUtils.NewTestAPIServer(t)
	defer cleanup2()

	nodeRegistry := registry.NewFixedNodeRegistry([]registry.Node{
		{
			NodeID:      100,
			HttpAddress: formatAddress(server1.Addr().String()),
		},
		{
			NodeID:      200,
			HttpAddress: formatAddress(server2.Addr().String()),
		},
	})

	cm := payer.NewClientManager(testutils.NewLog(t), nodeRegistry)

	client1, err := cm.GetClient(100)
	require.NoError(t, err)
	require.NotNil(t, client1)

	healthClient := grpc_health_v1.NewHealthClient(client1)
	healthResponse, err := healthClient.Check(
		context.Background(),
		&grpc_health_v1.HealthCheckRequest{},
	)
	require.NoError(t, err)
	require.Equal(t, grpc_health_v1.HealthCheckResponse_SERVING, healthResponse.Status)

	client2, err := cm.GetClient(200)
	require.NoError(t, err)
	require.NotNil(t, client2)

	_, err = cm.GetClient(300)
	require.Error(t, err)

}
