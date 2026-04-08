package server

import (
	"context"
	"net/http"
	"testing"

	"connectrpc.com/connect"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/constants"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api"
	"github.com/xmtp/xmtpd/pkg/ratelimiter"
	"go.uber.org/zap/zaptest"
)

// fakeLimiter is a test double for ratelimiter.RateLimiter.
type fakeLimiter struct {
	lastSubject string
	lastCost    uint64
	result      *ratelimiter.Result
	err         error
}

func (f *fakeLimiter) Allow(
	_ context.Context,
	subject string,
	cost uint64,
) (*ratelimiter.Result, error) {
	f.lastSubject = subject
	f.lastCost = cost
	return f.result, f.err
}

func (f *fakeLimiter) ForceDebit(
	_ context.Context,
	_ string,
	_ uint64,
) (*ratelimiter.Result, error) {
	return &ratelimiter.Result{Allowed: true}, nil
}

func newTestInterceptor(limiter *fakeLimiter) *RateLimitInterceptor {
	logger := zaptest.NewLogger(&testing.T{})
	return NewRateLimitInterceptor(logger, limiter, limiter, nil, RateLimitInterceptorConfig{})
}

// mockConnectRequest satisfies connect.AnyRequest for testing.
type mockConnectRequest struct {
	connect.AnyRequest
	header http.Header
	peer   connect.Peer
	spec   connect.Spec
	body   any
}

func (m *mockConnectRequest) Header() http.Header { return m.header }
func (m *mockConnectRequest) Peer() connect.Peer  { return m.peer }
func (m *mockConnectRequest) Spec() connect.Spec  { return m.spec }
func (m *mockConnectRequest) Any() any            { return m.body }

// mockStreamingConn satisfies connect.StreamingHandlerConn for testing.
type mockStreamingConn struct {
	connect.StreamingHandlerConn
	header http.Header
	peer   connect.Peer
	spec   connect.Spec
}

func (m *mockStreamingConn) RequestHeader() http.Header { return m.header }
func (m *mockStreamingConn) Peer() connect.Peer         { return m.peer }
func (m *mockStreamingConn) Spec() connect.Spec         { return m.spec }

// --- Routing helper tests ---

func TestQueryApiMethod_FromProcedure(t *testing.T) {
	tests := []struct {
		name       string
		procedure  string
		wantOk     bool
		wantMethod QueryApiMethod
	}{
		{
			name:       "QueryEnvelopes",
			procedure:  "/xmtp.xmtpv4.message_api.QueryApi/QueryEnvelopes",
			wantOk:     true,
			wantMethod: MethodQueryEnvelopes,
		},
		{
			name:       "SubscribeTopics",
			procedure:  "/xmtp.xmtpv4.message_api.QueryApi/SubscribeTopics",
			wantOk:     true,
			wantMethod: MethodSubscribeTopics,
		},
		{
			name:       "GetInboxIds",
			procedure:  "/xmtp.xmtpv4.message_api.QueryApi/GetInboxIds",
			wantOk:     true,
			wantMethod: MethodGetInboxIds,
		},
		{
			name:       "GetNewestEnvelope",
			procedure:  "/xmtp.xmtpv4.message_api.QueryApi/GetNewestEnvelope",
			wantOk:     true,
			wantMethod: MethodGetNewestEnvelope,
		},
		{
			name:       "PublishApi is not QueryApi",
			procedure:  "/xmtp.xmtpv4.message_api.PublishApi/PublishEnvelopes",
			wantOk:     false,
			wantMethod: "",
		},
		{
			name:       "empty procedure",
			procedure:  "",
			wantOk:     false,
			wantMethod: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := QueryApiMethodFromProcedure(tt.procedure)
			assert.Equal(t, tt.wantOk, ok)
			assert.Equal(t, tt.wantMethod, got)
		})
	}
}

// --- Cost helper tests ---

