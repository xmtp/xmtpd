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

type MetadataApiClientConstructor interface {
	NewMetadataApiClient(nodeId uint32) (metadata_api.MetadataApiClient, error)
}
type DefaultMetadataApiClientConstructor struct {
	clientManager *ClientManager
}

func (c *DefaultMetadataApiClientConstructor) NewMetadataApiClient(
	nodeId uint32,
) (metadata_api.MetadataApiClient, error) {
	conn, err := c.clientManager.GetClient(nodeId)
	if err != nil {
		return nil, err
	}
	return metadata_api.NewMetadataApiClient(conn), nil
}

type NodeCursorTracker struct {
	ctx               context.Context
	log               *zap.Logger
	metadataApiClient MetadataApiClientConstructor
}

func NewNodeCursorTracker(ctx context.Context,
	log *zap.Logger, metadataApiClient MetadataApiClientConstructor,
) *NodeCursorTracker {
	return &NodeCursorTracker{ctx: ctx, log: log, metadataApiClient: metadataApiClient}
}

func (ct *NodeCursorTracker) BlockUntilDesiredCursorReached(
	ctx context.Context,
	nodeId uint32,
	desiredOriginatorId uint32,
	desiredSequenceId uint64,
) error {
	// TODO(mkysel) ideally we wouldn't create and tear down the stream for every request

	start := time.Now()

	client, err := ct.metadataApiClient.NewMetadataApiClient(nodeId)
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
			return status.Errorf(codes.Internal, "node terminated. Cancelled wait for cursor")
		case <-ctx.Done():
			return nil
		case err := <-errCh:
			if errors.Is(ctx.Err(), context.Canceled) {
				return nil
			}
			if errors.Is(ct.ctx.Err(), context.Canceled) {
				return status.Errorf(codes.Internal, "node terminated. Cancelled wait for cursor")
			}
			return err
		case resp := <-respCh:
			if resp == nil || resp.LatestSync == nil {
				return status.Errorf(codes.Internal, "error getting node cursor: response is nil")
			}

			derefMap := resp.LatestSync.NodeIdToSequenceId
			seqId, exists := derefMap[desiredOriginatorId]
			if !exists {
				continue
			}

			if seqId >= desiredSequenceId {
				metrics.EmitPayerBlockUntilDesiredCursorReached(
					desiredOriginatorId,
					time.Since(start).Seconds(),
				)
				return nil
			}
		}
	}
}
