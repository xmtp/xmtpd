package server

import (
	"context"
	"errors"
	"net"
	"time"

	"connectrpc.com/connect"
	"go.uber.org/zap"

	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api/message_apiconnect"
	"github.com/xmtp/xmtpd/pkg/ratelimiter"
	"github.com/xmtp/xmtpd/pkg/utils/clientip"
)

// StreamLimitInterceptor enforces a per-IP concurrent stream limit on
// NotificationApi/SubscribeAllEnvelopes. Tier 0 (authenticated nodes) bypass.
type StreamLimitInterceptor struct {
	logger          *zap.Logger
	limiter         ratelimiter.StreamLimiter
	trustedCIDRs    []*net.IPNet
	refreshInterval time.Duration
}

// NewStreamLimitInterceptor creates a StreamLimitInterceptor.
func NewStreamLimitInterceptor(
	logger *zap.Logger,
	limiter ratelimiter.StreamLimiter,
	trustedCIDRs []*net.IPNet,
	refreshInterval time.Duration,
) *StreamLimitInterceptor {
	return &StreamLimitInterceptor{
		logger:          logger.Named("xmtpd.stream-limiter"),
		limiter:         limiter,
		trustedCIDRs:    trustedCIDRs,
		refreshInterval: refreshInterval,
	}
}

// WrapUnary is a no-op — SubscribeAllEnvelopes is a streaming RPC.
func (i *StreamLimitInterceptor) WrapUnary(next connect.UnaryFunc) connect.UnaryFunc {
	return next
}

// WrapStreamingClient is a no-op for server interceptors.
func (i *StreamLimitInterceptor) WrapStreamingClient(
	next connect.StreamingClientFunc,
) connect.StreamingClientFunc {
	return next
}

// WrapStreamingHandler enforces the concurrent stream limit on
// NotificationApi/SubscribeAllEnvelopes.
func (i *StreamLimitInterceptor) WrapStreamingHandler(
	next connect.StreamingHandlerFunc,
) connect.StreamingHandlerFunc {
	return func(ctx context.Context, conn connect.StreamingHandlerConn) error {
		if conn.Spec().Procedure != message_apiconnect.NotificationApiSubscribeAllEnvelopesProcedure {
			return next(ctx, conn)
		}

		if ratelimiter.ClassifyTier(ctx) == ratelimiter.Tier0 {
			ratelimiter.StreamDecisionsTotal.WithLabelValues(
				"NotificationApi", "bypassed",
			).Inc()
			return next(ctx, conn)
		}

		subject := clientip.Extract(
			conn.Peer().Addr,
			conn.RequestHeader().Get("X-Forwarded-For"),
			i.trustedCIDRs,
		)

		allowed, err := i.limiter.Acquire(ctx, subject)
		if err != nil {
			i.logger.Warn("stream limiter error, allowing stream",
				zap.String("subject", subject),
				zap.Error(err),
			)
			ratelimiter.StreamDecisionsTotal.WithLabelValues(
				"NotificationApi", "failed_open",
			).Inc()
			return next(ctx, conn)
		}

		if !allowed {
			ratelimiter.StreamDecisionsTotal.WithLabelValues(
				"NotificationApi", "denied",
			).Inc()
			return connect.NewError(
				connect.CodeResourceExhausted,
				errors.New("concurrent stream limit exceeded"),
			)
		}

		ratelimiter.StreamDecisionsTotal.WithLabelValues(
			"NotificationApi", "allowed",
		).Inc()
		ratelimiter.StreamActiveStreams.Inc()

		// Ensure Release fires on any exit path.
		defer func() {
			ratelimiter.StreamActiveStreams.Dec()
			if releaseErr := i.limiter.Release(context.Background(), subject); releaseErr != nil {
				i.logger.Warn("stream limiter release error",
					zap.String("subject", subject),
					zap.Error(releaseErr),
				)
			}
		}()

		// Start TTL refresh goroutine.
		refreshCtx, cancelRefresh := context.WithCancel(ctx)
		defer cancelRefresh()
		go i.refreshLoop(refreshCtx, subject)

		return next(ctx, conn)
	}
}

func (i *StreamLimitInterceptor) refreshLoop(ctx context.Context, subject string) {
	ticker := time.NewTicker(i.refreshInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := i.limiter.RefreshTTL(ctx, subject); err != nil {
				i.logger.Warn("stream TTL refresh error",
					zap.String("subject", subject),
					zap.Error(err),
				)
			}
		}
	}
}
