package memstore

import "github.com/xmtp/xmtpd/pkg/crdt/types"

func reverseEvents(in []*types.Event) (out []*types.Event) {
	for i := range in {
		out = append(out, in[len(in)-1-i])
	}
	return out
}
