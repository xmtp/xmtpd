package subscribers

import (
	"context"

	crdttypes "github.com/xmtp/xmtpd/pkg/crdt/types"
)

type Manager interface {
	OnNewEvent(topicId string, ev *crdttypes.Event)
	Subscribe(ctx context.Context, topicId string) chan *crdttypes.Event
	Unsubscribe(topicId string, ch chan *crdttypes.Event)
	Close() error
}
