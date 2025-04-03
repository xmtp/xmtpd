package registry

import (
	"context"
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/xmtp/xmtpd/pkg/tracing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/xmtp/xmtpd/pkg/abi/noderegistry"
	"github.com/xmtp/xmtpd/pkg/config"
	"go.uber.org/zap"
)

const (
	CONTRACT_CALL_TIMEOUT = 10 * time.Second
)

/*
*
The SmartContractRegistry notifies listeners of changes to the nodes by polling the contract
and diffing the returned node list with what is currently in memory.

This allows it to operate statelessly and not require a database, with a trade-off for latency.

Given how infrequently this list changes, that trade-off seems acceptable.
*/
type SmartContractRegistry struct {
	ctx      context.Context
	wg       sync.WaitGroup
	contract NodeRegistryContract
	logger   *zap.Logger
	// How frequently to poll the smart contract
	refreshInterval time.Duration
	// Mapping of nodes from ID -> Node
	nodes      map[uint32]Node
	nodesMutex sync.RWMutex
	// Notifiers for new nodes and changed nodes
	newNodesNotifier          *notifier[[]Node]
	changedNodeNotifiers      map[uint32]*notifier[Node]
	changedNodeNotifiersMutex sync.RWMutex
	cancel                    context.CancelFunc
}

// Interface implementation guard.
var _ NodeRegistry = &SmartContractRegistry{}

func NewSmartContractRegistry(
	ctx context.Context,
	ethclient bind.ContractCaller,
	logger *zap.Logger,
	options config.ContractsOptions,
) (*SmartContractRegistry, error) {
	contract, err := noderegistry.NewNodeRegistryCaller(
		common.HexToAddress(options.NodesContractAddress),
		ethclient,
	)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(ctx)

	return &SmartContractRegistry{
		ctx:                  ctx,
		contract:             contract,
		refreshInterval:      options.RegistryRefreshInterval,
		logger:               logger.Named("smartContractRegistry"),
		newNodesNotifier:     newNotifier[[]Node](),
		nodes:                make(map[uint32]Node),
		changedNodeNotifiers: make(map[uint32]*notifier[Node]),
		cancel:               cancel,
	}, nil
}

/*
*
Loads the initial state from the contract and starts a background refresh loop.

To stop refreshing callers should cancel the context
*
*/
func (s *SmartContractRegistry) Start() error {
	// If we can't load the data at least once, fail to start the service
	if err := s.refreshData(); err != nil {
		return err
	}

	tracing.GoPanicWrap(
		s.ctx,
		&s.wg,
		"smart-contract-registry",
		func(ctx context.Context) { s.refreshLoop() },
	)

	return nil
}

func (s *SmartContractRegistry) OnNewNodes() (<-chan []Node, CancelSubscription) {
	return s.newNodesNotifier.register()
}

func (s *SmartContractRegistry) OnChangedNode(
	nodeId uint32,
) (<-chan Node, CancelSubscription) {
	s.changedNodeNotifiersMutex.Lock()
	defer s.changedNodeNotifiersMutex.Unlock()

	notifier, ok := s.changedNodeNotifiers[nodeId]
	if !ok {
		notifier = newNotifier[Node]()
		s.changedNodeNotifiers[nodeId] = notifier
	}
	return notifier.register()
}

func (s *SmartContractRegistry) GetNodes() ([]Node, error) {
	s.nodesMutex.RLock()
	defer s.nodesMutex.RUnlock()

	nodes := make([]Node, 0)
	for _, node := range s.nodes {
		nodes = append(nodes, node)
	}
	return nodes, nil
}

func (s *SmartContractRegistry) GetNode(nodeId uint32) (*Node, error) {
	s.nodesMutex.RLock()
	defer s.nodesMutex.RUnlock()

	node, ok := s.nodes[nodeId]
	if !ok {
		return nil, errors.New("node not found")
	}
	return &node, nil
}

func (s *SmartContractRegistry) refreshLoop() {
	ticker := time.NewTicker(s.refreshInterval)
	for {
		select {
		case <-s.ctx.Done():
			return
		case <-ticker.C:
			if err := s.refreshData(); err != nil {
				s.logger.Error("Failed to refresh data", zap.Error(err))
			}
		}
	}
}

func (s *SmartContractRegistry) refreshData() error {
	fromContract, err := s.loadUnfilteredFromContract()
	if err != nil {
		return err
	}

	newNodes := []Node{}
	for _, node := range fromContract {
		// nodes realistically start at 100, but the contract fills the array with empty nodes
		if !node.IsValidConfig {
			continue
		}
		existingValue, ok := s.nodes[node.NodeID]
		if !ok {
			// New node found
			newNodes = append(newNodes, node)
		} else if !node.Equals(existingValue) {
			s.processChangedNode(node)
		}
	}

	if len(newNodes) > 0 {
		s.processNewNodes(newNodes)
	}

	return nil
}

func (s *SmartContractRegistry) processNewNodes(nodes []Node) {
	s.logger.Debug(
		"processing new nodes",
		zap.Int("count", len(nodes)),
		zap.Any("nodes", nodes))

	s.newNodesNotifier.trigger(nodes)

	s.nodesMutex.Lock()
	defer s.nodesMutex.Unlock()
	for _, node := range nodes {
		s.nodes[node.NodeID] = node
	}
}

func (s *SmartContractRegistry) processChangedNode(node Node) {
	s.nodesMutex.Lock()
	defer s.nodesMutex.Unlock()
	s.changedNodeNotifiersMutex.RLock()
	defer s.changedNodeNotifiersMutex.RUnlock()

	s.nodes[node.NodeID] = node
	s.logger.Info("processing changed node", zap.Any("node", node))
	if registry, ok := s.changedNodeNotifiers[node.NodeID]; ok {
		registry.trigger(node)
	}
}

func (s *SmartContractRegistry) loadUnfilteredFromContract() ([]Node, error) {
	ctx, cancel := context.WithTimeout(s.ctx, CONTRACT_CALL_TIMEOUT)
	defer cancel()
	nodes, err := s.contract.GetAllNodes(&bind.CallOpts{Context: ctx})
	if err != nil {
		return nil, err
	}
	out := make([]Node, len(nodes))
	for idx, node := range nodes {
		out[idx] = convertNode(node)
	}

	return out, nil
}

func (s *SmartContractRegistry) SetContractForTest(contract NodeRegistryContract) {
	s.contract = contract
}

func convertNode(rawNode noderegistry.INodeRegistryNodeWithId) Node {
	// Unmarshal the signing key.
	// If invalid, mark the config as being invalid as well. Clients should treat the
	// node as unhealthy in this case
	signingKey, err := crypto.UnmarshalPubkey(rawNode.Node.SigningKeyPub)
	isValidConfig := err == nil

	httpAddress := rawNode.Node.HttpAddress

	// Ensure the httpAddress is well formed
	if !strings.HasPrefix(httpAddress, "https://") && !strings.HasPrefix(httpAddress, "http://") {
		isValidConfig = false
	}

	if !rawNode.Node.InCanonicalNetwork {
		isValidConfig = false
	}

	return Node{
		NodeID:                    uint32(rawNode.NodeId.Uint64()),
		SigningKey:                signingKey,
		HttpAddress:               httpAddress,
		InCanonicalNetwork:        rawNode.Node.InCanonicalNetwork,
		MinMonthlyFeeMicroDollars: rawNode.Node.MinMonthlyFeeMicroDollars,
		IsValidConfig:             isValidConfig,
	}
}

func (f *SmartContractRegistry) Stop() {
	f.cancel()
	f.wg.Wait()
}
