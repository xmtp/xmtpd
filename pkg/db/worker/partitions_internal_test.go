package worker

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
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

// TODO: Remove
func printPartitionSlice(p []partitionTableInfo) {
	for i, part := range p {
		fmt.Printf("[%d]:{nodeID:%d,start:%d,end:%d},", i, part.nodeID, part.start, part.end)
	}
	fmt.Printf("\n")
}
