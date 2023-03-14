package zap

import (
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
	mh "github.com/multiformats/go-multihash"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Fields avoids referencing zapcore when creating []zapcore.Field
func Fields(fields ...zapcore.Field) []zapcore.Field {
	return fields
}

// AnyFields creates a slice of Any fields from a key, value sequence
func AnyFields(keysAndValues ...interface{}) []zap.Field {
	fields := make([]zap.Field, 0, len(keysAndValues)/2)
	for i := 0; i < len(keysAndValues); i += 2 {
		key, ok := keysAndValues[i].(string)
		if !ok {
			continue
		}
		fields = append(fields, zap.Any(key, keysAndValues[i+1]))
	}
	return fields
}

// peerId provides uniform logging for individual peer IDs
type peerID peer.ID

func PeerID(key string, id peer.ID) zapcore.Field {
	return zap.Stringer(key, peerID(id))
}

func (id peerID) String() string {
	return peer.ID(id).Pretty()
}

// peerIDSlice provides uniform logging for lists of peer.ID
type peerIDSlice []peer.ID

func (ids peerIDSlice) MarshalLogArray(enc zapcore.ArrayEncoder) error {
	for _, id := range ids {
		enc.AppendString(id.Pretty())
	}
	return nil
}

func PeerIDs(key string, ids ...peer.ID) zapcore.Field {
	return zap.Array(key, peerIDSlice(ids))
}

// cid provides uniform logging for individual CIDs
type cid mh.Multihash

func Cid(key string, c mh.Multihash) zapcore.Field {
	return zap.Stringer(key, cid(c))
}

func (c cid) String() string {
	return ShortCid(mh.Multihash(c))
}

// cidSlice provides uniform logging for lists of CIDs
type cidSlice []mh.Multihash

func (cids cidSlice) MarshalLogArray(enc zapcore.ArrayEncoder) error {
	for _, cid := range cids {
		enc.AppendString(ShortCid(cid))
	}
	return nil
}

func Cids(key string, cids ...mh.Multihash) zapcore.Field {
	return zap.Array(key, cidSlice(cids))
}

// multiaddrSlice provides uniform logging for lists of multiaddrs
type multiaddrSlice []multiaddr.Multiaddr

func (mas multiaddrSlice) MarshalLogArray(enc zapcore.ArrayEncoder) error {
	for _, ma := range mas {
		enc.AppendString(ma.String())
	}
	return nil
}

func Multiaddrs(key string, mas ...multiaddr.Multiaddr) zapcore.Field {
	return zap.Array(key, multiaddrSlice(mas))
}
