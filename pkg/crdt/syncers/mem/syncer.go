package memsyncer

import (
	"github.com/multiformats/go-multihash"
	"github.com/xmtp/xmtpd/pkg/crdt/types"
	"github.com/xmtp/xmtpd/pkg/zap"
)

type MemorySyncer struct {
	log *zap.Logger
}

func New(log *zap.Logger) *MemorySyncer {
	return &MemorySyncer{
		log: log,
	}
}

func (s *MemorySyncer) Fetch(cids []multihash.Multihash) (results []*types.Event, err error) {
	return nil, types.ErrTODO
}
