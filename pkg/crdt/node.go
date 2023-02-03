package crdt

import (
	"context"
	"errors"
	"sync"

	mh "github.com/multiformats/go-multihash"
	messagev1 "github.com/xmtp/proto/v3/go/message_api/v1"
	"github.com/xmtp/xmtpd/pkg/zap"
	"golang.org/x/sync/errgroup"
)

var TODO = errors.New("Not Yet Implemented")
var ErrUnknownTopic = errors.New("Unknown Topic")

// Node represents a peer in the XMTP network.
// Node hosts a set of Topics and provides the required
// supporting facilities (store, syncer, broadcaster).
type Node struct {
	ctx    context.Context
	cancel context.CancelFunc
	log    *zap.Logger

	topicsLock sync.RWMutex
	topics     map[string]*Topic

	NodeStore
	NodeSyncer
	NodeBroadcaster
}

// NewNode creates a new network node.
func NewNode(ctx context.Context, log *zap.Logger, store NodeStore, syncer NodeSyncer, bc NodeBroadcaster) (*Node, error) {
	ctx, cancel := context.WithCancel(ctx)
	node := &Node{
		ctx:             ctx,
		cancel:          cancel,
		log:             log,
		topics:          make(map[string]*Topic),
		NodeStore:       store,
		NodeSyncer:      syncer,
		NodeBroadcaster: bc,
	}
	// Find pre-existing topics
	topics, err := store.Topics()
	if err != nil {
		return nil, err
	}
	// Bootstrap all the topics with some parallelization.
	grp, ctx := errgroup.WithContext(ctx)
	grp.SetLimit(1000) // up to 1000 topic bootstraps in parallel
	for _, name := range topics {
		topic := name
		grp.Go(func() (err error) {
			t := node.createTopic(topic)
			return t.bootstrap(ctx)
		})
	}
	// Do not return until all topics are bootstrapped successfully.
	// If any bootstrap fails, bail out.
	if err := grp.Wait(); err != nil {
		cancel()
		return nil, err
	}
	return node, nil
}

func (n *Node) Close() {
	n.cancel()
}

// Publish sends a new message out to the network.
func (n *Node) Publish(ctx context.Context, env *messagev1.Envelope) (*Event, error) {
	topic := n.getOrCreateTopic(env.ContentTopic)
	return topic.Publish(ctx, env)
}

func (n *Node) Query(ctx context.Context, req *messagev1.QueryRequest) ([]*messagev1.Envelope, *messagev1.PagingInfo, error) {
	if len(req.ContentTopics) != 1 {
		// Not supporting querying multiple topics
		return nil, nil, TODO
	}
	t := n.getTopic(req.ContentTopics[0])
	if t == nil {
		return nil, nil, ErrUnknownTopic
	}
	return t.Query(ctx, req)
}

// Get retrieves an Event for given Topic.
func (n *Node) Get(topic string, cid mh.Multihash) (*Event, error) {
	t := n.getTopic(topic)
	if t == nil {
		return nil, ErrUnknownTopic
	}
	return t.Get(cid)
}

// Count returns count of all events on the Node.
func (n *Node) Count() (count int, err error) {
	n.topicsLock.RLock()
	defer n.topicsLock.RUnlock()
	for _, t := range n.topics {
		tc, err := t.Count()
		if err != nil {
			return 0, err
		}
		count += tc
	}
	return count, nil
}

func (n *Node) getTopic(topic string) *Topic {
	n.topicsLock.RLock()
	defer n.topicsLock.RUnlock()
	return n.topics[topic]
}

func (n *Node) createTopic(topic string) *Topic {
	n.topicsLock.Lock()
	defer n.topicsLock.Unlock()
	return n.newTopic(topic)
}

// getOrCreateTopic MUST NOT be called before topic bootstrap is complete
// to avoid creating empty topics that weren't bootstrapped.
func (n *Node) getOrCreateTopic(topic string) *Topic {
	n.topicsLock.Lock()
	defer n.topicsLock.Unlock()
	t := n.topics[topic]
	if t == nil {
		t = n.newTopic(topic)
	}
	return t
}

// newTopic adds a topic to the Node.
// MUST be called with a write lock!
func (n *Node) newTopic(name string) *Topic {
	t := NewTopic(
		n.ctx,
		name,
		n.log.Named(name),
		n.NodeStore.NewTopic(name, n),
		n.NodeSyncer.NewTopic(name, n),
		n.NodeBroadcaster.NewTopic(name, n),
	)
	n.topics[name] = t
	return t
}