func TestComputeCost_QueryEnvelopes_FromTopics(t *testing.T) {
	req := &message_api.QueryEnvelopesRequest{
		Query: &message_api.EnvelopesQuery{
			Topics: make([][]byte, 100),
		},
	}
	// CostQuery(100) = ceil(sqrt(100)) = 10
	cost := ComputeCost(MethodQueryEnvelopes, req)
	assert.Equal(t, uint64(10), cost)
}

func TestComputeCost_SubscribeTopics_FromFilters(t *testing.T) {
	req := &message_api.SubscribeTopicsRequest{
		Filters: make([]*message_api.SubscribeTopicsRequest_TopicFilter, 4),
	}
	// CostQuery(4) = ceil(sqrt(4)) = 2
	cost := ComputeCost(MethodSubscribeTopics, req)
	assert.Equal(t, uint64(2), cost)
}

func TestComputeCost_GetInboxIds_Constant(t *testing.T) {
	cost := ComputeCost(MethodGetInboxIds, &message_api.GetInboxIdsRequest{})
	assert.Equal(t, uint64(1), cost)
}

func TestComputeCost_GetNewestEnvelope_Constant(t *testing.T) {
	cost := ComputeCost(MethodGetNewestEnvelope, &message_api.GetNewestEnvelopeRequest{})
	assert.Equal(t, uint64(1), cost)
}

// --- Unary interceptor tests ---

func TestRateLimitInterceptor_BypassesNonQueryApi(t *testing.T) {
	limiter := &fakeLimiter{
		result: &ratelimiter.Result{Allowed: false},
	}
	interceptor := newTestInterceptor(limiter)

	req := &mockConnectRequest{
		header: http.Header{},
		peer:   connect.Peer{Addr: "1.2.3.4:5678"},
		spec:   connect.Spec{Procedure: "/xmtp.xmtpv4.message_api.PublishApi/PublishEnvelopes"},
		body:   nil,
	}

	handlerCalled := false
	next := func(ctx context.Context, r connect.AnyRequest) (connect.AnyResponse, error) {
		handlerCalled = true
		return nil, nil
	}

	wrappedUnary := interceptor.WrapUnary(next)
	_, err := wrappedUnary(context.Background(), req)

	require.NoError(t, err)
	assert.True(t, handlerCalled, "handler should be called for non-QueryApi procedures")
	// Limiter should not have been called.
	assert.Empty(t, limiter.lastSubject)
}

func TestRateLimitInterceptor_Tier0BypassesLimiter(t *testing.T) {
	limiter := &fakeLimiter{
		result: &ratelimiter.Result{Allowed: false},
	}
	interceptor := newTestInterceptor(limiter)

	req := &mockConnectRequest{
		header: http.Header{},
		peer:   connect.Peer{Addr: "1.2.3.4:5678"},
		spec:   connect.Spec{Procedure: "/xmtp.xmtpv4.message_api.QueryApi/QueryEnvelopes"},
		body:   &message_api.QueryEnvelopesRequest{},
	}

	handlerCalled := false
	next := func(ctx context.Context, r connect.AnyRequest) (connect.AnyResponse, error) {
		handlerCalled = true
		return nil, nil
	}

	// Tier 0: inject the verified-node context flag.
	ctx := context.WithValue(context.Background(), constants.VerifiedNodeRequestCtxKey{}, true)

	wrappedUnary := interceptor.WrapUnary(next)
	_, err := wrappedUnary(ctx, req)

	require.NoError(t, err)
	assert.True(t, handlerCalled, "Tier0 handler should bypass rate limiting")
	assert.Empty(t, limiter.lastSubject, "limiter should not be called for Tier0")
}

