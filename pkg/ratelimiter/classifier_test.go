package ratelimiter

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/constants"
)

func TestClassify_Tier0WhenContextFlagSet(t *testing.T) {
	ctx := context.WithValue(context.Background(), constants.VerifiedNodeRequestCtxKey{}, true)
	require.Equal(t, Tier0, ClassifyTier(ctx))
}

func TestClassify_Tier2WhenNoContextFlag(t *testing.T) {
	require.Equal(t, Tier2, ClassifyTier(context.Background()))
}

func TestClassify_Tier2WhenContextFlagFalse(t *testing.T) {
	ctx := context.WithValue(context.Background(), constants.VerifiedNodeRequestCtxKey{}, false)
	require.Equal(t, Tier2, ClassifyTier(ctx))
}
