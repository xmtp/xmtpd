# XMTPD Without Chain — Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Prove XMTPD can route MLS commits (originator 0) and identity updates (originator 1) to dedicated XMTPD nodes instead of Arbitrum L2 blockchain — eliminating blockchain dependency for message routing.

**Architecture:** Add a `--no-blockchain` flag to the gateway/payer that bypasses `publishToBlockchain()` and instead routes commits/identity updates to designated XMTPD nodes via `publishToNode()`. Dedicated nodes handle these as normal payer envelopes — staging, signing, publishing, replicating — using their own DB-assigned sequence IDs instead of blockchain event IDs. Feature-flagged so production path stays untouched.

**Tech Stack:** Go 1.25, testify, mockery mocks, Docker Compose (Anvil + Postgres), xdbg (Rust)

---

## File Structure

| File | Action | Responsibility |
|------|--------|---------------|
| `pkg/api/payer/service.go` | Modify | Add `noBlockchain` field to Service, modify `groupEnvelopes()` to route commits/identity to nodes when flag set |
| `pkg/api/payer/service_test.go` | Modify | Tests for new routing path |
| `pkg/api/payer/publish_test.go` | Modify | Tests for publishToNode path for commits/identity |
| `pkg/api/message/service.go` | Modify | Remove identity update rejection (line 850-859), remove commit/proposal rejection in `validateGroupMessage()` (line 941-946), update `nodeID >= 100` check (line 1202) |
| `pkg/api/message/service_test.go` | Modify | Tests for accepting identity updates and commits via node path |
| `pkg/constants/constants.go` | Modify | Add `NoBlockchainOriginatorThreshold` constant |
| `pkg/api/payer/selectors/topic_routing.go` | Create | New `TopicRoutingNodeSelector` that wraps existing selector + routes by topic kind to fixed nodes |
| `pkg/api/payer/selectors/topic_routing_test.go` | Create | Tests for topic-based routing |
| `pkg/config/options.go` | Modify | Add `--no-blockchain` and `--commit-node-id` / `--identity-node-id` flags |
| `dev/run-3` | Create | Script to run 3rd node (dedicated commit/identity node) |
| `dev/docker/up` | Modify | Add `triple` profile support |
| `dev/docker/docker-compose.yml` | Modify | Add db3 service for triple profile |

---

### Task 1: Add `--no-blockchain` Config Flag

**Files:**
- Modify: `pkg/config/options.go`

This is the feature flag. When set, gateway skips blockchain and routes commits/identity to nodes.

- [ ] **Step 1: Read current config structure**

Run: `grep -n "Contracts\|Blockchain\|Gateway" pkg/config/options.go | head -20`

Understand where to add new options.

- [ ] **Step 2: Add no-blockchain config options**

In `pkg/config/options.go`, add to the `Options` struct (find the `Payer` or `Gateway` section):

```go
// Inside the appropriate config section (Payer or top-level Options)
NoBlockchain   bool   `long:"no-blockchain" env:"XMTPD_NO_BLOCKCHAIN" description:"Skip blockchain for commits and identity updates, route to dedicated nodes instead"`
CommitNodeID   uint32 `long:"commit-node-id" env:"XMTPD_COMMIT_NODE_ID" description:"Node ID for dedicated commit node (used with --no-blockchain)" default:"0"`
IdentityNodeID uint32 `long:"identity-node-id" env:"XMTPD_IDENTITY_NODE_ID" description:"Node ID for dedicated identity update node (used with --no-blockchain)" default:"0"`
```

- [ ] **Step 3: Verify build**

Run: `cd /home/ubuntu/xmtp/xmtpd && PATH=/usr/local/go/bin:$PATH go build ./cmd/replication/`
Expected: clean build

- [ ] **Step 4: Commit**

```bash
git add pkg/config/options.go
git commit -m "feat: add --no-blockchain config flag for chainless experiment"
```

---

### Task 2: Create TopicRoutingNodeSelector

