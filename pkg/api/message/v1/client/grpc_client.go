package client

import (
	"context"

	messagev1 "github.com/xmtp/proto/v3/go/message_api/v1"
	"google.golang.org/grpc"
)

type grpcClient struct {
	grpc messagev1.MessageApiClient
}

func NewGRPCClient(ctx context.Context, dialFn func(context.Context) (*grpc.ClientConn, error)) (*grpcClient, error) {
	conn, err := dialFn(ctx)
	if err != nil {
		return nil, err
	}
	return &grpcClient{
		grpc: messagev1.NewMessageApiClient(conn),
	}, nil
}

func (c *grpcClient) Close() error {
	return nil
}

func (c *grpcClient) Subscribe(ctx context.Context, r *messagev1.SubscribeRequest) (Stream, error) {
	ctx, cancel := context.WithCancel(ctx)
	stream, err := c.grpc.Subscribe(ctx, r)
	if err != nil {
		cancel()
		return nil, err
	}
	return &grpcStream{
		cancel: cancel,
		stream: stream,
	}, nil
}

func (c *grpcClient) SubscribeAll(ctx context.Context) (Stream, error) {
	ctx, cancel := context.WithCancel(ctx)
	stream, err := c.grpc.SubscribeAll(ctx, &messagev1.SubscribeAllRequest{})
	if err != nil {
		cancel()
		return nil, err
	}

	return &grpcStream{
		cancel: cancel,
		stream: stream,
	}, nil
}

func (c *grpcClient) Publish(ctx context.Context, r *messagev1.PublishRequest) (*messagev1.PublishResponse, error) {
	return c.grpc.Publish(ctx, r)
}

func (c *grpcClient) Query(ctx context.Context, q *messagev1.QueryRequest) (*messagev1.QueryResponse, error) {
	return c.grpc.Query(ctx, q)
}

func (c *grpcClient) BatchQuery(ctx context.Context, q *messagev1.BatchQueryRequest) (*messagev1.BatchQueryResponse, error) {
	return c.grpc.BatchQuery(ctx, q)
}
