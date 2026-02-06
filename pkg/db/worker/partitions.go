package worker

import (
	"errors"
	"fmt"
	"regexp"
	"slices"
	"strconv"
)

// Regex to parse partition information, for example: gateway_envelopes_meta_o100_s1000000_2000000
var partitionRe = regexp.MustCompile(`^gateway_envelopes_meta_o(\d+)_s(\d+)_(\d+)$`)

type nodePartitions struct {
	partitions map[uint32][]partitionTableInfo
}

type partitionTableInfo struct {
	name   string
	nodeID uint32
	start  uint64
	end    uint64
}

func parsePartitionInfo(table string) (partitionTableInfo, error) {
	fields := partitionRe.FindStringSubmatch(table)
	if len(fields) != 4 {
		return partitionTableInfo{}, errors.New("unexpected table name format")
	}

	nodeID, err := strconv.ParseUint(fields[1], 10, 32)
	if err != nil {
		return partitionTableInfo{}, fmt.Errorf(
			"could not parse node ID from table (field: %v): %w",
			fields[1],
			err,
		)
	}

	start, err := strconv.ParseUint(fields[2], 10, 64)
	if err != nil {
		return partitionTableInfo{}, fmt.Errorf(
			"could not parse partition start from table (field: %v): %w",
			fields[2],
			err,
		)
	}

	end, err := strconv.ParseUint(fields[3], 10, 64)
	if err != nil {
		return partitionTableInfo{}, fmt.Errorf(
			"could not parse partition end from table (field: %v): %w",
			fields[3],
			err,
		)
	}

	part := partitionTableInfo{
		name:   table,
		nodeID: uint32(nodeID),
		start:  start,
		end:    end,
	}

	return part, nil
}

// group partitions by originator and sort them
// NOTE: sort does NOT validate any overlapping ranges, non-contigious ranges or anything like that.
func groupAndSortPartitions(partitions []partitionTableInfo) nodePartitions {
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

	for i, p := range partitions {

		if p.start >= p.end {
			return fmt.Errorf(
				"invalid partition size (start: %v, end: %v)", p.start, p.end)
		}

		// Save the current partition info so we can compare it against the next one.
		if prev == nil {
			prev = &partitions[i]
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

		// Update previous
		prev = &partitions[i]
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

func calculateFillRatio(start uint64, end uint64, lastSequenceID uint64) float64 {
	if lastSequenceID == 0 {
		return 0
	}

	// Should never happen because of DB constraints, but check for the sake of arithmetic safety.
	// This might be even a good place to panic as it might indicate a DB problem.
	if lastSequenceID < start {
		return 0
	}

	ratio := float64(lastSequenceID-start) / float64(end-start)
	return ratio
}
