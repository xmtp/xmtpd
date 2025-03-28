package registry

import (
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"testing"

	registryMocks "github.com/xmtp/xmtpd/pkg/mocks/registry"
	r "github.com/xmtp/xmtpd/pkg/registry"
)

func CreateNode(nodeID uint32, port int, privateKey *ecdsa.PrivateKey) r.Node {
	return r.Node{
		NodeID:                    nodeID,
		SigningKey:                &privateKey.PublicKey,
		HttpAddress:               fmt.Sprintf("http://localhost:%d", port),
		InCanonicalNetwork:        true,
		MinMonthlyFeeMicroDollars: big.NewInt(0),
		IsValidConfig:             true,
	}
}

func CreateMockRegistry(t *testing.T, nodes []r.Node) *registryMocks.MockNodeRegistry {
	mockRegistry := registryMocks.NewMockNodeRegistry(t)
	mockRegistry.On("GetNodes").Maybe().Return(nodes, nil)

	nodesChan := make(chan []r.Node)
	mockRegistry.On("OnNewNodes").Maybe().
		Return((<-chan []r.Node)(nodesChan), r.CancelSubscription(func() {}))

	for _, node := range nodes {
		nodeChan := make(chan r.Node)
		mockRegistry.On("OnChangedNode", node.NodeID).
			Maybe().
			Return((<-chan r.Node)(nodeChan), r.CancelSubscription(func() {
				close(nodeChan)
			}))
		mockRegistry.On("GetNode", node.NodeID).Maybe().Return(&node, nil)
	}

	mockRegistry.On("Stop").Maybe().Return(nil)

	return mockRegistry
}
