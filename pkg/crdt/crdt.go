package crdt

import (
	"context"

	mh "github.com/multiformats/go-multihash"
	"github.com/xmtp/xmtpd/pkg/crdt/types"
	"github.com/xmtp/xmtpd/pkg/zap"
)

type NewEventFunc func(ev *types.Event)

// CRDT manages the DAG of a dataset replica.
type CRDT struct {
	log        *zap.Logger
	ctx        context.Context
	ctxCancel  context.CancelFunc
	onNewEvent NewEventFunc

	store       Store
	broadcaster Broadcaster
	syncer      Syncer

	pendingReceiveEvents chan *types.Event // broadcasted events that were received from the network but not processed yet
	pendingSyncEvents    chan *types.Event // missing events that were fetched from the network but not processed yet
	pendingLinks         chan mh.Multihash // missing links that were discovered but not successfully fetched yet
}

func New(ctx context.Context, log *zap.Logger, store Store, bc Broadcaster, syncer Syncer, onNewEvent NewEventFunc) (*CRDT, error) {
	ctx, ctxCancel := context.WithCancel(ctx)
	c := &CRDT{
		log:        log,
		ctx:        ctx,
		ctxCancel:  ctxCancel,
		onNewEvent: onNewEvent,

		store:       store,
		broadcaster: bc,
		syncer:      syncer,

		// TODO: tuning the channel sizes will likely be important
		// current implementation can lock up if the channels fill up.
		pendingReceiveEvents: make(chan *types.Event, 20),
		pendingSyncEvents:    make(chan *types.Event, 20),
		pendingLinks:         make(chan mh.Multihash, 20),
	}

	go c.receiveEventLoop(ctx)
	go c.syncEventLoop(ctx)
	go c.syncLinkLoop(ctx)
	go c.nextBroadcastedEventLoop(ctx)

	err := c.bootstrap(ctx)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (c *CRDT) Close() error {
	if c.ctxCancel != nil {
		c.ctxCancel()
	}
	return nil
}

func (c *CRDT) Broadcast(ctx context.Context, payload []byte) error {
	ev, err := types.NewEvent(payload, nil)
	if err != nil {
		return err
	}
	return c.broadcaster.Broadcast(ev)
}

func (c *CRDT) nextBroadcastedEventLoop(ctx context.Context) {
	for {
		ev, err := c.broadcaster.Next(ctx)
		if err != nil {
			if err == context.Canceled {
				return
			}
			c.log.Error("error getting next broadcasted event", zap.Error(err))
			return
		}
		c.log.Debug("received broadcasted event", zap.Cid("event_cid", ev.Cid))
		c.pendingReceiveEvents <- ev

		if c.onNewEvent != nil {
			c.onNewEvent(ev)
		}
	}
}

// receiveEventLoop processes incoming Events from broadcasts.
// It consumes pendingReceiveEvents and writes into pendingLinks.
func (c *CRDT) receiveEventLoop(ctx context.Context) {
loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		case ev := <-c.pendingReceiveEvents:
			// c.log.Debug("adding event", zap.Cid("event", ev.cid))
			added, err := c.store.AddHead(ev)
			if err != nil {
				// requeue for later
				// TODO: may need a delay
				// TODO: if the channel is full, this will lock up the loop
				c.pendingReceiveEvents <- ev
			}
			if added {
				for _, link := range ev.Links {
					c.pendingLinks <- link
				}
			}
		}
	}
}

// syncLoop fetches missing events from links.
// It consumes pendingLinks and writes into pendingSyncEvents
func (c *CRDT) syncLinkLoop(ctx context.Context) {
loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		case cid := <-c.pendingLinks:
			// c.log.Debug("checking link", zap.Cid("link", cid))
			// If the CID is in heads, it should be removed because
			// we have an event that points to it.
			// We also don't need to fetch it since we already have it.
			haveAlready, err := c.store.RemoveHead(cid)
			if err != nil {
				// requeue for later
				// TODO: may need a delay
				// TODO: if the channel is full, this will lock up the loop
				c.pendingLinks <- cid
				continue
			}
			if haveAlready {
				continue
			}
			c.log.Debug("fetching link", zap.Cid("link", cid))
			cids := []mh.Multihash{cid}
			evs, err := c.syncer.Fetch(cids)
			if err != nil {
				// requeue for later
				// TODO: this will need refinement for invalid, missing cids etc.
				// TODO: if the channel is full, this will lock up the loop
				c.pendingLinks <- cid
			}
			for i, ev := range evs {
				if ev == nil {
					// requeue missing links
					c.pendingLinks <- cids[i]
					continue
				}
				c.pendingSyncEvents <- ev
			}
		}
	}
}

// syncEventLoop processes missing events that were fetched from links.
// It consumes pendingSyncEvents and writes into pendingLinks.
// TODO: There is channel read/write cycle between the two sync loops,
// i.e. they could potentially lock up if both channels fill up.
func (c *CRDT) syncEventLoop(ctx context.Context) {
loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		case ev := <-c.pendingSyncEvents:
			// c.log.Debug("adding link event", zap.Cid("event", ev.cid))
			added, err := c.store.AddEvent(ev)
			if err != nil {
				// requeue for later
				// TODO: may need a delay
				// TODO: if the channel is full, this will lock up the loop
				c.pendingSyncEvents <- ev
			}
			if added {
				for _, link := range ev.Links {
					// TODO: if the channel is full, this will lock up the loop
					c.pendingLinks <- link
				}
			}
		}
	}
}

// Bootstrap from the contents of the store.
func (c *CRDT) bootstrap(ctx context.Context) error {
	links, err := c.store.FindMissingLinks()
	if err != nil {
		return err
	}
	for _, link := range links {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case c.pendingLinks <- link:
		}
	}
	return nil
}
