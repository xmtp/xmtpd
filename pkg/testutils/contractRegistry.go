package testutils

import "github.com/xmtp/xmtpd/pkg/registry"

func GetHealthyNode(nodeID uint32) registry.Node {
	return registry.Node{
		NodeID:               nodeID,
		IsDisabled:           false,
		IsReplicationEnabled: true,
		IsApiEnabled:         true,
		IsValidConfig:        true,
	}
}

func GetUnhealthyNode(nodeID uint32) registry.Node {
	return registry.Node{
		NodeID:               nodeID,
		IsDisabled:           true,
		IsReplicationEnabled: true,
		IsApiEnabled:         true,
		IsValidConfig:        true,
	}
}
