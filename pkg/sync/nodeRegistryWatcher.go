package sync

import (
	"context"
	"fmt"
	"sync"

	"github.com/xmtp/xmtpd/pkg/registry"
	"github.com/xmtp/xmtpd/pkg/tracing"
	"go.uber.org/zap"
)

type NodeRegistryWatcher struct {
	ctx          context.Context
	log          *zap.Logger
	wg           sync.WaitGroup
	nodeid       uint32
	cancelFn     func()
	fnLock       sync.Mutex
	nodeRegistry registry.NodeRegistry
}

func NewNodeRegistryWatcher(ctx context.Context,
	log *zap.Logger, nodeId uint32, nodeRegistry registry.NodeRegistry,
) *NodeRegistryWatcher {
	return &NodeRegistryWatcher{
		ctx:          ctx,
		log:          log,
		nodeid:       nodeId,
		nodeRegistry: nodeRegistry,
		cancelFn:     nil,
	}
}

func (w *NodeRegistryWatcher) RegisterCancelFunction(fn func()) {
	w.fnLock.Lock()
	defer w.fnLock.Unlock()
	w.cancelFn = fn
}

func (w *NodeRegistryWatcher) triggerCancel() {
	if w.cancelFn != nil {
		w.fnLock.Lock()
		defer w.fnLock.Unlock()
		w.cancelFn()
		w.cancelFn = nil
	}
}

func (w *NodeRegistryWatcher) Watch() {
	registryChan := w.nodeRegistry.OnChangedNode(w.nodeid)

	tracing.GoPanicWrap(
		w.ctx,
		&w.wg,
		fmt.Sprintf("node-subscribe-%d-notifier", w.nodeid),
		func(ctx context.Context) {
			for {
				select {
				case <-ctx.Done():
					// this indicates that the node is shutting down
					w.triggerCancel()
					return
				case _, ok := <-registryChan:
					// this indicates that the registry has changed, and we need to rebuild the connection
					w.log.Info(
						"Node has been updated in the registry, terminating and rebuilding...",
					)
					w.triggerCancel()

					if !ok {
						w.log.Info("Node registry channel closed")
						return
					}
				}
			}
		},
	)
}
