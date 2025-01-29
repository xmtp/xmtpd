package metadata

import (
	"context"
	"database/sql"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/envelopes"
	"go.uber.org/zap"
	"sync"
	"time"
)

type CursorUpdater struct {
	ctx           context.Context
	log           *zap.Logger
	store         *sql.DB
	cursorMu      sync.RWMutex
	cursor        map[uint32]uint64
	subscribersMu sync.RWMutex
	subscribers   map[string][]chan struct{}
}

func NewCursorUpdater(ctx context.Context, log *zap.Logger, store *sql.DB) *CursorUpdater {
	subscribers := make(map[string][]chan struct{})
	cu := CursorUpdater{
		ctx:         ctx,
		log:         log.Named("cursor-updater"),
		store:       store,
		subscribers: subscribers,
	}

	go cu.start()
	return &cu
}

func (cu *CursorUpdater) GetCursor() *envelopes.Cursor {
	cu.cursorMu.RLock()
	defer cu.cursorMu.RUnlock()
	return &envelopes.Cursor{NodeIdToSequenceId: cu.cursor}
}

func (cu *CursorUpdater) start() {
	ticker := time.NewTicker(100 * time.Millisecond) // Adjust the period as needed
	defer ticker.Stop()
	for {
		select {
		case <-cu.ctx.Done():
			return
		case <-ticker.C:
			updated, err := cu.read()
			if err != nil {
				//TODO proper error handling
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

func (cu *CursorUpdater) read() (bool, error) {
	rows, err := queries.New(cu.store).GetLatestCursor(cu.ctx)
	if err != nil {
		return false, err
	}

	nodeIdToSequenceId := make(map[uint32]uint64)
	for _, row := range rows {
		nodeIdToSequenceId[uint32(row.OriginatorNodeID)] = uint64(row.MaxSequenceID)
	}

	cu.cursorMu.Lock()
	defer cu.cursorMu.Unlock()

	if !equalCursors(cu.cursor, nodeIdToSequenceId) {
		cu.cursor = nodeIdToSequenceId
		return true, nil
	}

	return false, nil
}

func (cu *CursorUpdater) notifySubscribers() {
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

func (cu *CursorUpdater) AddSubscriber(clientID string, updateChan chan struct{}) {
	cu.subscribersMu.Lock()
	defer cu.subscribersMu.Unlock()
	cu.subscribers[clientID] = append(cu.subscribers[clientID], updateChan)
}

func (cu *CursorUpdater) RemoveSubscriber(clientID string) {
	cu.subscribersMu.Lock()
	defer cu.subscribersMu.Unlock()
	delete(cu.subscribers, clientID)
}
