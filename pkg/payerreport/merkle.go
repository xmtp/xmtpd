package payerreport

import (
	"errors"
	"math/big"

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

func generateMerkleTree(payers payerMap) (*merkle.MerkleTree, error) {
	leaves := make([]merkle.Leaf, 0, len(payers))
	var (
		leaf []byte
		err  error
	)
	for address, fees := range payers {
		if leaf, err = buildLeaf(address, fees); err != nil {
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
