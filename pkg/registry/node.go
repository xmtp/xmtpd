package registry

import (
	"context"
	"crypto/ecdsa"
	"fmt"

	"connectrpc.com/connect"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api/message_apiconnect"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/metadata_api/metadata_apiconnect"
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

func (n *Node) BuildReplicationAPIClient(
	extraDialOpts ...connect.ClientOption,
) (message_apiconnect.ReplicationApiClient, error) {
	ctx := context.Background()

	return utils.NewConnectGRPCReplicationAPIClient(ctx, n.HTTPAddress, extraDialOpts...)
}

func (n *Node) BuildMetadataAPIClient(
	extraDialOpts ...connect.ClientOption,
) (metadata_apiconnect.MetadataApiClient, error) {
	ctx := context.Background()

	target, isTLS, err := utils.HTTPAddressToGRPCTarget(n.HTTPAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to convert HTTP address to gRPC target: %w", err)
	}

	httpClient, err := utils.BuildHTTP2Client(ctx, isTLS)
	if err != nil {
		return nil, fmt.Errorf("failed to build HTTP client: %w", err)
	}

	dialOpts := utils.BuildConnectProtocolDialOptions(extraDialOpts...)

	return metadata_apiconnect.NewMetadataApiClient(
		httpClient,
		target,
		dialOpts...,
	), nil
}
