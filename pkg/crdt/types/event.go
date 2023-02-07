package types

import (
	"io"

	mh "github.com/multiformats/go-multihash"
)

// Event represents a node in the Merkle-Clock
// It captures a message and links to its preceding Events
type Event struct {
	Links   []mh.Multihash // cid's of direct ancestors
	Cid     mh.Multihash   // cid is computed by hashing the links and message together
	Payload []byte
}

// NewEvent creates an event from a message and a set of links to preceding events (heads)
func NewEvent(payload []byte, heads []mh.Multihash) (*Event, error) {
	ev := &Event{
		Links:   heads,
		Payload: payload,
	}
	var err error
	ev.Cid, err = mh.SumStream(ev.Reader(), mh.SHA2_256, -1)
	if err != nil {
		return nil, err
	}
	return ev, nil
}

// Reader creates a chunk reader for given Event.
func (ev *Event) Reader() *chunkReader {
	chunks := make([][]byte, len(ev.Links)+1)
	chunks[0] = ev.Payload
	for i, link := range ev.Links {
		chunks[i+1] = link
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
