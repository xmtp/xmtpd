// Package server implements the server authentication interceptors.
package server

import (
	"context"
	"errors"
	"net"

	"connectrpc.com/connect"
	"go.uber.org/zap"

	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api/message_apiconnect"
	"github.com/xmtp/xmtpd/pkg/ratelimiter"
	"github.com/xmtp/xmtpd/pkg/utils/clientip"
)

// QueryApiMethod identifies a QueryApi RPC method. The string value is the
// short method name used as a label on rate-limit metrics.
type QueryApiMethod string

const (
	MethodQueryEnvelopes    QueryApiMethod = "QueryEnvelopes"
	MethodSubscribeTopics   QueryApiMethod = "SubscribeTopics"
	MethodGetInboxIds       QueryApiMethod = "GetInboxIds"
	MethodGetNewestEnvelope QueryApiMethod = "GetNewestEnvelope"
)

// QueryApiMethodFromProcedure maps a Connect procedure path to a QueryApiMethod.
// Returns ("", false) for non-QueryApi procedures.
//
// Procedure strings are compared against the generated
// message_apiconnect.QueryApi*Procedure constants so that a rename of the
// proto path is a compile-time break here rather than a silent loss of rate
// limiting on the renamed procedure.
func QueryApiMethodFromProcedure(procedure string) (QueryApiMethod, bool) {
	switch procedure {
	case message_apiconnect.QueryApiQueryEnvelopesProcedure:
		return MethodQueryEnvelopes, true
	case message_apiconnect.QueryApiSubscribeTopicsProcedure:
		return MethodSubscribeTopics, true
	case message_apiconnect.QueryApiGetInboxIdsProcedure:
		return MethodGetInboxIds, true
	case message_apiconnect.QueryApiGetNewestEnvelopeProcedure:
		return MethodGetNewestEnvelope, true
	}
	return "", false
}

// ComputeCost returns the token cost for a QueryApi request.
func ComputeCost(method QueryApiMethod, req any) uint64 {
	switch method {
	case MethodQueryEnvelopes:
		if r, ok := req.(*message_api.QueryEnvelopesRequest); ok {
			return ratelimiter.CostQuery(len(r.GetQuery().GetTopics()))
		}
		return 1
	case MethodSubscribeTopics:
		if r, ok := req.(*message_api.SubscribeTopicsRequest); ok {
			return ratelimiter.CostQuery(len(r.GetFilters()))
		}
		return 1
	default:
		return 1
	}
}

// RateLimitInterceptor enforces rate limits on QueryApi requests.
// It uses two separate RateLimiter instances:
//   - queryLimiter: per-IP, per-minute/per-hour buckets for all unary QueryApi methods.
//   - opensLimiter: per-IP, opens-per-minute sub-limit for SubscribeTopics streaming opens.
type RateLimitInterceptor struct {
	logger       *zap.Logger
	queryLimiter ratelimiter.RateLimiter
	opensLimiter ratelimiter.RateLimiter
	trustedCIDRs []*net.IPNet
}

var _ connect.Interceptor = (*RateLimitInterceptor)(nil)

// NewRateLimitInterceptor creates a new RateLimitInterceptor.
func NewRateLimitInterceptor(
	logger *zap.Logger,
	queryLimiter ratelimiter.RateLimiter,
	opensLimiter ratelimiter.RateLimiter,
	trustedCIDRs []*net.IPNet,
) *RateLimitInterceptor {
	return &RateLimitInterceptor{
		logger:       logger.Named("xmtpd.rate-limiter"),
		queryLimiter: queryLimiter,
		opensLimiter: opensLimiter,
		trustedCIDRs: trustedCIDRs,
	}
}

