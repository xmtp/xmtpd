package worker

import (
	"fmt"
	"math/rand/v2"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/db"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/testutils"
	envelopeTestUtils "github.com/xmtp/xmtpd/pkg/testutils/envelopes"
	"github.com/xmtp/xmtpd/pkg/topic"
)

func TestParsePartitionInfo(t *testing.T) {
	tests := []struct {
		name    string
		table   string
		wantErr bool
		nodeID  uint32
		start   uint64
		end     uint64
	}{
		{
			name:   "parsed ok - 1",
			table:  "gateway_envelopes_meta_o100_s0_1000000",
			nodeID: 100,
			start:  0,
			end:    1000000,
		},
		{
			name:   "parsed ok - 2",
			table:  "gateway_envelopes_meta_o400_s1000000_2000000",
			nodeID: 400,
			start:  1_000_000,
			end:    2_000_000,
		},
		{
			name:   "parsed ok - 3",
			table:  "gateway_envelopes_meta_o0_s7000000_8000000",
			nodeID: 0,
			start:  7_000_000,
			end:    8_000_000,
		},
		{
			name:    "inaplicable table",
			table:   "gateway_envelopes_meta",
			wantErr: true,
		},
		{
			name:    "invalid nodeID",
			table:   "gateway_envelopes_meta_oXYZ_s0_1000000",
			wantErr: true,
		},
		{
			name:    "invalid start offset",
			table:   "gateway_envelopes_meta_o100_sA_1000000",
			wantErr: true,
		},
		{
			name:    "invalid end value",
			table:   "gateway_envelopes_meta_o100_s0_B",
			wantErr: true,
		},
		{
			name:    "table has an unexpected prefix",
			table:   "pre_gateway_envelopes_meta_o400_s1000000_2000000",
			wantErr: true,
		},
		{
			name:    "table has an unexpected suffix",
			table:   "gateway_envelopes_meta_o400_s1000000_2000000_wat",
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			info, err := parsePartitionInfo(test.table)
			if test.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)

			require.Equal(t, test.table, info.name)
			require.Equal(t, test.nodeID, info.nodeID)
			require.Equal(t, test.start, info.start)
			require.Equal(t, test.end, info.end)
		})
	}
}

