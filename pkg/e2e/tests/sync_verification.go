package tests

import (
	"context"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/e2e/types"
	"go.uber.org/zap"
)

// SyncVerificationTest verifies that all nodes in a cluster agree on the latest
// sequence IDs for each originator. After publishing envelopes and allowing time
// for replication, it compares vector clocks across all nodes to ensure
// consistency.
type SyncVerificationTest struct{}

func NewSyncVerificationTest() *SyncVerificationTest {
	return &SyncVerificationTest{}
}

func (t *SyncVerificationTest) Name() string {
	return "sync-verification"
}

func (t *SyncVerificationTest) Description() string {
	return "Cross-check last sequence ID per originator across all nodes"
}

func (t *SyncVerificationTest) Run(
	ctx context.Context,
	env *types.Environment,
) error {
	require := require.New(env.T())

	require.NoError(env.AddNode(ctx))
	require.NoError(env.AddNode(ctx))
	require.NoError(env.AddNode(ctx))
	require.NoError(env.AddGateway(ctx))

	// Publish from multiple nodes to create distinct originators.
	require.NoError(env.NewClient(100))
	require.NoError(env.NewClient(200))

	const envelopesPerNode = 20
	require.NoError(env.Client(100).PublishEnvelopes(ctx, envelopesPerNode))
	require.NoError(env.Client(200).PublishEnvelopes(ctx, envelopesPerNode))

	// Wait for all nodes to replicate all envelopes.
	checkCtx, cancel := context.WithTimeout(ctx, 90*time.Second)
	defer cancel()
	for _, n := range env.Nodes() {
		require.NoError(n.WaitForEnvelopes(checkCtx, envelopesPerNode))
	}

	// Allow a brief settling period for replication to converge.
	time.Sleep(5 * time.Second)

	// Get the reference vector clock from node 100.
	nodes := env.Nodes()
	refClock, err := nodes[0].GetVectorClock(ctx)
	require.NoError(err, "failed to get vector clock from reference node")
	require.NotEmpty(refClock, "reference vector clock should not be empty")

	env.Logger.Info("reference vector clock",
		zap.Uint32("node_id", nodes[0].ID()),
		zap.Int("originators", len(refClock)),
	)

	// Compare every other node's vector clock against the reference.
	for i := 1; i < len(nodes); i++ {
		nodeClock, err := nodes[i].GetVectorClock(ctx)
		require.NoError(err,
			"failed to get vector clock from node %d", nodes[i].ID(),
		)

		nodeClockMap := make(map[int32]int64)
		for _, entry := range nodeClock {
			nodeClockMap[entry.OriginatorNodeID] = entry.OriginatorSequenceID
		}

		for _, refEntry := range refClock {
			nodeSeq, exists := nodeClockMap[refEntry.OriginatorNodeID]
			require.True(exists,
				"node %d missing originator %d in vector clock",
				nodes[i].ID(), refEntry.OriginatorNodeID,
			)
			require.Equal(
				refEntry.OriginatorSequenceID, nodeSeq,
				"sequence ID mismatch for originator %d: "+
					"node %d has %d, node %d has %d",
				refEntry.OriginatorNodeID,
				nodes[0].ID(), refEntry.OriginatorSequenceID,
				nodes[i].ID(), nodeSeq,
			)
		}

		env.Logger.Info("vector clock matches reference",
			zap.Uint32("node_id", nodes[i].ID()),
		)
	}

	env.Logger.Info("sync verification test completed successfully")
	return nil
}

var _ types.Test = (*SyncVerificationTest)(nil)
