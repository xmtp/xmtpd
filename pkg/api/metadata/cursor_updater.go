package metadata

import (
	"context"
	"maps"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/xmtp/xmtpd/pkg/db"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/envelopes"
	"github.com/xmtp/xmtpd/pkg/tracing"
)

type CursorUpdater interface {
	GetCursor() *envelopes.Cursor
	AddSubscriber(clientID string, updateChan chan struct{})
	RemoveSubscriber(clientID string)
	Stop()
}

type DBBasedCursorUpdater struct {
	ctx           context.Context
	vc            db.VectorClock
	cancel        context.CancelFunc
	wg            sync.WaitGroup
	cursorMu      sync.RWMutex
	cursor        map[uint32]uint64
	subscribersMu sync.RWMutex
	subscribers   map[string][]chan struct{}
}

func NewCursorUpdater(
	ctx context.Context,
	logger *zap.Logger,
	vc db.VectorClock,
) CursorUpdater {
	subscribers := make(map[string][]chan struct{})
	ctx, cancel := context.WithCancel(ctx)
	cu := DBBasedCursorUpdater{
		ctx:         ctx,
		vc:          vc,
		cancel:      cancel,
		wg:          sync.WaitGroup{},
		subscribers: subscribers,
	}

	tracing.GoPanicWrap(
		cu.ctx,
		&cu.wg,
		"cursor-updater",
		func(ctx context.Context) {
			cu.start()
		})
	return &cu
}

func (cu *DBBasedCursorUpdater) GetCursor() *envelopes.Cursor {
	cu.cursorMu.RLock()
	defer cu.cursorMu.RUnlock()
	return &envelopes.Cursor{NodeIdToSequenceId: cu.cursor}
}

func (cu *DBBasedCursorUpdater) start() {
	// TODO: Check - this can now be more frequent, if that is something we want.

	ticker := time.NewTicker(100 * time.Millisecond) // Adjust the period as needed
	defer ticker.Stop()
	for {
		select {
		case <-cu.ctx.Done():
			return
		case <-ticker.C:
			updated := cu.read()
			if updated {
				cu.notifySubscribers()
			}
		}
	}
}

func (cu *DBBasedCursorUpdater) read() bool {
	// Read current vector clock.
	current := cu.vc.Values()

	cu.cursorMu.Lock()
	defer cu.cursorMu.Unlock()

	if !maps.Equal(cu.cursor, current) {
		cu.cursor = current
		return true
	}

	return false
}

func (cu *DBBasedCursorUpdater) notifySubscribers() {
	cu.subscribersMu.Lock()
	defer cu.subscribersMu.Unlock()

	for _, channels := range cu.subscribers {
		for _, ch := range channels {
			select {
			case ch <- struct{}{}:
			default:
				// if the channel already has a notification pending, we don't need to add another one
				continue
			}
		}
	}
}

func (cu *DBBasedCursorUpdater) AddSubscriber(clientID string, updateChan chan struct{}) {
	cu.subscribersMu.Lock()
	defer cu.subscribersMu.Unlock()
	cu.subscribers[clientID] = append(cu.subscribers[clientID], updateChan)
}

func (cu *DBBasedCursorUpdater) RemoveSubscriber(clientID string) {
	cu.subscribersMu.Lock()
	defer cu.subscribersMu.Unlock()
	delete(cu.subscribers, clientID)
}

func (cu *DBBasedCursorUpdater) Stop() {
	cu.cancel()
	cu.wg.Wait()
}
