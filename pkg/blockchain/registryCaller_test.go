package blockchain_test

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetAllNodes(t *testing.T) {
	registry, caller, ctx := buildRegistry(t)

	addRandomNode(t, registry, ctx)

	nodes, err := caller.GetAllNodes(ctx)
	require.NoError(t, err)
	require.Len(t, nodes, 1)
}

func TestGetNode(t *testing.T) {
	registry, caller, ctx := buildRegistry(t)

	addRandomNode(t, registry, ctx)

	node, err := caller.GetNode(ctx, 100)
	require.NoError(t, err)
	require.NotNil(t, node)
}

func TestGetNodeNotFound(t *testing.T) {
	_, caller, ctx := buildRegistry(t)

	_, err := caller.GetNode(ctx, 100)
	require.Error(t, err)
}

func TestOwnerOf(t *testing.T) {
	registry, caller, ctx := buildRegistry(t)

	addRandomNode(t, registry, ctx)

	_, err := caller.OwnerOf(ctx, 100)
	require.NoError(t, err)
}
