package blockchain

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetAllNodes(t *testing.T) {
	registry, caller, ctx, cleanup := buildRegistry(t)
	defer cleanup()

	addRandomNode(t, registry, ctx)

	nodes, err := caller.GetAllNodes(ctx)
	require.NoError(t, err)
	require.NotNil(t, nodes)
}