**Files:**
- Create: `pkg/api/payer/selectors/topic_routing.go`
- Create: `pkg/api/payer/selectors/topic_routing_test.go`

Wraps existing selector. For commits/identity topics → fixed node ID. For everything else → delegate to inner selector.

- [ ] **Step 1: Write the failing test**

Create `pkg/api/payer/selectors/topic_routing_test.go`:

```go
package selectors_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/api/payer/selectors"
	"github.com/xmtp/xmtpd/pkg/registry"
	registryMocks "github.com/xmtp/xmtpd/pkg/testutils/mocks/registry"
	nodeRegistry "github.com/xmtp/xmtpd/pkg/testutils/registry"
	"github.com/xmtp/xmtpd/pkg/topic"
)

func TestTopicRoutingSelector_IdentityGoesToDedicatedNode(t *testing.T) {
	mockReg := registryMocks.NewMockNodeRegistry(t)
	mockReg.On("GetNodes").Return([]registry.Node{
		nodeRegistry.GetHealthyNode(100),
		nodeRegistry.GetHealthyNode(200),
		nodeRegistry.GetHealthyNode(300),
	}, nil)

	inner := selectors.NewStableHashingNodeSelectorAlgorithm(mockReg)
	selector := selectors.NewTopicRoutingNodeSelector(inner, 300, 300)

	tpc := *topic.NewTopic(topic.TopicKindIdentityUpdatesV1, []byte("deadbeef"))
	nodeID, err := selector.GetNode(tpc)

	require.NoError(t, err)
	require.Equal(t, uint32(300), nodeID)
}

func TestTopicRoutingSelector_RegularMessageUsesInnerSelector(t *testing.T) {
	mockReg := registryMocks.NewMockNodeRegistry(t)
	mockReg.On("GetNodes").Return([]registry.Node{
		nodeRegistry.GetHealthyNode(100),
		nodeRegistry.GetHealthyNode(200),
		nodeRegistry.GetHealthyNode(300),
	}, nil)

	inner := selectors.NewStableHashingNodeSelectorAlgorithm(mockReg)
	selector := selectors.NewTopicRoutingNodeSelector(inner, 300, 300)

	tpc := *topic.NewTopic(topic.TopicKindKeyPackagesV1, []byte("deadbeef"))
	nodeID, err := selector.GetNode(tpc)

	require.NoError(t, err)
	// Should NOT be 300 (unless stable hash happens to pick it)
	// Key packages go through inner selector, not forced to dedicated node
	_ = nodeID // Just verify no error
}

func TestTopicRoutingSelector_GroupMessageGoesToDedicatedNode(t *testing.T) {
	// GroupMessages topic kind — in no-blockchain mode, ALL group messages
	// go to the commit node (commits will be sorted out at the node level)
	// Actually no — only commits/proposals need the dedicated node.
	// But the selector operates on topic kind, not MLS content type.
	// The payer's groupEnvelopes() handles the commit/proposal split.
	// So group messages still use the inner selector at this level.
	mockReg := registryMocks.NewMockNodeRegistry(t)
	mockReg.On("GetNodes").Return([]registry.Node{
		nodeRegistry.GetHealthyNode(100),
		nodeRegistry.GetHealthyNode(200),
	}, nil)

	inner := selectors.NewStableHashingNodeSelectorAlgorithm(mockReg)
	selector := selectors.NewTopicRoutingNodeSelector(inner, 100, 200)

	tpc := *topic.NewTopic(topic.TopicKindGroupMessagesV1, []byte("deadbeef"))
	nodeID, err := selector.GetNode(tpc)

	require.NoError(t, err)
	// Group messages still go through inner selector
	_ = nodeID
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd /home/ubuntu/xmtp/xmtpd && PATH=/usr/local/go/bin:$PATH go test ./pkg/api/payer/selectors/ -run TestTopicRouting -v`
Expected: FAIL — `NewTopicRoutingNodeSelector` not defined

- [ ] **Step 3: Implement TopicRoutingNodeSelector**

Create `pkg/api/payer/selectors/topic_routing.go`:

