package node

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"math/rand"
	"time"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multihash"
	"github.com/pkg/errors"
	"github.com/xmtp/xmtpd/pkg/context"
	"github.com/xmtp/xmtpd/pkg/crdt/types"
	"github.com/xmtp/xmtpd/pkg/zap"
)

const (
	syncProtocol      = "/xmtp/sync/0.0.1"
	reqFetch     byte = 1
)

var (
	syncEndian = binary.BigEndian

	ErrNoSyncPeer = errors.New("no peer to sync with")
)

type syncer struct {
	metrics *Metrics
	host    host.Host
	topic   string
}

func (n *Node) setSyncHandler() {
	n.host.SetStreamHandler(syncProtocol, func(s network.Stream) {
		defer s.Close()
		log := n.ctx.Logger().Named("syncHandler")
		failed := func(msg string, err error) bool {
			if err == nil {
				return false
			}
			log.Error(msg, zap.Error(err), zap.PeerID("peer", s.Conn().RemotePeer()))
			_ = s.Reset()
			return true
		}
		r := bufio.NewReader(s)
		code, err := r.ReadByte()
		if failed("reading request code", err) {
			return
		}
		switch code {
		case reqFetch:
			topic, cids, err := readFetchRequest(r)
			if failed("reading fetch request", err) {
				return
			}
			if failed("fetching", n.doFetch(s, topic, cids)) {
				return
			}
		default:
			failed("reading request code", fmt.Errorf("unknown request code (%d)", code))
			return
		}
	})
}

func (s *syncer) Fetch(ctx context.Context, cids []multihash.Multihash) (evs []*types.Event, err error) {
	start := time.Now()
	peer, err := randomPeer(s.host)
	if err != nil {
		return nil, err
	}
	// open stream
	stream, err := s.host.NewStream(ctx, peer, syncProtocol)
	if err != nil {
		return nil, err
	}
	defer stream.Close()
	// use bufio to handle short writes
	w := bufio.NewWriter(stream)
	// send FETCH request
	if err = s.writeFetchRequest(w, cids); err != nil {
		return nil, err
	}
	evs, err = s.readFetchResponse(bufio.NewReader(stream))
	if err != nil {
		return nil, err
	}
	s.metrics.recordFetch(ctx, s.topic, time.Since(start))
	return evs, err
}

func (s *syncer) Close() error {
	return nil
}

func (n *Node) doFetch(s network.Stream, topic string, cids []multihash.Multihash) error {
	r, err := n.getTopic(topic)
	if err != nil {
		return err
	}
	evs, err := r.GetEvents(n.ctx, cids...)
	if err != nil {
		return err
	}
	return writeFetchResponse(bufio.NewWriter(s), topic, evs)
}

// Fetch request format
//  1. byte = request code (1)
//  2. byte = topic length
//  3. bytes = topic
//  4. byte = count of CIDs
//  5. count time repeated:
//     byte = cid length
//     bytes = cid
func (s *syncer) writeFetchRequest(w *bufio.Writer, cids []multihash.Multihash) (err error) {
	if err = w.WriteByte(reqFetch); err != nil {
		return err
	}
	if err = writeByte(w, "topic length", len(s.topic)); err != nil {
		return err
	}
	if _, err = w.WriteString(s.topic); err != nil {
		return err
	}
	if err = writeByte(w, "cid count", len(cids)); err != nil {
		return err
	}
	for _, cid := range cids {
		if err = writeByte(w, "cid size", len(cid)); err != nil {
			return err
		}
		if _, err = w.Write(cid); err != nil {
			return err
		}
	}
	return w.Flush()
}

func readFetchRequest(r *bufio.Reader) (topic string, cids []multihash.Multihash, err error) {
	topicLen, err := r.ReadByte()
	if err != nil {
		return "", nil, err
	}
	topicBytes := make([]byte, topicLen)
	if _, err = r.Read(topicBytes); err != nil {
		return "", nil, err
	}
	count, err := r.ReadByte()
	if err != nil {
		return "", nil, err
	}
	for i := byte(0); i < count; i++ {
		cidLen, err := r.ReadByte()
		if err != nil {
			return "", nil, err
		}
		cidBytes := make([]byte, cidLen)
		if _, err := r.Read(cidBytes); err != nil {
			return "", nil, err
		}
		_, cid, err := multihash.MHFromBytes(cidBytes)
		if err != nil {
			return "", nil, err
		}
		cids = append(cids, cid)
	}
	return string(topicBytes), cids, nil
}

// Fetch response format
//  1. byte = topic length
//  2. bytes = topic
//  3. byte = count of events
//  4. count time repeated:
//     uint32 = event byte length
//     bytes = event bytes
func writeFetchResponse(w *bufio.Writer, topic string, evs []*types.Event) (err error) {
	if err = writeByte(w, "topic length", len(topic)); err != nil {
		return err
	}
	if _, err = w.WriteString(topic); err != nil {
		return err
	}
	if err = writeByte(w, "event count", len(evs)); err != nil {
		return err
	}
	for _, ev := range evs {
		evBytes, err := ev.ToBytes()
		if err != nil {
			return err
		}
		if err = binary.Write(w, syncEndian, uint32(len(evBytes))); err != nil {
			return err
		}
		if _, err = w.Write(evBytes); err != nil {
			return err
		}
	}
	return w.Flush()
}

func (s *syncer) readFetchResponse(r *bufio.Reader) (evs []*types.Event, err error) {
	topicLen, err := r.ReadByte()
	if err != nil {
		return nil, err
	}
	topicBytes := make([]byte, topicLen)
	if _, err = r.Read(topicBytes); err != nil {
		return nil, err
	}
	count, err := r.ReadByte()
	if err != nil {
		return nil, err
	}
	for i := byte(0); i < count; i++ {
		var evBytesLen uint32
		if err := binary.Read(r, syncEndian, &evBytesLen); err != nil {
			return nil, err
		}
		evBytes := make([]byte, evBytesLen)
		if _, err := r.Read(evBytes); err != nil {
			return nil, err
		}
		ev, err := types.EventFromBytes(evBytes)
		if err != nil {
			return nil, err
		}
		evs = append(evs, ev)
	}
	return evs, nil
}

func writeByte(w *bufio.Writer, field string, i int) error {
	if i > 255 {
		return fmt.Errorf("%s out of byte range (%d)", field, i)
	}
	return w.WriteByte(byte(i))
}

func randomPeer(h host.Host) (peer.ID, error) {
	var peers peer.IDSlice
	for _, p := range h.Network().Peers() {
		if p != h.ID() {
			peers = append(peers, p)
		}
	}
	if len(peers) == 0 {
		return h.ID(), ErrNoSyncPeer
	}
	return peers[rand.Intn(len(peers))], nil
}
