package node_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	proto "github.com/xmtp/proto/v3/go/message_api/v1"
	messagev1 "github.com/xmtp/xmtpd/pkg/api/message/v1"
	"github.com/xmtp/xmtpd/pkg/api/message/v1/client"
	"github.com/xmtp/xmtpd/pkg/node"
	memstore "github.com/xmtp/xmtpd/pkg/store/mem"
	test "github.com/xmtp/xmtpd/pkg/testing"
	"github.com/xmtp/xmtpd/pkg/zap"
)

func TestNetwork(t *testing.T) {
	ctx := context.Background()
	net := newTestNetwork(t)

	_, n1c := net.NewNode(t)
	envs := []*proto.Envelope{
		{ContentTopic: "topic1", Message: []byte("msg1")},
		{ContentTopic: "topic1", Message: []byte("msg2")},
		{ContentTopic: "topic1", Message: []byte("msg3")},
	}
	_, err := n1c.Publish(ctx, &proto.PublishRequest{
		Envelopes: envs,
	})
	require.NoError(t, err)

	_, n2c := net.NewNode(t)
	envs = []*proto.Envelope{
		{ContentTopic: "topic2", Message: []byte("msg1")},
		{ContentTopic: "topic2", Message: []byte("msg2")},
		{ContentTopic: "topic2", Message: []byte("msg3")},
	}
	_, err = n2c.Publish(ctx, &proto.PublishRequest{
		Envelopes: envs,
	})
	require.NoError(t, err)

	// err = n1.Connect(n2)
	// require.NoError(t, err)
}

// testNetwork is an in-memory simulation of a network of a given number of Nodes.
// It also captures events that were published to it for final analysis of the test results.
type testNetwork struct {
	ctx    context.Context
	cancel context.CancelFunc
	log    *zap.Logger

	nodes map[string]*node.Node // all the nodes in the network
}

func newTestNetwork(t *testing.T) *testNetwork {
	ctx, cancel := context.WithCancel(context.Background())
	log := test.NewLogger(t)

	return &testNetwork{
		ctx:    ctx,
		cancel: cancel,
		log:    log,

		nodes: map[string]*node.Node{},
	}
}

func (net *testNetwork) Logger() *zap.Logger {
	return net.log
}

func (net *testNetwork) NewNode(t *testing.T) (*node.Node, client.Client) {
	var name string
	for _, ok := net.nodes[name]; name == "" || ok; {
		name = "node-" + test.RandomStringLower(5)
	}
	log := net.log.Named(name)
	store := memstore.New(log)
	messagev1, err := messagev1.New(log, store)
	require.NoError(t, err)
	node, err := node.New(net.ctx, net.log, messagev1, &node.Options{})
	require.NoError(t, err)
	net.nodes[name] = node

	client := client.NewHTTPClient(net.Logger(), fmt.Sprintf("http://localhost:%d", node.APIHTTPListenPort()), "test", "test")

	return node, client
}

func (net *testNetwork) RemoveNode(t *testing.T, name string) *node.Node {
	node, ok := net.nodes[name]
	if !ok {
		return nil
	}
	delete(net.nodes, name)
	return node
}