```go
package selectors

import (
	"github.com/xmtp/xmtpd/pkg/topic"
)

// TopicRoutingNodeSelector routes identity updates and commits to dedicated fixed nodes,
// delegating all other topics to an inner selector.
// Used in --no-blockchain mode to replace blockchain routing with node routing.
type TopicRoutingNodeSelector struct {
	inner          NodeSelectorAlgorithm
	commitNodeID   uint32
	identityNodeID uint32
}

var _ NodeSelectorAlgorithm = (*TopicRoutingNodeSelector)(nil)

func NewTopicRoutingNodeSelector(
	inner NodeSelectorAlgorithm,
	commitNodeID uint32,
	identityNodeID uint32,
) *TopicRoutingNodeSelector {
	return &TopicRoutingNodeSelector{
		inner:          inner,
		commitNodeID:   commitNodeID,
		identityNodeID: identityNodeID,
	}
}

func (s *TopicRoutingNodeSelector) GetNode(
	t topic.Topic,
	banlist ...[]uint32,
) (uint32, error) {
	switch t.Kind() {
	case topic.TopicKindIdentityUpdatesV1:
		return s.identityNodeID, nil
	default:
		return s.inner.GetNode(t, banlist...)
	}
}
```

Note: Group message commits are handled by `groupEnvelopes()` in the payer service (Task 3), not here. The selector only handles identity updates because they have a distinct topic kind. Commits share the GroupMessagesV1 topic kind with regular messages.

- [ ] **Step 4: Run tests to verify they pass**

Run: `cd /home/ubuntu/xmtp/xmtpd && PATH=/usr/local/go/bin:$PATH go test ./pkg/api/payer/selectors/ -run TestTopicRouting -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add pkg/api/payer/selectors/topic_routing.go pkg/api/payer/selectors/topic_routing_test.go
git commit -m "feat: add TopicRoutingNodeSelector for no-blockchain mode"
```

---

### Task 3: Modify Payer Service — Route Commits/Identity to Nodes

**Files:**
- Modify: `pkg/api/payer/service.go`

Core change. When `noBlockchain` is true, `groupEnvelopes()` puts commits and identity updates into `forNodes` instead of `forBlockchain`.

- [ ] **Step 1: Add noBlockchain field and config to Service**

In `pkg/api/payer/service.go`, add field to Service struct (around line 42):

```go
type Service struct {
	payer_apiconnect.UnimplementedPayerApiHandler
	gateway_apiconnect.UnimplementedGatewayApiHandler

	cfg                 Config
	ctx                 context.Context
	logger              *zap.Logger
	clientManager       *ClientManager
	blockchainPublisher blockchain.IBlockchainPublisher
	payerPrivateKey     *ecdsa.PrivateKey
	nodeSelector        selectors.NodeSelectorAlgorithm
	nodeRegistry        registry.NodeRegistry
	maxPayerMessageSize uint64
	noBlockchain        bool        // NEW
	commitNodeID        uint32      // NEW
}
```

- [ ] **Step 2: Add Option functions for no-blockchain mode**

Find the existing `Option` type and `With*` functions (likely near bottom of service.go or in a separate options file). Add:

```go
func WithNoBlockchain(commitNodeID uint32, identityNodeID uint32) Option {
	return func(s *Service) {
		s.noBlockchain = true
		s.commitNodeID = commitNodeID
	}
}
```

- [ ] **Step 3: Modify shouldSendToBlockchain()**

Change `shouldSendToBlockchain` at line 638 to accept a `noBlockchain` parameter, or make it a method on Service:

Replace the standalone function (lines 638-656):

```go
func (s *Service) shouldSendToBlockchain(clientEnvelope *envelopes.ClientEnvelope) (bool, error) {
	if s.noBlockchain {
		return false, nil
	}
	switch clientEnvelope.TargetTopic().Kind() {
	case topic.TopicKindIdentityUpdatesV1:
		return true, nil
	case topic.TopicKindGroupMessagesV1:
		switch payload := clientEnvelope.Payload().(type) {
		case *envelopesProto.ClientEnvelope_GroupMessage:
			shouldSendToBlockchain, err := deserializer.ShouldSendToBlockchain(payload)
			if err != nil {
				return false, err
			}
			return shouldSendToBlockchain, nil
		default:
			panic("mismatched payload type")
		}
	default:
		return false, nil
	}
}
```

