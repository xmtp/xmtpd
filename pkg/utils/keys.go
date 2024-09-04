package utils

import (
	"crypto/ecdsa"
	"encoding/hex"
	"strings"

	"github.com/ethereum/go-ethereum/crypto"
)

func ParseEcdsaPrivateKey(key string) (*ecdsa.PrivateKey, error) {
	return crypto.HexToECDSA(strings.TrimPrefix(key, "0x"))
}

func EcdsaPrivateKeyToString(key *ecdsa.PrivateKey) string {
	keyBytes := crypto.FromECDSA(key)
	return "0x" + hex.EncodeToString(keyBytes)
}
