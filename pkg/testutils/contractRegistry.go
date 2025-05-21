package testutils

import (
	"fmt"

	"github.com/xmtp/xmtpd/pkg/registry"
)

func GetHealthyNode(nodeID uint32) registry.Node {
	return registry.Node{
		NodeID:        nodeID,
		IsCanonical:   true,
		IsValidConfig: true,
		HttpAddress:   fmt.Sprintf("http://localhost:%d", nodeID),
	}
}

func GetUnhealthyNode(nodeID uint32) registry.Node {
	return registry.Node{
		NodeID:        nodeID,
		IsCanonical:   false,
		IsValidConfig: false,
		HttpAddress:   fmt.Sprintf("http://localhost:%d", nodeID),
	}
}
