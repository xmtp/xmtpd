package store

import (
	"context"

	messagev1 "github.com/xmtp/proto/v3/go/message_api/v1"
	crdt "github.com/xmtp/xmtpd/pkg/crdt"
)

type Store interface {
	crdt.Store

	InsertEnvelope(ctx context.Context, env *messagev1.Envelope) error
	QueryEnvelopes(ctx context.Context, req *messagev1.QueryRequest) (*messagev1.QueryResponse, error)
	Close() error
}
