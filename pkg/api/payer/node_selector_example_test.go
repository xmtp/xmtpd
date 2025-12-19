package payer_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/xmtp/xmtpd/pkg/api/payer"
	"github.com/xmtp/xmtpd/pkg/config"
	"github.com/xmtp/xmtpd/pkg/registry"
	"github.com/xmtp/xmtpd/pkg/topic"
	"go.uber.org/zap"
)

func ExampleManualNodeSelector() {
	ctx := context.Background()
	logger, _ := zap.NewDevelopment()

	cfg := &config.ContractsOptions{
		SettlementChain: config.SettlementChainOptions{
			NodeRegistryAddress:         "0x1234567890123456789012345678901234567890",
			RPCURL:                      "http://localhost:8545",
			NodeRegistryRefreshInterval: 30 * time.Second,
		},
	}

	nodeRegistry, err := registry.NewSmartContractRegistry(ctx, nil, logger, cfg)
	if err != nil {
		logger.Fatal("Failed to create registry", zap.Error(err))
	}

	selector := payer.NewManualNodeSelectorAlgorithm(nodeRegistry, []uint32{100})

	tpc := *topic.NewTopic(topic.TopicKindIdentityUpdatesV1, []byte("test"))
	nodeID, err := selector.GetNode(tpc)
	if err != nil {
		logger.Fatal("Failed to get node", zap.Error(err))
	}

	fmt.Printf("Selected node: %d\n", nodeID)
}

func ExampleOrderedPreferenceNodeSelector() {
	ctx := context.Background()
	logger, _ := zap.NewDevelopment()

	cfg := &config.ContractsOptions{
		SettlementChain: config.SettlementChainOptions{
			NodeRegistryAddress:         "0x1234567890123456789012345678901234567890",
			RPCURL:                      "http://localhost:8545",
			NodeRegistryRefreshInterval: 30 * time.Second,
		},
	}

	nodeRegistry, err := registry.NewSmartContractRegistry(ctx, nil, logger, cfg)
	if err != nil {
		logger.Fatal("Failed to create registry", zap.Error(err))
	}

	selector := payer.NewOrderedPreferenceNodeSelectorAlgorithm(nodeRegistry, []uint32{100, 200, 300})

	tpc := *topic.NewTopic(topic.TopicKindIdentityUpdatesV1, []byte("test"))
	nodeID, err := selector.GetNode(tpc)
	if err != nil {
		logger.Fatal("Failed to get node", zap.Error(err))
	}

	fmt.Printf("Selected node: %d\n", nodeID)
}

func ExampleRandomNodeSelector() {
	ctx := context.Background()
	logger, _ := zap.NewDevelopment()

	cfg := &config.ContractsOptions{
		SettlementChain: config.SettlementChainOptions{
			NodeRegistryAddress:         "0x1234567890123456789012345678901234567890",
			RPCURL:                      "http://localhost:8545",
			NodeRegistryRefreshInterval: 30 * time.Second,
		},
	}

	nodeRegistry, err := registry.NewSmartContractRegistry(ctx, nil, logger, cfg)
	if err != nil {
		logger.Fatal("Failed to create registry", zap.Error(err))
	}

	selector := payer.NewRandomNodeSelectorAlgorithm(nodeRegistry)

	tpc := *topic.NewTopic(topic.TopicKindIdentityUpdatesV1, []byte("test"))
	nodeID, err := selector.GetNode(tpc)
	if err != nil {
		logger.Fatal("Failed to get node", zap.Error(err))
	}

	fmt.Printf("Selected node: %d\n", nodeID)
}

func ExampleClosestNodeSelector() {
	ctx := context.Background()
	logger, _ := zap.NewDevelopment()

	cfg := &config.ContractsOptions{
		SettlementChain: config.SettlementChainOptions{
			NodeRegistryAddress:         "0x1234567890123456789012345678901234567890",
			RPCURL:                      "http://localhost:8545",
			NodeRegistryRefreshInterval: 30 * time.Second,
		},
	}

	nodeRegistry, err := registry.NewSmartContractRegistry(ctx, nil, logger, cfg)
	if err != nil {
		logger.Fatal("Failed to create registry", zap.Error(err))
	}

	selector := payer.NewClosestNodeSelectorAlgorithm(nodeRegistry, 5*time.Minute, 2*time.Second)

	tpc := *topic.NewTopic(topic.TopicKindIdentityUpdatesV1, []byte("test"))
	nodeID, err := selector.GetNode(tpc)
	if err != nil {
		logger.Fatal("Failed to get node", zap.Error(err))
	}

	fmt.Printf("Selected node: %d\n", nodeID)
}

func TestNodeSelectorsWithRealRegistry(t *testing.T) {
	t.Skip("Enable this test when you have a real blockchain node running")

	ctx := context.Background()
	logger, _ := zap.NewDevelopment()

	cfg := &config.ContractsOptions{
		SettlementChain: config.SettlementChainOptions{
			NodeRegistryAddress:         "0xYourContractAddress",
			RPCURL:                      "http://localhost:8545",
			NodeRegistryRefreshInterval: 30 * time.Second,
		},
	}

	nodeRegistry, err := registry.NewSmartContractRegistry(ctx, nil, logger, cfg)
	if err != nil {
		t.Fatalf("Failed to create registry: %v", err)
	}
	defer nodeRegistry.Stop()

	if err := nodeRegistry.Start(); err != nil {
		t.Fatalf("Failed to start registry: %v", err)
	}

	time.Sleep(2 * time.Second)

	tpc := *topic.NewTopic(topic.TopicKindIdentityUpdatesV1, []byte("test"))

	t.Run("StableHashing", func(t *testing.T) {
		selector := payer.NewStableHashingNodeSelectorAlgorithm(nodeRegistry)
		nodeID, err := selector.GetNode(tpc)
		if err != nil {
			t.Errorf("StableHashing failed: %v", err)
		} else {
			t.Logf("StableHashing selected node: %d", nodeID)
		}
	})

	t.Run("Random", func(t *testing.T) {
		selector := payer.NewRandomNodeSelectorAlgorithm(nodeRegistry)
		nodeID, err := selector.GetNode(tpc)
		if err != nil {
			t.Errorf("Random failed: %v", err)
		} else {
			t.Logf("Random selected node: %d", nodeID)
		}
	})

	t.Run("Closest", func(t *testing.T) {
		selector := payer.NewClosestNodeSelectorAlgorithm(nodeRegistry, 5*time.Minute, 2*time.Second)
		nodeID, err := selector.GetNode(tpc)
		if err != nil {
			t.Logf("Closest failed (expected if nodes not reachable): %v", err)
		} else {
			t.Logf("Closest selected node: %d", nodeID)
		}
	})
}

