package testutils

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

func IntToHash(x int64) common.Hash {
	return common.BigToHash(big.NewInt(x))
}
