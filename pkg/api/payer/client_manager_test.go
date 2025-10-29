package payer_test

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	registry2 "github.com/xmtp/xmtpd/pkg/testutils/registry"

	"github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"

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
	server1, _, _ := apiTestUtils.NewTestAPIServer(t)
	server2, _, _ := apiTestUtils.NewTestAPIServer(t)

	mockRegistry := registry2.CreateMockRegistry(t, []registry.Node{
		{
			NodeID:      100,
			HTTPAddress: formatAddress(server1.Addr()),
		},
		{
			NodeID:      200,
			HTTPAddress: formatAddress(server2.Addr()),
		},
	})

	mockRegistry.On("GetNode", uint32(300)).Maybe().Return(nil, errors.New("node not found"))

	cm := payer.NewClientManager(testutils.NewLog(t), mockRegistry, prometheus.NewClientMetrics())

	client1, err := cm.GetClientConnection(100)
	require.NoError(t, err)
	require.NotNil(t, client1)

	healthClient := grpc_health_v1.NewHealthClient(client1)
	healthResponse, err := healthClient.Check(
		context.Background(),
		&grpc_health_v1.HealthCheckRequest{},
	)
	require.NoError(t, err)
	require.Equal(t, grpc_health_v1.HealthCheckResponse_SERVING, healthResponse.Status)

	client2, err := cm.GetClientConnection(200)
	require.NoError(t, err)
	require.NotNil(t, client2)

	_, err = cm.GetClientConnection(300)
	require.Error(t, err)
}
