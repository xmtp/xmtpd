package ratelimiter_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/config"
	"github.com/xmtp/xmtpd/pkg/ratelimiter"
	"go.uber.org/zap"
)

func TestBuild_DisabledReturnsNil(t *testing.T) {
	got, err := ratelimiter.Build(
		context.Background(),
		zap.NewNop(),
		config.RedisOptions{},
		config.RateLimitOptions{Enable: false},
	)
	require.NoError(t, err)
	require.Nil(t, got)
}