func TestSortPartitions(t *testing.T) {
	tests := []struct {
		name       string
		partitions []partitionTableInfo
		expected   map[uint32][]partitionTableInfo
	}{
		{
			name:       "empty",
			partitions: []partitionTableInfo{},
			expected:   map[uint32][]partitionTableInfo{},
		},
		{
			name: "sort single",
			partitions: []partitionTableInfo{
				{nodeID: 100, start: 100, end: 200},
				{nodeID: 100, start: 0, end: 10},
				{nodeID: 100, start: 200, end: 300},
				{nodeID: 100, start: 50, end: 80},
			},
			expected: map[uint32][]partitionTableInfo{
				100: {
					{nodeID: 100, start: 0, end: 10},
					{nodeID: 100, start: 50, end: 80},
					{nodeID: 100, start: 100, end: 200},
					{nodeID: 100, start: 200, end: 300},
				},
			},
		},
		{
			name: "sort multiple",
			partitions: []partitionTableInfo{
				{nodeID: 100, start: 100, end: 200},
				{nodeID: 100, start: 70, end: 100},
				{nodeID: 200, start: 1_000_000, end: 2_000_000},
				{nodeID: 200, start: 4_000_000, end: 5_000_000},
				{nodeID: 400, start: 25, end: 26},
				{nodeID: 100, start: 30, end: 70},
				{nodeID: 100, start: 0, end: 30},
				{nodeID: 1, start: 4_000_000, end: 5_000_000},
				{nodeID: 400, start: 10, end: 12},
				{nodeID: 1, start: 3_000_000, end: 4_000_000},
				{nodeID: 200, start: 0, end: 1_000_000},
				{nodeID: 1, start: 0, end: 1_000_000},
				{nodeID: 1, start: 2_000_000, end: 3_000_000},
				{nodeID: 400, start: 14, end: 16},
				{nodeID: 400, start: 16, end: 18},
				{nodeID: 200, start: 2_100_000, end: 2_500_000},
			},
			expected: map[uint32][]partitionTableInfo{
				100: {
					{nodeID: 100, start: 0, end: 30},
					{nodeID: 100, start: 30, end: 70},
					{nodeID: 100, start: 70, end: 100},
					{nodeID: 100, start: 100, end: 200},
				},
				200: {
					{nodeID: 200, start: 0, end: 1_000_000},
					{nodeID: 200, start: 1_000_000, end: 2_000_000},
					{nodeID: 200, start: 2_100_000, end: 2_500_000},
					{nodeID: 200, start: 4_000_000, end: 5_000_000},
				},
				1: {
					{nodeID: 1, start: 0, end: 1_000_000},
					{nodeID: 1, start: 2_000_000, end: 3_000_000},
					{nodeID: 1, start: 3_000_000, end: 4_000_000},
					{nodeID: 1, start: 4_000_000, end: 5_000_000},
				},
				400: {
					{nodeID: 400, start: 10, end: 12},
					{nodeID: 400, start: 14, end: 16},
					{nodeID: 400, start: 16, end: 18},
					{nodeID: 400, start: 25, end: 26},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			np := sortPartitions(test.partitions)
			require.Len(t, np.partitions, len(test.expected))

			for nodeID, expected := range test.expected {
				nodePartitions, ok := np.partitions[nodeID]
				require.True(t, ok)

				require.Len(t, nodePartitions, len(expected))

				for i, part := range nodePartitions {
					require.Equal(t, nodeID, expected[i].nodeID)
					require.Equal(t, part.start, expected[i].start)
					require.Equal(t, part.end, expected[i].end)
				}
			}
		})
	}
}

func TestValidatePartitionChain(t *testing.T) {
	tests := []struct {
		name      string
		chain     []partitionTableInfo
		isInvalid bool
	}{
		{
			name: "valid chain",
			chain: []partitionTableInfo{
				{nodeID: 100, start: 0, end: 100},
				{nodeID: 100, start: 100, end: 200},
				{nodeID: 100, start: 200, end: 300},
				{nodeID: 100, start: 300, end: 400},
				{nodeID: 100, start: 400, end: 450},
			},
		},
		{
			name: "single partition chain is valid",
			chain: []partitionTableInfo{
				{nodeID: 100, start: 0, end: 100},
			},
		},
		{
			name:  "chain with no partitions is technically valid",
			chain: []partitionTableInfo{},
		},
		{
			name:      "unsorted chain is invalid",
			isInvalid: true,
			chain: []partitionTableInfo{
				{nodeID: 100, start: 0, end: 100},
				{nodeID: 100, start: 100, end: 200},
				{nodeID: 100, start: 200, end: 300},
				{nodeID: 100, start: 400, end: 450}, // partition for higher range before lower one
				{nodeID: 100, start: 300, end: 400},
			},
		},
		{
			name:      "non-contigious chain is invalid",
			isInvalid: true,
			chain: []partitionTableInfo{
				{nodeID: 100, start: 0, end: 100},
				{nodeID: 100, start: 100, end: 200},
				// range 200-300 is missing
				{nodeID: 100, start: 300, end: 400},
				{nodeID: 100, start: 400, end: 450},
			},
		},
		{
			name:      "non-contigious chain is invalid",
			isInvalid: true,
			chain: []partitionTableInfo{
				{nodeID: 100, start: 0, end: 100},
				{nodeID: 100, start: 100, end: 199}, // till 199
				{nodeID: 100, start: 200, end: 300}, // from 200
				{nodeID: 100, start: 400, end: 450},
			},
		},
		{
			name:      "overlapping partitions are invalid",
			isInvalid: true,
			chain: []partitionTableInfo{
				{nodeID: 100, start: 0, end: 100},
				{nodeID: 100, start: 100, end: 200},
				{nodeID: 100, start: 200, end: 301}, // till 301
				{nodeID: 100, start: 300, end: 400}, // from 300
				{nodeID: 100, start: 400, end: 450},
			},
		},
		{
			name:      "mixed node IDs",
			isInvalid: true,
			chain: []partitionTableInfo{
				{nodeID: 100, start: 0, end: 100},
				{nodeID: 200, start: 100, end: 200},
				{nodeID: 300, start: 200, end: 300},
				{nodeID: 1, start: 300, end: 400},
				{nodeID: 400, start: 400, end: 450},
			},
		},
		{
			name:      "invalid partition size",
			isInvalid: true,
			chain: []partitionTableInfo{
				{nodeID: 100, start: 0, end: 100},
				{nodeID: 100, start: 100, end: 200},
				{nodeID: 100, start: 200, end: 300},
				{nodeID: 100, start: 300, end: 250}, // end smaller than start
				{nodeID: 100, start: 400, end: 450},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := validatePartitionChain(test.chain)
			if test.isInvalid {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}

func TestWorker_CreateAndListPartitionsForMultipleNodes(t *testing.T) {
	var (
		ctx            = t.Context()
		partitionSize  = uint64(1_000_000)
		ts             = newTestScaffold(t, partitionSize)
		partitionCount = 10 // Number of partitions to create.
	)

	inputParams := make(map[uint32]partitionParams)

	// Create new tables in the DB for random nodes and sequenceIDs.
	for range partitionCount {

		params := newPartitionParams(partitionSize)

		err := ts.worker.createPartition(ctx, params.nodeID, params.sequenceID)
		require.NoError(t, err)

		inputParams[params.nodeID] = params
	}

	// Now list partitions and make sure we have what we expect to have.
	partitions, err := ts.worker.getPartitionList(ctx)
	require.NoError(t, err)
	require.Len(t, partitions, partitionCount)

	for _, partition := range partitions {

		matchingParams, ok := inputParams[partition.nodeID]
		require.True(t, ok)

		t.Logf("partition for node %d for seqID %d of size %d: [%d,%d]: %v\n",
			matchingParams.nodeID,
			matchingParams.sequenceID,
			matchingParams.bandwidth,
			partition.start,
			partition.end,
			partition.name)

		require.Equal(t, matchingParams.nodeID, partition.nodeID)

		// Make sure partition is correct for the sequence ID we have.
		require.GreaterOrEqual(t, matchingParams.sequenceID, partition.start)
		require.Less(t, matchingParams.sequenceID, partition.end)

		// Make sure partition is of the correct size.
		size := partition.end - partition.start
		require.Equal(t, matchingParams.bandwidth, size)
	}
}

func TestWorker_CreateAndListPartitionsForSingleNode(t *testing.T) {
	var (
		ctx           = t.Context()
		partitionSize = uint64(1_000_000)
		ts            = newTestScaffold(t, partitionSize)

		// Sequence IDs that should force creation of several partitions.
		sequenceIDs = []uint64{
			1,         // 0-1M
			1_000_000, // 1M-2M
			2_000_001, // 2M-3M
			3_567_987, // 3M-4M
			4_123_567, // 4M-5M
			5_000_000, // 5M-6M
		}

		// Limit nodeID so we don't get errors when casting to int64.
		nodeID = rand.Uint32N(100_000)
	)

	for _, seqID := range sequenceIDs {
		err := ts.worker.createPartition(ctx, nodeID, uint64(seqID))
		require.NoError(t, err)
	}

	partitions, err := ts.worker.getPartitionList(ctx)
	require.NoError(t, err)

	require.NoError(t, err)

	require.Len(t, partitions, len(sequenceIDs))

	// Make sure we have a fitting partition for each sequence ID.
	for i, seqID := range sequenceIDs {
		partition := partitions[i]

		require.Equal(t, nodeID, partition.nodeID)

		// Make sure partition is correct for the sequence ID we have.
		require.GreaterOrEqual(t, seqID, partition.start)
		require.Less(t, seqID, partition.end)

		// Make sure partition is of the correct size.
		size := partition.end - partition.start
		require.Equal(t, partitionSize, size)
	}
}

func TestWorker_LastSequenceID(t *testing.T) {
	var (
		ctx = t.Context()
		ts  = newTestScaffold(t, 1_000_000)

		tableName = "gateway_envelopes_meta_o100_s0_1000000"
		count     = 10 + rand.UintN(20) // 10-30 envelopes
		nodeID    = uint32(100)
		envelopes = generateEnvelopes(t, nodeID, count)

		// Insert multiple batches to verify the largest sequence ID changes.
		firstBatch  = envelopes[:5]
		secondBatch = envelopes[5:]
	)

	// Make sure partition exists - just so our query doesn't fail.
	err := ts.worker.createPartition(ctx, nodeID, 1)
	require.NoError(t, err)

	// No envelopes now so largest should be 0.
	largest, err := ts.worker.getLastSequenceID(ctx, tableName)
	require.NoError(t, err)
	require.Zero(t, largest)

	// Insert first batch.
	testutils.InsertGatewayEnvelopes(t, ts.db.DB(), firstBatch)

	// Get the sequence ID of the last envelope in the first batch.
	last := firstBatch[len(firstBatch)-1].OriginatorSequenceID
	largest, err = ts.worker.getLastSequenceID(ctx, tableName)
	require.NoError(t, err)
	require.Equal(t, last, largest)

	// Insert second batch.
	testutils.InsertGatewayEnvelopes(t, ts.db.DB(), secondBatch)

	// Get the sequence ID of the last envelope in the first batch.
	last = secondBatch[len(secondBatch)-1].OriginatorSequenceID
	largest, err = ts.worker.getLastSequenceID(ctx, tableName)
	require.NoError(t, err)
	require.Equal(t, last, largest)
}

func TestWorker_PreparesPartition(t *testing.T) {
	var (
		ctx   = t.Context()
		db, _ = testutils.NewDB(t, ctx)
		log   = testutils.NewLog(t)
		cfg   = Config{
			// Super small partition with a size of 10, after 70% we create a new one.
			Partition: PartitionConfig{
				PartitionSize: 10,
				FillThreshold: 0.7,
			},
		}

		nodeID    = uint32(100)
		envelopes = generateEnvelopes(t, nodeID, 8)
	)

	worker := newWorkerWithConfig(cfg, log, db)

	// Have worker manually create a small partition.
	err := worker.createPartition(ctx, nodeID, 1)
	require.NoError(t, err)

	// Insert 7 envelopes - we're below threshold.
	testutils.InsertGatewayEnvelopes(t, db.DB(), envelopes[:7])

	err = worker.runDBCheck(ctx)
	require.NoError(t, err)

	// We should still have one partition.
	partitions, err := worker.getPartitionList(ctx)
	require.NoError(t, err)
	require.Len(t, partitions, 1)

	// Insert remaining envelope - we're now over fill threshold.
	testutils.InsertGatewayEnvelopes(t, db.DB(), envelopes[7:])

	// Run DB check again - this should see our partition is 80% full and create a new one.
	err = worker.runDBCheck(ctx)
	require.NoError(t, err)

	// We should now have two partitions.
	partitions, err = worker.getPartitionList(ctx)
	require.NoError(t, err)
	require.Len(t, partitions, 2)

	secondPartition := partitions[1]
	require.Equal(t, uint64(10), secondPartition.start)
	require.Equal(t, uint64(20), secondPartition.end)
	require.Equal(t, nodeID, secondPartition.nodeID)
}

func TestWorker_MonitorLoop(t *testing.T) {
	var (
		ctx   = t.Context()
		db, _ = testutils.NewDB(t, ctx)
		log   = testutils.NewLog(t)

		cfg = Config{
			Interval: 1 * time.Second,
			// Super small partition with a size of 10, after 70% we create a new one.
			Partition: PartitionConfig{
				PartitionSize: 10,
				FillThreshold: 0.7,
			},
		}

		nodeID    = uint32(100)
		envelopes = generateEnvelopes(t, nodeID, 8)
	)

	worker := newWorkerWithConfig(cfg, log, db)
	go worker.Start(ctx)

	// Have worker manually create a small partition.
	err := worker.createPartition(ctx, nodeID, 1)
	require.NoError(t, err)

	// Insert envelopes - to push us over the threshold.
	testutils.InsertGatewayEnvelopes(t, db.DB(), envelopes)

	time.Sleep(2 * time.Second)

	// We should now have two partitions.
	partitions, err := worker.getPartitionList(ctx)
	require.NoError(t, err)
	require.Len(t, partitions, 2)

	secondPartition := partitions[1]
	require.Equal(t, uint64(10), secondPartition.start)
	require.Equal(t, uint64(20), secondPartition.end)
	require.Equal(t, nodeID, secondPartition.nodeID)
}

func generateEnvelopes(
	t *testing.T,
	nodeID uint32,
	count uint,
) []queries.InsertGatewayEnvelopeParams {
	t.Helper()

	topic := topic.NewTopic(
		topic.TopicKindGroupMessagesV1,
		[]byte(fmt.Sprintf("generic-topic-%v", rand.Int())),
	)

	out := make([]queries.InsertGatewayEnvelopeParams, count)
	for i := range count {

		seqID := int64(1 + i)
		oe := testutils.Marshal(
			t,
			envelopeTestUtils.CreateOriginatorEnvelopeWithTopic(
				t,
				uint32(nodeID),
				uint64(seqID),
				topic.Bytes(),
			),
		)

		out[i] = queries.InsertGatewayEnvelopeParams{
			OriginatorNodeID:     int32(nodeID),
			OriginatorSequenceID: seqID,
			Topic:                topic.Bytes(),
			OriginatorEnvelope:   oe,
		}
	}

	return out
}

type testScaffold struct {
	db     *db.Handler
	worker *Worker
}

func newTestScaffold(t *testing.T, partitionSize uint64) testScaffold {
	t.Helper()

	var (
		db, _ = testutils.NewDB(t, t.Context())
		log   = testutils.NewLog(t)
		cfg   = Config{
			Partition: PartitionConfig{
				PartitionSize: partitionSize,
			},
		}
	)

	worker := newWorkerWithConfig(cfg, log, db)

	ts := testScaffold{
		db:     db,
		worker: worker,
	}

	return ts
}

type partitionParams struct {
	nodeID     uint32
	sequenceID uint64
	bandwidth  uint64
}

func newPartitionParams(partitionSize uint64) partitionParams {
	const (
		maxNodeID     = 100_000
		maxSequenceID = 300_000_000
	)

	return partitionParams{
		nodeID:     rand.Uint32N(maxNodeID),
		sequenceID: 1 + rand.Uint64N(maxSequenceID),
		bandwidth:  uint64(partitionSize),
	}
}
