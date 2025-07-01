package sync

import (
	"context"
	"fmt"
	"sync"

	"github.com/xmtp/xmtpd/pkg/registry"
	"github.com/xmtp/xmtpd/pkg/tracing"
	"go.uber.org/zap"
)

// NodeRegistryWatcher monitors registry updates for a specific node and invokes
// a user-provided cancel function when a change is detected.
//
// The cancel function is intended to trigger teardown or reinitialization of
// any state associated with that node (e.g., gRPC connections, sync loops).
//
// 📌 Usage:
//
//	watcher := NewNodeRegistryWatcher(ctx, log, nodeID, registry)
//	watcher.RegisterCancelFunction(func() {
//	    teardownOrRestart()
//	})
//	watcher.Watch()
//
// When the registry signals a change, the cancel function is invoked **exactly once**.
// If the caller wishes to handle **subsequent changes**, they must call
// `RegisterCancelFunction` again after each invocation.
//
// If no cancel function is registered when a change occurs, the update is silently ignored.
type NodeRegistryWatcher struct {
	ctx          context.Context
	log          *zap.Logger
	wg           sync.WaitGroup
	nodeid       uint32
	cancelFn     func()
	fnLock       sync.Mutex
	nodeRegistry registry.NodeRegistry
}

// NewNodeRegistryWatcher creates a new watcher tied to the provided node ID.
// It does not begin watching until Watch() is explicitly called.
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

// RegisterCancelFunction registers a cancel function to be invoked on the next
// registry update or shutdown event.
//
// ❗️Important: This function is called **once**, and then cleared. To respond
// to subsequent changes, you must register a new cancel function after each call.
func (w *NodeRegistryWatcher) RegisterCancelFunction(fn func()) {
	w.fnLock.Lock()
	defer w.fnLock.Unlock()
	w.cancelFn = fn
}

func (w *NodeRegistryWatcher) triggerCancel() {
	w.fnLock.Lock()
	defer w.fnLock.Unlock()
	if w.cancelFn != nil {
		w.cancelFn()
		w.cancelFn = nil
	}
}

// Watch starts a background goroutine to listen for registry changes for the node.
// It invokes the registered cancel function (if present) once per change.
//
// The watcher exits when:
//   - The parent context is cancelled
//   - The registry channel is closed
//
// ⚠️ Note: If no cancel function is registered when a change is detected,
// the update is **ignored**.
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
