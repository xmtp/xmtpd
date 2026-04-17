package networkwatcher

import (
	"context"
	"errors"
	"net/http"
	"time"

	"connectrpc.com/connect"
	"go.uber.org/zap"

	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/metadata_api"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/metadata_api/metadata_apiconnect"
)

// SubscriberConfig configures a per-node Subscriber.
type SubscriberConfig struct {
	NodeID     uint32
	BaseURL    string
	Aggregator *Aggregator
	Logger     *zap.Logger
	HTTPClient connect.HTTPClient

	MinBackoff time.Duration
	MaxBackoff time.Duration
}

// Subscriber maintains a single long-lived SubscribeSyncCursor stream to one
// publisher node, reconnecting with exponential backoff on failure.
type Subscriber struct {
	cfg    SubscriberConfig
	client metadata_apiconnect.MetadataApiClient
}

// NewSubscriber builds a Subscriber with a Connect client bound to cfg.BaseURL.
func NewSubscriber(cfg SubscriberConfig) *Subscriber {
	httpClient := cfg.HTTPClient
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	client := metadata_apiconnect.NewMetadataApiClient(httpClient, cfg.BaseURL, connect.WithGRPC())
	return &Subscriber{cfg: cfg, client: client}
}

// Run blocks until ctx is canceled. It opens the subscribe stream, forwards
// each cursor snapshot to the aggregator, and reconnects on error.
func (s *Subscriber) Run(ctx context.Context) {
	backoff := s.cfg.MinBackoff
	for {
		if err := ctx.Err(); err != nil {
			return
		}

		err := s.runOnce(ctx)
		reason := classifyError(err)
		s.cfg.Aggregator.SetNodeUp(s.cfg.NodeID, false)
		nodeStreamErrors.
			WithLabelValues(nodeIDLabel(s.cfg.NodeID), reason).
			Inc()
		s.cfg.Logger.Debug(
			"subscribe stream ended",
			zap.Uint32("node_id", s.cfg.NodeID),
			zap.String("reason", reason),
			zap.Error(err),
		)

		if ctx.Err() != nil {
			return
		}

		select {
		case <-ctx.Done():
			return
		case <-time.After(backoff):
		}

		backoff *= 2
		if backoff > s.cfg.MaxBackoff {
			backoff = s.cfg.MaxBackoff
		}
	}
}

// runOnce opens a single stream and returns when the stream ends.
func (s *Subscriber) runOnce(ctx context.Context) error {
	stream, err := s.client.SubscribeSyncCursor(
		ctx,
		connect.NewRequest(&metadata_api.GetSyncCursorRequest{}),
	)
	if err != nil {
		return err
	}
	defer func() { _ = stream.Close() }()

	s.cfg.Aggregator.SetNodeUp(s.cfg.NodeID, true)

	for stream.Receive() {
		msg := stream.Msg()
		cursor := msg.GetLatestSync()
		if cursor == nil {
			continue
		}
		s.cfg.Aggregator.Apply(s.cfg.NodeID, cursor.GetNodeIdToSequenceId())
	}
	return stream.Err()
}

func classifyError(err error) string {
	if err == nil {
		return "stream_recv"
	}
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return "context_canceled"
	}
	var connectErr *connect.Error
	if errors.As(err, &connectErr) {
		if connectErr.Code() == connect.CodeUnavailable {
			return "dial"
		}
		return "stream_recv"
	}
	return "other"
}