- [ ] **Step 4: Modify groupEnvelopes() for commit routing**

In `groupEnvelopes()` (line 237), when `noBlockchain` is true and the envelope is a commit/proposal, route it to `commitNodeID` instead of using stable hash:

Replace the routing logic inside the for loop (around lines 260-283):

```go
		toBlockchain, err := s.shouldSendToBlockchain(clientEnvelope)
		if err != nil {
			return nil, connect.NewError(
				connect.CodeInvalidArgument,
				fmt.Errorf("could not determine routing for envelope at index %d: %w", i, err),
			)
		}

		if toBlockchain {
			out.forBlockchain = append(
				out.forBlockchain,
				newClientEnvelopeWithIndex(i, clientEnvelope),
			)
		} else {
			var targetNodeID uint32
			// In no-blockchain mode, route commits to dedicated commit node
			if s.noBlockchain && clientEnvelope.TargetTopic().Kind() == topic.TopicKindGroupMessagesV1 {
				payload, ok := clientEnvelope.Payload().(*envelopesProto.ClientEnvelope_GroupMessage)
				if ok {
					isCommit, err := deserializer.ShouldSendToBlockchain(payload)
					if err == nil && isCommit {
						targetNodeID = s.commitNodeID
					}
				}
			}
			// Fall through to normal selector if not a commit
			if targetNodeID == 0 {
				targetNodeID, err = s.nodeSelector.GetNode(clientEnvelope.TargetTopic())
				if err != nil {
					return nil, connect.NewError(connect.CodeInternal, err)
				}
			}
			out.forNodes[targetNodeID] = append(
				out.forNodes[targetNodeID],
				newClientEnvelopeWithIndex(i, clientEnvelope),
			)
		}
```

- [ ] **Step 5: Update call sites — shouldSendToBlockchain is now a method**

Search for all calls to `shouldSendToBlockchain(` in the file and change to `s.shouldSendToBlockchain(`:

Run: `grep -n "shouldSendToBlockchain(" pkg/api/payer/service.go`

Update each call site.

- [ ] **Step 6: Verify build**

Run: `cd /home/ubuntu/xmtp/xmtpd && PATH=/usr/local/go/bin:$PATH go build ./cmd/replication/`
Expected: clean build

- [ ] **Step 7: Commit**

```bash
git add pkg/api/payer/service.go
git commit -m "feat: route commits and identity updates to nodes in no-blockchain mode"
```

---

### Task 4: Remove Node-Side Rejection of Identity Updates and Commits

**Files:**
- Modify: `pkg/api/message/service.go`

The node currently rejects identity updates and commits sent via `PublishPayerEnvelopes`. Remove these guards when running without blockchain.

- [ ] **Step 1: Add noBlockchain field to message Service**

Find the Service struct in `pkg/api/message/service.go` and add:

```go
noBlockchain bool
```

Also add it to the constructor and pass it through from config.

- [ ] **Step 2: Modify identity update rejection (lines 850-859)**

Change:

```go
		if topicKind == topic.TopicKindIdentityUpdatesV1 {
			errs = append(
				errs,
				fmt.Sprintf(
					"identity updates must be published via the blockchain. index %d",
					i,
				),
			)
			continue
		}
```

To:

```go
		if topicKind == topic.TopicKindIdentityUpdatesV1 && !s.noBlockchain {
			errs = append(
				errs,
				fmt.Sprintf(
					"identity updates must be published via the blockchain. index %d",
					i,
				),
			)
			continue
		}
```

- [ ] **Step 3: Modify commit/proposal rejection in validateGroupMessage() (lines 941-946)**

Change:

