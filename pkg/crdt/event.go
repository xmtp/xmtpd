package crdt

import (
	"bytes"
	"encoding/binary"
	"io"

	mh "github.com/multiformats/go-multihash"
	messagev1 "github.com/xmtp/proto/v3/go/message_api/v1"
)

// Event represents a node in the Merkle-Clock
// It captures a message and links to its preceding Events
type Event struct {
	*messagev1.Envelope
	Links []mh.Multihash // cid's of direct ancestors
	Cid   mh.Multihash   // cid is computed by hashing the links and message together
}

// NewEvent creates an event from a message and a set of links to preceding events (heads)
func NewEvent(env *messagev1.Envelope, heads []mh.Multihash) (*Event, error) {
	ev := &Event{Envelope: env, Links: heads}
	var err error
	ev.Cid, err = mh.SumStream(ev.Reader(), mh.SHA2_256, -1)
	if err != nil {
		return nil, err
	}
	return ev, nil
}

// Reader creates a chunk reader for given Event.
func (ev *Event) Reader() *chunkReader {
	// compose the chunks of the Event data
	var chunks [][]byte
	if ev.Envelope != nil {
		head := make([]byte, 8+len(ev.ContentTopic))
		binary.BigEndian.PutUint64(head, ev.TimestampNs) // timestamp
		copy(head[8:], ev.ContentTopic)                  // topic
		chunks = append(chunks, head, ev.Message)        // message payload
	}
	for _, link := range ev.Links { // links
		chunks = append(chunks, link)
	}
	return &chunkReader{chunks, 0}
}

// Compare returns an integer comparing two events based on their timestamps.
// The result will be negative if ev < ev2, and positive if ev > ev2.
// The result can only be 0 if ev and ev2 are the same event.
// TODO: total order should reflect the DAG first and foremost.
func (ev *Event) Compare(ev2 *Event) int {
	res := ev.TimestampNs - ev2.TimestampNs
	if res != 0 {
		return int(res)
	}
	return bytes.Compare(ev.Cid, ev2.Cid)
}

// chunkReader helps computing an Event CID efficiently by
// yielding the bytes composed of the various bits of the Event
// without having to concatenate them all.
// This allows passing the reader to mh.SumStream()
type chunkReader struct {
	unreadChunks [][]byte // chunks of the Event data to be hashed
	pos          int      // current position from the start of the next chunk
}

func (r *chunkReader) Read(b []byte) (n int, err error) {
	total := 0
	for len(b) > 0 && len(r.unreadChunks) > 0 {
		chunk := r.unreadChunks[0]
		n := copy(b, chunk[r.pos:])
		total += n
		b = b[n:]
		r.pos += n
		if r.pos == len(chunk) {
			r.pos = 0
			r.unreadChunks = r.unreadChunks[1:]
		}
	}
	if len(r.unreadChunks) > 0 {
		return total, nil
	}
	return total, io.EOF
}
