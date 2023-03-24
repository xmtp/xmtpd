package crdt

import (
	"bytes"

	mh "github.com/multiformats/go-multihash"
	messagev1 "github.com/xmtp/proto/v3/go/message_api/v1"
	"github.com/xmtp/xmtpd/pkg/context"
	"github.com/xmtp/xmtpd/pkg/crdt/types"
	"github.com/xmtp/xmtpd/pkg/zap"
)

type NewEventFunc func(ev *types.Event)

// Replica manages the DAG of a dataset replica.
type Replica struct {
	log     *zap.Logger
	ctx     context.Context
	metrics *Metrics

	onNewEvent NewEventFunc

	store       Store
	broadcaster Broadcaster
	syncer      Syncer

	pendingReceiveEvents chan *types.Event // broadcasted events that were received from the network but not processed yet
	pendingSyncEvents    chan *types.Event // missing events that were fetched from the network but not processed yet
	pendingLinks         chan mh.Multihash // missing links that were discovered but not successfully fetched yet
}

func NewReplica(ctx context.Context, metrics *Metrics, store Store, bc Broadcaster, syncer Syncer, onNewEvent NewEventFunc) (*Replica, error) {
	r := &Replica{
		log:        ctx.Logger(),
		ctx:        ctx,
		metrics:    metrics,
		onNewEvent: onNewEvent,

		store:       store,
		broadcaster: bc,
		syncer:      syncer,

		// TODO: tuning the channel sizes will likely be important
		// current implementation can lock up if the channels fill up.
		pendingReceiveEvents: make(chan *types.Event, 20),
		pendingSyncEvents:    make(chan *types.Event, 20),
		pendingLinks:         make(chan mh.Multihash, 200),
	}

	r.ctx.Go(r.receiveEventLoop)
	r.ctx.Go(r.syncEventLoop)
	r.ctx.Go(r.syncLinkLoop)
	r.ctx.Go(r.nextBroadcastedEventLoop)

	err := r.bootstrap()
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (r *Replica) GetEvents(ctx context.Context, cids ...mh.Multihash) ([]*types.Event, error) {
	return r.store.GetEvents(ctx, cids...)
}

func (r *Replica) BroadcastAppend(ctx context.Context, env *messagev1.Envelope) (*types.Event, error) {
	ev, err := r.store.AppendEvent(ctx, env)
	if err != nil {
		return nil, err
	}
	err = r.broadcaster.Broadcast(ctx, ev)
	if err != nil {
		return nil, err
	}
	return ev, err
}

func (r *Replica) Query(ctx context.Context, req *messagev1.QueryRequest) (*messagev1.QueryResponse, error) {
	return r.store.Query(ctx, req)
}

func (r *Replica) nextBroadcastedEventLoop(ctx context.Context) {
	log := r.log.Named("nextBroadcastedEventLoop")
	for {
		ev, err := r.broadcaster.Next(ctx)
		if err != nil {
			if err == context.Canceled {
				log.Debug("context closed", zap.Error(err))
				return
			}
			log.Error("error getting next broadcasted event", zap.Error(err))
			return
		}
		log.Debug("received broadcasted event", zap.Cid("event", ev.Cid))

		r.metrics.recordFreeSpaceInEvents(ctx, r.pendingReceiveEvents, false)
		select {
		case r.pendingReceiveEvents <- ev:
		case <-ctx.Done():
			log.Debug("context closed", zap.Error(ctx.Err()))
			return
		}

		r.metrics.recordReceivedEvent(ctx, ev, false)

		if r.onNewEvent != nil {
			r.onNewEvent(ev)
		}
	}
}

// receiveEventLoop processes incoming Events from broadcasts.
// It consumes pendingReceiveEvents and writes into pendingLinks.
func (r *Replica) receiveEventLoop(ctx context.Context) {
	log := r.log.Named("receiveEventLoop")
	for {
		select {
		case <-ctx.Done():
			log.Debug("context closed", zap.Error(ctx.Err()))
			return
		case ev := <-r.pendingReceiveEvents:
			added, err := r.store.InsertHead(ctx, ev)
			if err != nil {
				log.Error("error inserting head", zap.Cid("event", ev.Cid), zap.Error(err))
				// requeue for later
				// TODO: may need a delay
				// TODO: if the channel is full, this will lock up the loop
				r.metrics.recordFreeSpaceInEvents(ctx, r.pendingReceiveEvents, false)
				select {
				case r.pendingReceiveEvents <- ev:
				case <-ctx.Done():
					log.Debug("context closed", zap.Error(ctx.Err()))
					return
				}
			}
			if added {
				for _, link := range ev.Links {
					r.metrics.recordFreeSpaceInLinks(ctx, r.pendingLinks)
					select {
					case r.pendingLinks <- link:
					case <-ctx.Done():
						log.Debug("context closed", zap.Error(ctx.Err()))
						return
					}
				}
			}
		}
	}
}

// syncLoop fetches missing events from links.
// It consumes pendingLinks and writes into pendingSyncEvents
func (r *Replica) syncLinkLoop(ctx context.Context) {
	log := r.log.Named("syncLinkLoop")
	for {
		select {
		case <-ctx.Done():
			log.Debug("context closed", zap.Error(ctx.Err()))
			return
		case cid := <-r.pendingLinks:
			// r.log.Debug("checking link", zap.Cid("link", cid))
			// If the CID is in heads, it should be removed because
			// we have an event that points to it.
			// We also don't need to fetch it since we already have it.
			removed, err := r.store.RemoveHead(ctx, cid)
			if err != nil {
				log.Error("error removing head", zap.Cid("event", cid), zap.Error(err))
				r.metrics.recordFreeSpaceInLinks(ctx, r.pendingLinks)
				// requeue for later
				// TODO: may need a delay
				// TODO: if the channel is full, this will lock up the loop
				select {
				case r.pendingLinks <- cid:
					continue
				case <-ctx.Done():
					log.Debug("context closed", zap.Error(ctx.Err()))
					return
				}
			}
			if removed {
				continue
			}
			log.Debug("fetching link", zap.Cid("link", cid))
			cids := []mh.Multihash{cid}
			evs, err := r.syncer.Fetch(ctx, cids)
			if err != nil {
				log.Error("error fetching event", zap.Cids("event", cids...), zap.Error(err))
				r.metrics.recordFreeSpaceInLinks(ctx, r.pendingLinks)
				// requeue for later
				// TODO: this will need refinement for invalid, missing cids etc.
				// TODO: if the channel is full, this will lock up the loop
				select {
				case r.pendingLinks <- cid:
				case <-ctx.Done():
					log.Debug("context closed", zap.Error(ctx.Err()))
					return
				}
			}
			for _, cid := range cids {
				ev := findEvent(cid, evs)
				if ev == nil {
					r.metrics.recordFreeSpaceInLinks(ctx, r.pendingLinks)
					// requeue missing links
					select {
					case r.pendingLinks <- cid:
					case <-ctx.Done():
						log.Debug("context closed", zap.Error(ctx.Err()))
						return
					}
				} else {
					r.metrics.recordFreeSpaceInEvents(ctx, r.pendingSyncEvents, true)
					select {
					case r.pendingSyncEvents <- ev:
					case <-ctx.Done():
						log.Debug("context closed", zap.Error(ctx.Err()))
						return
					}
				}
			}
		}
	}
}

// syncEventLoop processes missing events that were fetched from links.
// It consumes pendingSyncEvents and writes into pendingLinks.
// TODO: There is channel read/write cycle between the two sync loops,
// i.e. they could potentially lock up if both channels fill up.
func (r *Replica) syncEventLoop(ctx context.Context) {
	log := r.log.Named("syncEventLoop")
	for {
		select {
		case <-ctx.Done():
			log.Debug("context closed", zap.Error(ctx.Err()))
			return
		case ev := <-r.pendingSyncEvents:
			added, err := r.store.InsertEvent(ctx, ev)
			if err != nil {
				log.Error("error inserting event", zap.Cid("event", ev.Cid), zap.Error(err))
				r.metrics.recordFreeSpaceInEvents(ctx, r.pendingSyncEvents, true)
				// requeue for later
				// TODO: may need a delay
				// TODO: if the channel is full, this will lock up the loop
				select {
				case r.pendingSyncEvents <- ev:
				case <-ctx.Done():
					log.Debug("context closed", zap.Error(ctx.Err()))
					return
				}
			}

			r.metrics.recordReceivedEvent(ctx, ev, true)

			if added {
				for _, link := range ev.Links {
					r.metrics.recordFreeSpaceInLinks(ctx, r.pendingLinks)
					// TODO: if the channel is full, this will lock up the loop
					select {
					case r.pendingLinks <- link:
					case <-ctx.Done():
						log.Debug("context closed", zap.Error(ctx.Err()))
						return
					}
				}
			}
		}
	}
}

// Bootstrap from the contents of the store.
func (r *Replica) bootstrap() error {
	log := r.log.Named("bootstrap")
	links, err := r.store.FindMissingLinks(r.ctx)
	if err != nil {
		return err
	}
	for _, link := range links {
		r.metrics.recordFreeSpaceInLinks(r.ctx, r.pendingLinks)
		select {
		case <-r.ctx.Done():
			log.Debug("context closed", zap.Error(r.ctx.Err()))
			return r.ctx.Err()
		case r.pendingLinks <- link:
		}
	}
	return nil
}

func findEvent(cid mh.Multihash, evs []*types.Event) *types.Event {
	for _, ev := range evs {
		if bytes.Equal(cid, ev.Cid) {
			return ev
		}
	}
	return nil
}
