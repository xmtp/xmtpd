// Package network implements the network test utils.
package network

import (
	"net"
	"testing"

	"github.com/stretchr/testify/require"
)

func OpenFreePort(t *testing.T) net.Listener {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = ln.Close()
	})
	return ln
}
