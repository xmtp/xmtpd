package merkle

import (
	"bytes"
	"errors"
	"fmt"
)

const (
	ErrGenerateProof string = "cannot generate proof: %w"
	ErrVerifyProof   string = "cannot verify proof: %w"
)

var (
	ErrNilProof             = errors.New("nil proof")
	ErrNilRoot              = errors.New("nil root")
	ErrNoProofs             = errors.New("no proofs provided")
	ErrInvalidStartingIndex = errors.New("invalid starting index")
	ErrInsufficientProofs   = errors.New("insufficient proofs provided")
)

type ProofElement []byte

type MultiProof struct {
	startingIndex int
	leaves        []Leaf
	proofElements []ProofElement
}

func (p *MultiProof) GetStartingIndex() int {
	return p.startingIndex
}

func (p *MultiProof) GetLeafCount() (int, error) {
	if len(p.proofElements) == 0 {
		return 0, ErrNoProofs
	}

	return Bytes32ToInt(p.proofElements[0])
}

func (p *MultiProof) GetProofElements() []ProofElement {
	return p.proofElements
}

func (p *MultiProof) GetLeaves() []Leaf {
	return p.leaves
}

func NewMerkleProof(startingIndex int, leaves []Leaf, proofElements []ProofElement) *MultiProof {
	return &MultiProof{
		startingIndex: startingIndex,
		leaves:        leaves,
		proofElements: proofElements,
	}
}

// Verify verifies a MultiProof against the given tree root.
func Verify(root []byte, proof *MultiProof) (bool, error) {
	if len(root) == 0 {
		return false, fmt.Errorf(ErrVerifyProof, ErrNilRoot)
	}

	if err := proof.validate(); err != nil {
		return false, fmt.Errorf(ErrVerifyProof, err)
	}

	computedRoot, err := proof.computeRoot()
	if err != nil {
		return false, fmt.Errorf(ErrVerifyProof, err)
	}

	return bytes.Equal(computedRoot, root), nil
}

// validate performs common validation for Merkle proofs.
func (p *MultiProof) validate() error {
	if p.startingIndex < 0 {
		return ErrInvalidStartingIndex
	}

	// Must as least have the first proof element, which is the leaf count.
	if len(p.proofElements) == 0 {
		return ErrNoProofs
	}

	leafCount, err := Bytes32ToInt(p.proofElements[0])
	if err != nil {
		return fmt.Errorf(ErrVerifyProof, err)
	}

	// If the tree is not empty, and we are not proving any leaves, then the proof elements must include the leaf count
	// and the sub-root.
	if leafCount > 0 && len(p.leaves) == 0 && len(p.proofElements) < 2 {
		return ErrInsufficientProofs
	}

	if p.startingIndex+len(p.leaves) > leafCount {
		return ErrIndicesOutOfBounds
	}

	for _, leaf := range p.leaves {
		if leaf == nil {
			return ErrNilLeaf
		}
	}

	for _, proofElement := range p.proofElements {
		if proofElement == nil {
			return ErrNilProof
		}
	}

	return nil
}

// getNextProofElement safely retrieves the next proof and increments the index.
func (p *MultiProof) getNextProofElement(index *int) (ProofElement, error) {
	if *index >= len(p.proofElements) {
		return nil, ErrNoProofs
	}
	proofElement := p.proofElements[*index]
	*index++
	return proofElement, nil
}

// nodeQueue represents a node in the computation queue with its tree index and hash value.
// It's used during proof verification to track nodes as they are processed.
type nodeQueue struct {
	index int
	hash  []byte
}

// buildNodeQueue builds the node queue for the proof computation.
func (p *MultiProof) buildNodeQueue(balancedLeafCount int) ([]nodeQueue, error) {
	nodes, err := makeLeafNodes(p.leaves)
	if err != nil {
		return nil, err
	}

	n := len(p.leaves)

	queue := make([]nodeQueue, n)
	for i := range p.leaves {
		queue[n-1-i] = nodeQueue{
			index: balancedLeafCount + p.startingIndex + i,
			hash:  nodes[i],
		}
	}

	return queue, nil
}

// isEven returns true if the given number is even.
func isEven(n int) bool {
	return n&1 == 0
}

