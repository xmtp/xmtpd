package types

import (
	"encoding/binary"
	"io"

	"github.com/multiformats/go-multihash"
	messagev1 "github.com/xmtp/proto/v3/go/message_api/v1"
)

// Event represents a node in the Merkle-Clock
// It captures a message and links to its preceding Events
type Event struct {
	*messagev1.Envelope
	Links []multihash.Multihash // cid's of direct ancestors
	Cid   multihash.Multihash   // cid is computed by hashing the links and message together
}

// NewEvent creates an event from a message and a set of links to preceding events (heads)
func NewEvent(env *messagev1.Envelope, heads []multihash.Multihash) (*Event, error) {
	ev := &Event{
		Envelope: env,
		Links:    heads,
	}
	var err error
	ev.Cid, err = multihash.SumStream(ev.Reader(), multihash.SHA2_256, -1)
	if err != nil {
		return nil, err
	}
	return ev, nil
}

// Reader creates a chunk reader for given Event.
func (ev *Event) Reader() *chunkReader {
	chunks := make([][]byte, 0, len(ev.Links)+1)
	if ev.Envelope != nil {
		head := make([]byte, 8+len(ev.ContentTopic))
		binary.BigEndian.PutUint64(head, ev.TimestampNs) // timestamp
		copy(head[8:], ev.ContentTopic)                  // topic
		chunks = append(chunks, head, ev.Message)        // message payload
	}
	for _, link := range ev.Links {
		chunks = append(chunks, link)
	}
	return &chunkReader{chunks, 0}
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
