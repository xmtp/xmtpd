package client

import (
	"context"

	messagev1 "github.com/xmtp/proto/v3/go/message_api/v1"
)

type Client interface {
	Publish(context.Context, *messagev1.PublishRequest) (*messagev1.PublishResponse, error)
	Subscribe(context.Context, *messagev1.SubscribeRequest) (Stream, error)
	SubscribeAll(context.Context) (Stream, error)
	Query(context.Context, *messagev1.QueryRequest) (*messagev1.QueryResponse, error)
	BatchQuery(ctx context.Context, req *messagev1.BatchQueryRequest) (*messagev1.BatchQueryResponse, error)
	Close() error
}

// Stream is an abstraction of the subscribe response stream
type Stream interface {
	// Next returns io.EOF when the stream ends or is closed from either side.
	Next(ctx context.Context) (*messagev1.Envelope, error)
	// Closing the stream terminates the subscription.
	Close() error
}
