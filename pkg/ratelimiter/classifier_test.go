package ratelimiter

import (
	"context"
	"net"
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

func TestExtractClientIP_NoXFFUsesPeer(t *testing.T) {
	got := ExtractClientIP("203.0.113.1:5001", "", nil)
	require.Equal(t, "203.0.113.1", got)
}

func TestExtractClientIP_TrustedProxyPeelsOneHop(t *testing.T) {
	trusted := mustParseCIDRs(t, []string{"10.0.0.0/8"})
	got := ExtractClientIP("10.0.0.5:5001", "203.0.113.1, 10.0.0.5", trusted)
	require.Equal(t, "203.0.113.1", got)
}

func TestExtractClientIP_UntrustedProxyIgnoresXFF(t *testing.T) {
	trusted := mustParseCIDRs(t, []string{"10.0.0.0/8"})
	got := ExtractClientIP("198.51.100.7:5001", "203.0.113.1", trusted)
	require.Equal(t, "198.51.100.7", got)
}

func TestExtractClientIP_IPv6NormalizedToSlash64(t *testing.T) {
	got := ExtractClientIP("[2001:db8:abcd:1234:5678:9abc:def0:1234]:5001", "", nil)
	require.Equal(t, "2001:db8:abcd:1234::/64", got)
}

func TestExtractClientIP_IPv4NotNormalized(t *testing.T) {
	got := ExtractClientIP("203.0.113.1:5001", "", nil)
	require.Equal(t, "203.0.113.1", got)
}

func mustParseCIDRs(t *testing.T, cidrs []string) []*net.IPNet {
	t.Helper()
	out := make([]*net.IPNet, 0, len(cidrs))
	for _, c := range cidrs {
		_, n, err := net.ParseCIDR(c)
		require.NoError(t, err)
		out = append(out, n)
	}
	return out
}
