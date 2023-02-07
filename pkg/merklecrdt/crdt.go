package merklecrdt

import (
	"context"

	mh "github.com/multiformats/go-multihash"
	"github.com/xmtp/xmtpd/pkg/merklecrdt/types"
	"github.com/xmtp/xmtpd/pkg/zap"
)

// MerkleCRDT manages the DAG of a dataset replica.
type MerkleCRDT struct {
	log       *zap.Logger
	ctx       context.Context
	ctxCancel context.CancelFunc

	store       Store
	broadcaster Broadcaster
	syncer      Syncer

	pendingReceiveEvents chan *types.Event // broadcasted events that were received from the network but not processed yet
	pendingSyncEvents    chan *types.Event // missing events that were fetched from the network but not processed yet
	pendingLinks         chan mh.Multihash // missing links that were discovered but not successfully fetched yet
}

func New(ctx context.Context, log *zap.Logger, store Store, bc Broadcaster, syncer Syncer) (*MerkleCRDT, error) {
	ctx, ctxCancel := context.WithCancel(ctx)
	m := &MerkleCRDT{
		log:       log,
		ctx:       ctx,
		ctxCancel: ctxCancel,

		store:       store,
		broadcaster: bc,
		syncer:      syncer,

		// TODO: tuning the channel sizes will likely be important
		// current implementation can lock up if the channels fill up.
		pendingReceiveEvents: make(chan *types.Event, 20),
		pendingSyncEvents:    make(chan *types.Event, 20),
		pendingLinks:         make(chan mh.Multihash, 20),
	}

	go m.receiveEventLoop(ctx)
	go m.syncEventLoop(ctx)
	go m.syncLinkLoop(ctx)
	go m.nextBroadcastedEventLoop(ctx)

	err := m.bootstrap(ctx)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (m *MerkleCRDT) Close() error {
	if m.ctxCancel != nil {
		m.ctxCancel()
	}
	return nil
}

func (m *MerkleCRDT) nextBroadcastedEventLoop(ctx context.Context) {
	for {
		ev, err := m.broadcaster.Next(ctx)
		if err != nil {
			if err == context.Canceled {
				return
			}
			m.log.Error("error getting next broadcasted event", zap.Error(err))
			return
		}
		m.log.Debug("received broadcasted event", zap.Cid("event_cid", ev.Cid))
		m.pendingReceiveEvents <- ev
	}
}

// receiveEventLoop processes incoming Events from broadcasts.
// It consumes pendingReceiveEvents and writes into pendingLinks.
func (m *MerkleCRDT) receiveEventLoop(ctx context.Context) {
loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		case ev := <-m.pendingReceiveEvents:
			// m.log.Debug("adding event", zap.Cid("event", ev.cid))
			added, err := m.store.AddHead(ev)
			if err != nil {
				// requeue for later
				// TODO: may need a delay
				// TODO: if the channel is full, this will lock up the loop
				m.pendingReceiveEvents <- ev
			}
			if added {
				for _, link := range ev.Links {
					m.pendingLinks <- link
				}
			}
		}
	}
}

// syncLoop fetches missing events from links.
// It consumes pendingLinks and writes into pendingSyncEvents
func (m *MerkleCRDT) syncLinkLoop(ctx context.Context) {
loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		case cid := <-m.pendingLinks:
			// m.log.Debug("checking link", zap.Cid("link", cid))
			// If the CID is in heads, it should be removed because
			// we have an event that points to it.
			// We also don't need to fetch it since we already have it.
			haveAlready, err := m.store.RemoveHead(cid)
			if err != nil {
				// requeue for later
				// TODO: may need a delay
				// TODO: if the channel is full, this will lock up the loop
				m.pendingLinks <- cid
				continue
			}
			if haveAlready {
				continue
			}
			m.log.Debug("fetching link", zap.Cid("link", cid))
			cids := []mh.Multihash{cid}
			evs, err := m.syncer.Fetch(cids)
			if err != nil {
				// requeue for later
				// TODO: this will need refinement for invalid, missing cids etc.
				// TODO: if the channel is full, this will lock up the loop
				m.pendingLinks <- cid
			}
			for i, ev := range evs {
				if ev == nil {
					// requeue missing links
					m.pendingLinks <- cids[i]
					continue
				}
				m.pendingSyncEvents <- ev
			}
		}
	}
}

// syncEventLoop processes missing events that were fetched from links.
// It consumes pendingSyncEvents and writes into pendingLinks.
// TODO: There is channel read/write cycle between the two sync loops,
// i.e. they could potentially lock up if both channels fill up.
func (m *MerkleCRDT) syncEventLoop(ctx context.Context) {
loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		case ev := <-m.pendingSyncEvents:
			// m.log.Debug("adding link event", zap.Cid("event", ev.cid))
			added, err := m.store.AddEvent(ev)
			if err != nil {
				// requeue for later
				// TODO: may need a delay
				// TODO: if the channel is full, this will lock up the loop
				m.pendingSyncEvents <- ev
			}
			if added {
				for _, link := range ev.Links {
					// TODO: if the channel is full, this will lock up the loop
					m.pendingLinks <- link
				}
			}
		}
	}
}

// Bootstrap from the contents of the store.
func (m *MerkleCRDT) bootstrap(ctx context.Context) error {
	links, err := m.store.FindMissingLinks()
	if err != nil {
		return err
	}
	for _, link := range links {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case m.pendingLinks <- link:
		}
	}
	return nil
}
