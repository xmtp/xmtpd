package networkwatcher

import (
	"context"
	"errors"
	"sync"
	"time"

	"connectrpc.com/connect"
	"go.uber.org/zap"

	"github.com/xmtp/xmtpd/pkg/registry"
)

// WatcherConfig configures a network-watcher Watcher.
type WatcherConfig struct {
	Registry   registry.NodeRegistry
	Logger     *zap.Logger
	HTTPClient connect.HTTPClient

	MinBackoff time.Duration
	MaxBackoff time.Duration
}

// Watcher orchestrates per-node Subscribers driven by the on-chain
// registry. It owns an Aggregator and spawns Subscribers for each node
// the registry reports. The registry is add-only (nodes are never
// removed), so this watcher is add-only too — Subscribers are torn down
// only when the Watcher stops.
type Watcher struct {
	cfg        WatcherConfig
	aggregator *Aggregator

	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup

	mu      sync.Mutex
	spawned map[uint32]struct{}
}

// NewWatcher validates cfg and returns a Watcher.
func NewWatcher(cfg WatcherConfig) (*Watcher, error) {
	if cfg.Registry == nil {
		return nil, errors.New("networkwatcher: Registry is required")
	}
	if cfg.Logger == nil {
		cfg.Logger = zap.NewNop()
	}
	if cfg.MinBackoff <= 0 {
		cfg.MinBackoff = time.Second
	}
	if cfg.MaxBackoff <= 0 {
		cfg.MaxBackoff = 30 * time.Second
	}
	if cfg.MaxBackoff < cfg.MinBackoff {
		cfg.MaxBackoff = cfg.MinBackoff
	}
	return &Watcher{
		cfg:        cfg,
		aggregator: NewAggregator(),
		spawned:    make(map[uint32]struct{}),
	}, nil
}

// Start begins watching the registry and spawning Subscribers. It returns
// after initial node reconciliation; node-add events continue to be
// processed in the background until Stop (or ctx cancel).
func (w *Watcher) Start(ctx context.Context) error {
	w.ctx, w.cancel = context.WithCancel(ctx)

	initial, err := w.cfg.Registry.GetNodes()
	if err != nil {
		registryErrors.Inc()
		return err
	}
	w.spawnMissing(initial)

	w.wg.Go(func() {
		w.watchLoop()
	})
	return nil
}

// Stop cancels all subscribers and waits for them to exit.
func (w *Watcher) Stop() {
	if w.cancel != nil {
		w.cancel()
	}
	w.wg.Wait()
}

func (w *Watcher) watchLoop() {
	ch := w.cfg.Registry.OnNewNodes()
	for {
		select {
		case <-w.ctx.Done():
			return
		case nodes, ok := <-ch:
			if !ok {
				return
			}
			w.spawnMissing(nodes)
		}
	}
}

// spawnMissing starts a Subscriber for every node in the slice we haven't
// already spawned for. The slice may be a delta (OnNewNodes) or the full
// initial set (GetNodes); both are handled the same — we only add.
func (w *Watcher) spawnMissing(nodes []registry.Node) {
	w.mu.Lock()
	defer w.mu.Unlock()

	for _, n := range nodes {
		if !n.IsValidConfig || n.HTTPAddress == "" {
			continue
		}
		id := n.NodeID
		if _, exists := w.spawned[id]; exists {
			continue
		}
		w.cfg.Logger.Info(
			"subscriber added",
			zap.Uint32("node_id", id),
			zap.String("url", n.HTTPAddress),
		)
		sub := NewSubscriber(SubscriberConfig{
			NodeID:     id,
			BaseURL:    n.HTTPAddress,
			Aggregator: w.aggregator,
			Logger:     w.cfg.Logger.With(zap.Uint32("node_id", id)),
			HTTPClient: w.cfg.HTTPClient,
			MinBackoff: w.cfg.MinBackoff,
			MaxBackoff: w.cfg.MaxBackoff,
		})
		w.spawned[id] = struct{}{}
		w.wg.Go(func() { sub.Run(w.ctx) })
	}

	knownNodes.Set(float64(len(w.spawned)))
}
