package testing

import (
	"bytes"
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	proto "github.com/xmtp/proto/v3/go/message_api/v1"
	"github.com/xmtp/xmtpd/pkg/context"
	"github.com/xmtp/xmtpd/pkg/node"
	memstore "github.com/xmtp/xmtpd/pkg/store/mem"
	test "github.com/xmtp/xmtpd/pkg/testing"
)

type storeMaker func(t testing.TB, ctx context.Context) node.NodeStore

type networkOption func(n *network)

func WithStoreMaker(sm storeMaker) networkOption {
	return func(n *network) {
		n.storeMaker = sm
	}
}

type envsByTopic map[string][]*proto.Envelope

type network struct {
	ctx        context.Context
	storeMaker storeMaker
	nodes      []*testNode
}

func NewNetwork(t *testing.T, count int, opts ...networkOption) *network {
	t.Helper()
	n := &network{
		ctx:        test.NewContext(t),
		storeMaker: func(t testing.TB, ctx context.Context) node.NodeStore { return memstore.NewNodeStore(ctx) },
	}
	for _, opt := range opts {
		opt(n)
	}
	nodes := make([]*testNode, count)
	for i := 0; i < count; i++ {
		name := fmt.Sprintf("node%d", i+1)
		nodes[i] = NewNode(t,
			WithContext(n.ctx),
			WithName(name),
			WithStore(n.storeMaker(t, context.WithLogger(n.ctx, n.ctx.Logger().Named(name)))))
	}

	var wg sync.WaitGroup
	for i, a := range nodes {
		for _, b := range nodes[i:] {
			if a == b {
				continue
			}
			wg.Add(1)
			go func(a, b *testNode) {
				defer wg.Done()
				a.Connect(t, b)
			}(a, b)
		}
	}
	wg.Wait()
	n.nodes = nodes
	return n
}

func (net *network) Close() error {
	for _, node := range net.nodes {
		node.Close()
	}
	return nil
}

func (net *network) PublishRandom(t *testing.T, topic string, count int) []*proto.Envelope {
	t.Helper()
	node := net.nodes[rand.Intn(len(net.nodes))]
	return node.PublishRandom(t, topic, count)
}

func (net *network) RequireEventuallyStoredEvents(t *testing.T, topic string, expected []*proto.Envelope) {
	for _, node := range net.nodes {
		node.RequireEventuallyStoredEvents(t, topic, expected)
	}
}

type networkSubscriber struct {
	Topic string
	subs  []*testSubscriber
}

func (net *network) Subscribe(t *testing.T, topic string) *networkSubscriber {
	t.Helper()
	subs := make([]*testSubscriber, len(net.nodes))
	for i, node := range net.nodes {
		subs[i] = node.Subscribe(t, topic)
	}
	return &networkSubscriber{
		Topic: topic,
		subs:  subs,
	}
}

func (s *networkSubscriber) RequireEventuallyCapturedEvents(t *testing.T, expected []*proto.Envelope) {
	t.Helper()
	for _, sub := range s.subs {
		sub.RequireEventuallyCapturedEvents(t, expected)
	}
}

// map of missing envelopes by topic for each node
type missingEnvs []map[string][]int

func (me missingEnvs) String() string {
	var buf bytes.Buffer
	for n, node := range me {
		for topic, envs := range node {
			if len(envs) > 0 {
				fmt.Fprintf(&buf, "n%d/%s: %v\n", n, topic, envs)
			}
		}
	}
	return buf.String()
}

type convergenceTracker struct {
	net       *network
	envelopes envsByTopic
	envCount  int
}

func newConvergenceTracker(net *network) *convergenceTracker {
	return &convergenceTracker{
		net:       net,
		envelopes: make(envsByTopic),
	}
}

func (tr *convergenceTracker) Publish(t *testing.T, node int, topic, msg string) {
	t.Helper()
	n := tr.net.nodes[node]
	assert.NotNil(t, n)
	tr.envCount++
	env := &proto.Envelope{
		TimestampNs:  uint64(tr.envCount),
		ContentTopic: topic,
		Message:      []byte(msg),
	}
	_, err := n.Publish(tr.net.ctx, &proto.PublishRequest{Envelopes: []*proto.Envelope{env}})
	assert.NoError(t, err)
	tr.envelopes[topic] = append(tr.envelopes[topic], env)
}

