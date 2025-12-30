// Package registry implements the node registry test utils.
package registry

import (
	"crypto/ecdsa"
	"fmt"
	"testing"

	registryMocks "github.com/xmtp/xmtpd/pkg/mocks/registry"
	r "github.com/xmtp/xmtpd/pkg/registry"
)

func CreateNode(nodeID uint32, port int, privateKey *ecdsa.PrivateKey) r.Node {
	return r.Node{
		NodeID:        nodeID,
		SigningKey:    &privateKey.PublicKey,
		HTTPAddress:   fmt.Sprintf("http://localhost:%d", port),
		IsCanonical:   true,
		IsValidConfig: true,
	}
}

func CreateMockRegistry(t *testing.T, nodes []r.Node) *registryMocks.MockNodeRegistry {
	mockRegistry := registryMocks.NewMockNodeRegistry(t)
	mockRegistry.On("GetNodes").Maybe().Return(nodes, nil)

	nodesChan := make(chan []r.Node)
	mockRegistry.On("OnNewNodes").Maybe().
		Return((<-chan []r.Node)(nodesChan))

	removedNodesChan := make(chan []uint32)
	mockRegistry.On("OnRemovedNodes").Maybe().
		Return((<-chan []uint32)(removedNodesChan))

	for _, node := range nodes {
		nodeChan := make(chan r.Node)
		mockRegistry.On("OnChangedNode", node.NodeID).
			Maybe().
			Return((<-chan r.Node)(nodeChan))
		mockRegistry.On("GetNode", node.NodeID).Maybe().Return(&node, nil)
	}

	mockRegistry.On("Stop").Maybe().Return(nil)

	return mockRegistry
}
