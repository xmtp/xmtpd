// Package network implements the network test utils.
package network

import (
	"net"
	"testing"

	"github.com/stretchr/testify/require"
)

func OpenListener(t *testing.T) net.Listener {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)

	t.Cleanup(func() {
		_ = ln.Close()
	})

	return ln
}

// OpenFreePort opens a free port on the local machine.
// It should be used only in tests, where an http server
// needs a port, but unlike grP
func OpenFreePort(t *testing.T) int {
	l, err := net.Listen("tcp", "localhost:0")
	require.NoError(t, err)

	defer func() {
		_ = l.Close()
	}()

	return l.Addr().(*net.TCPAddr).Port
}
