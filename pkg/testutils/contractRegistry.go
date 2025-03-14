package testutils

import "github.com/xmtp/xmtpd/pkg/registry"

func HealthyNode(nodeID uint32) registry.Node {
	return registry.Node{
		NodeID:               nodeID,
		IsDisabled:           false,
		IsReplicationEnabled: true,
		IsApiEnabled:         true,
		IsValidConfig:        true,
	}
}

func ApiDisabledNode(nodeID uint32) registry.Node {
	return registry.Node{
		NodeID:               nodeID,
		IsDisabled:           false,
		IsReplicationEnabled: true,
		IsApiEnabled:         false,
		IsValidConfig:        true,
	}
}