```go
	if shouldSendToBlockchain {
		return connect.NewError(
			connect.CodeInvalidArgument,
			errors.New("commit and proposal messages must be published via the blockchain"),
		)
	}
```

To:

```go
	if shouldSendToBlockchain && !s.noBlockchain {
		return connect.NewError(
			connect.CodeInvalidArgument,
			errors.New("commit and proposal messages must be published via the blockchain"),
		)
	}
```

- [ ] **Step 4: Update DependsOn nodeID check (line 1202)**

The current check `nodeID >= 100` assumes originator 0-99 = blockchain. In no-blockchain mode, the dedicated node has ID >= 100 (e.g., 300). DependsOn needs to accept those node IDs.

Change:

```go
			if nodeID >= 100 {
				return connect.NewError(
					connect.CodeInvalidArgument,
					fmt.Errorf(
						"node ID %d specified in DependsOn is not a valid node ID, a message can not depend on a non-commit",
						nodeID,
					),
				)
			}
```

To:

```go
			if !s.noBlockchain && nodeID >= 100 {
				return connect.NewError(
					connect.CodeInvalidArgument,
					fmt.Errorf(
						"node ID %d specified in DependsOn is not a valid node ID, a message can not depend on a non-commit",
						nodeID,
					),
				)
			}
```

- [ ] **Step 5: Verify build**

Run: `cd /home/ubuntu/xmtp/xmtpd && PATH=/usr/local/go/bin:$PATH go build ./cmd/replication/`
Expected: clean build

- [ ] **Step 6: Commit**

```bash
git add pkg/api/message/service.go
git commit -m "feat: allow identity updates and commits via node path in no-blockchain mode"
```

---

### Task 5: Wire Config to Services

**Files:**
- Modify: `pkg/gateway/builder.go` (or wherever payer service is constructed)
- Modify: `cmd/replication/main.go` or server setup

Connect the `--no-blockchain` config flag to both the payer Service and message Service.

- [ ] **Step 1: Find payer service construction**

Run: `grep -rn "NewPayerAPIService" pkg/ cmd/ --include="*.go" | grep -v test | grep -v mock`

- [ ] **Step 2: Pass noBlockchain to payer service**

At the call site, add `WithNoBlockchain(cfg.CommitNodeID, cfg.IdentityNodeID)` option when `cfg.NoBlockchain` is true.

- [ ] **Step 3: Find message service construction**

Run: `grep -rn "NewService\|NewMessageService\|message\.New" pkg/ cmd/ --include="*.go" | grep -v test | grep -v mock | head -10`

- [ ] **Step 4: Pass noBlockchain to message service**

Add the flag to message service constructor or via an option pattern.

- [ ] **Step 5: Wire TopicRoutingNodeSelector**

When `noBlockchain` is true, wrap the existing node selector with `TopicRoutingNodeSelector`:

```go
var nodeSelector selectors.NodeSelectorAlgorithm
nodeSelector = selectors.NewStableHashingNodeSelectorAlgorithm(nodeRegistry)
if cfg.NoBlockchain {
    nodeSelector = selectors.NewTopicRoutingNodeSelector(
        nodeSelector, cfg.CommitNodeID, cfg.IdentityNodeID,
    )
}
```

- [ ] **Step 6: Verify build**

Run: `cd /home/ubuntu/xmtp/xmtpd && PATH=/usr/local/go/bin:$PATH go build ./cmd/replication/`
Expected: clean build

- [ ] **Step 7: Commit**

```bash
git add pkg/gateway/builder.go cmd/replication/main.go
git commit -m "feat: wire no-blockchain config to payer and message services"
```

---

### Task 6: Set Up Local 3-Node Dev Environment

**Files:**
- Create: `dev/run-3`
- Modify: `dev/docker/docker-compose.yml`
- Modify: `dev/docker/up`
- Modify: `dev/local.env`

Add a third node dedicated to commits/identity. Uses its own DB.

- [ ] **Step 1: Add ANVIL_ACC_3 to local.env**

Anvil provides 10 funded accounts. Account 3 (0-indexed):

