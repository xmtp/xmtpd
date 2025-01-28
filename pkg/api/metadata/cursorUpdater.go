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
	cu := CursorUpdater{ctx: ctx, log: log, store: store, subscribers: subscribers}

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
			err := cu.read()
			if err != nil {
				//TODO proper error handling
				return
			}
			cu.notifySubscribers()
		}
	}
}

func (cu *CursorUpdater) read() error {

	rows, err := queries.New(cu.store).GetLatestCursor(cu.ctx)
	if err != nil {
		return err
	}

	nodeIdToSequenceId := make(map[uint32]uint64)
	for _, row := range rows {
		nodeIdToSequenceId[uint32(row.OriginatorNodeID)] = uint64(row.MaxSequenceID)
	}

	cu.cursorMu.Lock()
	defer cu.cursorMu.Unlock()

	cu.cursor = nodeIdToSequenceId

	return nil
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
