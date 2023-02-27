package apigateway

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	memstore "github.com/xmtp/xmtpd/pkg/crdt/stores/mem"
	test "github.com/xmtp/xmtpd/pkg/testing"
	healthgrpc "google.golang.org/grpc/health/grpc_health_v1"
)

const (
	testMaxMsgSize = 2 * 1024 * 1024
)

func Test_HTTPRootPath(t *testing.T) {
	t.Parallel()

	server, cleanup := newTestServer(t)
	defer cleanup()

	// Root path responds with 404.
	var rootRes map[string]interface{}
	resp, err := http.Post(server.HTTPListenAddr(), "application/json", nil)
	require.NoError(t, err)
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	err = json.Unmarshal(body, &rootRes)
	require.NoError(t, err)
	require.Equal(t, map[string]interface{}{
		"code":    float64(5),
		"message": "Not Found",
		"details": []interface{}{},
	}, rootRes)
}

func Test_Health(t *testing.T) {
	ctx := context.Background()
	server, cleanup := newTestServer(t)
	conn, err := server.dialGRPC(ctx)
	assert.NoError(t, err)
	healthClient := healthgrpc.NewHealthClient(conn)

	res, err := healthClient.Check(ctx, &healthgrpc.HealthCheckRequest{})
	assert.NoError(t, err)
	assert.Equal(t, res.Status, healthgrpc.HealthCheckResponse_SERVING)
	cleanup()
}

func newTestServer(t *testing.T) (*Server, func()) {
	ctx := context.Background()
	log := test.NewLogger(t)

	store := memstore.New(log)
	s, err := New(ctx, log, nil, &Options{
		GRPCAddress: "localhost",
		GRPCPort:    0,
		HTTPAddress: "localhost",
		HTTPPort:    0,
		MaxMsgSize:  testMaxMsgSize,
	})
	require.NoError(t, err)
	return s, func() {
		s.Close()
		store.Close()
	}
}
