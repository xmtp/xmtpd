package payer

import (
	"sync"

	"connectrpc.com/connect"

	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api/message_apiconnect"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/metadata_api/metadata_apiconnect"
	"github.com/xmtp/xmtpd/pkg/registry"
	"github.com/xmtp/xmtpd/pkg/utils"
	"go.uber.org/zap"
)

// ClientManager contains a mapping of nodeIDs to Replicaion API client connections.
// These client connections are safe to be shared and re-used and will automatically attempt
// to reconnect if the underlying socket connection is lost.
type ClientManager struct {
	logger             *zap.Logger
	nodeRegistry       registry.NodeRegistry
	replicationClients map[uint32]message_apiconnect.ReplicationApiClient
	metadataClients    map[uint32]metadata_apiconnect.MetadataApiClient
	clientsMu          sync.RWMutex
	clientMetrics      *utils.ConnectClientMetrics
}

func NewClientManager(
	logger *zap.Logger,
	nodeRegistry registry.NodeRegistry,
	clientMetrics *utils.ConnectClientMetrics,
) *ClientManager {
	return &ClientManager{
		logger:             logger,
		nodeRegistry:       nodeRegistry,
		replicationClients: make(map[uint32]message_apiconnect.ReplicationApiClient),
		metadataClients:    make(map[uint32]metadata_apiconnect.MetadataApiClient),
		clientMetrics:      clientMetrics,
	}
}

func (c *ClientManager) GetReplicationClient(
	nodeID uint32,
) (message_apiconnect.ReplicationApiClient, error) {
	c.clientsMu.RLock()
	existing, ok := c.replicationClients[nodeID]
	c.clientsMu.RUnlock()
	if ok {
		return existing, nil
	}

	c.clientsMu.Lock()
	defer c.clientsMu.Unlock()

	// Check again under the write lock
	if existing, ok := c.replicationClients[nodeID]; ok {
		return existing, nil
	}

	client, err := c.newReplicationClientConnection(nodeID)
	if err != nil {
		return nil, err
	}

	c.replicationClients[nodeID] = client

	return client, nil
}

func (c *ClientManager) newReplicationClientConnection(
	nodeID uint32,
) (message_apiconnect.ReplicationApiClient, error) {
	c.logger.Info(
		"connecting to replication API",
		utils.OriginatorIDField(nodeID),
	)

	node, err := c.nodeRegistry.GetNode(nodeID)
	if err != nil {
		return nil, err
	}

	conn, err := node.BuildReplicationAPIClient(
		connect.WithInterceptors(c.clientMetrics.Interceptor()),
	)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func (c *ClientManager) GetMetadataClient(
	nodeID uint32,
) (metadata_apiconnect.MetadataApiClient, error) {
	c.clientsMu.RLock()
	existing, ok := c.metadataClients[nodeID]
	c.clientsMu.RUnlock()
	if ok {
		return existing, nil
	}

	c.clientsMu.Lock()
	defer c.clientsMu.Unlock()

	// Check again under the write lock
	if existing, ok := c.metadataClients[nodeID]; ok {
		return existing, nil
	}

	client, err := c.newMetadataClientConnection(nodeID)
	if err != nil {
		return nil, err
	}

	c.metadataClients[nodeID] = client

	return client, nil
}

func (c *ClientManager) newMetadataClientConnection(
	nodeID uint32,
) (metadata_apiconnect.MetadataApiClient, error) {
	c.logger.Info(
		"connecting to metadata API",
		utils.OriginatorIDField(nodeID),
	)

	node, err := c.nodeRegistry.GetNode(nodeID)
	if err != nil {
		return nil, err
	}

	client, err := node.BuildMetadataAPIClient(
		connect.WithInterceptors(c.clientMetrics.Interceptor()),
	)
	if err != nil {
		return nil, err
	}

	return client, nil
}
