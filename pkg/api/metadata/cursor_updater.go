package metadata

import (
	"context"
	"database/sql"
	"sync"
	"time"

	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/envelopes"
	"github.com/xmtp/xmtpd/pkg/tracing"
	"go.uber.org/zap"
)

type CursorUpdater interface {
	GetCursor() *envelopes.Cursor
	AddSubscriber(clientID string, updateChan chan struct{})
	RemoveSubscriber(clientID string)
	Stop()
}

type DBBasedCursorUpdater struct {
	ctx           context.Context
	store         *sql.DB
	cancel        context.CancelFunc
	wg            sync.WaitGroup
	cursorMu      sync.RWMutex
	cursor        map[uint32]uint64
	subscribersMu sync.RWMutex
	subscribers   map[string][]chan struct{}
}

func NewCursorUpdater(ctx context.Context, logger *zap.Logger, store *sql.DB) CursorUpdater {
	subscribers := make(map[string][]chan struct{})
	ctx, cancel := context.WithCancel(ctx)
	cu := DBBasedCursorUpdater{
		ctx:         ctx,
		store:       store,
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
	ticker := time.NewTicker(100 * time.Millisecond) // Adjust the period as needed
	defer ticker.Stop()
	for {
		select {
		case <-cu.ctx.Done():
			return
		case <-ticker.C:
			updated, err := cu.read()
			if err != nil {
				// TODO proper error handling
				return
			}
			if updated {
				cu.notifySubscribers()
			}
		}
	}
}

func equalCursors(a, b map[uint32]uint64) bool {
	if len(a) != len(b) {
		return false
	}
	for key, valA := range a {
		if valB, ok := b[key]; !ok || valA != valB {
			return false
		}
	}
	return true
}

func (cu *DBBasedCursorUpdater) read() (bool, error) {
	rows, err := queries.New(cu.store).SelectVectorClock(cu.ctx)
	if err != nil {
		return false, err
	}

	nodeIDToSequenceID := make(map[uint32]uint64)
	for _, row := range rows {
		nodeIDToSequenceID[uint32(row.OriginatorNodeID)] = uint64(row.OriginatorSequenceID)
	}

	cu.cursorMu.Lock()
	defer cu.cursorMu.Unlock()

	if !equalCursors(cu.cursor, nodeIDToSequenceID) {
		cu.cursor = nodeIDToSequenceID
		return true, nil
	}

	return false, nil
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
