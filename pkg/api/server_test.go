package api

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	messagev1 "github.com/xmtp/xmtpd/pkg/api/message/v1"
	"github.com/xmtp/xmtpd/pkg/crdt"
	membroadcaster "github.com/xmtp/xmtpd/pkg/crdt/broadcasters/mem"
	memstore "github.com/xmtp/xmtpd/pkg/crdt/stores/mem"
	memsyncer "github.com/xmtp/xmtpd/pkg/crdt/syncers/mem"
	"github.com/xmtp/xmtpd/pkg/crdt/types"
	memsubs "github.com/xmtp/xmtpd/pkg/node/subscribers/mem"
	memtopics "github.com/xmtp/xmtpd/pkg/node/topics/mem"
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
	resp, err := http.Post(server.httpListenAddr(), "application/json", nil)
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
	bc := membroadcaster.New(log)
	syncer := memsyncer.New(log, store)
	subs := memsubs.New(log, 100)
	topics, err := memtopics.New(log, func(topicId string) (*crdt.Replica, error) {
		return crdt.NewReplica(ctx, log, store, bc, syncer,
			func(ev *types.Event) {
				subs.OnNewEvent(topicId, ev)
			},
		)
	})
	require.NoError(t, err)
	messagev1, err := messagev1.New(log, topics, subs, store, bc, syncer)
	require.NoError(t, err)
	s, err := New(ctx, log, messagev1, &Options{
		GRPCAddress: "localhost",
		GRPCPort:    0,
		HTTPAddress: "localhost",
		HTTPPort:    0,
		MaxMsgSize:  testMaxMsgSize,
	})
	require.NoError(t, err)
	return s, func() {
		s.Close()
		messagev1.Close()
		topics.Close()
		subs.Close()
		syncer.Close()
		bc.Close()
		store.Close()
	}
}
