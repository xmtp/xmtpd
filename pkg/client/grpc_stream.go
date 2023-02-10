package client

import (
	"context"
	"errors"
	"io"

	messagev1 "github.com/xmtp/proto/v3/go/message_api/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type grpcStream struct {
	cancel context.CancelFunc
	stream messagev1.MessageApi_SubscribeClient
}

func (s *grpcStream) Next(ctx context.Context) (*messagev1.Envelope, error) {
	envC := make(chan *messagev1.Envelope)
	errC := make(chan error)
	go func() {
		env, err := s.stream.Recv()
		if ctx.Err() != nil {
			// If the context has already closed, then just return out of this.
			return
		}
		if err != nil {
			grpcErr, ok := status.FromError(err)
			if ok {
				if status.Code(err) == codes.Canceled {
					err = io.EOF
				} else {
					err = errors.New(grpcErr.Message())
				}
			}
			errC <- err
			return
		}
		envC <- env
	}()

	var env *messagev1.Envelope
	select {
	case v := <-envC:
		env = v
	case err := <-errC:
		return nil, err
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	return env, nil
}

func (s *grpcStream) Close() error {
	s.cancel()
	return nil
}
