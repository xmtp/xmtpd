package registry

import (
	"crypto/ecdsa"
	"fmt"

	"github.com/xmtp/xmtpd/pkg/utils"
	"google.golang.org/grpc"
)

type DialOptionFunc func(node Node) []grpc.DialOption

type Node struct {
	NodeID        uint32
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

	return n.NodeID == other.NodeID &&
		n.HttpAddress == other.HttpAddress &&
		equalsSigningKey &&
		n.IsHealthy == other.IsHealthy &&
		n.IsValidConfig == other.IsValidConfig
}

func (node *Node) BuildClient(
	extraDialOpts ...grpc.DialOption,
) (*grpc.ClientConn, error) {
	target, isTLS, err := utils.HttpAddressToGrpcTarget(node.HttpAddress)
	if err != nil {
		return nil, fmt.Errorf("Failed to convert HTTP address to gRPC target: %v", err)
	}

	creds, err := utils.GetCredentialsForAddress(isTLS)
	if err != nil {
		return nil, fmt.Errorf("Failed to get credentials: %v", err)
	}

	dialOpts := append([]grpc.DialOption{
		grpc.WithTransportCredentials(creds),
		grpc.WithDefaultCallOptions(),
	}, extraDialOpts...)

	conn, err := grpc.NewClient(
		target,
		dialOpts...,
	)

	if err != nil {
		return nil, fmt.Errorf("Failed to create channel at %s: %v", target, err)
	}

	return conn, nil
}