func (tr *convergenceTracker) newMissingEnvs() missingEnvs {
	nodes := make(missingEnvs, 0, len(tr.net.nodes))
	for range tr.net.nodes {
		missing := make(map[string][]int)
		for topic, envs := range tr.envelopes {
			ids := make([]int, 0, len(envs))
			for _, env := range envs {
				ids = append(ids, int(env.TimestampNs))
			}
			missing[topic] = ids
		}
		nodes = append(nodes, missing)
	}
	return nodes
}

// Wait for all the network nodes to converge on the captured set of events.
func (tr *convergenceTracker) RequireEventuallyComplete(t *testing.T, timeout time.Duration) {
	t.Helper()
	missing := tr.newMissingEnvs()
	timer := time.NewTimer(timeout)
	defer timer.Stop()
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-timer.C:
			missing, _ := tr.checkEvents(t, missing)
			if missing != nil {
				t.Errorf("missing events:\n%s", missing)
			}
			t.Log("converged")
			return
		case <-ticker.C:
			var progress bool
			missing, progress = tr.checkEvents(t, missing)
			if missing == nil {
				t.Log("converged")
				return
			}
			if !progress {
				t.Errorf("progress stopped:\n%s", missing)
				return
			}
			t.Logf("progress made:\n%s", missing)
		}
	}
}

// Check the remaining missing envelopes across all nodes.
// Return updated missing envelopes map.
// Return nil if nothing is missing.
func (tr *convergenceTracker) checkEvents(t *testing.T, missing missingEnvs) (remaining missingEnvs, progress bool) {
	anyMissing, progress := false, false
	for ni, nodeMissing := range missing {
		node := tr.net.nodes[ni]
		for topic, topicMissing := range nodeMissing {
			topicAll := tr.envelopes[topic]
			topicPresent := node.RequireQuery(t, topic)
			if len(topicAll) == len(topicPresent) {
				progress = true
				delete(nodeMissing, topic)
				continue
			}
			anyMissing = true
			topicRemaining := subtractEnvs(topicAll, topicPresent)
			if len(topicRemaining) < len(topicMissing) {
				progress = true
			}
			nodeMissing[topic] = topicRemaining
		}
	}
	if anyMissing {
		return missing, progress
	}
	return nil, true
}

func subtractEnvs(a, b []*proto.Envelope) []int {
	remaining := make([]int, 0)
OUTER:
	for _, env := range a {
		for _, env2 := range b {
			if env.TimestampNs == env2.TimestampNs {
				continue OUTER
			}
		}
		remaining = append(remaining, int(env.TimestampNs))
	}
	return remaining
}

func RandomNodeAndTopicSpraying(t *testing.T, nodes, topics, messages int, opts ...networkOption) {
	// to emulate significant concurrent activity we want nodes to be adding
	// events concurrently, but we also want to allow propagation at the same time.
	// So we need to introduce short delays to allow the network make some propagation progress.
	// Given the random spraying approach injecting a delay at every (nodes*topics)th event
	// should allow most nodes inject an event to most topics, and then the random length of the delay
	// should allow some amount of propagation to happen before the next burst.
	delayEvery := nodes * topics
	net := NewNetwork(t, nodes, opts...)
	defer net.Close()
	tracker := newConvergenceTracker(net)
	for i := 0; i < messages; i++ {
		topic := fmt.Sprintf("t%d", rand.Intn(topics))
		msg := fmt.Sprintf("gm %d", i)
		tracker.Publish(t, rand.Intn(nodes), topic, msg)
		if i%delayEvery == 0 {
			time.Sleep(time.Duration(rand.Intn(100)) * time.Microsecond)
		}
	}
	tracker.RequireEventuallyComplete(t, time.Duration(nodes*messages/100)*time.Second)
}
