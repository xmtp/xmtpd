package message

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/ratelimiter"
)

type spyLimiter struct {
	allowSubject string
	allowCost    uint64
	allowResult  *ratelimiter.Result
	debitSubject string
	debitCost    uint64
}

func (s *spyLimiter) Allow(
	ctx context.Context,
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

func (s *spyLimiter) ForceDebit(
	ctx context.Context,
	subject string,
	cost uint64,
) (*ratelimiter.Result, error) {
	s.debitSubject = subject
	s.debitCost = cost
	return &ratelimiter.Result{Allowed: true}, nil
}

func TestApplySubscribeAdmissionAndDrain_ChargesAdmissionAndDrainsOnCleanup(t *testing.T) {
	limiter := &spyLimiter{}
	cfg := RateLimitConfig{
		Enabled:              true,
		DrainIntervalMinutes: 5,
		DrainAmount:          1,
	}
	cleanup, err := applySubscribeAdmissionAndDrain(
		context.Background(),
		limiter,
		cfg,
		"203.0.113.1",
		4,
	)
	require.NoError(t, err)
	require.Equal(t, uint64(2), limiter.allowCost) // ceil(sqrt(4))
	require.Equal(t, "203.0.113.1", limiter.allowSubject)

	cleanup()
	require.Equal(t, "203.0.113.1", limiter.debitSubject)
}

func TestApplySubscribeAdmissionAndDrain_DisabledIsNoOp(t *testing.T) {
	limiter := &spyLimiter{}
	cleanup, err := applySubscribeAdmissionAndDrain(
		context.Background(), limiter, RateLimitConfig{Enabled: false}, "subj", 4,
	)
	require.NoError(t, err)
	require.Empty(t, limiter.allowSubject)
	cleanup()
	require.Empty(t, limiter.debitSubject)
}

func TestApplySubscribeAdmissionAndDrain_NilLimiterIsNoOp(t *testing.T) {
	cfg := RateLimitConfig{Enabled: true, DrainIntervalMinutes: 5, DrainAmount: 1}
	cleanup, err := applySubscribeAdmissionAndDrain(context.Background(), nil, cfg, "subj", 4)
	require.NoError(t, err)
	require.NotNil(t, cleanup)
	cleanup() // must not panic
}

func TestApplySubscribeAdmissionAndDrain_DenialReturnsError(t *testing.T) {
	limiter := &spyLimiter{allowResult: &ratelimiter.Result{Allowed: false}}
	cfg := RateLimitConfig{Enabled: true, DrainIntervalMinutes: 5, DrainAmount: 1}
	_, err := applySubscribeAdmissionAndDrain(context.Background(), limiter, cfg, "subj", 4)
	require.Error(t, err)
}

// silence unused-import warnings if any:
var _ = time.Second
