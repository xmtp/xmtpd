package types

import (
	"encoding/binary"
	"errors"
	"io"
	"math"

	"github.com/multiformats/go-multihash"
	messagev1 "github.com/xmtp/proto/v3/go/message_api/v1"
	"google.golang.org/protobuf/proto"
)

var ErrInvalidCids = errors.New("Event CID or Link CIDs are invalid")

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

// Event bytes are encoded as a list of event cid and link cids followed by the Envelope.
// The list of cids is encoded as total size of the list in bytes encoded as uvarint,
// followed by the multihash bytes of each cid (which are themselves length prefixed).
// The total size prefix allows efficiently skipping to the Envelope bytes.
// Envelope bytes are stored as prescribed by github.com/xmtp/proto.

func (ev *Event) ToBytes() ([]byte, error) {
	envBytes, err := proto.Marshal(ev.Envelope)
	if err != nil {
		return nil, err
	}
	linksSize := len(ev.Cid)
	for _, link := range ev.Links {
		linksSize += len(link)
	}
	prefix := binary.AppendUvarint(nil, uint64(linksSize))
	b := make([]byte, len(prefix)+linksSize+len(envBytes))
	copy(b, prefix)
	pos := len(prefix)
	copy(b[pos:], ev.Cid)
	pos += len(ev.Cid)
	for _, link := range ev.Links {
		copy(b[pos:], link)
		pos += len(link)
	}
	copy(b[pos:], envBytes)
	return b, nil
}

func EventFromBytes(evBytes []byte) (*Event, error) {
	cid, links, remainder, err := readLinks(evBytes)
	if err != nil {
		return nil, err
	}
	env, err := unmarshalEnvelope(remainder)
	if err != nil {
		return nil, err
	}
	return &Event{
		Cid:      cid,
		Links:    links,
		Envelope: env,
	}, nil
}

func readLinks(evBytes []byte) (cid multihash.Multihash, links []multihash.Multihash, remainder []byte, err error) {
	cidsSize, n := binary.Uvarint(evBytes)
	if n <= 0 || cidsSize == 0 || cidsSize > math.MaxInt {
		return nil, nil, nil, ErrInvalidCids
	}
	pos := n
	end := int(cidsSize) + n
	for pos < end {
		n, link, err := multihash.MHFromBytes(evBytes[pos:])
		if err != nil {
			return nil, nil, nil, err
		}
		if n <= 0 {
			return nil, nil, nil, ErrInvalidCids
		}
		links = append(links, link)
		pos += n
	}
	return links[0], links[1:], evBytes[pos:], nil
}

func unmarshalEnvelope(envBytes []byte) (*messagev1.Envelope, error) {
	var env messagev1.Envelope
	if err := proto.Unmarshal(envBytes, &env); err != nil {
		return nil, err
	}
	return &env, nil
}
