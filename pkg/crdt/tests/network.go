package tests

import (
	"context"
	"fmt"
	"io"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	messagev1 "github.com/xmtp/proto/v3/go/message_api/v1"
	"github.com/xmtp/xmtpd/pkg/crdt"
	test "github.com/xmtp/xmtpd/pkg/testing"
	"github.com/xmtp/xmtpd/pkg/zap"
)

type storeMaker func(*zap.Logger) crdt.NodeStore
type syncerMaker func(*zap.Logger) crdt.NodeSyncer
type broadcasterMaker func(*zap.Logger) crdt.NodeBroadcaster
type networkConfig struct {
	storeMaker
	syncerMaker
	broadcasterMaker
}

func defaultNetworkConfig() *networkConfig {
	var cfg networkConfig
	(&cfg).modify(
		WithStore(NewMapStore),
		WithSyncer(NewRandomSyncer),
		WithBroadcaster(NewChanBroadcaster),
	)
	return &cfg
}

func (cfg *networkConfig) modify(modifiers ...configModifier) {
	for _, m := range modifiers {
		m(cfg)
	}
}

type configModifier func(*networkConfig)

func WithStore(m storeMaker) configModifier {
	return func(cfg *networkConfig) {
		cfg.storeMaker = m
	}
}

func WithBroadcaster(m broadcasterMaker) configModifier {
	return func(cfg *networkConfig) {
		cfg.broadcasterMaker = m
	}
}

func WithSyncer(m syncerMaker) configModifier {
	return func(cfg *networkConfig) {
		cfg.syncerMaker = m
	}
}

// network is an in-memory simulation of a network of a given number of Nodes.
// network also captures events that were published to it for final analysis of the test results.
type network struct {
	ctx    context.Context
	cancel context.CancelFunc
	log    *zap.Logger
	bc     crdt.NodeBroadcaster
	sync   crdt.NodeSyncer
	count  int
	nodes  map[int]*crdt.Node // all the nodes in the network
	events []*crdt.Event      // captures events published to the network
}

const t0 = "t0" // first topic
// const t1 = "t1" // second topic
// const t2 = "t2" // third topic
// ...

// Creates a network with given number of nodes
func NewNetwork(t *testing.T, nodes, topics int, modifiers ...configModifier) *network {
	ctx, cancel := context.WithCancel(context.Background())
	log := test.NewLogger(t)
	cfg := defaultNetworkConfig()
	net := &network{
		ctx:    ctx,
		cancel: cancel,
		log:    log,
		bc:     cfg.broadcasterMaker(log),
		sync:   cfg.syncerMaker(log),
		nodes:  make(map[int]*crdt.Node),
	}
	for i := 0; i < nodes; i++ {
		net.AddNode(t, cfg.storeMaker(log))
	}
	require.Len(t, net.nodes, nodes)
	return net
}

func (net *network) Close() {
	net.cancel()
}

func (net *network) AddNode(t *testing.T, store crdt.NodeStore) *crdt.Node {
	name := fmt.Sprintf("n%d", net.count)
	n, err := crdt.NewNode(net.ctx,
		net.log.Named(name),
		store,
		net.sync,
		net.bc)
	assert.NoError(t, err)
	net.nodes[net.count] = n
	net.count += 1
	return n
}

func (net *network) RemoveNode(t *testing.T, n int) *crdt.Node {
	node := net.nodes[n]
	assert.NotNil(t, node)
	delete(net.nodes, n)
	node.Close()
	return node
}

// Publishes msg into a topic from given node
func (net *network) Publish(t *testing.T, node int, topic, msg string) {
	t.Helper()
	n := net.nodes[node]
	assert.NotNil(t, n)
	ev, err := n.Publish(net.ctx, &messagev1.Envelope{TimestampNs: uint64(len(net.events) + 1), ContentTopic: topic, Message: []byte(msg)})
	assert.NoError(t, err)
	net.events = append(net.events, ev)
}

