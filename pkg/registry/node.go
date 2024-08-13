package registry

import "crypto/ecdsa"

type Node struct {
	NodeId        uint16
	SigningKey    *ecdsa.PublicKey
	HttpAddress   string
	IsHealthy     bool
	IsValidConfig bool
}

func (n *Node) Equals(other Node) bool {
	var equalsSigningKey bool
	if n.SigningKey == nil && other.SigningKey == nil {
		equalsSigningKey = true
	} else if n.SigningKey != nil && other.SigningKey != nil {
		equalsSigningKey = n.SigningKey.Equal(other.SigningKey)
	}

	return n.NodeId == other.NodeId &&
		n.HttpAddress == other.HttpAddress &&
		equalsSigningKey &&
		n.IsHealthy == other.IsHealthy &&
		n.IsValidConfig == other.IsValidConfig
}
