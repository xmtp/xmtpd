package payer

import (
	"sync"

	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"

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
	log              *zap.Logger
	connections      map[uint32]*grpc.ClientConn
	connectionsMutex sync.RWMutex
	nodeRegistry     registry.NodeRegistry
	clientMetrics    *grpcprom.ClientMetrics
}

func NewClientManager(
	log *zap.Logger,
	nodeRegistry registry.NodeRegistry,
	clientMetrics *grpcprom.ClientMetrics,
) *ClientManager {
	return &ClientManager{
		log:           log,
		nodeRegistry:  nodeRegistry,
		connections:   make(map[uint32]*grpc.ClientConn),
		clientMetrics: clientMetrics,
	}
}

func (c *ClientManager) GetClient(nodeID uint32) (*grpc.ClientConn, error) {
	c.connectionsMutex.RLock()
	existing, ok := c.connections[nodeID]
	c.connectionsMutex.RUnlock()
	if ok {
		return existing, nil
	}

	c.connectionsMutex.Lock()
	defer c.connectionsMutex.Unlock()

	// Check again under the write lock
	if existing, ok := c.connections[nodeID]; ok {
		return existing, nil
	}

	conn, err := c.newClientConnection(nodeID)
	if err != nil {
		return nil, err
	}
	// Store the connection
	c.connections[nodeID] = conn

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
	conn, err := node.BuildClient(
		grpc.WithUnaryInterceptor(c.clientMetrics.UnaryClientInterceptor()),
		grpc.WithStreamInterceptor(c.clientMetrics.StreamClientInterceptor()),
	)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
