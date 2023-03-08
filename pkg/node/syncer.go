package node

import (
	"github.com/multiformats/go-multihash"
	"github.com/xmtp/xmtpd/pkg/context"
	"github.com/xmtp/xmtpd/pkg/crdt/types"
)

type nilSyncer struct{}

func (s *nilSyncer) Fetch(context.Context, []multihash.Multihash) ([]*types.Event, error) {
	return nil, nil
}

func (s *nilSyncer) Close() error {
	return nil
}
