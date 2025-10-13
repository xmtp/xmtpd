package payerreport

import (
	"errors"
	"math/big"
	"sort"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/xmtp/xmtpd/pkg/currency"
	"github.com/xmtp/xmtpd/pkg/merkle"
)

var payerLeaf = abi.Arguments{
	{
		Name: "address",
		Type: abi.Type{T: abi.AddressTy},
	},
	{
		Name: "fee",
		Type: abi.Type{T: abi.UintTy, Size: 96},
	},
}

func GenerateMerkleTree(payers PayerMap) (*merkle.MerkleTree, error) {
	leaves := make([]merkle.Leaf, 0, len(payers))
	var (
		leaf []byte
		err  error
	)

	payerAddresses := make([]common.Address, 0, len(payers))
	for address := range payers {
		payerAddresses = append(payerAddresses, address)
	}

	sort.Slice(payerAddresses, func(i, j int) bool {
		return payerAddresses[i].Cmp(payerAddresses[j]) < 0
	})

	for _, address := range payerAddresses {
		if leaf, err = buildLeaf(address, payers[address]); err != nil {
			return nil, err
		}
		leaves = append(leaves, leaf)
	}

	return merkle.NewMerkleTree(leaves)
}

func buildLeaf(address common.Address, fees currency.PicoDollar) ([]byte, error) {
	if fees < 0 {
		return nil, errors.New("fee is negative")
	}

	leaf, err := payerLeaf.Pack(address, big.NewInt(int64(fees)))
	if err != nil {
		return nil, err
	}

	return leaf, nil
}
