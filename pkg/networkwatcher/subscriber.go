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
	cfg           SubscriberConfig
	client        metadata_apiconnect.MetadataApiClient
	loggedInitial bool // true once the first-ever cursor has been logged
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
	attempt := 0
	for {
		if err := ctx.Err(); err != nil {
			return
		}

		attempt++
		if attempt == 1 {
			s.cfg.Logger.Info("subscribing", zap.String("url", s.cfg.BaseURL))
		} else {
			s.cfg.Logger.Info(
				"reconnecting",
				zap.String("url", s.cfg.BaseURL),
				zap.Int("attempt", attempt),
			)
		}

		updates, sessionDuration, err := s.runOnce(ctx)
		s.cfg.Aggregator.SetNodeUp(s.cfg.NodeID, false)
		if err != nil {
			nodeStreamErrors.
				WithLabelValues(nodeIDLabel(s.cfg.NodeID), classifyError(err)).
				Inc()
		}

		if ctx.Err() != nil {
			s.cfg.Logger.Info(
				"stream closed",
				zap.String("reason", "context_canceled"),
				zap.Uint64("updates", updates),
				zap.Duration("session_duration", sessionDuration),
			)
			return
		}

		// Reset backoff after a session that actually received data —
		// a dial-then-drop loop shouldn't collapse backoff, but a
		// long-lived session that eventually drops should reconnect fast.
		if updates > 0 {
			backoff = s.cfg.MinBackoff
		}

		s.cfg.Logger.Warn(
			"stream ended, will reconnect",
			zap.String("reason", classifyError(err)),
			zap.Error(err),
			zap.Uint64("updates", updates),
			zap.Duration("session_duration", sessionDuration),
			zap.Duration("backoff", backoff),
		)

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

// runOnce opens a single stream and returns when the stream ends, with
// the number of cursor updates received and the session duration.
func (s *Subscriber) runOnce(ctx context.Context) (uint64, time.Duration, error) {
	startedAt := time.Now()
	stream, err := s.client.SubscribeSyncCursor(
		ctx,
		connect.NewRequest(&metadata_api.GetSyncCursorRequest{}),
	)
	if err != nil {
		return 0, time.Since(startedAt), err
	}
	defer func() { _ = stream.Close() }()

	s.cfg.Aggregator.SetNodeUp(s.cfg.NodeID, true)

	var updates uint64
	for stream.Receive() {
		msg := stream.Msg()
		cursor := msg.GetLatestSync()
		if cursor == nil {
			continue
		}
		snapshot := cursor.GetNodeIdToSequenceId()
		if updates == 0 {
			s.cfg.Logger.Info(
				"stream connected",
				zap.Int("originators", len(snapshot)),
			)
		}
		if !s.loggedInitial {
			s.cfg.Logger.Info("initial cursor", zap.Any("cursor", snapshot))
			s.loggedInitial = true
		}
		updates++
		s.cfg.Aggregator.Apply(s.cfg.NodeID, snapshot)
	}
	return updates, time.Since(startedAt), stream.Err()
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
