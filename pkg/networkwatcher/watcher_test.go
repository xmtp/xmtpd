package networkwatcher

import (
	"context"
	"crypto/ecdsa"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"connectrpc.com/connect"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/envelopes"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/metadata_api"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/metadata_api/metadata_apiconnect"
	"github.com/xmtp/xmtpd/pkg/registry"
)

// fakeRegistry implements registry.NodeRegistry for tests.
type fakeRegistry struct {
	mu           sync.Mutex
	nodes        []registry.Node
	newNodesCh   chan []registry.Node
	changedChans map[uint32]chan registry.Node
}

func newFakeRegistry(initial []registry.Node) *fakeRegistry {
	return &fakeRegistry{
		nodes:        initial,
		newNodesCh:   make(chan []registry.Node, 8),
		changedChans: map[uint32]chan registry.Node{},
	}
}

func (f *fakeRegistry) GetNodes() ([]registry.Node, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	out := make([]registry.Node, len(f.nodes))
	copy(out, f.nodes)
	return out, nil
}

func (f *fakeRegistry) GetNode(id uint32) (*registry.Node, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	for i := range f.nodes {
		if f.nodes[i].NodeID == id {
			n := f.nodes[i]
			return &n, nil
		}
	}
	return nil, nil
}

func (f *fakeRegistry) OnNewNodes() <-chan []registry.Node { return f.newNodesCh }

func (f *fakeRegistry) OnChangedNode(id uint32) <-chan registry.Node {
	f.mu.Lock()
	defer f.mu.Unlock()
	ch, ok := f.changedChans[id]
	if !ok {
		ch = make(chan registry.Node, 1)
		f.changedChans[id] = ch
	}
	return ch
}

func (f *fakeRegistry) Stop() {}

func (f *fakeRegistry) pushNodes(nodes []registry.Node) {
	f.mu.Lock()
	f.nodes = nodes
	f.mu.Unlock()
	f.newNodesCh <- nodes
}

// stubHandler responds to SubscribeSyncCursor with a single snapshot and
// keeps the stream open until ctx is done.
type stubHandler struct {
	metadata_apiconnect.UnimplementedMetadataApiHandler
	cursor map[uint32]uint64
}

func (s *stubHandler) SubscribeSyncCursor(
	ctx context.Context,
	_ *connect.Request[metadata_api.GetSyncCursorRequest],
	stream *connect.ServerStream[metadata_api.GetSyncCursorResponse],
) error {
	_ = stream.Send(&metadata_api.GetSyncCursorResponse{
		LatestSync: &envelopes.Cursor{NodeIdToSequenceId: s.cursor},
	})
	<-ctx.Done()
	return nil
}

func makeStubNode(t *testing.T, id uint32, cursor map[uint32]uint64) registry.Node {
	t.Helper()
	mux := http.NewServeMux()
	path, h := metadata_apiconnect.NewMetadataApiHandler(&stubHandler{cursor: cursor})
	mux.Handle(path, h)
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)
	return registry.Node{
		NodeID:        id,
		HTTPAddress:   srv.URL,
		IsCanonical:   true,
		IsValidConfig: true,
		SigningKey:    &ecdsa.PublicKey{},
	}
}

func TestWatcher_SpawnsSubscribers_ForInitialNodes(t *testing.T) {
	resetAggregatorMetrics()
	knownNodes.Set(0)

	n1 := makeStubNode(t, 1, map[uint32]uint64{100: 42})
	n2 := makeStubNode(t, 2, map[uint32]uint64{100: 99})

	reg := newFakeRegistry([]registry.Node{n1, n2})

	w, err := NewWatcher(WatcherConfig{
		Registry:   reg,
		Logger:     zap.NewNop(),
		MinBackoff: 5 * time.Millisecond,
		MaxBackoff: 20 * time.Millisecond,
		HTTPClient: http.DefaultClient,
	})
	require.NoError(t, err)

	ctx := t.Context()
	require.NoError(t, w.Start(ctx))
	defer w.Stop()

	require.Eventually(t, func() bool {
		return metricValue(t, cursorGauge.WithLabelValues("1", "100")) == 42 &&
			metricValue(t, cursorGauge.WithLabelValues("2", "100")) == 99 &&
			metricValue(t, knownNodes) == 2
	}, 2*time.Second, 10*time.Millisecond)
}

func TestWatcher_AddsSubscriber_OnNewNodeEvent(t *testing.T) {
	resetAggregatorMetrics()
	knownNodes.Set(0)

	n1 := makeStubNode(t, 1, map[uint32]uint64{100: 1})
	reg := newFakeRegistry([]registry.Node{n1})

	w, err := NewWatcher(WatcherConfig{
		Registry:   reg,
		Logger:     zap.NewNop(),
		MinBackoff: 5 * time.Millisecond,
		MaxBackoff: 20 * time.Millisecond,
		HTTPClient: http.DefaultClient,
	})
	require.NoError(t, err)

	ctx := t.Context()
	require.NoError(t, w.Start(ctx))
	defer w.Stop()

	require.Eventually(t, func() bool {
		return metricValue(t, nodeUp.WithLabelValues("1")) == 1
	}, 2*time.Second, 10*time.Millisecond)

	// OnNewNodes emits the delta of newly-added nodes — the watcher
	// should spawn for n2 without disturbing the already-running n1.
	n2 := makeStubNode(t, 2, map[uint32]uint64{100: 2})
	reg.pushNodes([]registry.Node{n2})

	require.Eventually(t, func() bool {
		return metricValue(t, nodeUp.WithLabelValues("2")) == 1 &&
			metricValue(t, nodeUp.WithLabelValues("1")) == 1 &&
			metricValue(t, knownNodes) == 2
	}, 2*time.Second, 10*time.Millisecond)
}
