package payerreport

import (
	"sort"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/xmtp/xmtpd/pkg/currency"
	"github.com/xmtp/xmtpd/pkg/merkle"
)

var leafArgs = abi.Arguments{
	{
		Name: "address",
		Type: abi.Type{T: abi.AddressTy},
	},
	{
		Name: "amount",
		Type: abi.Type{T: abi.UintTy, Size: 64},
	},
}

func NewPayerMerkleTree(payers map[common.Address]currency.PicoDollar) (*merkle.MerkleTree, error) {
	leaves := make([]merkle.Leaf, 0, len(payers))

	keys := []common.Address{}
	for payerAddress := range payers {
		keys = append(keys, payerAddress)
	}
	sort.SliceStable(keys, func(i int, j int) bool {
		return keys[i].String() < keys[j].String()
	})

	for _, key := range keys {
		leafBytes, err := buildPayerLeaf(key, payers[key])
		if err != nil {
			return nil, err
		}
		leaves = append(leaves, leafBytes)
	}

	merkleTree, err := merkle.NewMerkleTree(leaves)
	if err != nil {
		return nil, err
	}

	return merkleTree, nil
}

func buildPayerLeaf(payerAddress common.Address, amount currency.PicoDollar) (merkle.Leaf, error) {
	return leafArgs.Pack(payerAddress, uint64(amount))
}