// computeRoot computes the root of the Merkle tree from the given proof.
func (p *MultiProof) computeRoot() ([]byte, error) {
	leafCount, err := Bytes32ToInt(p.proofElements[0])
	if err != nil {
		return nil, err
	}

	if len(p.leaves) == 0 && leafCount == 0 {
		return EmptyTreeRoot, nil
	}

	if len(p.leaves) == 0 {
		root, err := HashRoot(leafCount, p.proofElements[1])
		if err != nil {
			return nil, err
		}
		return root, nil
	}

	// 1. Prepare the queue.
	balancedLeafCount, err := CalculateBalancedNodesCount(leafCount)
	if err != nil {
		return nil, err
	}

	// 2. Build the circular queue, starting with the leaf nodes, in reverse order
	queue, err := p.buildNodeQueue(balancedLeafCount)
	if err != nil {
		return nil, err
	}

	queueLen := len(queue)

	var (
		readIdx, writeIdx, proofIdx = 0, 0, 1
		upperBound                  = balancedLeafCount + leafCount - 1
		lowestTreeIndex             = balancedLeafCount + p.startingIndex
		right                       []byte
		left                        []byte
	)

	// 3. Process queue until we hit the root.
	for {
		nodeIdx := queue[readIdx%queueLen].index

		// If we reach the sub-root (i.e. `index == 1`), we can return the root (i.e. `index == 0`) by
		// hashing the tree's leaf count with the last computed hash.
		if nodeIdx == 1 {
			root, err := HashRoot(leafCount, queue[(writeIdx-1)%queueLen].hash)
			if err != nil {
				return nil, err
			}
			return root, nil
		}

		// If node index we are handling is the upper bound and is even, then it's sibling to the right does not
		// exist (since this is an unbalanced tree), so we can just copy the hash up one level.
		if nodeIdx == upperBound && isEven(nodeIdx) {
			queue[writeIdx%queueLen] = nodeQueue{
				index: nodeIdx >> 1,
				hash:  HashPairlessNode(queue[readIdx%queueLen].hash),
			}
			writeIdx++
			readIdx++

			// If we are not at the lowest tree index (i.e. there are nodes to the left that we have yet to process at
			// this level), then continue.
			if nodeIdx != lowestTreeIndex {
				continue
			}

			// If we are at the lowest tree index (i.e. there are no nodes to the left that we have yet to process at
			// this level), then we can update the lower bound and upper bound for the next level up.
			lowestTreeIndex >>= 1
			upperBound >>= 1

			continue
		}

		nextNodeIdx := queue[(readIdx+1)%queueLen].index

		// Since we are processing nodes from right to left, then if the current node index is even, and there exists
		// nodes to the right (or else the previous if-continue would have been hit), then the right part of the hash is
		// a decommitment. If the current node index is odd, then the right part of the hash we already have computed.
		if isEven(nodeIdx) {
			right, err = p.getNextProofElement(&proofIdx)
			if err != nil {
				return nil, err
			}
		} else {
			right = queue[readIdx%queueLen].hash
			readIdx++
		}

		// Based on the current node index and the next node index, we can determine if the left part of the hash is an
		// existing computed hash or a decommitment.
		if isEven(nodeIdx) || (nextNodeIdx == nodeIdx-1) {
			left = queue[readIdx%queueLen].hash
			readIdx++
		} else {
			left, err = p.getNextProofElement(&proofIdx)
			if err != nil {
				return nil, err
			}
		}

		queue[writeIdx%queueLen] = nodeQueue{
			index: nodeIdx >> 1,
			hash:  HashNodePair(left, right),
		}
		writeIdx++

		// If we are not at the lowest tree index (i.e. there are nodes to the left that we have yet to process at this
		// level), then continue.
		// NOTE: Technically, if only `nextNodeIndex_ == lowestTreeIndex_`, and we did not use the hash at that
		// `nextNodeIndex_` as part of this step's hashing, then it was a node not yet handled, but it will be handled
		// in the next iteration, so the process will continue normally even if we prematurely "leveled up".
		if nodeIdx != lowestTreeIndex && nextNodeIdx != lowestTreeIndex {
			continue
		}

		// If we are at the lowest tree index (i.e. there are no nodes to the left that we have yet to process
		// level), then we can update the lower bound and upper bound for the next level up.
		// NOTE: Again, see the NOTE above.
		lowestTreeIndex >>= 1
		upperBound >>= 1
	}

	return nil, ErrNilRoot
}
