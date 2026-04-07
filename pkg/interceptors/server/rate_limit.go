// Package server implements the server authentication interceptors.
package server

import (
	"context"
	"fmt"
	"net"
	"strings"

	"connectrpc.com/connect"
	"go.uber.org/zap"

	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api"
	"github.com/xmtp/xmtpd/pkg/ratelimiter"
)

// QueryApiMethod identifies a QueryApi RPC method.
type QueryApiMethod string

const (
	MethodQueryEnvelopes    QueryApiMethod = "QueryEnvelopes"
	MethodSubscribeTopics   QueryApiMethod = "SubscribeTopics"
	MethodGetInboxIds       QueryApiMethod = "GetInboxIds"
	MethodGetNewestEnvelope QueryApiMethod = "GetNewestEnvelope"
)

const queryApiPathPrefix = "/xmtp.xmtpv4.message_api.QueryApi/"

// QueryApiMethodFromProcedure maps a Connect procedure path to a QueryApiMethod.
// Returns ("", false) for non-QueryApi procedures.
func QueryApiMethodFromProcedure(procedure string) (QueryApiMethod, bool) {
	if !strings.HasPrefix(procedure, queryApiPathPrefix) {
		return "", false
	}
	name := QueryApiMethod(strings.TrimPrefix(procedure, queryApiPathPrefix))
	switch name {
	case MethodQueryEnvelopes, MethodSubscribeTopics, MethodGetInboxIds, MethodGetNewestEnvelope:
		return name, true
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

// RateLimitInterceptorConfig holds configuration for the rate-limit interceptor.
type RateLimitInterceptorConfig struct {
	DrainIntervalMinutes int
	DrainAmount          int
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
	cfg          RateLimitInterceptorConfig
}

var _ connect.Interceptor = (*RateLimitInterceptor)(nil)

// NewRateLimitInterceptor creates a new RateLimitInterceptor.
func NewRateLimitInterceptor(
	logger *zap.Logger,
	queryLimiter ratelimiter.RateLimiter,
	opensLimiter ratelimiter.RateLimiter,
	trustedCIDRs []*net.IPNet,
	cfg RateLimitInterceptorConfig,
) *RateLimitInterceptor {
	return &RateLimitInterceptor{
		logger:       logger.Named("xmtpd.rate-limiter"),
		queryLimiter: queryLimiter,
		opensLimiter: opensLimiter,
		trustedCIDRs: trustedCIDRs,
		cfg:          cfg,
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
			return next(ctx, req)
		}

		subject := ratelimiter.ExtractClientIP(
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
			return next(ctx, req)
		}

		if !result.Allowed {
			return nil, connect.NewError(connect.CodeResourceExhausted, fmt.Errorf("rate limit exceeded"))
		}

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
			return next(ctx, conn)
		}

		if method != MethodSubscribeTopics {
			return next(ctx, conn)
		}

		clientIP := ratelimiter.ExtractClientIP(
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
			return next(ctx, conn)
		}

		if !result.Allowed {
			return connect.NewError(connect.CodeResourceExhausted, fmt.Errorf("subscribe rate limit exceeded"))
		}

		return next(ctx, conn)
	}
}
