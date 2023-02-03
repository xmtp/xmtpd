package crdt

import (
	"context"

	mh "github.com/multiformats/go-multihash"
	messagev1 "github.com/xmtp/proto/v3/go/message_api/v1"
	"github.com/xmtp/xmtpd/pkg/zap"
)

// Topic manages the DAG of a topic replica.
// It implements the topic API, as well as the
// replication mechanism using the store, broadcaster and syncer.
type Topic struct {
	name                 string            // the topic name
	pendingReceiveEvents chan *Event       // broadcasted events that were received from the network but not processed yet
	pendingSyncEvents    chan *Event       // missing events that were fetched from the network but not processed yet
	pendingLinks         chan mh.Multihash // missing links that were discovered but not successfully fetched yet
	log                  *zap.Logger

	TopicStore
	TopicSyncer
	TopicBroadcaster
}

// Creates a new topic replica
func NewTopic(ctx context.Context, name string, log *zap.Logger, store TopicStore, syncer TopicSyncer, bc TopicBroadcaster) *Topic {
	t := &Topic{
		name: name,
		// TODO: tuning the channel sizes will likely be important
		// current implementation can lock up if the channels fill up.
		pendingReceiveEvents: make(chan *Event, 20),
		pendingSyncEvents:    make(chan *Event, 20),
		pendingLinks:         make(chan mh.Multihash, 20),
		log:                  log,
		TopicStore:           store,
		TopicSyncer:          syncer,
		TopicBroadcaster:     bc,
	}
	go t.receiveEventLoop(ctx)
	go t.syncEventLoop(ctx)
	go t.syncLinkLoop(ctx)
	return t
}

// Publish adopts a new message into a topic and broadcasts it to the network.
func (t *Topic) Publish(ctx context.Context, env *messagev1.Envelope) (*Event, error) {
	ev, err := t.NewEvent(env)
	if err != nil {
		return nil, err
	}
	t.Broadcast(ev)
	return ev, nil
}

// receiveEventLoop processes incoming Events from broadcasts.
// It consumes pendingReceiveEvents and writes into pendingLinks.
func (t *Topic) receiveEventLoop(ctx context.Context) {
loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		case ev := <-t.pendingReceiveEvents:
			// t.log.Debug("adding event", zap.Cid("event", ev.cid))
			added, err := t.AddHead(ev)
			if err != nil {
				// requeue for later
				// TODO: may need a delay
				// TODO: if the channel is full, this will lock up the loop
				t.pendingReceiveEvents <- ev
			}
			if added {
				for _, link := range ev.links {
					t.pendingLinks <- link
				}
			}
		}
	}
}

// syncLoop fetches missing events from links.
// It consumes pendingLinks and writes into pendingSyncEvents
func (t *Topic) syncLinkLoop(ctx context.Context) {
loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		case cid := <-t.pendingLinks:
			// t.log.Debug("checking link", zap.Cid("link", cid))
			// If the CID is in heads, it should be removed because
			// we have an event that points to it.
			// We also don't need to fetch it since we already have it.
			haveAlready, err := t.RemoveHead(cid)
			if err != nil {
				// requeue for later
				// TODO: may need a delay
				// TODO: if the channel is full, this will lock up the loop
				t.pendingLinks <- cid
				continue
			}
			if haveAlready {
				continue
			}
			t.log.Debug("fetching link", zap.Cid("link", cid))
			cids := []mh.Multihash{cid}
			evs, err := t.Fetch(cids)
			if err != nil {
				// requeue for later
				// TODO: this will need refinement for invalid, missing cids etc.
				// TODO: if the channel is full, this will lock up the loop
				t.pendingLinks <- cid
			}
			for i, ev := range evs {
				if ev == nil {
					// requeue missing links
					t.pendingLinks <- cids[i]
					continue
				}
				t.pendingSyncEvents <- ev
			}
		}
	}
}

// syncEventLoop processes missing events that were fetched from links.
// It consumes pendingSyncEvents and writes into pendingLinks.
// TODO: There is channel read/write cycle between the two sync loops,
// i.e. they could potentially lock up if both channels fill up.
func (t *Topic) syncEventLoop(ctx context.Context) {
loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		case ev := <-t.pendingSyncEvents:
			// t.log.Debug("adding link event", zap.Cid("event", ev.cid))
			added, err := t.AddEvent(ev)
			if err != nil {
				// requeue for later
				// TODO: may need a delay
				// TODO: if the channel is full, this will lock up the loop
				t.pendingSyncEvents <- ev
			}
			if added {
				for _, link := range ev.links {
					// TODO: if the channel is full, this will lock up the loop
					t.pendingLinks <- link
				}
			}
		}
	}
}

// Bootstrap the topic from the contents of the topic store.
// This is called from a goroutine group during node creation.
func (t *Topic) bootstrap(ctx context.Context) error {
	links, err := t.FindMissingLinks()
	if err != nil {
		return err
	}
	for _, link := range links {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case t.pendingLinks <- link:
		}
	}
	return nil
}
