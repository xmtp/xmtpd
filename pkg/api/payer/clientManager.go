package payer

import (
	"sync"

	"github.com/xmtp/xmtpd/pkg/registry"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

/*
*
The ClientManager contains a mapping of nodeIDs to gRPC client connections.

These client connections are safe to be shared and re-used and will automatically attempt
to reconnect if the underlying socket connection is lost.
*
*/
type ClientManager struct {
	log          *zap.Logger
	connections  sync.Map // map[uint32]*grpc.ClientConn
	nodeRegistry registry.NodeRegistry
}

func NewClientManager(log *zap.Logger, nodeRegistry registry.NodeRegistry) *ClientManager {
	return &ClientManager{log: log, nodeRegistry: nodeRegistry}
}

func (c *ClientManager) GetClient(nodeID uint32) (*grpc.ClientConn, error) {
	existing, ok := c.connections.Load(nodeID)
	if ok {
		return existing.(*grpc.ClientConn), nil
	}

	conn, err := c.newClientConnection(nodeID)
	if err != nil {
		return nil, err
	}
	// Store the connection
	c.connections.Store(nodeID, conn)

	return conn, nil
}

func (c *ClientManager) newClientConnection(
	nodeID uint32,
) (*grpc.ClientConn, error) {
	c.log.Info("connecting to node", zap.Uint32("nodeID", nodeID))
	node, err := c.nodeRegistry.GetNode(nodeID)
	if err != nil {
		return nil, err
	}
	conn, err := node.BuildClient()
	if err != nil {
		return nil, err
	}

	return conn, nil
}
