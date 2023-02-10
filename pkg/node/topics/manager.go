package topics

import (
	"context"

	"github.com/xmtp/xmtpd/pkg/crdt"
)

type Manager interface {
	GetOrCreateTopic(ctx context.Context, topicId string) (*crdt.Replica, error)
	Close() error
}
