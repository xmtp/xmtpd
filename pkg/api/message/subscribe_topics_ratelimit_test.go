package message

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/ratelimiter"
)

type spyLimiter struct {
	allowSubject string
	allowCost    uint64
	allowResult  *ratelimiter.Result
}

func (s *spyLimiter) Allow(
	_ context.Context,
	subject string,
	cost uint64,
) (*ratelimiter.Result, error) {
	s.allowSubject = subject
	s.allowCost = cost
	if s.allowResult == nil {
		return &ratelimiter.Result{Allowed: true}, nil
	}
	return s.allowResult, nil
}

func TestApplySubscribeAdmission_ChargesCeilSqrtFilters(t *testing.T) {
	limiter := &spyLimiter{}
	err := applySubscribeAdmission(
		context.Background(),
		limiter,
		RateLimitConfig{Enabled: true},
		"203.0.113.1",
		4,
	)
	require.NoError(t, err)
	require.Equal(t, uint64(2), limiter.allowCost) // ceil(sqrt(4))
	require.Equal(t, "203.0.113.1", limiter.allowSubject)
}

func TestApplySubscribeAdmission_DisabledIsNoOp(t *testing.T) {
	limiter := &spyLimiter{}
	err := applySubscribeAdmission(
		context.Background(), limiter, RateLimitConfig{Enabled: false}, "subj", 4,
	)
	require.NoError(t, err)
	require.Empty(t, limiter.allowSubject)
}

func TestApplySubscribeAdmission_NilLimiterIsNoOp(t *testing.T) {
	err := applySubscribeAdmission(
		context.Background(), nil, RateLimitConfig{Enabled: true}, "subj", 4,
	)
	require.NoError(t, err)
}

func TestApplySubscribeAdmission_DenialReturnsError(t *testing.T) {
	limiter := &spyLimiter{allowResult: &ratelimiter.Result{Allowed: false}}
	err := applySubscribeAdmission(
		context.Background(), limiter, RateLimitConfig{Enabled: true}, "subj", 4,
	)
	require.Error(t, err)
}
