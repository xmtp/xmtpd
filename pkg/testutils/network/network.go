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
// needs a port, but unlike gRPC, cannot reuse a listener.
func OpenFreePort(t *testing.T) int {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)

	port := ln.Addr().(*net.TCPAddr).Port

	err = ln.Close()
	require.NoError(t, err)

	return port
}