// WrapUnary applies rate limiting to unary QueryApi requests.
func (i *RateLimitInterceptor) WrapUnary(next connect.UnaryFunc) connect.UnaryFunc {
	return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		method, ok := QueryApiMethodFromProcedure(req.Spec().Procedure)
		if !ok {
			return next(ctx, req)
		}

		if ratelimiter.ClassifyTier(ctx) == ratelimiter.Tier0 {
			ratelimiter.DecisionsTotal.WithLabelValues(
				"QueryApi", string(method), "tier0", "bypassed",
			).Inc()
			return next(ctx, req)
		}

		subject := clientip.Extract(
			req.Peer().Addr,
			req.Header().Get("X-Forwarded-For"),
			i.trustedCIDRs,
		)

		cost := ComputeCost(method, req.Any())

		result, err := i.queryLimiter.Allow(ctx, subject, cost)
		if err != nil {
			i.logger.Warn("rate limiter error, allowing request",
				zap.String("procedure", req.Spec().Procedure),
				zap.String("subject", subject),
				zap.Error(err),
			)
			ratelimiter.DecisionsTotal.WithLabelValues(
				"QueryApi", string(method), "tier2", "failed_open",
			).Inc()
			return next(ctx, req)
		}

		if !result.Allowed {
			ratelimiter.DecisionsTotal.WithLabelValues(
				"QueryApi", string(method), "tier2", "denied",
			).Inc()
			return nil, connect.NewError(
				connect.CodeResourceExhausted,
				errors.New("rate limit exceeded"),
			)
		}

		ratelimiter.DecisionsTotal.WithLabelValues(
			"QueryApi", string(method), "tier2", "allowed",
		).Inc()
		return next(ctx, req)
	}
}

// WrapStreamingClient is a no-op for server interceptors.
// It's only implemented to satisfy the connect.Interceptor interface.
// This method is never called on the server side.
func (i *RateLimitInterceptor) WrapStreamingClient(
	next connect.StreamingClientFunc,
) connect.StreamingClientFunc {
	return next
}

// WrapStreamingHandler applies rate limiting to the SubscribeTopics streaming opens.
func (i *RateLimitInterceptor) WrapStreamingHandler(
	next connect.StreamingHandlerFunc,
) connect.StreamingHandlerFunc {
	return func(ctx context.Context, conn connect.StreamingHandlerConn) error {
		method, ok := QueryApiMethodFromProcedure(conn.Spec().Procedure)
		if !ok {
			return next(ctx, conn)
		}

		if ratelimiter.ClassifyTier(ctx) == ratelimiter.Tier0 {
			ratelimiter.DecisionsTotal.WithLabelValues(
				"QueryApi", string(method), "tier0", "bypassed",
			).Inc()
			return next(ctx, conn)
		}

		if method != MethodSubscribeTopics {
			return next(ctx, conn)
		}

		clientIP := clientip.Extract(
			conn.Peer().Addr,
			conn.RequestHeader().Get("X-Forwarded-For"),
			i.trustedCIDRs,
		)
		opensSubject := clientIP + ":opens"

		result, err := i.opensLimiter.Allow(ctx, opensSubject, 1)
		if err != nil {
			i.logger.Warn("subscribe opens rate limiter error, allowing request",
				zap.String("procedure", conn.Spec().Procedure),
				zap.String("subject", opensSubject),
				zap.Error(err),
			)
			ratelimiter.DecisionsTotal.WithLabelValues(
				"QueryApi", "SubscribeTopics", "tier2", "failed_open",
			).Inc()
			return next(ctx, conn)
		}

		if !result.Allowed {
			ratelimiter.DecisionsTotal.WithLabelValues(
				"QueryApi", "SubscribeTopics", "tier2", "denied",
			).Inc()
			return connect.NewError(
				connect.CodeResourceExhausted,
				errors.New("subscribe rate limit exceeded"),
			)
		}

		ratelimiter.DecisionsTotal.WithLabelValues(
			"QueryApi", "SubscribeTopics", "tier2", "allowed",
		).Inc()
		return next(ctx, conn)
	}
}
