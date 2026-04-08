package clientip

import (
	"net"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExtract_NoXFFUsesPeer(t *testing.T) {
	require.Equal(t, "203.0.113.1", Extract("203.0.113.1:5001", "", nil))
}

func TestExtract_TrustedProxyPeelsOneHop(t *testing.T) {
	trusted := mustParseCIDRs(t, []string{"10.0.0.0/8"})
	require.Equal(t, "203.0.113.1",
		Extract("10.0.0.5:5001", "203.0.113.1, 10.0.0.5", trusted))
}

func TestExtract_UntrustedProxyIgnoresXFF(t *testing.T) {
	trusted := mustParseCIDRs(t, []string{"10.0.0.0/8"})
	require.Equal(t, "198.51.100.7", Extract("198.51.100.7:5001", "203.0.113.1", trusted))
}

func TestExtract_IPv6NormalizedToSlash64(t *testing.T) {
	require.Equal(t, "2001:db8:abcd:1234::/64",
		Extract("[2001:db8:abcd:1234:5678:9abc:def0:1234]:5001", "", nil))
}

func TestExtract_IPv4NotNormalized(t *testing.T) {
	require.Equal(t, "203.0.113.1", Extract("203.0.113.1:5001", "", nil))
}

func TestExtract_UnparseablePeerReturnsSentinel(t *testing.T) {
	require.Equal(t, InvalidClientIPKey, Extract("not-an-ip", "", nil))
}

func TestExtract_UnparseableXFFEntryReturnsSentinel(t *testing.T) {
	trusted := mustParseCIDRs(t, []string{"10.0.0.0/8"})
	require.Equal(t, InvalidClientIPKey, Extract("10.0.0.5:5001", "evil-string", trusted))
}

func TestParseTrustedProxyCIDRs_Empty(t *testing.T) {
	out, err := ParseTrustedProxyCIDRs("")
	require.NoError(t, err)
	require.Nil(t, out)
}

func TestParseTrustedProxyCIDRs_Valid(t *testing.T) {
	out, err := ParseTrustedProxyCIDRs("10.0.0.0/8, 192.168.0.0/16")
	require.NoError(t, err)
	require.Len(t, out, 2)
}

func TestParseTrustedProxyCIDRs_Invalid(t *testing.T) {
	_, err := ParseTrustedProxyCIDRs("10.0.0.0/8,not-a-cidr")
	require.Error(t, err)
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
