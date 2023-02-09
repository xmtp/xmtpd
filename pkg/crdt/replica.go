package crdt

import (
	"context"

	mh "github.com/multiformats/go-multihash"
	"github.com/xmtp/xmtpd/pkg/crdt/types"
	"github.com/xmtp/xmtpd/pkg/zap"
)

type NewEventFunc func(ev *types.Event)

// Replica manages the DAG of a dataset replica.
type Replica struct {
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

func NewReplica(ctx context.Context, log *zap.Logger, store Store, bc Broadcaster, syncer Syncer, onNewEvent NewEventFunc) (*Replica, error) {
	ctx, ctxCancel := context.WithCancel(ctx)
	r := &Replica{
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

	go r.receiveEventLoop(ctx)
	go r.syncEventLoop(ctx)
	go r.syncLinkLoop(ctx)
	go r.nextBroadcastedEventLoop(ctx)

	err := r.bootstrap(ctx)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (r *Replica) Close() error {
	if r.ctxCancel != nil {
		r.ctxCancel()
	}
	return nil
}

func (r *Replica) BroadcastAppend(ctx context.Context, payload []byte) (*types.Event, error) {
	ev, err := r.store.AppendEvent(payload)
	if err != nil {
		return nil, err
	}
	return ev, r.broadcaster.Broadcast(ev)
}

func (r *Replica) nextBroadcastedEventLoop(ctx context.Context) {
	for {
		ev, err := r.broadcaster.Next(ctx)
		if err != nil {
			if err == context.Canceled {
				return
			}
			r.log.Error("error getting next broadcasted event", zap.Error(err))
			return
		}
		r.log.Debug("received broadcasted event", zap.Cid("event_cid", ev.Cid))
		r.pendingReceiveEvents <- ev

		if r.onNewEvent != nil {
			r.onNewEvent(ev)
		}
	}
}

// receiveEventLoop processes incoming Events from broadcasts.
// It consumes pendingReceiveEvents and writes into pendingLinks.
func (r *Replica) receiveEventLoop(ctx context.Context) {
loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		case ev := <-r.pendingReceiveEvents:
			// r.log.Debug("adding event", zap.Cid("event", ev.cid))
			added, err := r.store.InsertHead(ev)
			if err != nil {
				// requeue for later
				// TODO: may need a delay
				// TODO: if the channel is full, this will lock up the loop
				r.pendingReceiveEvents <- ev
			}
			if added {
				for _, link := range ev.Links {
					r.pendingLinks <- link
				}
			}
		}
	}
}

// syncLoop fetches missing events from links.
// It consumes pendingLinks and writes into pendingSyncEvents
func (r *Replica) syncLinkLoop(ctx context.Context) {
loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		case cid := <-r.pendingLinks:
			// r.log.Debug("checking link", zap.Cid("link", cid))
			// If the CID is in heads, it should be removed because
			// we have an event that points to it.
			// We also don't need to fetch it since we already have it.
			haveAlready, err := r.store.RemoveHead(cid)
			if err != nil {
				// requeue for later
				// TODO: may need a delay
				// TODO: if the channel is full, this will lock up the loop
				r.pendingLinks <- cid
				continue
			}
			if haveAlready {
				continue
			}
			r.log.Debug("fetching link", zap.Cid("link", cid))
			cids := []mh.Multihash{cid}
			evs, err := r.syncer.Fetch(ctx, cids)
			if err != nil {
				// requeue for later
				// TODO: this will need refinement for invalid, missing cids etc.
				// TODO: if the channel is full, this will lock up the loop
				r.pendingLinks <- cid
			}
			for i, ev := range evs {
				if ev == nil {
					// requeue missing links
					r.pendingLinks <- cids[i]
					continue
				}
				r.pendingSyncEvents <- ev
			}
		}
	}
}

// syncEventLoop processes missing events that were fetched from links.
// It consumes pendingSyncEvents and writes into pendingLinks.
// TODO: There is channel read/write cycle between the two sync loops,
// i.e. they could potentially lock up if both channels fill up.
func (r *Replica) syncEventLoop(ctx context.Context) {
loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		case ev := <-r.pendingSyncEvents:
			// r.log.Debug("adding link event", zap.Cid("event", ev.cid))
			added, err := r.store.InsertEvent(ev)
			if err != nil {
				// requeue for later
				// TODO: may need a delay
				// TODO: if the channel is full, this will lock up the loop
				r.pendingSyncEvents <- ev
			}
			if added {
				for _, link := range ev.Links {
					// TODO: if the channel is full, this will lock up the loop
					r.pendingLinks <- link
				}
			}
		}
	}
}

// Bootstrap from the contents of the store.
func (r *Replica) bootstrap(ctx context.Context) error {
	links, err := r.store.FindMissingLinks()
	if err != nil {
		return err
	}
	for _, link := range links {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case r.pendingLinks <- link:
		}
	}
	return nil
}
