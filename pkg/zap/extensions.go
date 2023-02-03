package zap

import (
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

// cid provides uniform logging for individual CIDs
type cid mh.Multihash

func Cid(key string, c mh.Multihash) zapcore.Field {
	return zap.Stringer(key, cid(c))
}

func (c cid) String() string {
	return ShortenedCid(mh.Multihash(c))
}

// cidSlice provides uniform logging for lists of CIDs
type cidSlice []mh.Multihash

func (cids cidSlice) MarshalLogArray(enc zapcore.ArrayEncoder) error {
	for _, cid := range cids {
		enc.AppendString(ShortenedCid(cid))
	}
	return nil
}

func Cids(key string, cids ...mh.Multihash) zapcore.Field {
	return zap.Array(key, cidSlice(cids))
}