func (net *network) Query(t *testing.T, node int, topic string, modifiers ...queryModifier) ([]*messagev1.Envelope, *messagev1.PagingInfo, error) {
	t.Helper()
	n := net.nodes[node]
	assert.NotNil(t, n)
	q := &messagev1.QueryRequest{
		ContentTopics: []string{topic},
	}
	for _, m := range modifiers {
		m(q)
	}
	return n.Query(net.ctx, q)
}

// Suspends topic broadcast delivery to the given node while fn runs
func (net *network) WithSuspendedTopic(t *testing.T, node int, topic string, fn func(*crdt.Node)) {
	n := net.nodes[node]
	assert.NotNil(t, n)
	bc := n.NodeBroadcaster.(*chanBroadcaster)
	bc.RemoveNode(n)
	defer bc.AddNode(n)
	fn(n)
}

// Wait for all the network nodes to converge on the captured set of events.
func (net *network) AssertEventuallyConsistent(t *testing.T, timeout time.Duration, ignore ...int) {
	t.Helper()
	timer := time.NewTimer(timeout)
	defer timer.Stop()
	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-timer.C:
			missing := net.checkEvents(t, ignore)
			if len(missing) > 0 {
				t.Errorf("Missing events: %v", missing)
			}
			return
		case <-ticker.C:
			if len(net.checkEvents(t, ignore)) == 0 {
				return
			}
		}
	}
}

// Check that the resulting envelops match the expected list of timestamps.
// net.Publish generates timestamps incrementally starting from 1 so they are unique and match the publishing order.
func (net *network) AssertQueryResult(t *testing.T, envelopes []*messagev1.Envelope, expected ...int) {
	t.Helper()
	var result []int
	for _, env := range envelopes {
		result = append(result, int(env.TimestampNs))
	}
	assert.Equal(t, expected, result, "timestamps")
}

// Check that the cursor timestamp matches the expected timestamp
// net.Publish generates timestamps incrementally starting from 1 so they are unique and match the publishing order.
func (net *network) AssertQueryCursor(t *testing.T, expected int, cursor *messagev1.Cursor) {
	t.Helper()
	require.NotNil(t, cursor, "cursor")
	actual := int(cursor.GetIndex().SenderTimeNs)
	assert.Equal(t, expected, actual, "timestamp")
}

// Check that all nodes except the ignored ones have all events.
// Returns map of nodes that have missing events,
// the key is the node number
// the value is a string listing present events by number and _ for missing events.
func (net *network) checkEvents(t *testing.T, ignore []int) (missing map[int]string) {
	missing = make(map[int]string)
	for j, n := range net.nodes {
		if ignored(j, ignore) {
			continue
		}
		count, err := n.Count()
		assert.NoError(t, err)
		if count == len(net.events) {
			continue // shortcut
		}
		result := ""
		pass := true
		for i, ev := range net.events {
			ev2, err := n.Get(ev.ContentTopic, ev.Cid)
			if err != nil || ev2 == nil {
				result = result + "_"
				pass = false
			} else {
				result = result + strconv.FormatInt(int64(i), 36)
			}
		}
		assert.False(t, pass)
		missing[j] = result
	}
	return missing
}

// emit a graphvis depiction of the topic contents
// showing the individual events and their links
func (net *network) visualiseTopic(w io.Writer, topic string) {
	fmt.Fprintf(w, "strict digraph %s {\n", topic)
	for i := len(net.events) - 1; i >= 0; i-- {
		ev := net.events[i]
		if ev.ContentTopic != topic {
			continue
		}
		fmt.Fprintf(w, "\t\"%s\" [label=\"%d: \\N\"]\n", zap.ShortCid(ev.Cid), i)
		fmt.Fprintf(w, "\t\"%s\" -> { ", zap.ShortCid(ev.Cid))
		for _, l := range ev.Links {
			fmt.Fprintf(w, "\"%s\" ", zap.ShortCid(l))
		}
		fmt.Fprintf(w, "}\n")
	}
	fmt.Fprintf(w, "}\n")
}

func ignored(i int, ignore []int) bool {
	for _, j := range ignore {
		if i == j {
			return true
		}
	}
	return false
}