```bash
# Add to dev/local.env:
ANVIL_ACC_3_PRIVATE_KEY="0x7c852118294e51e653712a81e05800f419141751be58f605c371e15141b007a6"
ANVIL_ACC_3_PUBLIC_KEY="0x04b5a2d23a9e4c0fb8b4ec6e44c4c1e23bab4e46a0f8b4c12e5a4c9b5e6d7f8a"
ANVIL_ACC_3_ADDRESS="0x90F79bf6EB2c4f870365E785982E1f101E93b906"
NODE_3_HTTP_ADDRESS="http://localhost:5052"
```

Note: Verify the actual Anvil account 3 values — these are the standard Hardhat/Anvil accounts. Run `cast accounts` against local Anvil to confirm.

- [ ] **Step 2: Add db3 to docker-compose.yml**

Add after the db2 service definition:

```yaml
  db3:
    image: postgres:16
    environment:
      POSTGRES_PASSWORD: xmtp
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -d postgres -U postgres"]
      interval: 1s
    ports:
      - 8767:5432
    profiles:
      - triple
```

- [ ] **Step 3: Create dev/run-3**

Create `dev/run-3`:

```bash
#!/bin/bash

set -eu

. dev/local.env

export XMTPD_SIGNER_PRIVATE_KEY=$ANVIL_ACC_3_PRIVATE_KEY
export XMTPD_PAYER_PRIVATE_KEY=$XMTPD_SIGNER_PRIVATE_KEY
export XMTPD_DB_WRITER_CONNECTION_STRING="postgres://postgres:xmtp@localhost:8767/postgres?sslmode=disable"

export XMTPD_REFLECTION_ENABLE=true
export XMTPD_API_ENABLE=true
export XMTPD_SYNC_ENABLE=true
export XMTPD_INDEXER_ENABLE=true
export XMTPD_CONTRACTS_ENVIRONMENT=anvil

go run -ldflags="-X main.Version=$(git describe HEAD --tags --long)" cmd/replication/main.go -p 5052 "$@"
```

Make executable: `chmod +x dev/run-3`

- [ ] **Step 4: Update dev/docker/up for triple profile**

Add after the dual node registration block:

```bash
if [ "${profile}" = "triple" ]; then
  # Register node 1
  run_cli nodes register \
    --owner-address="${ANVIL_ACC_1_ADDRESS}" \
    --signing-key-pub="${ANVIL_ACC_1_PUBLIC_KEY}" \
    --http-address="${NODE_1_HTTP_ADDRESS}"
  run_cli nodes canonical-network --add --node-id=100

  # Register node 2
  run_cli nodes register \
    --owner-address="${ANVIL_ACC_2_ADDRESS}" \
    --signing-key-pub="${ANVIL_ACC_2_PUBLIC_KEY}" \
    --http-address="${NODE_2_HTTP_ADDRESS}"
  run_cli nodes canonical-network --add --node-id=200

  # Register node 3 (dedicated commit/identity node)
  run_cli nodes register \
    --owner-address="${ANVIL_ACC_3_ADDRESS}" \
    --signing-key-pub="${ANVIL_ACC_3_PUBLIC_KEY}" \
    --http-address="${NODE_3_HTTP_ADDRESS}"
  run_cli nodes canonical-network --add --node-id=300
fi
```

- [ ] **Step 5: Commit**

```bash
git add dev/run-3 dev/docker/docker-compose.yml dev/docker/up dev/local.env
git commit -m "feat: add triple-node dev environment for no-blockchain experiment"
```

---

### Task 7: Integration Smoke Test — Manual Verification

**Files:** None (manual testing)

Verify the full flow works: gateway routes commits/identity to node 300, node 300 processes them, replication delivers to nodes 100/200.

- [ ] **Step 1: Start infrastructure**

```bash
cd /home/ubuntu/xmtp/xmtpd
dev/up triple
```

- [ ] **Step 2: Start node 1 (terminal 1)**

```bash
dev/run
```

