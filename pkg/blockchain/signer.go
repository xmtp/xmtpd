package blockchain

import (
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

/*
*
PrivateKeySigner is a naive and not secure implementation of the TransactionSigner interface.

It is meant to be used in tests only
*
*/
type PrivateKeySigner struct {
	accountAddress common.Address
	signFunction   bind.SignerFn
}

func NewPrivateKeySigner(privateKeyString string, chainID int) (*PrivateKeySigner, error) {
	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(privateKeyString, "0x"))
	if err != nil {
		return nil, fmt.Errorf("unable to parse private key: %v", err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("Failed to cast to ECDSA public key %v", err)
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	auth, err := bind.NewKeyedTransactorWithChainID(
		privateKey,
		big.NewInt(int64(chainID)),
	)
	if err != nil {
		return nil, fmt.Errorf("Failed to create transactor: %v", err)
	}

	return &PrivateKeySigner{
		accountAddress: fromAddress,
		signFunction:   auth.Signer,
	}, nil
}

func (s *PrivateKeySigner) FromAddress() common.Address {
	return s.accountAddress
}

func (s *PrivateKeySigner) SignerFunc() bind.SignerFn {
	return s.signFunction
}
