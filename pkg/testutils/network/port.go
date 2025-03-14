package network

import (
	"net"
	"testing"

	"github.com/stretchr/testify/require"
)

func FindFreePort(t *testing.T) int {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	defer ln.Close()
	return ln.Addr().(*net.TCPAddr).Port
}