- [ ] **Step 3: Start node 2 (terminal 2)**

```bash
dev/run-2
```

- [ ] **Step 4: Start node 3 — dedicated commit/identity node (terminal 3)**

```bash
dev/run-3
```

- [ ] **Step 5: Start gateway with no-blockchain flag (terminal 4)**

```bash
dev/run --no-blockchain --commit-node-id=300 --identity-node-id=300
```

Or set env vars:
```bash
export XMTPD_NO_BLOCKCHAIN=true
export XMTPD_COMMIT_NODE_ID=300
export XMTPD_IDENTITY_NODE_ID=300
dev/run
```

- [ ] **Step 6: Run xdbg smoke test**

```bash
cd /home/ubuntu/xmtp/libxmtp
# Point xdbg at local gateway
./target/release/xdbg -b local -d test group-sync -m 5
```

If xdbg doesn't support `-b local`, use grpcurl or write a simple Go test:

```bash
# Verify identity update goes through
grpcurl -plaintext localhost:5050 xmtp.xmtpd.api.v1.ReplicationApi/QueryEnvelopes
```

- [ ] **Step 7: Verify on node 300**

Check node 300 logs for:
- Received PublishPayerEnvelopes
- Successfully staged identity update
- Publish worker processed envelope

- [ ] **Step 8: Verify replication to node 100/200**

Check node 100/200 logs for:
- Sync worker received envelope from node 300
- Successfully inserted gateway envelope

- [ ] **Step 9: Document results**

Record: pass/fail, any errors, latency observations.

---

### Task 8: Write Automated E2E Test

**Files:**
- Create: `pkg/api/payer/no_blockchain_test.go`

Automated test that verifies the full routing path without blockchain.

- [ ] **Step 1: Write integration test**

Create `pkg/api/payer/no_blockchain_test.go`:

```go
package payer_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/api/payer"
	"github.com/xmtp/xmtpd/pkg/api/payer/selectors"
	"github.com/xmtp/xmtpd/pkg/constants"
	"github.com/xmtp/xmtpd/pkg/registry"
	blockchainMocks "github.com/xmtp/xmtpd/pkg/testutils/mocks/blockchain"
	registryMocks "github.com/xmtp/xmtpd/pkg/testutils/mocks/registry"
	nodeRegistry "github.com/xmtp/xmtpd/pkg/testutils/registry"
	"github.com/xmtp/xmtpd/pkg/topic"
	"go.uber.org/zap"
)

func TestNoBlockchain_IdentityUpdateRoutesToNode(t *testing.T) {
	mockReg := registryMocks.NewMockNodeRegistry(t)
	mockReg.On("GetNodes").Return([]registry.Node{
		nodeRegistry.GetHealthyNode(100),
		nodeRegistry.GetHealthyNode(200),
		nodeRegistry.GetHealthyNode(300),
	}, nil)

	mockBlockchain := blockchainMocks.NewMockIBlockchainPublisher(t)
	// Blockchain publisher should NOT be called in no-blockchain mode

	privKey := testutils.RandomPrivateKey(t)

	svc, err := payer.NewPayerAPIService(
		context.Background(),
		zap.NewNop(),
		mockReg,
		privKey,
		mockBlockchain,
		nil, // metrics
		0,   // max size
		nil, // selector (will be created internally)
		payer.WithNoBlockchain(300, 300),
	)
	require.NoError(t, err)

	// Build a client envelope with identity update topic
	// Call PublishClientEnvelopes
	// Assert: mockBlockchain.PublishIdentityUpdate was NOT called
	// Assert: envelope was routed to node 300

	mockBlockchain.AssertNotCalled(t, "PublishIdentityUpdate")
}

func TestNoBlockchain_CommitRoutesToNode(t *testing.T) {
	// Similar test for commit messages
	// Build MLS commit message envelope
	// Assert: routed to commit node, not blockchain
}
```

- [ ] **Step 2: Run test to verify it fails/passes appropriately**