func TestRateLimitInterceptor_Tier2DeniedReturnsResourceExhausted(t *testing.T) {
	limiter := &fakeLimiter{
		result: &ratelimiter.Result{Allowed: false},
	}
	interceptor := newTestInterceptor(limiter)

	req := &mockConnectRequest{
		header: http.Header{},
		peer:   connect.Peer{Addr: "1.2.3.4:5678"},
		spec:   connect.Spec{Procedure: "/xmtp.xmtpv4.message_api.QueryApi/GetInboxIds"},
		body:   &message_api.GetInboxIdsRequest{},
	}

	handlerCalled := false
	next := func(ctx context.Context, r connect.AnyRequest) (connect.AnyResponse, error) {
		handlerCalled = true
		return nil, nil
	}

	wrappedUnary := interceptor.WrapUnary(next)
	_, err := wrappedUnary(context.Background(), req)

	require.Error(t, err)
	assert.False(t, handlerCalled, "handler should not be called when rate limited")

	var connectErr *connect.Error
	require.ErrorAs(t, err, &connectErr)
	assert.Equal(t, connect.CodeResourceExhausted, connectErr.Code())
}

// --- Streaming interceptor tests ---

func TestRateLimitInterceptor_Streaming_BypassesNonQueryApi(t *testing.T) {
	limiter := &fakeLimiter{
		result: &ratelimiter.Result{Allowed: false},
	}
	interceptor := newTestInterceptor(limiter)

	conn := &mockStreamingConn{
		header: http.Header{},
		peer:   connect.Peer{Addr: "1.2.3.4:5678"},
		spec:   connect.Spec{Procedure: "/xmtp.xmtpv4.message_api.PublishApi/PublishEnvelopes"},
	}

	handlerCalled := false
	next := func(ctx context.Context, c connect.StreamingHandlerConn) error {
		handlerCalled = true
		return nil
	}

	wrappedStream := interceptor.WrapStreamingHandler(next)
	err := wrappedStream(context.Background(), conn)

	require.NoError(t, err)
	assert.True(t, handlerCalled, "non-QueryApi streaming should pass through")
	assert.Empty(t, limiter.lastSubject, "limiter should not be called for non-QueryApi")
}

func TestRateLimitInterceptor_Streaming_Tier0Bypass(t *testing.T) {
	limiter := &fakeLimiter{
		result: &ratelimiter.Result{Allowed: false},
	}
	interceptor := newTestInterceptor(limiter)

	conn := &mockStreamingConn{
		header: http.Header{},
		peer:   connect.Peer{Addr: "1.2.3.4:5678"},
		spec:   connect.Spec{Procedure: "/xmtp.xmtpv4.message_api.QueryApi/SubscribeTopics"},
	}

	handlerCalled := false
	next := func(ctx context.Context, c connect.StreamingHandlerConn) error {
		handlerCalled = true
		return nil
	}

	// Tier 0: inject the verified-node context flag.
	ctx := context.WithValue(context.Background(), constants.VerifiedNodeRequestCtxKey{}, true)

	wrappedStream := interceptor.WrapStreamingHandler(next)
	err := wrappedStream(ctx, conn)

	require.NoError(t, err)
	assert.True(t, handlerCalled, "Tier0 streaming should bypass rate limiting")
	assert.Empty(t, limiter.lastSubject, "limiter should not be called for Tier0")
}

func TestRateLimitInterceptor_Streaming_OpensSubLimit_Denied(t *testing.T) {
	limiter := &fakeLimiter{
		result: &ratelimiter.Result{Allowed: false},
	}
	interceptor := newTestInterceptor(limiter)

	conn := &mockStreamingConn{
		header: http.Header{},
		peer:   connect.Peer{Addr: "1.2.3.4:5678"},
		spec:   connect.Spec{Procedure: "/xmtp.xmtpv4.message_api.QueryApi/SubscribeTopics"},
	}

	handlerCalled := false
	next := func(ctx context.Context, c connect.StreamingHandlerConn) error {
		handlerCalled = true
		return nil
	}

	wrappedStream := interceptor.WrapStreamingHandler(next)
	err := wrappedStream(context.Background(), conn)

	require.Error(t, err)
	assert.False(t, handlerCalled, "handler should not be called when subscribe rate limited")

	var connectErr *connect.Error
	require.ErrorAs(t, err, &connectErr)
	assert.Equal(t, connect.CodeResourceExhausted, connectErr.Code())

	// Verify the opens subject suffix is used.
	assert.Equal(t, "1.2.3.4:opens", limiter.lastSubject)
}
