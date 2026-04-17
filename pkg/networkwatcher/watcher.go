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
// registry. It owns an Aggregator and spawns/cancels Subscribers in
// response to registry change events.
type Watcher struct {
	cfg        WatcherConfig
	aggregator *Aggregator

	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup

	mu      sync.Mutex
	handles map[uint32]*subHandle
}

type subHandle struct {
	cancel context.CancelFunc
	done   chan struct{}
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
	if cfg.MaxBackoff < cfg.MinBackoff {
		cfg.MaxBackoff = 30 * time.Second
	}
	return &Watcher{
		cfg:        cfg,
		aggregator: NewAggregator(),
		handles:    make(map[uint32]*subHandle),
	}, nil
}

// Start begins watching the registry and spawning Subscribers. It returns
// after initial node reconciliation; node updates continue in the background
// until Stop (or ctx cancel).
func (w *Watcher) Start(ctx context.Context) error {
	w.ctx, w.cancel = context.WithCancel(ctx)

	initial, err := w.cfg.Registry.GetNodes()
	if err != nil {
		registryErrors.Inc()
		return err
	}
	w.reconcile(initial)

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

	w.mu.Lock()
	handles := w.handles
	w.handles = map[uint32]*subHandle{}
	w.mu.Unlock()

	for _, h := range handles {
		h.cancel()
		<-h.done
	}
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
			w.reconcile(nodes)
		}
	}
}

func (w *Watcher) reconcile(nodes []registry.Node) {
	wanted := make(map[uint32]registry.Node, len(nodes))
	for _, n := range nodes {
		if !n.IsValidConfig || n.HTTPAddress == "" {
			continue
		}
		wanted[n.NodeID] = n
	}

	w.mu.Lock()
	defer w.mu.Unlock()

	// Cancel subscribers for removed nodes.
	for id, h := range w.handles {
		if _, keep := wanted[id]; !keep {
			w.cfg.Logger.Info("subscriber removed", zap.Uint32("node_id", id))
			h.cancel()
			delete(w.handles, id)
		}
	}

	// Spawn subscribers for new nodes.
	for id, n := range wanted {
		if _, exists := w.handles[id]; exists {
			continue
		}
		w.cfg.Logger.Info(
			"subscriber added",
			zap.Uint32("node_id", id),
			zap.String("url", n.HTTPAddress),
		)
		subCtx, subCancel := context.WithCancel(w.ctx)
		done := make(chan struct{})
		sub := NewSubscriber(SubscriberConfig{
			NodeID:     id,
			BaseURL:    n.HTTPAddress,
			Aggregator: w.aggregator,
			Logger:     w.cfg.Logger.With(zap.Uint32("node_id", id)),
			HTTPClient: w.cfg.HTTPClient,
			MinBackoff: w.cfg.MinBackoff,
			MaxBackoff: w.cfg.MaxBackoff,
		})
		go func() {
			defer close(done)
			sub.Run(subCtx)
		}()
		w.handles[id] = &subHandle{cancel: subCancel, done: done}
	}

	knownNodes.Set(float64(len(wanted)))
}