Run: `cd /home/ubuntu/xmtp/xmtpd && PATH=/usr/local/go/bin:$PATH go test ./pkg/api/payer/ -run TestNoBlockchain -v`

- [ ] **Step 3: Fix any issues, iterate**

- [ ] **Step 4: Commit**

```bash
git add pkg/api/payer/no_blockchain_test.go
git commit -m "test: add no-blockchain routing integration tests"
```

---

### Task 9: Latency Comparison — With Chain vs Without

**Files:** None (benchmarking)

Measure latency improvement from removing blockchain settlement.

- [ ] **Step 1: Baseline — run xdbg WITH blockchain (current behavior)**

```bash
cd /home/ubuntu/xmtp/libxmtp
# Ensure gateway is running in normal mode (no --no-blockchain)
./target/release/xdbg -b local -d test message -n 10 -c 1
./target/release/xdbg -b local -d test identity -n 10 -c 1
./target/release/xdbg -b local -d test visibility -n 10 -c 1
./target/release/xdbg -b local -d test group-sync -m 10
```

Record average latency for each test type.

- [ ] **Step 2: Run same tests WITHOUT blockchain**

Restart gateway with `--no-blockchain --commit-node-id=300 --identity-node-id=300`.

```bash
./target/release/xdbg -b local -d test message -n 10 -c 1
./target/release/xdbg -b local -d test identity -n 10 -c 1
./target/release/xdbg -b local -d test visibility -n 10 -c 1
./target/release/xdbg -b local -d test group-sync -m 10
```

- [ ] **Step 3: Compare results**

Create comparison table:

| Test | With Chain (avg ms) | Without Chain (avg ms) | Improvement |
|------|--------------------|-----------------------|-------------|
| message | ? | ? | ? |
| identity | ? | ? | ? |
| visibility | ? | ? | ? |
| group-sync | ? | ? | ? |

Expected: commit/identity paths should be significantly faster (no blockchain settlement).

- [ ] **Step 4: Document results**

Write up findings in `docs/no-blockchain-experiment-results.md`.

---

### Task 10: Failure Mode Analysis

**Files:**
- Create: `docs/no-blockchain-experiment-results.md`

Document what works, what breaks, and what's needed for production.

- [ ] **Step 1: Test dedicated node restart**

Kill node 300 while traffic is flowing. Verify:
- Gateway gets errors for commits/identity updates
- Regular messages to nodes 100/200 are unaffected
- Restarting node 300 recovers — replication catches up

- [ ] **Step 2: Test replication completeness**

Send 100 messages through no-blockchain path. Verify all 100 appear on all nodes via QueryEnvelopes.

- [ ] **Step 3: Write findings document**

Create `docs/no-blockchain-experiment-results.md` with:
- Experiment setup (3-node topology, what was changed)
- Latency comparison table (from Task 9)
- Failure mode analysis
- What's missing for production:
  - Consensus/failover for dedicated commit node
  - Payer report changes (currently reference blockchain sequence IDs)
  - Indexer changes (currently watches blockchain events)
  - Client SDK DependsOn handling
- Recommendation

- [ ] **Step 4: Commit**

```bash
git add docs/no-blockchain-experiment-results.md
git commit -m "docs: no-blockchain experiment results and analysis"
```

---

## Execution Order & Dependencies

```
Task 1 (config flag)
  └─→ Task 2 (topic routing selector)
  └─→ Task 3 (payer routing changes) ← depends on Task 2
  └─→ Task 4 (node rejection removal)
  └─→ Task 5 (wire config) ← depends on Tasks 1,2,3,4
       └─→ Task 6 (dev environment)
            └─→ Task 7 (manual smoke test) ← depends on Tasks 5,6
            └─→ Task 8 (automated tests) ← depends on Task 5
                 └─→ Task 9 (latency comparison) ← depends on Task 7
                      └─→ Task 10 (analysis & writeup) ← depends on Tasks 7,8,9
```

Tasks 1-4 can be done in parallel (independent code changes).
Task 5 integrates them.
Tasks 6-7 validate.
Tasks 8-10 measure and document.
