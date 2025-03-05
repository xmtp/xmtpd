package merkle

import (
	"fmt"
	"math/big"
	"sort"

	smt "github.com/FantasyJony/openzeppelin-merkle-tree-go/standard_merkle_tree"
	"github.com/ethereum/go-ethereum/common"
)

var leafEncodings = []string{
	smt.SOL_UINT256,
	smt.SOL_ADDRESS,
	smt.SOL_UINT256,
}

type PayerReportInput struct {
	PayerAddress common.Address
	Amount       *big.Int
}

type PayerReportTree struct {
	tree      *smt.StandardTree
	numInputs int
	inputs    []PayerReportInput
}

func NewPayerReportTree(inputs []PayerReportInput) (*PayerReportTree, error) {
	if len(inputs) == 0 {
		return nil, fmt.Errorf("inputs cannot be empty")
	}
	values := buildValues(inputs)
	tree, err := smt.Of(values, leafEncodings)
	if err != nil {
		return nil, err
	}

	return &PayerReportTree{tree: tree, numInputs: len(inputs), inputs: inputs}, nil
}

func (t *PayerReportTree) Root() []byte {
	return t.tree.GetRoot()
}

func (t *PayerReportTree) GetMultiProof(offset int, count int) (*MultiProof, error) {
	if offset+count > t.numInputs {
		return nil, fmt.Errorf("offset + count is greater than the number of inputs")
	}
	leaves := t.tree.Entries()
	valueIndices := make([]int, count)
	for i := range count {
		valueIndices[i] = offset + i
	}

	// You can only get a multi proof if the indices are sorted in the order they exist in the tree
	// This will be different from the order of the inputs
	sort.Slice(valueIndices, func(x int, y int) bool {
		treeIndexX := findTreeIndex(leaves, valueIndices[x])
		treeIndexY := findTreeIndex(leaves, valueIndices[y])
		return treeIndexX > treeIndexY
	})

	proof, err := t.tree.GetMultiProofWithIndices(valueIndices)
	if err != nil {
		return nil, err
	}

	values := make([]*PayerReportValue, count)
	for i, valueIndex := range valueIndices {
		values[i] = &PayerReportValue{
			PayerAddress: t.inputs[valueIndex].PayerAddress,
			Amount:       t.inputs[valueIndex].Amount,
			Index:        big.NewInt(int64(valueIndex)),
		}
	}

	return &MultiProof{proof: proof, leaves: values, offset: offset, root: [32]byte(t.Root())}, nil
}

func findTreeIndex(leaves []*smt.LeafValue, valueIndex int) int {
	for _, leaf := range leaves {
		if leaf.ValueIndex == valueIndex {
			return leaf.TreeIndex
		}
	}
	return -1
}

func buildValues(inputs []PayerReportInput) [][]any {
	values := make([][]any, len(inputs))
	for i, input := range inputs {
		values[i] = []any{
			big.NewInt(int64(i)),
			input.PayerAddress,
			input.Amount,
		}
	}
	return values
}

type PayerReportValue struct {
	PayerAddress common.Address
	Amount       *big.Int
	Index        *big.Int
}

type MultiProof struct {
	proof  *smt.ValueMultiProof
	leaves []*PayerReportValue
	offset int
	root   [32]byte
}

func (p *MultiProof) Amounts() []*big.Int {
	out := make([]*big.Int, len(p.leaves))
	for i, leaf := range p.leaves {
		out[i] = leaf.Amount
	}
	return out
}

func (p *MultiProof) PayerAddresses() []common.Address {
	out := make([]common.Address, len(p.leaves))
	for i, leaf := range p.leaves {
		out[i] = leaf.PayerAddress
	}
	return out
}

func (p *MultiProof) Proof() [][32]byte {
	out := make([][32]byte, len(p.proof.Proof))
	for i, proof := range p.proof.Proof {
		out[i] = [32]byte(proof)
	}
	return out
}

func (p *MultiProof) ProofFlags() []bool {
	return p.proof.ProofFlags
}

func (p *MultiProof) Offset() *big.Int {
	return big.NewInt(int64(p.offset))
}

func (p *MultiProof) Root() [32]byte {
	return p.root
}

func (p *MultiProof) Indices() []*big.Int {
	out := make([]*big.Int, len(p.leaves))
	for i, leaf := range p.leaves {
		out[i] = leaf.Index
	}
	return out
}
