package testutils

import "github.com/xmtp/xmtpd/pkg/registry"

func GetHealthyNode(nodeID uint32) registry.Node {
	return registry.Node{
		NodeID:             nodeID,
		InCanonicalNetwork: true,
		IsValidConfig:      true,
	}
}

func GetUnhealthyNode(nodeID uint32) registry.Node {
	return registry.Node{
		NodeID:             nodeID,
		InCanonicalNetwork: false,
		IsValidConfig:      false,
	}
}
