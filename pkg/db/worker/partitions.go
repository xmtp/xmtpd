package worker

import (
	"errors"
	"fmt"
	"slices"
)

type partitionTableInfo struct {
	name   string
	nodeID uint32
	start  uint64
	end    uint64
}

type nodePartitions struct {
	partitions map[uint32][]partitionTableInfo
}

// group partitions by originator and sort them
func sortPartitions(partitions []partitionTableInfo) nodePartitions {
	out := make(map[uint32][]partitionTableInfo)
	for _, partition := range partitions {
		_, ok := out[partition.nodeID]
		if !ok {
			out[partition.nodeID] = make([]partitionTableInfo, 0)
		}

		out[partition.nodeID] = append(out[partition.nodeID], partition)
	}

	for nodeID := range out {
		// Sort partitions in ascending order
		slices.SortFunc(out[nodeID], partitionSortFunc)
	}

	np := nodePartitions{
		partitions: out,
	}

	return np
}

func (p *nodePartitions) validate() error {
	var errs []error
	for nodeID, partitions := range p.partitions {

		err := validatePartitionChain(partitions)
		if err != nil {
			errs = append(errs, fmt.Errorf("partition chain for %v is invalid: %w", nodeID, err))
		}
	}

	return errors.Join(errs...)
}

func validatePartitionChain(partitions []partitionTableInfo) error {
	var prev *partitionTableInfo

	for _, p := range partitions {

		if p.nodeID == 0 {
			return errors.New("invalid node ID")
		}

		if p.start >= p.end {
			return fmt.Errorf(
				"invalid partition size (start: %v, end: %v)", p.start, p.end)
		}

		// Save the current partition info so we can compare it against the next one.
		if prev == nil {
			prev = &p
			continue
		}

		if prev.nodeID != p.nodeID {
			return fmt.Errorf(
				"partitions refer to different nodes (prev: %v, current: %v)",
				prev.nodeID,
				p.nodeID,
			)
		}

		if prev.end != p.start {
			return fmt.Errorf(
				"partitions not contigious (previous_end: %v, current_start: %v)",
				prev.end,
				p.start,
			)
		}
	}

	return nil
}

func partitionSortFunc(a, b partitionTableInfo) int {
	if a.start < b.start && a.end < b.end {
		return -1
	}

	if a.start > b.start && a.end > b.end {
		return 1
	}

	return 0
}
