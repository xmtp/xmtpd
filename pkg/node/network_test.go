package node_test

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"testing"

	proto "github.com/xmtp/proto/v3/go/message_api/v1"
	test "github.com/xmtp/xmtpd/pkg/testing"
	"github.com/xmtp/xmtpd/pkg/zap"
)

func TestNetwork(t *testing.T) {
	net := newTestNetwork(t, 3)
	defer net.Close()

	sub := net.subscribe(t, "topic1")
	envs := net.publishRandom(t, sub.topic, 1)
	sub.requireEventuallyCapturedEvents(t, envs)
	net.requireStoredEvents(t, "topic1", envs)
}

type testNetwork struct {
	ctx    context.Context
	cancel context.CancelFunc
	log    *zap.Logger

	nodes []*testNode
}

func newTestNetwork(t *testing.T, count int) *testNetwork {
	t.Helper()

	ctx, cancel := context.WithCancel(context.Background())
	log := test.NewLogger(t)

	nodes := make([]*testNode, count)
	for i := 0; i < count; i++ {
		nodes[i] = newTestNodeWithName(t, fmt.Sprintf("node%d", i+1))
	}

	var wg sync.WaitGroup
	for _, a := range nodes {
		a := a
		for _, b := range nodes {
			b := b
			if a == b {
				continue
			}
			wg.Add(2)
			go func() {
				defer wg.Done()
				a.connect(t, b)
			}()
			go func() {
				defer wg.Done()
				b.connect(t, a)
			}()
		}
	}
	wg.Wait()

	return &testNetwork{
		ctx:    ctx,
		cancel: cancel,
		log:    log,

		nodes: nodes,
	}
}

func (net *testNetwork) Close() error {
	for _, node := range net.nodes {
		node.Close()
	}
	return nil
}

func (net *testNetwork) publishRandom(t *testing.T, topic string, count int) []*proto.Envelope {
	t.Helper()
	node := net.nodes[rand.Intn(len(net.nodes))]
	return node.publishRandom(t, topic, count)
}

func (net *testNetwork) requireStoredEvents(t *testing.T, topic string, expected []*proto.Envelope) {
	for _, node := range net.nodes {
		node.requireStoredEvents(t, topic, expected)
	}
}

type testNetworkSubscriber struct {
	topic string
	subs  []*testSubscriber
}

func (net *testNetwork) subscribe(t *testing.T, topic string) *testNetworkSubscriber {
	t.Helper()
	subs := make([]*testSubscriber, len(net.nodes))
	for i, node := range net.nodes {
		subs[i] = node.subscribe(t, topic)
	}
	return &testNetworkSubscriber{
		topic: topic,
		subs:  subs,
	}
}

func (s *testNetworkSubscriber) requireEventuallyCapturedEvents(t *testing.T, expected []*proto.Envelope) {
	t.Helper()
	for _, sub := range s.subs {
		sub.requireEventuallyCapturedEvents(t, expected)
	}
}
