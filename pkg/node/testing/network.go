package testing

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"

	proto "github.com/xmtp/proto/v3/go/message_api/v1"
)

type testNetwork struct {
	nodes []*testNode
}

func NewTestNetwork(t *testing.T, count int) *testNetwork {
	t.Helper()
	nodes := make([]*testNode, count)
	for i := 0; i < count; i++ {
		nodes[i] = NewTestNodeWithName(t, fmt.Sprintf("node%d", i+1))
	}

	var wg sync.WaitGroup
	for _, a := range nodes {
		a := a
		for _, b := range nodes {
			b := b
			if a == b {
				continue
			}
			wg.Add(1)
			go func() {
				defer wg.Done()
				a.Connect(t, b)
			}()
		}
	}
	wg.Wait()

	return &testNetwork{
		nodes: nodes,
	}
}

func (net *testNetwork) Close() error {
	for _, node := range net.nodes {
		node.Close()
	}
	return nil
}

func (net *testNetwork) PublishRandom(t *testing.T, topic string, count int) []*proto.Envelope {
	t.Helper()
	node := net.nodes[rand.Intn(len(net.nodes))]
	return node.PublishRandom(t, topic, count)
}

func (net *testNetwork) RequireStoredEvents(t *testing.T, topic string, expected []*proto.Envelope) {
	for _, node := range net.nodes {
		node.RequireStoredEvents(t, topic, expected)
	}
}

type testNetworkSubscriber struct {
	Topic string
	subs  []*testSubscriber
}

func (net *testNetwork) Subscribe(t *testing.T, topic string) *testNetworkSubscriber {
	t.Helper()
	subs := make([]*testSubscriber, len(net.nodes))
	for i, node := range net.nodes {
		subs[i] = node.Subscribe(t, topic)
	}
	return &testNetworkSubscriber{
		Topic: topic,
		subs:  subs,
	}
}

func (s *testNetworkSubscriber) RequireEventuallyCapturedEvents(t *testing.T, expected []*proto.Envelope) {
	t.Helper()
	for _, sub := range s.subs {
		sub.RequireEventuallyCapturedEvents(t, expected)
	}
}
