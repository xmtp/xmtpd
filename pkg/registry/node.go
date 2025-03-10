package registry

import (
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"time"

	"github.com/xmtp/xmtpd/pkg/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

type DialOptionFunc func(node Node) []grpc.DialOption

type Node struct {
	NodeID                    uint32
	SigningKey                *ecdsa.PublicKey
	HttpAddress               string
	IsReplicationEnabled      bool
	IsApiEnabled              bool
	IsDisabled                bool
	MinMonthlyFeeMicroDollars *big.Int
	IsValidConfig             bool
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
		n.IsReplicationEnabled == other.IsReplicationEnabled &&
		n.IsApiEnabled == other.IsApiEnabled &&
		n.IsDisabled == other.IsDisabled &&
		n.MinMonthlyFeeMicroDollars.Cmp(other.MinMonthlyFeeMicroDollars) == 0 &&
		n.IsValidConfig == other.IsValidConfig
}

func (node *Node) BuildClient(
	extraDialOpts ...grpc.DialOption,
) (*grpc.ClientConn, error) {
	target, isTLS, err := utils.HttpAddressToGrpcTarget(node.HttpAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to convert HTTP address to gRPC target: %v", err)
	}

	creds, err := utils.GetCredentialsForAddress(isTLS)
	if err != nil {
		return nil, fmt.Errorf("failed to get credentials: %v", err)
	}

	dialOpts := append([]grpc.DialOption{
		grpc.WithTransportCredentials(creds),
		grpc.WithDefaultCallOptions(),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                30 * time.Second,
			Timeout:             10 * time.Second,
			PermitWithoutStream: true,
		}),
	}, extraDialOpts...)

	conn, err := grpc.NewClient(
		target,
		dialOpts...,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create channel at %s: %v", target, err)
	}

	return conn, nil
}
