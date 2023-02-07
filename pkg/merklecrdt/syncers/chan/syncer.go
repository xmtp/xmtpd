package chansyncer

import (
	"github.com/multiformats/go-multihash"
	"github.com/xmtp/xmtpd/pkg/merklecrdt/types"
	"github.com/xmtp/xmtpd/pkg/zap"
)

type ChannelSyncer struct {
	log *zap.Logger
}

func New(log *zap.Logger) *ChannelSyncer {
	return &ChannelSyncer{
		log: log,
	}
}

func (s *ChannelSyncer) Fetch(cids []multihash.Multihash) (results []*types.Event, err error) {
	return nil, types.ErrTODO
}
