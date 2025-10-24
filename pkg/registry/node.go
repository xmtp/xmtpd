package registry

import (
	"crypto/ecdsa"

	"github.com/xmtp/xmtpd/pkg/utils"
	"google.golang.org/grpc"
)

type DialOptionFunc func(node Node) []grpc.DialOption

type Node struct {
	NodeID        uint32
	SigningKey    *ecdsa.PublicKey
	HTTPAddress   string
	IsCanonical   bool
	IsValidConfig bool
}

func (n *Node) Equals(other Node) bool {
	var equalsSigningKey bool
	if n.SigningKey == nil && other.SigningKey == nil {
		equalsSigningKey = true
	} else if n.SigningKey != nil && other.SigningKey != nil {
		equalsSigningKey = n.SigningKey.Equal(other.SigningKey)
	}

	return n.NodeID == other.NodeID &&
		n.HTTPAddress == other.HTTPAddress &&
		equalsSigningKey &&
		n.IsCanonical == other.IsCanonical &&
		n.IsValidConfig == other.IsValidConfig
}

func (n *Node) BuildClient(
	extraDialOpts ...grpc.DialOption,
) (*grpc.ClientConn, error) {
	_, conn, err := utils.NewGRPCReplicationAPIClientAndConn(n.HTTPAddress)
	return conn, err
}
