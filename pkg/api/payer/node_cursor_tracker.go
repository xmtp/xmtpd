package payer

import (
	"context"
	"errors"
	"time"

	"github.com/xmtp/xmtpd/pkg/metrics"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/metadata_api"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type MetadataAPIClientConstructor interface {
	NewMetadataAPIClient(nodeID uint32) (metadata_api.MetadataApiClient, error)
}
type DefaultMetadataAPIClientConstructor struct {
	clientManager *ClientManager
}

func (c *DefaultMetadataAPIClientConstructor) NewMetadataAPIClient(
	nodeID uint32,
) (metadata_api.MetadataApiClient, error) {
	conn, err := c.clientManager.GetClient(nodeID)
	if err != nil {
		return nil, err
	}
	return metadata_api.NewMetadataApiClient(conn), nil
}

type NodeCursorTracker struct {
	ctx               context.Context
	log               *zap.Logger
	metadataAPIClient MetadataAPIClientConstructor
}

func NewNodeCursorTracker(ctx context.Context,
	log *zap.Logger, metadataAPIClient MetadataAPIClientConstructor,
) *NodeCursorTracker {
	return &NodeCursorTracker{ctx: ctx, log: log, metadataAPIClient: metadataAPIClient}
}

func (ct *NodeCursorTracker) BlockUntilDesiredCursorReached(
	parentCtx context.Context,
	nodeID uint32,
	desiredOriginatorID uint32,
	desiredSequenceID uint64,
) error {
	// TODO(mkysel) ideally we wouldn't create and tear down the stream for every request

	ctx, cancel := context.WithTimeout(parentCtx, 10*time.Second)
	defer cancel()

	start := time.Now()

	client, err := ct.metadataAPIClient.NewMetadataAPIClient(nodeID)
	if err != nil {
		return err
	}
	stream, err := client.SubscribeSyncCursor(ctx, &metadata_api.GetSyncCursorRequest{})
	if err != nil {
		return err
	}

	// Channel to receive responses/errors from stream.Recv()
	respCh := make(chan *metadata_api.GetSyncCursorResponse)
	errCh := make(chan error)

	// Goroutine to read from the stream and send messages to channels
	go func() {
		defer close(respCh)
		defer close(errCh)
		for {
			resp, err := stream.Recv()
			if err != nil {
				errCh <- err
				return
			}
			respCh <- resp
		}
	}()

	for {
		select {
		case <-ct.ctx.Done():
			return status.Errorf(codes.Aborted, "node terminated. Cancelled wait for cursor")
		case <-ctx.Done():
			// client has shut down
			if parentCtx.Err() != nil {
				return nil
			}
			return status.Errorf(
				codes.DeadlineExceeded,
				"Wait for cursor was unsuccessful after %s",
				time.Since(start),
			)
		case err, ok := <-errCh:
			if !ok {
				return status.Errorf(
					codes.Internal,
					"error getting node cursor: error channel closed",
				)
			}

			if errors.Is(ctx.Err(), context.Canceled) {
				return nil
			}
			if errors.Is(ct.ctx.Err(), context.Canceled) {
				return status.Errorf(codes.Aborted, "node terminated. Cancelled wait for cursor")
			}
			return err
		case resp, ok := <-respCh:
			if !ok {
				return status.Errorf(
					codes.Internal,
					"error getting node cursor: response channel closed",
				)
			}

			if resp == nil || resp.LatestSync == nil {
				return status.Errorf(codes.Internal, "error getting node cursor: response is nil")
			}

			derefMap := resp.LatestSync.NodeIdToSequenceId
			seqID, exists := derefMap[desiredOriginatorID]
			if !exists {
				continue
			}

			if seqID >= desiredSequenceID {
				metrics.EmitPayerBlockUntilDesiredCursorReached(
					desiredOriginatorID,
					time.Since(start).Seconds(),
				)
				return nil
			}
		}
	}
}
